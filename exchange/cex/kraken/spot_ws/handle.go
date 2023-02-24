package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"clients/transform"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"

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

	base.IncrementDepthConf
	subscribleInfoChan                      chan Resp_Info
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
	pPmap                                   map[string]int32
	aPmap                                   map[string]int32
	rMap                                    map[string]bool
	firstFullSentMap                        map[string]bool
	CheckSendStatus                         *base.CheckDataSendStatus
}

func NewWebSocketSpotHandle(chanCap int64) *WebSocketSpotHandle {
	d := &WebSocketSpotHandle{}
	d.reSubscribe = make(chan req, chanCap)
	d.subscribleInfoChan = make(chan Resp_Info, chanCap)
	d.tradeChan = make(chan *client.WsTradeRsp, chanCap)
	d.depthIncrementChan = make(chan *client.WsDepthRsp, chanCap)
	d.depthLimitChan = make(chan *client.WsDepthRsp, chanCap)
	d.bookTickerChan = make(chan *client.WsBookTickerRsp, chanCap)
	d.accountChan = make(chan *client.WsAccountRsp, chanCap)
	d.balanceChan = make(chan *client.WsBalanceRsp, chanCap)
	d.orderChan = make(chan *client.WsOrderRsp, chanCap)
	d.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp)
	d.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp)
	d.aPmap = make(map[string]int32, chanCap)
	d.pPmap = make(map[string]int32, chanCap)
	d.rMap = make(map[string]bool, chanCap)
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

func (b *WebSocketSpotHandle) DepthIncrementGroupHandle(data []byte) error {
	var (
		t            = time.Now().UnixMicro()
		ExchangeTime int64
	)
	if strings.Contains(string(data), "heartbeat") {
		return nil
	}
	// 跳过前两个确认信息
	if strings.Contains(string(data), "event") && strings.Contains(string(data), "book") && strings.Contains(string(data), "subscription") {
		var init_resp map[string]interface{}
		err := json.Unmarshal(data, &init_resp)
		if err != nil {
			log.Fatal(err)
		}
		return err
	}
	var err error
	var resp interface{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Logger.Error("解析错误:", string(data))
		return err

	}
	if resp, ok := resp.([]interface{}); ok {
		var symbol string
		// 得到名称
		if len(resp) == 4 {
			symbol = resp[3].(string)
		} else {
			symbol = resp[4].(string)
		}
		if resp, ok := resp[1].(map[string]interface{}); ok {
			// asks
			var asks [][]float64
			var bids [][]float64
			if ax, ok := resp["as"]; ok {
				tmp := ax.([]interface{})
				asks, ExchangeTime = b.parseInfo(tmp, ExchangeTime)
			}

			// bsks
			if bx, ok := resp["bs"]; ok {
				tmp := bx.([]interface{})
				bids, ExchangeTime = b.parseInfo(tmp, ExchangeTime)
			}

			// a增量
			if a, ok := resp["a"]; ok {
				tmp := a.([]interface{})
				asks, ExchangeTime = b.parseInfo(tmp, ExchangeTime)
			}

			// b增量
			if b_, ok := resp["b"]; ok {
				tmp := b_.([]interface{})
				bids, ExchangeTime = b.parseInfo(tmp, ExchangeTime)
			}
			res := client.WsDepthRsp{
				Symbol:       strings.ReplaceAll(symbol, "XBT", "BTC"),
				Asks:         b.DepthLevelParse(asks),
				Bids:         b.DepthLevelParse(bids),
				ReceiveTime:  t,
				ExchangeTime: ExchangeTime,
			}
			if _, ok := b.depthIncrementGroupChanMap[res.Symbol]; ok {
				base.SendChan(b.depthIncrementGroupChanMap[res.Symbol], &res, "depthIncrement")
				//b.depthIncrementGroupChanMap[res.Symbol] <- &res
			} else {
				logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
			}
			b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
			return err
		}
	}
	return nil
}

