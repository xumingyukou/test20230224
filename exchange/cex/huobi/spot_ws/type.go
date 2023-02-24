package spot_ws

import "math/big"

type Sub struct {
	Sub      string `json:"sub,omitempty"`
	Req      string `json:"req,omitempty"`
	ID       string `json:"id"`
	DataType string `json:"data_type,omitempty"`
}

type SubFund struct {
	Op    string `json:"op"`
	CID   string `json:"cid"`
	Topic string `json:"topic"`
}

type Req struct {
	Req string `json:"req"`
	ID  string `json:"id"`
}

type Unsub struct {
	Unsub string `json:"unsub"`
	ID    string `json:"id"`
}

type TradeResponse struct {
	Ping int64  `json:"ping"`
	Ch   string `json:"ch"`
	Ts   int64  `json:"ts"`
	Tick struct {
		ID   int64       `json:"id"`
		Ts   int64       `json:"ts"`
		Data []TradeInfo `json:"data"`
	} `json:"tick"`
}

type TradeInfo struct {
	ID        big.Int `json:"id"`
	Ts        int64   `json:"ts"`
	TradeID   int64   `json:"tradeId"`
	Amount    float64 `json:"amount"`
	Price     float64 `json:"price"`
	Direction string  `json:"direction"`
}

type TickerResponse struct {
	Ping int64  `json:"ping"`
	Ch   string `json:"ch"`
	Ts   int64  `json:"ts"`
	Tick struct {
		SeqId     int64   `json:"seqId"`
		Bid       float64 `json:"bid"`
		BidSize   float64 `json:"bidSize"`
		Ask       float64 `json:"ask"`
		AskSize   float64 `json:"askSize"`
		QuoteTime int64   `json:"quoteTime"`
		Symbol    string  `json:"symbol"`
	} `json:"tick"`
}

type FutureOrSwapTickerResp struct {
	Ping int64  `json:"ping"`
	Ch   string `json:"ch"`
	Ts   int64  `json:"ts"`
	Tick struct {
		Mrid    int64     `json:"mrid"`
		Id      int64     `json:"id"`
		Bid     []float64 `json:"bid"`
		Ask     []float64 `json:"ask"`
		Ts      int64     `json:"ts"`
		Version int64     `json:"version"`
		Ch      string    `json:"ch"`
	} `json:"tick"`
}

type RespLimitDepthStream struct {
	Ping int64          `json:"ping"`
	Ch   string         `json:"ch"`
	Ts   int64          `json:"ts"`
	Tick RespLimitDepth `json:"tick"`
}
type RespLimitDepth struct {
	SeqNum int64       `json:"seqNum"`
	Bids   [][]float64 `json:"bids"`
	Asks   [][]float64 `json:"asks"`
}

type RespIncrementDepthStream struct {
	Ping int64              `json:"ping"`
	Ch   string             `json:"ch"`
	Ts   int64              `json:"ts"`
	Tick RespIncrementDepth `json:"tick"`
}

type RespIncrementDepth struct {
	SeqNum     int64       `json:"seqNum"`
	PrevSeqNum int64       `json:"prevSeqNum"`
	Bids       [][]float64 `json:"bids"`
	Asks       [][]float64 `json:"asks"`
}

type RespSnapShotTemp struct {
	Rep    string `json:"rep"`
	Ping   int64  `json:"ping"`
	Ch     string `json:"ch"`
	Subbed string `json:"Subbed"`
}

type RespLimitSnapShotStream struct {
	ID     string         `json:"id"`
	Rep    string         `json:"rep"`
	Status string         `json:"status"`
	Ts     int64          `json:"ts"`
	Data   RespLimitDepth `json:"data"`
}

type RespIncrementSnapShotStream struct {
	Ch   string             `json:"ch"`
	Ts   int64              `json:"ts"`
	Tick RespIncrementDepth `json:"tick"`
}

type RespFutureOrSwapIncrementSnapShotStream struct {
	Ping   int64            `json:"ping"`
	Ch     string           `json:"ch"`
	Ts     int64            `json:"ts"`
	Subbed string           `json:"Subbed"`
	Tick   FutureOrSwapTick `json:"tick"`
}

type FutureOrSwapTick struct {
	Asks    [][]float64 `json:"asks"`
	Bids    [][]float64 `json:"bids"`
	Ch      string      `json:"ch"`
	Event   string      `json:"event"`
	Id      int64       `json:"id"`
	Mrid    int64       `json:"mrid"`
	Ts      int64       `json:"ts"`
	Version int64       `json:"version"`
}

type RespFundRate struct {
	Ping  int64       `json:"ping"`
	Op    string      `json:"op"`
	Topic string      `json:"topic"`
	Ts    interface{} `json:"ts"`
	Data  []struct {
		Symbol         string `json:"symbol"`
		ContractCode   string `json:"contract_code"`
		FeeAsset       string `json:"fee_asset"`
		FundingTime    string `json:"funding_time"`
		FundingRate    string `json:"funding_rate"`
		EstimatedRate  string `json:"estimated_rate"`
		SettlementTime string `json:"settlement_time"`
	} `json:"data"`
}
