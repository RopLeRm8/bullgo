package main

import (
	queue_domain "bullgo/internal/domain/queue"
	redis_domain "bullgo/internal/domain/redis"
	"bullgo/jobs"
	"bullgo/queue"
	"bullgo/worker"
	"errors"
	"fmt"
)

type RandomData struct {
	Value int
}

func main() {
	queueData := []queue_domain.QueueOption{queue.NameLayer("Test_shit"), queue.RedisLayer(redis_domain.RedisConnection{Host: "127.0.0.1", Port: 6379, Password: nil})}
	queue := queue.CreateQueue(queueData...)

	jobAttemps := 3
	removeOnComplete := true
	removeOnFailure := false
	delay := 1

	jobData := jobs.JobInput{ID: "1", Name: "Test", Data: RandomData{Value: 5}, Delay: &delay}
	jobOptions := jobs.Options(jobs.OptionsConfig{Attempts: &jobAttemps, Backoff: &jobs.BackoffType{Type: "exponential", Delay: 10}, RemoveOnComplete: &removeOnComplete, RemoveOnFailure: &removeOnFailure})
	job := jobs.CreateJob(jobData, jobOptions)
	queue.Add(job)

	workerConcurrency := 5

	workerData := worker.CreateConfiguredWorker(worker.WorkerConfig{
		QueueName:   "Test_shit",
		Concurrency: &workerConcurrency,
		Cb: func(job jobs.Job) error {
			fmt.Print("Ayo testing first worker ever\n")
			return errors.New("Test err")
		},
	})

	worker.CreateWorker(workerData)

	select {}
}

// Workers -> always working per queue
// Jobs -> All stack up in 1 Queue (the queue is created by the user TOO!)
// Queues -> are the root of bullgo system
