package spot_api

const (
	SPOT_PUBLIC_BASE_URL  = "https://public.bitbank.cc"
	SPOT_PRIVATE_BASE_URL = "https://api.bitbank.cc/v1"
	// WS_BASE_URL           = "wss://stream.bitbank.cc/socket.io/?EIO=2&transport=websocket"
	WS_BASE_URL = "wss://stream.bitbank.cc/socket.io/?EIO=2&transport=websocket"
)

type ReqUrl struct {
}

//Check Order history
func NewSpotReqUrl() *ReqUrl {
	return &ReqUrl{}
}

type WsReqUrl struct {
	WS_BASE_URL string
	TICKER_URL  string
}

func NewWsReqUrl() *WsReqUrl {
	return &WsReqUrl{
		WS_BASE_URL: WS_BASE_URL,
	}
}
