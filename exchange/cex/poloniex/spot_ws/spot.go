package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/poloniex"
	"clients/exchange/cex/poloniex/spot_api"
	"clients/logger"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
)

type PoloniexSpotWebsocket struct {
	base.WsConf
	apiClient     *spot_api.ApiClient
	listenTimeout int64
	WsReqUrl      *spot_api.WsReqUrl
	handler       *WebSocketSpotHandle
	pingInterval  int64
	reqId         int
	token         string
}

func NewPoloniexWebsocket(conf base.WsConf) *PoloniexSpotWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	b := &PoloniexSpotWebsocket{
		WsConf: conf,
	}
	if b.ReadTimeout == 0 {
		b.ReadTimeout = 300
	}
	if b.listenTimeout == 0 {
		b.listenTimeout = 1800
	}
	if b.pingInterval == 0 {
		b.pingInterval = 25
	}
	b.apiClient = spot_api.NewApiClient(conf.APIConf)

	b.handler = NewWebSocketSpotHandle(b.ChanCap)
	b.WsReqUrl = spot_api.NewWsReqUrl()

	return b
}

func NewPoloniexWebsocket2(conf base.WsConf, cli *http.Client) *PoloniexSpotWebsocket {

	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}

	b := &PoloniexSpotWebsocket{
		WsConf: conf,
	}
	if b.ReadTimeout == 0 {
		b.ReadTimeout = 300
	}
	if b.listenTimeout == 0 {
		b.listenTimeout = 1800
	}
	if b.pingInterval == 0 {
		b.pingInterval = 25
	}
	b.apiClient = spot_api.NewApiClient2(conf.APIConf, cli)

	b.handler = NewWebSocketSpotHandle(b.ChanCap)
	b.WsReqUrl = spot_api.NewWsReqUrl()

	return b
}

func (b *PoloniexSpotWebsocket) baseUrl() string {
	return b.WsReqUrl.WS_PUBLIC_BASE_URL
}

func (b *PoloniexSpotWebsocket) EstablishConn(url string, subMsg BaseSubscribe, symbolInfos []*client.SymbolInfo, handler func([]byte) error, ctx context.Context) error {
	heartBeat := func() []byte {
		return []byte("{\"event\":\"ping\"}")
	}

	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
	)

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).
		ProtoHandleFunc(handler).
		AutoReconnect().
		Heartbeat(heartBeat, time.Duration(b.pingInterval)*time.Second).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			return b.Subscribe(wsClient, subMsg)
		})

	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}

	dataSendStatus := b.handler.GetChan("send_status").(*base.CheckDataSendStatus)
	dataSendStatus.Init(3600, ctx, symbolInfos...)

	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}

	err := b.Subscribe(wsClient, subMsg)
	if err != nil {
		return err
	}
	go b.start(wsClient, ctx, subMsg)
	go b.resubscribeAfterSilence(wsClient, ctx, subMsg, dataSendStatus.UpdateTimeoutChMap, nil)
	return err
}

func (b *PoloniexSpotWebsocket) resubscribeAfterSilence(wsClient *conn.WsConn, ctx context.Context, subMsg BaseSubscribe, ch chan []*client.SymbolInfo, reconnectSignalMap map[string]chan struct{}) {
	var err error
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err = b.UnSubscribe(wsClient, subMsg); err != nil {
				logger.Logger.Error("unsubscribe err: ", subMsg)
			}
			break LOOP
		case resubSymbols := <-ch:
			symbolExStrs := make([]string, 0, len(resubSymbols))
			logger.Logger.Warn("silent symbols:", resubSymbols, " need resubscribe")
			for _, symbol := range resubSymbols {
				symbolExStrs = append(symbolExStrs, poloniex.Canonical2Exchange(symbol.Symbol))
			}
			msg := subMsg.WithSymbols(symbolExStrs)
			msg.Unsub()

			if reconnectSignalMap != nil {
				// 取消订阅后通知update goroutine正在重新订阅
				for _, symbol := range resubSymbols {
					if ch, ok := reconnectSignalMap[symbol.Symbol]; ok {
						ch <- struct{}{}
					} else {
						logger.Logger.Errorf("get reconnect signal channel: %s %s", symbol.Symbol, err.Error())
					}
				}
			}

			if err = b.UnSubscribe(wsClient, msg); err != nil {
				logger.Logger.Error("unsubscribe err:", msg)
			}
			msg.Sub()
			if err = b.Subscribe(wsClient, msg); err != nil {
				logger.Logger.Error("subscribe err:", msg)
			}
		}
	}
}

