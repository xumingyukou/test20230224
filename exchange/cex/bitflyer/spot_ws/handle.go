package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/common"
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
	firstFullSentMap                        map[string]bool
	Lock                                    sync.Mutex // 锁

	CheckSendStatus *base.CheckDataSendStatus
}

func NewWebSocketSpotHandle(chanCap int64) *WebSocketSpotHandle {
	d := &WebSocketSpotHandle{}
	d.reSubscribe = make(chan req, chanCap)
	d.tradeChan = make(chan *client.WsTradeRsp, chanCap)
	d.depthIncrementChan = make(chan *client.WsDepthRsp, chanCap)
	d.depthLimitChan = make(chan *client.WsDepthRsp, chanCap)
	d.bookTickerChan = make(chan *client.WsBookTickerRsp, chanCap)
	d.accountChan = make(chan *client.WsAccountRsp, chanCap)
	d.balanceChan = make(chan *client.WsBalanceRsp, chanCap)
	d.orderChan = make(chan *client.WsOrderRsp, chanCap)
	d.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp)
	d.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp)
	d.firstFullSentMap = make(map[string]bool)
	d.CheckSendStatus = base.NewCheckDataSendStatus()

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
		b.firstFullSentMap[symbol.Symbol] = false
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
	var (
		resp   Resp_Depth
		err    error
		symbol string
		t      = time.Now().UnixMicro()
	)

	if err = b.HandleRespErr(data, &resp); err != nil {
		logger.Logger.Error(err)
		return err
	}
	asks := b.DepthLevelParse(resp.Params.Message.Asks)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	bids := b.DepthLevelParse(resp.Params.Message.Bids)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	symbol = ParseSymbol(strings.TrimLeft(resp.Params.Channel, "lightning_board_"))
	res := client.WsDepthRsp{
		Symbol:      symbol,
		Asks:        asks,
		Bids:        bids,
		ReceiveTime: t,
	}

	if _, ok := b.depthIncrementGroupChanMap[res.Symbol]; ok {
		base.SendChan(b.depthIncrementGroupChanMap[res.Symbol], &res, "depthIncrement")
		//b.depthIncrementGroupChanMap[res.Symbol] <- &res
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(symbol, time.Now().UnixMicro())
	return nil
}

// []interface转[][]float64

