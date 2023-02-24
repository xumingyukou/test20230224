package spot_ws

import (
	"clients/conn"
	"clients/crypto"
	"clients/exchange/cex/base"
	"clients/exchange/cex/huobi/c_api"
	"clients/exchange/cex/huobi/spot_api"
	"clients/exchange/cex/huobi/u_api"
	"clients/logger"
	"clients/transform"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/common"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
)

var (
	symbolNameMap      sync.Map
	contractSizeMap    sync.Map
	getSpotSymbolOrNot bool
)

type HuobiSpotWebsocket struct {
	base.WsConf
	apiClient     *spot_api.SpotApiClient
	UApiClient    *u_api.UApiClient
	CApiClient    *c_api.CApiClient
	readTimeout   int64 // can delete
	listenTimeout int64 // second
	lock          sync.Mutex
	reqId         int
	isStart       bool
	WsReqUrl      *spot_api.WsReqUrl
	handler       *WebSocketSpotHandle
}

func NewHuobiWebsocket(conf base.WsConf, endPoint string) *HuobiSpotWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &HuobiSpotWebsocket{
		WsConf: conf,
	}
	if conf.EndPoint == "" {
		d.EndPoint = endPoint
	}
	if d.ReadTimeout == 0 {
		d.ReadTimeout = 300
	}
	if d.listenTimeout == 0 {
		d.listenTimeout = 1800
	}

	d.handler = NewWebSocketSpotHandle(d.ChanCap)
	d.handler.Exchange = common.Exchange_HUOBI
	d.handler.Market = common.Market_SPOT

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

	for !getSpotSymbolOrNot {
		getSpotSymbolOrNot = true
		logger.Logger.Info("map initializing......")
		data, err := d.apiClient.ExchangeInfo()
		if err != nil {
			getSpotSymbolOrNot = false
			logger.Logger.Error("map initialize error:", err)
			continue
		}
		for _, pair := range data.Data {
			if pair.State == "online" {
				symbolNameMap.Store(pair.Symbol, strings.ToUpper(pair.BaseCurrency+"/"+pair.QuoteCurrency))
			}
		}
	}

	return d
}

func (b *HuobiSpotWebsocket) Subscribe(wsClient *conn.WsConn, subs []Sub) error {
	var (
		err error
	)
	for _, sub := range subs {
		logger.Logger.Info("Huobi subscribe:", sub.Sub)
		err = wsClient.Subscribe(sub)
		if err != nil {
			logger.Logger.Error("subscribe err:", err)
		}
	}
	return err
}

func (b *HuobiSpotWebsocket) SubscribeFund(wsClient *conn.WsConn, subs []SubFund) error {
	var (
		err error
	)
	for _, sub := range subs {
		logger.Logger.Info("Huobi subscribe:", sub.Topic)
		err = wsClient.Subscribe(sub)
		if err != nil {
			logger.Logger.Error("subscribe err:", err)
		}
	}
	return err
}

// func (b *HuobiSpotWebsocket) SubscribeSnapshotCoon(wsClient *conn.WsConn, subs []Sub) error {
// 	var (
// 		err error
// 	)
// 	for _, sub := range subs {
// 		fmt.Println("Huobi subscribe:", sub.Sub)
// 		req := Req{
// 			Req: sub.Sub,
// 			ID:  sub.ID,
// 		}
// 		err = wsClient.Subscribe(req)
// 		time.Sleep(time.Duration(110) * time.Millisecond)
// 		if err != nil {
// 			logger.Logger.Error("request err:", err)
// 		}
// 		err = wsClient.Subscribe(sub)
// 		if err != nil {
// 			logger.Logger.Error("subscribe err:", err)
// 		}
// 	}
// 	return err
// }

