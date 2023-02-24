package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/bybit/spot_api"
	"clients/logger"
	"context"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
	"strings"
	"sync"
	"time"
)

var symbolNameMap = make(map[string]string)

type BybitSpotWebsocket struct {
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

func NewBybitSpotWebsocket(conf base.WsConf) *BybitSpotWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &BybitSpotWebsocket{
		WsConf: conf,
	}
	if conf.EndPoint == "" {
		d.EndPoint = spot_api.WS_API_IP_URL
	}
	if d.ReadTimeout == 0 {
		d.ReadTimeout = 300
	}
	if d.listenTimeout == 0 {
		d.listenTimeout = 1800
	}
	d.pingInterval = 19
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
		logger.Logger.Error("map initializing......")
		data, err := d.apiClient.GetSymbols()
		if err != nil {
			logger.Logger.Error("map initialize error:", err)
		}
		for _, pair := range data.Result {
			if pair.ShowStatus == true {
				symbolNameMap[pair.Name] = pair.BaseCurrency + "/" + pair.QuoteCurrency
			}
		}
		//fmt.Println(symbolNameMap)
	}
	return d
}

func (b *BybitSpotWebsocket) Reconnect(wsClient *conn.WsConn, version int, request WSRequest, ctx context.Context, symbols []*client.SymbolInfo, ch chan bool) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := b.Unsubscribe(wsClient, request, symbols); err != nil {
				logger.Logger.Error("unsubscribe err:", symbols, err)
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
			fmt.Println("Reconnecting:", symbols)
			if err := b.Unsubscribe(wsClient, request, symbols); err != nil {
				logger.Logger.Error("unsubscribe err:", symbols, err)
			}
			if err := b.Subscribe(wsClient, version, request, symbols); err != nil {
				logger.Logger.Error("subscribe err:", symbols, err)
			}
		}
	}
}

func (b *BybitSpotWebsocket) reConnect2(wsClient *conn.WsConn, version int, request WSRequest, ctx context.Context, symbols []*client.SymbolInfo, ch chan []*client.SymbolInfo) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := b.Unsubscribe(wsClient, request, symbols); err != nil {
				logger.Logger.Error("unsubscribe err:", symbols)
			}
			wsClient.CloseWs()
			break LOOP
		case reconnectData := <-ch:
			var reconnectSymbols []*client.SymbolInfo
			for _, symbol := range reconnectData {
				reconnectSymbols = append(reconnectSymbols, symbol)
			}
			if err := b.Unsubscribe(wsClient, request, reconnectSymbols); err != nil {
				logger.Logger.Error("unsubscribe err:", reconnectSymbols)
			}
			if err := b.Subscribe(wsClient, version, request, reconnectSymbols); err != nil {
				logger.Logger.Error("unsubscribe err:", reconnectSymbols)
			}
		}
	}
}

func (b *BybitSpotWebsocket) EstablishCoon(version int, request WSRequest, url string, handler func([]byte) error, ctx context.Context, symbols []*client.SymbolInfo) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
	)

	heartBeat := func() []byte {
		data, _ := json.Marshal(PingMsg{
			Ping: time.Now().UnixMicro(),
		})
		return data
	}
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(url).AutoReconnect().Heartbeat(heartBeat, time.Duration(b.pingInterval)*time.Second).ProtoHandleFunc(handler).PostReconnectSuccess(func(wsClient *conn.WsConn) error { return b.Subscribe(wsClient, version, request, symbols) })
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = b.Subscribe(wsClient, version, request, symbols)
	if err != nil {
		return err
	}

	go b.Start(wsClient, ctx, request, symbols)
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.reConnect2(wsClient, version, request, ctx, symbols, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)
	return err
}

