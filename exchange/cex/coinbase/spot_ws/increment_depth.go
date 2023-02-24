package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/warmplanet/proto/go/common"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/sdk"
)

func (c *WebSocketSpotHandle) SetDepthIncrementSnapShotConf(symbols []*client.SymbolInfo, conf *base.IncrementDepthConf) {
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
	//conf.DepthNotMatchChanMap = make(map[*client.SymbolInfo]chan bool)
	c.IncrementDepthConf = conf
	c.symbolMap = make(map[string]*client.SymbolInfo)
	for _, symbol := range symbols {
		var (
			CheckDepthCacheChan = make(chan *base.OrderBook, conf.DepthCapLevel)
		)
		conf.DepthCacheMap.Set(symbol.Symbol, nil)
		conf.DepthDeltaUpdateMap.Set(symbol.Symbol, make(chan *base.DeltaDepthUpdate, 50))
		conf.CheckDepthCacheChanMap.Set(symbol.Symbol, CheckDepthCacheChan)
		//conf.DepthNotMatchChanMap[symbol] = make(chan bool, conf.DepthCapLevel)
		conf.CheckStates.Set(symbol.Symbol, false)
		c.symbolMap[symbol.Symbol] = symbol
	}
	for _, symbol := range symbols {
		go c.UpdateDeltaDepth(symbol.Symbol)
	}
	//go c.Check()
}

func (c *WebSocketSpotHandle) UpdateDeltaDepth(symbol string) {
	var (
		content    interface{}
		ok         bool
		ch         chan *base.DeltaDepthUpdate
		diffDepth  *base.DeltaDepthUpdate
		depthCache = &base.OrderBook{Symbol: symbol}
	)
	if content, ok = c.DepthDeltaUpdateMap.Get(symbol); ok {
		ch, ok = content.(chan *base.DeltaDepthUpdate)
	}
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
			if diffDepth.IsFullDepth {
				depthCache = &base.OrderBook{Symbol: symbol}
			}
			base.UpdateBidsAndAsks(diffDepth, depthCache, c.DepthCapLevel, nil)
			c.AfterUpdateHandle(symbol, depthCache)
		case <-c.Ctx.Done():
			break LOOP
		}

	}
}

