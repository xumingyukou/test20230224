package spot_api

import "time"

const (
	CONTANTTYPE = "application/json"
)

type RespError struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

type RespSymbols struct {
	RetCode int          `json:"ret_code"`
	RetMsg  string       `json:"ret_msg"`
	ExtCode interface{}  `json:"ext_code"`
	ExtInfo interface{}  `json:"ext_info"`
	Result  []SymbolInfo `json:"result"`
}

type RespWithdrawFee struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	Result  struct {
		Rows []struct {
			Name         string `json:"name"`
			Coin         string `json:"coin"`
			RemainAmount string `json:"remain_amount"`
			Chains       []struct {
				ChainType    string `json:"chain_type"`
				Confirmation string `json:"confirmation"`
				WithdrawFee  string `json:"withdraw_fee"`
				DepositMin   string `json:"deposit_min"`
				WithdrawMin  string `json:"withdraw_min"`
				Chain        string `json:"chain"`
			} `json:"chains"`
		} `json:"rows"`
	} `json:"result"`
	ExtInfo          interface{} `json:"ext_info"`
	TimeNow          int64       `json:"time_now"`
	RateLimitStatus  int         `json:"rate_limit_status"`
	RateLimitResetMs int64       `json:"rate_limit_reset_ms"`
	RateLimit        int         `json:"rate_limit"`
}

type SymbolInfo struct {
	Name              string `json:"name"`
	Alias             string `json:"alias"`
	BaseCurrency      string `json:"baseCurrency"`
	QuoteCurrency     string `json:"quoteCurrency"`
	BasePrecision     string `json:"basePrecision"`
	QuotePrecision    string `json:"quotePrecision"`
	MinTradeQuantity  string `json:"minTradeQuantity"`
	MinTradeAmount    string `json:"minTradeAmount"`
	MinPricePrecision string `json:"minPricePrecision"`
	MaxTradeQuantity  string `json:"maxTradeQuantity"`
	MaxTradeAmount    string `json:"maxTradeAmount"`
	Category          int    `json:"category"`
	Innovation        bool   `json:"innovation"`
	ShowStatus        bool   `json:"showStatus"`
}

type RespOrderBook struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		Time int64      `json:"time"`
		Bids [][]string `json:"bids"`
		Asks [][]string `json:"asks"`
	} `json:"result"`
	RetExtInfo struct {
	} `json:"retExtInfo"`
	Time int64 `json:"time"`
}

type RespTime struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		TimeSecond string `json:"timeSecond"`
		TimeNano   string `json:"timeNano"`
	} `json:"result"`
	RetExtInfo struct {
	} `json:"retExtInfo"`
	Time int64 `json:"time"`
}

type RespFundFee struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		List []struct {
			FundFee          string `json:"fundFee"`
			FundFeeTime      int64  `json:"fundFeeTime"`
			LtCode           string `json:"ltCode"`
			LtName           string `json:"ltName"`
			ManageFeeRate    string `json:"manageFeeRate"`
			ManageFeeTime    int64  `json:"manageFeeTime"`
			MaxPurchase      string `json:"maxPurchase"`
			MaxPurchaseDaily string `json:"maxPurchaseDaily"`
			MaxRedeem        string `json:"maxRedeem"`
			MaxRedeemDaily   string `json:"maxRedeemDaily"`
			MinPurchase      string `json:"minPurchase"`
			MinRedeem        string `json:"minRedeem"`
			NetValue         string `json:"netValue"`
			PurchaseFeeRate  string `json:"purchaseFeeRate"`
			RedeemFeeRate    string `json:"redeemFeeRate"`
			Status           string `json:"status"`
			Total            string `json:"total"`
			Value            string `json:"value"`
		} `json:"list"`
	} `json:"result"`
	RetExtInfo interface{} `json:"retExtInfo"`
	Time       int64       `json:"time"`
}

