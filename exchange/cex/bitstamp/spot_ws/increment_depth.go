package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/bitstamp/spot_api"
	"clients/logger"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/sdk"
)

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
		conf.DepthCheckLevel = 20
	}
	if conf.DepthCacheMap == nil {
		conf.DepthCacheMap = sdk.NewCmapI()
	}
	if conf.DepthDeltaUpdateMap == nil {
		conf.DepthDeltaUpdateMap = sdk.NewCmapI()
	}
	if conf.CheckDepthCacheChanMap == nil {
		conf.CheckDepthCacheChanMap = sdk.NewCmapI()
	}
	if conf.CheckStates == nil {
		conf.CheckStates = sdk.NewCmapI()
	}
	conf.DepthNotMatchChanMap = make(map[*client.SymbolInfo]chan bool)
	b.IncrementDepthConf = conf
	b.symbolMap = make(map[string]*client.SymbolInfo)
	for _, symbol := range symbols {
		conf.DepthCacheMap.Set(symbol.Symbol, nil)
		conf.DepthDeltaUpdateMap.Set(symbol.Symbol, make(chan *base.DeltaDepthUpdate, conf.DepthCapLevel))
		conf.DepthNotMatchChanMap[symbol] = make(chan bool, conf.DepthCapLevel)
		b.symbolMap[symbol.Symbol] = symbol
	}
	//fmt.Println("We got here")
	for _, symbol := range symbols {
		go b.UpdateDeltaDepth(symbol.Symbol)
	}
}
func (b *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte) error {
	var (
		resp       RespDepthStream
		deltaDepth = &base.DeltaDepthUpdate{}
		ok         bool
		content    interface{}
		chUpdate   chan *base.DeltaDepthUpdate
		err        error
	)
	b.Lock.Lock()
	defer b.Lock.Unlock()
	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if resp.Event == "bts:subscription_succeeded" {
		return nil
	} else if resp.Event == "bts:heartbeat" {
		return nil
	}
	transferDiffDepth(&resp, deltaDepth)
	if b.IsPublishDelta {
		if _, ok := b.depthIncrementSnapshotDeltaGroupChanMap[deltaDepth.Symbol]; ok {
			base.SendChan(b.depthIncrementSnapshotDeltaGroupChanMap[deltaDepth.Symbol], deltaDepth.Transfer2Depth(), "depthIncrementSnapshotDeltaGroupChanMap")
		} else {
			logger.Logger.Warn("get symbol from channel map err:", deltaDepth.Symbol)
		}
	}
	if b.IsPublishFull {
		if content, ok = b.DepthDeltaUpdateMap.Get(deltaDepth.Symbol); ok {
			chUpdate, ok = content.(chan *base.DeltaDepthUpdate)
		}
		base.SendChan(chUpdate, deltaDepth, "DeltaDepthUpdate", deltaDepth.Symbol)
	}
	//if rand.Intn(100) == 2 {
	//	b.CheckSendStatus.CheckUpdateTimeMap.Set(deltaDepth.Symbol, time.Date(2000, 04, 01, 01, 01, 01, 01, time.Local).UnixMicro())
	//	fmt.Println("We got here")
	//} else {
	b.CheckSendStatus.CheckUpdateTimeMap.Set(deltaDepth.Symbol, time.Now().UnixMicro())
	//}
	return nil
}

