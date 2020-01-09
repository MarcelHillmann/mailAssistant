package actions

import "time"

type metrics struct {
	jobName   string
	disabled  bool
	lastRun   int64
	runs      uint64
	results   uint64
	stoppedAt int64
}

func (m *metrics) run(){
	m.lastRun = time.Now().Unix()
	m.runs ++
}

func (m *metrics) stopped() {
	m.stoppedAt = time.Now().Unix()
}

func (m *metrics) result(res int){
	m.results += uint64(res)
}

// LastRun returns the epoch from last execution
func (m metrics) LastRun() int64 {
	return m.lastRun
}

// JobName returns the internal jobName
func (m metrics) JobName() string {
	return m.jobName
}

// Runs returns the number of runes
func (m metrics) Runs() uint64 {
	return m.runs
}

// Results return the number of executed mails
func (m metrics) Results() uint64 {
	return m.results
}

// StoppedAt returns the epoch from stopping
func (m metrics) StoppedAt() int64 {
	return m.stoppedAt
}

// IsDisabled returns the internal disabled
func (m metrics) IsDisabled() bool {
	return m.disabled
}

func metricsDummy(_ int) {
	// noop
}
