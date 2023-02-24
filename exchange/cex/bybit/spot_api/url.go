package spot_api

import (
	"clients/exchange/cex/base"
	"net/url"
	"strings"
)

var (
	R_IP_1S            base.RateLimitType = "ip_1s"
	IP_1S_LIMIT_MAP_50 *base.RateLimitConsume
	IP_1S_LIMIT_MAP_20 *base.RateLimitConsume
)

const (
	SPOT_API_BASE_URL    = "https://api.bybit.com"
	WS_API_IP_URL        = "wss://stream.bybit.com/realtime"
	WS_API_Public_Topic1 = "wss://stream.bybit.com/spot/quote/ws/v1"
	WS_API_Public_Topic2 = "wss://stream.bybit.com/spot/quote/ws/v2"
)

type ReqUrl struct {
	SYMBOLS_URL   base.ReqUrlInfo
	ORDERBOOK_URL base.ReqUrlInfo
	TIME_URL      base.ReqUrlInfo

	// todo 没有接口
	FUNDFEE_URL     base.ReqUrlInfo
	TRANSFERFEE_URL base.ReqUrlInfo

	ACCOUNT_URL base.ReqUrlInfo

	PLACEORDER_URL   base.ReqUrlInfo
	CANCLEORDER_URL  base.ReqUrlInfo
	ORDERINFO_URL    base.ReqUrlInfo
	ORDERINFO2_URL   base.ReqUrlInfo
	ORDERHISTORY_URL base.ReqUrlInfo

	WITHDRAW_URL        base.ReqUrlInfo
	WITHDRAWFEE_URL     base.ReqUrlInfo
	WITHDRAWHISTORY_URL base.ReqUrlInfo

	TRANSFER_URL           base.ReqUrlInfo
	TRANSFER_M2S_URL       base.ReqUrlInfo
	ALLTRANSFER_URL        base.ReqUrlInfo
	TRANSFERHISTORY_URL    base.ReqUrlInfo
	TRANSFERHISTORYM2S_URL base.ReqUrlInfo

	RECORDHISTORY_URL base.ReqUrlInfo

	LOAN_URL         base.ReqUrlInfo
	LOANHISTORY_URL  base.ReqUrlInfo
	REPAY_URL        base.ReqUrlInfo
	REPAYHISTORY_URL base.ReqUrlInfo

	FUTURE_SYMBOL_URL,
	FUTURE_ORDERBOOK_URL,
	FUTURE_MARKET_URL,
	FUTURE_POSITION_URL,
	FUTURE_PLACE_COIN_URL,
	FUTURE_PLACE_URL,
	FUTURE_PLACE_COIN_DELIVERY_URL base.ReqUrlInfo
}

func NewIpLimit(url, duration string, limit int64, optionKeys ...string) base.ReqUrlInfo {
	return base.ReqUrlInfo{
		Url: SPOT_API_BASE_URL + url,
		RateLimitConsumeMap: []*base.RateLimitConsume{
			{Count: 1, Limit: limit, LimitTypeName: base.RateLimitType(strings.Join([]string{"ip", duration, url}, "_"))},
		},
		OptionKeys: optionKeys,
	}
}

func NewGetIpLimit(url string) base.ReqUrlInfo {
	return NewIpLimit(url, "1s", 50)
}

func (r *ReqUrl) GetURLRateLimit(url base.ReqUrlInfo, params url.Values) []*base.RateLimitConsume {
	rates := url.RateLimitConsumeMap
	if len(rates) == 0 {
		return rates
	}
	if len(url.OptionKeys) > 0 {
		rateName := rates[0].LimitTypeName
		for _, item := range url.OptionKeys {
			rateName += base.RateLimitType(params.Get(item))
		}
		return []*base.RateLimitConsume{
			{Count: 1, Limit: rates[0].Limit, LimitTypeName: rateName},
		}
	}
	return rates
}

