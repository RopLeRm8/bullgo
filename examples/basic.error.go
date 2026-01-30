package examples

import (
	queue_domain "bullgo/internal/domain/queue"
	redis_domain "bullgo/internal/domain/redis"
	"bullgo/jobs"
	"bullgo/queue"
	"bullgo/worker"
	"errors"
	"fmt"
)

// Queue = {Name:"TestQueue", Redis: Default Conn}

// Job = {ID: "1", Name: "Test", Data: RandomData{Value: 5}, Delay: 1}
// Job Options = {Attempts: 3, Backoff: jobs.BackoffType{Type: "exponential", Delay: 10}, RemoveOnComplete: true, RemoveOnFailure: false}
// Delay in seconds

// Worker = {QueueName: "TestQueue", Concurrency: 5, Cb: print out "Testing first worker\n" and error}

func BasicError() {
	queueData := []queue_domain.QueueOption{queue.NameLayer("TestQueue"), queue.RedisLayer(redis_domain.RedisConnection{Host: "127.0.0.1", Port: 6379, Password: nil})}
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
		QueueName:   "TestQueue",
		Concurrency: &workerConcurrency,
		Cb: func(job jobs.Job) error {
			fmt.Print("Testing first worker\n")
			return errors.New("Test err")
		},
	})

	worker.CreateWorker(workerData)

	select {} // To prevent Go app to exit gracefully
}
