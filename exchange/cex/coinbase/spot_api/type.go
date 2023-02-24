package spot_api

import (
	"encoding/json"
)

type RespBook struct {
	Bids        []DepthInfo `json:"bids"`
	Asks        []DepthInfo `json:"asks"`
	Sequence    int64       `json:"sequence"`
	AuctionMode bool        `json:"auction_mode"`
	Auction     bool        `json:"auction"`
}

//0-Price, 1-Amount, 2-Number of orders [List of size 3]
//type DepthInfo [3]string
type DepthInfo struct {
	Price          string
	Size           string
	NumberOfOrders int
	OrderID        string
}

func (d *DepthInfo) UnmarshalJSON(data []byte) error {
	var tmp []json.RawMessage
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[0], &d.Price); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[1], &d.Size); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[2], &d.NumberOfOrders); err != nil {
		return err
	}
	if len(tmp) > 3 {
		if err := json.Unmarshal(tmp[3], &d.OrderID); err != nil {
			return err
		}
	}
	return nil
}

type RespProducts struct {
	ProductInfos []ProductInfo
}

func (r *RespProducts) UnmarshalJSON(data []byte) error {
	var (
		tmp  []json.RawMessage
		temp ProductInfo
	)
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	for _, next := range tmp {
		if err := json.Unmarshal(next, &temp); err != nil {
			return err
		}
		r.ProductInfos = append(r.ProductInfos, temp)
	}
	return nil
}

type ProductInfo struct {
	ID              string `json:"id"`
	BaseCurrency    string `json:"base_currency"`
	QuoteCurrency   string `json:"quote_currency"`
	BaseMinSize     string `json:"base_min_size"`
	BaseMaxSize     string `json:"base_max_size"`
	QuoteIncrement  string `json:"quote_increment"`
	BaseIncrement   string `json:"base_increment"`
	DisplayName     string `json:"display_name"`
	MinMarketFunds  string `json:"min_market_funds"`
	MaxMarketFunds  string `json:"max_market_funds"`
	MarginEnabled   bool   `json:"margin_enabled"`
	PostOnly        bool   `json:"post_only"`
	LimitOnly       bool   `json:"limit_only"`
	CancelOnly      bool   `json:"cancel_only"`
	TradingDisabled bool   `json:"trading_disabled"`
	Status          string `json:"status"`
	StatusMessage   string `json:"status_message"`
}
