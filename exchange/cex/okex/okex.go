package okex

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/binance/u_api"
	"clients/exchange/cex/okex/ok_api"
	"clients/logger"
	"clients/transform"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
)

type ClientOkex struct {
	api *ok_api.ClientOkex

	tradeFeeMap    map[string]*client.TradeFeeItem    //key: symbol
	transferFeeMap map[string]*client.TransferFeeItem //key: network+token
	precisionMap   map[string]*client.PrecisionItem   //key: symbol

	optionMap map[string]interface{}
}

func NewClientOkex(conf base.APIConf, maps ...interface{}) *ClientOkex {
	c := &ClientOkex{
		api: ok_api.NewClientOkex(conf),
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

func NewClientOkex2(conf base.APIConf, cli *http.Client, maps ...interface{}) *ClientOkex {
	// 使用自定义http client
	c := &ClientOkex{
		api: ok_api.NewClientOkex2(conf, cli),
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

func (c *ClientOkex) GetExchange() common.Exchange {
	return common.Exchange_OKEX
}

func (c *ClientOkex) GetSymbols() []string {
	var symbols []string
	exchangeInfoRes, err := c.api.Instrument_Info("SPOT", nil)
	if err != nil {
		logger.Logger.Error("get exchange in error:", err)
		return symbols
	}
	for _, symbol := range exchangeInfoRes.Data {
		if symbol.State == "live" {
			symbols = append(symbols, symbol.BaseCcy+"/"+symbol.QuoteCcy)
		}
	}
	return symbols
}

func (c *ClientOkex) GetFutureSymbols(market common.Market) []*client.SymbolInfo {
	var (
		type_   common.SymbolType
		name    string
		symbol_ string
	)
	res := make([]*client.SymbolInfo, 0)
	set := make(map[string]struct{})
	exchangeInfoRes, err := c.api.Instrument_Info("FUTURES", nil)
	exchangeInfoRes1, err := c.api.Instrument_Info("SWAP", nil)
	exchangeInfoRes.Data = append(exchangeInfoRes.Data, exchangeInfoRes1.Data...)
	if err != nil {
		logger.Logger.Error("get exchange in error:", err)
		return res
	}
	for _, symbol := range exchangeInfoRes.Data {
		if symbol.State == "live" {
			symbol_, name, type_ = ParseInstId(symbol.InstId)
			if _, ok := set[symbol_]; !ok {
				res = append(res, &client.SymbolInfo{
					Symbol: name,
					Name:   symbol_,
					Type:   type_,
				})
				set[symbol_] = struct{}{}
			}
		}
	}
	return res
}

// 用于解算instid中的类型
// name:币对
// symbol_ instId转换
func ParseInstId(instId string) (symbol_, name string, type_ common.SymbolType) {
	var u bool
	name = strings.ReplaceAll(instId, "-", "/")
	names := strings.Split(name, "/")
	symbol_ = names[0] + "/" + names[1]
	if names[1] == "USDT" {
		u = true
	} else {
		u = false
	}
	type_ = ok_api.ParseDate(names[2], u)
	return
}

func (c *ClientOkex) GetUBaseSymbols() []*client.SymbolInfo {
	res := make([]*client.SymbolInfo, 0)
	set := make(map[string]struct{})
	exchangeInfoRes, err := c.api.Instrument_Info("FUTURES", nil)
	if err != nil {
		logger.Logger.Error("get exchange in error:", err)
		return res
	}
	for _, symbol := range exchangeInfoRes.Data {
		if symbol.State == "live" {
			symbol_ := strings.ReplaceAll(symbol.InstId, "-", "/")
			if _, ok := set[symbol_]; !ok {

				res = append(res, &client.SymbolInfo{
					Symbol: symbol_})
				set[symbol_] = struct{}{}
			}
		}
	}
	return res
}

func (c *ClientOkex) ParseOrder(orders [][]string, slice *[]*depth.DepthLevel) error {
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

func (c *ClientOkex) GetDepth(symbol *client.SymbolInfo, limit int) (*depth.Depth, error) {
	sym := GetInstId(symbol)
	params := url.Values{}
	params.Add("sz", fmt.Sprint(limit))
	repDepth, err := c.api.Market_Books_Info(sym, &params)
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
	err = c.ParseOrder(repDepth.Data[0].Bids, &dep.Bids)
	if err != nil {
		return dep, err
	}

	err = c.ParseOrder(repDepth.Data[0].Asks, &dep.Asks)
	if err != nil {
		return dep, err
	}
	intNum, _ := strconv.Atoi(repDepth.Data[0].Ts)
	dep.TimeReceive = uint64(intNum)
	return dep, err
}

func (c *ClientOkex) GetUBaseDepth(symbol *client.SymbolInfo, limit int) (*depth.Depth, error) {
	return c.GetDepth(symbol, limit)
}

func (c *ClientOkex) GetFutureDepth(symbol *client.SymbolInfo, limit int) (*depth.Depth, error) {
	return c.GetDepth(symbol, limit)
}

func (c *ClientOkex) IsExchangeEnable() bool {
	res, err := c.api.Public_ServerTime_Info(nil)
	if err != nil {
		logger.Logger.Error(err)
	}
	serverTime, err := strconv.ParseInt(res.Data[0].Ts, 10, 64)
	return err == nil && time.Since(time.UnixMilli(serverTime)) < time.Second*60
}

func GetSymbol(symbol string) string {
	return strings.ToUpper(strings.ReplaceAll(symbol, "/", "-"))
}

func (c *ClientOkex) GetTradeFee(symbols ...string) (*client.TradeFee, error) { //查询交易手续费
	var (
		res             = &client.TradeFee{}
		searchSymbolMap = make(map[string]bool)
	)
	for _, symbol := range symbols {
		searchSymbolMap[GetSymbol(symbol)] = true
	}
	tradeFeeRes, err := c.api.Account_TradeFee_Info("SPOT", nil)
	if err != nil {
		return nil, err
	}
	if len(tradeFeeRes.Data) != 1 {
		err = errors.New("network error")
		if err != nil {
			return nil, err
		}
	}
	for _, symbol := range symbols {
		var (
			takerFee, makerFee float64
		)
		symbol = GetSymbol(symbol)
		takerFee, err = strconv.ParseFloat(tradeFeeRes.Data[0].Taker, 64)
		if err != nil {
			logger.Logger.Error("convert taker fee err:", err, tradeFeeRes.Data[0].Taker)
			continue
		}
		makerFee, err = strconv.ParseFloat(tradeFeeRes.Data[0].Maker, 64)
		if err != nil {
			logger.Logger.Error("convert maker fee err:", tradeFeeRes.Data[0].Maker)
			continue
		}
		res.TradeFeeList = append(res.TradeFeeList, &client.TradeFeeItem{
			Symbol: ok_api.ParseSymbolName(symbol),
			Maker:  makerFee,
			Taker:  takerFee,
		})
	}
	return res, nil
}

func (c *ClientOkex) GetUBaseTradeFee(symbols ...*client.SymbolInfo) (*client.TradeFee, error) { //查询交易手续费
	//symbols为空则返回常规
	var (
		res             = &client.TradeFee{}
		searchSymbolMap = make(map[string]bool)
	)
	if symbols == nil || len(symbols) <= 0 {
		tmp := c.GetUBaseSymbols()
		for _, symbol := range tmp {
			searchSymbolMap[symbol.Symbol] = true
		}
	} else {
		for _, symbol := range symbols {
			symbol_ := GetInstId(symbol)
			searchSymbolMap[strings.ReplaceAll(symbol_, "-", "/")] = true
		}
	}

	tradeFeeRes, err := c.api.Account_TradeFee_Info("FUTURES", nil)
	if err != nil {
		return nil, err
	}
	if len(tradeFeeRes.Data) != 1 {
		err = errors.New("network error")
		if err != nil {
			return nil, err
		}
	}
	for key, _ := range searchSymbolMap {
		var (
			takerFee, makerFee float64
		)
		takerFee, err = strconv.ParseFloat(tradeFeeRes.Data[0].Taker, 64)
		if err != nil {
			logger.Logger.Error("convert taker fee err:", err, tradeFeeRes.Data[0].Taker)
			continue
		}
		makerFee, err = strconv.ParseFloat(tradeFeeRes.Data[0].Maker, 64)
		if err != nil {
			logger.Logger.Error("convert maker fee err:", tradeFeeRes.Data[0].Maker)
			continue
		}
		res.TradeFeeList = append(res.TradeFeeList, &client.TradeFeeItem{
			Symbol: key,
			Maker:  makerFee,
			Taker:  takerFee,
		})
	}
	return res, nil
}

func (c *ClientOkex) GetFutureTradeFee(market common.Market, symbols ...*client.SymbolInfo) (*client.TradeFee, error) { //查询交易手续费
	var (
		res         = &client.TradeFee{}
		symbol      string
		tradeFeeRes *ok_api.Resp_Account_TradeFee
		err         error
		instId      string
		type_       common.SymbolType
	)
	if symbols != nil && symbols[0].Symbol != "" {
		instId = GetInstId(symbols[0])
		_, symbol, type_ = ParseInstId(instId)
		if type_ == common.SymbolType_SWAP_FOREVER || type_ == common.SymbolType_SWAP_COIN_FOREVER {
			tradeFeeRes, err = c.api.Account_TradeFee_Info("SWAP", nil)
		} else {
			tradeFeeRes, err = c.api.Account_TradeFee_Info("FUTURES", nil)
		}
	} else {
		if market == common.Market_FUTURE_COIN || market == common.Market_FUTURE {
			tradeFeeRes, err = c.api.Account_TradeFee_Info("FUTURES", nil)
		} else {
			tradeFeeRes, err = c.api.Account_TradeFee_Info("SWAP", nil)
		}
		for _, instrument := range tradeFeeRes.Data {
			a, _ := json.Marshal(instrument)
			fmt.Println(string(a))
		}
	}

	if err != nil {
		return nil, err
	}
	if len(tradeFeeRes.Data) != 1 {
		err = errors.New("network error")
		if err != nil {
			return nil, err
		}
	}
	for _, v := range tradeFeeRes.Data {
		var (
			takerFee, makerFee float64
		)
		takerFee, err = strconv.ParseFloat(v.Taker, 64)
		if err != nil {
			logger.Logger.Error("convert taker fee err:", err, tradeFeeRes.Data[0].Taker)
			continue
		}
		makerFee, err = strconv.ParseFloat(v.Maker, 64)
		if err != nil {
			logger.Logger.Error("convert maker fee err:", tradeFeeRes.Data[0].Maker)
			continue
		}
		res.TradeFeeList = append(res.TradeFeeList, &client.TradeFeeItem{
			Symbol: symbol,
			Maker:  makerFee,
			Taker:  takerFee,
		})
	}
	return res, nil
}

func (c *ClientOkex) GetTransferFee(chain common.Chain, tokens ...string) (*client.TransferFee, error) {
	var (
		res            = &client.TransferFee{}
		searchTokenMap = make(map[string]bool)
	)
	for _, token := range tokens {
		searchTokenMap[token] = true
	}
	configAll, err := c.api.Asset_Currencies_Info(nil)
	if err != nil {
		return nil, err
	}
	for _, item := range configAll.Data {
		if len(searchTokenMap) > 0 {
			_, ok := searchTokenMap[item.Ccy]
			if !ok || ok_api.GetChainFromNetWork(item.Chain) != chain {
				continue
			}
		}
		var fee float64
		//获取接口中的手续费值
		fee, err = strconv.ParseFloat(item.MinFee, 64)
		if err != nil {
			logger.Logger.Info("convert withdraw fee err:", item.Ccy, item.Chain, item.Ccy)
			continue
		}
		res.TransferFeeList = append(res.TransferFeeList, &client.TransferFeeItem{
			Token:   item.Ccy,
			Network: ok_api.GetChainFromNetWork(item.Chain),
			Fee:     fee,
		})
	}
	return res, nil
}

// precision定义
// AmountPrecision PricePrecision
func (c *ClientOkex) GetPrecision(symbols ...string) (*client.Precision, error) {
	var (
		res             = &client.Precision{}
		searchSymbolMap = make(map[string]bool)
	)
	for _, symbol := range symbols {
		searchSymbolMap[GetSymbol(symbol)] = true
	}
	exchangeInfoRes, err := c.api.Instrument_Info("SPOT", nil)
	if err != nil {
		return nil, err
	}
	for _, symbol := range exchangeInfoRes.Data {
		//fmt.Println(111, symbol, symbol.BaseCcy+"/"+symbol.QuoteCcy)
		if len(searchSymbolMap) > 0 {
			_, ok := searchSymbolMap[symbol.InstId]
			if !ok {
				continue
			}
		}
		var (
			price, amount int
			amountMin     float64
		)
		if strings.Contains(symbol.TickSz, ".") {
			price = len(strings.Split(symbol.TickSz, ".")[1]) // 下单价格精度
		} else {
			price = 1 - len(strings.Split(symbol.TickSz, ".")[0]) // 下单价格精度
		}
		if strings.Contains(symbol.LotSz, ".") {
			amount = len(strings.Split(symbol.LotSz, ".")[1]) // 下单价格精度
		} else {
			amount = 1 - len(strings.Split(symbol.LotSz, ".")[0]) // 下单价格精度
		}
		amountMin, err = strconv.ParseFloat(symbol.MinSz, 64)
		if err != nil {
			return nil, err
		}
		res.PrecisionList = append(res.PrecisionList, &client.PrecisionItem{
			Symbol:    symbol.BaseCcy + "/" + symbol.QuoteCcy,
			Type:      common.SymbolType_SPOT_NORMAL,
			Amount:    int64(amount),
			Price:     int64(price),
			AmountMin: amountMin,
		})
	}
	return res, nil
}

func (c *ClientOkex) GetUBasePrecision(symbols ...*client.SymbolInfo) (*client.Precision, error) {
	var (
		res             = &client.Precision{}
		searchSymbolMap = make(map[string]bool)
	)
	if symbols == nil || len(symbols) <= 0 {
		tmp := c.GetUBaseSymbols()
		for _, symbol := range tmp {
			searchSymbolMap[symbol.Symbol] = true
		}
	} else {
		for _, symbol := range symbols {
			symbol_ := GetInstId(symbol)
			searchSymbolMap[strings.ReplaceAll(symbol_, "-", "/")] = true
		}
	}
	exchangeInfoRes, err := c.api.Instrument_Info("FUTURES", nil)
	if err != nil {
		return nil, err
	}
	for _, symbol := range exchangeInfoRes.Data {
		//fmt.Println(111, symbol, symbol.BaseCcy+"/"+symbol.QuoteCcy)
		instId := strings.ReplaceAll(symbol.InstId, "-", "/")

		if len(searchSymbolMap) > 0 {
			_, ok := searchSymbolMap[instId]
			if !ok {
				continue
			}
		}
		var (
			price, amount int
			amountMin     float64
		)
		if strings.Contains(symbol.TickSz, ".") {
			price = len(strings.Split(symbol.TickSz, ".")[1]) // 下单价格精度
		} else {
			price = 1 - len(strings.Split(symbol.TickSz, ".")[0]) // 下单价格精度
		}
		if strings.Contains(symbol.LotSz, ".") {
			amount = len(strings.Split(symbol.LotSz, ".")[1]) // 下单价格精度
		} else {
			amount = 1 - len(strings.Split(symbol.LotSz, ".")[0]) // 下单价格精度
		}
		amountMin, err = strconv.ParseFloat(symbol.MinSz, 64)
		if err != nil {
			return nil, err
		}
		res.PrecisionList = append(res.PrecisionList, &client.PrecisionItem{
			Symbol:    instId,
			Type:      common.SymbolType_SPOT_NORMAL,
			Amount:    int64(amount),
			Price:     int64(price),
			AmountMin: amountMin,
		})
	}
	return res, nil
}

func (c *ClientOkex) GetFuturePrecision(market common.Market, symbols ...*client.SymbolInfo) (*client.Precision, error) {
	var (
		exchangeInfoRes *ok_api.RespInstruments
		err             error
		res             = &client.Precision{}
		searchSymbolMap = make(map[string]bool)
	)
	if symbols == nil || len(symbols) == 0 || symbols[0] == nil || symbols[0].Symbol == "" {
		if market == common.Market_FUTURE_COIN || market == common.Market_FUTURE {
			exchangeInfoRes, err = c.api.Instrument_Info("FUTURES", nil)
		} else {
			exchangeInfoRes, err = c.api.Instrument_Info("SWAP", nil)
		}
	} else {
		instId := GetInstId(symbols...)
		_, _, type_ := ParseInstId(instId)
		params := url.Values{}
		params.Add("instId", instId)
		if type_ == common.SymbolType_SWAP_FOREVER || type_ == common.SymbolType_SWAP_COIN_FOREVER {
			exchangeInfoRes, err = c.api.Instrument_Info("SWAP", &params)
		} else {
			exchangeInfoRes, err = c.api.Instrument_Info("FUTURES", &params)
		}
	}
	if err != nil {
		return nil, err
	}
	for _, symbol := range exchangeInfoRes.Data {
		//fmt.Println(111, symbol, symbol.BaseCcy+"/"+symbol.QuoteCcy)
		instId := strings.ReplaceAll(symbol.InstId, "-", "/")

		if len(searchSymbolMap) > 0 {
			_, ok := searchSymbolMap[instId]
			if !ok {
				continue
			}
		}
		var (
			price, amount int
			amountMin     float64
		)
		if strings.Contains(symbol.TickSz, ".") {
			price = len(strings.Split(symbol.TickSz, ".")[1]) // 下单价格精度
		} else {
			price = 1 - len(strings.Split(symbol.TickSz, ".")[0]) // 下单价格精度
		}
		if strings.Contains(symbol.LotSz, ".") {
			amount = len(strings.Split(symbol.LotSz, ".")[1]) // 下单价格精度
		} else {
			amount = 1 - len(strings.Split(symbol.LotSz, ".")[0]) // 下单价格精度
		}
		amountMin, err = strconv.ParseFloat(symbol.MinSz, 64)
		if err != nil {
			return nil, err
		}
		res.PrecisionList = append(res.PrecisionList, &client.PrecisionItem{
			Symbol:    instId,
			Type:      common.SymbolType_SPOT_NORMAL,
			Amount:    int64(amount),
			Price:     int64(price),
			AmountMin: amountMin,
		})
	}
	return res, nil
}

func (c *ClientOkex) GetFutureMarkPrice(market common.Market, symbol ...*client.SymbolInfo) (*client.RspMarkPrice, error) { //标记价格
	var (
		rep          *ok_api.Resp_MarkPrice
		err          error
		value        = &url.Values{}
		instId, name string
		res          = &client.RspMarkPrice{}
		type_        common.SymbolType
	)
	if symbol != nil && symbol[0].Symbol != "" {
		instId = GetInstId(symbol[0])
		value.Add("instId", instId)
		_, _, type_ = ParseInstId(instId)
		if type_ == common.SymbolType_SWAP_FOREVER || type_ == common.SymbolType_SWAP_COIN_FOREVER {
			rep, err = c.api.Public_MarkPrice_Info("SWAP", value)
		} else {
			rep, err = c.api.Public_MarkPrice_Info("FUTURES", value)
		}
	} else {
		if market == common.Market_FUTURE_COIN || market == common.Market_FUTURE {
			rep, err = c.api.Public_MarkPrice_Info("FUTURES", nil)
		} else {
			rep, err = c.api.Public_MarkPrice_Info("SWAP", nil)
		}
	}
	if err != nil {
		return nil, err
	}
	for _, i := range rep.Data {
		_, name, type_ = ParseInstId(i.InstId)
		time, err := strconv.ParseInt(i.Ts, 10, 64)
		price, err := strconv.ParseFloat(i.MarkPx, 64)
		if err != nil {
			return nil, err
		}
		res.Item = append(res.Item, &client.MarkPriceItem{
			Symbol:     name,
			Type:       type_,
			UpdateTime: time * 1000,
			MarkPrice:  price,
		})
	}
	return res, err
}

// lock的字段
func (c *ClientOkex) GetBalance() (*client.SpotBalance, error) {
	var (
		respAsset *ok_api.Resp_Asset_Balances
		res       = new(client.SpotBalance)
		err       error
	)
	respAsset, err = c.api.Asset_Balances_Info(nil)
	if err != nil {
		return res, err
	}
	for _, balance := range respAsset.Data {
		res.WalletList = append(res.WalletList, &client.SpotBalanceItem{
			Asset:  balance.Ccy,
			Total:  transform.StringToX[float64](balance.Bal).(float64),
			Free:   transform.StringToX[float64](balance.AvailBal).(float64),
			Frozen: transform.StringToX[float64](balance.FrozenBal).(float64),
		})
	}
	return res, nil
}

func (c *ClientOkex) GetUBaseBalance() (*client.UBaseBalance, error) {

	var (
		respAccount *ok_api.Resp_Accout_Balance
		res         = &client.UBaseBalance{}
		err         error
	)
	respAccount, err = c.api.Account_Balance_Info(nil)
	if err != nil {
		return res, err
	}

	res.UpdateTime = time.Now().UnixMicro()
	for _, asset := range respAccount.Data[0].Details {
		var (
			walletBalance, crossUnPnl, availableBalance float64
		)
		if walletBalance, err = strconv.ParseFloat(asset.Eq, 64); err != nil {
			logger.Logger.Error("get free err", asset)
			return res, err
		}
		if crossUnPnl, err = strconv.ParseFloat(asset.IsoUpl, 64); err != nil {
			logger.Logger.Error("get free err", asset)
			return res, err
		}
		if availableBalance, err = strconv.ParseFloat(asset.AvailBal, 64); err != nil {
			logger.Logger.Error("get free err", asset)
			return res, err
		}
		res.UBaseBalanceList = append(res.UBaseBalanceList, &client.UBaseBalanceItem{
			Asset:     asset.Ccy,
			Balance:   walletBalance,
			Market:    common.Market_FUTURE,
			Unprofit:  crossUnPnl,
			Available: availableBalance,
			Used:      walletBalance + crossUnPnl - availableBalance,
		})
	}
	return res, err
}

func (c *ClientOkex) GetFutureBalance(market common.Market) (*client.UBaseBalance, error) {

	var (
		respAccountTotal *ok_api.Resp_Accout_Balance
		respAccount      *ok_api.Resp_Account_Position
		res              = &client.UBaseBalance{}
		err              error
	)
	params := url.Values{}

	params.Add("instType", "FUTURES")
	respAccount, err = c.api.Account_Positions_Info(&params)
	if err != nil {
		return res, err
	}
	params.Set("instType", "SWAP")
	respAccount1, err := c.api.Account_Positions_Info(&params)
	if err != nil {
		return res, err
	}
	respAccountTotal, err = c.api.Account_Balance_Info(nil)
	if err != nil {
		return nil, err
	}

	res.Rights = transform.StringToX[float64](respAccountTotal.Data[0].TotalEq).(float64)
	res.TotalMarginBalance = transform.StringToX[float64](respAccountTotal.Data[0].Imr).(float64)
	res.Available = transform.StringToX[float64](respAccountTotal.Data[0].AdjEq).(float64)
	//res.= transform.StringToX[float64](respAccountTotal.Data[0].TotalEq).(float64)
	res.UpdateTime = time.Now().UnixMicro()

	respAccount.Data = append(respAccount.Data, respAccount1.Data...)
	res.UpdateTime = time.Now().UnixMicro()
	for _, position := range respAccount.Data { //仓位
		symbol, _, type_ := ParseInstId(position.InstId)
		names := strings.Split(position.InstId, "-")
		var closeDate string
		if len(names) == 3 {
			if names[2] != "SWAP" {
				closeDate = names[2]
			}
		}
		side := order.TradeSide_BUY
		if strings.Contains(position.Pos, "-") {
			side = order.TradeSide_SELL
		}
		res.UBasePositionList = append(res.UBasePositionList, &client.UBasePositionItem{
			Symbol:         symbol,
			Type:           type_,
			CloseDate:      closeDate,
			MaintainMargin: transform.StringToX[float64](position.Mmr).(float64),
			InitialMargin:  transform.StringToX[float64](position.Imr).(float64),
			Notional:       math.Abs(transform.StringToX[float64](position.NotionalUsd).(float64)),
			Leverage:       transform.StringToX[float64](position.Lever).(float64),
			Position:       math.Abs(transform.StringToX[float64](position.Pos).(float64)),
			Side:           side,
			Unprofit:       transform.StringToX[float64](position.Upl).(float64),
		})
	}
	return res, err
}

func Type2Book(type_ common.SymbolType) common.Market {
	switch type_ {
	case common.SymbolType_SWAP_FOREVER:
		return common.Market_SWAP
	case common.SymbolType_SWAP_COIN_FOREVER:
		return common.Market_SWAP_COIN
	case common.SymbolType_FUTURE_THIS_WEEK:
		return common.Market_FUTURE
	case common.SymbolType_FUTURE_NEXT_WEEK:
		return common.Market_FUTURE
	case common.SymbolType_FUTURE_THIS_QUARTER:
		return common.Market_FUTURE
	case common.SymbolType_FUTURE_NEXT_QUARTER:
		return common.Market_FUTURE
	default:
		return common.Market_FUTURE_COIN
	}
}

// Margin资产解读
func (c *ClientOkex) getMarginAsset(balance *ok_api.MarginAssetItem) (*client.MarginBalanceItem, error) {
	var (
		free, frozen, borrowed, netAsset float64
		err                              error
	)
	// 传回来的数据包含""代表0,测试账号为空，处理为0看结果
	if balance.CashBal == "" {
		balance.CashBal = "0"
	}
	if balance.FrozenBal == "" {
		balance.FrozenBal = "0"
	}
	if balance.CrossLiab == "" {
		balance.CrossLiab = "0"
	}
	if balance.AvailEq == "" {
		balance.AvailEq = "0"
	}

	if free, err = strconv.ParseFloat(balance.CashBal, 64); err != nil {
		logger.Logger.Error("get free err", balance)
		return nil, err
	}
	if frozen, err = strconv.ParseFloat(balance.FrozenBal, 64); err != nil {
		logger.Logger.Error("get frozen err", balance)
		return nil, err
	}
	if borrowed, err = strconv.ParseFloat(balance.CrossLiab, 64); err != nil {
		logger.Logger.Error("get borrowed err", balance)
		return nil, err
	}
	if netAsset, err = strconv.ParseFloat(balance.AvailEq, 64); err != nil {
		logger.Logger.Error("get netAsset err", balance)
		return nil, err
	}
	// todo netasset是哪个值
	res := &client.MarginBalanceItem{Asset: balance.Ccy, Total: free + frozen, Borrowed: borrowed, Free: free, Frozen: frozen, NetAsset: netAsset}
	return res, err
}

// 没有值
func (c *ClientOkex) GetMarginBalance() (*client.MarginBalance, error) {
	var (
		respAccount                                    *ok_api.Resp_Accout_Balance
		res                                            = &client.MarginBalance{}
		marginLevel, totalNetAsset, totalLiablityAsset float64
		err                                            error
	)
	respAccount, err = c.api.Account_Balance_Info(nil)
	if err != nil {
		return res, err
	}
	// todo 去掉判空
	if respAccount.Data[0].AdjEq == "" {
		respAccount.Data[0].AdjEq = "0"
	}
	if totalNetAsset, err = strconv.ParseFloat(respAccount.Data[0].AdjEq, 64); err != nil {
		logger.Logger.Error("get free err", respAccount)
		return res, err
	}
	if respAccount.Data[0].NotionalUsd == "" {
		respAccount.Data[0].NotionalUsd = "0"
	}
	if totalLiablityAsset, err = strconv.ParseFloat(respAccount.Data[0].NotionalUsd, 64); err != nil {
		logger.Logger.Error("get free err", respAccount)
		return res, err
	}
	res.MarginLevel = marginLevel
	res.TotalNetAsset = totalNetAsset
	res.TotalLiabilityAsset = totalLiablityAsset
	res.QuoteAsset = "USD"

	for _, balance := range respAccount.Data[0].Details {
		var item *client.MarginBalanceItem
		item, err = c.getMarginAsset(balance)
		if err != nil {
			return nil, err
		}
		res.MarginBalanceList = append(res.MarginBalanceList, item)
	}
	return res, err
}

func (c *ClientOkex) GetMarginIsolatedBalance(symbols ...string) (*client.MarginIsolatedBalance, error) {
	return nil, nil
}

func (c *ClientOkex) PlaceOrder(o *order.OrderTradeCEX) (*client.OrderRsp, error) {
	var (
		symbol       string
		side         ok_api.SideType
		type_        ok_api.OrderType
		options      = &url.Values{}
		resp         = &client.OrderRsp{}
		precisionIns interface{}
		precision    *client.PrecisionItem
		tdMode       string
		ok           bool
		err          error
	)
	symbol = strings.ToUpper(strings.Replace(string(o.Base.Symbol), "/", "-", -1))
	if precisionIns, ok = c.precisionMap[string(o.Base.Symbol)]; !ok {
		return nil, errors.New("get precision err")
	}
	if precision, ok = precisionIns.(*client.PrecisionItem); !ok {
		return nil, errors.New("get precision err")
	}
	if o.Amount < precision.AmountMin {
		return nil, errors.New(fmt.Sprint("less amount in err:", o.Amount, "<", precision.AmountMin))
	}
	side = ok_api.GetSideTypeToExchange(o.Side)
	//switch o.Base.Type {
	//case common.SymbolType_MARGIN_NORMAL:
	//	isIsolated = spot_api.MARGIN_TYPE_NORMAL
	//case common.SymbolType_MARGIN_ISOLATED:
	//	isIsolated = spot_api.MARGIN_TYPE_ISOLATED
	//}
	if o.Base.Id != 0 {
		options.Add("clOrdId", transform.IdToClientId(o.Hdr.Producer, o.Base.Id))
	}
	// 交易类型
	//tdMode := "cash"

	if c.optionMap != nil {
		if mode, ok := c.optionMap["tdMode"]; ok {
			tdMode, _ = mode.(string)
		}
	}
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
		type_ = ok_api.ORDER_TYPE_POST_ONLY
		options.Add("px", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
	} else if o.Tif == order.TimeInForce_FOK {
		type_ = ok_api.ORDER_TYPE_FOK
		options.Add("px", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
	} else if o.Tif == order.TimeInForce_IOC {
		type_ = ok_api.ORDER_TYPE_IOC
		options.Add("px", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
	} else if o.Tif == order.TimeInForce_OPTIMAL_LIMIT_IOC {
		type_ = ok_api.ORDER_TYPE_OPTIMAL_LIMIT_IOC
	} else {
		type_ = ok_api.GetOrderTypeToExchange(o.OrderType)
		if type_ == ok_api.ORDER_TYPE_LIMIT {
			options.Add("px", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
			// todo 为什么？
		} else if type_ == ok_api.ORDER_TYPE_MARKET && side == ok_api.SIDE_TYPE_BUY {
			//amount = amount * o.Price
		}
	}
	var (
		res *ok_api.Resp_Trde_Order
	)
	//
	if o.Base.Market == common.Market_SPOT {
		res, err = c.api.Trade_Order_Info(symbol, tdMode, strings.ToLower(string(side)), strings.ToLower(string(type_)), strconv.FormatFloat(amount, 'f', int(precision.Amount), 64), options)
	} else if o.Base.Market == common.Market_MARGIN {
		res, err = c.api.Trade_Order_Info(symbol, tdMode, string(side), string(type_), strconv.FormatFloat(amount, 'f', int(precision.Amount), 64), options)
	}
	if err != nil {
		return nil, err
	}
	resp.Producer, resp.Id = transform.ClientIdToId(res.Data[0].ClOrdId)
	resp.OrderId = res.Data[0].OrdId
	resp.Symbol = symbol
	resp.Status = GetOrderStatusFromExchange(res.Data[0].SCode)
	return resp, nil
}

// Operation is not supported under the current account mode
func (c *ClientOkex) PlaceUBaseOrder(o *order.OrderTradeCEX) (*client.OrderRsp, error) {
	var (
		symbol_      string
		symbol       client.SymbolInfo
		instId       string
		side         ok_api.SideType
		type_        ok_api.OrderType
		options      = &url.Values{}
		resp         = &client.OrderRsp{}
		precisionIns interface{}
		precision    *client.PrecisionItem
		tdMode       string
		ok           bool
		err          error
	)
	symbol_ = strings.ToUpper(strings.Replace(string(o.Base.Symbol), "/", "-", -1))
	symbol.Symbol = symbol_
	symbol.Type = o.Base.Type
	instId = GetInstId(&symbol)
	if precisionIns, ok = c.precisionMap[string(o.Base.Symbol)]; !ok {
		return nil, errors.New("get precision err")
	}
	if precision, ok = precisionIns.(*client.PrecisionItem); !ok {
		return nil, errors.New("get precision err")
	}
	if o.Amount < precision.AmountMin {
		return nil, errors.New(fmt.Sprint("less amount in err:", o.Amount, "<", precision.AmountMin))
	}
	side = ok_api.GetSideTypeToExchange(o.Side)

	if o.Base.Id != 0 {
		options.Add("clOrdId", transform.IdToClientId(o.Hdr.Producer, o.Base.Id))
	}
	// 交易类型
	//tdMode := "cash"

	if c.optionMap != nil {
		if mode, ok := c.optionMap["tdMode"]; ok {
			tdMode, _ = mode.(string)
		}
	}
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
		type_ = ok_api.ORDER_TYPE_POST_ONLY
		options.Add("px", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
	} else if o.Tif == order.TimeInForce_FOK {
		type_ = ok_api.ORDER_TYPE_FOK
		options.Add("px", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
	} else if o.Tif == order.TimeInForce_IOC {
		type_ = ok_api.ORDER_TYPE_IOC
		options.Add("px", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
	} else if o.Tif == order.TimeInForce_OPTIMAL_LIMIT_IOC {
		type_ = ok_api.ORDER_TYPE_OPTIMAL_LIMIT_IOC
	} else {
		type_ = ok_api.GetOrderTypeToExchange(o.OrderType)
		if type_ == ok_api.ORDER_TYPE_LIMIT {
			//options.Add("px", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
			options.Add("px", "1")
		} else if type_ == ok_api.ORDER_TYPE_MARKET && side == ok_api.SIDE_TYPE_BUY {
			amount = o.Price
		}
	}
	var (
		res *ok_api.Resp_Trde_Order
	)

	res, err = c.api.Trade_Order_Info(instId, tdMode, string(side), string(type_), strconv.FormatFloat(amount, 'f', int(precision.Amount), 64), options)
	if err != nil {
		return nil, err
	}
	names := strings.Split(instId, "-")
	var closeDate string
	if len(names) == 3 {
		if names[2] != "SWAP" {
			closeDate = names[2]
		}
	}
	resp.Producer, resp.Id = transform.ClientIdToId(res.Data[0].ClOrdId)
	resp.OrderId = res.Data[0].OrdId
	resp.Symbol = symbol_
	resp.Status = GetOrderStatusFromExchange(res.Data[0].SCode)
	resp.CloseDate = closeDate
	return resp, nil
}
func GetOrderStatusFromExchange(s string) order.OrderStatusCode {
	if s == "0" {
		return order.OrderStatusCode_OPENED
	} else {
		return order.OrderStatusCode_FAILED
	}
}

func (c *ClientOkex) PlaceFutureOrder(o *order.OrderTradeCEX) (*client.OrderRsp, error) {
	return c.PlaceUBaseOrder(o)
}

func (c *ClientOkex) CancelOrder(o *order.OrderCancelCEX) (*client.OrderRsp, error) {
	var (
		resp    = &client.OrderRsp{}
		err     error
		symbol  client.SymbolInfo
		orderId = string(o.Base.IdEx)
		params  = url.Values{}
		res     *ok_api.Resp_Trade_CancelOrder
	)
	// 只有ordeid,其他的是什么意思
	if string(o.Base.IdEx) != "" {
		params.Add("ordId", orderId)
	} else if o.Base.Id != 0 {
		params.Add("clOrdId", transform.IdToClientId(o.Hdr.Producer, o.Base.Id))
	} else {
		return nil, errors.New("id can not be empty")
	}
	symbol = client.SymbolInfo{
		Symbol: string(o.Base.Symbol),
		Type:   o.Base.Type,
	}
	instId := GetInstId(&symbol)
	res, err = c.api.Trade_CancelOrder_Info(instId, &params)

	if err != nil {
		return nil, err
	}
	resp.OrderId = res.Data[0].OrdId
	resp.Producer, resp.Id = transform.ClientIdToId(res.Data[0].ClOrdId)
	resp.RespType = client.OrderRspType_RESULT
	// symbol就是请求的symbol
	resp.Symbol = instId
	// 返回没有以下的
	if res.Data[0].SMsg == "" {
		resp.Status = ok_api.GetOrderStatusFromExchange(ok_api.OrderStatus("CANCELED"))
	} else {
		resp.Status = ok_api.GetOrderStatusFromExchange(ok_api.OrderStatus("FAILED"))
	}
	names := strings.Split(instId, "-")
	var closeDate string
	if len(names) == 3 {
		if names[2] != "SWAP" {
			closeDate = names[2]
		}
	}
	resp.CloseDate = closeDate
	return resp, nil
}

func (c *ClientOkex) GetOrderHistory(req *client.OrderHistoryReq) ([]*client.OrderRsp, error) {
	var (
		resp []*client.OrderRsp
		// todo Market是int32类型
		instType  = req.Market    //symbol	STRING	YES
		startTime = req.StartTime //startTime	LONG	NO
		endTime   = req.EndTime   //endTime	LONG	NO
		limit     = 100           //limit	INT	NO	默认 100; 最大 100.
		params    = url.Values{}
	)

	if endTime > 0 {
		params.Add("after", strconv.FormatInt(endTime, 10))
	}
	if startTime > 0 {
		params.Add("before", strconv.FormatInt(startTime, 10))
	}
	params.Add("limit", strconv.Itoa(limit))

	// 近三个月
	res, err := c.api.Trade_OrdersHistoryArchive_Info(ParseMarket(instType), &params)
	if err != nil {
		return nil, err
	}
	for _, item := range res.Data {
		producer, id := transform.ClientIdToId(item.ClOrdId)
		price, _ := strconv.ParseFloat(item.FillPx, 64)
		symbol := strings.ReplaceAll(item.InstId, "-", "")
		executed, _ := strconv.ParseFloat(item.FillSz, 64)
		avgPrice, _ := strconv.ParseFloat(item.AvgPx, 64)
		accumAmount, _ := strconv.ParseFloat(item.AccFillSz, 64)
		ts, err := strconv.ParseInt(res.Data[0].CTime, 10, 64)
		if err != nil {
			return nil, err
		}
		names := strings.Split(item.InstId, "-")
		var closeDate string
		if len(names) == 3 {
			if names[2] != "SWAP" {
				closeDate = names[2]
			}
		}
		resp = append(resp, &client.OrderRsp{
			Producer:    producer,
			Id:          id,
			OrderId:     item.OrdId,
			Timestamp:   ts,
			Symbol:      symbol,
			RespType:    client.OrderRspType_RESULT,
			Status:      ok_api.GetOrderStatusFromExchange(ok_api.OrderStatus(strings.ToUpper(item.State))),
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

// instId和币对的关系
func (c *ClientOkex) GetOrder(req *order.OrderQueryReq) (*client.OrderRsp, error) {
	var (
		symbol  client.SymbolInfo
		resp    = &client.OrderRsp{}
		instlId = strings.ReplaceAll(string(req.Symbol), "/", "-")
		orderId = req.IdEx //orderId	LONG	NO
		params  = url.Values{}
	)
	symbol = client.SymbolInfo{
		Symbol: string(req.Symbol),
		Type:   req.Type,
	}
	instlId = GetInstId(&symbol)
	if len(req.Producer) > 0 && req.Id != 0 {
		params.Add("clOrdId", transform.IdToClientId(req.Producer, req.Id))
	}

	res, err := c.api.Trade_OrderInfo_Info(instlId, orderId, &params)
	if err != nil {
		return nil, err
	}
	if len(res.Data) == 0 {
		err = errors.New("订单返回为空")
		return nil, err
	}
	// 创建时间
	ts, err := strconv.ParseInt(res.Data[0].CTime, 10, 64)
	if err != nil {
		return nil, err
	}
	resp.Producer, resp.Id = transform.ClientIdToId(res.Data[0].ClOrdId)
	resp.OrderId = res.Data[0].OrdId
	resp.Timestamp = ts

	//resp. = ParseInstType(res.Data[0].InstType)
	resp.RespType = client.OrderRspType_RESULT
	resp.Symbol = ok_api.ParseSymbolName(res.Data[0].InstId)
	resp.Status = ok_api.GetOrderStatusFromExchange(ok_api.OrderStatus(strings.ToUpper(res.Data[0].State)))
	resp.Price, _ = strconv.ParseFloat(res.Data[0].FillPx, 64)
	resp.Executed, _ = strconv.ParseFloat(res.Data[0].FillSz, 64)

	resp.AvgPrice, _ = strconv.ParseFloat(res.Data[0].AvgPx, 64)
	resp.AccumAmount, _ = strconv.ParseFloat(res.Data[0].AccFillSz, 64)
	resp.AccumQty = resp.AvgPrice * resp.AccumAmount
	names := strings.Split(res.Data[0].InstId, "-")
	var closeDate string
	if len(names) == 3 {
		if names[2] != "SWAP" {
			closeDate = names[2]
		}
	}
	resp.CloseDate = closeDate
	return resp, nil
}

// 将string转换为规定的market
func ParseInstType(instType string) common.Market {
	switch instType {
	case "SPOT":
		return common.Market_SPOT
	case "MARGIN":
		return common.Market_MARGIN
	case "SWAP":
		return common.Market_SWAP
	case "FUTURES":
		return common.Market_FUTURE
	case "OPTION":
		return common.Market_OPTION
	default:
		return common.Market_INVALID_MARKET
	}
}

func ParseMarket(market common.Market) string {
	switch market {
	case common.Market_SPOT:
		return "SPOT"
	case common.Market_MARGIN:
		return "MARGIN"
	case common.Market_SWAP:
		return "SWAP"
	case common.Market_FUTURE:
		return "FUTURES"
	case common.Market_OPTION:
		return "OPTION"
	default:
		err := errors.New("参数错误")
		panic(err)
	}
}

func (c *ClientOkex) GetProcessingOrders(req *client.OrderHistoryReq) ([]*client.OrderRsp, error) {
	var (
		resp   []*client.OrderRsp
		symbol client.SymbolInfo
		params = url.Values{}
		instId string
	)
	symbol = client.SymbolInfo{
		Symbol: req.Asset,
		Type:   req.Type,
	}
	instId = GetInstId(&symbol)
	params.Add("instId", instId)
	res, err := c.api.Trade_OrdersPending_Info(&params)
	if err != nil {
		return nil, err
	}
	for _, item := range res.Data {
		producer, id := transform.ClientIdToId(item.ClOrdId)
		price, _ := strconv.ParseFloat(item.FillPx, 64)
		executed, _ := strconv.ParseFloat(item.FillSz, 64)
		avgPrice, _ := strconv.ParseFloat(item.AvgPx, 64)
		accumAmount, _ := strconv.ParseFloat(item.AccFillSz, 64)

		// 创建时间
		ts, err := strconv.ParseInt(res.Data[0].CTime, 10, 64)
		if err != nil {
			return nil, err
		}
		names := strings.Split(res.Data[0].InstId, "-")
		var closeDate string
		if len(names) == 3 {
			if names[2] != "SWAP" {
				closeDate = names[2]
			}
		}
		resp = append(resp, &client.OrderRsp{
			Producer:    producer,
			Id:          id,
			OrderId:     item.OrdId,
			Timestamp:   ts,
			Symbol:      strings.ReplaceAll(item.InstId, "-", "/"),
			RespType:    client.OrderRspType_RESULT,
			Status:      ok_api.GetOrderStatusFromExchange(ok_api.OrderStatus(strings.ToUpper(item.State))),
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

func (c *ClientOkex) Transfer(o *order.OrderTransfer) (*client.OrderRsp, error) {
	//默认都是从现货提币
	var (
		resp            = &client.OrderRsp{}
		err             error
		coin            = string(o.ExchangeToken)
		withdrawOrderId = o.Base.Id                                        //自定义提币ID
		network         = coin + "-" + ok_api.GetNetWorkFromChain(o.Chain) //提币网络
		amount          = o.Amount                                         //数量
		params          = url.Values{}
		fee             string
	)

	params.Add("ccy", coin)
	b, _ := c.api.Asset_Currencies_Info(&params)
	for _, instrument := range b.Data {
		if instrument.Chain == network {
			fee = instrument.MinFee
			break
		}
	}

	params.Del("ccy")
	params.Add("chain", network)
	params.Add("clientId", transform.IdToClientId(o.Hdr.Producer, withdrawOrderId))
	fee_ := transform.StringToX[float64](fee).(float64)
	amt := fmt.Sprintf("%f", amount-fee_)
	var dest string
	if o.Base.Exchange == o.ExchangeTo {
		dest = "3"
	} else {
		dest = "4"
	}

	address := o.TransferAddress
	if len(o.Tag) != 0 {
		address = append(address, ':')
		address = append(address, o.Tag...)
	}

	res, err := c.api.Asset_Withdrawal_Info(coin, amt, dest, string(address), fee, &params)
	if err != nil {
		return nil, err
	}

	resp.OrderId = res.Data[0].WdId
	resp.RespType = client.OrderRspType_ACK
	return resp, err
}

// Internal Server Error
func (c *ClientOkex) GetTransferHistory(req *client.TransferHistoryReq) (*client.TransferHistoryRsp, error) {
	var (
		resp            = &client.TransferHistoryRsp{}
		err             error
		coin            = req.Asset                                    //coin	STRING	NO
		status          = ok_api.GetTransferTypeToExchange(req.Status) //status	INT	NO	0(0:已发送确认Email,1:已被用户取消 2:等待确认 3:被拒绝 4:处理中 5:提现交易失败 6 提现完成)
		withdrawOrderId string                                         //withdrawOrderId	STRING	NO
		offset          = 0                                            //offset	INT	NO
		limit           = 1000                                         //limit	INT	NO	默认：1000， 最大：1000
		startTime       = req.StartTime                                //startTime	LONG	NO	默认当前时间90天前的时间戳
		endTime         = req.EndTime                                  //endTime	LONG	NO	默认当前时间戳
		params          = url.Values{}
		res             *ok_api.Resp_Asset_WithdrawalHistory
	)

	if req.Id != 0 {
		withdrawOrderId = transform.IdToClientId(req.Producer, req.Id)
	}
	//if status == -1 {
	//	return nil, errors.New("transfer status err")
	//}
	if coin != "" {
		params.Add("ccy", coin)
	}
	if withdrawOrderId != "" {
		params.Add("clientId", withdrawOrderId)
	}
	if status >= 0 {
		params.Add("state", strconv.FormatInt(int64(status), 10))
	}
	if startTime > 0 {
		params.Add("before", strconv.FormatInt(startTime, 10))
	}
	if endTime > 0 {
		params.Add("after", strconv.FormatInt(endTime, 10))
	}
	for i := 0; i < 100; i++ {
		params.Add("limit", strconv.Itoa(limit))
		res, err = c.api.Asset_WithdrawalHistory_Info(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Data {
			fee, _ := strconv.ParseFloat(item.Fee, 64)
			timestamp, _ := strconv.ParseInt(item.Ts, 10, 64)
			amount, _ := strconv.ParseFloat(item.Amt, 64)
			stat, err := strconv.Atoi(item.State)
			if err != err {
				return nil, err
			}

			resp.TransferList = append(resp.TransferList, &client.TransferHistoryItem{
				Asset:   item.Ccy,
				Amount:  amount,
				OrderId: item.WdId,
				Network: ok_api.GetChainFromNetWork(item.Chain),
				Status:  ok_api.OkGetTransferTypeFromExchange(ok_api.OkexTransferStatus(stat)),
				Fee:     fee,
				// 没有确认数
				Address:   item.From,
				TxId:      item.TxId,
				Timestamp: timestamp,
			})
		}
		if len(res.Data) < limit {
			break
		}
		offset += limit
	}

	return resp, err
}

func (c *ClientOkex) MoveAsset(o *order.OrderMove) (*client.OrderRsp, error) {
	var (
		resp   = &client.OrderRsp{}
		err    error
		asset  = strings.ReplaceAll(o.Asset, "/", "-") // asset	STRING	YES
		amount = o.Amount                              // amount	DECIMAL	YES
		from   string
		target string
		params = url.Values{}
		type_  = "0"
	)
	// 母账户发起
	if o.ActionUser == order.OrderMoveUserType_Master {
		if o.AccountSource != "" {
			if o.AccountTarget == "" {
				type_ = "2"
			}
		} else if o.AccountSource == "" {
			if o.AccountTarget != "" {
				type_ = "1"
			} else {
				return nil, errors.New("划转参数出错，AccountTarget不应该为空")
			}
		} else {
			return nil, errors.New("划转参数出错，母账号不支持子转子")
		}
	} else if o.ActionUser == order.OrderMoveUserType_Sub {
		// 子账户发起
		if o.AccountTarget == "" {
			type_ = "3"
		} else {
			type_ = "4"
		}
	}
	params.Add("type", type_)
	if type_ == "1" || type_ == "4" {
		params.Add("subAcct", o.AccountTarget)
	} else if type_ == "2" {
		params.Add("subAcct", o.AccountSource)
	}
	if o.Source == common.Market_SPOT {
		from = "18"
	} else if o.Source == common.Market_WALLET {
		from = "6"
	} else {
		return nil, errors.New("划转只能用资金，交易账户")
	}
	if o.Target == common.Market_SPOT {
		target = "18"
	} else if o.Target == common.Market_WALLET {
		target = "6"
	} else {
		return nil, errors.New("划转只能用资金，交易账户")
	}
	res, err := c.api.Asset_Transfer_Info(asset, fmt.Sprintf("%f", amount), from, target, &params)
	if err != nil {
		return nil, err
	}
	resp.OrderId = res.Data[0].TransId
	resp.RespType = client.OrderRspType_ACK
	return resp, err
}

func moveDirectionConvert(tstr string) client.MoveType {
	i, _ := strconv.ParseInt(tstr, 10, 64)
	if i%2 == 0 {
		return client.MoveType_MOVETYPE_OUT
	} else {
		return client.MoveType_MOVETYPE_IN
	}
}
func moveTypeConvert(req *client.MoveHistoryReq) string {
	if req.ActionUser == order.OrderMoveUserType_Master {
		if req.AccountSource != "" {
			// 子母
			return "2"
		} else {
			// 母子
			return "1"
		}
	} else if req.ActionUser == order.OrderMoveUserType_Sub {
		if req.AccountTarget != "" {
			// 子子
			return "4"
		} else {
			// 子母
			return "3"
		}
	} else {
		return "0"
	}
}
func (c *ClientOkex) GetMoveHistory(req *client.MoveHistoryReq) (*client.MoveHistoryRsp, error) {
	var (
		resp      = &client.MoveHistoryRsp{}
		err       error
		startTime = req.StartTime //startTime	LONG	NO
		endTime   = req.EndTime   //endTime	LONG	NO
		params    = url.Values{}
	)
	// 如果指定了idEx，直接进行单点查询
	if req.IdEx != "" {
		params.Set("type", moveTypeConvert(req))
		res, err := c.api.Asset_TransferState_Info(req.IdEx, &params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Data {
			if item.State != "success" {
				continue
			}
			amount, _ := strconv.ParseFloat(item.Amt, 64)
			resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
				Asset:     item.Ccy,
				Amount:    amount,
				Timestamp: time.Now().UnixMilli(),
				Type:      client.MoveType_MOVETYPE_IN,
				Id:        item.TransId,
				Status:    client.MoveStatus_MOVESTATUS_CONFIRMED,
			})
		}
		return resp, err
	}
	if startTime > 0 {
		params.Add("before", strconv.FormatInt(startTime, 10))
	}
	if endTime > 0 {
		params.Add("after", strconv.FormatInt(endTime, 10))
	}

	if req.ActionUser == order.OrderMoveUserType_Master {
		// 子母账户万能划转历史
		// 从子账户转入
		params.Set("type", "21")
		res, err := c.api.Asset_Bills_Info(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Data {
			amount, _ := strconv.ParseFloat(item.BalChg, 64)
			resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
				Asset:     item.Ccy,
				Amount:    amount,
				Timestamp: transform.StringToX[int64](item.Ts).(int64),
				Type:      moveDirectionConvert(item.Type),
				Id:        item.BillId, //FIXME: billis is not moveId??
				Status:    client.MoveStatus_MOVESTATUS_CONFIRMED,
			})
		}
		// 转出至子账户
		params.Set("type", "20")
		res, err = c.api.Asset_Bills_Info(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Data {
			amount, _ := strconv.ParseFloat(item.BalChg, 64)
			resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
				Asset:     item.Ccy,
				Amount:    amount,
				Timestamp: transform.StringToX[int64](item.Ts).(int64),
				Type:      moveDirectionConvert(item.Type),
				Id:        item.BillId,
				Status:    client.MoveStatus_MOVESTATUS_CONFIRMED,
			})
		}
	} else if req.ActionUser == order.OrderMoveUserType_Sub {
		// 转出到母账户
		params.Set("type", "22")
		res, err := c.api.Asset_Bills_Info(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Data {
			amount, _ := strconv.ParseFloat(item.BalChg, 64)
			resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
				Asset:     item.Ccy,
				Amount:    amount,
				Timestamp: transform.StringToX[int64](item.Ts).(int64),
				Type:      moveDirectionConvert(item.Type),
				Id:        item.BillId,
				Status:    client.MoveStatus_MOVESTATUS_CONFIRMED,
			})
		}
		// 母账户转入
		params.Set("type", "23")
		res, err = c.api.Asset_Bills_Info(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Data {
			amount, _ := strconv.ParseFloat(item.BalChg, 64)
			resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
				Asset:     item.Ccy,
				Amount:    amount,
				Timestamp: transform.StringToX[int64](item.Ts).(int64),
				Type:      moveDirectionConvert(item.Type),
				Id:        item.BillId,
				Status:    client.MoveStatus_MOVESTATUS_CONFIRMED,
			})
		}
	} else { //OrderMoveUserType_Internal
		params = url.Values{}
		params.Set("type", "1") // 换转类型
		// 内部划转
		if startTime > 0 {
			params.Add("begin", strconv.FormatInt(startTime, 10))
		}
		if endTime > 0 {
			params.Add("end", strconv.FormatInt(endTime, 10))
		}
		res, err := c.api.Account_Bills_Info(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Data {
			amount, _ := strconv.ParseFloat(item.BalChg, 64)
			resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
				Asset:     item.Ccy,
				Amount:    amount,
				Timestamp: transform.StringToX[int64](item.Ts).(int64),
				Type:      moveDirectionConvert(item.SubType),
				Id:        item.BillId,
				Status:    client.MoveStatus_MOVESTATUS_CONFIRMED,
			})
		}
	}
	return resp, err
}

func (c *ClientOkex) GetDepositHistory(req *client.DepositHistoryReq) (*client.DepositHistoryRsp, error) {
	var (
		resp      = &client.DepositHistoryRsp{}
		err       error
		coin      = req.Asset     //coin	STRING	NO
		startTime = req.StartTime //startTime	LONG	NO	默认当前时间90天前的时间戳
		endTime   = req.EndTime   //endTime	LONG	NO	默认当前时间戳
		offset    = 0             //offset	INT	NO	默认:0
		limit     = 1000          //limit	INT	NO	默认：1000，最大1000
		params    = url.Values{}
		res       *ok_api.Resp_Asset_DepositHistory
	)

	if coin != "" {
		params.Add("ccy", coin)
	}
	if startTime > 0 {
		params.Add("after", strconv.FormatInt(endTime, 10))
	}
	if endTime > 0 {
		params.Add("before", strconv.FormatInt(startTime, 10))
	}

	for i := 0; i < 100; i++ {
		//params.Add("type", "1") // 1为充值类型流水
		params.Add("limit", strconv.Itoa(limit))

		res, err = c.api.Asset_DepositHistory_Info(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Data {

			ts, err := strconv.ParseInt(item.Ts, 10, 64)
			if err != nil {
				return nil, err
			}
			amount, err := strconv.ParseFloat(item.Amt, 64)
			state, err := strconv.Atoi(item.State)
			if err != nil {
				return nil, err
			}
			resp.DepositList = append(resp.DepositList, &client.DepositHistoryItem{
				Asset:     item.Ccy,
				Amount:    amount,
				Network:   ok_api.GetChainFromNetWork(item.Chain),
				Status:    ok_api.GetDepositTypeFromExchange(ok_api.DepositStatus(state)),
				Address:   item.From,
				TxId:      item.TxId,
				Timestamp: ts,
			})
		}
		if len(res.Data) < limit {
			break
		}
		offset += limit
	}
	return resp, err
}

// todo 无订单ID
// res: &{59310 [] Your account does not support VIP loan}
func (c *ClientOkex) Loan(o *order.OrderLoan) (*client.OrderRsp, error) {
	var (
		resp   = &client.OrderRsp{}
		err    error
		asset  = o.Asset  //asset	STRING	YES
		amount = o.Amount //amount	DECIMAL	YES
		params = url.Values{}
	)

	res, err := c.api.Account_BorrowRepay_Info(string(asset), "borrow", fmt.Sprintf("%f", amount), &params)
	if err != nil {
		return nil, err
	}
	_ = res
	//resp.OrderId = strconv.FormatInt(res.TranId, 10)
	resp.RespType = client.OrderRspType_ACK
	return resp, err
}

func (c *ClientOkex) GetLoanOrders(o *client.LoanHistoryReq) (*client.LoanHistoryRsp, error) {
	var (
		resp      = &client.LoanHistoryRsp{}
		err       error
		asset     = o.Asset     //asset	STRING	YES
		startTime = o.StartTime //startTime	LONG	NO
		endTime   = o.EndTime   //endTime	LONG	NO
		current   = 1           //current	LONG	NO	当前查询页。 开始值 1。 默认:1
		size      = 100         //size	LONG	NO	默认:100 最大:100
		params    = url.Values{}
		res       *ok_api.Resp_Account_BorrowRepayHistory
	)

	if o.Asset != "" {
		params.Add("ccy", asset)
	}

	if startTime > 0 {
		params.Add("after", strconv.FormatInt(endTime, 10))
	}
	if endTime > 0 {
		params.Add("before", strconv.FormatInt(startTime, 10))
	}

	for i := 0; i < 100; i++ {
		params.Add("limit", strconv.Itoa(size))
		res, err = c.api.Account_BorrowRepayHistory_Info(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Data {
			// 1 为借币
			if item.Type != "1" {
				continue
			}
			ts, err := strconv.ParseInt(item.Ts, 10, 64)
			principle, err := strconv.ParseFloat(item.Ts, 64)
			if err != nil {
				return nil, err
			}
			resp.LoadList = append(resp.LoadList, &client.LoanHistoryItem{
				Principal: principle,
				Asset:     item.Ccy,
				Timestamp: ts,
			})
		}
		if len(res.Data) < size {
			break
		}
		current++
	}
	return resp, err
}

func (c *ClientOkex) Return(o *order.OrderReturn) (*client.OrderRsp, error) {
	var (
		resp   = &client.OrderRsp{}
		err    error
		asset  = o.Asset  //asset	STRING	YES
		amount = o.Amount //amount	DECIMAL	YES
		params = url.Values{}
	)

	res, err := c.api.Account_BorrowRepay_Info(asset, "repay", fmt.Sprintf("%f", amount), &params)
	if err != nil {
		return nil, err
	}
	_ = res
	resp.RespType = client.OrderRspType_ACK
	return resp, err
}

func (c *ClientOkex) GetReturnOrders(o *client.LoanHistoryReq) (*client.ReturnHistoryRsp, error) {
	var (
		resp      = &client.ReturnHistoryRsp{}
		err       error
		asset     = o.Asset     //asset	STRING	YES
		startTime = o.StartTime //startTime	LONG	NO
		endTime   = o.EndTime   //endTime	LONG	NO
		current   = 1           //current	LONG	NO	当前查询页。 开始值 1。 默认:1
		size      = 100         //size	LONG	NO	默认:100 最大:100
		params    = url.Values{}
		res       *ok_api.Resp_Account_BorrowRepayHistory
	)

	if o.Asset != "" {
		params.Add("ccy", asset)
	}

	if startTime > 0 {
		params.Add("after", strconv.FormatInt(startTime, 10))
	}
	if endTime > 0 {
		params.Add("before", strconv.FormatInt(endTime, 10))
	}

	for i := 0; i < 100; i++ {
		params.Add("limit", strconv.Itoa(size))
		res, err = c.api.Account_BorrowRepayHistory_Info(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Data {
			// 2， 3为还币
			if item.Type == "1" {
				continue
			}
			ts, err := strconv.ParseInt(item.Ts, 10, 64)
			principle, err := strconv.ParseFloat(item.Ts, 64)
			if err != nil {
				return nil, err
			}
			resp.ReturnList = append(resp.ReturnList, &client.ReturnHistoryItem{
				Principal: principle,
				Asset:     item.Ccy,
				Timestamp: ts,
			})
		}
		if len(res.Data) < size {
			break
		}
		current++
	}
	return resp, err
}

func GetInstId(symbols ...*client.SymbolInfo) string {
	var type_ string
	symbol := GetSymbol(symbols[0].Symbol)
	if symbols[0].Type != common.SymbolType_SPOT_NORMAL {
		tmp := u_api.GetFutureTypeFromNats(symbols[0].Type)
		type_ = u_api.GetUBaseSymbol(symbol, tmp)
		symbol = strings.ReplaceAll(type_, "_", "-")
	}
	return symbol
}
