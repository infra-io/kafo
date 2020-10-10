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
	"sync/atomic"
	"time"
)

// Cache is a struct with caching functions.
type Cache struct {

	// segmentSize is the size of segments.
	// This value will affect the performance of concurrency.
	segmentSize int

	// segments is a slice stores the real data.
	segments []*segment

	// options stores all options.
	options *Options

	// dumping means if cache is in dumping status.
	// 1 is dumping.
	dumping int32
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
		segmentSize: options.SegmentSize,
		segments:    newSegments(&options),
		options:     &options,
		dumping:     0,
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

// newSegments returns a slice of initialized segments.
func newSegments(options *Options) []*segment {
	segments := make([]*segment, options.SegmentSize)
	for i := 0; i < options.SegmentSize; i++ {
		segments[i] = newSegment(options)
	}
	return segments
}

// index returns a position in segments of this key.
func index(key string) int {
	index := 0
	keyBytes := []byte(key)
	for _, b := range keyBytes {
		index = 31*index + int(b&0xff)
	}
	return index
}

// segmentOf returns the segment of this key.
func (c *Cache) segmentOf(key string) *segment {
	return c.segments[index(key)&(c.segmentSize-1)]
}

// Get returns the value of specified key.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.waitForDumping()
	return c.segmentOf(key).get(key)
}

// Set sets an entry of specified key and value.
func (c *Cache) Set(key string, value []byte) error {
	return c.SetWithTTL(key, value, NeverDie)
}

// SetWithTTL sets an entry of specified key and value which has ttl.
func (c *Cache) SetWithTTL(key string, value []byte, ttl int64) error {
	c.waitForDumping()
	return c.segmentOf(key).set(key, value, ttl)
}

// Delete deletes the specified key and value.
func (c *Cache) Delete(key string) {
	c.waitForDumping()
	c.segmentOf(key).delete(key)
}

// Status returns the status of cache.
func (c *Cache) Status() Status {
	result := newStatus()
	for _, segment := range c.segments {
		status := segment.status()
		result.Count += status.Count
		result.KeySize += status.KeySize
		result.ValueSize += status.ValueSize
	}
	return *result
}

// gc will clean up the dead entries in cache.
func (c *Cache) gc() {
	c.waitForDumping()
	wg := &sync.WaitGroup{}
	for _, seg := range c.segments {
		wg.Add(1)
		go func(s *segment) {
			defer wg.Done()
			s.gc()
		}(seg)
	}
	wg.Wait()
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
	atomic.StoreInt32(&c.dumping, 1)
	defer atomic.StoreInt32(&c.dumping, 0)
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

// waitForDumping will wait for dumping.
func (c *Cache) waitForDumping() {
	for atomic.LoadInt32(&c.dumping) != 0 {
		time.Sleep(time.Duration(c.options.CasSleepTime) * time.Microsecond)
	}
}
