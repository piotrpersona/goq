package goq

import (
	"fmt"
)

func (q *Queue[K, V]) Topics() (topics []Topic) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for topic, _ := range q.topics {
		topics = append(topics, topic)
	}

	return
}

func (q *Queue[K, V]) CreateTopic(topic Topic) (err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if existsErr := q.topicExists(topic); existsErr == nil {
		err = fmt.Errorf("cannot create topic, err: topic %s already exists", topic)
		return
	}

	q.topics[topic] = make([]*subscriberGroup[K, V], 0)
	return
}

func (q Queue[K, V]) topicExists(topic Topic) (err error) {
	_, exists := q.topics[topic]
	if !exists {
		err = fmt.Errorf("topic %s does not exist", topic)
		return
	}
	return 
}
