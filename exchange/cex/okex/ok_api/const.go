package ok_api

import (
	"clients/transform"
	"fmt"
	"github.com/warmplanet/proto/go/common"
	"strings"
)

type OkexTransferStatus int
type OrderType string
type MoveType string
type SideType string
type TimeInForceType string //有效方式
type OrderRespType string   //响应
type MarginType string
type OrderStatus string
type DepositStatus int
type TransferStatus int
type MoveStatus string
type LoanStatus string

const (
	//提币状态
	//-3：撤销中 -2：已撤销 -1：失败
	//0：等待提现 1：提现中 2：已汇出
	//3：邮箱确认 4：人工审核中 5：等待身份认证
	OK_TRANSFER_TYPE_REVOCATION             OkexTransferStatus = -3
	OK_TRANSFER_TYPE_WITHDRWAN              OkexTransferStatus = -2
	OK_TRANSFER_TYPE_FAIL                   OkexTransferStatus = -1
	OK_TRANSFER_TYPE_CONFIRMING             OkexTransferStatus = 0
	OK_TRANSFER_TYPE_PROCESSING             OkexTransferStatus = 1
	OK_TRANSFER_RYPE_REMITTED               OkexTransferStatus = 2
	OK_TRANSFER_TYPE_CREATED                OkexTransferStatus = 3
	OK_TRANSFER_TYPE_MANUALREVIEW           OkexTransferStatus = 4
	OK_TRANSFER_TYPE_IDENTITYAUTHENTICATION OkexTransferStatus = 5
)

const (
	ORDER_TYPE_LIMIT             OrderType = "LIMIT"             //限价单
	ORDER_TYPE_MARKET            OrderType = "MARKET"            //市价单
	ORDER_TYPE_STOP_LOSS         OrderType = "STOP_LOSS"         //止损单
	ORDER_TYPE_STOP_LOSS_LIMIT   OrderType = "STOP_LOSS_LIMIT"   //限价止损单
	ORDER_TYPE_TAKE_PROFIT       OrderType = "TAKE_PROFIT"       //止盈单
	ORDER_TYPE_TAKE_PROFIT_LIMIT OrderType = "TAKE_PROFIT_LIMIT" //限价止盈单
	ORDER_TYPE_LIMIT_MAKER       OrderType = "LIMIT_MAKER"       //限价只挂单
	ORDER_TYPE_POST_ONLY         OrderType = "POST_ONLY"
	ORDER_TYPE_FOK               OrderType = "FOK"
	ORDER_TYPE_IOC               OrderType = "IOC"
	ORDER_TYPE_OPTIMAL_LIMIT_IOC OrderType = "OPTIMAL_LIMIT_IOC"
)

const (
	TIME_IN_FORCE_GTC               TimeInForceType = "GTC"               //成交为止 订单会一直有效，直到被成交或者取消。
	TIME_IN_FORCE_IOC               TimeInForceType = "IOC"               //无法立即成交的部分就撤销 订单在失效前会尽量多的成交。
	TIME_IN_FORCE_FOK               TimeInForceType = "FOK"               //无法全部立即成交就撤销 如果无法全部成交，订单会失效。
	TIME_IN_FORCE_GTX               TimeInForceType = "GTX"               //- Good Till Crossing 无法成为挂单方就撤销
	TIME_IN_FORCE_OPTIMAL_LIMIT_IOC TimeInForceType = "OPTIMAL_LIMIT_IOC" //市价委托立即成交并取消剩余（仅适用交割、永续）
)

const (
	SIDE_TYPE_BUY  SideType = "buy"  // 买入
	SIDE_TYPE_SELL SideType = "sell" //卖出
)

const (
	ORDER_STATUS_NEW              OrderStatus = "NEW"
	ORDER_STATUS_PARTIALLY_FILLED OrderStatus = "PARTIALLY_FILLED"
	ORDER_STATUS_FILLED           OrderStatus = "FILLED"
	ORDER_STATUS_CANCELED         OrderStatus = "CANCELED"
	ORDER_STATUS_LIVE             OrderStatus = "LIVE"
	ORDER_STATUS_PENDING_CANCEL   OrderStatus = "PENDING_CANCEL"
	ORDER_STATUS_REJECTED         OrderStatus = "REJECTED"
	ORDER_STATUS_EXPIRED          OrderStatus = "EXPIRED"
)

const (
	MARGIN_TYPE_NORMAL   MarginType = "TRUE"
	MARGIN_TYPE_ISOLATED MarginType = "FALSE"
)

