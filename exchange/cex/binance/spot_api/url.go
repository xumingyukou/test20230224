package spot_api

import (
	"clients/exchange/cex/base"
	"clients/transform"
	"net/url"
	"sync"
)

const (
	SPOT_API_BASE_URL      = "https://api.binance.com"
	SPOT_TEST_API_BASE_URL = "https://testnet.binance.vision"
	SPOT_API_US_URL        = "https://api.binance.us"
	SPOT_API_JE_URL        = "https://api.binance.je"
	UBASE_API_BASE_URL     = "https://fapi.binance.com"
	UBASE_API_US_URL       = "https://fapi.binance.us"
	UBASE_API_JE_URL       = "https://fapi.binance.je"
	UBASE_TEST_BASE_URL    = "https://testnet.binancefuture.com"
	CBASE_API_BASE_URL     = "https://dapi.binance.com"
	CBASE_API_US_URL       = "https://dapi.binance.us"
	CBASE_API_JE_URL       = "https://dapi.binance.je"
	API_V1_URL             = "/api/v1/"
	API_V3_URL             = "/api/v3/"
	SAPI_V1_URL            = "/sapi/v1/"
	FAPI_V1_URL            = "/fapi/v1/"
	FAPI_V2_URL            = "/fapi/v2/"
	DAPI_V1_URL            = "/dapi/v1/"
	DAPI_V2_URL            = "/dapi/v2/"

	WS_API_BASE_URL  = "wss://stream.binance.com:9443"
	WS_API_US_URL    = "wss://stream.binance.us:9443"
	WS_API_JE_URL    = "wss://stream.binance.je:9443"
	WS_UAPI_BASE_URL = "wss://fstream.binance.com"
	WS_UAPI_US_URL   = "wss://fstream.binance.us"
	WS_UAPI_JE_URL   = "wss://fstream.binance.je"
	WS_CAPI_BASE_URL = "wss://dstream.binance.com"
	WS_CAPI_US_URL   = "wss://dstream.binance.us"
	WS_CAPI_JE_URL   = "wss://dstream.binance.je"
	SINGLE_API_URL   = "/ws/"
	STREAM_API_URL   = "/stream?streams="
)

type ReqUrl struct {
	DEPTH_URL,
	OPENORDERS_URL,
	MARGIN_OPENORDERS_URL,
	MYTRADES_URL,
	TIME_URL,
	PING_URL,
	TRADE_URL,
	HISTORICAL_TRADES_URL,
	EXCHANGEINFO_URL,
	CAPITAL_CONFIG_GETALL_URL,
	ASSET_TRADEFEE_URL,
	CAPITAL_WITHDRAW_APPLY_URL,
	CAPITAL_DEPOSIT_HISREC_URL,
	ASSET_TRANSFER_URL,
	SUBACCOUNT_UNIVERSALTRANSFER_URL,
	SUBACCOUNT_TRANSFER_SUBTOMASTER_URL,
	SUBACCOUNT_TRANSFER_SUBTOSUB_URL,
	SUBACCOUNT_TRANSFER_SUBUSERHISTORY_URL,
	ACCOUNT_URL,
	ASSET_URL,
	MARGIN_ACCOUNT_URL,
	MARGIN_ISOLATED_ACCOUNT_URL,
	ORDER_URL,
	ALLORDERS_URL,
	MARGIN_ORDER_URL,
	CAPITAL_WITHDRAW_HISTORY_URL,
	MARGIN_LOAN_URL,
	MARGIN_REPAY_URL,
	MARGIN_ALLORDERS_URL,
	USERDATASTREAM_URL,
	MARGINUSERDATASTREAM_URL,
	MARGINISOLATEDUSERDATASTREAM_URL,
	MARGIN_ISOLATED_PAIR_URL,
	MARGIN_ISOLATED_ALLPAIRS_URL,
	SUBACCOUNT_VIRTUALSUBACCOUNT_URL,
	SUBACCOUNT_LIST_URL,
	//ubase
	PREMIUMINDEX_URL,
	TICKER_PRICE_URL,
	POSITIONSIDE_DUAL_URL,
	BATCHORDERS_URL,
	ALLOPENORDERS_URL,
	OPENORDER_URL,
	BALANCE_URL,
	LEVERAGE_URL,
	MARGINTYPE_URL,
	POSITIONMARGIN_URL,
	POSITIONRISK_URL,
	COMMISSIONRATE_URL base.ReqUrlInfo
}