func (b *BybitSpotWebsocket) EstablishSnapshotConn(version int, request WSRequest, symbols []*client.SymbolInfo, url string, handler func([]byte) error, ctx context.Context, conf *base.IncrementDepthConf) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		symbolNames []string
		err         error
	)

	heartBeat := func() []byte {
		data, _ := json.Marshal(PingMsg{
			Ping: time.Now().UnixMicro(),
		})
		return data
	}
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(url).AutoReconnect().Heartbeat(heartBeat, time.Duration(b.pingInterval)*time.Second).ProtoHandleFunc(handler).PostReconnectSuccess(func(wsClient *conn.WsConn) error { return b.Subscribe(wsClient, version, request, symbols) })
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	b.handler.SetDepthIncrementSnapShotConf(symbols, conf)
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}

	for _, symbol := range symbols {
		symbolNames = append(symbolNames, GetSymbolName(symbol.Symbol))
	}

	err = b.Subscribe(wsClient, version, request, symbols)
	if err != nil {
		return err
	}

	go b.Start(wsClient, ctx, request, symbols)
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.reConnect2(wsClient, version, request, ctx, symbols, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)
	//for symbolInfo, ch := range conf.DepthNotMatchChanMap {
	//	go b.Reconnect(wsClient, version, request, ctx, symbolInfo, ch)
	//}
	return err
}

func (b *BybitSpotWebsocket) Subscribe(wsClient *conn.WsConn, version int, request WSRequest, symbols []*client.SymbolInfo) error {
	var (
		err error
	)
	for _, symbol := range symbols {
		if version == 1 {
			fmt.Println("Bybit subscribe:", GetSymbolName(symbol.Symbol))
			param := Param1{
				Binary: request.Binary,
			}
			err = wsClient.Subscribe(SendMsg1{
				Symbol: GetSymbolName(symbol.Symbol),
				Topic:  request.Topic,
				Event:  "sub",
				Params: param,
			})
		} else if version == 2 {
			fmt.Println("Bybit subscribe:", GetSymbolName(symbol.Symbol))
			param := Param2{
				Symbol: GetSymbolName(symbol.Symbol),
				Binary: request.Binary,
			}
			err = wsClient.Subscribe(SendMsg2{
				Topic:  request.Topic,
				Event:  "sub",
				Params: param,
			})
		}
		if err != nil {
			logger.Logger.Error("subscribe err:", err)
		}
	}

	return err
}

func (b *BybitSpotWebsocket) Unsubscribe(wsClient *conn.WsConn, request WSRequest, symbols []*client.SymbolInfo) error {
	var (
		err error
	)
	for _, symbol := range symbols {
		fmt.Println("Bybit unsubscribe", GetSymbolName(symbol.Symbol))
		param := Param1{
			Binary: request.Binary,
		}
		err = wsClient.Subscribe(SendMsg1{
			Symbol: GetSymbolName(symbol.Symbol),
			Topic:  request.Topic,
			Event:  "cancel",
			Params: param,
		})
		if err != nil {
			logger.Logger.Error("unsubscribe err:", err)
		}
	}
	return err
}

func (b *BybitSpotWebsocket) Start(wsClient *conn.WsConn, ctx context.Context, request WSRequest, symbols []*client.SymbolInfo) {
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

func (b *BybitSpotWebsocket) FundingRateGroup(context.Context, map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	logger.Logger.Error("spot do not have funding rate")
	return nil
}

func (b *BybitSpotWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
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
		Topic:  "trade",
		Binary: false,
	}
	err = b.EstablishCoon(2, request, spot_api.WS_API_Public_Topic2, b.handler.GeneralHandle("trades"), ctx, symbols)
	return err
}

func (b *BybitSpotWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
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
		Topic:  "bookTicker",
		Binary: false,
	}
	err = b.EstablishCoon(2, request, spot_api.WS_API_Public_Topic2, b.handler.GeneralHandle("book"), ctx, symbols)
	return err
}

func (b *BybitSpotWebsocket) DepthLimitGroup(ctx context.Context, interval, limit int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
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
		Topic:  "depth",
		Binary: false,
	}
	err = b.EstablishCoon(2, request, spot_api.WS_API_Public_Topic2, b.handler.GeneralHandle("limit"), ctx, symbols)
	return err
}

