package ftx

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/ftx/future_api"
	"clients/exchange/cex/ftx/spot_api"
	"clients/logger"
	"clients/transform"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
)

type ClientFTX struct {
	api  *spot_api.ApiClient
	uapi *future_api.UApiClient

	tradeFeeMap    map[string]*client.TradeFeeItem    //key: symbol
	transferFeeMap map[string]*client.TransferFeeItem //key: network+token
	precisionMap   map[string]*client.PrecisionItem   //key: symbol
}

func NewClientFTX(conf base.APIConf, maps ...interface{}) *ClientFTX {
	c := &ClientFTX{
		api:  spot_api.NewApiClient(conf, maps...),
		uapi: future_api.NewUApiClient(conf, maps...),
	}

	for _, m := range maps {
		switch t := m.(type) {
		case map[string]*client.TradeFeeItem:
			c.tradeFeeMap = t

		case map[string]*client.TransferFeeItem:
			c.transferFeeMap = t

		case map[string]*client.PrecisionItem:
			c.precisionMap = t
		}
	}
	return c
}

func NewClientFTX2(conf base.APIConf, cli *http.Client, maps ...interface{}) *ClientFTX {
	c := &ClientFTX{
		api:  spot_api.NewApiClient2(conf, cli, maps...),
		uapi: future_api.NewUApiClient2(conf, cli, maps...),
	}

	for _, m := range maps {
		switch t := m.(type) {
		case map[string]*client.TradeFeeItem:
			c.tradeFeeMap = t

		case map[string]*client.TransferFeeItem:
			c.transferFeeMap = t

		case map[string]*client.PrecisionItem:
			c.precisionMap = t
		}
	}

	return c
}

func (c *ClientFTX) GetExchange() common.Exchange {
	return common.Exchange_FTX
}

func (c *ClientFTX) GetSymbols() []string {
	var symbols []string
	exchangeInfoRes, err := c.api.GetMarkets()
	if err != nil {
		logger.Logger.Error("get exchange in error:", err)
		return symbols
	}
	for _, symbol := range exchangeInfoRes.Result {
		if symbol.Enabled && symbol.BaseCurrency != "" && symbol.QuoteCurrency != "" {
			symbols = append(symbols, fmt.Sprint(symbol.BaseCurrency, "/", symbol.QuoteCurrency))
		}
	}
	return symbols
}
func (c *ClientFTX) GetDepth(market *client.SymbolInfo, limit int) (*depth.Depth, error) {
	resDepth, err := c.api.GetOrderbook(market.Symbol, limit)
	if err != nil {
		//utils.Logger.Error(d.Symbol, d.getSymbolName(), "get full depth err", err)
		return nil, err
	}
	dep := &depth.Depth{
		Exchange:    common.Exchange_FTX,
		Market:      common.Market_SPOT,
		Type:        common.SymbolType_SPOT_NORMAL,
		Symbol:      market.Name,
		TimeReceive: uint64(time.Now().UnixMicro()),
		TimeOperate: uint64(time.Now().UnixMicro()),
	}
	err = c.ParseOrder(resDepth.Result.Bids, &dep.Bids)
	if err != nil {
		return dep, err
	}
	err = c.ParseOrder(resDepth.Result.Asks, &dep.Asks)
	if err != nil {
		return dep, err
	}
	return dep, err
} //获取行情
// How to get server time?
func (c *ClientFTX) IsExchangeEnable() bool {
	return true
}

func Weight() {

} // 权值信息

func (c *ClientFTX) GetTradeFee(symbols ...string) (*client.TradeFee, error) {
	var (
		res             = &client.TradeFee{}
		searchSymbolMap = make(map[string]bool)
	)
	for _, symbol := range symbols {
		searchSymbolMap[symbol] = true
	}
	tradeFees, err := c.api.GetAccountInfo()
	if err != nil {
		return nil, err
	}
	makerFee := tradeFees.Result.MakerFee
	takerFee := tradeFees.Result.TakerFee
	markets, err := c.api.GetMarkets()
	if err != nil {
		return nil, err
	}
	for _, marketName := range markets.Result {
		if len(searchSymbolMap) > 0 {
			_, ok := searchSymbolMap[marketName.Name]
			if !ok {
				continue
			}
		}
		res.TradeFeeList = append(res.TradeFeeList, &client.TradeFeeItem{
			Symbol: marketName.Name,
			Maker:  makerFee,
			Taker:  takerFee,
		})
	}
	return res, nil
}

// Tranfer Fee needs more parameters [size and address] TO-DO
func (c *ClientFTX) GetTransferFee(chain common.Chain, tokens ...string) (*client.TransferFee, error) {
	var (
		res     = &client.TransferFee{}
		options = make(map[string]interface{})
	)
	//for _, symbol := range tokens {
	//	searchSymbolMap[symbol] = true
	//}
	size := 0.0   //withdraw size
	address := "" //location
	options["method"] = spot_api.GetNetWorkFromChain(chain)
	for _, token := range tokens {
		withdrawalFee, err := c.api.GetWithdrawalFee(token, size, address, &options)
		if err != nil {
			continue
		}
		res.TransferFeeList = append(res.TransferFeeList, &client.TransferFeeItem{
			Token:   token,
			Network: chain,
			Fee:     withdrawalFee.Result.Fee,
		})
	}
	return res, nil
} //查询转账手续费

