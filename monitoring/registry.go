package monitoring

import (
	"mailAssistant/cntl"
	"net/http"
	"strings"
)

type observable interface {
	Run()
	GetMetric() IMetric
}
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
func Observe(name string, j observable){
	jobsCollector[strings.ToUpper(name)] = &j
}

// StartServer is launching the monitoring http server
func StartServer() error {
	server := http.Server{Addr: ":8080", Handler:jobMonitoring{}}
	server.SetKeepAlivesEnabled(true)

	go func() {
		if err := server.ListenAndServe();err != nil {
			panic(err)
		}
	}()
	go func(){
		cntl.WaitForNotify()
		server.Close()
	}()
	return nil
}

