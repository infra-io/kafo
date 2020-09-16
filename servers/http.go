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
	"strconv"

	"github.com/FishGoddess/Lighter/caches"
	"github.com/julienschmidt/httprouter"
)

// HTTPServer is a http type server.
type HTTPServer struct {
	// Cache is the real cache used inside.
	cache *caches.Cache
}

// NewHTTPServer returns a http server holder.
func NewHTTPServer(cache *caches.Cache) *HTTPServer {
	return &HTTPServer{
		cache: cache,
	}
}

// Run runs the server at address and returns an error if something wrong.
func (hs *HTTPServer) Run(address string) error {
	return http.ListenAndServe(address, hs.routerHandler())
}

// =======================================================================

// wrapUriWithVersion wraps uri with api version.
// If version is "v1" and uri is "/cache", the result will be like "/v1/cache".
func wrapUriWithVersion(uri string) string {
	return "/" + APIVersion + uri
}

// routerHandler returns a Handler registering routers.
func (hs *HTTPServer) routerHandler() http.Handler {
	router := httprouter.New()
	router.GET(wrapUriWithVersion("/cache/:key"), hs.getHandler)
	router.PUT(wrapUriWithVersion("/cache/:key"), hs.setHandler)
	router.DELETE(wrapUriWithVersion("/cache/:key"), hs.deleteHandler)
	router.GET(wrapUriWithVersion("/status"), hs.statusHandler)
	return router
}

// getHandler is a handler for getting value of specified key.
func (hs *HTTPServer) getHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	key := params.ByName("key")
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
	hs.cache.SetWithTTL(key, value, ttl)
}

// ttlOf returns ttl of this value in request and an error.
func ttlOf(request *http.Request) (int64, error) {
	ttls, ok := request.Header["Ttl"]
	if !ok || len(ttls) < 1 {
		return 0, nil
	}
	return strconv.ParseInt(ttls[0], 10, 64)
}

// deleteHandler is a handler for deleting the entry of specified key.
func (hs *HTTPServer) deleteHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	hs.cache.Delete(key)
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
