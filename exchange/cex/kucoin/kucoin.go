package kucoin

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/kucoin/spot_api"
	"clients/logger"
	"clients/transform"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
)

type void struct{}

var member void

type ClientKucoin struct {
	api              *spot_api.ApiClient
	futurePrecisionM map[string]*client.PrecisionItem
}

func GetCurrentTimestamp() int64 {
	return time.Now().UnixMicro()
}

func ParsePrecision(increment string) int64 {
	reg := regexp.MustCompile(`^0\.0*1$`)
	result := reg.FindAllString(increment, -1)
	return int64(len(result[0]) - 2)
}

func NewClientKucoin(conf base.APIConf) *ClientKucoin {
	c := &ClientKucoin{
		api: spot_api.NewApiClient(conf),
	}
	c.initialFuturePrecision()
	return c
}

func (cli *ClientKucoin) initialFuturePrecision() error {
	cli.futurePrecisionM = make(map[string]*client.PrecisionItem)
	futurePrecisions, err := cli.GetFuturePrecision(common.Market_FUTURE)
	if err != nil {
		return err
	}
	for _, v := range futurePrecisions.PrecisionList {
		cli.futurePrecisionM[v.Symbol] = v
	}
	return nil
}

func Exchange2Canonical(symbol string) string {
	return strings.Replace(symbol, "-", "/", 1)
}

func Canonical2Exchange(cs string) string {
	return strings.Replace(cs, "/", "-", 1)
}

func (cli *ClientKucoin) GetExchange() common.Exchange {
	return common.Exchange_KUCOIN
}

func (cli *ClientKucoin) GetSymbols() []string {
	res, err := cli.api.GetSymbols()
	var symbols []string
	if err != nil {
		logger.Logger.Error("get symbols ", err.Error())
		return symbols
	}

	for _, symbol := range res.Data {
		if symbol.EnableTrading {
			symbols = append(symbols, Exchange2Canonical(symbol.Symbol))
		}
	}
	return symbols
}

func (cli *ClientKucoin) IsExchangeEnable() bool {
	res, err := cli.api.GetStatus()
	if err != nil || res.Data.Status != "open" {
		return false
	}
	return true
}

func (cli *ClientKucoin) GetTradeFee(symbols ...string) (*client.TradeFee, error) {
	res := &client.TradeFee{
		TradeFeeList: make([]*client.TradeFeeItem, 0, len(symbols)),
	}
	if len(symbols) == 0 {
		return res, nil
	}

	newSymbols := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		newSymbols = append(newSymbols, Canonical2Exchange(symbol))
	}
	tradeFees, err := cli.api.GetTradeFees(newSymbols)
	if err != nil {
		return res, err
	}
	for _, tradeFee := range tradeFees {

		makerFee, err := transform.Str2Float64(tradeFee.MakerFeeRate)
		if err != nil {
			return res, err
		}
		takerFee, err := transform.Str2Float64(tradeFee.TakerFeeRate)

		if err != nil {
			return res, err
		}

		tradeFeeItem := &client.TradeFeeItem{
			Symbol: Exchange2Canonical(tradeFee.Symbol),
			Type:   common.SymbolType_SPOT_NORMAL,
			Maker:  makerFee,
			Taker:  takerFee,
		}

		res.TradeFeeList = append(res.TradeFeeList, tradeFeeItem)
	}

	return res, nil
}

func (cli *ClientKucoin) GetDepth(symbolInfo *client.SymbolInfo, limit int) (*depth.Depth, error) {
	symbol := symbolInfo.Symbol
	res, err := cli.api.GetMarkets(Canonical2Exchange(symbol), limit)
	if err != nil {
		return nil, err
	}

	dep := &depth.Depth{
		Exchange:     common.Exchange_KUCOIN,
		Market:       common.Market_SPOT,
		Symbol:       symbol,
		TimeReceive:  uint64(GetCurrentTimestamp()),
		TimeExchange: res.Data.Time,
	}

	sequence, err := transform.Str2Int64(res.Data.Sequence)
	if err != nil {
		return dep, err
	}

	dep.Hdr = &common.MsgHeader{
		Sequence: sequence,
	}

	bids := make([]*depth.DepthLevel, 0, len(res.Data.Bids))
	asks := make([]*depth.DepthLevel, 0, len(res.Data.Asks))

	for i, pricesize := range res.Data.Bids {
		if i >= limit {
			break
		}
		price, err := transform.Str2Float64(pricesize[0])
		if err != nil {
			return dep, err
		}
		amount, err := transform.Str2Float64(pricesize[1])
		if err != nil {
			return dep, err
		}
		bids = append(bids, &depth.DepthLevel{Price: price, Amount: amount})
	}

	for i, pricesize := range res.Data.Asks {
		if i >= limit {
			break
		}
		price, err := transform.Str2Float64(pricesize[0])
		if err != nil {
			return dep, err
		}
		amount, err := transform.Str2Float64(pricesize[1])
		if err != nil {
			return dep, err
		}
		asks = append(asks, &depth.DepthLevel{Price: price, Amount: amount})
	}

	dep.Bids = bids
	dep.Asks = asks
	dep.TimeOperate = uint64(GetCurrentTimestamp())

	return dep, nil
}

func (cli *ClientKucoin) GetTransferFe(chain common.Chain, tokens ...string) (*client.TransferFee, error) {
	fee := &client.TransferFee{}

	res, err := cli.api.GetCurrencies()
	if err != nil {
		return fee, err
	}
	searchTokenMap := make(map[string]void)
	for _, token := range tokens {
		searchTokenMap[token] = member
	}

	for _, currencyInfo := range res.Data {
		_, ok := searchTokenMap[currencyInfo.Currency]
		if !ok {
			continue
		}

		tokenFee, err := transform.Str2Float64(currencyInfo.WithdrawalMinFee)
		if err != nil {
			logger.Logger.Error("parse fee into float64 ", currencyInfo.WithdrawalMinFee)
			return fee, err
		}
		transferFeeItem := &client.TransferFeeItem{
			Token:   currencyInfo.Currency,
			Network: chain,
			Fee:     tokenFee,
		}

		fee.TransferFeeList = append(fee.TransferFeeList, transferFeeItem)
		break
	}

	return fee, nil
}

