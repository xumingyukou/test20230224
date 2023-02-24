package spot_api

import (
	"strconv"
	"strings"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/order"
)

type RespSymbol struct {
	//现货、杠杆
	Symbol                 string    `json:"symbol"`
	Status                 string    `json:"status"`             // 交易对状态
	BaseAsset              string    `json:"baseAsset"`          // 标的资产
	BaseAssetPrecision     int       `json:"baseAssetPrecision"` // 标的资产精度
	QuoteAsset             string    `json:"quoteAsset"`         // 报价资产
	QuotePrecision         int       `json:"quotePrecision"`     // 报价资产精度
	QuoteAssetPrecision    int       `json:"quoteAssetPrecision"`
	OrderTypes             []string  `json:"orderTypes"`
	IcebergAllowed         bool      `json:"icebergAllowed"`
	OcoAllowed             bool      `json:"ocoAllowed"`
	IsSpotTradingAllowed   bool      `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed bool      `json:"isMarginTradingAllowed"`
	AllowTrailingStop      bool      `json:"allowTrailingStop"`
	Filters                []*Filter `json:"filters"`
	Permissions            []string  `json:"permissions"`
	//u本位
	Pair                  string   `json:"pair"`                  // 标的交易对
	ContractType          string   `json:"contractType"`          // 合约类型
	DeliveryDate          int64    `json:"deliveryDate"`          // 交割日期
	OnboardDate           int64    `json:"onboardDate"`           // 上线日期
	MaintMarginPercent    string   `json:"maintMarginPercent"`    // 请忽略
	RequiredMarginPercent string   `json:"requiredMarginPercent"` // 请忽略
	MarginAsset           string   `json:"marginAsset"`           // 保证金资产
	PricePrecision        int      `json:"pricePrecision"`        // 价格小数点位数(仅作为系统精度使用，注意同tickSize 区分）
	QuantityPrecision     int      `json:"quantityPrecision"`     // 数量小数点位数(仅作为系统精度使用，注意同stepSize 区分）
	UnderlyingType        string   `json:"underlyingType"`        //"COIN"
	UnderlyingSubType     []string `json:"underlyingSubType"`     //["STORAGE"]
	SettlePlan            int      `json:"settlePlan"`            //0
	TriggerProtect        string   `json:"triggerProtect"`        // 开启"priceProtect"的条件订单的触发阈值
	OrderType             []string `json:"OrderType"`             // 订单类型
	TimeInForce           []string `json:"timeInForce"`           // 有效方式
	LiquidationFee        string   `json:"liquidationFee"`        // 强平费率
	MarketTakeBound       string   `json:"marketTakeBound"`       // 市价吃单(相对于标记价格)允许可造成的最大价格偏离比例
}

type Precision struct {
	MinAmount       float64 `json:"min_amount"`
	AmountPrecision int     `json:"amount_precision"`
	MinPrice        float64 `json:"min_price"`
	MinValue        float64 `json:"min_value"`
	PricePrecision  int     `json:"price_precision"`
}

func (ts RespSymbol) GetPrecision() Precision {
	precision := Precision{}
	for _, v := range ts.Filters {
		if v.FilterType == "LOT_SIZE" {
			precision.MinAmount = v.MinQty
		}
		if v.FilterType == "LOT_SIZE" {
			step := strconv.FormatFloat(v.StepSize, 'f', -1, 64)
			pres := strings.Split(step, ".")
			if len(pres) != 1 {
				precision.AmountPrecision = len(pres[1])
			}
		}
		if v.FilterType == "PRICE_FILTER" {
			precision.MinPrice = v.MinPrice
		}
		if v.FilterType == "MIN_NOTIONAL" {
			precision.MinValue = v.MinNotional
		}
		if v.FilterType == "PRICE_FILTER" {
			step := strconv.FormatFloat(v.TickSize, 'f', -1, 64)
			pres := strings.Split(step, ".")
			if len(pres) != 1 {
				precision.PricePrecision = len(pres[1])
			}
		}
	}
	return precision
}

type RespRateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

type RespExchangeInfo struct {
	Timezone        string           `json:"timezone"`
	ServerTime      int64            `json:"serverTime"`
	RateLimits      []*RespRateLimit `json:"rateLimits"`
	Assets          []AssetItem      `json:"assets"`
	ExchangeFilters interface{}      `json:"exchangeFilters"`
	Symbols         []*RespSymbol    `json:"symbols"`
}

type AssetItem struct {
	Asset             string `json:"asset"`
	MarginAvailable   bool   `json:"marginAvailable"`   // 是否可用作保证金
	AutoAssetExchange string `json:"autoAssetExchange"` // 保证金资产自动兑换阈值
}

type ResultError interface {
	ErrorCode() int32
	ErrorMsg() string
}

type RespError struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

func (re *RespError) ErrorCode() int32 {
	return re.Code
}

func (re *RespError) ErrorMsg() string {
	return re.Msg
}

type RespTime struct {
	ServerTime int64 `json:"serverTime"`
}

type RespPing struct{}

type RespDepth struct {
	LastUpdateId int64      `json:"lastUpdateId"`
	Symbol       string     `json:"symbol"` // 交易对（仅币本位合约）
	Pair         string     `json:"pair"`   // 标的交易对（仅币本位合约）
	E            int64      `json:"E"`      // 消息时间（合约）
	T            int64      `json:"T"`      // 撮合时间（合约）
	Bids         [][]string `json:"bids"`   // 买单
	Asks         [][]string `json:"asks"`   // 卖单
}

type TradeInfo struct {
	Id           int    `json:"id"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
	IsBestMatch  bool   `json:"isBestMatch"`
}

