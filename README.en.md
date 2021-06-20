# asy-cache-client

#### Introduce
Example asy-cache client demo

#### Software architecture
Example asy-cache client:

    1. Configure the service center
    
    2. Fetch the node diagram of the service center regularly
    
    3. If it cannot be connected according to the node diagram, it will be postponed and take the initiative to fetch a new one


#### Direction for use
```go
// Designated Service Center
client.CenterPath = "http://localhost:8000"
cli := client.New() // New Client
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