func (cli *ClientKucoin) GetPrecision(symbols ...string) (*client.Precision, error) {
	prec := &client.Precision{}

	searchSymbolMap := make(map[string]void)
	for _, symbol := range symbols {
		searchSymbolMap[Canonical2Exchange(symbol)] = member
	}

	res, err := cli.api.GetSymbols()

	if err != nil {
		return prec, err
	}

	for _, symbolInfo := range res.Data {
		if _, ok := searchSymbolMap[symbolInfo.Symbol]; !ok {
			continue
		}

		baseMinSize, err := transform.Str2Float64(symbolInfo.BaseMinSize)
		if err != nil {
			return prec, err
		}
		basePrecision := ParsePrecision(symbolInfo.BaseIncrement)
		pricePrecision := ParsePrecision(symbolInfo.PriceIncrement)

		precisionItem := &client.PrecisionItem{
			Symbol:    Exchange2Canonical(symbolInfo.Symbol),
			Type:      common.SymbolType_SPOT_NORMAL,
			Amount:    basePrecision,
			Price:     pricePrecision,
			AmountMin: baseMinSize,
		}

		prec.PrecisionList = append(prec.PrecisionList, precisionItem)
	}

	return prec, nil
}

func (cli *ClientKucoin) GetBalance() (*client.SpotBalance, error) {
	spotBalance := &client.SpotBalance{
		UpdateTime: GetCurrentTimestamp(),
	}
	res, err := cli.api.GetAccounts("", "trade")
	PrintRes(res)
	if err != nil {
		return spotBalance, err
	}

	for _, accountInfo := range res.Data {

		balance, err := transform.Str2Float64(accountInfo.Balance)
		if err != nil {
			return spotBalance, err
		}
		available, err := transform.Str2Float64(accountInfo.Available)
		if err != nil {
			return spotBalance, err
		}
		holds, err := transform.Str2Float64(accountInfo.Holds)
		if err != nil {
			return spotBalance, err
		}

		balanceItem := &client.SpotBalanceItem{
			Asset:  accountInfo.Currency,
			Total:  balance,
			Free:   available,
			Frozen: holds,
		}
		spotBalance.BalanceList = append(spotBalance.BalanceList, balanceItem)
	}
	return spotBalance, nil
}

func (cli *ClientKucoin) GetMarginBalance() (*client.MarginBalance, error) {
	marginBalance := &client.MarginBalance{
		UpdateTime: GetCurrentTimestamp(),
	}

	res, err := cli.api.GetMarginAcounts()
	PrintRes(res)
	if err != nil {
		return marginBalance, err
	}

	for _, accountInfo := range res.Data.Accounts {

		marginBalanceItem, err := cli.ParseMarginBalanceItem(&accountInfo)
		if err != nil {
			return marginBalance, err
		}
		marginBalance.MarginBalanceList = append(marginBalance.MarginBalanceList, marginBalanceItem)
	}

	return marginBalance, nil
}

func (cli *ClientKucoin) ParseMarginBalanceItem(accountInfo *spot_api.MarginAccountInfo) (*client.MarginBalanceItem, error) {
	available, err := transform.Str2Float64(accountInfo.AvailableBalance)
	if err != nil {
		return nil, err
	}

	hold, err := transform.Str2Float64(accountInfo.HoldBalance)
	if err != nil {
		return nil, err
	}

	total, err := transform.Str2Float64(accountInfo.TotalBalance)
	if err != nil {
		return nil, err
	}

	liability, err := transform.Str2Float64(accountInfo.Liability)
	if err != nil {
		return nil, err
	}

	marginBalanceItem := &client.MarginBalanceItem{
		Asset:    accountInfo.Currency,
		Total:    total,
		Free:     available,
		Frozen:   hold,
		Borrowed: liability,
		NetAsset: available - liability,
	}
	return marginBalanceItem, nil
}

func (cli *ClientKucoin) GetMarginIsolatedBalance(symbols ...string) (*client.MarginIsolatedBalance, error) {

	isolatedBalance := &client.MarginIsolatedBalance{
		UpdateTime: GetCurrentTimestamp(),
	}
	res, err := cli.api.GetIsolatedAccounts("BTC")

	searchSymbolMap := make(map[string]void)
	for _, symbol := range symbols {
		searchSymbolMap[Canonical2Exchange(symbol)] = member
	}

	if err != nil {
		return isolatedBalance, err
	}

	totalConversionBalance, err := transform.Str2Float64(res.Data.TotalConversionBalance)
	if err != nil {
		return isolatedBalance, err
	}
	liabilityConversionBalance, err := transform.Str2Float64(res.Data.LiabilityConversionBalance)
	if err != nil {
		return isolatedBalance, err
	}

	isolatedBalance.TotalNetAsset = totalConversionBalance
	isolatedBalance.QuoteAsset = "BTC"
	isolatedBalance.TotalLiabilityAsset = liabilityConversionBalance

	for _, asset := range res.Data.Assets {
		if _, ok := searchSymbolMap[asset.Symbol]; !ok {
			continue
		}

		baseAsset, err := cli.ParseMarginBalanceItem(&asset.BaseAsset)
		if err != nil {
			return isolatedBalance, err
		}
		quoteAsset, err := cli.ParseMarginBalanceItem(&asset.QuoteAsset)
		if err != nil {
			return isolatedBalance, err
		}
		isolatedBalanceItem := &client.MarginIsolatedBalanceItem{
			BaseAsset:  baseAsset,
			QuoteAsset: quoteAsset,
		}

		isolatedBalance.MarginIsolatedBalanceList = append(isolatedBalance.MarginIsolatedBalanceList, isolatedBalanceItem)
	}

	return isolatedBalance, nil
}

