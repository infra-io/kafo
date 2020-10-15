// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/15 22:45:56

package servers

import (
	"net"
	"strings"
	"sync"

	"github.com/FishGoddess/kafo/caches"
)

// TCPServer is a tcp type server.
type TCPServer struct {
	// cache is the real cache used inside.
	cache *caches.Cache

	// listener is the listener for accepting connections.
	listener net.Listener
}

// NewTCPServer returns a tcp server holder.
func NewTCPServer(cache *caches.Cache) *TCPServer {
	return &TCPServer{
		cache: cache,
	}
}

// Run runs the server at address and returns an error if something wrong.
func (ts *TCPServer) Run(address string) error {

	var err error
	ts.listener, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for {
		conn, err := ts.listener.Accept()
		if err != nil {
			// This error means listener has been closed
			// See src/internal/poll/fd.go@ErrNetClosing
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
			continue
		}

		wg.Add(1)
		go func(c net.Conn) {
			defer wg.Done()
			ts.handleConn(c)
		}(conn)
	}

	wg.Wait()
	return nil
}

// handleConn handles an accepted connection.
func (ts *TCPServer) handleConn(conn net.Conn) {

	defer conn.Close()
	for {
		command, args, err := readRequestFrom(conn)
		if err != nil {
			writeErrorResponseTo(conn, err.Error())
			return
		}

		handle, ok := handleFuncOf(command)
		if !ok {
			writeErrorResponseTo(conn, "the command is invalid")
			continue
		}

		reply, body, err := handle(ts, args)
		if err != nil {
			writeErrorResponseTo(conn, err.Error())
			continue
		}
		writeResponseTo(conn, reply, body)
	}
}
