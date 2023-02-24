package spot_api

type RespError struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type RespGetSymbols struct {
	Code string       `json:"code"`
	Data []SymbolInfo `json:"data"`
}

type SymbolInfo struct {
	Symbol          string `json:"symbol"`
	Name            string `json:"name"`
	BaseCurrency    string `json:"baseCurrency"`
	QuoteCurrency   string `json:"quoteCurrency"`
	FeeCurrency     string `json:"feeCurrency"`
	Market          string `json:"market"`
	BaseMinSize     string `json:"baseMinSize"`
	QuoteMinSize    string `json:"quoteMinSize"`
	BaseMaxSize     string `json:"baseMaxSize"`
	QuoteMaxSize    string `json:"quoteMaxSize"`
	BaseIncrement   string `json:"baseIncrement"`
	QuoteIncrement  string `json:"quoteIncrement"`
	PriceIncrement  string `json:"priceIncrement"`
	PriceLimitRate  string `json:"priceLimitRate"`
	IsMarginEnabled bool   `json:"isMarginEnabled"`
	EnableTrading   bool   `json:"enableTrading"`
}

type RespGetOrderbookInfo struct {
	Code string `json:"code"`
	Data struct {
		Sequence    string `json:"sequence"`
		Price       string `json:"price"`
		Size        string `json:"size"`
		BestAsk     string `json:"bestAsk"`
		BestAskSize string `json:"bestAskSize"`
		BestBid     string `json:"bestBid"`
		BestBidSize string `json:"bestBidSize"`
		Time        uint64 `json:"time"`
	} `json:"data"`
}

type RespGetTradeFees struct {
	Code string          `json:"code"`
	Data []*TradeFeeInfo `json:"data"`
}

type TradeFeeInfo struct {
	Symbol       string `json:"symbol"`
	TakerFeeRate string `json:"takerFeeRate"`
	MakerFeeRate string `json:"makerFeeRate"`
}

type RespGetMarkets struct {
	Code string `json:"code"`
	Data struct {
		Sequence string     `json:"sequence"`
		Bids     [][]string `json:"bids"`
		Asks     [][]string `json:"asks"`
		Time     uint64     `json:"time"`
	}
}

type RespFutureMarkets struct {
	Code string `json:"code"`
	Data struct {
		Symbol   string          `json:"symbol"`
		Sequence int             `json:"sequence"`
		Asks     [][]interface{} `json:"asks"`
		Bids     [][]interface{} `json:"bids"`
		Ts       int64           `json:"ts"`
	}
}

type RespGetCurrencies struct {
	Code string         `json:"code"`
	Data []CurrencyInfo `json:"data"`
}

type CurrencyInfo struct {
	Currency          string `json:"currency"`
	Name              string `json:"name"`
	FullName          string `json:"fullName"`
	Precision         int64  `json:"precision"`
	Confirms          int64  `json:"confirms"`
	ContractAddress   string `json:"contractAddress"`
	WithdrawalMinSize string `json:"withdrawalMinSize"`
	WithdrawalMinFee  string `json:"withdrawalMinFee"`
	IsWithdrawEnabled bool   `json:"isWithdrawEnabled"`
	IsDepositEnabled  bool   `json:"isDepositEnabled"`
	IsMarginEnabled   bool   `json:"isMarginEnabled"`
	IsDebitEnabled    bool   `json:"isDebitEnabled"`
}

type RespGetAccounts struct {
	Code string        `json:"code"`
	Data []AccountInfo `json:"data"`
}

type AccountInfo struct {
	Id        string `json:"id"`
	Currency  string `json:"currency"`
	Type      string `json:"type"`
	Balance   string `json:"balance"`
	Available string `json:"available"`
	Holds     string `json:"holds"`
}

type RespGetMarginAccounts struct {
	Code string                `json:"code"`
	Data MarginAccountInfoWrap `json:"data"`
}

type MarginAccountInfoWrap struct {
	Accounts  []MarginAccountInfo `json:"accounts"`
	DebtRatio string              `json:"debtRatio"`
}

type MarginAccountInfo struct {
	Currency         string `json:"currency"`
	AvailableBalance string `json:"availableBalance"`
	HoldBalance      string `json:"holdBalance"`
	Liability        string `json:"liability"`
	MaxBorrowSize    string `json:"maxBorrowSize"`
	TotalBalance     string `json:"totalBalance"`
	Interest         string `json:"interest"`
	BorrowableAmount string `json:"borrowableAmount"`
}

