// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/21 23:38:03

package caches

import (
	"os"
	"path/filepath"
	"testing"
)

// go test -cover -run=^TestDump$
func TestDump(t *testing.T) {

	cache := NewCache()
	if err := cache.Set("key", []byte("value")); err != nil {
		t.Fatal(err)
	}

	dumpFile := filepath.Join(os.TempDir(), "TestDump.dump")
	if err := newDump(cache).to(dumpFile); err != nil {
		t.Fatal(err)
	}

	cache, err := newEmptyDump().from(dumpFile)
	if err != nil {
		t.Fatal(err)
	}

	value, ok := cache.Get("key")
	if !ok || string(value) != "value" {
		t.Fatal("Get key from cache is wrong!", ok, string(value))
	}

	if cache.Status().Count != 1 {
		t.Fatalf("The status of cache is wrong! Status is %+v.", cache.Status())
	}
}
