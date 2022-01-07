package queue

import (
	"sync"
	"fmt"
)

const (
	channelMaxSize = 128
)

type Message[K, V any] struct {
	Key K
	Value V
}

type Topic string
type Group string

type subscriberGroup[K, V any] struct {
	name Group
	channel chan Message[K, V]
}

type Queue[K, V any] struct {
	mu   sync.RWMutex
	subs map[Topic][]subscriberGroup[K, V]
}

func New[K, V any]() *Queue[K, V] {
	return &Queue[K, V]{
		subs: make(map[Topic][]subscriberGroup[K, V]),
	}
}

func (q *Queue[K, V]) CreateTopic(topic Topic) (err error) {
	if q.topicExists(topic) {
		err = fmt.Errorf("cannot create topic, err: topic %s already exists", topic)
		return
	}

	q.subs[topic] = make([]subscriberGroup[K, V], 0)
	return
}

func (q Queue[K, V]) topicExists(topic Topic) bool {
	_, exists := q.subs[topic]
	return exists
}

func (q *Queue[K, V]) Publish(topic Topic, msg Message[K, V]) (err error) {
	return
}

func (q *Queue[K, V]) Subscribe(topic Topic, group Group) (channel <-chan Message[K, V], err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.topicExists(topic) {
		err = fmt.Errorf("cannot subscribe to a topic, err: topic %s does not exist", topic)
		return
	}

	subscriberGroup := subscriberGroup[K, V]{
		name: group,
		channel: make(chan Message[K, V], channelMaxSize),
	}

	q.subs[topic] = append(q.subs[topic], subscriberGroup)
	channel = subscriberGroup.channel
	return
}

func (q *Queue[K, V]) Ubsubscribe(group Group) (err error) {
	return
}
