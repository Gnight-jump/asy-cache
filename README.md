# asy-cache-client

#### 介绍
分布式缓存客户端demo


#### 软件架构
asy-cache客户端实现:

    1. 配置服务中心
    
    2. 定时拉取服务中心的节点图
    
    3. 根据节点图查询，失败返回提示，同时主动拉取新的



#### 使用说明
```go
// 指定服务中心
cli := client.New("http://localhost:8000") // 新建客户端
// 设置键值对
err := cli.Set("example", "x")
if err != nil {
    fmt.Println("can't set kv")
}
// 获取键值对
fmt.Println("get =", cli.Get("example"))
// 删除键值对
err = cli.Del("example")
if err != nil {
    fmt.Println("can't del kv")
}
fmt.Println("get =", cli.Get("example"))
```

#### 分析
对照组 1（设置一台asy-cache，一个服务中心）：

1. 1w次请求 -> 运行时间1050ms，丢失0
2. 10w次请求 -> 运行时间11680ms，丢失0
```go
var lost int64 = 0

func Test_ClientMap(t *testing.T) {
	cli := client.New("http://localhost:8000") // 新建客户端
	cli.Set("example", "result")
	start := time.Now()
	for i := 0; i < 100000; i++ { // 请求10w次
		get := cli.Get("example")
		if get == nil {
			atomic.AddInt64(&lost, 1)
		}
	}
	fmt.Println(time.Since(start))
	fmt.Println("lost", lost)
}
```

<hr/>

对照组 2（设置一台redis）：

1. 1w次请求 -> 运行时间580ms，丢失0
2. 10w次请求 -> 运行时间5892ms，丢失0
```go
var lost int64 = 0

func Test_mysql(t *testing.T) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = conn.Do("set", "example", "result")
	if err != nil {
		log.Println("[log] error")
		return
	}

	start := time.Now()
	for i := 0; i < 100000; i++ {
		_, err := redis.String(conn.Do("get", "example"))
		if err != nil {
			atomic.AddInt64(&lost, 1)
		}
	}
	fmt.Println(time.Since(start))
    fmt.Println("lost", lost)
}
```
