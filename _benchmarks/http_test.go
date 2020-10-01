// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/01 16:31:41

package _benchmarks

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// go test http_test.go -bench=^BenchmarkKafoHttpSet$ -benchtime=1s
func BenchmarkKafoHttpSet(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		data := strconv.Itoa(i)
		request, err := http.NewRequest("PUT", "http://localhost:5837/v1/cache/"+data, strings.NewReader(data))
		if err != nil {
			b.Fatal(err)
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			b.Fatal(err)
		}
		response.Body.Close()
	}
}

// go test http_test.go -bench=^BenchmarkKafoHttpGet$ -benchtime=1s
func BenchmarkKafoHttpGet(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		data := strconv.Itoa(i)
		request, err := http.NewRequest("GET", "http://localhost:5837/v1/cache/"+data, nil)
		if err != nil {
			b.Fatal(err)
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			b.Fatal(err)
		}
		response.Body.Close()
	}
}
