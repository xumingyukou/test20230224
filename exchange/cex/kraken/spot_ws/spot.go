package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/kraken/kraken_api"
	"clients/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/depth"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/sdk"
)

type KKWebsocket struct {
	base.WsConf
	symbolMap     sdk.ConcurrentMapI //string:string --> btcusdt:BTC/USDT
	handler       WebSocketHandleInterface
	apiClient     *kraken_api.ClientKraken
	listenTimeout int64 //second
	lock          sync.Mutex
	reqId         int
	isStart       bool
}

func NewKKWebsocket(conf base.WsConf) *KKWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &KKWebsocket{
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

	d.apiClient = kraken_api.NewClientKraken(conf.APIConf)
	return d
}

func NewKKWebsocket2(conf base.WsConf, cli *http.Client) *KKWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &KKWebsocket{
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

	d.apiClient = kraken_api.NewClientKraken2(conf.APIConf, cli)
	return d
}

func (b *KKWebsocket) EstablishConn(handler func([]byte) error, ctx context.Context, reqx req, url string, symbols []*client.SymbolInfo) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
	)

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(url).ProtoHandleFunc(handler).AutoReconnect().PostReconnectSuccess(
		func(wsClient *conn.WsConn) error {
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
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.reConnect2(wsClient, reqx, ctx, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)

	return err
}

func (b *KKWebsocket) reConnect2(wsClient *conn.WsConn, param req, ctx context.Context, ch chan []*client.SymbolInfo) {
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
			param.Pair = symbols
			param.Event = "unsubscribe"
			Subscribe(wsClient, param)
			param.Event = "subscribe"
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

func (b *KKWebsocket) UnSubscribe(wsClient *conn.WsConn, reqx interface{}) (err error) {
	r := reqx.(req)
	r.Event = "unsubscribe"
	err = wsClient.Subscribe(r)
	if err != nil {
		logger.Logger.Error("unsubscribe err:", err)
	}
	return
}

func (b *KKWebsocket) start(wsClient *conn.WsConn, reqx interface{}, ctx context.Context) {
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
			info.Event = WsSubscrible
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

// todo 统计一下r的数据有多少
// 通过pair传多个币对
func (b *KKWebsocket) DepthIncrementGroup(ctx context.Context, interval int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {
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
	var (
		url  string
		reqx = req{
			Event: WsSubscrible,
			Pair:  symbols_,
			Subscription: subscription{
				Name:  "book",
				Depth: 1000,
			},
		}
	)
	url = b.EndPoint
	b.handler.SetDepthIncrementGroupChannel(chMap)

	err = b.EstablishConn(b.handler.DepthIncrementGroupHandle, ctx, reqx, url, symbols)
	return err
}

func GetInstId(symbols ...*client.SymbolInfo) string {
	symbol := symbols[0].Symbol
	return symbol
}

func (b *KKWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
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
	var (
		url  string
		reqx = req{
			Event: WsSubscrible,
			Pair:  symbols_,
			Subscription: subscription{
				Name: "trade",
			},
		}
	)
	url = b.EndPoint
	b.handler.SetTradeGroupChannel(chMap)

	err = b.EstablishConn(b.handler.TradeGroupHandle, ctx, reqx, url, symbols)
	return err
}

func (b *KKWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
	var err error
	symbols_ := []string{}
	symbols := []*client.SymbolInfo{}
	for k, _ := range chMap {
		symbols = append(symbols, k)
		symbols_ = append(symbols_, GetInstId(k))
	}
	if err != nil {
		panic(err)
	}
	var (
		url  string
		reqx = req{
			Event: WsSubscrible,
			Pair:  symbols_,
			Subscription: subscription{
				Name: "ticker",
			},
		}
	)
	url = b.EndPoint
	b.handler.SetBookTickerGroupChannel(chMap)

	err = b.EstablishConn(b.handler.BookTickerGroupHandle, ctx, reqx, url, symbols)
	return err
}

func (b *KKWebsocket) DepthLimitGroup(ctx context.Context, interval, gear int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	return errors.New("kraken do not have limit websocket, the snap will store the full amount of data")
}

func (b *KKWebsocket) FundingRateGroup(context.Context, map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	return errors.New("spot do not have funding rate")
}

func (b *KKWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, interval, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error { //增量合并为全量副本，并校验
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
			symbols_ = append(symbols_, GetInstId(k))
			symbols = append(symbols, k)
		}
	} else {
		for k, _ := range chFullMap {
			symbols_ = append(symbols_, GetInstId(k))
			symbols = append(symbols, k)
		}
	}
	if err != nil {
		panic(err)
	}
	var (
		url  string
		reqx = req{
			Event: WsSubscrible,
			Pair:  symbols_,
			Subscription: subscription{
				Name:  "book",
				Depth: 1000,
			},
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
	// fmt.Println(url)
	err = b.EstablishConn(b.handler.DepthIncrementSnapShotGroupHandle, ctx, reqx, url, symbols)
	return err

}
