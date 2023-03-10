// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.9
// source: client/global.proto

package client

import (
	common "github.com/warmplanet/proto/go/common"
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

type WeightType int32

const (
	WeightType_WeightType_INVALID WeightType = 0
	WeightType_REQUEST_WEIGHT     WeightType = 1
	WeightType_ORDERS             WeightType = 2
	WeightType_RAW_REQUESTS       WeightType = 3
)

// Enum value maps for WeightType.
var (
	WeightType_name = map[int32]string{
		0: "WeightType_INVALID",
		1: "REQUEST_WEIGHT",
		2: "ORDERS",
		3: "RAW_REQUESTS",
	}
	WeightType_value = map[string]int32{
		"WeightType_INVALID": 0,
		"REQUEST_WEIGHT":     1,
		"ORDERS":             2,
		"RAW_REQUESTS":       3,
	}
)

func (x WeightType) Enum() *WeightType {
	p := new(WeightType)
	*p = x
	return p
}

func (x WeightType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WeightType) Descriptor() protoreflect.EnumDescriptor {
	return file_client_global_proto_enumTypes[0].Descriptor()
}

func (WeightType) Type() protoreflect.EnumType {
	return &file_client_global_proto_enumTypes[0]
}

func (x WeightType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WeightType.Descriptor instead.
func (WeightType) EnumDescriptor() ([]byte, []int) {
	return file_client_global_proto_rawDescGZIP(), []int{0}
}

type TradeFeeItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol string            `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Type   common.SymbolType `protobuf:"varint,5,opt,name=type,proto3,enum=SymbolType" json:"type,omitempty"` // ??????
	Maker  float64           `protobuf:"fixed64,2,opt,name=maker,proto3" json:"maker,omitempty"`
	Taker  float64           `protobuf:"fixed64,3,opt,name=taker,proto3" json:"taker,omitempty"`
}

func (x *TradeFeeItem) Reset() {
	*x = TradeFeeItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_global_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TradeFeeItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TradeFeeItem) ProtoMessage() {}

func (x *TradeFeeItem) ProtoReflect() protoreflect.Message {
	mi := &file_client_global_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TradeFeeItem.ProtoReflect.Descriptor instead.
func (*TradeFeeItem) Descriptor() ([]byte, []int) {
	return file_client_global_proto_rawDescGZIP(), []int{0}
}

func (x *TradeFeeItem) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *TradeFeeItem) GetType() common.SymbolType {
	if x != nil {
		return x.Type
	}
	return common.SymbolType(0)
}

func (x *TradeFeeItem) GetMaker() float64 {
	if x != nil {
		return x.Maker
	}
	return 0
}

func (x *TradeFeeItem) GetTaker() float64 {
	if x != nil {
		return x.Taker
	}
	return 0
}

type TransferFeeItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token   string       `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Network common.Chain `protobuf:"varint,2,opt,name=network,proto3,enum=Chain" json:"network,omitempty"`
	Fee     float64      `protobuf:"fixed64,3,opt,name=fee,proto3" json:"fee,omitempty"`
	FeeRate float64
}

func (x *TransferFeeItem) Reset() {
	*x = TransferFeeItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_global_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransferFeeItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransferFeeItem) ProtoMessage() {}

func (x *TransferFeeItem) ProtoReflect() protoreflect.Message {
	mi := &file_client_global_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransferFeeItem.ProtoReflect.Descriptor instead.
func (*TransferFeeItem) Descriptor() ([]byte, []int) {
	return file_client_global_proto_rawDescGZIP(), []int{1}
}

func (x *TransferFeeItem) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *TransferFeeItem) GetNetwork() common.Chain {
	if x != nil {
		return x.Network
	}
	return common.Chain(0)
}

func (x *TransferFeeItem) GetFee() float64 {
	if x != nil {
		return x.Fee
	}
	return 0
}

type PrecisionItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol    string            `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Type      common.SymbolType `protobuf:"varint,5,opt,name=type,proto3,enum=SymbolType" json:"type,omitempty"` // ??????
	Amount    int64             `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
	Price     int64             `protobuf:"varint,3,opt,name=price,proto3" json:"price,omitempty"`
	AmountMin float64           `protobuf:"fixed64,4,opt,name=amount_min,json=amountMin,proto3" json:"amount_min,omitempty"`
}

