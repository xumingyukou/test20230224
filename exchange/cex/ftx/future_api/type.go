package future_api

import (
	"time"
)

type RespListAllFutures struct {
	Success bool          `json:"success"`
	Result  []*FutureItem `json:"result"`
}

type FutureItem struct {
	Ask                 float64   `json:"ask"`
	Bid                 float64   `json:"bid"`
	Change1H            float64   `json:"change1h"`
	Change24H           float64   `json:"change24h"`
	ChangeBod           float64   `json:"changeBod"`
	VolumeUsd24H        float64   `json:"volumeUsd24h"`
	Volume              float64   `json:"volume"`
	Description         string    `json:"description"`
	Enabled             bool      `json:"enabled"`
	Expired             bool      `json:"expired"`
	Expiry              time.Time `json:"expiry"`
	Index               float64   `json:"index"`
	ImfFactor           float64   `json:"imfFactor"`
	Last                float64   `json:"last"`
	LowerBound          float64   `json:"lowerBound"`
	Mark                float64   `json:"mark"`
	Name                string    `json:"name"`
	OpenInterest        float64   `json:"openInterest"`
	OpenInterestUsd     float64   `json:"openInterestUsd"`
	Perpetual           bool      `json:"perpetual"`
	PositionLimitWeight float64   `json:"positionLimitWeight"`
	PostOnly            bool      `json:"postOnly"`
	PriceIncrement      float64   `json:"priceIncrement"`
	SizeIncrement       float64   `json:"sizeIncrement"`
	Underlying          string    `json:"underlying"`
	UpperBound          float64   `json:"upperBound"`
	Type                string    `json:"type"`
}

type RespGetFuture struct {
	Success bool        `json:"success"`
	Result  FutureItem2 `json:"result"`
}

type FutureItem2 struct {
	Ask            float64   `json:"ask"`
	Bid            float64   `json:"bid"`
	Change1H       float64   `json:"change1h"`
	Change24H      float64   `json:"change24h"`
	Description    string    `json:"description"`
	Enabled        bool      `json:"enabled"`
	Expired        bool      `json:"expired"`
	Expiry         time.Time `json:"expiry"`
	Index          float64   `json:"index"`
	Last           float64   `json:"last"`
	LowerBound     float64   `json:"lowerBound"`
	Mark           float64   `json:"mark"`
	Name           string    `json:"name"`
	Perpetual      bool      `json:"perpetual"`
	PostOnly       bool      `json:"postOnly"`
	PriceIncrement float64   `json:"priceIncrement"`
	SizeIncrement  float64   `json:"sizeIncrement"`
	Underlying     string    `json:"underlying"`
	UpperBound     float64   `json:"upperBound"`
	Type           string    `json:"type"`
}

type RespFundingRates struct {
	Success bool        `json:"success"`
	Result  []*RateItem `json:"result"`
}

type RateItem struct {
	Future string    `json:"future"`
	Rate   float64   `json:"rate"`
	Time   time.Time `json:"time"`
}
