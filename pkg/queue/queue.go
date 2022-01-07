package queue

import (
	"sync"
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

func (q *Queue[K, V]) Publish(topic Topic, msg Message[K, V]) (err error) {
	return
}

func (q *Queue[K, V]) Subscribe(topic Topic, group Group) (channel <-chan Message[K, V], err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

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