type RespGetIsolatedAccounts struct {
	Code string `json:"code"`
	Data struct {
		TotalConversionBalance     string `json:"totalConversionBalance"`
		LiabilityConversionBalance string `json:"liabilityConversionBalance"`
		Assets                     []struct {
			Symbol     string            `json:"symbol"`
			Status     string            `json:"status"`
			DebtRatio  string            `json:"debtRatio"`
			BaseAsset  MarginAccountInfo `json:"baseAsset"`
			QuoteAsset MarginAccountInfo `json:"quoteAsset"`
		} `json:"assets"`
	} `json:"data"`
}

type RespPlaceOrder struct {
	Code string `json:"code"`
	Data struct {
		OrderId string `json:"orderId"`
	} `json:"data"`
}

type RespGetOrder struct {
	Code string    `json:"code"`
	Data OrderInfo `json:"data"`
}

type OrderInfo struct {
	Id            string `json:"id"`
	Symbol        string `json:"symbol"`
	OpType        string `json:"opType"`
	Type          string `json:"type"`
	Side          string `json:"side"`
	Price         string `json:"price"`
	Size          string `json:"size"`
	Funds         string `json:"funds"`
	DealFunds     string `json:"dealFunds"`
	DealSize      string `json:"dealSize"`
	Fee           string `json:"fee"`
	FeeCurrency   string `json:"feeCurrency"`
	Stp           string `json:"stp"`
	Stop          string `json:"stop"`
	StopTriggered bool   `json:"stopTriggered"`
	StopPrice     string `json:"stopPrice"`
	TimeInForce   string `json:"timeInForce"`
	PostOnly      bool   `json:"postOnly"`
	Hidden        bool   `json:"hidden"`
	Iceberg       bool   `json:"iceberg"`
	VisibleSize   string `json:"visibleSize"`
	CancelAfter   int64  `json:"cancelAfter"`
	Channel       string `json:"channel"`
	ClientOid     string `json:"clientOid"`
	Remark        string `json:"remark"`
	Tags          string `json:"tags"`
	IsActive      bool   `json:"isActive"`
	CancelExist   bool   `json:"cancelExist"`
	CreatedAt     int64  `json:"createdAt"`
	TradeType     string `json:"tradeType"`
}

type RespGetFills struct {
	Code string `json:"code"`
	Data struct {
		CurrentPage int64       `json:"currentPage"`
		PageSize    int64       `json:"pageSize"`
		TotalNum    int64       `json:"totalNum"`
		TotalPage   int64       `json:"totalPage"`
		Items       []*FillInfo `json:"items"`
	} `json:"data"`
}

type FillInfo struct {
	Symbol         string `json:"symbol"`
	TradeId        string `json:"tradeId"`
	OrderId        string `json:"orderId"`
	CounterOrderId string `json:"counterOrderId"`
	Side           string `json:"side"`
	ForceTaker     bool   `json:"forceTaker"`
	Liquidity      string `json:"liquidity"`
	Price          string `json:"price"`
	Size           string `json:"size"`
	Funds          string `json:"funds"`
	Fee            string `json:"fee"`
	FeeRate        string `json:"feeRate"`
	FeeCurrency    string `json:"feeCurrency"`
	Stop           string `json:"stop"`
	Type           string `json:"type"`
	CreatedAt      int64  `json:"createdAt"`
	TradeType      string `json:"tradeType"`
}

type RespPlaceMarginOrder struct {
	Code string `json:"code"`
	Data struct {
		OrderId     string  `json:"orderId"`
		BorrowSize  float64 `json:"borrowSize"`
		LoanApplyId string  `json:"loanApplyId"`
	} `json:"data"`
}

type RespGetStatus struct {
	Code string `json:"code"`
	Data struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	} `json:"data"`
}

type RespInnerTransfer struct {
	Code string `json:"code"`
	Data struct {
		OrderId string `json:"orderId"`
	} `json:"data"`
}

type RespSubTransfer struct {
	Code string `json:"code"`
	Data struct {
		OrderId string `json:"orderId"`
	} `json:"data"`
}