func (c *ClientFTX) GetPrecision(symbols ...string) (*client.Precision, error) {
	var (
		res             = &client.Precision{}
		searchSymbolMap = make(map[string]bool)
	)
	for _, symbol := range symbols {
		searchSymbolMap[symbol] = true
	}
	exchangeInfoRes, err := c.api.GetMarkets()
	if err != nil {
		return nil, err
	}
	for _, symbol := range exchangeInfoRes.Result {
		if len(searchSymbolMap) > 0 {
			_, ok := searchSymbolMap[symbol.Name]
			if !ok {
				continue
			}
		}
		var (
			price, amount int
		)
		if strings.Contains(fmt.Sprintf("%v", symbol.PriceIncrement), ".") {
			price = len(strings.Split(fmt.Sprintf("%v", symbol.PriceIncrement), ".")[1]) // 下单价格精度
		} else {
			price = -len(strings.Split(fmt.Sprintf("%v", symbol.PriceIncrement), ".")[0]) // 下单价格精度
		}
		if strings.Contains(fmt.Sprintf("%v", symbol.SizeIncrement), ".") {
			amount = len(strings.Split(fmt.Sprintf("%v", symbol.SizeIncrement), ".")[1]) // 下单价格精度
		} else {
			amount = -len(strings.Split(fmt.Sprintf("%v", symbol.SizeIncrement), ".")[0]) // 下单价格精度
		}
		if err != nil {
			return nil, err
		}
		res.PrecisionList = append(res.PrecisionList, &client.PrecisionItem{
			Symbol:    symbol.BaseCurrency + "/" + symbol.QuoteCurrency,
			Type:      common.SymbolType_SPOT_NORMAL,
			Amount:    int64(amount),
			Price:     int64(price),
			AmountMin: symbol.MinProvideSize,
		})
	}
	return res, nil
} //查询交易对精读信息
// 资产查询
func (c *ClientFTX) GetBalance() (*client.SpotBalance, error) {
	return nil, nil
} //获得现货的balance信息
// Incomplete function TO-Fill
func (c *ClientFTX) GetMarginBalance() (*client.MarginBalance, error) {
	var (
		respBalance *spot_api.RespGetBalances
		respAccount *spot_api.RespGetAccountInfo
		res         = &client.MarginBalance{}
		//marginLevel, totalNetAsset, , totalLiablityAsset float64
		err error
	)

	respAccount, err = c.api.GetAccountInfo()
	res.MarginLevel = 0
	res.TotalAsset = respAccount.Result.TotalAccountValue
	res.TotalNetAsset = respAccount.Result.TotalAccountValue - respAccount.Result.TotalPositionSize
	res.TotalLiabilityAsset = respAccount.Result.TotalPositionSize
	res.QuoteAsset = "USD"

	respBalance, err = c.api.GetBalances()
	if err != nil {
		return res, err
	}
	for _, asset := range respBalance.Result {
		if err != nil {
			return nil, err
		}
		res.MarginBalanceList = append(res.MarginBalanceList, &client.MarginBalanceItem{
			Asset:    asset.Coin,                                                    //币种
			Total:    asset.Total + asset.SpotBorrow,                                //总资产=已借+用户资产+锁定
			Borrowed: asset.SpotBorrow,                                              //已借
			Free:     asset.AvailableWithoutBorrow,                                  //可用=已借+用户资产
			Frozen:   asset.Total + asset.SpotBorrow - asset.AvailableWithoutBorrow, //锁定
			NetAsset: asset.AvailableWithoutBorrow - asset.SpotBorrow,
		})
	}
	return res, err
} //Need to check value correctness
func (c *ClientFTX) GetMarginIsolatedBalance(symbols ...string) (*client.MarginIsolatedBalance, error) {
	return nil, nil
} //Empty

