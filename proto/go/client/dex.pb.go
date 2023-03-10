// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.9
// source: client/dex.proto

package client

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

type BlockType int32

const (
	BlockType_INVALID_BLOCK_TYPE BlockType = 0
	BlockType_LATEST             BlockType = 1
	BlockType_PENDING            BlockType = 2
)

// Enum value maps for BlockType.
var (
	BlockType_name = map[int32]string{
		0: "INVALID_BLOCK_TYPE",
		1: "LATEST",
		2: "PENDING",
	}
	BlockType_value = map[string]int32{
		"INVALID_BLOCK_TYPE": 0,
		"LATEST":             1,
		"PENDING":            2,
	}
)

func (x BlockType) Enum() *BlockType {
	p := new(BlockType)
	*p = x
	return p
}

func (x BlockType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BlockType) Descriptor() protoreflect.EnumDescriptor {
	return file_client_dex_proto_enumTypes[0].Descriptor()
}

func (BlockType) Type() protoreflect.EnumType {
	return &file_client_dex_proto_enumTypes[0]
}

func (x BlockType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BlockType.Descriptor instead.
func (BlockType) EnumDescriptor() ([]byte, []int) {
	return file_client_dex_proto_rawDescGZIP(), []int{0}
}

type TokenInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`                 //地址
	Token   string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`                     //token名称，大写：ETH
	IsErc20 bool   `protobuf:"varint,3,opt,name=is_erc20,json=isErc20,proto3" json:"is_erc20,omitempty"` //是否是代币
	Decimal int64  `protobuf:"varint,4,opt,name=decimal,proto3" json:"decimal,omitempty"`                //精度
	Chain   int64  `protobuf:"varint,5,opt,name=chain,proto3" json:"chain,omitempty"`                    // chain id
	Name    string `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty"`                       // token详细名称
}

func (x *TokenInfo) Reset() {
	*x = TokenInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_dex_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenInfo) ProtoMessage() {}

func (x *TokenInfo) ProtoReflect() protoreflect.Message {
	mi := &file_client_dex_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenInfo.ProtoReflect.Descriptor instead.
func (*TokenInfo) Descriptor() ([]byte, []int) {
	return file_client_dex_proto_rawDescGZIP(), []int{0}
}

func (x *TokenInfo) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *TokenInfo) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *TokenInfo) GetIsErc20() bool {
	if x != nil {
		return x.IsErc20
	}
	return false
}

func (x *TokenInfo) GetDecimal() int64 {
	if x != nil {
		return x.Decimal
	}
	return 0
}

func (x *TokenInfo) GetChain() int64 {
	if x != nil {
		return x.Chain
	}
	return 0
}

func (x *TokenInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type DexBalanceItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info    *TokenInfo `protobuf:"bytes,6,opt,name=info,proto3" json:"info,omitempty"`         //代币信息
	Balance float64    `protobuf:"fixed64,1,opt,name=balance,proto3" json:"balance,omitempty"` //余额
}

func (x *DexBalanceItem) Reset() {
	*x = DexBalanceItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_dex_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DexBalanceItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DexBalanceItem) ProtoMessage() {}

func (x *DexBalanceItem) ProtoReflect() protoreflect.Message {
	mi := &file_client_dex_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DexBalanceItem.ProtoReflect.Descriptor instead.
func (*DexBalanceItem) Descriptor() ([]byte, []int) {
	return file_client_dex_proto_rawDescGZIP(), []int{1}
}

func (x *DexBalanceItem) GetInfo() *TokenInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

func (x *DexBalanceItem) GetBalance() float64 {
	if x != nil {
		return x.Balance
	}
	return 0
}

type ReqBalance struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"` //地址
	Token   string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`     //token名称，大写：ETH
}

func (x *ReqBalance) Reset() {
	*x = ReqBalance{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_dex_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqBalance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqBalance) ProtoMessage() {}

func (x *ReqBalance) ProtoReflect() protoreflect.Message {
	mi := &file_client_dex_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqBalance.ProtoReflect.Descriptor instead.
func (*ReqBalance) Descriptor() ([]byte, []int) {
	return file_client_dex_proto_rawDescGZIP(), []int{2}
}

func (x *ReqBalance) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *ReqBalance) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type RspDexBalance struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UpdateTime  int64             `protobuf:"varint,2,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`   //更新时间 microsecond
	Owner       string            `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`                                //所有者地址
	BalanceList []*DexBalanceItem `protobuf:"bytes,6,rep,name=balance_list,json=balanceList,proto3" json:"balance_list,omitempty"` //余额信息
}