type RespGetTrade []*TradeInfo

type RespTickerPriceItem struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Time   int64  `json:"time"`
}

type RespTickerPrice []*RespTickerPriceItem

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
	MultiplierDecimal string `json:"multiplierDecimal,omitempty"`
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

type RespAccount struct {
	MakerCommission  int            `json:"makerCommission"`
	TakerCommission  int            `json:"takerCommission"`
	BuyerCommission  int            `json:"buyerCommission"`
	SellerCommission int            `json:"sellerCommission"`
	CanTrade         bool           `json:"canTrade"`
	CanWithdraw      bool           `json:"canWithdraw"`
	CanDeposit       bool           `json:"canDeposit"`
	UpdateTime       int64          `json:"updateTime"`
	AccountType      string         `json:"accountType"`
	Balances         []*BalanceItem `json:"balances"`
	Permissions      []string       `json:"permissions"`
}

type SpotAssetItem struct {
	Asset        string `json:"asset"`
	Free         string `json:"free"`
	Locked       string `json:"locked"`
	Freeze       string `json:"freeze"`
	Withdrawing  string `json:"withdrawing"`
	BtcValuation string `json:"btcValuation"`
}

type RespAsset []*SpotAssetItem

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
	ClientOrderId       string `json:"clientOrderId"` // 用户自定义的订单号
	ExecutedQty         string `json:"executedQty"`   // 成交量
	OrderId             int    `json:"orderId"`       // 系统订单号
	OrigQty             string `json:"origQty"`       // 原始委托数量
	Price               string `json:"price"`         // 委托价格
	Side                string `json:"side"`          // 买卖方向
	Status              string `json:"status"`        // 订单状态
	Symbol              string `json:"symbol"`        // 交易对
	TimeInForce         string `json:"timeInForce"`   // 有效方法
	Type                string `json:"type"`          // 订单类型
	OrigClientOrderId   string `json:"origClientOrderId"`
	OrderListId         int    `json:"orderListId"` // OCO订单ID，否则为 -1
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	IsIsolated          bool   `json:"isIsolated"`
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
	UnlockConfirm int    `json:"unlockConfirm"` // 解锁需要的网络确认次数
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

