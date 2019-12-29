package actions

import (
	"mailAssistant/logging"
)

func newDummy(job Job, wg *int32) {
	logging.NewLogger().Debug(job)
	_ = wg
}