const (
	ORDER_RESP_TYPE_ACK    OrderRespType = "ACK" //MARKET 和 LIMIT 订单类型默认为 FULL, 所有其他订单默认为 ACK.
	ORDER_RESP_TYPE_RESULT OrderRespType = "RESULT"
	ORDER_RESP_TYPE_FULL   OrderRespType = "FULL"
)

const (
	MOVE_TYPE_INVALID                       MoveType = ""                              //失败
	MOVE_TYPE_MAIN_UMFUTURE                 MoveType = "MAIN_UMFUTURE"                 //现货钱包转向U本位合约钱包
	MOVE_TYPE_MAIN_CMFUTURE                 MoveType = "MAIN_CMFUTURE"                 //现货钱包转向币本位合约钱包
	MOVE_TYPE_MAIN_MARGIN                   MoveType = "MAIN_MARGIN"                   //现货钱包转向杠杆全仓钱包
	MOVE_TYPE_UMFUTURE_MAIN                 MoveType = "UMFUTURE_MAIN"                 //U本位合约钱包转向现货钱包
	MOVE_TYPE_UMFUTURE_MARGIN               MoveType = "UMFUTURE_MARGIN"               //U本位合约钱包转向杠杆全仓钱包
	MOVE_TYPE_CMFUTURE_MAIN                 MoveType = "CMFUTURE_MAIN"                 //币本位合约钱包转向现货钱包
	MOVE_TYPE_MARGIN_MAIN                   MoveType = "MARGIN_MAIN"                   //杠杆全仓钱包转向现货钱包
	MOVE_TYPE_MARGIN_UMFUTURE               MoveType = "MARGIN_UMFUTURE"               //杠杆全仓钱包转向U本位合约钱包
	MOVE_TYPE_MARGIN_CMFUTURE               MoveType = "MARGIN_CMFUTURE"               //杠杆全仓钱包转向币本位合约钱包
	MOVE_TYPE_CMFUTURE_MARGIN               MoveType = "CMFUTURE_MARGIN"               //币本位合约钱包转向杠杆全仓钱包
	MOVE_TYPE_ISOLATEDMARGIN_MARGIN         MoveType = "ISOLATEDMARGIN_MARGIN"         //杠杆逐仓钱包转向杠杆全仓钱包
	MOVE_TYPE_MARGIN_ISOLATEDMARGIN         MoveType = "MARGIN_ISOLATEDMARGIN"         //杠杆全仓钱包转向杠杆逐仓钱包
	MOVE_TYPE_ISOLATEDMARGIN_ISOLATEDMARGIN MoveType = "ISOLATEDMARGIN_ISOLATEDMARGIN" //杠杆逐仓钱包转向杠杆逐仓钱包
	MOVE_TYPE_MAIN_FUNDING                  MoveType = "MAIN_FUNDING"                  //现货钱包转向资金钱包
	MOVE_TYPE_FUNDING_MAIN                  MoveType = "FUNDING_MAIN"                  //资金钱包转向现货钱包
	MOVE_TYPE_FUNDING_UMFUTURE              MoveType = "FUNDING_UMFUTURE"              //资金钱包转向U本位合约钱包
	MOVE_TYPE_UMFUTURE_FUNDING              MoveType = "UMFUTURE_FUNDING"              //U本位合约钱包转向资金钱包
	MOVE_TYPE_MARGIN_FUNDING                MoveType = "MARGIN_FUNDING"                //杠杆全仓钱包转向资金钱包
	MOVE_TYPE_FUNDING_MARGIN                MoveType = "FUNDING_MARGIN"                //资金钱包转向杠杆全仓钱包
	MOVE_TYPE_FUNDING_CMFUTURE              MoveType = "FUNDING_CMFUTURE"              //资金钱包转向币本位合约钱包
	MOVE_TYPE_CMFUTURE_FUNDING              MoveType = "CMFUTURE_FUNDING"              //币本位合约钱包转向资金钱包
)

