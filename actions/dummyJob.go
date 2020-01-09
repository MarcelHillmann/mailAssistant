package actions

import (
	"mailAssistant/logging"
)

func newDummy(job Job, _ *int32, result func(int)) {
	logging.NewLogger().Debug(job)
	result(1)
}
