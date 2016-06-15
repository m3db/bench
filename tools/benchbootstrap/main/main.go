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
	"strconv"
	"strings"
	"time"

	"github.com/m3db/m3db"
	"github.com/m3db/m3db/bootstrap"
	"github.com/m3db/m3db/bootstrap/bootstrapper/fs"
	"github.com/m3db/m3db/storage"
)

var (
	bootstrapShards  = flag.String("bootstrapShards", "", "shards that need to be bootstrapped, separated by comma")
	bootstrapperList = flag.String("bootstrapperList", "filesystem", "data sources used for bootstrapping")
	filePathPrefix   = flag.String("filePathPrefix", "", "file path prefix for the raw data files stored on disk")
)

func main() {
	flag.Parse()

	if len(*bootstrapShards) == 0 || len(*bootstrapperList) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	dbOptions := storage.NewDatabaseOptions()

	log := dbOptions.GetLogger()

	shardStrings := strings.Split(*bootstrapShards, ",")
	shards := make([]uint32, len(shardStrings))
	for i, ss := range shardStrings {
		sn, err := strconv.Atoi(ss)
		if err != nil {
			log.Fatalf("illegal shard string %d: %v", sn, err)
		}
		shards[i] = uint32(sn)
	}

	bootstrapperNames := strings.Split(*bootstrapperList, ",")
	var bs memtsdb.Bootstrapper
	for i := len(bootstrapperNames) - 1; i >= 0; i-- {
		switch bootstrapperNames[i] {
		case fs.FileSystemBootstrapperName:
			bs = fs.NewFileSystemBootstrapper(*filePathPrefix, dbOptions, bs)
		default:
			log.Fatalf("unrecognized bootstrapper name %s", bootstrapperNames[i])
		}
	}

	writeStart := time.Now()
	opts := bootstrap.NewOptions()
	b := bootstrap.NewBootstrapProcess(opts, dbOptions, bs)

	for _, shard := range shards {
		res, err := b.Run(writeStart, shard)
		if err != nil {
			log.Fatalf("unable to bootstrap shard %d: %v", shard, err)
		}
		log.Infof("shard %d has %d series", shard, len(res.GetAllSeries()))
	}

	log.Info("bootstrapping is done")
}
