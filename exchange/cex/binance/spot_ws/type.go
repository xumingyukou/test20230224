package spot_ws

import (
	"clients/exchange/cex/binance/c_api"
	"clients/exchange/cex/binance/spot_api"
	"clients/exchange/cex/binance/u_api"
	"clients/logger"
	"fmt"
	"github.com/warmplanet/proto/go/common"
	"strconv"
	"strings"
)

type req struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int      `json:"id"`
}

type RespTrade struct {
	E  string `json:"e"` // 事件类型
	E1 int64  `json:"E"` // 事件时间
	S  string `json:"s"` // 交易对
	T  int64  `json:"t"` // 交易ID
	P  string `json:"p"` // 成交价格
	Q  string `json:"q"` // 成交数量
	B  int64  `json:"b"` // 买方的订单ID
	A  int64  `json:"a"` // 卖方的订单ID
	T1 int64  `json:"T"` // 成交时间
	M  bool   `json:"m"` // 买方是否是做市方。如true，则此次成交是一个主动卖出单，否则是一个主动买入单。
	M1 bool   `json:"M"` // 请忽略该字段
}

type RespTradeStream struct {
	Stream string    `json:"stream"`
	Data   RespTrade `json:"data"`
}

type RespAggTrade struct {
	E  string `json:"e"` // 事件类型
	E1 int64  `json:"E"` // 事件时间
	S  string `json:"s"` // 交易对
	A  int64  `json:"a"` // 归集交易ID
	P  string `json:"p"` // 成交价格
	Q  string `json:"q"` // 成交数量
	F  int64  `json:"f"` // 被归集的首个交易ID
	L  int64  `json:"l"` // 被归集的末次交易ID
	T  int64  `json:"T"` // 成交时间
	M  bool   `json:"m"` // 买方是否是做市方。如true，则此次成交是一个主动卖出单，否则是一个主动买入单。
	M1 bool   `json:"M"` // 请忽略该字段
}

type RespMarkPrice struct {
	E  string `json:"e"` // 事件类型
	E1 int64  `json:"E"` // 事件时间
	S  string `json:"s"` // 交易对
	P  string `json:"p"` // 标记价格
	I  string `json:"i"` // 现货指数价格
	P1 string `json:"P"` // 预估结算价,仅在结算前最后一小时有参考价值
	R  string `json:"r"` // 资金费率
	T  int64  `json:"T"` // 下次资金时间
}

type RespAggTradeStream struct {
	Stream string       `json:"stream"`
	Data   RespAggTrade `json:"data"`
}

type RespMarkPriceStream struct {
	Stream string        `json:"stream"`
	Data   RespMarkPrice `json:"data"`
}

type RespBookTicker struct {
	U  int64  `json:"u"` // order book updateId
	S  string `json:"s"` // 交易对
	E  string `json:"e"` // 事件类型
	E1 int64  `json:"E"` // 事件推送时间
	T  int64  `json:"T"` // 撮合时间
	B  string `json:"b"` // 买单最优挂单价格
	B1 string `json:"B"` // 买单最优挂单数量
	A  string `json:"a"` // 卖单最优挂单价格
	A1 string `json:"A"` // 卖单最优挂单数量
}

type RespBookTickerStream struct {
	Stream string         `json:"stream"`
	Data   RespBookTicker `json:"data"`
}

type RespLimitDepth struct { //有限档深度信息
	LastUpdateId int64      `json:"lastUpdateId"` // Last update ID
	Bids         [][]string `json:"bids"`         // Bids to be updated:[Price level to be updated, Quantity]
	Asks         [][]string `json:"asks"`         //Asks to be updated[Price level to be updated, Quantity]
}

type RespUBaseLimitDepth struct {
	E  string     `json:"e"` // 事件类型
	E1 int64      `json:"E"` // 事件时间
	T  int64      `json:"T"` // 交易时间
	S  string     `json:"s"`
	U  int64      `json:"U"`
	U1 int64      `json:"u"`
	Pu int64      `json:"pu"`
	Ps string     `json:"ps"` // 标的交易对
	B  [][]string `json:"b"`  // 买方
	A  [][]string `json:"a"`  // 卖方
}

