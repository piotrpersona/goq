package main

import (
	"github.com/piotrpersona/goq/pkg/queue"
)

func main() {
	q := queue.Queue[int, string]{}
	q.Ubsubscribe("abc")
}