func (b *PoloniexSpotWebsocket) resubscribeOnceError(wsClient *conn.WsConn, ctx context.Context, subMsg BaseSubscribe, ch chan string) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case symbol := <-ch:
			logger.Logger.Warn("depth snapshot update error, resubscribe ", symbol)
			msg := subMsg.WithSymbols([]string{poloniex.Canonical2Exchange(symbol)})
			msg.Unsub()

			if err := b.UnSubscribe(wsClient, msg); err != nil {
				logger.Logger.Error("unsubscribe err:", msg)
			}
			msg.Sub()
			if err := b.Subscribe(wsClient, msg); err != nil {
				logger.Logger.Error("subscribe err:", msg)
			}
		}
	}
}

func (b *PoloniexSpotWebsocket) EstablishSnapshotConn(url string, subMsg BaseSubscribe, symbolInfos []*client.SymbolInfo, handler func([]byte) error, ctx context.Context, conf *base.IncrementDepthConf) error {
	var (
		err error
	)
	heartBeat := func() []byte {
		return []byte("{\"event\":\"ping\"}")
	}

	var (
		readTimeout     = time.Duration(b.ReadTimeout) * time.Second * 1000
		reconnectSymbol chan string
	)

	reconnectSignalChannelMap := make(map[string]chan struct{})
	for _, symbolInfo := range symbolInfos {
		reconnectSignalChannelMap[symbolInfo.Symbol] = make(chan struct{})
	}

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).
		ProtoHandleFunc(handler).
		AutoReconnect().
		Heartbeat(heartBeat, time.Duration(b.pingInterval)*time.Second).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			// 重连时提醒update goroutine 正在重连
			for _, ch := range reconnectSignalChannelMap {
				ch <- struct{}{}
			}
			b.Subscribe(wsClient, subMsg)
			return err
		})

	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}

	dataSendStatus := b.handler.GetChan("send_status").(*base.CheckDataSendStatus)
	dataSendStatus.Init(3600, ctx, symbolInfos...)

	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}

	b.handler.SetReconnectSignalChannel(reconnectSignalChannelMap)
	b.handler.SetDepthIncrementSnapShotConf(symbolInfos, conf)

	if b.handler.DepthIncrementSnapshotReconnectSymbol == nil {
		reconnectSymbol = make(chan string)
		b.handler.SetDepthIncrementSnapshotReconnectChan(reconnectSymbol)
	} else {
		reconnectSymbol = b.handler.DepthIncrementSnapshotReconnectSymbol
	}

	err = b.Subscribe(wsClient, subMsg)
	if err != nil {
		return err
	}
	go b.start(wsClient, ctx, subMsg)
	go b.resubscribeOnceError(wsClient, ctx, subMsg, reconnectSymbol)
	go b.resubscribeAfterSilence(wsClient, ctx, subMsg, dataSendStatus.UpdateTimeoutChMap, reconnectSignalChannelMap)
	return err
}

func (b *PoloniexSpotWebsocket) start(wsClient *conn.WsConn, ctx context.Context, subMsg BaseSubscribe) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			subMsg.Unsub()
			if err := b.UnSubscribe(wsClient, subMsg); err != nil {
				logger.Logger.Error("unsubscribe err: ", subMsg)
			}
			wsClient.CloseWs()
			break LOOP
		}
	}
}

func (b *PoloniexSpotWebsocket) Subscribe(client *conn.WsConn, subMsg BaseSubscribe) (err error) {
	logger.Logger.Info("Poloniex subscribe:", subMsg)
	err = client.Subscribe(subMsg)
	if err != nil {
		logger.Logger.Error("subscribe err:", err)
	}
	return err
}

func (b *PoloniexSpotWebsocket) UnSubscribe(client *conn.WsConn, subMsg BaseSubscribe) (err error) {
	subMsg.Unsub()
	logger.Logger.Info("Poloniex unsubscribe:", subMsg)
	err = client.Subscribe(subMsg)
	if err != nil {
		logger.Logger.Error("unsubscribe err:", err)
	}
	return
}

func makeSubscribeMessage(channel string, symbols []string) BaseSubscribe {
	return &SubscribeMessage{
		Event:   "subscribe",
		Channel: []string{channel},
		Symbols: symbols,
	}
}

func makeBookSubscribeMessage(channel string, symbols []string, depth int) BaseSubscribe {
	return &BookSubscribeMessage{
		Event:   "subscribe",
		Channel: []string{channel},
		Symbols: symbols,
		Depth:   depth,
	}
}

