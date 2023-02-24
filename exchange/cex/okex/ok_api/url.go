package ok_api

import (
	"clients/exchange/cex/base"
	"net/url"
	"strings"
)

const (
	GLOBAL_API_BASE_URL = "https://aws.okex.com"
	PUBLIC_API_V5_URL   = "/api/v5"
)

type ReqUrl struct {
	PUBLIC_MARK_PRICE,
	PUBLIC_SERVER_TIME,
	INSTRUMENT_URL,
	DELIVERY_EXERCISE_HISTORY_URL,
	OPEN_INTEREST_URL,
	MARKET_BOOKS,
	MARKET_TRADES,
	MARKET_INDEX_TICKERS,
	MARKET_HISTORY_TRADES,
	ACCOUNT_TRADE_INFO,
	ACCOUNT_BALANCE_INFO,
	ASSET_BALANCES_INFO,
	ACCOUNT_POSITION,
	ACCOUNT_POSITION_HISTORY,
	ACCOUNT_ACCOUNT_ACCOUNT_POSITION_RISK,
	ACCOUNT_ACCOUNT_BILLS,
	ACCOUNT_ACCOUNT_BILLS_ARCHIVE,
	ACCOUNT_CONFIG,
	ACCOUNT_SET_POSITION_MODE,
	ACCOUNT_SET_LEVERAGE,
	ACCOUNT_MAX_SIZE,
	ACCOUNT_MAX_AVAIL_SIZE,
	ACCOUNT_POSITION_MATGIN_BALANCE,
	ACCOUNT_LEVERAGE_INFO,
	ACCOUNT_MAX_LOAN,
	ACCOUNT_TRADE_FEE,
	ACCOUNT_INTEREST_ACCRUED,
	ACCOUNT_INTEREST_RATE,
	ACCOUNT_SET_GREEKS,
	ACCOUNT_SET_ISOLATED_MODE,
	ACCOUNT_MAX_WITHDRAWAL,
	ACCOUNT_RISK_STATE,
	ACCOUNT_BORROW_REPAY,
	ACCOUNT_BORROW_REPAY_HISTORY,
	ACCOUNT_INTEREST_LIMITS,
	ACCOUNT_SIMULATED_MARGIN,
	ACCOUNT_GREEKS,
	ACCOUNT_POSITION_TIERS,
	TRAD_ORDER,
	TRAD_BATCH_ORDERS,
	TRAD_CANCEL_ORDER,
	TRAD_CANCEL_BATCH_ORDERS,
	TRAD_CLOSE_POSITION,
	TRAD_ORDERS_PENDING,
	TRAD_ORDERS_HISTORY_WEEK,
	TRAD_ORDERS_HISTORY_ARCHIVE,
	TRAD_FILLS_HISTORY_THREE_DAYS,
	TRAD_FILLS_HISTORY_ARCHIVE,
	TRAD_ORDER_ALGO,
	TRAD_CANCEL_ALGO,
	TRAD_CANCLE_ADVANCE_ALGOS,
	TRAD_ORDERS_ALGO_PENDING,
	TRAD_ORDERS_ALGO_HISTORY,
	RFQ_COUNTERPARTIES,
	ASSET_CURRENCIES,
	ASSET_TRANSFER,
	ASSET_TRANSFER_STATE,
	ASSET_BILLS,
	ASSET_WITHDRAWAL,
	ASSET_WITHDRAWAL_LIGHTNING,
	ASSET_CANCEL_WITHDRAWAL,
	ASSET_WITHDRAWAL_HISTORY,
	ASSET_DEPOSIT_HISTORY,
	ASSET_SUBACCOUNT_BILLS,
	USER_SET_TRANSGER_OUT base.ReqUrlInfo
}

func NewIpLimit(url, duration string, limit int64, optionKeys ...string) base.ReqUrlInfo {
	return base.ReqUrlInfo{
		Url: PUBLIC_API_V5_URL + url,
		RateLimitConsumeMap: []*base.RateLimitConsume{
			{Count: 1, Limit: limit, LimitTypeName: base.RateLimitType(strings.Join([]string{"ip", duration, url}, "_"))},
		},
		OptionKeys: optionKeys,
	}
}

