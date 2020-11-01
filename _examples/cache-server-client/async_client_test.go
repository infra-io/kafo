// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/20 23:13:45

package client

import (
	"strconv"
	"testing"
	"time"
)

const (
	// keySize is the key size of test.
	keySize = 10000
)

// testTask is a wrapper wraps task to testTask.
func testTask(task func(no int)) string {
	beginTime := time.Now()
	for i := 0; i < keySize; i++ {
		task(i)
	}
	return time.Now().Sub(beginTime).String()
}

// go test -v -cover -run=^TestNewAsyncClient$
func TestNewAsyncClient(t *testing.T) {

	client, err := NewAsyncClient(":5837")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	for i := 0; i < 100; i++ {
		data := strconv.Itoa(i)
		response := <-client.Set(data, []byte(data), 0)
		if response.Err != nil {
			t.Fatal(response.Err)
		}
	}

	for i := 0; i < 100; i++ {
		data := strconv.Itoa(i)
		response := <-client.Get(data)
		if response.Err != nil {
			t.Fatal(response.Err)
		}

		if string(response.Body) != data {
			t.Fatalf("Get key %s returns wrong value %s!", data, string(response.Body))
		}
	}

	status, err := (<-client.Status()).ToStatus()
	if err != nil {
		t.Fatal(err)
	}

	if status.Count != 100 {
		t.Fatalf("Count in status %d is wrong!", status.Count)
	}

	if status.KeySize != status.ValueSize || status.KeySize != 190 {
		t.Fatalf("KeySize %d or valueSize %d is wrong!", status.KeySize, status.ValueSize)
	}

	for i := 0; i < 100; i++ {
		data := strconv.Itoa(i)
		response := <-client.Delete(data)
		if response.Err != nil {
			t.Fatal(response.Err)
		}
	}

	status, err = (<-client.Status()).ToStatus()
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

// go test -v -count=1 -run=^TestAsyncClientPerformance$
func TestAsyncClientPerformance(t *testing.T) {

	client, err := NewAsyncClient(":5837")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	writeTime := testTask(func(no int) {
		data := strconv.Itoa(no)
		client.Set(data, []byte(data), 0)
	})

	t.Logf("写入消耗时间为 %s！", writeTime)

	time.Sleep(3 * time.Second)

	readTime := testTask(func(no int) {
		data := strconv.Itoa(no)
		client.Get(data)
	})

	t.Logf("读取消耗时间为 %s！", readTime)

	time.Sleep(time.Second)
}
