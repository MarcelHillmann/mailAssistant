package monitoring

import (
	"bytes"
	"encoding/json"
	"github.com/openzipkin/zipkin-go"
	"gopkg.in/yaml.v2"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const htmlTemplate = `Report Generation Time: {{ .TimeGenerated }}
Active: {{ .Active }}
Passive: {{ .Passive }}
{{ range $index, $job := .Jobs }}
{{$job.Name }}
---------------------------------------
      Disabled: {{ $job.Disabled }}
       stopped: {{ $job.Stopped }}
  last execute: {{ $job.LastExec }}
      executed: {{ $job.Runs}} time
        result: {{ $job.Results }}
{{ end }}`

type jobMonitoring struct {
	*zipkin.Tracer
	zipkin.SpanOption
}

type data struct {
	TimeGenerated   string
	Jobs            []jobWrapper
	Active, Passive int
}

func newData(j []jobWrapper, active, passive int) data {
	return data{time.Now().Format(time.RFC3339), j, active, passive}
}

func (jobMonitoring) query(request *http.Request) (disabled, enabled bool, filter string) {
	uri, _ := url.ParseRequestURI(request.RequestURI)
	query := uri.Query()
	qDisabled, qEnabled, filter := query.Get("disabled"), query.Get("enabled"), query.Get("filter")
	disabled, enabled = strings.EqualFold(qDisabled, "j") || strings.EqualFold(qDisabled, "1") || strings.EqualFold(qDisabled, "y"),
		strings.EqualFold(qEnabled, "j") || strings.EqualFold(qEnabled, "1") || strings.EqualFold(qEnabled, "y")
	return
}
func (j jobMonitoring) textPlain(response http.ResponseWriter, request *http.Request) { // ,onlyDisabled, onlyEnabled bool) {
	span := j.StartSpan("AsPlain", j.SpanOption)
	result, err := template.New("monitoring.tmpl").Parse(htmlTemplate)
	if err != nil {
		_, _ = response.Write([]byte(err.Error()))
		return
	}

	disabled, enabled, filter := j.query(request)
	j.execute(response, disabled, enabled, filter, "application/json", func(wrappers []jobWrapper, a int, p int) ([]byte, error) {
		buffer := bytes.NewBuffer([]byte{})
		err = result.Execute(buffer, newData(wrappers, a, p))
		return buffer.Bytes(), err
	})
	span.Finish()
}

func (j jobMonitoring) applicationJSON(response http.ResponseWriter, request *http.Request) { // ,onlyDisabled, onlyEnabled bool) {
	span := j.StartSpan("AsJson", j.SpanOption)
	disabled, enabled, filter := j.query(request)
	j.execute(response, disabled, enabled, filter,
		"application/json", func(wrappers []jobWrapper, a int, p int) (bytes []byte, err error) {
			return json.MarshalIndent(newData(wrappers, a, p), "", "   ")
		})
	span.Finish()
}
func (j jobMonitoring) textYAML(response http.ResponseWriter, request *http.Request) { // ,onlyDisabled, onlyEnabled bool) {
	span := j.StartSpan("AsYaml", j.SpanOption)
	disabled, enabled, filter := j.query(request)
	j.execute(response, disabled, enabled, filter, "text/yaml", func(wrappers []jobWrapper, a int, p int) (bytes []byte, err error) {
		return yaml.Marshal(newData(wrappers, a, p))
	})
	span.Finish()
}

func (j jobMonitoring) favicon(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(http.StatusNotFound)
}

func (j jobMonitoring) execute(response http.ResponseWriter, onlyDisabled, onlyEnabled bool, filter, contentType string, writer func([]jobWrapper, int, int) ([]byte, error)) {
	keys := make([]string, 0)
	for key := range jobsCollector {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	active, passive := 0, 0
	jobs := make([]jobWrapper, 0)
	for _, jobName := range keys {
		job := jobsCollector[jobName]
		m := (*job).GetMetric()
		if m.IsDisabled() {
			passive++
		} else {
			active++
		}

		if filter != "" && (strings.EqualFold(m.JobName(), filter) || strings.Contains(strings.ToLower(m.JobName()), strings.ToLower(filter))) {
			jobs = append(jobs, newJobWrapper(m))
		} else if onlyDisabled && m.IsDisabled() ||
			onlyEnabled && !m.IsDisabled() ||
			onlyDisabled == onlyEnabled {
			jobs = append(jobs, newJobWrapper(m))
		}
	}

	if out, err := writer(jobs, active, passive); err == nil && out == nil {
		response.WriteHeader(http.StatusInsufficientStorage)
	} else if err == nil {
		response.Header().Add("Content-Type", contentType)
		_, _ = response.Write(out)
	} else {
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(err.Error()))
	}
}

func (j jobMonitoring) runJob(response http.ResponseWriter, request *http.Request) { //, query url.Values) {
	span := j.StartSpan("run_job", j.SpanOption)
	uri, _ := url.ParseRequestURI(request.RequestURI)
	name := uri.Query().Get("name")

	for jobName, job := range jobsCollector {
		if strings.EqualFold(name, jobName) {
			(*job).Run()
			response.WriteHeader(http.StatusOK)
			span.Finish()
			return
		}
	}

	response.WriteHeader(http.StatusNotImplemented)
	span.Finish()
}

func (j jobMonitoring) missingJobName(writer http.ResponseWriter, _ *http.Request) {
	span := j.StartSpan("missing_job_name", j.SpanOption)
	writer.Write([]byte(`missing parameter name`))
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusBadRequest)
	span.Finish()
}
