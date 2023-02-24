package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/bitflyer/bitFlyer_api"
	"clients/logger"
	"context"
	"encoding/json"
	"errors"
	"github.com/warmplanet/proto/go/depth"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/sdk"
)

type BitFlyerWebsocket struct {
	base.WsConf
	symbolMap sdk.ConcurrentMapI //string:string --> btcusdt:BTC/USDT
	handler   WebSocketHandleInterface

	//tradeChanMap          sdk.ConcurrentMapI //string:chan *client.WsTradeRsp
	//depthIncrementChanMap sdk.ConcurrentMapI //string:chan *client.WsDepthRsp
	//depthLimitChanMap     sdk.ConcurrentMapI //string:chan *client.WsDepthRsp
	//bookTickerChanMap     sdk.ConcurrentMapI //string:chan *client.WsDepthRsp
	//accountChanMap        sdk.ConcurrentMapI //string:chan *client.WsAccountRsp
	//balanceChanMap        sdk.ConcurrentMapI //string:chan *client.WsAccountRsp
	//orderChanMap          sdk.ConcurrentMapI //string:chan *client.WsOrderRsp
	//websocketMap  sdk.ConcurrentMapI //string:conn.WsConn

	apiClient     *bitFlyer_api.ClientBitFlyer
	listenTimeout int64 //second
	lock          sync.Mutex
	reqId         int
	isStart       bool
}

func NewBitFlyerWebsocket(conf base.WsConf) *BitFlyerWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &BitFlyerWebsocket{
		WsConf: conf,
	}
	if conf.EndPoint == "" {
		d.EndPoint = WS_PUBLIC_BASE_URL
	}
	// 未更改
	if d.ReadTimeout == 0 {
		d.ReadTimeout = 3
	}

	d.symbolMap = sdk.NewCmapI()
	d.handler = NewWebSocketSpotHandle(d.ChanCap)

	d.apiClient = bitFlyer_api.NewClientBitFlyer(conf.APIConf)
	return d
}

func NewBitFlyerWebsocket2(conf base.WsConf, cli *http.Client) *BitFlyerWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &BitFlyerWebsocket{
		WsConf: conf,
	}
	if conf.EndPoint == "" {
		d.EndPoint = WS_PUBLIC_BASE_URL
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

	d.apiClient = bitFlyer_api.NewClientBitFlyer2(conf.APIConf, cli)
	return d
}

func (b *BitFlyerWebsocket) EstablishConn(handler func([]byte) error, ctx context.Context, reqx req, url string, params []string, symbols []*client.SymbolInfo, prefix string) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
	)
	reqx.Params = map[string]string{}

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(url).ProtoHandleFunc(handler).AutoReconnect().PostReconnectSuccess(func(wsClient *conn.WsConn) error {
		return Subscribe(wsClient, reqx)
	})
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	for _, v := range params {
		reqx.Params["channel"] = v
		//fmt.Println(v)
		err = Subscribe(wsClient, reqx)
	}
	if err != nil {
		//fmt.Println(err)
		return err
	}
	go b.start(wsClient, reqx, ctx)
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.reConnect2(wsClient, reqx, ctx, prefix, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)
	return err
}

