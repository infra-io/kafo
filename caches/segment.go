// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/09 21:42:13

package caches

import (
	"errors"
	"sync"
)

// segment is the struct storing the real data.
type segment struct {

	// Data stores the real things of segment.
	Data map[string]*value

	// Status stores the status of segment.
	Status *Status

	// options stores all options.
	options *Options

	// lock is for concurrency.
	lock *sync.RWMutex
}

// newSegment returns a segment holder with options.
func newSegment(options *Options) *segment {
	return &segment{
		Data:    make(map[string]*value, options.MapSizeOfSegment),
		Status:  newStatus(),
		options: options,
		lock:    &sync.RWMutex{},
	}
}

// get returns the value of specified key.
func (s *segment) get(key string) ([]byte, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	value, ok := s.Data[key]
	if !ok {
		return nil, false
	}

	if !value.alive() {
		s.lock.RUnlock()
		s.delete(key)
		s.lock.RLock()
		return nil, false
	}
	return value.visit(), true
}

// set sets an entry of specified key and value which has ttl.
func (s *segment) set(key string, value []byte, ttl int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if oldValue, ok := s.Data[key]; ok {
		s.Status.subEntry(key, oldValue.Data)
	}

	if !s.checkEntrySize(key, value) {
		if oldValue, ok := s.Data[key]; ok {
			s.Status.addEntry(key, oldValue.Data)
		}
		return errors.New("the entry size will exceed if you set this entry")
	}

	s.Status.addEntry(key, value)
	s.Data[key] = newValue(value, ttl)
	return nil
}

// delete deletes the specified key and value.
func (s *segment) delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if oldValue, ok := s.Data[key]; ok {
		s.Status.subEntry(key, oldValue.Data)
		delete(s.Data, key)
	}
}

// checkEntrySize checks the entry size and guarantees it will not exceed.
func (s *segment) checkEntrySize(newKey string, newValue []byte) bool {
	return s.Status.entrySize()+int64(len(newKey))+int64(len(newValue)) <= int64((s.options.MaxEntrySize*1024*1024) / s.options.SegmentSize)
}

// gc will clean up the dead entries in segment.
func (s *segment) gc() {
	s.lock.Lock()
	defer s.lock.Unlock()
	count := 0
	for key, value := range s.Data {
		if !value.alive() {
			s.Status.subEntry(key, value.Data)
			delete(s.Data, key)
			count++
			if count >= s.options.MaxGcCount {
				break
			}
		}
	}
}
