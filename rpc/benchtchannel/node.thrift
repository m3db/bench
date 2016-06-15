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

namespace java com.memtsdb.node.benchmarks

service MemTSDB {
	FetchResult fetch(1: FetchRequest req)
	FetchBatchResult fetchBatch(1: FetchBatchRequest req)
}

struct FetchRequest {
	1: required i64 startUnixMs
	2: required i64 endUnixMs
	3: required string id
}

struct Segment {
	1: required binary value
}

struct FetchResult {
	1: required list<Segment> segments
}

struct FetchBatchRequest {
	1: required i64 startUnixMs
	2: required i64 endUnixMs
	3: required list<string> ids
}

struct FetchBatchResult {
	1: required list<FetchResult> results
}
