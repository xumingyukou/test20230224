// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.2
// source: client/header.proto

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

type ClientHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WeightUsed  int64  `protobuf:"varint,1,opt,name=weight_used,json=weightUsed,proto3" json:"weight_used,omitempty"`    // 已经使用的权重
	WeightType  string `protobuf:"bytes,2,opt,name=weight_type,json=weightType,proto3" json:"weight_type,omitempty"`     // 权重类型：ip/uid等，每秒/分钟/小时
	WeightValue int64  `protobuf:"varint,3,opt,name=weight_value,json=weightValue,proto3" json:"weight_value,omitempty"` // 单次请求权重消耗值
}

func (x *ClientHeader) Reset() {
	*x = ClientHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_header_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientHeader) ProtoMessage() {}

func (x *ClientHeader) ProtoReflect() protoreflect.Message {
	mi := &file_client_header_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientHeader.ProtoReflect.Descriptor instead.
func (*ClientHeader) Descriptor() ([]byte, []int) {
	return file_client_header_proto_rawDescGZIP(), []int{0}
}

func (x *ClientHeader) GetWeightUsed() int64 {
	if x != nil {
		return x.WeightUsed
	}
	return 0
}

func (x *ClientHeader) GetWeightType() string {
	if x != nil {
		return x.WeightType
	}
	return ""
}

func (x *ClientHeader) GetWeightValue() int64 {
	if x != nil {
		return x.WeightValue
	}
	return 0
}

var File_client_header_proto protoreflect.FileDescriptor

var file_client_header_proto_rawDesc = []byte{
	0x0a, 0x13, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x22, 0x73, 0x0a,
	0x0c, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x1f, 0x0a,
	0x0b, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x55, 0x73, 0x65, 0x64, 0x12, 0x1f,
	0x0a, 0x0b, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x21, 0x0a, 0x0c, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x42, 0x40, 0x0a, 0x17, 0x77, 0x61, 0x72, 0x6d, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5a, 0x25, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x61, 0x72, 0x6d, 0x70, 0x6c,
	0x61, 0x6e, 0x65, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_client_header_proto_rawDescOnce sync.Once
	file_client_header_proto_rawDescData = file_client_header_proto_rawDesc
)

func file_client_header_proto_rawDescGZIP() []byte {
	file_client_header_proto_rawDescOnce.Do(func() {
		file_client_header_proto_rawDescData = protoimpl.X.CompressGZIP(file_client_header_proto_rawDescData)
	})
	return file_client_header_proto_rawDescData
}

var file_client_header_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_client_header_proto_goTypes = []interface{}{
	(*ClientHeader)(nil), // 0: client.ClientHeader
}
var file_client_header_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_client_header_proto_init() }
func file_client_header_proto_init() {
	if File_client_header_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_client_header_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientHeader); i {
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
			RawDescriptor: file_client_header_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_client_header_proto_goTypes,
		DependencyIndexes: file_client_header_proto_depIdxs,
		MessageInfos:      file_client_header_proto_msgTypes,
	}.Build()
	File_client_header_proto = out.File
	file_client_header_proto_rawDesc = nil
	file_client_header_proto_goTypes = nil
	file_client_header_proto_depIdxs = nil
}
