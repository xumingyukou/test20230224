package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/sdk"
)

//func ParseOrder(orders [][]float64, slice *base.DepthItemSlice) {
//	for _, order := range orders {
//		*slice = append(*slice, &depth.DepthLevel{
//			Price:  order[0],
//			Amount: order[1],
//		})
//	}
//}

func (b *WebSocketSpotHandle) SetDepthIncrementSnapShotConf(symbols []*client.SymbolInfo, conf *base.IncrementDepthConf) {
	if conf.DepthCapLevel <= 0 {
		conf.DepthCapLevel = 1000
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
		var (
			CheckDepthCacheChan = make(chan *base.OrderBook, 100)
		)
		conf.DepthCacheMap.Set(base.SymInfoToString(symbol), nil)
		conf.DepthDeltaUpdateMap.Set(base.SymInfoToString(symbol), make(chan *base.DeltaDepthUpdate, 50))
		conf.CheckDepthCacheChanMap.Set(base.SymInfoToString(symbol), CheckDepthCacheChan)
		conf.DepthNotMatchChanMap[symbol] = make(chan bool, conf.DepthCapLevel)
		conf.CheckStates.Set(base.SymInfoToString(symbol), false)
		b.symbolMap[base.SymInfoToString(symbol)] = symbol
	}
	for _, symbol := range symbols {
		go b.UpdateDeltaDepth(symbol)
	}
}

func (b *WebSocketSpotHandle) SetDepthIncrementSnapshotReqChan(ch chan *client.SymbolInfo) {
	b.DepthIncrementSnapshotReqSymbol = ch
}

func (b *WebSocketSpotHandle) Reset() {
	b.CheckDepthCacheChanMap = sdk.NewCmapI()
	b.DepthCacheMap = sdk.NewCmapI()
	b.DepthDeltaUpdateMap = sdk.NewCmapI()
}

func (b *WebSocketSpotHandle) UpdateDeltaDepth(symbol *client.SymbolInfo) {
	var (
		deltaCh     chan *base.DeltaDepthUpdate
		fullCh      chan *base.OrderBook
		deltaSendCh chan *client.WsDepthRsp
		fullSendCh  chan *depth.Depth

		ok  bool
		obj interface{}
		//err                 error
		funcName = "updateDeltaDepth"
	)

	if obj, ok = b.DepthDeltaUpdateMap.Get(base.SymInfoToString(symbol)); ok {
		deltaCh, ok = obj.(chan *base.DeltaDepthUpdate)
	}
	if !ok {
		logger.Logger.Error("get DepthDeltaUpdateMap ", funcName, " ", base.SymInfoToString(symbol))
		return
	}

	if obj, ok = b.CheckDepthCacheChanMap.Get(base.SymInfoToString(symbol)); ok {
		fullCh, ok = obj.(chan *base.OrderBook)
	}
	if !ok {
		logger.Logger.Error("get DepthCacheChanMap ", funcName, " ", base.SymInfoToString(symbol))
		return
	}
	if b.IsPublishDelta {
		deltaSendCh, ok = b.depthIncrementSnapshotDeltaGroupChanMap[base.SymInfoToString(symbol)]
		if !ok {
			logger.Logger.Error("get delta send channel ", funcName, " ", symbol.Symbol)
			return
		}
	}
	if b.IsPublishFull {
		fullSendCh, ok = b.depthIncrementSnapshotFullGroupChanMap[base.SymInfoToString(symbol)]
		if !ok {
			logger.Logger.Error("get full send channel ", funcName, " ", symbol.Symbol)
			return
		}
	}
	if symbol.Market == common.Market_SPOT {
		go b.ProcessSpotDeltaDepth(symbol, deltaCh, fullCh, deltaSendCh, fullSendCh)
	} else {
		go b.ProcessFutureOrSwapDeltaDepth(symbol, deltaCh, fullCh, deltaSendCh, fullSendCh)
	}
	return
}

