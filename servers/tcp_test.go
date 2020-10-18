// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/17 00:18:48

package servers

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/FishGoddess/kafo/caches"
)

// go test -v -cover -run=^TestTCPServer$
func TestTCPServer(t *testing.T) {

	server := NewTCPServer(caches.NewCache())
	go func() {
		err := server.Run(":5837")
		if err != nil {
			t.Fatal(err)
		}
	}()
	defer server.Close()

	time.Sleep(time.Second)

	client, err := NewTCPClient(":5837")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	t.Log("Start setting...")
	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		func(data string) {
			defer wg.Done()
			err := client.Set(data, []byte(data), caches.NeverDie)
			if err != nil {
				t.Fatal(err)
			}
		}(strconv.Itoa(i))
	}
	wg.Wait()

	t.Log("Start getting value...")
	for i := 0; i < 100; i++ {
		wg.Add(1)
		func(data string) {
			defer wg.Done()
			value, err := client.Get(data)
			if err != nil {
				t.Fatal(err)
			}

			if string(value) != data {
				t.Fatalf("Get key %s returns wrong value %s!", data, string(value))
			}
		}(strconv.Itoa(i))
	}
	wg.Wait()

	t.Log("Start getting status...")
	status, err := client.Status()
	if err != nil {
		t.Fatal(err)
	}

	if status.Count != 100 {
		t.Fatalf("Count in status %d is wrong!", status.Count)
	}

	if status.KeySize != status.ValueSize || status.KeySize != 190 {
		t.Fatalf("KeySize %d or valueSize %d is wrong!", status.KeySize, status.ValueSize)
	}

	t.Log("Start deleting...")
	for i := 0; i < 100; i++ {
		wg.Add(1)
		func(data string) {
			defer wg.Done()
			err := client.Delete(data)
			if err != nil {
				t.Fatal(err)
			}
		}(strconv.Itoa(i))
	}
	wg.Wait()

	t.Log("Start checking status...")
	status, err = client.Status()
	if err != nil {
		t.Fatal(err)
	}

	if status.Count != 0 {
		t.Fatalf("Count in status %d is wrong!", status.Count)
	}

	if status.KeySize != status.ValueSize || status.KeySize != 0 {
		t.Fatalf("KeySize %d or valueSize %d is wrong!", status.KeySize, status.ValueSize)
	}
}
