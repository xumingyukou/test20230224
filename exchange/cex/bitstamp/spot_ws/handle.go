package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"encoding/json"
	"fmt"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
	"strconv"
	"strings"
	"sync"
	"time"
)

type WebSocketHandleInterface interface {
	WebSocketPublicHandleInterface
	WebSocketPrivateHandleInterface
}

type WebSocketPublicHandleInterface interface {
	FundingRateGroupHandle([]byte) error
	AggTradeGroupHandle([]byte) error
	TradeGroupHandle([]byte) error
	BookTickerGroupHandle([]byte) error
	DepthIncrementGroupHandle([]byte) error
	DepthLimitGroupHandle([]byte) error
	DepthIncrementSnapShotGroupHandle([]byte) error

	//设置chan map
	SetFundingRateGroupChannel(map[*client.SymbolInfo]chan *client.WsFundingRateRsp)
	SetTradeGroupChannel(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	SetBookTickerGroupChannel(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	SetDepthLimitGroupChannel(map[*client.SymbolInfo]chan *depth.Depth)
	SetDepthIncrementGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	SetDepthIncrementSnapshotGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp, map[*client.SymbolInfo]chan *depth.Depth)
	SetAggTradeGroupChannel(map[*client.SymbolInfo]chan *client.WsAggTradeRsp) // binance

	SetDepthIncrementSnapShotConf([]*client.SymbolInfo, *base.IncrementDepthConf)
}

type WebSocketPrivateHandleInterface interface {
	AccountHandle([]byte) error
	BalanceHandle([]byte) error
	MarginAccountHandle([]byte) error
	MarginBalanceHandle([]byte) error
	OrderHandle([]byte) error
	GetChan(chName string) interface{}
}

type WebSocketSpotHandle struct {
	*base.IncrementDepthConf

	accountChan chan *client.WsAccountRsp
	balanceChan chan *client.WsBalanceRsp
	orderChan   chan *client.WsOrderRsp

	fundingRateGroupChanMap                 map[string]chan *client.WsFundingRateRsp
	aggTradeGroupChanMap                    map[string]chan *client.WsAggTradeRsp
	tradeGroupChanMap                       map[string]chan *client.WsTradeRsp
	bookTickerGroupChanMap                  map[string]chan *client.WsBookTickerRsp
	depthLimitGroupChanMap                  map[string]chan *depth.Depth       //全量
	depthIncrementGroupChanMap              map[string]chan *client.WsDepthRsp //增量
	depthIncrementSnapshotDeltaGroupChanMap map[string]chan *client.WsDepthRsp //增量合并数据
	depthIncrementSnapshotFullGroupChanMap  map[string]chan *depth.Depth       //增量合并数据

	symbolMap       map[string]*client.SymbolInfo
	Lock            sync.Mutex
	CheckSendStatus *base.CheckDataSendStatus
}

func NewWebSocketSpotHandle(chanCap int64) *WebSocketSpotHandle {
	b := &WebSocketSpotHandle{}

	b.accountChan = make(chan *client.WsAccountRsp, chanCap)
	b.balanceChan = make(chan *client.WsBalanceRsp, chanCap)
	b.orderChan = make(chan *client.WsOrderRsp, chanCap)
	b.CheckSendStatus = base.NewCheckDataSendStatus()

	return b
}

/*Public Interface Methods*/
func (b *WebSocketSpotHandle) TradeGroupHandle(data []byte) error {
	var (
		resp      RespTradeStream
		timeStart = time.Now().UnixMicro()
		err       error
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if resp.Event == "bts:subscription_succeeded" {
		return nil
	} else if resp.Event == "bts:heartbeat" {
		return nil
	}

	timeStamp, err := strconv.ParseInt(resp.Data.Microtimestamp, 10, 64)
	if err != nil {
		logger.Logger.Error("parse time err:", string(data))
		return err
	}
	symbol := totalSymbolMap[strings.Split(resp.Channel, "_")[2]]

	res := &client.WsTradeRsp{
		OrderId:       strconv.Itoa(resp.Data.Id),
		Symbol:        symbol,
		ExchangeTime:  timeStamp,
		ReceiveTime:   timeStart,
		Price:         resp.Data.Price,
		Amount:        resp.Data.Amount,
		TakerSide:     GetSide(resp.Data.Type),
		BuyerOrderId:  strconv.FormatInt(resp.Data.BuyOrderId, 10),
		SellerOrderId: strconv.FormatInt(resp.Data.SellOrderId, 10),
		DealTime:      timeStamp,
	}

	if _, ok := b.tradeGroupChanMap[res.Symbol]; ok {
		base.SendChan(b.tradeGroupChanMap[res.Symbol], res, "trade")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}
func (b *WebSocketSpotHandle) BookTickerGroupHandle(data []byte) error {
	fmt.Println(string(data))
	fmt.Println("Bitstamp does not support Ticker")
	return nil
}
func (b *WebSocketSpotHandle) FundingRateGroupHandle(data []byte) error {
	fmt.Println(string(data))
	fmt.Println("Bitstamp does not support Funding Rate")
	return nil
}
func (b *WebSocketSpotHandle) DepthLimitGroupHandle(data []byte) error {
	var (
		resp          RespDepthStream
		asks          []*depth.DepthLevel
		bids          []*depth.DepthLevel
		timeStart     = time.Now().UnixMicro()
		amount, price float64
		err           error
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		fmt.Println(err)
		logger.Logger.Error("receive data err:", string(data))
		return err
	}

	if resp.Event == "bts:subscription_succeeded" {
		return nil
	} else if resp.Event == "bts:heartbeat" {
		return nil
	}

	timeStamp, err := strconv.ParseInt(resp.Data.MicroTimestamp, 10, 64)
	if err != nil {
		logger.Logger.Error("parse time err:", string(data))
		return err
	}
	symbol := totalSymbolMap[strings.Split(resp.Channel, "_")[2]]

	for _, info := range resp.Data.Asks {
		if price, err = strconv.ParseFloat(info.Price, 64); err != nil {
			return err
		}
		if amount, err = strconv.ParseFloat(info.Amount, 64); err != nil {
			return err
		}
		asks = append(asks, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}

	for _, info := range resp.Data.Bids {
		if price, err = strconv.ParseFloat(info.Price, 64); err != nil {
			return err
		}
		if amount, err = strconv.ParseFloat(info.Amount, 64); err != nil {
			return err
		}
		bids = append(bids, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}

	res := &depth.Depth{
		Exchange:     common.Exchange_BITSTAMP,
		Symbol:       symbol,
		TimeExchange: uint64(timeStamp),
		TimeReceive:  uint64(timeStart),
		TimeOperate:  uint64(timeStamp),
		Asks:         asks,
		Bids:         bids,
	}

	if _, ok := b.depthLimitGroupChanMap[res.Symbol]; ok {
		base.SendChan(b.depthLimitGroupChanMap[res.Symbol], res, "depthLimitGroup")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}
func (b *WebSocketSpotHandle) DepthIncrementGroupHandle(data []byte) error {
	var (
		resp          RespDepthStream
		asks          []*depth.DepthLevel
		bids          []*depth.DepthLevel
		timeStart     = time.Now().UnixMicro()
		amount, price float64
		err           error
	)

	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if resp.Event == "bts:subscription_succeeded" {
		return nil
	} else if resp.Event == "bts:heartbeat" {
		return nil
	}

	timeStamp, err := strconv.ParseInt(resp.Data.MicroTimestamp, 10, 64)
	if err != nil {
		logger.Logger.Error("parse time err:", string(data))
		return err
	}
	symbol := totalSymbolMap[strings.Split(resp.Channel, "_")[3]]

	for _, info := range resp.Data.Asks {
		if price, err = strconv.ParseFloat(info.Price, 64); err != nil {
			return err
		}
		if amount, err = strconv.ParseFloat(info.Amount, 64); err != nil {
			return err
		}
		asks = append(asks, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}

	for _, info := range resp.Data.Bids {
		if price, err = strconv.ParseFloat(info.Price, 64); err != nil {
			return err
		}
		if amount, err = strconv.ParseFloat(info.Amount, 64); err != nil {
			return err
		}
		bids = append(bids, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}

	res := &client.WsDepthRsp{
		ExchangeTime: timeStamp,
		ReceiveTime:  timeStart,
		Symbol:       symbol,
		Asks:         asks,
		Bids:         bids,
	}

	if _, ok := b.depthIncrementGroupChanMap[res.Symbol]; ok {
		base.SendChan(b.depthIncrementGroupChanMap[res.Symbol], res, "depthIncrement")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}

/*TODO*/
func (b *WebSocketSpotHandle) GetChan(chName string) interface{} {
	switch chName {
	case "account":
		return b.accountChan
	case "balance":
		return b.balanceChan
	case "order":
		return b.orderChan
	case "send_status":
		return b.CheckSendStatus
	default:
		return nil
	}
}

/*Unused Functions*/
func (b *WebSocketSpotHandle) AggTradeGroupHandle(data []byte) error {
	//TODO implement me
	panic("implement me")
}
func (b *WebSocketSpotHandle) SetAggTradeGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsAggTradeRsp) {
	//TODO implement me
	panic("implement me")
}

/*Channel Setters*/
func (b *WebSocketSpotHandle) SetTradeGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) {
	if chMap != nil {
		b.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp)
		for info, ch := range chMap {
			b.tradeGroupChanMap[info.Symbol] = ch
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
func (b *WebSocketSpotHandle) SetFundingRateGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsFundingRateRsp) {
	//TODO implement me
	panic("implement me")
}
func (b *WebSocketSpotHandle) SetDepthLimitGroupChannel(chMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chMap != nil {
		b.depthLimitGroupChanMap = make(map[string]chan *depth.Depth)
		for info, ch := range chMap {
			b.depthLimitGroupChanMap[info.Symbol] = ch
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

/*TODO Private Functions*/
func (b *WebSocketSpotHandle) AccountHandle(data []byte) error {
	//TODO implement me
	panic("implement me")
}
func (b *WebSocketSpotHandle) BalanceHandle(data []byte) error {
	//TODO implement me
	panic("implement me")
}
func (b *WebSocketSpotHandle) MarginAccountHandle(data []byte) error {
	//TODO implement me
	panic("implement me")
}
func (b *WebSocketSpotHandle) MarginBalanceHandle(data []byte) error {
	//TODO implement me
	panic("implement me")
}
func (b *WebSocketSpotHandle) OrderHandle(data []byte) error {
	//TODO implement me
	panic("implement me")
}

/*Other Helper Functions*/
func GetSide(side TradeType) order.TradeSide {
	switch side {
	case BUY:
		return order.TradeSide_BUY
	case SELL:
		return order.TradeSide_SELL
	default:
		return order.TradeSide_InvalidSide
	}
}
