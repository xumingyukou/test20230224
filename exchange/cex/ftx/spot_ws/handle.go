package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/ftx/future_api"
	"clients/exchange/cex/ftx/spot_api"
	"clients/logger"
	"clients/transform"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
	"github.com/warmplanet/proto/go/sdk"
)

var loc, _ = time.LoadLocation("Asia/Shanghai")

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

type WebSocketHandleInterface interface {
	WebSocketPublicHandleInterface
	WebSocketPrivateHandleInterface
}
type WebSocketPrivateHandleInterface interface {
	AccountHandle([]byte) error
	BalanceHandle([]byte) error
	MarginAccountHandle([]byte) error
	MarginBalanceHandle([]byte) error
	OrderHandle([]byte) error
	GetChan(chName string) interface{}
	SetOrderMarketConf([]*client.WsAccountReq)
}

type WebSocketSpotHandle struct {
	base.IncrementDepthConf
	orderConf []*client.WsAccountReq
	uapiCient *future_api.UApiClient

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

	Lock            sync.Mutex
	CheckSendStatus *base.CheckDataSendStatus
}

func NewWebSocketSpotHandle(chanCap int64, conf base.APIConf) *WebSocketSpotHandle {
	d := &WebSocketSpotHandle{}
	d.DepthCacheMap = sdk.NewCmapI()
	d.DepthCacheListMap = sdk.NewCmapI()
	d.CheckDepthCacheChanMap = sdk.NewCmapI()
	d.accountChan = make(chan *client.WsAccountRsp, chanCap)
	d.balanceChan = make(chan *client.WsBalanceRsp, chanCap)
	d.orderChan = make(chan *client.WsOrderRsp, chanCap)
	d.CheckSendStatus = base.NewCheckDataSendStatus()
	d.uapiCient = future_api.NewUApiClient(conf)

	return d
}

func (w *WebSocketSpotHandle) SetOrderMarketConf(reqs []*client.WsAccountReq) {
	w.orderConf = append(w.orderConf, reqs...)
}

// TradeGroupHandle Done
func (w *WebSocketSpotHandle) TradeGroupHandle(bytes []byte) error {
	var (
		resp TradeResponse
		t    = time.Now().UnixMicro()
		err  error
	)

	if err = json.Unmarshal(bytes, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(bytes))
		return err
	}
	if resp.Type == "error" {
		logger.Logger.Error("data error with code:", resp.Code, resp.Message)
		return nil
	} else if resp.Type == "pong" {
		return nil
	}
	for _, tradeItem := range resp.Data {
		res := &client.WsTradeRsp{
			OrderId:      strconv.FormatInt(tradeItem.Id, 10),
			Symbol:       resp.Market,
			ExchangeTime: tradeItem.Time.UnixMicro(),
			Price:        tradeItem.Price,
			Amount:       tradeItem.Size,
			TakerSide:    GetSide(tradeItem.Side),
			ReceiveTime:  t,
		}
		SymbolInfoStr := base.SymInfoToString(getSymbolInfo(res.Symbol))
		if _, ok := w.tradeGroupChanMap[SymbolInfoStr]; ok {
			base.SendChan(w.tradeGroupChanMap[SymbolInfoStr], res, "trade")
		} else {
			logger.Logger.Warn("[trade] get symbol from channel map err:", res.Symbol)
		}
	}
	w.CheckSendStatus.CheckUpdateTimeMap.Set(resp.Market, time.Now().UnixMicro())
	return nil
}

// SetTradeGroupChannel Done
func (w *WebSocketSpotHandle) SetTradeGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) {
	if chMap != nil {
		w.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp)
		for info, ch := range chMap {
			w.tradeGroupChanMap[base.SymInfoToString(info)] = ch
		}
	}
}

// AccountHandle TODO
func (w *WebSocketSpotHandle) AccountHandle(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}

// BalanceHandle TODO
func (w *WebSocketSpotHandle) BalanceHandle(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}

