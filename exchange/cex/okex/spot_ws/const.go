package spot_ws

type DepositStatus int

const (
	WS_API_BASE_URL      = "wss://wsaws.okx.com:8443/ws/v5"
	TEST_WS_API_BASE_URL = "wss://wspap.okx.com:8443/ws/v5"
	WS_API_PUBLIC        = "/public"
	WS_API_PRIVATE       = "/private"
)

const (
	//0(0:pending,6: credited but cannot withdraw, 1:success)

	DOPSIT_TYPE_PENDING   DepositStatus = 0
	DOPSIT_TYPE_SUCCESS   DepositStatus = 1
	DOPSIT_TYPE_CONFIRMED DepositStatus = 2
	DOPSIT_TYPE_FAIL      DepositStatus = 3
)

const (
	OP_ORDER         = "order"
	OP_ORDERS        = "batch-orders"
	OP_CANCEL_ORDER  = "cancel-order"
	OP_CANCEL_ORDERS = "batch-cancel-orders"
)