type WsReqUrl struct {
	AGGTRADE_URL         string
	MARKPRICE_URL        string
	TRADE_URL            string
	BOOK_TICKER_URL      string
	DEPTH_LIMIT_FULL_URL string
	DEPTH_INCRE_URL      string
	USER_DATA_URL        string
}

func NewSpotWsUrl() *WsReqUrl {
	return &WsReqUrl{
		MARKPRICE_URL:        "<symbol>@markPrice@1s",
		AGGTRADE_URL:         "<symbol>@aggTrade",
		TRADE_URL:            "<symbol>@trade",
		BOOK_TICKER_URL:      "<symbol>@bookTicker",
		DEPTH_LIMIT_FULL_URL: "<symbol>@depth",
		DEPTH_INCRE_URL:      "<symbol>@depth",
		USER_DATA_URL:        "<listenKey>",
	}
}

var (
	mtx sync.Mutex
	// 基于ip的限制
	IP_1M_LIMIT_MAP, IP_1M_LIMIT_MAP_2, IP_1M_LIMIT_MAP_3, IP_1M_LIMIT_MAP_5, IP_1M_LIMIT_MAP_6, IP_1M_LIMIT_MAP_10, IP_1M_LIMIT_MAP_20, IP_1M_LIMIT_MAP_40, IP_1M_LIMIT_MAP_50, IP_1M_LIMIT_MAP_200 *base.RateLimitConsume
)

func (r *ReqUrl) GetURLRateLimit(url base.ReqUrlInfo, params url.Values) []*base.RateLimitConsume {
	if url.Url == API_V3_URL+"depth" {
		limit := transform.StringToX[int64](params.Get("limit")).(int64)
		if limit > 101 && limit <= 500 {
			return []*base.RateLimitConsume{IP_1M_LIMIT_MAP_5}
		} else if limit > 501 && limit <= 1000 {
			return []*base.RateLimitConsume{IP_1M_LIMIT_MAP_10}
		} else if limit > 1000 && limit <= 5000 {
			return []*base.RateLimitConsume{IP_1M_LIMIT_MAP_50}
		} else {
			return []*base.RateLimitConsume{IP_1M_LIMIT_MAP}
		}
	} else if url.Url == API_V3_URL+"ticker/price" {
		if params.Get("symbol") != "" && params.Get("symbols") == "" {
			return []*base.RateLimitConsume{IP_1M_LIMIT_MAP}
		} else {
			return []*base.RateLimitConsume{IP_1M_LIMIT_MAP_2}
		}
	} else if url.Url == FAPI_V1_URL+"depth" || url.Url == DAPI_V1_URL+"depth" {
		limit := transform.StringToX[int64](params.Get("limit")).(int64)
		if limit == 100 {
			return []*base.RateLimitConsume{IP_1M_LIMIT_MAP_5}
		} else if limit == 500 {
			return []*base.RateLimitConsume{IP_1M_LIMIT_MAP_10}
		} else if limit == 1000 {
			return []*base.RateLimitConsume{IP_1M_LIMIT_MAP_20}
		} else {
			return []*base.RateLimitConsume{IP_1M_LIMIT_MAP_2}
		}
	} else if url.Url == FAPI_V1_URL+"ticker/price" || url.Url == DAPI_V1_URL+"ticker/price" {
		if params.Get("symbol") != "" {
			return []*base.RateLimitConsume{IP_1M_LIMIT_MAP}
		} else {
			return []*base.RateLimitConsume{IP_1M_LIMIT_MAP_2}
		}
	} else if url.Url == DAPI_V1_URL+"allOrders" {
		if params.Get("pair") != "" {
			return []*base.RateLimitConsume{IP_1M_LIMIT_MAP_40}
		} else {
			return []*base.RateLimitConsume{IP_1M_LIMIT_MAP_20}
		}
	} else {
		return url.RateLimitConsumeMap
	}
}