// 必须有producer
func (cli *ClientKucoin) PlaceOrder(req *order.OrderTradeCEX) (*client.OrderRsp, error) {
	var (
		symbol    string
		clientOid string
		side      string
		options   = make(map[string]interface{})
		placeRes  *spot_api.RespPlaceOrder
		resCode   int
		resp      = &client.OrderRsp{}
		err       error
	)

	orderBase := req.Base
	symbol = Canonical2Exchange(string(orderBase.Symbol))
	clientOid = transform.IdToClientId(req.Hdr.Producer, orderBase.Id)

	switch req.Side {
	case order.TradeSide_BUY:
		side = "buy"
	case order.TradeSide_SELL:
		side = "sell"
	default:
		return nil, errors.New("unsupported tradeside " + req.Side.String())
	}

	if req.OrderType == order.OrderType_LIMIT || req.OrderType == order.OrderType_LIMIT_MAKER { // 限价单
		options["type"] = "limit"
		options["price"] = fmt.Sprintf("%f", req.Price)
		options["size"] = fmt.Sprintf("%f", req.Amount)
		switch req.Tif {
		case order.TimeInForce_GTC:
			options["timeInForce"] = "GTC"
		case order.TimeInForce_IOC:
			options["timeInForce"] = "IOC"
		case order.TimeInForce_FOK:
			options["timeInForce"] = "FOK"
		default:
			return nil, errors.New("unsupported tif " + req.Tif.String())
		}

		switch req.TradeType {
		case order.TradeType_MAKER:
			options["postOnly"] = true
		case order.TradeType_TAKER:
			options["postOnly"] = false
		default:
			return nil, errors.New("invalid tradetype " + req.TradeType.String())
		}

	} else if req.OrderType == order.OrderType_MARKET { // 市价单
		options["type"] = "market"
		options["size"] = fmt.Sprintf("%f", req.Amount)
	} else {
		return nil, errors.New("unsupported order type " + req.OrderType.String())
	}

	if orderBase.Market == common.Market_SPOT {
		placeRes, err = cli.api.PlaceOrder(clientOid, side, symbol, options) // 下单
		resCode = transform.StringToX[int](placeRes.Code).(int)
		PrintRes(placeRes)
		if resCode != 200000 {
			return nil, errors.New("购买出错")
		}
	} else if orderBase.Market == common.Market_MARGIN {
		switch orderBase.Type {
		case common.SymbolType_MARGIN_NORMAL:
			options["marginMode"] = "cross"
		case common.SymbolType_MARGIN_ISOLATED:
			options["marginMode"] = "isolated"
		default:
			return nil, errors.New("unsupported symbol type " + orderBase.Type.String())
		}
		_, err = cli.api.PlaceMarginOrder(clientOid, side, symbol, options)
	}

	if err != nil {
		return nil, err
	}

	res, err := cli.api.GetOrderByClientOid(clientOid) // 查询订单信息
	if err != nil {
		return nil, err
	}
	orderInfo := res.Data

	resp.Producer, resp.Id = transform.ClientIdToId(clientOid)
	resp.OrderId = orderInfo.Id
	resp.Timestamp = orderInfo.CreatedAt
	resp.RespType = client.OrderRspType_FULL
	resp.Symbol = Exchange2Canonical(orderInfo.Symbol)
	if orderInfo.IsActive { // 判断订单状态
		resp.Status = order.OrderStatusCode_OPENED
	} else {
		if orderInfo.CancelExist {
			resp.Status = order.OrderStatusCode_CANCELED
		} else {
			resp.Status = order.OrderStatusCode_FILLED
		}
	}

	resp.AccumQty, err = transform.Str2Float64(orderInfo.DealFunds)
	if err != nil {
		return nil, err
	}
	resp.AccumAmount, err = transform.Str2Float64(orderInfo.DealSize)
	if err != nil {
		return nil, err
	}
	resp.Fee, err = transform.Str2Float64(orderInfo.Fee)
	if err != nil {
		return nil, err
	}

	resp.FeeAsset = orderInfo.FeeCurrency
	if resp.AccumAmount > 0 {
		getFillParams := make(map[string]interface{})
		getFillParams["orderId"] = orderInfo.Id
		res, err := cli.api.GetFills(true, getFillParams) // 查询所有成交信息
		if err != nil {
			return resp, err
		}
		for _, fill := range res {
			fillItem := &client.FillItem{}
			fillItem.Price, err = transform.Str2Float64(fill.Price)
			if err != nil {
				return resp, err
			}
			fillItem.Qty, err = transform.Str2Float64(fill.Size)
			if err != nil {
				return resp, err
			}
			fillItem.Commission, err = transform.Str2Float64(fill.Fee)
			if err != nil {
				return resp, err
			}
			fillItem.CommissionAsset = fill.FeeCurrency
			resp.Fills = append(resp.Fills, fillItem)
		}
	}
	return resp, nil
}

func PrintRes(x interface{}) {
	res, _ := json.Marshal(x)
	fmt.Println("结果为:", string(res))
}

func (cli *ClientKucoin) CancelOrder(req *order.OrderCancelCEX) (*client.OrderRsp, error) {
	clientOid := transform.IdToClientId(req.Hdr.Producer, req.Base.Id)
	res, err := cli.api.CancelOrderByClientID(clientOid)
	if err != nil {
		return nil, err
	}
	resp := &client.OrderRsp{}

	resp.Producer, resp.Id = transform.ClientIdToId(clientOid)
	resp.OrderId = res.Data.CancelledOrderId
	resp.RespType = client.OrderRspType_FULL
	resp.Timestamp = GetCurrentTimestamp()
	return resp, nil
}

// 使用的是交易网站返回的id
func (cli *ClientKucoin) GetOrder(req *order.OrderQueryReq) (*client.OrderRsp, error) {
	var (
		resp    = &client.OrderRsp{}
		orderId = req.IdEx //orderId	LONG	NO
	)
	if orderId == "" {
		return nil, errors.New("orderId不能为空")
	}

	res, err := cli.api.GetOrder(orderId)
	if err != nil {
		return nil, err
	}

	// 创建时间
	if err != nil {
		return nil, err
	}
	resp.Producer, resp.Id = transform.ClientIdToId(res.Data.ClientOid)
	resp.OrderId = res.Data.Id
	resp.Timestamp = res.Data.CreatedAt

	//resp. = ParseInstType(res.Data[0].InstType)
	resp.RespType = client.OrderRspType_RESULT
	resp.Symbol = strings.ReplaceAll(res.Data.Symbol, "-", "/")
	resp.Status = parseStatusOrder(res.Data.IsActive)

	resp.AccumAmount, _ = strconv.ParseFloat(res.Data.DealSize, 64)
	resp.AccumQty, _ = strconv.ParseFloat(res.Data.DealFunds, 64)
	//resp.Executed, _ = strconv.ParseFloat(res.Data.DealSize, 64)
	resp.AvgPrice = resp.AccumQty / resp.AccumAmount
	return resp, nil
}

func parseStatusOrder(isActive bool) order.OrderStatusCode {
	if isActive {
		return order.OrderStatusCode_OPENED
	} else {
		return order.OrderStatusCode_FILLED
	}
}

func (cli *ClientKucoin) GetOrderHistory(req *client.OrderHistoryReq) ([]*client.OrderRsp, error) {
	var (
		resp []*client.OrderRsp
		// todo Market是int32类型
		instType  = req.Market    //symbol	STRING	YES
		startTime = req.StartTime //startTime	LONG	NO
		endTime   = req.EndTime   //endTime	LONG	NO
		params    = url.Values{}
	)

	if endTime > 0 {
		params.Add("startAt", strconv.FormatInt(endTime, 10))
	}
	if startTime > 0 {
		params.Add("endAt", strconv.FormatInt(startTime, 10))
	}

	if instType != common.Market_INVALID_MARKET {
		params.Add("tradeType", parseType(instType))
	}

	// 近三个月
	res, err := cli.api.GetOrderHistory(params)
	//fmt.Println("res:", res)
	if err != nil {
		return nil, err
	}
	for _, item := range res.Data.Items {
		producer, id := transform.ClientIdToId(item.ClientOid)
		price, _ := strconv.ParseFloat(item.Price, 64)
		symbol := strings.ReplaceAll(item.Symbol, "-", "")
		executed, _ := strconv.ParseFloat(item.DealSize, 64)
		avgPrice := transform.StringToX[float64](item.DealFunds).(float64) / transform.StringToX[float64](item.DealSize).(float64)
		accumAmount, _ := strconv.ParseFloat(item.DealSize, 64)
		ts := item.CreatedAt
		names := strings.Split(item.Symbol, "-")
		var closeDate string
		if len(names) == 3 {
			if names[2] != "SWAP" {
				closeDate = names[2]
			}
		}
		resp = append(resp, &client.OrderRsp{
			Producer:    producer,
			Id:          id,
			OrderId:     item.Id,
			Timestamp:   ts,
			Symbol:      symbol,
			RespType:    client.OrderRspType_RESULT,
			Status:      GetOrderStatusFromExchange(item.IsActive),
			Price:       price,
			Executed:    executed,
			AvgPrice:    avgPrice,
			AccumAmount: accumAmount,
			AccumQty:    avgPrice * accumAmount,
			CloseDate:   closeDate,
		})
	}
	return resp, nil
}

