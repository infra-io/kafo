// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/05 22:19:56

package caches

import (
	"sync"
)

// Cache is a struct with caching functions.
type Cache struct {

	// data stores the real things in cache.
	data map[string]*value

	status *Status

	// lock is for concurrency.
	lock *sync.RWMutex
}

// NewCache returns a new Cache holder.
func NewCache() *Cache {
	return &Cache{
		data:   make(map[string]*value, 256),
		status: newStatus(),
		lock:   &sync.RWMutex{},
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
func (c *Cache) Set(key string, value []byte) {
	c.SetWithTTL(key, value, NeverDie)
}

// SetWithTTL sets an entry of specified key and value which has ttl.
func (c *Cache) SetWithTTL(key string, value []byte, ttl int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if oldValue, ok := c.data[key]; ok {
		c.status.subEntry(key, oldValue.data)
	}
	c.status.addEntry(key, value)
	c.data[key] = newValue(value, ttl)
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