type RespCancelOrder struct {
	Code string `json:"code"`
	Data struct {
		CancelledOrderId string `json:"cancelledOrderId"`
		ClientOid        string `json:"clientOid"`
	} `json:"data"`
}

type RespOrderHistory struct {
	Code string `json:"code"`
	Data struct {
		CurrentPage int `json:"currentPage"`
		PageSize    int `json:"pageSize"`
		TotalNum    int `json:"totalNum"`
		TotalPage   int `json:"totalPage"`
		Items       []struct {
			Id            string      `json:"id"`
			Symbol        string      `json:"symbol"`
			OpType        string      `json:"opType"`
			Type          string      `json:"type"`
			Side          string      `json:"side"`
			Price         string      `json:"price"`
			Size          string      `json:"size"`
			Funds         string      `json:"funds"`
			DealFunds     string      `json:"dealFunds"`
			DealSize      string      `json:"dealSize"`
			Fee           string      `json:"fee"`
			FeeCurrency   string      `json:"feeCurrency"`
			Stp           string      `json:"stp"`
			Stop          string      `json:"stop"`
			StopTriggered bool        `json:"stopTriggered"`
			StopPrice     string      `json:"stopPrice"`
			TimeInForce   string      `json:"timeInForce"`
			PostOnly      bool        `json:"postOnly"`
			Hidden        bool        `json:"hidden"`
			Iceberg       bool        `json:"iceberg"`
			VisibleSize   string      `json:"visibleSize"`
			CancelAfter   int         `json:"cancelAfter"`
			Channel       string      `json:"channel"`
			ClientOid     string      `json:"clientOid"`
			Remark        interface{} `json:"remark"`
			Tags          interface{} `json:"tags"`
			IsActive      bool        `json:"isActive"`
			CancelExist   bool        `json:"cancelExist"`
			CreatedAt     int64       `json:"createdAt"`
			TradeType     string      `json:"tradeType"`
		} `json:"items"`
	} `json:"data"`
}

type RespGetWsToken struct {
	Code string `json:"code"`
	Data struct {
		InstanceServers []struct {
			Endpoint     string `json:"endpoint"`
			Protocol     string `json:"protocol"`
			Encrypt      bool   `json:"encrypt"`
			PingInterval int64  `json:"pingInterval"`
			PingTimeout  int64  `json:"pingTimeout"`
		} `json:"instanceServers"`
		Token string `json:"token"`
	} `json:"data"`
}

type RespWithdraw struct {
	Code string `json:"code"`
	Data struct {
		WithdrawalId string `json:"withdrawalId"`
	} `json:"data"`
}

type RespWithdrawHistory struct {
	Code string `json:"code"`
	Data struct {
		CurrentPage int `json:"currentPage"`
		PageSize    int `json:"pageSize"`
		TotalNum    int `json:"totalNum"`
		TotalPage   int `json:"totalPage"`
		Items       []struct {
			Currency   string `json:"currency"`
			CreateAt   int    `json:"createAt"`
			Amount     string `json:"amount"`
			Address    string `json:"address"`
			WalletTxId string `json:"walletTxId"`
			IsInner    bool   `json:"isInner"`
			Status     string `json:"status"`
		} `json:"items"`
	} `json:"data"`
}

type RespWithdrawFee struct {
	Code string `json:"code"`
	Data struct {
		Currency            string `json:"currency"`
		LimitBTCAmount      string `json:"limitBTCAmount"`
		UsedBTCAmount       string `json:"usedBTCAmount"`
		RemainAmount        string `json:"remainAmount"`
		AvailableAmount     string `json:"availableAmount"`
		WithdrawMinFee      string `json:"withdrawMinFee"`
		InnerWithdrawMinFee string `json:"innerWithdrawMinFee"`
		WithdrawMinSize     string `json:"withdrawMinSize"`
		IsWithdrawEnabled   bool   `json:"isWithdrawEnabled"`
		Precision           int    `json:"precision"`
		Chain               string `json:"chain"`
	} `json:"data"`
}

type RespGetHistWithdraw struct {
	Code string `json:"code"`
	Data struct {
		CurrentPage int64           `json:"currentPage"`
		PageSize    int64           `json:"pageSize"`
		TotalNum    int64           `json:"totalNum"`
		TotalPage   int64           `json:"totalPage"`
		Items       []*WithdrawInfo `json:"items"`
	} `json:"data"`
}

