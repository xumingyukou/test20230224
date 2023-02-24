package spot_api

type RespError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type RespGetMarkets []struct {
	Symbol            string `json:"symbol"`
	BaseCurrencyName  string `json:"baseCurrencyName"`
	QuoteCurrencyName string `json:"quoteCurrencyName"`
	DisplayName       string `json:"displayName"`
	State             string `json:"state"`
	VisibleStartTime  int64  `json:"visibleStartTime"`
	TradableStartTime int64  `json:"tradableStartTime"`
	SymbolTradeLimit  struct {
		Symbol        string `json:"symbol"`
		PriceScale    int    `json:"priceScale"`
		QuantityScale int    `json:"quantityScale"`
		AmountScale   int    `json:"amountScale"`
		MinQuantity   string `json:"minQuantity"`
		MinAmount     string `json:"minAmount"`
		HighestBid    string `json:"highestBid"`
		LowestAsk     string `json:"lowestAsk"`
	} `json:"symbolTradeLimit"`
}
