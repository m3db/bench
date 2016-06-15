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

package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/m3db/bench/rpc/benchgrpc"
	"github.com/m3db/bench/rpc/benchtchannel"
	"github.com/m3db/bench/rpc/benchtchannelgogoprotobuf"
	"github.com/m3db/m3db/x/logging"

	"google.golang.org/grpc"
)

var log = logging.SimpleLogger

type keyValue struct {
	key   string
	value string
}

var (
	grpcConcurrencyArg    = flag.Int("grpcConcurrency", 1000000 /* C1M */, "number of concurrent requests")
	grpcAddrArg           = flag.String("grpcaddr", "0.0.0.0:8888", "benchmark GRPC server address")
	tchannelAddrArg       = flag.String("tchanneladdr", "0.0.0.0:8889", "benchmark TChannel server address")
	tchannelGogopbAddrArg = flag.String("tchannelgogopbaddr", "0.0.0.0:8890", "benchmark TChannel gogoprotobuf server address")
)

func main() {
	flag.Parse()

	if *grpcConcurrencyArg <= 0 ||
		*grpcAddrArg == "" ||
		*tchannelAddrArg == "" ||
		*tchannelGogopbAddrArg == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	grpcConcurrency := uint32(*grpcConcurrencyArg)
	grpcAddr := *grpcAddrArg
	tchannelAddr := *tchannelAddrArg
	tchannelGogopbAddr := *tchannelGogopbAddrArg

	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("could not create GRPC TCP Listener: %v", err)
	}

	defer grpcListener.Close()

	grpcClose, err := benchgrpc.StartTestGRPCServer(grpcListener, grpc.MaxConcurrentStreams(grpcConcurrency))
	if err != nil {
		log.Fatalf("could not create GRPC server: %v", err)
	}

	defer grpcClose()

	grpcParams := []keyValue{
		{"grpcConcurrency", fmt.Sprintf("%d", grpcConcurrency)},
	}

	description := ""
	for _, kv := range grpcParams {
		description += fmt.Sprintf("\n%s=%s", kv.key, kv.value)
	}
	log.Infof("started GRPC server: %v %s", grpcListener.Addr().String(), description)

	tchannelClose, err := benchtchannel.StartTestTChannelServer(tchannelAddr, nil)
	if err != nil {
		log.Fatalf("could not create TChannel server: %v", err)
	}

	defer tchannelClose()

	log.Infof("started TChannel server: %v", tchannelAddr)

	tchannelGogopbClose, err := benchtchannelgogoprotobuf.StartTestTChannelServer(tchannelGogopbAddr, nil)
	if err != nil {
		log.Fatalf("could not create TChannel server: %v", err)
	}

	defer tchannelGogopbClose()

	log.Infof("started TChannelGogopb server: %v", tchannelGogopbAddr)

	debugAddr := "0.0.0.0:8887"
	log.Infof("started debug server: %v", debugAddr)
	go func() {
		log.Infof("debug server error: %v", http.ListenAndServe(debugAddr, nil))
	}()

	log.Fatalf("fatal: %v", interrupt())
}

func interrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("%s", <-c)
}
