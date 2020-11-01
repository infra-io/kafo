// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/27 00:17:17

package servers

// Options is the options of servers.
type Options struct {

	// Address is the address used to listen.
	Address string

	// Port is the port used to listen.
	Port int

	// ServerType is the type of server.
	ServerType string

	// VirtualNodeCount is the number of virtual nodes, which is set to consistent hash.
	VirtualNodeCount int

	// UpdateCircleDuration is the duration between two circle updating operations.
	// The unit is second.
	UpdateCircleDuration int

	// cluster is all nodes in cluster that will be joined.
	Cluster []string
}

// DefaultOptions returns a default options.
func DefaultOptions() Options {
	return Options{
		Address:              "127.0.0.1",
		Port:                 5837,
		ServerType:           "tcp",
		VirtualNodeCount:     1024,
		UpdateCircleDuration: 3, // 3 Seconds
	}
}