func (b *HuobiSpotWebsocket) Unsubscribe(wsClient *conn.WsConn, subs []Sub) error {
	var (
		err error
	)
	for _, unsub := range subs {
		fmt.Println("Huobi unsubscribe", unsub)
		err = wsClient.Subscribe(Unsub{
			Unsub: unsub.Sub,
			ID:    unsub.ID,
		})
		if err != nil {
			logger.Logger.Error("unsubscribe err:", err)
		}
	}
	return err
}

func (b *HuobiSpotWebsocket) UnsubscribeFund(wsClient *conn.WsConn, subs []SubFund) error {
	var (
		err error
	)
	for _, unsub := range subs {
		fmt.Println("Huobi unsubscribe", unsub)
		err = wsClient.Subscribe(Unsub{
			Unsub: unsub.Topic,
			ID:    unsub.CID,
		})
		if err != nil {
			logger.Logger.Error("unsubscribe err:", err)
		}
	}
	return err
}

func (b *HuobiSpotWebsocket) Start(wsClient *conn.WsConn, subs []Sub, ctx context.Context) {
LOOP:
	for {
		select {
		case pong := <-b.handler.GetChan("pingpong").(chan int64):
			pongtoJson := "{\"pong\":" + strconv.FormatInt(pong, 10) + "}"
			// fmt.Println(pongtoJson)
			wsClient.SendMessage([]byte(pongtoJson))
		case <-ctx.Done():
			wsClient.CloseWs()
			if err := b.Unsubscribe(wsClient, subs); err != nil {
				logger.Logger.Error("unsubscribe err:", subs)
			}
			break LOOP
		}
	}
}

func (b *HuobiSpotWebsocket) StartFund(wsClient *conn.WsConn, subs []SubFund, ctx context.Context) {
LOOP:
	for {
		select {
		case pong := <-b.handler.GetChan("contractpingpong").(chan int64):
			pongtoJson := "{\"op\":" + "pong" + "," + "\"ts\":" + strconv.FormatInt(pong, 10) + "}"
			wsClient.SendMessage([]byte(pongtoJson))
		case pong := <-b.handler.GetChan("pingpong").(chan int64):
			pongtoJson := "{\"pong\":" + strconv.FormatInt(pong, 10) + "}"
			// fmt.Println(pongtoJson)
			wsClient.SendMessage([]byte(pongtoJson))
		case <-ctx.Done():
			wsClient.CloseWs()
			if err := b.UnsubscribeFund(wsClient, subs); err != nil {
				logger.Logger.Error("unsubscribe err:", subs)
			}
			break LOOP
		}
	}
}

func (b *HuobiSpotWebsocket) ResubscribeAfterSilence(wsClient *conn.WsConn, subs []Sub, ctx context.Context, ch chan []*client.SymbolInfo) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := b.Unsubscribe(wsClient, subs); err != nil {
				logger.Logger.Error("unsubscribe err:", subs)
			}
			wsClient.CloseWs()
			break LOOP
		case reconnectData := <-ch:
			var reconnectSubs []Sub
			for _, symbol := range reconnectData {
				reconnectSub := Sub{
					Sub: strings.Replace(subs[0].Sub, strings.Split(subs[0].Sub, ".")[1], GetSymbolName(symbol), -1),
					ID:  strconv.Itoa(b.reqId),
				}
				reconnectSubs = append(reconnectSubs, reconnectSub)
			}
			if err := b.Unsubscribe(wsClient, reconnectSubs); err != nil {
				logger.Logger.Error("unsubscribe err:", reconnectSubs)
			}
			if err := b.Subscribe(wsClient, reconnectSubs); err != nil {
				logger.Logger.Error("unsubscribe err:", reconnectSubs)
			}
		}
	}
}
func (b *HuobiSpotWebsocket) ResubscribeAfterSilenceFund(wsClient *conn.WsConn, subs []SubFund, ctx context.Context, ch chan []*client.SymbolInfo) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := b.UnsubscribeFund(wsClient, subs); err != nil {
				logger.Logger.Error("unsubscribe err:", subs)
			}
			wsClient.CloseWs()
			break LOOP
		case reconnectData := <-ch:
			var reconnectSubs []Sub
			for _, symbol := range reconnectData {
				reconnectSub := Sub{
					Sub: strings.Replace(subs[0].Topic, strings.Split(subs[0].Topic, ".")[1], GetSymbolName(symbol), -1),
					ID:  strconv.Itoa(b.reqId),
				}
				reconnectSubs = append(reconnectSubs, reconnectSub)
			}
			if err := b.Unsubscribe(wsClient, reconnectSubs); err != nil {
				logger.Logger.Error("unsubscribe err:", reconnectSubs)
			}
			if err := b.Subscribe(wsClient, reconnectSubs); err != nil {
				logger.Logger.Error("unsubscribe err:", reconnectSubs)
			}
		}
	}
}

