package main

import (
	"asy-cache-client/client"
	"fmt"
)

func main() {
	// 指定服务中心
	client.CenterPath = "http://localhost:8000"
	cli := client.New()
	err := cli.Set("example", "x")
	if err != nil {
		fmt.Println("can't set kv")
	}
	fmt.Println("get =", cli.Get("example"))
	err = cli.Del("example")
	if err != nil {
		fmt.Println("can't del kv")
	}
	fmt.Println("get =", cli.Get("example"))
}
