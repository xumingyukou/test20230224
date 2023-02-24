package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"encoding/json"
	"errors"
	"strings"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/sdk"
)

func (h *WebSocketSpotHandle) SetDepthIncrementSnapShotConf(symbolInfo []*client.SymbolInfo, conf *base.IncrementDepthConf) {
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
	h.IncrementDepthConf = conf
	for _, symbol := range symbolInfo {
		var (
			CheckDepthCacheChan = make(chan *base.OrderBook, conf.DepthCapLevel)
		)
		conf.DepthCacheMap.Set(symbol.Symbol, nil)
		conf.DepthDeltaUpdateMap.Set(symbol.Symbol, make(chan *base.DeltaDepthUpdate))
		conf.CheckDepthCacheChanMap.Set(symbol.Symbol, CheckDepthCacheChan)
		conf.DepthNotMatchChanMap[symbol] = make(chan bool, conf.DepthCapLevel)
		conf.CheckStates.Set(symbol.Symbol, false)
	}
	for _, symbol := range symbolInfo {
		go h.updateDeltaDepth(symbol.Symbol)
	}
	// go h.Check()
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
		funcName        = "updateDeltaDepth"
		firstFullSent   = false
		getFirstFull    = false
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
				// logger.Logger.Infof("get delta %s %d %d", symbol, delta.UpdateStartId, delta.UpdateEndId)
				deltaCacheList = append(deltaCacheList, delta)
				if depthCache.UpdateId == 0 || (!getFirstFull && depthCache.UpdateId < deltaCacheList[0].UpdateStartId) {
					// 如果还没有确定拿到了第一次时间在增量间的全量，需要继续拿
					logger.Logger.Tracef("early cache %s %d %d-%d", symbol, depthCache.UpdateId, deltaCacheList[0].UpdateStartId, deltaCacheList[len(deltaCacheList)-1].UpdateEndId)
					firstFullSent = false
					continue
				} else if depthCache.UpdateId > deltaCacheList[len(deltaCacheList)-1].UpdateEndId {
					getFirstFull = true
					deltaCacheList = append([]*base.DeltaDepthUpdate{})
					logger.Logger.Tracef("late cache %s %d %d-%d", symbol, depthCache.UpdateId, deltaCacheList[0].UpdateStartId, deltaCacheList[len(deltaCacheList)-1].UpdateEndId)
					firstFullSent = false
					continue
				} else {
					logger.Logger.Tracef("target cache %s %d %d-%d", symbol, depthCache.UpdateId, deltaCacheList[0].UpdateStartId, deltaCacheList[len(deltaCacheList)-1].UpdateEndId)
					// 如果拿到过可能可以更新
					target := 0
					for i, deltaCache := range deltaCacheList {
						target = i
						if depthCache.UpdateId >= deltaCache.UpdateStartId {
							logger.Logger.Infof("target cache %s skip %d", symbol, deltaCache.UpdateStartId)
						} else {
							break
						}
					}
					if target == 0 {
						// 没有找到全量的id
						if !getFirstFull {
							logger.Logger.Error("deltaList error", symbol)
							firstFullSent = false
							continue
						}
					}
					getFirstFull = true
					deltaCacheList = deltaCacheList[target:]
					for _, deltaCache := range deltaCacheList {
						logger.Logger.Tracef("target cache %s update %d", symbol, deltaCache.UpdateStartId)
						base.UpdateBidsAndAsks(deltaCache, depthCache, h.DepthCapLevel, nil)
					}
					fullToSend := &depth.Depth{}
					depthCache.Transfer2Depth(h.DepthLevel, fullToSend)
					if !firstFullSent {
						fullToSend.Hdr = base.MakeFirstDepthHdr()
						firstFullSent = true
					}
					base.SendChan(fullSendCh, fullToSend, "depth full send", symbol)
					// 增量使用完，清空
					deltaCacheList = append([]*base.DeltaDepthUpdate{})
				}
			}
		case full, ok = <-fullCh:
			// 首次收到全量，变为本地缓存
			// if !getFirstFull {
			depthCache = full
			firstFullSent = false
			// logger.Logger.Infof("get first full %s %d", symbol, full.UpdateId)
			// }

		// case full, ok = <-fullCh:
		// 	// 首次收到全量，变为本地缓存
		// 	if !getFirstFull {
		// 		depthCache = full
		// 		firstFullSent = false
		// 		logger.Logger.Infof("get first full %s %d", symbol, full.UpdateId)
		// 	}
		case _, ok = <-reconnectSignal:
			// 重连时清空当前缓存全量
			if ok {
				getFirstFull = false
				firstFullSent = false
				depthCache = &base.OrderBook{}
			}
		}
	}
	return
}

func (h *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte, receiveTime int64) error {
	var (
		obj     interface{}
		deltaCh chan *base.DeltaDepthUpdate
		fullCh  chan *base.OrderBook
		ok      = false
		err     error
		msg     = &BaseMsg{}
	)
	err = json.Unmarshal(data, msg)
	if err != nil {
		logger.Logger.Error("parse base msg in depthIncrementSnapshotGroupHandle ", err.Error())
	}
	if strings.HasPrefix(msg.RoomName, "depth_whole") {
		logger.Logger.Trace(string(data))
		res, err := h.parseDepthWhole(data, receiveTime)
		if err != nil {
			logger.Logger.Error("parse whole depth in snapshot", err, " data:", string(data))
			return err
		}
		if obj, ok = h.CheckDepthCacheChanMap.Get(res.Symbol); ok {
			if fullCh, ok = obj.(chan *base.OrderBook); ok {
				base.SendChan(fullCh, res, "depthCacheChan", res.Symbol)
			}
		}
		if !ok {
			logger.Logger.Error("get DepthCacheChanMap err, symbol:", res.Symbol)
			return errors.New("get DepthCacheChanMap err, symbol:" + res.Symbol)
		}
		h.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, receiveTime)
	} else if strings.HasPrefix(msg.RoomName, "depth_diff") {
		logger.Logger.Trace(string(data))
		res, err := h.parseDepthDelta(data, receiveTime)
		if err != nil {
			logger.Logger.Error("parse depth delta in snapshot", err, " data:", string(data))
			return err
		}
		if obj, ok = h.DepthDeltaUpdateMap.Get(res.Symbol); ok {
			if deltaCh, ok = obj.(chan *base.DeltaDepthUpdate); ok {
				base.SendChan(deltaCh, res, "depthDeltaUpdate", res.Symbol)
			}
		}
		if !ok {
			logger.Logger.Error("get DepthDeltaUpdateMap err, symbol:", res.Symbol)
			return errors.New("get DepthDeltaUpdateMap err, symbol:" + res.Symbol)
		}
		h.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, receiveTime)
	}
	return nil
}
