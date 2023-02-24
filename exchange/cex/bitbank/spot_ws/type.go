package spot_ws

type TickerMsg struct {
	RoomName string `json:"room_name"`
	Message  struct {
		Pid  int `json:"pid"`
		Data struct {
			Last      string `json:"last"`
			Open      string `json:"open"`
			Timestamp int64  `json:"timestamp"`
			Sell      string `json:"sell"`
			Vol       string `json:"vol"`
			Buy       string `json:"buy"`
			High      string `json:"high"`
			Low       string `json:"low"`
		} `json:"data"`
	} `json:"message"`
}

type TradeMsg struct {
	RoomName string `json:"room_name"`
	Message  struct {
		Pid  int `json:"pid"`
		Data struct {
			Transactions []*Transaction `json:"transactions"`
		} `json:"data"`
	} `json:"message"`
}

type Transaction struct {
	Side          string `json:"side"`
	ExecutedAt    int64  `json:"executed_at"`
	Amount        string `json:"amount"`
	Price         string `json:"price"`
	TransactionID int64  `json:"transaction_id"`
}

type DepthDiffMsg struct {
	RoomName string `json:"room_name"`
	Message  struct {
		Data struct {
			T int64      `json:"t"`
			B [][]string `json:"b"`
			A [][]string `json:"a"`
			S string     `json:"s"`
		}
	} `json:"message"`
}

type DepthWholeMsg struct {
	RoomName string `json:"room_name"`
	Message  struct {
		Data struct {
			Timestamp  int64      `json:"timestamp"`
			Bids       [][]string `json:"bids"`
			Asks       [][]string `json:"asks"`
			SequenceId string     `json:"sequenceId"`
		}
	} `json:"message"`
}

type BaseMsg struct {
	RoomName string `json:"room_name"`
}
