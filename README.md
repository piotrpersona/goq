# goq

Golang interprocess, in-memory message queue based solution.
It implements publish-subscribe pattern.

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
q.Stop()
```

> Publish, Subscribe, Unsubscribe and Stop works asynchronously, and they should 