package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"sort"
	"time"

	"github.com/goccy/go-json"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/sdk"
)

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
		conf.DepthCacheMap.Set(symbol.Symbol, nil)
		conf.DepthDeltaUpdateMap.Set(symbol.Symbol, make(chan *base.DeltaDepthUpdate, 50))
		conf.CheckDepthCacheChanMap.Set(symbol.Symbol, CheckDepthCacheChan)
		conf.DepthNotMatchChanMap[symbol] = make(chan bool, conf.DepthCapLevel)
		conf.CheckStates.Set(symbol.Symbol, false)
		b.symbolMap[symbol.Symbol] = symbol
	}
	for _, symbol := range symbols {
		go b.UpdateDeltaDepth(symbol.Symbol)
	}
}

func (b *WebSocketSpotHandle) UpdateDeltaDepth(symbol string) {
	var (
		content    interface{}
		ok         bool
		ch         chan *base.DeltaDepthUpdate
		diffDepth  *base.DeltaDepthUpdate
		depthCache = &base.OrderBook{Symbol: symbol}
		//err                 error
		firstSnapshotSent = false
		fullCh            chan *depth.Depth
	)
	if content, ok = b.DepthDeltaUpdateMap.Get(symbol); ok {
		ch, ok = content.(chan *base.DeltaDepthUpdate)
	}
	if !ok {
		logger.Logger.Error("get DepthDeltaUpdateMap err:", symbol)
		return
	}

	if fullCh, ok = b.depthIncrementSnapshotFullGroupChanMap[symbol]; !ok {
		logger.Logger.Warn("get symbol from channel map err:", symbol)
	}
LOOP:
	for {
		select {
		case diffDepth, ok = <-ch:
			if !ok {
				logger.Logger.Error("DepthDeltaUpdateMap channel close:", symbol)
				break LOOP
			}
			// 如果来的是全量，需要把depthCache初始化，避免新全量去merge老全量
			if diffDepth.IsFullDepth {
				depthCache = &base.OrderBook{Symbol: symbol}
				firstSnapshotSent = false
			}
			base.UpdateBidsAndAsks(diffDepth, depthCache, b.DepthCapLevel, nil)

			full := &depth.Depth{}
			depthCache.Transfer2Depth(b.DepthLevel, full)
			if !firstSnapshotSent {
				full.Hdr = base.MakeFirstDepthHdr()
				firstSnapshotSent = true
			}
			base.SendChan(fullCh, full, "DepthIncrementSnapShotGroupHandle")
		case <-b.Ctx.Done():
			break LOOP
		}

	}
}

func (b *WebSocketSpotHandle) AfterUpdateHandle(sym string, depthCache *base.OrderBook) {
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

func (b *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte, t int64, isFirstMsg bool) error {
	var (
		respL2Updates L2Update
		resp          base.OrderBook
		err           error
		deltaDepth    = &base.DeltaDepthUpdate{}
		chUpdate      chan *base.DeltaDepthUpdate
		ok            bool
		content       interface{}
		asks          []*depth.DepthLevel
		bids          []*depth.DepthLevel
	)
	if err = json.Unmarshal(data, &respL2Updates); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	for _, change := range respL2Updates.Changes {
		if change[0] == "sell" {
			asks = append(asks, &depth.DepthLevel{
				Price:  StringToFloat(change[1]),
				Amount: StringToFloat(change[2]),
			})
		} else {
			bids = append(bids, &depth.DepthLevel{
				Price:  StringToFloat(change[1]),
				Amount: StringToFloat(change[2]),
			})
		}
	}
	resp = base.OrderBook{
		Exchange:     common.Exchange_GEMINI,
		Market:       common.Market_SPOT,
		Type:         common.SymbolType_SPOT_NORMAL,
		Symbol:       symbolNameMap[respL2Updates.Symbol],
		TimeExchange: 0,
		TimeReceive:  uint64(t),
		Asks:         asks,
		Bids:         bids,
	}
	if isFirstMsg == true {
		deltaDepth.Symbol = resp.Symbol
		deltaDepth.Market = resp.Market
		deltaDepth.Type = resp.Type
		deltaDepth.Asks = resp.Asks
		deltaDepth.Bids = resp.Bids
		deltaDepth.TimeReceive = t
		deltaDepth.IsFullDepth = true
		if content, ok = b.DepthDeltaUpdateMap.Get(deltaDepth.Symbol); ok {
			chUpdate, ok = content.(chan *base.DeltaDepthUpdate)
		}
		base.SendChan(chUpdate, deltaDepth, "Initialization", deltaDepth.Symbol)
	} else {
		transferDiffDepth(&resp, deltaDepth)
		sym := deltaDepth.Symbol
		if b.IsPublishDelta {
			if _, ok = b.depthIncrementSnapshotDeltaGroupChanMap[sym]; ok {
				base.SendChan(b.depthIncrementSnapshotDeltaGroupChanMap[sym], deltaDepth.Transfer2Depth(), "depthIncrementSnapshotDeltaGroupChanMap")
			} else {
				logger.Logger.Warn("get symbol from channel map err:", sym)
			}
		}
		if b.IsPublishFull {
			if content, ok = b.DepthDeltaUpdateMap.Get(sym); ok {
				chUpdate, ok = content.(chan *base.DeltaDepthUpdate)
			}
			base.SendChan(chUpdate, deltaDepth, "DeltaDepthUpdate", sym)
		}
		b.CheckSendStatus.CheckUpdateTimeMap.Set(deltaDepth.Symbol, time.Now().UnixMicro())
	}
	return nil
}
