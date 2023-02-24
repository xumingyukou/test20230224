package spot_ws

type BaseSubscribe interface {
	Unsub()
	Sub()
	WithSymbols(symbols []string) BaseSubscribe
}

type SubscribeMessage struct {
	Event   string   `json:"event"`
	Channel []string `json:"channel"`
	Symbols []string `json:"symbols"`
}

func (m *SubscribeMessage) Unsub() {
	m.Event = "unsubscribe"
}

func (m *SubscribeMessage) Sub() {
	m.Event = "subscribe"
}

func (m *SubscribeMessage) WithSymbols(symbols []string) BaseSubscribe {
	return &SubscribeMessage{
		Event:   m.Event,
		Channel: m.Channel,
		Symbols: symbols,
	}
}

type BookSubscribeMessage struct {
	Event   string   `json:"event"`
	Channel []string `json:"channel"`
	Symbols []string `json:"symbols"`
	Depth   int      `json:"depth"`
}

func (m *BookSubscribeMessage) Unsub() {
	m.Event = "unsubscribe"
}

func (m *BookSubscribeMessage) Sub() {
	m.Event = "subscribe"
}

func (m BookSubscribeMessage) WithSymbols(symbols []string) BaseSubscribe {
	return &BookSubscribeMessage{
		Event:   m.Event,
		Channel: m.Channel,
		Symbols: symbols,
		Depth:   m.Depth,
	}
}

type EventMessage struct {
	Event   string `json:"event"`
	Message string `json:"message"`
}

type TradeMsg struct {
	Channel string `json:"channel"`
	Data    []struct {
		Symbol     string `json:"symbol"`
		Amount     string `json:"amount"`
		TakerSide  string `json:"takerSide"`
		Quantity   string `json:"quantity"`
		CreateTime int64  `json:"createTime"`
		Price      string `json:"price"`
		ID         string `json:"id"`
		Ts         int64  `json:"ts"`
	} `json:"data"`
}

type BookMsg struct {
	Channel string `json:"channel"`
	Data    []struct {
		Symbol     string     `json:"symbol"`
		CreateTime int64      `json:"createTime"`
		Asks       [][]string `json:"asks"`
		Bids       [][]string `json:"bids"`
		ID         int        `json:"id"`
		Ts         int64      `json:"ts"`
	} `json:"data"`
}

type BookLv2Msg struct {
	Channel string     `json:"channel"`
	Action  string     `json:"action"`
	Data    []DepthMsg `json:"data"`
}

type DepthMsg struct {
	Symbol     string     `json:"symbol"`
	Asks       [][]string `json:"asks"`
	Bids       [][]string `json:"bids"`
	LastID     int        `json:"lastId"`
	ID         int        `json:"id"`
	Ts         int64      `json:"ts"`
	CreateTime int64      `json:"createTime"`
}
