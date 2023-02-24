package spot_api

const (
	SPOT_API_BASE_URL   = "https://api.gemini.com/"
	WS_API_Public_Topic = "wss://api.gemini.com/v2/marketdata"
)

type ReqUrl struct {
	SYMBOLS_URL string
}

func NewSpotReqUrl() *ReqUrl {
	return &ReqUrl{
		SYMBOLS_URL: "v1/symbols",
	}
}

type WsReqUrl struct {
}

func NewSpotWsUrl() *WsReqUrl {
	return &WsReqUrl{}
}
