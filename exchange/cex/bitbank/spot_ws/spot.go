package spot_ws

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/bitbank"
	"clients/exchange/cex/bitbank/spot_api"
	"clients/logger"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
)

type BitbankSpotWebsocket struct {
	base.WsConf
	apiClient     *spot_api.ApiClient
	listenTimeout int64
	// pingTimeout   int64
	// pingInterval  int64
	WsReqUrl *spot_api.WsReqUrl
	handler  *WebSocketSpotHandle
	reqId    int
	token    string
}

func NewBitbankWebsocket(conf base.WsConf) *BitbankSpotWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}

	b := &BitbankSpotWebsocket{
		WsConf: conf,
	}

	if b.ReadTimeout == 0 {
		b.ReadTimeout = 300
	}
	if b.listenTimeout == 0 {
		b.listenTimeout = 1800
	}
	b.apiClient = spot_api.NewApiClient(conf.APIConf)

	b.handler = NewWebSocketSpotHandle(b.ChanCap)
	b.WsReqUrl = spot_api.NewWsReqUrl()

	return b
}

func NewBitbankWebsocket2(conf base.WsConf, cli *http.Client) *BitbankSpotWebsocket {

	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}

	b := &BitbankSpotWebsocket{
		WsConf: conf,
	}
	if b.ReadTimeout == 0 {
		b.ReadTimeout = 300
	}
	if b.listenTimeout == 0 {
		b.listenTimeout = 1800
	}
	b.apiClient = spot_api.NewApiClient2(conf.APIConf, cli)

	b.handler = NewWebSocketSpotHandle(b.ChanCap)
	b.WsReqUrl = spot_api.NewWsReqUrl()

	return b
}

func (b *BitbankSpotWebsocket) baseUrl() string {
	return b.WsReqUrl.WS_BASE_URL
}

func (b *BitbankSpotWebsocket) EstablishConn(url string, room string, symbolInfo *client.SymbolInfo, handler func([]byte) error, ctx context.Context) error {
	heartBeatStream := func() [][]byte {
		return [][]byte{
			[]byte(fmt.Sprintf("42[\"join-room\", \"%s\"]", room)),
			[]byte("2probe"),
		}
	}

	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
	)

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).
		ProtoHandleFunc(handler).
		AutoReconnect().
		HeartbeatStream(heartBeatStream, time.Duration(15)*time.Second).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			b.Subscribe(wsClient, room)
			return nil
		})
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}

	dataSendStatus := b.handler.GetChan("send_status").(*base.CheckDataSendStatus)
	dataSendStatus.Init(3600, ctx, symbolInfo)

	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}

	b.Subscribe(wsClient, room)
	go b.start(wsClient, ctx)
	go b.reconnectAfterSilence(wsClient, ctx, dataSendStatus.UpdateTimeoutChMap)
	return nil
}

func (b *BitbankSpotWebsocket) EstablishSnapshotConn(url string, symbol string, handler func([]byte) error, ctx context.Context, symbolInfo *client.SymbolInfo, conf *base.IncrementDepthConf) error {
	heartBeatStream := func() [][]byte {
		return [][]byte{
			[]byte(fmt.Sprintf("42[\"join-room\", \"depth_diff_%s\"]", symbol)),
			[]byte(fmt.Sprintf("42[\"join-room\", \"depth_whole_%s\"]", symbol)),
			[]byte("2probe"),
		}
	}

	var (
		readTimeout = time.Duration(b.ReadTimeout) * time.Second * 1000
	)

	reconnectSignalChannel := make(chan struct{})
	reconnectSignalChannelMap := make(map[string]chan struct{})
	reconnectSignalChannelMap[symbolInfo.Symbol] = reconnectSignalChannel

	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).
		WsUrl(url).
		ProtoHandleFunc(handler).
		AutoReconnect().
		HeartbeatStream(heartBeatStream, time.Duration(15)*time.Second).
		PostReconnectSuccess(func(wsClient *conn.WsConn) error {
			reconnectSignalChannel <- struct{}{}
			b.Subscribe(wsClient, "depth_diff_"+symbol)
			b.Subscribe(wsClient, "depth_whole_"+symbol)
			return nil
		})
	if b.ProxyUrl != "" {
		wsBuilder = wsBuilder.ProxyUrl(b.ProxyUrl)
	}
	b.handler.SetReconnectSignalChannel(reconnectSignalChannelMap)

	b.handler.SetDepthIncrementSnapShotConf([]*client.SymbolInfo{symbolInfo}, conf)

	dataSendStatus := b.handler.GetChan("send_status").(*base.CheckDataSendStatus)
	dataSendStatus.Init(3600, ctx, symbolInfo)

	wsClient := wsBuilder.Build()
	if wsClient == nil {
		return errors.New("build websocket connection error")
	}
	b.Subscribe(wsClient, "depth_diff_"+symbol)
	b.Subscribe(wsClient, "depth_whole_"+symbol)
	go b.start(wsClient, ctx)
	go b.reconnectAfterSilence(wsClient, ctx, dataSendStatus.UpdateTimeoutChMap)
	return nil
}