func NewSpotReqUrl(types map[string]int64, account string) *ReqUrl {
	//设置默认值
	var (
		R_IP_1M_limit      int64 = 1200
		R_ACCOUNT_1M_limit int64 = 1200
		R_ACCOUNT_1D_limit int64 = 160000
	)
	if types != nil {
		if value, ok := types[R_IP_1M]; ok {
			R_IP_1M_limit = value
		}
		if value, ok := types[R_ACCOUNT_1M]; ok {
			R_ACCOUNT_1M_limit = value
		}
		if value, ok := types[R_ACCOUNT_1D]; ok {
			R_ACCOUNT_1D_limit = value
		}
	}
	initIpLimit(R_IP_1M_limit)

	ORDER_LIMIT_MAP_3000 := []*base.RateLimitConsume{base.NewRLConsume(base.RateLimitType(R_ACCOUNT_1M+"_"+account), 3000, R_ACCOUNT_1M_limit), base.NewRLConsume(base.RateLimitType(R_ACCOUNT_1D+"_"+account), 3000, R_ACCOUNT_1D_limit)}
	ORDER_LIMIT_MAP := []*base.RateLimitConsume{base.NewRLConsume(R_IP_1M, 1, R_IP_1M_limit), base.NewRLConsume(base.RateLimitType(R_ACCOUNT_1M+"_"+account), 1, R_ACCOUNT_1M_limit), base.NewRLConsume(base.RateLimitType(R_ACCOUNT_1D+"_"+account), 1, R_ACCOUNT_1D_limit)}

	return &ReqUrl{
		DEPTH_URL:                              base.NewReqUrlInfo(API_V3_URL + "depth"),
		TICKER_PRICE_URL:                       base.NewReqUrlInfo(API_V3_URL + "ticker/price"),
		OPENORDERS_URL:                         base.NewReqUrlInfo(API_V3_URL+"openOrders", IP_1M_LIMIT_MAP_3),
		MARGIN_OPENORDERS_URL:                  base.NewReqUrlInfo(SAPI_V1_URL+"margin/openOrders", IP_1M_LIMIT_MAP),
		MYTRADES_URL:                           base.NewReqUrlInfo(API_V3_URL+"myTrades", IP_1M_LIMIT_MAP_10),
		TIME_URL:                               base.NewReqUrlInfo(API_V3_URL+"time", IP_1M_LIMIT_MAP),
		PING_URL:                               base.NewReqUrlInfo(API_V3_URL+"ping", IP_1M_LIMIT_MAP),
		TRADE_URL:                              base.NewReqUrlInfo(API_V3_URL+"trades", IP_1M_LIMIT_MAP),
		HISTORICAL_TRADES_URL:                  base.NewReqUrlInfo(API_V3_URL+"historicalTrades", IP_1M_LIMIT_MAP_5),
		EXCHANGEINFO_URL:                       base.NewReqUrlInfo(API_V3_URL+"exchangeInfo", IP_1M_LIMIT_MAP_10),
		CAPITAL_CONFIG_GETALL_URL:              base.NewReqUrlInfo(SAPI_V1_URL+"capital/config/getall", IP_1M_LIMIT_MAP_10),
		ASSET_TRADEFEE_URL:                     base.NewReqUrlInfo(SAPI_V1_URL+"asset/tradeFee", IP_1M_LIMIT_MAP),
		CAPITAL_WITHDRAW_APPLY_URL:             base.NewReqUrlInfo(SAPI_V1_URL+"capital/withdraw/apply", IP_1M_LIMIT_MAP),
		CAPITAL_DEPOSIT_HISREC_URL:             base.NewReqUrlInfo(SAPI_V1_URL+"capital/deposit/hisrec", IP_1M_LIMIT_MAP),
		ASSET_TRANSFER_URL:                     base.NewReqUrlInfo(SAPI_V1_URL+"asset/transfer", IP_1M_LIMIT_MAP),
		SUBACCOUNT_UNIVERSALTRANSFER_URL:       base.NewReqUrlInfo(SAPI_V1_URL+"sub-account/universalTransfer", IP_1M_LIMIT_MAP),
		SUBACCOUNT_TRANSFER_SUBTOMASTER_URL:    base.NewReqUrlInfo(SAPI_V1_URL+"sub-account/transfer/subToMaster", IP_1M_LIMIT_MAP),
		SUBACCOUNT_TRANSFER_SUBTOSUB_URL:       base.NewReqUrlInfo(SAPI_V1_URL+"sub-account/transfer/subToSub", IP_1M_LIMIT_MAP),
		SUBACCOUNT_TRANSFER_SUBUSERHISTORY_URL: base.NewReqUrlInfo(SAPI_V1_URL+"sub-account/transfer/subUserHistory", IP_1M_LIMIT_MAP),
		ACCOUNT_URL:                            base.NewReqUrlInfo(API_V3_URL+"account", IP_1M_LIMIT_MAP_10),
		ASSET_URL:                              base.NewReqUrlInfo(SAPI_V1_URL+"asset/get-funding-asset", IP_1M_LIMIT_MAP),
		MARGIN_ACCOUNT_URL:                     base.NewReqUrlInfo(SAPI_V1_URL+"margin/account", IP_1M_LIMIT_MAP_10),
		MARGIN_ISOLATED_ACCOUNT_URL:            base.NewReqUrlInfo(SAPI_V1_URL+"margin/isolated/account", IP_1M_LIMIT_MAP_10),
		ORDER_URL:                              base.NewReqUrlInfo(API_V3_URL+"order", ORDER_LIMIT_MAP...),
		ALLORDERS_URL:                          base.NewReqUrlInfo(API_V3_URL+"allOrders", IP_1M_LIMIT_MAP_10),
		MARGIN_ORDER_URL:                       base.NewReqUrlInfo(SAPI_V1_URL+"margin/order", IP_1M_LIMIT_MAP_6),
		CAPITAL_WITHDRAW_HISTORY_URL:           base.NewReqUrlInfo(SAPI_V1_URL+"capital/withdraw/history", IP_1M_LIMIT_MAP),
		MARGIN_LOAN_URL:                        base.NewReqUrlInfo(SAPI_V1_URL+"margin/loan", ORDER_LIMIT_MAP_3000...),
		MARGIN_REPAY_URL:                       base.NewReqUrlInfo(SAPI_V1_URL+"margin/repay", ORDER_LIMIT_MAP_3000...),
		MARGIN_ALLORDERS_URL:                   base.NewReqUrlInfo(SAPI_V1_URL+"margin/allOrders", IP_1M_LIMIT_MAP_200),
		USERDATASTREAM_URL:                     base.NewReqUrlInfo(API_V3_URL+"userDataStream", IP_1M_LIMIT_MAP),
		MARGINUSERDATASTREAM_URL:               base.NewReqUrlInfo(SAPI_V1_URL+"userDataStream", IP_1M_LIMIT_MAP),
		MARGINISOLATEDUSERDATASTREAM_URL:       base.NewReqUrlInfo(SAPI_V1_URL+"userDataStream/isolated", IP_1M_LIMIT_MAP),
		MARGIN_ISOLATED_PAIR_URL:               base.NewReqUrlInfo(SAPI_V1_URL+"margin/isolated/pair", IP_1M_LIMIT_MAP_10),
		MARGIN_ISOLATED_ALLPAIRS_URL:           base.NewReqUrlInfo(SAPI_V1_URL+"margin/isolated/allPairs", IP_1M_LIMIT_MAP_10),
		SUBACCOUNT_VIRTUALSUBACCOUNT_URL:       base.NewReqUrlInfo(SAPI_V1_URL+"sub-account/virtualSubAccount", IP_1M_LIMIT_MAP),
		SUBACCOUNT_LIST_URL:                    base.NewReqUrlInfo(SAPI_V1_URL+"sub-account/list", IP_1M_LIMIT_MAP),
	}
}

