// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.9
// source: common/constant.proto

package common

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Exchange int32

const (
	Exchange_INVALID_EXCHANGE      Exchange = 0
	Exchange_BINANCE               Exchange = 1
	Exchange_HUOBI                 Exchange = 2
	Exchange_OKEX                  Exchange = 3
	Exchange_FTX                   Exchange = 4
	Exchange_KUCOIN                Exchange = 5
	Exchange_BYBIT                 Exchange = 6
	Exchange_COINBASE              Exchange = 7
	Exchange_KRAKEN                Exchange = 8
	Exchange_BITBANK               Exchange = 9
	Exchange_GATE                  Exchange = 10
	Exchange_FTXUS                 Exchange = 11
	Exchange_GEMINI                Exchange = 12
	Exchange_BITFLYER              Exchange = 13
	Exchange_BITSTAMP              Exchange = 14
	Exchange_BITFINEX              Exchange = 15
	Exchange_BITHUMB               Exchange = 16
	Exchange_BITMEX                Exchange = 17
	Exchange_POLONIEX              Exchange = 18
	Exchange_BINANCEUS             Exchange = 19
	Exchange_ALL                   Exchange = 60 // 特殊标识，用于风控中代表所有CEX
	Exchange_DEX_ETH_UNISWAPV2     Exchange = 64
	Exchange_DEX_ETH_UNISWAPV3     Exchange = 65
	Exchange_DEX_ETH_SUSHISWAP     Exchange = 66
	Exchange_DEX_POLYGON_UNISWAPV2 Exchange = 67
	Exchange_DEX_POLYGON_UNISWAPV3 Exchange = 68
	//  ETH_DEX = 64;
	//  BSC_DEX = 65;
	//  AVALANCHE_DEX = 66;
	//  SOLANA_DEX = 67;
	//  FANTOM_DEX = 68;
	//  TRON_DEX = 69;
	//  POLYGON_DEX = 70;
	//  ARBITRUM_DEX = 71;
	//  CRONOS_DEX = 72;
	//  HECO_DEX = 73;
	//  OPTIMISM_DEX = 74;
	Exchange_ALL_DEX Exchange = 127 // 特殊标识，用于风控中代表所有DEX
)

// Enum value maps for Exchange.
var (
	Exchange_name = map[int32]string{
		0:   "INVALID_EXCHANGE",
		1:   "BINANCE",
		2:   "HUOBI",
		3:   "OKEX",
		4:   "FTX",
		5:   "KUCOIN",
		6:   "BYBIT",
		7:   "COINBASE",
		8:   "KRAKEN",
		9:   "BITBANK",
		10:  "GATE",
		11:  "FTXUS",
		12:  "GEMINI",
		13:  "BITFLYER",
		14:  "BITSTAMP",
		15:  "BITFINEX",
		16:  "BITHUMB",
		17:  "BITMEX",
		18:  "POLONIEX",
		19:  "BINANCEUS",
		60:  "ALL",
		64:  "DEX_ETH_UNISWAPV2",
		65:  "DEX_ETH_UNISWAPV3",
		66:  "DEX_ETH_SUSHISWAP",
		67:  "DEX_POLYGON_UNISWAPV2",
		68:  "DEX_POLYGON_UNISWAPV3",
		127: "ALL_DEX",
	}
	Exchange_value = map[string]int32{
		"INVALID_EXCHANGE":      0,
		"BINANCE":               1,
		"HUOBI":                 2,
		"OKEX":                  3,
		"FTX":                   4,
		"KUCOIN":                5,
		"BYBIT":                 6,
		"COINBASE":              7,
		"KRAKEN":                8,
		"BITBANK":               9,
		"GATE":                  10,
		"FTXUS":                 11,
		"GEMINI":                12,
		"BITFLYER":              13,
		"BITSTAMP":              14,
		"BITFINEX":              15,
		"BITHUMB":               16,
		"BITMEX":                17,
		"POLONIEX":              18,
		"BINANCEUS":             19,
		"ALL":                   60,
		"DEX_ETH_UNISWAPV2":     64,
		"DEX_ETH_UNISWAPV3":     65,
		"DEX_ETH_SUSHISWAP":     66,
		"DEX_POLYGON_UNISWAPV2": 67,
		"DEX_POLYGON_UNISWAPV3": 68,
		"ALL_DEX":               127,
	}
)

