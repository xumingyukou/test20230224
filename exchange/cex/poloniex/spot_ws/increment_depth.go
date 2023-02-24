package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"encoding/json"
	"errors"
	"math/rand"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/sdk"
)

func (h *WebSocketSpotHandle) SetDepthIncrementSnapShotConf(symbols []*client.SymbolInfo, conf *base.IncrementDepthConf) {
	if conf.DepthCapLevel <= 0 {
		conf.DepthCapLevel = 100
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
			CheckDepthCacheChan = make(chan *base.OrderBook, conf.DepthCapLevel)
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

func (h *WebSocketSpotHandle) SetDepthIncrementSnapshotReconnectChan(ch chan string) {
	h.DepthIncrementSnapshotReconnectSymbol = ch
}

func (h *WebSocketSpotHandle) SetReconnectSignalChannel(chMap map[string]chan struct{}) {
	if chMap != nil {
		if h.reconnectSignalChannel == nil {
			h.reconnectSignalChannel = make(map[string]chan struct{})
		}
		for symbol, ch := range chMap {
			h.reconnectSignalChannel[symbol] = ch
		}
	}
}

func (h *WebSocketSpotHandle) updateDeltaDepth(symbol string) {
	var (
		obj             interface{}
		deltaCh         chan *base.DeltaDepthUpdate
		delta           *base.DeltaDepthUpdate
		full            *base.OrderBook
		deltaCacheList  = make([]*base.DeltaDepthUpdate, 0, 10)
		depthCache      = &base.OrderBook{}
		fullCh          chan *base.OrderBook
		deltaSendCh     chan *client.WsDepthRsp
		fullSendCh      chan *depth.Depth
		reconnectSignal chan struct{}
		ok              bool
		firstDelta      = true
		isRescribing    = false
		firstFullSent   = false
		funcName        = "updateDeltaDepth"
	)
	if obj, ok = h.DepthDeltaUpdateMap.Get(symbol); ok {
		deltaCh, ok = obj.(chan *base.DeltaDepthUpdate)
	}
	if !ok {
		logger.Logger.Error("get DepthDeltaUpdateMap ", funcName, " ", symbol)
		return
	}

	if obj, ok = h.CheckDepthCacheChanMap.Get(symbol); ok {
		fullCh, ok = obj.(chan *base.OrderBook)
	}
	if !ok {
		logger.Logger.Error("get DepthCacheChanMap ", funcName, " ", symbol)
		return
	}
	if h.IsPublishDelta {
		deltaSendCh, ok = h.depthIncrementSnapshotDeltaGroupChanMap[symbol]
		if !ok {
			logger.Logger.Error("get delta send channel ", funcName, " ", symbol)
			return
		}
	}
	if h.IsPublishFull {
		fullSendCh, ok = h.depthIncrementSnapshotFullGroupChanMap[symbol]
		if !ok {
			logger.Logger.Error("get full send channel ", funcName, " ", symbol)
			return
		}
	}
	reconnectSignal, ok = h.reconnectSignalChannel[symbol]
	if !ok {
		logger.Logger.Error("get reconnect signal channel ", funcName, " ", symbol)
		return
	}
LOOP:
	for {
		select {
		case <-h.Ctx.Done():
			break LOOP
		case delta, ok = <-deltaCh:
			if h.IsPublishDelta {
				res := delta.Transfer2Depth()
				base.SendChan(deltaSendCh, res, "deltaPublish", symbol)
			}
			if h.IsPublishFull {
				deltaCacheList = append(deltaCacheList, delta)
				if firstDelta {
					logger.Logger.Infof("receive delta %s %d %d-%d", symbol, full.UpdateId, deltaCacheList[0].UpdateStartId, deltaCacheList[0].UpdateEndId)
					firstDelta = false
				}
				if depthCache.UpdateId == 0 || depthCache.UpdateId < deltaCacheList[0].UpdateStartId {
					if !isRescribing {
						logger.Logger.Errorf("missing delta %s %d %d", symbol, depthCache.UpdateId, deltaCacheList[0].UpdateStartId)
						// 数据丢失，重新订阅
						h.DepthIncrementSnapshotReconnectSymbol <- symbol
						deltaCacheList = append([]*base.DeltaDepthUpdate{})
						isRescribing = true
					}
					firstFullSent = false
					continue
				} else if depthCache.UpdateId >= deltaCacheList[len(deltaCacheList)-1].UpdateEndId {
					// 等待增量
					logger.Logger.Warnf("waiting delta %s %d %d", symbol, depthCache.UpdateId, deltaCacheList[len(deltaCacheList)-1].UpdateEndId)
					deltaCacheList = append([]*base.DeltaDepthUpdate{})
					firstFullSent = false
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
						logger.Logger.Error("deltaList error", symbol)
						h.DepthIncrementSnapshotReconnectSymbol <- symbol
						deltaCacheList = append([]*base.DeltaDepthUpdate{})
						firstFullSent = false
						continue
					}
					deltaCacheList = deltaCacheList[target:]
				UPDATE_LOOP:
					for _, deltaCache := range deltaCacheList {
						if depthCache.UpdateId != delta.UpdateStartId {
							if !isRescribing {
								logger.Logger.Error("depthCache updateId != delta.LastId ", depthCache.UpdateId, delta.UpdateStartId)
								h.DepthIncrementSnapshotReconnectSymbol <- symbol
								isRescribing = true
								firstFullSent = false
								break UPDATE_LOOP
							}
						}
						base.UpdateBidsAndAsks(deltaCache, depthCache, h.DepthCapLevel, nil)
					}
					fullToSend := &depth.Depth{}
					depthCache.Transfer2Depth(h.DepthLevel, fullToSend)
					if !firstFullSent {
						fullToSend.Hdr = base.MakeFirstDepthHdr()
						firstFullSent = true
					}
					base.SendChan(fullSendCh, fullToSend, "depth full send", symbol)
					deltaCacheList = append([]*base.DeltaDepthUpdate{})
				}
			}
		case full, ok = <-fullCh:
			logger.Logger.Info("receive full ", full.UpdateId, " ", symbol)
			depthCache = full
			isRescribing = false
			firstDelta = true
			firstFullSent = false
		case _, ok = <-reconnectSignal:
			if ok {
				isRescribing = true
				depthCache = &base.OrderBook{}
				firstFullSent = false
			}
		}
	}
	return
}

func (h *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte, receiveTime int64) error {
	var (
		resp    = &BookLv2Msg{}
		err     error
		ok      = false
		deltaCh chan *base.DeltaDepthUpdate
		fullCh  chan *base.OrderBook
		obj     interface{}
	)
	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Logger.Error("depth increment parse ", err, "data:", string(data))
		return err
	}
	// full
	if resp.Action == "snapshot" {
		for _, t := range resp.Data {
			full, err := h.parseDepthWhole(&t, receiveTime)
			if err != nil {
				logger.Logger.Error("parse whole depth in snapshot", err)
				return err
			}
			if obj, ok = h.CheckDepthCacheChanMap.Get(full.Symbol); ok {
				if fullCh, ok = obj.(chan *base.OrderBook); ok {
					base.SendChan(fullCh, full, "depthCacheChan", full.Symbol)
				}
			}
			if !ok {
				logger.Logger.Error("get DepthCacheChanMap err, symbol:", full.Symbol)
				return errors.New("get DepthCacheChanMap err, symbol:" + full.Symbol)
			}
			h.CheckSendStatus.CheckUpdateTimeMap.Set(full.Symbol, receiveTime)
		}
	} else if resp.Action == "update" {
		for _, t := range resp.Data {
			delta, err := h.parseDepthDelta(&t, receiveTime)

			if err != nil {
				logger.Logger.Error("parse depth delta in snapshot", err)
				return err
			}
			if obj, ok = h.DepthDeltaUpdateMap.Get(delta.Symbol); ok {
				if deltaCh, ok = obj.(chan *base.DeltaDepthUpdate); ok {
					base.SendChan(deltaCh, delta, "depthDeltaUpdate", delta.Symbol)
				}
			}
			if !ok {
				logger.Logger.Error("get DepthDeltaUpdateMap err, symbol:", delta.Symbol)
				return errors.New("get DepthDeltaUpdateMap err, symbol:" + delta.Symbol)
			}
			h.CheckSendStatus.CheckUpdateTimeMap.Set(delta.Symbol, receiveTime)
		}
	} else {
		logger.Logger.Warn("unknown depth type:", string(data))
		return errors.New("unknown depth type:" + resp.Action)
	}
	return nil
}
