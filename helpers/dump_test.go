// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/20 22:05:45

package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

// go test -cover -run=^TestDump$
func TestDump(t *testing.T) {

	m := map[string]int{
		"A": 1,
		"B": 2,
		"C": 3,
	}

	dumpFile := filepath.Join(os.TempDir(), "m.dump")
	err := Marshal(m, dumpFile)
	if err != nil {
		t.Fatal(err)
	}

	var mm map[string]int
	err = Unmarshal(&mm, dumpFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(mm) != 3 || mm["A"] != 1 || mm["B"] != 2 || mm["C"] != 3 {
		t.Fatal("Unmarshal is failed!", mm)
	}
}
