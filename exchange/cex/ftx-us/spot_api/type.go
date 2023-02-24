package spot_api

import (
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/order"
)

type RespError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type MarketItem struct {
	Name                  string  `json:"name"`                  //e.g. "BTC/USD" for spot, "BTC-PERP" for futures
	BaseCurrency          string  `json:"baseCurrency"`          //spot markets only
	QuoteCurrency         string  `json:"quoteCurrency"`         //spot markets only
	QuoteVolume24h        float64 `json:"quoteVolume24h"`        //
	Change1h              float64 `json:"change1h"`              //change in the past hour
	Change24h             float64 `json:"change24h"`             //change in the past 24 hours
	ChangeBod             float64 `json:"changeBod"`             //change since start of day (00:00 UTC)
	HighLeverageFeeExempt bool    `json:"highLeverageFeeExempt"` //false
	MinProvideSize        float64 `json:"minProvideSize"`        //Minimum maker order size (if >10 orders per hour fall below this size)
	Type                  string  `json:"type"`                  //"future" or "spot"
	FutureType            string  `json:"futureType"`            // "future" or "perpetual"
	Underlying            string  `json:"underlying"`            //future markets only
	Enabled               bool    `json:"enabled"`               //
	Ask                   float64 `json:"ask"`                   //best ask
	Bid                   float64 `json:"bid"`                   //best bid
	Last                  float64 `json:"last"`                  //last traded price
	PostOnly              bool    `json:"postOnly"`              //if the market is in post-only mode (all orders get modified to be post-only, in addition to other settings they may have)
	Price                 float64 `json:"price"`                 //current price
	PriceIncrement        float64 `json:"priceIncrement"`        //
	SizeIncrement         float64 `json:"sizeIncrement"`         //
	Restricted            bool    `json:"restricted"`            //if the market has nonstandard restrictions on which jurisdictions can trade it
	VolumeUsd24h          float64 `json:"volumeUsd24h"`          //USD volume in past 24 hours
	LargeOrderThreshold   float64 `json:"largeOrderThreshold"`   //threshold above which an order is considered large (for VIP rate limits)
	IsEtfMarket           bool    `json:"isEtfMarket"`           //if the market has an ETF as its baseCurrency
}

type RespGetMarkets struct {
	Success bool          `json:"success"`
	Result  []*MarketItem `json:"result"`
}

type PositionItem struct {
	Cost                         float64 `json:"cost"`
	EntryPrice                   float64 `json:"entryPrice"`
	Future                       string  `json:"future"`
	InitialMarginRequirement     float64 `json:"initialMarginRequirement"`
	LongOrderSize                float64 `json:"longOrderSize"`
	MaintenanceMarginRequirement float64 `json:"maintenanceMarginRequirement"`
	NetSize                      float64 `json:"netSize"`
	OpenSize                     float64 `json:"openSize"`
	RealizedPnl                  float64 `json:"realizedPnl"`
	ShortOrderSize               float64 `json:"shortOrderSize"`
	Side                         string  `json:"side"`
	Size                         float64 `json:"size"`
	UnrealizedPnl                float64 `json:"unrealizedPnl"`
}

type AccountInfo struct {
	BackstopProvider             bool            `json:"backstopProvider"`
	Collateral                   float64         `json:"collateral"`
	FreeCollateral               float64         `json:"freeCollateral"`
	InitialMarginRequirement     float64         `json:"initialMarginRequirement"`
	Leverage                     float64         `json:"leverage"`
	Liquidating                  bool            `json:"liquidating"`
	MaintenanceMarginRequirement float64         `json:"maintenanceMarginRequirement"`
	MakerFee                     float64         `json:"makerFee"`
	MarginFraction               float64         `json:"marginFraction"`
	OpenMarginFraction           float64         `json:"openMarginFraction"`
	TakerFee                     float64         `json:"takerFee"`
	TotalAccountValue            float64         `json:"totalAccountValue"`
	TotalPositionSize            float64         `json:"totalPositionSize"`
	Username                     string          `json:"username"`
	Positions                    []*PositionItem `json:"positions"`
}

type RespGetAccountInfo struct {
	Success bool        `json:"success"`
	Result  AccountInfo `json:"result"`
}

type RespGetOrderbook struct {
	Success bool `json:"success"`
	Result  struct {
		Asks [][]float64 `json:"asks"`
		Bids [][]float64 `json:"bids"`
	} `json:"result"`
}

type TradeInfo struct {
	Id          int       `json:"id"`
	Liquidation bool      `json:"liquidation"`
	Price       float64   `json:"price"`
	Side        string    `json:"side"`
	Size        float64   `json:"size"`
	Time        time.Time `json:"time"`
}

type RespGetTrades struct {
	Success bool         `json:"success"`
	Result  []*TradeInfo `json:"result"`
}

