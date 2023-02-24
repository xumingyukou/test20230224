package bitMex_api

const (
	GLOBAL_API_BASE_URL = "https://www.bitmex.com/api"
	ProxyUrl            = "http://127.0.0.1:7890"
)

type ReqUrl struct {
	SYMBOLS string
	DEPTH   string
}

func NewSpotReqUrl() *ReqUrl {
	return &ReqUrl{
		SYMBOLS: "/v1/instrument/active",
		DEPTH:   "/v1/orderBook/L2",
	}
}
