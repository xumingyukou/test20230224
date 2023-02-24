package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/poloniex"
	"clients/logger"
	"clients/transform"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
)

type WebSocketPublicHandleInterface interface {
	FundingRateGroupHandle([]byte) error
	AggTradeGroupHandle([]byte) error
	TradeGroupHandle([]byte) error
	BookTickerGroupHandle([]byte) error
	DepthIncrementGroupHandle([]byte) error
	DepthLimitGroupHandle([]byte) error
	DepthIncrementSnapShotGroupHandle([]byte) error

	SetFundingRateGroupChannel(map[*client.SymbolInfo]chan *client.WsFundingRateRsp)
	SetTradeGroupChannel(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	SetBookTickerGroupChannel(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	SetDepthLimitGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	SetDepthIncrementGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	SetDepthIncrementSnapshotGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp, map[*client.SymbolInfo]chan *depth.Depth)
	SetAggTradeGroupChannel(map[*client.SymbolInfo]chan *client.WsAggTradeRsp) // binance

	SetDepthIncrementSnapShotConf([]*client.SymbolInfo, *base.IncrementDepthConf)
	// GetChan(chName string) interface{}
}

type WebSocketSpotHandle struct {
	*base.IncrementDepthConf

	accountChan     chan *client.WsAccountRsp
	balanceChan     chan *client.WsBalanceRsp
	orderChan       chan *client.WsOrderRsp
	depthWholeLimit int

	fundingRateGroupChanMap                 map[string]chan *client.WsFundingRateRsp
	aggTradeGroupChanMap                    map[string]chan *client.WsAggTradeRsp
	tradeGroupChanMap                       map[string]chan *client.WsTradeRsp
	bookTickerGroupChanMap                  map[string]chan *client.WsBookTickerRsp
	depthLimitGroupChanMap                  map[string]chan *depth.Depth       //全量
	depthIncrementGroupChanMap              map[string]chan *client.WsDepthRsp //增量
	depthIncrementSnapshotDeltaGroupChanMap map[string]chan *client.WsDepthRsp //增量合并数据
	depthIncrementSnapshotFullGroupChanMap  map[string]chan *depth.Depth       //增量合并数据
	Lock                                    sync.Mutex

	symbolMap                             map[string]*client.SymbolInfo
	DepthIncrementSnapshotReconnectSymbol chan string
	reconnectSignalChannel                map[string]chan struct{}
	CheckSendStatus                       *base.CheckDataSendStatus
}

func NewWebSocketSpotHandle(chanCap int64) *WebSocketSpotHandle {
	h := &WebSocketSpotHandle{}
	// h.DepthCacheMap = sdk.NewCmapI()
	// h.DepthCacheListMap = sdk.NewCmapI()
	// h.CheckDepthCacheChanMap = sdk.NewCmapI()
	// h.DepthNotMatchChanMap = make(map[*client.SymbolInfo]chan bool)

	h.accountChan = make(chan *client.WsAccountRsp, chanCap)
	h.balanceChan = make(chan *client.WsBalanceRsp, chanCap)
	h.orderChan = make(chan *client.WsOrderRsp, chanCap)
	h.CheckSendStatus = base.NewCheckDataSendStatus()

	return h
}

func (b *WebSocketSpotHandle) GetChan(chName string) interface{} {
	switch chName {
	case "send_status":
		return b.CheckSendStatus
	default:
		return nil
	}
}

func (h *WebSocketSpotHandle) GetHandler(name string) func([]byte) error {
	switch name {
	case "trade":
		return h.HandlerWrapper(name, h.TradeGroupHandle)
	case "ticker":
		return h.HandlerWrapper(name, h.BookTickerGroupHandle)
	case "depthDiff":
		return h.HandlerWrapper(name, h.DepthIncrementGroupHandle)
	case "depthWhole":
		return h.HandlerWrapper(name, h.DepthLimitGroupHandle)
	case "depthIncrementSnapshot":
		return h.HandlerWrapper(name, h.DepthIncrementSnapShotGroupHandle)
	}
	return nil
}

func (h *WebSocketSpotHandle) TradeGroupHandle(data []byte, receiveTime int64) error {
	var (
		resp      = &TradeMsg{}
		err       error
		tradeSide order.TradeSide
	)

	err = json.Unmarshal(data, resp)
	if err != nil {
		return err
	}
	for _, t := range resp.Data {
		symbol := poloniex.Exchange2Canonical(t.Symbol)
		price, err := transform.Str2Float64(t.Price)
		if err != nil {
			logger.Logger.Error("trade parse ", err, " data:", string(data))
			return err
		}
		volume, err := transform.Str2Float64(t.Quantity)
		if err != nil {
			logger.Logger.Error("trade parse ", err, " data:", string(data))
			return err
		}
		if t.TakerSide == "buy" {
			tradeSide = order.TradeSide_BUY
		} else if t.TakerSide == "sell" {
			tradeSide = order.TradeSide_SELL
		} else {
			logger.Logger.Error("unknown trade side", " data:", string(data))
			return errors.New("unknown trade side")
		}

		res := &client.WsTradeRsp{
			Symbol:       symbol,
			OrderId:      t.ID,
			ExchangeTime: t.Ts * 1000,
			Price:        price,
			Amount:       volume,
			TakerSide:    tradeSide,
			ReceiveTime:  receiveTime,
			DealTime:     t.CreateTime * 1000,
		}

		if ch, ok := h.tradeGroupChanMap[symbol]; ok {
			base.SendChan(ch, res, "trade")
		} else {
			logger.Logger.Warn("get symbol from channel map err:", symbol)
		}
		h.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, receiveTime)
	}
	return nil
}

