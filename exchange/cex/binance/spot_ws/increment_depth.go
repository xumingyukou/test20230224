package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/binance/c_api"
	"clients/exchange/cex/binance/spot_api"
	"clients/exchange/cex/binance/u_api"
	"clients/logger"
	"clients/transform"
	"fmt"
	"github.com/warmplanet/proto/go/common"
	"sort"
	"time"

	"github.com/warmplanet/proto/go/client"
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
	//if conf.DepthCacheListMap == nil {
	//	conf.DepthCacheListMap = sdk.NewCmapI()
	//}
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
		var (
			CheckDepthCacheChan = make(chan *base.OrderBook, 100)
		)
		sym := SymbolKeyGen(symbol.Symbol, symbol.Market, symbol.Type)
		conf.DepthCacheMap.Set(sym, symbol.Type)
		conf.DepthDeltaUpdateMap.Set(sym, make(chan *base.DeltaDepthUpdate, conf.DepthCapLevel))
		conf.CheckDepthCacheChanMap.Set(sym, CheckDepthCacheChan)
		conf.DepthNotMatchChanMap[symbol] = make(chan bool, conf.DepthCapLevel)
		conf.CheckStates.Set(sym, false)
		b.symbolMap[sym] = symbol
	}
	for _, symbol := range symbols {
		go b.updateDeltaDepth(symbol.Symbol, symbol.Market, symbol.Type)
	}
	go b.Check()
}

