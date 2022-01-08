package main

import (
	"fmt"
	"time"
	"github.com/piotrpersona/goq"
)

type Logger struct {
	topic goq.Topic
	q *goq.Queue[time.Time, string]
}

func NewLogger(q *goq.Queue[time.Time, string]) (logger *Logger) {
	var topic goq.Topic = "logs"	
	q.CreateTopic(topic)

	logger = &Logger{
		topic: topic,
		q: q,
	}

	q.Subscribe(topic, "logger", logger)

	return
}

func (l *Logger) Info(message string) {
	msg := goq.Message[time.Time, string]{Key: time.Now(), Value: message}
	l.q.Publish(l.topic, msg)
}

func (l *Logger) Handle(msg goq.Message[time.Time, string]) {
	fmt.Printf("%s %s\n", msg.Key.UTC().Format(time.RFC3339), msg.Value)
}

func main() {
	q := goq.New[time.Time, string]()
	logger := NewLogger(q)

	logger.Info("hello!")
	
	time.Sleep(time.Millisecond*1)
}