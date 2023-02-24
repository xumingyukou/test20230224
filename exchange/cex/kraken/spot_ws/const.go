package spot_ws

type DepositStatus int

const (
	WS_PUBLIC_BASE_URL = "wss://ws.kraken.com"
	WS_AUTH_BASE_URL   = "ws-auth.kraken.com"
)

const (
	//0(0:pending,6: credited but cannot withdraw, 1:success)

	DOPSIT_TYPE_PENDING   DepositStatus = 0
	DOPSIT_TYPE_SUCCESS   DepositStatus = 1
	DOPSIT_TYPE_CONFIRMED DepositStatus = 2
	DOPSIT_TYPE_FAIL      DepositStatus = 3
)

const WsSubscrible = "subscribe"