func NewUBaseReqUrl(types map[string]int64, account string) *ReqUrl {
	var (
		R_IP_1M_limit       int64 = 2400
		R_ACCOUNT_10S_limit int64 = 300
		R_ACCOUNT_1M_limit  int64 = 1200
	)
	if types != nil {
		if value, ok := types[R_IP_1M]; ok {
			R_IP_1M_limit = value
		}
		if value, ok := types[R_ACCOUNT_10S]; ok {
			R_ACCOUNT_10S_limit = value
		}
		if value, ok := types[R_ACCOUNT_1M]; ok {
			R_ACCOUNT_1M_limit = value
		}
	}

	initIpLimit(R_IP_1M_limit)
	ORDER_LIMIT_MAP := []*base.RateLimitConsume{base.NewRLConsume(R_IP_1M, 1, R_IP_1M_limit), base.NewRLConsume(base.RateLimitType(R_ACCOUNT_1M+"_"+account), 1, R_ACCOUNT_1M_limit), base.NewRLConsume(base.RateLimitType(R_ACCOUNT_10S+"_"+account), 1, R_ACCOUNT_10S_limit)}

	return &ReqUrl{
		DEPTH_URL:             base.NewReqUrlInfo(FAPI_V1_URL + "depth"),
		OPENORDERS_URL:        base.NewReqUrlInfo(FAPI_V1_URL+"openOrders", IP_1M_LIMIT_MAP),
		TIME_URL:              base.NewReqUrlInfo(FAPI_V1_URL+"time", IP_1M_LIMIT_MAP),
		PING_URL:              base.NewReqUrlInfo(FAPI_V1_URL+"ping", IP_1M_LIMIT_MAP),
		TRADE_URL:             base.NewReqUrlInfo(FAPI_V1_URL+"trades", IP_1M_LIMIT_MAP_5),
		HISTORICAL_TRADES_URL: base.NewReqUrlInfo(FAPI_V1_URL+"historicalTrades", IP_1M_LIMIT_MAP_20),
		EXCHANGEINFO_URL:      base.NewReqUrlInfo(FAPI_V1_URL+"exchangeInfo", IP_1M_LIMIT_MAP),
		ACCOUNT_URL:           base.NewReqUrlInfo(FAPI_V2_URL+"account", IP_1M_LIMIT_MAP_5),
		ORDER_URL:             base.NewReqUrlInfo(FAPI_V1_URL+"order", ORDER_LIMIT_MAP...),
		ALLORDERS_URL:         base.NewReqUrlInfo(FAPI_V1_URL+"allOrders", IP_1M_LIMIT_MAP_10),
		USERDATASTREAM_URL:    base.NewReqUrlInfo(FAPI_V1_URL+"listenKey", IP_1M_LIMIT_MAP),
		PREMIUMINDEX_URL:      base.NewReqUrlInfo(FAPI_V1_URL+"premiumIndex", IP_1M_LIMIT_MAP),
		TICKER_PRICE_URL:      base.NewReqUrlInfo(FAPI_V1_URL + "ticker/price"),
		POSITIONSIDE_DUAL_URL: base.NewReqUrlInfo(FAPI_V1_URL+"positionSide/dual", IP_1M_LIMIT_MAP),
		BATCHORDERS_URL:       base.NewReqUrlInfo(FAPI_V1_URL+"batchOrders", IP_1M_LIMIT_MAP_5),
		ALLOPENORDERS_URL:     base.NewReqUrlInfo(FAPI_V1_URL+"allOpenOrders", IP_1M_LIMIT_MAP),
		OPENORDER_URL:         base.NewReqUrlInfo(FAPI_V1_URL+"openOrder", IP_1M_LIMIT_MAP),
		BALANCE_URL:           base.NewReqUrlInfo(FAPI_V2_URL+"balance", IP_1M_LIMIT_MAP_5),
		LEVERAGE_URL:          base.NewReqUrlInfo(FAPI_V1_URL+"leverage", IP_1M_LIMIT_MAP),
		MARGINTYPE_URL:        base.NewReqUrlInfo(FAPI_V1_URL+"marginType", IP_1M_LIMIT_MAP),
		POSITIONMARGIN_URL:    base.NewReqUrlInfo(FAPI_V1_URL+"positionMargin", IP_1M_LIMIT_MAP),
		POSITIONRISK_URL:      base.NewReqUrlInfo(FAPI_V1_URL+"positionRisk", IP_1M_LIMIT_MAP_5),
		COMMISSIONRATE_URL:    base.NewReqUrlInfo(FAPI_V1_URL+"commissionRate", IP_1M_LIMIT_MAP_20),
		//ORDER_URL:                    FAPI_V1_URL + "order/test",  // 测试下单
	}
}

