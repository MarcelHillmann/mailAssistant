package actions

import (
	"mailAssistant/account"
	"mailAssistant/arguments"
	"mailAssistant/conditions"
	"mailAssistant/logging"
	"mailAssistant/monitoring"
	"reflect"
	"runtime"
	"time"
)

// NewJob is a job factory
func NewJob(jobName, name string, args map[string]interface{}, accounts *account.Accounts, disabled bool) (job Job) {
	log := logging.NewNamedLogger("${project}.actions")
	fcc, ok := actions[name]
	if !ok {
		log.Severe("unknown action ", name, "for", jobName)
		fcc = newDummy
	}
	loggerName := runtime.FuncForPC(reflect.ValueOf(fcc).Pointer()).Name()
	log.Info("action ", loggerName, "for", jobName)

	semaphore[jobName] = semaphoreNull()
	if disabled {
		args["disabled"] = true
	}
	job = Job{Args: arguments.NewArgs(args), log: logging.NewNamedLogger(loggerName), callback: fcc, accounts: accounts, jobName: jobName}
	monitoring.Observe(&job)
	return
}

func semaphoreNull() *int32 {
	result := Released
	return &result
}

var semaphore = make(map[string]*int32)

// Job represents a job for scheduling
type Job struct {
	*arguments.Args
	log      *logging.Logger
	callback jobCallBack
	accounts *account.Accounts
	jobName  string
	saveTo   string
	lastRun  int64
	runs     uint64
	results  uint64
	stoppedAt  int64
}

// Run is called by clockwerk framework
func (j Job) Run() {
	j.log.Enter()
	j.callback(j, semaphore[j.jobName])
	j.runs++
	j.lastRun = time.Now().Unix()
	j.log.Leave()
}

// GetAccount is checking and returning the searched account
func (j Job) GetAccount(name string) *account.Account {
	if ! j.accounts.HasAccount(name) {
		j.log.Severe(name, "is not defined")
	}
	return j.accounts.GetAccount(name)
}
func (j *Job) getSaveTo() string {
	return saveTo(j)
}
func (j Job) getSearchParameter() []interface{} {
	result := conditions.ParseYaml(j.GetArg("search"))
	return result.Get()
}
// GetLogger is returning the job logger
func (j Job) GetLogger() *logging.Logger {
	return j.log
}

// JobName returns the internal jobName
func (j Job) JobName() string {
	return j.jobName
}

// LastRun returns the epoch from last execution
func (j Job) LastRun() int64 {
	return j.lastRun
}

// Runs returns the number of runes
func (j Job) Runs() uint64 {
	return j.runs
}

// Results return the number of executed mails
func (j Job) Results() uint64 {
	return j.results
}

// IsDisabled returns the internal disabled
func (j Job) IsDisabled() bool {
	return j.GetBool("disabled")
}

// Stopped that the epoch for descheduling
func (j Job) Stopped() {
	j.stoppedAt = time.Now().Unix()
}

// StoppedAt returns the epoch from stopping
func (j Job) StoppedAt() int64 {
	return j.stoppedAt
}
