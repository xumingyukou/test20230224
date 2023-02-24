package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/gate/gate_api"
	"clients/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/depth"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/sdk"
)

type GaWebsocket struct {
	base.WsConf
	symbolMap sdk.ConcurrentMapI //string:string --> btcusdt:BTC/USDT
	handler   WebSocketHandleInterface

	apiClient     *gate_api.ClientGate
	listenTimeout int64 //second
	lock          sync.Mutex
	reqId         int
	isStart       bool
	ping          PingInfo
}

func NewGaWebsocket(conf base.WsConf) *GaWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &GaWebsocket{
		WsConf: conf,
	}
	if conf.EndPoint == "" {
		d.EndPoint = WS_PUBLIC_BASE_URL
	}
	// 未更改
	if d.ReadTimeout == 0 {
		d.ReadTimeout = 3
	}

	d.ping = PingInfo{Channel: "spot.ping"}
	d.symbolMap = sdk.NewCmapI()
	d.handler = NewWebSocketSpotHandle(d.ChanCap)

	d.apiClient = gate_api.NewClientGate(conf.APIConf)
	d.apiClient.EndPoint = gate_api.GLOBAL_API_BASE_URL
	return d
}

func NewGaWebsocket2(conf base.WsConf, cli *http.Client) *GaWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &GaWebsocket{
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

	d.ping = PingInfo{Channel: "spot.ping"}
	d.symbolMap = sdk.NewCmapI()
	d.handler = NewWebSocketSpotHandle(d.ChanCap)

	d.apiClient = gate_api.NewClientGate2(conf.APIConf, cli)
	d.apiClient.EndPoint = gate_api.GLOBAL_API_BASE_URL

	return d
}

func (b *GaWebsocket) EstablishConn(handler func([]byte) error, ctx context.Context, reqx req, url_ string, symbols []*client.SymbolInfo, interval ...string) error {
	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
		err         error
		interval_   string
	)

	heartBeat := func() []byte {
		t := time.Now().UnixMicro()
		b.ping.Time = t
		date, err := json.Marshal(b.ping)
		if err != nil {
			logger.Logger.Error("gate ping 解析错误:", err)
		}
		return date
	}
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(url_).ProtoHandleFunc(handler).AutoReconnect().Heartbeat(heartBeat, 10*time.Second).PostReconnectSuccess(func(wsClient *conn.WsConn) error {
		return Subscribe(wsClient, reqx)
	})
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	if len(interval) == 0 {
		reqx.Time = time.Now().UnixMicro()
		err = Subscribe(wsClient, reqx)
	} else {
		for _, symbol := range reqx.Payload {
			s := []string{}
			tmp := reqx
			tmp.Time = time.Now().UnixMicro()
			tmp.Payload = append(s, symbol, interval[0])
			err = TSubscribe(wsClient, tmp)
		}
	}
	if interval == nil {
		interval_ = ""
	} else {
		interval_ = interval[0]
	}
	if err != nil {
		return err
	}
	go b.start(wsClient, reqx, ctx)
	b.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go b.reConnect2(wsClient, reqx, ctx, interval_, b.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)

	return err
}

func (b *GaWebsocket) reConnect2(wsClient *conn.WsConn, param req, ctx context.Context, interval string, ch chan []*client.SymbolInfo) {
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
			param.Payload = symbols
			SubComprehen(WsUnSubscribe, interval, param, wsClient)
			SubComprehen(WsSubscribe, interval, param, wsClient)
		}
	}
}