type RespCBaseLimitDepth struct {
	E  string     `json:"e"` // 事件类型
	E1 int64      `json:"E"` // 事件时间
	T  int64      `json:"T"` // 交易时间
	S  string     `json:"s"`
	U  int64      `json:"U"`
	U1 int64      `json:"u"`
	Pu int64      `json:"pu"`
	Ps string     `json:"ps"` // 标的交易对
	B  [][]string `json:"b"`  // 买方
	A  [][]string `json:"a"`  // 卖方
}

type RespLimitDepthStream struct {
	Stream string         `json:"stream"`
	Data   RespLimitDepth `json:"data"`
}

type RespUBaseLimitDepthStream struct {
	Stream string              `json:"stream"`
	Data   RespUBaseLimitDepth `json:"data"`
}

type RespCBaseLimitDepthStream struct {
	Stream string              `json:"stream"`
	Data   RespCBaseLimitDepth `json:"data"`
}

type RespIncrementDepth struct { //增量深度信息
	E  string     `json:"e"`  // 事件类型
	E1 int64      `json:"E"`  // 事件时间
	S  string     `json:"s"`  // 交易对
	U  int64      `json:"U"`  // 从上次推送至今新增的第一个 update Id
	U1 int64      `json:"u"`  // 从上次推送至今新增的最后一个 update Id
	Pu int64      `json:"pu"` // 上次推送的最后一个update Id(即上条消息的‘u’)
	Ps string     `json:"ps"` // 标的交易对
	B  [][]string `json:"b"`  // 变动的买单深度:[变动的价格档位, 数量]
	A  [][]string `json:"a"`  // 变动的卖单深度:[变动的价格档位, 数量]
}
type RespIncrementDepthStream struct {
	Stream string             `json:"stream"`
	Data   RespIncrementDepth `json:"data"`
}
type RespUBaseIncrementDepth struct {
	E  string     `json:"e"`  // 事件类型
	E1 int64      `json:"E"`  // 事件时间
	T  int64      `json:"T"`  // 撮合时间
	S  string     `json:"s"`  // 交易对
	U  int64      `json:"U"`  // 从上次推送至今新增的第一个 update Id
	U1 int64      `json:"u"`  // 从上次推送至今新增的最后一个 update Id
	Pu int64      `json:"pu"` // 上次推送的最后一个update Id(即上条消息的‘u’)
	B  [][]string `json:"b"`  // 变动的买单深度:[变动的价格档位, 数量]
	A  [][]string `json:"a"`  // 变动的卖单深度:[变动的价格档位, 数量]
}

type RespCBaseIncrementDepth struct {
	RespUBaseIncrementDepth
	Ps string `json:"ps"` // 标的交易对
}

type RespUBaseIncrementDepthStream struct {
	Stream string                  `json:"stream"`
	Data   RespUBaseIncrementDepth `json:"data"`
}
type RespCBaseIncrementDepthStream struct {
	Stream string                  `json:"stream"`
	Data   RespCBaseIncrementDepth `json:"data"`
}

type AccountPosition struct {
	A string `json:"a"` // 资产名称
	F string `json:"f"` // 可用余额
	L string `json:"l"` // 冻结余额
}

type RespUserAccount struct {
	E  string             `json:"e"` // 事件类型
	E1 int64              `json:"E"` // 事件时间
	U  int64              `json:"u"` // 账户末次更新时间戳
	B  []*AccountPosition `json:"B"` // 余额
}

type RespUserBalance struct {
	E  string `json:"e"` //Event Type
	E1 int64  `json:"E"` //Event Time
	A  string `json:"a"` //Asset
	D  string `json:"d"` //Balance Delta
	T  int64  `json:"T"` //Clear Time
}

