// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/15 23:59:31

package servers

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"net"
)

// The protocol isn't concurrency-safe in request, so clients should be responsible for this.
// Notice: All number used is BigEndian.

// A request will be like this:
// ***********************************************************
// version    command    argsLength    {argLength    arg}    *
//  1byte      1byte       1byte         4byte     unknown   *
// ***********************************************************

// A response will be like this:
// ***************************
// version    reply    body  *
//  1byte     1byte   1byte  *
// ***************************

const (
	// ProtocolVersion is the version of protocol.
	ProtocolVersion = byte(1)
)

// All supported commands here.
const (
	// getCommand is the command of get operation.
	getCommand = byte(1)

	// getCommand is the command of get operation.
	setCommand = byte(2)

	// getCommand is the command of get operation.
	deleteCommand = byte(3)

	// getCommand is the command of get operation.
	statusCommand = byte(101)
)

// All supported replies here.
const (
	// successReply is the reply of success.
	successReply = byte(1)

	// errorReply is the reply of error.
	errorReply = byte(2)
)

var (
	// protocolVersionMismatchErr means the version between client and server is mismatched.
	protocolVersionMismatchErr = errors.New("the version between client and server is mismatched")
)

// readRequestFrom reads a request from conn and returns the command and arguments.
// If the version between client and server is mismatched, return an protocolVersionMismatchErr.
// Otherwise, return an error if read failed.
func readRequestFrom(conn net.Conn) (command byte, args [][]byte, err error) {

	reader := bufio.NewReader(conn)

	// header has version, command and argsLength
	header := make([]byte, 3)
	_, err = io.ReadFull(reader, header)
	if err != nil {
		return 0, nil, err
	}

	// Check the version
	if header[0] != ProtocolVersion {
		return 0, nil, protocolVersionMismatchErr
	}

	// Read all arguments
	argsLength := header[2]
	args = make([][]byte, argsLength)

	argLength := make([]byte, 4)
	for i := byte(0); i < argsLength; i++ {

		// Read argument's length first
		_, err := io.ReadFull(reader, argLength)
		if err != nil {
			return 0, nil, err
		}

		// Then read argument, using BigEndian
		arg := make([]byte, binary.BigEndian.Uint32(argLength))
		_, err = io.ReadFull(reader, arg)
		if err != nil {
			return 0, nil, err
		}

		args[i] = arg
	}
	return header[1], args, nil
}

func writeResponseTo(conn net.Conn, reply byte, body []byte) error {
	return nil
}

func writeErrorResponseTo(conn net.Conn, msg string) error {
	return nil
}