func parseType(instType common.Market) string {
	switch instType {
	case common.Market_SPOT:
		return "TRADE"
	case common.Market_MARGIN:
		return "MARGIN_TRADE"
	default:
		return "TRADE"
	}
}

func GetOrderStatusFromExchange(status bool) order.OrderStatusCode {
	if status {
		return order.OrderStatusCode_OPENED
	}
	return order.OrderStatusCode_FILLED
}

func (cli ClientKucoin) GetProcessingOrders(req *client.OrderHistoryReq) ([]*client.OrderRsp, error) {
	var (
		resp []*client.OrderRsp
		// todo Market是int32类型
		instType  = req.Market    //symbol	STRING	YES
		startTime = req.StartTime //startTime	LONG	NO
		endTime   = req.EndTime   //endTime	LONG	NO
		params    = url.Values{}
	)

	if endTime > 0 {
		params.Add("startAt", strconv.FormatInt(endTime, 10))
	}
	if startTime > 0 {
		params.Add("endAt", strconv.FormatInt(startTime, 10))
	}

	if instType != common.Market_INVALID_MARKET {
		params.Add("tradeType", parseType(instType))
	}

	// 近三个月
	res, err := cli.api.GetOrderHistory(params)
	//fmt.Println("res:", res)
	if err != nil {
		return nil, err
	}
	for _, item := range res.Data.Items {
		if item.IsActive {
			producer, id := transform.ClientIdToId(item.ClientOid)
			price, _ := strconv.ParseFloat(item.Price, 64)
			symbol := strings.ReplaceAll(item.Symbol, "-", "")
			executed, _ := strconv.ParseFloat(item.DealSize, 64)
			avgPrice := transform.StringToX[float64](item.DealFunds).(float64) / transform.StringToX[float64](item.DealSize).(float64)
			accumAmount, _ := strconv.ParseFloat(item.DealSize, 64)
			ts := item.CreatedAt
			names := strings.Split(item.Symbol, "-")
			var closeDate string
			if len(names) == 3 {
				if names[2] != "SWAP" {
					closeDate = names[2]
				}
			}
			resp = append(resp, &client.OrderRsp{
				Producer:    producer,
				Id:          id,
				OrderId:     item.Id,
				Timestamp:   ts,
				Symbol:      symbol,
				RespType:    client.OrderRspType_RESULT,
				Status:      GetOrderStatusFromExchange(item.IsActive),
				Price:       price,
				Executed:    executed,
				AvgPrice:    avgPrice,
				AccumAmount: accumAmount,
				AccumQty:    avgPrice * accumAmount,
				CloseDate:   closeDate,
			})
		}
	}
	return resp, nil
}

func (cli ClientKucoin) Transfer(o *order.OrderTransfer) (*client.OrderRsp, error) {
	//默认都是从现货提币
	var (
		resp = &client.OrderRsp{}
		err  error
		coin = string(o.ExchangeToken)
		//withdrawOrderId = o.Base.Id //自定义提币ID
		//network         = coin + "-" + ok_api.GetNetWorkFromChain(o.Chain) //提币网络
		amount = o.Amount //数量
		params = url.Values{}
		fee    string
	)

	params.Add("currency", coin)
	b, _ := cli.api.WithdrawFee(coin)
	fee = b.Data.WithdrawMinFee

	params.Add("address", string(o.TransferAddress))
	//params.Add("clientId", transform.IdToClientId(o.Hdr.Producer, withdrawOrderId))
	fee_ := transform.StringToX[float64](fee).(float64)
	amt := fmt.Sprintf("%f", amount-fee_)
	//var dest string
	//if o.Base.Exchange == o.ExchangeTo {
	//	dest = "3"
	//} else {
	//	dest = "4"
	//}
	params.Add("amount", transform.XToString(amt))

	address := o.TransferAddress
	if len(o.Tag) != 0 {
		address = append(address, ':')
		address = append(address, o.Tag...)
	}

	res, err := cli.api.Withdraw(coin, string(address), amount, params)
	if err != nil {
		return nil, err
	}

	resp.OrderId = res.Data.WithdrawalId
	resp.RespType = client.OrderRspType_ACK
	return resp, err

}

func (cli *ClientKucoin) MoveAsset(req *order.OrderMove) (*client.OrderRsp, error) {

	var (
		clientOid string
		currency  string
		amount    string
		options   = make(map[string]interface{})
		from      string
		to        string
		resp      = &client.OrderRsp{}
	)
	clientOid = transform.IdToClientId(req.Hdr.Producer, req.Base.Id)
	currency = req.Asset
	amount = fmt.Sprintf("%f", req.Amount)
	switch req.Source {
	case common.Market_SPOT:
		from = "trade"
	case common.Market_MARGIN:
		from = "margin"
	default:
		return nil, errors.New("unsupported MoveAsset source " + req.Source.String())
	}
	switch req.Source {
	case common.Market_SPOT:
		to = "trade"
	case common.Market_MARGIN:
		to = "margin"
	case common.Market_FUTURE:
		to = "contract"
	default:
		return nil, errors.New("unsupported MoveAsset source " + req.Source.String())
	}

	if from == "margin" && req.SymbolSource != "" {
		from = "isolated"
		options["fromTag"] = Canonical2Exchange(req.SymbolSource)
	}

	if to == "margin" && req.SymbolSource != "" {
		to = "isolated"
		options["toTag"] = Canonical2Exchange(req.SymbolSource)
	}
	res, err := cli.api.InnerTransfer(clientOid, currency, from, to, amount, options)
	if err != nil {
		return resp, err
	}

	resp.Producer, resp.Id = transform.ClientIdToId(clientOid)
	resp.OrderId = res.Data.OrderId
	resp.RespType = client.OrderRspType_ACK
	resp.Timestamp = GetCurrentTimestamp()

	return resp, nil
}

