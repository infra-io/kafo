// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/05 22:19:56

package caches

import (
	"errors"
	"sync"
	"time"
)

// Cache is a struct with caching functions.
type Cache struct {

	// data stores the real things in cache.
	data map[string]*value

	// options stores all options.
	options Options

	// status stores the status of cache.
	status *Status

	// lock is for concurrency.
	lock *sync.RWMutex
}

// NewCache returns a new Cache holder with default options.
func NewCache() *Cache {
	return NewCacheWith(DefaultOptions())
}

// NewCacheWith returns a new Cache holder with given options.
func NewCacheWith(options Options) *Cache {
	if cache, ok := recoverFromDumpFile(options.DumpFile); ok {
		return cache
	}
	return &Cache{
		data:    make(map[string]*value, 256),
		options: options,
		status:  newStatus(),
		lock:    &sync.RWMutex{},
	}
}

// recoverFromDumpFile recovers the cache from a dump file.
// Return a false if failed.
func recoverFromDumpFile(dumpFile string) (*Cache, bool) {
	cache, err := newEmptyDump().from(dumpFile)
	if err != nil {
		return nil, false
	}
	return cache, true
}

// Get returns the value of specified key.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.data[key]
	if !ok {
		return nil, false
	}

	if !value.alive() {
		c.lock.RUnlock()
		c.Delete(key)
		c.lock.RLock()
		return nil, false
	}
	return value.visit(), true
}

// Set sets an entry of specified key and value.
func (c *Cache) Set(key string, value []byte) error {
	return c.SetWithTTL(key, value, NeverDie)
}

// SetWithTTL sets an entry of specified key and value which has ttl.
func (c *Cache) SetWithTTL(key string, value []byte, ttl int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if oldValue, ok := c.data[key]; ok {
		c.status.subEntry(key, oldValue.Data)
	}

	if !c.checkEntrySize(key, value) {
		if oldValue, ok := c.data[key]; ok {
			c.status.addEntry(key, oldValue.Data)
		}
		return errors.New("the entry size will exceed if you set this entry")
	}

	c.status.addEntry(key, value)
	c.data[key] = newValue(value, ttl)
	return nil
}

// Delete deletes the specified key and value.
func (c *Cache) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if oldValue, ok := c.data[key]; ok {
		c.status.subEntry(key, oldValue.Data)
		delete(c.data, key)
	}
}

// Status returns the status of cache.
func (c *Cache) Status() Status {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return *c.status
}

// checkEntrySize checks the entry size and guarantees it will not exceed.
func (c *Cache) checkEntrySize(newKey string, newValue []byte) bool {
	return c.status.entrySize()+int64(len(newKey))+int64(len(newValue)) <= c.options.MaxEntrySize*1024*1024
}

// gc will clean up the dead entries in cache.
func (c *Cache) gc() {
	c.lock.Lock()
	defer c.lock.Unlock()
	count := 0
	for key, value := range c.data {
		if !value.alive() {
			c.status.subEntry(key, value.Data)
			delete(c.data, key)
			count++
			if count >= c.options.MaxGcCount {
				break
			}
		}
	}
}

// AutoGc starts a goroutine and run the gc task at fixed duration.
func (c *Cache) AutoGc() {
	go func() {
		ticker := time.NewTicker(time.Duration(c.options.GcDuration) * time.Minute)
		for {
			select {
			case <-ticker.C:
				c.gc()
			}
		}
	}()
}

// dump dumps c to dumpFile and returns an error if failed.
func (c *Cache) dump() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return newDump(c).to(c.options.DumpFile)
}

// AutoDump starts a goroutine and run the dump task at fixed duration.
func (c *Cache) AutoDump() {
	go func() {
		ticker := time.NewTicker(time.Duration(c.options.DumpDuration) * time.Minute)
		for {
			select {
			case <-ticker.C:
				c.dump()
			}
		}
	}()
}