func (c *WebSocketSpotHandle) AfterUpdateHandle(sym string, depthCache *base.OrderBook) {
	var (
		res = &depth.Depth{}
		//content         interface{}
		//ok, checkStatus bool
	)
	//if content, ok = c.CheckStates.Get(sym); ok {
	//	checkStatus, ok = content.(bool)
	//}
	//if ok && checkStatus {
	//	var (
	//		ch chan *base.OrderBook
	//	)
	//	if content, ok = c.CheckDepthCacheChanMap.Get(depthCache.Symbol); ok {
	//		if ch, ok = content.(chan *base.OrderBook); ok {
	//			ch <- depthCache.Copy()
	//		}
	//	}
	//	if !ok {
	//		logger.Logger.Error("cannot get CheckDepthCacheChanMap err, symbol:", depthCache.Symbol)
	//	}
	//}
	depthCache.Transfer2Depth(c.DepthLevel, res)
	if _, ok := c.depthIncrementSnapshotFullGroupChanMap[sym]; ok {
		firstFullTagIf, ok := c.firstFullSentMap.Load(sym)
		if ok {
			firstFullTag, _ := firstFullTagIf.(bool)
			if firstFullTag {
				c.firstFullSentMap.Store(sym, false)
				res.Hdr = base.MakeFirstDepthHdr()
			}
		}
		base.SendChan(c.depthIncrementSnapshotFullGroupChanMap[sym], res, "DepthIncrementSnapShotGroupHandle")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
}

func (c *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte) error {
	var (
		resp       RespIncrementSnapshot
		deltaDepth = &base.DeltaDepthUpdate{}
		ok         bool
		content    interface{}
		chUpdate   chan *base.DeltaDepthUpdate
		err        error
		t          = time.Now().UnixMicro()
	)
	c.Lock.Lock()
	defer c.Lock.Unlock()

	if err = json.Unmarshal(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}

	if resp.Type == "subscriptions" {
		return nil
	} else if resp.Type == "heartbeat" {
		return nil
	} else if resp.Type == "error" {
		fmt.Println("Error:", resp.Message, resp.Reason)
		return err
	} else if resp.Type == "snapshot" {
		asks, err := DepthLevelParse(resp.Asks)
		if err != nil {
			logger.Logger.Error("depth level parse err", err, string(data))
			return err
		}
		bids, err := DepthLevelParse(resp.Bids)
		if err != nil {
			logger.Logger.Error("depth level parse err", err, string(data))
			return err
		}
		deltaDepth.Symbol = FormatSymbols(resp.ProductId)
		deltaDepth.Market = common.Market_SPOT
		deltaDepth.Type = common.SymbolType_SPOT_NORMAL
		deltaDepth.Asks = asks
		deltaDepth.Bids = bids
		deltaDepth.TimeReceive = t
		deltaDepth.IsFullDepth = true
		if content, ok = c.DepthDeltaUpdateMap.Get(deltaDepth.Symbol); ok {
			chUpdate, ok = content.(chan *base.DeltaDepthUpdate)
		}
		c.firstFullSentMap.Store(deltaDepth.Symbol, true)
		base.SendChan(chUpdate, deltaDepth, "Initialization", deltaDepth.Symbol)
	}

	if len(resp.Changes) == 0 {
		return nil
	}
	transferDiffDepth(&resp, deltaDepth)
	sym := deltaDepth.Symbol
	if c.IsPublishDelta {
		if _, ok := c.depthIncrementSnapshotDeltaGroupChanMap[sym]; ok {
			base.SendChan(c.depthIncrementSnapshotDeltaGroupChanMap[sym], deltaDepth.Transfer2Depth(), "depthIncrementSnapshotDeltaGroupChanMap")
		} else {
			logger.Logger.Warn("get symbol from channel map err:", sym)
		}
	}
	if c.IsPublishFull {

		if content, ok = c.DepthDeltaUpdateMap.Get(sym); ok {
			chUpdate, ok = content.(chan *base.DeltaDepthUpdate)
		}
		base.SendChan(chUpdate, deltaDepth, "DeltaDepthUpdate", deltaDepth.Symbol)
	}
	c.CheckSendStatus.CheckUpdateTimeMap.Set(sym, time.Now().UnixMicro())
	return nil
}

func transferDiffDepth(resp *RespIncrementSnapshot, diff *base.DeltaDepthUpdate) {
	var (
		amount, price float64
		err           error
	)

	diff.Symbol = FormatSymbols(resp.ProductId)
	diff.Market = common.Market_SPOT
	diff.Type = common.SymbolType_SPOT_NORMAL
	diff.TimeExchange = resp.Time.UnixMicro()
	diff.TimeReceive = time.Now().UnixMicro()
	//清空bids，asks
	diff.Bids = diff.Bids[:0]
	diff.Asks = diff.Asks[:0]

	for _, incrementInfo := range resp.Changes {
		price, err = strconv.ParseFloat(incrementInfo.Price, 64)
		amount, err = strconv.ParseFloat(incrementInfo.Size, 64)
		if err != nil {
			fmt.Println("Parse Error")
			return
		}
		switch incrementInfo.Side {
		case BUY:
			diff.Bids = append(diff.Bids, &depth.DepthLevel{
				Price:  price,
				Amount: amount,
			})
		case SELL:
			diff.Asks = append(diff.Asks, &depth.DepthLevel{
				Price:  price,
				Amount: amount,
			})
		}
	}
	sort.Stable(diff.Asks)
	sort.Stable(sort.Reverse(diff.Bids))
	return
}

/*Deprecated Function
func (c *WebSocketSpotHandle) CheckFullDepth(symbol string) {
	var (
		fullDepth *base.OrderBook
		ch        chan *base.OrderBook
		content   interface{}
		ok        bool
		err       error
	)
	c.CheckStates.Set(symbol, true)
	defer func() { c.CheckStates.Set(symbol, false) }()
	fullDepth = &base.OrderBook{}
	if content, ok = c.CheckDepthCacheChanMap.Get(symbol); ok {
		ch, ok = content.(chan *base.OrderBook)
	}
	if !ok {
		logger.Logger.Error("get CheckDepthCacheChanMap err, symbol:", symbol)
		return
	}

MAINLOOP:
	for {
		select {
		case dep := <-ch:
			if fullDepth.TimeExchange > dep.TimeExchange {
				continue
			} else if fullDepth.TimeExchange == dep.TimeExchange && dep.TimeExchange > 0 {
				if dep.Equal(fullDepth, c.DepthCheckLevel) {
					logger.Logger.Info("check full depth: correct", symbol)
					break MAINLOOP
				} else {
					logger.Logger.Error(symbol, "depth level:", c.DepthCheckLevel, "need reset")
					logger.SaveToFile(fmt.Sprintf("depth_error_%s.json", strings.Replace(symbol, "/", "_", 1)), "depCache", dep, "fullDepth", fullDepth)
					errCh := c.DepthNotMatchChanMap[c.symbolMap[symbol]]
					errCh <- true
					break MAINLOOP
				}
			} else {
				if fullDepth, err = c.GetFullDepth(symbol); err != nil {
					logger.Logger.Error("CheckFullDepth get full depth error", err)
					time.Sleep(time.Second * 30)
					continue
				}
			}
		}
	}
}
*/
/* Deprecated Function
func (c *WebSocketSpotHandle) UpdateFullDepth(depthCache *base.OrderBook) (err error) {
	fullDepth := &base.OrderBook{}
	fullDepth, err = c.GetFullDepth(depthCache.Symbol)
	if err == nil {
		depthCache.Bids = fullDepth.Bids
		depthCache.Asks = fullDepth.Asks
		depthCache.TimeExchange = fullDepth.TimeExchange
		depthCache.UpdateId = fullDepth.UpdateId
	}
	return
}
*/
/*Deprecated Function
func (c *WebSocketSpotHandle) Check() {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger.Error("incr depth check panic:", logger.PanicTrace(err))
			time.Sleep(time.Second)
			go c.Check()
		}
	}()
	heartTimer := time.NewTimer(time.Duration(c.CheckTimeSec) * time.Second)
LOOP:
	for {
		select {
		case <-heartTimer.C:
			for _, symbol := range c.DepthCacheMap.Keys() {
				go c.CheckFullDepth(symbol)
				time.Sleep(time.Second * 10)
			}
			heartTimer.Reset(time.Duration(c.CheckTimeSec) * time.Second)
		case <-c.Ctx.Done():
			c.Reset()
			break LOOP
		}
	}
}
*/
/*Deprecated Function
func (c *WebSocketSpotHandle) Reset() {
	c.CheckDepthCacheChanMap = sdk.NewCmapI()
	c.DepthCacheMap = sdk.NewCmapI()
	c.DepthDeltaUpdateMap = sdk.NewCmapI()
}*/
