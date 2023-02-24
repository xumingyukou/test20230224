package spot_ws

type CommandMessage struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

type TradeMsg struct {
	Code string `json:"code"`
	Data struct {
		P      string `json:"p"`
		S      string `json:"s"`
		V      string `json:"v"`
		T      string `json:"t"`
		Symbol string `json:"symbol"`
		Ver    string `json:"ver"`
	} `json:"data"`
	Timestamp int64  `json:"timestamp"`
	Topic     string `json:"topic"`
}

type TickerMsg struct {
	Code string `json:"code"`
	Data struct {
		C      string `json:"c"`
		H      string `json:"h"`
		L      string `json:"l"`
		P      string `json:"p"`
		Symbol string `json:"symbol"`
		V      string `json:"v"`
		Ver    string `json:"ver"`
	} `json:"data"`
	Timestamp int64  `json:"timestamp"`
	Topic     string `json:"topic"`
}

type BaseMsg struct {
	Code string
}

type OrderBookMsg struct {
	Code string `json:"code"`
	Data struct {
		B      [][]string `json:"b"`
		S      [][]string `json:"s"`
		Symbol string     `json:"symbol"`
		Ver    string     `json:"ver"`
	} `json:"data"`
	Timestamp int64  `json:"timestamp"`
	Topic     string `json:"topic"`
}
