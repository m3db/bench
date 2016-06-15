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

package benchtchannel

import (
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/m3db/bench/rpc/bench"
	"github.com/m3db/bench/rpc/benchtchannel/gen-go/node"
	"github.com/m3db/m3db/x/logging"

	"github.com/uber/tchannel-go"
	"github.com/uber/tchannel-go/thrift"
)

var log = logging.SimpleLogger

// Close is a method to call to close a resource or procedure
type Close func()

type service struct {
	segment      []byte
	bytesPerHour int
}

func newDefaultService() *service {
	// [2048]byte fixed data segment
	// 900 bytesPerHour (2.5 bytes per datapoint, 6 data points per minute, 60 minutes)
	data := make([]byte, 2048)
	for i := range data {
		data[i] = byte(i % 255)
	}
	return newService(data, 900)
}

func newService(segment []byte, bytesPerHour int) *service {
	return &service{segment, bytesPerHour}
}

func (s *service) Fetch(ctx thrift.Context, req *node.FetchRequest) (*node.FetchResult_, error) {
	durationMs := math.Abs(float64(req.EndUnixMs) - float64(req.StartUnixMs))
	hours := durationMs / 1000 / 60 / 60

	bytesLen := int(math.Ceil(hours * float64(s.bytesPerHour)))
	segmentsLen := int(math.Ceil(float64(bytesLen) / float64(len(s.segment))))
	lastSegmentLen := bytesLen % len(s.segment)

	segments := make([]*node.Segment, segmentsLen)
	for i := range segments {
		if i == segmentsLen-1 {
			segments[i] = &node.Segment{Value: s.segment[:lastSegmentLen]}
		} else {
			segments[i] = &node.Segment{Value: s.segment}
		}
	}

	return &node.FetchResult_{Segments: segments}, nil
}

func (s *service) FetchBatch(ctx thrift.Context, req *node.FetchBatchRequest) (*node.FetchBatchResult_, error) {
	durationMs := math.Abs(float64(req.EndUnixMs) - float64(req.StartUnixMs))
	hours := durationMs / 1000 / 60 / 60

	bytesLen := int(math.Ceil(hours * float64(s.bytesPerHour)))
	segmentsLen := int(math.Ceil(float64(bytesLen) / float64(len(s.segment))))
	lastSegmentLen := bytesLen % len(s.segment)

	segments := make([]*node.Segment, segmentsLen)
	for i := range segments {
		if i == segmentsLen-1 {
			segments[i] = &node.Segment{Value: s.segment[:lastSegmentLen]}
		} else {
			segments[i] = &node.Segment{Value: s.segment}
		}
	}

	result := &node.FetchBatchResult_{Results: make([]*node.FetchResult_, len(req.Ids))}
	singleResult := &node.FetchResult_{Segments: segments}
	for i := range result.Results {
		result.Results[i] = singleResult
	}

	return result, nil
}

func serveTestTChannelServer(address string, impl node.TChanMemTSDB, opts *tchannel.ChannelOptions) (Close, error) {
	channel, err := tchannel.NewChannel("benchtchannel-server", opts)
	if err != nil {
		return nil, err
	}

	server := thrift.NewServer(channel)
	server.Register(node.NewTChanMemTSDBServer(impl))

	channel.ListenAndServe(address)

	return func() {
		channel.Close()
	}, nil
}

// StartTestTChannelServer starts a test server
func StartTestTChannelServer(address string, opts *tchannel.ChannelOptions) (Close, error) {
	return serveTestTChannelServer(address, newDefaultService(), opts)
}

type clientLBEntry struct {
	channel *tchannel.Channel
	client  node.TChanMemTSDB
}

type clientRoundRobinLB struct {
	pool    []*clientLBEntry
	poolLen int64
	used    int64
}

