package worker

import "bullgo/queue"

type Limiter struct {
	Max      int
	Duration int
}

type Worker struct {
	Name        string
	Connection  queue.RedisConnection
	Concurrency int
	Limiter
}

type WorkerOption func(*Worker)

func CreateWorker(opts ...WorkerOption) *Worker {
	w := &Worker{
		Connection: queue.RedisConnection{Host: "127.0.0.1", Port: 6379, MaxRetriesPerRequest: nil},
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}
