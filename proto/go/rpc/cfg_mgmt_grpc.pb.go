// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: rpc/cfg_mgmt.proto

package rpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CfgMgmtClient is the client API for CfgMgmt service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CfgMgmtClient interface {
	// 添加/删除/获取item的最新值
	AddItems(ctx context.Context, in *CfgMgmtItemAdd, opts ...grpc.CallOption) (*CfgMgmtRsp, error)
	DelItems(ctx context.Context, in *CfgMgmtItemDel, opts ...grpc.CallOption) (*CfgMgmtRsp, error)
	GetItems(ctx context.Context, in *CfgMgmtItemGet, opts ...grpc.CallOption) (*CfgMgmtRsp, error)
	// 获取item修改历史
	GetItemHis(ctx context.Context, in *CfgMgmtHisGet, opts ...grpc.CallOption) (*CfgMgmtHisRsp, error)
	// 获取item列表
	GetKeys(ctx context.Context, in *CfgMgmtKeysGet, opts ...grpc.CallOption) (*CfgMgmtKeysRsp, error)
	// tenant/group运维接口
	CreateTenant(ctx context.Context, in *CfgMgmtTenantConfig, opts ...grpc.CallOption) (*CfgMgmtRsp, error)
	CreateGroup(ctx context.Context, in *CfgMgmtGroupConfig, opts ...grpc.CallOption) (*CfgMgmtRsp, error)
	DeleteGroup(ctx context.Context, in *CfgMgmtTenantGroup, opts ...grpc.CallOption) (*CfgMgmtRsp, error)
	BackupTenant(ctx context.Context, in *CfgMgmtTenantBackup, opts ...grpc.CallOption) (*CfgMgmtRsp, error)
}

type cfgMgmtClient struct {
	cc grpc.ClientConnInterface
}

func NewCfgMgmtClient(cc grpc.ClientConnInterface) CfgMgmtClient {
	return &cfgMgmtClient{cc}
}