func newClientRoundRobinLB(address string, n int) *clientRoundRobinLB {
	pool := make([]*clientLBEntry, n)
	for i := 0; i < n; i++ {
		channel, err := tchannel.NewChannel("benchtchannel-client", nil)
		if err != nil {
			log.Fatalf("TChannel construct channel to %s error: %v", address, err)
		}

		endpoint := &thrift.ClientOptions{HostPort: address}
		thriftClient := thrift.NewClient(channel, "benchtchannel-server", endpoint)
		client := node.NewTChanMemTSDBClient(thriftClient)

		// Force a request to open the connection
		tctx, _ := thrift.NewContext(time.Second)
		_, err = client.Fetch(tctx, &node.FetchRequest{StartUnixMs: 0, EndUnixMs: 0, ID: "health"})
		if err != nil {
			log.Fatalf("could not connect to tchannel benchmark server: %v", err)
		}

		pool[i] = &clientLBEntry{channel, client}
	}
	return &clientRoundRobinLB{pool, int64(len(pool)), 0}
}

func (p *clientRoundRobinLB) get() *clientLBEntry {
	n := atomic.AddInt64(&p.used, 1)
	return p.pool[n%p.poolLen]
}

func (p *clientRoundRobinLB) closeAll() {
	for i := range p.pool {
		p.pool[i].channel.Close()
	}
}

// BenchmarkTChannelFetch benchmarks fetching one by one
func BenchmarkTChannelFetch(
	address string,
	n, connectionsPerHost, concurrency int,
	nextRequest bench.RequestGenerator,
	ready chan<- struct{},
	start <-chan struct{},
	done chan<- []error,
) {
	lb := newClientRoundRobinLB(address, connectionsPerHost)
	defer lb.closeAll()

	workers := bench.NewWorkerPool(concurrency)

	fetchRequestPool := sync.Pool{
		New: func() interface{} {
			return &node.FetchRequest{}
		},
	}

	var errLock sync.Mutex
	var errs []error

	ready <- struct{}{}
	<-start

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		desc := nextRequest()

		for j := range desc.IDs {
			req := fetchRequestPool.Get().(*node.FetchRequest)
			req.StartUnixMs = desc.StartUnixMs
			req.EndUnixMs = desc.EndUnixMs
			req.ID = desc.IDs[j]

			wg.Add(1)
			workers.Go(func() {
				tctx, _ := thrift.NewContext(24 * time.Hour)
				_, err := lb.get().client.Fetch(tctx, req)
				if err != nil {
					errLock.Lock()
					errs = append(errs, err)
					errLock.Unlock()
				}
				fetchRequestPool.Put(req)

				wg.Done()
			})
		}
	}

	wg.Wait()

	done <- errs
}

// BenchmarkTChannelFetchBatch benchmarks fetching one by one with a stream
func BenchmarkTChannelFetchBatch(
	address string,
	n, connectionsPerHost, concurrency, batchLen int,
	nextRequest bench.RequestGenerator,
	ready chan<- struct{},
	start <-chan struct{},
	done chan<- []error,
) {
	lb := newClientRoundRobinLB(address, connectionsPerHost)
	defer lb.closeAll()

	workers := bench.NewWorkerPool(concurrency)

	fetchRequestPool := sync.Pool{
		New: func() interface{} {
			return &node.FetchBatchRequest{}
		},
	}

	var errLock sync.Mutex
	var errs []error

	ready <- struct{}{}
	<-start

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		desc := nextRequest()

		wg.Add(1)
		workers.Go(func() {
			var batchesWg sync.WaitGroup

			idsLen := len(desc.IDs)

			for j := 0; j < idsLen; j += batchLen {
				desc := desc
				j := j

				batchesWg.Add(1)
				workers.Go(func() {
					req := fetchRequestPool.Get().(*node.FetchBatchRequest)
					req.StartUnixMs = desc.StartUnixMs
					req.EndUnixMs = desc.EndUnixMs

					k := minInt(j+batchLen, idsLen)
					req.Ids = desc.IDs[j:k]

					tctx, _ := thrift.NewContext(24 * time.Hour)
					_, err := lb.get().client.FetchBatch(tctx, req)
					if err != nil {
						errLock.Lock()
						errs = append(errs, err)
						errLock.Unlock()
					}
					fetchRequestPool.Put(req)
					batchesWg.Done()
				})
			}

			batchesWg.Wait()
			wg.Done()
		})
	}

	wg.Wait()

	done <- errs
}

func minInt(lhs, rhs int) int {
	if lhs < rhs {
		return lhs
	}
	return rhs
}
