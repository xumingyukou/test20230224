package spot_ws

import (
	"clients/exchange/cex/base"
	bitMex "clients/exchange/cex/bitmex"
	"clients/exchange/cex/bitmex/bitMex_api"
	"clients/logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/sdk"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
)

type WebSocketHandleInterface interface {
	TradeGroupHandle([]byte) error
	SetTradeGroupChannel(map[*client.SymbolInfo]chan *client.WsTradeRsp)

	DepthIncrementGroupHandle([]byte) error
	DepthIncrementSnapShotGroupHandle([]byte) error

	SetDepthIncrementSnapshotGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp, map[*client.SymbolInfo]chan *depth.Depth)
	SetDepthIncrementSnapShotConf([]*client.SymbolInfo, *base.IncrementDepthConf)
	SetDepthIncrementGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp)

	BookTickerGroupHandle([]byte) error
	SetBookTickerGroupChannel(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)

	GetChan(chName string) interface{}
}

type WebSocketSpotHandle struct {
	reSubscribe chan req

	symbolMap map[string]*client.SymbolInfo
	base.IncrementDepthConf
	aggTradeChan                            chan *client.WsAggTradeRsp
	tradeChan                               chan *client.WsTradeRsp
	depthIncrementChan                      chan *client.WsDepthRsp
	depthLimitChan                          chan *client.WsDepthRsp
	bookTickerChan                          chan *client.WsBookTickerRsp
	accountChan                             chan *client.WsAccountRsp
	balanceChan                             chan *client.WsBalanceRsp
	orderChan                               chan *client.WsOrderRsp
	bookTickerGroupChanMap                  map[string]chan *client.WsBookTickerRsp
	tradeGroupChanMap                       map[string]chan *client.WsTradeRsp
	depthIncrementGroupChanMap              map[string]chan *client.WsDepthRsp //增量
	depthIncrementSnapshotDeltaGroupChanMap map[string]chan *client.WsDepthRsp //增量数据
	depthIncrementSnapshotFullGroupChanMap  map[string]chan *depth.Depth       //全量合并数据
	firstMap                                map[string]bool
	firstFullSentMap                        map[string]bool
	Lock                                    sync.Mutex // 锁
	apiClient                               *bitMex.ClientBitMex
	nameMap                                 map[string]string

	CheckSendStatus *base.CheckDataSendStatus
}

func NewWebSocketSpotHandle(chanCap int64) *WebSocketSpotHandle {
	conf := base.APIConf{}
	//conf = base.APIConf{
	//	ProxyUrl: "http://127.0.0.1:7890",
	//	IsTest:   false,
	//}
	d := &WebSocketSpotHandle{}
	d.reSubscribe = make(chan req, chanCap)
	d.tradeChan = make(chan *client.WsTradeRsp, chanCap)
	d.depthIncrementChan = make(chan *client.WsDepthRsp, chanCap)
	d.depthLimitChan = make(chan *client.WsDepthRsp, chanCap)
	d.bookTickerChan = make(chan *client.WsBookTickerRsp, chanCap)
	d.accountChan = make(chan *client.WsAccountRsp, chanCap)
	d.balanceChan = make(chan *client.WsBalanceRsp, chanCap)
	d.orderChan = make(chan *client.WsOrderRsp, chanCap)
	d.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp)
	d.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp)
	d.firstMap = make(map[string]bool)
	d.apiClient = bitMex.NewClientBitMex(conf)
	d.nameMap = make(map[string]string)
	symbols := d.apiClient.GetSymbols()
	for _, symbol := range symbols {
		d.nameMap[strings.ReplaceAll(symbol, "/", "")] = strings.ReplaceAll(symbol, "XBT", "BTC")
	}
	d.CheckSendStatus = base.NewCheckDataSendStatus()
	d.firstFullSentMap = make(map[string]bool)

	return d
}

func (b *WebSocketSpotHandle) GetChan(chName string) interface{} {
	switch chName {
	case "aggTrade":
		return b.aggTradeChan
	case "trade":
		return b.tradeChan
	case "depthIncrement":
		return b.depthIncrementChan
	case "depthLimit":
		return b.depthLimitChan
	case "bookTicker":
		return b.bookTickerChan
	case "account":
		return b.accountChan
	case "balance":
		return b.balanceChan
	case "order":
		return b.orderChan
	case "reSubscribe":
		return b.reSubscribe
	case "send_status":
		return b.CheckSendStatus
	default:
		return nil
	}
}

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
		conf.CheckTimeSec = 20
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
		)
		conf.DepthCacheMap.Set(GetInstId(symbol), nil)
		conf.DepthCacheListMap.Set(GetInstId(symbol), DepthCacheList)
		conf.CheckDepthCacheChanMap.Set(GetInstId(symbol), CheckDepthCacheChan)
		conf.DepthNotMatchChanMap[symbol] = make(chan bool, conf.DepthCapLevel)
		//go b.Check()
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

