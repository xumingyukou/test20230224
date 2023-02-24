package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/kucoin"
	"clients/exchange/cex/kucoin/spot_api"
	"clients/logger"
	"clients/transform"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
)

type KucoinSpotWebsocket struct {
	base.WsConf
	apiClient     *spot_api.ApiClient
	listenTimeout int64
	pingTimeout   int64
	pingInterval  int64
	WsReqUrl      *spot_api.WsReqUrl
	handler       *WebSocketSpotHandle
	reqId         int64
	token         string
	hasToken      bool
}

func NewKucoinWebsocket(conf base.WsConf) *KucoinSpotWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}

	k := &KucoinSpotWebsocket{
		WsConf: conf,
	}
	logger.Logger.Info("new socket ", k)
	if k.ReadTimeout == 0 {
		k.ReadTimeout = 300
	}
	if k.listenTimeout == 0 {
		k.listenTimeout = 1800
	}
	k.apiClient = spot_api.NewApiClient(conf.APIConf)

	k.handler = NewWebSocketSpotHandle(k.ChanCap)
	k.WsReqUrl = spot_api.NewWsReqUrl()

	return k
}

func NewKucoinWebsocket2(conf base.WsConf, cli *http.Client) *KucoinSpotWebsocket {

	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}

	k := &KucoinSpotWebsocket{
		WsConf: conf,
	}
	if k.ReadTimeout == 0 {
		k.ReadTimeout = 300
	}
	if k.listenTimeout == 0 {
		k.listenTimeout = 1800
	}
	k.apiClient = spot_api.NewApiClient2(conf.APIConf, cli)

	k.handler = NewWebSocketSpotHandle(k.ChanCap)
	k.WsReqUrl = spot_api.NewWsReqUrl()

	return k
}

func (k *KucoinSpotWebsocket) requestToken(force bool) error {
	if !k.hasToken || force {
		resp, err := k.apiClient.GetWsToken(k.WsConf.IsPrivate)
		if err != nil {
			logger.Logger.Error("get token:", err.Error())
			return errors.New("request websocket token " + err.Error())
		}
		k.token = resp.Data.Token
		k.WsReqUrl.WS_BASE_URL = resp.Data.InstanceServers[0].Endpoint
		k.pingInterval = resp.Data.InstanceServers[0].PingInterval
		k.pingTimeout = resp.Data.InstanceServers[0].PingTimeout
		k.hasToken = true
	}
	return nil
}

func (k *KucoinSpotWebsocket) preReconnect(connectId int64, topic string) func(*conn.WsConn, int) error {
	return func(wsClient *conn.WsConn, retry int) error {
		// 每十次重新请求token
		if retry%10 == 1 {
			originToken := k.token
			if err := k.requestToken(true); err != nil {
				return err
			}
			wsClient.WsUrl = k.GetUrl(connectId)
			logger.Logger.Infof("Kucoin reconnect: origintoken:%s curtoken:%s , retry:%d", originToken, k.token, retry)
			logger.Logger.Infof("Kucoin reconnect topic:%s", topic)
		}
		return nil
	}
}

func (k *KucoinSpotWebsocket) GetUrl(connectId int64) string {
	return fmt.Sprintf("%s?token=%s&connectId=%d", k.WsReqUrl.WS_BASE_URL, k.token, connectId)
}

func (k *KucoinSpotWebsocket) EstablishConn(id int64, subAddr string, symbols []*client.SymbolInfo, handler func([]byte) error, ctx context.Context) error {
	var (
		readTimeout = time.Duration(k.ReadTimeout) * time.Second * 1000
	)
	heartBeat := func() []byte {
		data, _ := json.Marshal(PingMessage{Id: fmt.Sprintf("%d", kucoin.GetCurrentTimestamp()), Type: "ping"})
		return data
	}
	topic := k.makeTopic(subAddr, symbols)

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(k.GetUrl(id)).
		ProtoHandleFunc(handler).
		AutoReconnect().
		Heartbeat(heartBeat, time.Duration(k.pingInterval)*time.Millisecond).
		PreReconnect(k.preReconnect(id, topic)).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			return k.Subscribe(id, wsClient, topic)
		})

	if k.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(k.ProxyUrl)
	}

	dataSendStatus := k.handler.GetChan("send_status").(*base.CheckDataSendStatus)
	dataSendStatus.Init(3600, ctx, symbols...)

	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err := k.Subscribe(id, wsClient, topic)
	if err != nil {
		return err
	}
	go k.start(wsClient, topic, id, ctx)
	go k.resubscribeAfterSilence(id, wsClient, subAddr, topic, ctx, 0, dataSendStatus.UpdateTimeoutChMap)
	return err
}