// 交易相关
func (c *ClientFTX) PlaceOrder(o *order.OrderTradeCEX) (*client.OrderRsp, error) {
	var (
		symbol       string
		side         spot_api.SideType
		type_        spot_api.OrderType
		size         float64
		options      = make(map[string]interface{})
		resp         = &client.OrderRsp{}
		precisionIns interface{}
		precision    *client.PrecisionItem
		ok           bool
		err          error
	)

	if precisionIns, ok = c.precisionMap[string(o.Base.Symbol)]; !ok {
		return nil, errors.New("get precision err")
	}
	if precision, ok = precisionIns.(*client.PrecisionItem); !ok {
		return nil, errors.New("get precision err")
	}
	if o.Amount < precision.AmountMin {
		return nil, errors.New(fmt.Sprint("less amount in err:", o.Amount, "<", precision.AmountMin))
	}

	side = spot_api.GetSideTypeToExchange(o.Side)

	switch o.Tif {
	case order.TimeInForce_IOC:
		options["ioc"] = true
	default:
		options["ioc"] = false
	}

	if o.Base.Id != 0 {
		options["clientId"] = transform.IdToClientId(o.Hdr.Producer, o.Base.Id)
	}
	if o.OrderType == order.OrderType_MARKET {
		options["price"] = nil
	} else if o.OrderType == order.OrderType_LIMIT {
		options["price"] = o.Price
	}
	if o.TradeType == order.TradeType_MAKER {
		options["postOnly"] = true
	}

	size = o.Amount
	type_ = spot_api.GetOrderTypeToExchange(o.OrderType)
	symbol = string(o.Base.Symbol)

	var (
		res *spot_api.RespPlaceOrder
	)
	if o.Base.Market == common.Market_SPOT {
		res, err = c.api.PlaceOrder(symbol, side, type_, size, &options)
	}
	if err != nil {
		return nil, err
	}
	var closeData string
	names := strings.Split(res.Result.Market, "-")
	if len(names) > 1 {
		if names[1] != "PERP" && names[1] != "MOVE" {
			closeData = names[1]
		}
	}
	resp.Producer, resp.Id = transform.ClientIdToId(res.Result.ClientId)
	resp.OrderId = strconv.Itoa(res.Result.Id)
	resp.Timestamp = res.Result.CreatedAt.UnixMilli()
	resp.RespType = client.OrderRspType_FULL
	resp.Symbol = res.Result.Market //Remember to change to /
	resp.Status = spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(res.Result.Status), res.Result.FilledSize, res.Result.Size)
	resp.Price = res.Result.Price
	resp.Executed = res.Result.FilledSize
	resp.CloseDate = closeData
	return resp, nil
} //下单

