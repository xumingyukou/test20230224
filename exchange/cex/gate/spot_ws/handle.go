package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/gate/gate_api"
	"clients/logger"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/warmplanet/proto/go/sdk"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
)

type WebSocketHandleInterface interface {
	TradeGroupHandle([]byte) error
	SetTradeGroupChannel(map[*client.SymbolInfo]chan *client.WsTradeRsp)

	DepthIncrementGroupHandle([]byte) error
	DepthIncrementSnapShotGroupHandle([]byte) error

	SetDepthIncrementSnapshotGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp, map[*client.SymbolInfo]chan *depth.Depth)
	SetDepthIncrementSnapShotConf([]*client.SymbolInfo, *base.IncrementDepthConf)
	SetDepthIncrementGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp)

	BookTickerGroupHandle([]byte) error
	SetBookTickerGroupChannel(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)

	GetChan(chName string) interface{}
}

type WebSocketSpotHandle struct {
	reSubscribe chan req

	symbolMap map[string]*client.SymbolInfo
	base.IncrementDepthConf
	subscribleInfoChan                      chan Resp_Info
	aggTradeChan                            chan *client.WsAggTradeRsp
	tradeChan                               chan *client.WsTradeRsp
	depthIncrementChan                      chan *client.WsDepthRsp
	depthLimitChan                          chan *client.WsDepthRsp
	bookTickerChan                          chan *client.WsBookTickerRsp
	accountChan                             chan *client.WsAccountRsp
	balanceChan                             chan *client.WsBalanceRsp
	orderChan                               chan *client.WsOrderRsp
	bookTickerGroupChanMap                  map[string]chan *client.WsBookTickerRsp
	tradeGroupChanMap                       map[string]chan *client.WsTradeRsp
	depthIncrementGroupChanMap              map[string]chan *client.WsDepthRsp //增量
	depthIncrementSnapshotDeltaGroupChanMap map[string]chan *client.WsDepthRsp //增量数据
	depthIncrementSnapshotFullGroupChanMap  map[string]chan *depth.Depth       //全量合并数据
	firstMap                                map[string]bool
	firstFullSentMap                        map[string]bool
	Lock                                    sync.Mutex // 锁
	CheckSendStatus                         *base.CheckDataSendStatus
}

func NewWebSocketSpotHandle(chanCap int64) *WebSocketSpotHandle {
	d := &WebSocketSpotHandle{}
	d.reSubscribe = make(chan req, chanCap)
	d.subscribleInfoChan = make(chan Resp_Info, chanCap)
	d.tradeChan = make(chan *client.WsTradeRsp, chanCap)
	d.depthIncrementChan = make(chan *client.WsDepthRsp, chanCap)
	d.depthLimitChan = make(chan *client.WsDepthRsp, chanCap)
	d.bookTickerChan = make(chan *client.WsBookTickerRsp, chanCap)
	d.accountChan = make(chan *client.WsAccountRsp, chanCap)
	d.balanceChan = make(chan *client.WsBalanceRsp, chanCap)
	d.orderChan = make(chan *client.WsOrderRsp, chanCap)
	d.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp)
	d.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp)
	d.firstMap = make(map[string]bool)
	d.CheckSendStatus = base.NewCheckDataSendStatus()
	d.firstFullSentMap = make(map[string]bool)
	return d
}

func (b *WebSocketSpotHandle) GetChan(chName string) interface{} {
	switch chName {
	case "aggTrade":
		return b.aggTradeChan
	case "trade":
		return b.tradeChan
	case "depthIncrement":
		return b.depthIncrementChan
	case "depthLimit":
		return b.depthLimitChan
	case "bookTicker":
		return b.bookTickerChan
	case "account":
		return b.accountChan
	case "balance":
		return b.balanceChan
	case "order":
		return b.orderChan
	case "subscribeInfo":
		return b.subscribleInfoChan
	case "reSubscribe":
		return b.reSubscribe
	case "send_status":
		return b.CheckSendStatus
	default:
		return nil
	}
}