func (b *BitFlyerWebsocket) reConnect2(wsClient *conn.WsConn, param req, ctx context.Context, prefix string, ch chan []*client.SymbolInfo) {
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
			symbols := []string{}
			for _, v := range info {
				symbols = append(symbols, prefix+GetInstId(v))
			}
			param.Method = WsUnSubscribe
			for _, v := range symbols {
				param.Params["channel"] = v
				//fmt.Println(v)
				err := Subscribe(wsClient, param)
				if err != nil {
					logger.Logger.Error("unsubscribe err:", param)
				}
			}
			param.Method = WsSubscribe
			for _, v := range symbols {
				param.Params["channel"] = v
				//fmt.Println(v)
				err := Subscribe(wsClient, param)
				if err != nil {
					logger.Logger.Error("subscribe err:", param)
				}
			}
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

func TSubscribe(client *conn.WsConn, reqx interface{}) error {
	err := client.Subscribe(reqx)
	if err != nil {
		logger.Logger.Error("subscribe err:", err)
	}
	return err
}

func (b *BitFlyerWebsocket) UnSubscribe(wsClient *conn.WsConn, reqx interface{}) (err error) {
	r := reqx.(req)
	r.Method = "unsubscribe"
	err = wsClient.Subscribe(r)
	if err != nil {
		logger.Logger.Error("unsubscribe err:", err)
	}
	return
}

func (b *BitFlyerWebsocket) start(wsClient *conn.WsConn, reqx interface{}, ctx context.Context) {
LOOP:
	for {
		select {
		case info := <-b.handler.GetChan("reSubscribe").(chan req):
			data, err := json.Marshal(info)
			wsClient.SendMessage(data)
			logger.Logger.Warn("开始重连")
			if err != nil {
				logger.Logger.Error(err)
			}
			info.Method = WsSubscribe
			Subscribe(wsClient, info)
		case <-ctx.Done():
			if err := b.UnSubscribe(wsClient, reqx); err != nil {
				logger.Logger.Error("unsubscribe err:", reqx)
			}
			wsClient.CloseWs()
			break LOOP
		}
	}
}

// 通过pair传多个币对
func (b *BitFlyerWebsocket) DepthIncrementGroup(ctx context.Context, interval int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {
	var err error
	symbols := []*client.SymbolInfo{}
	symbols_ := []string{}
	for k, _ := range chMap {
		symbols = append(symbols, k)
		symbols_ = append(symbols_, GetInstId(k))
	}
	params := []string{}
	for _, v := range symbols_ {
		params = append(params, "lightning_board_"+v)
	}
	if err != nil {
		panic(err)
	}
	var (
		url  string
		reqx = req{
			Method: WsSubscribe,
			Params: map[string]string{
				//"channel": "lightning_board_snapshot_BTC_USD",
			},
		}
	)
	url = b.EndPoint
	b.handler.SetDepthIncrementGroupChannel(chMap)

	err = b.EstablishConn(b.handler.DepthIncrementGroupHandle, ctx, reqx, url, params, symbols, "lightning_board_")
	return err
}

// 将当前币对转化为服务器需要的格式
func GetInstId(symbols ...*client.SymbolInfo) string {

	symbol := symbols[0].Symbol
	// spot转换
	symbol = strings.ReplaceAll(symbol, "/", "_")

	return symbol
}

// 把格式转换回客户端的格式
func TranInstId(symbols ...*client.SymbolInfo) string {
	return strings.ReplaceAll(GetInstId(symbols...), "_", "/")
}

func (b *BitFlyerWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
	var err error
	symbols := []*client.SymbolInfo{}
	symbols_ := []string{}
	for k, _ := range chMap {
		symbols = append(symbols, k)
		symbols_ = append(symbols_, GetInstId(k))
	}
	params := []string{}
	for _, v := range symbols_ {
		params = append(params, "lightning_executions_"+v)
	}
	if err != nil {
		panic(err)
	}
	var (
		url  string
		reqx = req{
			Method: WsSubscribe,
			Params: map[string]string{
				//"channel": "lightning_board_snapshot_BTC_USD",
			},
		}
	)
	url = b.EndPoint
	b.handler.SetTradeGroupChannel(chMap)

	err = b.EstablishConn(b.handler.TradeGroupHandle, ctx, reqx, url, params, symbols, "lightning_executions_")
	return err
}

func (b *BitFlyerWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
	var err error
	prefix := "lightning_ticker_"
	symbols := []*client.SymbolInfo{}
	symbols_ := []string{}
	for k, _ := range chMap {
		symbols_ = append(symbols_, GetInstId(k))
	}
	params := []string{}
	for _, v := range symbols_ {
		params = append(params, "lightning_ticker_"+v)
	}
	if err != nil {
		panic(err)
	}
	var (
		url  string
		reqx = req{
			Method: WsSubscribe,
			Params: map[string]string{
				//"channel": "lightning_board_snapshot_BTC_USD",
			},
		}
	)
	url = b.EndPoint
	b.handler.SetBookTickerGroupChannel(chMap)

	err = b.EstablishConn(b.handler.BookTickerGroupHandle, ctx, reqx, url, params, symbols, prefix)
	return err
}

func (b *BitFlyerWebsocket) DepthLimitGroup(ctx context.Context, interval, gear int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	return errors.New("kraken do not have limit websocket, the snap will store the full amount of data")
}

func (b *BitFlyerWebsocket) FundingRateGroup(context.Context, map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	return errors.New("spot do not have funding rate")
}

func (b *BitFlyerWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, interval, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error { //增量合并为全量副本，并校验
	var (
		err     error
		symbols []*client.SymbolInfo
		params  = []string{}
	)
	symbols_ := []string{}
	prefix := "lightning_board_snapshot_"
	if isDelta && isFull && len(chDeltaMap) != len(chFullMap) {
		return errors.New("chDeltaMap and chFullMap not match")
	}
	if isDelta {
		for k, _ := range chDeltaMap {
			symbols = append(symbols, k)
			symbols_ = append(symbols_, GetInstId(k))
		}
	} else {
		for k, _ := range chFullMap {
			symbols = append(symbols, k)
			symbols_ = append(symbols_, GetInstId(k))
		}
	}
	var (
		url  string
		reqx = req{
			Method: WsSubscribe,
			Params: map[string]string{
				//"channel": "lightning_board_snapshot_BTC_USD",
			},
		}
	)
	for _, v := range symbols_ {
		params = append(params, "lightning_board_snapshot_"+v)
	}
	conf := &base.IncrementDepthConf{
		IsPublishDelta: isDelta,
		IsPublishFull:  isFull,
		CheckTimeSec:   3600,
		DepthCapLevel:  1000,
		DepthLevel:     limit,
		Ctx:            ctx,
	}
	url = b.EndPoint
	b.handler.SetDepthIncrementSnapshotGroupChannel(chDeltaMap, chFullMap)
	b.handler.SetDepthIncrementSnapShotConf(symbols, conf)
	err = b.EstablishConn(b.handler.DepthIncrementSnapShotGroupHandle, ctx, reqx, url, params, symbols, prefix)
	return err

}
