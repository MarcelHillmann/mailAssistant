package monitoring

import (
	bytes2 "bytes"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"html/template"
	"net/http"
	"sort"
	"time"
)

const htmlTemplate = `Report Generation Time: {{ .TimeGenerated }}
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
}

type data struct {
	Jobs          []jobWrapper
	TimeGenerated string
}

func newData(j []jobWrapper) data {
	return data{j, time.Now().Format(time.RFC3339)}
}

func (j jobMonitoring) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	switch request.RequestURI {
	case "/favicon.ico":
		response.WriteHeader(http.StatusNotFound)
	case "/yaml":
		j.textYAML(response)
	case "/json":
		j.applicationJSON(response)
	default:
		j.textPlain(response)
	}

}

func (j jobMonitoring) textPlain(response http.ResponseWriter) {
	result, err := template.New("monitoring.tmpl").Parse(htmlTemplate)
	if err != nil {
		_, _ = response.Write([]byte(err.Error()))
		return
	}

	j.execute(response, "application/json", func(wrappers []jobWrapper) ([]byte, error) {
		buffer := bytes2.NewBuffer([]byte{})
		err = result.Execute(buffer, newData(wrappers))
		return buffer.Bytes(), err
	})
}

func (j jobMonitoring) applicationJSON(response http.ResponseWriter) {
	j.execute(response, "application/json", func(wrappers []jobWrapper) (bytes []byte, err error) {
		return json.MarshalIndent(wrappers, "","   ")
	})
}
func (j jobMonitoring) textYAML(response http.ResponseWriter) {
	j.execute(response, "text/yaml", func(wrappers []jobWrapper) (bytes []byte, err error) {
		return yaml.Marshal(wrappers)
	})
}

func (j jobMonitoring) FavIcon(writer http.ResponseWriter) {
	writer.WriteHeader(404)
}

func (j jobMonitoring) execute(response http.ResponseWriter, contentType string, writer func([]jobWrapper) ([]byte, error)) {
	keys := make([]string, 0)
	for key := range jobsCollector {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	jobs := make([]jobWrapper, len(jobsCollector))
	i := 0
	for _, jobName := range keys {
		jobs[i] = newJobWrapper(jobsCollector[jobName])
		i++
	}

	if out, err := writer(jobs); err == nil && out == nil {
		response.WriteHeader(http.StatusInsufficientStorage)
	} else if err == nil {
		response.Header().Add("Content-Type", contentType)
		_, _ = response.Write(out)
	} else {
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(err.Error()))
	}
}
