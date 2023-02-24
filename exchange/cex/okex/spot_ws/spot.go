package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/okex/ok_api"
	"clients/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/order"
	"github.com/warmplanet/proto/go/sdk"
)

var (
	contractSizeMap sync.Map
	once            sync.Once
)

type OkWebsocket struct {
	base.WsConf
	symbolMap sdk.ConcurrentMapI //string:string --> btcusdt:BTC/USDT
	handler   WebSocketHandleInterface

	apiClient     *ok_api.ClientOkex
	WsReqUrl      *ok_api.ReqUrl
	listenTimeout int64 //second
	lock          sync.Mutex
	reqId         int
	recon         bool
}

func NewOkWebsocket(conf base.WsConf, market common.Market) *OkWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &OkWebsocket{
		WsConf: conf,
	}
	if conf.EndPoint == "" {
		d.EndPoint = WS_API_BASE_URL
	}
	// 未更改
	if d.ReadTimeout == 0 {
		d.ReadTimeout = 3
	}
	d.symbolMap = sdk.NewCmapI()

	okexHandler := NewWebSocketSpotHandle(d.ChanCap)
	okexHandler.Market = market
	d.handler = okexHandler

	d.apiClient = ok_api.NewClientOkex(conf.APIConf)
	once.Do(func() {
		InitOkContractInfo(d.apiClient)
	})
	return d
}

func InitOkContractInfo(apiClient *ok_api.ClientOkex) {
	futureSuccess, swapSuccess := false, false
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for !futureSuccess {
			res, err := apiClient.Instrument_Info("FUTURES", nil)
			if err != nil {
				logger.Logger.Error("request futures instrument info error: ", err)
				continue
			}
			for _, d := range res.Data {
				ctVal, err := strconv.ParseFloat(d.CtVal, 64)
				if err != nil {
					logger.Logger.Error("futures instrument string parse float error: ", err)
					continue
				}
				symbolInfoStr := base.SymInfoToString(transSymbolToSymbolInfo(d.InstId))
				contractSizeMap.Store(symbolInfoStr, ctVal)
			}
			futureSuccess = true
		}
	}()

	go func() {
		defer wg.Done()
		for !swapSuccess {
			res, err := apiClient.Instrument_Info("SWAP", nil)
			if err != nil {
				logger.Logger.Error("request swap instrument info error: ", err)
				continue
			}
			for _, d := range res.Data {
				ctVal, err := strconv.ParseFloat(d.CtVal, 64)
				if err != nil {
					logger.Logger.Error("futures instrument string parse float error: ", err)
					continue
				}
				symbolInfoStr := base.SymInfoToString(transSymbolToSymbolInfo(d.InstId))
				contractSizeMap.Store(symbolInfoStr, ctVal)
			}
			swapSuccess = true
		}
	}()
	wg.Wait()
}

func NewOkWebsocket2(conf base.WsConf, cli *http.Client) *OkWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &OkWebsocket{
		WsConf: conf,
	}
	if conf.EndPoint == "" {
		d.EndPoint = WS_API_BASE_URL
	}
	// 未更改
	if d.ReadTimeout == 0 {
		d.ReadTimeout = 300
	}
	if d.listenTimeout == 0 {
		d.listenTimeout = 30
	}
	d.symbolMap = sdk.NewCmapI()
	d.handler = NewWebSocketSpotHandle(d.ChanCap)

	d.apiClient = ok_api.NewClientOkex2(conf.APIConf, cli)
	return d
}

func (b *OkWebsocket) GetSymbolName(symbol string) string {
	return strings.ToLower(strings.Replace(symbol, "/", "-", -1))
}

func (b *OkWebsocket) Heartbeat() []byte {
	return []byte("ping")
}

func (b *OkWebsocket) EstablishConn(reqx req, url string, handler func([]byte) error, ctx context.Context, symbols []*client.SymbolInfo, channel string) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
	)
	if b.IsTest {
		url += "?brokerId=9999"
	}
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).
		ProtoHandleFunc(handler).
		AutoReconnect().
		Heartbeat(b.Heartbeat, 29*time.Second).
		PostReconnectSuccess(
			func(wsClient *conn.WsConn) error {
				argDateList := make([]*argDates, 0)
				for _, symbol := range symbols {
					instId := GetInstId(symbol)
					tmp := &argDates{
						Channel: DeepGear(channel),
						InstId:  instId,
					}
					argDateList = append(argDateList, tmp)
				}
				reqx.Args = argDateList
				return Subscribe(wsClient, reqx)
			})
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = Subscribe(wsClient, reqx)
	if err != nil {
		return err
	}
	go b.start(wsClient, reqx, ctx)
	checkDataSendStatus := b.handler.GetChan("send_status").(*base.CheckDataSendStatus)
	checkDataSendStatus.Init(3600, ctx, symbols...)
	go b.reConnect2(wsClient, reqx, ctx, checkDataSendStatus.UpdateTimeoutChMap, checkDataSendStatus.UpdateDateChMap, channel)

	return err
}