func (b *BitbankSpotWebsocket) reconnectAfterSilence(wsClient *conn.WsConn, ctx context.Context, ch chan []*client.SymbolInfo) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			wsClient.CloseWs()
			break LOOP
		case <-ch:
			wsClient.Reconnect()
		}
	}
}

func (b *BitbankSpotWebsocket) start(wsClient *conn.WsConn, ctx context.Context) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			wsClient.CloseWs()
			break LOOP
		}
	}
}

func (b *BitbankSpotWebsocket) Subscribe(client *conn.WsConn, room string) {
	logger.Logger.Info("Bitbank subscribe:", room)
	client.SendMessage([]byte(fmt.Sprintf("42[\"join-room\", \"%s\"]", room)))
}

// func (b *BitbankSpotWebsocket) UnSubscribe(client *conn.WsConn) {
// 	client.SendMessage([]byte("exit"))
// }

func (b *BitbankSpotWebsocket) TradeGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) error {
	if len(chMap) != 1 {
		return errors.New("can only subscribe one symbol each time")
	}

	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}

	b.handler.SetTradeGroupChannel(chMap)
	err := b.EstablishConn(b.baseUrl(), "transactions_"+bitbank.Canonical2Exchange(symbols[0].Symbol), symbols[0], b.handler.GetHandler("trade"), ctx)
	return err
}

func (b *BitbankSpotWebsocket) BookTickerGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error {
	if len(chMap) != 1 {
		return errors.New("can only subscribe one symbol each time")
	}

	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}

	b.handler.SetBookTickerGroupChannel(chMap)
	err := b.EstablishConn(b.baseUrl(), "ticker_"+bitbank.Canonical2Exchange(symbols[0].Symbol), symbols[0], b.handler.GetHandler("ticker"), ctx)
	return err
}

func (b *BitbankSpotWebsocket) DepthIncrementGroup(ctx context.Context, interval int, chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) error {
	if len(chMap) != 1 {
		return errors.New("can only subscribe one symbol each time")
	}

	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}

	b.handler.SetDepthIncrementGroupChannel(chMap)
	err := b.EstablishConn(b.baseUrl(), "depth_diff_"+bitbank.Canonical2Exchange(symbols[0].Symbol), symbols[0], b.handler.GetHandler("depthDiff"), ctx)
	return err
}

func (b *BitbankSpotWebsocket) DepthLimitGroup(ctx context.Context, interval, limit int, chMap map[*client.SymbolInfo]chan *depth.Depth) error {
	if len(chMap) != 1 {
		return errors.New("can only subscribe one symbol each time")
	}

	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}

	b.handler.SetDepthLimitGroupChannel(chMap)
	b.handler.SetDepthLimit(limit)
	err := b.EstablishConn(b.baseUrl(), "depth_whole_"+bitbank.Canonical2Exchange(symbols[0].Symbol), symbols[0], b.handler.GetHandler("depthWhole"), ctx)
	return err
}

func (b *BitbankSpotWebsocket) DepthIncrementSnapshotGroup(ctx context.Context, interval, limit int, isDelta bool, isFull bool, chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) error {
	if isDelta && len(chDeltaMap) != 1 {
		return errors.New("can only subscribe one symbol each time")
	}

	if isFull && len(chFullMap) != 1 {
		return errors.New("can only subscribe one symbol each time")
	}

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

	b.handler.SetDepthIncrementSnapshotGroupChannel(chDeltaMap, chFullMap)

	conf := &base.IncrementDepthConf{
		IsPublishDelta:    isDelta,
		IsPublishFull:     isFull,
		CheckTimeSec:      3000 + rand.Intn(1200),
		DepthCapLevel:     1000,
		DepthLevel:        limit,
		GetFullDepth:      nil,
		GetFullDepthLimit: nil,
		Ctx:               ctx,
	}

	err := b.EstablishSnapshotConn(b.baseUrl(), bitbank.Canonical2Exchange(symbols[0].Symbol), b.handler.GetHandler("depthIncrementSnapshot"), ctx, symbols[0], conf)
	return err
}

func (k *BitbankSpotWebsocket) FundingRateGroup(ctx context.Context, chFundMap map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	return nil
}