type HistPriceItem struct {
	Close     float64   `json:"close"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Open      float64   `json:"open"`
	StartTime time.Time `json:"startTime"`
	Volume    float64   `json:"volume"`
}

type RespGetHistPrices struct {
	Success bool             `json:"success"`
	Result  []*HistPriceItem `json:"result"`
}

type HistBalanceItem struct {
	Account string  `json:"account"`
	Ticker  string  `json:"ticker"`
	Size    float64 `json:"size"`
	Price   float64 `json:"price"`
}

type HistBalanceInfo struct {
	Id       int                `json:"id"`
	Accounts []string           `json:"accounts"`
	Time     int                `json:"time"`
	EndTime  int                `json:"endTime"`
	Status   string             `json:"status"`
	Error    bool               `json:"error"`
	Results  []*HistBalanceItem `json:"results"`
}

type RespGetHistBalance struct {
	Success bool            `json:"success"`
	Result  HistBalanceInfo `json:"result"`
}

func GetOrderTypeToExchange(ot order.OrderType) OrderType {
	switch ot {
	case order.OrderType_MARKET:
		return ORDER_TYPE_MARKET
	case order.OrderType_LIMIT:
		return ORDER_TYPE_LIMIT
	default:
		return ""
	}
}

type PositionInfo2 struct {
	Cost                         float64 `json:"cost"`
	CumulativeBuySize            float64 `json:"cumulativeBuySize"`
	CumulativeSellSize           float64 `json:"cumulativeSellSize"`
	EntryPrice                   float64 `json:"entryPrice"`
	EstimatedLiquidationPrice    float64 `json:"estimatedLiquidationPrice"`
	Future                       string  `json:"future"`
	InitialMarginRequirement     float64 `json:"initialMarginRequirement"`
	LongOrderSize                float64 `json:"longOrderSize"`
	MaintenanceMarginRequirement float64 `json:"maintenanceMarginRequirement"`
	NetSize                      float64 `json:"netSize"`
	OpenSize                     float64 `json:"openSize"`
	RealizedPnl                  float64 `json:"realizedPnl"`
	RecentAverageOpenPrice       float64 `json:"recentAverageOpenPrice"`
	RecentBreakEvenPrice         float64 `json:"recentBreakEvenPrice"`
	RecentPnl                    float64 `json:"recentPnl"`
	ShortOrderSize               float64 `json:"shortOrderSize"`
	Side                         string  `json:"side"`
	Size                         float64 `json:"size"`
	UnrealizedPnl                int     `json:"unrealizedPnl"`
	CollateralUsed               float64 `json:"collateralUsed"`
}

type RespGetPositions struct {
	Success bool             `json:"success"`
	Result  []*PositionInfo2 `json:"result"`
}

type RespGetOrderHistory struct {
	Success bool `json:"success"`
	Result  []*struct {
		AvgFillPrice  float64   `json:"avgFillPrice"`
		ClientId      string    `json:"clientId"`
		CreatedAt     time.Time `json:"createdAt"`
		FilledSize    float64   `json:"filledSize"`
		Future        string    `json:"future"`
		Id            int       `json:"id"`
		Ioc           bool      `json:"ioc"`
		Market        string    `json:"market"`
		PostOnly      bool      `json:"postOnly"`
		Price         float64   `json:"price"`
		ReduceOnly    bool      `json:"reduceOnly"`
		RemainingSize float64   `json:"remainingSize"`
		Side          string    `json:"side"`
		Size          float64   `json:"size"`
		Status        string    `json:"status"`
		Type          string    `json:"type"`
	} `json:"result"`
	HasMoreData bool `json:"hasMoreData"`
}

type RespGetBalances struct {
	Success bool `json:"success"`
	Result  []struct {
		Coin                   string  `json:"coin"`
		Free                   float64 `json:"free"`
		SpotBorrow             float64 `json:"spotBorrow"`
		Total                  float64 `json:"total"`
		UsdValue               float64 `json:"usdValue"`
		AvailableWithoutBorrow float64 `json:"availableWithoutBorrow"`
	} `json:"result"`
}

type RespPlaceOrder struct {
	Success bool `json:"success"`
	Result  struct {
		CreatedAt     time.Time `json:"createdAt"`
		FilledSize    float64   `json:"filledSize"`
		Future        string    `json:"future"`
		Id            int       `json:"id"`
		Market        string    `json:"market"`
		Price         float64   `json:"price"`
		RemainingSize float64   `json:"remainingSize"`
		Side          string    `json:"side"`
		Size          float64   `json:"size"`
		Status        string    `json:"status"`
		Type          string    `json:"type"`
		ReduceOnly    bool      `json:"reduceOnly"`
		Ioc           bool      `json:"ioc"`
		PostOnly      bool      `json:"postOnly"`
		ClientId      string    `json:"clientId"`
	} `json:"result"`
}

type Precision struct {
	MinAmount       float64 `json:"min_amount"`
	AmountPrecision int     `json:"amount_precision"`
	MinPrice        float64 `json:"min_price"`
	MinValue        float64 `json:"min_value"`
	PricePrecision  int     `json:"price_precision"`
}

type FuturePositionRiskVosItem struct {
	EntryPrice       string `json:"entryPrice"`
	Leverage         string `json:"leverage"`
	MaxNotional      string `json:"maxNotional"`
	LiquidationPrice string `json:"liquidationPrice"`
	MarkPrice        string `json:"markPrice"`
	PositionAmount   string `json:"positionAmount"`
	Symbol           string `json:"symbol"`
	UnrealizedProfit string `json:"unrealizedProfit"`
}

type RespRequestWithdrawal struct {
	Success bool `json:"success"`
	Result  struct {
		Coin    string    `json:"coin"`
		Address string    `json:"address"`
		Tag     string    `json:"tag"`
		Fee     float64   `json:"fee"`
		Id      int       `json:"id"`
		Size    float64   `json:"size"`
		Status  string    `json:"status"`
		Time    time.Time `json:"time"`
		Txid    string    `json:"txid"`
	} `json:"result"`
}

type RespGetWithdrawalFee struct {
	Success bool `json:"success"`
	Result  struct {
		Method    string  `json:"method"`
		Address   string  `json:"address"`
		Fee       float64 `json:"fee"`
		Congested bool    `json:"congested"`
	} `json:"result"`
}

type RespGetOpenOrders struct {
	Success bool `json:"success"`
	Result  []struct {
		CreatedAt     time.Time `json:"createdAt"`
		FilledSize    float64   `json:"filledSize"`
		Future        string    `json:"future"`
		Id            int       `json:"id"`
		Market        string    `json:"market"`
		Price         float64   `json:"price"`
		AvgFillPrice  float64   `json:"avgFillPrice"`
		RemainingSize float64   `json:"remainingSize"`
		Side          string    `json:"side"`
		Size          float64   `json:"size"`
		Status        string    `json:"status"`
		Type          string    `json:"type"`
		ReduceOnly    bool      `json:"reduceOnly"`
		Ioc           bool      `json:"ioc"`
		PostOnly      bool      `json:"postOnly"`
		ClientId      string    `json:"clientId"`
	} `json:"result"`
}

type RespGetOrderStatus struct {
	Success bool `json:"success"`
	Result  struct {
		CreatedAt     time.Time `json:"createdAt"`
		FilledSize    float64   `json:"filledSize"`
		Future        string    `json:"future"`
		Id            int       `json:"id"`
		Market        string    `json:"market"`
		Price         float64   `json:"price"`
		AvgFillPrice  float64   `json:"avgFillPrice"`
		RemainingSize float64   `json:"remainingSize"`
		Side          string    `json:"side"`
		Size          float64   `json:"size"`
		Status        string    `json:"status"`
		Type          string    `json:"type"`
		ReduceOnly    bool      `json:"reduceOnly"`
		Ioc           bool      `json:"ioc"`
		PostOnly      bool      `json:"postOnly"`
		ClientId      string    `json:"clientId"`
		Liquidation   bool      `json:"liquidation"`
	} `json:"result"`
}

type RespFuturesPositionRisk struct {
	FuturePositionRiskVos []*FuturePositionRiskVosItem `json:"futurePositionRiskVos"`
}

type RespFills struct {
	Success bool `json:"success"`
	Result  []struct {
		Fee           float64   `json:"fee"`
		FeeCurrency   string    `json:"feeCurrency"`
		FeeRate       float64   `json:"feeRate"`
		Future        string    `json:"future"`
		Id            int       `json:"id"`
		Liquidity     string    `json:"liquidity"`
		Market        string    `json:"market"`
		BaseCurrency  string    `json:"baseCurrency"`
		QuoteCurrency string    `json:"quoteCurrency"`
		OrderId       int       `json:"orderId"`
		TradeId       int       `json:"tradeId"`
		Price         float64   `json:"price"`
		Side          string    `json:"side"`
		Size          float64   `json:"size"`
		Time          time.Time `json:"time"`
		Type          string    `json:"type"`
	} `json:"result"`
}

type RespGetWithdrawalHistory struct {
	Success bool `json:"success"`
	Result  []struct {
		Coin    string    `json:"coin"`
		Address string    `json:"address"`
		Tag     string    `json:"tag"`
		Fee     float64   `json:"fee"`
		Id      int64     `json:"id"`
		Size    float64   `json:"size"`
		Status  string    `json:"status"`
		Time    time.Time `json:"time"`
		Method  string    `json:"method"`
		Txid    string    `json:"txid"`
	} `json:"result"`
}

type RespGetDepositHistory struct {
	Success bool `json:"success"`
	Result  []struct {
		Coin          string    `json:"coin"`
		Confirmations int       `json:"confirmations"`
		ConfirmedTime time.Time `json:"confirmedTime"`
		Fee           float64   `json:"fee"`
		Id            int       `json:"id"`
		SentTime      time.Time `json:"sentTime"`
		Size          float64   `json:"size"`
		Status        string    `json:"status"`
		Time          time.Time `json:"time"`
		Txid          string    `json:"txid"`
	} `json:"result"`
}

type RespGetLendingHistory struct {
	Success bool `json:"success"`
	Result  []struct {
		Coin string    `json:"coin"`
		Time time.Time `json:"time"`
		Rate float64   `json:"rate"`
		Size float64   `json:"size"`
	} `json:"result"`
}

type RespGetBorrowRates struct {
	Success bool `json:"success"`
	Result  []struct {
		Coin     string  `json:"coin"`
		Estimate float64 `json:"estimate"`
		Previous float64 `json:"previous"`
	} `json:"result"`
}

type RespGetLendingRates struct {
	Success bool `json:"success"`
	Result  []struct {
		Coin     string  `json:"coin"`
		Estimate float64 `json:"estimate"`
		Previous float64 `json:"previous"`
	} `json:"result"`
}

type RespGetMarketInfo struct {
	Success bool `json:"success"`
	Result  []struct {
		Coin          string  `json:"coin"`
		Borrowed      float64 `json:"borrowed"`
		Free          float64 `json:"free"`
		EstimatedRate float64 `json:"estimatedRate"`
		PreviousRate  float64 `json:"previousRate"`
	} `json:"result"`
}

type RespGetLendingOffers struct {
	Success bool `json:"success"`
	Result  []struct {
		Coin string  `json:"coin"`
		Rate float64 `json:"rate"`
		Size float64 `json:"size"`
	} `json:"result"`
}

type RespGetDepositAddress struct {
	Success bool `json:"success"`
	Result  struct {
		Address string `json:"address"`
		Tag     string `json:"tag"`
	} `json:"result"`
}

type RespSubmitLendingOffer struct {
	Success bool   `json:"success"`
	Result  string `json:"result"`
}

type RespRateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

type AssetItem struct {
	Asset             string `json:"asset"`
	MarginAvailable   bool   `json:"marginAvailable"`   // 是否可用作保证金
	AutoAssetExchange *int   `json:"autoAssetExchange"` // 保证金资产自动兑换阈值
}

type HistoricalTradeInfo struct {
	Id           int    `json:"id"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
	IsBestMatch  bool   `json:"isBestMatch"`
}

