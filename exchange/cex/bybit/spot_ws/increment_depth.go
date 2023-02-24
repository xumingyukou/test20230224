package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"fmt"
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
	//go b.Check()
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

func (b *WebSocketSpotHandle) UpdateDeltaDepth(symbol string) {
	var (
		content       interface{}
		ok            bool
		ch            chan *base.DeltaDepthUpdate
		diffDepth     *base.DeltaDepthUpdate
		depthCache    = &base.OrderBook{Symbol: symbol}
		fullCh        chan *depth.Depth
		firstFullSent = false
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
			if diffDepth.IsFullDepth {
				depthCache = &base.OrderBook{Symbol: symbol}
				firstFullSent = false
			}
			base.UpdateBidsAndAsks(diffDepth, depthCache, b.DepthCapLevel, nil)
			full := &depth.Depth{}
			depthCache.Transfer2Depth(b.DepthLevel, full)
			if !firstFullSent {
				full.Hdr = base.MakeFirstDepthHdr()
				firstFullSent = true
			}
			base.SendChan(fullCh, full, "DepthIncrementSnapShotGroupHandle")
		case <-b.Ctx.Done():
			break LOOP
		}

	}
}

func (b *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte, t int64) error {
	b.Lock.Lock()
	var (
		respTemp   RespIncrementDepthStream
		resp       base.OrderBook
		err        error
		asks       []*depth.DepthLevel
		bids       []*depth.DepthLevel
		deltaDepth = &base.DeltaDepthUpdate{}
		chUpdate   chan *base.DeltaDepthUpdate
		ok         bool
		content    interface{}
	)
	defer b.Lock.Unlock()
	if err = json.Unmarshal(data, &respTemp); err != nil {
		fmt.Println("Err", err)
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if respTemp.F == true {
		logger.Logger.Info("Subscribed Success: ", symbolNameMap[respTemp.Symbol])
		asks, err = DepthLevelParse(respTemp.Data[0].A)
		if err != nil {
			logger.Logger.Error("book ticker parse err", err, string(data))
			return err
		}
		bids, err = DepthLevelParse(respTemp.Data[0].B)
		if err != nil {
			logger.Logger.Error("book ticker parse err", err, string(data))
			return err
		}
		resp = base.OrderBook{
			Exchange:     common.Exchange_BYBIT,
			Market:       common.Market_SPOT,
			Type:         common.SymbolType_SPOT_NORMAL,
			Symbol:       symbolNameMap[respTemp.Symbol],
			TimeExchange: uint64(respTemp.Data[0].T * 1000),
			TimeReceive:  uint64(t),
			Asks:         asks,
			Bids:         bids,
		}
		deltaDepth.Symbol = resp.Symbol
		deltaDepth.Market = resp.Market
		deltaDepth.Type = resp.Type
		deltaDepth.Asks = resp.Asks
		deltaDepth.Bids = resp.Bids
		deltaDepth.IsFullDepth = true
		deltaDepth.TimeReceive = time.Now().UnixMicro()
		if content, ok = b.DepthDeltaUpdateMap.Get(deltaDepth.Symbol); ok {
			chUpdate, ok = content.(chan *base.DeltaDepthUpdate)
		}
		base.SendChan(chUpdate, deltaDepth, "Initialization", deltaDepth.Symbol)
		return nil
	}
	if respTemp.Data == nil {
		return nil
	}
	asks, err = DepthLevelParse(respTemp.Data[0].A)
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	bids, err = DepthLevelParse(respTemp.Data[0].B)
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	resp = base.OrderBook{
		Exchange:     common.Exchange_BYBIT,
		Market:       common.Market_SPOT,
		Type:         common.SymbolType_SPOT_NORMAL,
		Symbol:       symbolNameMap[respTemp.Symbol],
		TimeExchange: uint64(respTemp.Data[0].T * 1000),
		TimeReceive:  uint64(t),
		Asks:         asks,
		Bids:         bids,
	}
	transferDiffDepth(&resp, deltaDepth)
	sym := deltaDepth.Symbol
	if b.IsPublishDelta {
		if _, ok = b.depthIncrementSnapshotDeltaGroupChanMap[sym]; ok {
			base.SendChan(b.depthIncrementSnapshotDeltaGroupChanMap[sym], deltaDepth.Transfer2Depth(), "depthIncrementSnapshotDeltaGroupChanMap")
		} else {
			//fmt.Println("Error1: ", resp)
			logger.Logger.Warn("get symbol from channel map err:", sym)
		}
	}
	if b.IsPublishFull {
		if content, ok = b.DepthDeltaUpdateMap.Get(sym); ok {
			chUpdate, ok = content.(chan *base.DeltaDepthUpdate)
		}
		base.SendChan(chUpdate, deltaDepth, "DeltaDepthUpdate", deltaDepth.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(deltaDepth.Symbol, time.Now().UnixMicro())

	return nil
}