func (k *KucoinSpotWebsocket) EstablishSnapshotConn(id int64, subAddr string, symbols []*client.SymbolInfo, handler func([]byte) error, ctx context.Context, conf *base.IncrementDepthConf) error {
	var (
		readTimeout = time.Duration(k.ReadTimeout) * time.Second * 1000
	)

	heartBeat := func() []byte {
		data, _ := json.Marshal(PingMessage{Id: fmt.Sprintf("%d", kucoin.GetCurrentTimestamp()), Type: "ping"})
		return data
	}
	topic := k.makeTopic(subAddr, symbols)

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(k.GetUrl(id)).
		ProtoHandleFunc(handler).
		AutoReconnect().
		Heartbeat(heartBeat, time.Duration(k.pingInterval)*time.Millisecond).
		PreReconnect(k.preReconnect(id, topic)).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			return k.Subscribe(id, wsClient, topic)
		})

	if k.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(k.ProxyUrl)
	}

	k.handler.SetDepthIncrementSnapShotConf(symbols, conf)
	dataSendStatus := k.handler.GetChan("send_status").(*base.CheckDataSendStatus)
	dataSendStatus.Init(3600, ctx, symbols...)

	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err := k.Subscribe(id, wsClient, topic)
	if err != nil {
		return err
	}
	go k.start(wsClient, topic, id, ctx)
	go k.resubscribeAfterSilence(id, wsClient, subAddr, topic, ctx, 0, dataSendStatus.UpdateTimeoutChMap)
	return err
}

func (k *KucoinSpotWebsocket) resubscribeAfterSilence(id int64, wsClient *conn.WsConn, subAddr string, topic string, ctx context.Context, interval int, ch chan []*client.SymbolInfo) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := k.UnSubscribe(id, wsClient, topic); err != nil {
				logger.Logger.Error("unsubscribe err:", topic, id)
			}
			break LOOP
		case resubSymbols := <-ch:
			resubTopic := k.makeTopic(subAddr, resubSymbols)
			if err := k.UnSubscribe(id, wsClient, resubTopic); err != nil {
				logger.Logger.Error("unsubscribe err:", resubTopic, id)
			}
			if err := k.Subscribe(id, wsClient, resubTopic); err != nil {
				logger.Logger.Error("subscribe err:", resubTopic, id)
			}
		}
	}
}

func (k *KucoinSpotWebsocket) start(wsClient *conn.WsConn, param string, id int64, ctx context.Context) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := k.UnSubscribe(id, wsClient, param); err != nil {
				logger.Logger.Error("unsubscribe err:", param, id)
			}
			break LOOP
		}
	}
	wsClient.CloseWs()
}

func (k *KucoinSpotWebsocket) Subscribe(id int64, client *conn.WsConn, topic string) error {
	fmt.Println("kucoin subscribe:", topic)
	err := client.Subscribe(SubscribeMessage{
		Id:             fmt.Sprintf("%d", id),
		Type:           "subscribe",
		Topic:          topic,
		PrivateChannel: false,
		Response:       true,
	})
	if err != nil {
		logger.Logger.Error("subscribe err:", err)
	}
	return err
}

func (k *KucoinSpotWebsocket) UnSubscribe(reqId int64, wsClient *conn.WsConn, topic string) (err error) {
	err = wsClient.Subscribe(UnSubscribeMessage{
		Id:             fmt.Sprintf("%d", reqId),
		Type:           "unsubscribe",
		Topic:          topic,
		PrivateChannel: false,
		Response:       false,
	})
	if err != nil {
		logger.Logger.Error("unsubscribe err:", err)
	}
	return
}

func (k *KucoinSpotWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
	if len(chMap) == 0 {
		return errors.New("no symbol channel")
	}
	var err error
	err = k.requestToken(false)
	if err != nil {
		return err
	}

	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	atomic.AddInt64(&k.reqId, 1)
	k.handler.SetTradeGroupChannel(chMap)
	err = k.EstablishConn(atomic.LoadInt64(&k.reqId), k.WsReqUrl.MARKET_MATCH_TOPIC, symbols, k.handler.GetHandler("tradeGroup"), ctx)
	return err
}

func (k *KucoinSpotWebsocket) makeTopic(baseUrl string, symbols []*client.SymbolInfo) string {
	topic := baseUrl + ":"
	for i, symbol := range symbols {
		if i == 0 {
			topic += kucoin.Canonical2Exchange(symbol.Symbol)
		} else {
			topic += "," + kucoin.Canonical2Exchange(symbol.Symbol)
		}
	}
	return topic
}

func (k *KucoinSpotWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
	if len(chMap) == 0 {
		return errors.New("no symbol channel")
	}
	var err error
	err = k.requestToken(false)
	if err != nil {
		return err
	}

	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	atomic.AddInt64(&k.reqId, 1)
	k.handler.SetBookTickerGroupChannel(chMap)
	err = k.EstablishConn(atomic.LoadInt64(&k.reqId), k.WsReqUrl.MARKET_TICKER_TOPIC, symbols, k.handler.GetHandler("bookTickerGroup"), ctx)
	return err
}

