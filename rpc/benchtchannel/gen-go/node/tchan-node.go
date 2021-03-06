// @generated Code generated by thrift-gen. Do not modify.

// Package node is generated code used to make or handle TChannel calls using Thrift.
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

package node

import (
	"fmt"

	athrift "github.com/apache/thrift/lib/go/thrift"
	"github.com/uber/tchannel-go/thrift"
)

// Interfaces for the service and client for the services defined in the IDL.

// TChanMemTSDB is the interface that defines the server handler and client interface.
type TChanMemTSDB interface {
	Fetch(ctx thrift.Context, req *FetchRequest) (*FetchResult_, error)
	FetchBatch(ctx thrift.Context, req *FetchBatchRequest) (*FetchBatchResult_, error)
}

// Implementation of a client and service handler.

type tchanMemTSDBClient struct {
	thriftService string
	client        thrift.TChanClient
}

func NewTChanMemTSDBInheritedClient(thriftService string, client thrift.TChanClient) *tchanMemTSDBClient {
	return &tchanMemTSDBClient{
		thriftService,
		client,
	}
}

// NewTChanMemTSDBClient creates a client that can be used to make remote calls.
func NewTChanMemTSDBClient(client thrift.TChanClient) TChanMemTSDB {
	return NewTChanMemTSDBInheritedClient("MemTSDB", client)
}

func (c *tchanMemTSDBClient) Fetch(ctx thrift.Context, req *FetchRequest) (*FetchResult_, error) {
	var resp MemTSDBFetchResult
	args := MemTSDBFetchArgs{
		Req: req,
	}
	success, err := c.client.Call(ctx, c.thriftService, "fetch", &args, &resp)
	if err == nil && !success {
	}

	return resp.GetSuccess(), err
}

func (c *tchanMemTSDBClient) FetchBatch(ctx thrift.Context, req *FetchBatchRequest) (*FetchBatchResult_, error) {
	var resp MemTSDBFetchBatchResult
	args := MemTSDBFetchBatchArgs{
		Req: req,
	}
	success, err := c.client.Call(ctx, c.thriftService, "fetchBatch", &args, &resp)
	if err == nil && !success {
	}

	return resp.GetSuccess(), err
}

type tchanMemTSDBServer struct {
	handler TChanMemTSDB
}

// NewTChanMemTSDBServer wraps a handler for TChanMemTSDB so it can be
// registered with a thrift.Server.
func NewTChanMemTSDBServer(handler TChanMemTSDB) thrift.TChanServer {
	return &tchanMemTSDBServer{
		handler,
	}
}

func (s *tchanMemTSDBServer) Service() string {
	return "MemTSDB"
}

func (s *tchanMemTSDBServer) Methods() []string {
	return []string{
		"fetch",
		"fetchBatch",
	}
}

func (s *tchanMemTSDBServer) Handle(ctx thrift.Context, methodName string, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	switch methodName {
	case "fetch":
		return s.handleFetch(ctx, protocol)
	case "fetchBatch":
		return s.handleFetchBatch(ctx, protocol)

	default:
		return false, nil, fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
}

func (s *tchanMemTSDBServer) handleFetch(ctx thrift.Context, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	var req MemTSDBFetchArgs
	var res MemTSDBFetchResult

	if err := req.Read(protocol); err != nil {
		return false, nil, err
	}

	r, err :=
		s.handler.Fetch(ctx, req.Req)

	if err != nil {
		return false, nil, err
	} else {
		res.Success = r
	}

	return err == nil, &res, nil
}

func (s *tchanMemTSDBServer) handleFetchBatch(ctx thrift.Context, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	var req MemTSDBFetchBatchArgs
	var res MemTSDBFetchBatchResult

	if err := req.Read(protocol); err != nil {
		return false, nil, err
	}

	r, err :=
		s.handler.FetchBatch(ctx, req.Req)

	if err != nil {
		return false, nil, err
	} else {
		res.Success = r
	}

	return err == nil, &res, nil
}
