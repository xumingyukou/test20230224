package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/bitbank"
	"clients/logger"
	"clients/transform"
	"encoding/json"
	"fmt"
	"sort"
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
	CheckSendStatus                         *base.CheckDataSendStatus
	reconnectSignalChannel                  map[string]chan struct{} // 重连时发送保证记录最新全量, 保证不会同时读写
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

func (h *WebSocketSpotHandle) GetChan(chName string) interface{} {
	switch chName {
	case "send_status":
		return h.CheckSendStatus
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

func parseTrade(symbol string, tsc *Transaction, receiveTime int64) (*client.WsTradeRsp, error) {
	var (
		side   order.TradeSide
		err    error
		price  float64
		amount float64
	)
	if tsc.Side == "buy" {
		side = order.TradeSide_BUY
	} else {
		side = order.TradeSide_SELL
	}

	price, err = transform.Str2Float64(tsc.Price)
	if err != nil {
		return nil, err
	}
	amount, err = transform.Str2Float64(tsc.Amount)
	if err != nil {
		return nil, err
	}
	res := &client.WsTradeRsp{
		OrderId:     fmt.Sprintf("%d", tsc.TransactionID),
		Symbol:      symbol,
		Price:       price,
		Amount:      amount,
		TakerSide:   side,
		ReceiveTime: receiveTime,
		DealTime:    tsc.ExecutedAt * 1000,
	}

	return res, nil
}

func (h *WebSocketSpotHandle) TradeGroupHandle(data []byte, receiveTime int64) error {
	var (
		resp = &TradeMsg{}
		err  error
	)

	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Logger.Error("trade group handle:", err.Error(), " data:", string(data))
		return err
	}
	symbol := bitbank.Exchange2Canonical(removePrefix(resp.RoomName, "transactions_"))

	if ch, ok := h.tradeGroupChanMap[symbol]; ok {
		for _, tsc := range resp.Message.Data.Transactions {
			res, err := parseTrade(symbol, tsc, receiveTime)
			if err != nil {
				logger.Logger.Warn("parse trade err:", tsc)
			} else {
				base.SendChan(ch, res, "trade")
			}
		}
	} else {
		logger.Logger.Warn("get symbol from channel map err:", symbol)
	}
	h.CheckSendStatus.CheckUpdateTimeMap.Set(symbol, receiveTime)
	return nil
}

func removePrefix(s string, prefix string) string {
	return s[len(prefix):]
}

func (h *WebSocketSpotHandle) BookTickerGroupHandle(data []byte, receiveTime int64) error {
	var (
		resp = &TickerMsg{}
		err  error
	)

	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Logger.Error("book ticker group handle:", err.Error(), " data:", string(data))
		return err
	}

	sell, err := transform.Str2Float64(resp.Message.Data.Sell)
	if err != nil {
		logger.Logger.Error("book ticker group handle:", err.Error(), "resp: ", resp)
		return err
	}
	buy, err := transform.Str2Float64(resp.Message.Data.Buy)
	if err != nil {
		logger.Logger.Error("book ticker group handle:", err.Error(), "resp: ", resp)
		return err
	}

	res := &client.WsBookTickerRsp{
		ExchangeTime: resp.Message.Data.Timestamp * 1000,
		ReceiveTime:  receiveTime,
		Symbol:       bitbank.Exchange2Canonical(removePrefix(resp.RoomName, "ticker_")),
		Ask: &depth.DepthLevel{
			Price: sell,
		},
		Bid: &depth.DepthLevel{
			Price: buy,
		},
	}

	if ch, ok := h.bookTickerGroupChanMap[res.Symbol]; ok {
		base.SendChan(ch, res, "ticker")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	h.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, receiveTime)
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

func (h *WebSocketSpotHandle) DepthIncrementGroupHandle(data []byte, receiveTime int64) error {

	delta, err := h.parseDepthDelta(data, receiveTime)
	if err != nil {
		logger.Logger.Error("depth increment group handle:", err.Error(), " data:", string(data))
		return err
	}
	res := delta.Transfer2Depth()

	if ch, ok := h.depthIncrementGroupChanMap[res.Symbol]; ok {
		base.SendChan(ch, res, "depthIncrement")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	h.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, receiveTime)
	return nil
}

func (h *WebSocketSpotHandle) parseDepthDelta(data []byte, receiveTime int64) (*base.DeltaDepthUpdate, error) {
	var (
		resp = &DepthDiffMsg{}
		err  error
	)
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	asks, err := h.DepthLevelParse(resp.Message.Data.A)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return nil, err
	}
	bids, err := h.DepthLevelParse(resp.Message.Data.B)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return nil, err
	}
	updateId, err := transform.Str2Int64(resp.Message.Data.S)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return nil, err
	}

	res := &base.DeltaDepthUpdate{
		UpdateStartId: updateId,
		UpdateEndId:   updateId,
		Market:        common.Market_SPOT,
		Type:          common.SymbolType_SPOT_NORMAL,
		Symbol:        bitbank.Exchange2Canonical(removePrefix(resp.RoomName, "depth_diff_")),
		TimeReceive:   receiveTime,
		TimeExchange:  resp.Message.Data.T * 1000,
		Asks:          asks,
		Bids:          bids,
	}

	sort.Stable(res.Asks)
	sort.Stable(sort.Reverse(res.Bids))
	return res, nil
}

func (h *WebSocketSpotHandle) parseDepthWhole(data []byte, receiveTime int64) (*base.OrderBook, error) {
	var (
		resp = &DepthWholeMsg{}
		err  error
	)

	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	asks, err := h.DepthLevelParse(resp.Message.Data.Asks)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return nil, err
	}

	bids, err := h.DepthLevelParse(resp.Message.Data.Bids)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return nil, err
	}

	sequenceId, err := transform.Str2Int64(resp.Message.Data.SequenceId)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return nil, err
	}

	res := &base.OrderBook{
		Exchange:     common.Exchange_BITBANK,
		Market:       common.Market_SPOT,
		Type:         common.SymbolType_SPOT_NORMAL,
		Symbol:       bitbank.Exchange2Canonical(removePrefix(resp.RoomName, "depth_whole_")),
		TimeExchange: uint64(resp.Message.Data.Timestamp * 1000),
		TimeReceive:  uint64(receiveTime),
		Asks:         asks,
		Bids:         bids,
		UpdateId:     sequenceId,
	}

	sort.Stable(res.Asks)
	sort.Stable(sort.Reverse(res.Bids))
	return res, nil
}

func (h *WebSocketSpotHandle) SetDepthLimit(limit int) {
	h.depthWholeLimit = limit
}

func (h *WebSocketSpotHandle) DepthLimitGroupHandle(data []byte, receiveTime int64) error {
	var (
		err error
	)
	orderBook, err := h.parseDepthWhole(data, receiveTime)
	if err != nil {
		logger.Logger.Error("depth limit group handle:", err.Error(), " data:", string(data))
		return err
	}

	res := &depth.Depth{}

	orderBook.Transfer2Depth(h.depthWholeLimit, res)

	if ch, ok := h.depthLimitGroupChanMap[res.Symbol]; ok {
		base.SendChan(ch, res, "depthLimit")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	h.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, receiveTime)
	return nil
}

func (h *WebSocketSpotHandle) HandlerWrapper(handleName string, handler func([]byte, int64) error) func([]byte) error {
	return func(data []byte) error {
		receiveTime := time.Now().UnixMicro()
		l := len(data)
		if l >= 12 && strings.HasPrefix(string(data[:12]), "42[\"message\"") {
			return handler(data[13:l-1], receiveTime)
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