func (k *KucoinSpotWebsocket) DepthLimitGroup(ctx context.Context, interval int, limit int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	if len(chMap) == 0 {
		return errors.New("no symbol channel")
	}
	if interval != 0 && interval != 100 {
		return errors.New("only internal = 100ms is supported")
	}
	var err error
	err = k.requestToken(false)
	if err != nil {
		return err
	}

	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	atomic.AddInt64(&k.reqId, 1)
	k.handler.SetDepthLimitGroupChannel(chMap)
	err = k.EstablishConn(atomic.LoadInt64(&k.reqId), k.WsReqUrl.SPOTMARKET_LEVEL2_DEPTH_TOPIC, symbols, k.handler.GetHandler("depthLimitGroup"), ctx)
	return err
}

func (k *KucoinSpotWebsocket) DepthIncrementGroup(ctx context.Context, interval int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {
	if len(chMap) == 0 {
		return errors.New("no symbol channel")
	}
	if interval != 0 && interval != 100 {
		return errors.New("only internal = 100ms is supported")
	}
	var err error
	err = k.requestToken(false)
	if err != nil {
		return err
	}

	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}

	atomic.AddInt64(&k.reqId, 1)
	k.handler.SetDepthIncrementGroupChannel(chMap)
	err = k.EstablishConn(atomic.LoadInt64(&k.reqId), k.WsReqUrl.MARKET_LEVEL2, symbols, k.handler.GetHandler("depthIncrementGroup"), ctx)
	return err
}

func (k *KucoinSpotWebsocket) GetDepthGivenLimit(symbol string, limit int) (*base.OrderBook, error) {
	resp, err := k.apiClient.GetMarkets(kucoin.Canonical2Exchange(symbol), limit)
	timeReceive := kucoin.GetCurrentTimestamp()
	if err != nil {
		return nil, err
	}
	sequence, err := transform.Str2Int64(resp.Data.Sequence)
	if err != nil {
		return nil, err
	}
	asks, err := DepthLevelParse(resp.Data.Asks)
	if err != nil {
		return nil, err
	}
	bids, err := DepthLevelParse(resp.Data.Bids)
	if err != nil {
		return nil, err
	}
	asksSlice := base.DepthItemSlice(asks)
	bidsSlice := base.DepthItemSlice(bids)

	res := &base.OrderBook{
		Exchange:     common.Exchange_KUCOIN,
		Symbol:       symbol,
		Market:       common.Market_SPOT,
		Type:         common.SymbolType_SPOT_NORMAL,
		TimeReceive:  uint64(timeReceive),
		TimeExchange: resp.Data.Time * 1000,
		UpdateId:     sequence,
		Asks:         asksSlice,
		Bids:         bidsSlice,
	}
	return res, nil
}

func (k *KucoinSpotWebsocket) GetFullDepth(symbol string) (*base.OrderBook, error) {
	return k.GetDepthGivenLimit(symbol, 1000)
}

func (k *KucoinSpotWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, interval int, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error {
	var symbols []*client.SymbolInfo
	if isDelta && isFull {
		// check delta and snapshot
		if len(chDeltaMap) != len(chFullMap) || chDeltaMap == nil || chFullMap == nil {
			return errors.New("symbols for delmap and fullmap are not equal")
		}
		for symbol := range chDeltaMap {
			if _, ok := chFullMap[symbol]; !ok {
				return errors.New("symbols for delmap and fullmap are not equal")
			}
		}
	}
	if isDelta {
		for symbol := range chDeltaMap {
			symbols = append(symbols, symbol)
		}
	} else if isFull {
		for symbol := range chFullMap {
			symbols = append(symbols, symbol)
		}
	} else {
		return errors.New("no symbol channel")
	}

	var err error
	err = k.requestToken(false)
	if err != nil {
		return err
	}

	atomic.AddInt64(&k.reqId, 1)
	k.handler.SetDepthIncrementSnapshotGroupChannel(chDeltaMap, chFullMap)

	conf := &base.IncrementDepthConf{
		IsPublishDelta:    isDelta,
		IsPublishFull:     isFull,
		CheckTimeSec:      3000 + rand.Intn(1200),
		DepthCapLevel:     1000,
		DepthLevel:        limit,
		GetFullDepth:      k.GetFullDepth,
		GetFullDepthLimit: k.GetDepthGivenLimit,
		Ctx:               ctx,
	}
	err = k.EstablishSnapshotConn(atomic.LoadInt64(&k.reqId), k.WsReqUrl.MARKET_LEVEL2, symbols, k.handler.GetHandler("depthIncrementSnapshotGroup"), ctx, conf)
	return err

}

func (k *KucoinSpotWebsocket) FundingRateGroup(ctx context.Context, chFundMap map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	return nil
}
