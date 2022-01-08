package goq


type QueueOption[K, V any] interface {
	apply(q *Queue[K, V])
}

type withTopicMaxSize[K, V any] struct{
	size int
}

func WithTopicMaxSize[K, V any](size int) QueueOption[K, V] {
	return withTopicMaxSize[K, V]{size: size}
}

func (w withTopicMaxSize[K, V]) apply(q *Queue[K, V]) {
	q.topicMaxSize = w.size
}

type SubscribeOption[K, V any] interface {
	apply(s *subscriberGroup[K, V])
}

type withRetry[K, V any] struct{
	retries int
}

func WithRetry[K, V any](retries int) SubscribeOption[K, V] {
	return withRetry[K, V]{retries: retries}
}

func (o withRetry[K, V]) apply(s *subscriberGroup[K, V]) {
	s.maxRetries = o.retries
}