func (c *cfgMgmtClient) AddItems(ctx context.Context, in *CfgMgmtItemAdd, opts ...grpc.CallOption) (*CfgMgmtRsp, error) {
	out := new(CfgMgmtRsp)
	err := c.cc.Invoke(ctx, "/rpc.CfgMgmt/AddItems", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cfgMgmtClient) DelItems(ctx context.Context, in *CfgMgmtItemDel, opts ...grpc.CallOption) (*CfgMgmtRsp, error) {
	out := new(CfgMgmtRsp)
	err := c.cc.Invoke(ctx, "/rpc.CfgMgmt/DelItems", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cfgMgmtClient) GetItems(ctx context.Context, in *CfgMgmtItemGet, opts ...grpc.CallOption) (*CfgMgmtRsp, error) {
	out := new(CfgMgmtRsp)
	err := c.cc.Invoke(ctx, "/rpc.CfgMgmt/GetItems", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cfgMgmtClient) GetItemHis(ctx context.Context, in *CfgMgmtHisGet, opts ...grpc.CallOption) (*CfgMgmtHisRsp, error) {
	out := new(CfgMgmtHisRsp)
	err := c.cc.Invoke(ctx, "/rpc.CfgMgmt/GetItemHis", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cfgMgmtClient) GetKeys(ctx context.Context, in *CfgMgmtKeysGet, opts ...grpc.CallOption) (*CfgMgmtKeysRsp, error) {
	out := new(CfgMgmtKeysRsp)
	err := c.cc.Invoke(ctx, "/rpc.CfgMgmt/GetKeys", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cfgMgmtClient) CreateTenant(ctx context.Context, in *CfgMgmtTenantConfig, opts ...grpc.CallOption) (*CfgMgmtRsp, error) {
	out := new(CfgMgmtRsp)
	err := c.cc.Invoke(ctx, "/rpc.CfgMgmt/CreateTenant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cfgMgmtClient) CreateGroup(ctx context.Context, in *CfgMgmtGroupConfig, opts ...grpc.CallOption) (*CfgMgmtRsp, error) {
	out := new(CfgMgmtRsp)
	err := c.cc.Invoke(ctx, "/rpc.CfgMgmt/CreateGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cfgMgmtClient) DeleteGroup(ctx context.Context, in *CfgMgmtTenantGroup, opts ...grpc.CallOption) (*CfgMgmtRsp, error) {
	out := new(CfgMgmtRsp)
	err := c.cc.Invoke(ctx, "/rpc.CfgMgmt/DeleteGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cfgMgmtClient) BackupTenant(ctx context.Context, in *CfgMgmtTenantBackup, opts ...grpc.CallOption) (*CfgMgmtRsp, error) {
	out := new(CfgMgmtRsp)
	err := c.cc.Invoke(ctx, "/rpc.CfgMgmt/BackupTenant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CfgMgmtServer is the server API for CfgMgmt service.
// All implementations must embed UnimplementedCfgMgmtServer
// for forward compatibility
type CfgMgmtServer interface {
	// 添加/删除/获取item的最新值
	AddItems(context.Context, *CfgMgmtItemAdd) (*CfgMgmtRsp, error)
	DelItems(context.Context, *CfgMgmtItemDel) (*CfgMgmtRsp, error)
	GetItems(context.Context, *CfgMgmtItemGet) (*CfgMgmtRsp, error)
	// 获取item修改历史
	GetItemHis(context.Context, *CfgMgmtHisGet) (*CfgMgmtHisRsp, error)
	// 获取item列表
	GetKeys(context.Context, *CfgMgmtKeysGet) (*CfgMgmtKeysRsp, error)
	// tenant/group运维接口
	CreateTenant(context.Context, *CfgMgmtTenantConfig) (*CfgMgmtRsp, error)
	CreateGroup(context.Context, *CfgMgmtGroupConfig) (*CfgMgmtRsp, error)
	DeleteGroup(context.Context, *CfgMgmtTenantGroup) (*CfgMgmtRsp, error)
	BackupTenant(context.Context, *CfgMgmtTenantBackup) (*CfgMgmtRsp, error)
	mustEmbedUnimplementedCfgMgmtServer()
}

// UnimplementedCfgMgmtServer must be embedded to have forward compatible implementations.
type UnimplementedCfgMgmtServer struct {
}

func (UnimplementedCfgMgmtServer) AddItems(context.Context, *CfgMgmtItemAdd) (*CfgMgmtRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddItems not implemented")
}
func (UnimplementedCfgMgmtServer) DelItems(context.Context, *CfgMgmtItemDel) (*CfgMgmtRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelItems not implemented")
}
func (UnimplementedCfgMgmtServer) GetItems(context.Context, *CfgMgmtItemGet) (*CfgMgmtRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItems not implemented")
}
func (UnimplementedCfgMgmtServer) GetItemHis(context.Context, *CfgMgmtHisGet) (*CfgMgmtHisRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItemHis not implemented")
}
func (UnimplementedCfgMgmtServer) GetKeys(context.Context, *CfgMgmtKeysGet) (*CfgMgmtKeysRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKeys not implemented")
}
func (UnimplementedCfgMgmtServer) CreateTenant(context.Context, *CfgMgmtTenantConfig) (*CfgMgmtRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTenant not implemented")
}
func (UnimplementedCfgMgmtServer) CreateGroup(context.Context, *CfgMgmtGroupConfig) (*CfgMgmtRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGroup not implemented")
}
func (UnimplementedCfgMgmtServer) DeleteGroup(context.Context, *CfgMgmtTenantGroup) (*CfgMgmtRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteGroup not implemented")
}
func (UnimplementedCfgMgmtServer) BackupTenant(context.Context, *CfgMgmtTenantBackup) (*CfgMgmtRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BackupTenant not implemented")
}
func (UnimplementedCfgMgmtServer) mustEmbedUnimplementedCfgMgmtServer() {}

// UnsafeCfgMgmtServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CfgMgmtServer will
// result in compilation errors.
type UnsafeCfgMgmtServer interface {
	mustEmbedUnimplementedCfgMgmtServer()
}

func RegisterCfgMgmtServer(s grpc.ServiceRegistrar, srv CfgMgmtServer) {
	s.RegisterService(&CfgMgmt_ServiceDesc, srv)
}

func _CfgMgmt_AddItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CfgMgmtItemAdd)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CfgMgmtServer).AddItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.CfgMgmt/AddItems",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CfgMgmtServer).AddItems(ctx, req.(*CfgMgmtItemAdd))
	}
	return interceptor(ctx, in, info, handler)
}

func _CfgMgmt_DelItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CfgMgmtItemDel)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CfgMgmtServer).DelItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.CfgMgmt/DelItems",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CfgMgmtServer).DelItems(ctx, req.(*CfgMgmtItemDel))
	}
	return interceptor(ctx, in, info, handler)
}

func _CfgMgmt_GetItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CfgMgmtItemGet)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CfgMgmtServer).GetItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.CfgMgmt/GetItems",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CfgMgmtServer).GetItems(ctx, req.(*CfgMgmtItemGet))
	}
	return interceptor(ctx, in, info, handler)
}

