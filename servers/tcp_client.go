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
	"errors"
	"strings"
	"time"

	"github.com/FishGoddess/cachego"
	"github.com/FishGoddess/vex"
	"github.com/avino-plan/kafo/caches"
	"stathat.com/c/consistent"
)

const (
	// ttlOfClient is the ttl of Client.
	ttlOfClient = 15 * 60

	// redirectPrefix is the prefix of redirect error.
	redirectPrefix = "redirect to node "

	// maxRedirectTimes is the max redirect times.
	maxRedirectTimes = 5

	// updateCircleDuration is the duration between two times of updating circle task.
	updateCircleDuration = 5 * time.Minute
)

var (
	// noClientIsAvailableErr means no client is available.
	noClientIsAvailableErr = errors.New("no client is available")

	// reachMaxRetriedTimesErr means one operation has reached max redirect times.
	reachedMaxRetriedTimesErr = errors.New("reached max redirect times")
)

// TCPClient is a tcp client for tcp server.
type TCPClient struct {

	// clients stores all clients mapping to node.
	clients *cachego.Cache

	// circle stores the relation of data and node.
	circle *consistent.Consistent
}

// NewTCPClient returns a tcp client holder connected to address.
// Returns an error if failed.
func NewTCPClient(address string) (*TCPClient, error) {

	client, err := vex.NewClient("tcp", address)
	if err != nil {
		return nil, err
	}

	circle := consistent.New()
	circle.NumberOfReplicas = 1024 // Should equals to server
	circle.Set([]string{address})

	clients := cachego.NewCache()
	clients.AutoGc(10 * time.Minute)
	clients.SetWithTTL(address, client, ttlOfClient)

	tc := &TCPClient{
		clients: clients,
		circle:  circle,
	}
	tc.updateCircleAtFixedDuration(updateCircleDuration)
	return tc, tc.updateCircleAndClients()
}

// updateCircleAtFixedDuration will update circle at fixed duration.
func (tc *TCPClient) updateCircleAtFixedDuration(duration time.Duration) {
	go func() {
		ticker := time.NewTicker(duration)
		for {
			select {
			case <-ticker.C:
				nodes, err := tc.nodes()
				if err == nil {
					tc.circle.Set(nodes)
				}
			}
		}
	}()
}

// nodes returns all nodes in cluster and an error if failed.
func (tc *TCPClient) nodes() ([]string, error) {

	nodes := tc.circle.Members()
	for _, node := range nodes {
		client, err := tc.getOrCreateClient(node)
		if err != nil {
			continue
		}
		body, err := client.Do(nodesCommand, nil)
		if err != nil {
			return nil, err
		}
		var nodes []string
		err = json.Unmarshal(body, &nodes)
		return nodes, err
	}
	return nil, noClientIsAvailableErr
}

// getOrCreateClient will get client first and create an new client if failed.
func (tc *TCPClient) getOrCreateClient(node string) (*vex.Client, error) {

	client, ok := tc.clients.Get(node)
	if !ok {
		var err error
		client, err = vex.NewClient("tcp", node)
		if err != nil {
			return nil, err
		}
		tc.clients.SetWithTTL(node, client, ttlOfClient)
	}
	return client.(*vex.Client), nil
}

// updateCircleAndClients updates circle and clients with nodes.
func (tc *TCPClient) updateCircleAndClients() error {

	nodes, err := tc.nodes()
	if err != nil {
		return err
	}

	tc.circle.Set(nodes)
	for _, node := range nodes {
		tc.getOrCreateClient(node)
	}
	return nil
}

// clientOf returns the right client of key and an error if failed.
func (tc *TCPClient) clientOf(key string) (*vex.Client, error) {
	node, err := tc.circle.Get(key)
	if err != nil {
		return nil, err
	}
	return tc.getOrCreateClient(node)
}

// doCommand will execute command with args and retry if failed.
func (tc *TCPClient) doCommand(client *vex.Client, command byte, args [][]byte) (body []byte, err error) {

	for i := 0; i < maxRedirectTimes; i++ {
		body, err := client.Do(command, args)
		if err != nil && strings.HasPrefix(err.Error(), redirectPrefix) {
			node := strings.TrimPrefix(err.Error(), redirectPrefix)
			rightClient, err := tc.getOrCreateClient(node)
			if err != nil {
				continue
			}
			client = rightClient
			continue
		}

		// An existing connection was forcibly closed by the remote host
		if err != nil && strings.HasSuffix(err.Error(), "closed by the remote host.") {
			nodes, err := tc.nodes()
			if err == nil {
				tc.circle.Set(nodes)
			}
		}
		return body, err
	}
	return nil, reachedMaxRetriedTimesErr
}

// Get returns the value of key and an error if failed.
func (tc *TCPClient) Get(key string) ([]byte, error) {
	client, err := tc.clientOf(key)
	if err != nil {
		return nil, err
	}
	return tc.doCommand(client, getCommand, [][]byte{[]byte(key)})
}

// Set adds the key and value with given ttl to cache.
// Returns an error if failed.
func (tc *TCPClient) Set(key string, value []byte, ttl int64) error {

	client, err := tc.clientOf(key)
	if err != nil {
		return err
	}

	ttlBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(ttlBytes, uint64(ttl))
	_, err = tc.doCommand(client, setCommand, [][]byte{
		ttlBytes, []byte(key), value,
	})
	return err
}

// Delete deletes the value of key and returns an error if failed.
func (tc *TCPClient) Delete(key string) error {

	client, err := tc.clientOf(key)
	if err != nil {
		return err
	}

	_, err = tc.doCommand(client, deleteCommand, [][]byte{[]byte(key)})
	return err
}

// Status returns the status of cache and an error if failed.
func (tc *TCPClient) Status() (*caches.Status, error) {

	totalStatus := caches.NewStatus()
	nodes := tc.circle.Members()
	for _, node := range nodes {
		client, err := tc.getOrCreateClient(node)
		if err != nil {
			continue
		}
		body, err := client.Do(statusCommand, nil)
		if err != nil {
			return nil, err
		}
		status := caches.NewStatus()
		err = json.Unmarshal(body, status)
		if err != nil {
			return nil, err
		}
		totalStatus.Count += status.Count
		totalStatus.KeySize += status.KeySize
		totalStatus.ValueSize += status.ValueSize
	}
	return totalStatus, nil
}

// Nodes returns the nodes of cluster and an error if failed.
func (tc *TCPClient) Nodes() ([]string, error) {
	return tc.nodes()
}

// Close closes this client and returns an error if failed.
func (tc *TCPClient) Close() (err error) {

	nodes := tc.circle.Members()
	for _, node := range nodes {
		client, ok := tc.clients.Get(node)
		if ok {
			err = client.(*vex.Client).Close()
		}
	}
	tc.clients.RemoveAll()
	return err
}
