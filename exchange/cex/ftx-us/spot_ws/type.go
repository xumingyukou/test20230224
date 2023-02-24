package spot_ws

import (
	"clients/exchange/cex/ftx-us/spot_api"
	"time"
)

type WSRequest struct {
	ChannelType ChannelType `json:"channel"`
	Market      string      `json:"market"`
	Op          Operation   `json:"op"`
}

type ChannelType string

const (
	OrderbookChannel = ChannelType("orderbook")
	TradesChannel    = ChannelType("trades")
	TickerChannel    = ChannelType("ticker")
	MarketsChannel   = ChannelType("markets")
	FillsChannel     = ChannelType("fills")
	OrdersChannel    = ChannelType("orders")
)

type Operation string

const (
	Subscribe   = Operation("subscribe")
	UnSubscribe = Operation("unsubscribe")
)

type ResponseType string

const (
	Error        = ResponseType("error")
	Subscribed   = ResponseType("subscribed")
	UnSubscribed = ResponseType("unsubscribed")
	Info         = ResponseType("info")
	Partial      = ResponseType("partial")
	Update       = ResponseType("update")
)

type WSAuthorizationRequest struct {
	Args WSAuthorizationArgs `json:"args"`
	Op   string              `json:"op"`
}

type WSAuthorizationArgs struct {
	Key  string `json:"key"`
	Sign string `json:"sign"`
	Time int64  `json:"time"`
}

type TickerResponse struct {
	Channel ChannelType  `json:"channel"`
	Market  string       `json:"market"`
	Type    ResponseType `json:"type"`
	Data    TickerItem   `json:"data"`
	ErrMsg
}

type TickerItem struct {
	Bid     float64 `json:"bid"`
	Ask     float64 `json:"ask"`
	BidSize float64 `json:"bidSize"`
	AskSize float64 `json:"askSize"`
	Last    float64 `json:"last"`
	Time    float64 `json:"time"`
}

type TradeResponse struct {
	Channel ChannelType  `json:"channel"`
	Market  string       `json:"market"`
	Type    ResponseType `json:"type"`
	Data    []TradeItem  `json:"data"`
	ErrMsg
}

type ErrMsg struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type TradeItem struct {
	Id          int64     `json:"id"`
	Price       float64   `json:"price"`
	Size        float64   `json:"size"`
	Side        string    `json:"side"`
	Liquidation bool      `json:"liquidation"`
	Time        time.Time `json:"time"`
}

type OrderbooksResponse struct {
	Channel ChannelType   `json:"channel"`
	Market  string        `json:"market"`
	Type    ResponseType  `json:"type"`
	Data    OrderbookItem `json:"data"`
	ErrMsg
}

type OrderbookItem struct {
	Time     float64     `json:"time"`
	Checksum uint32      `json:"checksum"`
	Bids     [][]float64 `json:"bids"`
	Asks     [][]float64 `json:"asks"`
	Action   string      `json:"action"`
}

type OrdersResponse struct {
	Channel ChannelType  `json:"channel"`
	Data    OrderItem    `json:"data"`
	Type    ResponseType `json:"type"`
	ErrMsg
}

type OrderItem struct {
	Id            int64                `json:"id"`
	ClientId      string               `json:"clientId"`
	Market        string               `json:"market"`
	Type          string               `json:"type"`
	Side          string               `json:"side"`
	Price         float64              `json:"price"`
	Size          float64              `json:"size"`
	Status        spot_api.OrderStatus `json:"status"`
	FilledSize    float64              `json:"filledSize"`
	RemainingSize float64              `json:"remainingSize"`
	ReduceOnly    bool                 `json:"reduceOnly"`
	Liquidation   bool                 `json:"liquidation"`
	AvgFillPrice  float64              `json:"avgFillPrice"`
	PostOnly      bool                 `json:"postOnly"`
	Ioc           bool                 `json:"ioc"`
	CreatedAt     time.Time            `json:"createdAt"`
}

type PingMessage struct {
	Op string `json:"op"`
}
