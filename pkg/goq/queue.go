package goq

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
	topic Topic
	channel chan Message[K, V]
}

type Queue[K, V any] struct {
	mu   sync.RWMutex
	topics map[Topic][]*subscriberGroup[K, V]
	subs map[Group]*subscriberGroup[K, V]
}

func New[K, V any]() *Queue[K, V] {
	return &Queue[K, V]{
		topics: make(map[Topic][]*subscriberGroup[K, V]),
		subs: make(map[Group]*subscriberGroup[K, V]),
	}
}
