// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/15 23:34:44

package caches

import (
	"encoding/json"
	"testing"
)

// go test -cover -run=^TestStatus$
func TestStatus(t *testing.T) {

	status := NewStatus()
	if status.Count != 0 || status.KeySize != 0 || status.ValueSize != 0 {
		t.Fatalf("The new status should be zero! Count is %d. KeySize is %d. ValueSize is %d.",
			status.Count, status.KeySize, status.ValueSize)
	}

	status.addEntry("key", []byte("value"))

	if status.Count != 1 || status.KeySize != 3 || status.ValueSize != 5 {
		t.Fatalf("The status is wrong! Count is %d. KeySize is %d. ValueSize is %d.",
			status.Count, status.KeySize, status.ValueSize)
	}

	if status.entrySize() != 8 {
		t.Fatalf("The entrySize is wrong! Got %d size.", status.entrySize())
	}

	statusJson, err := json.Marshal(status)
	if err != nil {
		t.Fatal(err)
	}

	if string(statusJson) != `{"count":1,"keySize":3,"valueSize":5}` {
		t.Fatal(string(statusJson))
	}
}