func GetMoveSide(src, dst common.Market) MoveType {
	switch {
	case src == common.Market_SPOT && dst == common.Market_MARGIN:
		return MOVE_TYPE_MAIN_MARGIN
	case src == common.Market_SPOT && dst == common.Market_FUTURE:
		return MOVE_TYPE_MAIN_UMFUTURE
	case src == common.Market_SPOT && dst == common.Market_FUTURE_COIN:
		return MOVE_TYPE_MAIN_CMFUTURE
	case src == common.Market_FUTURE && dst == common.Market_SPOT:
		return MOVE_TYPE_UMFUTURE_MAIN
	case src == common.Market_FUTURE && dst == common.Market_MARGIN:
		return MOVE_TYPE_UMFUTURE_MARGIN
	case src == common.Market_FUTURE_COIN && dst == common.Market_SPOT:
		return MOVE_TYPE_CMFUTURE_MAIN
	case src == common.Market_MARGIN && dst == common.Market_SPOT:
		return MOVE_TYPE_MARGIN_MAIN
	case src == common.Market_MARGIN && dst == common.Market_FUTURE:
		return MOVE_TYPE_MARGIN_UMFUTURE
	case src == common.Market_MARGIN && dst == common.Market_FUTURE_COIN:
		return MOVE_TYPE_MARGIN_CMFUTURE
	case src == common.Market_FUTURE_COIN && dst == common.Market_MARGIN:
		return MOVE_TYPE_CMFUTURE_MARGIN
	/*
		case src == common.SymbolType_MARGIN_ISOLATED && dst == common.Market_MARGIN:
			return MOVE_TYPE_ISOLATEDMARGIN_MARGIN
		case src == common.Market_MARGIN && dst == common.SymbolType_MARGIN_ISOLATED:
			return MOVE_TYPE_MARGIN_ISOLATEDMARGIN
		case src == common.SymbolType_MARGIN_ISOLATED && dst == common.SymbolType_MARGIN_ISOLATED:
			return MOVE_TYPE_ISOLATEDMARGIN_ISOLATEDMARGIN
	*/
	case src == common.Market_SPOT && dst == common.Market_WALLET:
		return MOVE_TYPE_MAIN_FUNDING
	case src == common.Market_WALLET && dst == common.Market_SPOT:
		return MOVE_TYPE_FUNDING_MAIN
	case src == common.Market_WALLET && dst == common.Market_FUTURE:
		return MOVE_TYPE_FUNDING_UMFUTURE
	case src == common.Market_FUTURE && dst == common.Market_WALLET:
		return MOVE_TYPE_UMFUTURE_FUNDING
	case src == common.Market_MARGIN && dst == common.Market_WALLET:
		return MOVE_TYPE_MARGIN_FUNDING
	case src == common.Market_WALLET && dst == common.Market_MARGIN:
		return MOVE_TYPE_FUNDING_MARGIN
	case src == common.Market_WALLET && dst == common.Market_FUTURE_COIN:
		return MOVE_TYPE_FUNDING_CMFUTURE
	case src == common.Market_FUTURE_COIN && dst == common.Market_WALLET:
		return MOVE_TYPE_CMFUTURE_FUNDING
	default:
		return MOVE_TYPE_INVALID
	}
}

