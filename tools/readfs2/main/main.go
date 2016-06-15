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
	"os"

	"github.com/m3db/bench/fs2"
	"github.com/m3db/m3db/x/logging"
)

var log = logging.SimpleLogger

var (
	indexFile = flag.String("indexFile", "", "input index file")
	dataFile  = flag.String("dataFile", "", "input data file")
)

func main() {
	flag.Parse()

	if len(*indexFile) == 0 || len(*dataFile) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	log.Infof("creating input reader")
	reader, err := fs2.NewReader(*indexFile, *dataFile)
	if err != nil {
		log.Fatalf("unable to create a new input reader: %v", err)
	}

	defer reader.Close()

	iter, err := reader.Iter()
	if err != nil {
		log.Fatalf("unable to create iterator: %v", err)
	}

	log.Infof("reading input")
	for iter.Next() {
		id, ts, value := iter.Value()
		log.Infof("timestamp: %16d, value: %.10f, id: %s", ts, value, id)
	}
	if err := iter.Err(); err != nil {
		log.Fatalf("error reading: %v", err)
	}
	log.Infof("done")
}