func (x Exchange) Enum() *Exchange {
	p := new(Exchange)
	*p = x
	return p
}

func (x Exchange) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Exchange) Descriptor() protoreflect.EnumDescriptor {
	return file_common_constant_proto_enumTypes[0].Descriptor()
}

func (Exchange) Type() protoreflect.EnumType {
	return &file_common_constant_proto_enumTypes[0]
}

func (x Exchange) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Exchange.Descriptor instead.
func (Exchange) EnumDescriptor() ([]byte, []int) {
	return file_common_constant_proto_rawDescGZIP(), []int{0}
}

type Chain int32

const (
	Chain_INVALID_CAHIN Chain = 0
	Chain_BTC           Chain = 1
	Chain_ETH           Chain = 2
	Chain_BSC           Chain = 3
	Chain_AVALANCHE     Chain = 4
	Chain_SOLANA        Chain = 5
	Chain_FANTOM        Chain = 6
	Chain_TRON          Chain = 7
	Chain_POLYGON       Chain = 8
	Chain_ARBITRUM      Chain = 9
	Chain_CRONOS        Chain = 10
	Chain_HECO          Chain = 11
	Chain_OKC           Chain = 12
	Chain_BNB           Chain = 13
	Chain_OPTIMISM      Chain = 14
)

// Enum value maps for Chain.
var (
	Chain_name = map[int32]string{
		0:  "INVALID_CAHIN",
		1:  "BTC",
		2:  "ETH",
		3:  "BSC",
		4:  "AVALANCHE",
		5:  "SOLANA",
		6:  "FANTOM",
		7:  "TRON",
		8:  "POLYGON",
		9:  "ARBITRUM",
		10: "CRONOS",
		11: "HECO",
		12: "OKC",
		13: "BNB",
		14: "OPTIMISM",
	}
	Chain_value = map[string]int32{
		"INVALID_CAHIN": 0,
		"BTC":           1,
		"ETH":           2,
		"BSC":           3,
		"AVALANCHE":     4,
		"SOLANA":        5,
		"FANTOM":        6,
		"TRON":          7,
		"POLYGON":       8,
		"ARBITRUM":      9,
		"CRONOS":        10,
		"HECO":          11,
		"OKC":           12,
		"BNB":           13,
		"OPTIMISM":      14,
	}
)

func (x Chain) Enum() *Chain {
	p := new(Chain)
	*p = x
	return p
}

func (x Chain) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Chain) Descriptor() protoreflect.EnumDescriptor {
	return file_common_constant_proto_enumTypes[1].Descriptor()
}

func (Chain) Type() protoreflect.EnumType {
	return &file_common_constant_proto_enumTypes[1]
}

func (x Chain) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Chain.Descriptor instead.
func (Chain) EnumDescriptor() ([]byte, []int) {
	return file_common_constant_proto_rawDescGZIP(), []int{1}
}

type Market int32

const (
	Market_INVALID_MARKET Market = 0
	Market_SPOT           Market = 1
	Market_FUTURE         Market = 2 // U本位交割
	Market_SWAP           Market = 3 // U本位永续
	Market_MARGIN         Market = 4
	Market_OPTION         Market = 5
	Market_OTC            Market = 6
	Market_WALLET         Market = 7
	Market_FUTURE_COIN    Market = 8 // 币本位交割
	Market_SWAP_COIN      Market = 9 // 币本位永续
	Market_ALL_MARKET     Market = 64 // 用于风控中代表所有market
)

// Enum value maps for Market.
var (
	Market_name = map[int32]string{
		0:  "INVALID_MARKET",
		1:  "SPOT",
		2:  "FUTURE",
		3:  "SWAP",
		4:  "MARGIN",
		5:  "OPTION",
		6:  "OTC",
		7:  "WALLET",
		8:  "FUTURE_COIN",
		9:  "SWAP_COIN",
		64: "ALL_MARKET",
	}
	Market_value = map[string]int32{
		"INVALID_MARKET": 0,
		"SPOT":           1,
		"FUTURE":         2,
		"SWAP":           3,
		"MARGIN":         4,
		"OPTION":         5,
		"OTC":            6,
		"WALLET":         7,
		"FUTURE_COIN":    8,
		"SWAP_COIN":      9,
		"ALL_MARKET":     64,
	}
)

func (x Market) Enum() *Market {
	p := new(Market)
	*p = x
	return p
}