func (b *WebSocketSpotHandle) UpdateDeltaDepth(symbol string) {
	var (
		content             interface{}
		ok                  bool
		ch                  chan *base.DeltaDepthUpdate
		deltaDepthCacheList []*base.DeltaDepthUpdate
		diffDepth           *base.DeltaDepthUpdate
		depthCache          = &base.OrderBook{Symbol: symbol}
		runningExchangeTime int64
		err                 error
		firstSnapshotSent   = false
	)
	if content, ok = b.DepthDeltaUpdateMap.Get(symbol); ok {
		ch, ok = content.(chan *base.DeltaDepthUpdate)
	}
	if !ok {
		logger.Logger.Error("get DepthDeltaUpdateMap err:", symbol)
		return
	}
	runningExchangeTime = -1

LOOP:
	for {
		select {
		case diffDepth, ok = <-ch:
			if !ok {
				logger.Logger.Error("DepthDeltaUpdateMap channel close:", symbol)
				break LOOP
			}
			// 缓存depth的增量数据
			deltaDepthCacheList = append(deltaDepthCacheList, diffDepth)
			if len(deltaDepthCacheList) > b.DepthCapLevel {
				deltaDepthCacheList = deltaDepthCacheList[len(deltaDepthCacheList)-b.DepthCapLevel:]
			}
			// 根据UpdateId，判断是否更新depth，或者更新全量数据
			if depthCache.TimeExchange == uint64(runningExchangeTime) {
				//fmt.Println("Connected:", depthCache.TimeExchange, runningExchangeTime)
				// 清空缓存，更新depthCache
				deltaDepthCacheList = append([]*base.DeltaDepthUpdate{})
				base.UpdateBidsAndAsks(diffDepth, depthCache, b.DepthCapLevel, nil)
				runningExchangeTime = diffDepth.TimeExchange
				if !firstSnapshotSent {
					b.AfterUpdateHandle(symbol, depthCache.Copy(), true)
					firstSnapshotSent = true
				} else {
					b.AfterUpdateHandle(symbol, depthCache.Copy(), false)
				}
			} else if depthCache.TimeExchange < uint64(deltaDepthCacheList[0].TimeExchange) {
				//DepthCache的更新id比deltaList中都早，需要重新获取一次全量快照
				if err = b.UpdateFullDepth(depthCache); err != nil {
					logger.Logger.Error("get full depth error", err)
					time.Sleep(time.Second * 30)
				}
				runningExchangeTime = diffDepth.TimeExchange
				firstSnapshotSent = false
				continue
			} else if depthCache.TimeExchange > uint64(deltaDepthCacheList[len(deltaDepthCacheList)-1].TimeExchange) {
				//DepthCache的更新id比deltaList中都晚，需要继续更新deltaList
				deltaDepthCacheList = append([]*base.DeltaDepthUpdate{})
				firstSnapshotSent = false
				continue
			} else {
				//DepthCache的更新id在deltaList中可能有能对应的，遍历若没有发现，则同样需要重新获取全量
				for _, deltaDepth := range deltaDepthCacheList {
					//fmt.Println("Time counter:", depthCache.TimeExchange, deltaDepth.TimeExchange)
					if depthCache.TimeExchange > uint64(deltaDepth.TimeExchange) {
						continue
					} else if depthCache.TimeExchange >= uint64(deltaDepth.TimeExchange) {
						//fmt.Println("Bigger")
						base.UpdateBidsAndAsks(deltaDepth, depthCache, b.DepthCapLevel, nil)
						if !firstSnapshotSent {
							b.AfterUpdateHandle(symbol, depthCache.Copy(), true)
							firstSnapshotSent = true
						} else {
							b.AfterUpdateHandle(symbol, depthCache.Copy(), false)
						}
						runningExchangeTime = deltaDepth.TimeExchange
						continue
					}
				}
			}
		case <-b.Ctx.Done():
			break LOOP
		}

	}
}
func (b *WebSocketSpotHandle) AfterUpdateHandle(symbol string, depthCache *base.OrderBook, setFirstSnapshotHdr bool) {
	var (
		res = &depth.Depth{}
	)
	depthCache.Transfer2Depth(b.DepthLevel, res)
	if _, ok := b.depthIncrementSnapshotFullGroupChanMap[symbol]; ok {
		if setFirstSnapshotHdr {
			res.Hdr = base.MakeFirstDepthHdr()
		}
		base.SendChan(b.depthIncrementSnapshotFullGroupChanMap[symbol], res, "DepthIncrementSnapShotGroupHandle")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
}
func (b *WebSocketSpotHandle) UpdateFullDepth(depthCache *base.OrderBook) (err error) {
	fullDepth := &base.OrderBook{}
	fullDepth, err = b.GetFullDepth(depthCache.Symbol)
	if err == nil {
		depthCache.Bids = fullDepth.Bids
		depthCache.Asks = fullDepth.Asks
		depthCache.TimeExchange = fullDepth.TimeExchange
	}
	return
}
func transferDiffDepth(resp *RespDepthStream, diff *base.DeltaDepthUpdate) {

	timeStamp, err := strconv.ParseInt(resp.Data.MicroTimestamp, 10, 64)
	if err != nil {
		logger.Logger.Error("parse time err:", string(resp.Data.MicroTimestamp))
	}
	if chList := strings.Split(resp.Channel, "_"); len(chList) > 3 {
		diff.Symbol = totalSymbolMap[chList[3]]
	}
	diff.Market = common.Market_SPOT
	diff.Type = common.SymbolType_SPOT_NORMAL
	diff.TimeExchange = timeStamp
	diff.TimeReceive = time.Now().UnixMicro()
	//清空bids，asks
	diff.Bids = diff.Bids[:0]
	diff.Asks = diff.Asks[:0]
	ParseOrder(resp.Data.Bids, &diff.Bids)
	ParseOrder(resp.Data.Asks, &diff.Asks)
	sort.Stable(diff.Asks)
	sort.Stable(sort.Reverse(diff.Bids))
	return
}
func ParseOrder(orders []*spot_api.DepthItem, slice *base.DepthItemSlice) {
	var (
		price, amount float64
		err           error
	)
	for _, order := range orders {
		if price, err = strconv.ParseFloat(order.Price, 64); err != nil {
			logger.Logger.Errorf("order float parse price error [%s] , response data = %s", err, order)
			continue
		}
		if amount, err = strconv.ParseFloat(order.Amount, 64); err != nil {
			logger.Logger.Errorf("order float parse amount error [%s] , response data = %s", err, order)
			continue
		}
		*slice = append(*slice, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
}
