package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"clients/transform"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
)

type WebSocketHandleInterface interface {
	FundingRateGroupHandle([]byte) error
	TradeGroupHandle([]byte) error
	BookTickerGroupHandle([]byte) error
	DepthLimitGroupHandle([]byte) error
	DepthIncrementGroupHandle([]byte) error
	DepthIncrementSnapShotGroupHandle([]byte) error

	SetFundingRateGroupChannel(map[*client.SymbolInfo]chan *client.WsFundingRateRsp)
	SetTradeGroupChannel(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	SetBookTickerGroupChannel(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	SetDepthLimitGroupChannel(map[*client.SymbolInfo]chan *depth.Depth)
	SetDepthIncrementGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	SetDepthIncrementSnapshotGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp, map[*client.SymbolInfo]chan *depth.Depth)

	SetDepthIncrementSnapShotConf([]*client.SymbolInfo, *base.IncrementDepthConf)

	GetChan(chName string) interface{}
}

type WebSocketSpotHandle struct {
	*base.IncrementDepthConf

	pingpong                                chan int64
	contractpingpong                        chan int64
	tradeGroupChanMap                       map[string]chan *client.WsTradeRsp
	bookTickerGroupChanMap                  map[string]chan *client.WsBookTickerRsp
	depthLimitGroupChanMap                  map[string]chan *depth.Depth       //全量
	depthIncrementGroupChanMap              map[string]chan *client.WsDepthRsp //增量
	depthIncrementSnapshotDeltaGroupChanMap map[string]chan *client.WsDepthRsp //增量合并数据
	depthIncrementSnapshotFullGroupChanMap  map[string]chan *depth.Depth       //增量合并数据
	fundingRateGroupChanMap                 map[string]chan *client.WsFundingRateRsp

	symbolMap                       map[string]*client.SymbolInfo
	DepthIncrementSnapshotReqSymbol chan *client.SymbolInfo
	Lock                            sync.Mutex
	CheckSendStatus                 *base.CheckDataSendStatus
}

func NewWebSocketSpotHandle(chanCap int64) *WebSocketSpotHandle {
	d := &WebSocketSpotHandle{
		IncrementDepthConf: &base.IncrementDepthConf{},
	}
	d.pingpong = make(chan int64, chanCap)
	d.contractpingpong = make(chan int64, chanCap)
	d.CheckSendStatus = base.NewCheckDataSendStatus()
	return d
}

func (c *WebSocketSpotHandle) TradeGroupHandle(data []byte) error {
	t := time.Now().UnixMicro()
	var (
		resp      TradeResponse
		err       error
		takerSide order.TradeSide
	)
	//var str bytes.Buffer
	////_ = json.Indent(&str, []byte(data), "", "    ")
	////fmt.Println("formated: ", str.String())
	////fmt.Println("data: ", data)
	if err = json.Unmarshal(data, &resp); err != nil {
		// fmt.Println("Err", err)
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	// fmt.Println(resp)
	if resp.Ping != 0 {
		//fmt.Println(resp)
		c.pingpong <- resp.Ping
		return nil
	}
	if resp.Ch == "" {
		//fmt.Println(resp)
		return nil
	}
	for _, tradeItem := range resp.Tick.Data {
		//fmt.Println("TradeItem:", tradeItem)
		if tradeItem.Direction == "buy" {
			takerSide = 1
		} else {
			takerSide = 2
		}
		chName := strings.Split(resp.Ch, ".")[1]
		var symbolName string
		symIf, ok := symbolNameMap.Load(chName)
		if ok {
			symbolName = symIf.(string)
		}
		res := &client.WsTradeRsp{
			OrderId:      strconv.FormatInt(tradeItem.TradeID, 10),
			Symbol:       symbolName,
			ExchangeTime: resp.Ts * 1000,
			ReceiveTime:  t,
			Price:        tradeItem.Price,
			Amount:       tradeItem.Amount,
			TakerSide:    takerSide,
		}
		//fmt.Println(res.Symbol)
		// fmt.Println(c.tradeGroupChanMap)
		if _, ok := c.tradeGroupChanMap[res.Symbol]; ok {
			// 张数转换
			TransContractSize(chName, c.Market, res)
			base.SendChan(c.tradeGroupChanMap[res.Symbol], res, "trade")
		} else {
			logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
		}
		c.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	}
	return nil
}

func (c *WebSocketSpotHandle) BookTickerGroupHandle(data []byte) error {
	t := time.Now().UnixMicro()
	var (
		resp TickerResponse
		err  error
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		fmt.Println("Err", err)
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if resp.Ping != 0 {
		// fmt.Println(resp)
		c.pingpong <- resp.Ping
		return nil
	}
	if resp.Ch == "" {
		//fmt.Println(resp)
		return nil
	}
	asks, err := DepthLevelParse([][]float64{{resp.Tick.Ask, resp.Tick.AskSize}})
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	bids, err := DepthLevelParse([][]float64{{resp.Tick.Bid, resp.Tick.BidSize}})
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}

	var symbolName string
	symIf, ok := symbolNameMap.Load(strings.Split(resp.Ch, ".")[1])
	if ok {
		symbolName = symIf.(string)
	}

	ask, bid := &depth.DepthLevel{}, &depth.DepthLevel{}
	if len(asks) > 0 {
		ask = asks[0]
	}
	if len(bids) > 0 {
		bid = bids[0]
	}
	res := &client.WsBookTickerRsp{
		ExchangeTime: resp.Ts * 1000,
		ReceiveTime:  t,
		Symbol:       symbolName,
		Ask:          ask,
		Bid:          bid,
	}
	// fmt.Println("res: ", res)
	if _, ok := c.bookTickerGroupChanMap[res.Symbol]; ok {
		base.SendChan(c.bookTickerGroupChanMap[res.Symbol], res, "bookTickerGroup")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	c.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}

func (c *WebSocketSpotHandle) FutureOrSwapBookTickerGroupHandle(data []byte) error {
	t := time.Now().UnixMicro()
	var (
		resp FutureOrSwapTickerResp
		err  error
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		fmt.Println("Err", err)
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if resp.Ping != 0 {
		// fmt.Println(resp)
		c.pingpong <- resp.Ping
		return nil
	}
	if resp.Ch == "" {
		//fmt.Println(resp)
		return nil
	}

	asks, err := DepthLevelParse([][]float64{resp.Tick.Ask})
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	bids, err := DepthLevelParse([][]float64{resp.Tick.Bid})
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}

	var symbolName string
	chName := strings.Split(resp.Ch, ".")[1]
	symIf, ok := symbolNameMap.Load(chName)
	if ok {
		symbolName = symIf.(string)
	}

	ask, bid := &depth.DepthLevel{}, &depth.DepthLevel{}
	if len(asks) > 0 {
		ask = asks[0]
	}
	if len(bids) > 0 {
		bid = bids[0]
	}
	res := &client.WsBookTickerRsp{
		ExchangeTime: resp.Ts * 1000,
		ReceiveTime:  t,
		Symbol:       symbolName,
		Ask:          ask,
		Bid:          bid,
	}
	// fmt.Println("res: ", res)
	if _, ok := c.bookTickerGroupChanMap[res.Symbol]; ok {
		// 张数转换
		TransContractSize(chName, c.Market, nil, res.Ask)
		TransContractSize(chName, c.Market, nil, res.Bid)
		base.SendChan(c.bookTickerGroupChanMap[res.Symbol], res, "bookTickerGroup")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	c.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}

func (c *WebSocketSpotHandle) DepthLimitGroupHandle(data []byte) error {
	t := time.Now().UnixMicro()
	var (
		resp RespLimitDepthStream
		err  error
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		fmt.Println("Err", err)
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if resp.Ping != 0 {
		c.pingpong <- resp.Ping
		return nil
	}
	if resp.Ch == "" {
		//fmt.Println(resp)
		return nil
	}
	asks, err := DepthLevelParse(resp.Tick.Asks)
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	bids, err := DepthLevelParse(resp.Tick.Bids)
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}

	var symbolName string
	symIf, ok := symbolNameMap.Load(strings.Split(resp.Ch, ".")[1])
	if ok {
		symbolName = symIf.(string)
	}
	res := &depth.Depth{
		Exchange:     common.Exchange_HUOBI,
		Symbol:       symbolName,
		TimeExchange: uint64(resp.Ts * 1000),
		TimeReceive:  uint64(t),
		TimeOperate:  uint64(time.Now().UnixMicro()),
		Asks:         asks,
		Bids:         bids,
	}
	// fmt.Println("res: ", res)
	if _, ok := c.depthLimitGroupChanMap[res.Symbol]; ok {
		base.SendChan(c.depthLimitGroupChanMap[res.Symbol], res, "depthLimitGroup")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	c.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}

func (c *WebSocketSpotHandle) DepthIncrementGroupHandle(data []byte) error {
	t := time.Now().UnixMicro()
	var (
		resp RespIncrementDepthStream
		err  error
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		fmt.Println("Err", err)
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	// fmt.Println("DepthIncrementGroupHandle: ", resp)
	if resp.Ping != 0 {
		// fmt.Println(resp)
		c.pingpong <- resp.Ping
		return nil
	}
	if resp.Ch == "" {
		//fmt.Println(resp)
		return nil
	}
	asks, err := DepthLevelParse(resp.Tick.Asks)
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	bids, err := DepthLevelParse(resp.Tick.Bids)
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}

	var symbolName string
	symIf, ok := symbolNameMap.Load(strings.Split(resp.Ch, ".")[1])
	if ok {
		symbolName = symIf.(string)
	}
	res := &client.WsDepthRsp{
		UpdateIdStart: resp.Tick.PrevSeqNum,
		UpdateIdEnd:   resp.Tick.SeqNum,
		ExchangeTime:  resp.Ts * 1000,
		ReceiveTime:   t,
		Symbol:        symbolName,
		Asks:          asks,
		Bids:          bids,
	}
	// fmt.Println("res: ", res)
	if _, ok := c.depthIncrementGroupChanMap[res.Symbol]; ok {
		base.SendChan(c.depthIncrementGroupChanMap[res.Symbol], res, "depthIncrement")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	c.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}

func (c *WebSocketSpotHandle) SetTradeGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) {
	if chMap != nil {
		c.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp)
		for info, ch := range chMap {
			sym := base.SymInfoToString(info)
			c.tradeGroupChanMap[sym] = ch
		}
	}
}

func (b *WebSocketSpotHandle) SetBookTickerGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) {
	if chMap != nil {
		b.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp)
		for info, ch := range chMap {
			sym := base.SymInfoToString(info)
			b.bookTickerGroupChanMap[sym] = ch
		}
	}
}

func (c *WebSocketSpotHandle) SetDepthLimitGroupChannel(chMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chMap != nil {
		c.depthLimitGroupChanMap = make(map[string]chan *depth.Depth)
		for info, ch := range chMap {
			sym := base.SymInfoToString(info)
			c.depthLimitGroupChanMap[sym] = ch
		}
	}
}

func (c *WebSocketSpotHandle) SetDepthIncrementGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) {
	if chMap != nil {
		c.depthIncrementGroupChanMap = make(map[string]chan *client.WsDepthRsp)
		for info, ch := range chMap {
			c.depthIncrementGroupChanMap[info.Symbol] = ch
		}
	}
}

func (c *WebSocketSpotHandle) SetDepthIncrementSnapshotGroupChannel(chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chDeltaMap != nil {
		c.depthIncrementSnapshotDeltaGroupChanMap = make(map[string]chan *client.WsDepthRsp)
		for info, ch := range chDeltaMap {
			sym := base.SymInfoToString(info)
			c.depthIncrementSnapshotDeltaGroupChanMap[sym] = ch
		}
	}
	if chFullMap != nil {
		c.depthIncrementSnapshotFullGroupChanMap = make(map[string]chan *depth.Depth)
		for info, ch := range chFullMap {
			sym := base.SymInfoToString(info)
			c.depthIncrementSnapshotFullGroupChanMap[sym] = ch
		}
	}
}

func DepthLevelParse(levelList [][]float64) ([]*depth.DepthLevel, error) {
	var (
		res []*depth.DepthLevel
		err error
	)
	for _, level := range levelList {
		if len(level) <= 0 {
			continue
		}
		res = append(res, &depth.DepthLevel{
			Price:  level[0],
			Amount: level[1],
		})
	}
	return res, err
}

func TransContractSize(chName string, market common.Market, tradeRsp *client.WsTradeRsp, levels ...*depth.DepthLevel) {
	var transFunc func(price, amount, contractSize float64) float64
	switch market {
	case common.Market_FUTURE_COIN, common.Market_SWAP_COIN:
		transFunc = func(price, amount, contractSize float64) float64 {
			if price == 0 {
				return amount
			}
			return amount * contractSize / price
		}
	case common.Market_FUTURE, common.Market_SWAP:
		transFunc = func(price, amount, contractSize float64) float64 {
			return amount * contractSize
		}
	case common.Market_SPOT:
		fallthrough
	default:
		return
	}
	contractSizeIf, ok := contractSizeMap.Load(chName)
	if !ok {
		logger.Logger.Error("contractSizeMap key error: ", chName)
		return
	}
	contractSize, okk := contractSizeIf.(float64)
	if !okk {
		logger.Logger.Error("contractSize assert error: ", chName)
		return
	}

	// depth和ticker数据
	for _, level := range levels {
		level.Amount = transFunc(level.Price, level.Amount, contractSize)
	}

	// trade数据
	if tradeRsp != nil {
		tradeRsp.Amount = transFunc(tradeRsp.Price, tradeRsp.Amount, contractSize)
	}
}

func (c *WebSocketSpotHandle) GetChan(chName string) interface{} {
	switch chName {
	case "pingpong":
		return c.pingpong
	case "send_status":
		return c.CheckSendStatus
	case "contractpingpong":
		return c.contractpingpong
	default:
		return nil
	}
}

func (b *WebSocketSpotHandle) FundingRateGroupHandle(data []byte) error {
	//fmt.Println("date", string(data))
	var (
		resp        RespFundRate
		err         error
		fundingRate float64
		type_       common.SymbolType
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if resp.Op == "ping" && resp.Ts != nil {
		b.contractpingpong <- transform.StringToX[int64](resp.Ts.(string)).(int64)
		return nil
	}
	if len(resp.Data) <= 0 {
		return nil
	}
	if strings.Contains(string(data), "error") {
		logger.Logger.Error(string(data))
		err = errors.New(string(data))
		return err
	}

	fundingRate, err = strconv.ParseFloat(resp.Data[0].FundingRate, 64)
	if b.Market == common.Market_SWAP {
		type_ = common.SymbolType_SWAP_FOREVER
	} else if b.Market == common.Market_SWAP_COIN {
		type_ = common.SymbolType_SWAP_COIN_FOREVER
	}
	var symbolName string
	symIf, ok := symbolNameMap.Load(resp.Data[0].ContractCode)
	if ok {
		symbolName = symIf.(string)
	}
	exchangeTime, err := strconv.ParseInt(resp.Data[0].FundingTime, 10, 64)
	if err != nil {
		logger.Logger.Error("funding rate exchangeTime convert error:", err)
		return err
	}
	res := &client.WsFundingRateRsp{
		Symbol:       symbolName,
		Type:         type_,
		FundingRate:  fundingRate,
		ReceiveTime:  time.Now().UnixMicro(),
		ExchangeTime: exchangeTime * 1000,
	}
	if _, ok := b.fundingRateGroupChanMap[symbolName]; ok {
		base.SendChan(b.fundingRateGroupChanMap[symbolName], res, "fundingRateGroup")
	} else {
		logger.Logger.Warn("FundingRateGroupHandle get symbol from channel map err:", symbolName)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(symbolName, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) SetFundingRateGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsFundingRateRsp) {
	if chMap != nil {
		b.fundingRateGroupChanMap = make(map[string]chan *client.WsFundingRateRsp)
		for info, ch := range chMap {
			sym := base.SymInfoToString(info)
			b.fundingRateGroupChanMap[sym] = ch
		}
	}
}