type RespHistoricalTrade []*HistoricalTradeInfo

type Filter struct {
	FilterType          string  `json:"filterType"`
	MaxPrice            float64 `json:"maxPrice,string"`
	MinPrice            float64 `json:"minPrice,string"`
	TickSize            float64 `json:"tickSize,string"`
	MultiplierUp        float64 `json:"multiplierUp,string"`
	MultiplierDown      float64 `json:"multiplierDown,string"`
	AvgPriceMins        int     `json:"avgPriceMins"`
	MinQty              float64 `json:"minQty,string"`
	MaxQty              float64 `json:"maxQty,string"`
	StepSize            float64 `json:"stepSize,string"`
	MinNotional         float64 `json:"minNotional,string"`
	ApplyToMarket       bool    `json:"applyToMarket"`
	Limit               int     `json:"limit"`
	MaxNumAlgoOrders    int     `json:"maxNumAlgoOrders"`
	MaxNumIcebergOrders int     `json:"maxNumIcebergOrders"`
	MaxNumOrders        int     `json:"maxNumOrders"`

	Notional          string `json:"notional,omitempty"`
	MultiplierDecimal int    `json:"multiplierDecimal,omitempty"`
}

type NetWork struct {
	AddressRegex            string `json:"addressRegex"`
	Coin                    string `json:"coin"`
	DepositDesc             string `json:"depositDesc,omitempty"`
	DepositEnable           bool   `json:"depositEnable"`
	IsDefault               bool   `json:"isDefault"`
	MemoRegex               string `json:"memoRegex"`
	MinConfirm              int    `json:"minConfirm"`
	Name                    string `json:"name"`
	Network                 string `json:"network"`
	ResetAddressStatus      bool   `json:"resetAddressStatus"`
	SpecialTips             string `json:"specialTips"`
	UnLockConfirm           int    `json:"unLockConfirm"`
	WithdrawDesc            string `json:"withdrawDesc,omitempty"`
	WithdrawEnable          bool   `json:"withdrawEnable"`
	WithdrawFee             string `json:"withdrawFee"`
	WithdrawIntegerMultiple string `json:"withdrawIntegerMultiple"`
	WithdrawMax             string `json:"withdrawMax"`
	WithdrawMin             string `json:"withdrawMin"`
	SameAddress             bool   `json:"sameAddress"`
}

