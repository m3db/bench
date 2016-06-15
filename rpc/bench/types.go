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

package bench

// RequestGenerator generates request descriptions on demand
type RequestGenerator func() *RequestDescription

// RequestDescription describes a request
type RequestDescription struct {
	StartUnixMs int64
	EndUnixMs   int64
	IDs         []string
}

// WorkerPool is a bounded pool of work
type WorkerPool struct {
	ch chan struct{}
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(size int) *WorkerPool {
	pool := &WorkerPool{
		ch: make(chan struct{}, size),
	}
	for i := 0; i < size; i++ {
		pool.ch <- struct{}{}
	}
	return pool
}

// Go will launch work and on reaching bounds will block until running
func (p *WorkerPool) Go(f func()) {
	s := <-p.ch
	go func() {
		f()
		p.ch <- s
	}()
}
