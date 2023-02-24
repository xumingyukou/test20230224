package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"errors"
	"math/rand"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/sdk"
)

var (
	numberCheckRoutine int64 = 0
)

func (h *WebSocketSpotHandle) SetDepthIncrementSnapShotConf(symbols []*client.SymbolInfo, conf *base.IncrementDepthConf) {
	if conf.DepthCapLevel <= 0 {
		conf.DepthCapLevel = 1000
	}
	if conf.CheckTimeSec <= 0 {
		conf.CheckTimeSec = 3600 + rand.Intn(1200)
	}
	if conf.DepthCheckLevel <= 0 {
		conf.DepthCheckLevel = 100
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
	h.IncrementDepthConf = conf
	h.symbolMap = make(map[string]*client.SymbolInfo)

	for _, symbol := range symbols {
		var (
			CheckDepthCacheChan = make(chan *base.OrderBook, 100)
		)
		conf.DepthCacheMap.Set(symbol.Symbol, nil)
		conf.DepthDeltaUpdateMap.Set(symbol.Symbol, make(chan *base.DeltaDepthUpdate))
		conf.CheckDepthCacheChanMap.Set(symbol.Symbol, CheckDepthCacheChan)
		conf.DepthNotMatchChanMap[symbol] = make(chan bool, conf.DepthCapLevel)
		conf.CheckStates.Set(symbol.Symbol, false)
		h.symbolMap[symbol.Symbol] = symbol
	}
	for _, symbol := range symbols {
		go h.updateDeltaDepth(symbol.Symbol)
	}
}

func (h *WebSocketSpotHandle) Reset() {
	h.CheckDepthCacheChanMap = sdk.NewCmapI()
	h.DepthCacheMap = sdk.NewCmapI()
	h.DepthDeltaUpdateMap = sdk.NewCmapI()
}

func (h *WebSocketSpotHandle) splitIncrementBySequence(inc *KucoinIncrement) []*base.DeltaDepthUpdate {
	baseId := inc.SequenceStart
	length := inc.SequenceEnd - baseId + 1
	depths := make([]*base.DeltaDepthUpdate, length, length)
	for i := range depths {
		depths[i] = &base.DeltaDepthUpdate{
			UpdateStartId: baseId + int64(i),
			UpdateEndId:   baseId + int64(i),
			Market:        common.Market_SPOT,
			Type:          common.SymbolType_SPOT_NORMAL,
			Symbol:        inc.Symbol,
			Bids:          make([]*depth.DepthLevel, 0, 5),
			Asks:          make([]*depth.DepthLevel, 0, 5),
			TimeReceive:   inc.TimeReceive,
		}
	}
	for _, ask := range inc.Asks {
		depths[ask.SequenceId-baseId].Asks = append(depths[ask.SequenceId-baseId].Asks, &ask.DepthLevel)
	}
	for _, bid := range inc.Bids {
		depths[bid.SequenceId-baseId].Bids = append(depths[bid.SequenceId-baseId].Bids, &bid.DepthLevel)
	}
	return depths
}

func (h *WebSocketSpotHandle) updateDeltaDepth(symbol string) {
	var (
		obj           interface{}
		ch            chan *base.DeltaDepthUpdate
		ok            bool
		delta         *base.DeltaDepthUpdate
		deltaList     = make([]*base.DeltaDepthUpdate, 0, 10)
		depthCache    *base.OrderBook // 本地cache
		fullCh        chan *depth.Depth
		err           error
		retryCount    = 0
		match         bool
		found         = false
		firstFullSent = false
	)
	if obj, ok = h.DepthDeltaUpdateMap.Get(symbol); ok {
		ch, ok = obj.(chan *base.DeltaDepthUpdate)
	}

	if !ok {
		logger.Logger.Error("get DepthDeltaUpdateMap", symbol)
		return
	}

	symbolInfo, ok := h.symbolMap[symbol]
	if !ok {
		logger.Logger.Error("get symbolinfo in updateDeltaDepth, symbol:", symbol)
		return
	}
	notMatchChan, ok := h.DepthNotMatchChanMap[symbolInfo]
	if !ok {
		logger.Logger.Error("get depth not match channel in updateDeltaDepth, symbol:", symbol)
		return
	}

	if fullCh, ok = h.depthIncrementSnapshotFullGroupChanMap[symbol]; !ok {
		logger.Logger.Error("get symbol from channel map err:", symbol)
	}

LOOP:
	for {
		select {
		case match, ok = <-notMatchChan:
			if !ok {
				logger.Logger.Error("not match channel close: ", symbol)
				continue
			}
			if !match {
				depthCache = nil // 清空缓存 重新拉全量
				found = false
			}
		case delta, ok = <-ch:
			if !ok {
				logger.Logger.Error("DepthDeltaUpdateMap channel close:", symbol)
				break LOOP
			}
			deltaList = append(deltaList, delta)
			if depthCache == nil {
				firstFullSent = false
				if depthCache, err = h.GetFullDepth(symbol); err != nil {
					logger.Logger.Error("get full depth in update delta depth", symbol, " ", retryCount, " ", err)
					retryCount++
					// if retryCount > 3 {
					// 	break LOOP
					// }
					continue
				}
			} else {
				retryCount = 0
			}

			if depthCache.UpdateId+1 < deltaList[0].UpdateStartId {
				logger.Logger.Warnf("repull snapshot %s %d %d-%d length %d, wait for 5s", symbol, depthCache.UpdateId, deltaList[0].UpdateStartId, deltaList[len(deltaList)-1].UpdateEndId, len(deltaList))
				time.Sleep(time.Second * 5) // wait for 5s
				//全量的id比deltaList中都早，需要重新获取一次全量快照
				if depthCache, err = h.GetFullDepth(symbol); err != nil {
					logger.Logger.Error("get full depth error", err)
				}
				found = false
				firstFullSent = false
				continue
			} else if depthCache.UpdateId >= deltaList[len(deltaList)-1].UpdateStartId {
				//全量的id比deltaList中都晚，需要等待deltaList更新, 直接清空
				logger.Logger.Tracef("outdated delta list %s %d %d-%d length %d", symbol, depthCache.UpdateId, deltaList[0].UpdateStartId, deltaList[len(deltaList)-1].UpdateEndId, len(deltaList))
				deltaList = append([]*base.DeltaDepthUpdate{})
				found = false
				firstFullSent = false
				continue
			} else {
				// 全量的id在deltaList中可能可以找到
				if !found {
					logger.Logger.Infof("found delta list %s %d %d-%d length %d", symbol, depthCache.UpdateId, deltaList[0].UpdateStartId, deltaList[len(deltaList)-1].UpdateEndId, len(deltaList))
					found = true
				}
				target := -1
				for i, delta := range deltaList {
					if delta.UpdateStartId == depthCache.UpdateId+1 {
						target = i
						break
					}
				}
				if target == -1 {
					// 没有找到全量的id
					logger.Logger.Error("deltaList is not successive", symbol)
					if depthCache, err = h.GetFullDepth(symbol); err != nil {
						logger.Logger.Error("update get full depth error", err)
					}
					firstFullSent = false
					continue
				}
				//判断增量是否为空
				zeroDelta := true
				for _, deltaCache := range deltaList[target:] {
					// h.MergeDeltaIntoSnapshot(delta, depthCache)
					base.UpdateBidsAndAsks(deltaCache, depthCache, h.DepthCapLevel, nil)
					zeroDelta = zeroDelta && isEmptyDelta(deltaCache)
				}
				// h.afterUpdateHandle(symbol, depthCache, zeroDelta)

				if !zeroDelta {
					fullToSend := &depth.Depth{}
					depthCache.Transfer2Depth(h.DepthLevel, fullToSend)
					if !firstFullSent {
						fullToSend.Hdr = base.MakeFirstDepthHdr()
						firstFullSent = true
					}
					base.SendChan(fullCh, fullToSend, "afterUpdateHandle", symbol)
				}
				// 增量使用完，清空增量
				deltaList = append([]*base.DeltaDepthUpdate{})
			}
		}
	}
}

func isEmptyDelta(delta *base.DeltaDepthUpdate) bool {
	if len(delta.Asks) == 0 && len(delta.Bids) == 0 {
		return true
	}

	allZero := true
	for _, ask := range delta.Asks {
		allZero = allZero && ask.Amount == 0 && ask.Price == 0
		if !allZero {
			break
		}
	}

	for _, bid := range delta.Bids {
		allZero = allZero && bid.Amount == 0 && bid.Price == 0
		if !allZero {
			break
		}
	}
	return allZero
}

func (h *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte, receiveTime int64) error {
	h.Lock.Lock()
	defer h.Lock.Unlock()

	inc, err := ParseKucoinIncrement(data, receiveTime)
	if err != nil {
		logger.Logger.Error("parse kucoin increment:", err.Error(), " data:", string(data))
		return err
	}

	delta := inc.toDeltaDepthUpdate()
	symbol := delta.Symbol

	if h.IsPublishDelta {
		if ch, ok := h.depthIncrementSnapshotDeltaGroupChanMap[symbol]; ok {
			if !isEmptyDelta(delta) {
				base.SendChan(ch, delta.Transfer2Depth(), "depthIncrementSnapshotGroup")
			}
		} else {
			logger.Logger.Error("get symbol from channel map err:", symbol)
			return errors.New("depth delta update map " + symbol)
		}
	}

	if h.IsPublishFull {
		var (
			ch  chan *base.DeltaDepthUpdate
			ok  bool
			obj interface{}
		)
		if obj, ok = h.DepthDeltaUpdateMap.Get(symbol); ok {
			ch, ok = obj.(chan *base.DeltaDepthUpdate)
		}
		if ok {
			deltaList := h.splitIncrementBySequence(inc)
			for _, d := range deltaList {
				base.SendChan(ch, d, "DeltaDepthUpdate", symbol)
			}
		} else {
			logger.Logger.Error("depth delta update map", symbol)
			return errors.New("depth delta update map " + symbol)
		}
	}
	h.CheckSendStatus.CheckUpdateTimeMap.Set(symbol, receiveTime)
	return nil
}

// func (h *WebSocketSpotHandle) MergeDeltaIntoSnapshot(delta *base.DeltaDepthUpdate, snapshot *base.OrderBook) {
// 	// asks: normal
// 	// bids: reverse
// 	snapshot.UpdateId = delta.UpdateEndId
// 	snapshot.TimeReceive = uint64(delta.TimeReceive)
// 	snapshot.Asks = MergeDepthList(snapshot.Asks, delta.Asks, false, snapshot.UpdateId)
// 	snapshot.Bids = MergeDepthList(snapshot.Bids, delta.Bids, true, snapshot.UpdateId)
// }

// func MergeDepthList(shot []*depth.DepthLevel, delta []*depth.DepthLevel, reverse bool, updateId int64) []*depth.DepthLevel {
// 	if len(delta) == 0 {
// 		return shot
// 	}
// 	newDepths := make([]*depth.DepthLevel, 0, len(shot)+len(delta))

// 	for i, j := 0, 0; i < len(shot) || j < len(delta); {
// 		if i == len(shot) {
// 			// 价格不为零更新
// 			if delta[j].Price != 0 {
// 				newDepths = append(newDepths, delta[j])
// 			}
// 			j++
// 			continue
// 		}
// 		if j == len(delta) {
// 			// 价格不为零更新
// 			if shot[i].Price != 0 {
// 				newDepths = append(newDepths, shot[i])
// 			}
// 			i++
// 			continue
// 		}

// 		if shot[i].Price == delta[j].Price {
// 			// 全量和增量价格记录相同
// 			// 增量数量为零则删除, 非零则更新
// 			if delta[j].Amount != 0 {
// 				newDepths = append(newDepths, delta[j])
// 			}
// 			i++
// 			j++
// 		} else if (shot[i].Price < delta[j].Price && !reverse) || (shot[i].Price > delta[j].Price && reverse) {
// 			// 全量价格在增量前, 保留全量
// 			newDepths = append(newDepths, shot[i])
// 			i++
// 		} else {
// 			// 全量价格在增量后, 加入增量
// 			if delta[j].Amount != 0 {
// 				newDepths = append(newDepths, delta[j])
// 			}
// 			j++
// 		}
// 	}

// 	return newDepths
// }

// func (h *WebSocketSpotHandle) updateBidsAndAsks(deltaDepth *base.DeltaDepthUpdate, DepthCache *base.OrderBook, depthLevel int) {
// 	base.MergeDepth(base.Ask, &DepthCache.Asks, deltaDepth.Asks)
// 	base.MergeDepth(base.Bid, &DepthCache.Bids, deltaDepth.Bids)
// 	// 调整最大容量
// 	if len(DepthCache.Asks) > depthLevel {
// 		DepthCache.Asks = DepthCache.Asks[:depthLevel]
// 	}
// 	if len(DepthCache.Bids) > depthLevel {
// 		DepthCache.Bids = DepthCache.Bids[:depthLevel]
// 	}
// 	DepthCache.Symbol = deltaDepth.Symbol
// 	DepthCache.TimeExchange = uint64(deltaDepth.TimeExchange)
// 	DepthCache.TimeReceive = uint64(deltaDepth.TimeReceive)
// 	DepthCache.UpdateId = deltaDepth.UpdateEndId
// }

// func (h *WebSocketSpotHandle) Check() {

// 	/**
// 	定期检查全量数据
// 	*/
// 	defer func() {
// 		if err := recover(); err != nil {
// 			logger.Logger.Error("incr depth check panic:", logger.PanicTrace(err))
// 			time.Sleep(time.Second)
// 			go h.Check()
// 		}
// 	}()
// 	heartTimer := time.NewTimer(time.Duration(h.CheckTimeSec) * time.Second)
// LOOP:
// 	for {
// 		select {
// 		case <-heartTimer.C:
// 			for _, symbol := range h.DepthCacheMap.Keys() {
// 				// go h.checkFullDepth(symbol)
// 				h.checkFullDepth(symbol)
// 				time.Sleep(time.Second * 10)
// 			}
// 			heartTimer.Reset(time.Duration(h.CheckTimeSec) * time.Second)
// 		case <-h.Ctx.Done():
// 			h.Reset()
// 			break LOOP
// 		}
// 	}
// }

// func (h *WebSocketSpotHandle) checkFullDepth(symbol string) {
// 	var (
// 		fullDepth *base.OrderBook
// 		ch        chan *base.OrderBook
// 		content   interface{}
// 		ok        bool
// 		err       error
// 	)
// 	h.CheckStates.Set(symbol, true)
// 	defer func() { h.CheckStates.Set(symbol, false) }()
// 	fullDepth = &base.OrderBook{}
// 	if content, ok = h.CheckDepthCacheChanMap.Get(symbol); ok {
// 		ch, ok = content.(chan *base.OrderBook)
// 	}
// 	if !ok {
// 		logger.Logger.Error("get CheckDepthCacheChanMap err, symbol:", symbol)
// 		return
// 	}
// 	atomic.AddInt64(&numberCheckRoutine, 1)
// 	logger.Logger.Info(">>>>> check full depth ", symbol, " time:", time.Now(), "number of routines:", atomic.LoadInt64(&numberCheckRoutine))
// 	start := time.Now()

// 	cnt := 0
// 	symbolInfo, ok := h.symbolMap[symbol]
// 	if !ok {
// 		logger.Logger.Error("get symbolinfo in checkFullDepth, symbol:", symbol)
// 		return
// 	}
// 	notMatchChan, ok := h.DepthNotMatchChanMap[symbolInfo]
// 	if !ok {
// 		logger.Logger.Error("get depth not match channel err in checkFullDepth, symbol:", symbol)
// 		return
// 	}

// MAINLOOP:
// 	for {
// 		logger.Logger.Info("----- in check full depth for loop ", symbol, "loop ", cnt)
// 		cnt++
// 		select {
// 		case dep := <-ch:
// 			if fullDepth == nil || fullDepth.UpdateId == 0 {
// 				if fullDepth, err = h.GetFullDepthLimit(symbol, h.DepthCheckLevel); err != nil {
// 					logger.Logger.Error("checkFullDepth get full depth error", err)
// 					time.Sleep(time.Second * 30)
// 					continue
// 				}
// 			}

// 			logger.Logger.Info("full ", fullDepth.UpdateId, "cache ", dep.UpdateId, " symbol ", symbol)
// 			if fullDepth.UpdateId > dep.UpdateId {
// 				continue
// 			} else if fullDepth.UpdateId == dep.UpdateId {
// 				if dep.Equal(fullDepth, h.DepthCheckLevel) {
// 					logger.Logger.Info("check full depth: correct", symbol)
// 					break MAINLOOP
// 				} else {
// 					logger.Logger.Error(symbol, "depth level:", h.DepthCheckLevel, "need reset")
// 					logger.SaveToFile(fmt.Sprintf("depth_error_%s.json", strings.Replace(symbol, "/", "_", 1)), "depCache", dep, "fullDepth", fullDepth)
// 					notMatchChan <- false
// 					break MAINLOOP
// 				}
// 			} else {
// 				logger.Logger.Info("outdated full", fullDepth.UpdateId, "cache ", dep.UpdateId, " symbol ", symbol)
// 				if fullDepth, err = h.GetFullDepthLimit(symbol, h.DepthCheckLevel); err != nil {
// 					logger.Logger.Error(symbol, "checkFullDepth < get full depth error", err)
// 					time.Sleep(time.Second * 30)
// 					continue
// 				}
// 			}
// 		}
// 	}
// 	atomic.AddInt64(&numberCheckRoutine, -1)
// 	logger.Logger.Info("<<<<< exit check full depth ", symbol, " time cost:", time.Now().Sub(start), "remaining number of routines:", atomic.LoadInt64(&numberCheckRoutine))
// }
