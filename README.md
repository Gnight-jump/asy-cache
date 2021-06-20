# asy-cache-client

#### 介绍
分布式缓存客户端demo


#### 软件架构
asy-cache客户端实现:

    1. 配置服务中心
    
    2. 定时拉取服务中心的节点图
    
    3. 根据节点图查询，无法连接，则顺延，同时主动拉取新的



#### 使用说明
```go
// 指定服务中心
client.CenterPath = "http://localhost:8000"
cli := client.New() // 新建客户端
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