package goq

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

	if err = q.topicExists(topic); err == nil {
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

func (q *Queue[K, V]) Publish(topic Topic, msg Message[K, V]) (err error) {
	return
}

func (q *Queue[K, V]) Subscribers(topic Topic) (subs []Group, err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if err = q.topicExists(topic); err != nil {
		err = fmt.Errorf("cannot get subscribers list, err: %w", err)
		return
	}
	
	for _, sub := range q.topics[topic] {
		subs = append(subs, sub.name)
	}

	return
}

func (q *Queue[K, V]) Subscribe(topic Topic, group Group) (channel <-chan Message[K, V], err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if err = q.topicExists(topic); err != nil {
		err = fmt.Errorf("cannot subscribe to a topic, err: %w", err)
		return
	}

	subscriberGroup := &subscriberGroup[K, V]{
		name: group,
		topic: topic,
		channel: make(chan Message[K, V], channelMaxSize),
	}

	q.topics[topic] = append(q.topics[topic], subscriberGroup)
	q.subs[group] = subscriberGroup

	channel = subscriberGroup.channel
	return
}

func (q *Queue[K, V]) Ubsubscribe(group Group) (err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	subGroup, groupExists := q.subs[group]
	if !groupExists {
		err = fmt.Errorf("cannot unsubscribe, err: group %s does not exist")
		return
	}

	delete(q.subs, group)
	q.deleteGroup(subGroup)

	return
}

func (q *Queue[K, V]) deleteGroup(group *subscriberGroup[K, V]) {
	var deletePos int
	for index, value := range q.topics[group.topic] {
		if value.name == group.name {
			deletePos = index
			break
		}
	}
	q.topics[group.topic] = append(q.topics[group.topic][:deletePos], q.topics[group.topic][deletePos+1:]...)
}