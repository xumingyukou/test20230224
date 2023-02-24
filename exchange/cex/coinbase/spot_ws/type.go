package spot_ws

import (
	"encoding/json"
	"time"
)

//This should be standard format
type Request struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
	Options
}

type SignedRequest struct {
	Request
	Signature  string `json:"signature"`
	Key        string `json:"key"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
}

type Options struct {
	Extensions string `json:"Sec-Websocket-Extensions"`
}

type Error struct {
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

//This is more advanced format
type Request2 struct {
	Type     string           `json:"type"`
	Channels []ChannelRequest `json:"channels"`
}

type ChannelRequest struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
}

//This is untested
type Request3 struct {
	Type       string           `json:"type"`
	ProductIds []string         `json:"product_ids"`
	Channels   []ChannelRequest `json:"channels"`
}

type RespBookTickerStream struct {
	Type      string    `json:"type"`
	Sequence  int       `json:"sequence"`
	ProductId string    `json:"product_id"`
	Price     string    `json:"price"`
	Open24H   string    `json:"open_24h"`
	Volume24H string    `json:"volume_24h"`
	Low24H    string    `json:"low_24h"`
	High24H   string    `json:"high_24h"`
	Volume30D string    `json:"volume_30d"`
	BestBid   string    `json:"best_bid"`
	BestAsk   string    `json:"best_ask"`
	Side      string    `json:"side"`
	Time      time.Time `json:"time"`
	TradeId   int       `json:"trade_id"`
	LastSize  string    `json:"last_size"`
	Error
}

type RespTradeStream struct {
	Type         string    `json:"type"`
	TradeId      int       `json:"trade_id"`
	MakerOrderId string    `json:"maker_order_id"`
	TakerOrderId string    `json:"taker_order_id"`
	Side         Side      `json:"side"`
	Size         string    `json:"size"`
	Price        string    `json:"price"`
	ProductId    string    `json:"product_id"`
	Sequence     int64     `json:"sequence"`
	Time         time.Time `json:"time"`
	Error
}

type RespLimitDepthStream struct {
	Type      string     `json:"type"`
	ProductId string     `json:"product_id"`
	Bids      [][]string `json:"bids,omitempty"`
	Asks      [][]string `json:"asks,omitempty"`
	Time      time.Time  `json:"time"`
	Error
}

type RespIncrementDepthStream struct {
	Type      string        `json:"type"`
	ProductId string        `json:"product_id"`
	Changes   []ChangeTuple `json:"changes,omitempty"`
	Time      time.Time     `json:"time"`
	Error
}

type RespIncrementSnapshot struct {
	Type      string        `json:"type"`
	ProductId string        `json:"product_id"`
	Bids      [][]string    `json:"bids,omitempty"`
	Asks      [][]string    `json:"asks,omitempty"`
	Changes   []ChangeTuple `json:"changes,omitempty"`
	Time      time.Time     `json:"time"`
	Error
}

//0-Side, 1-Price, 2-Size [List of size 3]
type ChangeTuple struct {
	Side  Side
	Price string
	Size  string
}

func (c *ChangeTuple) UnmarshalJSON(data []byte) error {
	var tmp []json.RawMessage
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[0], &c.Side); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[1], &c.Price); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[2], &c.Size); err != nil {
		return err
	}
	return nil
}

type FullChannelResponse struct {
	OrderId   string    `json:"order_id"`
	OrderType string    `json:"order_type"`
	Size      string    `json:"size"`
	Price     string    `json:"price"`
	ClientOid string    `json:"client_oid"`
	Type      string    `json:"type"`
	Side      Side      `json:"side"`
	ProductId string    `json:"product_id"`
	Time      time.Time `json:"time"`
	Sequence  int64     `json:"sequence"`
}

type Side string

const (
	BUY  = Side("buy")
	SELL = Side("sell")
)

type Operation string

const (
	Subscribe   = "subscribe"
	Unsubscribe = "unsubscribe"
)
