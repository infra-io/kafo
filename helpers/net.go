// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/27 21:31:02

package helpers

import "strconv"

// JoinAddressAndPort joins address and port with ":".
func JoinAddressAndPort(address string, port int) string {
	return address + ":" + strconv.Itoa(port)
}