// MarginAccountHandle TODO
func (w *WebSocketSpotHandle) MarginAccountHandle(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}

// MarginBalanceHandle TODO
func (w *WebSocketSpotHandle) MarginBalanceHandle(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}

// Done
func (w *WebSocketSpotHandle) OrderHandle(data []byte) error {
	var (
		resp      OrdersResponse
		validData = false
		err       error
	)

	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}

	if resp.Type == "subscribed" {
		fmt.Println("Subscribed to FTX Order Channel")
		return nil
	} else if resp.Type == "pong" {
		return nil
	}

	for _, req := range w.orderConf {
		// no market set means any market is fine
		if req.Market == common.Market_INVALID_MARKET || req.Market == getMarket(resp.Data.Market) {
			validData = true
		}
	}

	if !validData {
		return nil
	}

	producer, id := transform.ClientIdToId(resp.Data.ClientId)
	res := &client.WsOrderRsp{
		Producer:     producer,
		Id:           id,
		IdEx:         resp.Data.ClientId,
		Status:       spot_api.GetOrderStatusFromExchange(resp.Data.Status, resp.Data.FilledSize, resp.Data.Size),
		TimeFilled:   resp.Data.CreatedAt.UnixMilli(),
		Symbol:       resp.Data.Market,
		Market:       getMarket(resp.Data.Market),
		Type:         getSymbolType(resp.Data.Market),
		AmountFilled: resp.Data.FilledSize,
		PriceFilled:  resp.Data.AvgFillPrice,
		QtyFilled:    resp.Data.FilledSize * resp.Data.AvgFillPrice,
	}

	base.SendChan(w.orderChan, res, "order")
	return nil
}

// GetChan TODO
func (w *WebSocketSpotHandle) GetChan(chName string) interface{} {
	switch chName {
	case "balance":
		return w.balanceChan
	case "order":
		return w.orderChan
	case "send_status":
		return w.CheckSendStatus
	default:
		return nil
	}
}

func (w *WebSocketSpotHandle) FundingRateGroupHandle(bytes []byte) error {
	timeInterval := 10
	for {
		now := time.Now().In(loc)
		start := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, loc)
		end := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 59, 59, 0, loc)
		rates, err := w.uapiCient.GetFundingRates(start.Unix(), end.Unix())
		if err != nil {
			logger.Logger.Error("get funding rates error: ", err)
			continue
		}
		if len(rates.Result) == 0 {
			logger.Logger.Warn("get funding rates is empty")
			continue
		}
		logger.Logger.Info("get funding rate success, size=", len(rates.Result))
		for _, rateInfo := range rates.Result {
			sym := getSymbolInfo(rateInfo.Future)
			res := &client.WsFundingRateRsp{
				Symbol:       sym.Symbol,
				Type:         sym.Type,
				FundingRate:  rateInfo.Rate,
				ReceiveTime:  now.UnixMicro(),
				ExchangeTime: rateInfo.Time.UnixMicro(),
			}
			SymbolInfoStr := base.SymInfoToString(sym)
			if _, ok := w.fundingRateGroupChanMap[SymbolInfoStr]; ok {
				base.SendChan(w.fundingRateGroupChanMap[SymbolInfoStr], res, "FundingRateGroup")
			}
		}
		time.Sleep(time.Duration(timeInterval) * time.Second)
	}
}

// AggTradeGroupHandle TODO
func (w *WebSocketSpotHandle) AggTradeGroupHandle(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}

