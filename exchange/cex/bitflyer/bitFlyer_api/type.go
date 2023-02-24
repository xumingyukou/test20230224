package bitFlyer_api

type RespSymbols []date

type date struct {
	ProductCode string `json:"product_code"`
	MarketType  string `json:"market_type"`
	Alias       string `json:"alias,omitempty"`
}