func SubComprehen(event, interval string, reqx req, wsClient *conn.WsConn) {
	var err error
	reqx.Event = event
	if interval == "" {
		reqx.Time = time.Now().UnixMicro()
		err = Subscribe(wsClient, reqx)
		if err != nil {
			logger.Logger.Error("重订阅错误", err)
		}
	} else {
		for _, symbol := range reqx.Payload {
			s := []string{}
			tmp := reqx
			tmp.Time = time.Now().UnixMicro()
			tmp.Payload = append(s, symbol, interval)
			err = TSubscribe(wsClient, tmp)
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

func (b *GaWebsocket) UnSubscribe(wsClient *conn.WsConn, reqx interface{}) (err error) {
	r := reqx.(req)
	r.Event = "unsubscribe"
	err = wsClient.Subscribe(r)
	if err != nil {
		logger.Logger.Error("unsubscribe err:", err)
	}
	return
}

func (b *GaWebsocket) start(wsClient *conn.WsConn, reqx interface{}, ctx context.Context) {
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
			info.Event = WsSubscribe
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
func (b *GaWebsocket) DepthIncrementGroup(ctx context.Context, interval int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {
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
			Event:   WsSubscribe,
			Payload: symbols_,
			Channel: "spot.order_book_update",
		}
	)
	url = b.EndPoint
	b.handler.SetDepthIncrementGroupChannel(chMap)

	err = b.EstablishConn(b.handler.DepthIncrementGroupHandle, ctx, reqx, url, symbols, "100ms")
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

func (b *GaWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
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
			Event:   WsSubscribe,
			Payload: symbols_,
			Channel: "spot.trades",
		}
	)
	url = b.EndPoint
	b.handler.SetTradeGroupChannel(chMap)

	err = b.EstablishConn(b.handler.TradeGroupHandle, ctx, reqx, url, symbols)
	return err
}

func (b *GaWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
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
			Event:   WsSubscribe,
			Payload: symbols_,
			Channel: "spot.tickers",
		}
	)
	url = b.EndPoint
	b.handler.SetBookTickerGroupChannel(chMap)

	err = b.EstablishConn(b.handler.BookTickerGroupHandle, ctx, reqx, url, symbols)
	return err
}

func (b *GaWebsocket) DepthLimitGroup(ctx context.Context, interval, gear int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	return errors.New("kraken do not have limit websocket, the snap will store the full amount of data")
}

func (b *GaWebsocket) FundingRateGroup(context.Context, map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	return errors.New("spot do not have funding rate")
}

func (b *GaWebsocket) GetFullDepth(exSymbol string) (*base.OrderBook, error) {
	params := url.Values{}
	params.Add("limit", "100")
	//proxyUrl, err := url.Parse("http://127.0.0.1:7890")
	//transport := http.Transport{
	//	Proxy: http.ProxyURL(proxyUrl),
	//}
	//b.apiClient.HttpClient = &http.Client{
	//	Transport: &transport,
	//}
	fullRes, err := b.apiClient.Market_Books_Info(exSymbol, &params)

	if err != nil {
		logger.Logger.Error("获取全量问题 ", err.Error())
		return nil, err
	}
	tmpt := time.Now().UnixMicro()

	// 全量转换为depthUpdate进行排序
	msg_ := &Resp_Depth{}
	full := &base.DeltaDepthUpdate{}
	full2Depth(fullRes, msg_)
	// 清零，排序
	depth2Delta(msg_, full)
	full.Symbol = ParseSymbol(exSymbol)
	// 全量转换为book
	book := &base.OrderBook{}
	delta2Book(full, book)
	book.TimeReceive = uint64(tmpt)
	return book, nil
}

func (b *GaWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, interval, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error { //增量合并为全量副本，并校验
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
			symbols_ = append(symbols_, GetInstId(k))
		}
	} else {
		for k, _ := range chFullMap {
			symbols = append(symbols, k)
			symbols_ = append(symbols_, GetInstId(k))
		}
	}
	if err != nil {
		panic(err)
	}
	var (
		url  string
		reqx = req{
			Event:   WsSubscribe,
			Payload: symbols_,
			Channel: "spot.order_book_update",
		}
	)
	conf := &base.IncrementDepthConf{
		IsPublishDelta: isDelta,
		IsPublishFull:  isFull,
		CheckTimeSec:   3600,
		DepthCapLevel:  1000,
		DepthLevel:     limit,
		Ctx:            ctx,
		GetFullDepth:   b.GetFullDepth,
	}
	url = b.EndPoint
	b.handler.SetDepthIncrementSnapshotGroupChannel(chDeltaMap, chFullMap)
	b.handler.SetDepthIncrementSnapShotConf(symbols, conf)
	err = b.EstablishConn(b.handler.DepthIncrementSnapShotGroupHandle, ctx, reqx, url, symbols, "100ms")
	return err

}
