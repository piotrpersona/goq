package main

import (
	"fmt"
	"time"
	"sync"
	"runtime"
	"github.com/piotrpersona/goq"
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

func (c *cb[K, V]) Handle(msg goq.Message[K, V]){
	c.mu.Lock()
	defer c.mu.Unlock()

	c.arr = append(c.arr, msg)
}

func (c *cb[K, V]) read() []goq.Message[K, V] {
	return c.arr
}

func main() {
	var topic goq.Topic = "topic-words"
	fmt.Printf("Start main goroutine #go: %d\n", runtime.NumGoroutine())
	
	q := goq.New[int, string]()
	q.CreateTopic(topic)
	fmt.Printf("goq.New will spawn another goroutine #go: %d\n", runtime.NumGoroutine())

	cbA := newCb[int, string]()
	startedA, _ := q.Subscribe(topic, "A", cbA)
	<-startedA

	cbB := newCb[int, string]()
	startedB, _ := q.Subscribe(topic, "B", cbB)
	<-startedB

	fmt.Printf("2 subscribers were spawned #go: %d\n", runtime.NumGoroutine())

	fmt.Println(q.Topics())
	fmt.Println(q.Groups(topic))

	q.Publish(topic, goq.Message[int, string]{1, "hello!"})
	q.Publish(topic, goq.Message[int, string]{3, "abc!"})

	time.Sleep(time.Second)

	fmt.Printf("A: %+v\n", cbA.read())
	fmt.Printf("B: %+v\n", cbB.read())

	q.Unsubscribe("A")

	q.Publish(topic, goq.Message[int, string]{4, "who's listening?"})

	time.Sleep(time.Second)

	fmt.Printf("Group A unsubscribed #go: %d\n", runtime.NumGoroutine())

	fmt.Printf("A: %+v\n", cbA.read())
	fmt.Printf("B: %+v\n", cbB.read())

	<-q.Stop()

	fmt.Printf("Close queue #go: %d\n", runtime.NumGoroutine())
}