type RespUserOrder struct {
	E  string      `json:"e"` // 事件类型
	E1 int64       `json:"E"` // 事件时间
	S  string      `json:"s"` // 交易对
	C  string      `json:"c"` // clientOrderId
	S1 string      `json:"S"` // 订单方向
	O  string      `json:"o"` // 订单类型
	F  string      `json:"f"` // 有效方式
	Q  string      `json:"q"` // 订单原始数量
	P  string      `json:"p"` // 订单原始价格
	P1 string      `json:"P"` // 止盈止损单触发价格
	D  int64       `json:"d"` // 追踪止损(Trailing Delta) 只有在追踪止损订单中才会推送.
	F1 string      `json:"F"` // 冰山订单数量
	G  int64       `json:"g"` // OCO订单 OrderListId
	C1 string      `json:"C"` // 原始订单自定义ID(原始订单，指撤单操作的对象。撤单本身被视为另一个订单)
	X  string      `json:"x"` // 本次事件的具体执行类型
	X1 string      `json:"X"` // 订单的当前状态
	R  string      `json:"r"` // 订单被拒绝的原因
	I  int64       `json:"i"` // orderId
	L  string      `json:"l"` // 订单末次成交量
	Z  string      `json:"z"` // 订单累计已成交量
	L1 string      `json:"L"` // 订单末次成交价格
	N  string      `json:"n"` // 手续费数量
	N1 interface{} `json:"N"` // 手续费资产类别
	T  int64       `json:"T"` // 成交时间
	T1 int64       `json:"t"` // 成交ID
	I1 int64       `json:"I"` // 请忽略
	W  bool        `json:"w"` // 订单是否在订单簿上？
	M  bool        `json:"m"` // 该成交是作为挂单成交吗？
	M1 bool        `json:"M"` // 请忽略
	O1 int64       `json:"O"` // 订单创建时间
	Z1 string      `json:"Z"` // 订单累计已成交金额
	Y  string      `json:"Y"` // 订单末次成交金额
	Q1 string      `json:"Q"` // Quote Order Qty
}

type BalanceEvent struct {
	A  string `json:"a"`  // 资产名称
	Wb string `json:"wb"` // 钱包余额
	Cw string `json:"cw"` // 除去逐仓仓位保证金的钱包余额
	Bc string `json:"bc"` // 除去盈亏与交易手续费以外的钱包余额改变量
}

type PositionEvent struct {
	S  string `json:"s"`  // 交易对
	Pa string `json:"pa"` // 仓位
	Ep string `json:"ep"` // 入仓价格
	Cr string `json:"cr"` // (费前)累计实现损益
	Up string `json:"up"` // 持仓未实现盈亏
	Mt string `json:"mt"` // 保证金模式
	Iw string `json:"iw"` // 若为逐仓，仓位保证金
	Ps string `json:"ps"` // 持仓方向
}

type AccountEvent struct {
	M string           `json:"m"` // 事件推出原因
	B []*BalanceEvent  `json:"B"` // 余额信息
	P []*PositionEvent `json:"P"` // 仓位信息
}

type RespUBaseAccount struct {
	E  string       `json:"e"` // 事件类型
	E1 int64        `json:"E"` // 事件时间
	T  int64        `json:"T"` // 撮合时间
	A  AccountEvent `json:"a"` // 账户更新事件
}

