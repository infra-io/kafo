// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/16 00:57:44

package servers

import "errors"

var (
	// commandHandlers stores all handlers mapping commands.
	commandHandlers = map[byte]func(server *TCPServer, args [][]byte) (reply byte, body []byte, err error){
		getCommand: getHandle,
		setCommand: setHandle,
	}

	// notFoundErr means not found.
	notFoundErr = errors.New("not found")

	// getCommandNeedAtLeastOneArgumentErr means get command need at least one argument.
	getCommandNeedAtLeastOneArgumentErr = errors.New("get command need at least one argument")
)

// handleFuncOf returns the handle function of command.
// Returns a false if not found.
func handleFuncOf(command byte) (func(server *TCPServer, args [][]byte) (reply byte, body []byte, err error), bool) {
	handler, ok := commandHandlers[command]
	return handler, ok
}

// getHandle handles get command with given conn and args.
func getHandle(server *TCPServer, args [][]byte) (reply byte, body []byte, err error) {

	if len(args) < 1 {
		return errorReply, nil, getCommandNeedAtLeastOneArgumentErr
	}

	value, ok := server.cache.Get(string(args[0]))
	if !ok {
		return errorReply, nil, notFoundErr
	}
	return successReply, value, nil
}

// setHandle handles set command with given conn and args.
func setHandle(server *TCPServer, args [][]byte) (reply byte, body []byte, err error) {
	return successReply, nil, nil
}