func (x Market) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Market) Descriptor() protoreflect.EnumDescriptor {
	return file_common_constant_proto_enumTypes[2].Descriptor()
}

func (Market) Type() protoreflect.EnumType {
	return &file_common_constant_proto_enumTypes[2]
}

func (x Market) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Market.Descriptor instead.
func (Market) EnumDescriptor() ([]byte, []int) {
	return file_common_constant_proto_rawDescGZIP(), []int{2}
}

// 注意格式，错了出大事：SymbolType = Market + xxx
type SymbolType int32

const (
	SymbolType_INVALID_TYPE SymbolType = 0
	// 普通现货
	SymbolType_SPOT_NORMAL SymbolType = 1
	// 杠杆
	SymbolType_MARGIN_NORMAL   SymbolType = 11 // 全仓
	SymbolType_MARGIN_ISOLATED SymbolType = 12 //逐仓
	// 期货
	SymbolType_FUTURE_COIN_THIS_WEEK    SymbolType = 101 // 币
	SymbolType_FUTURE_THIS_WEEK         SymbolType = 102 // U
	SymbolType_FUTURE_COIN_NEXT_WEEK    SymbolType = 103
	SymbolType_FUTURE_NEXT_WEEK         SymbolType = 104
	SymbolType_FUTURE_COIN_THIS_MONTH   SymbolType = 201
	SymbolType_FUTURE_THIS_MONTH        SymbolType = 202
	SymbolType_FUTURE_COIN_NEXT_MONTH   SymbolType = 203
	SymbolType_FUTURE_NEXT_MONTH        SymbolType = 204
	SymbolType_FUTURE_COIN_THIS_QUARTER SymbolType = 301
	SymbolType_FUTURE_THIS_QUARTER      SymbolType = 302
	SymbolType_FUTURE_COIN_NEXT_QUARTER SymbolType = 303
	SymbolType_FUTURE_NEXT_QUARTER      SymbolType = 304
	//合约
	SymbolType_SWAP_COIN_FOREVER SymbolType = 1001
	SymbolType_SWAP_FOREVER      SymbolType = 1002
	//OTC
	SymbolType_OTC_NORMAL SymbolType = 2001
	//WALLET
	SymbolType_WALLET_NORMAL SymbolType = 3001
)

// Enum value maps for SymbolType.
var (
	SymbolType_name = map[int32]string{
		0:    "INVALID_TYPE",
		1:    "SPOT_NORMAL",
		11:   "MARGIN_NORMAL",
		12:   "MARGIN_ISOLATED",
		101:  "FUTURE_COIN_THIS_WEEK",
		102:  "FUTURE_THIS_WEEK",
		103:  "FUTURE_COIN_NEXT_WEEK",
		104:  "FUTURE_NEXT_WEEK",
		201:  "FUTURE_COIN_THIS_MONTH",
		202:  "FUTURE_THIS_MONTH",
		203:  "FUTURE_COIN_NEXT_MONTH",
		204:  "FUTURE_NEXT_MONTH",
		301:  "FUTURE_COIN_THIS_QUARTER",
		302:  "FUTURE_THIS_QUARTER",
		303:  "FUTURE_COIN_NEXT_QUARTER",
		304:  "FUTURE_NEXT_QUARTER",
		1001: "SWAP_COIN_FOREVER",
		1002: "SWAP_FOREVER",
		2001: "OTC_NORMAL",
		3001: "WALLET_NORMAL",
	}
	SymbolType_value = map[string]int32{
		"INVALID_TYPE":             0,
		"SPOT_NORMAL":              1,
		"MARGIN_NORMAL":            11,
		"MARGIN_ISOLATED":          12,
		"FUTURE_COIN_THIS_WEEK":    101,
		"FUTURE_THIS_WEEK":         102,
		"FUTURE_COIN_NEXT_WEEK":    103,
		"FUTURE_NEXT_WEEK":         104,
		"FUTURE_COIN_THIS_MONTH":   201,
		"FUTURE_THIS_MONTH":        202,
		"FUTURE_COIN_NEXT_MONTH":   203,
		"FUTURE_NEXT_MONTH":        204,
		"FUTURE_COIN_THIS_QUARTER": 301,
		"FUTURE_THIS_QUARTER":      302,
		"FUTURE_COIN_NEXT_QUARTER": 303,
		"FUTURE_NEXT_QUARTER":      304,
		"SWAP_COIN_FOREVER":        1001,
		"SWAP_FOREVER":             1002,
		"OTC_NORMAL":               2001,
		"WALLET_NORMAL":            3001,
	}
)