func NewAccountLimit(url, account, duration string, limit int64, optionKeys ...string) base.ReqUrlInfo {
	return base.ReqUrlInfo{
		Url: PUBLIC_API_V5_URL + url,
		RateLimitConsumeMap: []*base.RateLimitConsume{
			{Count: 1, Limit: limit, LimitTypeName: base.RateLimitType(strings.Join([]string{"account", duration, url, account}, "_"))},
		},
		OptionKeys: optionKeys,
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

func NewSpotReqUrl(types map[string]int64, account string) *ReqUrl {
	return &ReqUrl{
		PUBLIC_MARK_PRICE:                     NewIpLimit("/public/mark-price", "2s", 10, "instId"),
		PUBLIC_SERVER_TIME:                    NewIpLimit("/public/time", "2s", 10),
		INSTRUMENT_URL:                        NewIpLimit("/public/instruments", "11s", 10, "instType"),
		DELIVERY_EXERCISE_HISTORY_URL:         NewIpLimit("/public/delivery-exercise-history", "2s", 40, "instType", "uly"),
		OPEN_INTEREST_URL:                     NewIpLimit("/public/open-interest", "2s", 20, "instId"),
		MARKET_BOOKS:                          NewIpLimit("/market/books", "2s", 20),
		MARKET_TRADES:                         NewIpLimit("/market/trades", "2s", 20),
		MARKET_INDEX_TICKERS:                  NewIpLimit("/market/index-tickers", "2s", 20),
		MARKET_HISTORY_TRADES:                 NewIpLimit("/market/history-trades", "2s", 10),
		ACCOUNT_TRADE_INFO:                    NewAccountLimit("/account/trade-fee", account, "2s", 5),
		ACCOUNT_BALANCE_INFO:                  NewAccountLimit("/account/balance", account, "2s", 10),
		ASSET_BALANCES_INFO:                   NewAccountLimit("/asset/balances", account, "2s", 10),
		ACCOUNT_POSITION:                      NewAccountLimit("/account/positions", account, "2s", 10),
		ACCOUNT_POSITION_HISTORY:              NewAccountLimit("/account/positions-history", account, "10s", 1),
		ACCOUNT_ACCOUNT_ACCOUNT_POSITION_RISK: NewAccountLimit("/account/account-position-risk", account, "2s", 10),
		ACCOUNT_ACCOUNT_BILLS:                 NewAccountLimit("/account/bills", account, "2s", 10),
		ACCOUNT_ACCOUNT_BILLS_ARCHIVE:         NewAccountLimit("/account/bills-archive", account, "2s", 5),
		ACCOUNT_CONFIG:                        NewAccountLimit("/account/config", account, "2s", 5),
		ACCOUNT_SET_POSITION_MODE:             NewAccountLimit("/account/set-position-mode", account, "2s", 5),
		ACCOUNT_SET_LEVERAGE:                  NewAccountLimit("/account/set-leverage", account, "2s", 20),
		ACCOUNT_MAX_SIZE:                      NewAccountLimit("/account/max-size", account, "2s", 20),
		ACCOUNT_MAX_AVAIL_SIZE:                NewAccountLimit("/account/max-avail-size", account, "2s", 20),
		ACCOUNT_POSITION_MATGIN_BALANCE:       NewAccountLimit("/account/position/margin-balance", account, "2s", 20),
		ACCOUNT_LEVERAGE_INFO:                 NewAccountLimit("/account/leverage-info", account, "2s", 20),
		ACCOUNT_MAX_LOAN:                      NewAccountLimit("/account/max-loan", account, "2s", 20),
		ACCOUNT_TRADE_FEE:                     NewAccountLimit("/account/trade-fee", account, "2s", 5),
		ACCOUNT_INTEREST_ACCRUED:              NewAccountLimit("/account/interest-accrued", account, "2s", 5),
		ACCOUNT_INTEREST_RATE:                 NewAccountLimit("/account/interest-rate", account, "2s", 5),
		ACCOUNT_SET_GREEKS:                    NewAccountLimit("/account/set-greeks", account, "2s", 5),
		ACCOUNT_SET_ISOLATED_MODE:             NewAccountLimit("/account/set-isolated-mode", account, "2s", 5),
		ACCOUNT_MAX_WITHDRAWAL:                NewAccountLimit("/account/max-withdrawal", account, "2s", 20),
		ACCOUNT_RISK_STATE:                    NewAccountLimit("/account/risk-state", account, "2s", 10),
		ACCOUNT_BORROW_REPAY:                  NewAccountLimit("/account/borrow-repay", account, "1s", 6),
		ACCOUNT_BORROW_REPAY_HISTORY:          NewAccountLimit("/account/borrow-repay-history", account, "2s", 5),
		ACCOUNT_INTEREST_LIMITS:               NewAccountLimit("/account/interest-limits", account, "2s", 5),
		ACCOUNT_SIMULATED_MARGIN:              NewAccountLimit("/account/simulated_margin", account, "2s", 2),
		ACCOUNT_GREEKS:                        NewAccountLimit("/account/greeks", account, "2s", 10),
		ACCOUNT_POSITION_TIERS:                NewAccountLimit("/account/position-tiers", account, "2s", 10),
		TRAD_ORDER:                            NewAccountLimit("/trade/order", account, "2s", 60, "instId"),
		TRAD_BATCH_ORDERS:                     NewAccountLimit("/trade/batch-orders", account, "2s", 300, "instId"),
		TRAD_CANCEL_ORDER:                     NewAccountLimit("/trade/cancel-order", account, "2s", 60, "instId"),
		TRAD_CANCEL_BATCH_ORDERS:              NewAccountLimit("/trade/cancel-batch-orders", account, "2s", 300, "instId"),
		TRAD_CLOSE_POSITION:                   NewAccountLimit("/trade/close-position", account, "2s", 20, "instId"),
		TRAD_ORDERS_PENDING:                   NewAccountLimit("/trade/orders-pending", account, "2s", 60),
		TRAD_ORDERS_HISTORY_WEEK:              NewAccountLimit("/trade/orders-history", account, "2s", 40),
		TRAD_ORDERS_HISTORY_ARCHIVE:           NewAccountLimit("/trade/orders-history-archive", account, "2s", 20),
		TRAD_FILLS_HISTORY_THREE_DAYS:         NewAccountLimit("/trade/fills", account, "2s", 60),
		TRAD_FILLS_HISTORY_ARCHIVE:            NewAccountLimit("/trade/fills-history", account, "2s", 10),
		TRAD_ORDER_ALGO:                       NewAccountLimit("/trade/order-algo", account, "2s", 20, "instId"),
		TRAD_CANCEL_ALGO:                      NewAccountLimit("/trade/cancel-algos", account, "2s", 20, "instId"),
		TRAD_CANCLE_ADVANCE_ALGOS:             NewAccountLimit("/trade/cancel-advance-algos", account, "2s", 20, "instId"),
		TRAD_ORDERS_ALGO_PENDING:              NewAccountLimit("/trade/orders-algo-pending", account, "2s", 20),
		TRAD_ORDERS_ALGO_HISTORY:              NewAccountLimit("/trade/orders-algo-history", account, "2s", 20),
		RFQ_COUNTERPARTIES:                    NewAccountLimit("/rfq/counterparties", account, "2s", 5),
		ASSET_CURRENCIES:                      NewAccountLimit("/asset/currencies", account, "1s", 6),
		ASSET_TRANSFER:                        NewAccountLimit("/asset/transfer", account, "1s", 1, "ccy"),
		ASSET_TRANSFER_STATE:                  NewAccountLimit("/asset/transfer-state", account, "1s", 10),
		ASSET_BILLS:                           NewAccountLimit("/asset/bills", account, "1s", 6),
		ASSET_WITHDRAWAL:                      NewAccountLimit("/asset/withdrawal", account, "1s", 6),
		ASSET_WITHDRAWAL_LIGHTNING:            NewAccountLimit("/asset/withdrawal-lightning", account, "1s", 2),
		ASSET_CANCEL_WITHDRAWAL:               NewAccountLimit("/asset/cancel-withdrawal", account, "1s", 6),
		ASSET_WITHDRAWAL_HISTORY:              NewAccountLimit("/asset/withdrawal-history", account, "1s", 6),
		ASSET_DEPOSIT_HISTORY:                 NewAccountLimit("/asset/deposit-history", account, "1s", 6),
		ASSET_SUBACCOUNT_BILLS:                NewAccountLimit("/asset/subaccount/bills", account, "1s", 6),
		USER_SET_TRANSGER_OUT:                 NewAccountLimit("/users/subaccount/set-transfer-out", account, "1s", 1),
	}
}