// DepthIncrementGroupHandle To Verify
func (w *WebSocketSpotHandle) DepthIncrementGroupHandle(data []byte) error {
	var (
		resp OrderbooksResponse
		t    = time.Now().UnixMicro()
		err  error
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if resp.Type == "error" {
		logger.Logger.Error("data error with code:", resp.Code, resp.Message)
		return nil
	} else if resp.Type == "partial" {
		return nil
	} else if resp.Type == "pong" {
		return nil
	}
	asks, err := DepthLevelParse(resp.Data.Asks)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	bids, err := DepthLevelParse(resp.Data.Bids)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	res := &client.WsDepthRsp{
		ExchangeTime: int64(resp.Data.Time * float64(time.Microsecond)),
		ReceiveTime:  t,
		Symbol:       resp.Market,
		Asks:         asks,
		Bids:         bids,
	}

	SymbolInfoStr := base.SymInfoToString(getSymbolInfo(res.Symbol))
	if _, ok := w.depthIncrementGroupChanMap[SymbolInfoStr]; ok {
		base.SendChan(w.depthIncrementGroupChanMap[SymbolInfoStr], res, "depthIncrement")
	} else {
		logger.Logger.Warn("[depth] get symbol from channel map err:", res.Symbol)
	}
	w.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}

// DepthLimitGroupHandle TO Verify
func (w *WebSocketSpotHandle) DepthLimitGroupHandle(data []byte) error {
	var (
		resp OrderbooksResponse
		t    = time.Now().UnixMicro()
		err  error
	)

	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if resp.Type == "error" {
		logger.Logger.Error("data error with code:", resp.Code, resp.Message)
		return nil
	} else if resp.Type == "pong" {
		return nil
	}
	asks, err := DepthLevelParse(resp.Data.Asks)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	bids, err := DepthLevelParse(resp.Data.Bids)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	res := &depth.Depth{
		Exchange:     common.Exchange_FTX,
		Symbol:       resp.Market,
		TimeExchange: uint64(resp.Data.Time * float64(time.Microsecond)),
		TimeReceive:  uint64(t),
		TimeOperate:  uint64(time.Now().UnixMicro()),
		Asks:         asks,
		Bids:         bids,
	}

	SymbolInfoStr := base.SymInfoToString(getSymbolInfo(res.Symbol))
	if _, ok := w.depthLimitGroupChanMap[SymbolInfoStr]; ok {
		base.SendChan(w.depthLimitGroupChanMap[SymbolInfoStr], res, "depthLimitGroup")
	} else {
		logger.Logger.Warn("[depth_limit] get symbol from channel map err:", res.Symbol)
	}
	w.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}

func (w *WebSocketSpotHandle) SetFundingRateGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsFundingRateRsp) {
	if chMap != nil {
		w.fundingRateGroupChanMap = make(map[string]chan *client.WsFundingRateRsp)
		for info, ch := range chMap {
			w.fundingRateGroupChanMap[base.SymInfoToString(info)] = ch
		}
	}
}

// SetDepthLimitGroupChannel DONE
func (w *WebSocketSpotHandle) SetDepthLimitGroupChannel(chMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chMap != nil {
		w.depthLimitGroupChanMap = make(map[string]chan *depth.Depth)
		for info, ch := range chMap {
			w.depthLimitGroupChanMap[base.SymInfoToString(info)] = ch
		}
	}
}

// SetDepthIncrementGroupChannel DONE
func (w *WebSocketSpotHandle) SetDepthIncrementGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) {
	if chMap != nil {
		w.depthIncrementGroupChanMap = make(map[string]chan *client.WsDepthRsp)
		for info, ch := range chMap {
			w.depthIncrementGroupChanMap[base.SymInfoToString(info)] = ch
		}
	}
}

// SetDepthIncrementSnapshotGroupChannel Done
func (w *WebSocketSpotHandle) SetDepthIncrementSnapshotGroupChannel(chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chDeltaMap != nil {
		w.depthIncrementSnapshotDeltaGroupChanMap = make(map[string]chan *client.WsDepthRsp)
		for info, ch := range chDeltaMap {
			w.depthIncrementSnapshotDeltaGroupChanMap[base.SymInfoToString(info)] = ch
		}
	}
	if chFullMap != nil {
		w.depthIncrementSnapshotFullGroupChanMap = make(map[string]chan *depth.Depth)
		for info, ch := range chFullMap {
			w.depthIncrementSnapshotFullGroupChanMap[base.SymInfoToString(info)] = ch
		}
	}
}