type CapitalConfigGetAllItem struct {
	Coin              string     `json:"coin"`
	DepositAllEnable  bool       `json:"depositAllEnable"`
	Free              string     `json:"free"`
	Freeze            string     `json:"freeze"`
	Ipoable           string     `json:"ipoable"`
	Ipoing            string     `json:"ipoing"`
	IsLegalMoney      bool       `json:"isLegalMoney"`
	Locked            string     `json:"locked"`
	Name              string     `json:"name"`
	NetworkList       []*NetWork `json:"networkList"`
	Storage           string     `json:"storage"`
	Trading           bool       `json:"trading"`
	WithdrawAllEnable bool       `json:"withdrawAllEnable"`
	Withdrawing       string     `json:"withdrawing"`
}

type RespCapitalConfigGetAll []*CapitalConfigGetAllItem

type TradeFeeItem struct {
	Symbol          string `json:"symbol"`
	MakerCommission string `json:"makerCommission"`
	TakerCommission string `json:"takerCommission"`
}

type RespAssetTradeFee []*TradeFeeItem

type RespCapitalWithdrawApply struct {
	Id string `json:"id"`
}

type RespAssetTransfer struct {
	TranId int64 `json:"tranId"`
}

type BalanceItem struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

type MarginIsolatedAssetItem struct {
	Asset         string `json:"asset"`
	BorrowEnabled bool   `json:"borrowEnabled"`
	Borrowed      string `json:"borrowed"`
	Free          string `json:"free"`
	Interest      string `json:"interest"`
	Locked        string `json:"locked"`
	NetAsset      string `json:"netAsset"`
	NetAssetOfBtc string `json:"netAssetOfBtc"`
	RepayEnabled  bool   `json:"repayEnabled"`
	TotalAsset    string `json:"totalAsset"`
}

type MarginAssetItem struct {
	Asset    string `json:"asset"`
	Borrowed string `json:"borrowed"`
	Free     string `json:"free"`
	Interest string `json:"interest"`
	Locked   string `json:"locked"`
	NetAsset string `json:"netAsset"`
}

type RespMarginAccount struct {
	BorrowEnabled       bool               `json:"borrowEnabled"`
	MarginLevel         string             `json:"marginLevel"`
	TotalAssetOfBtc     string             `json:"totalAssetOfBtc"`
	TotalLiabilityOfBtc string             `json:"totalLiabilityOfBtc"`
	TotalNetAssetOfBtc  string             `json:"totalNetAssetOfBtc"`
	TradeEnabled        bool               `json:"tradeEnabled"`
	TransferEnabled     bool               `json:"transferEnabled"`
	UserAssets          []*MarginAssetItem `json:"userAssets"`
}

type MarginIsolatedAssertsItem struct {
	BaseAsset         *MarginIsolatedAssetItem `json:"baseAsset"`
	QuoteAsset        *MarginIsolatedAssetItem `json:"quoteAsset"`
	Symbol            string                   `json:"symbol"`
	IsolatedCreated   bool                     `json:"isolatedCreated"`
	Enabled           bool                     `json:"enabled"`
	MarginLevel       string                   `json:"marginLevel"`
	MarginLevelStatus string                   `json:"marginLevelStatus"`
	MarginRatio       string                   `json:"marginRatio"`
	IndexPrice        string                   `json:"indexPrice"`
	LiquidatePrice    string                   `json:"liquidatePrice"`
	LiquidateRate     string                   `json:"liquidateRate"`
	TradeEnabled      bool                     `json:"tradeEnabled"`
}

