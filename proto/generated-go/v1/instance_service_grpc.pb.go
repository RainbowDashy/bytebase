// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: instance_service.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	InstanceService_GetInstance_FullMethodName       = "/bytebase.v1.InstanceService/GetInstance"
	InstanceService_ListInstances_FullMethodName     = "/bytebase.v1.InstanceService/ListInstances"
	InstanceService_SearchInstances_FullMethodName   = "/bytebase.v1.InstanceService/SearchInstances"
	InstanceService_CreateInstance_FullMethodName    = "/bytebase.v1.InstanceService/CreateInstance"
	InstanceService_UpdateInstance_FullMethodName    = "/bytebase.v1.InstanceService/UpdateInstance"
	InstanceService_DeleteInstance_FullMethodName    = "/bytebase.v1.InstanceService/DeleteInstance"
	InstanceService_UndeleteInstance_FullMethodName  = "/bytebase.v1.InstanceService/UndeleteInstance"
	InstanceService_SyncInstance_FullMethodName      = "/bytebase.v1.InstanceService/SyncInstance"
	InstanceService_BatchSyncInstance_FullMethodName = "/bytebase.v1.InstanceService/BatchSyncInstance"
	InstanceService_AddDataSource_FullMethodName     = "/bytebase.v1.InstanceService/AddDataSource"
	InstanceService_RemoveDataSource_FullMethodName  = "/bytebase.v1.InstanceService/RemoveDataSource"
	InstanceService_UpdateDataSource_FullMethodName  = "/bytebase.v1.InstanceService/UpdateDataSource"
	InstanceService_SyncSlowQueries_FullMethodName   = "/bytebase.v1.InstanceService/SyncSlowQueries"
)

