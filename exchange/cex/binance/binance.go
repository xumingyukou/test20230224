package binance

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
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
)

type ClientBinance struct {
	conf base.APIConf

	Api  *spot_api.ApiClient
	UApi *u_api.UApiClient
	CApi *c_api.CApiClient

	tradeFeeMap    map[string]*client.TradeFeeItem    //key: symbol
	transferFeeMap map[string]*client.TransferFeeItem //key: network+token
	precisionMap   map[string]*client.PrecisionItem   //key: symbol

	transferFeeUpdateTime time.Time
	lock                  sync.Mutex
}

func NewClientBinance(conf base.APIConf, maps ...interface{}) *ClientBinance {
	c := &ClientBinance{
		conf: conf,
		Api:  spot_api.NewApiClient(conf, maps...),
		UApi: u_api.NewUApiClient(conf, maps...),
		CApi: c_api.NewCApiClient(conf, maps...),
	}
	c.InitConfMap(maps...)
	return c
}

func NewClientBinance2(conf base.APIConf, cli *http.Client, maps ...interface{}) *ClientBinance {
	// 使用自定义http client
	c := &ClientBinance{
		conf: conf,
		Api:  spot_api.NewApiClient2(conf, cli, maps...),
		UApi: u_api.NewUApiClient2(conf, cli, maps...),
		CApi: c_api.NewCApiClient2(conf, cli, maps...),
	}
	c.InitConfMap(maps...)
	return c
}

func (c *ClientBinance) InitConfMap(maps ...interface{}) {
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
}

