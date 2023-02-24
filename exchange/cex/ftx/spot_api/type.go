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
	UnrealizedPnl                float64 `json:"unrealizedPnl"`
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
		Coin               string    `json:"coin"`
		Address            string    `json:"address"`
		Tag                string    `json:"tag"`
		Fee                float64   `json:"fee"`
		Id                 int       `json:"id"`
		Size               float64   `json:"size"`
		Status             string    `json:"status"`
		Time               time.Time `json:"time"`
		Method             string    `json:"method"`
		Txid               string    `json:"txid,omitempty"`
		Notes              string    `json:"notes,omitempty"`
		ProposedTransferId string    `json:"proposedTransferId,omitempty"`
		DestinationEmail   string    `json:"destinationEmail,omitempty"`
	} `json:"result"`
}

type RespGetDepositHistory struct {
	Success bool `json:"success"`
	Result  []struct {
		Id      int       `json:"id"`
		Coin    string    `json:"coin"`
		Size    float64   `json:"size"`
		Time    time.Time `json:"time"`
		Notes   string    `json:"notes,omitempty"`
		Status  string    `json:"status"`
		Txid    string    `json:"txid,omitempty"`
		Address struct {
			Address string `json:"address"`
			Tag     string `json:"tag"`
			Method  string `json:"method"`
			Coin    string `json:"coin"`
		} `json:"address,omitempty"`
		Fee           float64   `json:"fee,omitempty"`
		SentTime      time.Time `json:"sentTime,omitempty"`
		ConfirmedTime time.Time `json:"confirmedTime,omitempty"`
		Confirmations int       `json:"confirmations,omitempty"`
		Method        string    `json:"method,omitempty"`
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
	Result  []struct {
		Id      int    `json:"id"`
		Coin    string `json:"coin"`
		Txid    string `json:"txid"`
		Address struct {
			Address string `json:"address"`
			Tag     string `json:"tag"`
			Method  string `json:"method"`
			Coin    string `json:"coin"`
		} `json:"address"`
		Size          float64   `json:"size"`
		Fee           float64   `json:"fee"`
		Status        string    `json:"status"`
		Time          time.Time `json:"time"`
		SentTime      time.Time `json:"sentTime"`
		ConfirmedTime time.Time `json:"confirmedTime"`
		Confirmations int       `json:"confirmations"`
		Method        string    `json:"method"`
	} `json:"result"`
}

type RespTransferSubaccount struct {
	Success bool `json:"success"`
	Result  struct {
		Id     int       `json:"id"`
		Coin   string    `json:"coin"`
		Size   float64   `json:"size"`
		Time   time.Time `json:"time"`
		Notes  string    `json:"notes"`
		Status string    `json:"status"`
	} `json:"result"`
}

type RespSubmitLendingOffer struct {
	Success bool   `json:"success"`
	Result  string `json:"result"`
}

type RespCancelOrder struct {
	Success bool   `json:"success"`
	Result  string `json:"result"`
}

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
		return "erc20"
	case common.Chain_BSC:
		return "bsc"
	case common.Chain_AVALANCHE:
		return "avax"
	case common.Chain_SOLANA:
		return "sol"
	case common.Chain_TRON:
		return "trc20"
	case common.Chain_POLYGON:
		return "matic"
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

func GetSideFromString(side string) order.TradeSide {
	switch side {
	case "buy":
		return order.TradeSide_BUY
	case "sell":
		return order.TradeSide_SELL
	default:
		return order.TradeSide_InvalidSide
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
