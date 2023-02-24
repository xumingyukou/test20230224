package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/bithumb"
	"clients/logger"
	"clients/transform"
	"encoding/json"
	"errors"
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

	symbolMap       map[string]*client.SymbolInfo
	CheckSendStatus *base.CheckDataSendStatus
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
		logger.Logger.Error("trade group handle:", err.Error(), " data:", string(data))
		return err
	}
	symbol := bithumb.Exchange2Canonical(resp.Data.Symbol)
	price, err := transform.Str2Float64(resp.Data.P)
	if err != nil {
		logger.Logger.Error("trade parse ", err, string(data))
		return err
	}
	volume, err := transform.Str2Float64(resp.Data.V)
	if err != nil {
		logger.Logger.Error("trade parse ", err, string(data))
		return err
	}
	if resp.Data.S == "buy" {
		tradeSide = order.TradeSide_BUY
	} else if resp.Data.S == "sell" {
		tradeSide = order.TradeSide_SELL
	} else {
		logger.Logger.Error("unknown trade side", string(data))
		return errors.New("unknown trade side")
	}

	dealTime, err := transform.Str2Int64(resp.Data.T)
	if err != nil {
		logger.Logger.Error("trade parse ", err, string(data))
		return err
	}

	res := &client.WsTradeRsp{
		Symbol:       symbol,
		ExchangeTime: resp.Timestamp * 1000,
		Price:        price,
		Amount:       volume,
		TakerSide:    tradeSide,
		ReceiveTime:  receiveTime,
		DealTime:     dealTime * 1000 * 1000,
	}

	if ch, ok := h.tradeGroupChanMap[symbol]; ok {
		base.SendChan(ch, res, "trade")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", symbol)
	}
	h.CheckSendStatus.CheckUpdateTimeMap.Set(symbol, receiveTime)
	return nil
}

func (h *WebSocketSpotHandle) BookTickerGroupHandle(data []byte, receiveTime int64) error {
	var (
		resp = &TickerMsg{}
		err  error
	)

	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Logger.Error("ticker group handle:", err.Error(), " data:", string(data))
		return err
	}

	price, err := transform.Str2Float64(resp.Data.C)
	if err != nil {
		logger.Logger.Error("ticker parse ", err, string(data))
		return err
	}
	updateId, err := transform.Str2Int64(resp.Data.Ver)
	if err != nil {
		logger.Logger.Error("ticker parse ", err, string(data))
		return err
	}

	res := &client.WsBookTickerRsp{
		ExchangeTime: resp.Timestamp * 1000,
		ReceiveTime:  receiveTime,
		Symbol:       bithumb.Exchange2Canonical(resp.Data.Symbol),
		Ask: &depth.DepthLevel{
			Price: price,
		},
		Bid: &depth.DepthLevel{
			Price: price,
		},
		UpdateIdStart: updateId,
		UpdateIdEnd:   updateId,
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

	var (
		resp = &OrderBookMsg{}
		err  error
	)
	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Logger.Error("depthincrement group handle:", err.Error(), " data:", string(data))
		return err
	}
	// full
	if resp.Code == "00006" {
		return nil
	}

	delta, err := h.parseDepthDelta(resp, receiveTime)
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
	return nil
}

func (h *WebSocketSpotHandle) parseDepthDelta(resp *OrderBookMsg, receiveTime int64) (*base.DeltaDepthUpdate, error) {

	asks, err := h.DepthLevelParse(resp.Data.S)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, resp)
		return nil, err
	}
	bids, err := h.DepthLevelParse(resp.Data.B)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, resp)
		return nil, err
	}
	updateId, err := transform.Str2Int64(resp.Data.Ver)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, resp)
		return nil, err
	}

	res := &base.DeltaDepthUpdate{
		UpdateStartId: updateId,
		UpdateEndId:   updateId,
		Market:        common.Market_SPOT,
		Type:          common.SymbolType_SPOT_NORMAL,
		Symbol:        bithumb.Exchange2Canonical(resp.Data.Symbol),
		TimeReceive:   receiveTime,
		TimeExchange:  resp.Timestamp * 1000,
		Asks:          asks,
		Bids:          bids,
	}

	return res, nil
}

func (h *WebSocketSpotHandle) parseDepthWhole(resp *OrderBookMsg, receiveTime int64) (*base.OrderBook, error) {
	asks, err := h.DepthLevelParse(resp.Data.S)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, resp)
		return nil, err
	}
	bids, err := h.DepthLevelParse(resp.Data.B)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, resp)
		return nil, err
	}
	updateId, err := transform.Str2Int64(resp.Data.Ver)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, resp)
		return nil, err
	}

	res := &base.OrderBook{
		Exchange:     common.Exchange_BITHUMB,
		Market:       common.Market_SPOT,
		Type:         common.SymbolType_SPOT_NORMAL,
		Symbol:       bithumb.Exchange2Canonical(resp.Data.Symbol),
		TimeExchange: uint64(resp.Timestamp * 1000),
		TimeReceive:  uint64(receiveTime),
		Asks:         asks,
		Bids:         bids,
		UpdateId:     updateId,
	}
	return res, nil
}

func (h *WebSocketSpotHandle) SetDepthLimit(limit int) {
	h.depthWholeLimit = limit
}

func (h *WebSocketSpotHandle) HandlerWrapper(handleName string, handler func([]byte, int64) error) func([]byte) error {
	return func(data []byte) error {
		var (
			err error
			msg = &BaseMsg{}
		)
		receiveTime := time.Now().UnixMicro()
		err = json.Unmarshal(data, msg)
		if err != nil {
			logger.Logger.Error("General Handle:", err.Error())
			return errors.New("json unmarshal data err:" + string(data))
		}
		switch msg.Code {
		case "0":
		case "00001":
		case "00002":
		case "00003":
		case "00006":
			if handleName == "depthIncrementSnapshot" {
				handler(data, receiveTime)
			}
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