func (cli *ClientKucoin) GetTransferHistory(req *client.TransferHistoryReq) (*client.TransferHistoryRsp, error) {
	var (
		resp            = &client.TransferHistoryRsp{}
		err             error
		coin            = req.Asset //coin	STRING	NO
		withdrawOrderId string      //withdrawOrderId	STRING	NO
		//offset          = 0             //offset	INT	NO
		startTime = req.StartTime //startTime	LONG	NO	默认当前时间90天前的时间戳
		endTime   = req.EndTime   //endTime	LONG	NO	默认当前时间戳
		params    = url.Values{}
		res       *spot_api.RespWithdrawHistory
	)

	if req.Id != 0 {
		withdrawOrderId = transform.IdToClientId(req.Producer, req.Id)
	}
	if coin != "" {
		params.Add("coin", coin)
	}
	if withdrawOrderId != "" {
		params.Add("withdraw_id", withdrawOrderId)
	}
	if startTime > 0 {
		params.Add("start_time", strconv.FormatInt(startTime, 10))
	}
	if endTime > 0 {
		params.Add("end_time", strconv.FormatInt(endTime, 10))
	}
	res, err = cli.api.WithdrawHistory(params)
	if err != nil {
		return nil, err
	}
	for _, item := range res.Data.Items {
		timestamp := item.CreateAt
		amount, _ := strconv.ParseFloat(item.Amount, 64)
		if err != err {
			return nil, err
		}

		resp.TransferList = append(resp.TransferList, &client.TransferHistoryItem{
			Asset:   item.Currency,
			Amount:  amount,
			OrderId: item.WalletTxId,
			Status:  parseStatus(item.Status),
			// 没有确认数
			Timestamp: int64(timestamp),
		})
	}
	return resp, err
}

func parseStatus(status string) client.TransferStatus {
	switch status {
	case "SUCCESS":
		return client.TransferStatus_TRANSFERSTATUS_COMPLETE
	case "FAILURE":
		return client.TransferStatus_TRANSFERSTATUS_FAILED
	case "PROCESSING":
		return client.TransferStatus_TRANSFERSTATUS_PROCESSING
	default:
		return client.TransferStatus_TRANSFERSTATUS_INVALID
	}
}

func (cli *ClientKucoin) GetMoveHistory(req *client.MoveHistoryReq) (*client.MoveHistoryRsp, error) {
	var (
		resp      = &client.MoveHistoryRsp{}
		err       error
		startTime = req.StartTime //startTime	LONG	NO
		endTime   = req.EndTime   //endTime	LONG	NO
		params    = url.Values{}
	)
	if startTime > 0 {
		params.Add("startAt", strconv.FormatInt(startTime, 10))
	}
	if endTime > 0 {
		params.Add("endAt", strconv.FormatInt(endTime, 10))
	}

	if req.ActionUser == order.OrderMoveUserType_Master {
		// 子母账户万能划转历史
		params.Add("bizType", "SUB_TRANSFER")
		res, err := cli.api.MoveHistory(params)
		if err != nil {
			return nil, err
		}

		for _, item := range res.Data.Items {
			type_ := client.MoveType_MOVETYPE_IN
			if item.Direction == "out" {
				type_ = client.MoveType_MOVETYPE_OUT
			}
			amount, _ := strconv.ParseFloat(item.Amount, 64)
			resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
				Asset:     item.Currency,
				Amount:    amount,
				Timestamp: item.CreatedAt,
				Type:      type_,
			})
		}
	} else if req.ActionUser == order.OrderMoveUserType_Sub {
		// 子账户划转历史
		// 没有子子划转
		return nil, nil
	} else { //OrderMoveUserType_Internal
		// 主账户划转历史
		//  转入划转
		params.Add("bizType", "TRANSFER")
		res, err := cli.api.MoveHistory(params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Data.Items {
			amount, _ := strconv.ParseFloat(item.Amount, 64)
			type_ := client.MoveType_MOVETYPE_IN
			if item.Direction == "out" {
				type_ = client.MoveType_MOVETYPE_OUT
			}
			resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
				Asset:     item.Currency,
				Amount:    amount,
				Timestamp: item.CreatedAt,
				Type:      type_,
			})
		}
	}
	return resp, err
}

func (cli *ClientKucoin) GetDepositHistory(req *client.DepositHistoryReq) (*client.DepositHistoryRsp, error) {
	var (
		resp = &client.DepositHistoryRsp{}
		err  error
		//status    = spot_api.GetDepositTypeToExchange(req.Status) //status	INT	NO	0(0:pending,6: credited but cannot withdraw, 1:success)
		startTime = req.StartTime //startTime	LONG	NO	默认当前时间90天前的时间戳
		endTime   = req.EndTime   //endTime	LONG	NO	默认当前时间戳
		offset    = 0             //offset	INT	NO	默认:0
		limit     = 50            //limit	INT	NO	默认：50，最大50
		params    = url.Values{}
		res       *spot_api.RespMoveHistory
	)

	if startTime > 0 {
		params.Add("startAt", strconv.FormatInt(startTime, 10))
	}
	if endTime > 0 {
		params.Add("endAt", strconv.FormatInt(endTime, 10))
	}

	params.Add("bizType", "DEPOSIT")

	res, err = cli.api.MoveHistory(params)
	if err != nil {
		return nil, err
	}
	for i := 0; i < 100; i++ {
		for _, item := range res.Data.Items {

			if err != nil {
				return nil, err
			}
			//state := item.Status
			resp.DepositList = append(resp.DepositList, &client.DepositHistoryItem{
				Asset:  item.Currency,
				Amount: transform.StringToX[float64](item.Amount).(float64),
				//Status:    client.DepositStatus_DEPOSITSTATUS_SUCCESS,
				//Address:   item.ToAddress,
				//TxId:      item.TxId,
				Timestamp: item.CreatedAt,
			})
		}
		if len(res.Data.Items) < limit {
			break
		}
		offset += limit
	}
	return resp, err
}

func (cli *ClientKucoin) Loan(o *order.OrderLoan) (*client.OrderRsp, error) {
	var (
		resp   = &client.OrderRsp{}
		err    error
		asset  = o.Asset  //asset	STRING	YES
		amount = o.Amount //amount	DECIMAL	YES
		params = url.Values{}
	)

	_, err = cli.api.Loan(asset, "FOK", amount, params)
	if err != nil {
		return nil, err
	}
	//resp.OrderId = strconv.FormatInt(res.TranId, 10)
	resp.RespType = client.OrderRspType_ACK
	return resp, err
}