const (
	//充值状态
	//0：等待确认
	//1：确认到账
	//2：充值成功
	//8：因该币种暂停充值而未到账，恢复充值后自动到账
	//12：账户或充值被冻结
	//13：子账户充值拦截
	DOPSIT_TYPE_PENDING   DepositStatus = 0
	DOPSIT_TYPE_CONFIRMED DepositStatus = 1
	DOPSIT_TYPE_SUCCESS   DepositStatus = 2
	DOPSIT_TYPE_PAUSE     DepositStatus = 8
	DOPSIT_TYPE_LOCK      DepositStatus = 12
	DOPSIT_TYPE_INTERCEPT DepositStatus = 13
)
const (
	//0(0:已发送确认Email,1:已被用户取消 2:等待确认 3:被拒绝 4:处理中 5:提现交易失败 6 提现完成)
	TRANSFER_TYPE_CREATED    TransferStatus = 0
	TRANSFER_TYPE_CANCELLED  TransferStatus = 1
	TRANSFER_TYPE_CONFIRMING TransferStatus = 2
	TRANSFER_TYPE_REJECTED   TransferStatus = 3
	TRANSFER_TYPE_PROCESSING TransferStatus = 4
	TRANSFER_TYPE_FAILED     TransferStatus = 5
	TRANSFER_TYPE_SUCCESS    TransferStatus = 6
)

const (
	//PENDING (等待执行), CONFIRMED (成功划转), FAILED (执行失败);
	MOVE_STATUS_PENDING   MoveStatus = "PENDING"
	MOVE_STATUS_CONFIRMED MoveStatus = "CONFIRMED"
	MOVE_STATUS_FAILED    MoveStatus = "FAILED"
)

const (
	//状态: PENDING (等待执行), CONFIRMED (成功借贷), FAILED (执行失败);
	LOAN_STATUS_PENDING   LoanStatus = "PENDING"
	LOAN_STATUS_CONFIRMED LoanStatus = "CONFIRMED"
	LOAN_STATUS_FAILED    LoanStatus = "FAILED"
)

type ContractType string
type ContractStatus string
type PositionSide string
type WorkingType string

const (
	INVALID_CONTRACTTYPE ContractType = ""
	PERPETUAL            ContractType = "PERPETUAL" // 永续合约
	CURRENT_WEEK         ContractType = "CURRENT_WEEK"
	NEXT_WEEK            ContractType = "NEXT_WEEK"
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
	switch type_ {
	case PERPETUAL:
		return fmt.Sprint(symbol, "-", "SWAP")
	case CURRENT_MONTH:
		date := GetDate(transform.THISMONTH)
		return fmt.Sprint(symbol, "-", date)
	case NEXT_MONTH:
		date := GetDate(transform.NEXTMONTH)
		return fmt.Sprint(symbol, "-", date)
	case CURRENT_QUARTER:
		date := GetDate(transform.THISQUARTER)
		return fmt.Sprint(symbol, "-", date)
	case NEXT_QUARTER:
		date := GetDate(transform.NEXTQUARTER)
		return fmt.Sprint(symbol, "-", date)
	case CURRENT_WEEK:
		date := GetDate(transform.THISWEEK)
		return fmt.Sprint(symbol, "-", date)
	case NEXT_WEEK:
		date := GetDate(transform.NEXTWEEK)
		return fmt.Sprint(symbol, "-", date)
	case PERPETUAL_DELIVERING:
		fallthrough
	default:
		return symbol
	}
}

func ParseDate(date string, u bool) common.SymbolType {
	/*
		https://www.binance.com/zh-CN/support/faq/3ae441db4ae740e19af3fe9228eb6619
		季度交割合约是具备固定到期日和交割日的衍生品合约，以每个季度的最后一个周五作爲交割日。例如：“BTCUSDT 当季 0326”代表 2021年03月26日16:00（香港时间）进行交割。
		当季度合约结算交割后，会产生新的季度合约，例如：“BTCUSDT 当季 0326” 于 2021年03月26日16:00（香港时间）交割下架后将会生成新的“BTCUSDT 当季 0625”合约，以此类推；
		交割日结算时将收取交割费，交割费与Taker(吃单)费率相同。
	*/
	switch date {
	case transform.GetDate(transform.THISWEEK):
		if u {
			return common.SymbolType_FUTURE_THIS_WEEK
		} else {
			return common.SymbolType_FUTURE_COIN_NEXT_WEEK
		}
	case transform.GetDate(transform.NEXTWEEK):
		if u {
			return common.SymbolType_FUTURE_NEXT_WEEK
		} else {
			return common.SymbolType_FUTURE_COIN_NEXT_QUARTER
		}
	case transform.GetDate(transform.THISQUARTER):
		if u {
			return common.SymbolType_FUTURE_THIS_QUARTER
		} else {
			return common.SymbolType_FUTURE_COIN_THIS_QUARTER
		}
	case transform.GetDate(transform.NEXTQUARTER):
		if u {
			return common.SymbolType_FUTURE_NEXT_QUARTER
		} else {
			return common.SymbolType_FUTURE_COIN_NEXT_QUARTER
		}
	case "SWAP":
		if u {
			return common.SymbolType_SWAP_FOREVER
		} else {
			return common.SymbolType_SWAP_COIN_FOREVER
		}
	default:
		return common.SymbolType_INVALID_TYPE
	}
}

