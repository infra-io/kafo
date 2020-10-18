// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/15 22:45:56

package servers

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/FishGoddess/kafo/caches"
	"github.com/FishGoddess/vex"
)

const (
	// getCommand is the command of get operation.
	getCommand = byte(1)

	// setCommand is the command of get operation.
	setCommand = byte(2)

	// deleteCommand is the command of get operation.
	deleteCommand = byte(3)

	// statusCommand is the command of get operation.
	statusCommand = byte(4)
)

var (
	// commandNeedsMoreArgumentsErr means command needs more arguments.
	commandNeedsMoreArgumentsErr = errors.New("command needs more arguments")

	// notFoundErr means not found.
	notFoundErr = errors.New("not found")
)

// TCPServer is a tcp type server.
type TCPServer struct {
	// cache is the real cache used inside.
	cache *caches.Cache

	// server is the real tcp server used inside.
	server *vex.Server
}

// NewTCPServer returns a tcp server holder.
func NewTCPServer(cache *caches.Cache) *TCPServer {
	return &TCPServer{
		cache:  cache,
		server: vex.NewServer(),
	}
}

// Run runs the server at address and returns an error if something wrong.
func (ts *TCPServer) Run(address string) error {
	ts.server.RegisterHandler(getCommand, ts.getHandler)
	ts.server.RegisterHandler(setCommand, ts.setHandler)
	ts.server.RegisterHandler(deleteCommand, ts.deleteHandler)
	ts.server.RegisterHandler(statusCommand, ts.statusHandler)
	return ts.server.ListenAndServe("tcp", address)
}

// Close closes the server and releases resources.
func (ts *TCPServer) Close() error {
	return ts.server.Close()
}

// =======================================================================

// getHandler is a handler for getting value of specified key.
func (ts *TCPServer) getHandler(args [][]byte) (body []byte, err error) {
	if len(args) < 1 {
		return nil, commandNeedsMoreArgumentsErr
	}

	value, ok := ts.cache.Get(string(args[0]))
	if !ok {
		return value, notFoundErr
	}
	return value, nil
}

// setHandler is a handler for setting an entry of specified key and value.
func (ts *TCPServer) setHandler(args [][]byte) (body []byte, err error) {
	if len(args) < 3 {
		return nil, commandNeedsMoreArgumentsErr
	}

	ttl := int64(binary.BigEndian.Uint64(args[0]))
	err = ts.cache.SetWithTTL(string(args[1]), args[2], ttl)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// deleteHandler is a handler for deleting the entry of specified key.
func (ts *TCPServer) deleteHandler(args [][]byte) (body []byte, err error) {
	if len(args) < 1 {
		return nil, commandNeedsMoreArgumentsErr
	}

	err = ts.cache.Delete(string(args[0]))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// statusHandler is handler for fetching the status of cache.
func (ts *TCPServer) statusHandler(args [][]byte) (body []byte, err error) {
	return json.Marshal(ts.cache.Status())
}