func (b *HuobiSpotWebsocket) EstablishConn(subs []Sub, url string, handler func([]byte) error, ctx context.Context, symbols []*client.SymbolInfo) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
	)
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).
		AutoReconnect().
		ProtoHandleFunc(handler).
		DecompressFunc(crypto.GZIPDecompress).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			return b.Subscribe(wsClient, subs)
		})
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = b.Subscribe(wsClient, subs)
	if err != nil {
		return err
	}
	go b.Start(wsClient, subs, ctx)
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.ResubscribeAfterSilence(wsClient, subs, ctx, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)
	return err
}

func (b *HuobiSpotWebsocket) EstablishConnFund(subs []SubFund, url string, handler func([]byte) error, ctx context.Context, symbols []*client.SymbolInfo) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
	)
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).
		AutoReconnect().
		ProtoHandleFunc(handler).
		DecompressFunc(crypto.GZIPDecompress).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			return b.SubscribeFund(wsClient, subs)
		})
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = b.SubscribeFund(wsClient, subs)
	if err != nil {
		return err
	}
	go b.StartFund(wsClient, subs, ctx)
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.ResubscribeAfterSilenceFund(wsClient, subs, ctx, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)
	return err
}

func (b *HuobiSpotWebsocket) EstablishSnapshotConn(symbols []*client.SymbolInfo, subs []Sub, url string,
	handler func([]byte) error, ctx context.Context, conf *base.IncrementDepthConf) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
		reqSymbol   chan *client.SymbolInfo
	)
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).
		AutoReconnect().
		ProtoHandleFunc(handler).
		DecompressFunc(crypto.GZIPDecompress).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			return b.Subscribe(wsClient, subs)
		})
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}

	if b.handler.DepthIncrementSnapshotReqSymbol == nil {
		reqSymbol = make(chan *client.SymbolInfo)
		b.handler.SetDepthIncrementSnapshotReqChan(reqSymbol)
	} else {
		reqSymbol = b.handler.DepthIncrementSnapshotReqSymbol
	}

	b.handler.SetDepthIncrementSnapShotConf(symbols, conf)
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}

	err = b.Subscribe(wsClient, subs)
	if err != nil {
		return err
	}
	go b.Start(wsClient, subs, ctx)
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.ResubscribeAfterSilence(wsClient, subs, ctx, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)
	go b.SendReq(wsClient, ctx, reqSymbol, conf.DepthLevel)
	return err
}

