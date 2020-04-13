package monitoring

import (
	"github.com/gorilla/mux"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	logreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"mailAssistant/cntl"
	"net/http"
	"os"
	"time"
)

// StartServer is launching the monitoring http server
func StartServer(zipKinServer string) error {
	reporter := logreporter.NewReporter(zipKinServer+"/api/v2/spans",
		/*                           */ logreporter.Timeout(time.Second),
		/*                           */ logreporter.Logger(log.New(os.Stderr, "", log.LstdFlags)))
	endpoint, err := zipkin.NewEndpoint("mailAssistant", "localhost:8080")
	if err != nil {
		log.Panic(err)
	}
	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint), zipkin.WithSharedSpans(true))
	if err != nil {
		log.Panic(err)
	}
	srvMiddleware := zipkinhttp.NewServerMiddleware(tracer, zipkinhttp.TagResponseSize(true))
	//
	job := jobMonitoring{tracer, zipkin.RemoteEndpoint(endpoint)}
	//
	router := mux.NewRouter()
	router.Methods("GET").Path("/metrics").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		span := tracer.StartSpan("metrics", job.SpanOption)
		promhttp.Handler().ServeHTTP(writer, request)
		span.Finish()
	})
	router.Methods("GET").Path("/").Headers("Accept", "application/json").HandlerFunc(job.applicationJSON)
	router.Methods("GET").Path("/").Headers("Accept", "text/yaml").HandlerFunc(job.textYAML)
	router.Methods("GET").Path("/").HandlerFunc(job.textPlain)
	router.Methods("GET").Path("/fire").Queries("name", "{name:.*}").HandlerFunc(job.runJob)
	router.Methods("GET").Path("/fire").HandlerFunc(job.missingJobName)
	router.Methods("GET").Path("/favicon.ico").HandlerFunc(job.favicon)
	//
	server := http.Server{Addr: ":8080", Handler: srvMiddleware(router)}
	server.SetKeepAlivesEnabled(true)
	//
	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	go func() {
		cntl.WaitForNotify()
		server.Close()
		reporter.Close()
	}()
	return nil
}
