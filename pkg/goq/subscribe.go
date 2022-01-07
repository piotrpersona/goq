package goq

import (
	"fmt"
)

type Callback[K, V any] interface {
	Handle(msg Message[K, V]) (err error)
}

type subscriberGroup[K, V any] struct {
	group Group
	topic Topic
	channel chan Message[K, V] // channel to receive messages
	stop chan struct{}
}

type Subscriber[K, V any] struct {
	subGroup *subscriberGroup[K, V]
}

func (s *Subscriber[K, V]) Run(cb Callback[K, V]) (err error) {
	for {
		select {
		case <- s.subGroup.stop:
			return
		default:
		}

		select {
		case <- s.subGroup.stop:
			return
		case msg := <- s.subGroup.channel:
			handleErr := cb.Handle(msg)
			if handleErr != nil {
				err = fmt.Errorf("subscriber %s error while handling message from %s, err: %w", s.subGroup.group, s.subGroup.topic, handleErr)
				return
			}
		}
	}
	return
}

func (q *Queue[K, V]) Subscribe(topic Topic, group Group) (subscriber Subscriber[K, V], err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if err = q.topicExists(topic); err != nil {
		err = fmt.Errorf("cannot subscribe to a topic, err: %w", err)
		return
	}

	if _, exists := q.subs[group]; exists {
		err = fmt.Errorf("cannot subscribe to a topic, err: group %s already exists", group)
		return
	}

	subGroup := &subscriberGroup[K, V]{
		group: group,
		topic: topic,
		channel: make(chan Message[K, V], channelMaxSize),
		stop: make(chan struct{}),
	}

	q.subAdd <- subGroup

	subscriber = Subscriber[K, V]{
		subGroup: subGroup,
	}
	return
}


func (q *Queue[K, V]) Unsubscribe(sub Subscriber[K, V]) (err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	_, groupExists := q.subs[sub.subGroup.group]
	if !groupExists {
		err = fmt.Errorf("cannot unsubscribe, err: group %s does not exist")
		return
	}

	q.subDel <- sub.subGroup

	return
}