func (x *PrecisionItem) Reset() {
	*x = PrecisionItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_global_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrecisionItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrecisionItem) ProtoMessage() {}

func (x *PrecisionItem) ProtoReflect() protoreflect.Message {
	mi := &file_client_global_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrecisionItem.ProtoReflect.Descriptor instead.
func (*PrecisionItem) Descriptor() ([]byte, []int) {
	return file_client_global_proto_rawDescGZIP(), []int{2}
}

func (x *PrecisionItem) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *PrecisionItem) GetType() common.SymbolType {
	if x != nil {
		return x.Type
	}
	return common.SymbolType(0)
}

func (x *PrecisionItem) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *PrecisionItem) GetPrice() int64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *PrecisionItem) GetAmountMin() float64 {
	if x != nil {
		return x.AmountMin
	}
	return 0
}

type WeightInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type        WeightType `protobuf:"varint,1,opt,name=type,proto3,enum=balance.WeightType" json:"type,omitempty"`
	Value       int64      `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
	Limit       int64      `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	IntervalSec int64      `protobuf:"varint,4,opt,name=interval_sec,json=intervalSec,proto3" json:"interval_sec,omitempty"`
}

func (x *WeightInfo) Reset() {
	*x = WeightInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_global_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WeightInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WeightInfo) ProtoMessage() {}

func (x *WeightInfo) ProtoReflect() protoreflect.Message {
	mi := &file_client_global_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WeightInfo.ProtoReflect.Descriptor instead.
func (*WeightInfo) Descriptor() ([]byte, []int) {
	return file_client_global_proto_rawDescGZIP(), []int{3}
}

func (x *WeightInfo) GetType() WeightType {
	if x != nil {
		return x.Type
	}
	return WeightType_WeightType_INVALID
}

func (x *WeightInfo) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *WeightInfo) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *WeightInfo) GetIntervalSec() int64 {
	if x != nil {
		return x.IntervalSec
	}
	return 0
}

type TradeFee struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//???????????????
	TradeFeeList []*TradeFeeItem `protobuf:"bytes,1,rep,name=trade_fee_list,json=tradeFeeList,proto3" json:"trade_fee_list,omitempty"`
}

func (x *TradeFee) Reset() {
	*x = TradeFee{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_global_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TradeFee) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TradeFee) ProtoMessage() {}

func (x *TradeFee) ProtoReflect() protoreflect.Message {
	mi := &file_client_global_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TradeFee.ProtoReflect.Descriptor instead.
func (*TradeFee) Descriptor() ([]byte, []int) {
	return file_client_global_proto_rawDescGZIP(), []int{4}
}

func (x *TradeFee) GetTradeFeeList() []*TradeFeeItem {
	if x != nil {
		return x.TradeFeeList
	}
	return nil
}

type TransferFee struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//???????????????
	TransferFeeList []*TransferFeeItem `protobuf:"bytes,1,rep,name=transfer_fee_list,json=transferFeeList,proto3" json:"transfer_fee_list,omitempty"`
}

func (x *TransferFee) Reset() {
	*x = TransferFee{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_global_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransferFee) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransferFee) ProtoMessage() {}

func (x *TransferFee) ProtoReflect() protoreflect.Message {
	mi := &file_client_global_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransferFee.ProtoReflect.Descriptor instead.
func (*TransferFee) Descriptor() ([]byte, []int) {
	return file_client_global_proto_rawDescGZIP(), []int{5}
}

func (x *TransferFee) GetTransferFeeList() []*TransferFeeItem {
	if x != nil {
		return x.TransferFeeList
	}
	return nil
}

type Precision struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PrecisionList []*PrecisionItem `protobuf:"bytes,1,rep,name=precision_list,json=precisionList,proto3" json:"precision_list,omitempty"`
}

func (x *Precision) Reset() {
	*x = Precision{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_global_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Precision) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Precision) ProtoMessage() {}

