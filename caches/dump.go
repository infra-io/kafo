// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/20 23:40:29

package caches

import (
	"encoding/gob"
	"os"
	"sync"
)

// dump is for dumping the cache.
type dump struct {

	// Data stores the real things in cache.
	Data map[string]*value

	// Options stores all options.
	Options Options

	// Status stores the status of cache.
	Status *Status
}

// newEmptyDump return an empty dump holder.
func newEmptyDump() *dump {
	return &dump{}
}

// newDump returns a dump holder of c.
func newDump(c *Cache) *dump {
	return &dump{
		Data:    c.data,
		Options: c.options,
		Status:  c.status,
	}
}

// to dumps d to dumpFile and returns an error if failed.
func (d *dump) to(dumpFile string) error {
	file, err := os.OpenFile(dumpFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return gob.NewEncoder(file).Encode(d)
}

// from returns a Cache holder parsed from d of dumpFile.
func (d *dump) from(dumpFile string) (*Cache, error) {
	file, err := os.Open(dumpFile)
	if err != nil || gob.NewDecoder(file).Decode(d) != nil {
		return nil, err
	}
	return &Cache{
		data:    d.Data,
		options: d.Options,
		status:  d.Status,
		lock:    &sync.RWMutex{},
	}, nil
}
