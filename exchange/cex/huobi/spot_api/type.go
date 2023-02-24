package spot_api

import (
	"math/big"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/order"
)

const (
	CONTANTTYPE             = "application/json"
	SIGNATUREMETHOT         = "HmacSHA256"
	SPOT_ACCOUNT_ID         = "52330168"
	MARGIN_ACCOUNT_ID       = ""
	SUPER_MARGIN_ACCOUNT_ID = "55303818"
)

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

type EmptyPostParams struct{}

type RespExchangeInfo struct {
	Status string `json:"status"`
	Data   []struct {
		BaseCurrency                    string  `json:"base-currency"`
		QuoteCurrency                   string  `json:"quote-currency"`
		PricePrecision                  int     `json:"price-precision"`
		AmountPrecision                 int     `json:"amount-precision"`
		SymbolPartition                 string  `json:"symbol-partition"`
		Symbol                          string  `json:"symbol"`
		State                           string  `json:"state"`
		ValuePrecision                  int     `json:"value-precision"`
		MinOrderAmt                     float64 `json:"min-order-amt"`
		MaxOrderAmt                     float64 `json:"max-order-amt"`
		MinOrderValue                   float64 `json:"min-order-value"`
		LimitOrderMinOrderAmt           float64 `json:"limit-order-min-order-amt"`
		LimitOrderMaxOrderAmt           float64 `json:"limit-order-max-order-amt"`
		LimitOrderMaxBuyAmt             float64 `json:"limit-order-max-buy-amt"`
		LimitOrderMaxSellAmt            float64 `json:"limit-order-max-sell-amt"`
		BuyLimitMustLessThan            float64 `json:"buy-limit-must-less-than"`
		SellLimitMustGreaterThan        float64 `json:"sell-limit-must-greater-than"`
		SellMarketMinOrderAmt           float64 `json:"sell-market-min-order-amt"`
		SellMarketMaxOrderAmt           float64 `json:"sell-market-max-order-amt"`
		BuyMarketMaxOrderValue          float64 `json:"buy-market-max-order-value"`
		MarketSellOrderRateMustLessThan float64 `json:"market-sell-order-rate-must-less-than"`
		MarketBuyOrderRateMustLessThan  float64 `json:"market-buy-order-rate-must-less-than"`
		MaxOrderValue                   float64 `json:"max-order-value"`
		Underlying                      string  `json:"underlying"`
		MgmtFeeRate                     float64 `json:"mgmt-fee-rate"`
		ChargeTime                      string  `json:"charge-time"`
		RebalTime                       string  `json:"rebal-time"`
		RebalThreshold                  float64 `json:"rebal-threshold"`
		InitNav                         float64 `json:"init-nav"`
		APITrading                      string  `json:"api-trading"`
		Tags                            string  `json:"tags"`
	} `json:"data"`
}

type RespDepth struct {
	Ch     string  `json:"ch"`
	Status string  `json:"status"`
	Ts     big.Int `json:"ts"`
	Tick   struct {
		Ts      int64       `json:"ts"`
		Version int64       `json:"version"`
		Bids    [][]float64 `json:"bids"`
		Asks    [][]float64 `json:"asks"`
	} `json:"tick"`
}

type RespTime struct {
	Status string `json:"status"`
	Data   int64  `json:"data"`
}

type RespAssetTradeFee struct {
	Code         int64           `json:"code"`
	TradeFeeItem []*TradeFeeItem `json:"data"`
	Success      bool            `json:"success"`
}

type TradeFeeItem struct {
	Symbol          string `json:"symbol"`
	ActualMakerRate string `json:"actualMakerRate"`
	ActualTakerRate string `json:"actualTakerRate"`
	TakerFeeRate    string `json:"takerFeeRate"`
	MakerFeeRate    string `json:"makerFeeRate"`
}

// 测试接口
type RespGetMarketStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		MarketStatus    int     `json:"marketStatus"`
		HaltStartTime   big.Int `json:"haltStartTime"`
		HaltEndTime     big.Int `json:"haltEndTime"`
		HaltReason      int     `json:"haltReason"`
		AffectedSymbols string  `json:"affectedSymbols"`
	} `json:"data"`
}

type RespGetAllSymbols struct {
	Status string `json:"status"`
	Data   []struct {
		Tags  string `json:"tags"`
		State string `json:"state"`
		Wr    string `json:"wr"`
		Sc    string `json:"sc"`
		P     []struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Weight int    `json:"weight"`
		} `json:"p"`
		Bcdn string      `json:"bcdn"`
		Qcdn string      `json:"qcdn"`
		Elr  interface{} `json:"elr"`
		Tpp  int         `json:"tpp"`
		Tap  int         `json:"tap"`
		Fp   int         `json:"fp"`
		Smlr interface{} `json:"smlr"`
		Flr  interface{} `json:"flr"`
		Whe  bool        `json:"whe"`
		Cd   bool        `json:"cd"`
		Te   bool        `json:"te"`
		Sp   string      `json:"sp"`
		D    interface{} `json:"d"`
		Bc   string      `json:"bc"`
		Qc   string      `json:"qc"`
		Toa  int64       `json:"toa"`
		Ttp  int         `json:"ttp"`
		W    int         `json:"w"`
		Lr   int         `json:"lr"`
		Dn   string      `json:"dn"`
	} `json:"data"`
	Ts   string `json:"ts"`
	Full int    `json:"full"`
}

