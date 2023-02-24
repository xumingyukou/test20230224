package spot_api

const (
	SPOT_API_BASE_URL = "https://api.exchange.coinbase.com"
	WS_API_BASE_URL   = "wss://ws-feed.exchange.coinbase.com"
)

type ReqUrl struct {
	DEPTH_URL        string
	TICKER_PRICE_URL string
	TRADE_URL        string
	FILLS_URL        string
	ORDERS_URL       string
}

type WsReqUrl struct {

	//BOOK_TICKER_URL      string
	//DEPTH_LIMIT_FULL_URL string
	//DEPTH_INCRE_URL      string
	//USER_DATA_URL        string
}

func NewSpotWsUrl() *WsReqUrl {
	return &WsReqUrl{
		//MARKPRICE_URL:        "<symbol>@markPrice@1s",
		//AGGTRADE_URL:         "<symbol>@aggTrade",
		//TRADE_URL:            "<symbol>@trade",
		//BOOK_TICKER_URL:      "<symbol>@bookTicker",
		//DEPTH_LIMIT_FULL_URL: "<symbol>@depth",
		//DEPTH_INCRE_URL:      "<symbol>@depth",
		//USER_DATA_URL:        "<listenKey>",
	}
}

func NewSpotReqUrl() *ReqUrl {
	return &ReqUrl{
		DEPTH_URL:        "depth",
		TICKER_PRICE_URL: "ticker/price",
		ORDERS_URL:       "/orders",
		TRADE_URL:        "/trades",
		FILLS_URL:        "/fills",
	}
}
