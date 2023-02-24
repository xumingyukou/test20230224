package spot_ws

import (
	"clients/crypto"
	"clients/exchange/cex/base"
	"clients/logger"
	"clients/transform"
	"errors"
	"sort"
	"strings"

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
		conf.DepthDeltaUpdateMap.Set(symbol.Symbol, make(chan *base.DeltaDepthUpdate, 100))
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

func (b *WebSocketSpotHandle) SetDepthIncrementSnapshotReconnectChan(ch chan string) {
	b.DepthIncrementSnapshotReconnectSymbol = ch
}

func (b *WebSocketSpotHandle) UpdateDeltaDepth(symbol string) {
	var (
		deltaCh chan *base.DeltaDepthUpdate
		// fullCh      chan *base.OrderBook
		delta *base.DeltaDepthUpdate
		// full        *base.OrderBook
		deltaSendCh chan *client.WsDepthRsp
		fullSendCh  chan *depth.Depth
		depthCache  *base.OrderBook
		// deltaCacheList = make([]*base.DeltaDepthUpdate, 0, 10)
		ok  bool
		obj interface{}
		//err                 error
		funcName      = "updateDeltaDepth"
		firstFullSent = false
	)

	if obj, ok = b.DepthDeltaUpdateMap.Get(symbol); ok {
		deltaCh, ok = obj.(chan *base.DeltaDepthUpdate)
	}
	if !ok {
		logger.Logger.Error("get DepthDeltaUpdateMap ", funcName, " ", symbol)
		return
	}

	// if obj, ok = b.CheckDepthCacheChanMap.Get(symbol); ok {
	// 	fullCh, ok = obj.(chan *base.OrderBook)
	// }
	// if !ok {
	// 	logger.Logger.Error("get DepthCacheChanMap ", funcName, " ", symbol)
	// 	return
	// }
	if b.IsPublishDelta {
		deltaSendCh, ok = b.depthIncrementSnapshotDeltaGroupChanMap[symbol]
		if !ok {
			logger.Logger.Error("get delta send channel ", funcName, " ", symbol)
			return
		}
	}
	if b.IsPublishFull {
		fullSendCh, ok = b.depthIncrementSnapshotFullGroupChanMap[symbol]
		if !ok {
			logger.Logger.Error("get full send channel ", funcName, " ", symbol)
			return
		}
	}

LOOP:
	for {
		select {
		case <-b.Ctx.Done():
			break LOOP
		case delta, ok = <-deltaCh:
			// 如果来的是全量，需要把depthCache初始化，避免新全量去merge老全量
			if delta.IsFullDepth {
				firstFullSent = false
				depthCache = &base.OrderBook{Symbol: symbol}
			}
			if delta.UpdateEndId == 1 {
				// checksum
				correctChecksum := delta.UpdateStartId
				if depthCache != nil {
					checkStrs := make([]string, 0, 50)
					for i := 0; i < 25; i++ {
						if i < len(depthCache.Bids) {
							checkStrs = append(checkStrs, transform.FormatFloatE(depthCache.Bids[i].Price, 1e-6))
							checkStrs = append(checkStrs, transform.FormatFloatE(depthCache.Bids[i].Amount, 1e-6))
						}
						if i < len(depthCache.Asks) {
							checkStrs = append(checkStrs, transform.FormatFloatE(depthCache.Asks[i].Price, 1e-6))
							checkStrs = append(checkStrs, transform.FormatFloatE(-depthCache.Asks[i].Amount, 1e-6))
						}
					}
					checkStr := strings.Join(checkStrs, ":")
					checksum := int64(int32(crypto.CRC32(checkStr)))
					if checksum != correctChecksum {
						logger.Logger.Error("error check sum ", symbol, " ", correctChecksum, " ", checksum, " ", checkStr)
						firstFullSent = false
						b.DepthIncrementSnapshotReconnectSymbol <- symbol
					}
				}
				continue
			}
			if b.IsPublishDelta && !delta.IsFullDepth {
				res := delta.Transfer2Depth()
				base.SendChan(deltaSendCh, res, "deltaPublish", symbol)
			}
			if b.IsPublishFull {
				if depthCache != nil {
					fullToSend := &depth.Depth{}
					base.UpdateBidsAndAsks(delta, depthCache, b.DepthCapLevel, fullToSend)
					if !firstFullSent {
						fullToSend.Hdr = base.MakeFirstDepthHdr()
						firstFullSent = true
					}
					base.SendChan(fullSendCh, fullToSend, "depth full send", symbol)
				}
			}
		}
	}
	return
}

func (b *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte, t int64, isFirstMsg bool) error {
	b.Lock.Lock()
	var (
		res Response
		err error
		ok  bool
		// fullCh  chan *base.OrderBook
		deltaCh chan *base.DeltaDepthUpdate
		obj     interface{}
	)
	b.Lock.Unlock()
	if err = json.Unmarshal(data, &res); err != nil {
		logger.Logger.Error("receive data err:", res)
		return err
	}
	chanIDf64, idOk := res[0].(float64)
	chanID := int64(chanIDf64)
	// fmt.Println(string(data))
	if res[1] == "cs" {
		// checksum
		checksum, csOk := res[2].(float64)
		if !csOk {
			logger.Logger.Error("receive data err[checksum]:", res)
			return nil
		}

		delta := &base.DeltaDepthUpdate{
			Market:        common.Market_SPOT,
			Type:          common.SymbolType_SPOT_NORMAL,
			Symbol:        Exchange2Client[b.GetSymbol(chanID)],
			TimeReceive:   t,
			UpdateStartId: int64(checksum),
			UpdateEndId:   1,
		}

		if obj, ok = b.DepthDeltaUpdateMap.Get(delta.Symbol); ok {
			if deltaCh, ok = obj.(chan *base.DeltaDepthUpdate); ok {
				base.SendChan(deltaCh, delta, "deltaDepthUpdate", delta.Symbol)
			}
		}

		if !ok {
			logger.Logger.Error("get DepthDeltaUpdateMap err, symbol:", delta.Symbol)
			return errors.New("get DepthDeltaUpdateMap err, symbol:" + delta.Symbol)
		}

		b.CheckSendStatus.CheckUpdateTimeMap.Set(delta.Symbol, t)
		return nil
	}
	depths, depthOk := res[1].([]interface{})
	if idOk && depthOk {
		_, isArray := depths[0].([]interface{})
		if isArray { // full
			asks := make([]*depth.DepthLevel, 0, len(depths))
			bids := make([]*depth.DepthLevel, 0, len(depths))
			for _, dep := range depths {
				Item := dep.([]interface{})
				// price := Item[0].(float64)
				// amount := Item[2].(float64)

				price, pok := Item[0].(float64)
				count, cok := Item[1].(float64)
				amount, aok := Item[2].(float64)
				if !(pok && cok && aok) {
					logger.Logger.Error("parse depth message ", depths)
					return nil
				}
				if amount <= 0 {
					if count == 0 {
						amount = 0
					}
					asks = append(asks, &depth.DepthLevel{
						Price:  price,
						Amount: -amount,
					})
				} else {
					if count == 0 {
						amount = 0
					}
					bids = append(bids, &depth.DepthLevel{
						Price:  price,
						Amount: amount,
					})
				}

				// if amount <= 0 {
				// 	asks = append(asks, &depth.DepthLevel{
				// 		Price:  price,
				// 		Amount: -amount,
				// 	})
				// } else {
				// 	bids = append(bids, &depth.DepthLevel{
				// 		Price:  price,
				// 		Amount: amount,
				// 	})
				// }
			}

			depth := &base.DeltaDepthUpdate{
				Market:      common.Market_SPOT,
				Type:        common.SymbolType_SPOT_NORMAL,
				Symbol:      Exchange2Client[b.GetSymbol(chanID)],
				TimeReceive: t,
				Asks:        asks,
				Bids:        bids,
				IsFullDepth: isFirstMsg,
			}

			sort.Stable(depth.Asks)
			sort.Stable(sort.Reverse(depth.Bids))

			if obj, ok = b.DepthDeltaUpdateMap.Get(depth.Symbol); ok {
				if deltaCh, ok = obj.(chan *base.DeltaDepthUpdate); ok {
					if isFirstMsg {
						base.SendChan(deltaCh, depth, "Initialization", depth.Symbol)
					} else {
						base.SendChan(deltaCh, depth, "DeltaDepthUpdate", depth.Symbol)
					}
				}
			}
			if !ok {
				logger.Logger.Error("get DepthCacheChanMap err, symbol:", depth.Symbol)
				return errors.New("get DepthCacheChanMap err, symbol:" + depth.Symbol)
			}

			b.CheckSendStatus.CheckUpdateTimeMap.Set(depth.Symbol, t)
			return nil
		} else { // delta
			logger.Logger.Error("unexpected message not array", res)
			return nil
		}
	} else {
		logger.Logger.Error("unexpected message ", res)
	}
	return nil
}
