package monitoring

import (
	"mailAssistant/cntl"
	"net/http"
)

type observable interface {
	JobName() string
	LastRun() int64
	StoppedAt() int64
	Runs() uint64
	Results() uint64
	IsDisabled() bool
}

var jobsCollector = make([]*observable,0)

// Observe is the central registry method for monitoring
func Observe(j observable){
	jobsCollector = append(jobsCollector, &j)
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

