package bybit

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/bybit/spot_api"
	"clients/logger"
	"clients/transform"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
)

type ClientBybit struct {
	api               *spot_api.ApiClient
	tradeFeeMap       map[string]*client.TradeFeeItem    //key: symbol
	transferFeeMap    map[string]*client.TransferFeeItem //key: network+token
	precisionMap      map[string]*client.PrecisionItem   //key: symbol
	symbolMap         map[string]string
	withdrawFeeMap    map[string]chains // 提现手续费
	futureTradeFeeMap map[string]*client.PrecisionItem

	optionMap map[string]interface{}
}

type chains map[string]string

func NewClientBybit(conf base.APIConf, maps ...interface{}) *ClientBybit {
	c := &ClientBybit{
		api: spot_api.NewApiClient(conf),
	}
	c.initialMap()
	for _, m := range maps {
		switch t := m.(type) {
		case map[string]*client.TradeFeeItem:
			c.tradeFeeMap = t

		case map[string]*client.TransferFeeItem:
			c.transferFeeMap = t

		case map[string]*client.PrecisionItem:
			c.precisionMap = t

		case map[string]interface{}:
			c.optionMap = t
		}
	}
	return c
}

func (c *ClientBybit) initialMap() {
	symbols := c.GetSymbols()
	c.symbolMap = make(map[string]string)
	for _, v := range symbols {
		c.symbolMap[trans2Recive(v)] = v
	}
	c.GetWithdrawFee()
	c.initTradeFeeMap()
}

func (c *ClientBybit) initTradeFeeMap() error {
	c.futureTradeFeeMap = make(map[string]*client.PrecisionItem)
	symbolInfo, err := c.api.GetFutureSymbols()
	if err != nil {
		return err
	}
	for _, item := range symbolInfo.Result {
		amount := getPrecision(item.LotSizeFilter.MinTradingQty)
		price := getPrecision(item.PriceFilter.TickSize)
		v := &client.PrecisionItem{
			Symbol:    item.BaseCurrency + "/" + item.QuoteCurrency,
			Amount:    amount,
			Price:     price,
			AmountMin: item.LotSizeFilter.MinTradingQty,
		}
		c.futureTradeFeeMap[v.Symbol] = v
	}
	return nil
}

func getPrecision(x interface{}) int64 {
	s := transform.XToString(x)
	if strings.Contains(s, ".") {
		return int64(len(strings.Split(s, ".")[1]))
	} else {
		return 1 - int64(len(strings.Split(s, ".")[0]))
	}
}

func trans2Recive(symbol string) string {
	return strings.ReplaceAll(symbol, "/", "")
}

func (c *ClientBybit) GetWithdrawFee() error {
	c.withdrawFeeMap = make(map[string]chains)
	res, err := c.api.GetWithdrawFee()
	if err != nil {
		return err
	}
	for _, item := range res.Result.Rows {
		chainsInstance := chains{}
		for _, v := range item.Chains {
			chainsInstance[v.Chain] = v.WithdrawFee
		}
		c.withdrawFeeMap[item.Coin] = chainsInstance
	}
	return nil
}

func (c *ClientBybit) GetExchange() common.Exchange {
	return common.Exchange_BYBIT
}

func (c *ClientBybit) GetSymbols() []string {
	res, err := c.api.GetSymbols()
	var symbols []string
	if err != nil {
		logger.Logger.Error("get symbols ", err.Error())
		return symbols
	}

	for _, symbol := range res.Result {
		if symbol.ShowStatus {
			symbols = append(symbols, symbol.BaseCurrency+"/"+symbol.QuoteCurrency)
		}
	}
	return symbols
}

func (c *ClientBybit) GetDepth(symbol *client.SymbolInfo, limit int) (*depth.Depth, error) {
	sym := spot_api.Trans2Send(symbol.Symbol)
	params := url.Values{}
	params.Add("sz", fmt.Sprint(limit))
	repDepth, err := c.api.GetOrderBook(sym)
	if err != nil {
		//utils.Logger.Error(d.Symbol, d.getSymbolName(), "get full depth err", err)
		return nil, err
	}
	dep := &depth.Depth{
		Exchange:    common.Exchange_OKEX,
		Market:      common.Market_SPOT,
		Type:        common.SymbolType_SPOT_NORMAL,
		Symbol:      strings.ReplaceAll(symbol.Symbol, "-", "/"),
		TimeOperate: uint64(time.Now().UnixMicro()),
	}
	err = c.ParseOrder(repDepth.Result.Bids, &dep.Bids)
	if err != nil {
		return dep, err
	}

	err = c.ParseOrder(repDepth.Result.Asks, &dep.Asks)
	if err != nil {
		return dep, err
	}
	dep.TimeReceive = uint64(repDepth.Time)
	return dep, err
}

