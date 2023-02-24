package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/bithumb"
	"clients/exchange/cex/bithumb/spot_api"
	"clients/logger"
	"clients/transform"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
)

type BithumbSpotWebsocket struct {
	base.WsConf
	apiClient     *spot_api.ApiClient
	listenTimeout int64
	WsReqUrl      *spot_api.WsReqUrl
	handler       *WebSocketSpotHandle
	pingInterval  int64
	reqId         int
	token         string
}

func NewBithumbWebsocket(conf base.WsConf) *BithumbSpotWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	b := &BithumbSpotWebsocket{
		WsConf: conf,
	}
	if b.ReadTimeout == 0 {
		b.ReadTimeout = 300
	}
	if b.listenTimeout == 0 {
		b.listenTimeout = 1800
	}
	if b.pingInterval == 0 {
		b.pingInterval = 30
	}
	b.apiClient = spot_api.NewApiClient(conf.APIConf)

	b.handler = NewWebSocketSpotHandle(b.ChanCap)

	b.WsReqUrl = spot_api.NewWsReqUrl()

	return b
}

func NewBithumbWebsocket2(conf base.WsConf, cli *http.Client) *BithumbSpotWebsocket {

	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}

	b := &BithumbSpotWebsocket{
		WsConf: conf,
	}
	if b.ReadTimeout == 0 {
		b.ReadTimeout = 300
	}
	if b.listenTimeout == 0 {
		b.listenTimeout = 1800
	}
	if b.pingInterval == 0 {
		b.pingInterval = 30
	}
	b.apiClient = spot_api.NewApiClient2(conf.APIConf, cli)

	b.handler = NewWebSocketSpotHandle(b.ChanCap)
	b.WsReqUrl = spot_api.NewWsReqUrl()

	return b
}

func (b *BithumbSpotWebsocket) baseUrl() string {
	return b.WsReqUrl.WS_BASE_URL
}

func (b *BithumbSpotWebsocket) EstablishConn(url string, event string, symbolInfos []*client.SymbolInfo, handler func([]byte) error, ctx context.Context) error {
	heartBeat := func() []byte {
		return []byte("{\"cmd\":\"ping\"}")
	}

	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
	)

	topics := b.GetTopics(event, symbolInfos)

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).
		ProtoHandleFunc(handler).
		AutoReconnect().
		Heartbeat(heartBeat, time.Duration(b.pingInterval)*time.Second).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			return b.Subscribe(wsClient, topics)
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

	err := b.Subscribe(wsClient, topics)
	if err != nil {
		return err
	}
	go b.start(wsClient, ctx, topics)
	go b.resubscribeAfterSilence(wsClient, event, topics, ctx, 0, dataSendStatus.UpdateTimeoutChMap)
	return err
}

func (b *BithumbSpotWebsocket) EstablishSnapshotConn(url string, event string, symbolInfos []*client.SymbolInfo, handler func([]byte) error, ctx context.Context, conf *base.IncrementDepthConf) error {
	var (
		err error
	)
	heartBeat := func() []byte {
		return []byte("{\"cmd\":\"ping\"}")
	}

	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
	)

	topics := b.GetTopics(event, symbolInfos)

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).
		ProtoHandleFunc(handler).
		AutoReconnect().
		Heartbeat(heartBeat, time.Duration(b.pingInterval)*time.Second).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			return b.Subscribe(wsClient, topics)
		})
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}

	b.handler.SetDepthIncrementSnapShotConf(symbolInfos, conf)

	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = b.Subscribe(wsClient, topics)
	if err != nil {
		return err
	}

	dataSendStatus := b.handler.GetChan("send_status").(*base.CheckDataSendStatus)
	dataSendStatus.InitDecorator(b.reSubscribe, 3600, ctx, wsClient, event, symbolInfos...)
	go b.start(wsClient, ctx, topics)
	return err
}

func (b *BithumbSpotWebsocket) resubscribeAfterSilence(wsClient *conn.WsConn, event string, originTopics []string, ctx context.Context, interval int, ch chan []*client.SymbolInfo) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := b.UnSubscribe(wsClient, originTopics); err != nil {
				logger.Logger.Error("unsubscribe err:", originTopics)
			}
			break LOOP
		case resubSymbols := <-ch:
			b.reSubscribe(wsClient, resubSymbols, event)
		}
	}
}

