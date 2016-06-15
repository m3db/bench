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

package fs2

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type testEntry struct {
	id        string
	timestamp int64
	value     float64
}

func TestRoundTrip(t *testing.T) {
	indexTempFile, err := ioutil.TempFile("", "indexTempFile-")
	if err != nil {
		t.Fatal(err)
	}
	indexFilePath := indexTempFile.Name()
	defer os.Remove(indexFilePath)

	dataTempFile, err := ioutil.TempFile("", "dataTempFile-")
	if err != nil {
		t.Fatal(err)
	}
	dataFilePath := dataTempFile.Name()
	defer os.Remove(dataFilePath)

	require.NoError(t, indexTempFile.Close())
	require.NoError(t, dataTempFile.Close())
	entries := []testEntry{
		{"foo", 10, 3.6},
		{"foo", 20, 4.8},
		{"bar", 22, 9.1},
		{"bar", 24, 4.2},
	}

	writer, err := NewWriter(indexFilePath, dataFilePath)
	require.NoError(t, err)

	for i := range entries {
		require.NoError(t, writer.Write(entries[i].id, entries[i].timestamp, entries[i].value))
	}
	require.NoError(t, writer.Close())

	reader, err := NewReader(indexFilePath, dataFilePath)
	require.NoError(t, err)
	iter, err := reader.Iter()
	require.NoError(t, err)

	index := 0
	for iter.Next() {
		id, timestamp, value := iter.Value()
		require.Equal(t, entries[index].id, id)
		require.Equal(t, entries[index].timestamp, timestamp)
		require.Equal(t, entries[index].value, value)
		index++
	}
	require.NoError(t, iter.Err())
	require.Equal(t, len(entries), index)

	require.NoError(t, reader.Close())
}
