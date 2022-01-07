package main

import (
	"fmt"
	"github.com/piotrpersona/goq/pkg/queue"
)

func main() {
	q := queue.New[int, string]()

	var topic queue.Topic = "topic-words"

	_ = q.CreateTopic(topic)

	q.Subscribe(topic, "example-subscriber")

	fmt.Println(q.Topics())
	fmt.Println(q.Subscribers(topic))
}