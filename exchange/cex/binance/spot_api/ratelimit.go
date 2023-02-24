package spot_api

const (
	R_IP_1M       = "ip_1m"
	R_ACCOUNT_10S = "account_10s"
	R_ACCOUNT_1M  = "account_1m"
	R_ACCOUNT_1D  = "account_24h"
)

var (
	BinanceWeightHeaderMap = map[string]string{
		"X-Mbx-Used-Weight-1m":  R_IP_1M,
		"X-Mbx-Order-Count-10s": R_ACCOUNT_10S,
		"X-Mbx-Order-Count-1m":  R_ACCOUNT_1M,
		"X-Mbx-Order-Count-1d":  R_ACCOUNT_1D,
	}
)