func (b *WebSocketSpotHandle) ProcessSpotDeltaDepth(symbol *client.SymbolInfo, deltaCh chan *base.DeltaDepthUpdate,
	fullCh chan *base.OrderBook, deltaSendCh chan *client.WsDepthRsp, fullSendCh chan *depth.Depth) {
	var (
		depthCache        = &base.OrderBook{}
		deltaCacheList    = make([]*base.DeltaDepthUpdate, 0, 10)
		firstSnapshotSent = false
	)

	lastReqTime := time.Now().UnixMicro()
LOOP:
	for {
		select {
		case <-b.Ctx.Done():
			break LOOP
		case delta, ok := <-deltaCh:
			if !ok {
				continue
			}
			if b.IsPublishDelta {
				res := delta.Transfer2Depth()
				base.SendChan(deltaSendCh, res, "deltaPublish", base.SymInfoToString(symbol))
			}
			if b.IsPublishFull {
				deltaCacheList = append(deltaCacheList, delta)
				if depthCache.UpdateId == 0 || depthCache.UpdateId < deltaCacheList[0].UpdateStartId {
					if time.Now().Sub(time.UnixMicro(lastReqTime)) > time.Duration(5)*time.Second {
						logger.Logger.Warnf("missing delta %s %d %d", base.SymInfoToString(symbol), depthCache.UpdateId, deltaCacheList[0].UpdateStartId)
						// 数据丢失或者全量未请求，请求全量
						b.DepthIncrementSnapshotReqSymbol <- symbol
						firstSnapshotSent = false
						deltaCacheList = append([]*base.DeltaDepthUpdate{})
						lastReqTime = time.Now().UnixMicro()
					}
					continue
				} else if depthCache.UpdateId >= deltaCacheList[len(deltaCacheList)-1].UpdateEndId {
					// 等待增量
					logger.Logger.Warnf("waiting delta %s %d %d", base.SymInfoToString(symbol), depthCache.UpdateId, deltaCacheList[len(deltaCacheList)-1].UpdateEndId)
					firstSnapshotSent = false
					deltaCacheList = append([]*base.DeltaDepthUpdate{})
					continue
				} else {
					target := -1
					for i, delta := range deltaCacheList {
						if depthCache.UpdateId == delta.UpdateStartId {
							target = i
						}
					}
					if target == -1 {
						// 没有找到全量的id
						// 数据丢失，重新订阅
						firstSnapshotSent = false
						if time.Now().Sub(time.UnixMicro(lastReqTime)) > time.Duration(5)*time.Second {
							logger.Logger.Error("deltaList error", base.SymInfoToString(symbol))
							b.DepthIncrementSnapshotReqSymbol <- symbol
							deltaCacheList = append([]*base.DeltaDepthUpdate{})
							lastReqTime = time.Now().UnixMicro()
						}
						continue
					}
					deltaCacheList = deltaCacheList[target:]
					dataLost := false
					for _, deltaCache := range deltaCacheList {
						if depthCache.UpdateId != delta.UpdateStartId {
							// 数据丢失，重新订阅
							firstSnapshotSent = false
							dataLost = true
							if time.Now().Sub(time.UnixMicro(lastReqTime)) > time.Duration(5)*time.Second {
								b.DepthIncrementSnapshotReqSymbol <- symbol
								lastReqTime = time.Now().UnixMicro()
							}
							break
						}
						base.UpdateBidsAndAsks(deltaCache, depthCache, b.DepthCapLevel, nil)
					}
					if !dataLost {
						fullToSend := &depth.Depth{}
						depthCache.Transfer2Depth(b.DepthLevel, fullToSend)
						if !firstSnapshotSent {
							fullToSend.Hdr = base.MakeFirstDepthHdr()
							firstSnapshotSent = true
						}
						base.SendChan(fullSendCh, fullToSend, "depth full send", base.SymInfoToString(symbol))
					}
					deltaCacheList = append([]*base.DeltaDepthUpdate{})
				}
			}
		case full, ok := <-fullCh:
			if !ok {
				continue
			}
			logger.Logger.Info("receive full ", full.UpdateId, " ", base.SymInfoToString(symbol))
			firstSnapshotSent = false
			depthCache = full
		}
	}
}

