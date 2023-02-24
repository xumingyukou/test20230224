package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"fmt"
	"github.com/goccy/go-json"
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
	//AggTradeGroupHandle([]byte) error
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

	Tester([]byte) error
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

	firstFullSentMap sync.Map
	symbolMap        map[string]*client.SymbolInfo
	Lock             sync.Mutex
	CheckSendStatus  *base.CheckDataSendStatus
}

func NewWebSocketSpotHandle(chanCap int64) *WebSocketSpotHandle {
	c := &WebSocketSpotHandle{}

	c.accountChan = make(chan *client.WsAccountRsp, chanCap)
	c.balanceChan = make(chan *client.WsBalanceRsp, chanCap)
	c.orderChan = make(chan *client.WsOrderRsp, chanCap)
	c.CheckSendStatus = base.NewCheckDataSendStatus()

	return c
}

/*Public Interface Methods*/
func (c *WebSocketSpotHandle) FundingRateGroupHandle(data []byte) error {
	//TODO implement me
	panic("implement me")
}
func (c *WebSocketSpotHandle) TradeGroupHandle(data []byte) error {
	res, err := c.parseTrade(data)
	if err != nil {
		return err
	}
	if res == nil {
		return nil
	}

	if _, ok := c.tradeGroupChanMap[res.Symbol]; ok {
		base.SendChan(c.tradeGroupChanMap[res.Symbol], res, "trade")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	c.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}
func (c *WebSocketSpotHandle) BookTickerGroupHandle(data []byte) error {
	var (
		resp RespBookTickerStream
		t    = time.Now().UnixMicro()
		err  error
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}

	if resp.Type == "subscriptions" {
		return nil
	} else if resp.Type == "heartbeat" {
		//fmt.Println("Pong")
		return nil
	} else if resp.Type == "error" {
		fmt.Println("Error:", resp.Message, resp.Reason)
		return err
	}

	askPrice, err := strconv.ParseFloat(resp.BestAsk, 64)
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	bidPrice, err := strconv.ParseFloat(resp.BestBid, 64)
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	res := &client.WsBookTickerRsp{
		ExchangeTime: resp.Time.UnixMicro(),
		ReceiveTime:  t,
		Symbol:       DecodeSymbols(resp.ProductId),
		Ask: &depth.DepthLevel{
			Price:  askPrice,
			Amount: 0, //Coinbase doesn't provide size
		},
		Bid: &depth.DepthLevel{
			Price:  bidPrice,
			Amount: 0, //Coinbase doesn't provide size
		},
	}
	if _, ok := c.bookTickerGroupChanMap[res.Symbol]; ok {
		base.SendChan(c.bookTickerGroupChanMap[res.Symbol], res, "bookTickerGroup")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	c.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}
func (c *WebSocketSpotHandle) DepthLimitGroupHandle(data []byte) error {
	var (
		resp          RespLimitDepthStream
		asks          []*depth.DepthLevel
		bids          []*depth.DepthLevel
		amount, price float64
		t             = time.Now().UnixMicro()
		err           error
	)

	if err = json.Unmarshal(data, &resp); err != nil {
		fmt.Println("Err", err)
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if resp.Type == "subscriptions" {
		return nil
	} else if resp.Type == "heartbeat" {
		return nil
	} else if resp.Type == "error" {
		fmt.Println("Error:", resp.Message, resp.Reason)
		return err
	}
	for _, initialSnapshot := range resp.Asks {
		if price, err = strconv.ParseFloat(initialSnapshot[0], 64); err != nil {
			return err
		}
		if amount, err = strconv.ParseFloat(initialSnapshot[1], 64); err != nil {
			return err
		}
		asks = append(asks, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}

	for _, initialSnapshot := range resp.Bids {
		if price, err = strconv.ParseFloat(initialSnapshot[0], 64); err != nil {
			return err
		}
		if amount, err = strconv.ParseFloat(initialSnapshot[1], 64); err != nil {
			return err
		}
		bids = append(bids, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}

	res := &depth.Depth{
		Exchange:     common.Exchange_COINBASE, //Coinbase is currently missing from Exchange database
		Symbol:       DecodeSymbols(resp.ProductId),
		TimeExchange: uint64(resp.Time.UnixMicro()),
		TimeReceive:  uint64(t),
		TimeOperate:  uint64(time.Now().UnixMicro()),
		Asks:         asks,
		Bids:         bids,
	}

	if _, ok := c.depthLimitGroupChanMap[res.Symbol]; ok {
		base.SendChan(c.depthLimitGroupChanMap[res.Symbol], res, "depthLimitGroup")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	c.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}
func (c *WebSocketSpotHandle) DepthIncrementGroupHandle(data []byte) error {
	var (
		resp          RespIncrementDepthStream
		asks          []*depth.DepthLevel
		bids          []*depth.DepthLevel
		t             = time.Now().UnixMicro()
		amount, price float64
		err           error
	)

	if err = json.Unmarshal(data, &resp); err != nil {
		fmt.Println("Err", err)
		logger.Logger.Error("receive data err:", string(data))
		return err
	}

	if resp.Type == "subscriptions" {
		return nil
	} else if resp.Type == "heartbeat" {
		return nil
	} else if resp.Type == "error" {
		fmt.Println("Error:", resp.Message, resp.Reason)
		return err
	}

	for _, incrementInfo := range resp.Changes {
		price, err = strconv.ParseFloat(incrementInfo.Price, 64)
		amount, err = strconv.ParseFloat(incrementInfo.Size, 64)
		switch incrementInfo.Side {
		case BUY:
			bids = append(bids, &depth.DepthLevel{
				Price:  price,
				Amount: amount,
			})
		case SELL:
			asks = append(asks, &depth.DepthLevel{
				Price:  price,
				Amount: amount,
			})
		}
	}

	res := &client.WsDepthRsp{
		ExchangeTime: resp.Time.UnixMicro(),
		ReceiveTime:  t,
		Symbol:       DecodeSymbols(resp.ProductId),
		Asks:         asks,
		Bids:         bids,
	}

	if _, ok := c.depthIncrementGroupChanMap[res.Symbol]; ok {
		base.SendChan(c.depthIncrementGroupChanMap[res.Symbol], res, "depthIncrement")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	c.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}

func (c *WebSocketSpotHandle) SetFundingRateGroupChannel(m map[*client.SymbolInfo]chan *client.WsFundingRateRsp) {
	//TODO implement me
	panic("implement me")
}
func (c *WebSocketSpotHandle) SetBookTickerGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) {
	if chMap != nil {
		c.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp)
		for info, ch := range chMap {
			c.bookTickerGroupChanMap[info.Symbol] = ch
		}
	}
}
func (c *WebSocketSpotHandle) SetDepthLimitGroupChannel(chMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chMap != nil {
		c.depthLimitGroupChanMap = make(map[string]chan *depth.Depth)
		for info, ch := range chMap {
			c.depthLimitGroupChanMap[info.Symbol] = ch
		}
	}
}
func (c *WebSocketSpotHandle) SetDepthIncrementGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) {
	if chMap != nil {
		c.depthIncrementGroupChanMap = make(map[string]chan *client.WsDepthRsp, 50)
		for info, ch := range chMap {
			c.depthIncrementGroupChanMap[info.Symbol] = ch
		}
	}
}
func (c *WebSocketSpotHandle) SetDepthIncrementSnapshotGroupChannel(chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chDeltaMap != nil {
		c.depthIncrementSnapshotDeltaGroupChanMap = make(map[string]chan *client.WsDepthRsp, 50)
		for info, ch := range chDeltaMap {
			sym := info.Symbol
			c.depthIncrementSnapshotDeltaGroupChanMap[sym] = ch
		}
	}
	if chFullMap != nil {
		c.depthIncrementSnapshotFullGroupChanMap = make(map[string]chan *depth.Depth, 50)
		for info, ch := range chFullMap {
			sym := info.Symbol
			c.depthIncrementSnapshotFullGroupChanMap[sym] = ch
		}
	}
}
func (c *WebSocketSpotHandle) SetAggTradeGroupChannel(m map[*client.SymbolInfo]chan *client.WsAggTradeRsp) {
	//TODO implement me
	panic("implement me")
}

func (c *WebSocketSpotHandle) AccountHandle(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}
func (c *WebSocketSpotHandle) BalanceHandle(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}
func (c *WebSocketSpotHandle) MarginAccountHandle(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}
func (c *WebSocketSpotHandle) MarginBalanceHandle(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}
func (c *WebSocketSpotHandle) OrderHandle(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}

/*Old Functions*/
func (c *WebSocketSpotHandle) DepthLimitHandle(data []byte) error {
	var (
		resp RespLimitDepthStream
		err  error
	)

	if err = json.Unmarshal(data, &resp); err != nil {
		fmt.Println("Err", err)
		logger.Logger.Error("receive data err:", string(data))
		return err
	}

	//asks, err := c.DepthLevelParse(resp.Asks)
	//if err != nil {
	//	logger.Logger.Error("depth level parse err", err, string(data))
	//	return err
	//}
	//bids, err := c.DepthLevelParse(resp.Bids)
	//if err != nil {
	//	logger.Logger.Error("depth level parse err", err, string(data))
	//	return err
	//}
	//
	//res := client.WsDepthRsp{
	//	UpdateIdEnd: resp.Time.UnixMilli(),
	//	Asks:        asks,
	//	Bids:        bids,
	//	ReceiveTime: time.Now().UnixMicro(),
	//}
	//if len(c.depthLimitChan) > 0 {
	//	content := fmt.Sprint("channel slow:", len(c.depthLimitChan))
	//	logger.Logger.Warn(content)
	//}
	//c.depthLimitChan <- &res
	return nil
}

/*Handler Helper Functions*/
func (c *WebSocketSpotHandle) parseTrade(data []byte) (*client.WsTradeRsp, error) {
	var (
		resp                        RespTradeStream
		err                         error
		price, amount               float64
		buyerOrderId, sellerOrderId string
		t                           = time.Now().UnixMicro()
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return nil, err
	}
	if resp.Type == "subscriptions" {
		return nil, err
	} else if resp.Type == "heartbeat" {
		//fmt.Println("Pong")
		return nil, err
	} else if resp.Type == "error" {
		fmt.Println("Error:", resp.Message, resp.Reason)
		return nil, err
	}

	price, err = strconv.ParseFloat(resp.Price, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return nil, err
	}
	amount, err = strconv.ParseFloat(resp.Size, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return nil, err
	}

	switch resp.Side {
	case BUY:
		buyerOrderId = resp.MakerOrderId
		sellerOrderId = resp.TakerOrderId
	case SELL:
		buyerOrderId = resp.TakerOrderId
		sellerOrderId = resp.MakerOrderId
	}

	res := &client.WsTradeRsp{
		OrderId:       strconv.Itoa(resp.TradeId),
		Symbol:        DecodeSymbols(resp.ProductId),
		ExchangeTime:  resp.Time.UnixMicro(),
		ReceiveTime:   t,
		Price:         price,
		Amount:        amount,
		TakerSide:     GetSide(resp.Side),
		BuyerOrderId:  buyerOrderId,
		SellerOrderId: sellerOrderId,
		DealTime:      resp.Time.UnixMicro(),
	}
	return res, nil
}
func (c *WebSocketSpotHandle) DepthLevelParse(levelList [][]string) ([]*depth.DepthLevel, error) {
	var (
		res           []*depth.DepthLevel
		amount, price float64
		err           error
	)
	for _, level := range levelList {
		if price, err = strconv.ParseFloat(level[0], 64); err != nil {
			return nil, err
		}
		if amount, err = strconv.ParseFloat(level[1], 64); err != nil {
			return nil, err
		}
		res = append(res, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	return res, err
}
func (c *WebSocketSpotHandle) SetTradeGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) {
	if chMap != nil {
		c.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp)
		for info, ch := range chMap {
			c.tradeGroupChanMap[info.Symbol] = ch
		}
	}
}
func (c *WebSocketSpotHandle) GetChan(chName string) interface{} {
	switch chName {
	case "send_status":
		return c.CheckSendStatus
	//case "fundingRate":
	//	return b.fundingRateChan
	//case "aggTrade":
	//	return b.aggTradeChan
	//case "trade":
	//	return c.tradeChan
	//case "depthIncrement":
	//	return c.depthIncrementChan
	//case "depthLimit":
	//	return c.depthLimitChan
	//case "bookTicker":
	//	return c.bookTickerChan
	//case "account":
	//	return b.accountChan
	//case "balance":
	//	return b.balanceChan
	//case "order":
	//	return b.orderChan
	default:
		return nil
	}
}

func (c *WebSocketSpotHandle) Tester(data []byte) error {
	fmt.Println(string(data))
	return nil
}

/*Other Helper Functions*/
func DecodeSymbols(input string) string {
	return strings.Replace(input, "-", "/", -1)
}
func GetSide(side Side) order.TradeSide {
	switch side {
	case BUY:
		return order.TradeSide_BUY
	case SELL:
		return order.TradeSide_SELL
	default:
		return order.TradeSide_InvalidSide
	}
}
func DepthLevelParse(levelList [][]string) ([]*depth.DepthLevel, error) {
	var (
		res           []*depth.DepthLevel
		amount, price float64
		err           error
	)
	for _, level := range levelList {
		if price, err = strconv.ParseFloat(level[0], 64); err != nil {
			return nil, err
		}
		if amount, err = strconv.ParseFloat(level[1], 64); err != nil {
			return nil, err
		}
		res = append(res, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	return res, err
}
