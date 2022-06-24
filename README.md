# PubSub

This library provides a generic implementation of pub-sub pattern using channels in Go. An example usage would be like:

```go
ps := pubsub.New[string, int]()

// Produces a random integer each second and publishes
// it on the "random" topic.
go func() {
    for {
        ps.Publish("random", rand.Int())
        time.Sleep(time.Second)
    }
}()

// Consumes produced random integers by subscribing on
// the "random" topic.
go func() {
    ch := make(chan int)
    ps.Subscribe("random", ch)
    for r := range ch {
        fmt.Println(r)
    }
}()
```
