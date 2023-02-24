package spot_ws

type Msg struct {
	Type          string         `json:"type"`
	Subscriptions []Subcriptions `json:"subscriptions"`
}

type Subcriptions struct {
	Name    string   `json:"name"`
	Symbols []string `json:"symbols"`
}

type Response struct {
	Type   string `json:"type"`
	Symbol string `json:"symbol"`
}

type L2Update struct {
	Type    string       `json:"type"`
	Symbol  string       `json:"symbol"`
	Changes [][]string   `json:"changes"`
	Trades  []TradesItem `json:"trades"`
}

type TradesItem struct {
	Type      string `json:"type"`
	Symbol    string `json:"symbol"`
	EventID   int64  `json:"event_id"`
	Timestamp int64  `json:"timestamp"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	Side      string `json:"side"`
}
type Trade struct {
}

//type Heartbeat struct {
//	Type      string `json:"type"`
//	Timestamp int64  `json:"timestamp"`
//}
