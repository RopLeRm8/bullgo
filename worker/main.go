package worker

import (
	queue_domain "bullgo/internal/domain/queue"
	"bullgo/jobs"
	"bullgo/queue"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	Max      int
	Duration int
}

type Worker struct {
	QueueName   string
	Concurrency int
	Limiter     *Limiter
	Cb          func(job jobs.Job) error
}

type WorkerConfig struct {
	QueueName   string
	Concurrency *int
	Limiter     *Limiter
	Cb          func(job jobs.Job) error
}

type RateLimitWindowsT struct {
	Uses        int
	windowStart time.Time
}

type WorkerOption func(*Worker)

func CreateConfiguredWorker(conf WorkerConfig) WorkerOption {
	return func(w *Worker) {
		w.QueueName = conf.QueueName
		w.Cb = conf.Cb

		if conf.Concurrency != nil {
			w.Concurrency = *conf.Concurrency
		}

		w.Limiter = conf.Limiter
	}
}

func CreateWorker(opts ...WorkerOption) *Worker {
	w := &Worker{
		Concurrency: 10,
	}

	for _, opt := range opts {
		opt(w)
	}

	queues := queue.Queues

	var workerQueue *queue_domain.Queue

	for _, q := range queues {
		if q.Name == w.QueueName {
			workerQueue = &q
			break
		}
	}

	if workerQueue == nil {
		panic("Queue not found for the worker")
	}

	fmt.Printf("Creating worker for queue %s\n", workerQueue.Name)
	go w.Init(workerQueue)
	go w.CreateScheduler(workerQueue)

	return w
}

func (w *Worker) CreateScheduler(workerQ *queue_domain.Queue) {
	ctx := context.Background()
	delayedJobsKey := fmt.Sprintf("queue:%s:delayed", workerQ.Name)
	activeJobsKey := fmt.Sprintf("queue:%s:active", workerQ.Name)

	for {
		delayedJobs, _ := workerQ.Client.ZRangeArgsWithScores(ctx, redis.ZRangeArgs{
			Key:   delayedJobsKey,
			Start: 0,
			Stop:  0,
		}).Result()

		if len(delayedJobs) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		fmt.Printf("%+v", delayedJobs)

		nextJob := delayedJobs[0]
		runAt := float64(nextJob.Score)
		now := float64(time.Now().Unix())

		if runAt > now {
			time.Sleep(time.Duration(runAt-now) * time.Second)
			continue
		}

		jobID := nextJob.Member.(string)

		workerQ.Client.ZRem(ctx, delayedJobsKey, jobID) // Removing from zset
		workerQ.Client.LPush(ctx, activeJobsKey, jobID) // Pushing back to the active set
	}
}

func (w *Worker) Init(workerQ *queue_domain.Queue) {
	var RateLimitWindows = RateLimitWindowsT{Uses: 0, windowStart: time.Now()}
	ctx := context.Background()

	for {
		for i := 0; i < w.Concurrency; i++ {
			activeJobID, err := workerQ.Client.BLPop(
				ctx,
				0,
				fmt.Sprintf("queue:%s:active", workerQ.Name),
			).Result()

			fmt.Printf("%s\n", activeJobID)

			if err != nil {
				continue
			}

			now := time.Now()
			jobID := activeJobID[1] // [1] is the job id, [0] is the key (we dont want it)

			if w.Limiter != nil {
				windowEnd := RateLimitWindows.windowStart.Add(
					time.Second * time.Duration(w.Limiter.Duration),
				)

				if now.After(windowEnd) {
					RateLimitWindows = RateLimitWindowsT{Uses: 0}
				}

				if now.Before(windowEnd) {
					RateLimitWindows = RateLimitWindowsT{Uses: RateLimitWindows.Uses + 1}
				}

				if RateLimitWindows.Uses >= w.Limiter.Max {
					DelayJob(workerQ, 5, jobID)
					continue
				}
				RateLimitWindows = RateLimitWindowsT{windowStart: now}
			}

			go w.Process(workerQ, jobID)
		}
	}
}

func (w *Worker) Process(workerQ *queue_domain.Queue, jobID string) {
	ctx := context.Background()

	redisKey := fmt.Sprintf("queue:%s:job:%s", workerQ.Name, jobID)
	failedJob := fmt.Sprintf("queue:%s:failed", workerQ.Name)
	delayedJob := fmt.Sprintf("queue:%s:delayed", workerQ.Name)
	attemptKey := fmt.Sprintf("queue:%s:job:%s:attempts", workerQ.Name, jobID)

	jobData, fetchErr := workerQ.Client.Get(ctx, redisKey).Result()

	if fetchErr != nil {
		workerQ.Client.ZAdd(ctx, failedJob)
		return
	}

	var jobDataParsed jobs.Job
	parseErr := json.Unmarshal([]byte(jobData), &jobDataParsed)

	if parseErr != nil {
		return
	}

	if jobDataParsed.Delay != nil {
		time.Sleep(time.Second * time.Duration(*jobDataParsed.Delay))
	}

	err := w.Cb(jobDataParsed)

	delKeys := func() {
		workerQ.Client.Del(ctx, redisKey)   // Deleting payload for redis
		workerQ.Client.Del(ctx, attemptKey) // Deleting attempts to not pollute that data
	}

	if err == nil { // On success
		fmt.Printf("Succesfully executed job %s", jobID)
		if jobDataParsed.Options.RemoveOnComplete {
			delKeys()
		}
		return
	}

	retried, err := workerQ.Client.Incr(ctx, attemptKey).Result()

	if err != nil {
		return
	}

	if jobDataParsed.Options.RemoveOnFailure { // On fail if configured
		fmt.Printf("Removed job %s because remove on failure was configured", jobID)
		workerQ.Client.RPush(ctx, failedJob, jobID)
		delKeys()
		return
	}

	if retried >= int64(jobDataParsed.Options.Attempts) {
		fmt.Printf("Too many retries for job %s", jobID)
		delKeys()
		workerQ.Client.RPush(ctx, failedJob, jobID) // Push into failed jobs (for logs)
		return
	}

	fmt.Printf("Delaying job %s\n", jobID)
	runAt := time.Now().Add(time.Second * time.Duration(jobDataParsed.Options.Backoff.Delay)).Unix()

	workerQ.Client.ZAdd(ctx, delayedJob, redis.Z{ // On fail if not configured it delays the job
		Score:  float64(runAt),
		Member: jobID,
	})
}

func DelayJob(workerQ *queue_domain.Queue, delay int, jobID string) {
	ctx := context.Background()
	runAt := time.Now().Add(time.Second * time.Duration(delay)).Unix()
	delayedJob := fmt.Sprintf("queue:%s:delayed", workerQ.Name)

	workerQ.Client.ZAdd(ctx, delayedJob, redis.Z{ // On fail if not configured it delays the job
		Score:  float64(runAt),
		Member: jobID,
	})
}