// []interface转[][]float64
func (b *WebSocketSpotHandle) parseInfo(x []interface{}, exchangeTime int64) ([][]float64, int64) {
	// todo precision需要确定哪个是哪个
	res := [][]float64{}
	//[[0.0000153100 3000.00000000 1660616521.801961] [0.0000153000 7496.30389313 1660616572.254744]]
	for _, v := range x {
		tmp := []float64{}
		// [0.0000153100 3000.00000000 1660616521.801961]

		if v, ok := v.([]interface{}); ok {
			// [0.0000153100 3000.00000000 1660616521.801961]
			//if len(v) > 3 && v[3].(string) == "r" {
			//
			//	continue
			//}
			for i, v := range v {
				if i > 1 {
					t := ParseI(v)
					if t > exchangeTime {
						exchangeTime = t
					}
					continue
				}

				if v, ok := v.(string); ok {
					i, _ := strconv.ParseFloat(v, 64)
					tmp = append(tmp, i)
				}
			}
		}
		res = append(res, tmp)
	}
	return res, exchangeTime
}

func (b *WebSocketSpotHandle) DepthLevelParse(levelList [][]float64) []*depth.DepthLevel {
	var (
		res []*depth.DepthLevel
	)
	for _, level := range levelList {
		res = append(res, &depth.DepthLevel{
			Price:  level[0],
			Amount: level[1],
		})
	}
	return res
}

func (b *WebSocketSpotHandle) SetDepthIncrementGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) {
	if chMap != nil {
		b.depthIncrementGroupChanMap = make(map[string]chan *client.WsDepthRsp)
		for k, v := range chMap {
			// 使用斜杠分开的
			b.depthIncrementGroupChanMap[GetInstId(k)] = v
		}
	}
}

func (b *WebSocketSpotHandle) SetTradeGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) {
	if chMap != nil {
		b.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp)
		for k, v := range chMap {
			b.tradeGroupChanMap[GetInstId(k)] = v
		}
	}
}

func (b *WebSocketSpotHandle) TradeGroupHandle(data []byte) error {
	var (
		err    error
		symbol string
		t      = time.Now().UnixMicro()
	)
	if strings.Contains(string(data), "heartbeat") {
		return nil
	}
	// 跳过前两个确认信息
	if strings.Contains(string(data), "event") && strings.Contains(string(data), "status") && strings.Contains(string(data), "version") {
		var init_resp map[string]interface{}
		err := json.Unmarshal(data, &init_resp)
		if err != nil {
			log.Fatal(err)
		}
		if !(init_resp["event"] == "systemStatus" && init_resp["status"] == "online") {
			log.Fatal(init_resp)
		}
		return err
	}
	var resp interface{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}
	if resp, ok := resp.([]interface{}); ok {
		// 得到名称
		symbol = resp[3].(string)
		symbol = strings.ReplaceAll(symbol, "XBT", "BTC")
		//fmt.Println(symbol)
		// asks
		// [0,[[2,3,4,"p"],[3,4,5,"s"]],"trade","xbt/usd"]
		if resp, ok := resp[1].([]interface{}); ok {
			//[[2,3,4,"p"],[3,4,5,"s"]]
			for _, v := range resp {
				// [2,3,4,"p"]
				if resp, ok := v.([]interface{}); ok {

					price := ParseF(resp[0])
					amount := ParseF(resp[1])
					dealTime := ParseI(resp[2])

					res := &client.WsTradeRsp{
						Symbol:       symbol,
						Price:        price,
						Amount:       amount,
						TakerSide:    GetSide(resp[3].(string)),
						ExchangeTime: dealTime,
						ReceiveTime:  t,
					}
					if _, ok := b.tradeGroupChanMap[res.Symbol]; ok {
						base.SendChan(b.tradeGroupChanMap[res.Symbol], res, "tradeGroup")
					} else {
						logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
					}
					b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
				}
			}
			return nil
		}
	}
	return err
}

