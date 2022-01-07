package main

import (
	"fmt"
	"time"
	"github.com/piotrpersona/goq/pkg/goq"
)

func consume[K, V any](channel <-chan goq.Message[K, V], topic goq.Topic, group goq.Group) {
	select {
	case msg := <- channel:
		fmt.Printf("[%s] Received message from %s, key: %v value: %v\n", group, topic, msg.Key, msg.Value)
	}
}

func main() {
	q := goq.New[int, string]()

	var topic goq.Topic = "topic-words"

	_ = q.CreateTopic(topic)

	var groupA goq.Group = "A"
	channelA, _ := q.Subscribe(topic, groupA)
	var groupB goq.Group = "B"
	channelB, _ := q.Subscribe(topic, groupB)

	fmt.Println(q.Topics())
	fmt.Println(q.Subscribers(topic))

	go consume(channelA, topic, groupA)
	go consume(channelB, topic, groupB)

	q.Publish(topic, goq.Message[int, string]{1, "hello!"})

	time.Sleep(time.Second *3)
}