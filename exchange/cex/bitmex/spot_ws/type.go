package spot_ws

import (
	"github.com/warmplanet/proto/go/depth"
	"time"
)

type SIDE int8

const (
	Err SIDE = iota
	Ask
	Bid
)

type req struct {
	Op   string   `json:"op"`
	Args []string `json:"args"`
}

type Resp_Ticker struct {
	Table  string `json:"table"`
	Action string `json:"action"`
	Data   []struct {
		Symbol             string    `json:"symbol"`
		LastPriceProtected float64   `json:"lastPriceProtected"`
		BidPrice           float64   `json:"bidPrice"`
		MidPrice           float64   `json:"midPrice"`
		AskPrice           float64   `json:"askPrice"`
		ImpactBidPrice     float64   `json:"impactBidPrice"`
		ImpactMidPrice     float64   `json:"impactMidPrice"`
		ImpactAskPrice     float64   `json:"impactAskPrice"`
		Timestamp          time.Time `json:"timestamp"`
	} `json:"data"`
}

type Resp_Trade struct {
	Table  string `json:"table"`
	Action string `json:"action"`
	Data   []struct {
		Timestamp       time.Time `json:"timestamp"`
		Symbol          string    `json:"symbol"`
		Side            string    `json:"side"`
		Size            float64   `json:"size"`
		Price           float64   `json:"price"`
		TickDirection   string    `json:"tickDirection"`
		TrdMatchID      string    `json:"trdMatchID"`
		GrossValue      float64   `json:"grossValue"`
		HomeNotional    float64   `json:"homeNotional"`
		ForeignNotional float64   `json:"foreignNotional"`
	} `json:"data"`
}

type Resp_Depth struct {
	Table  string `json:"table"`
	Action string `json:"action"`
	Data   []struct {
		Symbol    string    `json:"symbol"`
		Id        int64     `json:"id"`
		Side      string    `json:"side"`
		Size      float64   `json:"size"`
		Price     float64   `json:"price"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"data"`
}

type levelList struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

type FullDepth struct {
	Asks map[int64]*depth.DepthLevel
	Bids map[int64]*depth.DepthLevel
}