func (x *RspDexBalance) Reset() {
	*x = RspDexBalance{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_dex_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RspDexBalance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RspDexBalance) ProtoMessage() {}

func (x *RspDexBalance) ProtoReflect() protoreflect.Message {
	mi := &file_client_dex_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RspDexBalance.ProtoReflect.Descriptor instead.
func (*RspDexBalance) Descriptor() ([]byte, []int) {
	return file_client_dex_proto_rawDescGZIP(), []int{3}
}

func (x *RspDexBalance) GetUpdateTime() int64 {
	if x != nil {
		return x.UpdateTime
	}
	return 0
}

func (x *RspDexBalance) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *RspDexBalance) GetBalanceList() []*DexBalanceItem {
	if x != nil {
		return x.BalanceList
	}
	return nil
}

type TxLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address     string   `protobuf:"bytes,1,opt,name=Address,proto3" json:"Address,omitempty"`
	Topics      []string `protobuf:"bytes,2,rep,name=Topics,proto3" json:"Topics,omitempty"`
	Data        []byte   `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
	BlockNumber uint64   `protobuf:"varint,4,opt,name=BlockNumber,proto3" json:"BlockNumber,omitempty"`
	TxHash      string   `protobuf:"bytes,5,opt,name=TxHash,proto3" json:"TxHash,omitempty"`
	TxIndex     uint64   `protobuf:"varint,6,opt,name=TxIndex,proto3" json:"TxIndex,omitempty"`
	BlockHash   string   `protobuf:"bytes,7,opt,name=BlockHash,proto3" json:"BlockHash,omitempty"`
	Index       uint64   `protobuf:"varint,8,opt,name=Index,proto3" json:"Index,omitempty"`
	Removed     bool     `protobuf:"varint,9,opt,name=Removed,proto3" json:"Removed,omitempty"`
}

func (x *TxLog) Reset() {
	*x = TxLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_dex_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TxLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TxLog) ProtoMessage() {}

func (x *TxLog) ProtoReflect() protoreflect.Message {
	mi := &file_client_dex_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TxLog.ProtoReflect.Descriptor instead.
func (*TxLog) Descriptor() ([]byte, []int) {
	return file_client_dex_proto_rawDescGZIP(), []int{4}
}

func (x *TxLog) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *TxLog) GetTopics() []string {
	if x != nil {
		return x.Topics
	}
	return nil
}

func (x *TxLog) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *TxLog) GetBlockNumber() uint64 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *TxLog) GetTxHash() string {
	if x != nil {
		return x.TxHash
	}
	return ""
}

func (x *TxLog) GetTxIndex() uint64 {
	if x != nil {
		return x.TxIndex
	}
	return 0
}

func (x *TxLog) GetBlockHash() string {
	if x != nil {
		return x.BlockHash
	}
	return ""
}

func (x *TxLog) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *TxLog) GetRemoved() bool {
	if x != nil {
		return x.Removed
	}
	return false
}

type RspTxInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status            bool   `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`                       //执行结果 1表示成功，0表示失败
	CumulativeGasUsed uint64 `protobuf:"varint,2,opt,name=CumulativeGasUsed,proto3" json:"CumulativeGasUsed,omitempty"` // 区块累积使用gas
	//bytes Bloom=3;//交易事件日志布隆信息
	//TxLog  Logs=4;//交易事件日志
	TxHash           string `protobuf:"bytes,5,opt,name=TxHash,proto3" json:"TxHash,omitempty"`                       //交易哈希
	ContractAddress  string `protobuf:"bytes,6,opt,name=ContractAddress,proto3" json:"ContractAddress,omitempty"`     //新合约地址
	GasUsed          uint64 `protobuf:"varint,7,opt,name=GasUsed,proto3" json:"GasUsed,omitempty"`                    //交易消耗的gas
	BlockHash        string `protobuf:"bytes,8,opt,name=BlockHash,proto3" json:"BlockHash,omitempty"`                 //交易所在区块哈希
	BlockNumber      uint64 `protobuf:"varint,9,opt,name=BlockNumber,proto3" json:"BlockNumber,omitempty"`            //交易所在区块高度
	TransactionIndex uint64 `protobuf:"varint,10,opt,name=TransactionIndex,proto3" json:"TransactionIndex,omitempty"` //交易在区块交易集中的索引
}

func (x *RspTxInfo) Reset() {
	*x = RspTxInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_dex_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RspTxInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RspTxInfo) ProtoMessage() {}

func (x *RspTxInfo) ProtoReflect() protoreflect.Message {
	mi := &file_client_dex_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RspTxInfo.ProtoReflect.Descriptor instead.
func (*RspTxInfo) Descriptor() ([]byte, []int) {
	return file_client_dex_proto_rawDescGZIP(), []int{5}
}

func (x *RspTxInfo) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

func (x *RspTxInfo) GetCumulativeGasUsed() uint64 {
	if x != nil {
		return x.CumulativeGasUsed
	}
	return 0
}

func (x *RspTxInfo) GetTxHash() string {
	if x != nil {
		return x.TxHash
	}
	return ""
}

func (x *RspTxInfo) GetContractAddress() string {
	if x != nil {
		return x.ContractAddress
	}
	return ""
}

func (x *RspTxInfo) GetGasUsed() uint64 {
	if x != nil {
		return x.GasUsed
	}
	return 0
}

func (x *RspTxInfo) GetBlockHash() string {
	if x != nil {
		return x.BlockHash
	}
	return ""
}

func (x *RspTxInfo) GetBlockNumber() uint64 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *RspTxInfo) GetTransactionIndex() uint64 {
	if x != nil {
		return x.TransactionIndex
	}
	return 0
}

type ReqBlockInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type     BlockType `protobuf:"varint,1,opt,name=type,proto3,enum=client.BlockType" json:"type,omitempty"`
	BlockNum uint64    `protobuf:"varint,2,opt,name=block_num,json=blockNum,proto3" json:"block_num,omitempty"` //type为0，查询指定block
}

func (x *ReqBlockInfo) Reset() {
	*x = ReqBlockInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_dex_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqBlockInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqBlockInfo) ProtoMessage() {}

func (x *ReqBlockInfo) ProtoReflect() protoreflect.Message {
	mi := &file_client_dex_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqBlockInfo.ProtoReflect.Descriptor instead.
func (*ReqBlockInfo) Descriptor() ([]byte, []int) {
	return file_client_dex_proto_rawDescGZIP(), []int{6}
}

func (x *ReqBlockInfo) GetType() BlockType {
	if x != nil {
		return x.Type
	}
	return BlockType_INVALID_BLOCK_TYPE
}

func (x *ReqBlockInfo) GetBlockNum() uint64 {
	if x != nil {
		return x.BlockNum
	}
	return 0
}

type RspBlockInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number     uint64 `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`        // - Number: The block number. null when its pending block.
	Hash       string `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`             // 32 Bytes - String: Hash of the block. null when its pending block.
	ParentHash string `protobuf:"bytes,3,opt,name=parentHash,proto3" json:"parentHash,omitempty"` // 32 Bytes - String: Hash of the parent block.
	GasLimit   uint64 `protobuf:"varint,4,opt,name=gasLimit,proto3" json:"gasLimit,omitempty"`    // - Number: The maximum gas allowed in this block.
	GasUsed    uint64 `protobuf:"varint,5,opt,name=gasUsed,proto3" json:"gasUsed,omitempty"`      // - Number: The total used gas by all transactions in this block.
	Timestamp  int64  `protobuf:"varint,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`  // - Number: The unix timestamp for when the block was collated.
}