func (b *OkWebsocket) reConnect2(wsClient *conn.WsConn, param req, ctx context.Context, ch chan []*client.SymbolInfo,
	dateCh chan map[string]*client.SymbolInfo, channel string) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := b.UnSubscribe(wsClient, param); err != nil {
				logger.Logger.Error("unsubscribe err:", param)
			}
			wsClient.CloseWs()
			break LOOP
		case info := <-ch:
			param.Op = WsSubscrible
			Args := []*argDates{}

			for _, symbol := range info {
				instId := GetInstId(symbol)
				argDate := &argDates{
					InstId:  instId,
					Channel: DeepGear(channel),
				}
				Args = append(Args, argDate)
			}
			param.Args = Args
			param.Op = "unsubscribe"
			Subscribe(wsClient, param)
			param.Op = "subscribe"
			Subscribe(wsClient, param)
		case info := <-dateCh:
			param.Op = WsSubscrible
			var originArgs []*argDates
			var Args []*argDates

			for originDate, symbol := range info {
				instId := GetInstId(symbol)
				argDate := &argDates{
					InstId:  instId,
					Channel: DeepGear(channel),
				}
				Args = append(Args, argDate)

				originInstId := GetOriginInstId(instId, originDate)
				originArgDate := &argDates{
					InstId:  originInstId,
					Channel: DeepGear(channel),
				}
				originArgs = append(originArgs, originArgDate)
			}

			param.Args = originArgs
			param.Op = "unsubscribe"
			Subscribe(wsClient, param)
			param.Args = Args
			param.Op = "subscribe"
			Subscribe(wsClient, param)
		}
	}
}

func Subscribe(client *conn.WsConn, reqx interface{}) error {
	err := client.Subscribe(reqx)
	if err != nil {
		logger.Logger.Error("subscribe err:", err)
	}
	return err
}

func (b *OkWebsocket) start(wsClient *conn.WsConn, reqx interface{}, ctx context.Context) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := b.UnSubscribe(wsClient, reqx); err != nil {
				logger.Logger.Error("unsubscribe err:", reqx)
			}
			wsClient.CloseWs()
			break LOOP
		}
	}
}

func (b *OkWebsocket) UnSubscribe(wsClient *conn.WsConn, reqx interface{}) (err error) {
	r := reqx.(req)
	r.Op = WsUnSubscrible
	err = wsClient.Subscribe(r)
	if err != nil {
		logger.Logger.Error("unsubscribe err:", err)
	}
	return
}

func GetSide(side string) order.TradeSide {
	switch side {
	case "buy":
		return order.TradeSide_BUY
	case "sell":
		return order.TradeSide_SELL
	default:
		return order.TradeSide_InvalidSide
	}
}

func (b *OkWebsocket) DepthLimitGroup(ctx context.Context, interval, gear int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	var (
		gearw DeepGear
		err   error
	)
	symbols := []*client.SymbolInfo{}
	if gear == 1 {
		gearw = DeepGear1
	} else if gear >= 5 || gear == 0 {
		gearw = DeepGear5
	} else {
		return errors.New("档位错误")
	}
	var (
		url  = b.EndPoint + WS_API_PUBLIC
		reqx = req{
			Op:   WsSubscrible,
			Args: []*argDates{},
		}
	)
	for key, _ := range chMap {
		symbols = append(symbols, key)
		symbol := GetInstId(key)
		argDate := &argDates{
			InstId:  symbol,
			Channel: gearw,
		}
		reqx.Args = append(reqx.Args, argDate)
	}
	b.handler.SetDepthLimitGroupChannel(chMap)
	err = b.EstablishUserConn(b.handler.DepthLimitGroupHandle, ctx, reqx, url, symbols, string(gearw))
	return err
}