func (b *BithumbSpotWebsocket) reSubscribe(wsClient *conn.WsConn, symbolInfoList []*client.SymbolInfo, event string) {
	resubTopic := b.GetTopics(event, symbolInfoList)
	if err := b.UnSubscribe(wsClient, resubTopic); err != nil {
		logger.Logger.Error("unsubscribe err:", resubTopic)
	}
	if err := b.Subscribe(wsClient, resubTopic); err != nil {
		logger.Logger.Error("subscribe err:", resubTopic)
	}
}

func (b *BithumbSpotWebsocket) start(wsClient *conn.WsConn, ctx context.Context, topics []string) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			wsClient.CloseWs()
			if err := b.UnSubscribe(wsClient, topics); err != nil {
				logger.Logger.Error("unsubscribe err")
			}
			break LOOP
		}
	}
}

func (b *BithumbSpotWebsocket) Subscribe(client *conn.WsConn, topics []string) error {
	logger.Logger.Info("Bithumb subscribe:", topics)
	err := client.Subscribe(CommandMessage{
		Cmd:  "subscribe",
		Args: topics,
	})
	if err != nil {
		logger.Logger.Error("subscribe err:", err)
	}
	return err
}

func (b *BithumbSpotWebsocket) UnSubscribe(wsClient *conn.WsConn, topics []string) (err error) {
	err = wsClient.Subscribe(CommandMessage{
		Cmd:  "unSubscribe",
		Args: topics,
	})
	if err != nil {
		logger.Logger.Error("unsubscribe err:", err)
	}
	return
}

func (k *BithumbSpotWebsocket) GetTopics(event string, symbols []*client.SymbolInfo) []string {
	topics := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		topics = append(topics, fmt.Sprintf("%s:%s", event, bithumb.Canonical2Exchange(symbol.Symbol)))
	}
	return topics
}

func (b *BithumbSpotWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {

	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}

	b.handler.SetTradeGroupChannel(chMap)
	err := b.EstablishConn(b.baseUrl(), "TRADE", symbols, b.handler.GetHandler("trade"), ctx)
	return err
}

func (b *BithumbSpotWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {

	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}

	b.handler.SetBookTickerGroupChannel(chMap)
	err := b.EstablishConn(b.baseUrl(), "TICKER", symbols, b.handler.GetHandler("ticker"), ctx)
	return err
}

func (b *BithumbSpotWebsocket) DepthIncrementGroup(ctx context.Context, interval int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {

	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}

	b.handler.SetDepthIncrementGroupChannel(chMap)
	err := b.EstablishConn(b.baseUrl(), "ORDERBOOK", symbols, b.handler.GetHandler("depthDiff"), ctx)
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

func (b *BithumbSpotWebsocket) GetFullDepth(symbol string) (*base.OrderBook, error) {
	resp, err := b.apiClient.GetOrderBook(bithumb.Canonical2Exchange(symbol))
	timeReceive := time.Now().UnixMicro()
	if err != nil {
		return nil, err
	}
	updateId, err := transform.Str2Int64(resp.Data.Ver)
	if err != nil {
		return nil, err
	}
	asks, err := DepthLevelParse(resp.Data.S)
	if err != nil {
		return nil, err
	}
	bids, err := DepthLevelParse(resp.Data.B)
	if err != nil {
		return nil, err
	}

	res := &base.OrderBook{
		Exchange:     common.Exchange_BITHUMB,
		Market:       common.Market_SPOT,
		Type:         common.SymbolType_SPOT_NORMAL,
		Symbol:       symbol,
		TimeExchange: uint64(resp.Timestamp) * 1000,
		TimeReceive:  uint64(timeReceive),
		Bids:         bids,
		Asks:         asks,
		UpdateId:     updateId,
	}
	return res, nil
}

func (b *BithumbSpotWebsocket) DepthLimitGroup(ctx context.Context, interval, limit int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	return errors.New("DepthLimitGroup not supported")
}

func (b *BithumbSpotWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, interval, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error {
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

	conf := &base.IncrementDepthConf{
		IsPublishDelta:    isDelta,
		IsPublishFull:     isFull,
		CheckTimeSec:      3000 + rand.Intn(1200),
		DepthCapLevel:     1000,
		DepthLevel:        limit,
		GetFullDepth:      b.GetFullDepth,
		GetFullDepthLimit: nil,
		Ctx:               ctx,
	}

	b.handler.SetDepthIncrementSnapshotGroupChannel(chDeltaMap, chFullMap)
	err := b.EstablishSnapshotConn(b.baseUrl(), "ORDERBOOK", symbols, b.handler.GetHandler("depthIncrementSnapshot"), ctx, conf)
	return err
}

func (k *BithumbSpotWebsocket) FundingRateGroup(ctx context.Context, chFundMap map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	return nil
}
