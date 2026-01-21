package main

import (
	"bullgo/jobs"
	"fmt"
)

func main() {

	jobData := []jobs.JobOption{
		jobs.Attempts(1),
		jobs.Backoff(jobs.BackoffType{Type: "fixed", Delay: 5000}),
	}

	jobId := "test"

	fullJobData := jobs.JobOptions(jobs.JobOptionsConfig{JobId: &jobId})

	job := jobs.CreateJob(jobData...)
	fulljob := jobs.CreateJob(fullJobData)

	fmt.Printf("%+v\n", job)
	fmt.Printf("%+v\n", fulljob)
}