func (b *WebSocketSpotHandle) SetDepthIncrementSnapShotConf(symbols []*client.SymbolInfo, conf *base.IncrementDepthConf) {
	if conf.DepthCapLevel <= 0 {
		conf.DepthCapLevel = 1000
	}
	if conf.DepthCapLevel <= 0 {
		conf.DepthCapLevel = 20
	}
	if conf.CheckTimeSec <= 0 {
		conf.CheckTimeSec = 3600
	}
	if conf.DepthCheckLevel <= 0 {
		conf.CheckTimeSec = 20
	}
	if conf.DepthCacheMap == nil {
		conf.DepthCacheMap = sdk.NewCmapI()
	}
	if conf.DepthCacheListMap == nil {
		conf.DepthCacheListMap = sdk.NewCmapI()
	}
	if conf.CheckDepthCacheChanMap == nil {
		conf.CheckDepthCacheChanMap = sdk.NewCmapI()
	}
	conf.DepthNotMatchChanMap = make(map[*client.SymbolInfo]chan bool)
	for _, symbol := range symbols {
		var (
			DepthCacheList      []*base.DeltaDepthUpdate
			CheckDepthCacheChan = make(chan *base.OrderBook, conf.DepthCapLevel)
		)
		conf.DepthCacheMap.Set(GetInstId(symbol), nil)
		conf.DepthCacheListMap.Set(GetInstId(symbol), DepthCacheList)
		conf.CheckDepthCacheChanMap.Set(GetInstId(symbol), CheckDepthCacheChan)
		conf.DepthNotMatchChanMap[symbol] = make(chan bool, conf.DepthCapLevel)
		//go b.Check()
	}
	b.IncrementDepthConf = base.IncrementDepthConf{
		IsPublishDelta:         conf.IsPublishDelta,
		IsPublishFull:          conf.IsPublishFull,
		DepthCacheMap:          conf.DepthCacheMap,
		DepthCacheListMap:      conf.DepthCacheListMap,
		CheckDepthCacheChanMap: conf.CheckDepthCacheChanMap,
		CheckTimeSec:           conf.CheckTimeSec,
		DepthCapLevel:          conf.DepthCapLevel,
		DepthLevel:             conf.DepthLevel,
		GetFullDepth:           conf.GetFullDepth,
		DepthNotMatchChanMap:   conf.DepthNotMatchChanMap,
		Ctx:                    conf.Ctx,
		CheckStates:            conf.CheckStates,
	}
}

func (b *WebSocketSpotHandle) DepthIncrementGroupHandle(data []byte) error {
	// 跳过前两个确认信息
	var (
		resp Resp_Depth
		err  error
	)
	if strings.Contains(string(data), "pong") {
		return nil
	}

	if err = b.HandleRespErr(data, &resp); err != nil {
		logger.Logger.Error(err, "错误返回为：", string(data))
		return err
	}
	if resp.Status != "" {
		if strings.Contains(resp.Result.Status, "success") {
			return nil
		} else {
			logger.Logger.Error("订阅错误：", string(data))
			return errors.New("订阅错误")
		}
	}
	asks := b.DepthLevelParse(resp.Result.Asks)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	bids := b.DepthLevelParse(resp.Result.Bids)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	res := client.WsDepthRsp{
		Symbol:       ParseSymbol(resp.Result.CurrencyPair),
		Asks:         asks,
		Bids:         bids,
		ExchangeTime: resp.Result.TimeInMilli * 1000,
		ReceiveTime:  time.Now().UnixMicro(),
	}

	if _, ok := b.depthIncrementGroupChanMap[res.Symbol]; ok {
		base.SendChan(b.depthIncrementGroupChanMap[res.Symbol], &res, "depthIncrement")
		//b.depthIncrementGroupChanMap[res.Symbol] <- &res
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol, string(data))
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())

	return nil
}

// []interface转[][]float64

