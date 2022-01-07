package main

import (
	"os"
	"fmt"
	"time"
	"github.com/piotrpersona/goq/pkg/goq"
)

func consume[K, V any](channel <-chan goq.Message[K, V], topic goq.Topic, group goq.Group) {
	for msg := range channel {
		fmt.Printf("[%s] Received message from %s, key: %v value: %v\n", group, topic, msg.Key, msg.Value)
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

	var groupB goq.Group = "B"
	channelB, err := q.Subscribe(topic, groupB)
	logErr(err)

	fmt.Println(q.Topics())
	fmt.Println(q.Subscribers(topic))

	go consume(channelA, topic, groupA)
	go consume(channelB, topic, groupB)

	q.Publish(topic, goq.Message[int, string]{1, "hello!"})

	time.Sleep(time.Second *3)
}