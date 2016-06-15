// Code generated by protoc-gen-go.
// source: node.proto
// DO NOT EDIT!

/*
Package benchgrpc is a generated protocol buffer package.

It is generated from these files:
	node.proto

It has these top-level messages:
	FetchRequest
	FetchResult
	FetchBatchRequest
	FetchBatchResult
*/
// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package benchgrpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type FetchRequest struct {
	StartUnixMs int64  `protobuf:"varint,1,opt,name=startUnixMs" json:"startUnixMs,omitempty"`
	EndUnixMs   int64  `protobuf:"varint,2,opt,name=endUnixMs" json:"endUnixMs,omitempty"`
	Id          string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
}

func (m *FetchRequest) Reset()                    { *m = FetchRequest{} }
func (m *FetchRequest) String() string            { return proto.CompactTextString(m) }
func (*FetchRequest) ProtoMessage()               {}
func (*FetchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type FetchResult struct {
	Segments []*FetchResult_Segment `protobuf:"bytes,1,rep,name=segments" json:"segments,omitempty"`
}

func (m *FetchResult) Reset()                    { *m = FetchResult{} }
func (m *FetchResult) String() string            { return proto.CompactTextString(m) }
func (*FetchResult) ProtoMessage()               {}
func (*FetchResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *FetchResult) GetSegments() []*FetchResult_Segment {
	if m != nil {
		return m.Segments
	}
	return nil
}

type FetchResult_Segment struct {
	Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *FetchResult_Segment) Reset()                    { *m = FetchResult_Segment{} }
func (m *FetchResult_Segment) String() string            { return proto.CompactTextString(m) }
func (*FetchResult_Segment) ProtoMessage()               {}
func (*FetchResult_Segment) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

type FetchBatchRequest struct {
	StartUnixMs int64    `protobuf:"varint,1,opt,name=startUnixMs" json:"startUnixMs,omitempty"`
	EndUnixMs   int64    `protobuf:"varint,2,opt,name=endUnixMs" json:"endUnixMs,omitempty"`
	Ids         []string `protobuf:"bytes,3,rep,name=ids" json:"ids,omitempty"`
}

func (m *FetchBatchRequest) Reset()                    { *m = FetchBatchRequest{} }
func (m *FetchBatchRequest) String() string            { return proto.CompactTextString(m) }
func (*FetchBatchRequest) ProtoMessage()               {}
func (*FetchBatchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type FetchBatchResult struct {
	Results []*FetchResult `protobuf:"bytes,1,rep,name=results" json:"results,omitempty"`
}

func (m *FetchBatchResult) Reset()                    { *m = FetchBatchResult{} }
func (m *FetchBatchResult) String() string            { return proto.CompactTextString(m) }
func (*FetchBatchResult) ProtoMessage()               {}
func (*FetchBatchResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *FetchBatchResult) GetResults() []*FetchResult {
	if m != nil {
		return m.Results
	}
	return nil
}

func init() {
	proto.RegisterType((*FetchRequest)(nil), "benchgrpc.FetchRequest")
	proto.RegisterType((*FetchResult)(nil), "benchgrpc.FetchResult")
	proto.RegisterType((*FetchResult_Segment)(nil), "benchgrpc.FetchResult.Segment")
	proto.RegisterType((*FetchBatchRequest)(nil), "benchgrpc.FetchBatchRequest")
	proto.RegisterType((*FetchBatchResult)(nil), "benchgrpc.FetchBatchResult")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion2

// Client API for MemTSDB service

type MemTSDBClient interface {
	Fetch(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (*FetchResult, error)
	FetchStream(ctx context.Context, opts ...grpc.CallOption) (MemTSDB_FetchStreamClient, error)
	FetchBatch(ctx context.Context, in *FetchBatchRequest, opts ...grpc.CallOption) (*FetchBatchResult, error)
	FetchBatchStream(ctx context.Context, opts ...grpc.CallOption) (MemTSDB_FetchBatchStreamClient, error)
}

type memTSDBClient struct {
	cc *grpc.ClientConn
}

func NewMemTSDBClient(cc *grpc.ClientConn) MemTSDBClient {
	return &memTSDBClient{cc}
}

func (c *memTSDBClient) Fetch(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (*FetchResult, error) {
	out := new(FetchResult)
	err := grpc.Invoke(ctx, "/benchgrpc.MemTSDB/Fetch", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *memTSDBClient) FetchStream(ctx context.Context, opts ...grpc.CallOption) (MemTSDB_FetchStreamClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_MemTSDB_serviceDesc.Streams[0], c.cc, "/benchgrpc.MemTSDB/FetchStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &memTSDBFetchStreamClient{stream}
	return x, nil
}

type MemTSDB_FetchStreamClient interface {
	Send(*FetchRequest) error
	Recv() (*FetchResult, error)
	grpc.ClientStream
}

type memTSDBFetchStreamClient struct {
	grpc.ClientStream
}

func (x *memTSDBFetchStreamClient) Send(m *FetchRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *memTSDBFetchStreamClient) Recv() (*FetchResult, error) {
	m := new(FetchResult)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *memTSDBClient) FetchBatch(ctx context.Context, in *FetchBatchRequest, opts ...grpc.CallOption) (*FetchBatchResult, error) {
	out := new(FetchBatchResult)
	err := grpc.Invoke(ctx, "/benchgrpc.MemTSDB/FetchBatch", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *memTSDBClient) FetchBatchStream(ctx context.Context, opts ...grpc.CallOption) (MemTSDB_FetchBatchStreamClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_MemTSDB_serviceDesc.Streams[1], c.cc, "/benchgrpc.MemTSDB/FetchBatchStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &memTSDBFetchBatchStreamClient{stream}
	return x, nil
}

type MemTSDB_FetchBatchStreamClient interface {
	Send(*FetchBatchRequest) error
	Recv() (*FetchBatchResult, error)
	grpc.ClientStream
}

type memTSDBFetchBatchStreamClient struct {
	grpc.ClientStream
}

func (x *memTSDBFetchBatchStreamClient) Send(m *FetchBatchRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *memTSDBFetchBatchStreamClient) Recv() (*FetchBatchResult, error) {
	m := new(FetchBatchResult)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for MemTSDB service

type MemTSDBServer interface {
	Fetch(context.Context, *FetchRequest) (*FetchResult, error)
	FetchStream(MemTSDB_FetchStreamServer) error
	FetchBatch(context.Context, *FetchBatchRequest) (*FetchBatchResult, error)
	FetchBatchStream(MemTSDB_FetchBatchStreamServer) error
}

func RegisterMemTSDBServer(s *grpc.Server, srv MemTSDBServer) {
	s.RegisterService(&_MemTSDB_serviceDesc, srv)
}

func _MemTSDB_Fetch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemTSDBServer).Fetch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/benchgrpc.MemTSDB/Fetch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemTSDBServer).Fetch(ctx, req.(*FetchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MemTSDB_FetchStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MemTSDBServer).FetchStream(&memTSDBFetchStreamServer{stream})
}

type MemTSDB_FetchStreamServer interface {
	Send(*FetchResult) error
	Recv() (*FetchRequest, error)
	grpc.ServerStream
}

type memTSDBFetchStreamServer struct {
	grpc.ServerStream
}

func (x *memTSDBFetchStreamServer) Send(m *FetchResult) error {
	return x.ServerStream.SendMsg(m)
}

func (x *memTSDBFetchStreamServer) Recv() (*FetchRequest, error) {
	m := new(FetchRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _MemTSDB_FetchBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemTSDBServer).FetchBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/benchgrpc.MemTSDB/FetchBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemTSDBServer).FetchBatch(ctx, req.(*FetchBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MemTSDB_FetchBatchStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MemTSDBServer).FetchBatchStream(&memTSDBFetchBatchStreamServer{stream})
}

type MemTSDB_FetchBatchStreamServer interface {
	Send(*FetchBatchResult) error
	Recv() (*FetchBatchRequest, error)
	grpc.ServerStream
}

type memTSDBFetchBatchStreamServer struct {
	grpc.ServerStream
}

func (x *memTSDBFetchBatchStreamServer) Send(m *FetchBatchResult) error {
	return x.ServerStream.SendMsg(m)
}

func (x *memTSDBFetchBatchStreamServer) Recv() (*FetchBatchRequest, error) {
	m := new(FetchBatchRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _MemTSDB_serviceDesc = grpc.ServiceDesc{
	ServiceName: "benchgrpc.MemTSDB",
	HandlerType: (*MemTSDBServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Fetch",
			Handler:    _MemTSDB_Fetch_Handler,
		},
		{
			MethodName: "FetchBatch",
			Handler:    _MemTSDB_FetchBatch_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "FetchStream",
			Handler:       _MemTSDB_FetchStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "FetchBatchStream",
			Handler:       _MemTSDB_FetchBatchStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
}

var fileDescriptor0 = []byte{
	// 299 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x92, 0x41, 0x4f, 0x83, 0x40,
	0x10, 0x85, 0x03, 0xa4, 0x22, 0x43, 0x63, 0xea, 0xc4, 0x68, 0x83, 0x8d, 0x36, 0x9c, 0x38, 0x91,
	0xa6, 0x5e, 0x8c, 0x47, 0xd2, 0xe8, 0xa9, 0x07, 0x17, 0xbd, 0x9a, 0x50, 0x98, 0xb4, 0x98, 0x02,
	0x95, 0x5d, 0x8c, 0xff, 0xc7, 0x3f, 0x2a, 0x6c, 0x69, 0xbb, 0x51, 0x7b, 0xd0, 0xf4, 0x36, 0xcc,
	0xbc, 0x79, 0xbc, 0x6f, 0xb2, 0x00, 0x79, 0x91, 0x90, 0xbf, 0x2a, 0x0b, 0x51, 0xa0, 0x35, 0xa3,
	0x3c, 0x5e, 0xcc, 0xcb, 0x55, 0xec, 0xbe, 0x40, 0xf7, 0x9e, 0x44, 0xbc, 0x60, 0xf4, 0x56, 0x11,
	0x17, 0x38, 0x04, 0x9b, 0x8b, 0xa8, 0x14, 0xcf, 0x79, 0xfa, 0x31, 0xe5, 0x7d, 0x6d, 0xa8, 0x79,
	0x06, 0x53, 0x5b, 0x38, 0x00, 0x8b, 0xf2, 0xa4, 0x9d, 0xeb, 0x72, 0xbe, 0x6b, 0xe0, 0x09, 0xe8,
	0x69, 0xd2, 0x37, 0xea, 0xb6, 0xc5, 0xea, 0xca, 0x7d, 0x05, 0xbb, 0xf5, 0xe7, 0xd5, 0x52, 0xe0,
	0x1d, 0x1c, 0x73, 0x9a, 0x67, 0x94, 0x8b, 0xc6, 0xdb, 0xf0, 0xec, 0xf1, 0x95, 0xbf, 0x0d, 0xe3,
	0x2b, 0x4a, 0x3f, 0x5c, 0xcb, 0xd8, 0x56, 0xef, 0x5c, 0x83, 0xd9, 0x36, 0xf1, 0x0c, 0x3a, 0xef,
	0xd1, 0xb2, 0x22, 0x99, 0xaf, 0xcb, 0xd6, 0x1f, 0x2e, 0xc1, 0xa9, 0x74, 0x08, 0xa2, 0x03, 0x02,
	0xf5, 0xc0, 0x48, 0x13, 0x5e, 0x13, 0x19, 0x35, 0x51, 0x53, 0xba, 0x13, 0xe8, 0xa9, 0xbf, 0x91,
	0x5c, 0x23, 0x30, 0x4b, 0x59, 0x6d, 0xb0, 0xce, 0x7f, 0xc7, 0x62, 0x1b, 0xd9, 0xf8, 0x53, 0x07,
	0x73, 0x4a, 0xd9, 0x53, 0x38, 0x09, 0xf0, 0x16, 0x3a, 0x52, 0x83, 0x17, 0x3f, 0xb7, 0x24, 0x85,
	0xb3, 0xc7, 0x0e, 0x83, 0xf6, 0xbc, 0xa1, 0x28, 0x29, 0xca, 0xfe, 0xbc, 0xef, 0x69, 0x23, 0x0d,
	0x1f, 0x00, 0x76, 0x3c, 0x38, 0xf8, 0xae, 0x54, 0xaf, 0xe9, 0x5c, 0xee, 0x99, 0xca, 0x30, 0x8f,
	0xea, 0x61, 0xda, 0x44, 0xff, 0xb7, 0x6b, 0xb2, 0xcd, 0x8e, 0xe4, 0x83, 0xbd, 0xf9, 0x0a, 0x00,
	0x00, 0xff, 0xff, 0x54, 0x51, 0x38, 0xc3, 0xbe, 0x02, 0x00, 0x00,
}
