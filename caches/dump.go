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
	"time"
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

// nowSuffix returns a string of current time formatted as 20060102150405.
func nowSuffix() string {
	return "." + time.Now().Format("20060102150405")
}

// to dumps d to dumpFile and returns an error if failed.
func (d *dump) to(dumpFile string) error {
	newDumpFile := dumpFile + nowSuffix()
	file, err := os.OpenFile(newDumpFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = gob.NewEncoder(file).Encode(d)
	if err != nil {
		file.Close()
		os.Remove(newDumpFile)
		return err
	}

	os.Remove(dumpFile)
	file.Close()
	return os.Rename(newDumpFile, dumpFile)
}

// from returns a Cache holder parsed from d of dumpFile.
func (d *dump) from(dumpFile string) (*Cache, error) {
	file, err := os.Open(dumpFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err = gob.NewDecoder(file).Decode(d); err != nil {
		return nil, err
	}

	return &Cache{
		data:    d.Data,
		options: d.Options,
		status:  d.Status,
		lock:    &sync.RWMutex{},
	}, nil
}
