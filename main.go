package main

import (
	"fmt"
	"time"
	"sync"
	"github.com/piotrpersona/goq/pkg/goq"
)

type cb[K, V any] struct {
	mu sync.Mutex

	arr []goq.Message[K, V]
}

func newCb[K, V any]() *cb[K, V] {
	return &cb[K, V]{
		arr: make([]goq.Message[K, V], 0),
	}
}

func (c *cb[K, V]) Handle(msg goq.Message[K, V]) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.arr = append(c.arr, msg)
	return
}


func main() {
	var topic goq.Topic = "topic-words"
	q := goq.New[int, string]()
	q.CreateTopic(topic)

	cbA := newCb[int, string]()
	q.Subscribe(topic, "A", cbA)

	cbB := newCb[int, string]()
	q.Subscribe(topic, "B", cbB)

	fmt.Println(q.Topics())
	fmt.Println(q.Groups(topic))

	q.Publish(topic, goq.Message[int, string]{1, "hello!"})
	q.Publish(topic, goq.Message[int, string]{3, "abc!"})

	time.Sleep(time.Second)

	fmt.Printf("A: %+v\n", cbA.arr)
	fmt.Printf("B: %+v\n", cbB.arr)

	q.Unsubscribe("A")

	q.Publish(topic, goq.Message[int, string]{4, "who's listening?"})

	time.Sleep(time.Second)

	fmt.Printf("A: %+v\n", cbA.arr)
	fmt.Printf("B: %+v\n", cbB.arr)
}