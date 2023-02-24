package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/bitstamp/spot_api"
	"clients/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"strconv"
	"strings"
	"sync"
	"time"
)

var totalSymbolMap = make(map[string]string)

type BitstampSpotWebsocket struct {
	base.WsConf
	apiClient      *spot_api.ApiClient
	pingInterval   int64
	readTimeout    int64 //can delete
	listenTimeout  int64 //second
	lock           sync.Mutex
	sequenceNumber int
	isStart        bool
	WsReqUrl       *spot_api.WsReqUrl
	handler        WebSocketHandleInterface
}

func NewBitstampWebsocket(conf base.WsConf) *BitstampSpotWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	b := &BitstampSpotWebsocket{
		WsConf: conf,
	}
	if conf.EndPoint == "" {
		b.EndPoint = spot_api.WS_API_BASE_URL
	}
	if b.ReadTimeout == 0 {
		b.ReadTimeout = 300
	}
	if b.listenTimeout == 0 {
		b.listenTimeout = 1800
	}
	b.pingInterval = 5
	b.handler = NewWebSocketSpotHandle(b.ChanCap)
	b.apiClient = spot_api.NewApiClient(conf.APIConf)

	if conf.AccessKey != "" && conf.SecretKey != "" {
		b.apiClient = spot_api.NewApiClient(base.APIConf{
			ProxyUrl:    conf.ProxyUrl,
			ReadTimeout: conf.ReadTimeout,
			AccessKey:   conf.AccessKey,
			SecretKey:   conf.SecretKey,
		})
	}
	b.WsReqUrl = spot_api.NewSpotWsUrl()

	if len(totalSymbolMap) == 0 {
		data, err := b.apiClient.GetTradingPairsInfo()
		if err != nil {
			logger.Logger.Error("map initialize error:", err)
		}
		for _, pair := range data.TradingPairs {
			totalSymbolMap[pair.UrlSymbol] = pair.Name
		}
	}
	return b
}

/* Connection Methods*/
func (b *BitstampSpotWebsocket) Subscribe(wsClient *conn.WsConn, request Request, symbols ...string) error {
	var (
		err error
	)

	request.Event = SUBSCRIBE
	channel := request.Data.Channel

	for _, symbol := range symbols {
		request.Data.Channel = fmt.Sprintf("%s%s", channel, ReformatSymbols(symbol))
		fmt.Println("1010101010101", request.Data.Channel)
		err = wsClient.Subscribe(request)
		if err != nil {
			logger.Logger.Error("subscribe err:", err)
		}
	}

	return err
}
func (b *BitstampSpotWebsocket) Unsubscribe(wsClient *conn.WsConn, request Request, symbols ...string) (err error) {

	request.Event = UNSUBSCRIBE
	channel := request.Data.Channel

	for _, symbol := range symbols {
		request.Data.Channel = fmt.Sprintf("%s%s", channel, ReformatSymbols(symbol))
		err = wsClient.Subscribe(request)
		if err != nil {
			logger.Logger.Error("subscribe err:", err)
		}
	}
	return
}
func (b *BitstampSpotWebsocket) Start(wsClient *conn.WsConn, ctx context.Context, request Request, symbols ...string) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := b.Unsubscribe(wsClient, request, symbols...); err != nil {
				logger.Logger.Error("unsubscribe err:", request, err)
			}
			wsClient.CloseWs()
			break LOOP
		}
	}
}
func (b *BitstampSpotWebsocket) EstablishConn(request Request, url string, handler func([]byte) error, ctx context.Context, symbols []*client.SymbolInfo) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		symbolList  []string
		err         error
	)

	for _, symbol := range symbols {
		symbolList = append(symbolList, symbol.Symbol)
	}

	heartbeat := func() []byte {
		data, _ := json.Marshal(PingMessage{Event: "bts:heartbeat"})
		return data
	}

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).
		ProtoHandleFunc(handler).
		AutoReconnect().
		Heartbeat(heartbeat, time.Duration(b.pingInterval)*time.Second).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			return b.Subscribe(wsClient, request, symbolList...)
		})

	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = b.Subscribe(wsClient, request, symbolList...)
	if err != nil {
		return err
	}
	go b.Start(wsClient, ctx, request, symbolList...)
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.Reconnect(wsClient, ctx, request, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap, symbolList...)
	return err
}

