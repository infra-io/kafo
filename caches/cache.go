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
)

const (
	// DefaultMaxEntrySize is the max memory size that entries can use in default.
	DefaultMaxEntrySize = int64(4294967296)
)

// Cache is a struct with caching functions.
type Cache struct {

	// data stores the real things in cache.
	data map[string]*value

	// maxEntrySize is the max memory size that entries can use.
	maxEntrySize int64

	// status stores the status of cache.
	status *Status

	// lock is for concurrency.
	lock *sync.RWMutex
}

// NewCache returns a new Cache holder.
func NewCache() *Cache {
	return NewCacheWithMaxEntrySize(DefaultMaxEntrySize)
}

func NewCacheWithMaxEntrySize(maxEntrySize int64) *Cache {
	return &Cache{
		data:         make(map[string]*value, 256),
		maxEntrySize: maxEntrySize,
		status:       newStatus(),
		lock:         &sync.RWMutex{},
	}
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
	return value.visit(), ok
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
		c.status.subEntry(key, oldValue.data)
	}

	if !c.checkEntrySize(key, value) {
		if oldValue, ok := c.data[key]; ok {
			c.status.addEntry(key, oldValue.data)
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
		c.status.subEntry(key, oldValue.data)
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
	return c.status.entrySize()+int64(len(newKey))+int64(len(newValue)) <= c.maxEntrySize
}