type RespAccountBalance struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		Balances []struct {
			Coin   string `json:"coin"`
			CoinId string `json:"coinId"`
			Total  string `json:"total"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
		} `json:"balances"`
	} `json:"result"`
	RetExtMap struct {
	} `json:"retExtMap"`
	RetExtInfo struct {
	} `json:"retExtInfo"`
	Time int64 `json:"time"`
}

type RespAccount struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		Balances []struct {
			Coin   string `json:"coin"`
			CoinId string `json:"coinId"`
			Total  string `json:"total"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
		} `json:"balances"`
	} `json:"result"`
	RetExtMap struct {
	} `json:"retExtMap"`
	RetExtInfo struct {
	} `json:"retExtInfo"`
	Time int64 `json:"time"`
}

type RespPlaceOrder struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		OrderId       string `json:"orderId"`
		OrderLinkId   string `json:"orderLinkId"`
		Symbol        string `json:"symbol"`
		CreateTime    string `json:"createTime"`
		OrderPrice    string `json:"orderPrice"`
		OrderQty      string `json:"orderQty"`
		OrderType     string `json:"orderType"`
		Side          string `json:"side"`
		Status        string `json:"status"`
		TimeInForce   string `json:"timeInForce"`
		AccountId     string `json:"accountId"`
		OrderCategory int    `json:"orderCategory"`
		TriggerPrice  string `json:"triggerPrice"`
	} `json:"result"`
	RetExtMap struct {
	} `json:"retExtMap"`
	RetExtInfo interface{} `json:"retExtInfo"`
	Time       int64       `json:"time"`
}

type RespCancleOrder struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		OrderId     string `json:"orderId"`
		OrderLinkId string `json:"orderLinkId"`
		Symbol      string `json:"symbol"`
		Status      string `json:"status"`
		AccountId   string `json:"accountId"`
		CreateTime  string `json:"createTime"`
		OrderPrice  string `json:"orderPrice"`
		OrderQty    string `json:"orderQty"`
		ExecQty     string `json:"execQty"`
		TimeInForce string `json:"timeInForce"`
		OrderType   string `json:"orderType"`
		Side        string `json:"side"`
	} `json:"result"`
	RetExtMap struct {
	} `json:"retExtMap"`
	RetExtInfo struct {
	} `json:"retExtInfo"`
	Time int64 `json:"time"`
}

type RespOrderInfo struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		AccountId           string `json:"accountId"`
		Symbol              string `json:"symbol"`
		OrderLinkId         string `json:"orderLinkId"`
		OrderId             string `json:"orderId"`
		OrderPrice          string `json:"orderPrice"`
		OrderQty            string `json:"orderQty"`
		ExecQty             string `json:"execQty"`
		CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
		AvgPrice            string `json:"avgPrice"`
		Status              string `json:"status"`
		TimeInForce         string `json:"timeInForce"`
		OrderType           string `json:"orderType"`
		Side                string `json:"side"`
		StopPrice           string `json:"stopPrice"`
		IcebergQty          string `json:"icebergQty"`
		CreateTime          string `json:"createTime"`
		UpdateTime          string `json:"updateTime"`
		IsWorking           string `json:"isWorking"`
		Locked              string `json:"locked"`
	} `json:"result"`
	RetExtMap struct {
	} `json:"retExtMap"`
	RetExtInfo struct {
	} `json:"retExtInfo"`
	Time int64 `json:"time"`
}

type RespOrderHistory struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		List []struct {
			AccountId           string `json:"accountId"`
			Symbol              string `json:"symbol"`
			OrderLinkId         string `json:"orderLinkId"`
			OrderId             string `json:"orderId"`
			OrderPrice          string `json:"orderPrice"`
			OrderQty            string `json:"orderQty"`
			ExecQty             string `json:"execQty"`
			CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
			AvgPrice            string `json:"avgPrice"`
			Status              string `json:"status"`
			TimeInForce         string `json:"timeInForce"`
			OrderType           string `json:"orderType"`
			Side                string `json:"side"`
			StopPrice           string `json:"stopPrice"`
			IcebergQty          string `json:"icebergQty"`
			CreateTime          int64  `json:"createTime"`
			UpdateTime          int64  `json:"updateTime"`
			IsWorking           string `json:"isWorking"`
		} `json:"list"`
	} `json:"result"`
	RetExtInfo struct {
	} `json:"retExtInfo"`
	Time int64 `json:"time"`
}

