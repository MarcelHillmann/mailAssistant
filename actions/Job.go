package actions

import (
	"mailAssistant/account"
	"mailAssistant/arguments"
	"mailAssistant/conditions"
	"mailAssistant/logging"
	"reflect"
	"runtime"
)

// NewJob is a job factory
func NewJob(jobName, name string, args map[string]interface{}, accounts *account.Accounts) Job {
	log := logging.NewNamedLogger("${project}.actions")
	fcc, ok := actions[name]
	if !ok {
		log.Severe("unknown action ", name, "for", jobName)
		fcc = newDummy
	}
	loggerName := runtime.FuncForPC(reflect.ValueOf(fcc).Pointer()).Name()
	log.Info("action ", loggerName, "for", jobName)

	semaphore[jobName] = semaphoreNull()
	return Job{Args: arguments.NewArgs(args), log: logging.NewNamedLogger(loggerName), callback: fcc, accounts: accounts, jobName: jobName}
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
}

// Run is called by clockwerk framework
func (j Job) Run() {
	j.log.Enter()
	j.callback(j, semaphore[j.jobName])
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
