// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/08 23:36:12

package caches

import (
	"sync/atomic"
	"time"

	"github.com/FishGoddess/kafo/helpers"
)

const (
	// NeverDie means value.alive() returns true forever.
	NeverDie = 0
)

// value is a box of data.
type value struct {

	// data stores the real thing inside.
	Data []byte

	// ttl is the life of value.
	// The unit is second.
	Ttl int64

	// ctime is the created time of value.
	Ctime int64
}

// newValue returns a new value with data and ttl.
func newValue(data []byte, ttl int64) *value {
	return &value{
		Data:  helpers.Copy(data),
		Ttl:   ttl,
		Ctime: time.Now().Unix(),
	}
}

// alive returns if this value is alive or not.
func (v *value) alive() bool {
	return v.Ttl == NeverDie || time.Now().Unix()-v.Ctime < v.Ttl
}

// visit updates the ctime of value to now.
func (v *value) visit() []byte {
	atomic.SwapInt64(&v.Ctime, time.Now().Unix())
	return v.Data
}