func (b *WebSocketSpotHandle) DepthIncrementGroupHandle(data []byte) error {
	var (
		fullRes  = &bitMex_api.RespDepth{}
		depth_   = &client.WsDepthRsp{}
		depthMsg = &Resp_Depth{}
		fullMap  FullDepth
		ok       bool
		sym      string
		err      error
		t        = time.Now().UnixMicro()
	)
	b.Lock.Lock()
	defer b.Lock.Unlock()
	if err = b.HandleRespErr(data, depthMsg); err != nil {
		logger.Logger.Error(err)
		return err
	}
	if strings.Contains(string(data), "subscribe") || strings.Contains(string(data), "success") || strings.Contains(string(data), "partial") ||
		strings.Contains(string(data), "Welcome to the BitMEX Realtime API") {
		return nil
	}

	sym = strings.ReplaceAll(b.nameMap[strings.ReplaceAll(depthMsg.Data[0].Symbol, "XBT", "BTC")], "XBT", "BTC")

	if !b.firstMap[sym] { // 第一次要获取全量
		conf := base.APIConf{}
		c := bitMex_api.NewClientBitMex(conf)
		params := url.Values{}
		params.Add("symbol", depthMsg.Data[0].Symbol)
		fullRes, err = c.DepthInfo(&params)
		if err != nil {
			logger.Logger.Error("获取全量错误", sym)
		}
		// 全量转换为depthUpdate进行排序
		fullMap = FullDepth{Asks: make(map[int64]*depth.DepthLevel), Bids: make(map[int64]*depth.DepthLevel)}
		for _, v := range *fullRes {
			if v.Side == "Sell" {
				fullMap.Asks[v.Id] = &depth.DepthLevel{Price: v.Price, Amount: v.Size}
			} else {
				fullMap.Bids[v.Id] = &depth.DepthLevel{Price: v.Price, Amount: v.Size}
			}
		}
		b.DepthCacheMap.Set(sym, fullMap)
		b.firstMap[sym] = true
	}

	fullMap = FullDepth{Asks: make(map[int64]*depth.DepthLevel), Bids: make(map[int64]*depth.DepthLevel)}
	tmp, ok := b.DepthCacheMap.Get(sym)
	if !ok {
		logger.Logger.Error("本地没有存储:", sym, "币对完整信息")
	}
	fullMap = tmp.(FullDepth)

	//通过id获得增量
	depth_ = getDeltaDepth(fullMap, depthMsg, sym)

	if depthMsg.Action == "delete" {
		for _, v := range depthMsg.Data {
			if v.Side == "Sell" {
				delete(fullMap.Asks, v.Id)
			} else {
				delete(fullMap.Bids, v.Id)
			}
		}
	} else {
		for _, v := range depthMsg.Data {
			if v.Side == "Sell" {
				if _, ok := fullMap.Asks[v.Id]; !ok {
					fullMap.Asks[v.Id] = &depth.DepthLevel{Price: v.Price, Amount: v.Size}
				} else {
					fullMap.Asks[v.Id].Amount = v.Size
				}
			} else {
				if _, ok := fullMap.Bids[v.Id]; !ok {
					fullMap.Bids[v.Id] = &depth.DepthLevel{Price: v.Price, Amount: v.Size}
				} else {
					fullMap.Bids[v.Id].Amount = v.Size
				}
			}
		}
	}
	depth_.ReceiveTime = t
	depth_.Symbol = sym
	if _, ok := b.depthIncrementGroupChanMap[depth_.Symbol]; ok {
		base.SendChan(b.depthIncrementGroupChanMap[depth_.Symbol], depth_, "depthIncrement")
		//b.depthIncrementGroupChanMap[res.Symbol] <- &res
	} else {
		logger.Logger.Warnf("get symbol from channel map err, symbol=%s, resp_symbol=%s", sym, depthMsg.Data[0].Symbol)
	}

	b.CheckSendStatus.CheckUpdateTimeMap.Set(sym, time.Now().UnixMicro())

	return nil
}

