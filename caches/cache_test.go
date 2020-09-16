// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/05 23:31:09

package caches

import (
	"testing"
	"time"
)

// go test -cover -run=^TestCache$
func TestCache(t *testing.T) {

	cache := NewCache()
	if cache.Status().Count != 0 {
		t.Fatal("cache.Count() != 0...")
	}

	k := "key"
	if value, ok := cache.Get(k); ok {
		t.Fatalf("cache.Get(\"key\") = %s, but this should not happen...", string(value))
	}

	v := "value"
	cache.Set(k, []byte(v))
	if cache.Status().Count != 1 {
		t.Fatal("cache.Count() != 1...")
	}

	value, ok := cache.Get(k)
	if !ok || string(value) != v {
		t.Fatalf("ok = %v, value = %s, but ok should be true and value should be %s...", ok, string(value), v)
	}
}

// go test -cover -run=^TestCacheTTL$
func TestCacheTTL(t *testing.T) {

	cache := NewCache()

	k := "key"
	v := "value"
	cache.SetWithTTL(k, []byte(v), 2)

	if value, ok := cache.Get(k); !ok || string(value) != v {
		t.Fatalf("ok = %v, value = %s, but ok should be true and value should be %s...", ok, string(value), v)
	}

	time.Sleep(3 * time.Second)
	if value, ok := cache.Get(k); ok {
		t.Fatalf("cache.Get(\"key\") = %s, but this should not happen...", string(value))
	}
}

// go test -cover -run=^TestCacheStatus$
func TestCacheStatus(t *testing.T) {

	cache := NewCache()
	status := cache.Status()
	if status.Count != 0 || status.KeySize != 0 || status.ValueSize != 0 {
		t.Fatal("cache.Status() should be zero!")
	}

	k := "key"
	v := "value"
	cache.Set(k, []byte(v))
	status = cache.Status()
	if status.Count != 1 || status.KeySize != 3 || status.ValueSize != 5 {
		t.Fatal("cache.Status() returns wrong values!")
	}

	status.Count = 999
	status = cache.Status()
	if status.Count != 1 {
		t.Fatal("cache.status can be changed outside!")
	}
}
