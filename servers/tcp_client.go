// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/17 21:00:30

package servers

import (
	"encoding/binary"
	"encoding/json"

	"github.com/FishGoddess/kafo/caches"
	"github.com/FishGoddess/vex"
)

// TCPClient is a tcp client for tcp server.
type TCPClient struct {
	// client is the real client to connect tcp server.
	client *vex.Client
}

// NewTCPClient returns a tcp client holder connected to address.
// Returns an error if failed.
func NewTCPClient(address string) (*TCPClient, error) {
	client, err := vex.NewClient("tcp", address)
	if err != nil {
		return nil, err
	}
	return &TCPClient{
		client: client,
	}, nil
}

// Get returns the value of key and an error if failed.
func (tc *TCPClient) Get(key string) ([]byte, error) {
	return tc.client.Do(getCommand, [][]byte{[]byte(key)})
}

// Set adds the key and value with given ttl to cache.
// Returns an error if failed.
func (tc *TCPClient) Set(key string, value []byte, ttl int64) error {
	ttlBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(ttlBytes, uint64(ttl))
	_, err := tc.client.Do(setCommand, [][]byte{
		ttlBytes, []byte(key), value,
	})
	return err
}

// Delete deletes the value of key and returns an error if failed.
func (tc *TCPClient) Delete(key string) error {
	_, err := tc.client.Do(deleteCommand, [][]byte{[]byte(key)})
	return err
}

// Status returns the status of cache and an error if failed.
func (tc *TCPClient) Status() (*caches.Status, error) {
	body, err := tc.client.Do(statusCommand, nil)
	if err != nil {
		return nil, err
	}
	status := caches.NewStatus()
	err = json.Unmarshal(body, status)
	return status, err
}

// Close closes this client and returns an error if failed.
func (tc *TCPClient) Close() error {
	return tc.client.Close()
}