type RespMarginIsolatedAccount struct {
	Assets              []*MarginIsolatedAssertsItem `json:"assets"`
	TotalAssetOfBtc     string                       `json:"totalAssetOfBtc"`
	TotalLiabilityOfBtc string                       `json:"totalLiabilityOfBtc"`
	TotalNetAssetOfBtc  string                       `json:"totalNetAssetOfBtc"`
}

type RespMarginOrderAck struct {
	Symbol        string `json:"symbol"`
	OrderId       int    `json:"orderId"`
	ClientOrderId string `json:"clientOrderId"`
	IsIsolated    bool   `json:"isIsolated"`
	TransactTime  int64  `json:"transactTime"`
}
type RespMarginOrderResult struct {
	Symbol              string `json:"symbol"`
	OrderId             int    `json:"orderId"`
	ClientOrderId       string `json:"clientOrderId"`
	TransactTime        int64  `json:"transactTime"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	IsIsolated          bool   `json:"isIsolated"`
	Side                string `json:"side"`
}

type FillItem struct {
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
}

type RespMarginOrderFull struct {
	Symbol                string      `json:"symbol"`
	OrderId               int         `json:"orderId"`
	ClientOrderId         string      `json:"clientOrderId"`
	TransactTime          int64       `json:"transactTime"`
	Price                 string      `json:"price"`
	OrigQty               string      `json:"origQty"`
	ExecutedQty           string      `json:"executedQty"`
	CummulativeQuoteQty   string      `json:"cummulativeQuoteQty"`
	Status                string      `json:"status"`
	TimeInForce           string      `json:"timeInForce"`
	Type                  string      `json:"type"`
	Side                  string      `json:"side"`
	MarginBuyBorrowAmount int         `json:"marginBuyBorrowAmount"`
	MarginBuyBorrowAsset  string      `json:"marginBuyBorrowAsset"`
	IsIsolated            bool        `json:"isIsolated"`
	Fills                 []*FillItem `json:"fills"`
}

type RespCancelOrder struct {
	Success bool   `json:"success"`
	Result  string `json:"result"`
}

type RespGetOrder struct {
	Symbol              string `json:"symbol"`              // 交易对
	OrderId             int    `json:"orderId"`             // 系统的订单ID
	OrderListId         int    `json:"orderListId"`         // OCO订单的ID，不然就是-1
	ClientOrderId       string `json:"clientOrderId"`       // 客户自己设置的ID
	Price               string `json:"price"`               // 订单价格
	OrigQty             string `json:"origQty"`             // 用户设置的原始订单数量
	ExecutedQty         string `json:"executedQty"`         // 交易的订单数量
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"` // 累计交易的金额
	Status              string `json:"status"`              // 订单状态
	TimeInForce         string `json:"timeInForce"`         // 订单的时效方式
	Type                string `json:"type"`                // 订单类型， 比如市价单，现价单等
	Side                string `json:"side"`                // 订单方向，买还是卖
	StopPrice           string `json:"stopPrice"`           // 止损价格
	IcebergQty          string `json:"icebergQty"`          // 冰山数量
	Time                int64  `json:"time"`                // 订单时间
	UpdateTime          int64  `json:"updateTime"`          // 最后更新时间
	IsWorking           bool   `json:"isWorking"`           // 订单是否出现在orderbook中
	OrigQuoteOrderQty   string `json:"origQuoteOrderQty"`   // 原始的交易金额
}

