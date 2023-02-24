package u_api

import "clients/exchange/cex/binance/spot_api"

type RespPremiumIndexItem struct {
	Symbol               string `json:"symbol"`               // 交易对
	MarkPrice            string `json:"markPrice"`            // 标记价格
	IndexPrice           string `json:"indexPrice"`           // 指数价格
	EstimatedSettlePrice string `json:"estimatedSettlePrice"` // 预估结算价,仅在交割开始前最后一小时有意义
	LastFundingRate      string `json:"lastFundingRate"`      // 最近更新的资金费率
	NextFundingTime      int64  `json:"nextFundingTime"`      // 下次资金费时间
	InterestRate         string `json:"interestRate"`         // 标的资产基础利率
	Time                 int64  `json:"time"`                 // 更新时间
	Pair                 string `json:"pair"`                 // 基础标的
}

type RespPremiumIndexList []*RespPremiumIndexItem

type RespTickerPriceItem struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Time   int64  `json:"time"`
}

type RespTickerPriceList []*RespTickerPriceItem

type RespUBaseOrderResult struct {
	spot_api.RespCancelOrder

	CumQty        string `json:"cumQty"`
	CumQuote      string `json:"cumQuote"`      // 成交金额
	AvgPrice      string `json:"avgPrice"`      // 平均成交价
	ReduceOnly    bool   `json:"reduceOnly"`    // 仅减仓
	PositionSide  string `json:"positionSide"`  // 持仓方向
	StopPrice     string `json:"stopPrice"`     // 触发价，对`TRAILING_STOP_MARKET`无效
	ClosePosition bool   `json:"closePosition"` // 是否条件全平仓
	OrigType      string `json:"origType"`      // 触发前订单类型
	ActivatePrice string `json:"activatePrice"` // 跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	PriceRate     string `json:"priceRate"`     // 跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	UpdateTime    int64  `json:"updateTime"`    // 更新时间
	WorkingType   string `json:"workingType"`   // 条件价格触发类型
	PriceProtect  bool   `json:"priceProtect"`  // 是否开启条件单触发保护
	Time          int64  `json:"time"`          // 订单时间:查询时候返回
	CumBase       string `json:"cumBase"`       // 成交额(标的数量)
	Pair          string `json:"pair"`          // 标的交易对
	spot_api.RespError
}

type RespPositionSideDual struct {
	DualSidePosition bool `json:"dualSidePosition"` // "true": 双向持仓模式；"false": 单向持仓模式
}

type BalanceItem struct {
	AccountAlias       string `json:"accountAlias"`       // 账户唯一识别码
	Asset              string `json:"asset"`              // 资产
	Balance            string `json:"balance"`            // 总余额
	CrossWalletBalance string `json:"crossWalletBalance"` // 全仓余额
	CrossUnPnl         string `json:"crossUnPnl"`         // 全仓持仓未实现盈亏
	AvailableBalance   string `json:"availableBalance"`   // 下单可用余额
	MaxWithdrawAmount  string `json:"maxWithdrawAmount"`  // 最大可转出余额
	MarginAvailable    bool   `json:"marginAvailable"`    // 是否可用作联合保证金
	UpdateTime         int64  `json:"updateTime"`
}

type RespBalance []*BalanceItem