type RespWithdraw struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	Result  struct {
		Id string `json:"id"`
	} `json:"result"`
	ExtInfo          interface{} `json:"ext_info"`
	TimeNow          int64       `json:"time_now"`
	RateLimitStatus  int         `json:"rate_limit_status"`
	RateLimitResetMs int64       `json:"rate_limit_reset_ms"`
	RateLimit        int         `json:"rate_limit"`
}

type RespWithdrawHistory struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	Result  struct {
		Rows []struct {
			Coin        string `json:"coin"`
			Chain       string `json:"chain"`
			Amount      string `json:"amount"`
			TxId        string `json:"tx_id"`
			Status      string `json:"status"`
			ToAddress   string `json:"to_address,omitempty"`
			Tag         string `json:"tag"`
			WithdrawFee string `json:"withdraw_fee"`
			CreateTime  string `json:"create_time"`
			UpdateTime  string `json:"update_time"`
			WithdrawId  string `json:"withdraw_id"`
			ToAddress1  string `json:"toAddress,omitempty"`
		} `json:"rows"`
		Cursor string `json:"cursor"`
	} `json:"result"`
	ExtInfo          interface{} `json:"ext_info"`
	TimeNow          int64       `json:"time_now"`
	RateLimitStatus  int         `json:"rate_limit_status"`
	RateLimitResetMs int64       `json:"rate_limit_reset_ms"`
	RateLimit        int         `json:"rate_limit"`
}

type RespTransfer struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	Result  struct {
		TransferId string `json:"transfer_id"`
	} `json:"result"`
	ExtInfo          interface{} `json:"ext_info"`
	TimeNow          int64       `json:"time_now"`
	RateLimitStatus  int         `json:"rate_limit_status"`
	RateLimitResetMs int64       `json:"rate_limit_reset_ms"`
	RateLimit        int         `json:"rate_limit"`
}

type RespTransferM2S struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	Result  struct {
		TransferId string `json:"transfer_id"`
	} `json:"result"`
	ExtInfo          interface{} `json:"ext_info"`
	TimeNow          int64       `json:"time_now"`
	RateLimitStatus  int         `json:"rate_limit_status"`
	RateLimitResetMs int64       `json:"rate_limit_reset_ms"`
	RateLimit        int         `json:"rate_limit"`
}

type RespAllTransfer struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	Result  struct {
		TransferId string `json:"transfer_id"`
	} `json:"result"`
	ExtInfo          interface{} `json:"ext_info"`
	TimeNow          int64       `json:"time_now"`
	RateLimitStatus  int         `json:"rate_limit_status"`
	RateLimitResetMs int64       `json:"rate_limit_reset_ms"`
	RateLimit        int         `json:"rate_limit"`
}

type RespTransferHistory struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	Result  struct {
		List []struct {
			TransferId      string `json:"transfer_id"`
			Coin            string `json:"coin"`
			Amount          string `json:"amount"`
			FromAccountType string `json:"from_account_type"`
			ToAccountType   string `json:"to_account_type"`
			Timestamp       string `json:"timestamp"`
			Status          string `json:"status"`
		} `json:"list"`
		Cursor string `json:"cursor"`
	} `json:"result"`
	ExtInfo          interface{} `json:"ext_info"`
	TimeNow          int64       `json:"time_now"`
	RateLimitStatus  int         `json:"rate_limit_status"`
	RateLimitResetMs int64       `json:"rate_limit_reset_ms"`
	RateLimit        int         `json:"rate_limit"`
}