func (b *HuobiSpotWebsocket) SendReq(wsClient *conn.WsConn, ctx context.Context, reqSymbol chan *client.SymbolInfo, level int) {
	reqCount := make(map[string]int)
	lastReqTime := make(map[string]int64)
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case symbol := <-reqSymbol:
			reqTime, time_ok := lastReqTime[base.SymInfoToString(symbol)]
			// 5s 内只发送一次
			if !time_ok || time.Now().Sub(time.UnixMicro(reqTime)) > time.Duration(5)*time.Second {
				logger.Logger.Info("request full depth ", base.SymInfoToString(symbol))
				id, id_ok := reqCount[base.SymInfoToString(symbol)]
				if !id_ok {
					id = 0
				} else {
					id++
				}

				reqFullMsg := GetDepthFullSubscription(symbol, level, id)
				errFull := wsClient.Subscribe(reqFullMsg)
				if errFull == nil {
					lastReqTime[base.SymInfoToString(symbol)] = time.Now().UnixMicro()
					reqCount[base.SymInfoToString(symbol)] = id
				} else {
					logger.Logger.Error("send request ", base.SymInfoToString(symbol), " ", errFull)
				}

				//reqIncrementMsg := GetDepthIncrementSubscription(symbol, level, id)
				//errInc := wsClient.Subscribe(reqIncrementMsg)
				//if errInc == nil {
				//	lastReqTime[base.SymInfoToString(symbol)] = time.Now().UnixMicro()
				//	reqCount[base.SymInfoToString(symbol)] = id
				//} else {
				//	logger.Logger.Error("send request ", base.SymInfoToString(symbol), " ", errInc)
				//}
			}
		}
	}
}

func (b *HuobiSpotWebsocket) FundingRateGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	var (
		subsC, subsU []SubFund
		err          error
		symbols      []*client.SymbolInfo
	)
	for symbol := range chMap {
		param := SubFund{
			Op:    "sub",
			Topic: GetFundSubbscription(symbol),
			CID:   strconv.Itoa(b.reqId),
		}
		if strings.Contains(symbol.Symbol, "USDT") {
			subsU = append(subsU, param)
		} else {
			subsC = append(subsC, param)
		}
	}
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	b.reqId++
	b.handler.SetFundingRateGroupChannel(chMap)
	if len(subsC) > 0 {
		err = b.EstablishConnFund(subsC, "wss://api.hbdm.com/swap-notification", b.handler.FundingRateGroupHandle, ctx, symbols)
	}
	if len(subsU) > 0 {
		err = b.EstablishConnFund(subsU, "wss://api.hbdm.com/linear-swap-notification", b.handler.FundingRateGroupHandle, ctx, symbols)
	}
	return err
}

// 交易数据
func (b *HuobiSpotWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
	var (
		subs    []Sub
		err     error
		symbols []*client.SymbolInfo
	)
	for symbol := range chMap {
		param := Sub{
			Sub: GetTradeSubscription(symbol),
			ID:  transform.XToString(b.reqId),
		}
		subs = append(subs, param)
	}
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	b.reqId++
	b.handler.SetTradeGroupChannel(chMap)
	err = b.EstablishConn(subs, b.EndPoint, b.handler.TradeGroupHandle, ctx, symbols)
	return err
}

func (b *HuobiSpotWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
	var (
		subs    []Sub
		err     error
		symbols []*client.SymbolInfo
		market  common.Market
	)
	for symbol := range chMap {
		market = symbol.Market
		param := Sub{
			Sub: GetTickerSubscription(symbol),
			ID:  strconv.Itoa(b.reqId),
		}
		subs = append(subs, param)
	}
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	b.reqId++
	b.handler.SetBookTickerGroupChannel(chMap)
	var handlerFunc func([]byte) error
	if market == common.Market_SPOT {
		handlerFunc = b.handler.BookTickerGroupHandle
	} else {
		handlerFunc = b.handler.FutureOrSwapBookTickerGroupHandle
	}
	err = b.EstablishConn(subs, b.EndPoint, handlerFunc, ctx, symbols)
	return err
}

func (b *HuobiSpotWebsocket) DepthLimitGroup(ctx context.Context, interval, limit int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	var (
		subs    []Sub
		err     error
		symbols []*client.SymbolInfo
	)
	if limit != 5 && limit != 10 && limit != 20 {
		limit = 20
	}
	for symbol := range chMap {
		sub := GetDepthLimitSubscription(symbol, limit, b.reqId)
		subs = append(subs, sub)
	}
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	b.reqId++
	b.handler.SetDepthLimitGroupChannel(chMap)
	err = b.EstablishConn(subs, b.EndPoint, b.handler.DepthLimitGroupHandle, ctx, symbols)
	return err
}