type UBaseOrderInfo struct {
	S  string `json:"s"`  // 交易对
	C  string `json:"c"`  // 客户端自定订单ID
	S1 string `json:"S"`  // 订单方向
	O  string `json:"o"`  // 订单类型
	F  string `json:"f"`  // 有效方式
	Q  string `json:"q"`  // 订单原始数量
	P  string `json:"p"`  // 订单原始价格
	Ap string `json:"ap"` // 订单平均价格
	Sp string `json:"sp"` // 条件订单触发价格，对追踪止损单无效
	X  string `json:"x"`  // 本次事件的具体执行类型
	X1 string `json:"X"`  // 订单的当前状态
	I  int64  `json:"i"`  // 订单ID
	L  string `json:"l"`  // 订单末次成交量
	Z  string `json:"z"`  // 订单累计已成交量
	L1 string `json:"L"`  // 订单末次成交价格
	N  string `json:"N"`  // 手续费资产类型
	N1 string `json:"n"`  // 手续费数量
	T  int64  `json:"T"`  // 成交时间
	T1 int    `json:"t"`  // 成交ID
	B  string `json:"b"`  // 买单净值
	A  string `json:"a"`  // 卖单净值
	M  bool   `json:"m"`  // 该成交是作为挂单成交吗？
	R  bool   `json:"R"`  // 是否是只减仓单
	Wt string `json:"wt"` // 触发价类型
	Ot string `json:"ot"` // 原始订单类型
	Ps string `json:"ps"` // 持仓方向
	Cp bool   `json:"cp"` // 是否为触发平仓单; 仅在条件订单情况下会推送此字段
	AP string `json:"AP"` // 追踪止损激活价格, 仅在追踪止损单时会推送此字段
	Cr string `json:"cr"` // 追踪止损回调比例, 仅在追踪止损单时会推送此字段
	PP bool   `json:"pP"` // 忽略
	Si int    `json:"si"` // 忽略
	Ss int    `json:"ss"` // 忽略
	Rp string `json:"rp"` // 该交易实现盈亏
}

type RespUBaseOrder struct {
	E  string         `json:"e"` // 事件类型
	E1 int64          `json:"E"` // 事件时间
	T  int64          `json:"T"` // 撮合时间
	O  UBaseOrderInfo `json:"o"`
}

func GetSymbolMarket(symbol string, marketType int) (string, common.Market, common.SymbolType) {
	switch marketType {
	case 0:
		res := spot_api.ParseSymbolName(symbol)
		return res, common.Market_SPOT, common.SymbolType_SPOT_NORMAL
	case 1:
		sym, market, type_ := u_api.GetContractType(symbol)
		return sym, market, type_
	case 2:
		sym, market, type_ := c_api.GetContractType(symbol)
		return sym, market, type_
	default:
		return symbol, common.Market_INVALID_MARKET, common.SymbolType_INVALID_TYPE
	}
}

func GetSymbolKey(symbol string, marketType int) (string, string) {
	switch marketType {
	case 0:
		res := spot_api.ParseSymbolName(symbol)
		return res, res
	case 1:
		sym, market, type_ := u_api.GetContractType(symbol)
		return sym, SymbolKeyGen(sym, market, type_)
	case 2:
		sym, market, type_ := c_api.GetContractType(symbol)
		return sym, SymbolKeyGen(sym, market, type_)
	default:
		return symbol, ""
	}
}

func SymbolKeyGen(symbol string, market common.Market, type_ common.SymbolType) string {
	if !spot_api.IsUBaseSymbolType(type_) && !spot_api.IsCBaseSymbolType(type_) {
		return symbol
	}
	return fmt.Sprint(symbol, "_", market.Number(), "_", type_.Number())
}

func GetSymbolFromKey(symbolKey string) (symbol string, type_ common.SymbolType) {
	symStr := strings.Split(symbolKey, "_")
	symbol = symStr[0]
	if len(symStr) > 1 {
		typeId, _ := strconv.Atoi(symStr[len(symStr)-1])
		type_ = common.SymbolType(typeId)
	} else {
		type_ = common.SymbolType_INVALID_TYPE
	}
	return
}

func getBaseQty(symbol string, qty float64, price float64, market common.Market) float64 {
	if market == common.Market_FUTURE_COIN || market == common.Market_SWAP_COIN {
		return getCbaseQty(symbol, qty, price)
	} else {
		return qty
	}
}

func getCbaseQty(symbol string, qty float64, price float64) float64 {
	if len(strings.Split(symbol, "/")) == 1 {
		logger.Logger.Error("get cbase qty error, symbol is invalid: ", symbol)
		return qty
	}
	if strings.Split(symbol, "/")[0] == "BTC" {
		return qty * 100 / price
	} else {
		return qty * 10 / price
	}
}