type RespTransferHistoryS2M struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	Result  struct {
		List []struct {
			TransferId string `json:"transfer_id"`
			Coin       string `json:"coin"`
			Amount     string `json:"amount"`
			UserId     int    `json:"user_id"`
			SubUserId  int    `json:"sub_user_id"`
			Timestamp  string `json:"timestamp"`
			Status     string `json:"status"`
			Type       string `json:"type"`
		} `json:"list"`
		Cursor string `json:"cursor"`
	} `json:"result"`
	ExtInfo          interface{} `json:"ext_info"`
	TimeNow          int64       `json:"time_now"`
	RateLimitStatus  int         `json:"rate_limit_status"`
	RateLimitResetMs int64       `json:"rate_limit_reset_ms"`
	RateLimit        int         `json:"rate_limit"`
}

type RespRecordHistory struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	Result  struct {
		Rows []struct {
			Coin          string `json:"coin"`
			Chain         string `json:"chain"`
			Amount        string `json:"amount"`
			TxId          string `json:"tx_id"`
			Status        int    `json:"status"`
			ToAddress     string `json:"to_address"`
			Tag           string `json:"tag"`
			DepositFee    string `json:"deposit_fee"`
			SuccessAt     string `json:"success_at"`
			Confirmations string `json:"confirmations"`
			TxIndex       string `json:"tx_index"`
			BlockHash     string `json:"block_hash"`
		} `json:"rows"`
		Cursor string `json:"cursor"`
	} `json:"result"`
	ExtInfo          interface{} `json:"ext_info"`
	TimeNow          int64       `json:"time_now"`
	RateLimitStatus  int         `json:"rate_limit_status"`
	RateLimitResetMs int64       `json:"rate_limit_reset_ms"`
	RateLimit        int         `json:"rate_limit"`
}

type RespLoan struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		TransactId string `json:"transactId"`
	} `json:"result"`
	RetExtInfo interface{} `json:"retExtInfo"`
	Time       int64       `json:"time"`
}

type RespLoanHistory struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		List []struct {
			AccountId       string `json:"accountId"`
			Coin            string `json:"coin"`
			CreatedTime     int64  `json:"createdTime"`
			Id              string `json:"id"`
			InterestAmount  string `json:"interestAmount"`
			InterestBalance string `json:"interestBalance"`
			LoanAmount      string `json:"loanAmount"`
			LoanBalance     string `json:"loanBalance"`
			Status          int    `json:"status"`
			Type            int    `json:"type"`
		} `json:"list"`
	} `json:"result"`
	RetExtInfo interface{} `json:"retExtInfo"`
	Time       int64       `json:"time"`
}

type RespRepay struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		RepayId string `json:"repayId"`
	} `json:"result"`
	RetExtInfo interface{} `json:"retExtInfo"`
	Time       int64       `json:"time"`
}

type RespRepayHistory struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		List []struct {
			AccountId          string `json:"accountId"`
			Coin               string `json:"coin"`
			RepaidAmount       string `json:"repaidAmount"`
			RepayId            string `json:"repayId"`
			RepayMarginOrderId string `json:"repayMarginOrderId"`
			RepayTime          string `json:"repayTime"`
			TransactIds        []struct {
				RepaidAmount       string `json:"repaidAmount"`
				RepaidInterest     string `json:"repaidInterest"`
				RepaidPrincipal    string `json:"repaidPrincipal"`
				RepaidSerialNumber string `json:"repaidSerialNumber"`
				TransactId         string `json:"transactId"`
			} `json:"transactIds"`
		} `json:"list"`
	} `json:"result"`
	RetExtInfo interface{} `json:"retExtInfo"`
	Time       int64       `json:"time"`
}