// InstanceServiceClient is the client API for InstanceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InstanceServiceClient interface {
	GetInstance(ctx context.Context, in *GetInstanceRequest, opts ...grpc.CallOption) (*Instance, error)
	ListInstances(ctx context.Context, in *ListInstancesRequest, opts ...grpc.CallOption) (*ListInstancesResponse, error)
	SearchInstances(ctx context.Context, in *SearchInstancesRequest, opts ...grpc.CallOption) (*SearchInstancesResponse, error)
	CreateInstance(ctx context.Context, in *CreateInstanceRequest, opts ...grpc.CallOption) (*Instance, error)
	UpdateInstance(ctx context.Context, in *UpdateInstanceRequest, opts ...grpc.CallOption) (*Instance, error)
	DeleteInstance(ctx context.Context, in *DeleteInstanceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UndeleteInstance(ctx context.Context, in *UndeleteInstanceRequest, opts ...grpc.CallOption) (*Instance, error)
	SyncInstance(ctx context.Context, in *SyncInstanceRequest, opts ...grpc.CallOption) (*SyncInstanceResponse, error)
	BatchSyncInstance(ctx context.Context, in *BatchSyncInstanceRequest, opts ...grpc.CallOption) (*BatchSyncInstanceResponse, error)
	AddDataSource(ctx context.Context, in *AddDataSourceRequest, opts ...grpc.CallOption) (*Instance, error)
	RemoveDataSource(ctx context.Context, in *RemoveDataSourceRequest, opts ...grpc.CallOption) (*Instance, error)
	UpdateDataSource(ctx context.Context, in *UpdateDataSourceRequest, opts ...grpc.CallOption) (*Instance, error)
	SyncSlowQueries(ctx context.Context, in *SyncSlowQueriesRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type instanceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewInstanceServiceClient(cc grpc.ClientConnInterface) InstanceServiceClient {
	return &instanceServiceClient{cc}
}

func (c *instanceServiceClient) GetInstance(ctx context.Context, in *GetInstanceRequest, opts ...grpc.CallOption) (*Instance, error) {
	out := new(Instance)
	err := c.cc.Invoke(ctx, InstanceService_GetInstance_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instanceServiceClient) ListInstances(ctx context.Context, in *ListInstancesRequest, opts ...grpc.CallOption) (*ListInstancesResponse, error) {
	out := new(ListInstancesResponse)
	err := c.cc.Invoke(ctx, InstanceService_ListInstances_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instanceServiceClient) SearchInstances(ctx context.Context, in *SearchInstancesRequest, opts ...grpc.CallOption) (*SearchInstancesResponse, error) {
	out := new(SearchInstancesResponse)
	err := c.cc.Invoke(ctx, InstanceService_SearchInstances_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instanceServiceClient) CreateInstance(ctx context.Context, in *CreateInstanceRequest, opts ...grpc.CallOption) (*Instance, error) {
	out := new(Instance)
	err := c.cc.Invoke(ctx, InstanceService_CreateInstance_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instanceServiceClient) UpdateInstance(ctx context.Context, in *UpdateInstanceRequest, opts ...grpc.CallOption) (*Instance, error) {
	out := new(Instance)
	err := c.cc.Invoke(ctx, InstanceService_UpdateInstance_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instanceServiceClient) DeleteInstance(ctx context.Context, in *DeleteInstanceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, InstanceService_DeleteInstance_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instanceServiceClient) UndeleteInstance(ctx context.Context, in *UndeleteInstanceRequest, opts ...grpc.CallOption) (*Instance, error) {
	out := new(Instance)
	err := c.cc.Invoke(ctx, InstanceService_UndeleteInstance_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instanceServiceClient) SyncInstance(ctx context.Context, in *SyncInstanceRequest, opts ...grpc.CallOption) (*SyncInstanceResponse, error) {
	out := new(SyncInstanceResponse)
	err := c.cc.Invoke(ctx, InstanceService_SyncInstance_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instanceServiceClient) BatchSyncInstance(ctx context.Context, in *BatchSyncInstanceRequest, opts ...grpc.CallOption) (*BatchSyncInstanceResponse, error) {
	out := new(BatchSyncInstanceResponse)
	err := c.cc.Invoke(ctx, InstanceService_BatchSyncInstance_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instanceServiceClient) AddDataSource(ctx context.Context, in *AddDataSourceRequest, opts ...grpc.CallOption) (*Instance, error) {
	out := new(Instance)
	err := c.cc.Invoke(ctx, InstanceService_AddDataSource_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instanceServiceClient) RemoveDataSource(ctx context.Context, in *RemoveDataSourceRequest, opts ...grpc.CallOption) (*Instance, error) {
	out := new(Instance)
	err := c.cc.Invoke(ctx, InstanceService_RemoveDataSource_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instanceServiceClient) UpdateDataSource(ctx context.Context, in *UpdateDataSourceRequest, opts ...grpc.CallOption) (*Instance, error) {
	out := new(Instance)
	err := c.cc.Invoke(ctx, InstanceService_UpdateDataSource_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instanceServiceClient) SyncSlowQueries(ctx context.Context, in *SyncSlowQueriesRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, InstanceService_SyncSlowQueries_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InstanceServiceServer is the server API for InstanceService service.
// All implementations must embed UnimplementedInstanceServiceServer
// for forward compatibility
type InstanceServiceServer interface {
	GetInstance(context.Context, *GetInstanceRequest) (*Instance, error)
	ListInstances(context.Context, *ListInstancesRequest) (*ListInstancesResponse, error)
	SearchInstances(context.Context, *SearchInstancesRequest) (*SearchInstancesResponse, error)
	CreateInstance(context.Context, *CreateInstanceRequest) (*Instance, error)
	UpdateInstance(context.Context, *UpdateInstanceRequest) (*Instance, error)
	DeleteInstance(context.Context, *DeleteInstanceRequest) (*emptypb.Empty, error)
	UndeleteInstance(context.Context, *UndeleteInstanceRequest) (*Instance, error)
	SyncInstance(context.Context, *SyncInstanceRequest) (*SyncInstanceResponse, error)
	BatchSyncInstance(context.Context, *BatchSyncInstanceRequest) (*BatchSyncInstanceResponse, error)
	AddDataSource(context.Context, *AddDataSourceRequest) (*Instance, error)
	RemoveDataSource(context.Context, *RemoveDataSourceRequest) (*Instance, error)
	UpdateDataSource(context.Context, *UpdateDataSourceRequest) (*Instance, error)
	SyncSlowQueries(context.Context, *SyncSlowQueriesRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedInstanceServiceServer()
}

// UnimplementedInstanceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedInstanceServiceServer struct {
}

func (UnimplementedInstanceServiceServer) GetInstance(context.Context, *GetInstanceRequest) (*Instance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInstance not implemented")
}
func (UnimplementedInstanceServiceServer) ListInstances(context.Context, *ListInstancesRequest) (*ListInstancesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListInstances not implemented")
}
func (UnimplementedInstanceServiceServer) SearchInstances(context.Context, *SearchInstancesRequest) (*SearchInstancesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchInstances not implemented")
}
func (UnimplementedInstanceServiceServer) CreateInstance(context.Context, *CreateInstanceRequest) (*Instance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateInstance not implemented")
}
func (UnimplementedInstanceServiceServer) UpdateInstance(context.Context, *UpdateInstanceRequest) (*Instance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateInstance not implemented")
}
func (UnimplementedInstanceServiceServer) DeleteInstance(context.Context, *DeleteInstanceRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteInstance not implemented")
}
func (UnimplementedInstanceServiceServer) UndeleteInstance(context.Context, *UndeleteInstanceRequest) (*Instance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UndeleteInstance not implemented")
}
func (UnimplementedInstanceServiceServer) SyncInstance(context.Context, *SyncInstanceRequest) (*SyncInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncInstance not implemented")
}
func (UnimplementedInstanceServiceServer) BatchSyncInstance(context.Context, *BatchSyncInstanceRequest) (*BatchSyncInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchSyncInstance not implemented")
}
func (UnimplementedInstanceServiceServer) AddDataSource(context.Context, *AddDataSourceRequest) (*Instance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddDataSource not implemented")
}
func (UnimplementedInstanceServiceServer) RemoveDataSource(context.Context, *RemoveDataSourceRequest) (*Instance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveDataSource not implemented")
}
func (UnimplementedInstanceServiceServer) UpdateDataSource(context.Context, *UpdateDataSourceRequest) (*Instance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDataSource not implemented")
}
func (UnimplementedInstanceServiceServer) SyncSlowQueries(context.Context, *SyncSlowQueriesRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncSlowQueries not implemented")
}
func (UnimplementedInstanceServiceServer) mustEmbedUnimplementedInstanceServiceServer() {}

// UnsafeInstanceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InstanceServiceServer will
// result in compilation errors.
type UnsafeInstanceServiceServer interface {
	mustEmbedUnimplementedInstanceServiceServer()
}

func RegisterInstanceServiceServer(s grpc.ServiceRegistrar, srv InstanceServiceServer) {
	s.RegisterService(&InstanceService_ServiceDesc, srv)
}

func _InstanceService_GetInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstanceServiceServer).GetInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InstanceService_GetInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstanceServiceServer).GetInstance(ctx, req.(*GetInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstanceService_ListInstances_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListInstancesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstanceServiceServer).ListInstances(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InstanceService_ListInstances_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstanceServiceServer).ListInstances(ctx, req.(*ListInstancesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstanceService_SearchInstances_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchInstancesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstanceServiceServer).SearchInstances(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InstanceService_SearchInstances_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstanceServiceServer).SearchInstances(ctx, req.(*SearchInstancesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstanceService_CreateInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstanceServiceServer).CreateInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InstanceService_CreateInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstanceServiceServer).CreateInstance(ctx, req.(*CreateInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstanceService_UpdateInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstanceServiceServer).UpdateInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InstanceService_UpdateInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstanceServiceServer).UpdateInstance(ctx, req.(*UpdateInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstanceService_DeleteInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstanceServiceServer).DeleteInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InstanceService_DeleteInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstanceServiceServer).DeleteInstance(ctx, req.(*DeleteInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstanceService_UndeleteInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UndeleteInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstanceServiceServer).UndeleteInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InstanceService_UndeleteInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstanceServiceServer).UndeleteInstance(ctx, req.(*UndeleteInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstanceService_SyncInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstanceServiceServer).SyncInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InstanceService_SyncInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstanceServiceServer).SyncInstance(ctx, req.(*SyncInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstanceService_BatchSyncInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchSyncInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstanceServiceServer).BatchSyncInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InstanceService_BatchSyncInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstanceServiceServer).BatchSyncInstance(ctx, req.(*BatchSyncInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstanceService_AddDataSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddDataSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstanceServiceServer).AddDataSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InstanceService_AddDataSource_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstanceServiceServer).AddDataSource(ctx, req.(*AddDataSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstanceService_RemoveDataSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveDataSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstanceServiceServer).RemoveDataSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InstanceService_RemoveDataSource_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstanceServiceServer).RemoveDataSource(ctx, req.(*RemoveDataSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstanceService_UpdateDataSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDataSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstanceServiceServer).UpdateDataSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InstanceService_UpdateDataSource_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstanceServiceServer).UpdateDataSource(ctx, req.(*UpdateDataSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstanceService_SyncSlowQueries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncSlowQueriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstanceServiceServer).SyncSlowQueries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InstanceService_SyncSlowQueries_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstanceServiceServer).SyncSlowQueries(ctx, req.(*SyncSlowQueriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// InstanceService_ServiceDesc is the grpc.ServiceDesc for InstanceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var InstanceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bytebase.v1.InstanceService",
	HandlerType: (*InstanceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetInstance",
			Handler:    _InstanceService_GetInstance_Handler,
		},
		{
			MethodName: "ListInstances",
			Handler:    _InstanceService_ListInstances_Handler,
		},
		{
			MethodName: "SearchInstances",
			Handler:    _InstanceService_SearchInstances_Handler,
		},
		{
			MethodName: "CreateInstance",
			Handler:    _InstanceService_CreateInstance_Handler,
		},
		{
			MethodName: "UpdateInstance",
			Handler:    _InstanceService_UpdateInstance_Handler,
		},
		{
			MethodName: "DeleteInstance",
			Handler:    _InstanceService_DeleteInstance_Handler,
		},
		{
			MethodName: "UndeleteInstance",
			Handler:    _InstanceService_UndeleteInstance_Handler,
		},
		{
			MethodName: "SyncInstance",
			Handler:    _InstanceService_SyncInstance_Handler,
		},
		{
			MethodName: "BatchSyncInstance",
			Handler:    _InstanceService_BatchSyncInstance_Handler,
		},
		{
			MethodName: "AddDataSource",
			Handler:    _InstanceService_AddDataSource_Handler,
		},
		{
			MethodName: "RemoveDataSource",
			Handler:    _InstanceService_RemoveDataSource_Handler,
		},
		{
			MethodName: "UpdateDataSource",
			Handler:    _InstanceService_UpdateDataSource_Handler,
		},
		{
			MethodName: "SyncSlowQueries",
			Handler:    _InstanceService_SyncSlowQueries_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "instance_service.proto",
}
