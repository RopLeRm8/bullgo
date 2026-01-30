package queue_domain

import (
	redis_domain "bullgo/internal/domain/redis"
	"bullgo/jobs"
	"bullgo/storage"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Queue struct {
	Name       string
	Connection redis_domain.RedisConnection
	Jobs       []*jobs.Job
	Client     *redis.Client
}

type QueueOption func(q *Queue)

// Adding jobs to a queue
func (q *Queue) Add(job *jobs.Job) error {
	q.Jobs = append(q.Jobs, job)
	payload, _ := json.Marshal(job)
	storage.StoreJobData(q.Client, q.Name, job.ID, payload)
	storage.PushToActive(q.Client, q.Name, job.ID)
	return nil
}

// Remove job from queue
func (q *Queue) Remove(job jobs.Job) error {
	jobInd := -1
	for i, j := range q.Jobs {
		if j.ID == job.ID {
			jobInd = i
			break
		}
	}

	if jobInd != -1 {
		q.Jobs = append(q.Jobs[:jobInd], q.Jobs[jobInd+1:]...)
	}

	return nil
}

func InitClient(q *Queue) {
	connectionStr := fmt.Sprintf("%s:%d", q.Connection.Host, q.Connection.Port)

	opts := &redis.Options{
		Addr: connectionStr,
		DB:   0,
	}

	if q.Connection.Password != nil {
		opts.Password = *q.Connection.Password
	}

	q.Client = redis.NewClient(opts)
}
