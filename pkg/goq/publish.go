package goq

func (q *Queue[K, V]) Publish(topic Topic, msg Message[K, V]) (err error) {
	q.mu.Lock()
	q.mu.Unlock()

	subscribers := q.topics[topic]

	for _, subGroup := range subscribers {
		go func(subGroup *subscriberGroup[K, V]) {
			subGroup.channel <- msg
		}(subGroup)
	}
	return
}