func (c *ClientBybit) ParseOrder(orders [][]string, slice *[]*depth.DepthLevel) error {
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

func ParsePriceAmountFloat(data []string) (price float64, amount float64, err error) {
	price, err = strconv.ParseFloat(data[0], 64)
	if err != nil {
		return
	}
	amount, err = strconv.ParseFloat(data[1], 64)
	if err != nil {
		return
	}
	return
}

func (c *ClientBybit) IsExchangeEnable() bool {
	res, err := c.api.GetTime(nil)
	if err != nil {
		logger.Logger.Error(err)
	}
	serverTime, err := strconv.ParseInt(res.Result.TimeNano, 10, 64)
	serverTime = serverTime / 1e6
	return err == nil && time.Since(time.UnixMilli(serverTime)) < time.Second*60
}

func (c *ClientBybit) GetTransferFee(chain common.Chain, tokens ...string) (*client.TransferFee, error) {
	//FIXME: 支持此功能
	//var (
	//	res            = &client.TransferFee{}
	//	searchTokenMap = make(map[string]bool)
	//)
	//for _, token := range tokens {
	//	searchTokenMap[token] = true
	//}
	//configAll, err := c.api.Asset_Currencies_Info(nil)
	//if err != nil {
	//	return nil, err
	//}
	//for _, item := range configAll.Data {
	//	if len(searchTokenMap) > 0 {
	//		_, ok := searchTokenMap[item.Ccy]
	//		if !ok || ok_api.GetChainFromNetWork(item.Chain) != chain {
	//			continue
	//		}
	//	}
	//	var fee float64
	//	//获取接口中的手续费值
	//	fee, err = strconv.ParseFloat(item.MinFee, 64)
	//	if err != nil {
	//		logger.Logger.Info("convert withdraw fee err:", item.Ccy, item.Chain, item.Ccy)
	//		continue
	//	}
	//	res.TransferFeeList = append(res.TransferFeeList, &client.TransferFeeItem{
	//		Token:   item.Ccy,
	//		Network: ok_api.GetChainFromNetWork(item.Chain),
	//		Fee:     fee,
	//	})
	//}
	return nil, nil
}

func (c *ClientBybit) GetPrecision(symbols ...string) (*client.Precision, error) {
	return nil, nil
}

func (c *ClientBybit) GetMarginBalance() (*client.MarginBalance, error) {
	return nil, nil
}

/*
返回手续费按照vip等级返回
https://www.bybit.com/zh-TW/help-center/bybitHC_Article?language=zh_TW&id=000001634
这里返回最高手续费
*/

func (c *ClientBybit) GetTradeFee(symbols ...string) (*client.TradeFee, error) { //查询交易手续费
	var res = &client.TradeFee{}

	makerFee := 0.0001
	takerFee := 0.0001

	for _, symbol := range symbols {
		res.TradeFeeList = append(res.TradeFeeList, &client.TradeFeeItem{
			Symbol: symbol,
			Maker:  makerFee,
			Taker:  takerFee,
		})
	}
	return res, nil
}

func (c *ClientBybit) GetBalance() (*client.SpotBalance, error) {
	var (
		respAccount *spot_api.RespAccountBalance
		res         = &client.SpotBalance{}
		err         error
	)
	respAccount, err = c.api.GetAccountBalance()
	if err != nil {
		return nil, err
	}
	res.UpdateTime = respAccount.Time
	for _, balance := range respAccount.Result.Balances {
		var (
			free, frozen, total float64
		)
		free = transform.StringToX[float64](balance.Free).(float64)
		frozen = transform.StringToX[float64](balance.Locked).(float64)
		total = transform.StringToX[float64](balance.Total).(float64)

		res.BalanceList = append(res.BalanceList, &client.SpotBalanceItem{
			Asset:  balance.Coin,
			Free:   free,
			Frozen: frozen,
			Total:  total,
		})
	}
	return res, nil
}

func (c *ClientBybit) PlaceOrder(o *order.OrderTradeCEX) (*client.OrderRsp, error) {
	var (
		symbol       string
		side         SideType
		type_        OrderType
		options      = url.Values{}
		resp         = &client.OrderRsp{}
		orderLinkId  string
		precisionIns interface{}
		ok           bool
		precision    *client.PrecisionItem
		tdMode       string
		err          error
	)
	symbol = strings.ToUpper(strings.Replace(string(o.Base.Symbol), "/", "", -1))
	if precisionIns, ok = c.precisionMap[string(o.Base.Symbol)]; !ok {
		return nil, errors.New("get precision err")
	}
	if precision, ok = precisionIns.(*client.PrecisionItem); !ok {
		return nil, errors.New("get precision err")
	}
	if o.Amount < precision.AmountMin {
		return nil, errors.New(fmt.Sprint("less amount in err:", o.Amount, "<", precision.AmountMin))
	}
	side = GetSideTypeToExchange(o.Side)

	if o.Base.Id != 0 {
		orderLinkId = transform.IdToClientId(o.Hdr.Producer, o.Base.Id)
	}
	// 交易类型
	//tdMode := "cash"

	if o.Base.Market == common.Market_SPOT {
		tdMode = "cash"
	} else {
		if tdMode == "" {
			tdMode = "cross"
		}
	}

	amount := o.Amount
	// postOnly
	if o.TradeType == order.TradeType_MAKER {
		type_ = ORDER_TYPE_POST_ONLY
		options.Add("orderPrice", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
	} else if o.Tif == order.TimeInForce_FOK {
		type_ = ORDER_TYPE_FOK
		options.Add("orderPrice", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
	} else if o.Tif == order.TimeInForce_IOC {
		type_ = ORDER_TYPE_IOC
		options.Add("orderPrice", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
	} else if o.Tif == order.TimeInForce_OPTIMAL_LIMIT_IOC {
		type_ = ORDER_TYPE_OPTIMAL_LIMIT_IOC
	} else {
		type_ = GetOrderTypeToExchange(o.OrderType)
		if type_ == ORDER_TYPE_LIMIT {
			options.Add("orderPrice", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
		} else if type_ == ORDER_TYPE_MARKET {
			return nil, errors.New("bybit不支持下市价单")
		}
	}
	var (
		res *spot_api.RespPlaceOrder
	)
	//
	if o.Base.Market == common.Market_SPOT {
		res, err = c.api.PostPlaceOrder(symbol, strconv.FormatFloat(amount, 'f', int(precision.Amount), 64), string(side), strings.ToUpper(string(type_)), orderLinkId, options)
	} else if o.Base.Market == common.Market_MARGIN {
		// todo 目前只有现货
		//res, err = c.api.PostPlaceOrder(symbol, tdMode, string(side), string(type_), strconv.FormatFloat(amount, 'f', int(precision.Amount), 64), options)
	}
	if err != nil {
		return nil, err
	}
	resp.Producer, resp.Id = transform.ClientIdToId(res.Result.OrderLinkId)
	resp.OrderId = res.Result.OrderId
	resp.Symbol = string(o.Base.Symbol)
	resp.Status = GetOrderStatusFromExchange(res.Result.Status)
	return resp, nil
}

func (c *ClientBybit) CancelOrder(o *order.OrderCancelCEX) (*client.OrderRsp, error) {
	var (
		resp    = &client.OrderRsp{}
		err     error
		symbol  client.SymbolInfo
		orderId = string(o.Base.IdEx)
		params  = url.Values{}
		res     *spot_api.RespCancleOrder
	)
	// 只有ordeid,其他的是什么意思
	if o.Base.Id != 0 {
		params.Add("orderLinkId", transform.IdToClientId(o.Hdr.Producer, o.Base.Id))
	} else if orderId == "" {
		return nil, errors.New("id can not be empty")
	}
	symbol = client.SymbolInfo{
		Symbol: string(o.Base.Symbol),
		Type:   o.Base.Type,
	}
	res, err = c.api.PostCancleOrder(orderId, params)
	if err != nil {
		return nil, err
	}
	resp.OrderId = res.Result.OrderId
	resp.Producer, resp.Id = transform.ClientIdToId(res.Result.OrderLinkId)
	resp.RespType = client.OrderRspType_RESULT
	// symbol就是请求的symbol
	resp.Symbol = symbol.Symbol
	// 返回没有以下的
	resp.Status = GetOrderStatusFromExchange(res.Result.Status)
	return resp, nil
}

func (c *ClientBybit) GetOrder(req *order.OrderQueryReq) (*client.OrderRsp, error) {
	var (
		symbol  client.SymbolInfo
		resp    = &client.OrderRsp{}
		orderId = req.IdEx //orderId	LONG	NO
		params  = url.Values{}
	)
	symbol = client.SymbolInfo{
		Symbol: string(req.Symbol),
		Type:   req.Type,
	}

	if len(req.Producer) > 0 && req.Id != 0 {
		params.Add("orderLinkId", transform.IdToClientId(req.Producer, req.Id))
	} else if orderId == "" {
		return nil, errors.New("传入orderId和自定义Id不能同时为空")
	}

	res, err := c.api.GetOrderInfo(orderId, params)
	if err != nil {
		return nil, err
	}

	// 创建时间
	ts, err := strconv.ParseInt(res.Result.UpdateTime, 10, 64)
	if err != nil {
		return nil, err
	}
	resp.Producer, resp.Id = transform.ClientIdToId(res.Result.OrderLinkId)
	resp.OrderId = res.Result.OrderId
	resp.Timestamp = ts

	//resp. = ParseInstType(res.Data[0].InstType)
	resp.RespType = client.OrderRspType_RESULT
	resp.Symbol = symbol.Symbol
	resp.Status = GetOrderStatusFromExchange(res.Result.Status)
	resp.Price, _ = strconv.ParseFloat(res.Result.StopPrice, 64)
	resp.Executed, _ = strconv.ParseFloat(res.Result.ExecQty, 64)

	resp.AvgPrice, _ = strconv.ParseFloat(res.Result.AvgPrice, 64)
	resp.AccumAmount, _ = strconv.ParseFloat(res.Result.ExecQty, 64)
	resp.AccumQty = resp.AvgPrice * resp.AccumAmount
	return resp, nil
}

func (c *ClientBybit) GetOrderHistory(req *client.OrderHistoryReq) ([]*client.OrderRsp, error) {
	var (
		resp      []*client.OrderRsp
		startTime = req.StartTime //startTime	LONG	NO
		endTime   = req.EndTime   //endTime	LONG	NO
		limit     = 100           //limit	INT	NO	默认 100; 最大 100.
		params    = url.Values{}
	)

	params.Add("symbol", Trans2Send(req.Asset))
	if endTime > 0 {
		params.Add("endTime", strconv.FormatInt(endTime, 10))
	}
	if startTime > 0 {
		params.Add("startTime", strconv.FormatInt(startTime, 10))
	}
	params.Add("limit", strconv.Itoa(limit))

	// 近三个月
	res, err := c.api.GetOrderHistory(params)
	if err != nil {
		return nil, err
	}
	for _, item := range res.Result.List {
		producer, id := transform.ClientIdToId(item.OrderLinkId)
		price, _ := strconv.ParseFloat(item.StopPrice, 64)
		symbol := strings.ReplaceAll(item.Symbol, "-", "")
		executed, _ := strconv.ParseFloat(item.ExecQty, 64)
		avgPrice, _ := strconv.ParseFloat(item.AvgPrice, 64)
		accumAmount, _ := strconv.ParseFloat(item.ExecQty, 64)
		ts := item.UpdateTime
		resp = append(resp, &client.OrderRsp{
			Producer:    producer,
			Id:          id,
			OrderId:     item.OrderId,
			Timestamp:   ts,
			Symbol:      c.symbolMap[symbol],
			RespType:    client.OrderRspType_RESULT,
			Status:      GetOrderStatusFromExchange(item.Status),
			Price:       price,
			Executed:    executed,
			AvgPrice:    avgPrice,
			AccumAmount: accumAmount,
			AccumQty:    avgPrice * accumAmount,
		})
	}
	return resp, nil
}

func Trans2Send(symbol string) string {
	return strings.ReplaceAll(symbol, "/", "")

}

func (c *ClientBybit) GetProcessingOrders(req *client.OrderHistoryReq) ([]*client.OrderRsp, error) {
	var (
		resp   []*client.OrderRsp
		limit  = 100 //limit	INT	NO	默认 100; 最大 100.
		params = url.Values{}
	)
	params.Add("symbol", Trans2Send(req.Asset))
	params.Add("limit", strconv.Itoa(limit))

	// 近三个月
	res, err := c.api.GetOrderList(params)
	if err != nil {
		return nil, err
	}
	for _, item := range res.Result.List {
		status := GetOrderStatusFromExchange(item.Status)
		producer, id := transform.ClientIdToId(item.OrderLinkId)
		price, _ := strconv.ParseFloat(item.StopPrice, 64)
		symbol := strings.ReplaceAll(item.Symbol, "-", "")
		executed, _ := strconv.ParseFloat(item.ExecQty, 64)
		avgPrice, _ := strconv.ParseFloat(item.AvgPrice, 64)
		accumAmount, _ := strconv.ParseFloat(item.ExecQty, 64)
		ts := item.UpdateTime
		resp = append(resp, &client.OrderRsp{
			Producer:    producer,
			Id:          id,
			OrderId:     item.OrderId,
			Timestamp:   ts,
			Symbol:      c.symbolMap[symbol],
			RespType:    client.OrderRspType_RESULT,
			Status:      status,
			Price:       price,
			Executed:    executed,
			AvgPrice:    avgPrice,
			AccumAmount: accumAmount,
			AccumQty:    avgPrice * accumAmount,
		})
	}
	return resp, nil
}

func (c *ClientBybit) Transfer(o *order.OrderTransfer) (*client.OrderRsp, error) {
	//默认都是从现货提币
	var (
		resp = &client.OrderRsp{}
		err  error
		coin = string(o.ExchangeToken)
		//withdrawOrderId = o.Base.Id //自定义提币ID
		// todo 需要指定chain address
		chain   = common.Chain_name[int32(o.Chain)]
		address = string(o.TransferAddress)
		amount  = o.Amount //数量
		params  = url.Values{}
		// todo 手续费未实现
		fee string
	)

	params.Add("symbol", Trans2Send(coin))
	params.Add("chain", chain)
	fee = c.withdrawFeeMap[Trans2Send(coin)][chain]
	amount_ := transform.XToString(amount + transform.StringToX[float64](fee).(float64))
	params.Add("amount", amount_)
	params.Add("address", address)
	params.Add("timestamp", transform.XToString(time.Now().UnixMicro()))
	res, err := c.api.Withdraw(params)
	if err != nil {
		return nil, err
	}

	resp.OrderId = res.Result.Id
	resp.RespType = client.OrderRspType_ACK
	return resp, err
}

func (c *ClientBybit) GetTransferHistory(req *client.TransferHistoryReq) (*client.TransferHistoryRsp, error) {
	var (
		resp            = &client.TransferHistoryRsp{}
		err             error
		coin            = req.Asset //coin	STRING	NO
		withdrawOrderId string      //withdrawOrderId	STRING	NO
		//offset          = 0             //offset	INT	NO
		limit     = 1000          //limit	INT	NO	默认：1000， 最大：1000
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
	params.Add("limit", strconv.Itoa(limit))
	res, err = c.api.WithdrawHistory(params)
	if err != nil {
		return nil, err
	}
	for _, item := range res.Result.Rows {
		fee, _ := strconv.ParseFloat(item.WithdrawFee, 64)
		timestamp, _ := strconv.ParseInt(item.UpdateTime, 10, 64)
		amount, _ := strconv.ParseFloat(item.Amount, 64)
		if err != err {
			return nil, err
		}

		resp.TransferList = append(resp.TransferList, &client.TransferHistoryItem{
			Asset:   item.Coin,
			Amount:  amount,
			OrderId: item.TxId,
			Network: common.Chain(common.Chain_value[item.Chain]),
			Status:  parseStatus(item.Status),
			Fee:     fee,
			// 没有确认数
			Address:   item.ToAddress1,
			TxId:      item.TxId,
			Timestamp: timestamp,
		})
	}
	//offset += limit

	return resp, err
}

func parseStatus(status string) client.TransferStatus {
	switch status {
	case "SecurityCheck":
		return client.TransferStatus_TRANSFERSTATUS_WAITCERTIFICATE
	case "Pending":
		return client.TransferStatus_TRANSFERSTATUS_PROCESSING
	case "success":
		return client.TransferStatus_TRANSFERSTATUS_COMPLETE
	case "CancelByUser":
		return client.TransferStatus_TRANSFERSTATUS_CANCELLED
	case "Reject":
		return client.TransferStatus_TRANSFERSTATUS_REJECTED
	case "Fail":
		return client.TransferStatus_TRANSFERSTATUS_FAILED
	case "BlockchainConfirmed":
		return client.TransferStatus_TRANSFERSTATUS_COMPLETE
	default:
		return client.TransferStatus_TRANSFERSTATUS_INVALID
	}
}

func (c *ClientBybit) MoveAsset(o *order.OrderMove) (*client.OrderRsp, error) {
	var (
		resp                           = &client.OrderRsp{}
		err                            error
		asset                          = strings.ReplaceAll(o.Asset, "/", "-") // asset	STRING	YES
		amount                         = o.Amount                              // amount	DECIMAL	YES
		from_account_type, transfer_id string
		to_account_type                string
		params                         = url.Values{}
		type_                          = "0"
	)
	// 母账户发起
	if o.ActionUser == order.OrderMoveUserType_Master {
		if o.AccountSource != "" {
			if o.AccountTarget == "" {
				// 子母划转
				type_ = "2"
			}
		} else if o.AccountSource == "" {
			if o.AccountTarget != "" {
				// 子子划转
				type_ = "1"
			} else {
				// 母子划转
				return nil, errors.New("划转参数出错，AccountTarget不应该为空")
			}
		} else {
			return nil, errors.New("划转参数出错，母账号不支持子转子")
		}
	} else if o.ActionUser == order.OrderMoveUserType_Sub {
		// 子账户发起
		if o.AccountTarget == "" {
			// 子母划转
			type_ = "3"
		} else {
			// 子子划转
			type_ = "4"
		}
	}
	params.Add("type", type_)
	if o.Source == common.Market_SPOT {
		from_account_type = "SPOT"
	} else if o.Source == common.Market_WALLET {
		from_account_type = "INVESTMENT"
	} else {
		return nil, errors.New("划转只能用资金，交易账户")
	}
	if o.Target == common.Market_SPOT {
		to_account_type = "SPOT"
	} else if o.Target == common.Market_WALLET {
		to_account_type = "INVESTMENT"
	} else {
		return nil, errors.New("划转只能用资金，交易账户")
	}
	if type_ == "1" || type_ == "4" {
		res, err := c.api.TransferM2S(transfer_id, asset, fmt.Sprintf("%f", amount), from_account_type, to_account_type, params)
		if err != nil {
			return nil, err
		}
		resp.OrderId = res.Result.TransferId
	} else {
		//res, err := c.api.Transfer(transfer_id, asset, fmt.Sprintf("%f", amount), from_account_type, to_account_type, params)
		res, err := c.api.ALLTransfer(transfer_id, asset, fmt.Sprintf("%f", amount), o.AccountSource, o.AccountTarget, from_account_type, to_account_type, params)
		if err != nil {
			return nil, err
		}
		resp.OrderId = res.Result.TransferId
	}
	resp.RespType = client.OrderRspType_ACK
	return resp, err
}

func (c *ClientBybit) GetMoveHistory(req *client.MoveHistoryReq) (*client.MoveHistoryRsp, error) {
	var (
		resp      = &client.MoveHistoryRsp{}
		err       error
		startTime = req.StartTime //startTime	LONG	NO
		endTime   = req.EndTime   //endTime	LONG	NO
		params    = url.Values{}
		type_     client.MoveType
	)
	if startTime > 0 {
		params.Add("start_time", strconv.FormatInt(startTime, 10))
	}
	if endTime > 0 {
		params.Add("end_time", strconv.FormatInt(endTime, 10))
	}

	if req.ActionUser == order.OrderMoveUserType_Master {
		// 子母账户万能划转历史
		// 子子，子母 转入 从子账户转入
		res, err := c.api.TransferHistoryS2M(params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Result.List {
			amount, _ := strconv.ParseFloat(item.Amount, 64)
			if item.Type == "IN" {
				type_ = client.MoveType_MOVETYPE_IN
			} else {
				type_ = client.MoveType_MOVETYPE_OUT
			}
			resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
				Asset:     item.Coin,
				Amount:    amount,
				Timestamp: transform.StringToX[int64](item.Timestamp).(int64) * 1000,
				Type:      type_,
			})
		}
	} else if req.ActionUser == order.OrderMoveUserType_Sub {
		// 子账户划转历史
		// 无法区分
		return nil, nil
		// 子子
	} else { //OrderMoveUserType_Internal
		// 主账户划转历史
		// 母子 转入划转
		res, err := c.api.TransferHistory(params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Result.List {
			amount, _ := strconv.ParseFloat(item.Amount, 64)
			resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
				Asset:     item.Coin,
				Amount:    amount,
				Timestamp: transform.StringToX[int64](item.Timestamp).(int64) * 1000,
			})
		}
	}
	return resp, err
}

func (c *ClientBybit) GetDepositHistory(req *client.DepositHistoryReq) (*client.DepositHistoryRsp, error) {
	var (
		resp = &client.DepositHistoryRsp{}
		err  error
		//status    = spot_api.GetDepositTypeToExchange(req.Status) //status	INT	NO	0(0:pending,6: credited but cannot withdraw, 1:success)
		startTime = req.StartTime //startTime	LONG	NO	默认当前时间90天前的时间戳
		endTime   = req.EndTime   //endTime	LONG	NO	默认当前时间戳
		offset    = 0             //offset	INT	NO	默认:0
		limit     = 50            //limit	INT	NO	默认：50，最大50
		params    = url.Values{}
		res       *spot_api.RespRecordHistory
	)

	if startTime > 0 {
		params.Add("start_time", strconv.FormatInt(endTime, 10))
	}
	if endTime > 0 {
		params.Add("end_time", strconv.FormatInt(startTime, 10))
	}

	params.Add("limit", strconv.Itoa(limit))

	res, err = c.api.RecordHistory(fmt.Sprint(time.Now().UnixMicro()), params)
	if err != nil {
		return nil, err
	}
	for i := 0; i < 100; i++ {
		for _, item := range res.Result.Rows {

			ts, err := strconv.ParseInt(item.SuccessAt, 10, 64)
			ts *= 1e3
			if err != nil {
				return nil, err
			}
			amount, err := strconv.ParseFloat(item.Amount, 64)
			state := item.Status
			resp.DepositList = append(resp.DepositList, &client.DepositHistoryItem{
				Asset:  item.Coin,
				Amount: amount,
				//FIXME: 增加chain->network转换函数
				//Network:   GetChainFromNetWork(item.Chain),
				Status:    client.DepositStatus(state),
				Address:   item.ToAddress,
				TxId:      item.TxId,
				Timestamp: ts,
			})
		}
		if len(res.Result.Rows) < limit {
			break
		}
		offset += limit
	}
	return resp, err
}

func (c *ClientBybit) Loan(o *order.OrderLoan) (*client.OrderRsp, error) {
	var (
		resp   = &client.OrderRsp{}
		err    error
		asset  = o.Asset  //asset	STRING	YES
		amount = o.Amount //amount	DECIMAL	YES
		params = url.Values{}
	)

	res, err := c.api.PostLoan(string(asset), fmt.Sprintf("%f", amount), params)
	if err != nil {
		return nil, err
	}
	_ = res
	//resp.OrderId = strconv.FormatInt(res.TranId, 10)
	resp.RespType = client.OrderRspType_ACK
	return resp, err
}

func (c *ClientBybit) GetLoanOrders(o *client.LoanHistoryReq) (*client.LoanHistoryRsp, error) {
	var (
		resp      = &client.LoanHistoryRsp{}
		err       error
		asset     = o.Asset     //asset	STRING	YES
		startTime = o.StartTime //startTime	LONG	NO
		endTime   = o.EndTime   //endTime	LONG	NO
		current   = 1           //current	LONG	NO	当前查询页。 开始值 1。 默认:1
		size      = 100         //size	LONG	NO	默认:100 最大:500
		params    = url.Values{}
		res       *spot_api.RespLoanHistory
	)

	if o.Asset != "" {
		params.Add("coin", asset)
	}

	if startTime > 0 {
		params.Add("endTime", strconv.FormatInt(endTime, 10))
	}
	if endTime > 0 {
		params.Add("startTime", strconv.FormatInt(startTime, 10))
	}

	for i := 0; i < 100; i++ {
		params.Add("limit", strconv.Itoa(size))
		res, err = c.api.GetLoanHistory(params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Result.List {
			// 1 为借币
			ts := item.CreatedTime
			principle, err := strconv.ParseFloat(item.LoanBalance, 64)
			if err != nil {
				return nil, err
			}
			resp.LoadList = append(resp.LoadList, &client.LoanHistoryItem{
				Principal: principle,
				Asset:     item.Coin,
				Timestamp: ts,
			})
		}
		if len(res.Result.List) < size {
			break
		}
		current++
	}
	return resp, err
}

func (c *ClientBybit) Return(o *order.OrderReturn) (*client.OrderRsp, error) {
	var (
		resp   = &client.OrderRsp{}
		err    error
		asset  = o.Asset  //asset	STRING	YES
		amount = o.Amount //amount	DECIMAL	YES
		params = url.Values{}
	)

	res, err := c.api.PostRepay(asset, fmt.Sprintf("%f", amount), params)
	if err != nil {
		return nil, err
	}
	_ = res
	resp.RespType = client.OrderRspType_ACK
	return resp, err
}

func (c *ClientBybit) GetReturnOrders(o *client.LoanHistoryReq) (*client.ReturnHistoryRsp, error) {
	var (
		resp      = &client.ReturnHistoryRsp{}
		err       error
		asset     = o.Asset     //asset	STRING	YES
		startTime = o.StartTime //startTime	LONG	NO
		endTime   = o.EndTime   //endTime	LONG	NO
		current   = 1           //current	LONG	NO	当前查询页。 开始值 1。 默认:1
		size      = 100         //size	LONG	NO	默认:100 最大:100
		params    = url.Values{}
		res       *spot_api.RespRepayHistory
	)

	if o.Asset != "" {
		params.Add("ccy", asset)
	}

	if startTime > 0 {
		params.Add("endTime", strconv.FormatInt(startTime, 10))
	}
	if endTime > 0 {
		params.Add("startTime", strconv.FormatInt(endTime, 10))
	}

	for i := 0; i < 100; i++ {
		params.Add("limit", strconv.Itoa(size))
		res, err = c.api.GetRepayHistory(params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Result.List {
			ts := transform.StringToX[int64](item.RepayTime).(int64)
			principle, err := strconv.ParseFloat(item.RepaidAmount, 64)
			if err != nil {
				return nil, err
			}
			resp.ReturnList = append(resp.ReturnList, &client.ReturnHistoryItem{
				Principal: principle,
				Asset:     item.Coin,
				Timestamp: ts,
			})
		}
		if len(res.Result.List) < size {
			break
		}
		current++
	}
	return resp, err
}
func PrintRes(x interface{}) {
	res, _ := json.Marshal(x)
	fmt.Println("结果为:", string(res))
}
func (c *ClientBybit) GetFutureSymbols(market common.Market) []*client.SymbolInfo { //所有交易对
	res := []*client.SymbolInfo{}
	symbolres, err := c.api.GetFutureSymbols()
	if err != nil {
		logger.Logger.Error("get symbols ", err.Error())
		return nil
	}
	usdPost := []common.SymbolType{common.SymbolType_SWAP_COIN_FOREVER, common.SymbolType_FUTURE_COIN_NEXT_QUARTER, common.SymbolType_FUTURE_COIN_THIS_QUARTER}
	usdtPost := []common.SymbolType{common.SymbolType_SWAP_FOREVER, common.SymbolType_FUTURE_NEXT_QUARTER, common.SymbolType_FUTURE_THIS_QUARTER}

	for _, symbol := range symbolres.Result {
		symbol_ := symbol.BaseCurrency + "/" + symbol.QuoteCurrency
		if market == common.Market_FUTURE_COIN && symbol.QuoteCurrency == "USD" {
			for _, v := range usdPost {
				res = append(res, &client.SymbolInfo{
					Symbol: symbol_,
					Name:   symbol_ + "/" + GetDate(v),
					Type:   v,
				})
			}
		} else if market == common.Market_FUTURE && symbol.QuoteCurrency == "USDT" {
			for _, v := range usdtPost {
				res = append(res, &client.SymbolInfo{
					Symbol: symbol_,
					Name:   symbol_ + "/" + GetDate(v),
					Type:   v,
				})
			}
		}
	}
	return res
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
		return "SWAP"
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
		return "SWAP"
	default:
		return ""
	}
}

func (c *ClientBybit) GetFutureDepth(symbolInfo *client.SymbolInfo, limit int) (*depth.Depth, error) { //获取行情
	post := getpostfix(symbolInfo.Type)
	if post == "-1" {
		return nil, errors.New("只支持永续和反向季度交割")
	}

	symbol := symbolInfo.Symbol + post
	symbol = strings.ReplaceAll(symbol, "/", "")
	symbol = symbol
	res, err := c.api.GetFutureOrderbook(symbol)
	if err != nil {
		return nil, err
	}
	dep := &depth.Depth{
		Exchange:    common.Exchange_BYBIT,
		Market:      symbolInfo.Market,
		Type:        symbolInfo.Type,
		Symbol:      symbolInfo.Symbol,
		TimeOperate: uint64(time.Now().UnixMicro()),
		Bids:        make([]*depth.DepthLevel, 0),
		Asks:        make([]*depth.DepthLevel, 0),
	}
	for _, item := range res.Result {
		price := transform.StringToX[float64](item.Price).(float64)
		size := float64(item.Size)
		if item.Side == "Buy" {
			dep.Bids = append(dep.Bids, &depth.DepthLevel{Price: price, Amount: size})
		}
		if item.Side == "Sell" {
			dep.Asks = append(dep.Asks, &depth.DepthLevel{Price: price, Amount: size})
		}
	}
	return dep, nil
}

// bybit只有1，2，3，4季度的usd交割
// 通过判断当前属于的季度数来确定下一个是第几
func getpostfix(symbolType common.SymbolType) string {
	postfix := []string{"H", "M", "U", "Z"}
	quarter := getQuarter()
	year := time.Now().Year() - 2000
	if symbolType == common.SymbolType_FUTURE_COIN_THIS_QUARTER || symbolType == common.SymbolType_FUTURE_THIS_QUARTER {
		return fmt.Sprintf("%s%d", postfix[quarter-1], year)
	}
	if symbolType == common.SymbolType_FUTURE_COIN_NEXT_QUARTER || symbolType == common.SymbolType_FUTURE_NEXT_QUARTER {
		if quarter == 4 {
			year += 1
		}
		return fmt.Sprintf("%s%d", postfix[quarter%4], year)
	}
	if symbolType == common.SymbolType_SWAP_FOREVER || symbolType == common.SymbolType_SWAP_COIN_FOREVER {
		return ""
	}
	return "-1"
}
func getQuarter() int {
	month := int(time.Now().Month())
	if month >= 1 && month <= 3 {
		return 1
	}
	if month >= 4 && month <= 6 {
		return 2
	}
	if month >= 7 && month <= 9 {
		return 3
	}
	return 4
}

func (c *ClientBybit) GetFutureMarkPrice(market common.Market, symbols ...*client.SymbolInfo) (*client.RspMarkPrice, error) { //标记价格
	res := &client.RspMarkPrice{}
	if symbols == nil {
		symbolInfo := c.GetFutureSymbols(common.Market_FUTURE_COIN)
		symbolInfo = append(symbolInfo, c.GetFutureSymbols(common.Market_FUTURE)...)
		for _, item := range symbolInfo {
			post := getpostfix(item.Type)
			name := item.Symbol + post
			name = strings.ReplaceAll(name, "/", "")
			res_, err := c.api.GetMarketPrice(name, "1", fmt.Sprintf("%d", time.Now().UnixMilli()/1000-60))
			if err != nil {
				return nil, err
			}
			if len(res_.Result) > 0 {
				res.Item = append(res.Item, &client.MarkPriceItem{
					Symbol:    item.Symbol,
					Type:      item.Type,
					MarkPrice: res_.Result[0].Close,
				})
			}
		}
	} else {
		for _, item := range symbols {
			post := getpostfix(item.Type)
			name := item.Symbol + post
			name = strings.ReplaceAll(name, "/", "")
			res_, err := c.api.GetMarketPrice(name, "1", fmt.Sprintf("%d", time.Now().UnixMilli()/1000-6000))
			if err != nil {
				return nil, err
			}
			if len(res_.Result) > 0 {
				res.Item = append(res.Item, &client.MarkPriceItem{
					Symbol:    item.Symbol,
					Type:      item.Type,
					MarkPrice: res_.Result[0].Close,
				})
			} else {
				return nil, errors.New("无此合约信息")
			}
		}
	}
	return res, nil
}
func (c *ClientBybit) GetFutureTradeFee(market common.Market, symbols ...*client.SymbolInfo) (*client.TradeFee, error) { //查询交易手续费
	res := &client.TradeFee{}
	if symbols == nil {
		symbolInfo, err := c.api.GetFutureSymbols()
		if err != nil {
			return nil, err
		}
		for _, item := range symbolInfo.Result {
			taker := transform.StringToX[float64](item.TakerFee).(float64)
			maker := transform.StringToX[float64](item.MakerFee).(float64)
			res.TradeFeeList = append(res.TradeFeeList, &client.TradeFeeItem{
				Symbol: item.BaseCurrency + "/" + item.QuoteCurrency,
				Taker:  taker,
				Maker:  maker,
			})
		}
	} else {
		m := make(map[string]*client.TradeFeeItem)
		symbolInfo, err := c.api.GetFutureSymbols()
		if err != nil {
			return nil, err
		}
		for _, item := range symbolInfo.Result {
			taker := transform.StringToX[float64](item.TakerFee).(float64)
			maker := transform.StringToX[float64](item.MakerFee).(float64)
			v := &client.TradeFeeItem{
				Symbol: item.BaseCurrency + "/" + item.QuoteCurrency,
				Taker:  taker,
				Maker:  maker,
			}
			m[v.Symbol] = v
		}

		for _, symbol := range symbols {
			if v, ok := m[symbol.Symbol]; ok {
				res.TradeFeeList = append(res.TradeFeeList, v)
			}
		}
	}
	return res, nil
}
func (c *ClientBybit) GetFuturePrecision(market common.Market, symbols ...*client.SymbolInfo) (*client.Precision, error) { //查询交易对精读信息
	res := &client.Precision{}
	if symbols == nil {
		usdPost := []common.SymbolType{common.SymbolType_SWAP_COIN_FOREVER, common.SymbolType_FUTURE_COIN_NEXT_QUARTER, common.SymbolType_FUTURE_COIN_THIS_QUARTER}
		usdtPost := []common.SymbolType{common.SymbolType_SWAP_FOREVER, common.SymbolType_FUTURE_NEXT_QUARTER, common.SymbolType_FUTURE_THIS_QUARTER}
		for symbol_, value := range c.futureTradeFeeMap {
			if market == common.Market_FUTURE_COIN && !strings.Contains(symbol_, "USDT") {
				for _, v := range usdPost {
					value.Type = v
					res.PrecisionList = append(res.PrecisionList, value)
				}
			} else if market == common.Market_FUTURE && strings.Contains(symbol_, "USDT") {
				for _, v := range usdtPost {
					value.Type = v
					res.PrecisionList = append(res.PrecisionList, value)
				}
			}
		}
	} else {
		for _, symbol := range symbols {
			if isValidSymbol(symbol) {
				item := c.precisionMap[symbol.Symbol]
				item.Type = symbol.Type
				res.PrecisionList = append(res.PrecisionList, item)
			} else {
				return nil, errors.New("bybit只支持永续和季度合约")
			}
		}
	}
	return res, nil
}

func isValidSymbol(symbol ...*client.SymbolInfo) bool {
	for _, v := range symbol {
		type_ := v.Type
		if type_ == common.SymbolType_SWAP_FOREVER || type_ == common.SymbolType_SWAP_COIN_FOREVER ||
			type_ == common.SymbolType_FUTURE_NEXT_QUARTER || type_ == common.SymbolType_FUTURE_COIN_NEXT_QUARTER ||
			type_ == common.SymbolType_FUTURE_THIS_QUARTER || type_ == common.SymbolType_FUTURE_COIN_THIS_QUARTER {
		} else {
			return false
		}
	}
	return true
}

func (c *ClientBybit) GetFutureBalance(market common.Market) (*client.UBaseBalance, error) { //todo 获得合约的balance信息
	res := &client.UBaseBalance{}
	positionInfo, err := c.api.GetPosition()
	if err != nil {
		return nil, err
	}
	for _, item := range positionInfo.Result {
		// usdt只有永续
		side := parseSide(item.Side)
		if strings.Contains(item.Symbol, "USDT") {
			if market == common.Market_FUTURE {
				symbol := strings.Split(item.Symbol, "USDT")[0] + "/" + "USDT"
				res.UBasePositionList = append(res.UBasePositionList, &client.UBasePositionItem{
					Symbol:   symbol,
					Type:     common.SymbolType_SWAP_FOREVER,
					Position: float64(item.Size),
					Side:     side,
					Unprofit: item.UnrealisedPnl,
				})
			}
		} else if strings.Contains(item.Symbol, "USD") {
			type_ := getType(item.Symbol)
			symbol := strings.Split(item.Symbol, "USD")[0] + "/" + "USD"
			res.UBasePositionList = append(res.UBasePositionList, &client.UBasePositionItem{
				Symbol:   symbol,
				Type:     type_,
				Position: float64(item.Size),
				Side:     side,
				Unprofit: item.UnrealisedPnl,
			})
		}
	}

	return res, nil
}
func parseSide(side string) order.TradeSide {
	if side == "buy" {
		return order.TradeSide_BUY
	}
	return order.TradeSide_SELL
}
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

var s map[byte]int = map[byte]int{'H': 1, 'M': 2, 'U': 3, 'Z': 4}

func getType(name string) common.SymbolType {
	if IsNum(name[len(name)-2:]) {
		// 如果季度和当前相等
		if s[name[len(name)-3]] == getQuarter() {
			return common.SymbolType_FUTURE_COIN_THIS_QUARTER
		} else {
			return common.SymbolType_FUTURE_COIN_NEXT_QUARTER
		}
	}
	return common.SymbolType_SWAP_COIN_FOREVER
}
func (c *ClientBybit) PlaceFutureOrder(req *order.OrderTradeCEX) (*client.OrderRsp, error) { //下单
	orderInfo := spot_api.OrderInfo{}
	res := &client.OrderRsp{}

	symbol := strings.ReplaceAll(string(req.Base.Symbol), "/", "")
	v, ok := c.futureTradeFeeMap[string(req.Base.Symbol)]
	if !ok {
		return nil, errors.New("不支持当前币对")
	}
	if !ok || req.Amount < v.AmountMin {
		return nil, errors.New(fmt.Sprintf("需要高于最少数量：%v", v.AmountMin))
	}
	side := parseClientSide(req.Side)
	qty := req.Amount
	var timeInforce string
	timeInforce = parseClientTimeInfore(req.Tif)
	options := url.Values{}
	type_ := "Market"

	orderInfo.Symbol = symbol
	orderInfo.Side = side
	orderInfo.Order_type = type_
	orderInfo.Qty = qty
	orderInfo.TimeInForce = timeInforce
	orderInfo.PositionIdx = 0
	orderInfo.CloseOnTrigger = false
	res.Symbol = string(req.Base.Symbol)
	if req.Base.Type == common.SymbolType_SWAP_FOREVER {
		placeInfo, err := c.api.PlaceFutureOrder(orderInfo, options)
		if err != nil {
			return nil, err
		}
		if placeInfo.RetCode == 30031 {
			return nil, errors.New("余额不足")
		}
		if strings.Contains(placeInfo.Result.OrderStatus, "Created") {
			res.Status = order.OrderStatusCode_OPENED
		}
		res.AccumAmount = float64(placeInfo.Result.CumExecQty)
		res.Status = parseSerberState(placeInfo.Result.OrderStatus)
	} else if req.Base.Type == common.SymbolType_SWAP_COIN_FOREVER {
		placeInfo, err := c.api.PlaceFutureCoinOrder(orderInfo)
		if err != nil {
			return nil, err
		}
		if placeInfo.RetCode == 30031 {
			return nil, errors.New("余额不足")
		}
		if strings.Contains(placeInfo.Result.OrderStatus, "Created") {
			res.Status = order.OrderStatusCode_OPENED
		}
		res.AccumAmount = float64(placeInfo.Result.CumExecQty)
		res.Status = parseSerberState(placeInfo.Result.OrderStatus)
	} else {
		symbolInfo := client.SymbolInfo{}
		symbolInfo.Symbol = string(req.Base.Symbol)
		symbolInfo.Type = req.Base.Type
		if !isValidSymbol(&symbolInfo) {
			return nil, errors.New("只支持usd季度合约")
		}
		symbol += getpostfix(symbolInfo.Type)
		orderInfo.Symbol = symbol
		placeInfo, err := c.api.PlaceFutureCoinDeliveryOrder(orderInfo)
		if err != nil {
			return nil, err
		}
		if placeInfo.RetCode == 30031 {
			return nil, errors.New("余额不足")
		}
		if strings.Contains(placeInfo.Result.OrderStatus, "Created") {
			res.Status = order.OrderStatusCode_OPENED
		}
		res.AccumAmount = float64(placeInfo.Result.CumExecQty)
		res.Status = parseSerberState(placeInfo.Result.OrderStatus)
	}
	return res, nil
}

func parseSerberState(state string) order.OrderStatusCode {
	switch state {
	case "Created":
		return order.OrderStatusCode_CREATED
	case "New":
		return order.OrderStatusCode_OPENEING
	case "Rejected":
		return order.OrderStatusCode_FAILED
	case "PartiallyFilled":
		return order.OrderStatusCode_PARTFILLED
	case "Filled":
		return order.OrderStatusCode_FILLED
	case "Cancelled":
		return order.OrderStatusCode_CANCELED
	case "PendingCancel":
		return order.OrderStatusCode_CANCELING
	default:
		return order.OrderStatusCode_OrderStatusInvalid
	}
}

func parseClientTimeInfore(tif order.TimeInForce) string {
	switch tif {
	case order.TimeInForce_GTC:
		return "GoodTillCancel"
	case order.TimeInForce_IOC:
		return "ImmediateOrCancel"
	case order.TimeInForce_FOK:
		return "FillOrKill"
	default:
		return ""
	}
}

func parseClientSide(side order.TradeSide) string {
	if side == order.TradeSide_BUY {
		return "Buy"
	}
	return "Sell"
}