type RespGetMarginOrder struct {
	Symbol              string `json:"symbol"`
	OrderId             int    `json:"orderId"`
	ClientOrderId       string `json:"clientOrderId"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	StopPrice           string `json:"stopPrice"`
	IcebergQty          string `json:"icebergQty"`
	Time                int64  `json:"time"`
	UpdateTime          int64  `json:"updateTime"`
	IsWorking           bool   `json:"isWorking"`
	IsIsolated          bool   `json:"isIsolated"`
}

type CapitalDepositHisRecItem struct {
	Amount        string `json:"amount"`
	Coin          string `json:"coin"`
	Network       string `json:"network"`
	Status        int    `json:"status"`
	Address       string `json:"address"`
	AddressTag    string `json:"addressTag"`
	TxId          string `json:"txId"`
	InsertTime    int64  `json:"insertTime"`
	TransferType  int    `json:"transferType"`  // 1: 站内转账, 0: 站外转账
	UnlockConfirm string `json:"unlockConfirm"` // 解锁需要的网络确认次数
	ConfirmTimes  string `json:"confirmTimes"`
}

type RespCapitalDepositHisRec []*CapitalDepositHisRecItem

func IsFutureCoin(symbolType common.SymbolType) bool {
	if symbolType == common.SymbolType_FUTURE_COIN_THIS_WEEK ||
		symbolType == common.SymbolType_FUTURE_COIN_THIS_MONTH ||
		symbolType == common.SymbolType_FUTURE_COIN_THIS_QUARTER ||
		symbolType == common.SymbolType_SWAP_COIN_FOREVER {
		return true
	} else {
		return false
	}
}

func GetMoveSide(src, dst common.SymbolType) MoveType {
	switch {
	case src == common.SymbolType_SPOT_NORMAL && dst == common.SymbolType_MARGIN_NORMAL:
		return MOVE_TYPE_MAIN_MARGIN
	case src == common.SymbolType_SPOT_NORMAL && !IsFutureCoin(dst):
		return MOVE_TYPE_MAIN_UMFUTURE
	case src == common.SymbolType_SPOT_NORMAL && IsFutureCoin(dst):
		return MOVE_TYPE_MAIN_CMFUTURE
	case !IsFutureCoin(src) && dst == common.SymbolType_SPOT_NORMAL:
		return MOVE_TYPE_UMFUTURE_MAIN
	case !IsFutureCoin(src) && dst == common.SymbolType_MARGIN_NORMAL:
		return MOVE_TYPE_UMFUTURE_MARGIN
	case IsFutureCoin(src) && dst == common.SymbolType_SPOT_NORMAL:
		return MOVE_TYPE_CMFUTURE_MAIN
	case src == common.SymbolType_MARGIN_NORMAL && dst == common.SymbolType_SPOT_NORMAL:
		return MOVE_TYPE_MARGIN_MAIN
	case src == common.SymbolType_MARGIN_NORMAL && !IsFutureCoin(dst):
		return MOVE_TYPE_MARGIN_UMFUTURE
	case src == common.SymbolType_MARGIN_NORMAL && IsFutureCoin(dst):
		return MOVE_TYPE_MARGIN_CMFUTURE
	case IsFutureCoin(src) && dst == common.SymbolType_MARGIN_NORMAL:
		return MOVE_TYPE_CMFUTURE_MARGIN
	case src == common.SymbolType_MARGIN_ISOLATED && dst == common.SymbolType_MARGIN_NORMAL:
		return MOVE_TYPE_ISOLATEDMARGIN_MARGIN
	case src == common.SymbolType_MARGIN_NORMAL && dst == common.SymbolType_MARGIN_ISOLATED:
		return MOVE_TYPE_MARGIN_ISOLATEDMARGIN
	case src == common.SymbolType_MARGIN_ISOLATED && dst == common.SymbolType_MARGIN_ISOLATED:
		return MOVE_TYPE_ISOLATEDMARGIN_ISOLATEDMARGIN
	case src == common.SymbolType_SPOT_NORMAL && dst == common.SymbolType_WALLET_NORMAL:
		return MOVE_TYPE_MAIN_FUNDING
	case src == common.SymbolType_WALLET_NORMAL && dst == common.SymbolType_SPOT_NORMAL:
		return MOVE_TYPE_FUNDING_MAIN
	case src == common.SymbolType_WALLET_NORMAL && !IsFutureCoin(dst):
		return MOVE_TYPE_FUNDING_UMFUTURE
	case !IsFutureCoin(src) && dst == common.SymbolType_WALLET_NORMAL:
		return MOVE_TYPE_UMFUTURE_FUNDING
	case src == common.SymbolType_MARGIN_NORMAL && dst == common.SymbolType_WALLET_NORMAL:
		return MOVE_TYPE_MARGIN_FUNDING
	case src == common.SymbolType_WALLET_NORMAL && dst == common.SymbolType_MARGIN_NORMAL:
		return MOVE_TYPE_FUNDING_MARGIN
	case src == common.SymbolType_WALLET_NORMAL && IsFutureCoin(dst):
		return MOVE_TYPE_FUNDING_CMFUTURE
	case IsFutureCoin(src) && dst == common.SymbolType_WALLET_NORMAL:
		return MOVE_TYPE_CMFUTURE_FUNDING
	default:
		return MOVE_TYPE_INVALID
	}
}

func GetDepositTypeToExchange(status client.DepositStatus) DepositStatus {
	switch status {
	case client.DepositStatus_DEPOSITSTATUS_PENDING:
		return DOPSIT_TYPE_PENDING
	case client.DepositStatus_DEPOSITSTATUS_CONFIRMED:
		return DOPSIT_TYPE_CONFIRMED
	case client.DepositStatus_DEPOSITSTATUS_SUCCESS:
		return DOPSIT_TYPE_SUCCESS
	default:
		return -1
	}
}

func GetDepositTypeFromExchange(status DepositStatus) client.DepositStatus {
	switch status {
	case DOPSIT_TYPE_PENDING:
		return client.DepositStatus_DEPOSITSTATUS_PENDING
	case DOPSIT_TYPE_CONFIRMED:
		return client.DepositStatus_DEPOSITSTATUS_CONFIRMED
	case DOPSIT_TYPE_SUCCESS:
		return client.DepositStatus_DEPOSITSTATUS_SUCCESS
	default:
		return client.DepositStatus_DEPOSITSTATUS_INVALID
	}
}

func GetTransferTypeToExchange(status client.TransferStatus) TransferStatus {
	switch status {
	case client.TransferStatus_TRANSFERSTATUS_CREATED:
		return TRANSFER_TYPE_CREATED
	case client.TransferStatus_TRANSFERSTATUS_CANCELLED:
		return TRANSFER_TYPE_CANCELLED
	case client.TransferStatus_TRANSFERSTATUS_CONFORMING:
		return TRANSFER_TYPE_CONFIRMING
	case client.TransferStatus_TRANSFERSTATUS_REJECTED:
		return TRANSFER_TYPE_REJECTED
	case client.TransferStatus_TRANSFERSTATUS_PROCESSING:
		return TRANSFER_TYPE_PROCESSING
	case client.TransferStatus_TRANSFERSTATUS_FAILED:
		return TRANSFER_TYPE_FAILED
	case client.TransferStatus_TRANSFERSTATUS_COMPLETE:
		return TRANSFER_TYPE_SUCCESS
	default:
		return -1
	}
}

func GetTransferTypeFromExchange(status TransferStatus) client.TransferStatus {
	switch status {
	case TRANSFER_TYPE_CREATED:
		return client.TransferStatus_TRANSFERSTATUS_CREATED
	case TRANSFER_TYPE_CANCELLED:
		return client.TransferStatus_TRANSFERSTATUS_CANCELLED
	case TRANSFER_TYPE_CONFIRMING:
		return client.TransferStatus_TRANSFERSTATUS_CONFORMING
	case TRANSFER_TYPE_REJECTED:
		return client.TransferStatus_TRANSFERSTATUS_REJECTED
	case TRANSFER_TYPE_PROCESSING:
		return client.TransferStatus_TRANSFERSTATUS_PROCESSING
	case TRANSFER_TYPE_FAILED:
		return client.TransferStatus_TRANSFERSTATUS_FAILED
	case TRANSFER_TYPE_SUCCESS:
		return client.TransferStatus_TRANSFERSTATUS_COMPLETE
	default:
		return client.TransferStatus_TRANSFERSTATUS_INVALID
	}
}

func GetMoveStatusFromExchange(status MoveStatus) client.MoveStatus {
	switch status {
	case MOVE_STATUS_PENDING:
		return client.MoveStatus_MOVESTATUS_PENDING
	case MOVE_STATUS_CONFIRMED:
		return client.MoveStatus_MOVESTATUS_CONFIRMED
	case MOVE_STATUS_FAILED:
		return client.MoveStatus_MOVESTATUS_FAILED
	default:
		return client.MoveStatus_MOVESTATUS_INVALID
	}
}

func GetLoadStatusFromExchange(status LoanStatus) client.LoanStatus {
	switch status {
	case LOAN_STATUS_PENDING:
		return client.LoanStatus_LOANSTATUS_PENDING
	case LOAN_STATUS_CONFIRMED:
		return client.LoanStatus_LOANSTATUS_CONFIRMED
	case LOAN_STATUS_FAILED:
		return client.LoanStatus_LOANSTATUS_FAILED
	default:
		return client.LoanStatus_LOANSTATUS_INVALID
	}
}
func GetTimeInForceFromExchange(tif TimeInForceType) order.TimeInForce {
	switch tif {
	case TIME_IN_FORCE_IOC:
		return order.TimeInForce_IOC
	case TIME_IN_FORCE_FOK:
		return order.TimeInForce_FOK
	case TIME_IN_FORCE_GTC:
		return order.TimeInForce_GTC
	case TIME_IN_FORCE_GTX:
		return order.TimeInForce_GTX
	default:
		return order.TimeInForce_InvalidTIF
	}
}
func GetTimeInForceToExchange(tif order.TimeInForce) TimeInForceType {
	switch tif {
	case order.TimeInForce_IOC:
		return TIME_IN_FORCE_IOC
	case order.TimeInForce_FOK:
		return TIME_IN_FORCE_FOK
	case order.TimeInForce_GTC:
		return TIME_IN_FORCE_GTC
	default:
		return ""
	}
}

func GetNetWorkFromChain(chain common.Chain) string {
	switch chain {
	case common.Chain_ETH:
		return "ETH"
	case common.Chain_BSC:
		return "BSC"
	case common.Chain_AVALANCHE:
		return "AVAXC"
	case common.Chain_SOLANA:
		return "SOL"
	case common.Chain_TRON:
		return "TRX"
	case common.Chain_POLYGON:
		return "MATIC"
	case common.Chain_ARBITRUM:
		return "ARBITRUM"
	default:
		return ""
	}
}
func GetChainFromNetWork(network string) common.Chain {
	switch network {
	case "erc20":
		return common.Chain_ETH
	case "bsc":
		return common.Chain_BSC
	case "avax":
		return common.Chain_AVALANCHE
	case "sol":
		return common.Chain_SOLANA
	case "trx":
		return common.Chain_TRON
	case "matic":
		return common.Chain_POLYGON
	case "ftm":
		return common.Chain_FANTOM
	default:
		return common.Chain_INVALID_CAHIN
	}
}

func GetTransferStatusFromResponse(status string) client.TransferStatus {
	switch status {
	case "requested":
		return client.TransferStatus_TRANSFERSTATUS_CREATED
	case "processing":
		return client.TransferStatus_TRANSFERSTATUS_PROCESSING
	case "sent":
		return client.TransferStatus_TRANSFERSTATUS_CONFORMING
	case "complete":
		return client.TransferStatus_TRANSFERSTATUS_COMPLETE
	case "cancelled":
		return client.TransferStatus_TRANSFERSTATUS_CANCELLED
	default:
		return client.TransferStatus_TRANSFERSTATUS_INVALID
	}
}

func GetDepositStatusFromResponse(status string) client.DepositStatus {
	switch status {
	case "confirmed":
		return client.DepositStatus_DEPOSITSTATUS_CONFIRMED
	case "unconfirmed":
		return client.DepositStatus_DEPOSITSTATUS_PENDING
	case "cancelled":
		return client.DepositStatus_DEPOSITSTATUS_INVALID
	default:
		return client.DepositStatus_DEPOSITSTATUS_SUCCESS
	}
}

func GetSideTypeFromExchange(side SideType) order.TradeSide {
	switch side {
	case SIDE_TYPE_BUY:
		return order.TradeSide_BUY
	case SIDE_TYPE_SELL:
		return order.TradeSide_SELL
	default:
		return order.TradeSide_InvalidSide
	}
}

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

func GetOrderStatusFromExchange(status OrderStatus, filled float64, size float64) order.OrderStatusCode {
	switch status {
	case ORDER_STATUS_NEW:
		return order.OrderStatusCode_OPENED
	case ORDER_STATUS_OPEN:
		if filled > 0 {
			return order.OrderStatusCode_PARTFILLED
		}
		return order.OrderStatusCode_OPENED
	case ORDER_STATUS_CLOSED: //filled or Cancelled
		if filled != size {
			return order.OrderStatusCode_CANCELED
		}
		return order.OrderStatusCode_FILLED
	default:
		return order.OrderStatusCode_OrderStatusInvalid
	}
}

type CapitalWithdrawHistoryItem struct {
	Address         string `json:"address"`
	Amount          string `json:"amount"`
	ApplyTime       string `json:"applyTime"`
	Coin            string `json:"coin"`
	Id              string `json:"id"`
	WithdrawOrderId string `json:"withdrawOrderId"`
	Network         string `json:"network"`
	TransferType    int    `json:"transferType"`
	Status          int    `json:"status"`
	TransactionFee  string `json:"transactionFee"`
	ConfirmNo       int    `json:"confirmNo"`
	Info            string `json:"info"`
	TxId            string `json:"txId"`
}
type RespCapitalWithdrawHistory []*CapitalWithdrawHistoryItem

type AssetTransferHistoryItem struct {
	Asset     string `json:"asset"`
	Amount    string `json:"amount"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	TranId    int64  `json:"tranId"`
	Timestamp int64  `json:"timestamp"`
}