type RespGetAllCurrencies struct {
	Status string `json:"status"`
	Data   []struct {
		Tags  string      `json:"tags"`
		Cawt  bool        `json:"cawt"`
		Fc    int         `json:"fc"`
		Sc    int         `json:"sc"`
		Dma   string      `json:"dma"`
		Wma   string      `json:"wma"`
		Ft    string      `json:"ft"`
		Whe   bool        `json:"whe"`
		Cd    bool        `json:"cd"`
		Qc    bool        `json:"qc"`
		Sp    string      `json:"sp"`
		Wp    int         `json:"wp"`
		Fn    string      `json:"fn"`
		At    int         `json:"at"`
		Cc    string      `json:"cc"`
		V     bool        `json:"v"`
		De    bool        `json:"de"`
		Wed   bool        `json:"wed"`
		W     int         `json:"w"`
		State string      `json:"state"`
		Dn    string      `json:"dn"`
		Dd    string      `json:"dd"`
		Svd   interface{} `json:"svd"`
		Swd   interface{} `json:"swd"`
		Sdd   interface{} `json:"sdd"`
		Wd    string      `json:"wd"`
	} `json:"data"`
	Ts   string `json:"ts"`
	Full int    `json:"full"`
}

type RespGetCurrencysSettings struct {
	Status string `json:"status"`
	Data   []struct {
		Tags  string      `json:"tags"`
		Name  string      `json:"name"`
		State string      `json:"state"`
		Cawt  bool        `json:"cawt"`
		Fc    int         `json:"fc"`
		Sc    int         `json:"sc"`
		Sp    string      `json:"sp"`
		Iqc   bool        `json:"iqc"`
		Ct    string      `json:"ct"`
		De    bool        `json:"de"`
		We    bool        `json:"we"`
		Cd    bool        `json:"cd"`
		Oe    int         `json:"oe"`
		V     bool        `json:"v"`
		Whe   bool        `json:"whe"`
		Wet   int64       `json:"wet"`
		Det   int64       `json:"det"`
		Cp    string      `json:"cp"`
		Vat   int64       `json:"vat"`
		Ss    []string    `json:"ss"`
		Fn    string      `json:"fn"`
		Wp    int         `json:"wp"`
		W     int         `json:"w"`
		Dma   string      `json:"dma"`
		Wma   string      `json:"wma"`
		Dn    string      `json:"dn"`
		Dd    string      `json:"dd"`
		Svd   interface{} `json:"svd"`
		Swd   interface{} `json:"swd"`
		Sdd   interface{} `json:"sdd"`
		Wd    string      `json:"wd"`
	} `json:"data"`
	Ts   string `json:"ts"`
	Full int    `json:"full"`
}

type RespGetSymbolsSettings struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol string      `json:"symbol"`
		Tags   string      `json:"tags"`
		State  string      `json:"state"`
		Bcdn   string      `json:"bcdn"`
		Qcdn   string      `json:"qcdn"`
		Elr    interface{} `json:"elr"`
		Tm     string      `json:"tm"`
		Sn     string      `json:"sn"`
		Ve     bool        `json:"ve"`
		Dl     bool        `json:"dl"`
		Te     bool        `json:"te"`
		Ce     bool        `json:"ce"`
		Cd     bool        `json:"cd"`
		Tet    int64       `json:"tet"`
		We     bool        `json:"we"`
		Toa    int64       `json:"toa"`
		Tca    int64       `json:"tca"`
		Voa    int64       `json:"voa"`
		Vca    int64       `json:"vca"`
		Bc     string      `json:"bc"`
		Qc     string      `json:"qc"`
		Sp     string      `json:"sp"`
		D      interface{} `json:"d"`
		Tpp    int         `json:"tpp"`
		Tap    int         `json:"tap"`
		Fp     int         `json:"fp"`
		W      int         `json:"w"`
		Ttp    int         `json:"ttp"`
	} `json:"data"`
	Ts   string `json:"ts"`
	Full int    `json:"full"`
}

type RespGetMarketSymbolsSettings struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol  string  `json:"symbol"`
		State   string  `json:"state"`
		Bc      string  `json:"bc"`
		Qc      string  `json:"qc"`
		Pp      int     `json:"pp"`
		Ap      int     `json:"ap"`
		Sp      string  `json:"sp"`
		Vp      int     `json:"vp"`
		Minoa   float64 `json:"minoa"`
		Maxoa   float64 `json:"maxoa"`
		Minov   int     `json:"minov"`
		Lominoa float64 `json:"lominoa"`
		Lomaxoa float64 `json:"lomaxoa"`
		Lomaxba float64 `json:"lomaxba"`
		Lomaxsa float64 `json:"lomaxsa"`
		Smminoa float64 `json:"smminoa"`
		Blmlt   float64 `json:"blmlt"`
		Slmgt   float64 `json:"slmgt"`
		Smmaxoa float64 `json:"smmaxoa"`
		Bmmaxov int     `json:"bmmaxov"`
		Msormlt float64 `json:"msormlt"`
		Mbormlt float64 `json:"mbormlt"`
		Maxov   int     `json:"maxov"`
		U       string  `json:"u"`
		Mfr     float64 `json:"mfr"`
		Ct      string  `json:"ct"`
		Rt      string  `json:"rt"`
		Rthr    int     `json:"rthr"`
		In      float64 `json:"in"`
		At      string  `json:"at"`
		Tags    string  `json:"tags"`
	} `json:"data"`
	Ts   string `json:"ts"`
	Full int    `json:"full"`
}

