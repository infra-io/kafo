// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/05 22:06:34

package main

import (
	"flag"

	"github.com/FishGoddess/kafo/caches"
	"github.com/FishGoddess/kafo/servers"
)

func main() {

	// Parse all flags
	address := flag.String("address", ":5837", "The address used to listen, such as 127.0.0.1:5837.")
	maxEntrySize := flag.Int64("maxEntrySize", caches.DefaultMaxEntrySize, "The max memory size that entries can use.")
	flag.Parse()

	// Initialize
	cache := caches.NewCacheWithMaxEntrySize(*maxEntrySize)
	err := servers.NewHTTPServer(cache).Run(*address)
	if err != nil {
		panic(err)
	}
}
