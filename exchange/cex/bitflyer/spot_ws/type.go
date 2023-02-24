package spot_ws

import "time"

type SIDE int8

const (
	Err SIDE = iota
	Ask
	Bid
)

type req struct {
	Jsonrpc string            `json:"jsonrpc,omitempty"`
	Method  string            `json:"method"`
	Params  map[string]string `json:"params"`
	ID      int               `json:"id,omitempty"`
}

type Resp_Ticker struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Channel string `json:"channel"`
		Message struct {
			ProductCode     string    `json:"product_code"`
			State           string    `json:"state"`
			Timestamp       time.Time `json:"timestamp"`
			TickId          int       `json:"tick_id"`
			BestBid         float64   `json:"best_bid"`
			BestAsk         float64   `json:"best_ask"`
			BestBidSize     float64   `json:"best_bid_size"`
			BestAskSize     float64   `json:"best_ask_size"`
			TotalBidDepth   float64   `json:"total_bid_depth"`
			TotalAskDepth   float64   `json:"total_ask_depth"`
			MarketBidSize   float64   `json:"market_bid_size"`
			MarketAskSize   float64   `json:"market_ask_size"`
			Ltp             float64   `json:"ltp"`
			Volume          float64   `json:"volume"`
			VolumeByProduct float64   `json:"volume_by_product"`
		} `json:"message"`
	} `json:"params"`
}

type Resp_Trade struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Channel string `json:"channel"`
		Message []struct {
			Id                         int64     `json:"id"`
			Side                       string    `json:"side"`
			Price                      float64   `json:"price"`
			Size                       float64   `json:"size"`
			ExecDate                   time.Time `json:"exec_date"`
			BuyChildOrderAcceptanceId  string    `json:"buy_child_order_acceptance_id"`
			SellChildOrderAcceptanceId string    `json:"sell_child_order_acceptance_id"`
		} `json:"message"`
	} `json:"params"`
}

type Resp_Depth struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Channel string `json:"channel"`
		Message struct {
			MidPrice float64     `json:"mid_price"`
			Bids     []levelList `json:"bids"`
			Asks     []levelList `json:"asks"`
		} `json:"message"`
	} `json:"params"`
}

type levelList struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}
