package c_api

import (
	"clients/exchange/cex/binance/spot_api"
	"clients/exchange/cex/binance/u_api"
	"clients/transform"
	"fmt"
	"github.com/warmplanet/proto/go/common"
	"strings"
)

func GetCBaseSymbol(symbol string, type_ u_api.ContractType) string {
	/*
		https://www.binance.com/zh-CN/support/faq/3ae441db4ae740e19af3fe9228eb6619
		季度交割合约是具备固定到期日和交割日的衍生品合约，以每个季度的最后一个周五作爲交割日。例如：“BTCUSDT 当季 0326”代表 2021年03月26日16:00（香港时间）进行交割。
		当季度合约结算交割后，会产生新的季度合约，例如：“BTCUSDT 当季 0326” 于 2021年03月26日16:00（香港时间）交割下架后将会生成新的“BTCUSDT 当季 0625”合约，以此类推；
		交割日结算时将收取交割费，交割费与Taker(吃单)费率相同。
	*/
	symbol = spot_api.GetSymbol(symbol)
	switch type_ {
	case u_api.PERPETUAL:
		return fmt.Sprint(symbol, "_PERP")
	case u_api.CURRENT_MONTH:
		date := transform.GetDate(transform.THISMONTH)
		return fmt.Sprint(symbol, "_", date)
	case u_api.NEXT_MONTH:
		date := transform.GetDate(transform.NEXTMONTH)
		return fmt.Sprint(symbol, "_", date)
	case u_api.CURRENT_QUARTER:
		date := transform.GetDate(transform.THISQUARTER)
		return fmt.Sprint(symbol, "_", date)
	case u_api.NEXT_QUARTER:
		date := transform.GetDate(transform.NEXTQUARTER)
		return fmt.Sprint(symbol, "_", date)
	case u_api.PERPETUAL_DELIVERING:
		fallthrough
	default:
		return symbol
	}
}

func GetContractType(symbol string) (string, common.Market, common.SymbolType) {
	/*
		https://www.binance.com/zh-CN/support/faq/3ae441db4ae740e19af3fe9228eb6619
		季度交割合约是具备固定到期日和交割日的衍生品合约，以每个季度的最后一个周五作爲交割日。例如：“BTCUSDT 当季 0326”代表 2021年03月26日16:00（香港时间）进行交割。
		当季度合约结算交割后，会产生新的季度合约，例如：“BTCUSDT 当季 0326” 于 2021年03月26日16:00（香港时间）交割下架后将会生成新的“BTCUSDT 当季 0625”合约，以此类推；
		交割日结算时将收取交割费，交割费与Taker(吃单)费率相同。
	*/
	symList := strings.Split(symbol, "_")
	if len(symList) == 1 {
		return spot_api.ParseSymbolName(strings.ToUpper(symList[0])), common.Market_SWAP, common.SymbolType_SWAP_FOREVER
	} else if strings.ToLower(symList[1]) == "perp" {
		return spot_api.ParseSymbolName(strings.ToUpper(symList[0])), common.Market_SWAP_COIN, common.SymbolType_SWAP_COIN_FOREVER
	} else {
		sym := spot_api.ParseSymbolName(strings.ToUpper(symList[0]))
		switch symbol {
		case GetCBaseSymbol(symList[0], u_api.PERPETUAL):
			return sym, common.Market_SWAP_COIN, common.SymbolType_SWAP_COIN_FOREVER
		//case GetCBaseSymbol(symList[0], u_api.CURRENT_MONTH):
		//	return sym, common.Market_FUTURE_COIN, common.SymbolType_FUTURE_COIN_THIS_MONTH
		//case GetCBaseSymbol(symList[0], u_api.NEXT_MONTH):
		//	return sym, common.Market_FUTURE_COIN, common.SymbolType_FUTURE_COIN_NEXT_MONTH
		case GetCBaseSymbol(symList[0], u_api.CURRENT_QUARTER):
			return sym, common.Market_FUTURE_COIN, common.SymbolType_FUTURE_COIN_THIS_QUARTER
		case GetCBaseSymbol(symList[0], u_api.NEXT_QUARTER):
			return sym, common.Market_FUTURE_COIN, common.SymbolType_FUTURE_COIN_NEXT_QUARTER
		default:
			return sym, common.Market_INVALID_MARKET, common.SymbolType_INVALID_TYPE
		}
	}
}
