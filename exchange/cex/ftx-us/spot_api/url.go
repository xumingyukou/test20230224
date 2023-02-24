package spot_api

const (
	SPOT_API_BASE_URL   = "https://ftx.us/api/"
	FUTURE_API_BASE_URL = "https://ftx.com/api/"
	WS_API_BASE_URL     = "wss://ftx.us/ws/"
)

type ReqUrl struct {
	MARKETS_URL                string
	ACCOUNT_URL                string
	ORDERS_URL                 string
	ORDERS_BY_CLIENT_ID_URL    string
	ORDERS_HISTORY_URL         string
	FILLS_URL                  string
	FUTURES_URL                string
	SPOT_MARGIN_HISTORY        string
	SPOT_MARGIN_BORROW_RATES   string
	SPOT_MARGIN_LENDING_RATES  string
	SPOT_MARGIN_MARKET_INFO    string
	SPOT_MARGIN_OFFERS         string
	WALLET_BALANCES_URL        string
	WALLET_DEPOSIT_ADDRESS_URL string
	WALLET_WITHDRAWALS_URL     string
	WALLET_DEPOSITS_URL        string
	WALLET_WITHDRAWAL_FEE_URL  string
}

//Check Order history
func NewSpotReqUrl() *ReqUrl {
	return &ReqUrl{
		MARKETS_URL:                "markets",
		ACCOUNT_URL:                "account",
		ORDERS_URL:                 "orders",
		ORDERS_BY_CLIENT_ID_URL:    "orders/by_client_id",
		ORDERS_HISTORY_URL:         "orders/history",
		FILLS_URL:                  "fills",
		SPOT_MARGIN_HISTORY:        "spot_margin/history",
		SPOT_MARGIN_BORROW_RATES:   "spot_margin/borrow_rates",
		SPOT_MARGIN_LENDING_RATES:  "spot_margin/lending_rates",
		SPOT_MARGIN_MARKET_INFO:    "spot_margin/market_info",
		SPOT_MARGIN_OFFERS:         "spot_margin/offers",
		WALLET_BALANCES_URL:        "wallet/balances",
		WALLET_DEPOSITS_URL:        "wallet/deposits",
		WALLET_DEPOSIT_ADDRESS_URL: "wallet/deposit_address",
		WALLET_WITHDRAWALS_URL:     "wallet/withdrawals",
		WALLET_WITHDRAWAL_FEE_URL:  "wallet/withdrawal_fee",
	}
}

func NewUBaseReqUrl() *ReqUrl {
	return &ReqUrl{
		FUTURES_URL: "futures",
		MARKETS_URL: "markets",
		ACCOUNT_URL: "account",
	}
}
