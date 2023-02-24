package u_api

import (
	"strings"

	"github.com/warmplanet/proto/go/common"
)

type FundingFeeRateResponse struct {
	Status string `json:"status"`
	Data   []struct {
		EstimatedRate   string `json:"estimated_rate"`
		FundingRate     string `json:"funding_rate"`
		ContractCode    string `json:"contract_code"`
		Symbol          string `json:"symbol"`
		FeeAsset        string `json:"fee_asset"`
		FundingTime     string `json:"funding_time"`
		NextFundingTime string `json:"next_funding_time"`
	} `json:"data"`
	TS int64 `json:"ts"`
}

type ContractInfo struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol            string  `json:"symbol"`
		ContractCode      string  `json:"contract_code"`
		ContractSize      float64 `json:"contract_size"`
		PriceTick         float64 `json:"price_tick"`
		DeliveryDate      string  `json:"delivery_date"`
		DeliveryTime      string  `json:"delivery_time"`
		CreateDate        string  `json:"create_date"`
		ContractStatus    int     `json:"contract_status"`
		SettlementTime    string  `json:"settlement_time"`
		SupportMarginMode string  `json:"support_margin_mode"`
		BusinessType      string  `json:"business_type"`
		Pair              string  `json:"pair"`
		ContractType      string  `json:"contract_type"`
	} `json:"data"`
	TS int64 `json:"ts"`
}

type ReqPostSwapCrossPositionLimit struct {
	ContractCode string `json:"contract_code"`
	Pair         string `json:"pair"`
	ContractType string `json:"contract_type"`
	BusinessType string `json:"business_type"`
}

type ResPostSwapCrossPositionLimit struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol         string  `json:"symbol"`
		ContractCode   string  `json:"contract_code"`
		MarginMode     string  `json:"margin_mode"`
		BuyLimit       int     `json:"buy_limit"`
		SellLimit      int     `json:"sell_limit"`
		BusinessType   string  `json:"business_type"`
		ContractType   string  `json:"contract_type"`
		Pair           string  `json:"pair"`
		LeverRate      int     `json:"lever_rate"`
		BuyLimitValue  float64 `json:"buy_limit_value"`
		SellLimitValue float64 `json:"sell_limit_value"`
		MarkPrice      float64 `json:"mark_price"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostSwapFee struct {
	ContractCode string `json:"contract_code"`
	Pair         string `json:"pair"`
	ContractType string `json:"contract_type"`
	BusinessType string `json:"business_type"`
}

type ResPostSwapFee struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol        string `json:"symbol"`
		ContractCode  string `json:"contract_code"`
		OpenMakerFee  string `json:"open_maker_fee"`
		OpenTakerFee  string `json:"open_taker_fee"`
		CloseMakerFee string `json:"close_maker_fee"`
		CloseTakerFee string `json:"close_taker_fee"`
		FeeAsset      string `json:"fee_asset"`
		DeliveryFee   string `json:"delivery_fee"`
		BusinessType  string `json:"business_type"`
		ContractType  string `json:"contract_type"`
		Pair          string `json:"pair"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostSwapBalanceValuation struct {
	ValuationAsset string `json:"valuation_asset"`
}

