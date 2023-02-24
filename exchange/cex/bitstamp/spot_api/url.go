package spot_api

const (
	SPOT_API_BASE_URL = "https://www.bitstamp.net/api/v2"

	WS_API_BASE_URL = "wss://ws.bitstamp.net"
)

type ReqUrl struct {
	DEPTH_URL              string
	WEBSOCKETS_TOKEN_URL   string
	TRADING_PAIRS_INFO_URL string
	ORDER_BOOK_URL         string
	//TICKER_PRICE_URL      string
}

type WsReqUrl struct {
	TRADE_URL string
	//WEBSOCKETS_TOKEN_URL string
}

func NewSpotReqUrl() *ReqUrl {
	return &ReqUrl{
		WEBSOCKETS_TOKEN_URL:   "/websockets_token/",
		TRADING_PAIRS_INFO_URL: "/trading-pairs-info/",
		ORDER_BOOK_URL:         "/order_book/",
	}
}

func NewSpotWsUrl() *WsReqUrl {
	return &WsReqUrl{}
}
