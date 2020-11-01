// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/06 01:07:08

package servers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

	"github.com/FishGoddess/kafo/caches"
	"github.com/FishGoddess/kafo/helpers"
	"github.com/julienschmidt/httprouter"
)

// HTTPServer is a http type server.
type HTTPServer struct {

	// node is an internal thing as a part of cluster.
	*node

	// cache is the real cache used inside.
	cache *caches.Cache

	// options stores all settings of server.
	options *Options
}

// NewHTTPServer returns a http server holder.
func NewHTTPServer(cache *caches.Cache, options *Options) (*HTTPServer, error) {

	n, err := newNode(options)
	if err != nil {
		return nil, err
	}

	return &HTTPServer{
		node: n,
		cache:   cache,
		options: options,
	}, nil
}

// Run runs the server and returns an error if something wrong.
func (hs *HTTPServer) Run() error {
	return http.ListenAndServe(helpers.JoinAddressAndPort(hs.options.Address, hs.options.Port), hs.routerHandler())
}

// =======================================================================

// wrapUriWithVersion wraps uri with api version.
// If version is "v1" and uri is "/cache", the result will be like "/v1/cache".
func wrapUriWithVersion(uri string) string {
	return path.Join("/", APIVersion, uri)
}

// routerHandler returns a Handler registering routers.
func (hs *HTTPServer) routerHandler() http.Handler {
	router := httprouter.New()
	router.GET(wrapUriWithVersion("/cache/:key"), hs.getHandler)
	router.PUT(wrapUriWithVersion("/cache/:key"), hs.setHandler)
	router.DELETE(wrapUriWithVersion("/cache/:key"), hs.deleteHandler)
	router.GET(wrapUriWithVersion("/status"), hs.statusHandler)
	router.GET(wrapUriWithVersion("/nodes"), hs.nodesHandler)
	return router
}

// getHandler is a handler for getting value of specified key.
func (hs *HTTPServer) getHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	key := params.ByName("key")
	node, err := hs.selectNode(key)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !hs.isCurrentNode(node) {
		writer.Header().Set("Location", node + request.RequestURI)
		writer.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	value, ok := hs.cache.Get(key)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.Write(value)
}

// setHandler is a handler for setting an entry of specified key and value.
func (hs *HTTPServer) setHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	key := params.ByName("key")
	node, err := hs.selectNode(key)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !hs.isCurrentNode(node) {
		writer.Header().Set("Location", node+request.RequestURI)
		writer.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	value, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	ttl, err := ttlOf(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = hs.cache.SetWithTTL(key, value, ttl)
	if err != nil {
		writer.WriteHeader(http.StatusRequestEntityTooLarge)
		writer.Write([]byte("Error: " + err.Error()))
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

// ttlOf returns ttl of this value in request and an error.
func ttlOf(request *http.Request) (int64, error) {
	ttls, ok := request.Header["Ttl"]
	if !ok || len(ttls) < 1 {
		return caches.NeverDie, nil
	}
	return strconv.ParseInt(ttls[0], 10, 64)
}

// deleteHandler is a handler for deleting the entry of specified key.
func (hs *HTTPServer) deleteHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	node, err := hs.selectNode(key)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !hs.isCurrentNode(node) {
		writer.Header().Set("Location", node+request.RequestURI)
		writer.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	err = hs.cache.Delete(key)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// statusHandler is handler for fetching the status of cache.
func (hs *HTTPServer) statusHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	status, err := json.Marshal(hs.cache.Status())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(status)
}

// nodesHandler is handler for fetching the nodes of cluster.
func (hs *HTTPServer) nodesHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	nodes, err := json.Marshal(hs.nodes())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(nodes)
}
