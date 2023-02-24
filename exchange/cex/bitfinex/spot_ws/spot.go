package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/bitfinex/spot_api"
	"clients/logger"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
)

var (
	Client2Exchange = make(map[string]string)
	Exchange2Client = make(map[string]string)
	// symbolNameMap   = make(map[string]string) --> Exchange2Client
	// getSymbolName   = make(map[string]string) --> Client2Exchange
	mainlock      sync.Mutex
	slidingWindow = conn.NewSlidingWindow(2, 6000)
	// cntOfWsConn   int64
	// preTime       int64
	// capOfWSConn   int64 = 20
)

type BitfinexSpotWebsocket struct {
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

func NewBitfinexSpotWebsocket(conf base.WsConf) *BitfinexSpotWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &BitfinexSpotWebsocket{
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
	mainlock.Lock()
	for len(Client2Exchange) == 0 {
		exchangeCurrencys, err := d.apiClient.GetExchange()
		if err != nil {
			logger.Logger.Error("get exchange ", err.Error())
			continue
		}
		var fullSymbols []string
		for _, exchangeSymbols := range *exchangeCurrencys {
			for _, fullSymbol := range exchangeSymbols {
				fullSymbols = append(fullSymbols, fullSymbol)
			}
		}
		baseCurrencys, err := d.apiClient.GetCurrency()
		if err != nil {
			logger.Logger.Error("get exchange ", err.Error())
			continue
		}
		var preSymbols []string
		for _, baseSymbols := range *baseCurrencys {
			for _, baseSymbol := range baseSymbols {
				preSymbols = append(preSymbols, baseSymbol)
			}
		}
		for _, symbol := range fullSymbols {
			if find := strings.Contains(symbol, ":"); find {
				s := strings.Split(symbol, ":")
				Exchange2Client[symbol] = s[0] + "/" + s[1]
				Client2Exchange[s[0]+"/"+s[1]] = symbol
				continue
			}
			for _, base := range preSymbols {
				if strings.HasPrefix(symbol, base) {
					quote := strings.Replace(symbol[len(base):], ":", "", -1)
					Exchange2Client[symbol] = base + "/" + quote
					Client2Exchange[base+"/"+quote] = symbol
				}
			}
		}
		// preTime = time.Now().Unix()
	}
	mainlock.Unlock()
	return d
}

func (b *BitfinexSpotWebsocket) EstablishConn(request WSRequest, url string, handler func([]byte) error, ctx context.Context, symbols []*client.SymbolInfo) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
	)
	for !slidingWindow.Allow() {
		time.Sleep(time.Second)
	}

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(url).AutoReconnect().ProtoHandleFunc(handler).PostReconnectSuccess(func(wsClient *conn.WsConn) error { return b.Subscribe(wsClient, request, symbols) })
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = b.Subscribe(wsClient, request, symbols)
	if err != nil {
		return err
	}

	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.reConnectAfterSilence(wsClient, request, symbols, ctx, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)
	go b.Start(wsClient, ctx, request, symbols)
	return err
}

func (b *BitfinexSpotWebsocket) EstablishSnapshotConn(request WSRequest, url string, handler func([]byte) error, ctx context.Context, conf *base.IncrementDepthConf, symbols []*client.SymbolInfo) error {
	var (
		readTimeout     = time.Duration(b.ReadTimeout) * time.Second * 1000
		err             error
		reconnectSymbol chan string
	)
	for !slidingWindow.Allow() {
		time.Sleep(time.Second)
	}

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(url).AutoReconnect().ProtoHandleFunc(handler).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			err := wsClient.Subscribe(ConfigureMsg{
				Event: "conf",
				Flags: 131072 + 536870912,
			})
			if err != nil {
				return err
			}
			return b.Subscribe(wsClient, request, symbols)
		})
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	b.handler.SetDepthIncrementSnapShotConf(symbols, conf)
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	// reconnect channel if checksum error
	if b.handler.DepthIncrementSnapshotReconnectSymbol == nil {
		reconnectSymbol = make(chan string)
		b.handler.SetDepthIncrementSnapshotReconnectChan(reconnectSymbol)
	} else {
		reconnectSymbol = b.handler.DepthIncrementSnapshotReconnectSymbol
	}

	err = wsClient.Subscribe(ConfigureMsg{
		Event: "conf",
		Flags: 131072 + 536870912,
	})
	if err != nil {
		return err
	}

	err = b.Subscribe(wsClient, request, symbols)
	if err != nil {
		return err
	}
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.reConnectAfterSilence(wsClient, request, symbols, ctx, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)
	go b.resubscribeOnceError(wsClient, ctx, request, reconnectSymbol)
	go b.Start(wsClient, ctx, request, symbols)
	return err
}

func (b *BitfinexSpotWebsocket) resubscribeOnceError(wsClient *conn.WsConn, ctx context.Context, request WSRequest, ch chan string) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case symbol := <-ch:
			logger.Logger.Warn("checksum error, resubscribe ", symbol)
			symbolInfo := &client.SymbolInfo{
				Symbol: symbol,
			}

			if err := b.Unsubscribe(wsClient, request, []*client.SymbolInfo{symbolInfo}); err != nil {
				logger.Logger.Error("unsubscribe err:", request, " ", symbol)
			}

			if err := b.Subscribe(wsClient, request, []*client.SymbolInfo{symbolInfo}); err != nil {
				logger.Logger.Error("subscribe err:", request, " ", symbol)
			}
		}
	}
}

