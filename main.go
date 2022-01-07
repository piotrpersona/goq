package main

import (
	"fmt"
	"github.com/piotrpersona/goq/pkg/goq"
)

func main() {
	q := goq.New[int, string]()

	var topic goq.Topic = "topic-words"

	_ = q.CreateTopic(topic)

	q.Subscribe(topic, "example-subscriber")

	fmt.Println(q.Topics())
	fmt.Println(q.Subscribers(topic))
}