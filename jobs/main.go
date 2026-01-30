package jobs

import "time"

type DelayType string

const (
	DelayBackoff DelayType = "exponential"
	DelayFixed   DelayType = "fixed"
)

type BackoffType struct {
	Type  DelayType
	Delay int
}

type OptionsType struct {
	Attempts         int
	Backoff          BackoffType
	RemoveOnComplete bool
	RemoveOnFailure  bool
}

type OptionsConfig struct {
	Attempts         *int
	Backoff          *BackoffType
	RemoveOnComplete *bool
	RemoveOnFailure  *bool
}

type Job struct {
	ID        string
	Name      string
	Data      any
	Delay     *int
	Options   OptionsType
	CreatedAt time.Time
}

type JobInput struct {
	ID    string
	Name  string
	Delay *int
	Data  any
}

type JobOption func(*Job)

func CreateJob(input JobInput, opts ...JobOption) *Job {
	job := &Job{
		Options: OptionsType{
			Attempts:         1,
			Backoff:          BackoffType{Type: "exponential", Delay: 2000},
			RemoveOnComplete: true,
			RemoveOnFailure:  true,
		},
	}

	for _, opt := range opts {
		opt(job)
	}

	job.ID = input.ID
	job.Name = input.Name
	job.Data = input.Data
	job.Delay = input.Delay
	job.CreatedAt = time.Now()
	return job
}

func JobId(jobid string) JobOption {
	return func(j *Job) {
		j.ID = jobid
	}
}

func JobName(jobName string) JobOption {
	return func(j *Job) {
		j.Name = jobName
	}
}

func JobData(data any) JobOption {
	return func(j *Job) {
		j.Data = data
	}
}

func Attempts(attemps int) JobOption {
	return func(j *Job) {
		j.Options.Attempts = attemps
	}
}

func Backoff(backoff BackoffType) JobOption {
	return func(j *Job) {
		j.Options.Backoff = backoff
	}
}

func Removals(RemoveOnComplete bool, RemoveOnFailure bool) JobOption {
	return func(j *Job) {
		j.Options.RemoveOnComplete = RemoveOnComplete
		j.Options.RemoveOnFailure = RemoveOnFailure
	}
}

func Options(cfg OptionsConfig) JobOption {
	return func(j *Job) {
		if cfg.Attempts != nil {
			j.Options.Attempts = *cfg.Attempts
		}
		if cfg.Backoff != nil {
			j.Options.Backoff = *cfg.Backoff
		}

		if cfg.RemoveOnComplete != nil {
			j.Options.RemoveOnComplete = *cfg.RemoveOnComplete
		}
		if cfg.RemoveOnFailure != nil {
			j.Options.RemoveOnFailure = *cfg.RemoveOnFailure
		}
	}
}
