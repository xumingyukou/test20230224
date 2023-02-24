package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/binance/c_api"
	"clients/exchange/cex/binance/spot_api"
	"clients/exchange/cex/binance/u_api"
	"clients/logger"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/depth"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/order"
)

type GlobalDepthWight struct {
	GlobalLock   sync.Mutex
	GlobalWeight int64
}

var (
	g          *GlobalDepthWight
	globalOnce sync.Once
)

func init() {
	rand.Seed(time.Now().Unix())
}

type BinanceSpotWebsocket struct {
	base.WsConf
	ApiClient        *spot_api.ApiClient
	UApiClient       *u_api.UApiClient
	CApiClient       *c_api.CApiClient
	listenTimeout    int64 //second
	lock             sync.Mutex
	reqId            int
	isStart          bool
	WsReqUrl         *spot_api.WsReqUrl
	GlobalDW         *GlobalDepthWight
	handler          WebSocketHandleInterface
	GetFullDepthFunc func(symbol string) (*base.OrderBook, error)
}

type UserAction string

const (
	LISTEN UserAction = "listen"
	PUT    UserAction = "put"
	DELETE UserAction = "delete"
)

func NewBinanceWebsocket(conf base.WsConf) *BinanceSpotWebsocket {
	if conf.EndPoint == "" {
		conf.EndPoint = spot_api.WS_API_BASE_URL
	}
	d := &BinanceSpotWebsocket{
		WsConf: conf,
	}
	apiConf := base.APIConf{
		ProxyUrl:    d.ProxyUrl,
		ReadTimeout: d.ReadTimeout,
		AccessKey:   d.AccessKey,
		SecretKey:   d.SecretKey,
	}
	d.Init(spot_api.NewApiClient(apiConf), u_api.NewUApiClient(apiConf), c_api.NewCApiClient(apiConf), NewWebSocketSpotHandle(d.ChanCap), spot_api.NewSpotWsUrl(), d.GetFullDepth)
	d.GlobalDW = g
	globalOnce.Do(func() {
		if g == nil {
			g = new(GlobalDepthWight)
			d.GlobalDW = g
		}
		go func() {
			timer := time.NewTimer(time.Duration(d.ApiClient.WeightInfo[client.WeightType_REQUEST_WEIGHT].IntervalSec) * time.Second)
			for {
				select {
				case <-timer.C:
					d.GlobalDW.GlobalLock.Lock()
					d.GlobalDW.GlobalWeight = 0
					d.GlobalDW.GlobalLock.Unlock()
					timer.Reset(time.Duration(d.ApiClient.WeightInfo[client.WeightType_REQUEST_WEIGHT].IntervalSec) * time.Second)
				}
			}
		}()
	})

	return d
}

func NewBinanceWebsocket2(conf base.WsConf, cli *http.Client) *BinanceSpotWebsocket {
	d := &BinanceSpotWebsocket{
		WsConf: conf,
	}
	// ws socket目前固定写死
	d.EndPoint = spot_api.WS_API_BASE_URL
	apiConf := base.APIConf{
		ProxyUrl:    d.ProxyUrl,
		ReadTimeout: d.ReadTimeout,
		AccessKey:   d.AccessKey,
		SecretKey:   d.SecretKey,
	}
	d.Init(spot_api.NewApiClient2(apiConf, cli), u_api.NewUApiClient2(apiConf, cli), c_api.NewCApiClient2(apiConf, cli), NewWebSocketSpotHandle(d.ChanCap), spot_api.NewSpotWsUrl(), d.GetFullDepth)
	return d
}

func (b *BinanceSpotWebsocket) Init(apiClient *spot_api.ApiClient, uApiClient *u_api.UApiClient, cApiClient *c_api.CApiClient, handler WebSocketHandleInterface, wsReqUrl *spot_api.WsReqUrl, fullDepthFunc func(string) (*base.OrderBook, error)) {
	if b.ReadTimeout == 0 {
		b.ReadTimeout = 300
	}
	if b.ChanCap < 1 {
		b.ChanCap = 1024
	}
	if b.listenTimeout == 0 {
		b.listenTimeout = 1800
	}
	b.ApiClient = apiClient
	b.UApiClient = uApiClient
	b.CApiClient = cApiClient
	b.ApiClient.WeightInfo = make(map[client.WeightType]*client.WeightInfo)
	b.CApiClient.WeightInfo = make(map[client.WeightType]*client.WeightInfo)
	b.UApiClient.WeightInfo = make(map[client.WeightType]*client.WeightInfo)
	b.ApiClient.WeightInfo[client.WeightType_REQUEST_WEIGHT] = &client.WeightInfo{
		Type:        client.WeightType_REQUEST_WEIGHT,
		IntervalSec: 60,
	}
	b.CApiClient.WeightInfo[client.WeightType_REQUEST_WEIGHT] = &client.WeightInfo{
		Type:        client.WeightType_REQUEST_WEIGHT,
		IntervalSec: 60,
	}
	b.UApiClient.WeightInfo[client.WeightType_REQUEST_WEIGHT] = &client.WeightInfo{
		Type:        client.WeightType_REQUEST_WEIGHT,
		IntervalSec: 60,
	}

	b.handler = handler
	b.WsReqUrl = wsReqUrl
	b.GetFullDepthFunc = fullDepthFunc
}