func (b *BitstampSpotWebsocket) EstablishSnapshotConn(symbols []*client.SymbolInfo, request Request, url string, handler func([]byte) error, ctx context.Context, conf *base.IncrementDepthConf, interval int) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		symbolList  []string
		err         error
	)

	for _, symbol := range symbols {
		symbolList = append(symbolList, symbol.Symbol)
	}

	heartbeat := func() []byte {
		data, _ := json.Marshal(PingMessage{Event: "bts:heartbeat"})
		return data
	}

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).
		ProtoHandleFunc(handler).
		AutoReconnect().
		Heartbeat(heartbeat, time.Duration(b.pingInterval)*time.Second).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			return b.Subscribe(wsClient, request, symbolList...)
		})
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	b.handler.SetDepthIncrementSnapShotConf(symbols, conf)
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = b.Subscribe(wsClient, request, symbolList...)
	if err != nil {
		return err
	}
	go b.Start(wsClient, ctx, request, symbolList...)
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.Reconnect(wsClient, ctx, request, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap, symbolList...)
	return err
}
func (b *BitstampSpotWebsocket) Reconnect(wsClient *conn.WsConn, ctx context.Context, request Request, ch chan []*client.SymbolInfo, params ...string) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := b.Unsubscribe(wsClient, request, params...); err != nil {
				logger.Logger.Error("unsubscribe err:", request)
			}
			wsClient.CloseWs()
			break LOOP
		case reconnectData := <-ch:
			var ReconnectSymbols []string
			for _, symbol := range reconnectData {
				ReconnectSymbols = append(ReconnectSymbols, symbol.Symbol)
			}
			logger.Logger.Info("Update Timeout Resubscribe:", ReconnectSymbols)
			if err := b.Unsubscribe(wsClient, request, ReconnectSymbols...); err != nil {
				logger.Logger.Error("unsubscribe err:", request, ReconnectSymbols)
			}
			if err := b.Subscribe(wsClient, request, ReconnectSymbols...); err != nil {
				logger.Logger.Error("unsubscribe err:", request, ReconnectSymbols)
			}
		}
	}
}

/*Main Interface Methods*/
func (b *BitstampSpotWebsocket) FundingRateGroup(context.Context, map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	logger.Logger.Error("spot does not have funding group data")
	return nil
}
func (b *BitstampSpotWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
	var (
		symbols []*client.SymbolInfo
		request Request
		err     error
	)

	for symbolInfo := range chMap {
		symbols = append(symbols, symbolInfo)
	}

	b.handler.SetTradeGroupChannel(chMap)

	request = Request{
		Data: Data{
			Channel: LiveTrades,
		},
	}
	err = b.EstablishConn(request, spot_api.WS_API_BASE_URL, b.handler.TradeGroupHandle, ctx, symbols)
	return err
}
func (b *BitstampSpotWebsocket) BookTickerGroup(context.Context, map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
	logger.Logger.Error("spot does not have ticker data")
	return nil
}
func (b *BitstampSpotWebsocket) DepthLimitGroup(ctx context.Context, interval, limit int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	var (
		symbols []*client.SymbolInfo
		request Request
		err     error
	)

	for symbolInfo := range chMap {
		symbols = append(symbols, symbolInfo)
	}

	b.handler.SetDepthLimitGroupChannel(chMap)
	request = Request{
		Data: Data{
			Channel: LiveOrderBook,
		},
	}
	err = b.EstablishConn(request, spot_api.WS_API_BASE_URL, b.handler.DepthLimitGroupHandle, ctx, symbols)
	return err
}
func (b *BitstampSpotWebsocket) DepthIncrementGroup(ctx context.Context, internal int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {
	var (
		symbols []*client.SymbolInfo
		request Request
		err     error
	)

	for symbolInfo := range chMap {
		symbols = append(symbols, symbolInfo)
	}

	b.handler.SetDepthIncrementGroupChannel(chMap)
	request = Request{
		Data: Data{
			Channel: LiveFullOrderBook,
		},
	}
	err = b.EstablishConn(request, spot_api.WS_API_BASE_URL, b.handler.DepthIncrementGroupHandle, ctx, symbols)
	return err
}

func (b *BitstampSpotWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, interval, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error {
	var (
		symbols []*client.SymbolInfo
		request Request
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

	b.handler.SetDepthIncrementSnapshotGroupChannel(chDeltaMap, chFullMap)

	conf := &base.IncrementDepthConf{
		IsPublishDelta: isDelta,
		IsPublishFull:  isFull,
		DepthCapLevel:  1000,
		DepthLevel:     limit,
		GetFullDepth:   b.GetFullDepth,
		Ctx:            ctx,
	}
	request = Request{
		Data: Data{
			Channel: LiveFullOrderBook,
		},
	}
	err = b.EstablishSnapshotConn(symbols, request, spot_api.WS_API_BASE_URL, b.handler.DepthIncrementSnapShotGroupHandle, ctx, conf, interval)
	return err
}

/*Helper Functions*/
func (b *BitstampSpotWebsocket) GetFullDepth(symbol string) (*base.OrderBook, error) {
	var (
		res = &base.OrderBook{}
	)

	resp, err := b.apiClient.GetOrderbook(symbol)
	if err != nil {
		return nil, err
	}

	timeStamp, err := strconv.ParseInt(resp.Microtimestamp, 10, 64)
	if err != nil {
		logger.Logger.Error("parse time err:", string(resp.Microtimestamp))
		return res, err
	}

	res = &base.OrderBook{
		Exchange:     common.Exchange_BITSTAMP,
		Symbol:       symbol,
		Market:       common.Market_SPOT,
		Type:         common.SymbolType_SPOT_NORMAL,
		TimeExchange: uint64(timeStamp),
	}
	ParseOrder(resp.Asks, &res.Asks)
	ParseOrder(resp.Bids, &res.Bids)
	return res, nil
}
func ReformatSymbols(input string) string {
	return strings.ToLower(strings.Replace(input, "/", "", -1))
}
