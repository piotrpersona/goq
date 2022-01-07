package goq

func (q *Queue[K, V]) Publish(topic Topic, msg Message[K, V]) (err error) {
	subscribers := q.topics[topic]

	for _, subGroup := range subscribers {
		go func() {
			subGroup.channel <- msg
		}()
	}
	return
}