func (b *BinanceSpotWebsocket) start(wsClient *conn.WsConn, param string, id int, ctx context.Context) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := b.UnSubscribe(wsClient, param, id); err != nil {
				logger.Logger.Error("unsubscribe err:", param, id)
			}
			wsClient.CloseWs()
			break LOOP
		}
	}
}

func (b *BinanceSpotWebsocket) reConnect2(wsClient *conn.WsConn, symbols []*client.SymbolInfo, interval, limit int, wsurl string, id int, ctx context.Context, ch chan []*client.SymbolInfo) {
	var (
		param = b.GetParam(wsurl, interval, limit, symbols...)
	)
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := b.UnSubscribe(wsClient, param, id); err != nil {
				logger.Logger.Error("unsubscribe err:", param, id)
			}
			wsClient.CloseWs()
			break LOOP
		case <-ch:
			if err := b.UnSubscribe(wsClient, param, id); err != nil {
				logger.Logger.Error("unsubscribe err:", param, id)
			}
			param = b.GetParam(wsurl, interval, limit, symbols...)
			if err := b.Subscribe(id, wsClient, param); err != nil {
				logger.Logger.Error("subscribe err:", param, id)
			}
		}
	}
}

func (b *BinanceSpotWebsocket) Subscribe(id int, client *conn.WsConn, params string) error {
	fmt.Println("binance subscribe:", params)
	err := client.Subscribe(req{
		Method: "SUBSCRIBE",
		Params: strings.Split(params, "/"),
		Id:     id,
	})
	if err != nil {
		logger.Logger.Error("subscribe err:", err)
	}
	return err
}

func (b *BinanceSpotWebsocket) ReSubscribe(id int, client *conn.WsConn, symbols []*client.SymbolInfo, interval, limit int, wsUrl string) error {
	params := b.GetParam(wsUrl, interval, limit, symbols...)
	return b.Subscribe(id, client, params)
}

func (b *BinanceSpotWebsocket) UnSubscribe(wsClient *conn.WsConn, params string, reqId int) (err error) {
	err = wsClient.Subscribe(req{
		Method: "UNSUBSCRIBE",
		Params: strings.Split(params, "/"),
		Id:     reqId,
	})
	if err != nil {
		logger.Logger.Error("unsubscribe err:", err)
	}
	return
}

func (b *BinanceSpotWebsocket) GetSymbolName(symbol *client.SymbolInfo) string {
	if spot_api.IsUBaseSymbolType(symbol.Type) {
		return spot_api.GetSpotSymbolName(u_api.GetUBaseSymbol(symbol.Symbol, u_api.GetFutureTypeFromNats(symbol.Type)))
	} else {
		return spot_api.GetSpotSymbolName(c_api.GetCBaseSymbol(symbol.Symbol, u_api.GetFutureTypeFromNats(symbol.Type)))
	}
}

func (b *BinanceSpotWebsocket) GetParam(url string, interval, limit int, symbols ...*client.SymbolInfo) string {
	params := ""
	switch url {
	case b.WsReqUrl.DEPTH_INCRE_URL:
		if limit == -1 { // increment
			for _, symbol := range symbols {
				params += strings.Replace(url, "<symbol>", b.GetSymbolName(symbol), -1)
				if interval > 0 {
					params += fmt.Sprint("@", interval, "ms")
				}
				params += "/"
			}
			return params[:len(params)-1]
		} else { // full
			if limit == 0 {
				limit = 20
			}
			for _, symbol := range symbols {
				params += fmt.Sprint(strings.Replace(url, "<symbol>", b.GetSymbolName(symbol), -1), limit)
				if interval > 0 {
					params += fmt.Sprint("@", interval, "ms")
				}
				params += "/"
			}
		}
		return params[:len(params)-1]
	default:
		for _, symbol := range symbols {
			params += strings.Replace(url, "<symbol>", b.GetSymbolName(symbol), -1) + "/"
		}
		return params[:len(params)-1]
	}
}