type WithdrawInfo struct {
	Currency   string `json:"currency"`
	CreateAt   int64  `json:"createAt"`
	Amount     string `json:"amount"`
	Address    string `json:"address"`
	WalletTxId string `json:"walletTxId"`
	IsInner    bool   `json:"isInner"`
	Status     string `json:"status"`
}

type RespMoveHistory struct {
	Code string `json:"code"`
	Data struct {
		CurrentPage int `json:"currentPage"`
		PageSize    int `json:"pageSize"`
		TotalNum    int `json:"totalNum"`
		TotalPage   int `json:"totalPage"`
		Items       []struct {
			Id          string `json:"id"`
			Currency    string `json:"currency"`
			Amount      string `json:"amount"`
			Fee         string `json:"fee"`
			Balance     string `json:"balance"`
			AccountType string `json:"accountType"`
			BizType     string `json:"bizType"`
			Direction   string `json:"direction"`
			CreatedAt   int64  `json:"createdAt"`
			Context     string `json:"context"`
		} `json:"items"`
	} `json:"data"`
}

type RespLoan struct {
	Code string `json:"code"`
	Data struct {
		OrderId  string `json:"orderId"`
		Currency string `json:"currency"`
	} `json:"data"`
}

type RespRepay struct {
	Code string      `json:"code"`
	Data interface{} `json:"data"`
}

type RespUnrepayHistory struct {
	Code string `json:"code"`
	Data struct {
		CurrentPage int `json:"currentPage"`
		PageSize    int `json:"pageSize"`
		TotalNum    int `json:"totalNum"`
		TotalPage   int `json:"totalPage"`
		Items       []struct {
			TradeId         string `json:"tradeId"`
			AccruedInterest string `json:"accruedInterest"`
			Currency        string `json:"currency"`
			DailyIntRate    string `json:"dailyIntRate"`
			Liability       string `json:"liability"`
			MaturityTime    string `json:"maturityTime"`
			Principal       string `json:"principal"`
			RepaidSize      string `json:"repaidSize"`
			Term            int    `json:"term"`
			CreatedAt       string `json:"createdAt"`
		} `json:"items"`
	} `json:"data"`
}

type RespRepayHistory struct {
	Code string `json:"code"`
	Data struct {
		CurrentPage int `json:"currentPage"`
		PageSize    int `json:"pageSize"`
		TotalNum    int `json:"totalNum"`
		TotalPage   int `json:"totalPage"`
		Items       []struct {
			TradeId      string `json:"tradeId"`
			Currency     string `json:"currency"`
			DailyIntRate string `json:"dailyIntRate"`
			Interest     string `json:"interest"`
			Principal    string `json:"principal"`
			RepaidSize   string `json:"repaidSize"`
			RepayTime    string `json:"repayTime"`
			Term         int    `json:"term"`
		} `json:"items"`
	} `json:"data"`
}

type RespFutureSymbols struct {
	Code string             `json:"code"`
	Data []FutureSymbolInfo `json:"data"`
}

