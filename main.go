package main

import (
	"os"
	"fmt"
	"time"
	"sync"
	"github.com/piotrpersona/goq/pkg/goq"
)

type cb[K, V any] struct {
	mu sync.Mutex

	arr []goq.Message[K, V]
}

func (c *cb[K, V]) Handle(msg goq.Message[K, V]) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.arr = append(c.arr, msg)
	return
}

func logErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	var topic goq.Topic = "topic-words"
	q := goq.New[int, string]()
	q.CreateTopic(topic)

	cbA := cb[int, string]{arr: make([]goq.Message[int, string], 0)}
	q.Subscribe(topic, "A", &cbA)

	cbB := cb[int, string]{arr: make([]goq.Message[int, string], 0)}
	q.Subscribe(topic, "B", &cbB)

	fmt.Println(q.Topics())
	fmt.Println(q.Subscribers(topic))

	q.Publish(topic, goq.Message[int, string]{1, "hello!"})
	q.Publish(topic, goq.Message[int, string]{3, "abc!"})

	time.Sleep(time.Second *1)

	fmt.Printf("A: %+v\n", cbA.arr)
	fmt.Printf("B: %+v\n", cbB.arr)

	q.Unsubscribe("A")

	q.Publish(topic, goq.Message[int, string]{4, "who's listening?"})

	time.Sleep(time.Second *1)

	fmt.Printf("A: %+v\n", cbA.arr)
	fmt.Printf("B: %+v\n", cbB.arr)
}