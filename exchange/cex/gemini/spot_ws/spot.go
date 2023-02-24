package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/gemini/spot_api"
	"clients/logger"
	"clients/transform"
	"context"
	"errors"
	"fmt"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
	"strings"
	"sync"
	"time"
)

var symbolNameMap = make(map[string]string)

type GeminiSpotWebsocket struct {
	base.WsConf
	apiClient     *spot_api.ApiClient
	readTimeout   int64 // can delete
	listenTimeout int64 // second
	lock          sync.Mutex
	isStart       bool
	WsReqUrl      *spot_api.WsReqUrl
	handler       *WebSocketSpotHandle
	pingInterval  int64
}

func NewGeminiSpotWebsocket(conf base.WsConf) *GeminiSpotWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &GeminiSpotWebsocket{
		WsConf: conf,
	}
	if conf.EndPoint == "" {
		d.EndPoint = spot_api.WS_API_Public_Topic
	}
	if d.ReadTimeout == 0 {
		d.ReadTimeout = 300
	}
	if d.listenTimeout == 0 {
		d.listenTimeout = 1800
	}
	d.apiClient = spot_api.NewApiClient(conf.APIConf)

	d.handler = NewWebSocketSpotHandle(d.ChanCap)

	if conf.AccessKey != "" && conf.SecretKey != "" {
		d.apiClient = spot_api.NewApiClient(base.APIConf{
			ProxyUrl:    conf.ProxyUrl,
			ReadTimeout: conf.ReadTimeout,
			AccessKey:   conf.AccessKey,
			SecretKey:   conf.SecretKey,
		})
	} else {
		d.apiClient = spot_api.NewApiClient(conf.APIConf)
	}

	d.WsReqUrl = spot_api.NewSpotWsUrl()

	for len(symbolNameMap) == 0 {
		logger.Logger.Info("map initializing......")
		data, err := d.apiClient.GetSymbols()
		if err != nil {
			logger.Logger.Error("map initialize error:", err)
			return nil
		}
		wg := sync.WaitGroup{}
		wg.Add(len(*data))
		symbolPerRoutine := 10
		for i := 0; i < (len(*data)-1)/symbolPerRoutine+1; i++ {
			start := symbolPerRoutine * i
			end := transform.Min(symbolPerRoutine*(i+1), len(*data))
			go func(symbols []string) {
				for _, sym := range symbols {
					d.AddFormatSym2Map(sym)
					wg.Done()
				}
			}([]string(*data)[start:end])
		}
		wg.Wait()
	}

	return d
}

func (b *GeminiSpotWebsocket) AddFormatSym2Map(symbol string) {
	res, err := b.apiClient.GetSymbolDetails(symbol)
	if err != nil {
		logger.Logger.Error("add symbol to map error:", err)
		return
	}
	b.lock.Lock()
	symbolNameMap[res.Symbol] = res.BaseCurrency + "/" + res.QuoteCurrency
	b.lock.Unlock()
}

func (b *GeminiSpotWebsocket) EstablishCoon(url string, handler func([]byte) error, ctx context.Context, symbols []*client.SymbolInfo) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
	)

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(url).AutoReconnect().ProtoHandleFunc(handler).PostReconnectSuccess(func(wsClient *conn.WsConn) error { return b.Subscribe(wsClient, symbols) })
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = b.Subscribe(wsClient, symbols)
	if err != nil {
		return err
	}
	go b.Start(wsClient, ctx, symbols)
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.reConnect2(wsClient, symbols, ctx, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)
	return err
}

func (b *GeminiSpotWebsocket) EstablishSnapshotConn(url string, handler func([]byte) error, ctx context.Context, conf *base.IncrementDepthConf, symbols []*client.SymbolInfo) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
	)

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(url).AutoReconnect().ProtoHandleFunc(handler).PostReconnectSuccess(func(wsClient *conn.WsConn) error { return b.Subscribe(wsClient, symbols) })
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	b.handler.SetDepthIncrementSnapShotConf(symbols, conf)
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = b.Subscribe(wsClient, symbols)
	if err != nil {
		return err
	}
	go b.Start(wsClient, ctx, symbols)
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.reConnect2(wsClient, symbols, ctx, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)
	return err
}

func (b *GeminiSpotWebsocket) reConnect2(wsClient *conn.WsConn, symbols []*client.SymbolInfo, ctx context.Context, ch chan []*client.SymbolInfo) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := b.Unsubscribe(wsClient, symbols); err != nil {
				logger.Logger.Error("unsubscribe err:", symbols)
			}
			wsClient.CloseWs()
			break LOOP
		case reconnectData := <-ch:
			var reconnectSymbols []*client.SymbolInfo
			for _, symbol := range reconnectData {
				reconnectSymbols = append(reconnectSymbols, symbol)
			}
			if err := b.Unsubscribe(wsClient, reconnectSymbols); err != nil {
				logger.Logger.Error("unsubscribe err:", reconnectSymbols)
			}
			if err := b.Subscribe(wsClient, reconnectSymbols); err != nil {
				logger.Logger.Error("unsubscribe err:", reconnectSymbols)
			}
		}
	}
}