func (x SymbolType) Enum() *SymbolType {
	p := new(SymbolType)
	*p = x
	return p
}

func (x SymbolType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SymbolType) Descriptor() protoreflect.EnumDescriptor {
	return file_common_constant_proto_enumTypes[3].Descriptor()
}

func (SymbolType) Type() protoreflect.EnumType {
	return &file_common_constant_proto_enumTypes[3]
}

func (x SymbolType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SymbolType.Descriptor instead.
func (SymbolType) EnumDescriptor() ([]byte, []int) {
	return file_common_constant_proto_rawDescGZIP(), []int{3}
}

type StreamType int32

const (
	StreamType_INVLAID_STREAM StreamType = 0
	StreamType_MARKET_DEPTH   StreamType = 1
	StreamType_MARKET_TRADE   StreamType = 2
	StreamType_MARKET_KLINE   StreamType = 3
	StreamType_MARKET_TICKER  StreamType = 4
	StreamType_USER_BALANCE   StreamType = 64
	StreamType_USER_ORDER     StreamType = 65
)

// Enum value maps for StreamType.
var (
	StreamType_name = map[int32]string{
		0:  "INVLAID_STREAM",
		1:  "MARKET_DEPTH",
		2:  "MARKET_TRADE",
		3:  "MARKET_KLINE",
		4:  "MARKET_TICKER",
		64: "USER_BALANCE",
		65: "USER_ORDER",
	}
	StreamType_value = map[string]int32{
		"INVLAID_STREAM": 0,
		"MARKET_DEPTH":   1,
		"MARKET_TRADE":   2,
		"MARKET_KLINE":   3,
		"MARKET_TICKER":  4,
		"USER_BALANCE":   64,
		"USER_ORDER":     65,
	}
)

func (x StreamType) Enum() *StreamType {
	p := new(StreamType)
	*p = x
	return p
}

func (x StreamType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StreamType) Descriptor() protoreflect.EnumDescriptor {
	return file_common_constant_proto_enumTypes[4].Descriptor()
}

func (StreamType) Type() protoreflect.EnumType {
	return &file_common_constant_proto_enumTypes[4]
}

func (x StreamType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StreamType.Descriptor instead.
func (StreamType) EnumDescriptor() ([]byte, []int) {
	return file_common_constant_proto_rawDescGZIP(), []int{4}
}

type PostionSide int32

const (
	PostionSide_INVALID_POSTION_SIDE PostionSide = 0
	PostionSide_LONG                 PostionSide = 1
	PostionSide_SHORT                PostionSide = 2
)

// Enum value maps for PostionSide.
var (
	PostionSide_name = map[int32]string{
		0: "INVALID_POSTION_SIDE",
		1: "LONG",
		2: "SHORT",
	}
	PostionSide_value = map[string]int32{
		"INVALID_POSTION_SIDE": 0,
		"LONG":                 1,
		"SHORT":                2,
	}
)

func (x PostionSide) Enum() *PostionSide {
	p := new(PostionSide)
	*p = x
	return p
}

func (x PostionSide) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PostionSide) Descriptor() protoreflect.EnumDescriptor {
	return file_common_constant_proto_enumTypes[5].Descriptor()
}

func (PostionSide) Type() protoreflect.EnumType {
	return &file_common_constant_proto_enumTypes[5]
}

func (x PostionSide) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PostionSide.Descriptor instead.
func (PostionSide) EnumDescriptor() ([]byte, []int) {
	return file_common_constant_proto_rawDescGZIP(), []int{5}
}

var File_common_constant_proto protoreflect.FileDescriptor

var file_common_constant_proto_rawDesc = []byte{
	0x0a, 0x15, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x9b, 0x03, 0x0a, 0x08, 0x45, 0x78, 0x63, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x10, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f,
	0x45, 0x58, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x42, 0x49,
	0x4e, 0x41, 0x4e, 0x43, 0x45, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x48, 0x55, 0x4f, 0x42, 0x49,
	0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x4f, 0x4b, 0x45, 0x58, 0x10, 0x03, 0x12, 0x07, 0x0a, 0x03,
	0x46, 0x54, 0x58, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x4b, 0x55, 0x43, 0x4f, 0x49, 0x4e, 0x10,
	0x05, 0x12, 0x09, 0x0a, 0x05, 0x42, 0x59, 0x42, 0x49, 0x54, 0x10, 0x06, 0x12, 0x0c, 0x0a, 0x08,
	0x43, 0x4f, 0x49, 0x4e, 0x42, 0x41, 0x53, 0x45, 0x10, 0x07, 0x12, 0x0a, 0x0a, 0x06, 0x4b, 0x52,
	0x41, 0x4b, 0x45, 0x4e, 0x10, 0x08, 0x12, 0x0b, 0x0a, 0x07, 0x42, 0x49, 0x54, 0x42, 0x41, 0x4e,
	0x4b, 0x10, 0x09, 0x12, 0x08, 0x0a, 0x04, 0x47, 0x41, 0x54, 0x45, 0x10, 0x0a, 0x12, 0x09, 0x0a,
	0x05, 0x46, 0x54, 0x58, 0x55, 0x53, 0x10, 0x0b, 0x12, 0x0a, 0x0a, 0x06, 0x47, 0x45, 0x4d, 0x49,
	0x4e, 0x49, 0x10, 0x0c, 0x12, 0x0c, 0x0a, 0x08, 0x42, 0x49, 0x54, 0x46, 0x4c, 0x59, 0x45, 0x52,
	0x10, 0x0d, 0x12, 0x0c, 0x0a, 0x08, 0x42, 0x49, 0x54, 0x53, 0x54, 0x41, 0x4d, 0x50, 0x10, 0x0e,
	0x12, 0x0c, 0x0a, 0x08, 0x42, 0x49, 0x54, 0x46, 0x49, 0x4e, 0x45, 0x58, 0x10, 0x0f, 0x12, 0x0b,
	0x0a, 0x07, 0x42, 0x49, 0x54, 0x48, 0x55, 0x4d, 0x42, 0x10, 0x10, 0x12, 0x0a, 0x0a, 0x06, 0x42,
	0x49, 0x54, 0x4d, 0x45, 0x58, 0x10, 0x11, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x4f, 0x4c, 0x4f, 0x4e,
	0x49, 0x45, 0x58, 0x10, 0x12, 0x12, 0x0d, 0x0a, 0x09, 0x42, 0x49, 0x4e, 0x41, 0x4e, 0x43, 0x45,
	0x55, 0x53, 0x10, 0x13, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x4c, 0x4c, 0x10, 0x3c, 0x12, 0x15, 0x0a,
	0x11, 0x44, 0x45, 0x58, 0x5f, 0x45, 0x54, 0x48, 0x5f, 0x55, 0x4e, 0x49, 0x53, 0x57, 0x41, 0x50,
	0x56, 0x32, 0x10, 0x40, 0x12, 0x15, 0x0a, 0x11, 0x44, 0x45, 0x58, 0x5f, 0x45, 0x54, 0x48, 0x5f,
	0x55, 0x4e, 0x49, 0x53, 0x57, 0x41, 0x50, 0x56, 0x33, 0x10, 0x41, 0x12, 0x15, 0x0a, 0x11, 0x44,
	0x45, 0x58, 0x5f, 0x45, 0x54, 0x48, 0x5f, 0x53, 0x55, 0x53, 0x48, 0x49, 0x53, 0x57, 0x41, 0x50,
	0x10, 0x42, 0x12, 0x19, 0x0a, 0x15, 0x44, 0x45, 0x58, 0x5f, 0x50, 0x4f, 0x4c, 0x59, 0x47, 0x4f,
	0x4e, 0x5f, 0x55, 0x4e, 0x49, 0x53, 0x57, 0x41, 0x50, 0x56, 0x32, 0x10, 0x43, 0x12, 0x19, 0x0a,
	0x15, 0x44, 0x45, 0x58, 0x5f, 0x50, 0x4f, 0x4c, 0x59, 0x47, 0x4f, 0x4e, 0x5f, 0x55, 0x4e, 0x49,
	0x53, 0x57, 0x41, 0x50, 0x56, 0x33, 0x10, 0x44, 0x12, 0x0b, 0x0a, 0x07, 0x41, 0x4c, 0x4c, 0x5f,
	0x44, 0x45, 0x58, 0x10, 0x7f, 0x2a, 0xb7, 0x01, 0x0a, 0x05, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x12,
	0x11, 0x0a, 0x0d, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x43, 0x41, 0x48, 0x49, 0x4e,
	0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x42, 0x54, 0x43, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x45,
	0x54, 0x48, 0x10, 0x02, 0x12, 0x07, 0x0a, 0x03, 0x42, 0x53, 0x43, 0x10, 0x03, 0x12, 0x0d, 0x0a,
	0x09, 0x41, 0x56, 0x41, 0x4c, 0x41, 0x4e, 0x43, 0x48, 0x45, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06,
	0x53, 0x4f, 0x4c, 0x41, 0x4e, 0x41, 0x10, 0x05, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x4e, 0x54,
	0x4f, 0x4d, 0x10, 0x06, 0x12, 0x08, 0x0a, 0x04, 0x54, 0x52, 0x4f, 0x4e, 0x10, 0x07, 0x12, 0x0b,
	0x0a, 0x07, 0x50, 0x4f, 0x4c, 0x59, 0x47, 0x4f, 0x4e, 0x10, 0x08, 0x12, 0x0c, 0x0a, 0x08, 0x41,
	0x52, 0x42, 0x49, 0x54, 0x52, 0x55, 0x4d, 0x10, 0x09, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x52, 0x4f,
	0x4e, 0x4f, 0x53, 0x10, 0x0a, 0x12, 0x08, 0x0a, 0x04, 0x48, 0x45, 0x43, 0x4f, 0x10, 0x0b, 0x12,
	0x07, 0x0a, 0x03, 0x4f, 0x4b, 0x43, 0x10, 0x0c, 0x12, 0x07, 0x0a, 0x03, 0x42, 0x4e, 0x42, 0x10,
	0x0d, 0x12, 0x0c, 0x0a, 0x08, 0x4f, 0x50, 0x54, 0x49, 0x4d, 0x49, 0x53, 0x4d, 0x10, 0x0e, 0x2a,
	0x99, 0x01, 0x0a, 0x06, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x0e, 0x49, 0x4e,
	0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x4d, 0x41, 0x52, 0x4b, 0x45, 0x54, 0x10, 0x00, 0x12, 0x08,
	0x0a, 0x04, 0x53, 0x50, 0x4f, 0x54, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x55, 0x54, 0x55,
	0x52, 0x45, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x53, 0x57, 0x41, 0x50, 0x10, 0x03, 0x12, 0x0a,
	0x0a, 0x06, 0x4d, 0x41, 0x52, 0x47, 0x49, 0x4e, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x4f, 0x50,
	0x54, 0x49, 0x4f, 0x4e, 0x10, 0x05, 0x12, 0x07, 0x0a, 0x03, 0x4f, 0x54, 0x43, 0x10, 0x06, 0x12,
	0x0a, 0x0a, 0x06, 0x57, 0x41, 0x4c, 0x4c, 0x45, 0x54, 0x10, 0x07, 0x12, 0x0f, 0x0a, 0x0b, 0x46,
	0x55, 0x54, 0x55, 0x52, 0x45, 0x5f, 0x43, 0x4f, 0x49, 0x4e, 0x10, 0x08, 0x12, 0x0d, 0x0a, 0x09,
	0x53, 0x57, 0x41, 0x50, 0x5f, 0x43, 0x4f, 0x49, 0x4e, 0x10, 0x09, 0x12, 0x0e, 0x0a, 0x0a, 0x41,
	0x4c, 0x4c, 0x5f, 0x4d, 0x41, 0x52, 0x4b, 0x45, 0x54, 0x10, 0x40, 0x2a, 0xe5, 0x03, 0x0a, 0x0a,
	0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x0c, 0x49, 0x4e,
	0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b,
	0x53, 0x50, 0x4f, 0x54, 0x5f, 0x4e, 0x4f, 0x52, 0x4d, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x11, 0x0a,
	0x0d, 0x4d, 0x41, 0x52, 0x47, 0x49, 0x4e, 0x5f, 0x4e, 0x4f, 0x52, 0x4d, 0x41, 0x4c, 0x10, 0x0b,
	0x12, 0x13, 0x0a, 0x0f, 0x4d, 0x41, 0x52, 0x47, 0x49, 0x4e, 0x5f, 0x49, 0x53, 0x4f, 0x4c, 0x41,
	0x54, 0x45, 0x44, 0x10, 0x0c, 0x12, 0x19, 0x0a, 0x15, 0x46, 0x55, 0x54, 0x55, 0x52, 0x45, 0x5f,
	0x43, 0x4f, 0x49, 0x4e, 0x5f, 0x54, 0x48, 0x49, 0x53, 0x5f, 0x57, 0x45, 0x45, 0x4b, 0x10, 0x65,
	0x12, 0x14, 0x0a, 0x10, 0x46, 0x55, 0x54, 0x55, 0x52, 0x45, 0x5f, 0x54, 0x48, 0x49, 0x53, 0x5f,
	0x57, 0x45, 0x45, 0x4b, 0x10, 0x66, 0x12, 0x19, 0x0a, 0x15, 0x46, 0x55, 0x54, 0x55, 0x52, 0x45,
	0x5f, 0x43, 0x4f, 0x49, 0x4e, 0x5f, 0x4e, 0x45, 0x58, 0x54, 0x5f, 0x57, 0x45, 0x45, 0x4b, 0x10,
	0x67, 0x12, 0x14, 0x0a, 0x10, 0x46, 0x55, 0x54, 0x55, 0x52, 0x45, 0x5f, 0x4e, 0x45, 0x58, 0x54,
	0x5f, 0x57, 0x45, 0x45, 0x4b, 0x10, 0x68, 0x12, 0x1b, 0x0a, 0x16, 0x46, 0x55, 0x54, 0x55, 0x52,
	0x45, 0x5f, 0x43, 0x4f, 0x49, 0x4e, 0x5f, 0x54, 0x48, 0x49, 0x53, 0x5f, 0x4d, 0x4f, 0x4e, 0x54,
	0x48, 0x10, 0xc9, 0x01, 0x12, 0x16, 0x0a, 0x11, 0x46, 0x55, 0x54, 0x55, 0x52, 0x45, 0x5f, 0x54,
	0x48, 0x49, 0x53, 0x5f, 0x4d, 0x4f, 0x4e, 0x54, 0x48, 0x10, 0xca, 0x01, 0x12, 0x1b, 0x0a, 0x16,
	0x46, 0x55, 0x54, 0x55, 0x52, 0x45, 0x5f, 0x43, 0x4f, 0x49, 0x4e, 0x5f, 0x4e, 0x45, 0x58, 0x54,
	0x5f, 0x4d, 0x4f, 0x4e, 0x54, 0x48, 0x10, 0xcb, 0x01, 0x12, 0x16, 0x0a, 0x11, 0x46, 0x55, 0x54,
	0x55, 0x52, 0x45, 0x5f, 0x4e, 0x45, 0x58, 0x54, 0x5f, 0x4d, 0x4f, 0x4e, 0x54, 0x48, 0x10, 0xcc,
	0x01, 0x12, 0x1d, 0x0a, 0x18, 0x46, 0x55, 0x54, 0x55, 0x52, 0x45, 0x5f, 0x43, 0x4f, 0x49, 0x4e,
	0x5f, 0x54, 0x48, 0x49, 0x53, 0x5f, 0x51, 0x55, 0x41, 0x52, 0x54, 0x45, 0x52, 0x10, 0xad, 0x02,
	0x12, 0x18, 0x0a, 0x13, 0x46, 0x55, 0x54, 0x55, 0x52, 0x45, 0x5f, 0x54, 0x48, 0x49, 0x53, 0x5f,
	0x51, 0x55, 0x41, 0x52, 0x54, 0x45, 0x52, 0x10, 0xae, 0x02, 0x12, 0x1d, 0x0a, 0x18, 0x46, 0x55,
	0x54, 0x55, 0x52, 0x45, 0x5f, 0x43, 0x4f, 0x49, 0x4e, 0x5f, 0x4e, 0x45, 0x58, 0x54, 0x5f, 0x51,
	0x55, 0x41, 0x52, 0x54, 0x45, 0x52, 0x10, 0xaf, 0x02, 0x12, 0x18, 0x0a, 0x13, 0x46, 0x55, 0x54,
	0x55, 0x52, 0x45, 0x5f, 0x4e, 0x45, 0x58, 0x54, 0x5f, 0x51, 0x55, 0x41, 0x52, 0x54, 0x45, 0x52,
	0x10, 0xb0, 0x02, 0x12, 0x16, 0x0a, 0x11, 0x53, 0x57, 0x41, 0x50, 0x5f, 0x43, 0x4f, 0x49, 0x4e,
	0x5f, 0x46, 0x4f, 0x52, 0x45, 0x56, 0x45, 0x52, 0x10, 0xe9, 0x07, 0x12, 0x11, 0x0a, 0x0c, 0x53,
	0x57, 0x41, 0x50, 0x5f, 0x46, 0x4f, 0x52, 0x45, 0x56, 0x45, 0x52, 0x10, 0xea, 0x07, 0x12, 0x0f,
	0x0a, 0x0a, 0x4f, 0x54, 0x43, 0x5f, 0x4e, 0x4f, 0x52, 0x4d, 0x41, 0x4c, 0x10, 0xd1, 0x0f, 0x12,
	0x12, 0x0a, 0x0d, 0x57, 0x41, 0x4c, 0x4c, 0x45, 0x54, 0x5f, 0x4e, 0x4f, 0x52, 0x4d, 0x41, 0x4c,
	0x10, 0xb9, 0x17, 0x2a, 0x8b, 0x01, 0x0a, 0x0a, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x12, 0x0a, 0x0e, 0x49, 0x4e, 0x56, 0x4c, 0x41, 0x49, 0x44, 0x5f, 0x53, 0x54,
	0x52, 0x45, 0x41, 0x4d, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x4d, 0x41, 0x52, 0x4b, 0x45, 0x54,
	0x5f, 0x44, 0x45, 0x50, 0x54, 0x48, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x4d, 0x41, 0x52, 0x4b,
	0x45, 0x54, 0x5f, 0x54, 0x52, 0x41, 0x44, 0x45, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c, 0x4d, 0x41,
	0x52, 0x4b, 0x45, 0x54, 0x5f, 0x4b, 0x4c, 0x49, 0x4e, 0x45, 0x10, 0x03, 0x12, 0x11, 0x0a, 0x0d,
	0x4d, 0x41, 0x52, 0x4b, 0x45, 0x54, 0x5f, 0x54, 0x49, 0x43, 0x4b, 0x45, 0x52, 0x10, 0x04, 0x12,
	0x10, 0x0a, 0x0c, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x42, 0x41, 0x4c, 0x41, 0x4e, 0x43, 0x45, 0x10,
	0x40, 0x12, 0x0e, 0x0a, 0x0a, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x10,
	0x41, 0x2a, 0x3c, 0x0a, 0x0b, 0x50, 0x6f, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x69, 0x64, 0x65,
	0x12, 0x18, 0x0a, 0x14, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x50, 0x4f, 0x53, 0x54,
	0x49, 0x4f, 0x4e, 0x5f, 0x53, 0x49, 0x44, 0x45, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x4c, 0x4f,
	0x4e, 0x47, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x53, 0x48, 0x4f, 0x52, 0x54, 0x10, 0x02, 0x42,
	0x40, 0x0a, 0x17, 0x77, 0x61, 0x72, 0x6d, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x61, 0x72, 0x6d, 0x70, 0x6c, 0x61, 0x6e, 0x65,
	0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_constant_proto_rawDescOnce sync.Once
	file_common_constant_proto_rawDescData = file_common_constant_proto_rawDesc
)

func file_common_constant_proto_rawDescGZIP() []byte {
	file_common_constant_proto_rawDescOnce.Do(func() {
		file_common_constant_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_constant_proto_rawDescData)
	})
	return file_common_constant_proto_rawDescData
}

var file_common_constant_proto_enumTypes = make([]protoimpl.EnumInfo, 6)
var file_common_constant_proto_goTypes = []interface{}{
	(Exchange)(0),    // 0: Exchange
	(Chain)(0),       // 1: Chain
	(Market)(0),      // 2: Market
	(SymbolType)(0),  // 3: SymbolType
	(StreamType)(0),  // 4: StreamType
	(PostionSide)(0), // 5: PostionSide
}
var file_common_constant_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_common_constant_proto_init() }
func file_common_constant_proto_init() {
	if File_common_constant_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_constant_proto_rawDesc,
			NumEnums:      6,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_constant_proto_goTypes,
		DependencyIndexes: file_common_constant_proto_depIdxs,
		EnumInfos:         file_common_constant_proto_enumTypes,
	}.Build()
	File_common_constant_proto = out.File
	file_common_constant_proto_rawDesc = nil
	file_common_constant_proto_goTypes = nil
	file_common_constant_proto_depIdxs = nil
}
