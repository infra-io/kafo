// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/06 01:04:11

package servers

import "github.com/FishGoddess/kafo/caches"

const (
	// APIVersion is the version of serving APIs.
	APIVersion = "v1"
)

// Server is an interface of servers.
type Server interface {

	// Run runs a server on specified address.
	Run(address string) error
}

// NewServer returns a server of serverType.
func NewServer(serverType string, cache *caches.Cache) Server {
	if serverType == "tcp" {
		return NewTCPServer(cache)
	}
	return NewHTTPServer(cache)
}
