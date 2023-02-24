package u_api

import (
	"clients/exchange/cex/binance/spot_api"
	"clients/transform"
	"fmt"
	"github.com/warmplanet/proto/go/common"
	"strings"
)

type ContractType string
type ContractStatus string
type PositionSide string
type WorkingType string

const (
	INVALID_CONTRACTTYPE ContractType = ""
	PERPETUAL            ContractType = "PERPETUAL"            // 永续合约
	CURRENT_MONTH        ContractType = "CURRENT_MONTH"        // 当月交割合约
	NEXT_MONTH           ContractType = "NEXT_MONTH"           // 次月交割合约
	CURRENT_QUARTER      ContractType = "CURRENT_QUARTER"      // 当季交割合约
	NEXT_QUARTER         ContractType = "NEXT_QUARTER"         // 次季交割合约
	PERPETUAL_DELIVERING ContractType = "PERPETUAL_DELIVERING" // 交割结算中合约
)

func GetUBaseSymbol(symbol string, type_ ContractType) string {
	/*
		https://www.binance.com/zh-CN/support/faq/3ae441db4ae740e19af3fe9228eb6619
		季度交割合约是具备固定到期日和交割日的衍生品合约，以每个季度的最后一个周五作爲交割日。例如：“BTCUSDT 当季 0326”代表 2021年03月26日16:00（香港时间）进行交割。
		当季度合约结算交割后，会产生新的季度合约，例如：“BTCUSDT 当季 0326” 于 2021年03月26日16:00（香港时间）交割下架后将会生成新的“BTCUSDT 当季 0625”合约，以此类推；
		交割日结算时将收取交割费，交割费与Taker(吃单)费率相同。
	*/
	symbol = spot_api.GetSymbol(symbol)
	switch type_ {
	case PERPETUAL:
		return fmt.Sprint(symbol)
	case CURRENT_MONTH:
		date := transform.GetDate(transform.THISMONTH)
		return fmt.Sprint(symbol, "_", date)
	case NEXT_MONTH:
		date := transform.GetDate(transform.NEXTMONTH)
		return fmt.Sprint(symbol, "_", date)
	case CURRENT_QUARTER:
		date := transform.GetDate(transform.THISQUARTER)
		return fmt.Sprint(symbol, "_", date)
	case NEXT_QUARTER:
		date := transform.GetDate(transform.NEXTQUARTER)
		return fmt.Sprint(symbol, "_", date)
	case PERPETUAL_DELIVERING:
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
	if len(symList) <= 1 {
		return spot_api.ParseSymbolName(strings.ToUpper(symList[0])), common.Market_SWAP, common.SymbolType_SWAP_FOREVER
	} else if strings.ToLower(symList[1]) == "perp" {
		return spot_api.ParseSymbolName(strings.ToUpper(symList[0])), common.Market_SWAP_COIN, common.SymbolType_SWAP_COIN_FOREVER
	} else {
		sym := spot_api.ParseSymbolName(strings.ToUpper(symList[0]))
		switch symbol {
		//case GetUBaseSymbol(symList[0], CURRENT_MONTH):
		//	return sym, common.Market_FUTURE, common.SymbolType_FUTURE_THIS_MONTH
		//case GetUBaseSymbol(symList[0], NEXT_MONTH):
		//	return sym, common.Market_FUTURE, common.SymbolType_FUTURE_NEXT_MONTH
		case GetUBaseSymbol(symList[0], CURRENT_QUARTER):
			return sym, common.Market_FUTURE, common.SymbolType_FUTURE_THIS_QUARTER
		case GetUBaseSymbol(symList[0], NEXT_QUARTER):
			return sym, common.Market_FUTURE, common.SymbolType_FUTURE_NEXT_QUARTER
		default:
			return sym, common.Market_INVALID_MARKET, common.SymbolType_INVALID_TYPE
		}
	}
}

func GetFutureTypeFromExchange(type_ ContractType) common.SymbolType {
	switch type_ {
	case PERPETUAL:
		return common.SymbolType_SWAP_FOREVER
	case CURRENT_MONTH:
		return common.SymbolType_FUTURE_THIS_MONTH
	case NEXT_MONTH:
		return common.SymbolType_FUTURE_NEXT_MONTH
	case CURRENT_QUARTER:
		return common.SymbolType_FUTURE_THIS_QUARTER
	case NEXT_QUARTER:
		return common.SymbolType_FUTURE_NEXT_QUARTER
	case PERPETUAL_DELIVERING:
		fallthrough
	default:
		return common.SymbolType_INVALID_TYPE
	}
}
func GetFutureTypeFromNats(type_ common.SymbolType) ContractType {
	switch type_ {
	case common.SymbolType_SWAP_FOREVER, common.SymbolType_SWAP_COIN_FOREVER:
		return PERPETUAL
	case common.SymbolType_FUTURE_THIS_MONTH, common.SymbolType_FUTURE_COIN_THIS_MONTH:
		return CURRENT_MONTH
	case common.SymbolType_FUTURE_NEXT_MONTH, common.SymbolType_FUTURE_COIN_NEXT_MONTH:
		return NEXT_MONTH
	case common.SymbolType_FUTURE_THIS_QUARTER, common.SymbolType_FUTURE_COIN_THIS_QUARTER:
		return CURRENT_QUARTER
	case common.SymbolType_FUTURE_NEXT_QUARTER, common.SymbolType_FUTURE_COIN_NEXT_QUARTER:
		return NEXT_QUARTER
	default:
		return INVALID_CONTRACTTYPE
	}
}

const (
	PENDING_TRADING ContractStatus = "PENDING_TRADING" // 待上市
	TRADING         ContractStatus = "TRADING"         // 交易中
	PRE_DELIVERING  ContractStatus = "PRE_DELIVERING"  // 预交割
	DELIVERING      ContractStatus = "DELIVERING"      // 交割中
	DELIVERED       ContractStatus = "DELIVERED"       // 已交割
	PRE_SETTLE      ContractStatus = "PRE_SETTLE"      // 预结算
	SETTLING        ContractStatus = "SETTLING"        // 结算中
	CLOSE           ContractStatus = "CLOSE"           // 已下架
)

const (
	BOTH  PositionSide = "BOTH"  // 单一持仓方向
	LONG  PositionSide = "LONG"  // 多头(双向持仓下)
	SHORT PositionSide = "SHORT" // 空头(双向持仓下)
)

const (
	MARK_PRICE     WorkingType = "MARK_PRICE"     //条件价格触发
	CONTRACT_PRICE WorkingType = "CONTRACT_PRICE" //条件价格触发
)