type AssertItem struct {
	Asset                  string `json:"asset"`                  //资产
	WalletBalance          string `json:"walletBalance"`          //余额
	UnrealizedProfit       string `json:"unrealizedProfit"`       // 未实现盈亏
	MarginBalance          string `json:"marginBalance"`          // 保证金余额
	MaintMargin            string `json:"maintMargin"`            // 维持保证金
	InitialMargin          string `json:"initialMargin"`          // 当前所需起始保证金
	PositionInitialMargin  string `json:"positionInitialMargin"`  // 持仓所需起始保证金(基于最新标记价格)
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"` // 当前挂单所需起始保证金(基于最新标记价格)
	CrossWalletBalance     string `json:"crossWalletBalance"`     //全仓账户余额
	CrossUnPnl             string `json:"crossUnPnl"`             // 全仓持仓未实现盈亏
	AvailableBalance       string `json:"availableBalance"`       // 可用余额
	MaxWithdrawAmount      string `json:"maxWithdrawAmount"`      // 最大可转出余额
	MarginAvailable        bool   `json:"marginAvailable"`        // 是否可用作联合保证金
	UpdateTime             int64  `json:"updateTime"`             //更新时间
}
type PositionItem struct {
	Symbol                 string `json:"symbol"`                 // 交易对
	InitialMargin          string `json:"initialMargin"`          // 当前所需起始保证金(基于最新标记价格)
	MaintMargin            string `json:"maintMargin"`            //维持保证金
	UnrealizedProfit       string `json:"unrealizedProfit"`       // 持仓未实现盈亏
	PositionInitialMargin  string `json:"positionInitialMargin"`  // 持仓所需起始保证金(基于最新标记价格)
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"` // 当前挂单所需起始保证金(基于最新标记价格)
	Leverage               string `json:"leverage"`               // 杠杆倍率
	Isolated               bool   `json:"isolated"`               // 是否是逐仓模式
	EntryPrice             string `json:"entryPrice"`             // 持仓成本价
	MaxNotional            string `json:"maxNotional"`            // 当前杠杆下用户可用的最大名义价值
	Notional               string `json:"notional"`               // 名义价值
	BidNotional            string `json:"bidNotional"`            // 买单净值，忽略
	AskNotional            string `json:"askNotional"`            // 卖单净值，忽略
	PositionSide           string `json:"positionSide"`           // 持仓方向
	PositionAmt            string `json:"positionAmt"`            // 持仓数量
	UpdateTime             int    `json:"updateTime"`             // 更新时间
	MaxQty                 string `json:"maxQty"`                 // 当前杠杆下最大可开仓数(标的数量)
}

type RespAccount struct {
	FeeTier                     int             `json:"feeTier"`                     // 手续费等级
	CanTrade                    bool            `json:"canTrade"`                    // 是否可以交易
	CanDeposit                  bool            `json:"canDeposit"`                  // 是否可以入金
	CanWithdraw                 bool            `json:"canWithdraw"`                 // 是否可以出金
	UpdateTime                  int64           `json:"updateTime"`                  // 保留字段，请忽略
	TotalInitialMargin          string          `json:"totalInitialMargin"`          // 当前所需起始保证金总额(存在逐仓请忽略), 仅计算usdt资产
	TotalMaintMargin            string          `json:"totalMaintMargin"`            // 维持保证金总额, 仅计算usdt资产
	TotalWalletBalance          string          `json:"totalWalletBalance"`          // 账户总余额, 仅计算usdt资产
	TotalUnrealizedProfit       string          `json:"totalUnrealizedProfit"`       // 持仓未实现盈亏总额, 仅计算usdt资产
	TotalMarginBalance          string          `json:"totalMarginBalance"`          // 保证金总余额, 仅计算usdt资产
	TotalPositionInitialMargin  string          `json:"totalPositionInitialMargin"`  // 持仓所需起始保证金(基于最新标记价格), 仅计算usdt资产
	TotalOpenOrderInitialMargin string          `json:"totalOpenOrderInitialMargin"` // 当前挂单所需起始保证金(基于最新标记价格), 仅计算usdt资产
	TotalCrossWalletBalance     string          `json:"totalCrossWalletBalance"`     // 全仓账户余额, 仅计算usdt资产
	TotalCrossUnPnl             string          `json:"totalCrossUnPnl"`             // 全仓持仓未实现盈亏总额, 仅计算usdt资产
	AvailableBalance            string          `json:"availableBalance"`            // 可用余额, 仅计算usdt资产
	MaxWithdrawAmount           string          `json:"maxWithdrawAmount"`           // 最大可转出余额, 仅计算usdt资产
	Assets                      []*AssertItem   `json:"assets"`                      // 资产
	Positions                   []*PositionItem `json:"positions"`                   // 头寸，将返回所有市场symbol。根据用户持仓模式展示持仓方向，即单向模式下只返回BOTH持仓情况，双向模式下只返回 LONG 和 SHORT 持仓情况
}

type RespLeverage struct {
	Leverage         int    `json:"leverage"`         // 杠杆倍数
	MaxNotionalValue string `json:"maxNotionalValue"` // 当前杠杆倍数下允许的最大名义价值
	Symbol           string `json:"symbol"`           // 交易对
}

type RespCommissionRate struct {
	Symbol              string `json:"symbol"`
	MakerCommissionRate string `json:"makerCommissionRate"`
	TakerCommissionRate string `json:"takerCommissionRate"`
}