type RespGetChainsSettings struct {
	Status string `json:"status"`
	Data   []struct {
		Chain                        string `json:"chain"`
		Currency                     string `json:"currency"`
		Code                         string `json:"code"`
		Ct                           string `json:"ct"`
		Ac                           string `json:"ac"`
		Default                      int    `json:"default"`
		Dma                          string `json:"dma"`
		Wma                          string `json:"wma"`
		De                           bool   `json:"de"`
		We                           bool   `json:"we"`
		Wp                           int    `json:"wp"`
		Ft                           string `json:"ft"`
		Dn                           string `json:"dn"`
		Fn                           string `json:"fn"`
		Awt                          bool   `json:"awt"`
		Adt                          bool   `json:"adt"`
		Ao                           bool   `json:"ao"`
		Fc                           int    `json:"fc"`
		Sc                           int    `json:"sc"`
		V                            bool   `json:"v"`
		Sda                          string `json:"sda"`
		Swa                          string `json:"swa"`
		DepositTipsDesc              string `json:"deposit-tips-desc"`
		WithdrawDesc                 string `json:"withdraw-desc"`
		SuspendDepositDesc           string `json:"suspend-deposit-desc"`
		SuspendWithdrawDesc          string `json:"suspend-withdraw-desc"`
		ReplaceChainInfoDesc         string `json:"replace-chain-info-desc"`
		ReplaceChainNotificationDesc string `json:"replace-chain-notification-desc"`
		ReplaceChainPopupDesc        string `json:"replace-chain-popup-desc"`
	} `json:"data"`
	Ts   string `json:"ts"`
	Full int    `json:"full"`
}

type RespGetCurrenciesChains struct {
	Code int `json:"code"`
	Data []struct {
		Chains []struct {
			Chain                   string `json:"chain"`
			DisplayName             string `json:"displayName"`
			BaseChain               string `json:"baseChain"`
			BaseChainProtocol       string `json:"baseChainProtocol"`
			IsDynamic               bool   `json:"isDynamic"`
			DepositStatus           string `json:"depositStatus"`
			MaxTransactFeeWithdraw  string `json:"maxTransactFeeWithdraw,omitempty"`
			MaxWithdrawAmt          string `json:"maxWithdrawAmt"`
			MinDepositAmt           string `json:"minDepositAmt"`
			MinTransactFeeWithdraw  string `json:"minTransactFeeWithdraw,omitempty"`
			MinWithdrawAmt          string `json:"minWithdrawAmt"`
			NumOfConfirmations      int    `json:"numOfConfirmations"`
			NumOfFastConfirmations  int    `json:"numOfFastConfirmations"`
			WithdrawFeeType         string `json:"withdrawFeeType"`
			WithdrawPrecision       int    `json:"withdrawPrecision"`
			WithdrawQuotaPerDay     string `json:"withdrawQuotaPerDay"`
			WithdrawQuotaPerYear    string `json:"withdrawQuotaPerYear"`
			WithdrawQuotaTotal      string `json:"withdrawQuotaTotal"`
			WithdrawStatus          string `json:"withdrawStatus"`
			TransactFeeRateWithdraw string `json:"transactFeeRateWithdraw,omitempty"`
			TransactFeeWithdraw     string `json:"transactFeeWithdraw,omitempty"`
		} `json:"chains"`
		Currency   string `json:"currency"`
		InstStatus string `json:"instStatus"`
	} `json:"data"`
}

type RespGetKline struct {
	Ch     string `json:"ch"`
	Status string `json:"status"`
	Ts     int64  `json:"ts"`
	Data   []struct {
		ID     int     `json:"id"`
		Open   float64 `json:"open"`
		Close  float64 `json:"close"`
		Low    float64 `json:"low"`
		High   float64 `json:"high"`
		Amount float64 `json:"amount"`
		Vol    float64 `json:"vol"`
		Count  int     `json:"count"`
	} `json:"data"`
}

type RespGetMerged struct {
	Ch     string `json:"ch"`
	Status string `json:"status"`
	Ts     int64  `json:"ts"`
	Tick   struct {
		ID      int64     `json:"id"`
		Version int64     `json:"version"`
		Open    float64   `json:"open"`
		Close   float64   `json:"close"`
		Low     float64   `json:"low"`
		High    float64   `json:"high"`
		Amount  float64   `json:"amount"`
		Vol     float64   `json:"vol"`
		Count   int       `json:"count"`
		Bid     []float64 `json:"bid"`
		Ask     []float64 `json:"ask"`
	} `json:"tick"`
}

type RespGetAllTickers struct {
	Status string `json:"status"`
	Ts     int64  `json:"ts"`
	Data   []struct {
		Symbol  string  `json:"symbol"`
		Open    float64 `json:"open"`
		High    float64 `json:"high"`
		Low     float64 `json:"low"`
		Close   float64 `json:"close"`
		Amount  float64 `json:"amount"`
		Vol     float64 `json:"vol"`
		Count   int     `json:"count"`
		Bid     float64 `json:"bid"`
		BidSize float64 `json:"bidSize"`
		Ask     float64 `json:"ask"`
		AskSize float64 `json:"askSize"`
	} `json:"data"`
}

