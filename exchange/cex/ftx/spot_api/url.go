package spot_api

import (
	"clients/exchange/cex/base"
	"strings"
)

const (
	SPOT_API_BASE_URL   = "https://ftx.com/api/"
	FUTURE_API_BASE_URL = "https://ftx.com/api/"
	WS_API_BASE_URL     = "wss://ftx.com/ws/"
)

type ReqUrl struct {
	MARKETS_URL,
	ACCOUNT_URL,
	ORDERS_URL,
	ORDERS_BY_CLIENT_ID_URL,
	ORDERS_HISTORY_URL,
	FILLS_URL,
	FUTURES_URL,
	POSITIONS_URL,
	SPOT_MARGIN_HISTORY,
	SPOT_MARGIN_BORROW_RATES,
	SPOT_MARGIN_LENDING_RATES,
	SPOT_MARGIN_MARKET_INFO,
	SPOT_MARGIN_OFFERS,
	SUBACCOUNTS_TRANSFER,
	WALLET_BALANCES_URL,
	WALLET_DEPOSIT_ADDRESS_URL,
	WALLET_WITHDRAWALS_URL,
	WALLET_DEPOSITS_URL,
	WALLET_WITHDRAWAL_FEE_URL,
	FUNDING_RATES_URL base.ReqUrlInfo
}

func NewRateLimit(url string) base.ReqUrlInfo {
	return base.ReqUrlInfo{
		Url: url,
	}
}

func (r *ReqUrl) GetURLRateLimit(url base.ReqUrlInfo, params map[string]interface{}) []*base.RateLimitConsume {
	var (
		rates      = url.RateLimitConsumeMap
		symbolName interface{}
		symbol     string
		ok         bool
	)
	if len(rates) == 0 {
		return rates
	}
	if symbolName, ok = params["market"]; ok {
		symbol, ok = symbolName.(string)
	}
	symbol = strings.ToUpper(symbol)
	if !ok {
		return []*base.RateLimitConsume{}
	}
	if strings.HasPrefix(symbol, "BTC") && strings.HasSuffix(symbol, "PERP") {
		return []*base.RateLimitConsume{
			url.RateLimitConsumeMap[0],
			url.RateLimitConsumeMap[1],
			url.RateLimitConsumeMap[4],
			url.RateLimitConsumeMap[5],
		}
	} else if strings.HasPrefix(symbol, "ETH") && strings.HasSuffix(symbol, "PERP") {
		return []*base.RateLimitConsume{
			url.RateLimitConsumeMap[0],
			url.RateLimitConsumeMap[2],
			url.RateLimitConsumeMap[4],
			url.RateLimitConsumeMap[6],
		}
	} else {
		return []*base.RateLimitConsume{
			url.RateLimitConsumeMap[0],
			{
				Count:         url.RateLimitConsumeMap[3].Count,
				Limit:         url.RateLimitConsumeMap[3].Limit,
				LimitTypeName: base.RateLimitType(string(url.RateLimitConsumeMap[3].LimitTypeName) + "_" + symbol),
			},
			url.RateLimitConsumeMap[4],
			{
				Count:         url.RateLimitConsumeMap[7].Count,
				Limit:         url.RateLimitConsumeMap[7].Limit,
				LimitTypeName: base.RateLimitType(string(url.RateLimitConsumeMap[7].LimitTypeName) + "_" + symbol),
			},
		}
	}
}

