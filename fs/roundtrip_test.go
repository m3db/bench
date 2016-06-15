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
	"io/ioutil"
	"os"
	"testing"

	pb "github.com/m3db/bench/fs/proto"
	"github.com/stretchr/testify/require"
)

func TestRoundTrip(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "tempFile-")
	if err != nil {
		t.Fatal(err)
	}
	filePath := tempFile.Name()
	defer os.Remove(filePath)

	require.NoError(t, tempFile.Close())
	entries := []pb.DataEntry{
		{
			Id: "foo",
			Values: []*pb.DataEntry_Datapoint{
				{Timestamp: 10, Value: 3.6},
				{Timestamp: 20, Value: 4.8},
			},
		},
		{
			Id: "bar",
			Values: []*pb.DataEntry_Datapoint{
				{Timestamp: 100, Value: 9.1},
				{Timestamp: 200, Value: 4.2},
			},
		},
	}

	writer, err := NewWriter(filePath)
	require.NoError(t, err)

	for i := range entries {
		require.NoError(t, writer.Write(&entries[i]))
	}
	require.NoError(t, writer.Close())

	reader, err := NewReader(filePath)
	require.NoError(t, err)
	iter := reader.Iter()
	for index := 0; iter.Next(); index++ {
		v := iter.Value()
		require.Equal(t, entries[index].Id, v.Id)
		require.Equal(t, entries[index].Values, v.Values)
	}
	require.NoError(t, iter.Err())
}
