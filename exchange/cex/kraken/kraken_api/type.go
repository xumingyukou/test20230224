package kraken_api

type Resp_Market_Books struct {
	Error  []interface{} `json:"error"`
	Result map[string]struct {
		Altname           string      `json:"altname"`
		Wsname            string      `json:"wsname"`
		AclassBase        string      `json:"aclass_base"`
		Base              string      `json:"base"`
		AclassQuote       string      `json:"aclass_quote"`
		Quote             string      `json:"quote"`
		Lot               string      `json:"lot"`
		PairDecimals      int         `json:"pair_decimals"`
		LotDecimals       int         `json:"lot_decimals"`
		LotMultiplier     int         `json:"lot_multiplier"`
		LeverageBuy       []int       `json:"leverage_buy"`
		LeverageSell      []int       `json:"leverage_sell"`
		Fees              [][]float64 `json:"fees"`
		FeesMaker         [][]float64 `json:"fees_maker"`
		FeeVolumeCurrency string      `json:"fee_volume_currency"`
		MarginCall        int         `json:"margin_call"`
		MarginStop        int         `json:"margin_stop"`
		Ordermin          string      `json:"ordermin"`
	} `json:"result"`
}