type RespFutureSymbol struct {
	Code string `json:"code"`
	Data struct {
		Symbol                  string      `json:"symbol"`
		RootSymbol              string      `json:"rootSymbol"`
		Type                    string      `json:"type"`
		FirstOpenDate           int64       `json:"firstOpenDate"`
		ExpireDate              interface{} `json:"expireDate"`
		SettleDate              interface{} `json:"settleDate"`
		BaseCurrency            string      `json:"baseCurrency"`
		QuoteCurrency           string      `json:"quoteCurrency"`
		SettleCurrency          string      `json:"settleCurrency"`
		MaxOrderQty             int         `json:"maxOrderQty"`
		MaxPrice                float64     `json:"maxPrice"`
		LotSize                 int         `json:"lotSize"`
		TickSize                float64     `json:"tickSize"`
		IndexPriceTickSize      float64     `json:"indexPriceTickSize"`
		Multiplier              float64     `json:"multiplier"`
		InitialMargin           float64     `json:"initialMargin"`
		MaintainMargin          float64     `json:"maintainMargin"`
		MaxRiskLimit            int         `json:"maxRiskLimit"`
		MinRiskLimit            int         `json:"minRiskLimit"`
		RiskStep                int         `json:"riskStep"`
		MakerFeeRate            float64     `json:"makerFeeRate"`
		TakerFeeRate            float64     `json:"takerFeeRate"`
		TakerFixFee             float64     `json:"takerFixFee"`
		MakerFixFee             float64     `json:"makerFixFee"`
		SettlementFee           interface{} `json:"settlementFee"`
		IsDeleverage            bool        `json:"isDeleverage"`
		IsQuanto                bool        `json:"isQuanto"`
		IsInverse               bool        `json:"isInverse"`
		MarkMethod              string      `json:"markMethod"`
		FairMethod              string      `json:"fairMethod"`
		FundingBaseSymbol       string      `json:"fundingBaseSymbol"`
		FundingQuoteSymbol      string      `json:"fundingQuoteSymbol"`
		FundingRateSymbol       string      `json:"fundingRateSymbol"`
		IndexSymbol             string      `json:"indexSymbol"`
		SettlementSymbol        string      `json:"settlementSymbol"`
		Status                  string      `json:"status"`
		FundingFeeRate          float64     `json:"fundingFeeRate"`
		PredictedFundingFeeRate float64     `json:"predictedFundingFeeRate"`
		OpenInterest            string      `json:"openInterest"`
		TurnoverOf24H           float64     `json:"turnoverOf24h"`
		VolumeOf24H             float64     `json:"volumeOf24h"`
		MarkPrice               float64     `json:"markPrice"`
		IndexPrice              float64     `json:"indexPrice"`
		LastTradePrice          float64     `json:"lastTradePrice"`
		NextFundingRateTime     int         `json:"nextFundingRateTime"`
		MaxLeverage             int         `json:"maxLeverage"`
		SourceExchanges         []string    `json:"sourceExchanges"`
		PremiumsSymbol1M        string      `json:"premiumsSymbol1M"`
		PremiumsSymbol8H        string      `json:"premiumsSymbol8H"`
		FundingBaseSymbol1M     string      `json:"fundingBaseSymbol1M"`
		FundingQuoteSymbol1M    string      `json:"fundingQuoteSymbol1M"`
		LowPrice                float64     `json:"lowPrice"`
		HighPrice               float64     `json:"highPrice"`
		PriceChgPct             float64     `json:"priceChgPct"`
		PriceChg                float64     `json:"priceChg"`
	} `json:"data"`
}