func (b *BybitSpotWebsocket) DepthIncrementGroup(ctx context.Context, interval int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {
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
		Topic:  "diffDepth",
		Binary: false,
	}
	err = b.EstablishCoon(1, request, spot_api.WS_API_Public_Topic1, b.handler.GeneralHandle("increment"), ctx, symbols)
	return err
}

func (b *BybitSpotWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, internal, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error { //增量合并为全量副本，并校验
	var (
		symbols []*client.SymbolInfo
		request WSRequest
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
	request = WSRequest{
		Topic:  "diffDepth",
		Binary: false,
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
	err = b.EstablishSnapshotConn(1, request, symbols, spot_api.WS_API_Public_Topic1, b.handler.GeneralHandle("snapshot"), ctx, conf)

	return err
}

func GetSymbolName(symbol string) string {
	symbol = strings.Replace(symbol, "/", "", 1)
	//symbol = strings.Replace(symbol, "-", "", 1)
	//symbol = strings.Replace(symbol, "_", "", 1)
	return symbol
}

//func (c *BybitSpotWebsocket) FullDepthStart(wsClient *conn.WsConn, ctx context.Context) {
//LOOP:
//	for {
//		select {
//		case <-ctx.Done():
//			wsClient.CloseWs()
//			break LOOP
//		}
//	}
//}

//func (c *BybitSpotWebsocket) GetFullDepth(symbol string) (*base.OrderBook, error) {
//	fmt.Println("Check —— spot.go")
//	var (
//		request     WSRequest
//		ch          = make(chan *RespLimitDepthStream, 1)
//		depthCache  = &base.OrderBook{}
//		err         error
//		readTimeout = time.Duration(c.ReadTimeout) * time.Second * 1000
//	)
//
//	c.lock.Lock()
//	defer c.lock.Unlock()
//
//	ctx, cancel := context.WithCancel(context.Background()) //Can pass cancel in lieu of _
//
//	request = WSRequest{
//		Topic:  "diffDepth",
//		Binary: false,
//	}
//
//	checker := func(data []byte) error {
//		iData := &RespLimitDepthStream{}
//		err := json.Unmarshal(data, iData)
//		if err == nil {
//			if iData.F == true {
//				ch <- iData
//			}
//		}
//		return err
//	}
//	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(spot_api.WS_API_Public_Topic1).ProtoHandleFunc(checker)
//	if c.ProxyUrl != "" {
//		wsBuilder = wsBuilder.ProxyUrl(c.ProxyUrl)
//	}
//	wsClient := wsBuilder.Build()
//	if wsClient == nil {
//		cancel()
//		return depthCache, errors.New("build websocket connection error")
//	}
//
//	err = c.Subscribe(wsClient, 1, request, symbol)
//	if err != nil {
//		cancel()
//		return depthCache, err
//	}
//	go c.FullDepthStart(wsClient, ctx)
//LOOP:
//	for {
//		select {
//		case respTemp, _ := <-ch: //stop on "type" = "l2update" message
//			asks, err := DepthLevelParse(respTemp.Data[0].A)
//			if err != nil {
//				logger.Logger.Error("book ticker parse err", err, respTemp)
//				cancel()
//				return depthCache, err
//			}
//			bids, err := DepthLevelParse(respTemp.Data[0].B)
//			if err != nil {
//				logger.Logger.Error("book ticker parse err", err, respTemp)
//				cancel()
//				return depthCache, err
//			}
//			depthCache = &base.OrderBook{
//				Exchange:     common.Exchange_BYBIT,
//				Market:       common.Market_SPOT,
//				Type:         common.SymbolType_SPOT_NORMAL,
//				Symbol:       symbolNameMap[respTemp.Symbol],
//				TimeExchange: uint64(respTemp.Data[0].T * 1000),
//				TimeReceive:  uint64(time.Now().UnixMicro()),
//				Asks:         asks,
//				Bids:         bids,
//			}
//			break LOOP
//		}
//	}
//	cancel()
//
//	sort.Sort(depthCache.Asks)
//	sort.Sort(sort.Reverse(depthCache.Bids))
//
//	return depthCache, nil
//}
