package spot_api

const (
	CONTANTTYPE = "application/json"
)

type RespError struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

type AllSymbols []string

type SymbolDetails struct {
	Symbol         string  `json:"symbol"`
	BaseCurrency   string  `json:"base_currency"`
	QuoteCurrency  string  `json:"quote_currency"`
	TickSize       float64 `json:"tick_size"`
	QuoteIncrement float64 `json:"quote_increment"`
	MinOrderSize   string  `json:"min_order_size"`
	Status         string  `json:"status"`
	WrapEnabled    bool    `json:"wrap_enabled"`
}