func (cli *ClientKucoin) GetLoanOrders(o *client.LoanHistoryReq) (*client.LoanHistoryRsp, error) {
	var (
		resp   = &client.LoanHistoryRsp{}
		err    error
		params = url.Values{}
	)

	if o.Asset != "" {
		params.Add("currency", o.Asset)
	}

	resRepay, err := cli.api.RepayHistory(params)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	for _, item := range resRepay.Data.Items {
		// 2， 3为还币
		ts, err := strconv.ParseInt(item.RepayTime, 10, 64)
		principle, err := strconv.ParseFloat(item.Principal, 64)
		if err != nil {
			return nil, err
		}
		resp.LoadList = append(resp.LoadList, &client.LoanHistoryItem{
			Principal: principle,
			OrderId:   item.TradeId,
			Asset:     item.Currency,
			Timestamp: ts,
		})
	}
	resUnrepay, err := cli.api.UnRepayHistory(params)
	for _, item := range resUnrepay.Data.Items {
		// 2， 3为还币
		ts, err := strconv.ParseInt(item.CreatedAt, 10, 64)
		principle, err := strconv.ParseFloat(item.Principal, 64)
		if err != nil {
			return nil, err
		}
		resp.LoadList = append(resp.LoadList, &client.LoanHistoryItem{
			Principal: principle,
			OrderId:   item.TradeId,
			Asset:     item.Currency,
			Timestamp: ts,
		})
	}

	return resp, err
}

func (cli *ClientKucoin) Return(o *order.OrderReturn) (*client.OrderRsp, error) {
	var (
		resp   = &client.OrderRsp{}
		err    error
		asset  = o.Asset  //asset	STRING	YES
		amount = o.Amount //amount	DECIMAL	YES
		params = url.Values{}
	)

	res, err := cli.api.Repay(asset, "tradId", amount, params)
	if err != nil {
		return nil, err
	}
	resp.RespType = client.OrderRspType_ACK
	if res.Code == "200000" {
		resp.Status = order.OrderStatusCode_FILLED
	} else {
		resp.Status = order.OrderStatusCode_FAILED
	}
	return resp, err
}

func (cli *ClientKucoin) GetReturnOrders(o *client.LoanHistoryReq) (*client.ReturnHistoryRsp, error) {
	var (
		resp   = &client.ReturnHistoryRsp{}
		err    error
		params = url.Values{}
	)

	if o.Asset != "" {
		params.Add("currency", o.Asset)
	}

	res, err := cli.api.RepayHistory(params)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	for _, item := range res.Data.Items {
		// 2， 3为还币
		ts, err := strconv.ParseInt(item.RepayTime, 10, 64)
		principle, err := strconv.ParseFloat(item.Principal, 64)
		if err != nil {
			return nil, err
		}
		resp.ReturnList = append(resp.ReturnList, &client.ReturnHistoryItem{
			Principal: principle,
			Interest:  transform.StringToX[float64](item.Interest).(float64),
			Asset:     item.Currency,
			Timestamp: ts,
		})
	}
	return resp, err
}

func (cli *ClientKucoin) GetTransferFee(exchange common.Chain, tokens ...string) (*client.TransferFee, error) {
	fee := &client.TransferFee{}

	res, err := cli.api.GetCurrencies()
	if err != nil {
		return fee, err
	}
	searchTokenMap := make(map[string]void)
	for _, token := range tokens {
		searchTokenMap[token] = member
	}

	for _, currencyInfo := range res.Data {
		_, ok := searchTokenMap[currencyInfo.Currency]
		if !ok {
			continue
		}

		tokenFee, err := transform.Str2Float64(currencyInfo.WithdrawalMinFee)
		if err != nil {
			logger.Logger.Error("parse fee into float64 ", currencyInfo.WithdrawalMinFee)
			return fee, err
		}
		transferFeeItem := &client.TransferFeeItem{
			Token:   currencyInfo.Currency,
			Network: exchange,
			Fee:     tokenFee,
		}

		fee.TransferFeeList = append(fee.TransferFeeList, transferFeeItem)
		break
	}

	return fee, nil
}

func (cli *ClientKucoin) GetFutureSymbols(common.Market) []*client.SymbolInfo { //所有交易对
	res := make([]*client.SymbolInfo, 0)
	exchangeInfoRes, err := cli.api.GetFutureSymbols()
	//PrintRes(res)
	if err != nil {
		logger.Logger.Error("get symbols ", err.Error())
		return nil
	}

	for _, symbol := range exchangeInfoRes.Data {
		symbol_ := symbol.BaseCurrency + "/" + symbol.QuoteCurrency
		type_, name_3 := ParseDate(ParseTime(symbol.ExpireDate), symbol.BaseCurrency)
		name := symbol_ + "/" + name_3
		res = append(res, &client.SymbolInfo{
			Symbol: symbol_,
			Name:   name,
			Type:   type_,
		})
	}
	return res
}

func ParseTime(ts interface{}) string {
	if ts == nil {
		return ""
	}
	t := strings.Split(time.Unix(int64(ts.(float64)), 0).Format("2006-01-02"), "-")
	// 20 + 01 + 02
	res := t[0][:2] + t[0] + t[1]
	return res
}

func ParseDate(f interface{}, base string) (common.SymbolType, string) {
	/*
		https://www.binance.com/zh-CN/support/faq/3ae441db4ae740e19af3fe9228eb6619
		季度交割合约是具备固定到期日和交割日的衍生品合约，以每个季度的最后一个周五作爲交割日。例如：“BTCUSDT 当季 0326”代表 2021年03月26日16:00（香港时间）进行交割。
		当季度合约结算交割后，会产生新的季度合约，例如：“BTCUSDT 当季 0326” 于 2021年03月26日16:00（香港时间）交割下架后将会生成新的“BTCUSDT 当季 0625”合约，以此类推；
		交割日结算时将收取交割费，交割费与Taker(吃单)费率相同。
	*/
	var u bool
	date := transform.XToString(f)
	if base == "USDT" {
		u = true
	}
	switch date {
	case transform.GetDate(transform.THISWEEK):
		if u {
			return common.SymbolType_FUTURE_THIS_WEEK, date
		} else {
			return common.SymbolType_FUTURE_COIN_NEXT_WEEK, date
		}
	case transform.GetDate(transform.NEXTWEEK):
		if u {
			return common.SymbolType_FUTURE_NEXT_WEEK, date
		} else {
			return common.SymbolType_FUTURE_COIN_NEXT_QUARTER, date
		}
		//221230
	case transform.GetDate(transform.THISQUARTER):
		if u {
			return common.SymbolType_FUTURE_THIS_QUARTER, date
		} else {
			return common.SymbolType_FUTURE_COIN_THIS_QUARTER, date
		}
	case transform.GetDate(transform.NEXTQUARTER):
		if u {
			return common.SymbolType_FUTURE_NEXT_QUARTER, date
		} else {
			return common.SymbolType_FUTURE_COIN_NEXT_QUARTER, date
		}
	case "":
		if u {
			return common.SymbolType_SWAP_FOREVER, "SWAP"
		} else {
			return common.SymbolType_SWAP_COIN_FOREVER, "SWAP"
		}
	default:
		return common.SymbolType_INVALID_TYPE, ""
	}
}

