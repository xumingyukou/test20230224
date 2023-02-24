package spot_api

type RespError struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type RespSpotConfig struct {
	Data struct {
		CoinConfig     []*CoinConfig     `json:"coinConfig"`
		ContractConfig []*ContractConfig `json:"contractConfig"`
		SpotConfig     []*SpotConfig     `json:"spotConfig"`
	} `json:"data"`
	Code      string `json:"code"`
	Msg       string `json:"msg"`
	Timestamp int64  `json:"timestamp"`
}

type CoinConfig struct {
	MakerFeeRate   interface{} `json:"makerFeeRate"`
	TakerFeeRate   interface{} `json:"takerFeeRate"`
	MinWithdraw    string      `json:"minWithdraw"`
	WithdrawFee    string      `json:"withdrawFee"`
	Name           string      `json:"name"`
	DepositStatus  string      `json:"depositStatus"`
	FullName       string      `json:"fullName"`
	WithdrawStatus string      `json:"withdrawStatus"`
}
type ContractConfig struct {
	Symbol       string `json:"symbol"`
	MakerFeeRate string `json:"makerFeeRate"`
	TakerFeeRate string `json:"takerFeeRate"`
}
type PercentPrice struct {
	MultiplierDown string `json:"multiplierDown"`
	MultiplierUp   string `json:"multiplierUp"`
}
type SpotConfig struct {
	Symbol       string       `json:"symbol"`
	Accuracy     []string     `json:"accuracy"`
	OpenTime     int64        `json:"openTime"`
	PercentPrice PercentPrice `json:"percentPrice"`
	OpenPrice    string       `json:"openPrice"`
}

type RespOrderBook struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Data    struct {
		B   [][]string `json:"b"`
		S   [][]string `json:"s"`
		Ver string     `json:"ver"`
	} `json:"data"`
	Timestamp int64 `json:"timestamp"`
}
