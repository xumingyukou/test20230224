package spot_ws

const (
	WS_API_BASE_URL = "wss://stream.binance.com:9443"

	SINGLE_API_URL = "/ws/"
	STREAM_API_URL = "/stream?streams="
)

const (
	TRADE_URL            = "<symbol>@trade"
	BOOK_TICKER_URL      = "<symbol>@bookTicker"
	DEPTH_LIMIT_FULL_URL = "<symbol>@depth20@100ms"
	DEPTH_INCRE_URL      = "<symbol>@depth@100ms"
	USER_DATA_URL        = "<listenKey>"
)
