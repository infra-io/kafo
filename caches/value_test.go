// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/08 23:49:27

package caches

import (
	"testing"
	"time"
)

// go test -cover -run=^TestValue$
func TestValue(t *testing.T) {

	v := newValue([]byte{}, 1)
	if !v.alive() {
		t.Fatalf("%v should be alive!", v)
	}

	time.Sleep(2 * time.Second)
	if v.alive() {
		t.Fatalf("%v should be dead!", v)
	}

	v = newValue([]byte{}, 1)
	time.Sleep(2 * time.Second)
	v.visit()
	if !v.alive() {
		t.Fatalf("%v should be alive after visiting!", v)
	}
}
