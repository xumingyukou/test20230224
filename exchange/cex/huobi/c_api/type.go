package c_api

import (
	"strings"

	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/order"
)

type ResGetFutureContractState struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol            string `json:"symbol"`
		Open              int    `json:"open"`
		Close             int    `json:"close"`
		Cancel            int    `json:"cancel"`
		TransferIn        int    `json:"transfer_in"`
		TransferOut       int    `json:"transfer_out"`
		MasterTransferSub int    `json:"master_transfer_sub"`
		SubTransferMaster int    `json:"sub_transfer_master"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

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
		Symbol         string  `json:"symbol"`
		ContractCode   string  `json:"contract_code"`
		ContractType   string  `json:"contract_type"`
		ContractSize   float64 `json:"contract_size"`
		PriceTick      float64 `json:"price_tick"`
		DeliveryDate   string  `json:"delivery_date"`
		DeliveryTime   string  `json:"delivery_time"`
		CreateDate     string  `json:"create_date"`
		ContractStatus int     `json:"contract_status"`
		SettlementTime string  `json:"settlement_time"`
	} `json:"data"`
	TS int64 `json:"ts"`
}

type DepthInfo struct {
	Ch     string `json:"ch"`
	Status string `json:"status"`
	Tick   struct {
		Asks    [][]float64 `json:"asks"`
		Bids    [][]float64 `json:"bids"`
		Ch      string      `json:"ch"`
		ID      int         `json:"id"`
		Mrid    int64       `json:"mrid"`
		Ts      int64       `json:"ts"`
		Version int         `json:"version"`
	} `json:"tick"`
	Ts int64 `json:"ts"`
}

type ReqPostFutureContractFee struct {
	Symbol string `json:"symbol"`
}

type ResPostFutureContractFee struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol        string `json:"symbol"`
		OpenMakerFee  string `json:"open_maker_fee"`
		OpenTakerFee  string `json:"open_taker_fee"`
		CloseMakerFee string `json:"close_maker_fee"`
		CloseTakerFee string `json:"close_taker_fee"`
		DeliveryFee   string `json:"delivery_fee"`
		FeeAsset      string `json:"fee_asset"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostFutureBalanceValuation struct {
	ValuationAsset string `json:"valuation_asset"`
}

