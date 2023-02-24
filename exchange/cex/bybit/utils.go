package bybit

import "github.com/warmplanet/proto/go/order"

type SideType string
type OrderType string

const (
	SIDE_TYPE_BUY  SideType = "Buy"  // 买入
	SIDE_TYPE_SELL SideType = "Sell" //卖出
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

func GetSideTypeToExchange(side order.TradeSide) SideType {
	switch side {
	case order.TradeSide_BUY:
		return SIDE_TYPE_BUY
	case order.TradeSide_SELL:
		return SIDE_TYPE_SELL
	default:
		return ""
	}
}

func GetOrderTypeToExchange(ot interface{}) OrderType {
	switch ot {
	case order.OrderType_MARKET:
		return ORDER_TYPE_MARKET
	case order.OrderType_LIMIT:
		return ORDER_TYPE_LIMIT
	default:
		return ""
	}
}

func GetOrderStatusFromExchange(s string) order.OrderStatusCode {
	switch s {
	case "NEW":
		return order.OrderStatusCode_OPENED
	case "PARTIALLY_FILLED":
		return order.OrderStatusCode_PARTFILLED
	case "FILLED":
		return order.OrderStatusCode_FILLED
	case "CANCELED":
		return order.OrderStatusCode_CANCELED
	case "REJECTED":
		return order.OrderStatusCode_FAILED
	default:
		return order.OrderStatusCode_FAILED
	}
}
