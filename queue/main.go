package queue

import (
	queue_domain "bullgo/internal/domain/queue"
	redis_domain "bullgo/internal/domain/redis"
)

var Queues []queue_domain.Queue

func CreateQueue(opts ...queue_domain.QueueOption) *queue_domain.Queue {
	q := &queue_domain.Queue{
		Connection: redis_domain.RedisConnection{Host: "127.0.0.1", Port: 6379, Password: nil},
	}

	for _, opt := range opts {
		opt(q)
	}

	queue_domain.InitClient(q)
	Queues = append(Queues, *q)
	return q
}

func NameLayer(name string) queue_domain.QueueOption {
	return func(q *queue_domain.Queue) {
		q.Name = name
	}
}

func RedisLayer(conn redis_domain.RedisConnection) queue_domain.QueueOption {
	return func(q *queue_domain.Queue) {
		q.Connection = redis_domain.RedisConnection{Host: conn.Host, Port: conn.Port, Password: conn.Password}
	}
}
