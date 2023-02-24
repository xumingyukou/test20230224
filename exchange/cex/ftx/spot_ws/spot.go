package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/ftx/future_api"
	"clients/exchange/cex/ftx/spot_api"
	"clients/logger"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/warmplanet/proto/go/common"
	"net/http"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
)

type FTXSpotWebsocket struct {
	base.WsConf
	apiClient     *spot_api.ApiClient
	uapiCient     *future_api.UApiClient
	pingInterval  int64
	listenTimeout int64 //second
	lock          sync.Mutex
	isLoggedIn    bool
	channelType   ChannelType
	handler       WebSocketHandleInterface
}

func NewFTXWebsocket(conf base.WsConf) *FTXSpotWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &FTXSpotWebsocket{
		WsConf: conf,
	}
	if conf.EndPoint == "" {
		d.EndPoint = spot_api.WS_API_BASE_URL
	}
	if d.ReadTimeout == 0 {
		d.ReadTimeout = 300
	}
	if d.listenTimeout == 0 {
		d.listenTimeout = 1800
	}
	d.pingInterval = 14
	d.isLoggedIn = false
	d.handler = NewWebSocketSpotHandle(d.ChanCap, conf.APIConf)
	d.apiClient = spot_api.NewApiClient(conf.APIConf)
	d.uapiCient = future_api.NewUApiClient(conf.APIConf)
	if conf.AccessKey != "" && conf.SecretKey != "" {
		d.apiClient = spot_api.NewApiClient(base.APIConf{
			ProxyUrl:    conf.ProxyUrl,
			ReadTimeout: conf.ReadTimeout,
			AccessKey:   conf.AccessKey,
			SecretKey:   conf.SecretKey,
		})
	}
	return d
}

func NewFTXWebsocket2(conf base.WsConf, cli *http.Client) *FTXSpotWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &FTXSpotWebsocket{
		WsConf: conf,
	}

	if d.ReadTimeout == 0 {
		d.ReadTimeout = 300
	}
	if d.listenTimeout == 0 {
		d.listenTimeout = 1800
	}
	// ws socket目前固定写死
	d.EndPoint = spot_api.WS_API_BASE_URL
	d.handler = NewWebSocketSpotHandle(d.ChanCap, conf.APIConf)
	d.pingInterval = 14
	d.isLoggedIn = false
	if conf.AccessKey != "" && conf.SecretKey != "" {
		d.apiClient = spot_api.NewApiClient2(base.APIConf{
			ReadTimeout: conf.ReadTimeout,
			AccessKey:   conf.AccessKey,
			SecretKey:   conf.SecretKey,
			EndPoint:    conf.EndPoint,
		}, cli)
	}
	return d
}

func (b *FTXSpotWebsocket) Subscribe(wsClient *conn.WsConn, request WSRequest, params ...string) error {
	var (
		err error
	)

	for _, param := range params {
		request.Market = param
		request.Op = Subscribe
		err = wsClient.Subscribe(request)
		if err != nil {
			logger.Logger.Error("subscribe err:", err)
		}
	}
	return err
}
func (b *FTXSpotWebsocket) Unsubscribe(wsClient *conn.WsConn, request WSRequest, params ...string) (err error) {
	for _, param := range params {
		request.Market = param
		request.Op = UnSubscribe
		err = wsClient.Subscribe(request)
		if err != nil {
			logger.Logger.Error("unsubscribe err:", err)
		}
	}
	return
}
func (b *FTXSpotWebsocket) Start(wsClient *conn.WsConn, ctx context.Context, request WSRequest, params ...string) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			wsClient.CloseWs()
			if err := b.Unsubscribe(wsClient, request, params...); err != nil {
				logger.Logger.Error("unsubscribe err:", params)
			}
			break LOOP
		}
	}
}
func (b *FTXSpotWebsocket) Reconnect(wsClient *conn.WsConn, request WSRequest, ctx context.Context, symbol *client.SymbolInfo, ch chan bool) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := b.Unsubscribe(wsClient, request, GetExchangeSymName(symbol)); err != nil {
				logger.Logger.Error("unsubscribe err:", GetExchangeSymName(symbol), err)
			}
			wsClient.CloseWs()
			break LOOP
		case status, ok := <-ch:
			if !ok {
				break LOOP
			}
			if !status {
				continue
			}
			logger.Logger.Error("reconnecting", GetExchangeSymName(symbol))
			if err := b.Unsubscribe(wsClient, request, GetExchangeSymName(symbol)); err != nil {
				logger.Logger.Error("unsubscribe err:", GetExchangeSymName(symbol), err)
			}
			if err := b.Subscribe(wsClient, request, GetExchangeSymName(symbol)); err != nil {
				logger.Logger.Error("subscribe err:", GetExchangeSymName(symbol), err)
			}
		}
	}
}

