package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/binance/c_api"
	"clients/exchange/cex/binance/spot_api"
	"clients/exchange/cex/binance/u_api"
	"clients/logger"
	"clients/transform"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
)

type WebSocketPublicHandleInterface interface {
	FundingRateGroupHandle([]byte) error
	AggTradeGroupHandle([]byte) error
	TradeGroupHandle([]byte) error
	BookTickerGroupHandle([]byte) error
	DepthIncrementGroupHandle([]byte) error
	DepthLimitGroupHandle([]byte) error
	DepthIncrementSnapShotGroupHandle([]byte) error

	//设置chan map
	SetFundingRateGroupChannel(map[*client.SymbolInfo]chan *client.WsFundingRateRsp)
	SetTradeGroupChannel(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	SetBookTickerGroupChannel(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	SetDepthLimitGroupChannel(map[*client.SymbolInfo]chan *depth.Depth)
	SetDepthIncrementGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	SetDepthIncrementSnapshotGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp, map[*client.SymbolInfo]chan *depth.Depth)
	SetAggTradeGroupChannel(map[*client.SymbolInfo]chan *client.WsAggTradeRsp) // binance

	SetDepthIncrementSnapShotConf([]*client.SymbolInfo, *base.IncrementDepthConf)
}

type WebSocketHandleInterface interface {
	WebSocketPublicHandleInterface
	WebSocketPrivateHandleInterface
}
type WebSocketPrivateHandleInterface interface {
	AccountHandle([]byte) error
	BalanceHandle([]byte) error
	MarginAccountHandle([]byte) error
	MarginBalanceHandle([]byte) error
	OrderHandle([]byte) error
	GetChan(chName string) interface{}
}

type WebSocketSpotHandle struct {
	*base.IncrementDepthConf

	accountChan chan *client.WsAccountRsp
	balanceChan chan *client.WsBalanceRsp
	orderChan   chan *client.WsOrderRsp

	fundingRateGroupChanMap                 map[string]chan *client.WsFundingRateRsp
	aggTradeGroupChanMap                    map[string]chan *client.WsAggTradeRsp
	tradeGroupChanMap                       map[string]chan *client.WsTradeRsp
	bookTickerGroupChanMap                  map[string]chan *client.WsBookTickerRsp
	depthLimitGroupChanMap                  map[string]chan *depth.Depth       //全量
	depthIncrementGroupChanMap              map[string]chan *client.WsDepthRsp //增量
	depthIncrementSnapshotDeltaGroupChanMap map[string]chan *client.WsDepthRsp //增量数据
	depthIncrementSnapshotFullGroupChanMap  map[string]chan *depth.Depth       //合并的全量数据

	symbolMap       map[string]*client.SymbolInfo
	Lock            sync.Mutex // 锁
	CheckSendStatus *base.CheckDataSendStatus
	MarketType      int
}

func NewWebSocketSpotHandle(chanCap int64) *WebSocketSpotHandle {
	if chanCap < 1 {
		chanCap = 1024
	}
	d := &WebSocketSpotHandle{}
	d.accountChan = make(chan *client.WsAccountRsp, chanCap)
	d.balanceChan = make(chan *client.WsBalanceRsp, chanCap)
	d.orderChan = make(chan *client.WsOrderRsp, chanCap)
	d.CheckSendStatus = base.NewCheckDataSendStatus()

	return d
}

func (b *WebSocketSpotHandle) GetChan(chName string) interface{} {
	switch chName {
	case "account":
		return b.accountChan
	case "balance":
		return b.balanceChan
	case "order":
		return b.orderChan
	case "send_status":
		return b.CheckSendStatus
	default:
		return nil
	}
}

func (b *WebSocketSpotHandle) FundingRateGroupHandle(data []byte) error {
	var (
		resp        RespMarkPriceStream
		err         error
		fundingRate float64
	)
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	fundingRate, err = strconv.ParseFloat(resp.Data.R, 64)
	symbol, _, type_ := u_api.GetContractType(resp.Data.S)
	_, sym := GetSymbolKey(resp.Data.S, b.MarketType)
	res := &client.WsFundingRateRsp{
		Symbol:       symbol,
		Type:         type_,
		FundingRate:  fundingRate,
		ExchangeTime: resp.Data.E1 * 1000,
		ReceiveTime:  time.Now().UnixMicro(),
	}
	if _, ok := b.fundingRateGroupChanMap[sym]; ok {
		base.SendChan(b.fundingRateGroupChanMap[sym], res, "fundingRateGroup")
	} else {
		logger.Logger.Warn("FundingRateGroupHandle get symbol from channel map err:", sym)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(sym, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) AggTradeGroupHandle(data []byte) error {
	var (
		resp          RespAggTradeStream
		err           error
		price, amount float64
	)
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	price, err = strconv.ParseFloat(resp.Data.P, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return err
	}
	amount, err = strconv.ParseFloat(resp.Data.Q, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return err
	}
	symbol, sym := GetSymbolKey(resp.Data.S, b.MarketType)
	res := &client.WsAggTradeRsp{
		AggId:        strconv.FormatInt(resp.Data.A, 10),
		Symbol:       symbol,
		ExchangeTime: resp.Data.E1 * 1000,
		Price:        price,
		Amount:       amount,
		TakerSide:    GetSide(resp.Data.M),
		DealTime:     resp.Data.T * 1000,
		ReceiveTime:  time.Now().UnixMicro(),
	}
	// 币本位需要转换amount
	if b.MarketType == 2 {
		res.Amount = getCbaseQty(res.Symbol, res.Amount, res.Price)
	}
	if _, ok := b.aggTradeGroupChanMap[sym]; ok {
		base.SendChan(b.aggTradeGroupChanMap[sym], res, "aggTradeGroup")
	} else {
		logger.Logger.Warn("AggTradeGroupHandle get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(sym, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) TradeGroupHandle(data []byte) error {
	res, symbolStr, err := b.parseTrade(data)
	if err != nil {
		return err
	}
	_, sym := GetSymbolKey(symbolStr, b.MarketType)
	// 币本位需要转换amount
	if b.MarketType == 2 {
		res.Amount = getCbaseQty(res.Symbol, res.Amount, res.Price)
	}
	if _, ok := b.tradeGroupChanMap[sym]; ok {
		base.SendChan(b.tradeGroupChanMap[sym], res, fmt.Sprintf("%s %s", "trade", res.Symbol))
	} else {
		logger.Logger.Warn("TradeGroupHandle get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(sym, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) parseTrade(data []byte) (*client.WsTradeRsp, string, error) {
	var (
		sym           string
		resp          RespTradeStream
		err           error
		price, amount float64
	)
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return nil, sym, err
	}
	sym = resp.Data.S
	price, err = strconv.ParseFloat(resp.Data.P, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return nil, sym, err
	}
	amount, err = strconv.ParseFloat(resp.Data.Q, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return nil, sym, err
	}
	symbol, _, _ := u_api.GetContractType(resp.Data.S)
	res := &client.WsTradeRsp{
		OrderId:       strconv.FormatInt(resp.Data.T, 10),
		Symbol:        symbol,
		ExchangeTime:  resp.Data.E1 * 1000,
		Price:         price,
		Amount:        amount,
		TakerSide:     GetSide(resp.Data.M),
		BuyerOrderId:  strconv.FormatInt(resp.Data.B, 10),
		SellerOrderId: strconv.FormatInt(resp.Data.A, 10),
		DealTime:      resp.Data.T1 * 1000,
		ReceiveTime:   time.Now().UnixMicro(),
	}
	return res, sym, nil
}

func (b *WebSocketSpotHandle) HandleRespErr(data []byte, resp interface{}) error {
	var (
		err error
	)
	if strings.Contains(string(data), "{\"code\":") {
		logger.Logger.Debug("binance spot ws data error: ", string(data))
		return errors.New(string(data))
	}
	if strings.Contains(string(data), "\"result\":null,\"id\"") {
		return errors.New("response")
	}
	if strings.Contains(string(data), "\"id\":1,\"result\":null") {
		return errors.New("response")
	}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return errors.New("json unmarshal data err:" + string(data))
	}
	logger.Logger.Debug("binance spot ws data: ", resp)
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

func (b *WebSocketSpotHandle) DepthIncrementGroupHandle(data []byte) error {
	var (
		resp RespIncrementDepthStream
		err  error
	)
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	asks, err := b.DepthLevelParse(resp.Data.A)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	bids, err := b.DepthLevelParse(resp.Data.B)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	symbol, sym := GetSymbolKey(resp.Data.S, b.MarketType)
	res := &client.WsDepthRsp{
		UpdateIdStart: resp.Data.U,
		UpdateIdEnd:   resp.Data.U1,
		ExchangeTime:  resp.Data.E1,
		ReceiveTime:   time.Now().UnixMicro(),
		Symbol:        symbol,
		Asks:          asks,
		Bids:          bids,
	}
	// 币本位需要转换amount
	if b.MarketType == 2 {
		for bidIdx, _ := range res.Bids {
			res.Bids[bidIdx].Amount = getCbaseQty(res.Symbol, res.Bids[bidIdx].Amount, res.Bids[bidIdx].Price)
		}
		for askIdx, _ := range res.Asks {
			res.Asks[askIdx].Amount = getCbaseQty(res.Symbol, res.Asks[askIdx].Amount, res.Asks[askIdx].Price)
		}
	}
	if _, ok := b.depthIncrementGroupChanMap[sym]; ok {
		base.SendChan(b.depthIncrementGroupChanMap[sym], res, fmt.Sprintf("%s %s", "depthIncrement", res.Symbol))
	} else {
		logger.Logger.Warn("DepthIncrementGroupHandle get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(sym, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) DepthLimitGroupHandle(data []byte) error {
	var (
		resp RespLimitDepthStream
		err  error
	)
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	asks, err := b.DepthLevelParse(resp.Data.Asks)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	bids, err := b.DepthLevelParse(resp.Data.Bids)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	symbol, sym := GetSymbolKey(strings.ToUpper(strings.Split(resp.Stream, "@")[0]), b.MarketType)
	res := &depth.Depth{
		Exchange:     common.Exchange_BINANCE,
		Symbol:       symbol,
		TimeExchange: 0,
		TimeReceive:  uint64(time.Now().UnixMicro()),
		TimeOperate:  uint64(time.Now().UnixMicro()),
		Asks:         asks,
		Bids:         bids,
	}
	// 币本位需要转换amount
	if b.MarketType == 2 {
		for bidIdx, _ := range res.Bids {
			res.Bids[bidIdx].Amount = getCbaseQty(res.Symbol, res.Bids[bidIdx].Amount, res.Bids[bidIdx].Price)
		}
		for askIdx, _ := range res.Asks {
			res.Asks[askIdx].Amount = getCbaseQty(res.Symbol, res.Asks[askIdx].Amount, res.Asks[askIdx].Price)
		}
	}
	if _, ok := b.depthLimitGroupChanMap[sym]; ok {
		base.SendChan(b.depthLimitGroupChanMap[sym], res, fmt.Sprintf("%s %s", "depthLimitGroup", sym))
	} else {
		logger.Logger.Warn("DepthLimitGroupHandle get symbol from channel map err:", sym)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(sym, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) BookTickerGroupHandle(data []byte) error {
	var (
		resp RespBookTickerStream
		err  error
	)
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	asks, err := b.DepthLevelParse([][]string{{resp.Data.A, resp.Data.A1}})
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	bids, err := b.DepthLevelParse([][]string{{resp.Data.B, resp.Data.B1}})
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	symbol, sym := GetSymbolKey(resp.Data.S, b.MarketType)
	res := &client.WsBookTickerRsp{
		UpdateIdEnd:  resp.Data.U,
		Symbol:       symbol,
		ExchangeTime: 0,
		ReceiveTime:  time.Now().UnixMicro(),
		Ask:          asks[0],
		Bid:          bids[0],
	}
	if _, ok := b.bookTickerGroupChanMap[sym]; ok {
		base.SendChan(b.bookTickerGroupChanMap[sym], res, fmt.Sprintf("%s %s", "bookTickerGroup", res.Symbol))
	} else {
		logger.Logger.Warn("BookTickerGroupHandle get symbol from channel map err:", res.Symbol)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(sym, time.Now().UnixMicro())
	return nil
}

func (b *WebSocketSpotHandle) SetFundingRateGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsFundingRateRsp) {
	if chMap != nil {
		b.fundingRateGroupChanMap = make(map[string]chan *client.WsFundingRateRsp)
		for info, ch := range chMap {
			sym := SymbolKeyGen(info.Symbol, info.Market, info.Type)
			b.fundingRateGroupChanMap[sym] = ch
		}
	}
}

func (b *WebSocketSpotHandle) SetAggTradeGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsAggTradeRsp) {
	if chMap != nil {
		b.aggTradeGroupChanMap = make(map[string]chan *client.WsAggTradeRsp)
		for info, ch := range chMap {
			sym := SymbolKeyGen(info.Symbol, info.Market, info.Type)
			b.aggTradeGroupChanMap[sym] = ch
		}
	}
}

func (b *WebSocketSpotHandle) SetTradeGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsTradeRsp) {
	if chMap != nil {
		b.tradeGroupChanMap = make(map[string]chan *client.WsTradeRsp)
		for info, ch := range chMap {
			sym := SymbolKeyGen(info.Symbol, info.Market, info.Type)
			b.tradeGroupChanMap[sym] = ch
		}
	}
}

func (b *WebSocketSpotHandle) SetDepthLimitGroupChannel(chMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chMap != nil {
		b.depthLimitGroupChanMap = make(map[string]chan *depth.Depth)
		for info, ch := range chMap {
			sym := SymbolKeyGen(info.Symbol, info.Market, info.Type)
			b.depthLimitGroupChanMap[sym] = ch
		}
	}
}

func (b *WebSocketSpotHandle) SetDepthIncrementGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsDepthRsp) {
	if chMap != nil {
		b.depthIncrementGroupChanMap = make(map[string]chan *client.WsDepthRsp)
		for info, ch := range chMap {
			sym := SymbolKeyGen(info.Symbol, info.Market, info.Type)
			b.depthIncrementGroupChanMap[sym] = ch
		}
	}
}

func (b *WebSocketSpotHandle) SetDepthIncrementSnapshotGroupChannel(chDeltaMap map[*client.SymbolInfo]chan *client.WsDepthRsp, chFullMap map[*client.SymbolInfo]chan *depth.Depth) {
	if chDeltaMap != nil {
		b.depthIncrementSnapshotDeltaGroupChanMap = make(map[string]chan *client.WsDepthRsp)
		for info, ch := range chDeltaMap {
			sym := SymbolKeyGen(info.Symbol, info.Market, info.Type)
			b.depthIncrementSnapshotDeltaGroupChanMap[sym] = ch
		}
	}
	if chFullMap != nil {
		b.depthIncrementSnapshotFullGroupChanMap = make(map[string]chan *depth.Depth)
		for info, ch := range chFullMap {
			sym := SymbolKeyGen(info.Symbol, info.Market, info.Type)
			b.depthIncrementSnapshotFullGroupChanMap[sym] = ch
		}
	}
}

func (b *WebSocketSpotHandle) SetBookTickerGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) {
	if chMap != nil {
		b.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp)
		for info, ch := range chMap {
			sym := SymbolKeyGen(info.Symbol, info.Market, info.Type)
			b.bookTickerGroupChanMap[sym] = ch
		}
	}
}

func (b *WebSocketSpotHandle) AccountHandle(data []byte) error {
	var (
		resp            RespUserAccount
		err             error
		available, lock float64
	)
	//过滤账户更新信息数据
	//fmt.Println(string(data))
	if !strings.Contains(string(data), "outboundAccountPosition") {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	res := &client.WsAccountRsp{
		ExchangeTime: resp.E1,
		UpdateTime:   resp.U,
	}
	for _, item := range resp.B {
		available, err = strconv.ParseFloat(item.F, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		lock, err = strconv.ParseFloat(item.L, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		res.BalanceList = append(res.BalanceList, &client.BalanceItem{
			Assert:    item.A,
			Available: available,
			Lock:      lock,
			Market:    common.Market_SPOT,
		})
	}
	base.SendChan(b.accountChan, res, "account")
	return nil
}

func (b *WebSocketSpotHandle) BalanceHandle(data []byte) error {
	var (
		resp  RespUserBalance
		err   error
		delta float64
	)
	//fmt.Println(string(data))
	if !strings.Contains(string(data), "balanceUpdate") {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	delta, err = strconv.ParseFloat(resp.D, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return err
	}
	res := &client.WsBalanceRsp{
		ExchangeTime: resp.E1,
		UpdateTime:   resp.T,
		Assert:       resp.A,
		Market:       common.Market_SPOT,
		Delta:        delta,
	}
	base.SendChan(b.balanceChan, res, "balance")
	return nil
}

func (b *WebSocketSpotHandle) MarginAccountHandle(data []byte) error {
	var (
		resp            RespUserAccount
		err             error
		available, lock float64
	)
	//过滤账户更新信息数据
	//fmt.Println(string(data))
	if !strings.Contains(string(data), "outboundAccountPosition") {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	res := &client.WsAccountRsp{
		ExchangeTime: resp.E1,
		UpdateTime:   resp.U,
	}
	for _, item := range resp.B {
		available, err = strconv.ParseFloat(item.F, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		lock, err = strconv.ParseFloat(item.L, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		res.BalanceList = append(res.BalanceList, &client.BalanceItem{
			Assert:    item.A,
			Available: available,
			Lock:      lock,
			Market:    common.Market_MARGIN,
		})
	}
	base.SendChan(b.accountChan, res, "account")
	return nil
}

func (b *WebSocketSpotHandle) MarginBalanceHandle(data []byte) error {
	var (
		resp  RespUserBalance
		err   error
		delta float64
	)
	//fmt.Println(string(data))
	if !strings.Contains(string(data), "balanceUpdate") {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	delta, err = strconv.ParseFloat(resp.D, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return err
	}
	res := &client.WsBalanceRsp{
		ExchangeTime: resp.E1,
		UpdateTime:   resp.T,
		Assert:       resp.A,
		Delta:        delta,
		Market:       common.Market_MARGIN,
	}
	base.SendChan(b.balanceChan, res, "balance")
	return nil
}

func (b *WebSocketSpotHandle) OrderHandle(data []byte) error {
	var (
		resp                                      RespUserOrder
		err                                       error
		clientId                                  string
		feeAsset                                  string
		amountFilled, priceFilled, qtyFilled, fee float64
		ok                                        bool
	)
	//fmt.Println(string(data))
	if !strings.Contains(string(data), "executionReport") {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	// 先看撤销的原始订单Id
	clientId = resp.C1
	if clientId == "" {
		// 撤销的原始订单id为空，再看clientOrderId
		clientId = resp.C
	}
	amountFilled, err = strconv.ParseFloat(resp.Z, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return err
	}
	qtyFilled, err = strconv.ParseFloat(resp.Z1, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return err
	}
	if amountFilled != 0 {
		priceFilled = qtyFilled / amountFilled
	}
	if resp.N1 != nil {
		feeAsset, ok = resp.N1.(string)
		if !ok {
			feeAsset = ""
		}
		fee, err = strconv.ParseFloat(resp.N, 64)
		if err != nil {
			logger.Logger.Error("parse fee error:" + err.Error())
			fee = 0
		}
	}

	producer, id := transform.ClientIdToId(clientId)
	var closeDate string
	names := strings.Split(resp.S, "_")
	if len(names) > 1 && len(names[1]) == 6 {
		closeDate = names[1]
	}
	res := &client.WsOrderRsp{
		Producer:     producer,
		Id:           id,
		IdEx:         strconv.FormatInt(resp.I, 10),
		Status:       spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(resp.X1)),
		Symbol:       resp.S,
		AmountFilled: amountFilled,
		PriceFilled:  priceFilled,
		QtyFilled:    qtyFilled,
		Fee:          fee,
		FeeAsset:     feeAsset,
		TimeFilled:   resp.T,
		CloseData:    closeDate,
	}
	base.SendChan(b.orderChan, res, "order")
	return nil
}

type WebSocketUBaseHandle struct {
	WebSocketSpotHandle
}

func NewWebSocketUBaseHandle(chanCap int64) *WebSocketUBaseHandle {
	d := &WebSocketUBaseHandle{}

	d.accountChan = make(chan *client.WsAccountRsp, chanCap)
	d.balanceChan = make(chan *client.WsBalanceRsp, chanCap)
	d.orderChan = make(chan *client.WsOrderRsp, chanCap)
	d.CheckSendStatus = base.NewCheckDataSendStatus()
	d.MarketType = 1
	return d
}

func (b *WebSocketUBaseHandle) GetChan(chName string) interface{} {
	switch chName {
	case "account":
		return b.accountChan
	case "balance":
		return b.balanceChan
	case "order":
		return b.orderChan
	case "send_status":
		return b.CheckSendStatus
	default:
		return nil
	}
}

func (b *WebSocketUBaseHandle) HandleRespErr(data []byte, resp interface{}) error {
	var (
		err error
	)
	if strings.Contains(string(data), "{\"code\":") {
		logger.Logger.Debug("binance spot ws data error: ", string(data))
		return errors.New(string(data))
	}
	if strings.Contains(string(data), "\"result\":null,\"id\":") {
		return errors.New("response")
	}
	if strings.Contains(string(data), "\"id\":1,\"result\":null") {
		return errors.New("response")
	}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return errors.New("json unmarshal data err:" + string(data))
	}
	logger.Logger.Debug("binance spot ws data: ", resp)
	return nil
}

func (b *WebSocketUBaseHandle) DepthLevelParse(levelList [][]string) ([]*depth.DepthLevel, error) {
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

func (b *WebSocketUBaseHandle) BookTickerGroupHandle(data []byte) error {
	var (
		resp RespBookTickerStream
		err  error
	)
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	asks, err := b.DepthLevelParse([][]string{{resp.Data.A, resp.Data.A1}})
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	bids, err := b.DepthLevelParse([][]string{{resp.Data.B, resp.Data.B1}})
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	symbol, sym := GetSymbolKey(resp.Data.S, b.MarketType)
	res := &client.WsBookTickerRsp{
		UpdateIdEnd:  resp.Data.U,
		ExchangeTime: 0,
		ReceiveTime:  time.Now().UnixMicro(),
		Symbol:       symbol,
		Ask:          asks[0],
		Bid:          bids[0],
	}

	if _, ok := b.bookTickerGroupChanMap[sym]; ok {
		base.SendChan(b.bookTickerGroupChanMap[sym], res, "bookTickerGroup")
	} else {
		logger.Logger.Warn("BookTickerGroupHandle get symbol from channel map err:", res.Symbol)
	}
	return nil
}

func (b *WebSocketUBaseHandle) AccountHandle(data []byte) error {
	var (
		resp               RespUBaseAccount
		err                error
		balance, available float64
	)
	//过滤账户更新信息数据
	if !strings.Contains(string(data), "ACCOUNT_UPDATE") {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	res := &client.WsAccountRsp{
		ExchangeTime: resp.E1,
		UpdateTime:   resp.T,
	}
	for _, item := range resp.A.B {
		available, err = strconv.ParseFloat(item.Cw, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		balance, err = strconv.ParseFloat(item.Wb, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		res.BalanceList = append(res.BalanceList, &client.BalanceItem{
			Assert:       item.A,
			Available:    available,
			Lock:         balance - available,
			ExchangeTime: resp.E1,
			Market:       common.Market_FUTURE,
			ChangeType:   spot_api.GetChangeType(resp.A.M),
		})
	}
	for _, item := range resp.A.P {
		symbol, market, symbolType := u_api.GetContractType(item.S)
		position, err := strconv.ParseFloat(item.Pa, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		positionPrice, err := strconv.ParseFloat(item.Ep, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		unprofit, err := strconv.ParseFloat(item.Up, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		tradeSide := order.TradeSide_BUY
		if position < 0 {
			tradeSide = order.TradeSide_SELL
		}
		res.PositionList = append(res.PositionList, &client.PositionItem{
			ExchangeTime: resp.E1,
			Symbol:       symbol,
			Market:       market,
			Type:         symbolType,
			Side:         tradeSide,
			Unprofit:     unprofit,
			Position:     math.Abs(position),
			Price:        positionPrice,
			ChangeType:   spot_api.GetChangeType(resp.A.M),
		})
	}
	base.SendChan(b.accountChan, res, "account")
	return nil
}

func GetCloseTime(symbol string) string {
	var closeDate string
	names := strings.Split(symbol, "_")
	if len(names) > 1 && len(names[1]) == 6 {
		closeDate = names[1]
	}
	return closeDate
}

func (b *WebSocketUBaseHandle) OrderHandle(data []byte) error {
	var (
		resp                                      RespUBaseOrder
		err                                       error
		clientId                                  string
		feeAsset                                  string
		amountFilled, priceFilled, qtyFilled, fee float64
	)
	if !strings.Contains(string(data), "ORDER_TRADE_UPDATE") {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	// 先看撤销的原始订单Id
	clientId = strconv.FormatInt(resp.O.I, 10)
	if clientId == "" {
		// 撤销的原始订单id为空，再看clientOrderId
		clientId = resp.O.C
	}
	symbol, market, symbolType := u_api.GetContractType(resp.O.S)
	amountFilled, err = strconv.ParseFloat(resp.O.Z, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return err
	}
	feeAsset = resp.O.N1
	fee, err = strconv.ParseFloat(resp.O.N, 64)
	if err != nil {
		logger.Logger.Error("parse fee error:" + err.Error())
		fee = 0
	}

	producer, id := transform.ClientIdToId(clientId)
	res := &client.WsOrderRsp{
		Producer:     producer,
		Id:           id,
		IdEx:         strconv.FormatInt(resp.O.I, 10),
		Status:       spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(resp.O.X1)),
		TimeFilled:   resp.T,
		Symbol:       symbol,
		Market:       market,
		Type:         symbolType,
		AmountFilled: amountFilled,
		PriceFilled:  priceFilled,
		QtyFilled:    qtyFilled,
		Fee:          fee,
		FeeAsset:     feeAsset,
		CloseData:    GetCloseTime(resp.O.S),
	}

	if len(b.orderChan) > cap(b.orderChan)/2 {
		content := fmt.Sprint("order channel slow:", len(b.orderChan))
		logger.Logger.Warn("OrderHandle", content)
	}
	base.SendChan(b.orderChan, res, "order")
	return nil
}

func (b *WebSocketUBaseHandle) SetBookTickerGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) {
	if chMap != nil {
		b.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp)
		for info, ch := range chMap {
			sym := SymbolKeyGen(info.Symbol, info.Market, info.Type)
			b.bookTickerGroupChanMap[sym] = ch
		}
	}
}

func (b *WebSocketUBaseHandle) DepthLimitGroupHandle(data []byte) error {
	var (
		resp RespUBaseLimitDepthStream
		err  error
	)
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	asks, err := b.DepthLevelParse(resp.Data.A)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	bids, err := b.DepthLevelParse(resp.Data.B)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	symbol, sym := GetSymbolKey(strings.ToUpper(strings.Split(resp.Stream, "@")[0]), b.MarketType)
	res := &depth.Depth{
		Exchange:     common.Exchange_BINANCE,
		Symbol:       symbol,
		TimeExchange: uint64(resp.Data.E1) * 1000,
		TimeReceive:  uint64(time.Now().UnixMicro()),
		TimeOperate:  uint64(time.Now().UnixMicro()),
		Asks:         asks,
		Bids:         bids,
	}
	// 币本位需要转换amount
	if b.MarketType == 2 {
		for bidIdx, _ := range res.Bids {
			res.Bids[bidIdx].Amount = getCbaseQty(res.Symbol, res.Bids[bidIdx].Amount, res.Bids[bidIdx].Price)
		}
		for askIdx, _ := range res.Asks {
			res.Asks[askIdx].Amount = getCbaseQty(res.Symbol, res.Asks[askIdx].Amount, res.Asks[askIdx].Price)
		}
	}
	if _, ok := b.depthLimitGroupChanMap[sym]; ok {
		base.SendChan(b.depthLimitGroupChanMap[sym], res, fmt.Sprintf("%s %s", "depthLimitGroup", sym))
	} else {
		logger.Logger.Warn("DepthLimitGroupHandle get symbol from channel map err:", sym)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(sym, time.Now().UnixMicro())
	return nil
}

type WebSocketCBaseHandle struct {
	WebSocketSpotHandle
}

func NewWebSocketCBaseHandle(chanCap int64) *WebSocketCBaseHandle {
	d := &WebSocketCBaseHandle{}

	d.accountChan = make(chan *client.WsAccountRsp, chanCap)
	d.balanceChan = make(chan *client.WsBalanceRsp, chanCap)
	d.orderChan = make(chan *client.WsOrderRsp, chanCap)
	d.CheckSendStatus = base.NewCheckDataSendStatus()
	d.MarketType = 2
	return d
}

func (b *WebSocketCBaseHandle) GetChan(chName string) interface{} {
	switch chName {
	case "account":
		return b.accountChan
	case "balance":
		return b.balanceChan
	case "order":
		return b.orderChan
	case "send_status":
		return b.CheckSendStatus

	default:
		return nil
	}
}

func (b *WebSocketCBaseHandle) HandleRespErr(data []byte, resp interface{}) error {
	var (
		err error
	)
	if strings.Contains(string(data), "{\"code\":") {
		logger.Logger.Debug("binance spot ws data error: ", string(data))
		return errors.New(string(data))
	}
	if strings.Contains(string(data), "\"result\":null,\"id\":") {
		return errors.New("response")
	}
	if strings.Contains(string(data), "\"id\":1,\"result\":null") {
		return errors.New("response")
	}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return errors.New("json unmarshal data err:" + string(data))
	}
	logger.Logger.Debug("binance spot ws data: ", resp)
	return nil
}

func (b *WebSocketCBaseHandle) DepthLevelParse(levelList [][]string) ([]*depth.DepthLevel, error) {
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

func (b *WebSocketCBaseHandle) BookTickerGroupHandle(data []byte) error {
	var (
		resp RespBookTickerStream
		err  error
	)
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	asks, err := b.DepthLevelParse([][]string{{resp.Data.A, resp.Data.A1}})
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	bids, err := b.DepthLevelParse([][]string{{resp.Data.B, resp.Data.B1}})
	if err != nil {
		logger.Logger.Error("book ticker parse err", err, string(data))
		return err
	}
	symbol, sym := GetSymbolKey(resp.Data.S, b.MarketType)
	res := &client.WsBookTickerRsp{
		UpdateIdEnd:  resp.Data.U,
		ExchangeTime: 0,
		ReceiveTime:  time.Now().UnixMicro(),
		Symbol:       symbol,
		Ask:          asks[0],
		Bid:          bids[0],
	}
	res.Ask.Amount = getCbaseQty(symbol, res.Ask.Amount, res.Ask.Price)
	res.Bid.Amount = getCbaseQty(symbol, res.Bid.Amount, res.Bid.Price)
	if _, ok := b.bookTickerGroupChanMap[sym]; ok {
		base.SendChan(b.bookTickerGroupChanMap[sym], res, fmt.Sprintf("%s %s", "bookTickerGroup", sym))
	} else {
		logger.Logger.Warn("BookTickerGroupHandle get symbol from channel map err:", sym)
	}
	return nil
}

func (b *WebSocketCBaseHandle) AccountHandle(data []byte) error {
	var (
		resp               RespUBaseAccount
		err                error
		balance, available float64
	)
	//过滤账户更新信息数据
	if !strings.Contains(string(data), "ACCOUNT_UPDATE") {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	res := &client.WsAccountRsp{
		ExchangeTime: resp.E1,
		UpdateTime:   resp.T,
	}
	for _, item := range resp.A.B {
		available, err = strconv.ParseFloat(item.Cw, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		balance, err = strconv.ParseFloat(item.Wb, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		res.BalanceList = append(res.BalanceList, &client.BalanceItem{
			Assert:       item.A,
			Available:    available,
			Lock:         balance - available,
			ExchangeTime: resp.E1,
			Market:       common.Market_FUTURE_COIN,
			ChangeType:   spot_api.GetChangeType(resp.A.M),
		})
	}
	for _, item := range resp.A.P {
		symbol, market, symbolType := c_api.GetContractType(item.S)
		position, err := strconv.ParseFloat(item.Pa, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		positionPrice, err := strconv.ParseFloat(item.Ep, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		unprofit, err := strconv.ParseFloat(item.Up, 64)
		if err != nil {
			logger.Logger.Error("trade parse err", err, string(data))
			return err
		}
		tradeSide := order.TradeSide_BUY
		if position < 0 {
			tradeSide = order.TradeSide_SELL
		}
		res.PositionList = append(res.PositionList, &client.PositionItem{
			ExchangeTime: resp.E1,
			Symbol:       symbol,
			Market:       market,
			Type:         symbolType,
			Side:         tradeSide,
			Unprofit:     unprofit,
			Position:     position,
			Price:        positionPrice,
			ChangeType:   spot_api.GetChangeType(resp.A.M),
		})
	}
	if len(b.accountChan) > cap(b.accountChan)/2 {
		content := fmt.Sprint("account channel slow:", len(b.accountChan))
		logger.Logger.Warn("AccountHandle", content)
	}
	base.SendChan(b.accountChan, res, "account")
	return nil
}

func (b *WebSocketCBaseHandle) OrderHandle(data []byte) error {
	var (
		resp                                      RespUBaseOrder
		err                                       error
		clientId                                  string
		feeAsset                                  string
		amountFilled, priceFilled, qtyFilled, fee float64
	)
	if !strings.Contains(string(data), "ORDER_TRADE_UPDATE") {
		return nil
	}
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	// 先看撤销的原始订单Id
	clientId = strconv.FormatInt(resp.O.I, 10)
	if clientId == "" {
		// 撤销的原始订单id为空，再看clientOrderId
		clientId = resp.O.C
	}
	symbol, market, symbolType := u_api.GetContractType(resp.O.S)
	amountFilled, err = strconv.ParseFloat(resp.O.Z, 64)
	if err != nil {
		logger.Logger.Error("trade parse err", err, string(data))
		return err
	}
	feeAsset = resp.O.N1
	fee, err = strconv.ParseFloat(resp.O.N, 64)
	if err != nil {
		logger.Logger.Error("parse fee error:" + err.Error())
		fee = 0
	}

	producer, id := transform.ClientIdToId(clientId)
	res := &client.WsOrderRsp{
		Producer:     producer,
		Id:           id,
		IdEx:         strconv.FormatInt(resp.O.I, 10),
		Status:       spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(resp.O.X1)),
		TimeFilled:   resp.T,
		Symbol:       symbol,
		Market:       market,
		Type:         symbolType,
		AmountFilled: amountFilled,
		PriceFilled:  priceFilled,
		QtyFilled:    qtyFilled,
		Fee:          fee,
		FeeAsset:     feeAsset,
		CloseData:    GetCloseTime(resp.O.S),
	}
	base.SendChan(b.orderChan, res, "order")
	return nil
}

func (b *WebSocketCBaseHandle) SetBookTickerGroupChannel(chMap map[*client.SymbolInfo]chan *client.WsBookTickerRsp) {
	if chMap != nil {
		b.bookTickerGroupChanMap = make(map[string]chan *client.WsBookTickerRsp)
		for info, ch := range chMap {
			sym := SymbolKeyGen(info.Symbol, info.Market, info.Type)
			b.bookTickerGroupChanMap[sym] = ch
		}
	}
}

func (b *WebSocketCBaseHandle) DepthLimitGroupHandle(data []byte) error {
	var (
		resp RespCBaseLimitDepthStream
		err  error
	)
	if err = b.HandleRespErr(data, &resp); err != nil {
		if err.Error() != "response" {
			logger.Logger.Error("receive data err:", string(data))
		}
		return err
	}
	asks, err := b.DepthLevelParse(resp.Data.A)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	bids, err := b.DepthLevelParse(resp.Data.B)
	if err != nil {
		logger.Logger.Error("depth level parse err", err, string(data))
		return err
	}
	symbol, sym := GetSymbolKey(strings.ToUpper(strings.Split(resp.Stream, "@")[0]), b.MarketType)
	res := &depth.Depth{
		Exchange:     common.Exchange_BINANCE,
		Symbol:       symbol,
		TimeExchange: uint64(resp.Data.E1) * 1000,
		TimeReceive:  uint64(time.Now().UnixMicro()),
		TimeOperate:  uint64(time.Now().UnixMicro()),
		Asks:         asks,
		Bids:         bids,
	}
	// 币本位需要转换amount
	if b.MarketType == 2 {
		for bidIdx, _ := range res.Bids {
			res.Bids[bidIdx].Amount = getCbaseQty(res.Symbol, res.Bids[bidIdx].Amount, res.Bids[bidIdx].Price)
		}
		for askIdx, _ := range res.Asks {
			res.Asks[askIdx].Amount = getCbaseQty(res.Symbol, res.Asks[askIdx].Amount, res.Asks[askIdx].Price)
		}
	}
	if _, ok := b.depthLimitGroupChanMap[sym]; ok {
		base.SendChan(b.depthLimitGroupChanMap[sym], res, fmt.Sprintf("%s %s", "depthLimitGroup", sym))
	} else {
		logger.Logger.Warn("DepthLimitGroupHandle get symbol from channel map err:", sym)
	}
	b.CheckSendStatus.CheckUpdateTimeMap.Set(sym, time.Now().UnixMicro())
	return nil
}
