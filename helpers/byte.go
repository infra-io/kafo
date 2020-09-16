// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/05 22:36:07

package helpers

// Copy returns a new byte slice copy from src.
func Copy(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