func (x *Precision) ProtoReflect() protoreflect.Message {
	mi := &file_client_global_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Precision.ProtoReflect.Descriptor instead.
func (*Precision) Descriptor() ([]byte, []int) {
	return file_client_global_proto_rawDescGZIP(), []int{6}
}

func (x *Precision) GetPrecisionList() []*PrecisionItem {
	if x != nil {
		return x.PrecisionList
	}
	return nil
}

type SymbolInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol string            `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"` // base/quote: BTC/USDT
	Name   string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`     // ??????????????????BTCUSDT_220629
	Market common.Market     `protobuf:"varint,4,opt,name=market,proto3,enum=Market" json:"market,omitempty"`
	Type   common.SymbolType `protobuf:"varint,5,opt,name=type,proto3,enum=SymbolType" json:"type,omitempty"` // ??????
}

func (x *SymbolInfo) Reset() {
	*x = SymbolInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_global_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SymbolInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SymbolInfo) ProtoMessage() {}

func (x *SymbolInfo) ProtoReflect() protoreflect.Message {
	mi := &file_client_global_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SymbolInfo.ProtoReflect.Descriptor instead.
func (*SymbolInfo) Descriptor() ([]byte, []int) {
	return file_client_global_proto_rawDescGZIP(), []int{7}
}

func (x *SymbolInfo) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *SymbolInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SymbolInfo) GetMarket() common.Market {
	if x != nil {
		return x.Market
	}
	return common.Market(0)
}

func (x *SymbolInfo) GetType() common.SymbolType {
	if x != nil {
		return x.Type
	}
	return common.SymbolType(0)
}

type MarkPriceItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol     string            `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`                            // base/quote: BTC/USDT
	Type       common.SymbolType `protobuf:"varint,5,opt,name=type,proto3,enum=SymbolType" json:"type,omitempty"`               // ??????
	UpdateTime int64             `protobuf:"varint,4,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"` // microsecond
	MarkPrice  float64           `protobuf:"fixed64,3,opt,name=mark_price,json=markPrice,proto3" json:"mark_price,omitempty"`   // ????????????
}

func (x *MarkPriceItem) Reset() {
	*x = MarkPriceItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_global_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarkPriceItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarkPriceItem) ProtoMessage() {}

func (x *MarkPriceItem) ProtoReflect() protoreflect.Message {
	mi := &file_client_global_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarkPriceItem.ProtoReflect.Descriptor instead.
func (*MarkPriceItem) Descriptor() ([]byte, []int) {
	return file_client_global_proto_rawDescGZIP(), []int{8}
}

func (x *MarkPriceItem) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *MarkPriceItem) GetType() common.SymbolType {
	if x != nil {
		return x.Type
	}
	return common.SymbolType(0)
}

func (x *MarkPriceItem) GetUpdateTime() int64 {
	if x != nil {
		return x.UpdateTime
	}
	return 0
}

func (x *MarkPriceItem) GetMarkPrice() float64 {
	if x != nil {
		return x.MarkPrice
	}
	return 0
}

type RspMarkPrice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Item []*MarkPriceItem `protobuf:"bytes,1,rep,name=item,proto3" json:"item,omitempty"`
}

func (x *RspMarkPrice) Reset() {
	*x = RspMarkPrice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_global_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RspMarkPrice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RspMarkPrice) ProtoMessage() {}

func (x *RspMarkPrice) ProtoReflect() protoreflect.Message {
	mi := &file_client_global_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RspMarkPrice.ProtoReflect.Descriptor instead.
func (*RspMarkPrice) Descriptor() ([]byte, []int) {
	return file_client_global_proto_rawDescGZIP(), []int{9}
}

func (x *RspMarkPrice) GetItem() []*MarkPriceItem {
	if x != nil {
		return x.Item
	}
	return nil
}

type TransferFeeReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token     string       `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`                          // ??????
	Network   common.Chain `protobuf:"varint,2,opt,name=network,proto3,enum=Chain" json:"network,omitempty"`          // ??????
	Value     float64      `protobuf:"fixed64,3,opt,name=value,proto3" json:"value,omitempty"`                        // ??????
	ToAddress string       `protobuf:"bytes,4,opt,name=to_address,json=toAddress,proto3" json:"to_address,omitempty"` // ????????????
}

func (x *TransferFeeReq) Reset() {
	*x = TransferFeeReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_global_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransferFeeReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransferFeeReq) ProtoMessage() {}

func (x *TransferFeeReq) ProtoReflect() protoreflect.Message {
	mi := &file_client_global_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransferFeeReq.ProtoReflect.Descriptor instead.
func (*TransferFeeReq) Descriptor() ([]byte, []int) {
	return file_client_global_proto_rawDescGZIP(), []int{10}
}

func (x *TransferFeeReq) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *TransferFeeReq) GetNetwork() common.Chain {
	if x != nil {
		return x.Network
	}
	return common.Chain(0)
}

func (x *TransferFeeReq) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *TransferFeeReq) GetToAddress() string {
	if x != nil {
		return x.ToAddress
	}
	return ""
}

var File_client_global_proto protoreflect.FileDescriptor

var file_client_global_proto_rawDesc = []byte{
	0x0a, 0x13, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x1a, 0x15,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x73, 0x0a, 0x0c, 0x54, 0x72, 0x61, 0x64, 0x65, 0x46, 0x65,
	0x65, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x1f, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x53, 0x79,
	0x6d, 0x62, 0x6f, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x6d, 0x61, 0x6b, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x6d,
	0x61, 0x6b, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x61, 0x6b, 0x65, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x05, 0x74, 0x61, 0x6b, 0x65, 0x72, 0x22, 0x5b, 0x0a, 0x0f, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x46, 0x65, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x20, 0x0a, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x06, 0x2e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x52, 0x07, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x10, 0x0a, 0x03, 0x66, 0x65, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x03, 0x66, 0x65, 0x65, 0x22, 0x95, 0x01, 0x0a, 0x0d, 0x50, 0x72, 0x65, 0x63,
	0x69, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d,
	0x62, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f,
	0x6c, 0x12, 0x1f, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x0b, 0x2e, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72,
	0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65,
	0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x6d, 0x69, 0x6e, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x4d, 0x69, 0x6e, 0x22,
	0x84, 0x01, 0x0a, 0x0a, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x27,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x62,
	0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x5f,
	0x73, 0x65, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x76, 0x61, 0x6c, 0x53, 0x65, 0x63, 0x22, 0x47, 0x0a, 0x08, 0x54, 0x72, 0x61, 0x64, 0x65, 0x46,
	0x65, 0x65, 0x12, 0x3b, 0x0a, 0x0e, 0x74, 0x72, 0x61, 0x64, 0x65, 0x5f, 0x66, 0x65, 0x65, 0x5f,
	0x6c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x62, 0x61, 0x6c,
	0x61, 0x6e, 0x63, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x64, 0x65, 0x46, 0x65, 0x65, 0x49, 0x74, 0x65,
	0x6d, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x64, 0x65, 0x46, 0x65, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x22,
	0x53, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x46, 0x65, 0x65, 0x12, 0x44,
	0x0a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x5f, 0x66, 0x65, 0x65, 0x5f, 0x6c,
	0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x62, 0x61, 0x6c, 0x61,
	0x6e, 0x63, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x46, 0x65, 0x65, 0x49,
	0x74, 0x65, 0x6d, 0x52, 0x0f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x46, 0x65, 0x65,
	0x4c, 0x69, 0x73, 0x74, 0x22, 0x4a, 0x0a, 0x09, 0x50, 0x72, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x3d, 0x0a, 0x0e, 0x70, 0x72, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x6c,
	0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x62, 0x61, 0x6c, 0x61,
	0x6e, 0x63, 0x65, 0x2e, 0x50, 0x72, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x74, 0x65,
	0x6d, 0x52, 0x0d, 0x70, 0x72, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74,
	0x22, 0x7a, 0x0a, 0x0a, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x06, 0x6d, 0x61,
	0x72, 0x6b, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x07, 0x2e, 0x4d, 0x61, 0x72,
	0x6b, 0x65, 0x74, 0x52, 0x06, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x12, 0x1f, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x53, 0x79, 0x6d, 0x62,
	0x6f, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x88, 0x01, 0x0a,
	0x0d, 0x4d, 0x61, 0x72, 0x6b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x1f, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x61, 0x72, 0x6b,
	0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x6d, 0x61,
	0x72, 0x6b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x22, 0x3a, 0x0a, 0x0c, 0x52, 0x73, 0x70, 0x4d, 0x61,
	0x72, 0x6b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x2a, 0x0a, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x2e,
	0x4d, 0x61, 0x72, 0x6b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x04, 0x69,
	0x74, 0x65, 0x6d, 0x22, 0x7d, 0x0a, 0x0e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x46,
	0x65, 0x65, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x20, 0x0a, 0x07, 0x6e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x06, 0x2e, 0x43,
	0x68, 0x61, 0x69, 0x6e, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x6f, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x2a, 0x56, 0x0a, 0x0a, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x16, 0x0a, 0x12, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x54, 0x79, 0x70, 0x65, 0x5f, 0x49,
	0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x52, 0x45, 0x51, 0x55,
	0x45, 0x53, 0x54, 0x5f, 0x57, 0x45, 0x49, 0x47, 0x48, 0x54, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06,
	0x4f, 0x52, 0x44, 0x45, 0x52, 0x53, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c, 0x52, 0x41, 0x57, 0x5f,
	0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x53, 0x10, 0x03, 0x42, 0x40, 0x0a, 0x17, 0x77, 0x61,
	0x72, 0x6d, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x77, 0x61, 0x72, 0x6d, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x74, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_client_global_proto_rawDescOnce sync.Once
	file_client_global_proto_rawDescData = file_client_global_proto_rawDesc
)

func file_client_global_proto_rawDescGZIP() []byte {
	file_client_global_proto_rawDescOnce.Do(func() {
		file_client_global_proto_rawDescData = protoimpl.X.CompressGZIP(file_client_global_proto_rawDescData)
	})
	return file_client_global_proto_rawDescData
}

var file_client_global_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_client_global_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_client_global_proto_goTypes = []interface{}{
	(WeightType)(0),         // 0: balance.WeightType
	(*TradeFeeItem)(nil),    // 1: balance.TradeFeeItem
	(*TransferFeeItem)(nil), // 2: balance.TransferFeeItem
	(*PrecisionItem)(nil),   // 3: balance.PrecisionItem
	(*WeightInfo)(nil),      // 4: balance.WeightInfo
	(*TradeFee)(nil),        // 5: balance.TradeFee
	(*TransferFee)(nil),     // 6: balance.TransferFee
	(*Precision)(nil),       // 7: balance.Precision
	(*SymbolInfo)(nil),      // 8: balance.SymbolInfo
	(*MarkPriceItem)(nil),   // 9: balance.MarkPriceItem
	(*RspMarkPrice)(nil),    // 10: balance.RspMarkPrice
	(*TransferFeeReq)(nil),  // 11: balance.TransferFeeReq
	(common.SymbolType)(0),  // 12: SymbolType
	(common.Chain)(0),       // 13: Chain
	(common.Market)(0),      // 14: Market
}
var file_client_global_proto_depIdxs = []int32{
	12, // 0: balance.TradeFeeItem.type:type_name -> SymbolType
	13, // 1: balance.TransferFeeItem.network:type_name -> Chain
	12, // 2: balance.PrecisionItem.type:type_name -> SymbolType
	0,  // 3: balance.WeightInfo.type:type_name -> balance.WeightType
	1,  // 4: balance.TradeFee.trade_fee_list:type_name -> balance.TradeFeeItem
	2,  // 5: balance.TransferFee.transfer_fee_list:type_name -> balance.TransferFeeItem
	3,  // 6: balance.Precision.precision_list:type_name -> balance.PrecisionItem
	14, // 7: balance.SymbolInfo.market:type_name -> Market
	12, // 8: balance.SymbolInfo.type:type_name -> SymbolType
	12, // 9: balance.MarkPriceItem.type:type_name -> SymbolType
	9,  // 10: balance.RspMarkPrice.item:type_name -> balance.MarkPriceItem
	13, // 11: balance.TransferFeeReq.network:type_name -> Chain
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_client_global_proto_init() }
func file_client_global_proto_init() {
	if File_client_global_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_client_global_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TradeFeeItem); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_global_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransferFeeItem); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_global_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrecisionItem); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_global_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WeightInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_global_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TradeFee); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_global_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransferFee); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_global_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Precision); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_global_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SymbolInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_global_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarkPriceItem); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_global_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RspMarkPrice); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_global_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransferFeeReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_client_global_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_client_global_proto_goTypes,
		DependencyIndexes: file_client_global_proto_depIdxs,
		EnumInfos:         file_client_global_proto_enumTypes,
		MessageInfos:      file_client_global_proto_msgTypes,
	}.Build()
	File_client_global_proto = out.File
	file_client_global_proto_rawDesc = nil
	file_client_global_proto_goTypes = nil
	file_client_global_proto_depIdxs = nil
}