func (b *BinanceSpotWebsocket) GetUrl(url, param string) string {
	return b.EndPoint + url + param
}

func (b *BinanceSpotWebsocket) EstablishConn(id int, symbols []*client.SymbolInfo, interval, limit int, wsUrl, apiUrl string, handler func([]byte) error, ctx context.Context) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
		params      = b.GetParam(wsUrl, interval, limit, symbols...)
		url         = b.GetUrl(apiUrl, params)
	)
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(url).ProtoHandleFunc(handler).AutoReconnect().PostReconnectSuccess(func(wsClient *conn.WsConn) error { return b.ReSubscribe(id, wsClient, symbols, interval, limit, wsUrl) })
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = b.Subscribe(id, wsClient, params)
	if err != nil {
		return err
	}
	go b.start(wsClient, params, id, ctx)
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.reConnect2(wsClient, symbols, interval, limit, wsUrl, id, ctx, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)
	return err
}

func (b *BinanceSpotWebsocket) EstablishSnapshotConn(id int, symbols []*client.SymbolInfo, interval, limit int, wsUrl, apiUrl string, handler func([]byte) error, ctx context.Context, conf *base.IncrementDepthConf) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
		params      = b.GetParam(wsUrl, interval, limit, symbols...)
		url         = b.GetUrl(apiUrl, params)
	)
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(url).ProtoHandleFunc(handler).AutoReconnect().PostReconnectSuccess(func(wsClient *conn.WsConn) error { return b.ReSubscribe(id, wsClient, symbols, interval, limit, wsUrl) })
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	b.handler.SetDepthIncrementSnapShotConf(symbols, conf)
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = b.Subscribe(id, wsClient, params)
	if err != nil {
		return err
	}
	go b.start(wsClient, params, id, ctx)
	//for symbolInfo, ch := range conf.DepthNotMatchChanMap {
	//	go b.reConnect(wsClient, params, id, ctx, interval, symbolInfo, ch)
	//}
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.reConnect2(wsClient, symbols, interval, limit, wsUrl, id, ctx, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)
	return err
}

func GetSide(isMaker bool) order.TradeSide {
	switch isMaker {
	case true:
		return order.TradeSide_SELL
	case false:
		return order.TradeSide_BUY
	default:
		return order.TradeSide_InvalidSide
	}
}

func GetSideByString(side string) order.TradeSide {
	switch side {
	case "SELL":
		return order.TradeSide_SELL
	case "BUY":
		return order.TradeSide_BUY
	default:
		return order.TradeSide_InvalidSide
	}
}

func (b *BinanceSpotWebsocket) FundingRateGroup(context.Context, map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	return errors.New("spot do not have funding rate")
}

func (b *BinanceSpotWebsocket) AggTradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsAggTradeRsp) error {
	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	wsUrl := b.WsReqUrl.AGGTRADE_URL
	apiUrl := spot_api.STREAM_API_URL
	b.reqId++
	b.handler.SetAggTradeGroupChannel(chMap)
	err := b.EstablishConn(b.reqId, symbols, 0, 0, wsUrl, apiUrl, b.handler.AggTradeGroupHandle, ctx)
	return err
}

func (b *BinanceSpotWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	b.reqId++
	b.handler.SetTradeGroupChannel(chMap)
	err := b.EstablishConn(b.reqId, symbols, 0, 0, b.WsReqUrl.TRADE_URL, spot_api.STREAM_API_URL, b.handler.TradeGroupHandle, ctx)
	return err
}

