package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"github.com/goccy/go-json"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
	"strconv"
	"sync"
	"time"
)

type WebSocketHandleInterface interface {
	TradeGroupHandle([]byte) error
	BookTickerGroupHandle([]byte) error
	DepthLimitGroupHandle([]byte) error
	DepthIncrementGroupHandle([]byte) error
	DepthIncrementSnapShotGroupHandle([]byte) error

	SetTradeGroupChannel(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	SetBookTickerGroupChannel(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	SetDepthLimitGroupChannel(map[*client.SymbolInfo]chan *depth.Depth)
	SetDepthIncrementGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	SetDepthIncrementSnapshotGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp, map[*client.SymbolInfo]chan *depth.Depth)

	SetDepthIncrementSnapShotConf([]*client.SymbolInfo, *base.IncrementDepthConf)

	GetChan(chName string) interface{}
}

type WebSocketSpotHandle struct {
	*base.IncrementDepthConf

	tradeGroupChanMap                       map[string]chan *client.WsTradeRsp
	bookTickerGroupChanMap                  map[string]chan *client.WsBookTickerRsp
	depthLimitGroupChanMap                  map[string]chan *depth.Depth       //全量
	depthIncrementGroupChanMap              map[string]chan *client.WsDepthRsp //增量
	depthIncrementSnapshotDeltaGroupChanMap map[string]chan *client.WsDepthRsp //增量合并数据
	depthIncrementSnapshotFullGroupChanMap  map[string]chan *depth.Depth       //增量合并数据
	depthIncrJudgeMap                       sync.Map

	symbolMap       map[string]*client.SymbolInfo
	Lock            sync.Mutex
	CheckSendStatus *base.CheckDataSendStatus
}

func NewWebSocketSpotHandle(chanCap int64) *WebSocketSpotHandle {
	d := &WebSocketSpotHandle{}
	d.CheckSendStatus = base.NewCheckDataSendStatus()
	return d

}

func (b *WebSocketSpotHandle) GetChan(chName string) interface{} {
	switch chName {
	case "send_status":
		return b.CheckSendStatus
	default:
		return nil
	}
}

func (b *WebSocketSpotHandle) GeneralHandle(name string) func([]byte) error {
	hand := func(data []byte) error {
		t := time.Now().UnixMicro()
		var (
			respFirst Response
			err       error
		)
		if err = json.Unmarshal(data, &respFirst); err != nil {
			// fmt.Println("Err", err)
			logger.Logger.Error("receive data err:", string(data))
			return err
		}
		if respFirst.Type == "heartbeat" {
			return nil
		}
		switch name {
		case "trade":
			return b.TradeGroupHandle(data, t, respFirst.Type)
		case "increment":
			if respFirst.Type == "trade" {
				return nil
			} else if respFirst.Type == "l2_updates" {
				if _, ok := b.depthIncrJudgeMap.Load(respFirst.Symbol); ok {
					return b.DepthIncrementGroupHandle(data, t)
				} else {
					b.depthIncrJudgeMap.Store(respFirst.Symbol, 1)
					return nil
				}
			}
		case "snapshot":
			if respFirst.Type == "trade" {
				return nil
			} else if respFirst.Type == "l2_updates" {
				if _, ok := b.depthIncrJudgeMap.Load(respFirst.Symbol); ok {
					return b.DepthIncrementSnapShotGroupHandle(data, t, false)
				} else {
					b.depthIncrJudgeMap.Store(respFirst.Symbol, 1)
					logger.Logger.Info("Subscribed Success: ", respFirst.Symbol)
					return b.DepthIncrementSnapShotGroupHandle(data, t, true)
				}
			}
		}
		return nil
	}
	return hand
}

func (b *WebSocketSpotHandle) TradeGroupHandle(data []byte, t int64, msgType string) error {
	//var str bytes.Buffer
	//_ = json.Indent(&str, []byte(data), "", "    ")
	//fmt.Println("formated: ", str.String())
	var (
		err           error
		respTrade     TradesItem
		respL2Updates L2Update
		takerSide     order.TradeSide
	)
	if msgType == "l2_updates" {
		if err = json.Unmarshal(data, &respL2Updates); err != nil {
			//fmt.Println("Err", err)
			logger.Logger.Error("receive data err:", string(data))
			return err
		}
		if respL2Updates.Trades != nil {
			for _, tradeItem := range respL2Updates.Trades {
				if tradeItem.Side == "buy" {
					takerSide = 1
				} else {
					takerSide = 2
				}
				res := &client.WsTradeRsp{
					OrderId:      strconv.FormatInt(tradeItem.EventID, 10),
					Symbol:       symbolNameMap[tradeItem.Symbol],
					ExchangeTime: tradeItem.Timestamp * 1000,
					ReceiveTime:  t,
					Price:        StringToFloat(tradeItem.Price),
					Amount:       StringToFloat(tradeItem.Quantity),
					TakerSide:    takerSide,
				}
				if _, ok := b.tradeGroupChanMap[res.Symbol]; ok {
					base.SendChan(b.tradeGroupChanMap[res.Symbol], res, "trade")
				} else {
					logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
				}
				b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
			}
		}
		return nil
	}
	if msgType == "trade" {
		if err = json.Unmarshal(data, &respTrade); err != nil {
			// fmt.Println("Err", err)
			logger.Logger.Error("receive data err:", string(data))
			return err
		}
		if respTrade.Side == "buy" {
			takerSide = 1
		} else {
			takerSide = 2
		}
		res := &client.WsTradeRsp{
			OrderId:      strconv.FormatInt(respTrade.EventID, 10),
			Symbol:       symbolNameMap[respTrade.Symbol],
			ExchangeTime: respTrade.Timestamp * 1000,
			ReceiveTime:  t,
			Price:        StringToFloat(respTrade.Price),
			Amount:       StringToFloat(respTrade.Quantity),
			TakerSide:    takerSide,
		}
		if _, ok := b.tradeGroupChanMap[res.Symbol]; ok {
			base.SendChan(b.tradeGroupChanMap[res.Symbol], res, "trade")
		} else {
			logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
		}
		b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
		return nil
	}
	return nil
}

func (b *WebSocketSpotHandle) BookTickerGroupHandle(data []byte) error {
	return nil
}

func (b *WebSocketSpotHandle) DepthLimitGroupHandle(data []byte) error {
	return nil
}

func (b *WebSocketSpotHandle) DepthIncrementGroupHandle(data []byte, t int64) error {
	var (
		respL2Updates L2Update
		err           error
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
	res := &client.WsDepthRsp{
		ExchangeTime: 0,
		ReceiveTime:  t,
		Symbol:       symbolNameMap[respL2Updates.Symbol],
		Asks:         asks,
		Bids:         bids,
	}
	if _, ok := b.depthIncrementGroupChanMap[res.Symbol]; ok {
		base.SendChan(b.depthIncrementGroupChanMap[res.Symbol], res, "depthIncrement")
	} else {
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) SetTradeGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) {
	if chMap != nil {
		b.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp)
		for info, ch := range chMap {
			b.tradeGroupChanMap[info.Symbol] = ch
		}
	}
}

func (b *WebSocketSpotHandle) SetDepthIncrementGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) {
	if chMap != nil {
		b.depthIncrementGroupChanMap = make(map[string]chan *client.WsDepthRsp)
		for info, ch := range chMap {
			b.depthIncrementGroupChanMap[info.Symbol] = ch
		}
	}
}

