package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"github.com/goccy/go-json"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
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
		//var str bytes.Buffer
		//_ = json.Indent(&str, []byte(data), "", "    ")
		//fmt.Println("formated: ", str.String())
		t := time.Now().UnixMicro()
		var (
			firstResp FirstResp
			err       error
		)
		if err = json.Unmarshal(data, &firstResp); err != nil {
			//fmt.Println("Err", err)
			logger.Logger.Error("receive data err:", string(data))
			return err
		}
		if firstResp.Msg != "" {
			return nil
		}
		if firstResp.Pong != 0 {
			return nil
		}
		switch name {
		case "trades":
			return b.TradeGroupHandle(data, t)
		case "book":
			return b.BookTickerGroupHandle(data, t)
		case "limit":
			return b.DepthLimitGroupHandle(data, t)
		case "increment":
			return b.DepthIncrementGroupHandle(data, t)
		case "snapshot":
			return b.DepthIncrementSnapShotGroupHandle(data, t)
		}
		return nil
	}
	return hand
}

func (b *WebSocketSpotHandle) TradeGroupHandle(data []byte, t int64) error {
	var (
		err       error
		resp      TradeResponse
		takerSide order.TradeSide
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		//fmt.Println("Err", err)
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if resp.Data.M == true {
		takerSide = 1
	} else {
		takerSide = 2
	}
	res := &client.WsTradeRsp{
		OrderId:      resp.Data.V,
		Symbol:       symbolNameMap[resp.Params.Symbol],
		ExchangeTime: resp.Data.T * 1000,
		ReceiveTime:  t,
		Price:        StringToFloat(resp.Data.P),
		Amount:       StringToFloat(resp.Data.Q),
		TakerSide:    takerSide,
	}
	//fmt.Println(res)
	if _, ok := b.tradeGroupChanMap[res.Symbol]; ok {
		base.SendChan(b.tradeGroupChanMap[res.Symbol], res, "trade")
	} else {
		//fmt.Println("Error1: ", res)
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) BookTickerGroupHandle(data []byte, t int64) error {
	var (
		err  error
		resp TickerResponse
		asks *depth.DepthLevel
		bids *depth.DepthLevel
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		//fmt.Println("Err", err)
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	asks = &depth.DepthLevel{
		Price:  StringToFloat(resp.Data.AskPrice),
		Amount: StringToFloat(resp.Data.AskQty),
	}
	bids = &depth.DepthLevel{
		Price:  StringToFloat(resp.Data.BidPrice),
		Amount: StringToFloat(resp.Data.BidQty),
	}
	res := &client.WsBookTickerRsp{
		ExchangeTime: resp.Data.Time * 1000,
		ReceiveTime:  t,
		Symbol:       symbolNameMap[resp.Params.Symbol],
		Ask:          asks,
		Bid:          bids,
	}
	//fmt.Println("res: ", res)
	if _, ok := b.bookTickerGroupChanMap[res.Symbol]; ok {
		base.SendChan(b.bookTickerGroupChanMap[res.Symbol], res, "bookTickerGroup")
	} else {
		//fmt.Println("Error1: ", res)
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) DepthIncrementGroupHandle(data []byte, t int64) error {
	var (
		resp RespIncrementDepthStream
		err  error
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		//fmt.Println("Err", err)
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	if resp.F == true {
		return nil
	}
	if resp.Data == nil {
		return nil
	}
	asks, err := DepthLevelParse(resp.Data[0].A)
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	bids, err := DepthLevelParse(resp.Data[0].B)
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	res := &client.WsDepthRsp{
		ExchangeTime: resp.Data[0].T * 1000,
		ReceiveTime:  t,
		Symbol:       symbolNameMap[resp.Symbol],
		Asks:         asks,
		Bids:         bids,
	}
	// fmt.Println("res: ", res)
	if _, ok := b.depthIncrementGroupChanMap[res.Symbol]; ok {
		base.SendChan(b.depthIncrementGroupChanMap[res.Symbol], res, "depthIncrementGroup")
	} else {
		// fmt.Println("Error1: ", res)
		logger.Logger.Warn("get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(res.Symbol, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) DepthLimitGroupHandle(data []byte, t int64) error {
	var (
		err  error
		resp RespLimitDepthStream
	)
	if err = json.Unmarshal(data, &resp); err != nil {
		//fmt.Println("Err", err)
		logger.Logger.Error("receive data err:", string(data))
		return err
	}
	asks, err := DepthLevelParse(resp.Data.A)
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	bids, err := DepthLevelParse(resp.Data.B)
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	res := &depth.Depth{
		Exchange:     common.Exchange_BYBIT,
		Symbol:       symbolNameMap[resp.Params.Symbol],
		TimeExchange: uint64(resp.Data.T * 1000),
		TimeReceive:  uint64(t),
		TimeOperate:  uint64(time.Now().UnixMicro()),
		Asks:         asks,
		Bids:         bids,
	}
	//fmt.Println("res: ", res)
	if _, ok := b.depthLimitGroupChanMap[res.Symbol]; ok {
		base.SendChan(b.depthLimitGroupChanMap[res.Symbol], res, "depthLimitGroup")
	} else {
		// fmt.Println("Error1: ", res)
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

//func GetChecksum(orderbook *base.OrderBook) uint32 {
//	var fields []string
//
//	la, lb := len(orderbook.Asks), len(orderbook.Bids)
//	for i := 0; i < 100; i++ {
//		if i < lb {
//			fields = append(fields, AppendDecimal(orderbook.Bids[i].Price)+":"+AppendDecimal(orderbook.Bids[i].Amount))
//		}
//		if i < la {
//			fields = append(fields, AppendDecimal(orderbook.Asks[i].Price)+":"+AppendDecimal(orderbook.Asks[i].Amount))
//		}
//	}
//
//	raw := strings.Join(fields, ":")
//	cs := crc32.ChecksumIEEE([]byte(raw))
//	return cs
//}

//func AppendDecimal(value float64) string {
//	//Case 1: Small number
//	if value < 0.0001 { //Really small values need to reformat scientific notation
//		return strconv.FormatFloat(value, 'e', -1, 64)
//	}
//
//	//Case 2: Not a small number
//	res := fmt.Sprintf("%v", value)
//	if strings.Contains(res, "e") { //Unpack large values formatted in scientific notation
//		res = strconv.FormatFloat(value, 'f', -1, 64)
//	}
//	//Return ok formatted number
//	if strings.Contains(res, ".") {
//		return res
//	} else { //Append .0 if not formatted correctly
//		res += ".0"
//		return res
//	}
//}

func StringToFloat(valueS string) float64 {
	valueF, err := strconv.ParseFloat(valueS, 64)
	if err != nil {
		logger.Logger.Error("String to Float err:", err)
	}
	return valueF
}

func DepthLevelParse(levelList [][]string) ([]*depth.DepthLevel, error) {
	var (
		res []*depth.DepthLevel
		err error
	)
	for _, level := range levelList {
		res = append(res, &depth.DepthLevel{
			Price:  StringToFloat(level[0]),
			Amount: StringToFloat(level[1]),
		})
	}
	return res, err
}
