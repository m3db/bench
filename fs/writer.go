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
	"os"

	"github.com/golang/protobuf/proto"
	pb "github.com/m3db/bench/fs/proto"
)

// Writer writes protobuf encoded data entries to a file.
type Writer struct {
	fd      *os.File
	sizeBuf *proto.Buffer
	dataBuf *proto.Buffer
}

// NewWriter creates a writer.
func NewWriter(filePath string) (*Writer, error) {
	fd, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return nil, err
	}
	return &Writer{fd: fd, sizeBuf: proto.NewBuffer(nil), dataBuf: proto.NewBuffer(nil)}, nil
}

// Write writes a data entry.
func (w *Writer) Write(entry *pb.DataEntry) error {
	w.sizeBuf.Reset()
	w.dataBuf.Reset()
	if err := w.dataBuf.Marshal(entry); err != nil {
		return err
	}
	entryBytes := w.dataBuf.Bytes()
	if err := w.sizeBuf.EncodeVarint(uint64(len(entryBytes))); err != nil {
		return err
	}
	if _, err := w.fd.Write(w.sizeBuf.Bytes()); err != nil {
		return err
	}
	if _, err := w.fd.Write(entryBytes); err != nil {
		return err
	}
	return nil
}

func (w *Writer) writeEOFMarker() error {
	if _, err := w.fd.Write(proto.EncodeVarint(eofMarker)); err != nil {
		return err
	}
	return nil
}

// Close closes the file with the eof marker appended.
func (w *Writer) Close() error {
	if err := w.writeEOFMarker(); err != nil {
		return err
	}
	return w.fd.Close()
}
