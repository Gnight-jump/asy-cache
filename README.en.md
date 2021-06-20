# asy-cache-client

#### Introduce
Example asy-cache client demo

#### Software architecture
Example asy-cache client:

    1. Configure the service center
    
    2. Fetch the node diagram of the service center regularly
    
    3. According to the node graph query, failure to return a prompt, and actively pull the new


#### Direction for use
```go
// Designated Service Center
cli := client.New("http://localhost:8000") // New Client
// Set key value pairs
err := cli.Set("example", "x")
if err != nil {
    fmt.Println("can't set kv")
}
// Gets the key-value pair
fmt.Println("get =", cli.Get("example"))
// Delete key-value pairs
err = cli.Del("example")
if err != nil {
    fmt.Println("can't del kv")
}
fmt.Println("get =", cli.Get("example"))
```

#### 
Control group 1 (set an ASY-Cache, a service center) :

1. 1W requests -> running time 1050ms, 0 lost
2. 10W requests -> running time 11680ms, 0 lost
```go
var lost int64 = 0

func Test_ClientMap(t *testing.T) {
	cli := client.New("http://localhost:8000") // new Cline
	cli.Set("example", "result")
	start := time.Now()
	for i := 0; i < 100000; i++ { 
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

Control group 2 (set one Redis) :

1. 1W requests -> running time 580ms, 0 lost
2. 10W requests -> run time 5892ms, 0 lost
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