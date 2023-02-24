package bitFlyer_api

const (
	GLOBAL_API_BASE_URL = "https://api.bitflyer.com"
	ProxyUrl            = "http://127.0.0.1:7890"
)

type ReqUrl struct {
	SYMBOLS string
}

func NewSpotReqUrl() *ReqUrl {
	return &ReqUrl{
		SYMBOLS: "/v1/markets",
	}
}
