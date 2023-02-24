package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"fmt"

	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
)

type BaseMessage struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type SubscribeMessage struct {
	Id             string `json:"id"`
	Type           string `json:"type"`
	Topic          string `json:"topic"`
	Response       bool   `json:"response"`
	PrivateChannel bool   `json:"privateChannel"`
}

type UnSubscribeMessage struct {
	Id             string `json:"id"`
	Type           string `json:"type"`
	Topic          string `json:"topic"`
	Response       bool   `json:"response"`
	PrivateChannel bool   `json:"privateChannel"`
}

type PingMessage struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type TradeMessage struct {
	// Id      string `json:"id"`
	Type    string `json:"type"`
	Topic   string `json:"topic"`
	Subject string `json:"subject"`
	Data    struct {
		Sequence     string `json:"sequence"`
		Type         string `json:"type"`
		Symbol       string `json:"symbol"`
		Side         string `json:"side"`
		Price        string `json:"price"`
		Size         string `json:"size"`
		TradeId      string `json:"tradeId"`
		TakerOrderId string `json:"takerOrderId"`
		MakerOrderId string `json:"makerOrderId"`
		Time         string `json:"time"`
	} `json:"data"`
}

type TickerMessage struct {
	Type    string `json:"type"`
	Topic   string `json:"topic"`
	Subject string `json:"subject"`
	Data    struct {
		Sequence    string `json:"sequence"`
		Price       string `json:"price"`
		Size        string `json:"size"`
		BestAsk     string `json:"bestAsk"`
		BestAskSize string `json:"bestAskSize"`
		BestBid     string `json:"bestBid"`
		BestBidSize string `json:"bestBidSize"`
		Time        int64  `json:"time"`
	} `json:"data"`
}

type DepthLimitMessage struct {
	Type    string `json:"type"`
	Topic   string `json:"topic"`
	Subject string `json:"subject"`
	Data    struct {
		Asks      [][]string `json:"asks"`
		Bids      [][]string `json:"bids"`
		Timestamp int64      `json:"timestamp"`
	} `json:"data"`
}

type DepthIncrementMessage struct {
	Type    string `json:"type"`
	Topic   string `json:"topic"`
	Subject string `json:"subject"`
	Data    struct {
		SequenceStart int64  `json:"sequenceStart"`
		SequenceEnd   int64  `json:"sequenceEnd"`
		Symbol        string `json:"symbol"`
		Changes       struct {
			Asks [][]string `json:"asks"`
			Bids [][]string `json:"bids"`
		} `json:"changes"`
	} `json:"data"`
}

type KucoinIncrement struct {
	SequenceStart int64
	SequenceEnd   int64
	Symbol        string
	TimeReceive   int64
	Asks          []*KucoinIncrementLevel
	Bids          []*KucoinIncrementLevel
}
type KucoinIncrementSlice []*KucoinIncrementLevel

func (d KucoinIncrementSlice) Len() int           { return len(d) }
func (d KucoinIncrementSlice) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d KucoinIncrementSlice) Less(i, j int) bool { return d[i].Price < d[j].Price }
func (d KucoinIncrementSlice) Print() {
	msg := ""
	for _, item := range d {
		msg += fmt.Sprintf("(%f ,%f)", item.Price, item.Amount)
	}
	logger.Logger.Info(msg)
}
func (d *KucoinIncrementSlice) Reverse() {
	for i := 0; i < d.Len()/2; i++ {
		(*d)[i], (*d)[d.Len()-i-1] = (*d)[d.Len()-i-1], (*d)[i]
	}
}

type KucoinIncrementLevel struct {
	depth.DepthLevel
	SequenceId int64
}

func (inc *KucoinIncrement) toDeltaDepthUpdate() *base.DeltaDepthUpdate {
	asks, bids := []*depth.DepthLevel{}, []*depth.DepthLevel{}
	for _, ask := range inc.Asks {
		asks = append(asks, &ask.DepthLevel)
	}

	for _, bid := range inc.Bids {
		bids = append(bids, &bid.DepthLevel)
	}

	return &base.DeltaDepthUpdate{
		UpdateStartId: inc.SequenceStart,
		UpdateEndId:   inc.SequenceEnd,
		Market:        common.Market_SPOT,
		Type:          common.SymbolType_SPOT_NORMAL,
		Symbol:        inc.Symbol,
		TimeReceive:   inc.TimeReceive,
		Asks:          asks,
		Bids:          bids,
	}
}
