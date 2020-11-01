// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/05 22:06:34

package main

import (
	"flag"
	"log"
	"strings"

	"github.com/avino-plan/kafo/caches"
	"github.com/avino-plan/kafo/servers"
)

func main() {

	// Prepare options
	serverOptions := servers.DefaultOptions()
	flag.StringVar(&serverOptions.Address, "address", serverOptions.Address, "The address used to listen, such as 127.0.0.1.")
	flag.IntVar(&serverOptions.Port, "port", serverOptions.Port, "The port used to listen, such as 5837.")
	flag.StringVar(&serverOptions.ServerType, "serverType", serverOptions.ServerType, "The type of server (http, tcp).")
	flag.IntVar(&serverOptions.VirtualNodeCount, "virtualNodeCount", serverOptions.VirtualNodeCount, "The number of virtual nodes in consistent hash.")
	flag.IntVar(&serverOptions.UpdateCircleDuration, "updateCircleDuration", serverOptions.UpdateCircleDuration, "The duration between two circle updating operations. The unit is second.")
	cluster := flag.String("cluster", "", "The cluster of servers. One node in cluster will be ok.")

	cacheOptions := caches.DefaultOptions()
	flag.IntVar(&cacheOptions.MaxEntrySize, "maxEntrySize", cacheOptions.MaxEntrySize, "The max memory size that entries can use. The unit is GB.")
	flag.IntVar(&cacheOptions.MaxGcCount, "maxGcCount", cacheOptions.MaxGcCount, "The max count of entries that gc will clean.")
	flag.IntVar(&cacheOptions.GcDuration, "gcDuration", cacheOptions.GcDuration, "The duration between two gc tasks. The unit is Minute.")
	flag.StringVar(&cacheOptions.DumpFile, "dumpFile", cacheOptions.DumpFile, "The file used to dump the cache.")
	flag.IntVar(&cacheOptions.DumpDuration, "dumpDuration", cacheOptions.DumpDuration, "The duration between two dump tasks. The unit is Minute.")
	flag.IntVar(&cacheOptions.MapSizeOfSegment, "mapSizeOfSegment", cacheOptions.MapSizeOfSegment, "The map size of segment.")
	flag.IntVar(&cacheOptions.SegmentSize, "segmentSize", cacheOptions.SegmentSize, "The number of segment in a cache. This value should be the pow of 2 for precision.")
	flag.IntVar(&cacheOptions.CasSleepTime, "casSleepTime", cacheOptions.CasSleepTime, "The time of sleep in one cas step. The unit is Microsecond.")
	flag.Parse()

	serverOptions.Cluster = nodesInCluster(*cluster)

	// Initialize
	cache := caches.NewCacheWith(cacheOptions)
	cache.AutoGc()
	cache.AutoDump()

	server, err := servers.NewServer(cache, serverOptions)
	if err != nil {
		panic(err)
	}

	log.Printf("Using server options %+v\n", serverOptions)
	log.Printf("Using cache options %+v\n", cacheOptions)
	log.Printf("Kafo is running on %s at %s:%d.", serverOptions.ServerType, serverOptions.Address, serverOptions.Port)
	err = server.Run()
	if err != nil {
		panic(err)
	}
}

// nodesInCluster parses cluster to nodes with comma (,) and returns nil if cluster is "".
func nodesInCluster(cluster string) []string {
	if cluster == "" {
		return nil
	}
	return strings.Split(cluster, ",")
}
