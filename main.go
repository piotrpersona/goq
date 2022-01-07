package main

import (
	"os"
	"fmt"
	"time"
	"sync"
	"github.com/piotrpersona/goq/pkg/goq"
)

type consumer[K, V any] struct {
	mu sync.Mutex

	arr []goq.Message[K, V]
}

func (c *consumer[K, V]) consume(channel <-chan goq.Message[K, V], topic goq.Topic, group goq.Group) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for msg := range channel {
		c.arr = append(c.arr, msg)
	}
}

func logErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	q := goq.New[int, string]()

	var topic goq.Topic = "topic-words"

	err := q.CreateTopic(topic)
	logErr(err)

	var groupA goq.Group = "A"
	channelA, err := q.Subscribe(topic, groupA)
	logErr(err)
	consumerA := consumer[int, string]{arr: make([]goq.Message[int, string], 0)}

	var groupB goq.Group = "B"
	channelB, err := q.Subscribe(topic, groupB)
	logErr(err)
	consumerB := consumer[int, string]{arr: make([]goq.Message[int, string], 0)}

	fmt.Println(q.Topics())
	fmt.Println(q.Subscribers(topic))

	go consumerA.consume(channelA, topic, groupA)
	go consumerB.consume(channelB, topic, groupB)

	q.Publish(topic, goq.Message[int, string]{1, "hello!"})
	q.Publish(topic, goq.Message[int, string]{3, "abc!"})

	time.Sleep(time.Second *1)

	fmt.Printf("A: %+v\n", consumerA.arr)
	fmt.Printf("B: %+v\n", consumerB.arr)
}