func (b *WebSocketSpotHandle) SetBookTickerGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) {
	if chMap != nil {
		b.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp)
		for k, v := range chMap {
			b.bookTickerGroupChanMap[GetInstId(k)] = v
		}
	}
}

// 无exchangetime
func (b *WebSocketSpotHandle) BookTickerGroupHandle(data []byte) error {
	// 跳过前两个确认信息
	var (
		t = time.Now().UnixMicro()
	)
	if strings.Contains(string(data), "heartbeat") {
		return nil
	}
	if strings.Contains(string(data), "event") && strings.Contains(string(data), "status") && strings.Contains(string(data), "version") {
		var init_resp map[string]interface{}
		err := json.Unmarshal(data, &init_resp)
		if err != nil {
			log.Fatal(err)
		}
		if !(init_resp["event"] == "systemStatus" && init_resp["status"] == "online") {
			log.Fatal(init_resp)
		}
		return err
	}
	var err error
	var resp interface{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Logger.Error("解析错误:", string(data))
		return err

	}
	//fmt.Println("res:", string(data))
	if resp, ok := resp.([]interface{}); ok {
		// 得到名称
		symbol := resp[3].(string)
		symbol = strings.ReplaceAll(symbol, "XBT", "BTC")
		//fmt.Println(symbol)
		if resp, ok := resp[1].(map[string]interface{}); ok {
			// ask
			var ask *depth.DepthLevel
			var bid *depth.DepthLevel
			if ax, ok := resp["a"]; ok {
				tmp := ax.([]interface{})
				ask = &depth.DepthLevel{
					Price:  ParseF(tmp[0]),
					Amount: ParseF(tmp[2]),
				}
			}

			// bid
			if bx, ok := resp["b"]; ok {
				tmp := bx.([]interface{})
				bid = &depth.DepthLevel{
					Price:  ParseF(tmp[0]),
					Amount: ParseF(tmp[2]),
				}
			}

			res := client.WsBookTickerRsp{
				Symbol:      symbol,
				Ask:         ask,
				Bid:         bid,
				ReceiveTime: t,
			}
			if _, ok := b.bookTickerGroupChanMap[res.Symbol]; ok {
				base.SendChan(b.bookTickerGroupChanMap[res.Symbol], &res, "depthIncrement")
				//b.depthIncrementGroupChanMap[res.Symbol] <- &res
			} else {
				logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
			}
			b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())

			return err
		}
	}
	return err
}

func (b *WebSocketSpotHandle) SetDepthIncrementSnapshotGroupChannel(chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chDeltaMap != nil {
		b.depthIncrementSnapshotDeltaGroupChanMap = make(map[string]chan *client.WsDepthRsp, 100)
		for k, v := range chDeltaMap {
			b.depthIncrementSnapshotDeltaGroupChanMap[GetInstId(k)] = v
		}
	}
	if chFullMap != nil {
		b.depthIncrementSnapshotFullGroupChanMap = make(map[string]chan *depth.Depth, 100)
		for k, v := range chFullMap {
			b.depthIncrementSnapshotFullGroupChanMap[GetInstId(k)] = v
		}
	}
}