func (cli *ClientKucoin) GetFutureDepth(symbolInfo *client.SymbolInfo, limit int) (*depth.Depth, error) { //获取行情
	if !isValidFutureSymbols(symbolInfo) {
		return nil, errors.New("kucoin不支持除BTC的季度交割外的其他交割合约")
	}
	symbolName := getSymbolName(symbolInfo)
	res, err := cli.api.GetFutureMarkets(symbolName, limit)
	if err != nil {
		return nil, err
	}
	s := strings.Split(symbolInfo.Symbol, "/")
	dep := &depth.Depth{
		Exchange:    common.Exchange_KUCOIN,
		Market:      symbolInfo.Market,
		Type:        symbolInfo.Type,
		Symbol:      XBT2BTC(s[0] + "/" + s[1]),
		TimeOperate: uint64(time.Now().UnixMicro()),
	}
	err = ParseOrder(res.Data.Asks, &dep.Bids)
	if err != nil {
		return dep, err
	}

	err = ParseOrder(res.Data.Bids, &dep.Asks)
	if err != nil {
		return dep, err
	}
	intNum := res.Data.Ts
	dep.TimeReceive = uint64(intNum)
	return dep, err
}

func getSymbolName(symbolInfo *client.SymbolInfo) string {
	symbol := symbolInfo.Symbol
	symbol = BTC2XBT(symbol)
	s := strings.Split(symbol, "/")
	symbolName := s[0] + s[1]
	if symbolInfo.Type == common.SymbolType_FUTURE_THIS_QUARTER && symbolInfo.Symbol == "BTC/USD" {
		symbolName = "XBTMZ22"
	} else {
		symbolName += "M"
	}
	return symbolName
}

func isValidFutureSymbols(symbolInfo *client.SymbolInfo) bool {
	if symbolInfo.Market == common.Market_SWAP || symbolInfo.Market == common.Market_SWAP_COIN {
		return true
	} else if symbolInfo.Type == common.SymbolType_FUTURE_THIS_QUARTER && symbolInfo.Symbol == "BTC/USD" {
		return true
	} else {
		return false
	}
}

