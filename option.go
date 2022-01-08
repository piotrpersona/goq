package goq


type Option[K, V any] interface {
	Apply(q *Queue[K, V])
}

type withTopicMaxSize[K, V any] struct{
	size int
}

func WithTopicMaxSize[K, V any](size int) Option[K, V] {
	return withTopicMaxSize[K, V]{size: size}
}

func (w withTopicMaxSize[K, V]) Apply(q *Queue[K, V]) {
	q.topicMaxSize = w.size
}