func (b *WebSocketSpotHandle) SetDepthIncrementSnapshotGroupChannel(chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chDeltaMap != nil {
		b.depthIncrementSnapshotDeltaGroupChanMap = make(map[string]chan *client.WsDepthRsp)
		for info, ch := range chDeltaMap {
			sym := info.Symbol
			b.depthIncrementSnapshotDeltaGroupChanMap[sym] = ch
		}
	}
	if chFullMap != nil {
		b.depthIncrementSnapshotFullGroupChanMap = make(map[string]chan *depth.Depth)
		for info, ch := range chFullMap {
			sym := info.Symbol
			b.depthIncrementSnapshotFullGroupChanMap[sym] = ch
		}
	}
}

func (b *WebSocketSpotHandle) SetBookTickerGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) {
	if chMap != nil {
		b.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp)
		for info, ch := range chMap {
			b.bookTickerGroupChanMap[info.Symbol] = ch
		}
	}
}

func (b *WebSocketSpotHandle) SetDepthLimitGroupChannel(chMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chMap != nil {
		b.depthLimitGroupChanMap = make(map[string]chan *depth.Depth)
		for info, ch := range chMap {
			b.depthLimitGroupChanMap[info.Symbol] = ch
		}
	}
}

func StringToFloat(valueS string) float64 {
	valueF, err := strconv.ParseFloat(valueS, 64)
	if err != nil {
		logger.Logger.Error("String to Float err:", err)
	}
	return valueF
}