func removePrefix(s string, prefix string) string {
	return s[len(prefix):]
}

func (h *WebSocketSpotHandle) BookTickerGroupHandle(data []byte, receiveTime int64) error {
	return nil
}

func (h *WebSocketSpotHandle) DepthLevelParse(levelList [][]string) ([]*depth.DepthLevel, error) {
	var (
		res           []*depth.DepthLevel
		amount, price float64
		err           error
	)
	for _, level := range levelList {
		if price, err = transform.Str2Float64(level[0]); err != nil {
			return nil, err
		}
		if amount, err = transform.Str2Float64(level[1]); err != nil {
			return nil, err
		}
		res = append(res, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	return res, err
}

func (h *WebSocketSpotHandle) parseDepthDelta(msg *DepthMsg, receiveTime int64) (*base.DeltaDepthUpdate, error) {
	asks, err := h.DepthLevelParse(msg.Asks)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, msg)
		return nil, err
	}
	bids, err := h.DepthLevelParse(msg.Bids)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, msg)
		return nil, err
	}

	res := &base.DeltaDepthUpdate{
		UpdateStartId: int64(msg.LastID),
		UpdateEndId:   int64(msg.ID),

		Market:       common.Market_SPOT,
		Type:         common.SymbolType_SPOT_NORMAL,
		Symbol:       poloniex.Exchange2Canonical(msg.Symbol),
		TimeExchange: msg.Ts * 1000,
		TimeReceive:  receiveTime,
		Asks:         asks,
		Bids:         bids,
	}
	return res, nil
}

func (h *WebSocketSpotHandle) parseDepthWhole(msg *DepthMsg, receiveTime int64) (*base.OrderBook, error) {
	asks, err := h.DepthLevelParse(msg.Asks)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, msg)
		return nil, err
	}
	bids, err := h.DepthLevelParse(msg.Bids)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, msg)
		return nil, err
	}

	res := &base.OrderBook{
		UpdateId:     int64(msg.ID),
		Exchange:     common.Exchange_POLONIEX,
		Market:       common.Market_SPOT,
		Type:         common.SymbolType_SPOT_NORMAL,
		Symbol:       poloniex.Exchange2Canonical(msg.Symbol),
		TimeExchange: uint64(msg.Ts) * 1000,
		TimeReceive:  uint64(receiveTime),
		Asks:         asks,
		Bids:         bids,
	}
	return res, nil
}

func (h *WebSocketSpotHandle) DepthIncrementGroupHandle(data []byte, receiveTime int64) error {

	var (
		resp = &BookLv2Msg{}
		err  error
	)
	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Logger.Error("depth increment parse ", err, " data:", string(data))
		return err
	}
	// full
	if resp.Action == "snapshot" {
		return nil
	}
	if resp.Action != "update" {
		return errors.New("unknown increment type:" + resp.Action)
	}

	for _, t := range resp.Data {
		delta, err := h.parseDepthDelta(&t, receiveTime)
		if err != nil {
			return err
		}
		res := delta.Transfer2Depth()

		if ch, ok := h.depthIncrementGroupChanMap[res.Symbol]; ok {
			base.SendChan(ch, res, "depthIncrement")
		} else {
			logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
		}
		h.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, receiveTime)
	}

	return nil
}

func (h *WebSocketSpotHandle) SetDepthLimit(limit int) {
	h.depthWholeLimit = limit
}