func (c *ClientBinance) InitTransferFeeMap() error {
	if time.Now().Sub(c.transferFeeUpdateTime) < time.Hour {
		return nil
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.transferFeeMap = make(map[string]*client.TransferFeeItem)
	configAll, err := c.Api.CapitalConfigGetAll()
	if err != nil {
		return err
	}
	for _, item := range *configAll {
		for _, network := range item.NetworkList {
			networkName := spot_api.GetChainFromNetWork(network.Network)
			if networkName == common.Chain_INVALID_CAHIN {
				continue
			}
			c.transferFeeMap[networkName.String()+strings.ToUpper(item.Coin)] = &client.TransferFeeItem{
				Token:   network.Coin,
				Network: networkName,
				Fee:     transform.StringToX[float64](network.WithdrawFee).(float64),
			}
		}
	}
	c.transferFeeUpdateTime = time.Now()
	return nil
}

func (c *ClientBinance) GetExchange() common.Exchange {
	return common.Exchange_BINANCE
}

func (c *ClientBinance) GetSymbols() []string {
	var symbols []string
	exchangeInfoRes, err := c.Api.ExchangeInfo()
	if err != nil {
		logger.Logger.Error("get exchange in error:", err)
		return symbols
	}
	for _, symbol := range exchangeInfoRes.Symbols {
		if symbol.Status == "TRADING" {
			symbols = append(symbols, fmt.Sprint(symbol.BaseAsset, "/", symbol.QuoteAsset))
		}
	}
	return symbols
}

func (c *ClientBinance) GetFutureSymbols(market common.Market) []*client.SymbolInfo {
	var (
		symbols         []*client.SymbolInfo
		exchangeInfoRes *spot_api.RespExchangeInfo
		err             error
	)
	if market == common.Market_SPOT || market == common.Market_MARGIN {
		exchangeInfoRes, err = c.Api.ExchangeInfo()
	} else if market == common.Market_FUTURE || market == common.Market_SWAP {
		exchangeInfoRes, err = c.UApi.ExchangeInfo()
	} else if market == common.Market_FUTURE_COIN || market == common.Market_SWAP_COIN {
		exchangeInfoRes, err = c.CApi.ExchangeInfo()
	} else {
		fmt.Printf("invalid market/type %v\n", market)
		return symbols
	}
	if err != nil {
		logger.Logger.Error("get exchange in error:", err)
		return symbols
	}
	for _, symbol := range exchangeInfoRes.Symbols {
		if symbol.Status == "TRADING" {
			symbols = append(symbols, &client.SymbolInfo{
				Symbol: fmt.Sprint(symbol.BaseAsset, "/", symbol.QuoteAsset),
				Name:   symbol.Symbol,
				Type:   u_api.GetFutureTypeFromExchange(u_api.ContractType(symbol.ContractType)),
			})
		}
	}
	return symbols
}

func (c *ClientBinance) ParseOrder(orders [][]string, slice *[]*depth.DepthLevel) error {
	for _, item := range orders {
		price, amount, err := transform.ParsePriceAmountFloat(item)
		if err != nil {
			logger.Logger.Errorf("order float parse price error [%s] , response data = %s", err, item)
			return err
		}
		*slice = append(*slice, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	return nil
}

func (c *ClientBinance) GetDepth(symbol *client.SymbolInfo, limit int) (*depth.Depth, error) {
	sym := spot_api.GetSymbol(symbol.Symbol)
	repDepth, err := c.Api.GetDepth(sym, limit)

	if err != nil {
		//utils.Logger.Error(d.Symbol, d.getSymbolName(), "get full depth err", err)
		return nil, err
	}
	dep := &depth.Depth{
		Hdr: &common.MsgHeader{
			Sequence: int64(repDepth.LastUpdateId),
		},
		Exchange:     common.Exchange_BINANCE,
		Market:       common.Market_SPOT,
		Type:         common.SymbolType_SPOT_NORMAL,
		Symbol:       symbol.Symbol,
		TimeExchange: uint64(repDepth.E),
		TimeReceive:  uint64(time.Now().UnixMicro()),
	}
	err = c.ParseOrder(repDepth.Bids, &dep.Bids)
	if err != nil {
		return dep, err
	}
	err = c.ParseOrder(repDepth.Asks, &dep.Asks)
	if err != nil {
		return dep, err
	}
	return dep, err
}

func (c *ClientBinance) GetFutureDepth(symbol *client.SymbolInfo, limit int) (*depth.Depth, error) {
	var (
		symbolName string
		repDepth   *spot_api.RespDepth
		err        error
	)
	if symbol != nil {
		if spot_api.IsUBaseSymbolType(symbol.Type) {
			symbolName = u_api.GetUBaseSymbol(symbol.Symbol, u_api.GetFutureTypeFromNats(symbol.Type))
			repDepth, err = c.UApi.GetDepth(symbolName, limit)
		} else if spot_api.IsCBaseSymbolType(symbol.Type) {
			symbolName = c_api.GetCBaseSymbol(symbol.Symbol, u_api.GetFutureTypeFromNats(symbol.Type))
			repDepth, err = c.CApi.GetDepth(symbolName, limit)
		} else {
			symbolName = c_api.GetCBaseSymbol(symbol.Symbol, u_api.GetFutureTypeFromNats(symbol.Type))
			repDepth, err = c.Api.GetDepth(symbolName, limit)
		}
		if err != nil {
			//utils.Logger.Error(d.Symbol, d.getSymbolName(), "get full depth err", err)
			return nil, err
		}
		dep := &depth.Depth{
			Hdr: &common.MsgHeader{
				Sequence: int64(repDepth.LastUpdateId),
			},
			Exchange:     common.Exchange_BINANCE,
			Market:       spot_api.GetNatsMarket(symbol.Type),
			Type:         symbol.Type,
			Symbol:       symbol.Symbol,
			TimeExchange: uint64(repDepth.E),
			TimeReceive:  uint64(time.Now().UnixMicro()),
		}
		err = c.ParseOrder(repDepth.Bids, &dep.Bids)
		if err != nil {
			return dep, err
		}
		err = c.ParseOrder(repDepth.Asks, &dep.Asks)
		if err != nil {
			return dep, err
		}
		return dep, err
	} else {
		return nil, errors.New("symbol is necessary")
	}
}

func (c *ClientBinance) GetFutureMarkPrice(market common.Market, symbols ...*client.SymbolInfo) (*client.RspMarkPrice, error) { //标记价格
	var (
		rep        *u_api.RespPremiumIndexList
		err        error
		value      = &url.Values{}
		symbolMap  = make(map[string]*client.SymbolInfo)
		symbolName string
		res        = &client.RspMarkPrice{}
	)
	for _, sym := range symbols {
		if spot_api.IsUBaseSymbolType(sym.Type) {
			symbolName = u_api.GetUBaseSymbol(sym.Symbol, u_api.GetFutureTypeFromNats(sym.Type))
		} else if spot_api.IsCBaseSymbolType(sym.Type) {
			symbolName = c_api.GetCBaseSymbol(sym.Symbol, u_api.GetFutureTypeFromNats(sym.Type))
		} else {
			symbolName = spot_api.GetSymbol(sym.Symbol)
		}
		symbolMap[symbolName] = sym
	}
	value.Add("symbol", "")
	if len(symbols) == 1 {
		value.Add("symbol", symbolName)
	}
	if spot_api.IsUBaseMarket(market) {
		rep, err = c.UApi.PremiumIndex(value)
	} else if spot_api.IsCBaseMarket(market) {
		rep, err = c.CApi.PremiumIndex(value)
	} else {
		return nil, errors.New("")
	}
	if err != nil {
		return nil, err
	}
	for _, i := range *rep {
		symbol := spot_api.ParseSymbolName(i.Symbol)
		type_ := symbolMap[i.Symbol].Type
		if _, ok := symbolMap[i.Symbol]; len(symbols) != 0 && !ok {
			symbol = symbolMap[i.Symbol].Symbol
			type_ = symbolMap[i.Symbol].Type
		} else {
			continue
		}
		res.Item = append(res.Item, &client.MarkPriceItem{
			Symbol:     symbol,
			Type:       type_,
			UpdateTime: i.Time * 1000,
			MarkPrice:  transform.StringToX[float64](i.MarkPrice).(float64),
		})
	}
	return res, err
}

func (c *ClientBinance) IsExchangeEnable() bool {
	serverTime1, err := c.Api.ServerTime()
	res1 := err == nil && time.Since(time.UnixMilli(serverTime1.ServerTime)) < time.Second*60
	serverTime2, err := c.UApi.ServerTime()
	res2 := err == nil && time.Since(time.UnixMilli(serverTime2.ServerTime)) < time.Second*60
	serverTime3, err := c.CApi.ServerTime()
	res3 := err == nil && time.Since(time.UnixMilli(serverTime3.ServerTime)) < time.Second*60
	return res1 && res2 && res3
}

func (c *ClientBinance) GetTradeFee(symbols ...string) (*client.TradeFee, error) { //查询交易手续费
	var (
		res             = &client.TradeFee{}
		searchSymbolMap = make(map[string]bool)
	)
	for _, symbol := range symbols {
		searchSymbolMap[spot_api.GetSymbol(symbol)] = true
	}
	trakeFeeRes, err := c.Api.AssetTradeFee()
	if err != nil {
		return nil, err
	}
	for _, tradeFee := range *trakeFeeRes {
		var (
			takerFee, makerFee float64
		)
		if len(searchSymbolMap) > 0 {
			_, ok := searchSymbolMap[tradeFee.Symbol]
			if !ok {
				continue
			}
		}
		takerFee, err = strconv.ParseFloat(tradeFee.TakerCommission, 64)
		if err != nil {
			fmt.Println("convert taker fee err:", tradeFee.Symbol, tradeFee.TakerCommission)
			continue
		}
		makerFee, err = strconv.ParseFloat(tradeFee.MakerCommission, 64)
		if err != nil {
			fmt.Println("convert maker fee err:", tradeFee.Symbol, tradeFee.MakerCommission)
			continue
		}
		res.TradeFeeList = append(res.TradeFeeList, &client.TradeFeeItem{
			Symbol: spot_api.ParseSymbolName(tradeFee.Symbol),
			Type:   common.SymbolType_SPOT_NORMAL,
			Maker:  makerFee,
			Taker:  takerFee,
		})
	}
	return res, nil
}

func (c *ClientBinance) GetFutureTradeFee(market common.Market, symbols ...*client.SymbolInfo) (*client.TradeFee, error) { //查询交易手续费
	//symbols为空则返回常规
	var (
		res = &client.TradeFee{}
	)
	if len(symbols) == 0 {
		type_ := common.SymbolType_SWAP_FOREVER
		if spot_api.IsCBaseMarket(market) {
			type_ = common.SymbolType_SWAP_COIN_FOREVER
		}
		for _, symbol := range []string{"BTC/USDT", "LTC/USDT", "ETH/USDT", "LINK/USDT", "BCH/USDT", "XRP/USDT", "EOS/USDT", "TRX/USDT", "ETC/USDT", "DOT/USDT", "ADA/USDT", "BNB/USDT", "FIL/USDT", "UNI/USDT", "XLM/USDT", "DOGE/USDT"} {
			symbols = append(symbols, &client.SymbolInfo{
				Symbol: symbol,
				Type:   type_,
			})
		}
	}
	for _, symbol := range symbols {
		var (
			takerFee, makerFee float64
			symbolName         string
			trakeFeeRes        *u_api.RespCommissionRate
			err                error
		)
		if spot_api.IsUBaseSymbolType(symbol.Type) {
			symbolName = u_api.GetUBaseSymbol(symbol.Symbol, u_api.GetFutureTypeFromNats(symbol.Type))
			trakeFeeRes, err = c.UApi.CommissionRate(symbolName)
		} else if spot_api.IsCBaseSymbolType(symbol.Type) {
			symbolName = c_api.GetCBaseSymbol(symbol.Symbol, u_api.GetFutureTypeFromNats(symbol.Type))
			trakeFeeRes, err = c.CApi.CommissionRate(symbolName)
		} else {
			return nil, errors.New("use GetTradeFee api")
		}
		if err != nil {
			return nil, err
		}
		takerFee, err = strconv.ParseFloat(trakeFeeRes.TakerCommissionRate, 64)
		if err != nil {
			fmt.Println("convert taker fee err:", trakeFeeRes.Symbol, trakeFeeRes.TakerCommissionRate)
			continue
		}
		makerFee, err = strconv.ParseFloat(trakeFeeRes.MakerCommissionRate, 64)
		if err != nil {
			fmt.Println("convert maker fee err:", trakeFeeRes.Symbol, trakeFeeRes.MakerCommissionRate)
			continue
		}
		res.TradeFeeList = append(res.TradeFeeList, &client.TradeFeeItem{
			Symbol: spot_api.ParseSymbolName(trakeFeeRes.Symbol),
			Type:   symbol.Type,
			Maker:  makerFee,
			Taker:  takerFee,
		})
	}
	return res, nil
}

func (c *ClientBinance) GetTransferFee(chain common.Chain, tokens ...string) (*client.TransferFee, error) {
	var (
		res            = &client.TransferFee{}
		searchTokenMap = make(map[string]bool)
	)
	for _, token := range tokens {
		searchTokenMap[token] = true
	}
	err := c.InitTransferFeeMap()
	if err != nil {
		return nil, err
	}
	for _, token := range tokens {
		res.TransferFeeList = append(res.TransferFeeList, c.transferFeeMap[chain.String()+strings.ToUpper(token)])
	}
	return res, nil
}

func (c *ClientBinance) GetPrecision(symbols ...string) (*client.Precision, error) {
	var (
		res             = &client.Precision{}
		searchSymbolMap = make(map[string]bool)
	)
	for _, symbol := range symbols {
		searchSymbolMap[spot_api.GetSymbol(symbol)] = true
	}
	exchangeInfoRes, err := c.Api.ExchangeInfo()
	if err != nil {
		return nil, err
	}
	for _, symbol := range exchangeInfoRes.Symbols {
		if len(searchSymbolMap) > 0 {
			_, ok := searchSymbolMap[symbol.Symbol]
			if !ok {
				continue
			}
		}
		precisionConf := symbol.GetPrecision()
		res.PrecisionList = append(res.PrecisionList, &client.PrecisionItem{
			Symbol:    spot_api.ParseSymbolName(symbol.Symbol),
			Amount:    int64(precisionConf.AmountPrecision),
			Price:     int64(precisionConf.PricePrecision),
			AmountMin: precisionConf.MinAmount,
		})
	}
	return res, nil
}

func (c *ClientBinance) GetFuturePrecision(market common.Market, symbols ...*client.SymbolInfo) (*client.Precision, error) {
	var (
		res             = &client.Precision{}
		searchSymbolMap = make(map[string]*client.SymbolInfo)
		exchangeInfoRes *spot_api.RespExchangeInfo
		err             error
	)
	for _, symbol := range symbols {
		if spot_api.IsUBaseMarket(market) {
			searchSymbolMap[u_api.GetUBaseSymbol(symbol.Symbol, u_api.GetFutureTypeFromNats(symbol.Type))] = symbol
		} else if spot_api.IsCBaseMarket(market) {
			searchSymbolMap[c_api.GetCBaseSymbol(symbol.Symbol, u_api.GetFutureTypeFromNats(symbol.Type))] = symbol
		} else {
			searchSymbolMap[spot_api.GetSymbol(symbol.Symbol)] = symbol
		}
	}
	if spot_api.IsUBaseMarket(market) {
		exchangeInfoRes, err = c.UApi.ExchangeInfo()
	} else if spot_api.IsCBaseMarket(market) {
		exchangeInfoRes, err = c.CApi.ExchangeInfo()
	} else {
		exchangeInfoRes, err = c.Api.ExchangeInfo()
	}
	if err != nil {
		return nil, err
	}
	for _, symbol := range exchangeInfoRes.Symbols {
		if len(searchSymbolMap) > 0 {
			_, ok := searchSymbolMap[symbol.Symbol]
			if !ok {
				continue
			}
		}
		precisionConf := symbol.GetPrecision()
		res.PrecisionList = append(res.PrecisionList, &client.PrecisionItem{
			Symbol:    symbol.BaseAsset + "/" + symbol.QuoteAsset,
			Type:      u_api.GetFutureTypeFromExchange(u_api.ContractType(symbol.ContractType)),
			Amount:    int64(precisionConf.AmountPrecision),
			Price:     int64(precisionConf.PricePrecision),
			AmountMin: precisionConf.MinAmount,
		})
	}
	return res, nil
}

func (c *ClientBinance) GetBalance() (*client.SpotBalance, error) {
	var (
		respAccount *spot_api.RespAccount
		respAsset   *spot_api.RespAsset
		res         = &client.SpotBalance{}
		err         error
	)
	respAccount, err = c.Api.Account()
	if err != nil {
		return res, err
	}
	res.UpdateTime = respAccount.UpdateTime * 1000
	for _, balance := range respAccount.Balances {
		var (
			free, frozen float64
		)
		if free, err = strconv.ParseFloat(balance.Free, 64); err != nil {
			logger.Logger.Error("get free err", balance)
			return res, err
		}
		if frozen, err = strconv.ParseFloat(balance.Locked, 64); err != nil {
			logger.Logger.Error("get frozen err", balance)
			return res, err
		}
		res.BalanceList = append(res.BalanceList, &client.SpotBalanceItem{
			Asset:  balance.Asset,
			Free:   free,
			Frozen: frozen,
			Total:  free + frozen,
		})
	}
	if c.conf.IsTest {
		return res, nil
	}
	respAsset, err = c.Api.Asset(nil)
	if err != nil {
		return res, err
	}
	for _, balance := range *respAsset {
		var (
			free   = transform.StringToX[float64](balance.Free).(float64)
			freeze = transform.StringToX[float64](balance.Freeze).(float64)
		)
		res.WalletList = append(res.WalletList, &client.SpotBalanceItem{
			Asset:  balance.Asset,
			Free:   free,
			Frozen: freeze,
			Total:  free + freeze,
		})
	}
	return res, err
}

func (c *ClientBinance) GetFutureBalance(market common.Market) (*client.UBaseBalance, error) {
	var (
		respAccount *u_api.RespAccount
		res         = &client.UBaseBalance{}
		err         error
	)
	if market == common.Market_FUTURE || market == common.Market_SWAP {
		respAccount, err = c.UApi.Account()
	} else if market == common.Market_FUTURE_COIN || market == common.Market_SWAP_COIN {
		respAccount, err = c.CApi.Account()
	}
	if err != nil {
		return res, err
	}
	res.UpdateTime = respAccount.UpdateTime * 1000
	res.Market = market
	res.Balance = transform.StringToX[float64](respAccount.TotalWalletBalance).(float64)
	res.Unprofit = transform.StringToX[float64](respAccount.TotalUnrealizedProfit).(float64)
	res.Rights = res.Balance + res.Unprofit
	res.Available = transform.StringToX[float64](respAccount.AvailableBalance).(float64)
	res.TotalMarginBalance = transform.StringToX[float64](respAccount.TotalMarginBalance).(float64)
	res.Used = transform.StringToX[float64](respAccount.TotalPositionInitialMargin).(float64)
	for _, asset := range respAccount.Assets {
		var (
			balance   = transform.StringToX[float64](asset.WalletBalance).(float64)
			unprofit  = transform.StringToX[float64](asset.UnrealizedProfit).(float64)
			available = transform.StringToX[float64](asset.AvailableBalance).(float64)
		)
		res.UBaseBalanceList = append(res.UBaseBalanceList, &client.UBaseBalanceItem{
			Asset:     asset.Asset,
			Balance:   balance,
			Market:    market,
			Unprofit:  unprofit,
			Rights:    balance + unprofit,
			Available: available,
			Used:      balance + unprofit - available,
		})
	}
	for _, position := range respAccount.Positions { //仓位
		var closeDate string
		symbol, market_, type_ := u_api.GetContractType(position.Symbol)
		symList := strings.Split(position.Symbol, "_")
		if len(symList) > 1 {
			closeDate = symList[1]
		}
		side := order.TradeSide_BUY
		if strings.Contains(position.PositionAmt, "-") {
			side = order.TradeSide_SELL
		}
		res.UBasePositionList = append(res.UBasePositionList, &client.UBasePositionItem{
			Symbol:         symbol,
			Market:         market_,
			Type:           type_,
			CloseDate:      closeDate,
			Price:          transform.StringToX[float64](position.EntryPrice).(float64),
			MaintainMargin: transform.StringToX[float64](position.MaintMargin).(float64),
			InitialMargin:  transform.StringToX[float64](position.InitialMargin).(float64),
			Notional:       math.Abs(transform.StringToX[float64](position.Notional).(float64)),
			Leverage:       transform.StringToX[float64](position.Leverage).(float64),
			Position:       math.Abs(transform.StringToX[float64](position.PositionAmt).(float64)),
			Side:           side,
			Unprofit:       transform.StringToX[float64](position.UnrealizedProfit).(float64),
		})
	}
	return res, err
}

func (c *ClientBinance) getMarginAsset(balance *spot_api.MarginAssetItem) (*client.MarginBalanceItem, error) {
	var (
		free, frozen, borrowed, netAsset float64
		err                              error
	)
	if free, err = strconv.ParseFloat(balance.Free, 64); err != nil {
		logger.Logger.Error("get free err", balance)
		return nil, err
	}
	if frozen, err = strconv.ParseFloat(balance.Locked, 64); err != nil {
		logger.Logger.Error("get frozen err", balance)
		return nil, err
	}
	if borrowed, err = strconv.ParseFloat(balance.Borrowed, 64); err != nil {
		logger.Logger.Error("get borrowed err", balance)
		return nil, err
	}
	if netAsset, err = strconv.ParseFloat(balance.NetAsset, 64); err != nil {
		logger.Logger.Error("get netAsset err", balance)
		return nil, err
	}
	res := &client.MarginBalanceItem{Asset: balance.Asset, Total: free + frozen, Borrowed: borrowed, Free: free, Frozen: frozen, NetAsset: netAsset}
	return res, err
}

func (c *ClientBinance) getMarginIsolatedAsset(balance *spot_api.MarginIsolatedAssetItem) (*client.MarginBalanceItem, error) {
	var (
		free, frozen, borrowed, netAsset float64
		err                              error
	)
	if free, err = strconv.ParseFloat(balance.Free, 64); err != nil {
		logger.Logger.Error("get free err", balance)
		return nil, err
	}
	if frozen, err = strconv.ParseFloat(balance.Locked, 64); err != nil {
		logger.Logger.Error("get frozen err", balance)
		return nil, err
	}
	if borrowed, err = strconv.ParseFloat(balance.Borrowed, 64); err != nil {
		logger.Logger.Error("get borrowed err", balance)
		return nil, err
	}
	if netAsset, err = strconv.ParseFloat(balance.NetAsset, 64); err != nil {
		logger.Logger.Error("get netAsset err", balance)
		return nil, err
	}
	res := &client.MarginBalanceItem{Asset: balance.Asset, Total: free + frozen, Borrowed: borrowed, Free: free, Frozen: frozen, NetAsset: netAsset}
	return res, err
}

func (c *ClientBinance) GetMarginBalance() (*client.MarginBalance, error) {
	var (
		respAccount                                                *spot_api.RespMarginAccount
		res                                                        = &client.MarginBalance{}
		marginLevel, totalNetAsset, totalAsset, totalLiablityAsset float64
		err                                                        error
	)
	respAccount, err = c.Api.MarginAccount()
	if err != nil {
		return res, err
	}
	if marginLevel, err = strconv.ParseFloat(respAccount.MarginLevel, 64); err != nil {
		logger.Logger.Error("get MarginLevel err", respAccount)
		return res, err
	}
	if totalNetAsset, err = strconv.ParseFloat(respAccount.TotalNetAssetOfBtc, 64); err != nil {
		logger.Logger.Error("get TotalNetAssetOfBtc err", respAccount)
		return res, err
	}
	if totalAsset, err = strconv.ParseFloat(respAccount.TotalAssetOfBtc, 64); err != nil {
		logger.Logger.Error("get TotalAssetOfBtc err", respAccount)
		return res, err
	}
	if totalLiablityAsset, err = strconv.ParseFloat(respAccount.TotalLiabilityOfBtc, 64); err != nil {
		logger.Logger.Error("get TotalLiabilityOfBtc err", respAccount)
		return res, err
	}
	res.MarginLevel = marginLevel
	res.TotalAsset = totalAsset
	res.TotalNetAsset = totalNetAsset
	res.TotalLiabilityAsset = totalLiablityAsset
	res.QuoteAsset = "BTC"
	res.UpdateTime = time.Now().UnixMicro()
	for _, balance := range respAccount.UserAssets {
		var item *client.MarginBalanceItem
		item, err = c.getMarginAsset(balance)
		if err != nil {
			return nil, err
		}
		res.MarginBalanceList = append(res.MarginBalanceList, item)
	}
	return res, err
}

func (c *ClientBinance) GetMarginIsolatedBalance(symbols ...string) (*client.MarginIsolatedBalance, error) {
	var (
		respAccount                       *spot_api.RespMarginIsolatedAccount
		res                               = &client.MarginIsolatedBalance{}
		totalNetAsset, totalLiablityAsset float64
		err                               error
		symbolList                        []string
	)
	for _, symbol := range symbols {
		symbolList = append(symbolList, strings.Replace(strings.ToUpper(symbol), "/", "", -1))
	}
	respAccount, err = c.Api.MarginIsolatedAccount(symbolList...)
	if err != nil {
		return res, err
	}
	if totalNetAsset, err = strconv.ParseFloat(respAccount.TotalNetAssetOfBtc, 64); err != nil {
		logger.Logger.Error("get free err", respAccount)
		return res, err
	}
	if totalLiablityAsset, err = strconv.ParseFloat(respAccount.TotalLiabilityOfBtc, 64); err != nil {
		logger.Logger.Error("get free err", respAccount)
		return res, err
	}
	res.TotalNetAsset = totalNetAsset
	res.TotalLiabilityAsset = totalLiablityAsset
	res.QuoteAsset = "BTC"
	for _, balance := range respAccount.Assets {
		var (
			baseAsset, quoteAsset *client.MarginBalanceItem
		)
		if baseAsset, err = c.getMarginIsolatedAsset(balance.BaseAsset); err != nil {
			return nil, err
		}
		if quoteAsset, err = c.getMarginIsolatedAsset(balance.BaseAsset); err != nil {
			return nil, err
		}
		res.MarginIsolatedBalanceList = append(res.MarginIsolatedBalanceList, &client.MarginIsolatedBalanceItem{
			BaseAsset:  baseAsset,
			QuoteAsset: quoteAsset,
		})
	}
	return res, err
}

func (c *ClientBinance) PlaceOrder(o *order.OrderTradeCEX) (*client.OrderRsp, error) {
	var (
		symbol       string
		side         spot_api.SideType
		type_        spot_api.OrderType
		isIsolated   spot_api.MarginType
		options      = &url.Values{}
		resp         = &client.OrderRsp{}
		precisionIns interface{}
		precision    *client.PrecisionItem
		ok           bool
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
	side = spot_api.GetSideTypeToExchange(o.Side)
	switch o.Base.Type {
	case common.SymbolType_MARGIN_NORMAL:
		isIsolated = spot_api.MARGIN_TYPE_NORMAL
	case common.SymbolType_MARGIN_ISOLATED:
		isIsolated = spot_api.MARGIN_TYPE_ISOLATED
	}
	options.Add("timeInForce", string(spot_api.GetTimeInForceToExchange(o.Tif)))
	if o.Base.Id != 0 {
		options.Add("newClientOrderId", transform.IdToClientId(o.Hdr.Producer, o.Base.Id))
	}
	options.Add("quantity", strconv.FormatFloat(o.Amount, 'f', int(precision.Amount), 64))
	type_ = spot_api.GetOrderTypeToExchange(o.OrderType)
	switch o.OrderType {
	case order.OrderType_LIMIT:
		options.Add("price", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
		if o.TradeType == order.TradeType_MAKER {
			type_ = spot_api.ORDER_TYPE_LIMIT_MAKER
			options.Del("timeInForce")
		} else {
			if options.Get("timeInForce") == "" {
				return nil, errors.New("need timeInForce param")
			}
		}
	case order.OrderType_MARKET:
		options.Del("timeInForce")
	case order.OrderType_STOP_LOSS:
		options.Del("timeInForce")
		//止损价格:stopPrice 或者 trailingDelta
		options.Add("stopPrice", transform.XToString(o.PriceLimit))
	case order.OrderType_STOP_LOSS_LIMIT:
		//price, stopPrice 或者 trailingDelta
		options.Add("stopPrice", transform.XToString(o.PriceLimit))
	case order.OrderType_TAKE_PROFIT:
		options.Del("timeInForce")
		//stopPrice 或者 trailingDelta
		options.Add("stopPrice", transform.XToString(o.PriceLimit))
	case order.OrderType_TAKE_PROFIT_LIMIT:
		//price, stopPrice 或者 trailingDelta
		options.Add("stopPrice", transform.XToString(o.PriceLimit))
		options.Add("price", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
	default:
		return nil, errors.New("order type error")
	}
	options.Add("newOrderRespType", string(spot_api.ORDER_RESP_TYPE_FULL))
	var (
		res     interface{}
		resFull *spot_api.RespMarginOrderFull
	)
	if o.Base.Market == common.Market_SPOT {
		res, err = c.Api.Order(symbol, side, type_, options)
	} else if o.Base.Market == common.Market_MARGIN {
		res, err = c.Api.MarginOrder(symbol, isIsolated, side, type_, options)
	}
	if err != nil {
		return nil, err
	}
	if resFull, ok = res.(*spot_api.RespMarginOrderFull); ok {
		resp.Producer, resp.Id = transform.ClientIdToId(resFull.ClientOrderId)
		resp.OrderId = strconv.Itoa(resFull.OrderId)
		resp.Timestamp = resFull.TransactTime
		resp.RespType = client.OrderRspType_FULL
		resp.Symbol = resFull.Symbol
		resp.Status = spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(resFull.Status))
		resp.AccumAmount, _ = strconv.ParseFloat(resFull.ExecutedQty, 64)
		resp.AccumQty, _ = strconv.ParseFloat(resFull.CummulativeQuoteQty, 64)
		// binance没有平均binance没有平均价，不用设置avgPrice
		for _, fill := range resFull.Fills {
			price, _ := strconv.ParseFloat(fill.Price, 64)
			qty, _ := strconv.ParseFloat(fill.Qty, 64)
			commission, _ := strconv.ParseFloat(fill.Commission, 64)
			resp.Fills = append(resp.Fills, &client.FillItem{
				Price:           price,
				Qty:             qty,
				Commission:      commission,
				CommissionAsset: fill.CommissionAsset,
			})
			resp.Price = price
			resp.Executed = qty
		}
		return resp, nil
	} else {
		return resp, errors.New("parse response error")
	}
}

func (c *ClientBinance) PlaceFutureOrder(o *order.OrderTradeCEX) (*client.OrderRsp, error) {
	var (
		symbol    string
		side      spot_api.SideType
		type_     spot_api.OrderType
		options   = &url.Values{}
		resp      = &client.OrderRsp{}
		err       error
		resFull   *u_api.RespUBaseOrderResult
		ok        bool
		precision *client.PrecisionItem
	)
	if spot_api.IsUBaseSymbolType(o.Base.Type) {
		symbol = u_api.GetUBaseSymbol(string(o.Base.Symbol), u_api.GetFutureTypeFromNats(o.Base.Type))
	} else if spot_api.IsCBaseSymbolType(o.Base.Type) {
		symbol = c_api.GetCBaseSymbol(string(o.Base.Symbol), u_api.GetFutureTypeFromNats(o.Base.Type))
	} else {
		symbol = spot_api.GetSymbol(string(o.Base.Symbol))
	}
	side = spot_api.GetSideTypeToExchange(o.Side)
	options.Add("timeInForce", string(spot_api.GetTimeInForceToExchange(o.Tif)))
	if o.Base.Id != 0 {
		options.Add("newClientOrderId", transform.IdToClientId(o.Hdr.Producer, o.Base.Id))
	}
	if precision, ok = c.precisionMap[string(o.Base.Symbol)]; !ok {
		return nil, errors.New("get precision err")
	}
	if o.Amount < precision.AmountMin {
		return nil, errors.New(fmt.Sprint("less amount in err:", o.Amount, "<", precision.AmountMin))
	}
	options.Add("quantity", strconv.FormatFloat(o.Amount, 'f', int(precision.Amount), 64))
	type_ = spot_api.GetOrderTypeToExchange(o.OrderType)
	switch o.OrderType {
	case order.OrderType_LIMIT:
		options.Add("price", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
		if o.TradeType == order.TradeType_MAKER {
			if options.Get("timeInForce") == "" {
				options.Set("timeInForce", string(spot_api.TIME_IN_FORCE_GTX))
			} else if options.Get("timeInForce") != string(spot_api.TIME_IN_FORCE_GTX) {
				return nil, errors.New("timeInForce must be GTX for limit maker")
			}
		}
		if options.Get("timeInForce") == "" {
			return nil, errors.New("need timeInForce param")
		}
	case order.OrderType_MARKET:
		options.Del("timeInForce")
	case order.OrderType_STOP_LOSS:
		options.Del("timeInForce")
		//止损价格:stopPrice 或者 trailingDelta
		options.Add("stopPrice", transform.XToString(o.PriceLimit))
	case order.OrderType_STOP_LOSS_LIMIT:
		//price, stopPrice 或者 trailingDelta
		options.Add("stopPrice", transform.XToString(o.PriceLimit))
	case order.OrderType_TAKE_PROFIT:
		options.Del("timeInForce")
		//stopPrice 或者 trailingDelta
		options.Add("stopPrice", transform.XToString(o.PriceLimit))
	case order.OrderType_TAKE_PROFIT_LIMIT:
		//price, stopPrice 或者 trailingDelta
		options.Add("stopPrice", transform.XToString(o.PriceLimit))
		options.Add("price", strconv.FormatFloat(o.Price, 'f', int(precision.Price), 64))
	default:
		return nil, errors.New("order type error")
	}
	if spot_api.IsUBaseSymbolType(o.Base.Type) {
		resFull, err = c.UApi.Order(symbol, side, type_, options)
	} else if spot_api.IsCBaseSymbolType(o.Base.Type) {
		resFull, err = c.CApi.Order(symbol, side, type_, options)
	} else {
		return nil, fmt.Errorf("invalid market/type %v:%v", o.Base.Market, o.Base.Type)
	}
	if err != nil {
		return nil, err
	}
	resp.Producer, resp.Id = transform.ClientIdToId(resFull.ClientOrderId)
	resp.OrderId = strconv.Itoa(resFull.OrderId)
	resp.Timestamp = resFull.UpdateTime
	resp.RespType = client.OrderRspType_FULL
	resp.Symbol = resFull.Symbol
	resp.Status = spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(resFull.Status))
	resp.AccumAmount, _ = strconv.ParseFloat(resFull.ExecutedQty, 64)
	resp.AccumQty, _ = strconv.ParseFloat(resFull.CumQuote, 64)
	resp.AvgPrice, _ = strconv.ParseFloat(resFull.AvgPrice, 64)
	var closeDate string
	names := strings.Split(resFull.Symbol, "_")
	if len(names) > 1 && len(names[1]) == 6 {
		closeDate = names[1]
	}
	resp.CloseDate = closeDate
	return resp, nil
}

func (c *ClientBinance) CancelOrder(o *order.OrderCancelCEX) (*client.OrderRsp, error) {
	var (
		resp    = &client.OrderRsp{}
		err     error
		symbol  = c_api.GetCBaseSymbol(string(o.Base.Symbol), u_api.GetFutureTypeFromNats(o.Base.Type))
		orderId = o.Base.IdEx
		params  = url.Values{}
		res     *spot_api.RespCancelOrder
		u_res   *u_api.RespUBaseOrderResult
	)
	if string(o.Base.IdEx) != "" {
		params.Add("orderId", orderId)
	}
	if o.Base.Id != 0 {
		params.Add("origClientOrderId", transform.IdToClientId(o.Hdr.Producer, o.Base.Id))
	}
	if o.CancelId != 0 {
		params.Add("newClientOrderId", transform.IdToClientId(o.Hdr.Producer, o.CancelId))
	}
	if o.Base.Type == common.SymbolType_SPOT_NORMAL {
		res, err = c.Api.CancelOrder(symbol, &params)
		resp.AccumQty, _ = strconv.ParseFloat(res.CummulativeQuoteQty, 64)
	} else if o.Base.Type == common.SymbolType_MARGIN_NORMAL || o.Base.Type == common.SymbolType_MARGIN_ISOLATED {
		res, err = c.Api.CancelMarginOrder(symbol, &params)
		resp.AccumQty, _ = strconv.ParseFloat(res.CummulativeQuoteQty, 64)
	} else if spot_api.IsUBaseSymbolType(o.Base.Type) {
		u_res, err = c.UApi.CancelOrder(symbol, &params)
		res = &u_res.RespCancelOrder
		resp.AccumQty, _ = strconv.ParseFloat(u_res.CumQuote, 64)
	} else if spot_api.IsCBaseSymbolType(o.Base.Type) {
		u_res, err = c.CApi.CancelOrder(symbol, &params)
		res = &u_res.RespCancelOrder
		resp.AccumQty, _ = strconv.ParseFloat(u_res.CumQuote, 64)
	} else {
		return nil, fmt.Errorf("invalid market/type %v:%v", o.Base.Market, o.Base.Type)
	}
	if err != nil {
		return nil, err
	}
	resp.OrderId = strconv.Itoa(res.OrderId)
	resp.Producer, resp.Id = transform.ClientIdToId(res.OrigClientOrderId)
	resp.RespType = client.OrderRspType_RESULT
	resp.Symbol = res.Symbol
	resp.Status = spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(res.Status))
	resp.AccumAmount, _ = strconv.ParseFloat(res.ExecutedQty, 64)
	// binance没有平均价，不用设置avgPrice
	return resp, nil
}

func (c *ClientBinance) Transfer(o *order.OrderTransfer) (*client.OrderRsp, error) {
	//默认都是从现货提币
	var (
		resp            = &client.OrderRsp{}
		err             error
		coin            = string(o.ExchangeToken)
		withdrawOrderId = transform.IdToClientId(o.Hdr.Producer, o.Base.Id) //自定义提币ID
		network         = spot_api.GetNetWorkFromChain(o.Chain)             //提币网络
		address         = string(o.TransferAddress)                         //提币地址
		addressTag      = string(o.Tag)                                     //某些币种例如 XRP,XMR 允许填写次级地址标签
		amount          = o.Amount                                          //数量
		name            = string(o.Comment)                                 //地址的备注，填写该参数后会加入该币种的提现地址簿。地址簿上限为20，超出后会造成提现失败。地址中的空格需要encode成%20
		//walletType =	INTEGER	NO	表示出金使用的钱包，0为现货钱包，1为资金钱包，默认为现货钱包
		params = url.Values{}
	)
	params.Add("name", name)
	params.Add("network", network)
	params.Add("withdrawOrderId", withdrawOrderId)
	params.Add("addressTag", addressTag)
	params.Add("transactionFeeFlag", "true") // 当站内转账时免手续费, true: 手续费归资金转入方; false: 手续费归资金转出方; . 默认 false.

	res, err := c.Api.CapitalWithdrawApply(coin, address, amount, &params)
	if err != nil {
		return nil, err
	}
	resp.OrderId = res.Id
	resp.RespType = client.OrderRspType_ACK
	return resp, err
}

func (c *ClientBinance) MoveAsset(o *order.OrderMove) (*client.OrderRsp, error) {
	var (
		resp       = &client.OrderRsp{}
		res        *spot_api.RespAssetTransfer
		err        error
		type_      = spot_api.GetMoveSide(o.Source, o.Target) // type	ENUM	YES
		asset      = o.Asset                                  // asset	STRING	YES
		amount     = o.Amount                                 // amount	DECIMAL	YES
		fromSymbol = o.SymbolSource                           // fromSymbol	STRING	NO
		toSymbol   = o.SymbolTarget                           // toSymbol	STRING	NO
		params     = url.Values{}
	)

	params.Add("fromSymbol", fromSymbol)
	params.Add("toSymbol", toSymbol)
	params.Add("fromEmail", o.AccountSource)
	params.Add("toEmail", o.AccountTarget)
	if o.ActionUser == order.OrderMoveUserType_Master {
		if o.AccountSource != "" || o.AccountTarget != "" {
			//主账户向子账户划转：子母账户万能划转
			res, err = c.Api.SubAccountUniversalTransfer(spot_api.GetSubMoveSide(o.Source), spot_api.GetSubMoveSide(o.Target), asset, amount, &params)
		} else {
			return nil, errors.New(fmt.Sprint("OrderMoveUserType_OrderMove_Sub not supply operation from:", o.AccountSource, " to:", o.AccountTarget))
		}
	} else if o.ActionUser == order.OrderMoveUserType_Sub {
		if o.AccountTarget == "" {
			//子账户划转向主账户划转
			res, err = c.Api.SubAccountTransferSubToMaster(asset, amount)
		} else if o.AccountTarget != "" {
			//子账户向共同主账户下的子账户划转
			res, err = c.Api.SubAccountTransferSubToSub(o.AccountTarget, asset, amount)
		} else {
			return nil, errors.New(fmt.Sprint("OrderMoveUserType_OrderMove_Sub not supply operation from:", o.AccountSource, " to:", o.AccountTarget))
		}
	} else { // OrderMoveUserType_Internal
		//主用户万向划转
		res, err = c.Api.AssetTransfer(type_, asset, amount, &params)
	}

	if err != nil {
		return nil, err
	}
	resp.OrderId = strconv.FormatInt(res.TranId, 10)
	resp.RespType = client.OrderRspType_ACK
	return resp, err
}

func (c *ClientBinance) GetDepositHistory(req *client.DepositHistoryReq) (*client.DepositHistoryRsp, error) {
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

		res, err = c.Api.CapitalDepositHisRec(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range *res {
			amount, _ := strconv.ParseFloat(item.Amount, 64)
			resp.DepositList = append(resp.DepositList, &client.DepositHistoryItem{
				Asset:      item.Coin,
				Amount:     amount,
				Network:    spot_api.GetChainFromNetWork(item.Network),
				Status:     spot_api.GetDepositTypeFromExchange(spot_api.DepositStatus(item.Status)),
				Address:    item.Address,
				AddressTag: item.AddressTag,
				TxId:       item.TxId,
				Timestamp:  item.InsertTime,
			})
		}
		if len(*res) < limit {
			break
		}
		offset += limit
	}
	return resp, err
}

func (c *ClientBinance) GetTransferHistory(req *client.TransferHistoryReq) (*client.TransferHistoryRsp, error) {
	var (
		resp            = &client.TransferHistoryRsp{}
		err             error
		coin            = req.Asset                                      //coin	STRING	NO
		status          = spot_api.GetTransferTypeToExchange(req.Status) //status	INT	NO	0(0:已发送确认Email,1:已被用户取消 2:等待确认 3:被拒绝 4:处理中 5:提现交易失败 6 提现完成)
		withdrawOrderId string                                           //withdrawOrderId	STRING	NO
		offset          = 0                                              //offset	INT	NO
		limit           = 1000                                           //limit	INT	NO	默认：1000， 最大：1000
		startTime       = req.StartTime                                  //startTime	LONG	NO	默认当前时间90天前的时间戳
		endTime         = req.EndTime                                    //endTime	LONG	NO	默认当前时间戳
		params          = url.Values{}
		res             *spot_api.RespCapitalWithdrawHistory
	)

	if req.Id != 0 {
		withdrawOrderId = transform.IdToClientId(req.Producer, req.Id)
	}

	if coin != "" {
		params.Add("coin", coin)
	}
	if withdrawOrderId != "" {
		params.Add("withdrawOrderId", withdrawOrderId)
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
		res, err = c.Api.CapitalWithdrawHistory(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range *res {
			//_, id := transform.ClientIdToId(item.WithdrawOrderId)
			fee, _ := strconv.ParseFloat(item.TransactionFee, 64)
			timestamp, _ := strconv.ParseInt(item.ApplyTime, 10, 64)
			amount, _ := strconv.ParseFloat(item.Amount, 64)

			resp.TransferList = append(resp.TransferList, &client.TransferHistoryItem{
				Asset:         item.Coin,
				Amount:        amount,
				OrderId:       item.Id,
				Network:       spot_api.GetChainFromNetWork(item.Network),
				Status:        spot_api.GetTransferTypeFromExchange(spot_api.TransferStatus(item.Status)),
				Fee:           fee,
				ConfirmNumber: int64(item.ConfirmNo),
				Info:          item.Info,
				Address:       item.Address,
				TxId:          item.TxId,
				Timestamp:     timestamp,
			})
		}
		if len(*res) < limit {
			break
		}
		offset += limit
	}

	return resp, err
}

func (c *ClientBinance) GetMoveHistory(req *client.MoveHistoryReq) (*client.MoveHistoryRsp, error) {
	var (
		resp       = &client.MoveHistoryRsp{}
		err        error
		type_      = spot_api.GetMoveSide(req.Source, req.Target) // type	ENUM	YES
		fromSymbol = req.SymbolSource                             // fromSymbol	STRING	NO
		toSymbol   = req.SymbolTarget                             // toSymbol	STRING	NO
		startTime  = req.StartTime                                //startTime	LONG	NO
		endTime    = req.EndTime                                  //endTime	LONG	NO
		current    = 1                                            //current	INT	NO	默认 1
		size       = 100                                          //size	INT	NO	默认 10, 最大 100
		params     = url.Values{}
		res        *spot_api.RespAssetTransferHistory
	)
	if startTime > 0 {
		params.Add("startTime", strconv.FormatInt(startTime, 10))
	}
	if endTime > 0 {
		params.Add("endTime", strconv.FormatInt(endTime, 10))
	}

	if req.ActionUser == order.OrderMoveUserType_Master {
		// 子母账户万能划转历史
		params.Add("fromEmail", req.AccountSource)
		params.Add("toEmail", req.AccountTarget)
		for i := 0; i < 100; i++ {
			params.Add("limit", strconv.Itoa(size))
			params.Add("page", strconv.Itoa(current))
			res, err := c.Api.GetSubAccountUniversalTransfer(&params)
			if err != nil {
				return nil, err
			}

			for _, item := range res.Result {
				amount, _ := strconv.ParseFloat(item.Amount, 64)
				moveType := client.MoveType_MOVETYPE_IN
				if item.FromEmail == "" {
					moveType = client.MoveType_MOVETYPE_OUT
				}
				resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
					Asset:     item.Asset,
					Id:        strconv.FormatInt(item.TranId, 10),
					Type:      moveType,
					Amount:    amount,
					Timestamp: item.CreateTimeStamp,
					Status:    spot_api.GetMoveStatusFromExchange(spot_api.MoveStatus(item.Status)),
				})
			}
			if len(res.Result) < size {
				break
			}
			current++
		}
	} else if req.ActionUser == order.OrderMoveUserType_Sub {
		// 子账户划转历史, 时间段内的500条
		params.Add("limit", "500")
		params.Add("type", strconv.Itoa(spot_api.TRANSFER_IN))
		res, err := c.Api.SubAccountTransferSubUserHistory(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range *res {
			amount, _ := strconv.ParseFloat(item.Qty, 64)
			resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
				Asset:     item.Asset,
				Id:        strconv.FormatInt(item.TranId, 10),
				Amount:    amount,
				Timestamp: item.Time,
				Status:    spot_api.GetMoveStatusFromExchange(spot_api.MoveStatus(item.Status)),
				Type:      client.MoveType_MOVETYPE_IN,
			})
		}
		params.Set("type", strconv.Itoa(spot_api.TRANSFER_OUT))
		res, err = c.Api.SubAccountTransferSubUserHistory(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range *res {
			amount, _ := strconv.ParseFloat(item.Qty, 64)
			resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
				Asset:     item.Asset,
				Id:        strconv.FormatInt(item.TranId, 10),
				Amount:    amount,
				Timestamp: item.Time,
				Status:    spot_api.GetMoveStatusFromExchange(spot_api.MoveStatus(item.Status)),
				Type:      client.MoveType_MOVETYPE_OUT,
			})
		}
	} else { //OrderMoveUserType_Internal
		// 主账户划转历史
		params.Add("fromSymbol", fromSymbol)
		params.Add("toSymbol", toSymbol)
		for i := 0; i < 100; i++ {
			params.Add("current", strconv.Itoa(current))
			params.Add("size", strconv.Itoa(size))
			res, err = c.Api.AssetTransferHistory(type_, &params)
			if err != nil {
				return nil, err
			}

			for _, item := range res.Rows {
				amount, _ := strconv.ParseFloat(item.Amount, 64)
				resp.MoveList = append(resp.MoveList, &client.MoveHistoryItem{
					Asset:     item.Asset,
					Id:        strconv.FormatInt(item.TranId, 10),
					Amount:    amount,
					Timestamp: item.Timestamp,
					Status:    spot_api.GetMoveStatusFromExchange(spot_api.MoveStatus(item.Status)),
				})
			}
			if len(res.Rows) < size {
				break
			}
			current++
		}
	}
	return resp, err
}

func (c *ClientBinance) Loan(o *order.OrderLoan) (*client.OrderRsp, error) {
	var (
		resp       = &client.OrderRsp{}
		err        error
		asset      = o.Asset               //asset	STRING	YES
		isIsolated = "false"               //isIsolated	STRING	NO	是否逐仓杠杆，"TRUE", "FALSE", 默认 "FALSE"
		symbol     = string(o.Base.Symbol) //symbol	STRING	NO	逐仓交易对，配合逐仓使用
		amount     = o.Amount              //amount	DECIMAL	YES
		params     = url.Values{}
	)
	if o.Base.Type == common.SymbolType_MARGIN_ISOLATED {
		isIsolated = "true"
	}
	params.Add("isIsolated", isIsolated)
	params.Add("symbol", symbol)

	res, err := c.Api.MarginLoan(asset, amount, &params)
	if err != nil {
		return nil, err
	}
	resp.OrderId = strconv.FormatInt(res.TranId, 10)
	resp.RespType = client.OrderRspType_ACK
	return resp, err
}

func (c *ClientBinance) Return(o *order.OrderReturn) (*client.OrderRsp, error) {
	var (
		resp       = &client.OrderRsp{}
		err        error
		asset      = o.Asset               //asset	STRING	YES
		isIsolated = "false"               //isIsolated	STRING	NO	是否逐仓杠杆，"TRUE", "FALSE", 默认 "FALSE"
		symbol     = string(o.Base.Symbol) //symbol	STRING	NO	逐仓交易对，配合逐仓使用
		amount     = o.Amount              //amount	DECIMAL	YES
		params     = url.Values{}
	)
	if o.Base.Type == common.SymbolType_MARGIN_ISOLATED {
		isIsolated = "true"
	}
	params.Add("isIsolated", isIsolated)
	params.Add("symbol", symbol)

	res, err := c.Api.MarginRepay(asset, amount, &params)
	if err != nil {
		return nil, err
	}
	resp.OrderId = strconv.FormatInt(res.TranId, 10)
	resp.RespType = client.OrderRspType_ACK
	return resp, err
}

func (c *ClientBinance) GetLoanOrders(o *client.LoanHistoryReq) (*client.LoanHistoryRsp, error) {
	var (
		resp           = &client.LoanHistoryRsp{}
		err            error
		asset          = o.Asset          //asset	STRING	YES
		isolatedSymbol = o.IsolatedSymbol //isolatedSymbol	STRING	NO	逐仓symbol
		txId           = o.IdEx           //txId	LONG	NO	tranId in POST /sapi/v1/margin/loan
		startTime      = o.StartTime      //startTime	LONG	NO
		endTime        = o.EndTime        //endTime	LONG	NO
		current        = 1                //current	LONG	NO	当前查询页。 开始值 1。 默认:1
		size           = 100              //size	LONG	NO	默认:10 最大:100
		archived       = "false"          //archived	STRING	NO	默认: false. 查询6个月以前的数据，需要设为 true
		params         = url.Values{}
		res            *spot_api.RespMarginLoanHistory
	)

	if isolatedSymbol != "" {
		params.Add("isolatedSymbol", isolatedSymbol)
	}
	if txId != "" {
		params.Add("txId", txId)
	}
	if startTime > 0 {
		params.Add("startTime", strconv.FormatInt(startTime, 10))
	}
	if endTime > 0 {
		params.Add("endTime", strconv.FormatInt(endTime, 10))
	}
	params.Add("archived", archived)

	for i := 0; i < 100; i++ {
		params.Add("current", strconv.Itoa(current))
		params.Add("size", strconv.Itoa(size))
		res, err = c.Api.MarginLoanHistory(asset, &params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Rows {
			principal, _ := strconv.ParseFloat(item.Principal, 64)
			resp.LoadList = append(resp.LoadList, &client.LoanHistoryItem{
				IsolatedSymbol: item.IsolatedSymbol,
				OrderId:        strconv.FormatInt(item.TxId, 10),
				Asset:          item.Asset,
				Principal:      principal,
				Timestamp:      item.Timestamp,
				Status:         spot_api.GetLoadStatusFromExchange(spot_api.LoanStatus(item.Status)),
			})
		}
		if len(res.Rows) < size {
			break
		}
		current++
	}
	return resp, err
}

func (c *ClientBinance) GetReturnOrders(o *client.LoanHistoryReq) (*client.ReturnHistoryRsp, error) {
	var (
		resp           = &client.ReturnHistoryRsp{}
		err            error
		asset          = o.Asset          //asset	STRING	YES
		isolatedSymbol = o.IsolatedSymbol //isolatedSymbol	STRING	NO	逐仓symbol
		txId           = o.IdEx           //txId	LONG	NO	tranId in POST /sapi/v1/margin/loan
		startTime      = o.StartTime      //startTime	LONG	NO
		endTime        = o.EndTime        //endTime	LONG	NO
		current        = 1                //current	LONG	NO	当前查询页。 开始值 1。 默认:1
		size           = 100              //size	LONG	NO	默认:10 最大:100
		archived       = "false"          //archived	STRING	NO	默认: false. 查询6个月以前的数据，需要设为 true
		params         = url.Values{}
		res            *spot_api.RespMarginRepayHistory
	)
	if isolatedSymbol != "" {
		params.Add("isolatedSymbol", isolatedSymbol)
	}
	if txId != "" {
		params.Add("txId", txId)
	}
	if startTime > 0 {
		params.Add("startTime", strconv.FormatInt(startTime, 10))
	}
	if endTime > 0 {
		params.Add("endTime", strconv.FormatInt(endTime, 10))
	}
	params.Add("archived", archived)

	for i := 0; i < 100; i++ {
		params.Add("current", strconv.Itoa(current))
		params.Add("size", strconv.Itoa(size))

		res, err = c.Api.MarginRepayHistory(asset, &params)
		if err != nil {
			return nil, err
		}
		for _, item := range res.Rows {
			principal, _ := strconv.ParseFloat(item.Principal, 64)
			amount, _ := strconv.ParseFloat(item.Amount, 64)
			interest, _ := strconv.ParseFloat(item.Interest, 64)
			resp.ReturnList = append(resp.ReturnList, &client.ReturnHistoryItem{
				IsolatedSymbol: item.IsolatedSymbol,
				Amount:         amount,
				OrderId:        strconv.FormatInt(item.TxId, 10),
				Interest:       interest,
				Asset:          item.Asset,
				Principal:      principal,
				Timestamp:      item.Timestamp,
				Status:         spot_api.GetLoadStatusFromExchange(spot_api.LoanStatus(item.Status)),
			})
		}
		if len(res.Rows) < size {
			break
		}
		current++
	}
	return resp, err
}

func GetCloseTime(symbol string) string {
	var closeDate string
	names := strings.Split(symbol, "_")
	if len(names) > 1 && len(names[1]) == 6 {
		closeDate = names[1]
	}
	return closeDate
}

func (c *ClientBinance) GetOrder(req *order.OrderQueryReq) (*client.OrderRsp, error) {
	var (
		resp    = &client.OrderRsp{}
		symbol  = strings.ToUpper(strings.Replace(string(req.Symbol), "/", "", -1)) //symbol	STRING	YES
		orderId = req.IdEx                                                          //orderId	LONG	NO
		params  = url.Values{}
	)
	if string(req.IdEx) != "" {
		params.Add("orderId", orderId)
	} else if req.Id != 0 {
		params.Add("origClientOrderId", transform.IdToClientId(req.Producer, req.Id))
	} else {
		return nil, errors.New("id can not be empty")
	}
	if req.Market == common.Market_SPOT {
		res, err := c.Api.GetOrder(symbol, &params)
		if err != nil {
			return nil, err
		}
		cummulativeQuoteQty, _ := strconv.ParseFloat(res.CummulativeQuoteQty, 64)
		executedQty, _ := strconv.ParseFloat(res.ExecutedQty, 64)
		if cummulativeQuoteQty > 0 {
			resp.Price = executedQty / cummulativeQuoteQty
		} else {
			resp.Price, _ = strconv.ParseFloat(res.Price, 64)
		}
		resp.Producer, resp.Id = transform.ClientIdToId(res.ClientOrderId)
		resp.OrderId = strconv.Itoa(res.OrderId)
		resp.Timestamp = res.Time * 1000
		resp.RespType = client.OrderRspType_RESULT
		resp.Symbol = spot_api.ParseSymbolName(res.Symbol)
		resp.Status = spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(res.Status))
		resp.AccumAmount, _ = strconv.ParseFloat(res.ExecutedQty, 64)
		resp.AccumQty, _ = strconv.ParseFloat(res.CummulativeQuoteQty, 64)
		resp.CloseDate = GetCloseTime(res.Symbol)
		// binance没有平均价，不用设置avgPrice
	} else if req.Market == common.Market_MARGIN {
		res, err := c.Api.GetMarginOrder(symbol, &params)
		if err != nil {
			return nil, err
		}
		cummulativeQuoteQty, _ := strconv.ParseFloat(res.CummulativeQuoteQty, 64)
		executedQty, _ := strconv.ParseFloat(res.ExecutedQty, 64)
		if cummulativeQuoteQty > 0 {
			resp.Price = executedQty / cummulativeQuoteQty
		} else {
			resp.Price, _ = strconv.ParseFloat(res.Price, 64)
		}
		resp.Producer, resp.Id = transform.ClientIdToId(res.ClientOrderId)
		resp.OrderId = strconv.Itoa(res.OrderId)
		resp.Timestamp = res.Time * 1000
		resp.RespType = client.OrderRspType_RESULT
		resp.Symbol = spot_api.ParseSymbolName(res.Symbol)
		resp.Status = spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(res.Status))
		resp.AccumAmount, _ = strconv.ParseFloat(res.ExecutedQty, 64)
		resp.AccumQty, _ = strconv.ParseFloat(res.CummulativeQuoteQty, 64)
		resp.CloseDate = GetCloseTime(res.Symbol)
		// binance没有平均价，不用设置avgPrice
	} else if spot_api.IsUBaseSymbolType(req.Type) {
		res, err := c.UApi.GetOrder(u_api.GetUBaseSymbol(symbol, u_api.GetFutureTypeFromNats(req.Type)), &params)
		if err != nil {
			return nil, err
		}
		cummulativeQuoteQty, _ := strconv.ParseFloat(res.CummulativeQuoteQty, 64)
		executedQty, _ := strconv.ParseFloat(res.ExecutedQty, 64)
		if cummulativeQuoteQty > 0 {
			resp.Price = executedQty / cummulativeQuoteQty
		} else {
			resp.Price, _ = strconv.ParseFloat(res.Price, 64)
		}
		resp.Producer, resp.Id = transform.ClientIdToId(res.ClientOrderId)
		resp.OrderId = strconv.Itoa(res.OrderId)
		resp.Timestamp = res.Time * 1000
		resp.RespType = client.OrderRspType_RESULT
		resp.Symbol = spot_api.ParseSymbolName(res.Symbol)
		resp.Status = spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(res.Status))
		resp.AccumAmount, _ = strconv.ParseFloat(res.ExecutedQty, 64)
		resp.AccumQty, _ = strconv.ParseFloat(res.CumQuote, 64)
		resp.AvgPrice, _ = strconv.ParseFloat(res.AvgPrice, 64)
		resp.CloseDate = GetCloseTime(res.Symbol)
	} else if spot_api.IsCBaseSymbolType(req.Type) {
		res, err := c.CApi.GetOrder(c_api.GetCBaseSymbol(symbol, u_api.GetFutureTypeFromNats(req.Type)), &params)
		if err != nil {
			return nil, err
		}
		cummulativeQuoteQty, _ := strconv.ParseFloat(res.CummulativeQuoteQty, 64)
		executedQty, _ := strconv.ParseFloat(res.ExecutedQty, 64)
		if cummulativeQuoteQty > 0 {
			resp.Price = executedQty / cummulativeQuoteQty
		} else {
			resp.Price, _ = strconv.ParseFloat(res.Price, 64)
		}
		resp.Producer, resp.Id = transform.ClientIdToId(res.ClientOrderId)
		resp.OrderId = strconv.Itoa(res.OrderId)
		resp.Timestamp = res.Time * 1000
		resp.RespType = client.OrderRspType_RESULT
		resp.Symbol = spot_api.ParseSymbolName(res.Symbol)
		resp.Status = spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(res.Status))
		resp.AccumAmount, _ = strconv.ParseFloat(res.ExecutedQty, 64)
		resp.AccumQty, _ = strconv.ParseFloat(res.CumQuote, 64)
		resp.AvgPrice, _ = strconv.ParseFloat(res.AvgPrice, 64)
		resp.CloseDate = GetCloseTime(res.Symbol)
	} else {
		return nil, fmt.Errorf("invalid market/type %v:%v", req.Market, req.Type)
	}

	return resp, nil
}