func (b *BitfinexSpotWebsocket) reConnectAfterSilence(wsClient *conn.WsConn, request WSRequest, symbols []*client.SymbolInfo, ctx context.Context, ch chan []*client.SymbolInfo) {
	// LOOP:
	for {
		select {
		// case <-ctx.Done():
		// 	if err := b.Unsubscribe(wsClient, request, symbols); err != nil {
		// 		logger.Logger.Error("unsubscribe err:", symbols)
		// 	}
		// 	wsClient.CloseWs()
		// 	break LOOP
		case reconnectData := <-ch:
			var reconnectSymbols []*client.SymbolInfo
			for _, symbol := range reconnectData {
				reconnectSymbols = append(reconnectSymbols, symbol)
			}
			if err := b.Unsubscribe(wsClient, request, reconnectSymbols); err != nil {
				logger.Logger.Error("unsubscribe err:", reconnectSymbols)
			}
			if err := b.Subscribe(wsClient, request, reconnectSymbols); err != nil {
				logger.Logger.Error("subscribe err:", reconnectSymbols)
			}
		}
	}
}

func (b *BitfinexSpotWebsocket) Subscribe(wsClient *conn.WsConn, request WSRequest, symbols []*client.SymbolInfo) error {
	var (
		err error
	)
	for _, symbol := range symbols {
		fmt.Println("Bitfinex subscribe:", request.Channel, symbol)
		if request.Channel == "book" {
			b.handler.depthIncrJudgeMap.Delete(symbol.Symbol)
			err = wsClient.Subscribe(BookSubMsg{
				Event:   "subscribe",
				Channel: "book",
				Symbol:  "t" + Client2Exchange[symbol.Symbol],
				Freq:    "F0",
				Len:     "250",
				Prec:    "P0",
			})
		} else {
			err = wsClient.Subscribe(SubMsg{
				Event:   "subscribe",
				Channel: request.Channel,
				Symbol:  "t" + Client2Exchange[symbol.Symbol],
			})
		}

		if err != nil {
			logger.Logger.Error("subscribe err:", err)
		}
	}
	return err
}

func (b *BitfinexSpotWebsocket) Unsubscribe(wsClient *conn.WsConn, request WSRequest, symbols []*client.SymbolInfo) error {
	var (
		err error
	)
	for _, symbol := range symbols {
		err = wsClient.Subscribe(UnsubMsg{
			Event:  "unsubscribe",
			ChanID: b.handler.GetChanID(symbol.Symbol),
		})
		if err != nil {
			logger.Logger.Error("unsubscribe err:", err)
		}
	}
	return err
}

func (b *BitfinexSpotWebsocket) Start(wsClient *conn.WsConn, ctx context.Context, request WSRequest, symbols []*client.SymbolInfo) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			wsClient.CloseWs()
			if err := b.Unsubscribe(wsClient, request, symbols); err != nil {
				logger.Logger.Error("unsubscribe err:", request)
			}
			break LOOP
		}
	}
}

func (b *BitfinexSpotWebsocket) FundingRateGroup(context.Context, map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	logger.Logger.Error("spot do not have funding rate")
	return errors.New("spot do not have funding rate")
}

func (b *BitfinexSpotWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
	if len(chMap) > 30 {
		logger.Logger.Error("All websocket connections have a limit of 30 subscriptions to public market data feed channels")
		return nil
	}
	var (
		symbols []*client.SymbolInfo
		request WSRequest
		err     error
	)
	request.Channel = "trades"
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}

	b.handler.SetTradeGroupChannel(chMap)

	err = b.EstablishConn(request, spot_api.WS_API_Public_Topic, b.handler.GeneralHandle("trades"), ctx, symbols)
	return err
}

func (b *BitfinexSpotWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
	if len(chMap) > 30 {
		logger.Logger.Error("All websocket connections have a limit of 30 subscriptions to public market data feed channels")
		return nil
	}
	var (
		symbols []*client.SymbolInfo
		request WSRequest
		err     error
	)
	request.Channel = "ticker"
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}

	b.handler.SetBookTickerGroupChannel(chMap)

	err = b.EstablishConn(request, spot_api.WS_API_Public_Topic, b.handler.GeneralHandle("ticker"), ctx, symbols)
	return err
}

func (b *BitfinexSpotWebsocket) DepthLimitGroup(ctx context.Context, interval, limit int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	logger.Logger.Error("spot do not have depth limit")
	return nil
}

func (b *BitfinexSpotWebsocket) DepthIncrementGroup(ctx context.Context, interval int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {
	if len(chMap) > 30 {
		logger.Logger.Error("All websocket connections have a limit of 30 subscriptions to public market data feed channels")
		return nil
	}
	var (
		symbols []*client.SymbolInfo
		request WSRequest
		err     error
	)
	request.Channel = "book"
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}

	b.handler.SetDepthIncrementGroupChannel(chMap)

	err = b.EstablishConn(request, spot_api.WS_API_Public_Topic, b.handler.GeneralHandle("increment"), ctx, symbols)
	return err
}

func (b *BitfinexSpotWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, interval, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error {
	if len(chDeltaMap) > 30 || len(chFullMap) > 30 {
		logger.Logger.Error("All websocket connections have a limit of 30 subscriptions to public market data feed channels")
		return nil
	}
	var (
		symbols []*client.SymbolInfo
		request WSRequest
		err     error
	)
	request.Channel = "book"
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
	b.handler.SetDepthIncrementSnapshotGroupChannel(chDeltaMap, chFullMap)
	conf := &base.IncrementDepthConf{
		IsPublishDelta: isDelta,
		IsPublishFull:  isFull,
		CheckTimeSec:   3600,
		DepthCapLevel:  1000,
		DepthLevel:     limit,
		Ctx:            ctx,
	}
	err = b.EstablishSnapshotConn(request, spot_api.WS_API_Public_Topic, b.handler.GeneralHandle("snapshot"), ctx, conf, symbols)
	return err
}