func (b *OkWebsocket) DepthIncrementGroup(ctx context.Context, interval int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {
	var err error
	var (
		url  = b.EndPoint + WS_API_PUBLIC
		reqx = req{
			Op:   WsSubscrible,
			Args: []*argDates{},
		}
	)
	symbols := []*client.SymbolInfo{}
	for key, _ := range chMap {
		symbols = append(symbols, key)
		symbol := GetInstId(key)
		argDate := &argDates{
			InstId:  symbol,
			Channel: DeepGear400,
		}
		reqx.Args = append(reqx.Args, argDate)
	}
	b.handler.SetDepthIncrementGroupChannel(chMap)
	err = b.EstablishUserConn(b.handler.DepthIncrementGroupHandle, ctx, reqx, url, symbols, string(DeepGear400))
	return err
}

func GetSymbol(s string) (string, error) {
	if len(s) < 1 {
		return "", errors.New("参数错误")
	}
	return strings.ReplaceAll(s, "/", "-"), nil
}

// 多组同时查询
func (b *OkWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
	var (
		err  error
		url  = b.EndPoint + WS_API_PUBLIC
		reqx = req{
			Op:   WsSubscrible,
			Args: []*argDates{},
		}
	)
	symbols := []*client.SymbolInfo{}

	args := make([]*argDates, 0)
	for symbol1 := range chMap {
		symbols = append(symbols, symbol1)
		symbol_ := GetInstId(symbol1)
		if err != nil {
			return err
		}
		tmp := argDates{
			Channel: DeepGear1,
			InstId:  symbol_,
		}
		args = append(args, &tmp)

	}
	reqx.Args = args
	b.handler.SetBookTickerGroupChannel(chMap)
	err = b.EstablishConn(reqx, url, b.handler.BookTickerGroupHandle, ctx, symbols, string(DeepGear1))
	return err
}

func (b *OkWebsocket) FundingRateGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	var symbols []*client.SymbolInfo
	var (
		err  error
		url  = b.EndPoint + WS_API_PUBLIC
		reqx = req{
			Op:   WsSubscrible,
			Args: []*argDates{},
		}
	)
	for symbol := range chMap {
		reqx.Args = append(reqx.Args, &argDates{Channel: "funding-rate", InstId: GetInstId(symbol)})
	}
	b.reqId++
	b.handler.SetFundingRateGroupChannel(chMap)
	err = b.EstablishConn(reqx, url, b.handler.FundingRateGroupHandle, ctx, symbols, "funding-rate")
	return err
}

// 交易频道
func (b *OkWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
	var err error
	var (
		url  = b.EndPoint + WS_API_PUBLIC
		reqx = req{
			Op:   WsSubscrible,
			Args: []*argDates{},
		}
	)
	symbols := []*client.SymbolInfo{}
	args := make([]*argDates, 0)
	for symbol1 := range chMap {
		symbols = append(symbols, symbol1)
		symbol_ := GetInstId(symbol1)
		if err != nil {
			return err
		}
		tmp := argDates{
			Channel: "trades",
			InstId:  symbol_,
		}
		args = append(args, &tmp)

	}
	reqx.Args = args
	b.handler.SetTradeGroupChannel(chMap)
	err = b.EstablishConn(reqx, url, b.handler.TradeGroupHandle, ctx, symbols, "trades")
	return err
}

func getUserId(id string, side, mod string) string {
	switch mod {
	case side:
		return id
	default:
		return "-1"
	}
}

// 登陆
func (b *OkWebsocket) Login(ctx context.Context, reqList ...*client.WsAccountReq) (<-chan *client.WsAccountRsp, error) {

	var (
		ch   = make(chan *client.WsAccountRsp, b.ChanCap)
		err  error
		reqx req
		url  string
	)
	url = b.EndPoint + "/private"
	err = b.EstablishUserConn(b.LoginHandle, ctx, reqx, url, []*client.SymbolInfo{}, "")
	return ch, err
}

func (b *OkWebsocket) writeMessage() *login {
	AccessKey := b.AccessKey
	b.listenTimeout = 10
	Passphrase := b.Passphrase
	tim := time.Now().Unix()
	timestamp := strconv.FormatInt(tim, 10)
	sign := ok_api.ComputeHmacSha256(timestamp+"GET"+"/users/self/verify", b.SecretKey)
	account := login{
		Op: "login",
		Args: []account{
			{
				ApiKey:     AccessKey,
				Passphrase: Passphrase,
				Sign:       sign,
				Timestamp:  timestamp,
			},
		},
	}
	return &account
}

