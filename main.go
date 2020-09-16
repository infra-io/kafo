// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/05 22:06:34

package main

import (
	"github.com/FishGoddess/Lighter/caches"
	"github.com/FishGoddess/Lighter/servers"
)

func main() {
	cache := caches.NewCache()
	err := servers.NewHTTPServer(cache).Run(":5837")
	if err != nil {
		panic(err)
	}
}
