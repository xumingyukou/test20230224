package spot_api

const (
	SPOT_API_Public_URL = "https://api-pub.bitfinex.com/"
	WS_API_Public_Topic = "wss://api-pub.bitfinex.com/ws/2"
)

type ReqUrl struct {
	EXCHANGE_URL string
	CURRENCY_URL string
}

func NewSpotReqUrl() *ReqUrl {
	return &ReqUrl{
		EXCHANGE_URL: "v2/conf/pub:list:pair:exchange",
		CURRENCY_URL: "v2/conf/pub:list:currency",
	}
}

type WsReqUrl struct {
}

func NewSpotWsUrl() *WsReqUrl {
	return &WsReqUrl{}
}
