package monitoring

import (
	"time"
)

const never = "NEVER"

type jobWrapper struct {
	Name     string
	Disabled bool
	LastExec string
	Runs     uint64
	Results  uint64
	Stopped  string
}

func newJobWrapper(job *observable) jobWrapper {
	nJob := *job
	result := jobWrapper{}
	result.Name = nJob.JobName()
	if nJob.LastRun() == 0 {
		result.LastExec = never
	}else {
		result.LastExec = time.Unix(nJob.LastRun(), 0).Format(time.StampMilli)
	}
	result.Runs = nJob.Runs()
	result.Results = nJob.Results()
	result.Disabled = nJob.IsDisabled()
	if nJob.StoppedAt() == 0 {
		result.Stopped = never
	}else{
		result.Stopped  = time.Unix(nJob.StoppedAt(), 0).Format(time.StampMilli)
	}
	return result
}