func (x *RspBlockInfo) Reset() {
	*x = RspBlockInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_dex_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RspBlockInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RspBlockInfo) ProtoMessage() {}

func (x *RspBlockInfo) ProtoReflect() protoreflect.Message {
	mi := &file_client_dex_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RspBlockInfo.ProtoReflect.Descriptor instead.
func (*RspBlockInfo) Descriptor() ([]byte, []int) {
	return file_client_dex_proto_rawDescGZIP(), []int{7}
}

func (x *RspBlockInfo) GetNumber() uint64 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *RspBlockInfo) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *RspBlockInfo) GetParentHash() string {
	if x != nil {
		return x.ParentHash
	}
	return ""
}

func (x *RspBlockInfo) GetGasLimit() uint64 {
	if x != nil {
		return x.GasLimit
	}
	return 0
}

func (x *RspBlockInfo) GetGasUsed() uint64 {
	if x != nil {
		return x.GasUsed
	}
	return 0
}

func (x *RspBlockInfo) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

var File_client_dex_proto protoreflect.FileDescriptor

var file_client_dex_proto_rawDesc = []byte{
	0x0a, 0x10, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x64, 0x65, 0x78, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x22, 0x9a, 0x01, 0x0a, 0x09, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x73, 0x5f, 0x65,
	0x72, 0x63, 0x32, 0x30, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x45, 0x72,
	0x63, 0x32, 0x30, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x12, 0x14, 0x0a,
	0x05, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x51, 0x0a, 0x0e, 0x44, 0x65, 0x78, 0x42, 0x61,
	0x6c, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x25, 0x0a, 0x04, 0x69, 0x6e, 0x66,
	0x6f, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f,
	0x12, 0x18, 0x0a, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x22, 0x3c, 0x0a, 0x0a, 0x52, 0x65,
	0x71, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x81, 0x01, 0x0a, 0x0d, 0x52, 0x73, 0x70,
	0x44, 0x65, 0x78, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6f,
	0x77, 0x6e, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x12, 0x39, 0x0a, 0x0c, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x6c, 0x69, 0x73,
	0x74, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x2e, 0x44, 0x65, 0x78, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52,
	0x0b, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x22, 0xef, 0x01, 0x0a,
	0x05, 0x54, 0x78, 0x4c, 0x6f, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x06, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x20, 0x0a, 0x0b,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0b, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x54, 0x78, 0x48, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x54, 0x78, 0x48, 0x61, 0x73, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x54, 0x78, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x54, 0x78, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x12, 0x1c, 0x0a, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x14,
	0x0a, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x12, 0x18, 0x0a, 0x07, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x22, 0x99,
	0x02, 0x0a, 0x09, 0x52, 0x73, 0x70, 0x54, 0x78, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x2c, 0x0a, 0x11, 0x43, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69,
	0x76, 0x65, 0x47, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x11, 0x43, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x76, 0x65, 0x47, 0x61, 0x73, 0x55, 0x73,
	0x65, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x54, 0x78, 0x48, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x54, 0x78, 0x48, 0x61, 0x73, 0x68, 0x12, 0x28, 0x0a, 0x0f, 0x43, 0x6f,
	0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x47, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x47, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x12, 0x1c,
	0x0a, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x20, 0x0a, 0x0b,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0b, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x2a,
	0x0a, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x52, 0x0a, 0x0c, 0x52, 0x65,
	0x71, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x25, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x22, 0xae,
	0x01, 0x0a, 0x0c, 0x52, 0x73, 0x70, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x1e, 0x0a, 0x0a, 0x70,
	0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x67,
	0x61, 0x73, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x67,
	0x61, 0x73, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x67, 0x61, 0x73, 0x55, 0x73,
	0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x67, 0x61, 0x73, 0x55, 0x73, 0x65,
	0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2a,
	0x3c, 0x0a, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x12,
	0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x4c, 0x41, 0x54, 0x45, 0x53, 0x54, 0x10, 0x01,
	0x12, 0x0b, 0x0a, 0x07, 0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x42, 0x40, 0x0a,
	0x17, 0x77, 0x61, 0x72, 0x6d, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x61, 0x72, 0x6d, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x74, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_client_dex_proto_rawDescOnce sync.Once
	file_client_dex_proto_rawDescData = file_client_dex_proto_rawDesc
)

