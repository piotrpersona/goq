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
	cb Callback[K, V]
	stop chan struct{}
}

func (s *subscriberGroup[K, V]) run() (err error) {
	for {
		select {
		case <- s.stop:
			return
		default:
		}

		select {
		case <- s.stop:
			return
		case msg := <- s.channel:
			handleErr := s.cb.Handle(msg)
			if handleErr != nil {
				err = fmt.Errorf("subscriber %s error while handling message from %s, err: %w", s.group, s.topic, handleErr)
				return
			}
		}
	}
	return
}

func (q *Queue[K, V]) Subscribe(topic Topic, group Group, cb Callback[K, V]) (err error) {
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
		cb: cb,
		stop: make(chan struct{}),
	}

	q.subAdd <- subGroup

	return
}


func (q *Queue[K, V]) Unsubscribe(group Group) (err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	subGroup, groupExists := q.subs[group]
	if !groupExists {
		err = fmt.Errorf("cannot unsubscribe, err: group %s does not exist")
		return
	}

	q.subDel <- subGroup

	return
}
