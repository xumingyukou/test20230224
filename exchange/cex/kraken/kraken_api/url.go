package kraken_api

const (
	GLOBAL_API_BASE_URL = "https://api.kraken.com/0/public"
	ProxyUrl            = "http://127.0.0.1:7890"
)

type ReqUrl struct {
	Asset_Pairs string
}

func NewSpotReqUrl() *ReqUrl {
	return &ReqUrl{
		Asset_Pairs: "/AssetPairs",
	}
}