func (b *HuobiSpotWebsocket) DepthIncrementGroup(ctx context.Context, interval int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {
	var (
		subs    []Sub
		err     error
		symbols []*client.SymbolInfo
	)
	for symbol := range chMap {
		param := GetDepthIncrementSubscription(symbol, 150, b.reqId)
		subs = append(subs, param)
	}
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	b.reqId++
	b.handler.SetDepthIncrementGroupChannel(chMap)
	err = b.EstablishConn(subs, spot_api.WS_API_BASE_URL, b.handler.DepthIncrementGroupHandle, ctx, symbols)
	return err
}

func (b *HuobiSpotWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, interval, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error {
	var (
		subs    []Sub
		symbols []*client.SymbolInfo
		err     error
		market  common.Market
	)
	if isDelta && isFull && len(chDeltaMap) != len(chFullMap) {
		return errors.New("chDeltaMap and chFullMap not match")
	}
	if limit >= 400 {
		limit = 400
	} else {
		limit = 150
	}
	if isDelta {
		for symbol := range chDeltaMap {
			market = symbol.Market
			symbols = append(symbols, symbol)
			param := GetDepthIncrementSubscription(symbol, limit, b.reqId)
			subs = append(subs, param)
		}
	} else {
		for symbol := range chFullMap {
			market = symbol.Market
			symbols = append(symbols, symbol)
			param := GetDepthIncrementSubscription(symbol, limit, b.reqId)
			subs = append(subs, param)
		}
	}
	b.reqId++
	b.handler.SetDepthIncrementSnapshotGroupChannel(chDeltaMap, chFullMap)
	conf := &base.IncrementDepthConf{
		Exchange:       b.handler.Exchange,
		Market:         b.handler.Market,
		IsPublishDelta: isDelta,
		IsPublishFull:  isFull,
		CheckTimeSec:   3600,
		DepthCapLevel:  1000,
		DepthLevel:     limit,
		Ctx:            ctx,
	}
	var handlerFunc func([]byte) error
	if market == common.Market_SPOT {
		handlerFunc = b.handler.DepthIncrementSnapShotGroupHandle
	} else {
		handlerFunc = b.handler.FutureOrSwapDepthIncrementGroupHandle
	}
	err = b.EstablishSnapshotConn(symbols, subs, b.EndPoint, handlerFunc, ctx, conf)
	return err
}

func NewHuobiWs(market common.Market, wsConf base.WsConf, streamType common.StreamType) base.CexWebsocketPublicInterface {
	var huobiWs base.CexWebsocketPublicInterface
	var endPoint string
	switch market {
	case common.Market_SPOT:
		if streamType == common.StreamType_MARKET_DEPTH {
			endPoint = spot_api.WS_API_MBP_INCREMENT
		} else {
			endPoint = spot_api.WS_API_BASE_URL
		}
		huobiWs = NewHuobiWebsocket(wsConf, endPoint)
	case common.Market_SWAP_COIN:
		endPoint = spot_api.WS_CBASE_SWAP_URL
		huobiWs = NewHuobiCBaseSwapWebsocket(wsConf, endPoint)
	case common.Market_FUTURE_COIN:
		endPoint = spot_api.WS_CBASE_FUTURE_URL
		huobiWs = NewHuobiCBaseFutureWebsocket(wsConf, endPoint)
	case common.Market_SWAP:
		endPoint = spot_api.WS_UBASE_URL
		huobiWs = NewHuoBiUBaseSwapWebsocket(wsConf, endPoint)
	case common.Market_FUTURE:
		endPoint = spot_api.WS_UBASE_URL
		huobiWs = NewHuobiUBaseFutureWebsocket(wsConf, endPoint)
	}
	return huobiWs
}