type RespFutureSymbols struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	Result  []struct {
		Name            string `json:"name"`
		Alias           string `json:"alias"`
		Status          string `json:"status"`
		BaseCurrency    string `json:"base_currency"`
		QuoteCurrency   string `json:"quote_currency"`
		PriceScale      int    `json:"price_scale"`
		TakerFee        string `json:"taker_fee"`
		MakerFee        string `json:"maker_fee"`
		FundingInterval int    `json:"funding_interval"`
		LeverageFilter  struct {
			MinLeverage  int    `json:"min_leverage"`
			MaxLeverage  int    `json:"max_leverage"`
			LeverageStep string `json:"leverage_step"`
		} `json:"leverage_filter"`
		PriceFilter struct {
			MinPrice string `json:"min_price"`
			MaxPrice string `json:"max_price"`
			TickSize string `json:"tick_size"`
		} `json:"price_filter"`
		LotSizeFilter struct {
			MaxTradingQty         float64 `json:"max_trading_qty"`
			MinTradingQty         float64 `json:"min_trading_qty"`
			QtyStep               float64 `json:"qty_step"`
			PostOnlyMaxTradingQty string  `json:"post_only_max_trading_qty"`
		} `json:"lot_size_filter"`
	} `json:"result"`
	ExtCode string `json:"ext_code"`
	ExtInfo string `json:"ext_info"`
	TimeNow string `json:"time_now"`
}

type RespFutureOrderbook struct {
	RetCode int         `json:"ret_code"`
	RetMsg  string      `json:"ret_msg"`
	Result  []OrderItem `json:"result"`
	ExtCode string      `json:"ext_code"`
	ExtInfo string      `json:"ext_info"`
	TimeNow string      `json:"time_now"`
}

type OrderItem struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Side   string `json:"side"`
	Size   int    `json:"size"`
}

type RespMarketPrice struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	ExtInfo string `json:"ext_info"`
	Result  []struct {
		Symbol  string  `json:"symbol"`
		Period  string  `json:"period"`
		StartAt int     `json:"start_at"`
		Open    float64 `json:"open"`
		High    float64 `json:"high"`
		Low     float64 `json:"low"`
		Close   float64 `json:"close"`
	} `json:"result"`
	TimeNow string `json:"time_now"`
}
type RespFuturePosition struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	ExtInfo string `json:"ext_info"`
	Result  []struct {
		UserId              int     `json:"user_id"`
		Symbol              string  `json:"symbol"`
		Side                string  `json:"side"`
		Size                int     `json:"size"`
		PositionValue       float64 `json:"position_value"`
		EntryPrice          float64 `json:"entry_price"`
		LiqPrice            float64 `json:"liq_price"`
		BustPrice           float64 `json:"bust_price"`
		Leverage            int     `json:"leverage"`
		AutoAddMargin       int     `json:"auto_add_margin"`
		IsIsolated          bool    `json:"is_isolated"`
		PositionMargin      float64 `json:"position_margin"`
		OccClosingFee       float64 `json:"occ_closing_fee"`
		RealisedPnl         int     `json:"realised_pnl"`
		CumRealisedPnl      float64 `json:"cum_realised_pnl"`
		FreeQty             int     `json:"free_qty"`
		TpSlMode            string  `json:"tp_sl_mode"`
		UnrealisedPnl       float64 `json:"unrealised_pnl"`
		DeleverageIndicator int     `json:"deleverage_indicator"`
		RiskId              int     `json:"risk_id"`
		StopLoss            float64 `json:"stop_loss"`
		TakeProfit          float64 `json:"take_profit"`
		TrailingStop        int     `json:"trailing_stop"`
		PositionIdx         int     `json:"position_idx"`
		Mode                string  `json:"mode"`
		TpTriggerBy         int     `json:"tp_trigger_by,omitempty"`
		SlTriggerBy         int     `json:"sl_trigger_by,omitempty"`
	} `json:"result"`
	TimeNow          string `json:"time_now"`
	RateLimitStatus  int    `json:"rate_limit_status"`
	RateLimitResetMs int64  `json:"rate_limit_reset_ms"`
	RateLimit        int    `json:"rate_limit"`
}