func (b *PoloniexSpotWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {

	symbols := make([]string, 0, len(chMap))
	symbolInfos := make([]*client.SymbolInfo, 0, len(chMap))

	for symbol := range chMap {
		symbols = append(symbols, poloniex.Canonical2Exchange(symbol.Symbol))
		symbolInfos = append(symbolInfos, symbol)
	}

	b.handler.SetTradeGroupChannel(chMap)
	msg := makeSubscribeMessage("trades", symbols)
	err := b.EstablishConn(b.baseUrl(), msg, symbolInfos, b.handler.GetHandler("trade"), ctx)
	return err
}

func (b *PoloniexSpotWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
	return errors.New("no ticker data")
}

func (b *PoloniexSpotWebsocket) DepthIncrementGroup(ctx context.Context, interval int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {

	symbols := make([]string, 0, len(chMap))
	symbolInfos := make([]*client.SymbolInfo, 0, len(chMap))

	for symbol := range chMap {
		symbols = append(symbols, poloniex.Canonical2Exchange(symbol.Symbol))
		symbolInfos = append(symbolInfos, symbol)
	}

	b.handler.SetDepthIncrementGroupChannel(chMap)
	msg := makeSubscribeMessage("book_lv2", symbols)
	err := b.EstablishConn(b.baseUrl(), msg, symbolInfos, b.handler.GetHandler("depthDiff"), ctx)
	return err
}

func DepthLevelParse(levelList [][]string) ([]*depth.DepthLevel, error) {

	var (
		res           []*depth.DepthLevel
		amount, price float64
		err           error
	)
	for _, level := range levelList {
		if price, err = strconv.ParseFloat(level[0], 64); err != nil {
			return nil, err
		}
		if amount, err = strconv.ParseFloat(level[1], 64); err != nil {
			return nil, err
		}
		res = append(res, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	return res, err
}

func (b *PoloniexSpotWebsocket) DepthLimitGroup(ctx context.Context, interval, limit int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	if limit != 0 && limit != 5 && limit != 10 && limit != 20 {
		return errors.New(fmt.Sprintf("unsupported limit, only support 0(default to 5), 5, 10, 20, got %d", limit))
	}
	if limit == 0 {
		limit = 5
	}

	symbols := make([]string, 0, len(chMap))
	symbolInfos := make([]*client.SymbolInfo, 0, len(chMap))

	for symbol := range chMap {
		symbols = append(symbols, poloniex.Canonical2Exchange(symbol.Symbol))
		symbolInfos = append(symbolInfos, symbol)
	}

	b.handler.SetDepthLimit(limit)
	b.handler.SetDepthLimitGroupChannel(chMap)
	msg := makeBookSubscribeMessage("book", symbols, limit)
	err := b.EstablishConn(b.baseUrl(), msg, symbolInfos, b.handler.GetHandler("depthWhole"), ctx)
	return err
}

func (b *PoloniexSpotWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, interval, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error {
	if isDelta && isFull {
		// check delta and snapshot
		if len(chDeltaMap) != len(chFullMap) || chDeltaMap == nil || chFullMap == nil {
			return errors.New("symbols for delmap and fullmap are not equal 1")
		}
		for symbol := range chDeltaMap {
			if _, ok := chFullMap[symbol]; !ok {
				return errors.New("symbols for delmap and fullmap are not equal 2")
			}
		}
	}

	var symbols []*client.SymbolInfo
	var symbolsString []string

	if isDelta {
		for symbol := range chDeltaMap {
			symbols = append(symbols, symbol)
			symbolsString = append(symbolsString, poloniex.Canonical2Exchange(symbol.Symbol))
		}
	} else if isFull {
		for symbol := range chFullMap {
			symbols = append(symbols, symbol)
			symbolsString = append(symbolsString, poloniex.Canonical2Exchange(symbol.Symbol))
		}
	} else {
		return errors.New("no symbol channel")
	}

	conf := &base.IncrementDepthConf{
		IsPublishDelta:    isDelta,
		IsPublishFull:     isFull,
		CheckTimeSec:      3000 + rand.Intn(1200),
		DepthCapLevel:     1000,
		DepthLevel:        limit,
		GetFullDepth:      nil,
		GetFullDepthLimit: nil,
		Ctx:               ctx,
	}

	b.handler.SetDepthIncrementSnapshotGroupChannel(chDeltaMap, chFullMap)

	msg := makeSubscribeMessage("book_lv2", symbolsString)
	err := b.EstablishSnapshotConn(b.baseUrl(), msg, symbols, b.handler.GetHandler("depthIncrementSnapshot"), ctx, conf)
	return err
}

func (b *PoloniexSpotWebsocket) FundingRateGroup(ctx context.Context, chFundMap map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	return nil
}