// 登陆使用，每隔30s就会刷新sign
func (b *OkWebsocket) EstablishUserConn(handler func([]byte) error, ctx context.Context, reqx req, url string, symbols []*client.SymbolInfo, channel string) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
	)
	// 更新登陆信息（时间点）
	login := b.writeMessage()

	if b.IsTest {
		url += "?brokerId=9999"
	}
	ws := conn.NewWsBuilderWithReadTimeout(readTimeout)
	wsBuilder := ws.WsUrl(url).ProtoHandleFunc(handler).AutoReconnect().Heartbeat(b.Heartbeat, 29*time.Second).PostReconnectSuccess(
		func(wsClient *conn.WsConn) error {
			argDateList := make([]*argDates, 0)
			for _, v := range symbols {
				symbol := GetInstId(v)
				argDate := &argDates{
					InstId:  symbol,
					Channel: DeepGear(channel),
				}
				argDateList = append(argDateList, argDate)
			}
			reqx.Args = argDateList
			return Subscribe(wsClient, reqx)
		})

	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	ws.CustomSub(func() []byte {
		login = b.writeMessage()
		data, err := json.Marshal(login)
		if err != nil {
			fmt.Println(err)
		}
		return data
	})
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = Subscribe(wsClient, login)
	if err != nil {
		return err
	}
	//判断登陆是否完成 重试, 超时
	retryCount := 1
LOOP:
	for {
		select {

		case loginRes := <-b.handler.GetChan("subscribeInfo").(chan Resp_Info):
			if loginRes.Event == "login" {
				logger.Logger.Infof("ws login success")
				break LOOP
			} else {
				if retryCount < 10 {
					logger.Logger.Infof("ws login failed, begin to retry %d", retryCount)
					login = b.writeMessage()
					err = Subscribe(wsClient, login)
					if err != nil {
						return err
					}
					retryCount++
					rand.Seed(time.Now().Unix())
					time.Sleep(time.Duration(rand.Intn(5)))
					continue LOOP
				} else {
					err = errors.New(fmt.Sprintf("ws login res failed"))
					return err
				}
			}
		case <-time.After(5 * time.Second):
			if retryCount < 10 {
				logger.Logger.Infof("wait ws login res timeout, begin to retry %d", retryCount)
				login = b.writeMessage()
				err = Subscribe(wsClient, login)
				if err != nil {
					return err
				}
				retryCount++
				rand.Seed(time.Now().Unix())
				time.Sleep(time.Duration(rand.Intn(5)))
				continue LOOP
			} else {
				err = errors.New(fmt.Sprintf("wait ws login res timeout"))
				return err
			}
		}
		//if info := <-b.handler.GetChan("subscribeInfo").(chan Resp_Info); info.Event == "login" {
		//	break
		//}
	}
	err = Subscribe(wsClient, reqx)

	if err != nil {
		return err
	}
	go b.loginStart(wsClient, *login, ctx, reqx)
	if len(symbols) != 0 {
		checkDataSendStatus := b.handler.GetChan("send_status").(*base.CheckDataSendStatus)
		checkDataSendStatus.InitGetDateFunc(ok_api.GetDate)
		checkDataSendStatus.Init(3600, ctx, symbols...)
		go b.reConnect2(wsClient, reqx, ctx, checkDataSendStatus.UpdateTimeoutChMap, checkDataSendStatus.UpdateDateChMap, string(reqx.Args[0].Channel))
	}

	return err
}

func (b *OkWebsocket) loginStart(wsClient *conn.WsConn, l login, ctx context.Context, reqx req) {
	var (
		heartTimer = time.NewTimer(time.Second * time.Duration(b.listenTimeout))
		//err        error
	)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic recover, Err =", err)

		}
	}()
LOOP:
	for {
		select {
		case loginRes := <-b.handler.GetChan("subscribeInfo").(chan Resp_Info):
			// 如果登陆成功
			if loginRes.Event == "login" {
				logger.Logger.Info("okex ws login success, resub: ", reqx)
				err := Subscribe(wsClient, reqx)
				if err != nil {
					logger.Logger.Error("订阅错误:", err)
				}
			}
		case info := <-b.handler.GetChan("reSubscribe").(chan req):
			b.UnSubscribe(wsClient, info)
			info.Op = WsSubscrible
			Subscribe(wsClient, info)
		case <-heartTimer.C:
			//续期
			wsClient.SendMessage([]byte("ping"))
			heartTimer.Reset(time.Second * time.Duration(b.listenTimeout))
		case <-ctx.Done():
			b.UnSubscribe(wsClient, reqx)
			wsClient.CloseWs()
			break LOOP
		}
	}
}

