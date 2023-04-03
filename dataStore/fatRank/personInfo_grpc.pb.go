// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.0
// source: personInfo.proto

package fatRank

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

const (
	RankService_Register_FullMethodName       = "/fatRank.RankService/Register"
	RankService_RegisterPerson_FullMethodName = "/fatRank.RankService/RegisterPerson"
)

// RankServiceClient is the client API for RankService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RankServiceClient interface {
	Register(ctx context.Context, in *PersonalInformation, opts ...grpc.CallOption) (*PersonalInformation, error)
	//单次发送，多次接收
	RegisterPerson(ctx context.Context, opts ...grpc.CallOption) (RankService_RegisterPersonClient, error)
}

type rankServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRankServiceClient(cc grpc.ClientConnInterface) RankServiceClient {
	return &rankServiceClient{cc}
}

func (c *rankServiceClient) Register(ctx context.Context, in *PersonalInformation, opts ...grpc.CallOption) (*PersonalInformation, error) {
	out := new(PersonalInformation)
	err := c.cc.Invoke(ctx, RankService_Register_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rankServiceClient) RegisterPerson(ctx context.Context, opts ...grpc.CallOption) (RankService_RegisterPersonClient, error) {
	stream, err := c.cc.NewStream(ctx, &RankService_ServiceDesc.Streams[0], RankService_RegisterPerson_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &rankServiceRegisterPersonClient{stream}
	return x, nil
}

type RankService_RegisterPersonClient interface {
	Send(*PersonalInformation) error
	CloseAndRecv() (*PersonalInformationList, error)
	grpc.ClientStream
}

type rankServiceRegisterPersonClient struct {
	grpc.ClientStream
}

func (x *rankServiceRegisterPersonClient) Send(m *PersonalInformation) error {
	return x.ClientStream.SendMsg(m)
}

func (x *rankServiceRegisterPersonClient) CloseAndRecv() (*PersonalInformationList, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(PersonalInformationList)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RankServiceServer is the server API for RankService service.
// All implementations must embed UnimplementedRankServiceServer
// for forward compatibility
type RankServiceServer interface {
	Register(context.Context, *PersonalInformation) (*PersonalInformation, error)
	//单次发送，多次接收
	RegisterPerson(RankService_RegisterPersonServer) error
	mustEmbedUnimplementedRankServiceServer()
}

// UnimplementedRankServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRankServiceServer struct {
}

func (UnimplementedRankServiceServer) Register(context.Context, *PersonalInformation) (*PersonalInformation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedRankServiceServer) RegisterPerson(RankService_RegisterPersonServer) error {
	return status.Errorf(codes.Unimplemented, "method RegisterPerson not implemented")
}
func (UnimplementedRankServiceServer) mustEmbedUnimplementedRankServiceServer() {}

// UnsafeRankServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RankServiceServer will
// result in compilation errors.
type UnsafeRankServiceServer interface {
	mustEmbedUnimplementedRankServiceServer()
}

func RegisterRankServiceServer(s grpc.ServiceRegistrar, srv RankServiceServer) {
	s.RegisterService(&RankService_ServiceDesc, srv)
}

func _RankService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PersonalInformation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RankServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RankService_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RankServiceServer).Register(ctx, req.(*PersonalInformation))
	}
	return interceptor(ctx, in, info, handler)
}

func _RankService_RegisterPerson_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RankServiceServer).RegisterPerson(&rankServiceRegisterPersonServer{stream})
}

type RankService_RegisterPersonServer interface {
	SendAndClose(*PersonalInformationList) error
	Recv() (*PersonalInformation, error)
	grpc.ServerStream
}

type rankServiceRegisterPersonServer struct {
	grpc.ServerStream
}

func (x *rankServiceRegisterPersonServer) SendAndClose(m *PersonalInformationList) error {
	return x.ServerStream.SendMsg(m)
}

func (x *rankServiceRegisterPersonServer) Recv() (*PersonalInformation, error) {
	m := new(PersonalInformation)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RankService_ServiceDesc is the grpc.ServiceDesc for RankService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RankService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fatRank.RankService",
	HandlerType: (*RankServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _RankService_Register_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RegisterPerson",
			Handler:       _RankService_RegisterPerson_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "personInfo.proto",
}