type ResPostSwapBalanceValuation struct {
	Status string `json:"status"`
	Data   []struct {
		ValuationAsset string `json:"valuation_asset"`
		Balance        string `json:"balance"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostSwapCrossOrder struct {
	ContractCode     string  `json:"contract_code"`
	Pair             string  `json:"pair"`
	ContractType     string  `json:"contract_type"`
	ReduceOnly       int     `json:"reduce_only"`
	ClientOrderID    string  `json:"client_order_id"`
	Direction        string  `json:"direction"`
	Offset           string  `json:"offset"`
	Price            string  `json:"price"`
	LeverRate        int     `json:"lever_rate"`
	Volume           int64   `json:"volume"`
	OrderPriceType   string  `json:"order_price_type"`
	TpTriggerPrice   float64 `json:"tp_trigger_price"`
	TpOrderPrice     float64 `json:"tp_order_price"`
	TpOrderPriceType string  `json:"tp_order_price_type"`
	//SlTriggerPrice   string  `json:"slt_rigger_price"`
	SlOrderPrice     string `json:"sl_order_price"`
	SlOrderPriceType string `json:"sl_order_price_type"`
}

type ResPostSwapCrossOrder struct {
	Status string `json:"status"`
	Data   struct {
		OrderID       int64  `json:"order_id"`
		OrderIDStr    string `json:"order_id_str"`
		ClientOrderID int64  `json:"client_order_id"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostSwapCrossCancelOrder struct {
	OrderID       string `json:"order_id"`
	ClientOrderID string `json:"client_order_id"`
	Contract_Code string `json:"contract_code"`
	Pair          string `json:"pair"`
	ContractType  string `json:"contract_type"`
}

type ResPostSwapCrossCancelOrder struct {
	Status string `json:"status"`
	Data   struct {
		Errors []struct {
			OrderID string `json:"order_id"`
			ErrCode int    `json:"err_code"`
			ErrMsg  string `json:"err_msg"`
		} `json:"errors"`
		Successes string `json:"successes"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostSwapCrossOrderInfo struct {
	OrderID       string `json:"order_id"`
	ClientOrderID string `json:"client_order_id"`
	ContractCode  string `json:"contract_code"`
	Pair          string `json:"pair"`
}

type ResPostSwapCrossOrderInfo struct {
	Status string `json:"status"`
	Data   []struct {
		BusinessType    string      `json:"business_type"`
		ContractType    string      `json:"contract_type"`
		Pair            string      `json:"pair"`
		Symbol          string      `json:"symbol"`
		ContractCode    string      `json:"contract_code"`
		Volume          float64     `json:"volume"`
		Price           float64     `json:"price"`
		OrderPriceType  string      `json:"order_price_type"`
		OrderType       int         `json:"order_type"`
		Direction       string      `json:"direction"`
		Offset          string      `json:"offset"`
		LeverRate       int         `json:"lever_rate"`
		OrderID         int64       `json:"order_id"`
		ClientOrderID   interface{} `json:"client_order_id"`
		CreatedAt       int64       `json:"created_at"`
		TradeVolume     int         `json:"trade_volume"`
		TradeTurnover   int         `json:"trade_turnover"`
		Fee             int         `json:"fee"`
		TradeAvgPrice   interface{} `json:"trade_avg_price"`
		MarginFrozen    float64     `json:"margin_frozen"`
		Profit          int         `json:"profit"`
		Status          int         `json:"status"`
		OrderSource     string      `json:"order_source"`
		OrderIDStr      string      `json:"order_id_str"`
		FeeAsset        string      `json:"fee_asset"`
		LiquidationType string      `json:"liquidation_type"`
		CanceledAt      int         `json:"canceled_at"`
		MarginAsset     string      `json:"margin_asset"`
		MarginAccount   string      `json:"margin_account"`
		MarginMode      string      `json:"margin_mode"`
		IsTpsl          int         `json:"is_tpsl"`
		RealProfit      int         `json:"real_profit"`
		ReduceOnly      int         `json:"reduce_only"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostSwapCrossMatchresults struct {
	Contract  string `json:"contract"`
	TradeType int    `json:"trade_type"`
	Pair      string `json:"pair"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Direct    string `json:"direct"`
	FromID    int    `json:"from_id"`
}

type ResPostSwapCrossMatchresults struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		QueryID          int     `json:"query_id"`
		ContractType     string  `json:"contract_type"`
		Pair             string  `json:"pair"`
		BusinessType     string  `json:"business_type"`
		MatchID          int     `json:"match_id"`
		OrderID          int64   `json:"order_id"`
		Symbol           string  `json:"symbol"`
		ContractCode     string  `json:"contract_code"`
		Direction        string  `json:"direction"`
		Offset           string  `json:"offset"`
		TradeVolume      float64 `json:"trade_volume"`
		TradePrice       float64 `json:"trade_price"`
		TradeTurnover    float64 `json:"trade_turnover"`
		TradeFee         float64 `json:"trade_fee"`
		OffsetProfitloss int     `json:"offset_profitloss"`
		CreateDate       int64   `json:"create_date"`
		Role             string  `json:"role"`
		OrderSource      string  `json:"order_source"`
		OrderIDStr       string  `json:"order_id_str"`
		ID               string  `json:"id"`
		FeeAsset         string  `json:"fee_asset"`
		MarginMode       string  `json:"margin_mode"`
		MarginAccount    string  `json:"margin_account"`
		RealProfit       int     `json:"real_profit"`
		ReduceOnly       int     `json:"reduce_only"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostSwapOpenOrders struct {
	ContractCode string `json:"contract_code"`
	PageIndex    int    `json:"page_index"`
	PageSize     int    `json:"page_size"`
	SortBy       string `json:"sort_by"`
	TradeType    int    `json:"trade_type"`
}

type ResPostSwapOpenOrders struct {
	Status string `json:"status"`
	Data   struct {
		Orders []struct {
			Symbol          string      `json:"symbol"`
			ContractCode    string      `json:"contract_code"`
			Volume          int         `json:"volume"`
			Price           int         `json:"price"`
			OrderPriceType  string      `json:"order_price_type"`
			OrderType       int         `json:"order_type"`
			Direction       string      `json:"direction"`
			Offset          string      `json:"offset"`
			LeverRate       int         `json:"lever_rate"`
			OrderID         int64       `json:"order_id"`
			ClientOrderID   int64       `json:"client_order_id"`
			CreatedAt       int64       `json:"created_at"`
			TradeVolume     int         `json:"trade_volume"`
			TradeTurnover   int         `json:"trade_turnover"`
			Fee             int         `json:"fee"`
			TradeAvgPrice   interface{} `json:"trade_avg_price"`
			MarginFrozen    float64     `json:"margin_frozen"`
			Profit          int         `json:"profit"`
			Status          int         `json:"status"`
			OrderSource     string      `json:"order_source"`
			OrderIDStr      string      `json:"order_id_str"`
			FeeAsset        string      `json:"fee_asset"`
			LiquidationType interface{} `json:"liquidation_type"`
			CanceledAt      interface{} `json:"canceled_at"`
			MarginAsset     string      `json:"margin_asset"`
			MarginMode      string      `json:"margin_mode"`
			MarginAccount   string      `json:"margin_account"`
			IsTpsl          int         `json:"is_tpsl"`
			UpdateTime      int64       `json:"update_time"`
			RealProfit      int         `json:"real_profit"`
			ReduceOnly      int         `json:"reduce_only"`
		} `json:"orders"`
		TotalPage   int `json:"total_page"`
		CurrentPage int `json:"current_page"`
		TotalSize   int `json:"total_size"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostSwapCrossAccountInfo struct {
	MarginAccount string `json:"margin_account"`
}

type ResPostSwapCrossAccountInfo struct {
	Status string `json:"status"`
	Data   []struct {
		FuturesContractDetail []struct {
			Symbol           string      `json:"symbol"`
			ContractCode     string      `json:"contract_code"`
			MarginPosition   float64     `json:"margin_position"`
			MarginFrozen     float64     `json:"margin_frozen"`
			MarginAvailable  float64     `json:"margin_available"`
			ProfitUnreal     int         `json:"profit_unreal"`
			LiquidationPrice interface{} `json:"liquidation_price"`
			LeverRate        int         `json:"lever_rate"`
			AdjustFactor     float64     `json:"adjust_factor"`
			ContractType     string      `json:"contract_type"`
			Pair             string      `json:"pair"`
			BusinessType     string      `json:"business_type"`
		} `json:"futures_contract_detail"`
		MarginMode        string      `json:"margin_mode"`
		MarginAccount     string      `json:"margin_account"`
		MarginAsset       string      `json:"margin_asset"`
		MarginBalance     float64     `json:"margin_balance"`
		MarginStatic      float64     `json:"margin_static"`
		MarginPosition    float64     `json:"margin_position"`
		MarginFrozen      float64     `json:"margin_frozen"`
		ProfitReal        float64     `json:"profit_real"`
		ProfitUnreal      float64     `json:"profit_unreal"`
		WithdrawAvailable float64     `json:"withdraw_available"`
		RiskRate          interface{} `json:"risk_rate"`
		PositionMode      string      `json:"position_mode"`
		ContractDetail    []struct {
			Symbol           string      `json:"symbol"`
			ContractCode     string      `json:"contract_code"`
			MarginPosition   float64     `json:"margin_position"`
			MarginFrozen     float64     `json:"margin_frozen"`
			MarginAvailable  float64     `json:"margin_available"`
			ProfitUnreal     float64     `json:"profit_unreal"`
			LiquidationPrice interface{} `json:"liquidation_price"`
			LeverRate        float64     `json:"lever_rate"`
			AdjustFactor     float64     `json:"adjust_factor"`
			ContractType     string      `json:"contract_type"`
			Pair             string      `json:"pair"`
			BusinessType     string      `json:"business_type"`
		} `json:"contract_detail"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostSwapCrossPositionInfo struct {
	ContractCode string `json:"contract_code"`
	Pair         string `json:"pair"`
	ContractType string `json:"contract_type"`
}

type ResPostSwapCrossPositionInfo struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol         string  `json:"symbol"`
		ContractCode   string  `json:"contract_code"`
		Volume         float64 `json:"volume"`
		Available      float64 `json:"available"`
		Frozen         int     `json:"frozen"`
		CostOpen       float64 `json:"cost_open"`
		CostHold       float64 `json:"cost_hold"`
		ProfitUnreal   float64 `json:"profit_unreal"`
		ProfitRate     float64 `json:"profit_rate"`
		LeverRate      int     `json:"lever_rate"`
		PositionMargin float64 `json:"position_margin"`
		Direction      string  `json:"direction"`
		Profit         float64 `json:"profit"`
		LastPrice      float64 `json:"last_price"`
		MarginAsset    string  `json:"margin_asset"`
		MarginMode     string  `json:"margin_mode"`
		MarginAccount  string  `json:"margin_account"`
		ContractType   string  `json:"contract_type"`
		Pair           string  `json:"pair"`
		BusinessType   string  `json:"business_type"`
		PositionMode   string  `json:"position_mode"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ResGetSwapContractInfo struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol            string  `json:"symbol"`
		ContractCode      string  `json:"contract_code"`
		ContractSize      float64 `json:"contract_size"`
		PriceTick         float64 `json:"price_tick"`
		DeliveryDate      string  `json:"delivery_date"`
		DeliveryTime      string  `json:"delivery_time"`
		CreateDate        string  `json:"create_date"`
		ContractStatus    int     `json:"contract_status"`
		SettlementDate    string  `json:"settlement_date"`
		SupportMarginMode string  `json:"support_margin_mode"`
		BusinessType      string  `json:"business_type"`
		Pair              string  `json:"pair"`
		ContractType      string  `json:"contract_type"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ResGetLinerSwapMarkPriceKline struct {
	Ch   string `json:"ch"`
	Data []struct {
		Amount        string `json:"amount"`
		Close         string `json:"close"`
		Count         string `json:"count"`
		High          string `json:"high"`
		ID            int    `json:"id"`
		Low           string `json:"low"`
		Open          string `json:"open"`
		TradeTurnover string `json:"trade_turnover"`
		Vol           string `json:"vol"`
	} `json:"data"`
	Status string `json:"status"`
	Ts     int64  `json:"ts"`
}

type ReqPostSwapSwitchAccountType struct {
	AccountType int `json:"account_type"`
}

type ResPostSwapSwitchAccountType struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		AccountType int `json:"account_type"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

func GetContractCodeFromSymbolAndTypeU(symbol string, type_ common.SymbolType) string {
	symbol = strings.Replace(symbol, "/", "-", 1)
	postFix := ""
	switch type_ {
	case common.SymbolType_FUTURE_THIS_WEEK:
		postFix = "-211210"
	case common.SymbolType_FUTURE_NEXT_WEEK:
		postFix = "-211217"
	case common.SymbolType_FUTURE_THIS_QUARTER:
		postFix = "-211231"
	case common.SymbolType_FUTURE_NEXT_QUARTER:
		postFix = "-210625"
	}
	return symbol + postFix
}
