package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
)

type WebSocketHandleInterface interface {
	TradeGroupHandle([]byte) error
	BookTickerGroupHandle([]byte) error
	DepthLimitGroupHandle([]byte) error
	DepthIncrementGroupHandle([]byte) error
	DepthIncrementSnapShotGroupHandle([]byte) error

	SetTradeGroupChannel(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	SetBookTickerGroupChannel(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	SetDepthLimitGroupChannel(map[*client.SymbolInfo]chan *depth.Depth)
	SetDepthIncrementGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	SetDepthIncrementSnapshotGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp, map[*client.SymbolInfo]chan *depth.Depth)

	SetDepthIncrementSnapShotConf([]*client.SymbolInfo, *base.IncrementDepthConf)

	//GetChan(chName string) interface{}
}

type WebSocketSpotHandle struct {
	*base.IncrementDepthConf

	tradeGroupChanMap                       map[string]chan *client.WsTradeRsp
	bookTickerGroupChanMap                  map[string]chan *client.WsBookTickerRsp
	depthLimitGroupChanMap                  map[string]chan *depth.Depth       //全量
	depthIncrementGroupChanMap              map[string]chan *client.WsDepthRsp //增量
	depthIncrementSnapshotDeltaGroupChanMap map[string]chan *client.WsDepthRsp //增量合并数据
	depthIncrementSnapshotFullGroupChanMap  map[string]chan *depth.Depth       //增量合并数据
	depthIncrJudgeMap                       sync.Map

	symbolMap                             map[string]*client.SymbolInfo
	Lock                                  sync.Mutex
	CheckSendStatus                       *base.CheckDataSendStatus
	symbol2chanID                         map[string]int64 // client_symbol -> chanID
	chanID2symbol                         map[int64]string // chanID -> exchange_symbol
	countMap                              map[int64]int64  // chanID -> msgcount
	DepthIncrementSnapshotReconnectSymbol chan string
}

func NewWebSocketSpotHandle(chanCap int64) *WebSocketSpotHandle {
	d := &WebSocketSpotHandle{}
	d.CheckSendStatus = base.NewCheckDataSendStatus()
	d.symbol2chanID = make(map[string]int64)
	d.chanID2symbol = make(map[int64]string)
	d.countMap = make(map[int64]int64)
	return d
}

func (b *WebSocketSpotHandle) GetChanID(symbol string) int64 {
	return b.symbol2chanID[symbol]
}

func (b *WebSocketSpotHandle) GetSymbol(chanID int64) string {
	return b.chanID2symbol[chanID]
}

func (b *WebSocketSpotHandle) GetChan(chName string) interface{} {
	switch chName {
	case "send_status":
		return b.CheckSendStatus
	default:
		return nil
	}
}

func (b *WebSocketSpotHandle) GeneralHandle(name string) func([]byte) error {
	hand := func(data []byte) error {
		t := time.Now().UnixMicro()
		var (
			respInitial initialResp
			respFirst   Response
			err         error
		)
		json.Unmarshal(data, &respInitial)
		if respInitial.Event == "info" {
			return nil
		}
		if respInitial.Event == "error" {
			logger.Logger.Info("Error:", string(data))
			return nil
		}
		if respInitial.Event == "subscribed" {
			logger.Logger.Info("Subscribed OK:", respInitial)
			chanID := int64(respInitial.ChanID)
			// b.countMap[chanID] = 1
			b.chanID2symbol[chanID] = respInitial.Pair
			b.symbol2chanID[Exchange2Client[respInitial.Pair]] = chanID
			return nil
		}
		if respInitial.Event == "unsubscribed" {
			chanID := int64(respInitial.ChanID)
			delete(b.countMap, chanID)
			delete(b.symbol2chanID, Exchange2Client[b.chanID2symbol[chanID]])
			delete(b.chanID2symbol, chanID)
			// fmt.Println(respInitial.ChanID, symbol2chanID)
			logger.Logger.Info("UnSubscribed:", string(data))
			return nil
		}
		if respInitial.Event == "conf" {
			logger.Logger.Info("Set conf:", string(data))
			return nil
		}
		if err = json.Unmarshal(data, &respFirst); err != nil {
			logger.Logger.Error("receive data err:", err.Error(), " data:", string(data))
			return err
		}
		if respFirst[1] == "hb" {
			return nil
		}
		switch name {
		case "trades":
			return b.TradeGroupHandle(data, t)
		case "ticker":
			return b.BookTickerGroupHandle(data, t)
		case "increment":
			return b.DepthIncrementGroupHandle(data, t)
		case "snapshot":
			chanID := int64(respFirst[0].(float64))
			symbol := Exchange2Client[b.chanID2symbol[chanID]]
			if _, ok := b.depthIncrJudgeMap.Load(symbol); ok {
				return b.DepthIncrementSnapShotGroupHandle(data, t, false)
			} else {
				b.depthIncrJudgeMap.Store(symbol, 1)
				// logger.Logger.Info("First snapshot: ", symbol)
				return b.DepthIncrementSnapShotGroupHandle(data, t, true)
			}
		}
		return nil
	}
	return hand
}

func (b *WebSocketSpotHandle) TradeGroupHandle(data []byte, t int64) error {
	var (
		takerSide order.TradeSide
		resp      Response
		err       error
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	s, ok := resp[1].(string)
	if !ok {
		// drop trade snapshot
		return nil
	}
	if s == "te" {
		return nil
	} else if s != "tu" {
		logger.Logger.Error("unknown resp, ", resp[1])
		return nil
	}
	ChanID, ok1 := resp[0].(float64)
	Item, ok2 := resp[2].([]interface{})
	if ok1 && ok2 {
		price := Item[3].(float64)
		amount := Item[2].(float64)
		exchangeTime := Item[1].(float64)
		orderId := Item[0].(float64)
		if amount <= 0 {
			takerSide = order.TradeSide_SELL
			amount = -amount
		} else {
			takerSide = order.TradeSide_BUY
		}
		res := &client.WsTradeRsp{
			OrderId:      strconv.FormatInt(int64(orderId), 10),
			Symbol:       Exchange2Client[b.GetSymbol(int64(ChanID))],
			ExchangeTime: int64(exchangeTime) * 1000,
			ReceiveTime:  t,
			Price:        price,
			Amount:       amount,
			TakerSide:    takerSide,
		}
		if _, ok := b.tradeGroupChanMap[res.Symbol]; ok {
			base.SendChan(b.tradeGroupChanMap[res.Symbol], res, "trade")
		} else {
			logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
		}
		b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	} else {
		fmt.Println(ok1, ok2, resp)
	}
	return nil
}

func (b *WebSocketSpotHandle) BookTickerGroupHandle(data []byte, t int64) error {
	var (
		resp Response
		err  error
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	ChanID, ok1 := resp[0].(float64)
	Item, ok2 := resp[1].([]interface{})
	if ok1 && ok2 {
		bid_price, ok := Item[0].(float64)
		if !ok {
			logger.Logger.Error("receive data bid_price err:", string(data))
			return nil
		}
		bid_size, ok := Item[1].(float64)
		if !ok {
			logger.Logger.Error("receive data bid_size err:", string(data))
			return nil
		}
		bid := &depth.DepthLevel{
			Price:  bid_price,
			Amount: bid_size,
		}
		ask_price, ok := Item[2].(float64)
		if !ok {
			logger.Logger.Error("receive data ask_price err:", string(data))
			return nil
		}
		ask_size, ok := Item[3].(float64)
		if !ok {
			logger.Logger.Error("receive data ask_size err:", string(data))
			return nil
		}
		ask := &depth.DepthLevel{
			Price:  ask_price,
			Amount: ask_size,
		}
		res := &client.WsBookTickerRsp{
			ReceiveTime: t,
			Symbol:      Exchange2Client[b.GetSymbol(int64(ChanID))],
			Ask:         ask,
			Bid:         bid,
		}
		if _, ok := b.bookTickerGroupChanMap[res.Symbol]; ok {
			base.SendChan(b.bookTickerGroupChanMap[res.Symbol], res, "bookTickerGroup")
		} else {
			logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
		}
		b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	}
	return nil
}

func (b *WebSocketSpotHandle) DepthLimitGroupHandle(data []byte, t int64) error {
	return nil
}

func (b *WebSocketSpotHandle) DepthIncrementGroupHandle(data []byte, t int64) error {
	var (
		resp Response
		err  error
		asks []*depth.DepthLevel
		bids []*depth.DepthLevel
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	ChanID, ok1 := resp[0].(float64)
	Item, ok2 := resp[1].([]interface{})

	_, isArray := Item[0].([]interface{})
	if isArray {
		// drop depth snapshot
		return nil
	}
	if ok1 && ok2 {
		price := Item[0].(float64)
		count := Item[1].(float64)
		amount := Item[2].(float64)
		if count == 0 {
			amount = 0
		}
		if amount <= 0 {
			asks = append(asks, &depth.DepthLevel{
				Price:  price,
				Amount: -amount,
			})
		} else {
			bids = append(bids, &depth.DepthLevel{
				Price:  price,
				Amount: amount,
			})
		}
		res := &client.WsDepthRsp{
			ReceiveTime: t,
			Symbol:      Exchange2Client[b.GetSymbol(int64(ChanID))],
			Asks:        asks,
			Bids:        bids,
		}
		if _, ok := b.depthIncrementGroupChanMap[res.Symbol]; ok {
			base.SendChan(b.depthIncrementGroupChanMap[res.Symbol], res, "depthIncrement")
		} else {
			logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
		}
		b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	}
	return nil
}
func (b *WebSocketSpotHandle) SetTradeGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) {
	if chMap != nil {
		b.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp)
		for info, ch := range chMap {
			b.tradeGroupChanMap[info.Symbol] = ch
		}
	}
}

func (b *WebSocketSpotHandle) SetDepthIncrementGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) {
	if chMap != nil {
		b.depthIncrementGroupChanMap = make(map[string]chan *client.WsDepthRsp)
		for info, ch := range chMap {
			b.depthIncrementGroupChanMap[info.Symbol] = ch
		}
	}
}

func (b *WebSocketSpotHandle) SetDepthIncrementSnapshotGroupChannel(chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chDeltaMap != nil {
		b.depthIncrementSnapshotDeltaGroupChanMap = make(map[string]chan *client.WsDepthRsp)
		for info, ch := range chDeltaMap {
			sym := info.Symbol
			b.depthIncrementSnapshotDeltaGroupChanMap[sym] = ch
		}
	}
	if chFullMap != nil {
		b.depthIncrementSnapshotFullGroupChanMap = make(map[string]chan *depth.Depth)
		for info, ch := range chFullMap {
			sym := info.Symbol
			b.depthIncrementSnapshotFullGroupChanMap[sym] = ch
		}
	}
}

func (b *WebSocketSpotHandle) SetBookTickerGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) {
	if chMap != nil {
		b.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp)
		for info, ch := range chMap {
			b.bookTickerGroupChanMap[info.Symbol] = ch
		}
	}
}

func (b *WebSocketSpotHandle) SetDepthLimitGroupChannel(chMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chMap != nil {
		b.depthLimitGroupChanMap = make(map[string]chan *depth.Depth)
		for info, ch := range chMap {
			b.depthLimitGroupChanMap[info.Symbol] = ch
		}
	}
}