func (c *ClientFTX) CancelOrder(o *order.OrderCancelCEX) (*client.OrderRsp, error) {
	var (
		resp         = &client.OrderRsp{}
		orderId, err = strconv.Atoi(o.Base.IdEx)
		statusErr    error
		clientId     = transform.IdToClientId(o.Hdr.Producer, o.Base.Id)
		res          *spot_api.RespCancelOrder
		status       *spot_api.RespGetOrderStatus
	)

	if orderId != 0 {
		res, err = c.api.CancelOrder(orderId)
		status, statusErr = c.api.GetOrderStatus(orderId)
	} else {
		res, err = c.api.CancelOrderByClientID(clientId)
		status, statusErr = c.api.GetOrderStatusByClientID(clientId)
	}

	if err != nil {
		return nil, err
	}
	if statusErr != nil {
		return nil, statusErr
	}
	if res.Success {
		var closeData string
		names := strings.Split(status.Result.Market, "-")
		if len(names) > 1 {
			if names[1] != "PERP" && names[1] != "MOVE" {
				closeData = names[1]
			}
		}
		resp.OrderId = strconv.Itoa(status.Result.Id)
		resp.Producer, resp.Id = transform.ClientIdToId(transform.IdToClientId(o.Hdr.Producer, o.Base.Id))
		resp.RespType = client.OrderRspType_RESULT
		resp.Symbol = status.Result.Market
		resp.Status = spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(status.Result.Status), status.Result.FilledSize, status.Result.Size)
		resp.Price = status.Result.AvgFillPrice
		resp.Executed = status.Result.FilledSize
		resp.CloseDate = closeData
	}
	return resp, nil
} //取消订单
func (c *ClientFTX) GetOrder(req *order.OrderQueryReq) (*client.OrderRsp, error) {
	var (
		resp       = &client.OrderRsp{}
		orderId, _ = strconv.Atoi(req.IdEx)
		statusErr  error
		clientId   = transform.IdToClientId(req.Producer, req.Id)
		status     *spot_api.RespGetOrderStatus
	)

	if clientId != "" {
		status, statusErr = c.api.GetOrderStatusByClientID(clientId)
	} else {
		status, statusErr = c.api.GetOrderStatus(orderId)
	}

	if statusErr != nil {
		return nil, statusErr
	}
	var closeData string
	names := strings.Split(status.Result.Market, "-")
	if len(names) > 1 {
		if names[1] != "PERP" && names[1] != "MOVE" {
			closeData = names[1]
		}
	}
	resp.OrderId = strconv.Itoa(status.Result.Id)
	resp.Producer, resp.Id = transform.ClientIdToId(status.Result.ClientId)
	resp.RespType = client.OrderRspType_RESULT
	resp.Symbol = status.Result.Market
	resp.Status = spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(status.Result.Status), status.Result.FilledSize, status.Result.Size)
	resp.Price = status.Result.AvgFillPrice
	resp.Executed = status.Result.FilledSize
	resp.CloseDate = closeData
	return resp, nil
} //查询订单信息。不支持ws回报订单行情的交易所，如bitbank, 需要频发查单获知订单状态
func (c *ClientFTX) GetOrderHistory(req *client.OrderHistoryReq) ([]*client.OrderRsp, error) {
	var (
		resp      []*client.OrderRsp
		symbol    string //symbol	STRING	YES
		market    = req.Market
		startTime = req.StartTime //startTime	LONG	NO
		endTime   = req.EndTime   //endTime	LONG	NO
		params    = make(map[string]interface{})
	)

	if req.Asset != "" {
		switch req.Market {
		case common.Market_SPOT:
			params["market"] = req.Asset
		case common.Market_FUTURE:
			currencies := strings.Split(req.Asset, "/")
			if strings.Contains(currencies[1], "USD") {
				symbol = currencies[0] + "-"
			} else {
				return nil, errors.New("market not supported on FTX")
			}
		case common.Market_SWAP:
			currencies := strings.Split(req.Asset, "/")
			if strings.Contains(currencies[1], "USD") {
				params["market"] = currencies[0] + "-PERP"
			} else {
				return nil, errors.New("market not supported on FTX")
			}
		default:
			return nil, errors.New("market not supported on FTX")
		}
	}

	if startTime > 0 {
		params["start_time"] = time.UnixMilli(startTime).Unix()
	}
	if endTime > 0 {
		params["end_time"] = time.UnixMilli(endTime).Unix()
	}

	res, err := c.api.GetOrderHistory(&params)
	if err != nil {
		return nil, err
	}

	if market == common.Market_SPOT || market == common.Market_SWAP {
		for _, item := range res.Result {
			producer, id := transform.ClientIdToId(item.ClientId)
			var closeData string
			names := strings.Split(item.Market, "-")
			if len(names) > 1 {
				if names[1] != "PERP" && names[1] != "MOVE" {
					closeData = names[1]
				}
			}
			resp = append(resp, &client.OrderRsp{
				Producer:  producer,
				Id:        id,
				OrderId:   strconv.Itoa(item.Id),
				Timestamp: item.CreatedAt.UnixMilli(),
				RespType:  client.OrderRspType_RESULT,
				Symbol:    item.Market,
				Status:    spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(item.Status), item.FilledSize, item.Size),
				Price:     item.AvgFillPrice,
				Executed:  item.FilledSize,
				CloseDate: closeData,
			})
		}
	} else {
		for _, item := range res.Result {
			if strings.Contains(item.Market, "PERP") || !strings.Contains(item.Market, symbol) {
				continue
			}
			producer, id := transform.ClientIdToId(item.ClientId)
			var closeData string
			names := strings.Split(item.Market, "-")
			if len(names) > 1 {
				if names[1] != "PERP" && names[1] != "MOVE" {
					closeData = names[1]
				}
			}
			resp = append(resp, &client.OrderRsp{
				Producer:  producer,
				Id:        id,
				OrderId:   strconv.Itoa(item.Id),
				Timestamp: item.CreatedAt.UnixMilli(),
				RespType:  client.OrderRspType_RESULT,
				Symbol:    item.Market,
				Status:    spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(item.Status), item.FilledSize, item.Size),
				Price:     item.AvgFillPrice,
				Executed:  item.FilledSize,
				CloseDate: closeData,
			})
		}
	}
	return resp, nil
} //查询订单信息。不支持ws回报订单行情的交易所，如bitbank, 需要频发查单获知订单状态
func (c *ClientFTX) GetProcessingOrders(req *client.OrderHistoryReq) ([]*client.OrderRsp, error) {
	var (
		resp   []*client.OrderRsp
		symbol = req.Asset //symbol	STRING	NO
		params = make(map[string]interface{})
	)
	if len(symbol) > 0 {
		params["market"] = symbol
	}
	res, err := c.api.GetOpenOrders(&params)
	if err != nil {
		return nil, err
	}
	for _, item := range res.Result {
		producer, id := transform.ClientIdToId(item.ClientId)
		var closeData string
		names := strings.Split(item.Market, "-")
		if len(names) > 1 {
			if names[1] != "PERP" && names[1] != "MOVE" {
				closeData = names[1]
			}
		}
		resp = append(resp, &client.OrderRsp{
			Producer:  producer,
			Id:        id,
			OrderId:   strconv.Itoa(item.Id),
			Timestamp: item.CreatedAt.UnixMilli(),
			RespType:  client.OrderRspType_RESULT,
			Symbol:    item.Market,
			Status:    spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(item.Status), item.FilledSize, item.Size),
			Price:     item.AvgFillPrice,
			Executed:  item.FilledSize,
			CloseDate: closeData,
		})
	}
	return resp, nil
} //查询订单信息。不支持ws回报订单行情的交易所，如bitbank, 需要频发查单获知订单状态
// 提币
func (c *ClientFTX) Transfer(o *order.OrderTransfer) (*client.OrderRsp, error) {
	//默认都是从现货提币
	var (
		resp    = &client.OrderRsp{}
		err     error
		coin    = string(o.ExchangeToken)
		address = string(o.TransferAddress) //提币地址
		tag     = string(o.Tag)             //某些币种例如 XRP,XMR 允许填写次级地址标签
		amount  = o.Amount                  //数量
		params  = make(map[string]interface{})
	)

	params["tag"] = tag
	switch o.Chain {
	case common.Chain_POLYGON:
		params["method"] = "matic"
	case common.Chain_ETH:
		params["method"] = "erc20"
	case common.Chain_TRON:
		params["method"] = "trx"
	case common.Chain_AVALANCHE:
		params["method"] = "avax"
	case common.Chain_FANTOM:
		params["method"] = "ftm"
	case common.Chain_BSC:
		params["method"] = "bsc"
	case common.Chain_SOLANA:
		params["method"] = "sol"
	default:
		params["method"] = ""
	}

	res, err := c.api.RequestWithdrawal(coin, amount, address, &params)
	if err != nil {
		return nil, err
	}
	resp.OrderId = strconv.Itoa(res.Result.Id)
	resp.Timestamp = res.Result.Time.UnixMilli()
	resp.Status = order.OrderStatusCode_PROCESSING
	resp.RespType = client.OrderRspType_ACK
	return resp, err
} //转账订单
func (c *ClientFTX) GetTransferHistory(req *client.TransferHistoryReq) (*client.TransferHistoryRsp, error) {
	var (
		resp      = &client.TransferHistoryRsp{}
		err       error
		startTime = req.StartTime //startTime	LONG	NO	默认当前时间90天前的时间戳
		endTime   = req.EndTime   //endTime	LONG	NO	默认当前时间戳
		params    = make(map[string]interface{})
		res       *spot_api.RespGetWithdrawalHistory
	)

	if startTime > 0 {
		params["start_time"] = time.UnixMilli(startTime).Unix()
	}
	if endTime > 0 {
		params["end_time"] = time.UnixMilli(endTime).Unix()
	}
	res, err = c.api.GetWithdrawalHistory(&params)
	if err != nil {
		return nil, err
	}

	for _, item := range res.Result {
		resp.TransferList = append(resp.TransferList, &client.TransferHistoryItem{
			Asset:     item.Coin,
			Amount:    item.Size,
			OrderId:   strconv.Itoa(item.Id),
			Network:   spot_api.GetChainFromNetWork(item.Method),
			Status:    spot_api.GetTransferStatusFromResponse(item.Status),
			Fee:       item.Fee,
			Address:   item.Address,
			TxId:      item.Txid,
			Timestamp: item.Time.UnixMilli(),
		})
	}
	return resp, err
} //查询提款记录
//划转

