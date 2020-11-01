// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/20 23:30:35

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/FishGoddess/kafo/servers"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Missing args!")
		os.Exit(1)
	}

	client, err := servers.NewTCPClient(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer client.Close()

	command := os.Args[2]
	args := os.Args[3:]
	switch command {
	case "get":
		doGet(client, args)
	case "set":
		doSet(client, args)
	case "delete":
		doDelete(client, args)
	case "status":
		doStatus(client, args)
	case "nodes":
		doNodes(client, args)
	default:
		fmt.Println("Command not found!")
		os.Exit(1)
	}
}

func doGet(client *servers.TCPClient, args []string) {

	if len(args) < 1 {
		fmt.Println("Get needs 1 argument!")
		os.Exit(1)
	}
	value, err := client.Get(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(value))
}

func doSet(client *servers.TCPClient, args []string) {

	if len(args) < 3 {
		fmt.Println("Set needs 3 arguments!")
		os.Exit(1)
	}

	ttl, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fmt.Println("TTL is not an integer!")
		os.Exit(1)
	}

	err = client.Set(args[1], []byte(args[2]), ttl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Done")
}

func doDelete(client *servers.TCPClient, args []string) {

	if len(args) < 1 {
		fmt.Println("Delete needs 1 argument!")
		os.Exit(1)
	}

	err := client.Delete(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Done")
}

func doStatus(client *servers.TCPClient, args []string) {
	status, err := client.Status()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("count: %d, keySize: %d, valueSize: %d", status.Count, status.KeySize, status.ValueSize)
}

func doNodes(client *servers.TCPClient, args []string) {
	nodes, err := client.Nodes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("nodes: %v", nodes)
}
