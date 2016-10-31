// Code generated by protoc-gen-go.
// source: hello.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	hello.proto

It has these top-level messages:
	HelloReq
	HelloResp
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type HelloReq struct {
	Name  string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Times int32  `protobuf:"varint,2,opt,name=times" json:"times,omitempty"`
}

func (m *HelloReq) Reset()                    { *m = HelloReq{} }
func (m *HelloReq) String() string            { return proto1.CompactTextString(m) }
func (*HelloReq) ProtoMessage()               {}
func (*HelloReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type HelloResp struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *HelloResp) Reset()                    { *m = HelloResp{} }
func (m *HelloResp) String() string            { return proto1.CompactTextString(m) }
func (*HelloResp) ProtoMessage()               {}
func (*HelloResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto1.RegisterType((*HelloReq)(nil), "proto.HelloReq")
	proto1.RegisterType((*HelloResp)(nil), "proto.HelloResp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Greeter service

type GreeterClient interface {
	SayHello(ctx context.Context, in *HelloReq, opts ...grpc.CallOption) (*HelloResp, error)
}

type greeterClient struct {
	cc *grpc.ClientConn
}

func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
	return &greeterClient{cc}
}

func (c *greeterClient) SayHello(ctx context.Context, in *HelloReq, opts ...grpc.CallOption) (*HelloResp, error) {
	out := new(HelloResp)
	err := grpc.Invoke(ctx, "/proto.Greeter/SayHello", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Greeter service

type GreeterServer interface {
	SayHello(context.Context, *HelloReq) (*HelloResp, error)
}

func RegisterGreeterServer(s *grpc.Server, srv GreeterServer) {
	s.RegisterService(&_Greeter_serviceDesc, srv)
}

func _Greeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Greeter/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SayHello(ctx, req.(*HelloReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Greeter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto1.RegisterFile("hello.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 147 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0x48, 0xcd, 0xc9,
	0xc9, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0x26, 0x5c, 0x1c, 0x1e,
	0x20, 0xd1, 0xa0, 0xd4, 0x42, 0x21, 0x21, 0x2e, 0x96, 0xbc, 0xc4, 0xdc, 0x54, 0x09, 0x46, 0x05,
	0x46, 0x0d, 0xce, 0x20, 0x30, 0x5b, 0x48, 0x84, 0x8b, 0xb5, 0x24, 0x33, 0x37, 0xb5, 0x58, 0x82,
	0x49, 0x81, 0x51, 0x83, 0x35, 0x08, 0xc2, 0x51, 0x52, 0xe5, 0xe2, 0x84, 0xea, 0x2a, 0x2e, 0x10,
	0x92, 0xe0, 0x62, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x87, 0xe9, 0x84, 0x71, 0x8d, 0xac, 0xb8,
	0xd8, 0xdd, 0x8b, 0x52, 0x53, 0x4b, 0x52, 0x8b, 0x84, 0xf4, 0xb9, 0x38, 0x82, 0x13, 0x2b, 0xc1,
	0x9a, 0x84, 0xf8, 0x21, 0x4e, 0xd0, 0x83, 0x59, 0x2c, 0x25, 0x80, 0x2a, 0x50, 0x5c, 0xa0, 0xc4,
	0x90, 0xc4, 0x06, 0x16, 0x32, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x09, 0xfe, 0xf3, 0xf2, 0xb5,
	0x00, 0x00, 0x00,
}