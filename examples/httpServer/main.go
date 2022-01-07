package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/piotrpersona/goq"
)

func writeHandler(q *goq.Queue[time.Time, http.Request], topic goq.Topic) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		q.Publish(topic, goq.Message[time.Time, http.Request]{Key: time.Now(), Value: *r})

		fmt.Fprintf(rw, "request saved in %s", time.Now().Sub(start))
	}
}

type requestsCallback struct {}

func (r *requestsCallback) Handle(msg goq.Message[time.Time, http.Request]) (err error) {
	time.Sleep(time.Second * 2) // Long message processing
	fmt.Printf("%s %+v\n", msg.Key.UTC().Format(time.RFC3339), msg.Value)
	return
}

func main() {
	q := goq.New[time.Time, http.Request]()
	q.CreateTopic("requests")
	q.Subscribe("requests", "requests-processor", &requestsCallback{})

	http.HandleFunc("/", writeHandler(q, "requests"))
	http.ListenAndServe(":8080", nil)
}