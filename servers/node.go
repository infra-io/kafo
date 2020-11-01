// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/26 23:34:49

package servers

import (
	"io/ioutil"
	"time"

	"github.com/FishGoddess/kafo/helpers"
	"github.com/hashicorp/memberlist"
	"stathat.com/c/consistent"
)

// node isn't only a node of cluster but also a node of consistent hash.
type node struct {

	// options stores all settings of node.
	options *Options

	// address contains host/ip and port.
	address string

	// circle is for handling consistent hash.
	circle *consistent.Consistent

	// nodeManager is for managing all nodes.
	nodeManager *memberlist.Memberlist
}

// newNode returns a new node for use with given options and an error if failed.
// This returned node will join to the cluster in options.
// If cluster is nil or len(cluster) == 0, address in options will be added to cluster,
// which means this node becomes a new cluster.
func newNode(options *Options) (*node, error) {

	if options.Cluster == nil || len(options.Cluster) == 0 {
		options.Cluster = []string{options.Address}
	}

	nodeManager, err := createNodeManager(options)
	if err != nil {
		return nil, err
	}

	node := &node{
		options:     options,
		address:     helpers.JoinAddressAndPort(options.Address, options.Port),
		circle:      consistent.New(),
		nodeManager: nodeManager,
	}
	node.circle.NumberOfReplicas = options.VirtualNodeCount
	node.autoUpdateCircle()
	return node, nil
}

// createNodeManager creates a new node manager and joins the cluster.
// Returns an error if failed.
func createNodeManager(options *Options) (*memberlist.Memberlist, error) {

	config := memberlist.DefaultLANConfig()
	config.Name = helpers.JoinAddressAndPort(options.Address, options.Port)
	config.BindAddr = options.Address
	config.LogOutput = ioutil.Discard // disable logging

	nodeManager, err := memberlist.Create(config)
	if err != nil {
		return nil, err
	}

	_, err = nodeManager.Join(options.Cluster)
	return nodeManager, err
}

// nodes returns all names of nodes in cluster.
func (n *node) nodes() []string {
	members := n.nodeManager.Members()
	nodes := make([]string, len(members))
	for i, member := range members {
		nodes[i] = member.Name
	}
	return nodes
}

// selectNode selects a node for name.
func (n *node) selectNode(name string) (string, error) {
	return n.circle.Get(name)
}

// isCurrentNode returns if address is current node or not.
func (n *node) isCurrentNode(address string) bool {
	return n.address == address
}

// updateCircle updates the consistent hash to new members.
func (n *node) updateCircle() {
	n.circle.Set(n.nodes())
}

// autoUpdateCircle starts a goroutine and runs updateCircle task at fixed duration.
func (n *node) autoUpdateCircle() {
	n.updateCircle()
	go func() {
		ticker := time.NewTicker(time.Duration(n.options.UpdateCircleDuration) * time.Second)
		for {
			select {
			case <-ticker.C:
				n.updateCircle()
			}
		}
	}()
}
