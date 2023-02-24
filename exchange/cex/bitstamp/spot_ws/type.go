package spot_ws

import "clients/exchange/cex/bitstamp/spot_api"

type Event string

const (
	SUBSCRIBE   = "bts:subscribe"
	UNSUBSCRIBE = "bts:unsubscribe"
)

type Data struct {
	Channel string `json:"channel"`
	Auth    string `json:"auth"`
}

type Request struct {
	Event Event `json:"event"`
	Data  Data  `json:"data"`
}

const (
	LiveTrades          = "live_trades_"
	LiveOrders          = "live_orders_"
	LiveOrderBook       = "order_book_"
	LiveDetailOrderBook = "detail_order_book_"
	LiveFullOrderBook   = "diff_order_book_"
)

type RespTradeStream struct {
	Data struct {
		Id             int       `json:"id"`
		Timestamp      string    `json:"timestamp"`
		Amount         float64   `json:"amount"`
		AmountStr      string    `json:"amount_str"`
		Price          float64   `json:"price"`
		PriceStr       string    `json:"price_str"`
		Type           TradeType `json:"type"`
		Microtimestamp string    `json:"microtimestamp"`
		BuyOrderId     int64     `json:"buy_order_id"`
		SellOrderId    int64     `json:"sell_order_id"`
	} `json:"data"`
	Channel string `json:"channel"`
	Event   string `json:"event"`
}

type TradeType int

const (
	BUY  = TradeType(0)
	SELL = TradeType(1)
)

type RespDepthStream struct {
	Data struct {
		Timestamp      string                `json:"timestamp"`
		MicroTimestamp string                `json:"microtimestamp"`
		Bids           []*spot_api.DepthItem `json:"bids"`
		Asks           []*spot_api.DepthItem `json:"asks"`
	} `json:"data"`
	Channel string `json:"channel"`
	Event   string `json:"event"`
}

type PingMessage struct {
	Event string `json:"event"`
}