func file_client_dex_proto_rawDescGZIP() []byte {
	file_client_dex_proto_rawDescOnce.Do(func() {
		file_client_dex_proto_rawDescData = protoimpl.X.CompressGZIP(file_client_dex_proto_rawDescData)
	})
	return file_client_dex_proto_rawDescData
}

var file_client_dex_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_client_dex_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_client_dex_proto_goTypes = []interface{}{
	(BlockType)(0),         // 0: client.BlockType
	(*TokenInfo)(nil),      // 1: client.TokenInfo
	(*DexBalanceItem)(nil), // 2: client.DexBalanceItem
	(*ReqBalance)(nil),     // 3: client.ReqBalance
	(*RspDexBalance)(nil),  // 4: client.RspDexBalance
	(*TxLog)(nil),          // 5: client.TxLog
	(*RspTxInfo)(nil),      // 6: client.RspTxInfo
	(*ReqBlockInfo)(nil),   // 7: client.ReqBlockInfo
	(*RspBlockInfo)(nil),   // 8: client.RspBlockInfo
}
var file_client_dex_proto_depIdxs = []int32{
	1, // 0: client.DexBalanceItem.info:type_name -> client.TokenInfo
	2, // 1: client.RspDexBalance.balance_list:type_name -> client.DexBalanceItem
	0, // 2: client.ReqBlockInfo.type:type_name -> client.BlockType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_client_dex_proto_init() }
func file_client_dex_proto_init() {
	if File_client_dex_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_client_dex_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenInfo); i {
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
		file_client_dex_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DexBalanceItem); i {
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
		file_client_dex_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqBalance); i {
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
		file_client_dex_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RspDexBalance); i {
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
		file_client_dex_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TxLog); i {
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
		file_client_dex_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RspTxInfo); i {
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
		file_client_dex_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqBlockInfo); i {
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
		file_client_dex_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RspBlockInfo); i {
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
			RawDescriptor: file_client_dex_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_client_dex_proto_goTypes,
		DependencyIndexes: file_client_dex_proto_depIdxs,
		EnumInfos:         file_client_dex_proto_enumTypes,
		MessageInfos:      file_client_dex_proto_msgTypes,
	}.Build()
	File_client_dex_proto = out.File
	file_client_dex_proto_rawDesc = nil
	file_client_dex_proto_goTypes = nil
	file_client_dex_proto_depIdxs = nil
}
