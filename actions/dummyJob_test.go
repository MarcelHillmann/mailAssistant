package actions

import (
	"mailAssistant/logging"
	"testing"
)

func TestDummyJob(t *testing.T){
	var wg int32
	job := Job{log: logging.NewLogger()}

	newDummy(job,&wg, metricsDummy)
}
