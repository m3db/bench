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

package fs

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	pb "github.com/m3db/bench/fs/proto"
)

// Reader reads protobuf encoded data entries from a file.
type Reader struct {
	buf []byte // raw bytes stored in the file
}

// NewReader creates a reader.
func NewReader(filePath string) (*Reader, error) {
	input, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open %s: %v", filePath, err)
	}

	defer input.Close()

	b, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("unable to read all bytes from %s: %v", filePath, err)
	}

	return &Reader{buf: b}, nil
}

// Iter creates an iterator.
func (r *Reader) Iter() *Iterator {
	return newIterator(r.buf)
}

// Iterator iterates over protobuf encoded data entries.
type Iterator struct {
	buf   []byte
	entry *pb.DataEntry
	done  bool
	err   error
}

func newIterator(buf []byte) *Iterator {
	return &Iterator{buf: buf, entry: &pb.DataEntry{}}
}

// Next moves to the next data entry.
func (it *Iterator) Next() bool {
	if it.done || it.err != nil {
		return false
	}
	it.entry.Reset()
	entrySize, consumed := proto.DecodeVarint(it.buf)
	if consumed == 0 {
		it.err = errReadDataEntryNoSize
		return false
	}
	it.buf = it.buf[consumed:]
	if entrySize == eofMarker {
		it.done = true
		return false
	}
	data := it.buf[:entrySize]
	if it.err = proto.Unmarshal(data, it.entry); it.err != nil {
		return false
	}
	it.buf = it.buf[entrySize:]
	return true
}

// Value returns the current entry.
func (it *Iterator) Value() *pb.DataEntry {
	return it.entry
}

// Err returns the current error.
func (it *Iterator) Err() error {
	return it.err
}
