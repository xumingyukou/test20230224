package config

type SymbolPrecision struct {
	Amount    int64   `json:"amount"`
	Price     int64   `json:"price"`
	AmountMin float64 `json:"amount_min"`
}

type TradeFeeItem struct {
	Maker float64 `json:"maker"`
	Taker float64 `json:"taker"`
}

type TransferFeeItem struct {
	NetWork string  `json:"network"`
	Fee     float64 `json:"fee"`
}

type SymbolInfo struct {
	Symbol    string          `json:"symbol"`
	Precision SymbolPrecision `json:"precision"`
	TradeFee  TradeFeeItem    `json:"trade_fee"`
}

type TokenInfo struct {
	Token       string             `json:"token"`
	WithdrawFee []*TransferFeeItem `json:"withdraw_fee"` //network:fee
}

type ExchangeSpotConfig struct {
	Symbols []*SymbolInfo `json:"symbols"`
	Tokens  []*TokenInfo  `json:"tokens"` //coin:
}

type ApiConfig struct {
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Passphrase string `json:"passphrase"`
}
