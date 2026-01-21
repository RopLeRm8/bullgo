package jobs

type DelayType string

const (
	DelayBackoff DelayType = "exponential"
	DelayFixed   DelayType = "fixed"
)

type BackoffType struct {
	Type  DelayType
	Delay int
}

type JobOptionsType struct {
	JobId            string
	Attempts         int
	Backoff          BackoffType
	Delay            int
	RemoveOnComplete bool
	RemoveOnFailure  bool
}

type JobOptionsConfig struct {
	JobId            *string
	Attempts         *int
	Backoff          *BackoffType
	Delay            *int
	RemoveOnComplete *bool
	RemoveOnFailure  *bool
}

type Job struct {
	JobName    string
	Data       any
	JobOptions JobOptionsType
}

type JobOption func(*Job)

func CreateJob(opts ...JobOption) *Job {
	job := &Job{
		JobOptions: JobOptionsType{
			Attempts:         1,
			Backoff:          BackoffType{Type: "exponential", Delay: 2000},
			RemoveOnComplete: true,
			RemoveOnFailure:  true,
		},
	}

	for _, opt := range opts {
		opt(job)
	}

	return job
}

func JobId(jobid string) JobOption {
	return func(j *Job) {
		j.JobOptions.JobId = jobid
	}
}

func Attempts(attemps int) JobOption {
	return func(j *Job) {
		j.JobOptions.Attempts = attemps
	}
}

func Backoff(backoff BackoffType) JobOption {
	return func(j *Job) {
		j.JobOptions.Backoff = backoff
	}
}

func Delay(delay int) JobOption {
	return func(j *Job) {
		j.JobOptions.Delay = delay
	}
}

func Removals(RemoveOnComplete bool, RemoveOnFailure bool) JobOption {
	return func(j *Job) {
		j.JobOptions.RemoveOnComplete = RemoveOnComplete
		j.JobOptions.RemoveOnFailure = RemoveOnFailure
	}
}

func JobOptions(cfg JobOptionsConfig) JobOption {
	return func(j *Job) {
		if cfg.JobId != nil {
			j.JobOptions.JobId = *cfg.JobId
		}
		if cfg.Attempts != nil {
			j.JobOptions.Attempts = *cfg.Attempts
		}
		if cfg.Backoff != nil {
			j.JobOptions.Backoff = *cfg.Backoff
		}
		if cfg.Delay != nil {
			j.JobOptions.Delay = *cfg.Delay
		}
		if cfg.RemoveOnComplete != nil {
			j.JobOptions.RemoveOnComplete = *cfg.RemoveOnComplete
		}
		if cfg.RemoveOnFailure != nil {
			j.JobOptions.RemoveOnFailure = *cfg.RemoveOnFailure
		}
	}
}