func (b *BinanceSpotWebsocket) DepthIncrementGroup(ctx context.Context, interval int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {
	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	b.reqId++
	b.handler.SetDepthIncrementGroupChannel(chMap)
	err := b.EstablishConn(b.reqId, symbols, interval, -1, b.WsReqUrl.DEPTH_INCRE_URL, spot_api.STREAM_API_URL, b.handler.DepthIncrementGroupHandle, ctx)
	return err
}

func (b *BinanceSpotWebsocket) DepthLimitGroup(ctx context.Context, interval, limit int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	b.reqId++
	b.handler.SetDepthLimitGroupChannel(chMap)
	err := b.EstablishConn(b.reqId, symbols, interval, limit, b.WsReqUrl.DEPTH_LIMIT_FULL_URL, spot_api.STREAM_API_URL, b.handler.DepthLimitGroupHandle, ctx)
	return err
}

func (b *BinanceSpotWebsocket) GetFullDepth(symbol string) (*base.OrderBook, error) {
	b.GlobalDW.GlobalLock.Lock()
	defer b.GlobalDW.GlobalLock.Unlock()
	if b.GlobalDW.GlobalWeight > 1000 {
		return nil, errors.New(fmt.Sprintf("api visit %s in high frequency %d", symbol, b.ApiClient.WeightInfo[client.WeightType_REQUEST_WEIGHT].Value))
	}
	resp, err := b.ApiClient.GetDepth(spot_api.GetSymbol(symbol), 1000)
	b.GlobalDW.GlobalWeight = b.ApiClient.WeightInfo[client.WeightType_REQUEST_WEIGHT].Value
	if err != nil {
		return nil, err
	}
	symbol_, market, type_ := u_api.GetContractType(symbol)
	res := &base.OrderBook{
		Exchange: common.Exchange_BINANCE,
		Symbol:   symbol_,
		Market:   market,
		Type:     type_,
		UpdateId: resp.LastUpdateId,
	}
	ParseOrder(resp.Asks, &res.Asks)
	ParseOrder(resp.Bids, &res.Bids)
	sort.Stable(res.Asks)
	sort.Stable(sort.Reverse(res.Bids))
	return res, nil
}

func (b *BinanceSpotWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, interval, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error { //增量合并为全量副本，并校验
	var symbols []*client.SymbolInfo
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
	b.reqId++
	b.handler.SetDepthIncrementSnapshotGroupChannel(chDeltaMap, chFullMap)
	conf := &base.IncrementDepthConf{
		IsPublishDelta: isDelta,
		IsPublishFull:  isFull,
		CheckTimeSec:   10800 + rand.Intn(3600),
		DepthCapLevel:  1000,
		DepthLevel:     limit,
		GetFullDepth:   b.GetFullDepthFunc,
		Ctx:            ctx,
	}
	err := b.EstablishSnapshotConn(b.reqId, symbols, interval, -1, b.WsReqUrl.DEPTH_INCRE_URL, spot_api.STREAM_API_URL, b.handler.DepthIncrementSnapShotGroupHandle, ctx, conf)
	return err
}

func (b *BinanceSpotWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	b.reqId++
	b.handler.SetBookTickerGroupChannel(chMap)
	err := b.EstablishConn(b.reqId, symbols, 0, 0, b.WsReqUrl.BOOK_TICKER_URL, spot_api.STREAM_API_URL, b.handler.BookTickerGroupHandle, ctx)
	return err
}

func (b *BinanceSpotWebsocket) AccountExecute(req *client.WsAccountReq, action UserAction, key string) (string, error) {
	if req.Market == common.Market_SPOT {
		switch action {
		case LISTEN:
			listenKey, err := b.ApiClient.UserDataStream()
			return listenKey.ListenKey, err
		case PUT:
			if key == "" {
				return "", errors.New("listen key is none")
			}
			_, err := b.ApiClient.PutUserDataStream(key)
			return "", err
		case DELETE:
			if key == "" {
				return "", errors.New("listen key is none")
			}
			_, err := b.ApiClient.DELETEUserDataStream(key)
			return "", err
		}
	} else if req.Market == common.Market_MARGIN && req.Type == common.SymbolType_MARGIN_NORMAL {
		switch action {
		case LISTEN:
			listenKey, err := b.ApiClient.MarginUserDataStream()
			return listenKey.ListenKey, err
		case PUT:
			_, err := b.ApiClient.PutMarginUserDataStream(key)
			return "", err
		case DELETE:
			_, err := b.ApiClient.DELETEMarginUserDataStream(key)
			return "", err
		}
	} else if req.Market == common.Market_MARGIN && req.Type == common.SymbolType_MARGIN_ISOLATED {
		switch action {
		case LISTEN:
			listenKey, err := b.ApiClient.MarginIsolatedUserDataStream(req.Symbol)
			return listenKey.ListenKey, err
		case PUT:
			_, err := b.ApiClient.PutMarginIsolatedUserDataStream(req.Symbol, key)
			return "", err
		case DELETE:
			_, err := b.ApiClient.DELETEMarginIsolatedUserDataStream(req.Symbol, key)
			return "", err
		}
	} else if req.Market == common.Market_FUTURE || req.Market == common.Market_SWAP {
		//仅在u本位websocket实例可用
		switch action {
		case LISTEN:
			listenKey, err := b.UApiClient.UserDataStream()
			return listenKey.ListenKey, err
		case PUT:
			_, err := b.UApiClient.PutUserDataStream(key)
			return "", err
		case DELETE:
			_, err := b.UApiClient.DELETEUserDataStream(key)
			return "", err
		}
	} else if req.Market == common.Market_FUTURE_COIN || req.Market == common.Market_SWAP_COIN {
		//仅在币本位websocket实例可用
		switch action {
		case LISTEN:
			listenKey, err := b.CApiClient.UserDataStream()
			return listenKey.ListenKey, err
		case PUT:
			_, err := b.CApiClient.PutUserDataStream(key)
			return "", err
		case DELETE:
			_, err := b.CApiClient.DELETEUserDataStream(key)
			return "", err
		}
	}
	return "", errors.New("action not support:" + req.Market.String() + req.Type.String() + string(action))
}

func (b *BinanceSpotWebsocket) EstablishUserConn(id int, req *client.WsAccountReq, params, url string, handler func([]byte) error, ctx context.Context) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
	)
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(url).ProtoHandleFunc(handler).AutoReconnect().PostReconnectSuccess(func(wsClient *conn.WsConn) error { return b.Subscribe(id, wsClient, params) })
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = b.Subscribe(id, wsClient, params)
	if err != nil {
		return err
	}
	go b.userStart(wsClient, req, params, id, ctx)
	return err
}

