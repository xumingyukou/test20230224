package c_api

import (
	"clients/exchange/cex/binance/u_api"
)

type RespTickerPriceItem struct {
	u_api.RespTickerPriceItem
	Ps string `json:"ps"` // 标的交易对
}

type RespTickerPriceList []*RespTickerPriceItem

type RespPositionSideDual struct {
	DualSidePosition bool `json:"dualSidePosition"` // "true": 双向持仓模式；"false": 单向持仓模式
}

type BalanceItem struct {
	u_api.BalanceItem
	WithdrawAvailable string `json:"withdrawAvailable"` // 最大可提款金额,同`GET /dapi/account`中"maxWithdrawAmount"
}

type RespBalance []*BalanceItem

type RespLeverage struct {
	Leverage int    `json:"leverage"` // 杠杆倍数
	Symbol   string `json:"symbol"`   // 交易对
	MaxQty   string `json:"maxQty"`   // 当前杠杆倍数下允许的最大base asset数量
}