// Not an option in FTX (unless subaccounts)
func (c *ClientFTX) MoveAsset(o *order.OrderMove) (*client.OrderRsp, error) {

	var (
		res         *spot_api.RespTransferSubaccount
		resp        = &client.OrderRsp{}
		err         error
		coin        = o.Asset
		size        = o.Amount
		source      = o.AccountSource
		destination = o.AccountTarget
	)

	if o.ActionUser == order.OrderMoveUserType_Master {
		if o.AccountTarget == "" && o.AccountSource == "" {
			return nil, errors.New(fmt.Sprint("OrderMoveUserType_OrderMove_Sub not supply operation from:", o.AccountSource, " to:", o.AccountTarget))
		} else if o.AccountSource == "" {
			res, err = c.api.TransferSubAccount(coin, size, "main", destination)
		} else if o.AccountTarget == "" {
			res, err = c.api.TransferSubAccount(coin, size, source, "main")
		} else {
			res, err = c.api.TransferSubAccount(coin, size, source, destination)
		}
	} else if o.ActionUser == order.OrderMoveUserType_Sub {
		if o.AccountTarget == "" {
			res, err = c.api.TransferSubAccount(coin, size, source, "main")
		} else {
			res, err = c.api.TransferSubAccount(coin, size, source, destination)
		}
	} else { //Internal case, return nothing
		return nil, errors.New("not support move on ftx")
	}

	if err != nil {
		return nil, err
	}

	resp.OrderId = strconv.Itoa(res.Result.Id)
	resp.RespType = client.OrderRspType_ACK
	resp.Timestamp = res.Result.Time.UnixMicro()
	return resp, err
} //资产划转订单
// Not an option
func (c *ClientFTX) GetMoveHistory(req *client.MoveHistoryReq) (*client.MoveHistoryRsp, error) {
	var (
		resp        = &client.MoveHistoryRsp{}
		err         error
		startTime   = req.StartTime //startTime	LONG	NO
		endTime     = req.EndTime   //endTime	LONG	NO
		params      = make(map[string]interface{})
		source      = req.AccountSource
		destination = req.AccountTarget
	)

	if startTime > 0 {
		params["start_time"] = time.UnixMilli(startTime).Unix()
	}
	if endTime > 0 {
		params["end_time"] = time.UnixMilli(endTime).Unix()
	}

	if req.ActionUser == order.OrderMoveUserType_Master {
		// 获取主账户信息
		params["_subAccount"] = ""
		if req.AccountTarget == "" && req.AccountSource == "" {
			return nil, errors.New(fmt.Sprint("OrderMoveUserType_Master not supply operation from:", req.AccountSource, " to:", req.AccountTarget))
		} else if req.AccountSource == "" {
			res, err := c.api.GetWithdrawalHistory(&params)
			if err != nil {
				return nil, err
			}
			for _, transaction := range res.Result {
				if transaction.Notes == fmt.Sprintf("Transfer from %v to %v", "main account", destination) {
					resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
						Asset:     transaction.Coin,
						Id:        strconv.Itoa(transaction.Id),
						Type:      client.MoveType_MOVETYPE_OUT,
						Amount:    transaction.Size,
						Timestamp: transaction.Time.UnixMilli(),
						Status:    client.MoveStatus_MOVESTATUS_CONFIRMED,
					})
				}
			}
		} else if req.AccountTarget == "" {
			res, err := c.api.GetDepositHistory(&params)
			if err != nil {
				return nil, err
			}
			for _, transaction := range res.Result {
				if transaction.Notes == fmt.Sprintf("Transfer from %v to %v", source, "main account") {
					resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
						Asset:     transaction.Coin,
						Id:        strconv.Itoa(transaction.Id),
						Type:      client.MoveType_MOVETYPE_IN,
						Amount:    transaction.Size,
						Timestamp: transaction.Time.UnixMilli(),
						Status:    client.MoveStatus_MOVESTATUS_CONFIRMED,
					})
				}
			}
		} else {
			return nil, errors.New(fmt.Sprint("OrderMoveUserType_Master not supply operation from:", req.AccountSource, " to:", req.AccountTarget))
		}
	} else if req.ActionUser == order.OrderMoveUserType_Sub {
		// 获取指定的子账户信息
		params["_subAccount"] = source
		res, err := c.api.GetWithdrawalHistory(&params)
		if err != nil {
			return nil, err
		}
		for _, transaction := range res.Result {
			if transaction.Notes == fmt.Sprintf("Transfer from %v to %v", source, destination) {
				resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
					Asset:     transaction.Coin,
					Id:        strconv.Itoa(transaction.Id),
					Type:      client.MoveType_MOVETYPE_OUT,
					Amount:    transaction.Size,
					Timestamp: transaction.Time.UnixMilli(),
					Status:    client.MoveStatus_MOVESTATUS_CONFIRMED,
				})
			}
		}
	} else { //Internal case, return nothing
		return nil, errors.New("not support move on ftx")
	}

	return resp, err
} //查询划转记录
//充值

