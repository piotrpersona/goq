package goq

import (
	"fmt"
)

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

	if _, exists := q.subs[group]; exists {
		err = fmt.Errorf("cannot subscribe to a topic, err: group %s already exists", group)
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