func ParseOrder(orders [][]interface{}, slice *[]*depth.DepthLevel) error {
	for _, order := range orders {
		price, amount, err := ParsePriceAmountFloat(order)
		if err != nil {
			logger.Logger.Errorf("order float parse price error [%s] , response data = %s", err, order)
			return err
		}
		*slice = append(*slice, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	return nil
}

func ParsePriceAmountFloat(data []interface{}) (price float64, amount float64, err error) {
	price = data[0].(float64)
	amount = data[1].(float64)
	return
}

func BTC2XBT(s string) string {
	return strings.ReplaceAll(s, "BTC", "XBT")
}

func XBT2BTC(s string) string {
	return strings.ReplaceAll(s, "XBT", "BTC")
}

func (cli *ClientKucoin) GetFutureMarkPrice(market common.Market, symbols ...*client.SymbolInfo) (*client.RspMarkPrice, error) { //标记价格
	res := &client.RspMarkPrice{}
	if symbols == nil && (market == common.Market_FUTURE || market == common.Market_FUTURE_COIN) {
		symbolInfo, err := cli.api.GetFutureSymbols()
		if err != nil {
			return nil, err
		}
		for _, info := range symbolInfo.Data {
			type_, _ := ParseDate(ParseTime(info.ExpireDate), info.BaseCurrency)
			res.Item = append(res.Item, &client.MarkPriceItem{
				Symbol:    XBT2BTC(info.BaseCurrency) + "/" + info.QuoteCurrency,
				Type:      type_,
				MarkPrice: info.MarkPrice,
			})
		}
	} else {
		// 用单独查询的
		for i := range symbols {
			symbolName := getSymbolName(symbols[i])
			info, err := cli.api.GetFutureSpecifySymbol(symbolName)
			if err != nil {
				return nil, err
			}
			type_, _ := ParseDate(ParseTime(info.Data.ExpireDate), info.Data.BaseCurrency)
			res.Item = append(res.Item, &client.MarkPriceItem{
				Symbol:    XBT2BTC(info.Data.BaseCurrency) + "/" + info.Data.QuoteCurrency,
				Type:      type_,
				MarkPrice: info.Data.MarkPrice,
			})
		}
	}
	return res, nil
}
func (cli *ClientKucoin) GetFutureTradeFee(market common.Market, symbols ...*client.SymbolInfo) (*client.TradeFee, error) { //查询交易手续费
	res := &client.TradeFee{}
	if symbols == nil && (market == common.Market_FUTURE || market == common.Market_FUTURE_COIN) {
		symbolInfo, err := cli.api.GetFutureSymbols()
		if err != nil {
			return nil, err
		}
		for _, info := range symbolInfo.Data {
			type_, _ := ParseDate(ParseTime(info.ExpireDate), info.BaseCurrency)
			res.TradeFeeList = append(res.TradeFeeList, &client.TradeFeeItem{
				Symbol: XBT2BTC(info.BaseCurrency) + "/" + info.QuoteCurrency,
				Type:   type_,
				Maker:  info.MakerFeeRate,
				Taker:  info.TakerFeeRate,
			})
		}
	} else {
		// 用单独查询的
		for i := range symbols {
			symbolName := getSymbolName(symbols[i])
			info, err := cli.api.GetFutureSpecifySymbol(symbolName)
			if err != nil {
				return nil, err
			}
			type_, _ := ParseDate(ParseTime(info.Data.ExpireDate), info.Data.BaseCurrency)
			res.TradeFeeList = append(res.TradeFeeList, &client.TradeFeeItem{
				Symbol: XBT2BTC(info.Data.BaseCurrency) + "/" + info.Data.QuoteCurrency,
				Type:   type_,
				Maker:  info.Data.MakerFeeRate,
				Taker:  info.Data.TakerFeeRate,
			})
		}
	}
	return res, nil
}
func (cli *ClientKucoin) GetFuturePrecision(market common.Market, symbols ...*client.SymbolInfo) (*client.Precision, error) { //查询交易对精读信息
	res := &client.Precision{}
	if symbols == nil && (market == common.Market_FUTURE || market == common.Market_FUTURE_COIN) {
		symbolInfo, err := cli.api.GetFutureSymbols()
		if err != nil {
			return nil, err
		}
		for _, info := range symbolInfo.Data {
			type_, _ := ParseDate(ParseTime(info.ExpireDate), info.BaseCurrency)
			res.PrecisionList = append(res.PrecisionList, &client.PrecisionItem{
				Symbol:    XBT2BTC(info.BaseCurrency) + "/" + info.QuoteCurrency,
				Type:      type_,
				Amount:    getPrecision(info.LotSize),
				Price:     getPrecision(info.TickSize),
				AmountMin: float64(info.LotSize),
			})
		}
	} else {
		// 用单独查询的
		for i := range symbols {
			symbolName := getSymbolName(symbols[i])
			info, err := cli.api.GetFutureSpecifySymbol(symbolName)
			if err != nil {
				return nil, err
			}
			type_, _ := ParseDate(ParseTime(info.Data.ExpireDate), info.Data.BaseCurrency)
			res.PrecisionList = append(res.PrecisionList, &client.PrecisionItem{
				Symbol:    XBT2BTC(info.Data.BaseCurrency) + "/" + info.Data.QuoteCurrency,
				Type:      type_,
				Amount:    getPrecision(info.Data.TickSize),
				Price:     getPrecision(info.Data.LotSize),
				AmountMin: float64(info.Data.LotSize),
			})
		}
	}
	return res, nil
}

func getPrecision(x interface{}) int64 {
	s := transform.XToString(x)
	if strings.Contains(s, ".") {
		return int64(len(strings.Split(s, ".")[1]))
	} else {
		return 1 - int64(len(strings.Split(s, ".")[0]))
	}
}
func (cli *ClientKucoin) GetFutureBalance(market common.Market) (*client.UBaseBalance, error) { // 获得合约的balance信息
	res := &client.UBaseBalance{}
	if market == common.Market_FUTURE || market == common.Market_FUTURE_COIN {
		positions, err := cli.api.GetFuturePositions()
		if err != nil {
			return nil, err
		}
		PrintRes(positions)
		for _, v := range positions.Data {
			symbol := strings.Split(v.Symbol, "USD")[0] + v.SettleCurrency
			type_ := common.SymbolType_SWAP_FOREVER
			if v.Symbol[len(v.Symbol)-1] == 'M' {
				if strings.Contains(v.Symbol, "USDT") {
					type_ = common.SymbolType_SWAP_COIN_FOREVER
				}
			} else {
				type_ = common.SymbolType_FUTURE_THIS_QUARTER
			}
			side := order.TradeSide_BUY
			if v.AvgEntryPrice < v.LiquidationPrice {
				side = order.TradeSide_SELL
			}
			res.UBasePositionList = append(res.UBasePositionList, &client.UBasePositionItem{
				Symbol:   symbol,
				Type:     type_,
				Position: float64(v.CurrentQty),
				Unprofit: v.UnrealisedPnl,
				Side:     side,
			})
		}
	} else {
		return nil, errors.New("只能返回FUTURE和FUTURE_COIN")
	}
	return res, nil
}

// todo 杠杆数 不返回orderID
var leverage = 1

func (cli *ClientKucoin) PlaceFutureOrder(req *order.OrderTradeCEX) (*client.OrderRsp, error) { //下单
	var (
		precision = &client.PrecisionItem{}
		ok        = false
		params    = url.Values{}
		side      = "buy"
	)
	resp := &client.OrderRsp{}
	symbolName := getSymbolName(&client.SymbolInfo{Symbol: string(req.Base.Symbol), Type: req.Base.Type})
	amount := req.Amount
	if precision, ok = cli.futurePrecisionM[string(req.Base.Symbol)]; !ok {
		return nil, errors.New("获取精度错误")
	}
	if amount < precision.AmountMin {
		s := fmt.Sprintf("下单数量:%f小于要求最小数量:%f", amount, precision.AmountMin)
		return nil, errors.New(s)
	}
	if req.Side != order.TradeSide_BUY {
		side = "sell"
	}
	orderBase := req.Base
	clientOid := transform.IdToClientId(req.Hdr.Producer, orderBase.Id)
	if req.OrderType == order.OrderType_LIMIT {
		if req.Tif == order.TimeInForce_GTC {
			params.Add("timeInForce", "GTC")
		} else if req.Tif == order.TimeInForce_IOC {
			params.Add("timeInForce", "IOC")
		}
	}
	orderRes, err := cli.api.PlaceFutureOrder(clientOid, side, symbolName, leverage, params)
	PrintRes(orderRes)
	if err != nil {
		return nil, err
	}
	resp.OrderId = orderRes.Data.OrderId
	if orderRes.Code == "200000" {
		resp.Status = order.OrderStatusCode_OPENED
	} else {
		resp.Status = order.OrderStatusCode_FAILED
	}
	resp.Symbol = string(req.Base.Symbol)
	resp.CloseDate = GetDate(req.Base.Type)
	//resp.Producer, resp.Id = transform.ClientIdToId(clientOid)
	return resp, nil
}

func GetDate(f common.SymbolType) string {
	switch f {
	case common.SymbolType_FUTURE_THIS_WEEK:
		return transform.GetDate(transform.THISWEEK)
	case common.SymbolType_FUTURE_NEXT_WEEK:
		return transform.GetDate(transform.NEXTWEEK)
	case common.SymbolType_FUTURE_THIS_MONTH:
		return transform.GetDate(transform.THISMONTH)
	case common.SymbolType_FUTURE_NEXT_MONTH:
		return transform.GetDate(transform.NEXTMONTH)
	case common.SymbolType_FUTURE_THIS_QUARTER:
		return transform.GetDate(transform.THISQUARTER)
	case common.SymbolType_FUTURE_NEXT_QUARTER:
		return transform.GetDate(transform.NEXTQUARTER)
	case common.SymbolType_SWAP_FOREVER:
		return "FOREVER"
	case common.SymbolType_FUTURE_COIN_THIS_WEEK:
		return transform.GetDate(transform.THISWEEK)
	case common.SymbolType_FUTURE_COIN_NEXT_WEEK:
		return transform.GetDate(transform.NEXTWEEK)
	case common.SymbolType_FUTURE_COIN_THIS_MONTH:
		return transform.GetDate(transform.THISMONTH)
	case common.SymbolType_FUTURE_COIN_NEXT_MONTH:
		return transform.GetDate(transform.NEXTMONTH)
	case common.SymbolType_FUTURE_COIN_THIS_QUARTER:
		return transform.GetDate(transform.THISQUARTER)
	case common.SymbolType_FUTURE_COIN_NEXT_QUARTER:
		return transform.GetDate(transform.NEXTQUARTER)
	case common.SymbolType_SWAP_COIN_FOREVER:
		return "FOREVER"
	default:
		return ""
	}
}
