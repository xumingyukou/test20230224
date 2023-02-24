package spot_ws

type SendMsg1 struct {
	Symbol string `json:"symbol"`
	Topic  string `json:"topic"`
	Event  string `json:"event"`
	Params Param1 `json:"params"`
}

type SendMsg2 struct {
	Topic  string `json:"topic"`
	Event  string `json:"event"`
	Params Param2 `json:"params"`
}

type Param1 struct {
	Binary bool `json:"binary"`
}

type Param2 struct {
	Symbol string `json:"symbol"`
	Binary bool   `json:"binary"`
}

type WSRequest struct {
	Topic  string `json:"topic"`
	Event  string `json:"event"`
	Binary bool   `json:"binary"`
}

type PingMsg struct {
	Ping int64 `json:"ping"`
}

type FirstResp struct {
	Msg  string `json:"msg"`
	Pong int64  `json:"pong"`
}

type TradeResponse struct {
	Topic  string `json:"topic"`
	Params struct {
		Symbol     string `json:"symbol"`
		Binary     string `json:"binary"`
		SymbolName string `json:"symbolName"`
	} `json:"params"`
	Data struct {
		V string `json:"v"`
		T int64  `json:"t"`
		P string `json:"p"`
		Q string `json:"q"`
		M bool   `json:"m"`
	} `json:"data"`
}

type TickerResponse struct {
	Pong   int64  `json:"pong"`
	Topic  string `json:"topic"`
	Params struct {
		Symbol     string `json:"symbol"`
		Binary     string `json:"binary"`
		SymbolName string `json:"symbolName"`
	} `json:"params"`
	Data struct {
		Symbol   string `json:"symbol"`
		BidPrice string `json:"bidPrice"`
		BidQty   string `json:"bidQty"`
		AskPrice string `json:"askPrice"`
		AskQty   string `json:"askQty"`
		Time     int64  `json:"time"`
	} `json:"data"`
}

type RespIncrementDepthStream struct {
	Pong       int64  `json:"pong"`
	Symbol     string `json:"symbol"`
	SymbolName string `json:"symbolName"`
	Topic      string `json:"topic"`
	Params     struct {
		RealtimeInterval string `json:"realtimeInterval"`
		Binary           string `json:"binary"`
	} `json:"params"`
	Data     []RespIncrementDepth `json:"data"`
	F        bool                 `json:"f"`
	SendTime int64                `json:"sendTime"`
}

type RespIncrementDepth struct {
	E int        `json:"e"`
	T int64      `json:"t"`
	V string     `json:"v"`
	B [][]string `json:"b"`
	A [][]string `json:"a"`
	O int        `json:"o"`
}

type RespLimitDepthStream struct {
	Pong   int64  `json:"pong"`
	Topic  string `json:"topic"`
	Params struct {
		Symbol     string `json:"symbol"`
		Binary     string `json:"binary"`
		SymbolName string `json:"symbolName"`
	} `json:"params"`
	Data RespLimitDepth `json:"data"`
}

type RespLimitDepth struct {
	S string     `json:"s"`
	T int64      `json:"t"`
	V string     `json:"v"`
	B [][]string `json:"b"`
	A [][]string `json:"a"`
}