func (c *ClientFTX) GetDepositHistory(req *client.DepositHistoryReq) (*client.DepositHistoryRsp, error) {
	var (
		resp      = &client.DepositHistoryRsp{}
		err       error
		startTime = req.StartTime //startTime	LONG	NO	默认当前时间90天前的时间戳
		endTime   = req.EndTime   //endTime	LONG	NO	默认当前时间戳

		params = make(map[string]interface{})
		res    *spot_api.RespGetDepositHistory
	)

	if startTime > 0 {
		params["start_time"] = int(startTime)
	}
	if endTime > 0 {
		params["end_time"] = int(endTime)
	}
	res, err = c.api.GetDepositHistory(&params)
	if err != nil {
		return nil, err
	}

	for _, transaction := range res.Result {
		if err != nil {
			return nil, err
		}
		resp.DepositList = append(resp.DepositList, &client.DepositHistoryItem{
			Asset:      transaction.Coin,
			Amount:     transaction.Size,
			Network:    common.Chain_INVALID_CAHIN, //There is no method in deposit history
			Status:     spot_api.GetDepositStatusFromResponse(transaction.Status),
			Address:    transaction.Address.Address,
			AddressTag: transaction.Address.Tag,
			TxId:       transaction.Txid,
			Timestamp:  transaction.Time.UnixMilli(),
		})
	}
	return resp, err
} //查询存款记录

// 借贷还款 (FTX Not available)
func (c *ClientFTX) Loan(*order.OrderLoan) (*client.OrderRsp, error) {
	return nil, nil
} //借贷
func (c *ClientFTX) GetLoanOrders(*client.LoanHistoryReq) (*client.LoanHistoryRsp, error) { //get borrow history
	return nil, nil
} //获取已放款订单
func (c *ClientFTX) Return(*order.OrderReturn) (*client.OrderRsp, error) {
	return nil, nil
} //还币
func (c *ClientFTX) GetReturnOrders(*client.LoanHistoryReq) (*client.ReturnHistoryRsp, error) {
	return nil, nil
} //获取已放款订单

func (c *ClientFTX) ParseOrder(orders [][]float64, slice *[]*depth.DepthLevel) error {
	for _, order := range orders {
		*slice = append(*slice, &depth.DepthLevel{
			Price:  order[0],
			Amount: order[1],
		})
	}
	return nil
}

func (c *ClientFTX) GetFutureSymbols(market common.Market) []*client.SymbolInfo {
	var (
		spotSymbols   []*client.SymbolInfo
		futureSymbols []*client.SymbolInfo
		swapSymbols   []*client.SymbolInfo
		marketType    common.SymbolType
		err           error
	)

	exchangeInfoRes, err := c.api.GetMarkets()

	if err != nil {
		logger.Logger.Error("get exchange in error:", err)
		return futureSymbols
	}

	for _, next := range exchangeInfoRes.Result {
		switch next.Type {
		case "spot":
			spotSymbols = append(spotSymbols, &client.SymbolInfo{
				Symbol: next.Name,
				Name:   next.Name,
				Type:   common.SymbolType_SPOT_NORMAL,
			})
		case "future":
			marketType = spot_api.GetSymbolType(next.Name)
			switch marketType {
			case common.SymbolType_SWAP_FOREVER:
				swapSymbols = append(swapSymbols, &client.SymbolInfo{
					Symbol: strings.Split(next.Name, "-")[0] + "/USD",
					Name:   next.Name,
					Type:   marketType,
				})
			default:
				futureSymbols = append(futureSymbols, &client.SymbolInfo{
					Symbol: strings.Split(next.Name, "-")[0] + "/USD",
					Name:   next.Name,
					Type:   marketType,
				})
			}
		}
	}

	if market == common.Market_SPOT {
		return spotSymbols
	} else if market == common.Market_FUTURE {
		return futureSymbols
	} else if market == common.Market_SWAP {
		return swapSymbols
	} else {
		fmt.Println("Error: Invalid FTX Market")
		return nil
	}
}
func (c *ClientFTX) GetFutureDepth(symbol *client.SymbolInfo, limit int) (*depth.Depth, error) {
	apiSymbol := future_api.GetBaseSymbol(symbol.Symbol, symbol.Type)
	resDepth, err := c.api.GetOrderbook(apiSymbol, limit)
	if err != nil {
		//utils.Logger.Error(d.Symbol, d.getSymbolName(), "get full depth err", err)
		return nil, err
	}
	dep := &depth.Depth{
		Exchange:    common.Exchange_FTX,
		Market:      spot_api.GetMarket(apiSymbol),
		Type:        spot_api.GetSymbolType(apiSymbol),
		Symbol:      symbol.Symbol,
		TimeReceive: uint64(time.Now().UnixMicro()),
		TimeOperate: uint64(time.Now().UnixMicro()),
	}
	err = c.ParseOrder(resDepth.Result.Bids, &dep.Bids)
	if err != nil {
		return dep, err
	}
	err = c.ParseOrder(resDepth.Result.Asks, &dep.Asks)
	if err != nil {
		return dep, err
	}
	return dep, err
} //获取行情