// SetAggTradeGroupChannel TODO
func (w *WebSocketSpotHandle) SetAggTradeGroupChannel(m map[*client.SymbolInfo]chan *client.WsAggTradeRsp) {
	//TODO implement me
	panic("implement me")
}

// SetBookTickerGroupChannel DONE
func (w *WebSocketSpotHandle) SetBookTickerGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) {
	if chMap != nil {
		w.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp)
		for info, ch := range chMap {
			w.bookTickerGroupChanMap[base.SymInfoToString(info)] = ch
		}
	}
}

// BookTickerGroupHandle DONE
func (w *WebSocketSpotHandle) BookTickerGroupHandle(data []byte) error {
	var (
		resp TickerResponse
		t    = time.Now().UnixMicro()
		err  error
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if resp.Type == "subscribed" {
		//Initialization ack message
		return nil
	} else if resp.Type == "error" {
		logger.Logger.Error("data error with code:", resp.Code, resp.Message)
		return nil
	} else if resp.Type == "pong" {
		return nil
	}

	asks, err := DepthLevelParse([][]float64{{resp.Data.Ask, resp.Data.AskSize}})
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	bids, err := DepthLevelParse([][]float64{{resp.Data.Bid, resp.Data.BidSize}})
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}

	res := &client.WsBookTickerRsp{
		ExchangeTime: int64(resp.Data.Time * 1000 * 1000),
		ReceiveTime:  t,
		Symbol:       resp.Market,
		Ask:          asks[0],
		Bid:          bids[0],
	}
	SymbolInfoStr := base.SymInfoToString(getSymbolInfo(res.Symbol))
	if _, ok := w.bookTickerGroupChanMap[SymbolInfoStr]; ok {
		base.SendChan(w.bookTickerGroupChanMap[SymbolInfoStr], res, "bookTickerGroup")
	} else {
		logger.Logger.Warn("[ticker] get symbol from channel map err:", res.Symbol)
	}
	w.CheckSendStatus.CheckUpdateTimeMap.Set(resp.Market, time.Now().UnixMicro())
	return nil
}

// DepthLevelParse DONE
func DepthLevelParse(levelList [][]float64) ([]*depth.DepthLevel, error) {
	var (
		res []*depth.DepthLevel
		err error
	)
	for _, level := range levelList {
		res = append(res, &depth.DepthLevel{
			Price:  level[0],
			Amount: level[1],
		})
	}
	return res, err
}

func GetSide(side string) order.TradeSide {
	switch side {
	case "sell":
		return order.TradeSide_SELL
	case "buy":
		return order.TradeSide_BUY
	default:
		return order.TradeSide_InvalidSide
	}
}

func GetChecksum(orderbook *base.OrderBook) uint32 {
	var fields []string

	la, lb := len(orderbook.Asks), len(orderbook.Bids)
	for i := 0; i < 100; i++ {
		if i < lb {
			fields = append(fields, AppendDecimal(orderbook.Bids[i].Price)+":"+AppendDecimal(orderbook.Bids[i].Amount))
		}
		if i < la {
			fields = append(fields, AppendDecimal(orderbook.Asks[i].Price)+":"+AppendDecimal(orderbook.Asks[i].Amount))
		}
	}

	raw := strings.Join(fields, ":")
	cs := crc32.ChecksumIEEE([]byte(raw))
	return cs
}

func AppendDecimal(value float64) string {
	//Case 1: Small number
	if value < 0.0001 { //Really small values need to reformat scientific notation
		return strconv.FormatFloat(value, 'e', -1, 64)
	}

	//Case 2: Not a small number
	res := fmt.Sprintf("%v", value)
	if strings.Contains(res, "e") { //Unpack large values formatted in scientific notation
		res = strconv.FormatFloat(value, 'f', -1, 64)
	}
	//Return ok formatted number
	if strings.Contains(res, ".") {
		return res
	} else { //Append .0 if not formatted correctly
		res += ".0"
		return res
	}
}
