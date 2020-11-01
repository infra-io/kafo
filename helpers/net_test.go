// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/27 21:33:28

package helpers

import "testing"

// go test -cover -run=^TestJoinAddressAndPort$
func TestJoinAddressAndPort(t *testing.T) {
	address := "127.0.0.1"
	port := 5837
	result := JoinAddressAndPort(address, port)
	if result != "127.0.0.1:5837" {
		t.Fatalf("Result %s is wrong!", result)
	}
}
