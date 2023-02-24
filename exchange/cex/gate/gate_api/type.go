package gate_api

type Resp_Market_Books struct {
	Id      int64      `json:"id"`
	Current int64      `json:"current"`
	Update  int64      `json:"update"`
	Asks    [][]string `json:"asks"`
	Bids    [][]string `json:"bids"`
}

type RespSymbols []date

type date struct {
	Id              string `json:"id"`
	Base            string `json:"base"`
	Quote           string `json:"quote"`
	Fee             string `json:"fee"`
	MinBaseAmount   string `json:"min_base_amount"`
	MinQuoteAmount  string `json:"min_quote_amount"`
	AmountPrecision int    `json:"amount_precision"`
	Precision       int    `json:"precision"`
	TradeStatus     string `json:"trade_status"`
	SellStart       int    `json:"sell_start"`
	BuyStart        int    `json:"buy_start"`
}
