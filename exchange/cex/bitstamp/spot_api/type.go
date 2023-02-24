package spot_api

import "encoding/json"

type RespError struct {
	Code   int32       `json:"code"`
	Errors interface{} `json:"errors"`
}

type RespWSToken struct {
	Token    string `json:"token"`
	ValidSec string `json:"valid_sec"`
	UserId   string `json:"user_id"`
}

type RespTradingPairs struct {
	TradingPairs []TradingPairItem
}

func (r *RespTradingPairs) UnmarshalJSON(data []byte) error {
	var (
		tmp  []json.RawMessage
		temp TradingPairItem
	)
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	for _, next := range tmp {
		if err := json.Unmarshal(next, &temp); err != nil {
			return err
		}
		r.TradingPairs = append(r.TradingPairs, temp)
	}
	return nil
}

type TradingPairItem struct {
	Trading                string `json:"trading"`
	BaseDecimals           int    `json:"base_decimals"`
	UrlSymbol              string `json:"url_symbol"`
	Name                   string `json:"name"`
	InstantAndMarketOrders string `json:"instant_and_market_orders"`
	MinimumOrder           string `json:"minimum_order"`
	CounterDecimals        int    `json:"counter_decimals"`
	Description            string `json:"description"`
}

type RespOrderbook struct {
	Timestamp      string       `json:"timestamp"`
	Microtimestamp string       `json:"microtimestamp"`
	Bids           []*DepthItem `json:"bids"`
	Asks           []*DepthItem `json:"asks"`
}

type DepthItem struct {
	Price  string
	Amount string
}

func (d *DepthItem) UnmarshalJSON(data []byte) error {
	var tmp []json.RawMessage
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[0], &d.Price); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[1], &d.Amount); err != nil {
		return err
	}
	return nil
}
