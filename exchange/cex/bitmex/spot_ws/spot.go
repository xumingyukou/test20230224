package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	bitMex "clients/exchange/cex/bitmex"
	"clients/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/warmplanet/proto/go/depth"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/sdk"
)

type BitMexWebsocket struct {
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

	apiClient     *bitMex.ClientBitMex
	listenTimeout int64 //second
	lock          sync.Mutex
	reqId         int
	isStart       bool
	nameMap       map[string]string
}

func NewBitMexWebsocket(conf base.WsConf) *BitMexWebsocket {

	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &BitMexWebsocket{
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

	return d
}

func NewBitMexWebsocket2(conf base.WsConf, cli *http.Client) *BitMexWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &BitMexWebsocket{
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

	d.apiClient = bitMex.NewClientBitMex(conf.APIConf, cli)
	return d
}

func (b *BitMexWebsocket) EstablishConn(handler func([]byte) error, ctx context.Context, reqx req, url string, action string, symbols []*client.SymbolInfo) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
	)
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
	err = Subscribe(wsClient, reqx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	go b.start(wsClient, reqx, ctx)
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.reConnect2(wsClient, reqx, ctx, 0, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap, action)
	return err
}

func (b *BitMexWebsocket) reConnect2(wsClient *conn.WsConn, param req, ctx context.Context, interval int, ch chan []*client.SymbolInfo, action string) {
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
				symbols = append(symbols, GetInstId(v))
			}
			param.Op = WsUnSubscribe
			args := []string{}
			for _, v := range symbols {
				args = append(args, action+v)
			}
			param.Args = args
			if err := b.UnSubscribe(wsClient, param); err != nil {
				logger.Logger.Error("unsubscribe err:", param)
			}
			param.Op = WsSubscribe
			if err := Subscribe(wsClient, param); err != nil {
				logger.Logger.Error("subscribe err:", param)
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

func (b *BitMexWebsocket) UnSubscribe(wsClient *conn.WsConn, reqx interface{}) (err error) {
	r := reqx.(req)
	err = wsClient.Subscribe(r)
	if err != nil {
		logger.Logger.Error("unsubscribe err:", err)
	}
	return
}

func (b *BitMexWebsocket) start(wsClient *conn.WsConn, reqx interface{}, ctx context.Context) {
LOOP:
	for {
		select {
		case info := <-b.handler.GetChan("reSubscribe").(chan req):
			data, err := json.Marshal(info)
			wsClient.SendMessage(data)
			logger.Logger.Warn("开始重连")
			if err != nil {
				fmt.Println(err)
			}
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
func (b *BitMexWebsocket) DepthIncrementGroup(ctx context.Context, interval int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {
	var (
		err     error
		symbols []*client.SymbolInfo
	)
	symbols_ := []string{}
	for k, _ := range chMap {
		symbols = append(symbols, k)
		symbols_ = append(symbols_, GetInstId(k))
	}
	args := []string{}
	for _, v := range symbols_ {
		v = strings.ReplaceAll(v, "BTC", "XBT")
		args = append(args, "orderBookL2_25:"+strings.ReplaceAll(v, "/", ""))
	}
	var (
		url  string
		reqx = req{
			Op:   WsSubscribe,
			Args: args,
		}
	)
	conf := &base.IncrementDepthConf{
		IsPublishDelta: true,
		IsPublishFull:  false,
		CheckTimeSec:   3600,
		DepthCapLevel:  1000,
		DepthLevel:     1000,
		Ctx:            ctx,
	}
	url = b.EndPoint
	b.handler.SetDepthIncrementGroupChannel(chMap)
	b.handler.SetDepthIncrementSnapShotConf(symbols, conf)

	err = b.EstablishConn(b.handler.DepthIncrementGroupHandle, ctx, reqx, url, "orderBookL2_25:", symbols)
	return err
}

// 将当前币对转化为服务器需要的格式
func GetInstId(symbols ...*client.SymbolInfo) string {

	symbol := symbols[0].Symbol
	// spot转换
	symbol = strings.ReplaceAll(symbol, "/", "")

	return symbol
}

// 把格式转换回客户端的格式
func TranInstId(symbols ...*client.SymbolInfo) string {
	return strings.ReplaceAll(GetInstId(symbols...), "_", "/")
}

func (b *BitMexWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
	var err error
	symbols := []*client.SymbolInfo{}
	symbols_ := []string{}
	for k, _ := range chMap {
		symbols = append(symbols, k)
		symbols_ = append(symbols_, GetInstId(k))
	}
	if err != nil {
		panic(err)
	}
	args := []string{}
	for _, v := range symbols_ {
		v = strings.ReplaceAll(v, "BTC", "XBT")
		args = append(args, "trade:"+strings.ReplaceAll(v, "/", ""))
	}
	var (
		url  string
		reqx = req{
			Op:   WsSubscribe,
			Args: args,
		}
	)
	url = b.EndPoint
	b.handler.SetTradeGroupChannel(chMap)

	err = b.EstablishConn(b.handler.TradeGroupHandle, ctx, reqx, url, "trade:", symbols)
	return err
}

func (b *BitMexWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
	var err error
	symbols_ := []string{}
	symbols := []*client.SymbolInfo{}
	for k, _ := range chMap {
		symbols = append(symbols, k)
		symbols_ = append(symbols_, GetInstId(k))
	}
	args := []string{}
	for _, v := range symbols_ {
		v = strings.ReplaceAll(v, "BTC", "XBT")
		args = append(args, "instrument:"+strings.ReplaceAll(v, "/", ""))
	}
	var (
		url  string
		reqx = req{
			Op:   WsSubscribe,
			Args: args,
		}
	)
	url = b.EndPoint
	b.handler.SetBookTickerGroupChannel(chMap)

	err = b.EstablishConn(b.handler.BookTickerGroupHandle, ctx, reqx, url, "instrument:", symbols)
	return err
}

func (b *BitMexWebsocket) DepthLimitGroup(ctx context.Context, interval, gear int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	return errors.New("kraken do not have limit websocket, the snap will store the full amount of data")
}

func (b *BitMexWebsocket) FundingRateGroup(context.Context, map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	return errors.New("spot do not have funding rate")
}

func (b *BitMexWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, interval, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error { //增量合并为全量副本，并校验
	var (
		err     error
		symbols []*client.SymbolInfo
	)
	symbols_ := []string{}

	if isDelta && isFull && len(chDeltaMap) != len(chFullMap) {
		return errors.New("chDeltaMap and chFullMap not match")
	}
	if isDelta {
		for k, _ := range chDeltaMap {
			symbols = append(symbols, k)
			symbols_ = append(symbols_, k.Symbol)
		}
	} else {
		for k, _ := range chFullMap {
			symbols = append(symbols, k)
			symbols_ = append(symbols_, k.Symbol)
		}
	}
	args := []string{}
	for _, v := range symbols_ {
		v = strings.ReplaceAll(v, "BTC", "XBT")
		args = append(args, "orderBookL2:"+strings.ReplaceAll(v, "/", ""))
	}
	var (
		url  string
		reqx = req{
			Op:   WsSubscribe,
			Args: args,
		}
	)
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
	err = b.EstablishConn(b.handler.DepthIncrementSnapShotGroupHandle, ctx, reqx, url, "orderBookL2:", symbols)
	return err

}
