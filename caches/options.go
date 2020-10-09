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
	MaxEntrySize int

	// MaxGcCount is the max count of entries that gc will clean.
	MaxGcCount int

	// GcDuration is the duration between two gc tasks.
	// The unit is Minute.
	GcDuration int

	// DumpFile is the file used to dump the cache.
	DumpFile string

	// DumpDuration is the duration between two dump tasks.
	// The unit is Minute.
	DumpDuration int

	// MapSizeOfSegment is the map size of segment.
	MapSizeOfSegment int

	// SegmentSize is the number of segment in a cache.
	// This value should be the pow of 2 for precision.
	SegmentSize int
}

// DefaultOptions returns a default options.
func DefaultOptions() Options {
	return Options{
		MaxEntrySize:     4, // 4 GB
		MaxGcCount:       1000,
		GcDuration:       60, // 1 hour
		DumpFile:         "kafo.dump",
		DumpDuration:     30, // 30 minutes
		MapSizeOfSegment: 256,
		SegmentSize:      1024,
	}
}
