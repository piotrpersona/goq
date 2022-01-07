package main

import (
	"fmt"
	"github.com/piotrpersona/goq/pkg/queue"
)

func main() {
	q := queue.New[int, string]()
	channel, _ := q.Subscribe("topic-words", "example-subscriber")

	fmt.Println(len(channel))
}