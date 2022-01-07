package main

import (
	"fmt"
	"time"
	"github.com/piotrpersona/goq/pkg/goq"
)

func main() {
	q := goq.New[int, string]()

	var topic goq.Topic = "topic-words"

	_ = q.CreateTopic(topic)

	channel, _ := q.Subscribe(topic, "example-subscriber")

	fmt.Println(q.Topics())
	fmt.Println(q.Subscribers(topic))

	go func() {
		select {
		case msg := <- channel:
			fmt.Printf("Received message from %s, key: %v value: %v\n", topic, msg.Key, msg.Value)
		}
	}()
	q.Publish(topic, goq.Message[int, string]{1, "hello!"})

	time.Sleep(time.Second *3)
}