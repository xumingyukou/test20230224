package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/okex/ok_api"
	"clients/logger"
	"clients/transform"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
	"github.com/warmplanet/proto/go/sdk"
)

type WebSocketHandleInterface interface {
	FundingRateGroupHandle([]byte) error
	SetFundingRateGroupChannel(map[*client.SymbolInfo]chan *client.WsFundingRateRsp)
	//AggTradeHandle([]byte) error
	TradeHandle([]byte) error
	TradeGroupHandle([]byte) error
	SetTradeGroupChannel(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	DepthIncrementHandle([]byte) error
	DepthLimitHandle([]byte) error
	DepthLimitGroupHandle([]byte) error
	DepthIncrementGroupHandle([]byte) error
	DepthIncrementSnapShotGroupHandle([]byte) error

	//DepthLimitHandle2([]byte) error
	BookTickerGroupHandle([]byte) error
	SetBookTickerGroupChannel(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	SetDepthLimitGroupChannel(map[*client.SymbolInfo]chan *depth.Depth)
	SetDepthIncrementGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	SetDepthIncrementSnapshotGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp, map[*client.SymbolInfo]chan *depth.Depth)
	SetDepthIncrementSnapShotConf([]*client.SymbolInfo, *base.IncrementDepthConf)

	BookTickerHandle([]byte) error
	AccountHandle([]byte) error
	OrderHandle([]byte) error

	GetChan(chName string) interface{}
}

type WebSocketSpotHandle struct {
	base.IncrementDepthConf
	reSubscribe                             chan req
	subscribleInfoChan                      chan Resp_Info
	aggTradeChan                            chan *client.WsAggTradeRsp
	tradeChan                               chan *client.WsTradeRsp
	depthIncrementChan                      chan *client.WsDepthRsp
	depthLimitChan                          chan *client.WsDepthRsp
	bookTickerChan                          chan *client.WsBookTickerRsp
	accountChan                             chan *client.WsAccountRsp
	balanceChan                             chan *client.WsBalanceRsp
	orderChan                               chan *client.WsOrderRsp
	fundingRateGroupChanMap                 map[string]chan *client.WsFundingRateRsp
	depthLimitGroupChanMap                  map[string]chan *depth.Depth
	depthIncrementGroupChanMap              map[string]chan *client.WsDepthRsp //增量
	depthIncrementSnapshotDeltaGroupChanMap map[string]chan *client.WsDepthRsp //增量数据
	depthIncrementSnapshotFullGroupChanMap  map[string]chan *depth.Depth       //全量合并数据
	bookTickerGroupChanMap                  map[string]chan *client.WsBookTickerRsp
	tradeGroupChanMap                       map[string]chan *client.WsTradeRsp
	Lock                                    sync.Mutex // 锁
	recon                                   bool
	reMap                                   map[string]bool
	firstFullSentMap                        map[string]bool

	CheckSendStatus *base.CheckDataSendStatus
}

func NewWebSocketSpotHandle(chanCap int64) *WebSocketSpotHandle {
	d := &WebSocketSpotHandle{
		IncrementDepthConf: base.IncrementDepthConf{},
	}
	d.reSubscribe = make(chan req, chanCap)
	d.subscribleInfoChan = make(chan Resp_Info, chanCap)
	d.tradeChan = make(chan *client.WsTradeRsp, chanCap)
	d.depthIncrementChan = make(chan *client.WsDepthRsp, chanCap)
	d.depthLimitChan = make(chan *client.WsDepthRsp, chanCap)
	d.bookTickerChan = make(chan *client.WsBookTickerRsp, chanCap)
	d.accountChan = make(chan *client.WsAccountRsp, chanCap)
	d.balanceChan = make(chan *client.WsBalanceRsp, chanCap)
	d.orderChan = make(chan *client.WsOrderRsp, chanCap)
	d.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp)
	d.depthLimitGroupChanMap = make(map[string]chan *depth.Depth)
	d.reMap = make(map[string]bool)
	d.CheckSendStatus = base.NewCheckDataSendStatus()
	d.firstFullSentMap = make(map[string]bool)
	return d
}

func (b *WebSocketSpotHandle) resetrecon() {
	b.recon = false
}

func (b *WebSocketSpotHandle) isrecon() bool {
	return b.recon
}
func (b *WebSocketSpotHandle) HandleRespErr(data []byte, resp interface{}) error {
	var (
		err error
	)
	if string(data) == "pong" {
		return nil
	}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return errors.New("json unmarshal data err: " + string(data))
	}
	return nil
}