type RespGetNewestTrade struct {
	Status string `json:"status"`
	Ch     string `json:"ch"`
	Ts     int64  `json:"ts"`
	Tick   struct {
		Id   int64 `json:"id"`
		Ts   int64 `json:"ts"`
		Data []struct {
			Id        big.Int `json:"id"`
			TradeId   int64   `json:"trade-id"`
			amount    float64 `json:"amount"`
			price     float64 `json:"price"`
			Ts        int64   `json:"ts"`
			Direction string  `json:"direction"`
		} `json:"data"`
	} `json:"tick"`
}

type RespGetHistoryTrade struct {
	Ch     string `json:"ch"`
	Status string `json:"status"`
	Ts     int64  `json:"ts"`
	Data   []struct {
		ID   int64 `json:"id"`
		Ts   int64 `json:"ts"`
		Data []struct {
			ID        big.Int `json:"id"`
			Ts        int64   `json:"ts"`
			TradeID   int64   `json:"trade-id"`
			Amount    float64 `json:"amount"`
			Price     float64 `json:"price"`
			Direction string  `json:"direction"`
		} `json:"data"`
	} `json:"data"`
}

type RespGetDetail struct {
	Ch     string `json:"ch"`
	Status string `json:"status"`
	Ts     int64  `json:"ts"`
	Tick   struct {
		ID      int64   `json:"id"`
		Low     float64 `json:"low"`
		High    float64 `json:"high"`
		Open    float64 `json:"open"`
		Close   float64 `json:"close"`
		Vol     float64 `json:"vol"`
		Amount  float64 `json:"amount"`
		Version int64   `json:"version"`
		Count   int     `json:"count"`
	} `json:"tick"`
}

type RespGetAccount struct {
	Status string `json:"status"`
	Data   []struct {
		Id      int64  `json:"id"`
		State   string `json:"state"`
		Type    string `json:"type"`
		subType string `json:"subtype"`
	} `json:"data"`
}