func NewCBaseReqUrl(types map[string]int64, account string) *ReqUrl {
	var (
		R_IP_1M_limit      int64 = 2400
		R_ACCOUNT_1M_limit int64 = 1200
	)
	if types != nil {
		if value, ok := types[R_IP_1M]; ok {
			R_IP_1M_limit = value
		}
		if value, ok := types[R_ACCOUNT_1M]; ok {
			R_ACCOUNT_1M_limit = value
		}
	}

	initIpLimit(R_IP_1M_limit)
	ORDER_LIMIT_MAP := []*base.RateLimitConsume{base.NewRLConsume(R_IP_1M, 1, R_IP_1M_limit), base.NewRLConsume(base.RateLimitType(R_ACCOUNT_1M+"_"+account), 1, R_ACCOUNT_1M_limit)}

	return &ReqUrl{
		DEPTH_URL:             base.NewReqUrlInfo(DAPI_V1_URL + "depth"),
		OPENORDERS_URL:        base.NewReqUrlInfo(DAPI_V1_URL+"openOrders", IP_1M_LIMIT_MAP),
		TIME_URL:              base.NewReqUrlInfo(DAPI_V1_URL+"time", IP_1M_LIMIT_MAP),
		PING_URL:              base.NewReqUrlInfo(DAPI_V1_URL+"ping", IP_1M_LIMIT_MAP),
		TRADE_URL:             base.NewReqUrlInfo(DAPI_V1_URL+"trades", IP_1M_LIMIT_MAP_5),
		HISTORICAL_TRADES_URL: base.NewReqUrlInfo(DAPI_V1_URL+"historicalTrades", IP_1M_LIMIT_MAP_20),
		EXCHANGEINFO_URL:      base.NewReqUrlInfo(DAPI_V1_URL+"exchangeInfo", IP_1M_LIMIT_MAP),
		ACCOUNT_URL:           base.NewReqUrlInfo(DAPI_V1_URL+"account", IP_1M_LIMIT_MAP_5),
		ORDER_URL:             base.NewReqUrlInfo(DAPI_V1_URL+"order", ORDER_LIMIT_MAP...),
		ALLORDERS_URL:         base.NewReqUrlInfo(DAPI_V1_URL + "allOrders"),
		USERDATASTREAM_URL:    base.NewReqUrlInfo(DAPI_V1_URL+"listenKey", IP_1M_LIMIT_MAP),
		PREMIUMINDEX_URL:      base.NewReqUrlInfo(DAPI_V1_URL+"premiumIndex", IP_1M_LIMIT_MAP_10),
		TICKER_PRICE_URL:      base.NewReqUrlInfo(DAPI_V1_URL + "ticker/price"),
		POSITIONSIDE_DUAL_URL: base.NewReqUrlInfo(DAPI_V1_URL+"positionSide/dual", IP_1M_LIMIT_MAP),
		BATCHORDERS_URL:       base.NewReqUrlInfo(DAPI_V1_URL+"batchOrders", IP_1M_LIMIT_MAP_5),
		ALLOPENORDERS_URL:     base.NewReqUrlInfo(DAPI_V1_URL+"allOpenOrders", IP_1M_LIMIT_MAP),
		OPENORDER_URL:         base.NewReqUrlInfo(DAPI_V1_URL+"openOrder", IP_1M_LIMIT_MAP),
		BALANCE_URL:           base.NewReqUrlInfo(DAPI_V1_URL+"balance", IP_1M_LIMIT_MAP),
		LEVERAGE_URL:          base.NewReqUrlInfo(DAPI_V1_URL+"leverage", IP_1M_LIMIT_MAP),
		MARGINTYPE_URL:        base.NewReqUrlInfo(DAPI_V1_URL+"marginType", IP_1M_LIMIT_MAP),
		POSITIONMARGIN_URL:    base.NewReqUrlInfo(DAPI_V1_URL+"positionMargin", IP_1M_LIMIT_MAP),
		POSITIONRISK_URL:      base.NewReqUrlInfo(DAPI_V1_URL+"positionRisk", IP_1M_LIMIT_MAP),
		COMMISSIONRATE_URL:    base.NewReqUrlInfo(DAPI_V1_URL+"commissionRate", IP_1M_LIMIT_MAP_20),
		//ORDER_URL:                    DAPI_V1_URL + "order/test",  // 测试下单
	}
}