func (b *OkWebsocket) updateMessage(wsClient *conn.WsConn) error {
	l := *b.writeMessage()
	err := Subscribe(wsClient, l)
	if err != nil {
		return err
	}
	return nil
}

func (b *OkWebsocket) Account(ctx context.Context, reqList ...*client.WsAccountReq) (<-chan *client.WsAccountRsp, error) {
	var (
		req_     req
		instType string
		url      string
	)

	req_ = req{
		Op: WsSubscrible,
		Args: []*argDates{
			{
				Channel: "account",
			},
		},
	}
	for _, r := range reqList {
		switch r.Market {
		case common.Market_ALL_MARKET:
			instType = "ANY"
		case common.Market_FUTURE, common.Market_FUTURE_COIN:
			instType = "FUTURES"
		case common.Market_SWAP, common.Market_SWAP_COIN:
			instType = "SWAP"
		case common.Market_MARGIN:
			instType = "MARGIN"
		default:
			instType = common.Market_name[int32(r.Market)]
		}
		req_.Args = append(req_.Args, &argDates{
			Channel:  "positions",
			InstType: instType,
		})
	}

	url = b.EndPoint + "/private"
	err := b.EstablishUserConn(b.handler.AccountHandle, ctx, req_, url, []*client.SymbolInfo{}, "positions")
	return b.handler.GetChan("account").(chan *client.WsAccountRsp), err
}

func (b *OkWebsocket) Balance(ctx context.Context, reqList ...*client.WsAccountReq) (<-chan *client.WsBalanceRsp, error) {
	var (
		req_ req
		url  string
	)
	req_ = req{
		Op: WsSubscrible,
		Args: []*argDates{
			{
				Channel: "account",
			},
		},
	}
	url = b.EndPoint + "/private"
	err := b.EstablishUserConn(b.handler.AccountHandle, ctx, req_, url, []*client.SymbolInfo{}, "account")
	return b.handler.GetChan("balance").(chan *client.WsBalanceRsp), err
}

func (b *OkWebsocket) Order(ctx context.Context, reqList ...*client.WsAccountReq) (<-chan *client.WsOrderRsp, error) {
	var (
		req_     req
		url      string
		instType string
	)
	req_ = req{
		Op:   WsSubscrible,
		Args: []*argDates{},
	}

	for _, r := range reqList {
		switch r.Market {
		case common.Market_INVALID_MARKET:
			instType = "ANY"
		case common.Market_FUTURE, common.Market_FUTURE_COIN:
			instType = "FUTURES"
		case common.Market_SWAP, common.Market_SWAP_COIN:
			instType = "SWAP"
		default:
			instType = common.Market_name[int32(r.Market)]
		}
		req_.Args = append(req_.Args, &argDates{
			Channel:  "orders",
			InstType: instType,
		})
	}

	url = b.EndPoint + "/private"
	err := b.EstablishUserConn(b.handler.OrderHandle, ctx, req_, url, []*client.SymbolInfo{}, "orders")
	return b.handler.GetChan("order").(chan *client.WsOrderRsp), err
}

func (b *OkWebsocket) LoginHandle(data []byte) error {
	// 按event区分，进行日志打印
	//fmt.Println(string(data))
	logger.Logger.Info("okex login info:", string(data))
	return nil
}

func (b *OkWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, interval, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error { //增量合并为全量副本，并校验
	var (
		symbols []*client.SymbolInfo
		url     = b.EndPoint + WS_API_PUBLIC
		reqx    = req{
			Op:   WsSubscrible,
			Args: []*argDates{},
		}
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
	for _, v := range symbols {
		symbol := GetInstId(v)
		argDate := &argDates{
			InstId:  symbol,
			Channel: DeepGear400,
		}
		reqx.Args = append(reqx.Args, argDate)
	}
	conf := &base.IncrementDepthConf{
		IsPublishDelta: isDelta,
		IsPublishFull:  isFull,
		CheckTimeSec:   3600,
		DepthCapLevel:  1000,
		DepthLevel:     limit,
		Ctx:            ctx,
	}
	b.handler.SetDepthIncrementSnapshotGroupChannel(chDeltaMap, chFullMap)
	b.handler.SetDepthIncrementSnapShotConf(symbols, conf)
	err := b.EstablishUserConn(b.handler.DepthIncrementSnapShotGroupHandle, ctx, reqx, url, symbols, string(DeepGear400))
	return err
}