// 获取交易所的交易对名称
func GetExchangeSymName(symbolInfo *client.SymbolInfo) string {
	if symbolInfo.Market == common.Market_SPOT {
		return symbolInfo.Symbol
	} else {
		return future_api.GetBaseSymbol(symbolInfo.Symbol, symbolInfo.Type)
	}
}

func (b *FTXSpotWebsocket) EstablishConn(request WSRequest, url string, handler func([]byte) error, ctx context.Context, symbols []*client.SymbolInfo) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		symbolNames []string
		err         error
	)
	heartBeat := func() []byte {
		data, _ := json.Marshal(PingMessage{Op: "ping"})
		return data
	}

	for _, symbol := range symbols {
		symbolNames = append(symbolNames, GetExchangeSymName(symbol))
	}

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).ProtoHandleFunc(handler).
		AutoReconnect().
		Heartbeat(heartBeat, time.Duration(b.pingInterval)*time.Second).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			return b.Subscribe(wsClient, request, symbolNames...)
		})
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = b.Subscribe(wsClient, request, symbolNames...)
	if err != nil {
		return err
	}
	go b.Start(wsClient, ctx, request, symbolNames...)
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.Reconnect2(wsClient, request, ctx, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap, symbolNames...)
	return err
}
func (b *FTXSpotWebsocket) Reconnect2(wsClient *conn.WsConn, request WSRequest, ctx context.Context, ch chan []*client.SymbolInfo, params ...string) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := b.Unsubscribe(wsClient, request, params...); err != nil {
				logger.Logger.Error("unsubscribe err:", params)
			}
			wsClient.CloseWs()
			break LOOP
		case x := <-ch:
			var ReconnectSymbols []string
			for _, symbol := range x {
				ReconnectSymbols = append(ReconnectSymbols, GetExchangeSymName(symbol))
			}
			logger.Logger.Info("update timeout reconnecting:", ReconnectSymbols)
			if err := b.Unsubscribe(wsClient, request, ReconnectSymbols...); err != nil {
				logger.Logger.Error("unsubscribe err:", request, ReconnectSymbols)
			}
			if err := b.Subscribe(wsClient, request, ReconnectSymbols...); err != nil {
				logger.Logger.Error("subscribe err:", request, ReconnectSymbols)
			}
		}
	}
}

func (b *FTXSpotWebsocket) EstablishSnapshotConn(symbols []*client.SymbolInfo, request WSRequest, url string, handler func([]byte) error, ctx context.Context, conf *base.IncrementDepthConf) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		symbolNames []string
		err         error
	)
	for _, symbol := range symbols {
		symbolNames = append(symbolNames, GetExchangeSymName(symbol))
	}
	heartBeat := func() []byte {
		data, _ := json.Marshal(PingMessage{Op: "ping"})
		return data
	}
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).ProtoHandleFunc(handler).
		AutoReconnect().
		Heartbeat(heartBeat, time.Duration(b.pingInterval)*time.Second).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			return b.Subscribe(wsClient, request, symbolNames...)
		})
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	b.handler.SetDepthIncrementSnapShotConf(symbols, conf)
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = b.Subscribe(wsClient, request, symbolNames...)
	if err != nil {
		return err
	}
	go b.Start(wsClient, ctx, request, symbolNames...)

	for symbolInfo, ch := range conf.DepthNotMatchChanMap {
		go b.Reconnect(wsClient, request, ctx, symbolInfo, ch)
	}
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.Reconnect2(wsClient, request, ctx, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap, symbolNames...)
	return err
}

func (b *FTXSpotWebsocket) FundingRateGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	b.handler.SetFundingRateGroupChannel(chMap)
	err := b.handler.FundingRateGroupHandle([]byte{})
	return err
}

func (b *FTXSpotWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
	var (
		symbols []*client.SymbolInfo
		request WSRequest
		err     error
	)
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	b.handler.SetTradeGroupChannel(chMap)
	request = WSRequest{
		ChannelType: TradesChannel,
	}
	err = b.EstablishConn(request, spot_api.WS_API_BASE_URL, b.handler.TradeGroupHandle, ctx, symbols)
	return err
}

func (b *FTXSpotWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
	var (
		symbols []*client.SymbolInfo
		request WSRequest
		err     error
	)
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	b.handler.SetBookTickerGroupChannel(chMap)

	request = WSRequest{
		ChannelType: TickerChannel,
	}
	err = b.EstablishConn(request, spot_api.WS_API_BASE_URL, b.handler.BookTickerGroupHandle, ctx, symbols)
	return err
}