func (h *WebSocketSpotHandle) DepthLimitGroupHandle(data []byte, receiveTime int64) error {
	var (
		err  error
		resp = &BookMsg{}
	)

	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Logger.Error("depth limit parse:", err, " data:", string(data))
		return err
	}

	for _, t := range resp.Data {
		asks, err := h.DepthLevelParse(t.Asks)
		if err != nil {
			logger.Logger.Error("depth level parse err", err, resp)
			return err
		}
		bids, err := h.DepthLevelParse(t.Bids)
		if err != nil {
			logger.Logger.Error("depth level parse err", err, resp)
			return err
		}

		orderBook := &base.OrderBook{
			Exchange:     common.Exchange_POLONIEX,
			Market:       common.Market_SPOT,
			Type:         common.SymbolType_SPOT_NORMAL,
			Symbol:       poloniex.Exchange2Canonical(t.Symbol),
			TimeExchange: uint64(t.Ts) * 1000,
			TimeReceive:  uint64(receiveTime),
			Asks:         asks,
			Bids:         bids,
			UpdateId:     int64(t.ID),
		}

		res := &depth.Depth{}
		orderBook.Transfer2Depth(h.depthWholeLimit, res)

		if ch, ok := h.depthLimitGroupChanMap[res.Symbol]; ok {
			base.SendChan(ch, res, "depthLimit")
		} else {
			logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
		}
		h.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, receiveTime)
	}

	return nil
}

func (h *WebSocketSpotHandle) HandlerWrapper(handleName string, handler func([]byte, int64) error) func([]byte) error {
	return func(data []byte) error {
		var (
			err error
			msg = &EventMessage{}
		)
		receiveTime := time.Now().UnixMicro()
		err = json.Unmarshal(data, msg)
		if err != nil {
			logger.Logger.Error("General Handle:", err.Error())
			return errors.New("json unmarshal data err:" + string(data))
		}
		// fmt.Println(msg, string(data))
		switch strings.ToLower(msg.Event) {
		case "pong":
		case "subscribe":
		case "unsubscribe":
		case "unsubscribe_all":
		case "error":
			logger.Logger.Error("subscribe error ", msg.Message)
		default:
			return handler(data, receiveTime)
		}
		return nil
	}
}

func (h *WebSocketSpotHandle) SetTradeGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) {
	if chMap != nil {
		if h.tradeGroupChanMap == nil {
			h.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp)
		}
		for symbol, ch := range chMap {
			h.tradeGroupChanMap[symbol.Symbol] = ch
		}
	}
}

func (h *WebSocketSpotHandle) SetBookTickerGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) {
	if chMap != nil {
		if h.bookTickerGroupChanMap == nil {
			h.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp)
		}
		for symbol, ch := range chMap {
			h.bookTickerGroupChanMap[symbol.Symbol] = ch
		}
	}
}

func (h *WebSocketSpotHandle) SetDepthLimitGroupChannel(chMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chMap != nil {
		if h.depthLimitGroupChanMap == nil {
			h.depthLimitGroupChanMap = make(map[string]chan *depth.Depth)
		}
		for symbol, ch := range chMap {
			h.depthLimitGroupChanMap[symbol.Symbol] = ch
		}

	}
}

func (h *WebSocketSpotHandle) SetDepthIncrementGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) {
	if chMap != nil {
		if h.depthIncrementGroupChanMap == nil {
			h.depthIncrementGroupChanMap = make(map[string]chan *client.WsDepthRsp)
		}
		for symbol, ch := range chMap {
			h.depthIncrementGroupChanMap[symbol.Symbol] = ch
		}
	}
}

func (h *WebSocketSpotHandle) SetDepthIncrementSnapshotGroupChannel(chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chDeltaMap != nil {
		if h.depthIncrementSnapshotDeltaGroupChanMap == nil {
			h.depthIncrementSnapshotDeltaGroupChanMap = make(map[string]chan *client.WsDepthRsp)
		}
		for symbol, ch := range chDeltaMap {
			h.depthIncrementSnapshotDeltaGroupChanMap[symbol.Symbol] = ch
		}
	}
	if chFullMap != nil {
		if h.depthIncrementSnapshotFullGroupChanMap == nil {
			h.depthIncrementSnapshotFullGroupChanMap = make(map[string]chan *depth.Depth)
		}
		for symbol, ch := range chFullMap {
			h.depthIncrementSnapshotFullGroupChanMap[symbol.Symbol] = ch
		}
	}
}