func GetContractType(symbol string) (sym string, market common.Market, symbolType common.SymbolType) {
	/*
		https://www.binance.com/zh-CN/support/faq/3ae441db4ae740e19af3fe9228eb6619
		季度交割合约是具备固定到期日和交割日的衍生品合约，以每个季度的最后一个周五作爲交割日。例如：“BTCUSDT 当季 0326”代表 2021年03月26日16:00（香港时间）进行交割。
		当季度合约结算交割后，会产生新的季度合约，例如：“BTCUSDT 当季 0326” 于 2021年03月26日16:00（香港时间）交割下架后将会生成新的“BTCUSDT 当季 0625”合约，以此类推；
		交割日结算时将收取交割费，交割费与Taker(吃单)费率相同。
	*/
	symList := strings.Split(symbol, "-")
	sym = fmt.Sprintf("%s/%s", symList[0], symList[1])
	if len(symList) <= 2 {
		return sym, common.Market_SPOT, common.SymbolType_SPOT_NORMAL
	} else {
		if symList[2] == "SWAP" && symList[1] == "USDT" {
			return sym, common.Market_SWAP, common.SymbolType_SWAP_FOREVER
		}
		if symList[2] == "SWAP" && symList[1] == "USD" {
			return sym, common.Market_SWAP_COIN, common.SymbolType_SWAP_COIN_FOREVER
		}

		ss := fmt.Sprintf("%s-%s", symList[0], symList[1])
		switch symbol {
		case GetUBaseSymbol(ss, CURRENT_WEEK):
			if symList[1] == "USDT" {
				return sym, common.Market_FUTURE, common.SymbolType_FUTURE_THIS_WEEK
			} else {
				return sym, common.Market_FUTURE_COIN, common.SymbolType_FUTURE_COIN_THIS_WEEK
			}
		case GetUBaseSymbol(ss, NEXT_WEEK):
			if symList[1] == "USDT" {
				return sym, common.Market_FUTURE, common.SymbolType_FUTURE_NEXT_WEEK
			} else {
				return sym, common.Market_FUTURE_COIN, common.SymbolType_FUTURE_COIN_NEXT_WEEK
			}
		//case GetUBaseSymbol(ss, CURRENT_MONTH):
		//	if symList[1] == "USDT" {
		//		return sym, common.Market_FUTURE, common.SymbolType_FUTURE_THIS_MONTH
		//	} else {
		//		return sym, common.Market_FUTURE_COIN, common.SymbolType_FUTURE_COIN_THIS_MONTH
		//	}
		//case GetUBaseSymbol(ss, NEXT_MONTH):
		//	if symList[1] == "USDT" {
		//		return sym, common.Market_FUTURE, common.SymbolType_FUTURE_NEXT_MONTH
		//	} else {
		//		return sym, common.Market_FUTURE_COIN, common.SymbolType_FUTURE_COIN_NEXT_MONTH
		//	}
		case GetUBaseSymbol(ss, CURRENT_QUARTER):
			if symList[1] == "USDT" {
				return sym, common.Market_FUTURE, common.SymbolType_FUTURE_THIS_QUARTER
			} else {
				return sym, common.Market_FUTURE_COIN, common.SymbolType_FUTURE_COIN_THIS_QUARTER
			}
		case GetUBaseSymbol(ss, NEXT_QUARTER):
			if symList[1] == "USDT" {
				return sym, common.Market_FUTURE, common.SymbolType_FUTURE_NEXT_QUARTER
			} else {
				return sym, common.Market_FUTURE_COIN, common.SymbolType_FUTURE_COIN_NEXT_QUARTER
			}
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
	case common.SymbolType_FUTURE_THIS_WEEK, common.SymbolType_FUTURE_COIN_THIS_WEEK:
		return CURRENT_WEEK
	case common.SymbolType_FUTURE_NEXT_WEEK, common.SymbolType_FUTURE_COIN_NEXT_WEEK:
		return NEXT_WEEK
		// ok没有月的
	//case common.SymbolType_FUTURE_THIS_MONTH:
	//	return CURRENT_MONTH
	//case common.SymbolType_FUTURE_NEXT_MONTH:
	//	return NEXT_MONTH
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
