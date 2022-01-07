# goq

Golang interprocess, in-memory pub-sub message queue.

Checkout `examples/` directory.

> It requires go1.18 or later

## Usage

Create topic:
```go
q := goq.New[int, string]()
q.CreateTopic("topic")
```

Create callback:
```go
type callback struct {}

func (c callback) Handle(msg goq.Message[int, string]) (err error) {
    fmt.Printf("key: %d value: %s", msg.Key, msg.Value)
    return
}
```

Subscribe to a topic:
```go
q.Subscribe(topic, "example-subscriber", &callback{})
```

Publish to a topic:
```go
q.Publish("topic", goq.Message[int, string]{1, "Hello world!"})
```

Stop queue:
```go
<-q.Stop() // <-chan struct{}
```

> Publish, Subscribe, Unsubscribe and Close works asynchronously. Close can be awaited.