type RespAssetTransferHistory struct {
	Total int                         `json:"total"`
	Rows  []*AssetTransferHistoryItem `json:"rows"`
}

type RespMarginLoan struct {
	TranId int64 `json:"tranId"`
}
type RespMarginRepay struct {
	TranId int64 `json:"tranId"`
}

type MarginLoanHistoryItem struct {
	IsolatedSymbol string `json:"isolatedSymbol"`
	TxId           int64  `json:"txId"`
	Asset          string `json:"asset"`
	Principal      string `json:"principal"`
	Timestamp      int64  `json:"timestamp"`
	Status         string `json:"status"`
}

type RespMarginLoanHistory struct {
	Rows  []*MarginLoanHistoryItem `json:"rows"`
	Total int                      `json:"total"`
}

type MarginRepayHistoryItem struct {
	IsolatedSymbol string `json:"isolatedSymbol"`
	Amount         string `json:"amount"`
	Asset          string `json:"asset"`
	Interest       string `json:"interest"`
	Principal      string `json:"principal"`
	Status         string `json:"status"`
	Timestamp      int64  `json:"timestamp"`
	TxId           int64  `json:"txId"`
}

type RespMarginRepayHistory struct {
	Rows  []*MarginRepayHistoryItem `json:"rows"`
	Total int                       `json:"total"`
}