type RespAccountBalance struct {
	Status string `json:"status"`
	Data   struct {
		ID    int    `json:"id"`
		Type  string `json:"type"`
		State string `json:"state"`
		List  []struct {
			Currency string `json:"currency"`
			Type     string `json:"type"`
			Balance  string `json:"balance"`
			SeqNum   string `json:"seq-num"`
			Debt     string `json:"debt"`
		} `json:"list"`
	} `json:"data"`
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

type RespCapitalDepositHisRec struct {
	Status string `json:"status"`
	Data   []struct {
		ID          int     `json:"id"`
		Type        string  `json:"type"`
		SubType     string  `json:"sub-type"`
		Currency    string  `json:"currency"`
		Chain       string  `json:"chain"`
		TxHash      string  `json:"tx-hash"`
		Amount      float64 `json:"amount"`
		FromAddrTag string  `json:"from-addr-tag"`
		Address     string  `json:"address"`
		AddressTag  string  `json:"address-tag"`
		Fee         int     `json:"fee"`
		State       string  `json:"state"`
		CreatedAt   int64   `json:"created-at"`
		UpdatedAt   int64   `json:"updated-at"`
	} `json:"data"`
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

// FIXME: 根据huobi实现GetChainFromNetWork和GetNetworkFromChain
func GetChainFromNetWork(network string) common.Chain {
	switch network {
	case "ETH":
		return common.Chain_ETH
	case "BSC":
		return common.Chain_BSC
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
	default:
		return common.Chain_INVALID_CAHIN
	}
}

func GetNetWorkFromChain(chain common.Chain) string {
	switch chain {
	case common.Chain_ETH:
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

type RespMarginOrderFull struct {
	AccountID     string `json:"account-id"`
	Amount        string `json:"amount"`
	Price         string `json:"price"`
	Source        string `json:"source"`
	Symbol        string `json:"symbol"`
	Type          string `json:"type"`
	ClientOrderID string `json:"client-order-id"`
}

type RespPostOrder struct {
	Status       string `json:"status"`
	Data         string `json:"data"`
	ErrorCode    string `json:"err-code"`
	ErrorMessage string `json:"err-msg"`
}

type RespGetValuation struct {
	Code int `json:"code"`
	Data struct {
		Updated struct {
			Success bool  `json:"success"`
			Time    int64 `json:"time"`
		} `json:"updated"`
		TodayProfitRate          string `json:"todayProfitRate"`
		TotalBalance             string `json:"totalBalance"`
		TodayProfit              string `json:"todayProfit"`
		ProfitAccountBalanceList []struct {
			DistributionType string  `json:"distributionType"`
			Balance          float64 `json:"balance"`
			Success          bool    `json:"success"`
			AccountBalance   string  `json:"accountBalance"`
		} `json:"profitAccountBalanceList"`
	} `json:"data"`
	Success bool `json:"success"`
}

type RespGetAssetValuation struct {
	Code int `json:"code"`
	Data struct {
		Balance   string `json:"balance"`
		Timestamp int64  `json:"timestamp"`
	} `json:"data"`
	Ok bool `json:"ok"`
}

type RespGetAccountHistory struct {
	Status string `json:"status"`
	Data   []struct {
		AccountID    int    `json:"account-id"`
		Currency     string `json:"currency"`
		RecordID     int64  `json:"record-id"`
		TransactAmt  string `json:"transact-amt"`
		TransactType string `json:"transact-type"`
		AvailBalance string `json:"avail-balance"`
		AcctBalance  string `json:"acct-balance"`
		TransactTime int64  `json:"transact-time"`
	} `json:"data"`
	NextID int64 `json:"next-id"`
}

type RespGetAccountLedger struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		AccountID    int     `json:"accountId"`
		Currency     string  `json:"currency"`
		TransactAmt  float64 `json:"transactAmt"`
		TransactType string  `json:"transactType"`
		TransferType string  `json:"transferType"`
		TransactID   int     `json:"transactId"`
		TransactTime int64   `json:"transactTime"`
		Transferer   int     `json:"transferer"`
		Transferee   int     `json:"transferee"`
	} `json:"data"`
	NextID int  `json:"nextId"`
	Ok     bool `json:"ok"`
}

type RespPostCFuturesTransfer struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	ErrCode string      `json:"err-code"`
	ErrMsg  string      `json:"err-msg"`
}

type RespGetPointAccount struct {
	Code int `json:"code"`
	Data struct {
		AccountID string `json:"accountId"`
		GroupIds  []struct {
			GroupID    int    `json:"groupId"`
			ExpiryDate int64  `json:"expiryDate"`
			RemainAmt  string `json:"remainAmt"`
		} `json:"groupIds"`
		AcctBalance   string `json:"acctBalance"`
		AccountStatus string `json:"accountStatus"`
	} `json:"data"`
	Success bool `json:"success"`
}

type RespPostPointTransfer struct {
	Code int `json:"code"`
	Data struct {
		TransactID   string `json:"transactId"`
		TransactTime int64  `json:"transactTime"`
	} `json:"data"`
	Success bool `json:"success"`
}

type RespGetDepositAddress struct {
	Code int `json:"code"`
	Data []struct {
		UserID     int    `json:"userId"`
		Currency   string `json:"currency"`
		Address    string `json:"address"`
		AddressTag string `json:"addressTag"`
		Chain      string `json:"chain"`
	} `json:"data"`
}

type RespGetWithdrawQuota struct {
	Code int `json:"code"`
	Data struct {
		Currency string `json:"currency"`
		Chains   []struct {
			Chain                      string `json:"chain"`
			MaxWithdrawAmt             string `json:"maxWithdrawAmt"`
			WithdrawQuotaPerDay        string `json:"withdrawQuotaPerDay"`
			RemainWithdrawQuotaPerDay  string `json:"remainWithdrawQuotaPerDay"`
			WithdrawQuotaPerYear       string `json:"withdrawQuotaPerYear"`
			RemainWithdrawQuotaPerYear string `json:"remainWithdrawQuotaPerYear"`
			WithdrawQuotaTotal         string `json:"withdrawQuotaTotal"`
			RemainWithdrawQuotaTotal   string `json:"remainWithdrawQuotaTotal"`
		} `json:"chains"`
	} `json:"data"`
}

type RespGetWithdrawAddress struct {
	Code int `json:"code"`
	Data []struct {
		Currency   string `json:"currency"`
		Chain      string `json:"chain"`
		Note       string `json:"note"`
		AddressTag string `json:"addressTag"`
		Address    string `json:"address"`
	} `json:"data"`
	NextID int `json:"next-id"`
}

type RespPostWithdrawCreate struct {
	Data int `json:"data"`
}

type RespGetWithdrawClientOrderId struct {
	Status string `json:"status"`
	Data   struct {
		ID            int     `json:"id"`
		ClientOrderID string  `json:"client-order-id"`
		Type          string  `json:"type"`
		SubType       string  `json:"sub-type"`
		Currency      string  `json:"currency"`
		Chain         string  `json:"chain"`
		TxHash        string  `json:"tx-hash"`
		Amount        float64 `json:"amount"`
		FromAddrTag   string  `json:"from-addr-tag"`
		Address       string  `json:"address"`
		AddressTag    string  `json:"address-tag"`
		Fee           int     `json:"fee"`
		State         string  `json:"state"`
		CreatedAt     int64   `json:"created-at"`
		UpdatedAt     int64   `json:"updated-at"`
	} `json:"data"`
}

type RespPostCancelWithdrawCreate struct {
	Data int `json:"data"`
}

type RespPostSubUserDeductMode struct {
	Code int `json:"code"`
	Data []struct {
		SubUID     string `json:"subUid"`
		DeductMode string `json:"deductMode"`
		ErrCode    int    `json:"errCode,omitempty"`
		ErrMessage string `json:"errMessage,omitempty"`
	} `json:"data"`
	Ok bool `json:"ok"`
}

type RespGetUserApiKey struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		AccessKey   string `json:"accessKey"`
		Status      string `json:"status"`
		Note        string `json:"note"`
		Permission  string `json:"permission"`
		IPAddresses string `json:"ipAddresses"`
		ValidDays   int    `json:"validDays"`
		CreateTime  int64  `json:"createTime"`
		UpdateTime  int64  `json:"updateTime"`
	} `json:"data"`
	Ok bool `json:"ok"`
}

type RespGetUserUid struct {
	Code int `json:"code"`
	Data int `json:"data"`
}

type ReqPostSubUserCreation struct {
	UserList []struct {
		UserName string `json:"userName"`
		Note     string `json:"note"`
	} `json:"userList"`
}

type RespPostSubUserCreation struct {
	Code int `json:"code"`
	Data []struct {
		UserName   string `json:"userName"`
		Note       string `json:"note"`
		UID        int    `json:"uid,omitempty"`
		ErrCode    string `json:"errCode,omitempty"`
		ErrMessage string `json:"errMessage,omitempty"`
	} `json:"data"`
}

type RespGetSubUserList struct {
	Code int `json:"code"`
	Data []struct {
		UID       int    `json:"uid"`
		UserState string `json:"userState"`
	} `json:"data"`
}

