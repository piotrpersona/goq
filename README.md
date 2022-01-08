# goq

Golang interprocess, in-memory pub-sub message queue.

## Install

goq requires go1.18 or later.

```bash
go get -u github.com/piotrpersona/goq
```

## Examples

Checkout `examples/` directory.
Run
```console
go run examples/pubsub/main.go
```

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
---
**NOTE**

Publish, Subscribe, Unsubscribe and Close works asynchronously. Close can be awaited.