func (b *WebSocketSpotHandle) DepthLevelParse(levelList [][]string) []*depth.DepthLevel {
	var (
		res []*depth.DepthLevel
	)
	for _, level := range levelList {
		price, err := strconv.ParseFloat(level[0], 64)
		if err != nil {
			logger.Logger.Error("转换错误", level[0])
		}
		amount, err := strconv.ParseFloat(level[1], 64)
		if err != nil {
			logger.Logger.Error("转换错误", level[1])
		}
		res = append(res, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	return res
}

func (b *WebSocketSpotHandle) SetDepthIncrementGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) {
	if chMap != nil {
		b.depthIncrementGroupChanMap = make(map[string]chan *client.WsDepthRsp, 1000)
		for k, v := range chMap {
			// 使用斜杠分开的
			b.depthIncrementGroupChanMap[TranInstId(k)] = v
		}
	}
}

func (b *WebSocketSpotHandle) SetTradeGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) {
	if chMap != nil {
		b.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp, 1000)
		for k, v := range chMap {
			b.tradeGroupChanMap[TranInstId(k)] = v
		}
	}
}

func (b *WebSocketSpotHandle) TradeGroupHandle(data []byte) error {
	var (
		t             = time.Now().UnixMicro()
		resp          Resp_Trade
		err           error
		price, amount float64
		dealtime      int64
	)
	if strings.Contains(string(data), "pong") {
		return nil
	}
	if strings.Contains(string(data), "success") {
		logger.Logger.Info(string(data))
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		logger.Logger.Error(err)
		return err
	}
	if resp.Error.Code != 0 {
		logger.Logger.Error("subscribe error: ", resp.Error.Message)
		return errors.New(resp.Error.Message)
	}
	// if resp.Result.Status != "" {
	// 	if strings.Contains(resp.Result.Status, "success") {
	// 		return nil
	// 	} else {
	// 		logger.Logger.Error("订阅错误：", string(data))
	// 		return errors.New("订阅错误")
	// 	}
	// }
	price, err = strconv.ParseFloat(resp.Result.Price, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return err
	}
	amount, err = strconv.ParseFloat(resp.Result.Amount, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return err
	}
	dealtimeS := strings.Split(resp.Result.CreateTimeMs, ".")[0]

	dealtime, err = strconv.ParseInt(dealtimeS, 10, 64)
	res := &client.WsTradeRsp{
		//OrderId: resp.Data[0].TradeId,
		Symbol: ParseSymbol(resp.Result.CurrencyPair),
		//ExchangeTime: resp.Data.Ts,
		Price:       price,
		Amount:      amount,
		TakerSide:   GetSide(resp.Result.Side),
		DealTime:    dealtime * 1000,
		ReceiveTime: t,
	}
	if _, ok := b.tradeGroupChanMap[res.Symbol]; ok {
		base.SendChan(b.tradeGroupChanMap[res.Symbol], res, "tradeGroup")
		//b.tradeGroupChanMap[res.Symbol] <- res
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) SetBookTickerGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) {
	if chMap != nil {
		b.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp, 1000)
		for k, v := range chMap {
			b.bookTickerGroupChanMap[TranInstId(k)] = v
		}
	}
}

