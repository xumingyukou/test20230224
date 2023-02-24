package spot_api

const (
	CONTANTTYPE = "application/json"
)

type RespError struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}
type ExchangeSymbols [][]string
type CurrencySymbols [][]string