func (b *BinanceSpotWebsocket) userStart(wsClient *conn.WsConn, req *client.WsAccountReq, param string, id int, ctx context.Context) {
	var (
		heartTimer = time.NewTimer(time.Second * time.Duration(b.listenTimeout))
	)
LOOP:
	for {
		select {
		case <-heartTimer.C:
			//续期
			if _, err := b.AccountExecute(req, PUT, param); err != nil {
				logger.Logger.Error("continue listen key error:" + err.Error())
			}
			heartTimer.Reset(time.Second * time.Duration(b.listenTimeout))
		case <-ctx.Done():
			if err := b.UnSubscribe(wsClient, param, id); err != nil {
				logger.Logger.Error("unsubscribe err:", param, id)
			}
			wsClient.CloseWs()
			if _, err := b.AccountExecute(req, DELETE, param); err != nil {
				logger.Logger.Error("delete listen key err:", param, id)
			}
			break LOOP
		}
	}
}

func (b *BinanceSpotWebsocket) Account(ctx context.Context, reqList ...*client.WsAccountReq) (<-chan *client.WsAccountRsp, error) {
	key, err := b.AccountExecute(reqList[0], LISTEN, "")
	if err != nil {
		logger.Logger.Error("get listen key err:", err)
		return nil, err
	}
	var (
		url = b.GetUrl(spot_api.SINGLE_API_URL, key)
	)
	b.reqId++
	if reqList[0].Market == common.Market_SPOT {
		err = b.EstablishUserConn(b.reqId, reqList[0], key, url, b.handler.AccountHandle, ctx)
	} else {
		err = b.EstablishUserConn(b.reqId, reqList[0], key, url, b.handler.MarginAccountHandle, ctx)
	}
	return b.handler.GetChan("account").(chan *client.WsAccountRsp), err
}

func (b *BinanceSpotWebsocket) Balance(ctx context.Context, reqList ...*client.WsAccountReq) (<-chan *client.WsBalanceRsp, error) {
	key, err := b.AccountExecute(reqList[0], LISTEN, "")
	if err != nil {
		logger.Logger.Error("get listen key err:", err)
		return nil, err
	}
	var (
		url = b.GetUrl(spot_api.SINGLE_API_URL, key)
	)
	b.reqId++
	if reqList[0].Market == common.Market_SPOT {
		err = b.EstablishUserConn(b.reqId, reqList[0], key, url, b.handler.BalanceHandle, ctx)
	} else {
		err = b.EstablishUserConn(b.reqId, reqList[0], key, url, b.handler.MarginBalanceHandle, ctx)
	}
	return b.handler.GetChan("balance").(chan *client.WsBalanceRsp), err
}

func (b *BinanceSpotWebsocket) Order(ctx context.Context, reqList ...*client.WsAccountReq) (<-chan *client.WsOrderRsp, error) {
	key, err := b.AccountExecute(reqList[0], LISTEN, "")
	if err != nil {
		logger.Logger.Error("get listen key err:", err)
		return nil, err
	}
	var (
		url = b.GetUrl(spot_api.SINGLE_API_URL, key)
	)
	b.reqId++
	err = b.EstablishUserConn(b.reqId, reqList[0], key, url, b.handler.OrderHandle, ctx)
	return b.handler.GetChan("order").(chan *client.WsOrderRsp), err
}
