package main

import (
	"asy-cache-client/client"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

var lost int64 = 0

func Test_ClientMap(t *testing.T) {
	cli := client.New("http://localhost:8000") // 新建客户端
	cli.Set("example", "result")
	start := time.Now()
	for i := 0; i < 100000; i++ { // 请求1w次
		get := cli.Get("example")
		if get == nil {
			atomic.AddInt64(&lost, 1)
		}
	}
	fmt.Println(time.Since(start)) //
	fmt.Println("lost", lost)
}