func _CfgMgmt_GetItemHis_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CfgMgmtHisGet)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CfgMgmtServer).GetItemHis(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.CfgMgmt/GetItemHis",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CfgMgmtServer).GetItemHis(ctx, req.(*CfgMgmtHisGet))
	}
	return interceptor(ctx, in, info, handler)
}

func _CfgMgmt_GetKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CfgMgmtKeysGet)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CfgMgmtServer).GetKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.CfgMgmt/GetKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CfgMgmtServer).GetKeys(ctx, req.(*CfgMgmtKeysGet))
	}
	return interceptor(ctx, in, info, handler)
}

func _CfgMgmt_CreateTenant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CfgMgmtTenantConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CfgMgmtServer).CreateTenant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.CfgMgmt/CreateTenant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CfgMgmtServer).CreateTenant(ctx, req.(*CfgMgmtTenantConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _CfgMgmt_CreateGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CfgMgmtGroupConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CfgMgmtServer).CreateGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.CfgMgmt/CreateGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CfgMgmtServer).CreateGroup(ctx, req.(*CfgMgmtGroupConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _CfgMgmt_DeleteGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CfgMgmtTenantGroup)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CfgMgmtServer).DeleteGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.CfgMgmt/DeleteGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CfgMgmtServer).DeleteGroup(ctx, req.(*CfgMgmtTenantGroup))
	}
	return interceptor(ctx, in, info, handler)
}

func _CfgMgmt_BackupTenant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CfgMgmtTenantBackup)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CfgMgmtServer).BackupTenant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.CfgMgmt/BackupTenant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CfgMgmtServer).BackupTenant(ctx, req.(*CfgMgmtTenantBackup))
	}
	return interceptor(ctx, in, info, handler)
}

// CfgMgmt_ServiceDesc is the grpc.ServiceDesc for CfgMgmt service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CfgMgmt_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.CfgMgmt",
	HandlerType: (*CfgMgmtServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddItems",
			Handler:    _CfgMgmt_AddItems_Handler,
		},
		{
			MethodName: "DelItems",
			Handler:    _CfgMgmt_DelItems_Handler,
		},
		{
			MethodName: "GetItems",
			Handler:    _CfgMgmt_GetItems_Handler,
		},
		{
			MethodName: "GetItemHis",
			Handler:    _CfgMgmt_GetItemHis_Handler,
		},
		{
			MethodName: "GetKeys",
			Handler:    _CfgMgmt_GetKeys_Handler,
		},
		{
			MethodName: "CreateTenant",
			Handler:    _CfgMgmt_CreateTenant_Handler,
		},
		{
			MethodName: "CreateGroup",
			Handler:    _CfgMgmt_CreateGroup_Handler,
		},
		{
			MethodName: "DeleteGroup",
			Handler:    _CfgMgmt_DeleteGroup_Handler,
		},
		{
			MethodName: "BackupTenant",
			Handler:    _CfgMgmt_BackupTenant_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc/cfg_mgmt.proto",
}
