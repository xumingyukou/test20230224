package spot_api

type MoveType string
type SideType string
type TimeInForceType string //有效方式
type OrderType string       //订单类型 (orderTypes, type):
type OrderRespType string   //响应
type MarginType string
type OrderStatus string
type DepositStatus int
type TransferStatus int
type MoveStatus string
type LoanStatus string

const (
	ORDER_TYPE_LIMIT                OrderType = "LIMIT"                //限价单
	ORDER_TYPE_MARKET               OrderType = "MARKET"               //市价单
	ORDER_TYPE_STOP_LOSS            OrderType = "STOP_LOSS"            //市价止损单
	ORDER_TYPE_STOP_LOSS_LIMIT      OrderType = "STOP_LOSS_LIMIT"      //限价止损单
	ORDER_TYPE_TAKE_PROFIT          OrderType = "TAKE_PROFIT"          //市价止盈单
	ORDER_TYPE_TAKE_PROFIT_LIMIT    OrderType = "TAKE_PROFIT_LIMIT"    //限价止盈单
	ORDER_TYPE_LIMIT_MAKER          OrderType = "LIMIT_MAKER"          //限价只挂单
	ORDER_TYPE_TRAILING_STOP_MARKET OrderType = "TRAILING_STOP_MARKET" // 跟踪止损单
)

const (
	TIME_IN_FORCE_GTC TimeInForceType = "GTC" //成交为止 订单会一直有效，直到被成交或者取消。
	TIME_IN_FORCE_IOC TimeInForceType = "IOC" //无法立即成交的部分就撤销 订单在失效前会尽量多的成交。
	TIME_IN_FORCE_FOK TimeInForceType = "FOK" //无法全部立即成交就撤销 如果无法全部成交，订单会失效。
	TIME_IN_FORCE_GTX TimeInForceType = "GTX" //- Good Till Crossing 无法成为挂单方就撤销

)

const (
	SIDE_TYPE_BUY  SideType = "BUY"  // 买入
	SIDE_TYPE_SELL SideType = "SELL" //卖出
)

const (
	ORDER_STATUS_NEW              OrderStatus = "NEW"
	ORDER_STATUS_PARTIALLY_FILLED OrderStatus = "PARTIALLY_FILLED"
	ORDER_STATUS_FILLED           OrderStatus = "FILLED"
	ORDER_STATUS_CANCELED         OrderStatus = "CANCELED"
	ORDER_STATUS_PENDING_CANCEL   OrderStatus = "PENDING_CANCEL"
	ORDER_STATUS_REJECTED         OrderStatus = "REJECTED"
	ORDER_STATUS_EXPIRED          OrderStatus = "EXPIRED"
	ORDER_STATUS_LIVE             OrderStatus = "LIVE"
	ORDER_STATUS_NEW_INSURANCE    OrderStatus = "NEW_INSURANCE" // 风险保障基金(强平)
	ORDER_STATUS_NEW_ADL          OrderStatus = "NEW_ADL"       // 自动减仓序列(强平)
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
	//0(0:pending,6: credited but cannot withdraw, 1:success)

	DOPSIT_TYPE_PENDING   DepositStatus = 0
	DOPSIT_TYPE_CONFIRMED DepositStatus = 6
	DOPSIT_TYPE_SUCCESS   DepositStatus = 1
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
