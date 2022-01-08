package goq

import (
	"fmt"
	"sync"
)

// Topic represents topic name that can be subscribed.
type Topic string

// Group represents subscriber group name that can subscribe to a topic.
type Group string


// Message that is published to and consumed from a topic.
type Message[K, V any] struct {
	Key K
	Value V
}

// Callback to be run when new message is consumed from a topic.
// All errors should be handled inside the Callback Handle method.
type Callback[K, V any] interface {
	Handle(msg Message[K, V])
}

// Queue represents message queue that can spawn new subscibers, accepts messages from publishers,
// broadcast messages to subscribers, unsubscribe subscribers.
type Queue[K, V any] struct {
	mu   sync.RWMutex

	topics map[Topic][]*subscriberGroup[K, V]
	subs map[Group]*subscriberGroup[K, V]

	topicMaxSize int
	
	subAdd chan *subscriberGroup[K, V]
	subDel chan *subscriberGroup[K, V]

	closeChan chan struct{}
	doneChan chan struct{}
}

func New[K, V any](opts ...queueOption[K, V]) *Queue[K, V] {
	queue := &Queue[K, V]{
		topics: make(map[Topic][]*subscriberGroup[K, V]),
		subs: make(map[Group]*subscriberGroup[K, V]),
		topicMaxSize: DefaultTopicMaxSize,
		subAdd: make(chan *subscriberGroup[K, V]),
		subDel: make(chan *subscriberGroup[K, V]),
		closeChan: make(chan struct{}),
		doneChan: make(chan struct{}),
	}

	for _, option := range opts {
		option.apply(queue)
	}

	go queue.listenSubscribe()

	return queue
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
			s.cb.Handle(msg)
		}
	}
	return
}


// Subscribe will spawn new subscriber goroutine and run cb Callback[K, V].
// It works asynchronously.
func (q *Queue[K, V]) Subscribe(topic Topic, group Group, cb Callback[K, V]) (err error) {
	q.mu.Lock()
	

	if err = q.topicExists(topic); err != nil {
		err = fmt.Errorf("cannot subscribe to a topic, err: %w", err)
		return
	}

	if _, exists := q.subs[group]; exists {
		err = fmt.Errorf("cannot subscribe to a topic, err: group %s already exists", group)
		return
	}

	q.mu.Unlock()

	subGroup := &subscriberGroup[K, V]{
		group: group,
		topic: topic,
		channel: make(chan Message[K, V], q.topicMaxSize),
		cb: cb,
		stop: make(chan struct{}),
	}

	q.subAdd <- subGroup

	return
}

// Unsubscribe will unsubscribe group.
// It works asynchronously.
func (q *Queue[K, V]) Unsubscribe(group Group) (err error) {
	q.mu.Lock()
	q.mu.Unlock()

	subGroup, groupExists := q.subs[group]
	if !groupExists {
		err = fmt.Errorf("cannot unsubscribe, err: group %s does not exist", group)
		return
	}

	q.subDel <- subGroup

	return
}

// Publish will send Message[K, V] to the internal queue topic channel.
func (q *Queue[K, V]) Publish(topic Topic, msg Message[K, V]) (err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	subscribers := q.topics[topic]

	for _, subGroup := range subscribers {
		go func(subGroup *subscriberGroup[K, V]) {
			select {
			case <- subGroup.stop:
				return
			default:
			}

			select {
			case <- subGroup.stop:
				return
			case subGroup.channel <- msg:
			}
		}(subGroup)
	}
	return
}

func (q *Queue[K, V]) listenSubscribe() {
	for {
		select {
		case subGroup := <-q.subAdd:
			q.subscribe(subGroup)

			go func() {
				subGroup.run()
			}()

		case subGroup := <-q.subDel:
			subGroup.stop <- struct{}{}

			close(subGroup.stop)
			close(subGroup.channel)

			q.deleteGroup(subGroup)
		case <- q.closeChan:
			q.doneChan <- struct{}{}
			return
		default:
		}
	}
}

func (q *Queue[K, V]) subscribe(subGroup *subscriberGroup[K, V]) {
	topic := subGroup.topic
	group := subGroup.group

	q.mu.Lock()
	q.topics[topic] = append(q.topics[topic], subGroup)
	q.subs[group] = subGroup
	q.mu.Unlock()
}

func (q *Queue[K, V]) deleteGroup(subGroup *subscriberGroup[K, V]) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var deletePos int
	for index, value := range q.topics[subGroup.topic] {
		if value.group == subGroup.group {
			deletePos = index
			break
		}
	}
	q.topics[subGroup.topic] = append(q.topics[subGroup.topic][:deletePos], q.topics[subGroup.topic][deletePos+1:]...)
	delete(q.subs, subGroup.group)
}

// Groups will return list of groups for a topic.
func (q *Queue[K, V]) Groups(topic Topic) (subs []Group, err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if err = q.topicExists(topic); err != nil {
		err = fmt.Errorf("cannot get subscribers list, err: %w", err)
		return
	}
	
	for _, sub := range q.topics[topic] {
		subs = append(subs, sub.group)
	}

	return
}

// Stop will terminate all subscriber processes.
// It works asynchronously and returns channel to await for queue to stop.
func (q *Queue[K, V]) Stop() <-chan struct{} {
	var wg sync.WaitGroup
	wg.Add(len(q.subs))

	q.mu.Lock()
	subscribers := q.subs

	for group, _ := range subscribers {
		go func() {
			defer wg.Done()
			q.Unsubscribe(group)
		}()
	}
	q.mu.Unlock()

	wg.Wait()

	q.closeChan <- struct{}{}
	return q.doneChan
}
