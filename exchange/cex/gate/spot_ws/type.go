package spot_ws

import (
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type req struct {
	Event   string   `json:"event"`
	Time    int64    `json:"time"`
	Id      int      `json:"id"`
	Channel string   `json:"channel"`
	Payload []string `json:"payload"`
	Auth    struct {
		Method string `json:"method"`
		KEY    string `json:"KEY"`
		SIGN   string `json:"SIGN"`
	} `json:"auth"`

	Pair         []string     `json:"pair"`
	Subscription subscription `json:"subscription"`
}

type subscription struct {
	Name  string `json:"name"`
	Depth int    `json:"depth,omitempty"`
}

type Resp_Info struct {
	RespError

	Time    int    `json:"time"`
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Result  struct {
		Status string `json:"status"`
	} `json:"result"`
}

type RespError struct {
	Event string `json:"event"`
	Code  string `json:"code"`
	Msg   string `json:"msg"`
}

type Level struct {
	Price  decimal.Decimal
	Volume decimal.Decimal
}

func remove(slice []Level, s int) []Level {
	return append(slice[:s], slice[s+1:]...) // preserves the order
}

type KKDepth struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol       string  `protobuf:"bytes,6,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Checksum     int64   `json:"checksum"`
	TimeExchange uint64  `protobuf:"fixed64,7,opt,name=time_exchange,json=timeExchange,proto3" json:"time_exchange,omitempty"`
	TimeReceive  uint64  `protobuf:"fixed64,8,opt,name=time_receive,json=timeReceive,proto3" json:"time_receive,omitempty"`
	TimeOperate  uint64  `protobuf:"fixed64,9,opt,name=time_operate,json=timeOperate,proto3" json:"time_operate,omitempty"`
	Bids         []Level `protobuf:"bytes,20,rep,name=bids,proto3" json:"bids,omitempty"`
	Asks         []Level `protobuf:"bytes,21,rep,name=asks,proto3" json:"asks,omitempty"`
}

type Resp_Ticker struct {
	Id      int    `json:"id"`
	Time    int    `json:"time"`
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Error   Err    `json:"error"`
	Result  struct {
		Status           string `json:"status"`
		CurrencyPair     string `json:"currency_pair"`
		Last             string `json:"last"`
		LowestAsk        string `json:"lowest_ask"`
		HighestBid       string `json:"highest_bid"`
		ChangePercentage string `json:"change_percentage"`
		BaseVolume       string `json:"base_volume"`
		QuoteVolume      string `json:"quote_volume"`
		High24H          string `json:"high_24h"`
		Low24H           string `json:"low_24h"`
	} `json:"result"`
}

type Err struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Resp_Trade struct {
	Time    int    `json:"time"`
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Error   Err    `json:"error"`
	Result  struct {
		Status       string `json:"status"`
		Id           int64  `json:"id"`
		CreateTime   int    `json:"create_time"`
		CreateTimeMs string `json:"create_time_ms"`
		Side         string `json:"side"`
		CurrencyPair string `json:"currency_pair"`
		Amount       string `json:"amount"`
		Price        string `json:"price"`
	} `json:"result"`
}

type Resp_Depth struct {
	Time    int    `json:"time"`
	Channel string `json:"channel"`
	Error   Err    `json:"error"`

	Result struct {
		Status       string     `json:"status"`
		Bids         [][]string `json:"b"`
		Asks         [][]string `json:"a"`
		TimeInMilli  int64      `json:"t"`
		E            string     `json:"e"`
		ETime        int64      `json:"E"`
		CurrencyPair string     `json:"s"`
		FirstId      int64      `json:"U"`
		LastId       int64      `json:"u"`
	} `json:"result"`

	Status string `json:"status"`
	Event  string `json:"event"`
}

type PingInfo struct {
	Time    int64  `json:"time"`
	Channel string `json:"channel"`
}