func (b *WebSocketSpotHandle) DepthLevelParse(levelList_ []levelList) []*depth.DepthLevel {
	var (
		res []*depth.DepthLevel
	)
	for _, level := range levelList_ {
		res = append(res, &depth.DepthLevel{
			Price:  level.Price,
			Amount: level.Price,
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
		t    = time.Now().UnixMicro()
		resp Resp_Trade
		err  error
	)

	if err = b.HandleRespErr(data, &resp); err != nil {
		logger.Logger.Error(err)
		return err
	}
	symbol := ParseSymbol(strings.TrimLeft(resp.Params.Channel, "lightning_executions_"))

	var res []*client.WsTradeRsp
	for _, v := range resp.Params.Message {
		tmpres := &client.WsTradeRsp{
			Symbol:       symbol,
			ReceiveTime:  t,
			Price:        v.Price,
			Amount:       v.Size,
			TakerSide:    GetSide(v.Side),
			ExchangeTime: v.ExecDate.UnixMicro(),
		}
		res = append(res, tmpres)
	}
	if _, ok := b.tradeGroupChanMap[symbol]; ok {
		for _, v := range res {
			base.SendChan(b.tradeGroupChanMap[symbol], v, "tradeGroup")
		}
	} else {
		logger.Logger.Warn("get symbol from channel map err:", symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(symbol, time.Now().UnixMicro())
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
	var (
		t    = time.Now().UnixMicro()
		resp Resp_Ticker
		err  error
		ask  *depth.DepthLevel
		bid  *depth.DepthLevel
	)

	if err = b.HandleRespErr(data, &resp); err != nil {
		logger.Logger.Error(err)
		return err
	}

	ask = &depth.DepthLevel{
		Price:  resp.Params.Message.BestAsk,
		Amount: resp.Params.Message.BestAskSize,
	}
	bid = &depth.DepthLevel{
		Price:  resp.Params.Message.BestBid,
		Amount: resp.Params.Message.BestBidSize,
	}
	time_ := resp.Params.Message.Timestamp.UnixMicro()
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}

	res := client.WsBookTickerRsp{
		Symbol:       ParseSymbol(resp.Params.Message.ProductCode),
		ExchangeTime: time_,
		Ask:          ask,
		Bid:          bid,
		ReceiveTime:  t,
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

// 通过全量数据反推增量
func (b *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte) error {
	var (
		fullRes  = &base.DeltaDepthUpdate{}
		depth_   = &client.WsDepthRsp{}
		depthMsg = &Resp_Depth{}
		ok       bool
		sym      string
		err      error
		t        = time.Now().UnixMicro()
	)
	b.Lock.Lock()
	defer b.Lock.Unlock()
	if err = b.HandleRespErr(data, depthMsg); err != nil {
		logger.Logger.Error(err)
		return err
	}
	sym = strings.TrimLeft(depthMsg.Params.Channel, "lightning_board_snapshot_")
	sym = ParseSymbol(sym)
	//拿到信息并排序
	depth2Delta(depthMsg, fullRes)
	fullRes.TimeReceive = t
	var m interface{}
	if m, ok = b.DepthCacheMap.Get(sym); !ok {
		b.DepthCacheMap.Set(sym, fullRes)
		return nil
	}
	depth_ = getDeltaDepth(fullRes, m.(*base.DeltaDepthUpdate))
	depth_.ReceiveTime = t
	b.DepthCacheMap.Set(sym, fullRes)

	// 增量
	if b.IsPublishDelta {
		if _, ok := b.depthIncrementSnapshotDeltaGroupChanMap[sym]; ok {
			base.SendChan(b.depthIncrementSnapshotDeltaGroupChanMap[sym], depth_, "deltaDepth")
			//b.depthIncrementSnapshotDeltaGroupChanMap[sym] <- deltaDepth.Transfer2Depth()
		} else {
			logger.Logger.Warn("get symbol from channel map err:", sym)
		}
	}
	//	//发送全量
	if b.IsPublishFull {
		if _, ok := b.depthIncrementSnapshotFullGroupChanMap[sym]; ok {

			fullToSent := update2depth(fullRes)
			fullToSent.Exchange = common.Exchange_BITFLYER
			fullToSent.Symbol = sym
			fullToSent.Market = common.Market_SPOT
			fullToSent.Type = common.SymbolType_SPOT_NORMAL
			if !b.firstFullSentMap[sym] {
				fullToSent.Hdr = base.MakeFirstDepthHdr()
				b.firstFullSentMap[sym] = true
			} else {
				fullToSent.Hdr = &common.MsgHeader{
					Type: uint32(common.MsgType_DPETH),
				}
			}
			base.SendChan(b.depthIncrementSnapshotFullGroupChanMap[sym], fullToSent, "fullDepth")
			//b.depthIncrementSnapshotFullGroupChanMap[sym] <- &depth_
		} else {
			logger.Logger.Warn("get symbol from channel map err:", sym)
		}
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(sym, time.Now().UnixMicro())
	return nil
}

func update2depth(update *base.DeltaDepthUpdate) (res *depth.Depth) {
	res = &depth.Depth{}
	res.Bids = update.Bids
	res.Asks = update.Asks
	res.TimeReceive = uint64(update.TimeReceive)
	return res
}

func getDeltaDepth(new *base.DeltaDepthUpdate, later *base.DeltaDepthUpdate) (res *client.WsDepthRsp) {
	res = &client.WsDepthRsp{}
	res.Asks = getDeltaLevel(new.Asks, later.Asks, Ask)
	res.Bids = getDeltaLevel(new.Bids, later.Bids, Bid)
	res.ReceiveTime = new.TimeReceive
	return
}

func getDeltaLevel(nLever base.DepthItemSlice, lLever base.DepthItemSlice, side SIDE) base.DepthItemSlice {
	res := base.DepthItemSlice{}
	lMap := make(map[float64]float64)
	nMap := make(map[float64]float64)
	for _, v := range lLever {
		lMap[v.Price] = v.Amount
	}
	for _, v := range nLever {
		nMap[v.Price] = v.Amount
		// 新的里面的price，旧的有，但是amount不一样，替换了
		if amount, ok := lMap[v.Price]; ok {
			if amount != v.Amount {
				res = append(res, &depth.DepthLevel{Price: v.Price, Amount: v.Amount})
			}
			// 新的里面price,旧的没有，表示插入
		} else {
			res = append(res, &depth.DepthLevel{Price: v.Price, Amount: v.Amount})
		}
	}
	// 如果旧的有，新的没有。表示删除了这个元素
	for price, _ := range lMap {
		if _, ok := nMap[price]; !ok {
			res = append(res, &depth.DepthLevel{Price: price, Amount: 0})
		}
	}
	if side == Ask {
		sort.Stable(res)
	} else if side == Bid {
		sort.Stable(sort.Reverse(res))
	}
	return res
}

func depth2Delta(depth *Resp_Depth, delta *base.DeltaDepthUpdate) {
	delta.Asks = delta.Asks[:0]
	delta.Bids = delta.Bids[:0]
	ParseOrder(depth.Params.Message.Asks, &delta.Asks)
	ParseOrder(depth.Params.Message.Bids, &delta.Bids)
	sort.Stable(delta.Asks)
	sort.Stable(sort.Reverse(delta.Bids))
	return
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