func (b *WebSocketSpotHandle) SetDepthIncrementSnapShotConf(symbols []*client.SymbolInfo, conf *base.IncrementDepthConf) {
	if conf.DepthCapLevel <= 0 {
		conf.DepthCapLevel = 1000
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

func (b *WebSocketSpotHandle) DepthIncrementSnapShotGroupHandle(data []byte) error {
	var (
		//res        = &base.OrderBook{}
		//deltaDepth = &base.DeltaDepthUpdate{}
		err                          error
		depth_                       KKDepth
		resp                         []interface{}
		symbol                       string
		bids, asks                   []Level
		bids_, asks_                 []Level // 增量数据
		depth                        = 1000
		t                            = time.Now().UnixMicro()
		exchangeTimeA, exchangeTimeB int64
	)
	if strings.Contains(string(data), "event") && strings.Contains(string(data), "subscription") || strings.Contains(string(data), "connectionID") {
		logger.Logger.Info("subscribe data ", string(data))
		return nil
	}
	// 得到名称
	if strings.Contains(string(data), "heartbeat") {
		return nil
		// not a heartbeat message
	} else {
		err = json.Unmarshal(data, &resp)
		if err != nil {
			logger.Logger.Error("解析错误:", string(data))
			return err

		}
		if len(resp) < 0 {
			return nil
		}
		//第一次全量处理
		if strings.Contains(string(data), "as") || strings.Contains(string(data), "bs") {
			symbol = resp[3].(string)
			symbol = strings.ReplaceAll(symbol, "XBT", "BTC")
			bids, exchangeTimeB = b.CreateInitial(resp[1].(map[string]interface{}), "bs")
			asks, exchangeTimeA = b.CreateInitial(resp[1].(map[string]interface{}), "as")

			// 获取精度
			if tmps, ok := resp[1].(map[string]interface{})["as"].([]interface{}); ok {
				if len(tmps) <= 0 {
					logger.Logger.Error("币对无交易信息:", string(data))
					return nil
				}
				if ttmps, ok := tmps[0].([]interface{}); ok {
					numList := strings.Split(ttmps[0].(string), ".")
					if len(numList) > 0 {
						b.pPmap[symbol] = int32(len(numList[1]))
					}
					numList = strings.Split(ttmps[1].(string), ".")
					if len(numList) > 0 {
						b.aPmap[symbol] = int32(len(numList[1]))
					}
				}
			} else {
				logger.Logger.Error("获取精度失败", string(data))
				panic("获取精度失败")
			}
			b.rMap[symbol] = false
			b.firstFullSentMap[symbol] = true

			if exchangeTimeB > exchangeTimeA {
				exchangeTimeA = exchangeTimeB
			}
			depth_ = KKDepth{
				Symbol:       symbol,
				Asks:         asks,
				Bids:         bids,
				TimeExchange: exchangeTimeA,
				TimeReceive:  t,
			}
			b.DepthCacheMap.Set(symbol, depth_)
			return nil
		}

		if len(resp) == 4 {
			symbol = resp[3].(string)
			symbol = strings.ReplaceAll(symbol, "XBT", "BTC")
			if b.rMap[symbol] {
				logger.Logger.Warn("重连中")
				return nil
			}
			if value, ok := b.DepthCacheMap.Get(symbol); !ok || value == nil {
				logger.Logger.Warn("本地没有对应数据存储", string(data))
				return errors.New("本地没有对应增量数据存储")
			} else {
				if depth_, ok = value.(KKDepth); !ok {
					depth_ = KKDepth{Symbol: symbol}
				}
			}
			bids = depth_.Bids
			asks = depth_.Asks
			var order_book_diff map[string]interface{} = resp[1].(map[string]interface{})
			checksum := order_book_diff["c"].(string)

			if val, ok := order_book_diff["b"]; ok {

				for _, el_interface := range val.([]interface{}) {
					price, volume, exchangeTime, length, priceF, volumeF := getPriceAndVolume(el_interface)
					if exchangeTime > exchangeTimeA {
						exchangeTimeA = exchangeTime
					}
					if volume.Equals(decimal.Zero) {
						bids_ = InsertPriceInBids(bids_, price, volume, priceF, volumeF)
						bids = RemovePriceFromBids(bids, price)
					} else {
						if length == 4 { // it has the 4th element "r" so we just re-append
							// bids = append(bids, Level{Price: price, Volume: volume, PriceF: priceF, VolumeF: volumeF})
						} else {
							bids = InsertPriceInBids(bids, price, volume, priceF, volumeF)
							bids_ = InsertPriceInBids(bids_, price, volume, priceF, volumeF)
							if len(bids) > depth {
								bids = bids[:depth]
							}
						}
					}
				}
			} else {

				for _, el_interface := range order_book_diff["a"].([]interface{}) {
					price, volume, exchangeTime, length, priceF, volumeF := getPriceAndVolume(el_interface)
					if exchangeTime > exchangeTimeA {
						exchangeTimeA = exchangeTime
					}
					if volume.Equals(decimal.Zero) {
						asks_ = InsertPriceInAsks(asks_, price, volume, priceF, volumeF)
						asks = RemovePriceFromAsks(asks, price)
					} else {
						if length == 4 { // it has the 4th element "r" so we just re-append
							// asks = append(asks, Level{Price: price, Volume: volume, PriceF: priceF, VolumeF: volumeF})
						} else {
							asks = InsertPriceInAsks(asks, price, volume, priceF, volumeF)
							asks_ = InsertPriceInAsks(asks_, price, volume, priceF, volumeF)
							if len(asks) > depth {
								asks = asks[:depth]
							}
						}
					}
				}
			}
			if !b.verifyOrderBookChecksum(bids, asks, checksum, symbol) {
				logger.Logger.Error("check错误币对:", string(symbol))
				b.recoon(symbol)
				return nil
			}
		} else {
			symbol = resp[4].(string)
			symbol = strings.ReplaceAll(symbol, "XBT", "BTC")
			if b.rMap[symbol] {
				logger.Logger.Warn("重连中")
				return nil
			}
			if value, ok := b.DepthCacheMap.Get(symbol); !ok || value == nil {
				logger.Logger.Warn("本地没有对应增量数据存储")
				return errors.New("本地没有对应增量数据存储")
			} else {
				depth_ = value.(KKDepth)
			}
			bids = depth_.Bids
			asks = depth_.Asks
			order_book_diff_asks := resp[1].(map[string]interface{})
			order_book_diff_bids := resp[2].(map[string]interface{})
			checksum := order_book_diff_bids["c"].(string)

			for _, el_interface := range order_book_diff_bids["b"].([]interface{}) {
				price, volume, exchangeTime, length, priceF, volumF := getPriceAndVolume(el_interface)
				if exchangeTime > exchangeTimeA {
					exchangeTimeA = exchangeTime
				}
				if volume.Equals(decimal.Zero) {
					bids_ = InsertPriceInBids(bids_, price, volume, priceF, volumF)
					bids = RemovePriceFromBids(bids, price)
				} else {
					if length == 4 { // it has the 4th element "r" so we just re-append
						// bids = append(bids, Level{Price: price, Volume: volume, PriceF: priceF, VolumeF: volumF})
					} else {
						bids = InsertPriceInBids(bids, price, volume, priceF, volumF)
						bids_ = InsertPriceInBids(bids_, price, volume, priceF, volumF)
						if len(bids) > depth {
							bids = bids[:depth]
						}
					}
				}
			}
			for _, el_interface := range order_book_diff_asks["a"].([]interface{}) {
				price, volume, exchangeTime, length, priceF, volumF := getPriceAndVolume(el_interface)
				if exchangeTime > exchangeTimeA {
					exchangeTimeA = exchangeTime
				}
				if volume.Equals(decimal.Zero) {
					asks_ = InsertPriceInAsks(asks_, price, volume, priceF, volumF)

					asks = RemovePriceFromAsks(asks, price)
				} else {
					if length == 4 { // it has the 4th element "r" so we just re-append
						// asks = append(asks, Level{Price: price, Volume: volume, PriceF: priceF, VolumeF: volumF})
					} else {
						asks = InsertPriceInAsks(asks, price, volume, priceF, volumF)
						asks_ = InsertPriceInAsks(asks_, price, volume, priceF, volumF)
						if len(asks) > depth {
							asks = asks[:depth]
						}
					}
				}
			}
			if !b.verifyOrderBookChecksum(bids, asks, checksum, symbol) {
				logger.Logger.Error("check错误币对:", string(data))
				b.recoon(symbol)
				return nil
			}
		}
		depth_.Asks = asks
		depth_.Bids = bids
		depth_.TimeExchange = exchangeTimeA
		depth_.TimeReceive = t

		b.DepthCacheMap.Set(symbol, depth_)
		tmp := &KKDepth{
			Bids:         bids_,
			Asks:         asks_,
			TimeExchange: exchangeTimeA,
			TimeReceive:  t,
		}

		if b.IsPublishDelta {
			if _, ok := b.depthIncrementSnapshotDeltaGroupChanMap[symbol]; ok {
				base.SendChan(b.depthIncrementSnapshotDeltaGroupChanMap[symbol], kk2clientWs(tmp), "deltaDepth")
			} else {
				logger.Logger.Warn("get symbol from channel map err:", symbol)
			}
		}
		//发送全量
		if b.IsPublishFull {
			if _, ok := b.depthIncrementSnapshotFullGroupChanMap[symbol]; ok {
				tmp := kkd2orderbook(depth_)
				if b.firstFullSentMap[symbol] {
					b.firstFullSentMap[symbol] = false
					tmp.Hdr = base.MakeFirstDepthHdr()
				}

				base.SendChan(b.depthIncrementSnapshotFullGroupChanMap[symbol], tmp, "fullDepth")
				//b.depthIncrementSnapshotFullGroupChanMap[sym] <- &depth_
			} else {
				logger.Logger.Warn("get symbol from channel map err:", symbol)
			}
		}
		b.CheckSendStatus.CheckUpdateTimeMap.Set(symbol, time.Now().UnixMicro())

	}
	return err
}

func (b *WebSocketSpotHandle) CreateInitial(book map[string]interface{}, key string) ([]Level, int64) {
	var (
		list         []Level = make([]Level, 0)
		exchangeTime int64
	)
	for _, element := range book[key].([]interface{}) {
		price_interface := element.([]interface{})[0]
		price_str := price_interface.(string)
		price, err := decimal.NewFromString(price_str)
		if err != nil {
			log.Fatal(err)
		}
		priceF := transform.StringToX[float64](price_str).(float64)
		vol_interface := element.([]interface{})[1]
		vol_str := vol_interface.(string)
		vol, err := decimal.NewFromString(vol_str)
		volF := transform.StringToX[float64](vol_str).(float64)
		if err != nil {
			log.Fatal(err)
		}
		tmp := ParseI(element.([]interface{})[2])
		if tmp > exchangeTime {
			exchangeTime = tmp
		}
		list = append(list, Level{Price: price, Volume: vol, PriceF: priceF, VolumeF: volF})
	}
	return list, exchangeTime
}

// 重新订阅
func (d *WebSocketSpotHandle) recoon(symbol string) {
	//d.DepthCacheMap.Set(symbol, nil)
	d.rMap[symbol] = true
	symbol = strings.ReplaceAll(symbol, "BTC", "XBT")
	reqx := req{
		Event: "unsubscribe",
		Pair:  []string{symbol},
		Subscription: subscription{
			Name: "book",
			//Depth: 10,
		},
	}
	base.SendChan(d.reSubscribe, reqx, "reSubscribe")
}

func (b *WebSocketSpotHandle) verifyOrderBookChecksum(bids []Level, asks []Level, checksum string, symbol string) bool {
	checksumInput := GetChecksumInput(bids, asks, b.aPmap[symbol], b.pPmap[symbol])
	crc := crc32.ChecksumIEEE([]byte(checksumInput))
	if fmt.Sprint(crc) != checksum {
		logger.Logger.Error(symbol, "not the same ", " ", crc, " ", checksum)
		return false
	} else {
		return true
	}
}