func GetSubMoveSide(market common.Market) MoveAccountType {
	switch {
	case market == common.Market_SPOT:
		return MOVE_ACCOUNT_TYPE_SPOT
	case market == common.Market_MARGIN:
		return MOVE_ACCOUNT_TYPE_MARGIN
	case market == common.Market_FUTURE || market == common.Market_SWAP:
		return MOVE_ACCOUNT_TYPE_FUTURE
	case market == common.Market_FUTURE_COIN || market == common.Market_SWAP_COIN:
		return MOVE_ACCOUNT_TYPE_C_FUTURE
	default:
		return MOVE_ACCOUNT_TYPE_INVALID
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
	case MOVE_STATUS_CONFIRMED, MOVE_STATUS_SUCCESS:
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

func GetOrderTypeToExchange(ot order.OrderType) OrderType {
	switch ot {
	case order.OrderType_MARKET:
		return ORDER_TYPE_MARKET
	case order.OrderType_LIMIT_MAKER:
		return ORDER_TYPE_LIMIT_MAKER
	case order.OrderType_LIMIT:
		return ORDER_TYPE_LIMIT
	case order.OrderType_STOP_LOSS:
		return ORDER_TYPE_STOP_LOSS
	case order.OrderType_STOP_LOSS_LIMIT:
		return ORDER_TYPE_STOP_LOSS_LIMIT
	case order.OrderType_TAKE_PROFIT:
		return ORDER_TYPE_TAKE_PROFIT
	case order.OrderType_TAKE_PROFIT_LIMIT:
		return ORDER_TYPE_TAKE_PROFIT_LIMIT
	case order.OrderType_TRAILING_STOP:
		return ORDER_TYPE_TRAILING_STOP_MARKET
	default:
		return ""
	}
}

func GetOrderTypeFromExchange(ot OrderType) order.OrderType {
	switch ot {
	case ORDER_TYPE_MARKET:
		return order.OrderType_MARKET
	case ORDER_TYPE_LIMIT_MAKER:
		return order.OrderType_LIMIT_MAKER
	case ORDER_TYPE_LIMIT:
		return order.OrderType_LIMIT
	case ORDER_TYPE_STOP_LOSS:
		return order.OrderType_STOP_LOSS
	case ORDER_TYPE_STOP_LOSS_LIMIT:
		return order.OrderType_STOP_LOSS_LIMIT
	case ORDER_TYPE_TAKE_PROFIT:
		return order.OrderType_TAKE_PROFIT
	case ORDER_TYPE_TAKE_PROFIT_LIMIT:
		return order.OrderType_TAKE_PROFIT_LIMIT
	case ORDER_TYPE_TRAILING_STOP_MARKET:
		return order.OrderType_TRAILING_STOP
	default:
		return order.OrderType_InvalidOrder
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
	case order.TimeInForce_GTX:
		return TIME_IN_FORCE_GTX
	default:
		return ""
	}
}

func GetNetWorkFromChain(chain common.Chain) string {
	switch chain {
	case common.Chain_CRONOS:
		return "ETH"
	case common.Chain_BSC:
		return "BSC"
	//case common.Chain_BNB:
	//	return "BNB"
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
	case common.Chain_OPTIMISM:
		return "OPTIMISM"
	default:
		return ""
	}
}

func GetChainFromNetWork(network string) common.Chain {
	switch network {
	case "ETH":
		return common.Chain_ETH
	case "BSC":
		return common.Chain_BSC
	//case "BNB":
	//	return common.Chain_BNB
	case "AVAXC":
		return common.Chain_AVALANCHE
	case "SOL":
		return common.Chain_SOLANA
	case "TRX":
		return common.Chain_TRON
	case "MATIC":
		return common.Chain_POLYGON
	case "ARBITRUM":
		return common.Chain_ARBITRUM
	case "OPTIMISM":
		return common.Chain_OPTIMISM
	default:
		return common.Chain_INVALID_CAHIN
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
	case order.TradeSide_BUY_TO_OPEN:
		return SIDE_TYPE_BUY
	case order.TradeSide_BUY_TO_CLOSE:
		return SIDE_TYPE_BUY
	case order.TradeSide_SELL_TO_OPEN:
		return SIDE_TYPE_SELL
	case order.TradeSide_SELL_TO_CLOSE:
		return SIDE_TYPE_SELL
	default:
		return ""
	}
}

func GetOrderStatusFromExchange(status OrderStatus) order.OrderStatusCode {
	switch status {
	case ORDER_STATUS_LIVE:
		return order.OrderStatusCode_OPENED
	case ORDER_STATUS_NEW:
		return order.OrderStatusCode_OPENED
	case ORDER_STATUS_PARTIALLY_FILLED:
		return order.OrderStatusCode_PARTFILLED
	case ORDER_STATUS_FILLED:
		return order.OrderStatusCode_FILLED
	case ORDER_STATUS_CANCELED:
		return order.OrderStatusCode_CANCELED
	case ORDER_STATUS_PENDING_CANCEL:
		return order.OrderStatusCode_CANCELED
	case ORDER_STATUS_REJECTED:
		return order.OrderStatusCode_FAILED
	case ORDER_STATUS_EXPIRED:
		return order.OrderStatusCode_EXPIRED
	default:
		return order.OrderStatusCode_OrderStatusInvalid
	}
}
func GetOrderStatusToExchange(status order.OrderStatusCode) OrderStatus {
	switch status {
	case order.OrderStatusCode_OPENED:
		return ORDER_STATUS_NEW
	case order.OrderStatusCode_PARTFILLED:
		return ORDER_STATUS_PARTIALLY_FILLED
	case order.OrderStatusCode_FILLED:
		return ORDER_STATUS_FILLED
	case order.OrderStatusCode_CANCELED:
		return ORDER_STATUS_CANCELED
	case order.OrderStatusCode_FAILED:
		return ORDER_STATUS_REJECTED
	case order.OrderStatusCode_EXPIRED:
		return ORDER_STATUS_EXPIRED
	default:
		return ""
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
type SubAccountUniversalTransfer struct {
	TranId          int64  `json:"tranId"`
	FromEmail       string `json:"fromEmail"`
	ToEmail         string `json:"toEmail"`
	Asset           string `json:"asset"`
	Amount          string `json:"amount"`
	CreateTimeStamp int64  `json:"createTimeStamp"`
	FromAccountType string `json:"fromAccountType"`
	ToAccountType   string `json:"toAccountType"`
	Status          string `json:"status"`
	ClientTranId    string `json:"clientTranId"`
}

type SubAccountTransferSubUserHistory struct {
	CounterParty    string `json:"counterParty"`
	Email           string `json:"email"`
	Type            int    `json:"type"` // 1 for transfer in , 2 for transfer out
	Asset           string `json:"asset"`
	Qty             string `json:"qty"`
	FromAccountType string `json:"fromAccountType"`
	ToAccountType   string `json:"toAccountType"`
	Status          string `json:"status"`
	TranId          int64  `json:"tranId"`
	Time            int64  `json:"time"`
}

type RespAssetTransferHistory struct {
	Total int                         `json:"total"`
	Rows  []*AssetTransferHistoryItem `json:"rows"`
}

type RespSubAccountUniversalTransfer struct {
	Result     []*SubAccountUniversalTransfer `json:"result"`
	TotalCount int                            `json:"totalCount"`
}

type RespSubAccountTransferSubUserHistory []*SubAccountTransferSubUserHistory

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

	AvgPrice      string `json:"avgPrice"`
	CumQuote      string `json:"cumQuote"`
	OrigType      string `json:"origType"`
	ReduceOnly    bool   `json:"reduceOnly"`
	PositionSide  string `json:"positionSide"`
	ClosePosition bool   `json:"closePosition"`
	ActivatePrice string `json:"activatePrice"`
	PriceRate     string `json:"priceRate"`
	WorkingType   string `json:"workingType"`
	PriceProtect  bool   `json:"priceProtect"`
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

func GetChangeType(event string) client.BalanceChangeType {
	switch event {
	case "DEPOSIT":
		return client.BalanceChangeType_CHANGE_DEPOSIT
	case "WITHDRAW":
		return client.BalanceChangeType_CHANGE_WITHDRAW
	case "ORDER":
		return client.BalanceChangeType_CHANGE_ORDER
	case "FUNDING_FEE":
		return client.BalanceChangeType_CHANGE_FUNDING_FEE
	case "WITHDRAW_REJECT":
		return client.BalanceChangeType_CHANGE_WITHDRAW_REJECT
	case "ADJUSTMENT":
		return client.BalanceChangeType_CHANGE_ADJUSTMENT
	case "INSURANCE_CLEAR":
		return client.BalanceChangeType_CHANGE_INSURANCE_CLEAR
	case "ADMIN_DEPOSIT":
		return client.BalanceChangeType_CHANGE_ADMIN_DEPOSIT
	case "ADMIN_WITHDRAW":
		return client.BalanceChangeType_CHANGE_ADMIN_WITHDRAW
	case "MARGIN_TRANSFER":
		return client.BalanceChangeType_CHANGE_MARGIN_TRANSFER
	case "MARGIN_TYPE_CHANGE":
		return client.BalanceChangeType_CHANGE_MARGIN_TYPE_CHANGE
	case "ASSET_TRANSFER":
		return client.BalanceChangeType_CHANGE_ASSET_TRANSFER
	case "OPTIONS_PREMIUM_FEE":
		return client.BalanceChangeType_CHANGE_OPTIONS_PREMIUM_FEE
	case "OPTIONS_SETTLE_PROFIT":
		return client.BalanceChangeType_CHANGE_OPTIONS_SETTLE_PROFIT
	case "AUTO_EXCHANGE":
		return client.BalanceChangeType_CHANGE_AUTO_EXCHANGE
	default:
		return client.BalanceChangeType_CHANGE_INVALID
	}
}

func GetOrderType(o OrderType) order.OrderType {
	switch o {
	case ORDER_TYPE_LIMIT:
		return order.OrderType_LIMIT
	case ORDER_TYPE_MARKET:
		return order.OrderType_MARKET
	case ORDER_TYPE_STOP_LOSS:
		return order.OrderType_STOP_LOSS
	case ORDER_TYPE_STOP_LOSS_LIMIT:
		return order.OrderType_STOP_LOSS_LIMIT
	case ORDER_TYPE_TAKE_PROFIT:
		return order.OrderType_TAKE_PROFIT
	case ORDER_TYPE_TAKE_PROFIT_LIMIT:
		return order.OrderType_TAKE_PROFIT_LIMIT
	case ORDER_TYPE_LIMIT_MAKER:
		return order.OrderType_LIMIT_MAKER
	case ORDER_TYPE_TRAILING_STOP_MARKET:
		return order.OrderType_STOP
	default:
		return order.OrderType_InvalidOrder
	}
}

type SubAccount struct {
	Email                       string `json:"email"`
	IsFreeze                    bool   `json:"isFreeze"`
	CreateTime                  int64  `json:"createTime"`
	IsManagedSubAccount         bool   `json:"isManagedSubAccount"`
	IsAssetManagementSubAccount bool   `json:"isAssetManagementSubAccount"`
}

type RespSubAccountList struct {
	SubAccounts []*SubAccount `json:"subAccounts"`
}