func ParseOrder(orders [][]string, slice *base.DepthItemSlice) {
	for _, order := range orders {
		price, amount, err := transform.ParsePriceAmountFloat(order)
		if err != nil {
			logger.Logger.Errorf("order float parse price error [%s] , response data = %s", err, order)
			continue
		}
		*slice = append(*slice, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
}

func GetSymbolName(symbol *client.SymbolInfo) string {
	if spot_api.IsUBaseSymbolType(symbol.Type) {
		return spot_api.GetSpotSymbolName(u_api.GetUBaseSymbol(symbol.Symbol, u_api.GetFutureTypeFromNats(symbol.Type)))
	} else {
		return spot_api.GetSpotSymbolName(c_api.GetCBaseSymbol(symbol.Symbol, u_api.GetFutureTypeFromNats(symbol.Type)))
	}
}

func (b *WebSocketSpotHandle) UpdateFullDepth(depthCache *base.OrderBook) (err error) {
	/*
		从api中获取全量的数据
	*/
	fullDepth := &base.OrderBook{}
	symbolInfo := &client.SymbolInfo{
		Symbol: depthCache.Symbol,
		Type:   depthCache.Type,
	}
	fullDepth, err = b.GetFullDepth(GetSymbolName(symbolInfo))
	if err == nil {
		depthCache.Bids = fullDepth.Bids
		depthCache.Asks = fullDepth.Asks
		depthCache.UpdateId = fullDepth.UpdateId
	}
	return
}

func transferDiffDepth(r *RespIncrementDepth, diff *base.DeltaDepthUpdate, marketType int) {
	// 将binance返回的结构，解析为DeltaDepthUpdate，并将bids和ask进行排序
	diff.Symbol, diff.Market, diff.Type = GetSymbolMarket(r.S, marketType)
	diff.TimeExchange = r.E1 * 1000
	if diff.TimeReceive-diff.TimeExchange > 10000 {
		//fmt.Println(diff.Symbol, diff.TimeExchange, diff.TimeReceive, diff.TimeReceive-diff.TimeExchange)
	}
	diff.UpdateEndId = r.U1
	if r.Pu != 0 {
		diff.UpdateNextId = r.Pu
	}
	diff.UpdateStartId = r.U
	//清空bids，asks
	diff.Bids = diff.Bids[:0]
	diff.Asks = diff.Asks[:0]
	ParseOrder(r.B, &diff.Bids)
	ParseOrder(r.A, &diff.Asks)
	// 币本位需要转换amount
	if marketType == 2 {
		for bidIdx, _ := range diff.Bids {
			diff.Bids[bidIdx].Amount = getCbaseQty(diff.Symbol, diff.Bids[bidIdx].Amount, diff.Bids[bidIdx].Price)
		}
		for askIdx, _ := range diff.Asks {
			diff.Asks[askIdx].Amount = getCbaseQty(diff.Symbol, diff.Asks[askIdx].Amount, diff.Asks[askIdx].Price)
		}
	}
	sort.Stable(diff.Asks)
	sort.Stable(sort.Reverse(diff.Bids))
	return
}

func (b *WebSocketSpotHandle) updateDeltaDepth(symbol string, market common.Market, symbolType common.SymbolType) {
	var (
		content             interface{}
		ok                  bool
		ch                  chan *base.DeltaDepthUpdate
		deltaDepthCacheList []*base.DeltaDepthUpdate
		diffDepth           *base.DeltaDepthUpdate
		depthCache          = &base.OrderBook{Symbol: symbol, Market: market, Type: symbolType}
		err                 error
		firstFullSent       = false
	)
	sym := SymbolKeyGen(symbol, market, symbolType)

	if content, ok = b.DepthDeltaUpdateMap.Get(sym); ok {
		ch, ok = content.(chan *base.DeltaDepthUpdate)
	}
	errCh := b.DepthNotMatchChanMap[b.symbolMap[symbol]]

	if !ok {
		logger.Logger.Error("get DepthDeltaUpdateMap err:", symbol)
		return
	}
LOOP:
	for {
		select {
		case diffDepth, ok = <-ch:
			if !ok {
				logger.Logger.Error("DepthDeltaUpdateMap channel close:", symbol)
				break LOOP
			}
			logger.Logger.Debug(depthCache.UpdateId+1, diffDepth.UpdateStartId, diffDepth.UpdateEndId)
			// 缓存depth的增量数据
			deltaDepthCacheList = append(deltaDepthCacheList, diffDepth)
			if len(deltaDepthCacheList) > b.DepthCapLevel {
				deltaDepthCacheList = deltaDepthCacheList[len(deltaDepthCacheList)-b.DepthCapLevel:]
			}
			// 根据UpdateId，判断是否更新depth，或者更新全量数据
			if (depthCache.UpdateId+1 >= diffDepth.UpdateStartId && depthCache.UpdateId+1 <= diffDepth.UpdateEndId) || (depthCache.UpdateId == diffDepth.UpdateNextId) {
				// 清空缓存，更新depthCache
				deltaDepthCacheList = append([]*base.DeltaDepthUpdate{})
				base.UpdateBidsAndAsks(diffDepth, depthCache, b.DepthCapLevel, nil)
				if !firstFullSent {
					firstFullSent = true
					b.afterUpdateHandle(sym, depthCache, true)
				} else {
					b.afterUpdateHandle(sym, depthCache, false)
				}
				// b.afterUpdateHandle(sym, depthCache.Copy())
			} else if depthCache.UpdateId+1 < deltaDepthCacheList[0].UpdateStartId {
				//DepthCache的更新id比deltaList中都早，需要重新获取一次全量快照
				firstFullSent = false
				// depthCache.Type = diffDepth.Type
				if err = b.UpdateFullDepth(depthCache); err != nil {
					logger.Logger.Error("get full depth error", err, depthCache.Symbol, depthCache.Market)
					fmt.Println("get full depth error", err, depthCache.Symbol, depthCache.Market)
					time.Sleep(time.Second * 30)
				}
				continue
			} else if depthCache.UpdateId+1 > deltaDepthCacheList[len(deltaDepthCacheList)-1].UpdateEndId {
				//DepthCache的更新id比deltaList中都晚，需要继续更新deltaList
				firstFullSent = false
				deltaDepthCacheList = append([]*base.DeltaDepthUpdate{})
				continue
			} else {
				//DepthCache的更新id在deltaList中可能有能对应的，遍历若没有发现，则同样需要重新获取全量
				dataLost := false
				for _, deltaDepth := range deltaDepthCacheList {
					if depthCache.UpdateId+1 > deltaDepth.UpdateEndId {
						continue
					} else if depthCache.UpdateId+1 >= deltaDepth.UpdateStartId && depthCache.UpdateId+1 <= deltaDepth.UpdateEndId {
						base.UpdateBidsAndAsks(deltaDepth, depthCache, b.DepthCapLevel, nil)
						continue
					} else {
						firstFullSent = false
						dataLost = true
						if err := b.UpdateFullDepth(depthCache); err != nil {
							logger.Logger.Error("update get full depth error", err)
							fmt.Println("get full depth error", err, depthCache.Symbol, depthCache.Market)
							time.Sleep(time.Second * 30)
							break
						}
					}
					logger.Logger.Error("need wait a while")
					continue
				}
				if !dataLost {
					if !firstFullSent {
						firstFullSent = true
						b.afterUpdateHandle(sym, depthCache, true)
					} else {
						b.afterUpdateHandle(sym, depthCache, false)
					}
				}
			}
		case <-errCh:
			firstFullSent = false
			logger.Logger.Info("check signal: reset full depth")
			depthCache.Type = diffDepth.Type
			if err = b.UpdateFullDepth(depthCache); err != nil {
				logger.Logger.Error("reset full depth error", err)
				fmt.Println("get full depth error", err, depthCache.Symbol, depthCache.Market)
				time.Sleep(time.Second * 30)
			}
		case <-b.Ctx.Done():
			break LOOP
		}
	}
}

func (b *WebSocketSpotHandle) afterUpdateHandle(sym string, depthCache *base.OrderBook, sendFirstFull bool) {
	var (
		res             = &depth.Depth{}
		content         interface{}
		ok, checkStatus bool
	)
	if content, ok = b.CheckStates.Get(sym); ok {
		checkStatus, ok = content.(bool)
	}
	depthCache.Transfer2Depth(b.DepthLevel, res)
	if sendFirstFull {
		res.Hdr = base.MakeFirstDepthHdr()
	}
	if _, ok := b.depthIncrementSnapshotFullGroupChanMap[sym]; ok {
		base.SendChan(b.depthIncrementSnapshotFullGroupChanMap[sym], res, fmt.Sprintf("%s %s", "DepthIncrementSnapShotGroupHandle", res.Symbol))
	} else {
		logger.Logger.Warn("[full depth] get symbol from channel map err:", res.Symbol)
	}
	if ok && checkStatus {
		var (
			ch chan *base.OrderBook
		)
		if content, ok = b.CheckDepthCacheChanMap.Get(sym); ok {
			if ch, ok = content.(chan *base.OrderBook); ok {
				depthCopy := depthCache.Copy()
				depthCopy.Limit(100)
				ch <- depthCache
				// depthCache.Limit(100)
				// ch <- depthCache
			}
		}
		if !ok {
			logger.Logger.Error("cannot get CheckDepthCacheChanMap err, symbol:", depthCache.Symbol)
		}
	}
}

func (b *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte) error {
	receiveTime := time.Now().UnixMicro()
	b.Lock.Lock()
	var (
		resp       RespIncrementDepthStream
		deltaDepth = &base.DeltaDepthUpdate{}
		ok         bool
		content    interface{}
		chUpdate   chan *base.DeltaDepthUpdate
		err        error
	)
	defer b.Lock.Unlock()
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	if len(resp.Data.B) == 0 && len(resp.Data.A) == 0 {
		return nil
	}
	transferDiffDepth(&resp.Data, deltaDepth, b.MarketType)
	deltaDepth.TimeReceive = receiveTime
	_, sym := GetSymbolKey(resp.Data.S, b.MarketType)

	if b.IsPublishDelta {
		if _, ok := b.depthIncrementSnapshotDeltaGroupChanMap[sym]; ok {
			base.SendChan(b.depthIncrementSnapshotDeltaGroupChanMap[sym], deltaDepth.Transfer2Depth(), "depthIncrementSnapshotDeltaGroupChanMap")
		} else {
			logger.Logger.Warn("[delta depth] get symbol from channel map err:", sym)
		}
	}
	if b.IsPublishFull {
		if content, ok = b.DepthDeltaUpdateMap.Get(sym); ok {
			chUpdate, ok = content.(chan *base.DeltaDepthUpdate)
		}
		base.SendChan(chUpdate, deltaDepth, "DeltaDepthUpdate", deltaDepth.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(sym, time.Now().UnixMicro())
	return nil
}

var checkGoroutineCount int64

func (b *WebSocketSpotHandle) checkFullDepth(symbol *client.SymbolInfo) {
	/**
	rest api返回1000条数据，但是由于程序开始也是1000条，中间数据抛弃后，会有新的数据在rest api中与增量数据不同，所以不能对比所有，要限定数量
	*/
	var (
		sym            = SymbolKeyGen(symbol.Symbol, symbol.Market, symbol.Type)
		depthCacheList []base.OrderBook
		fullDepth      *base.OrderBook
		ch             chan *base.OrderBook
		content        interface{}
		ok             bool
		err            error
	)
	b.CheckStates.Set(sym, true)
	defer func() { b.CheckStates.Set(sym, false) }()
	fullDepth = &base.OrderBook{}
	if content, ok = b.CheckDepthCacheChanMap.Get(sym); ok {
		ch, ok = content.(chan *base.OrderBook)
	}
	if !ok {
		logger.Logger.Error("get CheckDepthCacheChanMap err, symbol:", symbol)
		return
	}
	checkGoroutineCount++
	start := time.Now()
	defer func(start time.Time) {
		checkGoroutineCount--
		fmt.Println(symbol.Symbol, " check time use:", time.Now().Sub(start), checkGoroutineCount)
	}(start)

MAINLOOP:
	for {
		select {
		case dep := <-ch:
			b.Lock.Lock()
			depthCacheList = append(depthCacheList, *dep)
			b.Lock.Unlock()
			if len(depthCacheList) > b.DepthCapLevel {
				depthCacheList = depthCacheList[len(depthCacheList)-b.DepthCapLevel:]
			}
			if fullDepth == nil || fullDepth.UpdateId == 0 {
				if fullDepth, err = b.GetFullDepth(GetSymbolName(symbol)); err != nil {
					logger.Logger.Error("checkFullDepth get full depth error", err)
					time.Sleep(time.Second * 30)
					continue
				}
			}
			for _, depCache := range depthCacheList {
				if fullDepth == nil {
					logger.Logger.Warnf(symbol.Symbol+"checkfull depth address: %p,%p", depCache, fullDepth)
					continue MAINLOOP
				}
				if depCache.UpdateId == fullDepth.UpdateId {
					depthCacheList = []base.OrderBook{}
					if depCache.Equal(fullDepth, b.DepthCheckLevel) {
						break MAINLOOP
					} else {
						logger.Logger.Error(symbol.Symbol, "depth level:", b.DepthCheckLevel, "need reset")
						logger.SaveToFile("depth_error.json", "depCache", depCache, "fullDepth", fullDepth)
						errCh := b.DepthNotMatchChanMap[b.symbolMap[sym]]
						errCh <- false
						break MAINLOOP
					}
				} else if depCache.UpdateId < fullDepth.UpdateId {
					continue
				} else {
					if fullDepth, err = b.GetFullDepth(GetSymbolName(symbol)); err != nil {
						//logger.Logger.Error(symbol, "checkFullDepth < get full depth error", err)
						continue
					}
				}
			}
		}
	}
}

func (b *WebSocketSpotHandle) Check() {

	/**
	定期检查全量数据
	*/
	defer func() {
		if err := recover(); err != nil {
			logger.Logger.Error("incr depth check panic:", logger.PanicTrace(err))
			time.Sleep(time.Second)
			go b.Check()
		}
	}()
	heartTimer := time.NewTimer(time.Duration(b.CheckTimeSec) * time.Second)
LOOP:
	for {
		select {
		case <-heartTimer.C:
			for _, symbol := range b.DepthCacheMap.Keys() {
				if _, ok := b.DepthCacheMap.Get(symbol); ok {
					sym := b.symbolMap[symbol]
					b.checkFullDepth(sym)
				} else {
					logger.Logger.Error("get symbol from DepthCacheMap err")
				}
				time.Sleep(time.Second * 10)
			}
			heartTimer.Reset(time.Duration(b.CheckTimeSec) * time.Second)
		case <-b.Ctx.Done():
			b.Reset()
			break LOOP
		}
	}
}

func (b *WebSocketSpotHandle) Reset() {
	b.CheckDepthCacheChanMap = sdk.NewCmapI()
	//b.DepthCacheListMap = sdk.NewCmapI()
	b.DepthCacheMap = sdk.NewCmapI()
	b.DepthDeltaUpdateMap = sdk.NewCmapI()
}
