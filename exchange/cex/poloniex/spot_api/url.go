package spot_api

const (
	REST_PUBLIC_BASE_URL = "https://api.poloniex.com"
	WS_PUBLIC_BASE_URL   = "wss://ws.poloniex.com/ws/public"
)

type ReqUrl struct {
	MARKETS string
}

//Check Order history
func NewSpotReqUrl() *ReqUrl {
	return &ReqUrl{
		MARKETS: "/markets",
	}
}

type WsReqUrl struct {
	WS_PUBLIC_BASE_URL string
}

func NewWsReqUrl() *WsReqUrl {
	return &WsReqUrl{
		WS_PUBLIC_BASE_URL: WS_PUBLIC_BASE_URL,
	}
}
