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

func (h *WebSocketSpotHandle) updateDeltaDepth(symbol string) {
	var (
		obj            interface{}
		deltaCh        chan *base.DeltaDepthUpdate
		delta          *base.DeltaDepthUpdate
		full           *base.OrderBook
		deltaCacheList = make([]*base.DeltaDepthUpdate, 0, 10)
		depthCache     = &base.OrderBook{}
		fullCh         chan *base.OrderBook
		deltaSendCh    chan *client.WsDepthRsp
		fullSendCh     chan *depth.Depth
		ok             bool
		funcName       = "updateDeltaDepth"
		getFirstFull   = false
		getFirstDelta  = false
		firstFullSent  = false
		err            error
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
LOOP:
	for {
		select {
		case <-h.Ctx.Done():
			break LOOP
		case delta, ok = <-deltaCh:
			if !getFirstDelta {
				logger.Logger.Infof("receive first delta %d %d %s", delta.UpdateStartId, delta.UpdateEndId, symbol)
				getFirstDelta = true
			}
			if h.IsPublishDelta {
				res := delta.Transfer2Depth()
				base.SendChan(deltaSendCh, res, "deltaPublish", symbol)
			}
			if h.IsPublishFull {
				deltaCacheList = append(deltaCacheList, delta)
				if depthCache.UpdateId == 0 || depthCache.UpdateId+1 < deltaCacheList[0].UpdateStartId {
					firstFullSent = false
					if getFirstFull {
						// 已经收到了增量，仍然不满足条件
						// 重新拉全量
						if depthCache, err = h.GetFullDepth(symbol); err != nil {
							logger.Logger.Error("update get full depth error", err)
						}
					}
					continue
				} else if depthCache.UpdateId >= deltaCacheList[len(deltaCacheList)-1].UpdateEndId {
					firstFullSent = false
					deltaCacheList = append([]*base.DeltaDepthUpdate{})
					continue
				} else {
					target := -1
					for i, delta := range deltaCacheList {
						if depthCache.UpdateId < delta.UpdateStartId {
							target = i
						}
					}
					if target == -1 {
						// 没有找到全量的id
						logger.Logger.Error("deltaList error", symbol)
						if depthCache, err = h.GetFullDepth(symbol); err != nil {
							logger.Logger.Error("update get full depth error", err)
						}
						firstFullSent = false
						continue
					}
					getFirstFull = true
					deltaCacheList = deltaCacheList[target:]
					for _, deltaCache := range deltaCacheList {
						base.UpdateBidsAndAsks(deltaCache, depthCache, h.DepthCapLevel, nil)
					}
					fullToSend := &depth.Depth{}
					depthCache.Transfer2Depth(h.DepthLevel, fullToSend)
					if !firstFullSent { // 首次推送全量检测
						fullToSend.Hdr = base.MakeFirstDepthHdr()
						firstFullSent = true
					}
					base.SendChan(fullSendCh, fullToSend, "depth full send", symbol)
					deltaCacheList = append([]*base.DeltaDepthUpdate{})
				}
			}
		case full, ok = <-fullCh:
			if !getFirstFull {
				logger.Logger.Infof("receive first full %d %s", full.UpdateId, symbol)
				depthCache = full
				getFirstFull = true
			}
		}
	}
	return
}

func (h *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte, receiveTime int64) error {
	var (
		resp    = &OrderBookMsg{}
		err     error
		ok      = false
		deltaCh chan *base.DeltaDepthUpdate
		fullCh  chan *base.OrderBook
		obj     interface{}
	)
	err = json.Unmarshal(data, resp)
	if err != nil {
		logger.Logger.Error("depth increment parse ", err, " data:", string(data))
		return err
	}
	// full
	if resp.Code == "00006" {
		full, err := h.parseDepthWhole(resp, receiveTime)
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
	} else {
		delta, err := h.parseDepthDelta(resp, receiveTime)

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
	return nil
}