func (b *WebSocketSpotHandle) BookTickerGroupHandle(data []byte) error {
	// 跳过前两个确认信息
	var (
		resp Resp_Ticker
		err  error
		ask  *depth.DepthLevel
		bid  *depth.DepthLevel
	)
	if strings.Contains(string(data), "pong") {
		return nil
	}
	if strings.Contains(string(data), "success") {
		logger.Logger.Info(string(data))
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		logger.Logger.Error(err)
		return err
	}
	if resp.Error.Code != 0 {
		logger.Logger.Error("subscribe error: ", resp.Error.Message)
		return errors.New(resp.Error.Message)
	}
	// if resp.Result.Status != "" {
	// 	if strings.Contains(resp.Result.Status, "success") {
	// 		return nil
	// 	} else {
	// 		logger.Logger.Error("订阅错误：", string(data))
	// 		return errors.New("订阅错误")
	// 	}
	// }

	price, _ := strconv.ParseFloat(resp.Result.LowestAsk, 64)
	ask = &depth.DepthLevel{
		Price: price,
	}
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	price, _ = strconv.ParseFloat(resp.Result.HighestBid, 64)
	bid = &depth.DepthLevel{
		Price: price,
	}
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}

	res := client.WsBookTickerRsp{
		Symbol:       ParseSymbol(resp.Result.CurrencyPair),
		ExchangeTime: int64(resp.Time) * 1000000,
		Ask:          ask,
		Bid:          bid,
		ReceiveTime:  time.Now().UnixMicro(),
	}
	if len(b.depthLimitChan) > 0 {
		content := fmt.Sprint("channel slow:", len(b.depthLimitChan))
		logger.Logger.Warn(content)
	}
	if _, ok := b.bookTickerGroupChanMap[res.Symbol]; ok {
		base.SendChan(b.bookTickerGroupChanMap[res.Symbol], &res, "booktickerGroup")
		//b.bookTickerGroupChanMap[res.Symbol] <- &res
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil

}

func (b *WebSocketSpotHandle) SetDepthIncrementSnapshotGroupChannel(chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chDeltaMap != nil {
		b.depthIncrementSnapshotDeltaGroupChanMap = make(map[string]chan *client.WsDepthRsp, 100)
		for k, v := range chDeltaMap {
			b.depthIncrementSnapshotDeltaGroupChanMap[TranInstId(k)] = v
		}
	}
	if chFullMap != nil {
		b.depthIncrementSnapshotFullGroupChanMap = make(map[string]chan *depth.Depth, 100)
		for k, v := range chFullMap {
			b.depthIncrementSnapshotFullGroupChanMap[TranInstId(k)] = v
		}
	}
}