func initIpLimit(R_IP_1M_limit int64) {
	mtx.Lock()
	defer mtx.Unlock()
	oldLimit := int64(0)
	if IP_1M_LIMIT_MAP != nil {
		oldLimit = IP_1M_LIMIT_MAP.Limit
	}

	if oldLimit >= R_IP_1M_limit {
		return
	}
	//  IP限制使用所有实例中最大的limit
	IP_1M_LIMIT_MAP = base.NewRLConsume(R_IP_1M, 1, R_IP_1M_limit)
	IP_1M_LIMIT_MAP_2 = base.NewRLConsume(R_IP_1M, 2, R_IP_1M_limit)
	IP_1M_LIMIT_MAP_3 = base.NewRLConsume(R_IP_1M, 3, R_IP_1M_limit)
	IP_1M_LIMIT_MAP_5 = base.NewRLConsume(R_IP_1M, 5, R_IP_1M_limit)
	IP_1M_LIMIT_MAP_6 = base.NewRLConsume(R_IP_1M, 6, R_IP_1M_limit)
	IP_1M_LIMIT_MAP_10 = base.NewRLConsume(R_IP_1M, 10, R_IP_1M_limit)
	IP_1M_LIMIT_MAP_20 = base.NewRLConsume(R_IP_1M, 20, R_IP_1M_limit)
	IP_1M_LIMIT_MAP_40 = base.NewRLConsume(R_IP_1M, 40, R_IP_1M_limit)
	IP_1M_LIMIT_MAP_50 = base.NewRLConsume(R_IP_1M, 50, R_IP_1M_limit)
	IP_1M_LIMIT_MAP_200 = base.NewRLConsume(R_IP_1M, 200, R_IP_1M_limit)
}