func (b *WebSocketSpotHandle) DepthLevelParse(levelList [][]string) ([]*depth.DepthLevel, error) {
	var (
		res           []*depth.DepthLevel
		amount, price float64
		err           error
	)
	for _, level := range levelList {
		if price, err = strconv.ParseFloat(level[0], 64); err != nil {
			return nil, err
		}
		if amount, err = strconv.ParseFloat(level[1], 64); err != nil {
			return nil, err
		}
		res = append(res, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	return res, err
}

func (b *WebSocketSpotHandle) BookTickerHandle(data []byte) error {
	t := time.Now().UnixMicro()
	var (
		resp RespLimitDepthStreamp
		err  error
		ask  *depth.DepthLevel
		bid  *depth.DepthLevel
	)
	if string(data) == "pong" {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if len(resp.Data) == 0 {
		return nil
	}
	asks, err := b.DepthLevelParse(resp.Data[0].Asks)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	bids, err := b.DepthLevelParse(resp.Data[0].Bids)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	if len(asks) > 0 && len(bids) > 0 {
		ask = asks[0]
		bid = bids[0]
	} else {
		logger.Logger.Warn("asks or bids len is less than 0:", string(data))
		return nil
	}
	res := client.WsBookTickerRsp{
		Symbol:       ok_api.ParseSymbolName(resp.Arg.InstId),
		UpdateIdEnd:  int64(resp.Data[0].Ts),
		ExchangeTime: time.Now().UnixMicro(),
		Ask:          ask,
		Bid:          bid,
		ReceiveTime:  t,
	}
	if len(b.depthLimitChan) > 0 {
		content := fmt.Sprint("channel slow:", len(b.depthLimitChan))
		logger.Logger.Warn(content)
	}
	base.SendChan(b.bookTickerChan, &res, "bookTicker")
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, t)
	//b.bookTickerChan <- &res
	return nil
}

func (b *WebSocketSpotHandle) BookTickerGroupHandle(data []byte) error {
	t := time.Now().UnixMicro()
	var (
		resp RespLimitDepthStreamp
		err  error
		ask  *depth.DepthLevel
		bid  *depth.DepthLevel
	)
	if strings.Contains(string(data), "pong") {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		if string(data) != "pong" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	if len(resp.Data) == 0 {
		return nil
	}
	asks, err := b.DepthLevelParse(resp.Data[0].Asks)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	bids, err := b.DepthLevelParse(resp.Data[0].Bids)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	if len(asks) > 0 && len(bids) > 0 {
		ask = asks[0]
		bid = bids[0]
	} else {
		logger.Logger.Warn("asks or bids len is less than 0:", string(data))
		return nil
	}

	symbolInfo := transSymbolToSymbolInfo(resp.Arg.InstId)
	res := client.WsBookTickerRsp{
		Symbol:       base.SymInfoToString(symbolInfo),
		UpdateIdEnd:  0,
		ExchangeTime: int64(resp.Data[0].Ts) * 1000,
		Ask:          ask,
		Bid:          bid,
		ReceiveTime:  t,
	}
	if len(b.depthLimitChan) > 0 {
		content := fmt.Sprint("channel slow:", len(b.depthLimitChan))
		logger.Logger.Warn(content)
	}
	if _, ok := b.bookTickerGroupChanMap[res.Symbol]; ok {
		// 更新张数
		TransContractSize(res.Symbol, b.Market, nil, res.Ask)
		TransContractSize(res.Symbol, b.Market, nil, res.Bid)
		base.SendChan(b.bookTickerGroupChanMap[res.Symbol], &res, "booktickerGroup")
		//b.bookTickerGroupChanMap[res.Symbol] <- &res
	} else {
		logger.Logger.Warn("ticker get symbol from channel map err:", res.Symbol, b.bookTickerGroupChanMap)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, t)
	return nil
}

func transferDiffDepth(r *RespLimitDepthStreamp, diff *base.DeltaDepthUpdate) {
	// 将binance返回的结构，解析为DeltaDepthUpdate，并将bids和ask进行排序
	diff.Symbol, diff.Market, diff.Type = ok_api.GetContractType(r.Arg.InstId)
	diff.TimeExchange = int64(r.Data[0].Ts) * 1000
	diff.TimeReceive = time.Now().UnixMicro()
	if diff.TimeReceive-diff.TimeExchange > 10000 {
		//fmt.Println(diff.Symbol, diff.TimeExchange, diff.TimeReceive, diff.TimeReceive-diff.TimeExchange)
	}
	diff.UpdateEndId = int64(r.Data[0].Checksum)
	//清空bids，asks
	diff.Bids = diff.Bids[:0]
	diff.Asks = diff.Asks[:0]
	ParseOrder(r.Data[0].Bids, &diff.Bids)
	ParseOrder(r.Data[0].Asks, &diff.Asks)
	sort.Stable(diff.Asks)
	sort.Stable(sort.Reverse(diff.Bids))
	return
}

func ParseOrder(orders [][]string, slice *base.DepthItemSlice) {
	for _, order := range orders {
		price, amount, err := transform.ParsePriceAmountFloat(order)
		if err != nil {
			logger.Logger.Errorf("order float parse price error [%s] , response data = %s", err, order)
			continue
		}
		*slice = append(*slice, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
}

// DepthIncrementSnapShotGroupHandle
// @Description: 根据全量和增量返回数据，固定容量为1000。 数据校验/读取失败时会取消订阅并重新订阅。在获取全量数据时将数据清零
func (b *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte) error {
	t := time.Now().UnixMicro()
	var (
		res = &base.OrderBook{
			TimeReceive: uint64(time.Now().UnixMicro()),
		}
		resp       RespLimitDepthStreamp
		deltaDepth = &base.DeltaDepthUpdate{}
		err        error
		orderBook_ base.OrderBook
		depth_     depth.Depth
	)
	b.Lock.Lock()
	defer b.Lock.Unlock()
	if strings.Contains(string(data), "pong") {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		if string(data) != "pong" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	if resp.Event == "login" {
		base.SendChan(b.subscribleInfoChan, resp.Resp_Info, "respInfo")
		//b.subscribleInfoChan <- resp.Resp_Info
	}
	if len(resp.Data) == 0 {
		return nil
	}

	// 获取数据并且将数据进行排序
	transferDiffDepth(&resp, deltaDepth)
	symbolInfo := transSymbolToSymbolInfo(resp.Arg.InstId)
	symbolInfoStr := base.SymInfoToString(symbolInfo)
	if strings.Contains(string(data), "snapshot") {
		b.DepthCacheMap.Remove(base.SymInfoToString(symbolInfo))
		orderBook_ = base.OrderBook{
			Symbol:   symbolInfo.Symbol,
			Market:   symbolInfo.Market,
			Type:     symbolInfo.Type,
			Bids:     deltaDepth.Bids,
			Asks:     deltaDepth.Asks,
			UpdateId: deltaDepth.UpdateEndId,
		}
		b.reMap[symbolInfoStr] = false
		b.firstFullSentMap[symbolInfoStr] = true
		// 初始全量校验失败，使用通道传过去重新订阅
		if !b.checkFullDepth(&orderBook_) {
			b.recoon(GetInstId(symbolInfo))
			return errors.New("Check failure")
		}
		b.DepthCacheMap.Set(symbolInfoStr, &orderBook_)
		return nil
	}
	if b.reMap[symbolInfoStr] {
		return nil
	}

	res_, ok := b.DepthCacheMap.Get(symbolInfoStr)

	if !ok {
		b.recoon(GetInstId(symbolInfo))
		return errors.New("不能从副本中获取全量数据")
	}

	res = res_.(*base.OrderBook)
	// 更新数据给depth_, depth_转化为orderbook
	base.UpdateBidsAndAsks(deltaDepth, res, 1000, &depth_)
	transDepth2OrderBook(&depth_, res)
	b.DepthCacheMap.Set(symbolInfoStr, res)

	// 副本校验
	if !b.checkFullDepth(res) {
		// 使用通道传过去重新订阅
		b.recoon(GetInstId(symbolInfo))
		return errors.New("Check failure")
	}

	// 发送增量数据
	if b.IsPublishDelta {
		if _, ok := b.depthIncrementSnapshotDeltaGroupChanMap[symbolInfoStr]; ok {
			// 更新张数
			TransContractSize(symbolInfoStr, deltaDepth.Market, nil, deltaDepth.Asks...)
			TransContractSize(symbolInfoStr, deltaDepth.Market, nil, deltaDepth.Bids...)
			base.SendChan(b.depthIncrementSnapshotDeltaGroupChanMap[symbolInfoStr], deltaDepth.Transfer2Depth(), "deltaDepth")
			//b.depthIncrementSnapshotDeltaGroupChanMap[sym] <- deltaDepth.Transfer2Depth()
		} else {
			logger.Logger.Warn("depth incr snapshot get symbol from channel map err:", symbolInfoStr, b.depthIncrementSnapshotDeltaGroupChanMap)
		}
	}
	//发送全量
	if b.IsPublishFull {
		if _, ok := b.depthIncrementSnapshotFullGroupChanMap[symbolInfoStr]; ok {
			if b.firstFullSentMap[symbolInfoStr] {
				b.firstFullSentMap[symbolInfoStr] = false
				depth_.Hdr = base.MakeFirstDepthHdr()
			}
			// 更新张数
			depthCopy := base.DepthDeepCopyByCustom(&depth_)
			TransContractSize(symbolInfoStr, depthCopy.Market, nil, depthCopy.Asks...)
			TransContractSize(symbolInfoStr, depthCopy.Market, nil, depthCopy.Bids...)
			base.SendChan(b.depthIncrementSnapshotFullGroupChanMap[symbolInfoStr], depthCopy, "fullDepth")
		} else {
			logger.Logger.Warn("depth incr snapshot get symbol from channel map err:", symbolInfoStr, b.depthIncrementSnapshotDeltaGroupChanMap)
		}
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(symbolInfoStr, t)
	return nil
}

// 重新订阅
func (d *WebSocketSpotHandle) recoon(symbol string) {
	//d.DepthCacheMap.Set(transformSymbol(symbol), nil)
	d.recon = true
	d.reMap[symbol] = true
	reqx := req{Op: WsUnSubscrible,
		Args: []*argDates{
			{
				Channel: DeepGear400,
				InstId:  symbol,
			},
		},
	}
	base.SendChan(d.reSubscribe, reqx, "reSubscribe")
}

func transDepth2OrderBook(d *depth.Depth, res *base.OrderBook) {
	res.Asks = d.Asks
	res.Bids = d.Bids
	res.TimeReceive = d.TimeReceive
	res.TimeExchange = d.TimeExchange
}

func (d *WebSocketSpotHandle) checkFullDepth(dep *base.OrderBook) bool {
	/**
	rest api返回1000条数据，但是由于程序开始也是1000条，中间数据抛弃后，会有新的数据在rest api中与增量数据不同，所以不能对比所有，要限定数量
	*/
	lb, la := len(dep.Bids), len(dep.Asks)
	var fields []string
	for i := 0; i < 25; i++ {
		if i < lb {
			e := dep.Bids[i]
			fields = append(fields, ParseFloat(e.Price), ParseFloat(e.Amount))
		}
		if i < la {
			e := dep.Asks[i]
			fields = append(fields, ParseFloat(e.Price), ParseFloat(e.Amount))
		}
	}
	raw := strings.Join(fields, ":")
	cs := crc32.ChecksumIEEE([]byte(raw))
	return int32(dep.UpdateId) == int32(cs)
}

func (b *WebSocketSpotHandle) update(newDate *base.DeltaDepthUpdate, depthCache sdk.ConcurrentMapI) {
	/*
		1、从通道获取增量数据
		2、更新缓存的slice，缓存容量限制
		3、更新前对更新量对price排序，对相同price的depth查找并替换交易量，未找到则新增（归并）
		4、depth容量限制，减少1000条导致的性能占用
		5、如果cache中没有对应的量(第一次传入)直接赋值
		6、如果失败直接重新订阅
	*/
	var (
		content interface{}
		ok      bool
	)
	if content, ok = depthCache.Get(transformSymbol(newDate.Symbol)); !ok {
		tmp := depth.Depth{
			TimeReceive: uint64(newDate.TimeReceive),
			Asks:        newDate.Asks,
			Bids:        newDate.Bids,
		}
		depthCache.Set(transformSymbol(newDate.Symbol), &tmp)
		return
	}
	// 全量数据
	content_ := content.(*depth.Depth)
	// 有值，合并后进行校验
	content_.Asks = merge(content_.Asks, newDate.Asks)
	content_.Bids = merge(content_.Bids, newDate.Bids)
	// 翻转bids为降序
	for i, j := 0, len(content_.Bids)-1; i < j; {
		content_.Bids[i], content_.Bids[j] = content_.Bids[j], content_.Bids[i]
		i++
		j--
	}
	if len(content_.Asks) > 1000 {
		content_.Asks = content_.Asks[:1000]
	}
	if len(content_.Bids) > 1000 {
		content_.Bids = content_.Bids[:1000]
	}
	depthCache.Set(transformSymbol(newDate.Symbol), content_)
}

// merge
// @Description: 合并两个有序数组，返回一个新的数组
// @param asks 全量
// @param asks2 增量
func merge(asks []*depth.DepthLevel, asks2 base.DepthItemSlice) []*depth.DepthLevel {
	res := []*depth.DepthLevel{}
	len1 := len(asks)
	len2 := len(asks2)
	l1, l2 := 0, 0
	// 归并排序
	for {
		if l1 == len1 {
			res = append(res, asks2[l2:]...)
			break
		}
		if l2 == len2 {
			res = append(res, asks[l1:]...)
			break
		}
		tmp1 := asks[l1]
		tmp2 := asks2[l2]
		if tmp1.Price == tmp2.Price {
			if tmp2.Amount != 0 {
				tmp := depth.DepthLevel{
					Price:  tmp2.Price,
					Amount: tmp2.Amount,
				}
				res = append(res, &tmp)
			}
			l1++
			l2++
		} else if tmp1.Price < tmp2.Price {
			res = append(res, tmp1)
			l1++
		} else if tmp1.Price > tmp2.Price {
			res = append(res, tmp2)
			l2++
		}
	}
	return res
}

func (b *WebSocketSpotHandle) SetBookTickerGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) {
	if chMap != nil {
		b.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp)
		for k, v := range chMap {
			sym := base.SymInfoToString(k)
			b.bookTickerGroupChanMap[sym] = v
		}
	}
}

func (b *WebSocketSpotHandle) SetDepthLimitGroupChannel(chMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chMap != nil {
		for k, v := range chMap {
			sym := base.SymInfoToString(k)
			b.depthLimitGroupChanMap[sym] = v
		}
	}
}

func (b *WebSocketSpotHandle) SetDepthIncrementGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) {
	if chMap != nil {
		b.depthIncrementGroupChanMap = make(map[string]chan *client.WsDepthRsp)
		for k, v := range chMap {
			// 使用斜杠分开的
			sym := base.SymInfoToString(k)
			b.depthIncrementGroupChanMap[sym] = v
		}
	}
}

func (b *WebSocketSpotHandle) SetTradeGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) {
	if chMap != nil {
		for k, v := range chMap {
			sym := base.SymInfoToString(k)
			b.tradeGroupChanMap[sym] = v
		}
	}
}

func (b *WebSocketSpotHandle) SetDepthIncrementSnapshotGroupChannel(chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chDeltaMap != nil {
		b.depthIncrementSnapshotDeltaGroupChanMap = make(map[string]chan *client.WsDepthRsp)
		for k, v := range chDeltaMap {
			sym := base.SymInfoToString(k)
			b.depthIncrementSnapshotDeltaGroupChanMap[sym] = v
		}
	}
	if chFullMap != nil {
		b.depthIncrementSnapshotFullGroupChanMap = make(map[string]chan *depth.Depth)
		for k, v := range chFullMap {
			sym := base.SymInfoToString(k)
			b.depthIncrementSnapshotFullGroupChanMap[sym] = v
		}
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
		conf.DepthCacheMap.Set(base.SymInfoToString(symbol), nil)
		conf.DepthCacheListMap.Set(base.SymInfoToString(symbol), DepthCacheList)
		conf.CheckDepthCacheChanMap.Set(base.SymInfoToString(symbol), CheckDepthCacheChan)
		conf.DepthNotMatchChanMap[symbol] = make(chan bool, conf.DepthCapLevel)
		b.firstFullSentMap[base.SymInfoToString(symbol)] = true
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

func (b *WebSocketSpotHandle) parseTrade(data []byte) (*client.WsTradeRsp, error) {
	var (
		resp          RespTradeStream
		err           error
		price, amount float64
	)
	if string(data) == "pong" {
		return nil, nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		return nil, err
	}
	if len(resp.Data) <= 0 {
		return nil, nil
	}
	price, err = strconv.ParseFloat(resp.Data[0].Px, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return nil, err
	}
	amount, err = strconv.ParseFloat(resp.Data[0].Sz, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return nil, err
	}

	symbolInfo := transSymbolToSymbolInfo(resp.Data[0].InstId)
	res := &client.WsTradeRsp{
		//OrderId: resp.Data[0].TradeId,
		Symbol: base.SymInfoToString(symbolInfo),
		//ExchangeTime: resp.Data.Ts,
		Price:         price,
		Amount:        amount,
		TakerSide:     GetSide(resp.Data[0].Side),
		BuyerOrderId:  getUserId(resp.Data[0].TradeId, resp.Data[0].Side, "buy"),
		SellerOrderId: getUserId(resp.Data[0].TradeId, resp.Data[0].Side, "sell"),
		DealTime:      int64(resp.Data[0].Ts) * 1000,
		ReceiveTime:   time.Now().UnixMicro(),
	}
	// 更新张数
	TransContractSize(res.Symbol, b.Market, res)
	return res, nil
}

func (b *WebSocketSpotHandle) TradeGroupHandle(data []byte) error {
	if strings.Contains(string(data), "pong") {
		return nil
	}
	res, err := b.parseTrade(data)
	if err != nil {
		return err
	}
	if res == nil {
		return nil
	}
	if _, ok := b.tradeGroupChanMap[res.Symbol]; ok {
		base.SendChan(b.tradeGroupChanMap[res.Symbol], res, "tradeGroup")
		//b.tradeGroupChanMap[res.Symbol] <- res
	} else {
		logger.Logger.Warn("trade get symbol from channel map err:", res.Symbol, b.tradeGroupChanMap)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) TradeHandle(data []byte) error {
	res, err := b.parseTrade(data)
	if err != nil {
		return err
	}
	if len(b.tradeChan) > 0 {
		content := fmt.Sprint("channel slow:", len(b.tradeChan))
		logger.Logger.Warn(content)
	}
	b.tradeChan <- res
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) DepthIncrementHandle(data []byte) error {
	t := time.Now().UnixMicro()
	var (
		resp RespLimitDepthStreamp
		err  error
	)
	if err = b.HandleRespErr(data, &resp); err != nil {
		if string(data) != "pong" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	if resp.Event == "login" {
		base.SendChan(b.subscribleInfoChan, resp.Resp_Info, "subscribeInfo")
		//b.subscribleInfoChan <- resp.Resp_Info
	}
	if len(resp.Data) == 0 {
		return nil
	}
	asks, err := b.DepthLevelParse(resp.Data[0].Asks)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	bids, err := b.DepthLevelParse(resp.Data[0].Bids)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	res := client.WsDepthRsp{
		Symbol:      ok_api.ParseSymbolName(resp.Arg.InstId),
		Asks:        asks,
		Bids:        bids,
		ReceiveTime: t,
	}
	if len(b.depthLimitChan) > 0 {
		content := fmt.Sprint("channel slow:", len(b.depthLimitChan))
		logger.Logger.Warn(content)
	}
	b.depthIncrementChan <- &res

	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, t)
	return nil
}

func (b *WebSocketSpotHandle) DepthIncrementGroupHandle(data []byte) error {
	t := time.Now().UnixMicro()
	var (
		resp RespLimitDepthStreamp
		err  error
	)
	if strings.Contains(string(data), "pong") {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		return err
	}
	if resp.Event == "login" {
		base.SendChan(b.subscribleInfoChan, resp.Resp_Info, "subscribeInfo")
		//b.subscribleInfoChan <- resp.Resp_Info
	}
	if len(resp.Data) == 0 {
		return nil
	}
	asks, err := b.DepthLevelParse(resp.Data[0].Asks)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	bids, err := b.DepthLevelParse(resp.Data[0].Bids)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	res := client.WsDepthRsp{
		Symbol:       ok_api.ParseSymbolName(resp.Arg.InstId),
		Asks:         asks,
		Bids:         bids,
		ExchangeTime: t,
	}

	if _, ok := b.depthIncrementGroupChanMap[res.Symbol]; ok {
		base.SendChan(b.depthIncrementGroupChanMap[res.Symbol], &res, "depthIncrement")
		//b.depthIncrementGroupChanMap[res.Symbol] <- &res
	} else {
		logger.Logger.Warn("depth incr get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, t)
	return nil
}

func (b *WebSocketSpotHandle) DepthLimitHandle(data []byte) error {
	t := time.Now().UnixMicro()
	var (
		resp RespLimitDepthStreamp
		err  error
	)
	if err = b.HandleRespErr(data, &resp); err != nil {
		if string(data) != "pong" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	if resp.Event == "login" {
		base.SendChan(b.subscribleInfoChan, resp.Resp_Info, "subscribeInfo")
		//b.subscribleInfoChan <- resp.Resp_Info
	}
	if len(resp.Data) == 0 {
		return nil
	}

	asks, err := b.DepthLevelParse(resp.Data[0].Asks)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	bids, err := b.DepthLevelParse(resp.Data[0].Bids)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	res := client.WsDepthRsp{
		Symbol:      ok_api.ParseSymbolName(resp.Arg.InstId),
		Asks:        asks,
		Bids:        bids,
		ReceiveTime: t,
	}
	if len(b.depthLimitChan) > 0 {
		content := fmt.Sprint("channel slow:", len(b.depthLimitChan))
		logger.Logger.Warn(content)
	}
	b.depthLimitChan <- &res
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, t)
	return nil
}

func (b *WebSocketSpotHandle) DepthLimitGroupHandle(data []byte) error {
	t := time.Now().UnixMicro()
	var (
		resp RespLimitDepthStreamp
		err  error
	)
	if strings.Contains(string(data), "pong") {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		return err
	}
	if resp.Event == "login" {
		base.SendChan(b.subscribleInfoChan, resp.Resp_Info, "subscribeInfo")
		//b.subscribleInfoChan <- resp.Resp_Info
	}
	if len(resp.Data) == 0 {
		return nil
	}

	asks, err := b.DepthLevelParse(resp.Data[0].Asks)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	bids, err := b.DepthLevelParse(resp.Data[0].Bids)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	res := depth.Depth{
		Exchange:    common.Exchange_OKEX,
		Symbol:      ok_api.ParseSymbolName(resp.Arg.InstId),
		Asks:        asks,
		Bids:        bids,
		TimeReceive: uint64(t),
		TimeOperate: uint64(t),
	}
	if len(b.depthLimitChan) > 0 {
		content := fmt.Sprint("channel slow:", len(b.depthLimitChan))
		logger.Logger.Warn(content)
	}
	if _, ok := b.depthLimitGroupChanMap[res.Symbol]; ok {
		base.SendChan(b.depthLimitGroupChanMap[res.Symbol], &res, "depthLimit")
		//b.depthLimitGroupChanMap[res.Symbol] <- &res
	} else {
		logger.Logger.Warn("depth limit get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, t)
	return nil
}

func (b *WebSocketSpotHandle) AccountHandle(data []byte) error {
	var (
		resp                                                                        RespUserAccount
		err                                                                         error
		position, price, available, lock, unprofit, totalAsset, totalLiabilityAsset float64
		instType                                                                    common.Market
		side                                                                        order.TradeSide
	)
	// 返回ping的信息省略
	if string(data) == "pong" {
		return nil
	}

	if err = b.HandleRespErr(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if resp.Event == "login" {
		base.SendChan(b.subscribleInfoChan, resp.Resp_Info, "respInfo")
		//b.subscribleInfoChan <- resp.Resp_Info
	}
	b.subscribleInfoChan <- resp.Resp_Info
	// 省略刚订阅的时候返回的信息 {"op":"subscribe","args":[{"channel":"account","instId":""}]}
	if len(resp.Data) < 1 {
		return nil
	}
	res := client.WsAccountRsp{
		UpdateTime: int64(resp.Data[0].UTime) * 1000,
		// ok使用usd计价
		QuoteAsset: "USD",
	}
	if resp.Arg.Channel == "account" {
		totalAsset, err = strconv.ParseFloat(resp.Data[0].TotalEq, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		if resp.Data[0].AdjEq == "" {
			resp.Data[0].AdjEq = "0"
		}
		totalLiabilityAsset, err = strconv.ParseFloat(resp.Data[0].AdjEq, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		// 总权益
		res.TotalAsset = totalAsset
		// 美金层面有效保证金
		res.TotalLiabilityAsset = totalLiabilityAsset
		res.TotalNetAsset = totalAsset - totalLiabilityAsset

		for _, item := range resp.Data[0].Details {
			available, err = strconv.ParseFloat(item.CashBal, 64)
			if err != nil {
				logger.Logger.Error("trade parse err", err, string(data))
				return err
			}
			lock, err = strconv.ParseFloat(item.FrozenBal, 64)
			if err != nil {
				logger.Logger.Error("trade parse err", err, string(data))
				return err
			}
			res.BalanceList = append(res.BalanceList, &client.BalanceItem{
				Assert:       strings.ReplaceAll(item.Ccy, "-", "/"),
				Available:    available,
				Lock:         lock,
				Market:       common.Market_SPOT,
				ExchangeTime: int64(item.UTime) * 1000,
			})
		}
	} else if resp.Arg.Channel == "positions" {
		instType = ok_api.TransformMarket(resp.Arg.InstType)
		for _, item := range resp.Data {
			position, err = strconv.ParseFloat(item.Pos, 64)
			if err != nil {
				logger.Logger.Error("trade parse err", err, string(data))
				return err
			}
			price, err = strconv.ParseFloat(item.AvgPx, 64)
			if err != nil {
				logger.Logger.Error("trade parse err", err, string(data))
				return err
			}
			unprofit, err = strconv.ParseFloat(item.Upl, 64)
			if err != nil {
				logger.Logger.Error("trade parse err", err, string(data))
				return err
			}

			side = order.TradeSide_BUY
			if position < 0 {
				side = order.TradeSide_SELL
			}
			res.PositionList = append(res.PositionList, &client.PositionItem{
				Position:     position,
				Symbol:       strings.ReplaceAll(item.Ccy, "-", "/"),
				ExchangeTime: int64(item.CTime) * 1000,
				Price:        price,
				Unprofit:     unprofit,
				Market:       instType,
				Side:         side,
			})
		}
	}
	if len(b.accountChan) > 0 {
		content := fmt.Sprint("channel slow:", len(b.accountChan))
		logger.Logger.Warn(content)
	}
	b.accountChan <- &res
	return nil
}

func (b *WebSocketSpotHandle) BalanceHandle(data []byte) error {
	return nil
}

func (b *WebSocketSpotHandle) OrderHandle(data []byte) error {
	var (
		resp                                      RespUserOrder
		err                                       error
		clientId                                  string
		feeAsset                                  string
		amountFilled, priceFilled, qtyFilled, fee float64
	)
	if string(data) == "pong" {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if resp.Event == "login" {
		base.SendChan(b.subscribleInfoChan, resp.Resp_Info, "subscribeInfo")
		//b.subscribleInfoChan <- resp.Resp_Info
	}

	if err = b.HandleRespErr(data, &resp); err != nil {
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	for _, item := range resp.Data {
		clientId = item.ClOrdId
		amountFilled, err = strconv.ParseFloat(item.AccFillSz, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		priceFilled, err = strconv.ParseFloat(item.AvgPx, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		fee, err = strconv.ParseFloat(item.FillFee, 64)
		if fee >= 0 {
			feeAsset = item.RebateCcy
		} else {
			feeAsset = item.FeeCcy
		}
		producer, id := transform.ClientIdToId(clientId)
		var closeDate string
		names := strings.Split(item.InstId, "-")
		if len(names) == 3 {
			if names[2] != "SWAP" {
				closeDate = names[2]
			}
		}

		res := client.WsOrderRsp{
			Producer: producer,

			Id:           id,
			IdEx:         item.OrdId,
			Symbol:       item.InstId,
			Status:       ok_api.GetOrderStatusFromExchange(ok_api.OrderStatus(strings.ToUpper(item.State))),
			AmountFilled: amountFilled,
			PriceFilled:  priceFilled,
			QtyFilled:    qtyFilled,
			Fee:          fee,
			FeeAsset:     feeAsset,
			TimeFilled:   int64(item.UTime),
			CloseData:    closeDate,
		}

		if len(b.orderChan) > 0 {
			content := fmt.Sprint("channel slow:", len(b.orderChan))
			logger.Logger.Warn(content)
		}
		b.orderChan <- &res
	}
	return nil
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
	case "subscribeInfo":
		return b.subscribleInfoChan
	case "reSubscribe":
		return b.reSubscribe
	case "send_status":
		return b.CheckSendStatus
	default:
		return nil
	}
}

type WebSocketUBaseHandle struct {
	WebSocketSpotHandle
}

func NewWebSocketUBaseHandle(chanCap int64) *WebSocketUBaseHandle {
	d := &WebSocketUBaseHandle{}
	d.tradeChan = make(chan *client.WsTradeRsp, chanCap)
	d.depthIncrementChan = make(chan *client.WsDepthRsp, chanCap)
	d.depthLimitChan = make(chan *client.WsDepthRsp, chanCap)
	d.bookTickerChan = make(chan *client.WsBookTickerRsp, chanCap)
	d.accountChan = make(chan *client.WsAccountRsp, chanCap)
	d.balanceChan = make(chan *client.WsBalanceRsp, chanCap)
	d.orderChan = make(chan *client.WsOrderRsp, chanCap)

	return d
}

func (d *WebSocketSpotHandle) MergeDepth(side base.SIDE, fullDepth *base.DepthItemSlice, deltaDepth base.DepthItemSlice) {
	// 1. 与全量价格相等，增量amount不为0，则替换
	// 2. 与全量价格相等，增量amount为0，则删除
	// 3. 全量未找到，增量amount不为0，则插入
	var zeroIdx []int
	i, j := 0, 0
	for i < len(deltaDepth) && j < len(*fullDepth) {
		if (*fullDepth)[j].Price == deltaDepth[i].Price {
			if deltaDepth[i].Amount > 0 {
				(*fullDepth)[j].Amount = deltaDepth[i].Amount
			} else {
				zeroIdx = append(zeroIdx, j)
			}
			i++
			j++
		} else if (side == base.Ask && (*fullDepth)[j].Price < deltaDepth[i].Price) || (side == base.Bid && (*fullDepth)[j].Price > deltaDepth[i].Price) {
			// ask卖盘：全量的price小于增量的price，i不变，j增加
			j++
		} else {
			// ask卖盘：全量的price大于增量的price，需要在全量的j索引前面插入这个amount
			if deltaDepth[i].Amount > 0 {
				if j == 0 {
					*fullDepth = append(base.DepthItemSlice{
						&depth.DepthLevel{
							Price:  deltaDepth[i].Price,
							Amount: deltaDepth[i].Amount,
						},
					}, *fullDepth...)
				} else {
					tmp := append(base.DepthItemSlice{}, (*fullDepth)[:j]...)
					tmp = append(tmp, &depth.DepthLevel{
						Price:  deltaDepth[i].Price,
						Amount: deltaDepth[i].Amount,
					})
					*fullDepth = append(tmp, (*fullDepth)[j:]...)
				}
			} else {
				i++
				continue
			}
			i++
			j++
		}
	}
	//处理剩余的增量数据
	for k := 0; k < len(deltaDepth)-i; k++ {
		if deltaDepth[i+k].Amount > 0 {
			*fullDepth = append(*fullDepth, &depth.DepthLevel{
				Price:  deltaDepth[i+k].Price,
				Amount: deltaDepth[i+k].Amount,
			})
		}
	}
	for i1, idx := range zeroIdx {
		if idx == 0 {
			*fullDepth = append((*fullDepth)[0:0], (*fullDepth)[idx+1:]...)
		} else {
			*fullDepth = append((*fullDepth)[:idx-i1], (*fullDepth)[idx-i1+1:]...)
		}
	}
	return

}

func (b *WebSocketSpotHandle) FundingRateGroupHandle(data []byte) error {
	var (
		resp        RespFundingRate
		err         error
		fundingRate float64
	)
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	if len(resp.Data) <= 0 {
		return nil
	}
	fundingRate, err = strconv.ParseFloat(resp.Data[0].FundingRate, 64)
	symbolInfo := transSymbolToSymbolInfo(resp.Data[0].InstId)
	res := &client.WsFundingRateRsp{
		ExchangeTime: int64(resp.Data[0].FundingTime) * 1000,
		Symbol:       symbolInfo.Symbol,
		Type:         symbolInfo.Type,
		FundingRate:  fundingRate,
		ReceiveTime:  time.Now().UnixMicro(),
	}
	symStr := base.SymInfoToString(symbolInfo)
	if _, ok := b.fundingRateGroupChanMap[symStr]; ok {
		base.SendChan(b.fundingRateGroupChanMap[symStr], res, "fundingRateGroup")
	} else {
		logger.Logger.Warn("FundingRateGroupHandle get symbol from channel map err:", symStr)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(symStr, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) SetFundingRateGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsFundingRateRsp) {
	if chMap != nil {
		b.fundingRateGroupChanMap = make(map[string]chan *client.WsFundingRateRsp)
		for info, ch := range chMap {
			sym := base.SymInfoToString(info)
			b.fundingRateGroupChanMap[sym] = ch
		}
	}
}
