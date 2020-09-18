// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/18 21:32:01

package caches

// Options is the struct of options.
type Options struct {

	// MaxEntrySize is the max memory size that entries can use.
	// The unit is GB.
	MaxEntrySize int64

	// MaxGcCount is the max count of entries that gc will clean.
	MaxGcCount int

	// GcDuration is the duration between two gc tasks.
	// The unit is Minute.
	GcDuration int64
}

// DefaultOptions returns a default options.
func DefaultOptions() Options {
	return Options{
		MaxEntrySize: int64(4), // 4 GB
		MaxGcCount:   1000,
		GcDuration:   60, // 1 hour
	}
}