type RespFuturePlaceOrder struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	ExtInfo string `json:"ext_info"`
	Result  struct {
		OrderId        string    `json:"order_id"`
		UserId         int       `json:"user_id"`
		Symbol         string    `json:"symbol"`
		Side           string    `json:"side"`
		OrderType      string    `json:"order_type"`
		Price          float64   `json:"price"`
		Qty            int       `json:"qty"`
		TimeInForce    string    `json:"time_in_force"`
		OrderStatus    string    `json:"order_status"`
		LastExecPrice  int       `json:"last_exec_price"`
		CumExecQty     int       `json:"cum_exec_qty"`
		CumExecValue   int       `json:"cum_exec_value"`
		CumExecFee     int       `json:"cum_exec_fee"`
		ReduceOnly     bool      `json:"reduce_only"`
		CloseOnTrigger bool      `json:"close_on_trigger"`
		OrderLinkId    string    `json:"order_link_id"`
		CreatedTime    time.Time `json:"created_time"`
		UpdatedTime    time.Time `json:"updated_time"`
		TakeProfit     float64   `json:"take_profit"`
		StopLoss       float64   `json:"stop_loss"`
		TpTriggerBy    string    `json:"tp_trigger_by"`
		SlTriggerBy    string    `json:"sl_trigger_by"`
		PositionIdx    int       `json:"position_idx"`
	} `json:"result"`
	TimeNow          string `json:"time_now"`
	RateLimitStatus  int    `json:"rate_limit_status"`
	RateLimitResetMs int64  `json:"rate_limit_reset_ms"`
	RateLimit        int    `json:"rate_limit"`
}

type RespFutureCoinDeliveryPlaceOrder struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	ExtInfo string `json:"ext_info"`
	Result  struct {
		UserId        int       `json:"user_id"`
		OrderId       string    `json:"order_id"`
		Symbol        string    `json:"symbol"`
		Side          string    `json:"side"`
		OrderType     string    `json:"order_type"`
		Price         int       `json:"price"`
		Qty           int       `json:"qty"`
		TimeInForce   string    `json:"time_in_force"`
		OrderStatus   string    `json:"order_status"`
		LastExecTime  int       `json:"last_exec_time"`
		LastExecPrice int       `json:"last_exec_price"`
		LeavesQty     int       `json:"leaves_qty"`
		CumExecQty    int       `json:"cum_exec_qty"`
		CumExecValue  int       `json:"cum_exec_value"`
		CumExecFee    int       `json:"cum_exec_fee"`
		RejectReason  string    `json:"reject_reason"`
		OrderLinkId   string    `json:"order_link_id"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		TakeProfit    string    `json:"take_profit"`
		StopLoss      string    `json:"stop_loss"`
		TpTriggerBy   string    `json:"tp_trigger_by"`
		SlTriggerBy   string    `json:"sl_trigger_by"`
	} `json:"result"`
	TimeNow          string `json:"time_now"`
	RateLimitStatus  int    `json:"rate_limit_status"`
	RateLimitResetMs int64  `json:"rate_limit_reset_ms"`
	RateLimit        int    `json:"rate_limit"`
}

type RespFutureCoinOrder struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	ExtInfo string `json:"ext_info"`
	Result  struct {
		UserId        int       `json:"user_id"`
		OrderId       string    `json:"order_id"`
		Symbol        string    `json:"symbol"`
		Side          string    `json:"side"`
		OrderType     string    `json:"order_type"`
		Price         int       `json:"price"`
		Qty           int       `json:"qty"`
		TimeInForce   string    `json:"time_in_force"`
		OrderStatus   string    `json:"order_status"`
		LastExecTime  int       `json:"last_exec_time"`
		LastExecPrice int       `json:"last_exec_price"`
		LeavesQty     int       `json:"leaves_qty"`
		CumExecQty    int       `json:"cum_exec_qty"`
		CumExecValue  int       `json:"cum_exec_value"`
		CumExecFee    int       `json:"cum_exec_fee"`
		RejectReason  string    `json:"reject_reason"`
		OrderLinkId   string    `json:"order_link_id"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		TakeProfit    string    `json:"take_profit"`
		StopLoss      string    `json:"stop_loss"`
		TpTriggerBy   string    `json:"tp_trigger_by"`
		SlTriggerBy   string    `json:"sl_trigger_by"`
	} `json:"result"`
	TimeNow          string `json:"time_now"`
	RateLimitStatus  int    `json:"rate_limit_status"`
	RateLimitResetMs int64  `json:"rate_limit_reset_ms"`
	RateLimit        int    `json:"rate_limit"`
}