func (c *ClientBinance) GetOrderHistory(req *client.OrderHistoryReq) ([]*client.OrderRsp, error) {
	var (
		resp      []*client.OrderRsp
		symbol    = spot_api.GetSymbol(req.Asset) //symbol	STRING	YES
		market    = req.Market
		isolated  = "FALSE"
		startTime = req.StartTime //startTime	LONG	NO
		endTime   = req.EndTime   //endTime	LONG	NO
		limit     = 1000          //limit	INT	NO	默认 500; 最大 1000.
		params    = url.Values{}
	)
	if req.Isolated {
		isolated = "True"
	}
	if startTime > 0 {
		params.Add("startTime", strconv.FormatInt(startTime, 10))
	}
	if endTime > 0 {
		params.Add("endTime", strconv.FormatInt(endTime, 10))
	}
	params.Add("limit", strconv.Itoa(limit))

	for i := 0; i < 100; i++ {
		if market == common.Market_SPOT {
			res, err := c.Api.AllOrders(symbol, &params)
			if err != nil {
				return nil, err
			}
			for _, item := range *res {
				producer, id := transform.ClientIdToId(item.ClientOrderId)
				cummulativeQuoteQty, _ := strconv.ParseFloat(item.CummulativeQuoteQty, 64)
				executed, _ := strconv.ParseFloat(item.ExecutedQty, 64)
				resp = append(resp, &client.OrderRsp{
					Producer:    producer,
					Id:          id,
					OrderId:     strconv.Itoa(item.OrderId),
					Timestamp:   item.Time * 1000,
					RespType:    client.OrderRspType_RESULT,
					Symbol:      item.Symbol,
					Status:      spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(item.Status)),
					AccumAmount: executed,
					AccumQty:    cummulativeQuoteQty,
					CloseDate:   GetCloseTime(item.Symbol),
					// binance没有平均价，不用设置avgPrice
				})
			}
			params.Add("orderId", strconv.Itoa((*res)[len(*res)-1].OrderId))
			if len(*res) < limit {
				break
			}
		} else if market == common.Market_MARGIN {
			params.Add("isolated", isolated)
			res, err := c.Api.MarginAllOrders(symbol, &params)
			if err != nil {
				return nil, err
			}
			for _, item := range *res {
				producer, id := transform.ClientIdToId(item.ClientOrderId)
				cummulativeQuoteQty, _ := strconv.ParseFloat(item.CummulativeQuoteQty, 64)
				executed, _ := strconv.ParseFloat(item.ExecutedQty, 64)
				resp = append(resp, &client.OrderRsp{
					Producer:    producer,
					Id:          id,
					OrderId:     strconv.Itoa(item.OrderId),
					Timestamp:   item.Time * 1000,
					RespType:    client.OrderRspType_RESULT,
					Symbol:      item.Symbol,
					Status:      spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(item.Status)),
					AccumAmount: executed,
					AccumQty:    cummulativeQuoteQty,
					CloseDate:   GetCloseTime(item.Symbol),
					// binance没有平均价，不用设置avgPrice
				})
			}
			params.Add("orderId", strconv.Itoa((*res)[len(*res)-1].OrderId))
			if len(*res) < limit {
				break
			}
		} else if spot_api.IsUBaseSymbolType(req.Type) {
			res, err := c.UApi.AllOrders(symbol, &params)
			if err != nil {
				return nil, err
			}
			for _, item := range *res {
				producer, id := transform.ClientIdToId(item.ClientOrderId)
				cummulativeQuoteQty, _ := strconv.ParseFloat(item.CumQuote, 64)
				executed, _ := strconv.ParseFloat(item.ExecutedQty, 64)
				avgPrice, _ := strconv.ParseFloat(item.AvgPrice, 64)
				resp = append(resp, &client.OrderRsp{
					Producer:    producer,
					Id:          id,
					OrderId:     strconv.Itoa(item.OrderId),
					Timestamp:   item.Time * 1000,
					RespType:    client.OrderRspType_RESULT,
					Symbol:      item.Symbol,
					Status:      spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(item.Status)),
					AccumAmount: executed,
					AccumQty:    cummulativeQuoteQty,
					AvgPrice:    avgPrice,
					CloseDate:   GetCloseTime(item.Symbol),
				})
			}
			params.Add("orderId", strconv.Itoa((*res)[len(*res)-1].OrderId))
			if len(*res) < limit {
				break
			}
		} else if spot_api.IsCBaseSymbolType(req.Type) {
			res, err := c.CApi.AllOrders(symbol, &params)
			if err != nil {
				return nil, err
			}
			for _, item := range *res {
				producer, id := transform.ClientIdToId(item.ClientOrderId)
				cummulativeQuoteQty, _ := strconv.ParseFloat(item.CumQuote, 64)
				executed, _ := strconv.ParseFloat(item.ExecutedQty, 64)
				avgPrice, _ := strconv.ParseFloat(item.AvgPrice, 64)
				resp = append(resp, &client.OrderRsp{
					Producer:    producer,
					Id:          id,
					OrderId:     strconv.Itoa(item.OrderId),
					Timestamp:   item.Time * 1000,
					RespType:    client.OrderRspType_RESULT,
					Symbol:      item.Symbol,
					Status:      spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(item.Status)),
					AccumAmount: executed,
					AccumQty:    cummulativeQuoteQty,
					AvgPrice:    avgPrice,
					CloseDate:   GetCloseTime(item.Symbol),
				})
			}
			params.Add("orderId", strconv.Itoa((*res)[len(*res)-1].OrderId))
			if len(*res) < limit {
				break
			}
		} else {
			return nil, fmt.Errorf("invalid market/type %v:%v", req.Market, req.Type)
		}
	}

	return resp, nil
}

