package spot_api

import (
	"clients/exchange/cex/base"
	"net/url"
	"strings"
)

const (
	SPOT_API_BASE_URL    = "https://api.kucoin.com"
	FUTURUE_API_BASE_URL = "https://api-futures.kucoin.com"
	SANDBOX_BASE_URL     = "https://openapi-sandbox.kucoin.com"
	API_V1               = "/api/v1/"
	API_V2               = "/api/v2/"
	API_V3               = "/api/v3/"
)

type ReqUrl struct {
	SYMBOLS_URL,
	TRADEFEE_URL,
	MARKET_LV2_20_URL,
	MARKET_LV2_100_URL,
	MARKET_LV2_URL,
	CURRENCIES_URL,
	ACCOUNTS_URL,
	MARGIN_ACCOUNTS_URL,
	ISOLATED_ACCOUNTS_URL,
	ORDERS_URL,
	ORDER_CLIENT_URL,
	FILLS_URL,
	MARGIN_ORDER_URL,
	STATUS_URL,
	ACCOUNTS_INNER_TRANSFER_URL,
	ACCOUNTS_SUB_TRANSFER_URL,
	WITHDRAW_URL,
	WITHDRAW_FEE_URL,
	WITHDRAW_HIS_URL,
	HIST_WITHDRAW_URL,
	// ws
	TOKEN_PUBLIC_URL,
	TOKEN_PRIVATE_URL,
	ORDER_HISTORY_URL,
	MOVE_HISTORY_URL,
	LOAN_URL,
	UNREPAY_HISTORY_URL,
	REPAY_URL,
	REPAY_HISTORY_URL,
	CONTRACT_SYMBOL_URL,
	CONTRACT_DEPTH_URL,
	CONTRAC_SPECIFY_SYMBOL_URL,
	CONTRACT_POSITIONS_URL,
	CONTRACT_PALCE_ORDER_URL base.ReqUrlInfo
}

//Check Order history
// 交易所没有频率限制的API，这里自己定义 1s 10次
func NewSpotReqUrl(account string) *ReqUrl {
	return &ReqUrl{
		SYMBOLS_URL:                 NewAccountLimit(API_V1+"symbols", account, "1s", 10),
		TRADEFEE_URL:                NewAccountLimit(API_V1+"trade-fees", account, "1s", 10),
		MARKET_LV2_20_URL:           NewAccountLimit(API_V1+"market/orderbook/level2_20", account, "1s", 10),
		MARKET_LV2_100_URL:          NewAccountLimit(API_V1+"market/orderbook/level2_100", account, "1s", 10),
		MARKET_LV2_URL:              NewAccountLimit(API_V3+"market/orderbook/level2", account, "3s", 30),
		CURRENCIES_URL:              NewAccountLimit(API_V1+"currencies", account, "1s", 10),
		ACCOUNTS_URL:                NewAccountLimit(API_V1+"accounts", account, "1s", 10),
		MARGIN_ACCOUNTS_URL:         NewAccountLimit(API_V1+"margin/account", account, "1s", 10),
		ISOLATED_ACCOUNTS_URL:       NewAccountLimit(API_V1+"isolated/accounts", account, "1s", 10),
		ORDERS_URL:                  NewAccountLimit(API_V1+"orders", account, "1s", 10),
		ORDER_CLIENT_URL:            NewAccountLimit(API_V1+"order/client-order", account, "1s", 10),
		FILLS_URL:                   NewAccountLimit(API_V1+"fills", account, "3s", 9),
		MARGIN_ORDER_URL:            NewAccountLimit(API_V1+"margin/order", account, "3s", 45),
		STATUS_URL:                  NewAccountLimit(API_V1+"status", account, "1s", 10),
		ACCOUNTS_INNER_TRANSFER_URL: NewAccountLimit(API_V2+"accounts/inner-transfer", account, "1s", 10),
		ACCOUNTS_SUB_TRANSFER_URL:   NewAccountLimit(API_V2+"accounts/sub-transfer", account, "3s", 3),
		TOKEN_PUBLIC_URL:            NewAccountLimit(API_V1+"bullet-public", account, "1s", 10),
		TOKEN_PRIVATE_URL:           NewAccountLimit(API_V1+"bullet-private", account, "1s", 10),
		WITHDRAW_URL:                NewAccountLimit(API_V1+"withdrawals", account, "3s", 6),
		WITHDRAW_FEE_URL:            NewAccountLimit(API_V1+"withdrawals/quotas", account, "3s", 10),
		WITHDRAW_HIS_URL:            NewAccountLimit(API_V1+"hist-withdrawals", account, "3s", 6),
		ORDER_HISTORY_URL:           NewAccountLimit(API_V1+"orders", account, "3s", 30),
		MOVE_HISTORY_URL:            NewAccountLimit(API_V1+"accounts/ledgers", account, "3s", 18),
		LOAN_URL:                    NewAccountLimit(API_V1+"margin/borrow", account, "1s", 10),
		UNREPAY_HISTORY_URL:         NewAccountLimit(API_V1+"margin/borrow/outstanding", account, "1s", 10),
		REPAY_URL:                   NewAccountLimit(API_V1+"margin/repay/single", account, "1s", 10),
		REPAY_HISTORY_URL:           NewAccountLimit(API_V1+"margin/borrow/repaid", account, "1s", 10),
		CONTRACT_SYMBOL_URL:         NewAccountLimit(API_V1+"contracts/active", account, "1s", 10),
		CONTRACT_DEPTH_URL:          NewAccountLimit(API_V1+"level2/depth", account, "1s", 10),
		CONTRAC_SPECIFY_SYMBOL_URL:  NewAccountLimit(API_V1+"contracts/", account, "1s", 10),
		CONTRACT_POSITIONS_URL:      NewAccountLimit(API_V1+"positions", account, "1s", 10),
		CONTRACT_PALCE_ORDER_URL:    NewAccountLimit(API_V1+"orders", account, "3s", 30),
	}
}

func NewAccountLimit(url, account, duration string, limit int64, optionKeys ...string) base.ReqUrlInfo {
	return base.ReqUrlInfo{
		Url: url,
		RateLimitConsumeMap: []*base.RateLimitConsume{
			{1, limit, base.RateLimitType(strings.Join([]string{"account", duration, url, account}, "_"))},
		},
		OptionKeys: optionKeys,
	}
}

type WsReqUrl struct {
	WS_BASE_URL                   string
	MARKET_MATCH_TOPIC            string
	MARKET_TICKER_TOPIC           string
	SPOTMARKET_LEVEL2_DEPTH_TOPIC string
	MARKET_LEVEL2                 string
}

func NewWsReqUrl() *WsReqUrl {
	return &WsReqUrl{
		WS_BASE_URL:                   "",
		MARKET_MATCH_TOPIC:            "/market/match",
		MARKET_TICKER_TOPIC:           "/market/ticker",
		SPOTMARKET_LEVEL2_DEPTH_TOPIC: "/spotMarket/level2Depth50",
		MARKET_LEVEL2:                 "/market/level2",
	}
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