// []interface转[][]float64

func (b *WebSocketSpotHandle) DepthLevelParse(levelList_ []levelList) []*depth.DepthLevel {
	var (
		res []*depth.DepthLevel
	)
	for _, level := range levelList_ {
		res = append(res, &depth.DepthLevel{
			Price:  level.Price,
			Amount: level.Price,
		})
	}
	return res
}

func (b *WebSocketSpotHandle) SetDepthIncrementGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) {
	if chMap != nil {
		b.depthIncrementGroupChanMap = make(map[string]chan *client.WsDepthRsp, 1000)
		for k, v := range chMap {
			// 使用斜杠分开的
			b.depthIncrementGroupChanMap[k.Symbol] = v
		}
	}
}

func (b *WebSocketSpotHandle) SetTradeGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) {
	if chMap != nil {
		b.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp, 1000)
		for k, v := range chMap {
			b.tradeGroupChanMap[k.Symbol] = v
		}
	}
}

func (b *WebSocketSpotHandle) TradeGroupHandle(data []byte) error {
	var (
		t    = time.Now().UnixMicro()
		resp Resp_Trade
		err  error
	)

	if strings.Contains(string(data), "subscribe") || strings.Contains(string(data), "success") || strings.Contains(string(data), "partial") ||
		strings.Contains(string(data), "Welcome to the BitMEX Realtime API") {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		logger.Logger.Error(err)
		return err
	}
	symbol := strings.ReplaceAll(b.nameMap[strings.ReplaceAll(resp.Data[0].Symbol, "XBT", "BTC")], "XBT", "BTC")

	var res []*client.WsTradeRsp
	for _, v := range resp.Data {
		tmpres := &client.WsTradeRsp{
			Symbol:       symbol,
			ReceiveTime:  t,
			Price:        v.Price,
			Amount:       v.Size,
			TakerSide:    GetSide(v.Side),
			ExchangeTime: v.Timestamp.UnixMicro(),
		}
		res = append(res, tmpres)
	}

	if _, ok := b.tradeGroupChanMap[symbol]; ok {
		for _, v := range res {
			base.SendChan(b.tradeGroupChanMap[symbol], v, "tradeGroup")
		}
	} else {
		logger.Logger.Warnf("get symbol from channel map err, symbol=%s, resp_symbol=%s", symbol, resp.Data[0].Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(symbol, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) SetBookTickerGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) {
	if chMap != nil {
		b.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp, 1000)
		for k, v := range chMap {
			b.bookTickerGroupChanMap[k.Symbol] = v
		}
	}
}

func (b *WebSocketSpotHandle) BookTickerGroupHandle(data []byte) error {
	if strings.Contains(string(data), "askPrice") || strings.Contains(string(data), "bidPrice") {
		var (
			t      = time.Now().UnixMicro()
			resp   Resp_Ticker
			err    error
			ask    *depth.DepthLevel
			bid    *depth.DepthLevel
			symbol string
		)
		if err = b.HandleRespErr(data, &resp); err != nil {
			logger.Logger.Error(err)
			return err
		}
		symbol = strings.ReplaceAll(b.nameMap[strings.ReplaceAll(resp.Data[0].Symbol, "XBT", "BTC")], "XBT", "BTC")

		ask = &depth.DepthLevel{
			Price: resp.Data[0].AskPrice,
		}
		bid = &depth.DepthLevel{
			Price: resp.Data[0].BidPrice,
		}
		time_ := resp.Data[0].Timestamp.UnixMicro()
		if err != nil {
			logger.Logger.Error("depth level parse err", err, string(data))
			return err
		}

		res := client.WsBookTickerRsp{
			Symbol:       symbol,
			ExchangeTime: time_,
			Ask:          ask,
			Bid:          bid,
			ReceiveTime:  t,
		}
		if len(b.depthLimitChan) > 0 {
			content := fmt.Sprint("channel slow:", len(b.depthLimitChan))
			logger.Logger.Warn(content)
		}
		if _, ok := b.bookTickerGroupChanMap[res.Symbol]; ok {
			base.SendChan(b.bookTickerGroupChanMap[res.Symbol], &res, "booktickerGroup")
			//b.bookTickerGroupChanMap[res.Symbol] <- &res
		} else {
			logger.Logger.Warnf("get symbol from channel map err, symbol=%s, resp_symbol=%s", res.Symbol, resp.Data[0].Symbol)
		}
		b.CheckSendStatus.CheckUpdateTimeMap.Set(symbol, time.Now().UnixMicro())
	}
	return nil

}

func (b *WebSocketSpotHandle) SetDepthIncrementSnapshotGroupChannel(chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chDeltaMap != nil {
		b.depthIncrementSnapshotDeltaGroupChanMap = make(map[string]chan *client.WsDepthRsp, 100)
		for k, v := range chDeltaMap {
			b.depthIncrementSnapshotDeltaGroupChanMap[k.Symbol] = v
		}
	}
	if chFullMap != nil {
		b.depthIncrementSnapshotFullGroupChanMap = make(map[string]chan *depth.Depth, 100)
		for k, v := range chFullMap {
			b.depthIncrementSnapshotFullGroupChanMap[k.Symbol] = v
		}
	}
}

// 通过全量数据反推增量
// 所有的内部处理都是BTC 除了订阅和接受时是XBT
func (b *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte) error {
	var (
		fullRes  = &bitMex_api.RespDepth{}
		depth_   = &client.WsDepthRsp{}
		depthMsg = &Resp_Depth{}
		fullMap  FullDepth
		full     = &depth.Depth{}
		ok       bool
		sym      string
		err      error
		t        = time.Now().UnixMicro()
	)
	b.Lock.Lock()
	defer b.Lock.Unlock()
	if err = b.HandleRespErr(data, depthMsg); err != nil {
		logger.Logger.Error(err)
		return err
	}
	if strings.Contains(string(data), "subscribe") || strings.Contains(string(data), "success") || strings.Contains(string(data), "partial") ||
		strings.Contains(string(data), "Welcome to the BitMEX Realtime API") {
		return nil
	}

	sym = strings.ReplaceAll(b.nameMap[strings.ReplaceAll(depthMsg.Data[0].Symbol, "XBT", "BTC")], "XBT", "BTC")

	if !b.firstMap[sym] { // 第一次要获取全量
		conf := base.APIConf{}
		//conf = base.APIConf{
		//	ProxyUrl: "http://127.0.0.1:7890",
		//	IsTest:   false,
		//}
		c := bitMex_api.NewClientBitMex(conf)
		params := url.Values{}
		params.Add("symbol", depthMsg.Data[0].Symbol)
		fullRes, err = c.DepthInfo(&params)
		if err != nil {
			logger.Logger.Error("获取全量错误", sym)
		}
		t, _ := json.Marshal(params)
		fmt.Println(string(t))
		// 全量转换为depthUpdate进行排序
		fullMap = FullDepth{Asks: make(map[int64]*depth.DepthLevel), Bids: make(map[int64]*depth.DepthLevel)}
		for _, v := range *fullRes {
			if v.Side == "Sell" {
				fullMap.Asks[v.Id] = &depth.DepthLevel{Price: v.Price, Amount: v.Size}
			} else {
				fullMap.Bids[v.Id] = &depth.DepthLevel{Price: v.Price, Amount: v.Size}
			}
		}
		// 如果这次全量是和增量的firstid相同，那么就直接合并,否则直接存储全量
		b.DepthCacheMap.Set(sym, fullMap)
		b.firstMap[sym] = true
		b.firstFullSentMap[sym] = true
	}

	fullMap = FullDepth{Asks: make(map[int64]*depth.DepthLevel), Bids: make(map[int64]*depth.DepthLevel)}
	tmp, ok := b.DepthCacheMap.Get(sym)
	if !ok {
		logger.Logger.Error("本地没有存储:", sym, "币对完整信息")
	}
	fullMap = tmp.(FullDepth)

	//通过id获得增量
	depth_ = getDeltaDepth(fullMap, depthMsg, sym)

	if depthMsg.Action == "delete" {
		for _, v := range depthMsg.Data {
			if v.Side == "Sell" {
				delete(fullMap.Asks, v.Id)
			} else {
				delete(fullMap.Bids, v.Id)
			}
		}
	} else {
		for _, v := range depthMsg.Data {
			if v.Side == "Sell" {
				if _, ok := fullMap.Asks[v.Id]; !ok {
					fullMap.Asks[v.Id] = &depth.DepthLevel{Price: v.Price, Amount: v.Size}
				} else {
					fullMap.Asks[v.Id].Amount = v.Size
				}
				if v.Size == 0 {
					logger.Logger.Error("Amount为0", string(data))
				}
			} else {
				if _, ok := fullMap.Bids[v.Id]; !ok {
					fullMap.Bids[v.Id] = &depth.DepthLevel{Price: v.Price, Amount: v.Size}
				} else {
					fullMap.Bids[v.Id].Amount = v.Size
				}
				if v.Size == 0 {
					logger.Logger.Error("Amount为0", string(data))
				}
			}
		}
	}
	depth_.ReceiveTime = t
	depth_.Symbol = sym
	//增量
	if b.IsPublishDelta {
		if _, ok := b.depthIncrementSnapshotDeltaGroupChanMap[sym]; ok {
			base.SendChan(b.depthIncrementSnapshotDeltaGroupChanMap[sym], depth_, "deltaDepth")
		} else {
			logger.Logger.Warnf("get symbol from channel map err, symbol=%s, resp_symbol=%s", sym, depthMsg.Data[0].Symbol)
		}
	}

	//发送全量
	if b.IsPublishFull {
		full = full2fullDepth(fullMap)
		full.TimeReceive = uint64(t)
		full.Symbol = sym
		if _, ok := b.depthIncrementSnapshotFullGroupChanMap[sym]; ok {
			if b.firstFullSentMap[sym] {
				b.firstFullSentMap[sym] = false
				full.Hdr = base.MakeFirstDepthHdr()
			}
			base.SendChan(b.depthIncrementSnapshotFullGroupChanMap[sym], full, "fullDepth")
		} else {
			logger.Logger.Warnf("get symbol from channel map err, symbol=%s, resp_symbol=%s", sym, depthMsg.Data[0].Symbol)
		}
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(sym, time.Now().UnixMicro())
	return nil
}

func getDeltaDepth(fullMap FullDepth, msg *Resp_Depth, sym string) (res *client.WsDepthRsp) {
	res = &client.WsDepthRsp{}
	asks := base.DepthItemSlice{}
	bids := base.DepthItemSlice{}

	for _, v := range msg.Data {
		if msg.Action == "delete" {
			if v.Side == "Sell" {
				if _, ok := fullMap.Asks[v.Id]; ok {
					asks = append(asks, &depth.DepthLevel{Price: fullMap.Asks[v.Id].Price, Amount: 0})
				}
			} else if v.Side == "Buy" {
				if _, ok := fullMap.Bids[v.Id]; ok {
					bids = append(bids, &depth.DepthLevel{Price: fullMap.Bids[v.Id].Price, Amount: 0})
				}
			} else {
				logger.Logger.Error("未知side:", v.Side)
			}
		} else if msg.Action == "update" {
			if v.Side == "Sell" {
				if _, ok := fullMap.Asks[v.Id]; ok {
					asks = append(asks, &depth.DepthLevel{Price: fullMap.Asks[v.Id].Price, Amount: v.Size})
				}
			} else if v.Side == "Buy" {
				if _, ok := fullMap.Bids[v.Id]; ok {
					bids = append(bids, &depth.DepthLevel{Price: fullMap.Bids[v.Id].Price, Amount: v.Size})
				}
			} else {
				logger.Logger.Error("未知side:", v.Side)
			}
		} else {
			if v.Side == "Sell" {
				asks = append(asks, &depth.DepthLevel{Price: v.Price, Amount: v.Size})
			} else if v.Side == "Buy" {
				bids = append(bids, &depth.DepthLevel{Price: v.Price, Amount: v.Size})
			} else {
				logger.Logger.Error("未知side:", v.Side)
			}
		}
	}
	sort.Stable(asks)
	sort.Stable(sort.Reverse(bids))
	res.Asks = asks
	res.Bids = bids
	return res
}

func full2fullDepth(fullMap FullDepth) (res *depth.Depth) {
	res = &depth.Depth{}
	asks := base.DepthItemSlice{}
	bids := base.DepthItemSlice{}
	for _, v := range fullMap.Asks {
		asks = append(asks, &depth.DepthLevel{Price: v.Price, Amount: v.Amount})
	}
	for _, v := range fullMap.Bids {
		bids = append(bids, &depth.DepthLevel{Price: v.Price, Amount: v.Amount})
	}
	sort.Stable(asks)
	sort.Stable(sort.Reverse(bids))
	res.Bids = bids
	res.Asks = asks
	return
}

func (b *WebSocketSpotHandle) HandleRespErr(data []byte, resp interface{}) error {
	var (
		err error
	)
	err = json.Unmarshal(data, resp)
	if err != nil {
		return errors.New("json unmarshal data err:" + string(data))
	}
	return nil
}