type RespPostSubUserManagement struct {
	Code int `json:"code"`
	Data struct {
		SubUID    int    `json:"subUid"`
		UserState string `json:"userState"`
	} `json:"data"`
	Ok bool `json:"ok"`
}

type ReqPostCFuturesTransfer struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
	Type     string  `json:"type"`
}

type ReqPostOrder struct {
	AccountID        string `json:"account-id"`
	Amount           string `json:"amount"`
	Price            string `json:"price"`
	Source           string `json:"source"`
	Symbol           string `json:"symbol"`
	Type             string `json:"type"`
	ClientOrderID    string `json:"client-order-id"`
	SelfMatchPrevent int    `json:"self-match-prevent"`
	//StopPrice        string `json:"stop-price"`
	//Operator         string `json:"operator"`
}

type ReqPostPointTransfer struct {
	FromUid string `json:"fromUid"`
	ToUid   string `json:"toUid"`
	GroupId int64  `json:"groupId"`
	Amount  string `json:"amount"`
}

type ReqPostSubUserDeductMode struct {
	SubUids    int64  `json:"subUids"`
	DeductMode string `json:"deductMode"`
}

type ReqPostSubUserManagement struct {
	SubUids int64  `json:"subUids"`
	Action  string `json:"action"`
}

type ReqPostWithdrawCreate struct {
	Address       string `json:"address"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	Fee           string `json:"fee"`
	Chain         string `json:"chain"`
	AddrTag       string `json:"addr-tag"`
	ClientOrderId string `json:"client-order-id"`
}

type RespGetSubUserState struct {
	Code int `json:"code"`
	Data struct {
		UID       int    `json:"uid"`
		UserState string `json:"userState"`
	} `json:"data"`
}

type ReqPostBatchOrders []struct {
	AccountID     string `json:"account-id"`
	Symbol        string `json:"symbol"`
	Type          string `json:"type"`
	Amount        string `json:"amount"`
	Price         string `json:"price"`
	Source        string `json:"source"`
	ClientOrderID string `json:"client-order-id"`
}

type RespPostBatchOrders struct {
	Status string `json:"status"`
	Data   []struct {
		OrderID       int64  `json:"order-id"`
		ClientOrderID string `json:"client-order-id"`
		ErrCode       string `json:"err-code"`
		ErrMsg        string `json:"err-msg"`
	} `json:"data"`
}

type ReqPostSubmitOrder struct {
	Symbol string `json:""symbol`
}

type RespPostSubmitOrder struct {
	Status     string `json:"status"`
	ErrCode    string `json:"err-code"`
	ErrMsg     string `json:"err-msg"`
	Data       string `json:"data"`
	OrderState int    `json:"order-state"`
}

type RespGetOrder struct {
	Status string `json:"status"`
	Data   struct {
		ID              int64  `json:"id"`
		Symbol          string `json:"symbol"`
		AccountID       int    `json:"account-id"`
		ClientOrderID   string `json:"client-order-id"`
		Amount          string `json:"amount"`
		Price           string `json:"price"`
		CreatedAt       int64  `json:"created-at"`
		Type            string `json:"type"`
		FieldAmount     string `json:"field-amount"`
		FieldCashAmount string `json:"field-cash-amount"`
		FieldFees       string `json:"field-fees"`
		FinishedAt      int    `json:"finished-at"`
		Source          string `json:"source"`
		State           string `json:"state"`
		CanceledAt      int    `json:"canceled-at"`
	} `json:"data"`
}

type RespGetHistoryOrders struct {
	Status string `json:"status"`
	Data   []struct {
		ID              int64  `json:"id"`
		Symbol          string `json:"symbol"`
		AccountID       int    `json:"account-id"`
		ClientOrderID   string `json:"client-order-id"`
		Amount          string `json:"amount"`
		Price           string `json:"price"`
		CreatedAt       int64  `json:"created-at"`
		Type            string `json:"type"`
		FieldAmount     string `json:"field-amount"`
		FieldCashAmount string `json:"field-cash-amount"`
		FieldFees       string `json:"field-fees"`
		FinishedAt      int64  `json:"finished-at"`
		Source          string `json:"source"`
		State           string `json:"state"`
		CanceledAt      int    `json:"canceled-at"`
	} `json:"data"`
}

type RespGetOrderMatchresults struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol            string `json:"symbol"`
		FeeCurrency       string `json:"fee-currency"`
		Source            string `json:"source"`
		OrderID           int64  `json:"order-id"`
		Price             string `json:"price"`
		CreatedAt         int64  `json:"created-at"`
		Role              string `json:"role"`
		MatchID           int    `json:"match-id"`
		FilledAmount      string `json:"filled-amount"`
		FilledFees        string `json:"filled-fees"`
		FilledPoints      string `json:"filled-points"`
		FeeDeductCurrency string `json:"fee-deduct-currency"`
		FeeDeductState    string `json:"fee-deduct-state"`
		TradeID           int    `json:"trade-id"`
		ID                int64  `json:"id"`
		Type              string `json:"type"`
	} `json:"data"`
}

type ReqPostMarginOrders struct {
	Symbol   string `json:"symbol"`
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type RespPostMarginOrders struct {
	Status string `json:"status"`
	Data   int    `json:"data"`
}

type ReqPostCrossMarginOrders struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type RespPostCrossMarginOrders struct {
	Status string `json:"status"`
	Data   int    `json:"data"`
}

type RespGetMarginLoanOrders struct {
	Status string `json:"status"`
	Data   []struct {
		DeductRate       string `json:"deduct-rate"`
		CreatedAt        int64  `json:"created-at"`
		UpdatedAt        int64  `json:"updated-at"`
		AccruedAt        int64  `json:"accrued-at"`
		InterestAmount   string `json:"interest-amount"`
		LoanAmount       string `json:"loan-amount"`
		HourInterestRate string `json:"hour-interest-rate"`
		LoanBalance      string `json:"loan-balance"`
		InterestBalance  string `json:"interest-balance"`
		PaidCoin         string `json:"paid-coin"`
		DayInterestRate  string `json:"day-interest-rate"`
		InterestRate     string `json:"interest-rate"`
		UserID           int    `json:"user-id"`
		AccountID        int    `json:"account-id"`
		Currency         string `json:"currency"`
		Symbol           string `json:"symbol"`
		PaidPoint        string `json:"paid-point"`
		DeductCurrency   string `json:"deduct-currency"`
		DeductAmount     string `json:"deduct-amount"`
		ID               int    `json:"id"`
		State            string `json:"state"`
	} `json:"data"`
}

type RespGetCrossMarginLoanOrders struct {
	Status string `json:"status"`
	Data   []struct {
		LoanBalance     string `json:"loan-balance"`
		InterestBalance string `json:"interest-balance"`
		LoanAmount      string `json:"loan-amount"`
		AccruedAt       int64  `json:"accrued-at"`
		InterestAmount  string `json:"interest-amount"`
		FilledPoints    string `json:"filled-points"`
		FilledHt        string `json:"filled-ht"`
		Currency        string `json:"currency"`
		ID              int    `json:"id"`
		State           string `json:"state"`
		AccountID       int    `json:"account-id"`
		UserID          int    `json:"user-id"`
		CreatedAt       int64  `json:"created-at"`
	} `json:"data"`
}

type ReqPostMarginRepay struct {
	Amount string `json:"amount"`
}

type RespPostMarginRepay struct {
	Data int `json:"data"`
}

type ReqPostCrossMarginRepay struct {
	Amount string `json:"amount"`
}

type RespPostCrossMarginRepay struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type RespGetRepayment struct {
	Code int `json:"code"`
	Data []struct {
		RepayID      string `json:"repayId"`
		RepayTime    int64  `json:"repayTime"`
		AccountID    string `json:"accountId"`
		Currency     string `json:"currency"`
		RepaidAmount string `json:"repaidAmount"`
		TransactIds  struct {
			TransactID      int    `json:"transactId"`
			Repaidprincipal string `json:"repaidprincipal"`
			RepaidInterest  string `json:"repaidInterest"`
			PaidHt          string `json:"paidHt"`
			PaidPoint       string `json:"paidPoint"`
		} `json:"transactIds"`
	} `json:"data"`
}

type RespGetReferenceCurrencies struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		Currency string `json:"currency"`
		Chains   []struct {
			Chain                   string `json:"chain"`
			DisplayName             string `json:"displayName"`
			BaseChain               string `json:"baseChain"`
			BaseChainProtocol       string `json:"baseChainProtocol"`
			IsDynamic               bool   `json:"isDynamic"`
			DepositStatus           string `json:"depositStatus"`
			MaxTransactFeeWithdraw  string `json:"maxTransactFeeWithdraw,omitempty"`
			MaxWithdrawAmt          string `json:"maxWithdrawAmt"`
			MinDepositAmt           string `json:"minDepositAmt"`
			MinTransactFeeWithdraw  string `json:"minTransactFeeWithdraw,omitempty"`
			MinWithdrawAmt          string `json:"minWithdrawAmt"`
			NumOfConfirmations      int    `json:"numOfConfirmations"`
			NumOfFastConfirmations  int    `json:"numOfFastConfirmations"`
			WithdrawFeeType         string `json:"withdrawFeeType"`
			WithdrawPrecision       int    `json:"withdrawPrecision"`
			WithdrawQuotaPerDay     string `json:"withdrawQuotaPerDay"`
			WithdrawQuotaPerYear    string `json:"withdrawQuotaPerYear"`
			WithdrawQuotaTotal      string `json:"withdrawQuotaTotal"`
			WithdrawStatus          string `json:"withdrawStatus"`
			TransactFeeRateWithdraw string `json:"transactFeeRateWithdraw,omitempty"` // 提币手续费率
			TransactFeeWithdraw     string `json:"transactFeeWithdraw,omitempty"`     // 提币手续费
		} `json:"chains"`
		InstStatus string `json:"instStatus"`
	} `json:"data"`
}

type RespGetMatchresults struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol            string `json:"symbol"`
		FeeCurrency       string `json:"fee-currency"`
		Source            string `json:"source"`
		Price             string `json:"price"`
		CreatedAt         int64  `json:"created-at"`
		Role              string `json:"role"`
		OrderID           int64  `json:"order-id"`
		MatchID           int    `json:"match-id"`
		TradeID           int    `json:"trade-id"`
		FilledAmount      string `json:"filled-amount"`
		FilledFees        string `json:"filled-fees"`
		FilledPoints      string `json:"filled-points"`
		FeeDeductCurrency string `json:"fee-deduct-currency"`
		FeeDeductState    string `json:"fee-deduct-state"`
		ID                int64  `json:"id"`
		Type              string `json:"type"`
	} `json:"data"`
}

type RespGetClientOrder struct {
	Status string `json:"status"`
	Data   struct {
		ID              int64  `json:"id"`
		Symbol          string `json:"symbol"`
		AccountID       int    `json:"account-id"`
		ClientOrderID   string `json:"client-order-id"`
		Amount          string `json:"amount"`
		Price           string `json:"price"`
		CreatedAt       int64  `json:"created-at"`
		Type            string `json:"type"`
		FieldAmount     string `json:"field-amount"`
		FieldCashAmount string `json:"field-cash-amount"`
		FieldFees       string `json:"field-fees"`
		FinishedAt      int    `json:"finished-at"`
		Source          string `json:"source"`
		State           string `json:"state"`
		CanceledAt      int    `json:"canceled-at"`
	} `json:"data"`
}

type RespGetOpenOrders struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol           string `json:"symbol"`
		Source           string `json:"source"`
		Price            string `json:"price"`
		CreatedAt        int64  `json:"created-at"`
		Amount           string `json:"amount"`
		AccountID        int    `json:"account-id"`
		FilledCashAmount string `json:"filled-cash-amount"`
		ClientOrderID    string `json:"client-order-id"`
		FilledAmount     string `json:"filled-amount"`
		FilledFees       string `json:"filled-fees"`
		ID               int64  `json:"id"`
		State            string `json:"state"`
		Type             string `json:"type"`
	} `json:"data"`
}

type RespTestGetReferenceCurrencies struct {
	Currency string `json:"currency"`
	Chains   []Chain
}

type Chain struct {
	Chain       string `json:"chain"`
	DisplayName string `json:"displayName"`
	BaseChain   string `json:"baseChain"`
}

type RespGetDepositHistory struct {
	Status string `json:"status"`
	Data   []struct {
		ID          int     `json:"id"`
		Type        string  `json:"type"`
		SubType     string  `json:"sub-type"`
		Currency    string  `json:"currency"`
		Chain       string  `json:"chain"`
		TxHash      string  `json:"tx-hash"`
		Amount      float64 `json:"amount"`
		FromAddrTag string  `json:"from-addr-tag"`
		Address     string  `json:"address"`
		AddressTag  string  `json:"address-tag"`
		Fee         float64 `json:"fee"`
		State       string  `json:"state"`
		CreatedAt   int64   `json:"created-at"`
		UpdatedAt   int64   `json:"updated-at"`
	} `json:"data"`
}

type ReqPostTransfer struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
	Type     string  `json:"type"`
}

type RespPostTransfer struct {
	Data    int64  `json:"data"`
	Status  string `json:"status"`
	ErrCode string `json:"err-code"`
	ErrMsg  string `json:"err-msg"`
}

type RespResPostFuturesTransfer struct {
	Status  string `json:"status"`
	Data    int64  `json:"data"`
	ErrCode string `json:"err-code"`
	ErrMsg  string `json:"err-msg"`
}

type ReqPostCSwapTransfer struct {
	From     string  `json:"from"`
	To       string  `json:"to"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type RespPostCSwapTransfer struct {
	Code    int64  `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    int64  `json:"data"`
}

type ReqPostUTransfer struct {
	From          string  `json:"from"`
	To            string  `json:"to"`
	Currency      string  `json:"currency"`
	Amount        float64 `json:"amount"`
	MarginAccount string  `json:"margin-account"`
}

type RespPostUTransfer struct {
	Code    int64  `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    int64  `json:"data"`
}

func GetTypeFromTradeOrderSideTif(OrderType order.OrderType, Side order.TradeSide, TradeType order.TradeType, Tif order.TimeInForce) string {
	switch Side {
	case order.TradeSide_BUY:
		{
			if TradeType == order.TradeType_MAKER {
				return "buy-limit-maker"
			} else {
				switch OrderType {
				case order.OrderType_LIMIT:
					{
						if Tif == order.TimeInForce_IOC {
							return "buy-ioc"
						} else if Tif == order.TimeInForce_FOK {
							return "buy-limit-fok"
						} else {
							return "buy-limit"
						}
					}
				case order.OrderType_MARKET:
					return "buy-market"
				case order.OrderType_STOP_LOSS_LIMIT:
					{
						if Tif == order.TimeInForce_FOK {
							return "buy-stop-limit-fok"
						} else {
							return "buy-stop-limit"
						}
					}
				}
			}
		}
	case order.TradeSide_SELL:
		{
			if TradeType == order.TradeType_MAKER {
				return "sell-limit-maker"
			} else {
				switch OrderType {
				case order.OrderType_LIMIT:
					{
						if Tif == order.TimeInForce_IOC {
							return "sell-ioc"
						} else if Tif == order.TimeInForce_FOK {
							return "sell-limit-fok"
						} else {
							return "sell-limit"
						}
					}
				case order.OrderType_MARKET:
					return "sell-market"
				case order.OrderType_STOP_LOSS_LIMIT:
					{
						if Tif == order.TimeInForce_FOK {
							return "sell-stop-limit-fok"
						} else {
							return "sell-stop-limit"
						}
					}
				}
			}
		}
	}
	return ""
}

func GetTypeFromType(market common.Market) string {
	switch market {
	case common.Market_SPOT:
		return "spot"
	case common.Market_SWAP_COIN:
		return "swap"
	case common.Market_FUTURE, common.Market_SWAP:
		return "linear-swap"
	}
	return ""
}
