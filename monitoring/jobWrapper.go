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

func newJobWrapper(metric IMetric) jobWrapper {
	result := jobWrapper{}
	result.Name = metric.JobName()
	if metric.LastRun() == 0 {
		result.LastExec = never
	}else {
		result.LastExec = time.Unix(metric.LastRun(), 0).Format(time.StampMilli)
	}
	result.Runs = metric.Runs()
	result.Results = metric.Results()
	result.Disabled = metric.IsDisabled()
	if metric.StoppedAt() == 0 {
		result.Stopped = never
	}else{
		result.Stopped  = time.Unix(metric.StoppedAt(), 0).Format(time.StampMilli)
	}
	return result
}
