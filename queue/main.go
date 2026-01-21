package queue

import (
	"bullgo/jobs"
)

type RedisConnection struct {
	Host                 string
	Port                 int
	MaxRetriesPerRequest *int
}

type Queue struct {
	Name       string
	Connection RedisConnection
}

type QueueOption func(*Queue)

func CreateQueue(opts ...QueueOption) *Queue {
	q := &Queue{
		Connection: RedisConnection{Host: "127.0.0.1", Port: 6379, MaxRetriesPerRequest: nil},
	}

	for _, opt := range opts {
		opt(q)
	}

	return q
}

func NameLayer(name string) QueueOption {
	return func(q *Queue) {
		q.Name = name
	}
}

func RedisLayer(conn RedisConnection) QueueOption {
	return func(q *Queue) {
		q.Connection = RedisConnection{Host: conn.Host, Port: conn.Port, MaxRetriesPerRequest: conn.MaxRetriesPerRequest}
	}
}

func (q *Queue) Add(job jobs.Job) error {
	return nil
}
