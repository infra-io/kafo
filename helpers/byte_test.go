// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/05 23:21:36

package helpers

import "testing"

// go test -cover -run=^TestCopy$
func TestCopy(t *testing.T) {

	src := []byte{
		1, 2, 3, 4, 5,
	}

	dst := Copy(src)
	if &dst == &src {
		t.Fatal("&dst == &src")
	}

	if len(dst) != len(src) {
		t.Fatal("len(dst) != len(src)...")
	}

	for i := range dst {
		if dst[i] != src[i] {
			t.Fatalf("dst[%d] != src[%d]...", i, i)
		}
	}
}
