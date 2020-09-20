// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/20 23:40:29

package caches

import "sync"

// dump is for dumping the cache.
type dump struct {

	// Data stores the real things in cache.
	Data map[string]*value

	// Options stores all options.
	Options Options

	// Status stores the status of cache.
	Status *Status
}

// newDump returns a dump holder of c.
func newDump(c *Cache) *dump {
	return &dump{
		Data:    c.data,
		Options: c.options,
		Status:  c.status,
	}
}

// toCache returns a Cache holder parsed from d.
func (d *dump) toCache() *Cache {
	return &Cache{
		data:    d.Data,
		options: d.Options,
		status:  d.Status,
		lock:    &sync.RWMutex{},
	}
}
