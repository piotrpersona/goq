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
