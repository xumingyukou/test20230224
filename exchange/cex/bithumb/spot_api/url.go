package spot_api

const (
	API_BASE_URL = "https://global-openapi.bithumb.pro/openapi/v1"
	WS_BASE_URL  = "wss://global-api.bithumb.pro/message/realtime"
)

type ReqUrl struct {
	SPOT_CONFIG    string
	SPOT_ORDERBOOK string
}

//Check Order history
func NewSpotReqUrl() *ReqUrl {
	return &ReqUrl{
		SPOT_CONFIG:    "/spot/config",
		SPOT_ORDERBOOK: "/spot/orderBook",
	}
}

type WsReqUrl struct {
	WS_BASE_URL string
}

func NewWsReqUrl() *WsReqUrl {
	return &WsReqUrl{
		WS_BASE_URL: WS_BASE_URL,
	}
}
