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

type Queue[K, V any] struct {
	mu   sync.RWMutex

	topics map[Topic][]*subscriberGroup[K, V]
	subs map[Group]*subscriberGroup[K, V]
	
	subAdd chan *subscriberGroup[K, V]
	subDel chan *subscriberGroup[K, V]
}

func New[K, V any]() *Queue[K, V] {
	queue := &Queue[K, V]{
		topics: make(map[Topic][]*subscriberGroup[K, V]),
		subs: make(map[Group]*subscriberGroup[K, V]),
		subAdd: make(chan *subscriberGroup[K, V]),
		subDel: make(chan *subscriberGroup[K, V]),
	}

	go queue.listenSubscribe()

	return queue
}

func (q *Queue[K, V]) listenSubscribe() {
	for {
		select {
		case subGroup := <-q.subAdd:
			topic := subGroup.topic
			group := subGroup.group

			q.topics[topic] = append(q.topics[topic], subGroup)
			q.subs[group] = subGroup

			go func() {
				subGroup.run()
			}()

		case subGroup := <-q.subDel:
			subGroup.stop <- struct{}{}

			close(subGroup.stop)
			close(subGroup.channel)

			delete(q.subs, subGroup.group)
			q.deleteGroup(subGroup)
		default:
		}
	}
}

func (q *Queue[K, V]) deleteGroup(subGroup *subscriberGroup[K, V]) {
	var deletePos int
	for index, value := range q.topics[subGroup.topic] {
		if value.group == subGroup.group {
			deletePos = index
			break
		}
	}
	q.topics[subGroup.topic] = append(q.topics[subGroup.topic][:deletePos], q.topics[subGroup.topic][deletePos+1:]...)
}