func (b *WebSocketSpotHandle) ProcessFutureOrSwapDeltaDepth(symbol *client.SymbolInfo,
	deltaCh chan *base.DeltaDepthUpdate, fullCh chan *base.OrderBook, deltaSendCh chan *client.WsDepthRsp,
	fullSendCh chan *depth.Depth) {

	var (
		depthCache        = &base.OrderBook{}
		deltaCacheList    = make([]*base.DeltaDepthUpdate, 0, 10)
		firstSnapshotSent = false
	)

	lastReqTime := time.Now().UnixMicro()
LOOP:
	for {
		select {
		case <-b.Ctx.Done():
			break LOOP
		case delta, ok := <-deltaCh:
			if !ok {
				continue
			}
			if b.IsPublishDelta {
				res := delta.Transfer2Depth()
				base.SendChan(deltaSendCh, res, "deltaPublish", base.SymInfoToString(symbol))
			}
			if b.IsPublishFull {
				deltaCacheList = append(deltaCacheList, delta)
				if depthCache.UpdateId == 0 || depthCache.UpdateId+1 < deltaCacheList[0].UpdateEndId {
					if time.Now().Sub(time.UnixMicro(lastReqTime)) > time.Duration(5)*time.Second {
						logger.Logger.Warnf("missing delta %s %d %d", base.SymInfoToString(symbol), depthCache.UpdateId, deltaCacheList[0].UpdateStartId)
						// 数据丢失或者全量未请求，请求全量
						b.DepthIncrementSnapshotReqSymbol <- symbol
						firstSnapshotSent = false
						deltaCacheList = append([]*base.DeltaDepthUpdate{})
						lastReqTime = time.Now().UnixMicro()
					}
					continue
				} else if depthCache.UpdateId+1 > deltaCacheList[len(deltaCacheList)-1].UpdateEndId {
					// 等待增量
					logger.Logger.Warnf("waiting delta %s %d %d", base.SymInfoToString(symbol), depthCache.UpdateId, deltaCacheList[len(deltaCacheList)-1].UpdateEndId)
					firstSnapshotSent = false
					deltaCacheList = append([]*base.DeltaDepthUpdate{})
					continue
				} else {
					target := -1
					for i, delta := range deltaCacheList {
						if depthCache.UpdateId+1 == delta.UpdateEndId {
							target = i
						}
					}
					if target == -1 {
						// 没有找到全量的id
						// 数据丢失，重新订阅
						firstSnapshotSent = false
						if time.Now().Sub(time.UnixMicro(lastReqTime)) > time.Duration(5)*time.Second {
							logger.Logger.Error("deltaList error", base.SymInfoToString(symbol))
							b.DepthIncrementSnapshotReqSymbol <- symbol
							deltaCacheList = append([]*base.DeltaDepthUpdate{})
							lastReqTime = time.Now().UnixMicro()
						}
						continue
					}
					deltaCacheList = deltaCacheList[target:]
					for _, deltaCache := range deltaCacheList {
						base.UpdateBidsAndAsks(deltaCache, depthCache, b.DepthCapLevel, nil)
					}
					fullToSend := &depth.Depth{}
					depthCache.Transfer2Depth(b.DepthLevel, fullToSend)
					if !firstSnapshotSent {
						fullToSend.Hdr = base.MakeFirstDepthHdr()
						firstSnapshotSent = true
					}
					base.SendChan(fullSendCh, fullToSend, "depth full send", base.SymInfoToString(symbol))

					deltaCacheList = append([]*base.DeltaDepthUpdate{})
				}
			}
		case full, ok := <-fullCh:
			if !ok {
				continue
			}
			logger.Logger.Info("receive full ", full.UpdateId, " ", base.SymInfoToString(symbol))
			firstSnapshotSent = false
			depthCache = full
		}
	}
}