type MyTradeItem struct {
	Symbol          string `json:"symbol"`
	Id              int    `json:"id"`
	OrderId         int    `json:"orderId"`
	OrderListId     int    `json:"orderListId"`
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	QuoteQty        string `json:"quoteQty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	Time            int64  `json:"time"`
	IsBuyer         bool   `json:"isBuyer"`
	IsMaker         bool   `json:"isMaker"`
	IsBestMatch     bool   `json:"isBestMatch"`
}

type RespMyTrades []*MyTradeItem

type MarginAllOrdersItem struct {
	ClientOrderId       string `json:"clientOrderId"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	ExecutedQty         string `json:"executedQty"`
	IcebergQty          string `json:"icebergQty"`
	IsWorking           bool   `json:"isWorking"`
	OrderId             int    `json:"orderId"`
	OrigQty             string `json:"origQty"`
	Price               string `json:"price"`
	Side                string `json:"side"`
	Status              string `json:"status"`
	StopPrice           string `json:"stopPrice"`
	Symbol              string `json:"symbol"`
	IsIsolated          bool   `json:"isIsolated"`
	Time                int64  `json:"time"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	UpdateTime          int64  `json:"updateTime"`
}

type RespMarginAllOrders []*MarginAllOrdersItem

type AllOrderItem struct {
	ClientOrderId       string `json:"clientOrderId"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	ExecutedQty         string `json:"executedQty"`
	IcebergQty          string `json:"icebergQty"`
	IsWorking           bool   `json:"isWorking"`
	OrderId             int    `json:"orderId"`
	OrigQty             string `json:"origQty"`
	Price               string `json:"price"`
	Side                string `json:"side"`
	Status              string `json:"status"`
	StopPrice           string `json:"stopPrice"`
	Symbol              string `json:"symbol"`
	Time                int64  `json:"time"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	UpdateTime          int64  `json:"updateTime"`
}

type RespAllOrders []*AllOrderItem

type OpenOrderItem struct {
	Symbol              string `json:"symbol"`
	OrderId             int    `json:"orderId"`
	OrderListId         int    `json:"orderListId"`
	ClientOrderId       string `json:"clientOrderId"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	StopPrice           string `json:"stopPrice"`
	IcebergQty          string `json:"icebergQty"`
	Time                int64  `json:"time"`
	UpdateTime          int64  `json:"updateTime"`
	IsWorking           bool   `json:"isWorking"`
	OrigQuoteOrderQty   string `json:"origQuoteOrderQty"`
}

type MarginOpenOrderItem struct {
	ClientOrderId       string `json:"clientOrderId"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	ExecutedQty         string `json:"executedQty"`
	IcebergQty          string `json:"icebergQty"`
	IsWorking           bool   `json:"isWorking"`
	OrderId             int    `json:"orderId"`
	OrigQty             string `json:"origQty"`
	Price               string `json:"price"`
	Side                string `json:"side"`
	Status              string `json:"status"`
	StopPrice           string `json:"stopPrice"`
	Symbol              string `json:"symbol"`
	IsIsolated          bool   `json:"isIsolated"`
	Time                int64  `json:"time"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	UpdateTime          int64  `json:"updateTime"`
}

type RespOpenOrders []*OpenOrderItem
type RespMarginOpenOrders []*MarginOpenOrderItem
type RespUserDataStream struct {
	ListenKey string `json:"listenKey"`
}

type RespIsolatedPair struct {
	Symbol        string `json:"symbol"`
	Base          string `json:"base"`
	Quote         string `json:"quote"`
	IsMarginTrade bool   `json:"isMarginTrade"`
	IsBuyAllowed  bool   `json:"isBuyAllowed"`
	IsSellAllowed bool   `json:"isSellAllowed"`
}

type RespIsolatedAllPairs []*RespIsolatedPair
