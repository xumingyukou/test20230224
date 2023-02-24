package huobi

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/huobi/c_api"
	"clients/exchange/cex/huobi/spot_api"
	"clients/exchange/cex/huobi/u_api"
	"clients/logger"
	"clients/transform"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
)

type ClientHuobi struct {
	spot_api *spot_api.SpotApiClient
	c_api    *c_api.CApiClient
	u_api    *u_api.UApiClient

	tradeFeeMap    map[string]*client.TradeFeeItem    //key: symbol
	transferFeeMap map[string]*client.TransferFeeItem //key: network+token
	precisionMap   map[string]*client.PrecisionItem   //key: symbol
	optionMap      map[string]interface{}
}

func NewClientHuobi(conf base.APIConf, maps ...interface{}) *ClientHuobi {
	c := &ClientHuobi{
		spot_api: spot_api.NewApiClient(conf),
		c_api:    c_api.NewCApiClient(conf),
		u_api:    u_api.NewUApiClient(conf),
	}

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

func (c *ClientHuobi) GetExchange() common.Exchange {
	return common.Exchange_HUOBI
}

func (c *ClientHuobi) GetSymbols() []string {
	var symbols []string
	exchangeInfoRes, err := c.spot_api.ExchangeInfo()
	if err != nil {
		logger.Logger.Error("get exchange in error:", err)
		return symbols
	}
	for _, symbol := range exchangeInfoRes.Data {
		if symbol.State == "online" {
			symbols = append(symbols, strings.ToUpper(fmt.Sprint(symbol.BaseCurrency, "/", symbol.QuoteCurrency)))
		}
	}
	return symbols
}

func (c *ClientHuobi) ParseOrder(orders [][]float64, slice *[]*depth.DepthLevel) error {
	for _, item := range orders {
		price, amount := item[0], item[1]
		*slice = append(*slice, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	return nil
}

func (c *ClientHuobi) GetDepth(symbol *client.SymbolInfo, limit int) (*depth.Depth, error) {
	sym := strings.Replace(symbol.Symbol, "/", "", 1)
	sym = strings.ToLower(sym)

	repDepth, err := c.spot_api.GetDepth(sym, limit)

	if err != nil {
		//utils.Logger.Error(d.Symbol, d.getSymbolName(), "get full depth err", err)
		return nil, err
	}
	dep := &depth.Depth{
		Exchange:     common.Exchange_HUOBI,
		Market:       common.Market_SPOT,
		Type:         common.SymbolType_SPOT_NORMAL,
		Symbol:       sym,
		TimeExchange: uint64(time.Now().UnixMicro()),
		TimeReceive:  uint64(time.Now().UnixMicro()),
	}
	err = c.ParseOrder(repDepth.Tick.Bids, &dep.Bids)
	if err != nil {
		return dep, err
	}
	err = c.ParseOrder(repDepth.Tick.Asks, &dep.Asks)
	if err != nil {
		return dep, err
	}
	return dep, err
}

func (c *ClientHuobi) IsExchangeEnable() bool {
	serverTime1, err := c.spot_api.ServerTime()
	res1 := err == nil && time.Since(time.UnixMilli(serverTime1.Data)) < time.Second*60
	return res1
}

func (c *ClientHuobi) GetTradeFee(symbols ...string) (*client.TradeFee, error) {
	var (
		res             = &client.TradeFee{}
		searchSymbolMap = make(map[string]bool)
	)
	var symbolsRef string
	for _, symbol := range symbols {
		searchSymbolMap[spot_api.GetSymbolName(symbol)] = true
		if symbolsRef == "" {
			symbolsRef += symbol
		} else {
			symbolsRef += "," + symbol
		}
	}
	tradeFeeRes, err := c.spot_api.AssetTradeFee(symbolsRef)
	if err != nil {
		return nil, err
	}
	for _, tradeFee := range tradeFeeRes.TradeFeeItem {
		fmt.Println(tradeFee)
		var (
			takerFee, makerFee float64
		)
		if len(searchSymbolMap) > 0 {
			_, ok := searchSymbolMap[tradeFee.Symbol]
			if !ok {
				continue
			}
		}
		takerFee, err = strconv.ParseFloat(tradeFee.ActualTakerRate, 64)
		if err != nil {
			fmt.Println("convert taker fee err:", tradeFee.Symbol, tradeFee.ActualTakerRate)
			continue
		}
		makerFee, err = strconv.ParseFloat(tradeFee.MakerFeeRate, 64)
		if err != nil {
			fmt.Println("convert maker fee err:", tradeFee.Symbol, tradeFee.MakerFeeRate)
			continue
		}
		res.TradeFeeList = append(res.TradeFeeList, &client.TradeFeeItem{
			Symbol: tradeFee.Symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
			Maker:  makerFee,
			Taker:  takerFee,
		})
	}
	return res, nil
}

// 提币手续费，参数传递应该得改
func (c *ClientHuobi) GetTransferFee(chain common.Chain, tokens ...string) (*client.TransferFee, error) {
	var res = &client.TransferFee{}

	for _, currency := range tokens {
		options := url.Values{}
		options.Add("currency", strings.ToLower(currency))

		ReferenceCurrencies, err := c.spot_api.GetReferenceCurrencies(&options)
		if err != nil {
			fmt.Println("c.api.GetReferenceCurrencies error =", err)
			return nil, err
		}

		if len(ReferenceCurrencies.Data) == 0 {
			continue
		}

		for _, dataItem := range ReferenceCurrencies.Data[0].Chains {
			if dataItem.BaseChain == common.Chain_name[int32(chain)] || dataItem.DisplayName == common.Chain_name[int32(chain)] {
				var transferFee float64
				if dataItem.WithdrawFeeType == "ratio" {
					transferFeeRate, _ := strconv.ParseFloat(dataItem.TransactFeeRateWithdraw, 32)
					res.TransferFeeList = append(res.TransferFeeList, &client.TransferFeeItem{
						Token:   currency,
						Network: chain,
						FeeRate: transferFeeRate,
					})
				} else {
					transferFee, _ = strconv.ParseFloat(dataItem.TransactFeeWithdraw, 32)
					res.TransferFeeList = append(res.TransferFeeList, &client.TransferFeeItem{
						Token:   currency,
						Network: chain,
						Fee:     transferFee,
					})
				}
			}
		}
	}

	return res, nil
}

func (c *ClientHuobi) GetPrecision(symbols ...string) (*client.Precision, error) {
	var (
		res             = &client.Precision{}
		searchSymbolMap = make(map[string]bool)
	)
	for _, symbol := range symbols {
		searchSymbolMap[spot_api.GetSymbolName(symbol)] = true
	}
	exchangeInfoRes, err := c.spot_api.ExchangeInfo()
	if err != nil {
		return nil, err
	}
	fmt.Println(exchangeInfoRes.Status)
	for _, symbol := range exchangeInfoRes.Data {
		if len(searchSymbolMap) > 0 {
			_, ok := searchSymbolMap[symbol.Symbol]
			if !ok {
				continue
			}
		}
		//fmt.Println(symbol.Symbol)
		res.PrecisionList = append(res.PrecisionList, &client.PrecisionItem{
			Symbol:    symbol.Symbol,
			Amount:    int64(symbol.AmountPrecision),
			Price:     int64(symbol.PricePrecision),
			AmountMin: symbol.LimitOrderMinOrderAmt,
		})
		for _, i := range res.PrecisionList {
			fmt.Println(i.Symbol)
		}
	}
	return res, nil
}

func (c *ClientHuobi) GetBalance() (*client.SpotBalance, error) {
	var (
		respAccountBalance *spot_api.RespAccountBalance
		res                = &client.SpotBalance{}
		err                error
	)
	respAccountBalance, err = c.spot_api.AccountBalance("spot")
	if err != nil {
		return res, err
	}
	if respAccountBalance.Data.Type != "spot" {
		return nil, nil
	}
	type BalanceType struct {
		Free   float64
		Frozen float64
		Debt   float64
	}
	var balanceItems = make(map[string]*BalanceType)

	for _, balance := range respAccountBalance.Data.List {
		_, ok := balanceItems[balance.Currency]
		if !ok {
			balanceItems[balance.Currency] = &BalanceType{}
		}
		if balance.Type == "trade" {
			balanceItems[balance.Currency].Free, err = strconv.ParseFloat(balance.Balance, 64)
			if err != nil {
				return res, err
			}
		} else if balance.Type == "frozen" {
			balanceItems[balance.Currency].Frozen, err = strconv.ParseFloat(balance.Balance, 64)
			if err != nil {
				return res, err
			}
		}
		balanceItems[balance.Currency].Debt, _ = strconv.ParseFloat(balance.Balance, 64)
	}
	for coin, balance := range balanceItems {
		res.BalanceList = append(res.BalanceList, &client.SpotBalanceItem{
			Asset:  strings.ToUpper(coin),
			Free:   balance.Free,
			Frozen: balance.Frozen,
			Total:  balance.Free + balance.Frozen,
		})
	}
	return res, nil
}

func (c *ClientHuobi) GetMarginBalance() (*client.MarginBalance, error) {
	var (
		respAccountBalance *spot_api.RespAccountBalance
		res                = &client.MarginBalance{}
		err                error
	)
	respAccountBalance, err = c.spot_api.AccountBalance("super-margin")
	if err != nil {
		return res, err
	}
	if respAccountBalance.Data.Type != "super-margin" {
		return nil, nil
	}
	res.MarginLevel = 0
	res.MarginLevel = 0
	res.TotalAsset = 0
	res.TotalNetAsset = 0
	res.TotalLiabilityAsset = 0
	res.QuoteAsset = "BTC"

	type BalanceType struct {
		Free   float64
		Frozen float64
		Debt   float64
	}
	var balanceItems = make(map[string]*BalanceType)

	for _, balance := range respAccountBalance.Data.List {
		_, ok := balanceItems[balance.Currency]
		if !ok {
			balanceItems[balance.Currency] = &BalanceType{}
		}
		if balance.Type == "trade" {
			balanceItems[balance.Currency].Free, err = strconv.ParseFloat(balance.Balance, 64)
			if err != nil {
				return res, err
			}
		} else if balance.Type == "frozen" {
			balanceItems[balance.Currency].Frozen, err = strconv.ParseFloat(balance.Balance, 64)
			if err != nil {
				return res, err
			}
		}
		balanceItems[balance.Currency].Debt, _ = strconv.ParseFloat(balance.Balance, 64)
	}
	for coin, balance := range balanceItems {
		res.MarginBalanceList = append(res.MarginBalanceList, &client.MarginBalanceItem{
			Asset:    coin,
			Total:    balance.Free + balance.Frozen + balance.Debt,
			Borrowed: balance.Debt,
			Free:     balance.Free,
			Frozen:   balance.Frozen,
			NetAsset: balance.Free + balance.Frozen,
		})
	}
	return res, nil
}

func (c *ClientHuobi) GetMarginIsolatedBalance(...string) (*client.MarginIsolatedBalance, error) {
	var (
		respAccountBalance *spot_api.RespAccountBalance
		res                = &client.MarginIsolatedBalance{}
		err                error
	)
	respAccountBalance, err = c.spot_api.AccountBalance("margin")
	if err != nil {
		return res, err
	}
	if respAccountBalance.Data.Type != "margin" {
		return nil, nil
	}
	res.TotalNetAsset = 0
	res.TotalLiabilityAsset = 0
	res.QuoteAsset = "BTC"

	type BalanceType struct {
		Free   float64
		Frozen float64
		Debt   float64
	}
	var balanceItems = make(map[string]*BalanceType)

	for _, balance := range respAccountBalance.Data.List {
		_, ok := balanceItems[balance.Currency]
		if !ok {
			balanceItems[balance.Currency] = &BalanceType{}
		}
		if balance.Type == "trade" {
			balanceItems[balance.Currency].Free, err = strconv.ParseFloat(balance.Balance, 64)
			if err != nil {
				return res, err
			}
		} else if balance.Type == "frozen" {
			balanceItems[balance.Currency].Frozen, err = strconv.ParseFloat(balance.Balance, 64)
			if err != nil {
				return res, err
			}
		}
		balanceItems[balance.Currency].Debt, _ = strconv.ParseFloat(balance.Balance, 64)
	}
	for coin, balance := range balanceItems {
		fmt.Println(coin, balance)
		baseAsset := &client.MarginBalanceItem{Asset: coin, Total: balance.Free + balance.Frozen + balance.Debt, Borrowed: balance.Debt, Free: balance.Free, Frozen: balance.Frozen, NetAsset: balance.Free + balance.Frozen}
		quoteAsset := &client.MarginBalanceItem{Asset: coin, Total: balance.Free + balance.Frozen + balance.Debt, Borrowed: balance.Debt, Free: balance.Free, Frozen: balance.Frozen, NetAsset: balance.Free + balance.Frozen}
		res.MarginIsolatedBalanceList = append(res.MarginIsolatedBalanceList, &client.MarginIsolatedBalanceItem{
			BaseAsset:  baseAsset,
			QuoteAsset: quoteAsset,
		})
	}
	return res, nil
}

// 下订单
func (c *ClientHuobi) PlaceOrder(o *order.OrderTradeCEX) (*client.OrderRsp, error) {
	res := &client.OrderRsp{
		Id:     o.Base.Id,
		Symbol: string(o.Base.Symbol),
		Price:  o.Price,
	}

	post_params := spot_api.ReqPostOrder{
		AccountID:     spot_api.SPOT_ACCOUNT_ID,
		Amount:        transform.XToString(o.Amount),
		Type:          spot_api.GetTypeFromTradeOrderSideTif(o.OrderType, o.Side, o.TradeType, o.Tif),
		Price:         transform.XToString(o.Price),
		Source:        "spot-api",
		Symbol:        strings.ToLower(strings.Replace(string(o.Base.Symbol), "/", "", 1)),
		ClientOrderID: transform.XToString(o.Base.Id),
	}

	order_res, err := c.spot_api.PostOrder(post_params)
	if err != nil {
		fmt.Println("c.api.PostOrder error =", err)
		return res, err
	}
	if order_res.Status != "ok" {
		fmt.Println("order err:", order_res.ErrorCode, order_res.ErrorMessage)
		return res, err
	}
	orderId := order_res.Data
	res.OrderId = orderId

	return res, nil
}

// 取消订单
func (c *ClientHuobi) CancelOrder(o *order.OrderCancelCEX) (*client.OrderRsp, error) {
	res := &client.OrderRsp{
		Id: o.CancelId,
	}

	switch o.Base.Market {
	case common.Market_SPOT:
		{
			order_id := o.Base.IdEx
			res.OrderId = order_id
			res.Symbol = string(o.Base.Symbol)

			// 撤销订单
			post_params := spot_api.ReqPostSubmitOrder{
				Symbol: string(o.Base.Symbol),
			}
			resSubmitOrder, err1 := c.spot_api.PostSubmitOrder(order_id, post_params)
			if err1 != nil {
				fmt.Println("c.api.PostSubmitOrder =", err1)
				return res, err1
			}
			if resSubmitOrder.Status == "error" {
				fmt.Println("cancel error =", resSubmitOrder.ErrCode, resSubmitOrder.ErrMsg)
			}
		}
	case common.Market_FUTURE_COIN:
		{
			post_params := c_api.ReqPostFutureCancelOrder{
				OrderID:       o.Base.IdEx,
				ClientOrderId: strconv.FormatInt(o.CancelId, 10),
				Symbol:        string(o.Base.Symbol)[:3],
			}
			ResCancelOrder, err := c.c_api.PostFutureCancelOrder(post_params)
			if err != nil {
				fmt.Println("c.c_api.PostFutureCancelOrder =", err)
				return nil, err
			}
			res.Id, _ = strconv.ParseInt(ResCancelOrder.Data.Successes, 10, 64)
			res.Timestamp = ResCancelOrder.Ts
		}
	case common.Market_SWAP_COIN:
		{
			post_params := c_api.ReqPostSwapCancelOrder{
				OrderId:       o.Base.IdEx,
				ClientOrderId: strconv.FormatInt(o.CancelId, 10),
				ContractCode:  strings.Replace(string(o.Base.Symbol), "/", "-", 1),
			}
			ResCancelOrder, err := c.c_api.PostSwapCancelOrder(post_params)
			if err != nil {
				fmt.Println("c.c_api.PostFutureCancelOrder =", err)
				return nil, err
			}
			res.Id, _ = strconv.ParseInt(ResCancelOrder.Data.Successes, 10, 64)
			res.Timestamp = ResCancelOrder.Ts
		}
	case common.Market_FUTURE, common.Market_SWAP:
		{
			post_params := u_api.ReqPostSwapCrossCancelOrder{
				OrderID:       o.Base.IdEx,
				ClientOrderID: strconv.FormatInt(o.CancelId, 10),
				Pair:          strings.Replace(string(o.Base.Symbol), "/", "-", 1),
				ContractType:  c_api.GetContractTypeFromType(o.Base.Type),
			}
			ResCancelOrder, err := c.u_api.PostSwapCrossCancelOrder(post_params)
			if err != nil {
				fmt.Println("c.c_api.PostFutureCancelOrder =", err)
				return nil, err
			}
			res.Id, _ = strconv.ParseInt(ResCancelOrder.Data.Successes, 10, 64)
			res.Timestamp = ResCancelOrder.Ts
		}
	}

	return res, nil
}

// 查询订单信息
func (c *ClientHuobi) GetOrder(req *order.OrderQueryReq) (*client.OrderRsp, error) {
	res := &client.OrderRsp{
		Producer: req.Producer,
		OrderId:  req.IdEx,
		Symbol:   string(req.Symbol),
	}

	switch req.Market {
	case common.Market_SPOT:
		{
			matchresults, err1 := c.spot_api.GetOrderMatchresults(req.IdEx)
			if err1 != nil {
				fmt.Println("c.api.GetOrderMatchresults =", err1)
				return res, err1
			}
			res.FeeAsset = matchresults.Data[0].FeeCurrency
			newestTime := matchresults.Data[0].CreatedAt
			for _, data := range matchresults.Data {
				if data.CreatedAt >= newestTime {
					res.Price, _ = strconv.ParseFloat(data.Price, 64)
					res.Executed, _ = strconv.ParseFloat(data.FilledAmount, 64)
					newestTime = data.CreatedAt
				}
				nowAmount, _ := strconv.ParseFloat(data.FilledAmount, 64)
				nowPrice, _ := strconv.ParseFloat(data.Price, 64)
				res.AccumAmount += nowAmount
				res.AccumQty += nowAmount * nowPrice
				nowFee, _ := strconv.ParseFloat(data.FilledFees, 64)
				res.Fee += nowFee
				res.Fills = append(res.Fills, &client.FillItem{
					Price:           nowPrice,
					Qty:             nowAmount,
					Commission:      nowPrice * nowAmount,
					CommissionAsset: data.FeeCurrency,
				})
			}
			res.AvgPrice = res.AccumQty / res.AccumAmount
		}
	case common.Market_FUTURE_COIN:
		{
			// 币本位交割
			post_params := c_api.ReqPostFutureContractOrderInfo{
				OrderID:       req.IdEx,
				ClientOrderId: strconv.FormatInt(req.Id, 10),
				Symbol:        string(req.Symbol)[:3],
			}

			ResContractOrderInfo, err := c.c_api.PostFutureContractOrderInfo(post_params)
			if err != nil {
				fmt.Println("c.c_api.PostFutureContractOrderInfo err =", err)
				return nil, err
			}
			res.FeeAsset = ResContractOrderInfo.Data[0].FeeAsset
			newestTime := ResContractOrderInfo.Data[0].CreatedAt
			for _, dataOrderInfo := range ResContractOrderInfo.Data {
				if dataOrderInfo.CreatedAt >= int64(newestTime) {
					newestTime = dataOrderInfo.CreatedAt
					res.Price = dataOrderInfo.Price
					res.Executed = float64(dataOrderInfo.TradeVolume)
				}
				nowAmount := float64(dataOrderInfo.TradeVolume)
				nowPrice := dataOrderInfo.Price
				res.AccumAmount += nowAmount
				res.AccumQty += nowAmount * nowPrice
				nowFee := dataOrderInfo.Fee
				res.Fee += nowFee
				res.Fills = append(res.Fills, &client.FillItem{
					Price:           nowPrice,
					Qty:             nowAmount,
					Commission:      nowAmount * nowPrice,
					CommissionAsset: dataOrderInfo.FeeAsset,
				})
			}
		}
	case common.Market_SWAP_COIN:
		{
			// 币本位永续
			post_params := c_api.ReqPostSwapContractOrderInfo{
				OrderID:       req.IdEx,
				ClientOrderId: strconv.FormatInt(req.Id, 10),
				ContractCode:  strings.Replace(string(req.Symbol), "/", "-", 1),
			}

			ResContractOrderInfo, err := c.c_api.PostSwapContractOrderInfo(post_params)
			if err != nil {
				fmt.Println("c.c_api.PostSwapContractOrderInfo =", err)
				return nil, err
			}
			res.FeeAsset = ResContractOrderInfo.Data[0].FeeAsset
			newestTime := ResContractOrderInfo.Data[0].CreatedAt
			for _, dataOrderInfo := range ResContractOrderInfo.Data {
				if dataOrderInfo.CanceledAt >= int(newestTime) {
					newestTime = dataOrderInfo.CreatedAt
					res.Price = dataOrderInfo.Price
					res.Executed = float64(dataOrderInfo.TradeVolume)
				}
				nowAmount := float64(dataOrderInfo.TradeVolume)
				nowPrice := dataOrderInfo.Price
				res.AccumAmount += nowAmount
				res.AccumQty += nowAmount * nowPrice
				nowFee := dataOrderInfo.Fee
				res.Fee += float64(nowFee)
				res.Fills = append(res.Fills, &client.FillItem{
					Price:           nowPrice,
					Qty:             nowAmount,
					Commission:      nowAmount * nowPrice,
					CommissionAsset: dataOrderInfo.FeeAsset,
				})
			}
		}
	case common.Market_FUTURE, common.Market_SWAP:
		{
			// U本位
			post_params := u_api.ReqPostSwapCrossOrderInfo{
				OrderID:       req.IdEx,
				ClientOrderID: strconv.FormatInt(req.Id, 10),
				ContractCode:  u_api.GetContractCodeFromSymbolAndTypeU(string(req.Symbol), req.Type),
				Pair:          strings.Replace(string(req.Symbol), "/", "-", 1),
			}
			ResContractOrderInfo, err := c.u_api.PostSwapCrossOrderInfo(post_params)
			if err != nil {
				fmt.Println("c.c_api.PostSwapContractOrderInfo =", err)
				return nil, err
			}
			res.FeeAsset = ResContractOrderInfo.Data[0].FeeAsset
			newestTime := ResContractOrderInfo.Data[0].CreatedAt
			for _, dataOrderInfo := range ResContractOrderInfo.Data {
				if dataOrderInfo.CanceledAt >= int(newestTime) {
					newestTime = dataOrderInfo.CreatedAt
					res.Price = dataOrderInfo.Price
					res.Executed = float64(dataOrderInfo.TradeVolume)
				}
				nowAmount := float64(dataOrderInfo.TradeVolume)
				nowPrice := dataOrderInfo.Price
				res.AccumAmount += nowAmount
				res.AccumQty += nowAmount * nowPrice
				nowFee := dataOrderInfo.Fee
				res.Fee += float64(nowFee)
				res.Fills = append(res.Fills, &client.FillItem{
					Price:           nowPrice,
					Qty:             nowAmount,
					Commission:      nowAmount * nowPrice,
					CommissionAsset: dataOrderInfo.FeeAsset,
				})
			}

		}
	}

	return res, nil
} //查询订单信息。不支持ws回报订单行情的交易所，如bitbank, 需要频发查单获知订单状态

// 查询历史订单信息
func (c *ClientHuobi) GetOrderHistory(req *client.OrderHistoryReq) ([]*client.OrderRsp, error) {
	res := make([]*client.OrderRsp, 1)

	switch req.Market {
	case common.Market_SPOT:
		{
			symbol := req.Asset
			states := "filled,partial-canceled,canceled"
			options := &url.Values{}
			options.Add("start-time", strconv.FormatInt(req.StartTime, 10))
			options.Add("end-time", strconv.FormatInt(req.EndTime, 10))

			resHistory, err := c.spot_api.GetHistoryOrders(symbol, states, options)
			if err != nil {
				fmt.Println("c.api.GetHistoryOrders =", err)
				return res, nil
			}
			for _, data := range resHistory.Data {
				res_tmp := &client.OrderRsp{
					OrderId:   strconv.FormatInt(data.ID, 10),
					Timestamp: data.CreatedAt,
					Symbol:    data.Symbol,
				}
				immutableT := reflect.TypeOf(data)
				if _, ok := immutableT.FieldByName("ClientOrderID"); ok {
					res_tmp.Id, _ = strconv.ParseInt(data.ClientOrderID, 10, 0)
				}
				res_tmp.Price, _ = strconv.ParseFloat(data.Price, 64)
				res_tmp.AccumQty, _ = strconv.ParseFloat(data.FieldCashAmount, 64)
				res_tmp.AccumAmount, _ = strconv.ParseFloat(data.FieldAmount, 64)
				res_tmp.AvgPrice = res_tmp.AccumQty / res_tmp.AccumAmount
				res_tmp.Fee, _ = strconv.ParseFloat(data.FieldFees, 64)

				res = append(res, res_tmp)
			}
		}
	case common.Market_FUTURE_COIN:
		{
			post_params := c_api.ReqPostFutureContractMatchresults{
				Symbol:    req.Asset[:3],
				TradeType: 0,
				StartTime: req.StartTime,
				EndTime:   req.EndTime,
			}
			resHistory, err := c.c_api.PostFutureContractMatchresults(post_params)
			if err != nil {
				fmt.Println("c.c_api.PostFutureContractHisorders =", err)
				return res, nil
			}
			for _, data := range resHistory.Data {
				res_tmp := &client.OrderRsp{
					OrderId:   data.OrderIDStr,
					Timestamp: data.CreateDate,
					Symbol:    data.Symbol,
				}
				res_tmp.Price = float64(data.TradePrice)
				res_tmp.AccumQty = float64(data.TradeTurnover)
				res_tmp.AccumAmount = float64(data.TradeVolume)
				res_tmp.AvgPrice = res_tmp.AccumQty / res_tmp.AccumAmount
				res_tmp.Fee = float64(data.TradeFee)
				res_tmp.FeeAsset = data.FeeAsset
				res = append(res, res_tmp)
			}
		}
	case common.Market_SWAP_COIN:
		{
			post_params := c_api.ReqPostSwapContractMatchresults{
				Contract:  strings.Replace(req.Asset, "/", "-", 1),
				TradeType: 0,
				StartTime: req.StartTime,
				EndTime:   req.EndTime,
			}
			resHistory, err := c.c_api.PostSwapContractMatchresults(post_params)
			if err != nil {
				fmt.Println("c.c_api.PostSwapContractMatchresults =", err)
				return res, nil
			}
			for _, data := range resHistory.Data {
				res_tmp := &client.OrderRsp{
					OrderId:   data.OrderIDStr,
					Timestamp: data.CreateDate,
					Symbol:    data.Symbol,
				}
				res_tmp.Price = float64(data.TradePrice)
				res_tmp.AccumQty = float64(data.TradeTurnover)
				res_tmp.AccumAmount = float64(data.TradeVolume)
				res_tmp.AvgPrice = res_tmp.AccumQty / res_tmp.AccumAmount
				res_tmp.Fee = float64(data.TradeFee)
				res_tmp.FeeAsset = data.FeeAsset
				res = append(res, res_tmp)
			}
		}
	case common.Market_FUTURE, common.Market_SWAP:
		{
			post_params := u_api.ReqPostSwapCrossMatchresults{
				Contract:  u_api.GetContractCodeFromSymbolAndTypeU(req.Asset, req.Type),
				Pair:      strings.Replace(req.Asset, "/", "-", 1),
				TradeType: 0,
				StartTime: req.StartTime,
				EndTime:   req.EndTime,
			}
			resHistory, err := c.u_api.PostSwapCrossMatchresults(post_params)
			if err != nil {
				fmt.Println("c.u_api.PostSwapCrossMatchresults =", err)
				return res, nil
			}
			for _, data := range resHistory.Data {
				res_tmp := &client.OrderRsp{
					OrderId:   data.OrderIDStr,
					Timestamp: data.CreateDate,
					Symbol:    data.Symbol,
				}
				res_tmp.Price = float64(data.TradePrice)
				res_tmp.AccumQty = float64(data.TradeTurnover)
				res_tmp.AccumAmount = float64(data.TradeVolume)
				res_tmp.AvgPrice = res_tmp.AccumQty / res_tmp.AccumAmount
				res_tmp.Fee = float64(data.TradeFee)
				res_tmp.FeeAsset = data.FeeAsset
				res = append(res, res_tmp)
			}
		}
	}

	return res, nil
} //查询订单信息。不支持ws回报订单行情的交易所，如bitbank, 需要频发查单获知订单状态

// 查询当前未成交订单
func (c *ClientHuobi) GetProcessingOrders(req *client.OrderHistoryReq) ([]*client.OrderRsp, error) {
	var res []*client.OrderRsp

	switch req.Market {
	case common.Market_SPOT:
		{
			accountInfo, err := c.spot_api.Account()
			if err != nil {
				fmt.Println("c.api.Account err =", err)
				return res, err
			}

			var account_id string
			symbol := req.Asset
			for _, accountData := range accountInfo.Data {
				account_id = strconv.FormatInt(accountData.Id, 10)

				options := &url.Values{}
				options.Add("account-id", account_id)
				options.Add("symbol", symbol)
				openOrderResp, err1 := c.spot_api.GetOpenOrders(options)
				if err1 != nil {
					fmt.Println("account-id", account_id, "c.api.GetOpenOrders err =", err1)
					return res, err1
				}

				for _, openOrderData := range openOrderResp.Data {
					resTmp := &client.OrderRsp{
						Timestamp: openOrderData.CreatedAt,
						Symbol:    symbol,
					}
					resTmp.OrderId = strconv.FormatInt(openOrderData.ID, 10)
					resTmp.Id, _ = strconv.ParseInt(openOrderData.ClientOrderID, 10, 0)
					resTmp.Price, _ = strconv.ParseFloat(openOrderData.Price, 64)
					resTmp.AccumAmount, _ = strconv.ParseFloat(openOrderData.FilledAmount, 64)
					resTmp.AccumQty, _ = strconv.ParseFloat(openOrderData.FilledCashAmount, 64)
					resTmp.AvgPrice = resTmp.AccumQty / resTmp.AccumAmount
					resTmp.Fee, _ = strconv.ParseFloat(openOrderData.FilledFees, 64)

					matchresultsRsp, _ := c.spot_api.GetOrderMatchresults(resTmp.OrderId)
					resTmp.FeeAsset = matchresultsRsp.Data[0].FeeCurrency

					res = append(res, resTmp)
				}
			}
		}
	case common.Market_FUTURE_COIN:
		{
			post_params := c_api.ReqPostFutureContractOpenorders{
				Symbol: req.Asset[:3],
			}
			ResOpenOrders, err := c.c_api.PostFutureContractOpenorders(post_params)
			if err != nil {
				fmt.Println("c.c_api.PostFutureContractOpenorders err =", err)
				return nil, err
			}
			for _, DataOpenOrder := range ResOpenOrders.Data.Orders {
				res_tmp := &client.OrderRsp{
					Id:          DataOpenOrder.ClientOrderID,
					OrderId:     DataOpenOrder.OrderIDStr,
					Timestamp:   DataOpenOrder.CreatedAt,
					Symbol:      req.Asset,
					AccumQty:    float64(DataOpenOrder.TradeTurnover),
					AccumAmount: float64(DataOpenOrder.TradeVolume),
					FeeAsset:    DataOpenOrder.FeeAsset,
					Fee:         float64(DataOpenOrder.Fee),
				}
				res_tmp.AvgPrice = res_tmp.AccumQty / res_tmp.AccumAmount
				res = append(res, res_tmp)
			}
		}
	case common.Market_SWAP_COIN:
		{
			post_params := c_api.ReqPostSwapContractOpenorders{
				ContractCode: strings.Replace(req.Asset, "/", "-", 1),
			}
			ResOpenOrders, err := c.c_api.PostSwapContractOpenorders(post_params)
			if err != nil {
				fmt.Println("c.c_api.PostFutureContractOpenorders err =", err)
				return nil, err
			}
			for _, DataOpenOrder := range ResOpenOrders.Data.Orders {
				res_tmp := &client.OrderRsp{
					Id:          DataOpenOrder.ClientOrderID,
					OrderId:     DataOpenOrder.OrderIDStr,
					Timestamp:   DataOpenOrder.CreatedAt,
					Symbol:      req.Asset,
					AccumQty:    float64(DataOpenOrder.TradeTurnover),
					AccumAmount: float64(DataOpenOrder.TradeVolume),
					FeeAsset:    DataOpenOrder.FeeAsset,
					Fee:         float64(DataOpenOrder.Fee),
				}
				res_tmp.AvgPrice = res_tmp.AccumQty / res_tmp.AccumAmount
				res = append(res, res_tmp)
			}
		}
	case common.Market_FUTURE, common.Market_SWAP:
		{
			post_params := u_api.ReqPostSwapOpenOrders{
				ContractCode: strings.Replace(req.Asset, "/", "-", 1),
			}
			ResOpenOrders, err := c.u_api.PostSwapOpenOrders(post_params)
			if err != nil {
				fmt.Println("c.u_api.PostSwapOpenOrders err =", err)
				return nil, err
			}
			for _, DataOpenOrder := range ResOpenOrders.Data.Orders {
				res_tmp := &client.OrderRsp{
					Id:          DataOpenOrder.ClientOrderID,
					OrderId:     DataOpenOrder.OrderIDStr,
					Timestamp:   DataOpenOrder.CreatedAt,
					Symbol:      req.Asset,
					AccumQty:    float64(DataOpenOrder.TradeTurnover),
					AccumAmount: float64(DataOpenOrder.TradeVolume),
					FeeAsset:    DataOpenOrder.FeeAsset,
					Fee:         float64(DataOpenOrder.Fee),
				}
				res_tmp.AvgPrice = res_tmp.AccumQty / res_tmp.AccumAmount
				res = append(res, res_tmp)
			}
		}
	}

	return res, nil
} //查询订单信息。不支持ws回报订单行情的交易所，如bitbank, 需要频发查单获知订单状态

// 转账订单，提币
func (c *ClientHuobi) Transfer(o *order.OrderTransfer) (*client.OrderRsp, error) {

	// 先不管chain
	var (
		res            = &client.OrderRsp{}
		address        = string(o.TransferAddress)
		amountOriginal = o.Amount
		currency       = string(o.ExchangeToken)
		addressTag     = string(o.Tag)
		baseChain      = spot_api.GetNetWorkFromChain(o.Chain)
		Chain          = ""
		clientOrderId  = transform.IdToClientId(o.Hdr.Producer, o.Base.Id)

		post_params = spot_api.ReqPostWithdrawCreate{
			Address:       address,
			Currency:      currency,
			AddrTag:       addressTag,
			ClientOrderId: clientOrderId,
		}
	)

	// 计算手续费
	// 查询手续费or手续费率
	resGetTransferFee, _ := c.GetTransferFee(o.Chain, currency)
	var fee float64
	if resGetTransferFee.TransferFeeList[0].Fee != 0 {
		fee = resGetTransferFee.TransferFeeList[0].Fee
	} else {
		fee = resGetTransferFee.TransferFeeList[0].FeeRate * amountOriginal
	}
	post_params.Amount = transform.XToString(amountOriginal - fee)
	post_params.Fee = transform.XToString(fee)

	// 获取chain
	options := &url.Values{}
	options.Add("currency", strings.ToLower(currency))
	resGetReferenceCurrencies, _ := c.spot_api.GetReferenceCurrencies(options)

	for _, chainData := range resGetReferenceCurrencies.Data[0].Chains {
		if chainData.BaseChain == baseChain || chainData.DisplayName == baseChain {
			Chain = chainData.Chain
		}
	}
	post_params.Chain = Chain

	resWithdraw, err := c.spot_api.PostWithdrawCreate(post_params)
	if err != nil {
		fmt.Println("c.api.PostWithdrawCreate err =", err)
		return nil, err
	}
	res.OrderId = strconv.FormatInt(int64(resWithdraw.Data), 10)
	if res.OrderId == "0" {
		fmt.Println("withdraw failed")
	}

	return res, nil
}

// 提币记录
func (c *ClientHuobi) GetTransferHistory(req *client.TransferHistoryReq) (*client.TransferHistoryRsp, error) {

	var (
		resp     = &client.TransferHistoryRsp{}
		err      error
		currency = req.Asset
		type_    = "withdraw"
		options  = &url.Values{}
	)

	if currency != "" {
		options.Add("currency", currency)
	}

	resHistory, err := c.spot_api.GetDepositHistory(type_, options)
	if err != nil {
		fmt.Println("c.api.GetDepositHistory err =", err)
		return nil, err
	}

	optionsTmp := &url.Values{}
	optionsTmp.Add("currency", currency)
	resGetReferenceCurrencies, _ := c.spot_api.GetReferenceCurrencies(optionsTmp)

	for _, historyData := range resHistory.Data {

		Chain := historyData.Chain
		var baseChain string
		for _, chainData := range resGetReferenceCurrencies.Data[0].Chains {
			if chainData.Chain == Chain {
				baseChain = chainData.BaseChain
				break
			}
		}

		resp.TransferList = append(resp.TransferList, &client.TransferHistoryItem{
			OrderId: strconv.FormatInt(int64(historyData.ID), 10),
			Asset:   currency,
			Amount:  float64(historyData.Amount),
			Fee:     float64(historyData.Fee),
			Address: historyData.Address,
			TxId:    historyData.TxHash,
			Network: spot_api.GetChainFromNetWork(baseChain),
		})
	}

	return resp, nil
}

// 资产划转订单
func (c *ClientHuobi) MoveAsset(o *order.OrderMove) (*client.OrderRsp, error) {

	var (
		amount        = o.Amount
		currency      = strings.ToLower(o.Asset)
		transfer_type string
	)

	res := &client.OrderRsp{}

	if o.Source == common.Market_SPOT {
		transfer_type = "pro-to-futures"
	} else {
		transfer_type = "futures-to-pro"
	}

	if o.Source == common.Market_FUTURE_COIN || o.Target == common.Market_FUTURE_COIN {
		post_params := spot_api.ReqPostCFuturesTransfer{
			Currency: currency,
			Amount:   amount,
			Type:     transfer_type,
		}
		resTransfer, err := c.spot_api.PostCFuturesTransfer(post_params)
		if err != nil {
			fmt.Println("c.spot_api.PostCFuturesTransfer err =", err)
			return nil, err
		}
		res.OrderId = transform.XToString(resTransfer.Data)
	} else if o.Source == common.Market_SWAP_COIN || o.Target == common.Market_SWAP_COIN {
		post_params := spot_api.ReqPostCSwapTransfer{
			From:     spot_api.GetTypeFromType(o.Source),
			To:       spot_api.GetTypeFromType(o.Target),
			Currency: currency,
			Amount:   amount,
		}
		resTransfer, err := c.spot_api.PostCSwapTransfer(post_params)
		if err != nil {
			fmt.Println("c.spot_api.PostCSwapTransfer err =", err)
			return nil, err
		}
		res.OrderId = transform.XToString(resTransfer.Data)
	} else if o.Source == common.Market_FUTURE || o.Source == common.Market_SWAP || o.Target == common.Market_FUTURE || o.Target == common.Market_SWAP {
		post_params := spot_api.ReqPostUTransfer{
			From:          spot_api.GetTypeFromType(o.Source),
			To:            spot_api.GetTypeFromType(o.Target),
			Currency:      currency,
			Amount:        amount,
			MarginAccount: "USDT",
		}
		resTransfer, err := c.spot_api.PostUTransfer(post_params)
		if err != nil {
			fmt.Println("c.spot_api.PostUTransfer err =", err)
			return nil, err
		}
		res.OrderId = transform.XToString(resTransfer.Data)
	}

	return res, nil
}

// 查询划转记录，财务流水
func (c *ClientHuobi) GetMoveHistory(req *client.MoveHistoryReq) (*client.MoveHistoryRsp, error) {

	res := &client.MoveHistoryRsp{}

	source := req.AccountSource
	target := req.AccountTarget
	currency := req.SymbolSource
	transactypes := "transfer"
	startTime := strconv.FormatInt(int64(req.StartTime), 10)
	endTime := strconv.FormatInt(int64(req.EndTime), 10)
	options := &url.Values{}

	options.Add("currency", currency)
	options.Add("transactTypes", transactypes)
	options.Add("startTime", startTime)
	options.Add("endTime", endTime)

	account_id := source
	resGetAccountLedger, err := c.spot_api.GetAccountLedger(account_id, options)
	if err != nil {
		fmt.Println("c.api.GetAccountLedger err =", err)
		return res, err
	}
	for _, dataLedger := range resGetAccountLedger.Data {
		if strconv.FormatInt(int64(dataLedger.Transferee), 10) != target {
			continue
		}
		res.MoveList = append(res.MoveList, &client.MoveHistoryItem{
			Asset:     dataLedger.Currency,
			Type:      client.MoveType_MOVETYPE_OUT, // 有点问题，针对的账户对象不一样，划转方向也不一样
			Id:        strconv.FormatInt(int64(dataLedger.TransactID), 10),
			Amount:    dataLedger.TransactAmt,
			Timestamp: dataLedger.TransactTime,
			Status:    client.MoveStatus_MOVESTATUS_CONFIRMED,
		})
	}

	return res, nil
}

func StatusToInt(status string) int32 {
	// 0(0:pending,6: credited but cannot withdraw, 1:success)
	switch status {
	case "unknown": // 状态未知
		return 0
	case "confirming": //区块确认中
		return 6
	case "confirmed": //区块已完成，已经上账，可以划转和交易
		return 1
	case "safe":
		return 1 //区块已确认，可以提币
	case "orphan":
		return 0 //区块已被孤立
	default:
		return 0
	}
}

func (c *ClientHuobi) GetDepositHistory(req *client.DepositHistoryReq) (*client.DepositHistoryRsp, error) {
	var (
		resp      = &client.DepositHistoryRsp{}
		err       error
		coin      = req.Asset                                     //coin	STRING	NO
		status    = spot_api.GetDepositTypeToExchange(req.Status) //status	INT	NO	0(0:pending,6: credited but cannot withdraw, 1:success)
		startTime = req.StartTime                                 //startTime	LONG	NO	默认当前时间90天前的时间戳
		endTime   = req.EndTime                                   //endTime	LONG	NO	默认当前时间戳
		offset    = 0                                             //offset	INT	NO	默认:0
		limit     = 1000                                          //limit	INT	NO	默认：1000，最大1000
		params    = url.Values{}
		res       *spot_api.RespCapitalDepositHisRec
	)
	if status == -1 {
		return nil, errors.New("deposit status err")
	}

	if coin != "" {
		params.Add("coin", coin)
	}
	if status >= 0 {
		params.Add("status", strconv.FormatInt(int64(status), 10))
	}
	if startTime > 0 {
		params.Add("startTime", strconv.FormatInt(startTime, 10))
	}
	if endTime > 0 {
		params.Add("endTime", strconv.FormatInt(endTime, 10))
	}

	for i := 0; i < 100; i++ {
		params.Add("offset", strconv.Itoa(offset))
		params.Add("limit", strconv.Itoa(limit))

		res, err = c.spot_api.CapitalDepositHisRec(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range (*res).Data {

			resp.DepositList = append(resp.DepositList, &client.DepositHistoryItem{
				Asset:      item.Currency,
				Amount:     item.Amount,
				Network:    spot_api.GetChainFromNetWork(strings.ToUpper(item.Chain)),
				Status:     spot_api.GetDepositTypeFromExchange(spot_api.DepositStatus((StatusToInt(item.State)))),
				Address:    item.Address,
				AddressTag: item.AddressTag,
				TxId:       item.TxHash,
				Timestamp:  item.CreatedAt,
			})
		}
		if len((*res).Data) < limit {
			break
		}
		offset += limit
	}
	return resp, err
}

// 借贷
func (c *ClientHuobi) Loan(o *order.OrderLoan) (*client.OrderRsp, error) {
	res := &client.OrderRsp{}

	post_params := spot_api.ReqPostCrossMarginOrders{
		Currency: o.Asset,
		Amount:   strconv.FormatFloat(o.Amount, 'e', 3, 64),
	}

	resMarginLoan, err := c.spot_api.PostCrossMarginOrders(post_params)
	if err != nil {
		fmt.Println("c.api.PostCrossMarginOrders err =", err)
		return res, err
	}

	res.OrderId = strconv.FormatInt(int64(resMarginLoan.Data), 10)
	res.RespType = client.OrderRspType_ACK

	return nil, nil
}

// 获取已放款订单
func (c *ClientHuobi) GetLoanOrders(o *client.LoanHistoryReq) (*client.LoanHistoryRsp, error) {
	var (
		res        = &client.LoanHistoryRsp{}
		start_time = o.StartTime
		end_time   = o.EndTime
		currency   = o.Asset
		options    = &url.Values{}
	)
	options.Add("start-date", strconv.FormatInt(start_time, 10))
	options.Add("end-date", strconv.FormatInt(end_time, 10))
	options.Add("currency", currency)

	resGetCrossMarginLoan, err := c.spot_api.GetCrossMarginLoanOrders(options)
	if err != nil {
		fmt.Println("c.api.GetCrossMarginLoanOrders =", err)
		return nil, err
	}

	for _, dataResGetCrossMarginLoan := range resGetCrossMarginLoan.Data {
		principal, _ := strconv.ParseFloat(dataResGetCrossMarginLoan.LoanBalance, 64)
		res.LoadList = append(res.LoadList, &client.LoanHistoryItem{
			OrderId:   strconv.FormatInt(int64(dataResGetCrossMarginLoan.ID), 10),
			Asset:     dataResGetCrossMarginLoan.Currency,
			Principal: principal,
			Timestamp: dataResGetCrossMarginLoan.CreatedAt,
		})
	}

	return res, nil
}

// 还币，返回值过少
func (c *ClientHuobi) Return(o *order.OrderReturn) (*client.OrderRsp, error) {

	orderId := o.Base.IdEx
	post_params := spot_api.ReqPostCrossMarginRepay{
		Amount: strconv.FormatFloat(o.Amount, 'e', 5, 64),
	}

	resPostCrossMarginRepay, err := c.spot_api.PostCrossMarginRepay(orderId, post_params)
	if err != nil {
		fmt.Println("c.api.PostCrossMarginRepay err =", err)
		return nil, err
	}
	fmt.Println("status:", resPostCrossMarginRepay.Status)

	return nil, nil
}

// 获取已还款订单
func (c *ClientHuobi) GetReturnOrders(o *client.LoanHistoryReq) (*client.ReturnHistoryRsp, error) {

	var (
		startTime = strconv.FormatInt(o.StartTime, 10)
		endTime   = strconv.FormatInt(o.EndTime, 10)
		currency  = o.Asset
		options   = &url.Values{}
	)
	options.Add("startTime", startTime)
	options.Add("endTime", endTime)
	options.Add("currency", currency)

	resGetRepayment, err := c.spot_api.GetRepayment(options)
	if err != nil {
		fmt.Println("c.api.GetRepayment err =", err)
		return nil, err
	}

	res := client.ReturnHistoryRsp{}
	for _, dataResGetRepayment := range resGetRepayment.Data {
		amount, _ := strconv.ParseFloat(dataResGetRepayment.RepaidAmount, 64)
		principal, _ := strconv.ParseFloat(dataResGetRepayment.TransactIds.Repaidprincipal, 64)
		res.ReturnList = append(res.ReturnList, &client.ReturnHistoryItem{
			Amount:    amount,
			OrderId:   dataResGetRepayment.RepayID,
			Asset:     dataResGetRepayment.Currency,
			Principal: principal,
			Timestamp: dataResGetRepayment.RepayTime,
		})
	}

	return nil, nil
}

// 合约所有交易对
func (c *ClientHuobi) GetFutureSymbols(market common.Market) []*client.SymbolInfo {

	var res []*client.SymbolInfo
	options := &url.Values{}
	if market == common.Market_FUTURE_COIN {
		// future 币本位交割
		ResContractInfo, err := c.c_api.GetFutureContractInfo(&url.Values{})
		if err != nil {
			fmt.Println("c.c_api.GetFutureContractInfo err =", err)
			return nil
		}
		for _, StateData := range ResContractInfo.Data {
			res = append(res, &client.SymbolInfo{
				Symbol: StateData.Symbol + "/USD",
				Name:   StateData.ContractCode,
				Type:   c_api.GetContractType(market, StateData.ContractType),
			})
		}
	} else if market == common.Market_SWAP_COIN {
		// swap 币本位永续
		ResContractInfo, err := c.c_api.GetSwapContractInfo(options)
		if err != nil {
			fmt.Println("c.c_api.GetFutureContractInfo err =", err)
			return nil
		}
		for _, StateData := range ResContractInfo.Data {
			res = append(res, &client.SymbolInfo{
				Symbol: strings.Replace(StateData.ContractCode, "-", "/", 1),
				Name:   StateData.ContractCode,
				Type:   c_api.GetContractType(market, StateData.ContractType),
			})
		}
	} else if market == common.Market_FUTURE {
		// U本位交割
		ResContractInfo, err := c.u_api.GetContractInfo("future")
		if err != nil {
			fmt.Println("c.u_api.GetContractInfo err =", err)
			return nil
		}
		for _, StateData := range ResContractInfo.Data {
			res = append(res, &client.SymbolInfo{
				Symbol: strings.Replace(StateData.Pair, "-", "/", 1),
				Name:   StateData.ContractCode,
				Type:   c_api.GetContractType(market, StateData.ContractType),
			})
		}
	} else if market == common.Market_SWAP {
		//U本位永续
		ResContractInfo, err := c.u_api.GetContractInfo("swap")
		if err != nil {
			fmt.Println("c.u_api.GetContractInfo err =", err)
			return nil
		}
		for _, StateData := range ResContractInfo.Data {
			res = append(res, &client.SymbolInfo{
				Symbol: strings.Replace(StateData.Pair, "-", "/", 1),
				Name:   StateData.ContractCode,
				Type:   c_api.GetContractType(market, StateData.ContractType),
			})
		}
	}

	return res
}

// 获取行情
func (c *ClientHuobi) GetFutureDepth(o *client.SymbolInfo, i int) (*depth.Depth, error) {
	res := &depth.Depth{}
	ResDepth := &c_api.DepthInfo{}
	var err error
	contract_code := c_api.GetContractCodeFromSymbolAndType(o.Symbol, o.Type)
	switch o.Market {
	case common.Market_FUTURE_COIN:
		{
			// 币本位交割
			ResDepth, err = c.c_api.GetFutureDepth(contract_code, "step0")
			if err != nil {
				fmt.Println("c.c_api.GetFutureDepth err =", err)
				return nil, err
			}

		}
	case common.Market_SWAP_COIN:
		{
			// 币本位永续
			ResDepth, err = c.c_api.GetSwapDepth(contract_code, "step0")
			if err != nil {
				fmt.Println("c.c_api.GetFutureDepth err =", err)
				return nil, err
			}
		}
	case common.Market_FUTURE, common.Market_SWAP:
		{
			// U本位交割
			ResDepth, err = c.u_api.GetDepth(contract_code, "step0")
			if err != nil {
				fmt.Println("c.u_api.GetDepth err =", err)
				return nil, err
			}
		}
	}
	for _, ask := range ResDepth.Tick.Asks {
		res.Asks = append(res.Asks, &depth.DepthLevel{
			Price:  ask[0],
			Amount: ask[1],
		})
	}
	for _, bid := range ResDepth.Tick.Bids {
		res.Bids = append(res.Bids, &depth.DepthLevel{
			Price:  bid[0],
			Amount: bid[1],
		})
	}
	res.Symbol = o.Symbol
	res.Type = o.Type
	res.Exchange = common.Exchange_HUOBI
	res.TimeExchange = uint64(ResDepth.Tick.Ts)
	res.TimeOperate = uint64(time.Now().UnixNano()) / 1e6
	res.TimeReceive = uint64(time.Now().UnixNano()) / 1e6

	return res, nil
}

// 标记价格
func (c *ClientHuobi) GetFutureMarkPrice(market common.Market, symbols ...*client.SymbolInfo) (*client.RspMarkPrice, error) {

	res := &client.RspMarkPrice{}
	switch market {
	case common.Market_SWAP_COIN:
		{
			for _, symbolInfo := range symbols {
				contract_code := c_api.GetContractCodeFromSymbolAndType(symbolInfo.Symbol, symbolInfo.Type)
				period := "1min"
				size := 1
				ResMarkPriceKline, err := c.c_api.GetSwapMarkPriceKline(contract_code, period, int64(size))
				if err != nil {
					fmt.Println("c.c_api.GetSwapMarkPriceKline err =", err)
					continue
				}
				resTmp := &client.MarkPriceItem{
					Symbol:     symbolInfo.Symbol,
					Type:       symbolInfo.Type,
					UpdateTime: ResMarkPriceKline.Ts,
				}
				resTmp.MarkPrice, _ = strconv.ParseFloat(ResMarkPriceKline.Data[0].Close, 64)
				res.Item = append(res.Item, resTmp)
			}

		}
	case common.Market_SWAP:
		{
			// U本位永续
			for _, symbolInfo := range symbols {
				contract_code := c_api.GetContractCodeFromSymbolAndType(symbolInfo.Symbol, symbolInfo.Type)
				period := "1min"
				size := 1
				ResMarkPriceKline, err := c.u_api.GetLinerSwapMarkPriceKline(contract_code, period, int64(size))
				if err != nil {
					fmt.Println("c.u_api.GetLinerSwapMarkPriceKline err =", err)
					continue
				}
				resTmp := &client.MarkPriceItem{
					Symbol:     symbolInfo.Symbol,
					Type:       symbolInfo.Type,
					UpdateTime: ResMarkPriceKline.Ts,
				}
				resTmp.MarkPrice, _ = strconv.ParseFloat(ResMarkPriceKline.Data[0].Close, 64)
				res.Item = append(res.Item, resTmp)
			}
		}

	}

	return res, nil
}

// 交易手续费
func (c *ClientHuobi) GetFutureTradeFee(market common.Market, symbols ...*client.SymbolInfo) (*client.TradeFee, error) {
	res := &client.TradeFee{}
	switch market {
	case common.Market_FUTURE_COIN:
		{
			// 币本位交割
			for _, symbolInfo := range symbols {
				post_params := c_api.ReqPostFutureContractFee{
					Symbol: symbolInfo.Symbol[:3],
				}
				ResContractFee, err := c.c_api.PostFutureContractFee(post_params)
				if err != nil {
					fmt.Println("c.c_api.PostFutureContractFee err =", err)
					continue
				}
				maker, _ := strconv.ParseFloat(ResContractFee.Data[0].OpenMakerFee, 64)
				taker, _ := strconv.ParseFloat(ResContractFee.Data[0].OpenTakerFee, 64)
				res.TradeFeeList = append(res.TradeFeeList, &client.TradeFeeItem{
					Symbol: symbolInfo.Symbol,
					Type:   symbolInfo.Type,
					Maker:  maker,
					Taker:  taker,
				})

			}

		}
	case common.Market_SWAP_COIN:
		{
			// 币本位永续
			for _, symbolInfo := range symbols {
				post_params := c_api.ReqPostSwapContractFee{
					ContractCode: strings.Replace(symbolInfo.Symbol, "/", "-", 1),
				}
				ResContractFee, err := c.c_api.PostSwapContractFee(post_params)
				if err != nil {
					fmt.Println("c.c_api.PostSwapContractFee err =", err)
					continue
				}
				maker, _ := strconv.ParseFloat(ResContractFee.Data[0].OpenMakerFee, 64)
				taker, _ := strconv.ParseFloat(ResContractFee.Data[0].OpenTakerFee, 64)
				res.TradeFeeList = append(res.TradeFeeList, &client.TradeFeeItem{
					Symbol: symbolInfo.Symbol,
					Type:   symbolInfo.Type,
					Maker:  maker,
					Taker:  taker,
				})

			}
		}
	case common.Market_FUTURE, common.Market_SWAP:
		{
			// U本位
			for _, symbolInfo := range symbols {
				post_params := u_api.ReqPostSwapFee{
					ContractCode: u_api.GetContractCodeFromSymbolAndTypeU(symbolInfo.Symbol, symbolInfo.Type),
				}
				ResContractFee, err := c.u_api.PostSwapFee(post_params)
				if err != nil {
					fmt.Println("c.u_api.PostSwapFee err =", err)
					continue
				}
				maker, _ := strconv.ParseFloat(ResContractFee.Data[0].OpenMakerFee, 64)
				taker, _ := strconv.ParseFloat(ResContractFee.Data[0].OpenTakerFee, 64)
				res.TradeFeeList = append(res.TradeFeeList, &client.TradeFeeItem{
					Symbol: symbolInfo.Symbol,
					Type:   symbolInfo.Type,
					Maker:  maker,
					Taker:  taker,
				})

			}
		}
	}

	return res, nil
}

// 查询交易对精度信息
func (c *ClientHuobi) GetFuturePrecision(market common.Market, symbols ...*client.SymbolInfo) (*client.Precision, error) {
	res := &client.Precision{}

	switch market {
	case common.Market_FUTURE_COIN, common.Market_SWAP_COIN:
		{
			res.PrecisionList = append(res.PrecisionList, &client.PrecisionItem{
				AmountMin: 1,
			})
		}
	case common.Market_FUTURE, common.Market_SWAP:
		{
			business_type := ""
			if market == common.Market_FUTURE {
				business_type = "future"
			} else {
				business_type = "swap"
			}
			ResContractInfo, err := c.u_api.GetSwapContractInfo(business_type)
			if err != nil {
				fmt.Println("c.u_api.GetSwapContractInfo err =", err)
				return nil, err
			}
			for _, dataConInfo := range ResContractInfo.Data {
				res.PrecisionList = append(res.PrecisionList, &client.PrecisionItem{
					Symbol:    strings.Replace(dataConInfo.Pair, "-", "/", 1),
					Type:      c_api.GetContractType(market, dataConInfo.ContractType),
					AmountMin: dataConInfo.ContractSize,
				})
			}
		}
	}

	return nil, nil
}

// 资产查询
func (c *ClientHuobi) GetFutureBalance(market common.Market) (*client.UBaseBalance, error) {
	res := &client.UBaseBalance{}
	switch market {
	case common.Market_FUTURE_COIN:
		{
			post_params1 := c_api.ReqPostFutureContractAccountInfo{}
			ResAccountInfo, err1 := c.c_api.PostFutureContractAccountInfo(post_params1)
			if err1 != nil {
				fmt.Println("c.c_api.PostFutureContractAccountInfo err =", err1)
				return nil, err1
			}
			for _, dataAccInfo := range ResAccountInfo.Data {
				resUBaseBalanceTmp := &client.UBaseBalanceItem{
					Asset:        dataAccInfo.Symbol,
					Balance:      dataAccInfo.MarginBalance,
					Market:       market,
					Unprofit:     float64(dataAccInfo.ProfitUnreal),
					Max_Withdraw: dataAccInfo.WithdrawAvailable,
					Available:    dataAccInfo.MarginAvailable,
				}
				resUBaseBalanceTmp.Rights = resUBaseBalanceTmp.Balance + resUBaseBalanceTmp.Unprofit
				resUBaseBalanceTmp.Used = resUBaseBalanceTmp.Rights - resUBaseBalanceTmp.Available
				res.UBaseBalanceList = append(res.UBaseBalanceList, resUBaseBalanceTmp)
			}
			post_params2 := c_api.ReqPostFutureContractPositionInfo{}
			ResPositionInfo, err2 := c.c_api.PostFutureContractPositionInfo(post_params2)
			if err2 != nil {
				fmt.Println("c.c_api.PostFutureContractPositionInfo err =", err2)
				return nil, err2
			}
			for _, dataPosInfo := range ResPositionInfo.Data {
				res.UBasePositionList = append(res.UBasePositionList, &client.UBasePositionItem{
					Symbol:         dataPosInfo.Symbol + "/USD",
					Position:       float64(dataPosInfo.Volume),
					Side:           c_api.GetTradeSideFromDirection(dataPosInfo.Direction),
					Market:         market,
					Type:           c_api.GetContractType(market, dataPosInfo.ContractType),
					Leverage:       float64(dataPosInfo.LeverRate),
					Unprofit:       float64(dataPosInfo.ProfitUnreal),
					MaintainMargin: dataPosInfo.PositionMargin,
					Price:          dataPosInfo.LastPrice,
				})
			}
		}
	case common.Market_SWAP_COIN:
		{
			post_params1 := c_api.ReqPostSwapAccountInfo{}
			ResAccountInfo, err1 := c.c_api.PostSwapAccountInfo(post_params1)
			if err1 != nil {
				fmt.Println("c.c_api.PostSwapAccountInfo err =", err1)
				return nil, err1
			}
			for _, dataAccInfo := range ResAccountInfo.Data {
				resUBaseBalanceTmp := &client.UBaseBalanceItem{
					Asset:        dataAccInfo.Symbol,
					Balance:      dataAccInfo.MarginBalance,
					Market:       market,
					Unprofit:     float64(dataAccInfo.ProfitUnreal),
					Max_Withdraw: dataAccInfo.WithdrawAvailable,
				}
				resUBaseBalanceTmp.Rights = resUBaseBalanceTmp.Balance + resUBaseBalanceTmp.Unprofit
				res.UBaseBalanceList = append(res.UBaseBalanceList, resUBaseBalanceTmp)
			}
			post_params2 := c_api.ReqPostSwapPositionInfo{}
			ResPositionInfo, err2 := c.c_api.PostSwapPositionInfo(post_params2)
			if err2 != nil {
				fmt.Println("c.c_api.PostSwapPositionInfo err =", err2)
				return nil, err2
			}
			for _, dataPosInfo := range ResPositionInfo.Data {
				res.UBasePositionList = append(res.UBasePositionList, &client.UBasePositionItem{
					Symbol:         strings.Replace(dataPosInfo.ContractCode, "-", "/", 1),
					Position:       float64(dataPosInfo.Volume),
					Side:           c_api.GetTradeSideFromDirection(dataPosInfo.Direction),
					Market:         market,
					Type:           common.SymbolType_SWAP_COIN_FOREVER,
					Leverage:       float64(dataPosInfo.LeverRate),
					Unprofit:       float64(dataPosInfo.ProfitUnreal),
					MaintainMargin: dataPosInfo.PositionMargin,
					Price:          dataPosInfo.LastPrice,
				})
			}
		}
	case common.Market_FUTURE, common.Market_SWAP:
		{
			post_params1 := u_api.ReqPostSwapCrossAccountInfo{}
			ResAccountInfo, err1 := c.u_api.PostSwapCrossAccountInfo(post_params1)
			if err1 != nil {
				fmt.Println("c.c_api.PostSwapAccountInfo err =", err1)
				return nil, err1
			}
			for _, dataAccInfo := range ResAccountInfo.Data {
				resUBaseBalanceTmp := &client.UBaseBalanceItem{
					Asset:        dataAccInfo.MarginAsset,
					Balance:      dataAccInfo.MarginBalance,
					Market:       market,
					Unprofit:     float64(dataAccInfo.ProfitUnreal),
					Max_Withdraw: float64(dataAccInfo.WithdrawAvailable),
				}
				resUBaseBalanceTmp.Rights = resUBaseBalanceTmp.Balance + resUBaseBalanceTmp.Unprofit
				res.UBaseBalanceList = append(res.UBaseBalanceList, resUBaseBalanceTmp)
			}
			post_params2 := u_api.ReqPostSwapCrossPositionInfo{
				ContractType: c_api.GetContractTypeFromType(common.SymbolType(market)),
			}
			ResPositionInfo, err2 := c.u_api.PostSwapCrossPositionInfo(post_params2)
			if err2 != nil {
				fmt.Println("c.u_api.PostSwapCrossPositionInfo err =", err2)
				return nil, err2
			}
			for _, dataPosInfo := range ResPositionInfo.Data {
				res.UBasePositionList = append(res.UBasePositionList, &client.UBasePositionItem{
					Symbol:         strings.Replace(dataPosInfo.ContractCode, "-", "/", 1),
					Position:       float64(dataPosInfo.Volume),
					Side:           c_api.GetTradeSideFromDirection(dataPosInfo.Direction),
					Market:         market,
					Type:           common.SymbolType_SWAP_COIN_FOREVER,
					Leverage:       float64(dataPosInfo.LeverRate),
					Unprofit:       float64(dataPosInfo.ProfitUnreal),
					MaintainMargin: dataPosInfo.PositionMargin,
					Price:          dataPosInfo.LastPrice,
				})
			}
		}
	}
	return res, nil
}

// 交易相关
// 下单
func (c *ClientHuobi) PlaceFutureOrder(o *order.OrderTradeCEX) (*client.OrderRsp, error) {
	res := &client.OrderRsp{}

	contract_type := c_api.GetContractTypeFromType(o.Base.Type)
	client_order_id := o.Base.Id
	price := o.Price
	volumn := o.Amount
	direction := c_api.GetDirectionFromTradeSide(o.Side)
	offset := "open"
	lever_rate := 5
	order_price_type := c_api.GetOrderPriceType(o.OrderType, o.TradeType, o.Tif)
	switch o.Base.Market {
	case common.Market_FUTURE_COIN:
		{
			// 币本位交割
			symbol := o.Base.Symbol[:3]
			post_params := c_api.ReqPostFutureOrder{
				Symbol:         string(symbol),
				ContractType:   contract_type,
				ClientOrderID:  client_order_id,
				Price:          int(price),
				Volume:         int(volumn),
				Direction:      direction,
				Offset:         offset,
				LeverRate:      lever_rate,
				OrderPriceType: order_price_type,
				//TpTriggerPrice: int(tp_trigger_price),
				//SlTriggerPrice: int(sl_trigger_price),
			}

			ResOrder, err := c.c_api.PostFutureOrder(post_params)
			if err != nil {
				fmt.Println("c.c_api.PostFutureOrder err =", err)
				return nil, err
			}
			res.Id = ResOrder.Data.ClientOrderID
			res.OrderId = ResOrder.Data.OrderIDStr
			res.Timestamp = ResOrder.Ts
		}
	case common.Market_SWAP_COIN:
		{
			// 币本位永续
			contract_code := strings.Replace(string(o.Base.Symbol), "/", "-", 1)
			post_params := c_api.ReqPostSwapOrder{
				ContractCode:   contract_code,
				ClientOrderID:  strconv.FormatInt(client_order_id, 10),
				Price:          int(price),
				Volume:         int(volumn),
				Direction:      direction,
				Offset:         offset,
				LeverRate:      lever_rate,
				OrderPriceType: order_price_type,
			}
			ResOrder, err := c.c_api.PostSwapOrder(post_params)
			if err != nil {
				fmt.Println("c.c_api.PostFutureOrder err =", err)
				return nil, err
			}
			res.Id = ResOrder.Data.ClientOrderID
			res.OrderId = ResOrder.Data.OrderIDStr
			res.Timestamp = ResOrder.Ts
		}
	case common.Market_FUTURE, common.Market_SWAP:
		{
			// U本位
			pair := strings.Replace(string(o.Base.Symbol), "/", "-", 1)
			symbol_list := strings.Split(string(o.Base.Symbol), "/")
			post_params := u_api.ReqPostSwapCrossOrder{
				Pair:           pair,
				ContractType:   contract_type,
				ClientOrderID:  transform.XToString(client_order_id),
				Price:          transform.XToString(price),
				Volume:         int64(volumn / c.u_api.Currency_size[symbol_list[0]]),
				Direction:      direction,
				Offset:         offset,
				LeverRate:      lever_rate,
				OrderPriceType: order_price_type,
				TpTriggerPrice: 0.01,
				TpOrderPrice:   0.01,
			}
			ResOrder, err := c.u_api.PostSwapCrossOrder(post_params)
			if err != nil {
				fmt.Println("c.u_api.PostSwapCrossOrder err =", err)
				return nil, err
			}
			res.Id = ResOrder.Data.ClientOrderID
			res.OrderId = ResOrder.Data.OrderIDStr
			res.Timestamp = ResOrder.Ts
		}
	}

	return res, nil
}
