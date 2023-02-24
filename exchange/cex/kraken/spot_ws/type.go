package spot_ws

import (
	"github.com/shopspring/decimal"
)

type req struct {
	Event        string       `json:"event"`
	Pair         []string     `json:"pair"`
	Subscription subscription `json:"subscription"`
}

type subscription struct {
	Name  string `json:"name"`
	Depth int    `json:"depth,omitempty"`
}

type Resp_Info struct {
	RespError
	Event string `json:"event"`
	Arg   struct {
		Channel string `json:"channel"`
	} `json:"arg"`
}

type RespError struct {
	Event string `json:"event"`
	Code  string `json:"code"`
	Msg   string `json:"msg"`
}

type Level struct {
	Price   decimal.Decimal
	Volume  decimal.Decimal
	PriceF  float64
	VolumeF float64
}

func remove(slice []Level, s int) []Level {
	return append(slice[:s], slice[s+1:]...) // preserves the order
}

type KKDepth struct {
	// state         protoimpl.MessageState
	// sizeCache     protoimpl.SizeCache
	// unknownFields protoimpl.UnknownFields

	Symbol       string  `protobuf:"bytes,6,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Checksum     int64   `json:"checksum"`
	TimeExchange int64   `protobuf:"fixed64,7,opt,name=time_exchange,json=timeExchange,proto3" json:"time_exchange,omitempty"`
	TimeReceive  int64   `protobuf:"fixed64,8,opt,name=time_receive,json=timeReceive,proto3" json:"time_receive,omitempty"`
	TimeOperate  uint64  `protobuf:"fixed64,9,opt,name=time_operate,json=timeOperate,proto3" json:"time_operate,omitempty"`
	Bids         []Level `protobuf:"bytes,20,rep,name=bids,proto3" json:"bids,omitempty"`
	Asks         []Level `protobuf:"bytes,21,rep,name=asks,proto3" json:"asks,omitempty"`
}