func (b *GeminiSpotWebsocket) Subscribe(wsClient *conn.WsConn, symbols []*client.SymbolInfo) error {
	var (
		err error
	)
	for _, symbol := range symbols {
		// 先重置币对的推送标记
		b.handler.depthIncrJudgeMap.Delete(GetSymbolName(symbol.Symbol))
		fmt.Println("Gemini subscribe", symbol)
		var symbollist []string
		symbollist = append(symbollist, GetSymbolName(symbol.Symbol))
		var subscriptions []Subcriptions
		subscriptions = append(subscriptions, Subcriptions{
			Name:    "l2",
			Symbols: symbollist,
		})
		err = wsClient.Subscribe(Msg{
			Type:          "subscribe",
			Subscriptions: subscriptions,
		})
		if err != nil {
			logger.Logger.Error("subscribe err:", err)
		}
	}

	return err
}

func (b *GeminiSpotWebsocket) Unsubscribe(wsClient *conn.WsConn, symbols []*client.SymbolInfo) error {
	var (
		err error
	)
	for _, symbol := range symbols {
		fmt.Println("Gemini unsubscribe", symbol)
		var subscriptions []Subcriptions
		var unsubSymbol []string
		unsubSymbol = append(unsubSymbol, GetSymbolName(symbol.Symbol))
		subscriptions = append(subscriptions, Subcriptions{
			Name:    "l2",
			Symbols: unsubSymbol,
		})
		err = wsClient.Subscribe(Msg{
			Type:          "unsubscribe",
			Subscriptions: subscriptions,
		})
		if err != nil {
			logger.Logger.Error("unsubscribe err:", err)
		}
	}
	return err
}

func (b *GeminiSpotWebsocket) Start(wsClient *conn.WsConn, ctx context.Context, symbols []*client.SymbolInfo) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			wsClient.CloseWs()
			if err := b.Unsubscribe(wsClient, symbols); err != nil {
				logger.Logger.Error("unsubscribe err:")
			}
			break LOOP
		}
	}
}

func (b *GeminiSpotWebsocket) FundingRateGroup(context.Context, map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	logger.Logger.Error("spot do not have funding rate")
	return nil
}

func (b *GeminiSpotWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
	var (
		symbols []*client.SymbolInfo
		err     error
	)

	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	//fmt.Println(symbols)
	b.handler.SetTradeGroupChannel(chMap)

	err = b.EstablishCoon(spot_api.WS_API_Public_Topic, b.handler.GeneralHandle("trade"), ctx, symbols)
	return err
}

func (b *GeminiSpotWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
	logger.Logger.Error("spot do not  have book ticker")
	return nil
}

func (b *GeminiSpotWebsocket) DepthLimitGroup(ctx context.Context, interval, limit int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	logger.Logger.Error("spot do not have depth limit")
	return nil
}

func (b *GeminiSpotWebsocket) DepthIncrementGroup(ctx context.Context, interval int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {
	var (
		symbols []*client.SymbolInfo
		err     error
	)

	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	//fmt.Println(symbols)
	b.handler.SetDepthIncrementGroupChannel(chMap)

	err = b.EstablishCoon(spot_api.WS_API_Public_Topic, b.handler.GeneralHandle("increment"), ctx, symbols)
	return err
}

func (b *GeminiSpotWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, interval, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error {
	var (
		symbols []*client.SymbolInfo
		err     error
	)
	if isDelta && isFull && len(chDeltaMap) != len(chFullMap) {
		return errors.New("chDeltaMap and chFullMap not match")
	}
	if isDelta {
		for symbol := range chDeltaMap {
			symbols = append(symbols, symbol)
		}
	} else {
		for symbol := range chFullMap {
			symbols = append(symbols, symbol)
		}
	}

	//fmt.Println(symbols)
	b.handler.SetDepthIncrementSnapshotGroupChannel(chDeltaMap, chFullMap)
	conf := &base.IncrementDepthConf{
		IsPublishDelta: isDelta,
		IsPublishFull:  isFull,
		CheckTimeSec:   3600,
		DepthCapLevel:  1000,
		DepthLevel:     limit,
		Ctx:            ctx,
	}
	err = b.EstablishSnapshotConn(spot_api.WS_API_Public_Topic, b.handler.GeneralHandle("snapshot"), ctx, conf, symbols)
	return err
}

func GetSymbolName(symbol string) string {
	symbol = strings.Replace(symbol, "/", "", 1)
	//symbol = strings.Replace(symbol, "-", "", 1)
	//symbol = strings.Replace(symbol, "_", "", 1)
	return symbol
}
