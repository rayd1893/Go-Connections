// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: pkg/pb/snapshot.proto

package pb

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

// ShipmentsSnapshotClient is the client API for ShipmentsSnapshot service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShipmentsSnapshotClient interface {
	GetShipment(ctx context.Context, in *GetShipmentRequest, opts ...grpc.CallOption) (*Shipment, error)
	GetShipments(ctx context.Context, in *GetShipmentsRequest, opts ...grpc.CallOption) (*GetShipmentsResponse, error)
	GetEvents(ctx context.Context, in *GetEventsRequest, opts ...grpc.CallOption) (*GetEventsResponse, error)
}

type shipmentsSnapshotClient struct {
	cc grpc.ClientConnInterface
}

func NewShipmentsSnapshotClient(cc grpc.ClientConnInterface) ShipmentsSnapshotClient {
	return &shipmentsSnapshotClient{cc}
}

func (c *shipmentsSnapshotClient) GetShipment(ctx context.Context, in *GetShipmentRequest, opts ...grpc.CallOption) (*Shipment, error) {
	out := new(Shipment)
	err := c.cc.Invoke(ctx, "/ShipmentsSnapshot/GetShipment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shipmentsSnapshotClient) GetShipments(ctx context.Context, in *GetShipmentsRequest, opts ...grpc.CallOption) (*GetShipmentsResponse, error) {
	out := new(GetShipmentsResponse)
	err := c.cc.Invoke(ctx, "/ShipmentsSnapshot/GetShipments", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shipmentsSnapshotClient) GetEvents(ctx context.Context, in *GetEventsRequest, opts ...grpc.CallOption) (*GetEventsResponse, error) {
	out := new(GetEventsResponse)
	err := c.cc.Invoke(ctx, "/ShipmentsSnapshot/GetEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShipmentsSnapshotServer is the server API for ShipmentsSnapshot service.
// All implementations must embed UnimplementedShipmentsSnapshotServer
// for forward compatibility
type ShipmentsSnapshotServer interface {
	GetShipment(context.Context, *GetShipmentRequest) (*Shipment, error)
	GetShipments(context.Context, *GetShipmentsRequest) (*GetShipmentsResponse, error)
	GetEvents(context.Context, *GetEventsRequest) (*GetEventsResponse, error)
	mustEmbedUnimplementedShipmentsSnapshotServer()
}

// UnimplementedShipmentsSnapshotServer must be embedded to have forward compatible implementations.
type UnimplementedShipmentsSnapshotServer struct {
}

func (UnimplementedShipmentsSnapshotServer) GetShipment(context.Context, *GetShipmentRequest) (*Shipment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShipment not implemented")
}
func (UnimplementedShipmentsSnapshotServer) GetShipments(context.Context, *GetShipmentsRequest) (*GetShipmentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShipments not implemented")
}
func (UnimplementedShipmentsSnapshotServer) GetEvents(context.Context, *GetEventsRequest) (*GetEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEvents not implemented")
}
func (UnimplementedShipmentsSnapshotServer) mustEmbedUnimplementedShipmentsSnapshotServer() {}

// UnsafeShipmentsSnapshotServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShipmentsSnapshotServer will
// result in compilation errors.
type UnsafeShipmentsSnapshotServer interface {
	mustEmbedUnimplementedShipmentsSnapshotServer()
}

func RegisterShipmentsSnapshotServer(s grpc.ServiceRegistrar, srv ShipmentsSnapshotServer) {
	s.RegisterService(&ShipmentsSnapshot_ServiceDesc, srv)
}

func _ShipmentsSnapshot_GetShipment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetShipmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShipmentsSnapshotServer).GetShipment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ShipmentsSnapshot/GetShipment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShipmentsSnapshotServer).GetShipment(ctx, req.(*GetShipmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShipmentsSnapshot_GetShipments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetShipmentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShipmentsSnapshotServer).GetShipments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ShipmentsSnapshot/GetShipments",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShipmentsSnapshotServer).GetShipments(ctx, req.(*GetShipmentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShipmentsSnapshot_GetEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShipmentsSnapshotServer).GetEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ShipmentsSnapshot/GetEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShipmentsSnapshotServer).GetEvents(ctx, req.(*GetEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ShipmentsSnapshot_ServiceDesc is the grpc.ServiceDesc for ShipmentsSnapshot service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShipmentsSnapshot_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ShipmentsSnapshot",
	HandlerType: (*ShipmentsSnapshotServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetShipment",
			Handler:    _ShipmentsSnapshot_GetShipment_Handler,
		},
		{
			MethodName: "GetShipments",
			Handler:    _ShipmentsSnapshot_GetShipments_Handler,
		},
		{
			MethodName: "GetEvents",
			Handler:    _ShipmentsSnapshot_GetEvents_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/pb/snapshot.proto",
}
