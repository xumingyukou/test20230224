package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/coinbase/spot_api"
	"clients/logger"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CoinbaseSpotWebsocket struct {
	base.WsConf
	apiClient      *spot_api.ApiClient
	readTimeout    int64 //can delete
	listenTimeout  int64 //second
	lock           sync.Mutex
	sequenceNumber int
	isStart        bool
	WsReqUrl       *spot_api.WsReqUrl
	handler        WebSocketHandleInterface
}

func NewCoinbaseWebsocket(conf base.WsConf) *CoinbaseSpotWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &CoinbaseSpotWebsocket{
		WsConf: conf,
	}
	if conf.EndPoint == "" {
		d.EndPoint = spot_api.WS_API_BASE_URL
	}
	if d.ReadTimeout == 0 {
		d.ReadTimeout = 300
	}
	if d.listenTimeout == 0 {
		d.listenTimeout = 1800
	}
	d.handler = NewWebSocketSpotHandle(d.ChanCap)

	if conf.AccessKey != "" && conf.SecretKey != "" {
		d.apiClient = spot_api.NewApiClient(base.APIConf{
			ProxyUrl:    conf.ProxyUrl,
			ReadTimeout: conf.ReadTimeout,
			AccessKey:   conf.AccessKey,
			SecretKey:   conf.SecretKey,
		})
	}
	d.WsReqUrl = spot_api.NewSpotWsUrl()
	return d
}

/* Connection Methods*/
func (c *CoinbaseSpotWebsocket) Subscribe(wsClient *conn.WsConn, request Request) error {
	request.Type = Subscribe
	request.Extensions = "permessage-deflate" //We send this in opening handshake to allow for compression
	err := wsClient.Subscribe(request)
	if err != nil {
		logger.Logger.Error("subscribe err:", err)
	}
	return err
}
func (c *CoinbaseSpotWebsocket) Unsubscribe(wsClient *conn.WsConn, request Request) (err error) {
	request.Type = Unsubscribe
	err = wsClient.Subscribe(request)
	if err != nil {
		logger.Logger.Error("unsubscribe err:", err)
	}
	return
}
func (c *CoinbaseSpotWebsocket) Start(wsClient *conn.WsConn, ctx context.Context, request Request) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := c.Unsubscribe(wsClient, request); err != nil {
				logger.Logger.Error("unsubscribe err:", request, err)
			}
			wsClient.CloseWs()
			break LOOP
		}
	}
}
func (c *CoinbaseSpotWebsocket) EstablishConn(symbols []*client.SymbolInfo, request Request, url string, handler func([]byte) error, ctx context.Context) error {
	var (
		readTimeout = time.Duration(c.ReadTimeout) * time.Second * 1000
		err         error
	)
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).
		ProtoHandleFunc(handler).
		AutoReconnect().
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			return c.Subscribe(wsClient, request)
		})
	if c.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(c.ProxyUrl)
	}
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = c.Subscribe(wsClient, request)
	if err != nil {
		return err
	}
	go c.Start(wsClient, ctx, request)
	c.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go c.Reconnect(wsClient, ctx, request, c.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)
	return err
}
func (c *CoinbaseSpotWebsocket) EstablishSnapshotConn(symbols []*client.SymbolInfo, request Request, url string, handler func([]byte) error, ctx context.Context, conf *base.IncrementDepthConf, interval int) error {
	var (
		readTimeout = time.Duration(c.ReadTimeout) * time.Second * 1000
		//newRequest  Request
		err error
	)
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).
		ProtoHandleFunc(handler).
		AutoReconnect().
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			return c.Subscribe(wsClient, request)
		})
	if c.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(c.ProxyUrl)
	}
	c.handler.SetDepthIncrementSnapShotConf(symbols, conf)
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	err = c.Subscribe(wsClient, request)
	if err != nil {
		return err
	}
	go c.Start(wsClient, ctx, request)

	//newRequest = request
	//for symbolInfo, ch := range conf.DepthNotMatchChanMap {
	//	newRequest.ProductIds = []string{ReformatSymbols(symbolInfo.Symbol)}
	//	go c.Reconnect(wsClient, ctx, newRequest, symbolInfo, ch)
	//}
	c.handler.GetChan("send_status").(*base.CheckDataSendStatus).Init(3600, ctx, symbols...)
	go c.Reconnect(wsClient, ctx, request, c.handler.GetChan("send_status").(*base.CheckDataSendStatus).UpdateTimeoutChMap)
	return err
}
func (c *CoinbaseSpotWebsocket) Reconnect(wsClient *conn.WsConn, ctx context.Context, request Request, ch chan []*client.SymbolInfo) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := c.Unsubscribe(wsClient, request); err != nil {
				logger.Logger.Error("unsubscribe err:", request)
			}
			wsClient.CloseWs()
			break LOOP
		case reconnectData := <-ch:
			var ReconnectSymbols []string
			reconnectRequest := request
			for _, symbol := range reconnectData {
				ReconnectSymbols = append(ReconnectSymbols, ReformatSymbols(symbol.Symbol))
			}
			reconnectRequest.ProductIds = ReconnectSymbols
			logger.Logger.Info("Update Timeout Resubscribe:", reconnectRequest)
			if err := c.Unsubscribe(wsClient, reconnectRequest); err != nil {
				logger.Logger.Error("unsubscribe err:", reconnectRequest)
			}
			if err := c.Subscribe(wsClient, reconnectRequest); err != nil {
				logger.Logger.Error("unsubscribe err:", reconnectRequest)
			}
		}
	}
}