type ResPostFutureBalanceValuation struct {
	Status string `json:"status"`
	Data   []struct {
		ValuationAsset string `json:"valuation_asset"`
		Balance        string `json:"balance"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostFutureOrder struct {
	Symbol         string `json:"symbol"`
	ContractType   string `json:"contract_type"`
	ContractCode   string `json:"contract_code"`
	ClientOrderID  int64  `json:"client_order_id"`
	Direction      string `json:"direction"`
	Offset         string `json:"offset"`
	Price          int    `json:"price"`
	LeverRate      int    `json:"lever_rate"`
	Volume         int    `json:"volume"`
	OrderPriceType string `json:"order_price_type"`
	//TpTriggerPrice   int    `json:"tp_trigger_price"`
	//TpOrderPrice     int    `json:"tp_order_price"`
	//TpOrderPriceType string `json:"tp_order_price_type"`
	//SlTriggerPrice   int    `json:"sl_trigger_price"`
	//SlOrderPrice     int    `json:"sl_order_price"`
	//SlOrderPriceType string `json:"sl_order_price_type"`
}

type ResPostFutureOrder struct {
	Status string `json:"status"`
	Data   struct {
		OrderID       int64  `json:"order_id"`
		OrderIDStr    string `json:"order_id_str"`
		ClientOrderID int64  `json:"client_order_id"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostFutureCancelOrder struct {
	OrderID       string `json:"order_id"`
	ClientOrderId string `json:"client_order_id"`
	Symbol        string `json:"symbol"`
}

type ResPostFutureCancelOrder struct {
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

type ReqPostFutureContractOrderInfo struct {
	OrderID       string `json:"order_id"`
	ClientOrderId string `json:"client_order_id"`
	Symbol        string `json:"symbol"`
}

type ResPostFutureContractOrderInfo struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol          string      `json:"symbol"`
		ContractCode    string      `json:"contract_code"`
		ContractType    string      `json:"contract_type"`
		Volume          int         `json:"volume"`
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
		Fee             float64     `json:"fee"`
		TradeAvgPrice   float64     `json:"trade_avg_price"`
		MarginFrozen    int         `json:"margin_frozen"`
		Profit          int         `json:"profit"`
		Status          int         `json:"status"`
		OrderSource     string      `json:"order_source"`
		OrderIDStr      string      `json:"order_id_str"`
		FeeAsset        string      `json:"fee_asset"`
		LiquidationType string      `json:"liquidation_type"`
		CanceledAt      int         `json:"canceled_at"`
		IsTpsl          int         `json:"is_tpsl"`
		RealProfit      int         `json:"real_profit"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostFutureContractMatchresults struct {
	Contract  string `json:"contract"`
	TradeType int    `json:"trade_type"`
	Symbol    string `json:"symbol"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Direct    string `json:"direct"`
	FromID    int    `json:"from_id"`
}

type ResPostFutureContractMatchresults struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		QueryID          int     `json:"query_id"`
		MatchID          int64   `json:"match_id"`
		OrderID          int64   `json:"order_id"`
		Symbol           string  `json:"symbol"`
		ContractType     string  `json:"contract_type"`
		ContractCode     string  `json:"contract_code"`
		Direction        string  `json:"direction"`
		Offset           string  `json:"offset"`
		TradeVolume      int     `json:"trade_volume"`
		TradePrice       float64 `json:"trade_price"`
		TradeTurnover    int     `json:"trade_turnover"`
		TradeFee         float64 `json:"trade_fee"`
		OffsetProfitloss int     `json:"offset_profitloss"`
		CreateDate       int64   `json:"create_date"`
		Role             string  `json:"role"`
		OrderSource      string  `json:"order_source"`
		OrderIDStr       string  `json:"order_id_str"`
		FeeAsset         string  `json:"fee_asset"`
		RealProfit       int     `json:"real_profit"`
		ID               string  `json:"id"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostFutureContractOpenorders struct {
	Symbol    string `json:"symbol"`
	PageIndex int    `json:"page_index"`
	PageSize  int    `json:"page_size"`
	SortBy    string `json:"sort_by"`
	TradeType int    `json:"trade_type"`
}

type ResPostFutureContractOpenorders struct {
	Status string `json:"status"`
	Data   struct {
		Orders []struct {
			Symbol          string      `json:"symbol"`
			ContractCode    string      `json:"contract_code"`
			ContractType    string      `json:"contract_type"`
			Volume          int         `json:"volume"`
			Price           float64     `json:"price"`
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
			MarginFrozen    int         `json:"margin_frozen"`
			Profit          int         `json:"profit"`
			Status          int         `json:"status"`
			OrderSource     string      `json:"order_source"`
			OrderIDStr      string      `json:"order_id_str"`
			FeeAsset        string      `json:"fee_asset"`
			LiquidationType interface{} `json:"liquidation_type"`
			CanceledAt      interface{} `json:"canceled_at"`
			IsTpsl          int         `json:"is_tpsl"`
			UpdateTime      int64       `json:"update_time"`
			RealProfit      int         `json:"real_profit"`
		} `json:"orders"`
		TotalPage   int `json:"total_page"`
		CurrentPage int `json:"current_page"`
		TotalSize   int `json:"total_size"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ResGetSwapContractState struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol            string `json:"symbol"`
		ContractCode      string `json:"contract_code"`
		Open              int    `json:"open"`
		Close             int    `json:"close"`
		Cancel            int    `json:"cancel"`
		TransferIn        int    `json:"transfer_in"`
		TransferOut       int    `json:"transfer_out"`
		MasterTransferSub int    `json:"master_transfer_sub"`
		SubTransferMaster int    `json:"sub_transfer_master"`
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

type ReqPostSwapOrder struct {
	ContractCode     string `json:"contract_code"`
	ClientOrderID    string `json:"client_order_id"`
	Direction        string `json:"direction"`
	Offset           string `json:"offset"`
	Price            int    `json:"price"`
	LeverRate        int    `json:"lever_rate"`
	Volume           int    `json:"volume"`
	OrderPriceType   string `json:"order_price_type"`
	TpTriggerPrice   int    `json:"tp_trigger_price"`
	TpOrderPrice     int    `json:"tp_order_price"`
	TpOrderPriceType string `json:"tp_order_price_type"`
	//SlTriggerPrice   int    `json:"sl_trigger_price"`
	SlOrderPrice     int    `json:"sl_order_price"`
	SlOrderPriceType string `json:"sl_order_price_type"`
}

type ResPostSwapOrder struct {
	Status string `json:"status"`
	Data   struct {
		OrderID       int64  `json:"order_id"`
		OrderIDStr    string `json:"order_id_str"`
		ClientOrderID int64  `json:"client_order_id"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostSwapCancelOrder struct {
	OrderId       string `json:"order_id"`
	ClientOrderId string `json:"client_order_id"`
	ContractCode  string `json:"contract_code"`
}

type ResPostSwapCancelOrder struct {
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

type ReqPostSwapContractOrderInfo struct {
	OrderID       string `json:"order_id"`
	ClientOrderId string `json:"client_order_id"`
	ContractCode  string `json:"contract_code"`
}

type ResPostSwapContractOrderInfo struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol          string      `json:"symbol"`
		ContractCode    string      `json:"contract_code"`
		Volume          int         `json:"volume"`
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
		RealProfit      int         `json:"real_profit"`
		IsTpsl          int         `json:"is_tpsl"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostSwapContractMatchresults struct {
	Contract  string `json:"contract"`
	TradeType int    `json:"trade_type"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Direct    string `json:"direct"`
	FromID    int    `json:"from_id"`
}

type ResPostSwapContractMatchresults struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		QueryID          int     `json:"query_id"`
		MatchID          int64   `json:"match_id"`
		OrderID          int64   `json:"order_id"`
		Symbol           string  `json:"symbol"`
		ContractCode     string  `json:"contract_code"`
		Direction        string  `json:"direction"`
		Offset           string  `json:"offset"`
		TradeVolume      float64 `json:"trade_volume"`
		TradePrice       float64 `json:"trade_price"`
		TradeTurnover    float64 `json:"trade_turnover"`
		TradeFee         float64 `json:"trade_fee"`
		OffsetProfitloss float64 `json:"offset_profitloss"`
		CreateDate       int64   `json:"create_date"`
		Role             string  `json:"role"`
		OrderSource      string  `json:"order_source"`
		OrderIDStr       string  `json:"order_id_str"`
		ID               string  `json:"id"`
		FeeAsset         string  `json:"fee_asset"`
		RealProfit       int     `json:"real_profit"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostSwapContractOpenorders struct {
	ContractCode string `json:"contract_code"`
	PageIndex    int    `json:"page_index"`
	PageSize     int    `json:"page_size"`
	SortBy       string `json:"sort_by"`
	TradeType    int    `json:"trade_type"`
}

type ResPostSwapContractOpenorders struct {
	Status string `json:"status"`
	Data   struct {
		Orders []struct {
			Symbol          string      `json:"symbol"`
			ContractCode    string      `json:"contract_code"`
			Volume          int         `json:"volume"`
			Price           float64     `json:"price"`
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
			IsTpsl          int         `json:"is_tpsl"`
			UpdateTime      int64       `json:"update_time"`
			RealProfit      int         `json:"real_profit"`
		} `json:"orders"`
		TotalPage   int `json:"total_page"`
		CurrentPage int `json:"current_page"`
		TotalSize   int `json:"total_size"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ResGetFutureContractIndex struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol     string  `json:"symbol"`
		IndexPrice float64 `json:"index_price"`
		IndexTs    int64   `json:"index_ts"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ResGetSwapContractIndex struct {
	Status string `json:"status"`
	Data   []struct {
		IndexPrice   float64 `json:"index_price"`
		IndexTs      int64   `json:"index_ts"`
		ContractCode string  `json:"contract_code"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostSwapContractFee struct {
	ContractCode string `json:"contract_code"`
}

type ResPostSwapContractFee struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol        string `json:"symbol"`
		ContractCode  string `json:"contract_code"`
		OpenMakerFee  string `json:"open_maker_fee"`
		OpenTakerFee  string `json:"open_taker_fee"`
		CloseMakerFee string `json:"close_maker_fee"`
		CloseTakerFee string `json:"close_taker_fee"`
		FeeAsset      string `json:"fee_asset"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostFutureContractAccountInfo struct {
	Symbol string `json:"symbol"`
}

type ResPostFutureContractAccountInfo struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol            string      `json:"symbol"`
		MarginBalance     float64     `json:"margin_balance"`
		MarginPosition    float64     `json:"margin_position"`
		MarginFrozen      int         `json:"margin_frozen"`
		MarginAvailable   float64     `json:"margin_available"`
		ProfitReal        float64     `json:"profit_real"`
		ProfitUnreal      int         `json:"profit_unreal"`
		RiskRate          interface{} `json:"risk_rate"`
		WithdrawAvailable float64     `json:"withdraw_available"`
		LiquidationPrice  interface{} `json:"liquidation_price"`
		LeverRate         int         `json:"lever_rate"`
		AdjustFactor      float64     `json:"adjust_factor"`
		MarginStatic      float64     `json:"margin_static"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostSwapAccountInfo struct {
	ContractCode string `json:"contrct_code"`
}

type ResPostSwapAccountInfo struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol            string  `json:"symbol"`
		MarginBalance     float64 `json:"margin_balance"`
		MarginPosition    float64 `json:"margin_position"`
		MarginFrozen      float64 `json:"margin_frozen"`
		MarginAvailable   float64 `json:"margin_available"`
		ProfitReal        int     `json:"profit_real"`
		ProfitUnreal      float64 `json:"profit_unreal"`
		RiskRate          float64 `json:"risk_rate"`
		WithdrawAvailable float64 `json:"withdraw_available"`
		LiquidationPrice  float64 `json:"liquidation_price"`
		LeverRate         int     `json:"lever_rate"`
		AdjustFactor      float64 `json:"adjust_factor"`
		MarginStatic      float64 `json:"margin_static"`
		ContractCode      string  `json:"contract_code"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostFutureContractPositionInfo struct {
	Symbol string `json:"symbol"`
}

type ResPostFutureContractPositionInfo struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol         string  `json:"symbol"`
		ContractCode   string  `json:"contract_code"`
		ContractType   string  `json:"contract_type"`
		Volume         int     `json:"volume"`
		Available      int     `json:"available"`
		Frozen         int     `json:"frozen"`
		CostOpen       float64 `json:"cost_open"`
		CostHold       float64 `json:"cost_hold"`
		ProfitUnreal   int     `json:"profit_unreal"`
		ProfitRate     int     `json:"profit_rate"`
		LeverRate      int     `json:"lever_rate"`
		PositionMargin float64 `json:"position_margin"`
		Direction      string  `json:"direction"`
		Profit         int     `json:"profit"`
		LastPrice      float64 `json:"last_price"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ReqPostSwapPositionInfo struct {
	ContractCode string `json:"contract_code"`
}

type ResPostSwapPositionInfo struct {
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
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ResGetSwapMarkPriceKline struct {
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

type ReqPostFuturesTransfer struct {
	Currency string  `json"currency"`
	Amount   float64 `json:"amount"`
	Type     string  `json:"type"`
}

type ResPostFuturesTransfer struct {
	Status  string `json:"status"`
	Data    int64  `json:"data"`
	ErrCode string `json:"err-code"`
	ErrMsg  string `json:"err-msg"`
}

func GetContractType(market common.Market, contact_type_str string) common.SymbolType {
	switch market {
	case common.Market_FUTURE_COIN:
		{
			switch contact_type_str {
			case "this_week":
				return common.SymbolType_FUTURE_COIN_THIS_WEEK
			case "next_week":
				return common.SymbolType_FUTURE_COIN_NEXT_WEEK
			case "quarter":
				return common.SymbolType_FUTURE_COIN_THIS_QUARTER
			case "next_quarter":
				return common.SymbolType_FUTURE_COIN_NEXT_QUARTER
			}
		}
	case common.Market_SWAP_COIN:
		{
			return common.SymbolType_SWAP_COIN_FOREVER
		}
	case common.Market_FUTURE:
		{
			switch contact_type_str {
			case "this_week":
				return common.SymbolType_FUTURE_THIS_WEEK
			case "next_week":
				return common.SymbolType_FUTURE_NEXT_WEEK
			case "quarter":
				return common.SymbolType_FUTURE_THIS_QUARTER
			case "next_quarter":
				return common.SymbolType_FUTURE_NEXT_QUARTER
			}
		}
	case common.Market_SWAP:
		{
			return common.SymbolType_SWAP_FOREVER
		}
	}
	return 0
}

func GetContractCodeFromSymbolAndType(symbol string, type_ common.SymbolType) string {
	symbol = strings.Replace(symbol, "/", "-", 1)
	postFix := ""
	switch type_ {
	case common.SymbolType_FUTURE_THIS_WEEK:
		postFix = "-CW"
	case common.SymbolType_FUTURE_NEXT_WEEK:
		postFix = "-NW"
	case common.SymbolType_FUTURE_THIS_QUARTER:
		postFix = "-CQ"
	case common.SymbolType_FUTURE_NEXT_QUARTER:
		postFix = "-NQ"
	case common.SymbolType_FUTURE_COIN_THIS_WEEK:
		symbol = symbol[:3]
		postFix = "_CW"
	case common.SymbolType_FUTURE_COIN_NEXT_WEEK:
		symbol = symbol[:3]
		postFix = "_NW"
	case common.SymbolType_FUTURE_COIN_THIS_QUARTER:
		symbol = symbol[:3]
		postFix = "_CQ"
	case common.SymbolType_FUTURE_COIN_NEXT_QUARTER:
		symbol = symbol[:3]
		postFix = "_NQ"
	}

	return symbol + postFix
}

func GetContractTypeFromType(type_ common.SymbolType) string {
	switch type_ {
	case common.SymbolType_FUTURE_COIN_THIS_WEEK, common.SymbolType_FUTURE_THIS_WEEK:
		return "this_week"
	case common.SymbolType_FUTURE_COIN_NEXT_WEEK, common.SymbolType_FUTURE_NEXT_WEEK:
		return "next_week"
	case common.SymbolType_FUTURE_COIN_THIS_QUARTER, common.SymbolType_FUTURE_THIS_QUARTER:
		return "quarter"
	case common.SymbolType_FUTURE_COIN_NEXT_QUARTER, common.SymbolType_FUTURE_NEXT_QUARTER:
		return "next_quarter"
	}
	return "swap"
}

func GetDirectionFromTradeSide(tradeSide order.TradeSide) string {
	switch tradeSide {
	case order.TradeSide_BUY, order.TradeSide_BUY_TO_CLOSE, order.TradeSide_BUY_TO_OPEN:
		return "buy"
	case order.TradeSide_SELL, order.TradeSide_SELL_TO_CLOSE, order.TradeSide_SELL_TO_OPEN:
		return "sell"
	}
	return ""
}

func GetTradeSideFromDirection(direction string) order.TradeSide {
	switch direction {
	case "buy":
		return order.TradeSide_BUY
	case "sell":
		return order.TradeSide_SELL
	}
	return order.TradeSide_InvalidSide
}

func GetOrderPriceType(order_type order.OrderType, trade_type order.TradeType, tif order.TimeInForce) string {
	if order_type == order.OrderType_LIMIT {
		return "limit"
	} else {
		if trade_type == order.TradeType_MAKER {
			return "post_only"
		} else {
			switch tif {
			case order.TimeInForce_FOK:
				return "fok"
			case order.TimeInForce_IOC:
				return "ioc"
			}
		}
	}
	return ""
}
