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
type Callback[K, V any] interface {
	Handle(msg Message[K, V]) (err error)
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
	maxRetries int
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
			s.handle(msg)
		}
	}
	return
}

func (s *subscriberGroup[K, V]) handle(msg Message[K, V]) (err error) {
	handleErr := s.cb.Handle(msg)
	if handleErr == nil {
		return
	}

	retryNumber := 1

	for retryNumber <= s.maxRetries {
		handleErr = s.cb.Handle(msg)
		if handleErr == nil {
			return
		}

		retryNumber++
	}

	if handleErr != nil {
		err = fmt.Errorf("subscriber %s error while handling message from %s, err: %w", s.group, s.topic, handleErr)
		return
	}
	return
}

// Subscribe will spawn new subscriber goroutine and run cb Callback[K, V].
// It works asynchronously.
func (q *Queue[K, V]) Subscribe(topic Topic, group Group, cb Callback[K, V], opts ...subscribeOption[K, V]) (err error) {
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
		channel: make(chan Message[K, V], q.topicMaxSize),
		cb: cb,
		stop: make(chan struct{}),
		maxRetries: DefaultCallbackHandleMaxRetries,
	}

	for _, option := range opts {
		option.apply(subGroup)
	}

	q.subAdd <- subGroup

	return
}

// Unsubscribe will unsubscribe group.
// It works asynchronously.
func (q *Queue[K, V]) Unsubscribe(group Group) (err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

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
	q.mu.Unlock()

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
		case <- q.closeChan:
			q.doneChan <- struct{}{}
			return
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

	for group, _ := range q.subs {
		go func() {
			defer wg.Done()
			q.Unsubscribe(group)
		}()
	}

	wg.Wait()

	q.closeChan <- struct{}{}
	return q.doneChan
}