func (c *ClientFTX) GetFutureMarkPrice(market common.Market, symbols ...*client.SymbolInfo) (*client.RspMarkPrice, error) { //标记价格
	var (
		err error
		res = &client.RspMarkPrice{}
	)

	if len(symbols) == 0 {
		resMarkPrice, err := c.uapi.ListAllFutures()
		if err != nil {
			fmt.Println("Uapi get error")
			return nil, err
		}
		for _, futureItem := range resMarkPrice.Result {
			res.Item = append(res.Item, &client.MarkPriceItem{
				Symbol:     futureItem.Name,
				Type:       spot_api.GetSymbolType(futureItem.Name),
				UpdateTime: time.Now().UnixMicro(), //There is no update time listed
				MarkPrice:  futureItem.Mark,
			})
		}
		return res, err
	}

	for _, symbol := range symbols {
		resMarkPrice, err := c.uapi.GetFuture(future_api.GetBaseSymbol(symbol.Symbol, symbol.Type))
		if err != nil {
			fmt.Println("Uapi get error")
			return nil, err
		}
		res.Item = append(res.Item, &client.MarkPriceItem{
			Symbol:     resMarkPrice.Result.Name,
			Type:       spot_api.GetSymbolType(resMarkPrice.Result.Name),
			UpdateTime: time.Now().UnixMicro(), //There is no update time listed
			MarkPrice:  resMarkPrice.Result.Mark,
		})
	}

	return res, err
}

func (c *ClientFTX) GetFutureTradeFee(market common.Market, symbols ...*client.SymbolInfo) (*client.TradeFee, error) { //查询交易手续费
	var (
		res = &client.TradeFee{}
	)

	tradeFees, err := c.api.GetAccountInfo()
	if err != nil {
		return nil, err
	}
	makerFee := tradeFees.Result.MakerFee
	takerFee := tradeFees.Result.TakerFee

	if len(symbols) == 0 {
		markets, err := c.api.GetMarkets()
		if err != nil {
			return nil, err
		}
		for _, next := range markets.Result {
			if spot_api.GetMarket(next.Name) == market {
				symbols = append(symbols, &client.SymbolInfo{
					Symbol: next.Name, //Do we keep ETH-PERP or do we need ETH/USD here?
					Name:   next.Name,
					Type:   spot_api.GetSymbolType(next.Name),
				})
			}
		}
	}

	for _, symbol := range symbols {
		res.TradeFeeList = append(res.TradeFeeList, &client.TradeFeeItem{
			Symbol: symbol.Symbol,
			Type:   symbol.Type,
			Maker:  makerFee,
			Taker:  takerFee,
		})
	}
	return res, nil
}

func (c *ClientFTX) GetFuturePrecision(market common.Market, symbols ...*client.SymbolInfo) (*client.Precision, error) {
	var (
		res             = &client.Precision{}
		searchSymbolMap = make(map[string]*client.SymbolInfo)
		exchangeInfoRes *spot_api.RespGetMarkets
		err             error
	)

	for _, symbol := range symbols {
		searchSymbolMap[future_api.GetBaseSymbol(symbol.Symbol, symbol.Type)] = symbol
	}
	exchangeInfoRes, err = c.api.GetMarkets()
	if err != nil {
		return nil, err
	}
	for _, symbol := range exchangeInfoRes.Result {
		if len(searchSymbolMap) > 0 {
			_, ok := searchSymbolMap[symbol.Name]
			if !ok {
				continue
			}
		}
		var (
			price, amount int
		)
		if strings.Contains(fmt.Sprintf("%v", symbol.PriceIncrement), ".") {
			price = len(strings.Split(fmt.Sprintf("%v", symbol.PriceIncrement), ".")[1]) // 下单价格精度
		} else {
			price = -len(strings.Split(fmt.Sprintf("%v", symbol.PriceIncrement), ".")[0]) // 下单价格精度
		}
		if strings.Contains(fmt.Sprintf("%v", symbol.SizeIncrement), ".") {
			amount = len(strings.Split(fmt.Sprintf("%v", symbol.SizeIncrement), ".")[1]) // 下单价格精度
		} else {
			amount = -len(strings.Split(fmt.Sprintf("%v", symbol.SizeIncrement), ".")[0]) // 下单价格精度
		}
		if err != nil {
			return nil, err
		}
		res.PrecisionList = append(res.PrecisionList, &client.PrecisionItem{
			Symbol:    symbol.Name,
			Type:      spot_api.GetSymbolType(symbol.Name),
			Amount:    int64(amount),
			Price:     int64(price),
			AmountMin: symbol.MinProvideSize,
		})
	}
	return res, nil
}