// Check Order history
func NewSpotReqUrl(types map[string]int64, account string) *ReqUrl {
	var (
		account1s         int64 = 12
		account1sETH      int64 = 7
		account1sBTC      int64 = 7
		account1sOther    int64 = 7
		account200ms      int64 = 7
		account200msETH   int64 = 7
		account200msBTC   int64 = 7
		account200msOther int64 = 7
	)

	for wType, limit := range types {
		switch base.RateLimitType(wType) {
		case WEIGHT_TYPE_Account1s:
			account1s = limit
		case WEIGHT_TYPE_Account1sEth:
			account1sETH = limit
		case WEIGHT_TYPE_Account1sBtc:
			account1sBTC = limit
		case WEIGHT_TYPE_Account1sOther:
			account1sOther = limit
		case WEIGHT_TYPE_Account200ms:
			account200ms = limit
		case WEIGHT_TYPE_Account200msEth:
			account200msETH = limit
		case WEIGHT_TYPE_Account200msBtc:
			account200msBTC = limit
		case WEIGHT_TYPE_Account200msOther:
			account200msOther = limit
		}
	}

	orderConsumeMap := []*base.RateLimitConsume{
		{Count: 1, Limit: account1s, LimitTypeName: base.RateLimitType(string(WEIGHT_TYPE_Account1s) + "_" + account)},
		{Count: 1, Limit: account1sETH, LimitTypeName: base.RateLimitType(string(WEIGHT_TYPE_Account1sEth) + "_" + account)},
		{Count: 1, Limit: account1sBTC, LimitTypeName: base.RateLimitType(string(WEIGHT_TYPE_Account1sBtc) + "_" + account)},
		{Count: 1, Limit: account1sOther, LimitTypeName: base.RateLimitType(string(WEIGHT_TYPE_Account1sOther) + "_" + account)},
		{Count: 1, Limit: account200ms, LimitTypeName: base.RateLimitType(string(WEIGHT_TYPE_Account200ms) + "_" + account)},
		{Count: 1, Limit: account200msETH, LimitTypeName: base.RateLimitType(string(WEIGHT_TYPE_Account200msEth) + "_" + account)},
		{Count: 1, Limit: account200msBTC, LimitTypeName: base.RateLimitType(string(WEIGHT_TYPE_Account200msBtc) + "_" + account)},
		{Count: 1, Limit: account200msOther, LimitTypeName: base.RateLimitType(string(WEIGHT_TYPE_Account200msOther) + "_" + account)},
	}

	return &ReqUrl{
		MARKETS_URL:                NewRateLimit("markets"),
		ACCOUNT_URL:                NewRateLimit("account"),
		ORDERS_URL:                 base.ReqUrlInfo{Url: "orders", RateLimitConsumeMap: orderConsumeMap},
		ORDERS_BY_CLIENT_ID_URL:    NewRateLimit("orders/by_client_id"),
		ORDERS_HISTORY_URL:         NewRateLimit("orders/history"),
		FILLS_URL:                  NewRateLimit("fills"),
		POSITIONS_URL:              NewRateLimit("positions"),
		SPOT_MARGIN_HISTORY:        NewRateLimit("spot_margin/history"),
		SPOT_MARGIN_BORROW_RATES:   NewRateLimit("spot_margin/borrow_rates"),
		SPOT_MARGIN_LENDING_RATES:  NewRateLimit("spot_margin/lending_rates"),
		SPOT_MARGIN_MARKET_INFO:    NewRateLimit("spot_margin/market_info"),
		SPOT_MARGIN_OFFERS:         NewRateLimit("spot_margin/offers"),
		SUBACCOUNTS_TRANSFER:       NewRateLimit("subaccounts/transfer"),
		WALLET_BALANCES_URL:        NewRateLimit("wallet/balances"),
		WALLET_DEPOSITS_URL:        NewRateLimit("wallet/deposits"),
		WALLET_DEPOSIT_ADDRESS_URL: NewRateLimit("wallet/deposit_address"),
		WALLET_WITHDRAWALS_URL:     NewRateLimit("wallet/withdrawals"),
		WALLET_WITHDRAWAL_FEE_URL:  NewRateLimit("wallet/withdrawal_fee"),
	}
}

func NewUBaseReqUrl() *ReqUrl {
	return &ReqUrl{
		FUTURES_URL:       NewRateLimit("futures"),
		MARKETS_URL:       NewRateLimit("markets"),
		ACCOUNT_URL:       NewRateLimit("account"),
		FUNDING_RATES_URL: NewRateLimit("funding_rates"),
	}
}
