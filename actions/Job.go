package actions

import (
	"github.com/prometheus/client_golang/prometheus"
	"mailAssistant/account"
	"mailAssistant/arguments"
	"mailAssistant/conditions"
	"mailAssistant/logging"
	"mailAssistant/monitoring"
	"reflect"
	"runtime"
)

var (
	jobAnno = []string{"job"}
	runes = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: "mailassistant", Subsystem: "job",Name: "runes", Help: "job runes", ConstLabels: nil}, jobAnno)
	results = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: "mailassistant", Subsystem: "job",Name: "results", Help: "job results", ConstLabels: nil}, jobAnno)
)

func init(){
	prometheus.MustRegister(runes, results)
}

// NewJob is a job factory
func NewJob(jobName, name string, args map[string]interface{}, accounts *account.Accounts, disabled bool) (job Job) {
	log := logging.NewNamedLogger("${project}.actions")
	fcc, ok := actions[name]
	if !ok {
		log.Severe("unknown action ", name, "for", jobName)
		fcc = newDummy
	}
	loggerName := runtime.FuncForPC(reflect.ValueOf(fcc).Pointer()).Name()
	log.Infof("action '%s' for '%s'", loggerName, jobName)

	semaphore[jobName] = semaphoreNull()

	jobCounter := runes.WithLabelValues(jobName)
	jobResults := results.WithLabelValues(jobName)

	job = Job{	Args: arguments.NewArgs(args),
				Logger: logging.NewNamedLogger(loggerName+"%"+jobName),
				Accounts: accounts,
				callback: fcc,
				jobName: jobName,
				metrics: &metrics{disabled: disabled, promRuns: jobCounter, promResults: jobResults}}
	monitoring.Observe(jobName, job)
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
	logging.Logger
	*account.Accounts
	*metrics

	callback jobCallBack
	jobName  string
	saveTo   string
}

// Run is called by clockwerk framework
func (j Job) Run() {
	j.Enter()
	j.callback(j, semaphore[j.jobName], j.result)
	j.run()
	j.Leave()
}

// GetAccount is checking and returning the searched account
func (j Job) GetAccount(name string) *account.Account {
	if ! j.HasAccount(name) {
		j.Severe(name, "is not defined")
	}
	return j.Accounts.GetAccount(name)
}

func (j *Job) getSaveTo() string {
	return saveTo(j)
}

func (j Job) getSearchParameter() []interface{} {
	result := conditions.ParseYaml(j.GetArg("search"))
	return result.Get()
}

// GetLogger is returning the job logger
func (j Job) GetLogger() logging.Logger {
	return j
}

// Stopped that the epoch for descheduled
func (j Job) Stopped() {
	j.stopped()
}

// GetMetric is exporting a IMetric object
func (j Job) GetMetric() monitoring.IMetric {
	metric := j.metrics
	metric.jobName = j.jobName
	return metric
}