func (c *ClientBinance) GetProcessingOrders(req *client.OrderHistoryReq) ([]*client.OrderRsp, error) {
	var (
		resp     []*client.OrderRsp
		symbol   = req.Asset //symbol	STRING	NO
		market   = req.Market
		isolated = "FALSE"
		params   = url.Values{}
	)
	if req.Isolated {
		isolated = "True"
	}
	params.Add("symbol", symbol)
	if market == common.Market_SPOT {
		res, err := c.Api.OpenOrders(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range *res {
			producer, id := transform.ClientIdToId(item.ClientOrderId)
			cummulativeQuoteQty, _ := strconv.ParseFloat(item.CummulativeQuoteQty, 64)
			executed, _ := strconv.ParseFloat(item.ExecutedQty, 64)
			resp = append(resp, &client.OrderRsp{
				Producer:    producer,
				Id:          id,
				OrderId:     strconv.Itoa(item.OrderId),
				Timestamp:   item.Time * 1000,
				Symbol:      item.Symbol,
				RespType:    client.OrderRspType_RESULT,
				Status:      spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(item.Status)),
				AccumAmount: executed,
				AccumQty:    cummulativeQuoteQty,
				CloseDate:   GetCloseTime(item.Symbol),
				// binance没有平均价，不用设置avgPrice
			})
		}
	} else if market == common.Market_MARGIN {
		params.Add("isolated", isolated)
		res, err := c.Api.MarginOpenOrders(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range *res {
			producer, id := transform.ClientIdToId(item.ClientOrderId)
			cummulativeQuoteQty, _ := strconv.ParseFloat(item.CummulativeQuoteQty, 64)
			executed, _ := strconv.ParseFloat(item.ExecutedQty, 64)
			resp = append(resp, &client.OrderRsp{
				Producer:    producer,
				Id:          id,
				OrderId:     strconv.Itoa(item.OrderId),
				Timestamp:   item.Time * 1000,
				RespType:    client.OrderRspType_RESULT,
				Symbol:      item.Symbol,
				Status:      spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(item.Status)),
				AccumAmount: executed,
				AccumQty:    cummulativeQuoteQty,
				CloseDate:   GetCloseTime(item.Symbol),
				// binance没有平均价，不用设置avgPrice
			})
		}
	} else if spot_api.IsUBaseSymbolType(req.Type) {
		if symbol != "" {
			params.Set("symbol", u_api.GetUBaseSymbol(symbol, u_api.GetFutureTypeFromNats(req.Type)))
		}
		res, err := c.UApi.GetOpenOrders(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range res {
			producer, id := transform.ClientIdToId(item.ClientOrderId)
			cummulativeQuoteQty, _ := strconv.ParseFloat(item.CumQuote, 64)
			executed, _ := strconv.ParseFloat(item.ExecutedQty, 64)
			avgPrice, _ := strconv.ParseFloat(item.AvgPrice, 64)
			resp = append(resp, &client.OrderRsp{
				Producer:    producer,
				Id:          id,
				OrderId:     strconv.Itoa(item.OrderId),
				Timestamp:   item.Time * 1000,
				Symbol:      item.Symbol,
				RespType:    client.OrderRspType_RESULT,
				Status:      spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(item.Status)),
				AccumAmount: executed,
				AccumQty:    cummulativeQuoteQty,
				AvgPrice:    avgPrice,
				CloseDate:   GetCloseTime(item.Symbol),
			})
		}
	} else if spot_api.IsCBaseSymbolType(req.Type) {
		if symbol != "" {
			params.Set("symbol", c_api.GetCBaseSymbol(symbol, u_api.GetFutureTypeFromNats(req.Type)))
		}
		res, err := c.CApi.GetOpenOrders(&params)
		if err != nil {
			return nil, err
		}
		for _, item := range res {
			producer, id := transform.ClientIdToId(item.ClientOrderId)
			cummulativeQuoteQty, _ := strconv.ParseFloat(item.CumQuote, 64)
			executed, _ := strconv.ParseFloat(item.ExecutedQty, 64)
			avgPrice, _ := strconv.ParseFloat(item.AvgPrice, 64)
			resp = append(resp, &client.OrderRsp{
				Producer:    producer,
				Id:          id,
				OrderId:     strconv.Itoa(item.OrderId),
				Timestamp:   item.Time * 1000,
				Symbol:      item.Symbol,
				RespType:    client.OrderRspType_RESULT,
				Status:      spot_api.GetOrderStatusFromExchange(spot_api.OrderStatus(item.Status)),
				AccumAmount: executed,
				AccumQty:    cummulativeQuoteQty,
				AvgPrice:    avgPrice,
				CloseDate:   GetCloseTime(item.Symbol),
			})
		}
	} else {
		return nil, fmt.Errorf("invalid market/type %v:%v", req.Market, req.Type)
	}
	return resp, nil
}

func (c *ClientBinance) SubAccountList() (*spot_api.RespSubAccountList, error) {
	return c.Api.SubAccountList(nil)
}
