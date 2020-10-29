package monitoring

import (
	"strings"
)

type observable interface {
	Run()
	GetMetric() IMetric
}

// IMetric will be wrapped for jobMonitoring
type IMetric interface {
	JobName() string
	LastRun() int64
	StoppedAt() int64
	Runs() uint64
	Results() uint64
	IsDisabled() bool
}

var jobsCollector = make(map[string]*observable)

// Observe is the central registry method for monitoring
func Observe(name string, j observable) {
	jobsCollector[strings.ToUpper(name)] = &j
}