/*Main Interface Methods*/
func (c *CoinbaseSpotWebsocket) FundingRateGroup(context.Context, map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	return errors.New("spot do not have funding rate")
}
func (c *CoinbaseSpotWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
	var (
		symbolInfo []*client.SymbolInfo
		symbols    []string
		request    Request
		err        error
	)
	for symbol := range chMap {
		symbolInfo = append(symbolInfo, symbol)
		symbols = append(symbols, ReformatSymbols(symbol.Symbol))
	}
	c.handler.SetTradeGroupChannel(chMap)
	request = Request{
		ProductIds: symbols,
		Channels:   []string{"matches", "heartbeat"},
	}
	err = c.EstablishConn(symbolInfo, request, spot_api.WS_API_BASE_URL, c.handler.TradeGroupHandle, ctx)
	return err
}
func (c *CoinbaseSpotWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
	var (
		symbolInfo []*client.SymbolInfo
		symbols    []string
		request    Request
		err        error
	)
	for symbol := range chMap {
		symbolInfo = append(symbolInfo, symbol)
		symbols = append(symbols, ReformatSymbols(symbol.Symbol))
	}

	c.handler.SetBookTickerGroupChannel(chMap)

	request = Request{
		ProductIds: symbols,
		Channels:   []string{"ticker", "heartbeat"},
	}
	err = c.EstablishConn(symbolInfo, request, spot_api.WS_API_BASE_URL, c.handler.BookTickerGroupHandle, ctx)
	return err
}
func (c *CoinbaseSpotWebsocket) DepthLimitGroup(ctx context.Context, interval, limit int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	var (
		symbolInfo []*client.SymbolInfo
		symbols    []string
		request    Request
		err        error
	)
	for symbol := range chMap {
		symbolInfo = append(symbolInfo, symbol)
		symbols = append(symbols, ReformatSymbols(symbol.Symbol))
	}

	c.handler.SetDepthLimitGroupChannel(chMap)
	request = Request{
		ProductIds: symbols,
		Channels:   []string{"level2_batch", "heartbeat"},
	}
	err = c.EstablishConn(symbolInfo, request, spot_api.WS_API_BASE_URL, c.handler.DepthLimitGroupHandle, ctx)
	return err
}
func (c *CoinbaseSpotWebsocket) DepthIncrementGroup(ctx context.Context, internal int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {
	var (
		symbolInfo []*client.SymbolInfo
		symbols    []string
		request    Request
		err        error
	)
	for symbol := range chMap {
		symbolInfo = append(symbolInfo, symbol)
		symbols = append(symbols, ReformatSymbols(symbol.Symbol))
	}

	c.handler.SetDepthIncrementGroupChannel(chMap)
	request = Request{
		ProductIds: symbols,
		Channels:   []string{"level2_batch", "heartbeat"},
	}
	err = c.EstablishConn(symbolInfo, request, spot_api.WS_API_BASE_URL, c.handler.DepthIncrementGroupHandle, ctx)
	return err
}
func (c *CoinbaseSpotWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, interval, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error { //增量合并为全量副本，并校验
	var (
		symbols      []*client.SymbolInfo
		symbolParams []string
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

	for _, reqSymbol := range symbols {
		symbolParams = append(symbolParams, ReformatSymbols(reqSymbol.Symbol))
	}

	c.handler.SetDepthIncrementSnapshotGroupChannel(chDeltaMap, chFullMap)
	conf := &base.IncrementDepthConf{
		IsPublishDelta: isDelta,
		IsPublishFull:  isFull,
		CheckTimeSec:   5,
		DepthCapLevel:  1000,
		DepthLevel:     limit,
		//GetFullDepth:   c.GetFullDepth,
		Ctx: ctx,
	}
	request := Request{
		ProductIds: symbolParams,
		Channels:   []string{"level2_batch", "heartbeat"},
	}
	err := c.EstablishSnapshotConn(symbols, request, spot_api.WS_API_BASE_URL, c.handler.DepthIncrementSnapShotGroupHandle, ctx, conf, interval)
	return err
}

/*TODO Private Channels need authorization*/
func (c *CoinbaseSpotWebsocket) GetSignature(method string) string {
	//timestamp := time.Now().Unix
	//fmt.Println(timestamp)
	return ""
}

/*Helper Functions*/
func ReformatSymbols(input string) string {
	return strings.Replace(input, "/", "-", -1)
}
func FormatSymbols(input string) string {
	return strings.Replace(input, "-", "/", -1)
}
func generateSig(message, secret string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	signature := hmac.New(sha256.New, key)
	_, err = signature.Write([]byte(message))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature.Sum(nil)), nil
}
func (c *CoinbaseSpotWebsocket) Sign(initial Request) (SignedRequest, error) {

	method := "GET"
	url := "/users/self/verify"
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sig, err := generateSig(fmt.Sprintf("%s%s%s", timestamp, method, url), c.SecretKey)

	return SignedRequest{
		Request:    initial,
		Key:        c.AccessKey,
		Passphrase: c.Passphrase,
		Timestamp:  strconv.FormatInt(time.Now().Unix(), 10),
		Signature:  sig,
	}, err
}

/*Deprecated: Check Depth Functions*/
/*
func (c *CoinbaseSpotWebsocket) Reconnect(wsClient *conn.WsConn, ctx context.Context, request Request, symbol *client.SymbolInfo, ch chan bool) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			if err := c.Unsubscribe(wsClient, request); err != nil {
				logger.Logger.Error("unsubscribe err:", symbol.Symbol, err)
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
			fmt.Println("Reconnecting:", symbol.Symbol)
			if err := c.Unsubscribe(wsClient, request); err != nil {
				logger.Logger.Error("unsubscribe err:", symbol.Symbol, err)
			}
			if err := c.Subscribe(wsClient, request); err != nil {
				logger.Logger.Error("subscribe err:", symbol.Symbol, err)
			}
		}
	}
}
*/
/*
func (c *CoinbaseSpotWebsocket) GetFullDepth(symbol string) (*base.OrderBook, error) {
	var (
		request       Request
		ch            = make(chan *RespIncrementSnapshot, 1)
		depthCache    = &base.OrderBook{}
		incomingData  *RespIncrementSnapshot
		diff          = &base.DeltaDepthUpdate{}
		amount, price float64
		err           error
		readTimeout   = time.Duration(c.ReadTimeout) * time.Second * 1000
	)

	c.lock.Lock()
	defer c.lock.Unlock()

	ctx, cancel := context.WithCancel(context.Background()) //Can pass cancel in lieu of _

	request = Request{
		ProductIds: []string{ReformatSymbols(symbol)},
		Channels:   []string{"level2_batch"},
	}

	checker := func(data []byte) error {
		iData := &RespIncrementSnapshot{}
		err := json.Unmarshal(data, iData)
		ch <- iData
		return err
	}
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(spot_api.WS_API_BASE_URL).ProtoHandleFunc(checker)
	if c.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(c.ProxyUrl)
	}
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return depthCache, errors.New("build websocket connection error")
	}
	err = c.Subscribe(wsClient, request)
	if err != nil {
		return depthCache, err
	}
	go c.FullDepthStart(wsClient, ctx)
LOOP:
	for {
		select {
		case incomingData, _ = <-ch: //stop on "type" = "l2update" message
			if incomingData.Type == "snapshot" {
				depthCache.Asks, err = DepthLevelParse(incomingData.Asks)
				depthCache.Bids, err = DepthLevelParse(incomingData.Bids)
				depthCache.Symbol = FormatSymbols(incomingData.ProductId)
				depthCache.Market = common.Market_SPOT
			} else if incomingData.Type == "l2update" {
				break LOOP
			}
		}
	}
	cancel()

	for _, incrementInfo := range incomingData.Changes {
		price, err = strconv.ParseFloat(incrementInfo.Price, 64)
		amount, err = strconv.ParseFloat(incrementInfo.Size, 64)
		if err != nil {
			fmt.Println("Parse Error")
			return depthCache, err
		}
		switch incrementInfo.Side {
		case BUY:
			diff.Bids = append(diff.Bids, &depth.DepthLevel{
				Price:  price,
				Amount: amount,
			})
		case SELL:
			diff.Asks = append(diff.Asks, &depth.DepthLevel{
				Price:  price,
				Amount: amount,
			})
		}
	}
	sort.Sort(depthCache.Asks)
	sort.Sort(sort.Reverse(depthCache.Bids))
	sort.Sort(diff.Asks)
	sort.Sort(sort.Reverse(diff.Bids))
	base.MergeDepth(1, &depthCache.Asks, diff.Asks)
	base.MergeDepth(2, &depthCache.Bids, diff.Bids)

	depthCache.TimeExchange = uint64(incomingData.Time.UnixMicro())
	depthCache.Type = common.SymbolType_SPOT_NORMAL
	depthCache.Exchange = common.Exchange_COINBASE

	return depthCache, nil
}
func (c *CoinbaseSpotWebsocket) FullDepthStart(wsClient *conn.WsConn, ctx context.Context) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			wsClient.CloseWs()
			break LOOP
		}
	}
}
*/

/*Used to test different Channel requests*/
func (c *CoinbaseSpotWebsocket) Tester(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
	var (
		symbolInfo []*client.SymbolInfo
		symbols    []string
		request    Request
		err        error
	)
	for symbol := range chMap {
		symbolInfo = append(symbolInfo, symbol)
		symbols = append(symbols, ReformatSymbols(symbol.Symbol))
	}
	c.handler.SetTradeGroupChannel(chMap)
	request = Request{
		ProductIds: symbols,
		Channels:   []string{"full"},
	}
	err = c.EstablishConn(symbolInfo, request, spot_api.WS_API_BASE_URL, c.handler.Tester, ctx)
	return err
}
