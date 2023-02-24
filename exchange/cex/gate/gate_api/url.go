package gate_api

const (
	GLOBAL_API_BASE_URL = "https://api.gateio.ws/api/v4"
)

type ReqUrl struct {
	MARKET_BOOKS string
	SYMBOLS      string
}

func NewSpotReqUrl() *ReqUrl {
	return &ReqUrl{
		MARKET_BOOKS: "/spot/order_book",
		SYMBOLS:      "/spot/currency_pairs",
	}
}
