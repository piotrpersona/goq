package queue

type Key[K] interface {
	Key() K
}

type Value[V] interface {
	Value() V
}

type Message[K, V] interface {
	Key[K]
	Value[V]
}

type Queue[M Message] interface {
	Publish(topic string, M) (err error)
	Subscribe(topic, group string) (<-chan M, err error)
	Ubsubscribe(group string) (err error)
}