type RespFuturePositions struct {
	Code string `json:"code"`
	Data []struct {
		Id                string      `json:"id"`
		Symbol            string      `json:"symbol"`
		AutoDeposit       interface{} `json:"autoDeposit"`
		MaintMarginReq    float64     `json:"maintMarginReq"`
		RiskLimit         int         `json:"riskLimit"`
		RealLeverage      float64     `json:"realLeverage"`
		CrossMode         interface{} `json:"crossMode"`
		DelevPercentage   float64     `json:"delevPercentage"`
		OpeningTimestamp  int64       `json:"openingTimestamp"`
		CurrentTimestamp  int64       `json:"currentTimestamp"`
		CurrentQty        int         `json:"currentQty"`
		CurrentCost       float64     `json:"currentCost"`
		CurrentComm       float64     `json:"currentComm"`
		UnrealisedCost    float64     `json:"unrealisedCost"`
		RealisedGrossCost float64     `json:"realisedGrossCost"`
		RealisedCost      float64     `json:"realisedCost"`
		IsOpen            interface{} `json:"isOpen"`
		MarkPrice         float64     `json:"markPrice"`
		MarkValue         float64     `json:"markValue"`
		PosCost           float64     `json:"posCost"`
		PosCross          float64     `json:"posCross"`
		PosInit           float64     `json:"posInit"`
		PosComm           float64     `json:"posComm"`
		PosLoss           float64     `json:"posLoss"`
		PosMargin         float64     `json:"posMargin"`
		PosMaint          float64     `json:"posMaint"`
		MaintMargin       float64     `json:"maintMargin"`
		RealisedGrossPnl  float64     `json:"realisedGrossPnl"`
		RealisedPnl       float64     `json:"realisedPnl"`
		UnrealisedPnl     float64     `json:"unrealisedPnl"`
		UnrealisedPnlPcnt float64     `json:"unrealisedPnlPcnt"`
		UnrealisedRoePcnt float64     `json:"unrealisedRoePcnt"`
		AvgEntryPrice     float64     `json:"avgEntryPrice"`
		LiquidationPrice  float64     `json:"liquidationPrice"`
		BankruptPrice     float64     `json:"bankruptPrice"`
		SettleCurrency    string      `json:"settleCurrency"`
		IsInverse         bool        `json:"isInverse"`
		MaintainMargin    float64     `json:"maintainMargin"`
	} `json:"data"`
}
type FutureSymbolInfo struct {
	Symbol                  string      `json:"symbol"`
	RootSymbol              string      `json:"rootSymbol"`
	Type                    string      `json:"type"`
	FirstOpenDate           int64       `json:"firstOpenDate"`
	ExpireDate              interface{} `json:"expireDate"`
	SettleDate              interface{} `json:"settleDate"`
	BaseCurrency            string      `json:"baseCurrency"`
	QuoteCurrency           string      `json:"quoteCurrency"`
	SettleCurrency          string      `json:"settleCurrency"`
	MaxOrderQty             int         `json:"maxOrderQty"`
	MaxPrice                float64     `json:"maxPrice"`
	LotSize                 int         `json:"lotSize"`
	TickSize                float64     `json:"tickSize"`
	IndexPriceTickSize      float64     `json:"indexPriceTickSize"`
	Multiplier              float64     `json:"multiplier"`
	InitialMargin           float64     `json:"initialMargin"`
	MaintainMargin          float64     `json:"maintainMargin"`
	MaxRiskLimit            int         `json:"maxRiskLimit"`
	MinRiskLimit            int         `json:"minRiskLimit"`
	RiskStep                int         `json:"riskStep"`
	MakerFeeRate            float64     `json:"makerFeeRate"`
	TakerFeeRate            float64     `json:"takerFeeRate"`
	TakerFixFee             float64     `json:"takerFixFee"`
	MakerFixFee             float64     `json:"makerFixFee"`
	SettlementFee           interface{} `json:"settlementFee"`
	IsDeleverage            bool        `json:"isDeleverage"`
	IsQuanto                bool        `json:"isQuanto"`
	IsInverse               bool        `json:"isInverse"`
	MarkMethod              string      `json:"markMethod"`
	FairMethod              string      `json:"fairMethod"`
	FundingBaseSymbol       string      `json:"fundingBaseSymbol"`
	FundingQuoteSymbol      string      `json:"fundingQuoteSymbol"`
	FundingRateSymbol       string      `json:"fundingRateSymbol"`
	IndexSymbol             string      `json:"indexSymbol"`
	SettlementSymbol        string      `json:"settlementSymbol"`
	Status                  string      `json:"status"`
	FundingFeeRate          float64     `json:"fundingFeeRate"`
	PredictedFundingFeeRate float64     `json:"predictedFundingFeeRate"`
	OpenInterest            string      `json:"openInterest"`
	TurnoverOf24H           float64     `json:"turnoverOf24h"`
	VolumeOf24H             float64     `json:"volumeOf24h"`
	MarkPrice               float64     `json:"markPrice"`
	IndexPrice              float64     `json:"indexPrice"`
	LastTradePrice          float64     `json:"lastTradePrice"`
	NextFundingRateTime     int         `json:"nextFundingRateTime"`
	MaxLeverage             int         `json:"maxLeverage"`
	SourceExchanges         []string    `json:"sourceExchanges"`
	PremiumsSymbol1M        string      `json:"premiumsSymbol1M"`
	PremiumsSymbol8H        string      `json:"premiumsSymbol8H"`
	FundingBaseSymbol1M     string      `json:"fundingBaseSymbol1M"`
	FundingQuoteSymbol1M    string      `json:"fundingQuoteSymbol1M"`
	LowPrice                float64     `json:"lowPrice"`
	HighPrice               float64     `json:"highPrice"`
	PriceChgPct             float64     `json:"priceChgPct"`
	PriceChg                float64     `json:"priceChg"`
}
type RespFuturePlaceOrder struct {
	Code string `json:"code"`
	Data struct {
		OrderId string `json:"orderId"`
	} `json:"data"`
}
