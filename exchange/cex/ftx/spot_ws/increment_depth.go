package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"clients/transform"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/sdk"
)

func (b *WebSocketSpotHandle) SetDepthIncrementSnapShotConf(symbols []*client.SymbolInfo, conf *base.IncrementDepthConf) {
	if conf.DepthCapLevel < 0 {
		conf.DepthCapLevel = 1000
	}
	if conf.DepthCapLevel < 0 {
		conf.DepthCapLevel = 20
	}
	if conf.CheckTimeSec < 0 {
		conf.CheckTimeSec = 3600
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
			//Orderbook           *base.OrderBook //Modified by me (checkpoint)
		)
		SymbolInfoStr := base.SymInfoToString(symbol)
		conf.DepthCacheMap.Set(SymbolInfoStr, nil) //Modified by me (checkpoint)
		conf.DepthCacheListMap.Set(SymbolInfoStr, DepthCacheList)
		conf.CheckDepthCacheChanMap.Set(SymbolInfoStr, CheckDepthCacheChan)
		conf.DepthNotMatchChanMap[symbol] = make(chan bool, conf.DepthCapLevel)
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

func ParseOrder(orders [][]float64, slice *base.DepthItemSlice) {
	for _, order := range orders {
		*slice = append(*slice, &depth.DepthLevel{
			Price:  order[0],
			Amount: order[1],
		})
	}
}

func transferDiffDepth(resp *OrderbooksResponse, diff *base.DeltaDepthUpdate) {
	diff.Symbol = resp.Market
	diff.Market = getMarket(resp.Market)
	diff.Type = getSymbolType(resp.Market)
	diff.TimeExchange = int64(resp.Data.Time * float64(time.Millisecond))
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

func (b *WebSocketSpotHandle) update(deltaDepthCacheList []*base.DeltaDepthUpdate, diffDepth *base.DeltaDepthUpdate, depthCache *base.OrderBook, res *depth.Depth) (deltaList []*base.DeltaDepthUpdate, err error) {
	// 缓存depth的增量数据
	deltaList = append(deltaDepthCacheList, diffDepth)
	if len(deltaList) > b.DepthCapLevel {
		deltaList = deltaList[len(deltaList)-b.DepthCapLevel:]
	}
	deltaList = append([]*base.DeltaDepthUpdate{}) // clear delta list ???
	base.UpdateBidsAndAsks(diffDepth, depthCache, b.DepthCapLevel, res)

	return
}

func (b *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte) error {
	b.Lock.Lock()
	var (
		res                 = &depth.Depth{}
		resp                OrderbooksResponse
		deltaDepth          = &base.DeltaDepthUpdate{}
		deltaDepthCacheList []*base.DeltaDepthUpdate
		depthCache          *base.OrderBook
		t                   = time.Now().UnixMicro()
		ok                  bool
		content             interface{}
		err                 error
		isFirstSnapshot     = false
	)
	defer b.Lock.Unlock()
	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}

	if resp.Type == "error" {
		logger.Logger.Error("data error with code:", resp.Code, resp.Message)
		return nil
	} else if resp.Type == "pong" {
		return nil
	} else if resp.Type == "subscribed" {
		logger.Logger.Info("subscribe successful to:", resp.Market)
		return nil
	} else if resp.Type != "partial" && resp.Type != "update" {
		logger.Logger.Error(" error different type:", string(data))
		return nil
	}
	if len(resp.Data.Asks) == 0 && len(resp.Data.Bids) == 0 {
		return nil
	}
	if resp.Type == "partial" {

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

		depthCache = &base.OrderBook{
			Exchange:     common.Exchange_FTX,
			Market:       getMarket(resp.Market),
			Type:         getSymbolType(resp.Market),
			Symbol:       resp.Market,
			TimeExchange: uint64(resp.Data.Time * float64(time.Millisecond)),
			TimeReceive:  uint64(t),
			Asks:         asks,
			Bids:         bids,
		}
		isFirstSnapshot = true
	}
	SymbolInfoStr := base.SymInfoToString(getSymbolInfo(resp.Market))
	transferDiffDepth(&resp, deltaDepth)
	if b.IsPublishDelta {
		if _, ok := b.depthIncrementSnapshotDeltaGroupChanMap[SymbolInfoStr]; ok {
			base.SendChan(b.depthIncrementSnapshotDeltaGroupChanMap[SymbolInfoStr], deltaDepth.Transfer2Depth(), "depthIncrementSnapshotDeltaGroupChanMap")
		} else {
			logger.Logger.Warn("[depth_delta] get symbol from channel map err:", res.Symbol)
		}
	}
	if b.IsPublishFull {
		// get increment depth list
		if content, ok = b.DepthCacheListMap.Get(SymbolInfoStr); ok {
			deltaDepthCacheList, ok = content.([]*base.DeltaDepthUpdate)
		}
		if !ok {
			return errors.New(deltaDepth.Symbol + "get deltaDepthCacheList error")
		}
		// get cached snapshot depth
		if depthCache == nil {
			if content, ok = b.DepthCacheMap.Get(SymbolInfoStr); ok {
				depthCache, ok = content.(*base.OrderBook)
			}
		}

		if deltaDepthCacheList, err = b.update(deltaDepthCacheList, deltaDepth, depthCache, nil); err != nil { // merge this increment to cached snapshot
			if err.Error() == "need wait a while" {
				b.DepthCacheListMap.Set(SymbolInfoStr, deltaDepthCacheList)
				b.DepthCacheMap.Set(SymbolInfoStr, depthCache)
				return nil
			}
			logger.Logger.Error(deltaDepth.Symbol+"update depth err:", err)
		}
		for symbol, _ := range b.DepthNotMatchChanMap {
			if base.SymInfoToString(symbol) == SymbolInfoStr {
				// 防止阻塞
				select {
				case b.DepthNotMatchChanMap[symbol] <- !(resp.Data.Checksum == GetChecksum(depthCache)): // checksum
				default:
				}
			}
		}
		b.DepthCacheListMap.Set(SymbolInfoStr, deltaDepthCacheList)
		b.DepthCacheMap.Set(SymbolInfoStr, depthCache)
		depthCache.Transfer2Depth(b.DepthLevel, res)
		if _, ok := b.depthIncrementSnapshotFullGroupChanMap[SymbolInfoStr]; ok {
			if isFirstSnapshot {
				res.Hdr = base.MakeFirstDepthHdr()
			}
			base.SendChan(b.depthIncrementSnapshotFullGroupChanMap[SymbolInfoStr], res, "b.depthIncrementSnapshotFullGroupChanMap", resp.Market)
		} else {
			logger.Logger.Warn("[depth_full] get symbol from channel map err:", res.Symbol)
		}
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(deltaDepth.Symbol, time.Now().UnixMicro())
	return nil
}

func getSymbolInfo(symbol string) *client.SymbolInfo {
	return &client.SymbolInfo{Symbol: strings.Split(symbol, "-")[0], Market: getMarket(symbol), Type: getSymbolType(symbol)}
}

func getMarket(symbol string) common.Market {
	if strings.Contains(symbol, "-") {
		if strings.Contains(symbol, "PERP") {
			return common.Market_SWAP
		} else {
			return common.Market_FUTURE
		}
	} else {
		return common.Market_SPOT
	}
}

func getSymbolType(symbol string) common.SymbolType {
	if strings.Contains(symbol, "-") {
		if strings.Contains(symbol, "PERP") {
			return common.SymbolType_SWAP_FOREVER
		} else if strings.Split(symbol, "-")[1] == transform.GetThisQuarter(time.Now().UTC(), 5, 2).Format("0102") {
			return common.SymbolType_FUTURE_THIS_QUARTER
		} else {
			return common.SymbolType_FUTURE_NEXT_QUARTER
		}
	} else {
		return common.SymbolType_SPOT_NORMAL
	}
}