func (b *WebSocketSpotHandle) AfterUpdateHandle(sym string, depthCache *base.OrderBook) {
	var (
		res = &depth.Depth{}
		// content         interface{}
		// ok, checkStatus bool
	)
	//if content, ok = b.CheckStates.Get(sym); ok {
	//	checkStatus, ok = content.(bool)
	//}
	depthCache.Transfer2Depth(b.DepthLevel, res)
	if _, ok := b.depthIncrementSnapshotFullGroupChanMap[sym]; ok {
		base.SendChan(b.depthIncrementSnapshotFullGroupChanMap[sym], res, "DepthIncrementSnapShotGroupHandle")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	//if ok && checkStatus {
	//	var (
	//		ch chan *base.OrderBook
	//	)
	//	if content, ok = b.CheckDepthCacheChanMap.Get(spot_api.GetSymbol(depthCache.Symbol)); ok {
	//		if ch, ok = content.(chan *base.OrderBook); ok {
	//			depthCache.Limit(100)
	//			ch <- depthCache
	//		}
	//	}
	//	if !ok {
	//		logger.Logger.Error("cannot get CheckDepthCacheChanMap err, symbol:", depthCache.Symbol)
	//	}
	//}
}

func transferDiffDepth(resp *base.OrderBook, diff *base.DeltaDepthUpdate) {
	diff.Symbol = resp.Symbol
	diff.Market = resp.Market
	diff.Type = resp.Type
	diff.TimeExchange = int64(resp.TimeExchange)
	diff.TimeReceive = int64(resp.TimeReceive)
	//清空bids，asks
	diff.Bids = diff.Bids[:0]
	diff.Asks = diff.Asks[:0]
	diff.Bids = resp.Bids
	diff.Asks = resp.Asks
	sort.Stable(diff.Asks)
	sort.Stable(sort.Reverse(diff.Bids))
	return
}

func (b *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte) error {
	//var str bytes.Buffer
	//_ = json.Indent(&str, []byte(data), "", "    ")
	//fmt.Println("Receive: ", str.String())
	t := time.Now().UnixMicro()
	b.Lock.Lock()
	var (
		respTemp      RespSnapShotTemp
		respIncrement RespIncrementSnapShotStream
		respLimit     RespLimitSnapShotStream
		deltaCh       chan *base.DeltaDepthUpdate
		fullCh        chan *base.OrderBook
		asks          []*depth.DepthLevel
		bids          []*depth.DepthLevel
		ok            bool
		obj           interface{}
		err           error
	)
	defer b.Lock.Unlock()
	if err = json.Unmarshal(data, &respTemp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if respTemp.Ping != 0 {
		b.pingpong <- respTemp.Ping
		return nil
	}

	if respTemp.Rep != "" {
		// 全量信息
		if err = json.Unmarshal(data, &respLimit); err != nil {
			logger.Logger.Error("receive data err:", string(data))
			return err
		}
		if respLimit.Data.SeqNum != 0 {
			asks, err = DepthLevelParse(respLimit.Data.Asks)
			if err != nil {
				logger.Logger.Error("book ticker parse err", err, string(data))
				return err
			}
			bids, err = DepthLevelParse(respLimit.Data.Bids)
			if err != nil {
				logger.Logger.Error("book ticker parse err", err, string(data))
				return err
			}
		} else {
			asks = nil
			bids = nil
		}

		var symbolName string
		symIf, okk := symbolNameMap.Load(strings.Split(respLimit.Rep, ".")[1])
		if okk {
			symbolName = symIf.(string)
		}
		full := &base.OrderBook{
			Exchange:     common.Exchange_HUOBI,
			Market:       b.Market,
			Type:         b.Type,
			Symbol:       symbolName,
			TimeExchange: uint64(respLimit.Ts * 1000),
			TimeReceive:  uint64(t),
			Asks:         asks,
			Bids:         bids,
			UpdateId:     respLimit.Data.SeqNum,
		}
		if obj, ok = b.CheckDepthCacheChanMap.Get(full.Symbol); ok {
			if fullCh, ok = obj.(chan *base.OrderBook); ok {
				base.SendChan(fullCh, full, "depthCacheChan", full.Symbol)
			}
		}
		if !ok {
			logger.Logger.Error("get DepthCacheChanMap err, symbol:", full.Symbol)
			return errors.New("get DepthCacheChanMap err, symbol:" + full.Symbol)
		}
		b.CheckSendStatus.CheckUpdateTimeMap.Set(full.Symbol, t)
		return nil
		// deltaDepth.Symbol = resp.Symbol
		// deltaDepth.Market = resp.Market
		// deltaDepth.Type = resp.Type
		// deltaDepth.Asks = resp.Asks
		// deltaDepth.Bids = resp.Bids
		// deltaDepth.TimeReceive = t
		// if content, ok = b.DepthDeltaUpdateMap.Get(deltaDepth.Symbol); ok {
		// 	chUpdate, ok = content.(chan *base.DeltaDepthUpdate)
		// }
		// base.SendChan(chUpdate, deltaDepth, "Initialization", deltaDepth.Symbol)
		// return nil
	}
	if respTemp.Subbed != "" {
		logger.Logger.Info("Subscribed(SUB) Success:", string(data))
		return nil
	}
	if respTemp.Ch != "" {
		if err = json.Unmarshal(data, &respIncrement); err != nil {
			logger.Logger.Error("receive data err:", string(data))
			return err
		}
		asks, err = DepthLevelParse(respIncrement.Tick.Asks)
		if err != nil {
			logger.Logger.Error("book ticker parse err", err, string(data))
			return err
		}
		bids, err = DepthLevelParse(respIncrement.Tick.Bids)
		if err != nil {
			logger.Logger.Error("book ticker parse err", err, string(data))
			return err
		}

		var symbolName string
		symIf, ok := symbolNameMap.Load(strings.Split(respIncrement.Ch, ".")[1])
		if ok {
			symbolName = symIf.(string)
		}
		delta := &base.DeltaDepthUpdate{
			UpdateStartId: respIncrement.Tick.PrevSeqNum,
			UpdateEndId:   respIncrement.Tick.SeqNum,

			Market:       b.Market,
			Type:         b.Type,
			Symbol:       symbolName,
			TimeExchange: respIncrement.Ts * 1000,
			TimeReceive:  t,
			Asks:         asks,
			Bids:         bids,
		}
		symbol := delta.Symbol
		if obj, ok = b.DepthDeltaUpdateMap.Get(symbol); ok {
			if deltaCh, ok = obj.(chan *base.DeltaDepthUpdate); ok {
				base.SendChan(deltaCh, delta, "deltaDepthUpdate", symbol)
			}
		}

		if !ok {
			logger.Logger.Error("get DepthDeltaUpdateMap err, symbol:", delta.Symbol)
			return errors.New("get DepthDeltaUpdateMap err, symbol:" + delta.Symbol)
		}
		b.CheckSendStatus.CheckUpdateTimeMap.Set(delta.Symbol, t)

		// if b.IsPublishDelta {
		// 	if _, ok = b.depthIncrementSnapshotDeltaGroupChanMap[sym]; ok {
		// 		base.SendChan(b.depthIncrementSnapshotDeltaGroupChanMap[sym], deltaDepth.Transfer2Depth(), "depthIncrementSnapshotDeltaGroupChanMap")
		// 	} else {
		// 		logger.Logger.Warn("get symbol from channel map err:", sym)
		// 	}
		// }
		// if b.IsPublishFull {
		// 	if content, ok = b.DepthDeltaUpdateMap.Get(sym); ok {
		// 		chUpdate, ok = content.(chan *base.DeltaDepthUpdate)
		// 	}
		// 	base.SendChan(chUpdate, deltaDepth, "DeltaDepthUpdate", deltaDepth.Symbol)
		// }
		// b.CheckSendStatus.CheckUpdateTimeMap.Set(deltaDepth.Symbol, time.Now().UnixMicro())
		return nil
	}
	//fmt.Println("Pass Msg!!!!!!")
	return nil
}

func (b *WebSocketSpotHandle) FutureOrSwapDepthIncrementGroupHandle(data []byte) error {
	t := time.Now().UnixMicro()
	b.Lock.Lock()
	var (
		respTemp RespFutureOrSwapIncrementSnapShotStream
		deltaCh  chan *base.DeltaDepthUpdate
		fullCh   chan *base.OrderBook
		asks     []*depth.DepthLevel
		bids     []*depth.DepthLevel
		ok       bool
		obj      interface{}
		err      error
	)
	defer b.Lock.Unlock()
	if err = json.Unmarshal(data, &respTemp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if respTemp.Ping != 0 {
		b.pingpong <- respTemp.Ping
		return nil
	}
	if respTemp.Subbed != "" {
		logger.Logger.Info("Subscribed(SUB) Success:", string(data))
		return nil
	}
	if respTemp.Ch == "" {
		logger.Logger.Warn("resp ch fault: ", string(data))
		return nil
	}

	asks, err = DepthLevelParse(respTemp.Tick.Asks)
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	bids, err = DepthLevelParse(respTemp.Tick.Bids)
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	var symbolName string
	chName := strings.Split(respTemp.Ch, ".")[1]
	symIf, okk := symbolNameMap.Load(chName)
	if okk {
		symbolName = symIf.(string)
	}
	if respTemp.Tick.Event == "snapshot" {
		// 全量信息
		full := &base.OrderBook{
			Exchange:     b.Exchange,
			Market:       b.Market,
			Type:         b.Type,
			Symbol:       symbolName,
			TimeExchange: uint64(respTemp.Ts * 1000),
			TimeReceive:  uint64(t),
			Asks:         asks,
			Bids:         bids,
			UpdateId:     respTemp.Tick.Version,
		}
		if obj, ok = b.CheckDepthCacheChanMap.Get(full.Symbol); ok {
			if fullCh, ok = obj.(chan *base.OrderBook); ok {
				// 张数转换
				TransContractSize(chName, b.Market, nil, full.Asks...)
				TransContractSize(chName, b.Market, nil, full.Bids...)
				base.SendChan(fullCh, full, "depthCacheChan", full.Symbol)
			}
		}
		if !ok {
			logger.Logger.Error("get DepthCacheChanMap err, symbol:", full.Symbol)
			return errors.New("get DepthCacheChanMap err, symbol:" + full.Symbol)
		}
		b.CheckSendStatus.CheckUpdateTimeMap.Set(full.Symbol, t)
		return nil
	}
	if respTemp.Tick.Event == "incremental" || respTemp.Tick.Event == "update" {
		// 增量信息
		delta := &base.DeltaDepthUpdate{
			UpdateEndId: respTemp.Tick.Version,

			Market:       b.Market,
			Type:         b.Type,
			Symbol:       symbolName,
			TimeExchange: respTemp.Ts * 1000,
			TimeReceive:  t,
			Asks:         asks,
			Bids:         bids,
		}
		symbol := delta.Symbol
		if obj, ok = b.DepthDeltaUpdateMap.Get(symbol); ok {
			if deltaCh, ok = obj.(chan *base.DeltaDepthUpdate); ok {
				// 张数转换
				TransContractSize(chName, b.Market, nil, delta.Asks...)
				TransContractSize(chName, b.Market, nil, delta.Bids...)
				base.SendChan(deltaCh, delta, "deltaDepthUpdate", symbol)
			}
		}

		if !ok {
			logger.Logger.Error("get DepthDeltaUpdateMap err, symbol:", delta.Symbol)
			return errors.New("get DepthDeltaUpdateMap err, symbol:" + delta.Symbol)
		}
		b.CheckSendStatus.CheckUpdateTimeMap.Set(delta.Symbol, t)

		return nil
	}
	return nil
}