func (b *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte) error {
	var (
		// fullRes    *gate_api.Resp_Market_Books
		deltaDepth = &base.DeltaDepthUpdate{}
		book       = &base.OrderBook{}
		depth_     = &depth.Depth{}
		depthMsg   = &Resp_Depth{}
		sym        string
		err        error
		t          = time.Now().UnixMicro()
	)
	if strings.Contains(string(data), "pong") {
		return nil
	}
	if strings.Contains(string(data), "success") {
		logger.Logger.Info(string(data))
		return nil
	}
	b.Lock.Lock()
	defer b.Lock.Unlock()
	if err = b.HandleRespErr(data, depthMsg); err != nil {
		logger.Logger.Error(err)
		return err
	}
	if depthMsg.Result.Status != "" {
		if depthMsg.Error.Code != 0 {
			logger.Logger.Error("subscribe error: ", depthMsg.Error.Message)
			return errors.New(depthMsg.Error.Message)
		}
		return nil
	}

	sym = ParseSymbol(depthMsg.Result.CurrencyPair)
	depth2Delta(depthMsg, deltaDepth)
	// 给予时间
	deltaDepth.TimeReceive = t
	deltaDepth.TimeExchange = int64(depthMsg.Result.TimeInMilli) * 1000
	if !b.firstMap[sym] { // 第一次要获取全量
		for {

			book, err = b.GetFullDepth(depthMsg.Result.CurrencyPair)
			if err != nil {
				logger.Logger.Error("获取全量问题 ", err.Error())
				return err
			}
			if book.UpdateId < depthMsg.Result.FirstId {
				continue
			} else {
				break
			}
		}

		if book.UpdateId+1 >= depthMsg.Result.FirstId && book.UpdateId+1 <= depthMsg.Result.LastId {
			// 处理增量
			// 合并
			base.UpdateBidsAndAsks(deltaDepth, book, 1000, depth_)
			//transDepth2OrderBook(depth_, book)
		}
		b.DepthCacheMap.Set(sym, book)
		b.firstMap[sym] = true
		b.firstFullSentMap[sym] = true

		return nil
	}

	// 拿到book
	res_, ok := b.DepthCacheMap.Get(sym)
	if !ok {
		err = errors.New("提取全量数据错误")
		logger.Logger.Error("提取全量数据错误")
		b.firstMap[sym] = false
		return err
	}
	book = res_.(*base.OrderBook)
	if book.UpdateId+1 < depthMsg.Result.FirstId {
		err = errors.New("id对应不上")
		logger.Logger.Error("id对应不上")
		b.firstMap[sym] = false
		return err
	} else if book.UpdateId+1 > depthMsg.Result.LastId {
		// 如果当前全量id大于推送的增量id就跳过
		return nil
	}
	base.UpdateBidsAndAsks(deltaDepth, book, 1000, depth_)
	b.DepthCacheMap.Set(ParseSymbol(sym), book)
	// 发送增量数据
	if b.IsPublishDelta {
		if _, ok := b.depthIncrementSnapshotDeltaGroupChanMap[sym]; ok {
			base.SendChan(b.depthIncrementSnapshotDeltaGroupChanMap[sym], deltaDepth.Transfer2Depth(), "deltaDepth")
			//b.depthIncrementSnapshotDeltaGroupChanMap[sym] <- deltaDepth.Transfer2Depth()
		} else {
			logger.Logger.Warn("get symbol from channel map err:", book.Symbol)
		}
	}
	//发送全量
	if b.IsPublishFull {
		if _, ok := b.depthIncrementSnapshotFullGroupChanMap[sym]; ok {
			if b.firstFullSentMap[sym] {
				b.firstFullSentMap[sym] = false
				depth_.Hdr = base.MakeFirstDepthHdr()
			}
			base.SendChan(b.depthIncrementSnapshotFullGroupChanMap[sym], depth_, "fullDepth")
			//b.depthIncrementSnapshotFullGroupChanMap[sym] <- &depth_
		} else {
			logger.Logger.Warn("get symbol from channel map err:", book.Symbol)
		}
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(sym, time.Now().UnixMicro())
	return nil
}

func full2Depth(fullRes *gate_api.Resp_Market_Books, depth *Resp_Depth) {
	depth.Result.Asks = fullRes.Asks
	depth.Result.Bids = fullRes.Bids
	depth.Result.LastId = fullRes.Id
	depth.Result.TimeInMilli = fullRes.Current
	return
}

func depth2Delta(depth *Resp_Depth, delta *base.DeltaDepthUpdate) {
	delta.Asks = delta.Asks[:0]
	delta.Bids = delta.Bids[:0]
	delta.Symbol = ParseSymbol(depth.Result.CurrencyPair)
	delta.UpdateEndId = depth.Result.LastId
	delta.TimeExchange = depth.Result.TimeInMilli
	ParseOrder(depth.Result.Asks, &delta.Asks)
	ParseOrder(depth.Result.Bids, &delta.Bids)
	sort.Stable(delta.Asks)
	sort.Stable(sort.Reverse(delta.Bids))
	return
}

func delta2Book(delta *base.DeltaDepthUpdate, book *base.OrderBook) {
	book.Symbol = delta.Symbol
	book.Asks = delta.Asks
	book.Bids = delta.Bids
	book.TimeExchange = uint64(delta.TimeReceive)
	book.UpdateId = delta.UpdateEndId
}

func (b *WebSocketSpotHandle) CreateInitial(book map[string]interface{}, key string) []Level {
	var list []Level = make([]Level, 0)
	for _, element := range book[key].([]interface{}) {
		price_interface := element.([]interface{})[0]
		price_str := price_interface.(string)
		price, err := decimal.NewFromString(price_str)
		if err != nil {
			log.Fatal(err)
		}
		vol_interface := element.([]interface{})[1]
		vol_str := vol_interface.(string)
		vol, err := decimal.NewFromString(vol_str)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, Level{Price: price, Volume: vol})
	}
	return list
}

func (b *WebSocketSpotHandle) HandleRespErr(data []byte, resp interface{}) error {
	var (
		err error
	)
	err = json.Unmarshal(data, resp)
	if err != nil {
		return errors.New("json unmarshal data err:" + string(data))
	}
	return nil
}
