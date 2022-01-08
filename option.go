package goq


type queueOption[K, V any] interface {
	apply(q *Queue[K, V])
}

type withTopicMaxSize[K, V any] struct{
	size int
}

// WithRetry sets the limit of messages that can be stored in a topic.
func WithTopicMaxSize[K, V any](size int) queueOption[K, V] {
	return withTopicMaxSize[K, V]{size: size}
}

func (w withTopicMaxSize[K, V]) apply(q *Queue[K, V]) {
	q.topicMaxSize = w.size
}

type subscribeOption[K, V any] interface {
	apply(s *subscriberGroup[K, V])
}

type withRetry[K, V any] struct{
	retries int
}

// WithRetry sets number of retries for Callback.
func WithRetry[K, V any](retries int) subscribeOption[K, V] {
	return withRetry[K, V]{retries: retries}
}

func (o withRetry[K, V]) apply(s *subscriberGroup[K, V]) {
	s.maxRetries = o.retries
}