// 通过同一个名字进行绑定
func NewSpotReqUrl() *ReqUrl {
	var R_IPGET_1S_LIMIT_GET int64 = 50
	var R_IPGET_1S_LIMIT_POST int64 = 20
	IP_1S_LIMIT_MAP_50 = base.NewRLConsume(R_IP_1S, 1, R_IPGET_1S_LIMIT_GET)
	IP_1S_LIMIT_MAP_20 = base.NewRLConsume(R_IP_1S, 1, R_IPGET_1S_LIMIT_POST)
	return &ReqUrl{
		SYMBOLS_URL:                    base.NewReqUrlInfo("/spot/v1/symbols", IP_1S_LIMIT_MAP_50),
		ORDERBOOK_URL:                  base.NewReqUrlInfo("/spot/v3/public/quote/depth", IP_1S_LIMIT_MAP_50),
		TIME_URL:                       base.NewReqUrlInfo("/v3/public/time", IP_1S_LIMIT_MAP_50),
		FUNDFEE_URL:                    base.NewReqUrlInfo("/asset/v1/private/coin-info/query", IP_1S_LIMIT_MAP_50),
		ACCOUNT_URL:                    base.NewReqUrlInfo("/spot/v3/private/account", IP_1S_LIMIT_MAP_50),
		PLACEORDER_URL:                 base.NewReqUrlInfo("/spot/v3/private/order", IP_1S_LIMIT_MAP_20),
		CANCLEORDER_URL:                base.NewReqUrlInfo("/spot/v3/private/cancel-order", IP_1S_LIMIT_MAP_20),
		ORDERINFO_URL:                  base.NewReqUrlInfo("/spot/v3/private/order", IP_1S_LIMIT_MAP_50),
		ORDERINFO2_URL:                 base.NewReqUrlInfo("/spot/v3/private/open-orders", IP_1S_LIMIT_MAP_50),
		ORDERHISTORY_URL:               base.NewReqUrlInfo("/spot/v3/private/history-orders", IP_1S_LIMIT_MAP_50),
		WITHDRAW_URL:                   base.NewReqUrlInfo("/asset/v1/private/withdraw", IP_1S_LIMIT_MAP_20),
		WITHDRAWFEE_URL:                base.NewReqUrlInfo("/asset/v1/private/coin-info/query", IP_1S_LIMIT_MAP_50),
		WITHDRAWHISTORY_URL:            base.NewReqUrlInfo("/asset/v1/private/withdraw/record/query", IP_1S_LIMIT_MAP_50),
		TRANSFER_URL:                   base.NewReqUrlInfo("/asset/v1/private/transfer", IP_1S_LIMIT_MAP_20),
		TRANSFER_M2S_URL:               base.NewReqUrlInfo("/asset/v1/private/sub-member/transfer", IP_1S_LIMIT_MAP_20),
		ALLTRANSFER_URL:                base.NewReqUrlInfo("/asset/v1/private/universal/transfer", IP_1S_LIMIT_MAP_20),
		TRANSFERHISTORY_URL:            base.NewReqUrlInfo("/asset/v1/private/transfer/list", IP_1S_LIMIT_MAP_50),
		TRANSFERHISTORYM2S_URL:         base.NewReqUrlInfo("/asset/v1/private/sub-member/transfer/list", IP_1S_LIMIT_MAP_50),
		RECORDHISTORY_URL:              base.NewReqUrlInfo("/asset/v1/private/deposit/record/query", IP_1S_LIMIT_MAP_50),
		LOAN_URL:                       base.NewReqUrlInfo("/spot/v3/private/cross-margin-loan", IP_1S_LIMIT_MAP_20),
		LOANHISTORY_URL:                base.NewReqUrlInfo("/spot/v3/private/cross-margin-orders", IP_1S_LIMIT_MAP_50),
		REPAY_URL:                      base.NewReqUrlInfo("/spot/v3/private/cross-margin-repay", IP_1S_LIMIT_MAP_20),
		REPAYHISTORY_URL:               base.NewReqUrlInfo("/spot/v3/private/cross-margin-repay-history", IP_1S_LIMIT_MAP_50),
		FUTURE_SYMBOL_URL:              base.NewReqUrlInfo("/v2/public/symbols", IP_1S_LIMIT_MAP_50),
		FUTURE_ORDERBOOK_URL:           base.NewReqUrlInfo("/v2/public/orderBook/L2", IP_1S_LIMIT_MAP_50),
		FUTURE_MARKET_URL:              base.NewReqUrlInfo("/public/linear/mark-price-kline", IP_1S_LIMIT_MAP_50),
		FUTURE_POSITION_URL:            base.NewReqUrlInfo("/private/linear/position/list", IP_1S_LIMIT_MAP_50),
		FUTURE_PLACE_COIN_DELIVERY_URL: base.NewReqUrlInfo("/futures/private/order/create", IP_1S_LIMIT_MAP_50),
		FUTURE_PLACE_COIN_URL:          base.NewReqUrlInfo("/v2/private/order/create", IP_1S_LIMIT_MAP_50),
		FUTURE_PLACE_URL:               base.NewReqUrlInfo("/private/linear/order/create", IP_1S_LIMIT_MAP_50),
	}
}

type WsReqUrl struct {
}

func NewSpotWsUrl() *WsReqUrl {
	return &WsReqUrl{}
}