// Only first message is the Depth Limit at 100 snapshot
func (b *FTXSpotWebsocket) DepthLimitGroup(ctx context.Context, interval, limit int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	var (
		symbols []*client.SymbolInfo
		request WSRequest
		err     error
	)
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}

	b.handler.SetDepthLimitGroupChannel(chMap)
	request = WSRequest{
		ChannelType: OrderbookChannel,
	}
	err = b.EstablishConn(request, spot_api.WS_API_BASE_URL, b.handler.DepthLimitGroupHandle, ctx, symbols)
	return err
}

// All messages following the first message are increment updates
func (b *FTXSpotWebsocket) DepthIncrementGroup(ctx context.Context, internal int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {
	var (
		symbols []*client.SymbolInfo
		request WSRequest
		err     error
	)
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}

	b.handler.SetDepthIncrementGroupChannel(chMap)
	request = WSRequest{
		ChannelType: OrderbookChannel,
	}
	err = b.EstablishConn(request, spot_api.WS_API_BASE_URL, b.handler.DepthIncrementGroupHandle, ctx, symbols)
	return err
}

func (b *FTXSpotWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, internal, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error { //增量合并为全量副本，并校验
	var (
		symbols []*client.SymbolInfo
		params  WSRequest
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
		CheckTimeSec:   3600,
		DepthCapLevel:  1000,
		DepthLevel:     limit,
		Ctx:            ctx,
	}

	params = WSRequest{
		ChannelType: OrderbookChannel,
	}
	err = b.EstablishSnapshotConn(symbols, params, spot_api.WS_API_BASE_URL, b.handler.DepthIncrementSnapShotGroupHandle, ctx, conf)

	return err
}

func (b *FTXSpotWebsocket) EstablishPrivateConn(request WSRequest, url string, handler func([]byte) error, ctx context.Context, params ...string) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
	)
	heartBeat := func() []byte {
		data, _ := json.Marshal(PingMessage{Op: "ping"})
		return data
	}
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(url).ProtoHandleFunc(handler).AutoReconnect().Heartbeat(heartBeat, time.Duration(b.pingInterval)*time.Second)
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = b.PrivateSubscribe(wsClient, request)
	if err != nil {
		return err
	}
	go b.PrivateStart(wsClient, ctx, request)
	return err
}

func (b *FTXSpotWebsocket) Account(context.Context, ...*client.WsAccountReq) (<-chan *client.WsAccountRsp, error) {
	return nil, errors.New("ftx account ws not implement")
}
func (b *FTXSpotWebsocket) Balance(context.Context, ...*client.WsAccountReq) (<-chan *client.WsBalanceRsp, error) {
	return nil, errors.New("ftx account ws not implement")
}

func (b *FTXSpotWebsocket) Order(ctx context.Context, reqList ...*client.WsAccountReq) (<-chan *client.WsOrderRsp, error) {
	request := WSRequest{
		ChannelType: OrdersChannel,
	}
	b.handler.SetOrderMarketConf(reqList)
	err := b.EstablishPrivateConn(request, spot_api.WS_API_BASE_URL, b.handler.OrderHandle, ctx)
	return b.handler.GetChan("order").(chan *client.WsOrderRsp), err
}

func (b *FTXSpotWebsocket) Authorize(client *conn.WsConn) (err error) {
	if b.isLoggedIn {
		return nil
	}

	time := time.Now().UnixMilli()
	mac := hmac.New(sha256.New, []byte(b.SecretKey))
	mac.Write([]byte(fmt.Sprintf("%dwebsocket_login", time)))

	authentication := WSAuthorizationRequest{
		Args: WSAuthorizationArgs{
			Key:  b.AccessKey,
			Sign: hex.EncodeToString(mac.Sum(nil)),
			Time: time,
		},
		Op: "login",
	}

	client.SendJsonMessage(authentication)
	b.isLoggedIn = true
	return
}

func (b *FTXSpotWebsocket) PrivateSubscribe(wsClient *conn.WsConn, request WSRequest) error {
	var (
		err error
	)
	b.Authorize(wsClient)
	request.Op = Subscribe
	err = wsClient.Subscribe(request)
	if err != nil {
		logger.Logger.Error("subscribe err:", err)
	}
	return err
}

func (b *FTXSpotWebsocket) PrivateUnsubscribe(wsClient *conn.WsConn, request WSRequest) (err error) {
	request.Op = UnSubscribe
	err = wsClient.Subscribe(request)
	if err != nil {
		logger.Logger.Error("subscribe err:", err)
	}
	return
}

func (b *FTXSpotWebsocket) PrivateStart(wsClient *conn.WsConn, ctx context.Context, request WSRequest) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			wsClient.CloseWs()
			b.isLoggedIn = false
			if err := b.PrivateUnsubscribe(wsClient, request); err != nil {
				logger.Logger.Error("unsubscribe err:", err)
			}
			break LOOP
		}
	}
}