func (c *ClientFTX) GetFutureBalance(market common.Market) (*client.UBaseBalance, error) {
	var (
		respAccount   *spot_api.RespGetAccountInfo
		respPositions *spot_api.RespGetPositions
		res           = &client.UBaseBalance{}
		err           error
	)
	respAccount, err = c.uapi.GetAccountInfo()
	if err != nil {
		return res, err
	}

	res.UpdateTime = time.Now().UnixMicro()
	res.Market = common.Market_ALL_MARKET
	res.Balance = respAccount.Result.Collateral
	res.Rights = respAccount.Result.TotalAccountValue
	res.Used = respAccount.Result.Collateral - respAccount.Result.FreeCollateral
	res.Available = respAccount.Result.FreeCollateral
	res.TotalMarginBalance = 0 //TODO Fill with correct value

	respPositions, err = c.api.GetPositions()
	if err != nil {
		fmt.Println(err)
		return res, err
	}

	unrealizedPnl := 0.0

	for _, asset := range respPositions.Result {
		var closeDate string
		//TODO Check the math logic in the future
		//fmt.Printf("%#v\n", asset)
		names := strings.Split(asset.Future, "-")
		if len(names) > 1 {
			if names[1] != "PERP" && names[1] != "MOVE" {
				closeDate = names[1]
			}
		}
		res.UBasePositionList = append(res.UBasePositionList, &client.UBasePositionItem{
			Symbol:         asset.Future,
			Position:       asset.Size,
			Side:           spot_api.GetSideFromString(asset.Side),
			Market:         spot_api.GetMarket(asset.Future),
			Type:           spot_api.GetSymbolType(asset.Future),
			CloseDate:      closeDate,
			Price:          asset.EntryPrice,
			MaintainMargin: asset.RecentBreakEvenPrice * asset.Size * asset.MaintenanceMarginRequirement / respAccount.Result.Leverage,
			InitialMargin:  asset.RecentBreakEvenPrice * asset.Size / respAccount.Result.Leverage,
			Notional:       0, //TODO add in the future
			Leverage:       respAccount.Result.Leverage,
			Unprofit:       asset.UnrealizedPnl,
		})
		unrealizedPnl += asset.UnrealizedPnl
	}

	res.Unprofit = unrealizedPnl

	return res, err
} //Needs to be checked for math correctness

func (c *ClientFTX) PlaceFutureOrder(o *order.OrderTradeCEX) (*client.OrderRsp, error) {
	var (
		symbol  string
		side    spot_api.SideType
		type_   spot_api.OrderType
		size    float64
		options = make(map[string]interface{})
		resp    = &client.OrderRsp{}
		err     error
	)

	side = spot_api.GetSideTypeToExchange(o.Side)

	switch o.Tif {
	case order.TimeInForce_IOC:
		options["ioc"] = true
	default:
		options["ioc"] = false
	}

	if o.Base.Id != 0 {
		options["clientId"] = transform.IdToClientId(o.Hdr.Producer, o.Base.Id)
	}
	if o.OrderType == order.OrderType_MARKET {
		options["price"] = nil
	} else if o.OrderType == order.OrderType_LIMIT {
		options["price"] = o.Price
	}
	if o.TradeType == order.TradeType_MAKER {
		options["postOnly"] = true
	}

	size = o.Amount
	type_ = spot_api.GetOrderTypeToExchange(o.OrderType)
	symbol = string(o.Base.Symbol)

	var (
		res *spot_api.RespPlaceOrder
	)
	if o.Base.Market == common.Market_FUTURE || o.Base.Market == common.Market_SWAP {
		res, err = c.api.PlaceOrder(symbol, side, type_, size, &options)
	}
	if err != nil {
		return nil, err
	}
	var closeData string
	names := strings.Split(res.Result.Market, "-")
	if len(names) > 1 {
		if names[1] != "PERP" && names[1] != "MOVE" {
			closeData = names[1]
		}
	}
	resp.Producer, resp.Id = transform.ClientIdToId(res.Result.ClientId)
	resp.OrderId = strconv.Itoa(res.Result.Id)
	resp.Timestamp = res.Result.CreatedAt.UnixMilli()
	resp.RespType = client.OrderRspType_FULL
	resp.Symbol = res.Result.Market //Remember to change to /
	resp.Status = spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(res.Result.Status), res.Result.FilledSize, res.Result.Size)
	resp.Price = res.Result.Price
	resp.Executed = res.Result.FilledSize
	resp.CloseDate = closeData

	return resp, nil
}
