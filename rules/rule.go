package rules

import (
	"github.com/onatm/clockwerk"
	"mailAssistant/account"
	"mailAssistant/actions"
	"mailAssistant/arguments"
	"mailAssistant/cntl"
	"mailAssistant/logging"
	"mailAssistant/planning"
)

var ruleLogger *logging.Logger

// Rule represent a rule yaml file
type Rule struct {
	*arguments.Args
	name     string
	schedule string
	action   string
	clock    *clockwerk.Clockwerk
	disabled bool
}

func (r Rule) getLogger() *logging.Logger {
	if ruleLogger == nil {
		ruleLogger = logging.NewNamedLogger("${project}.rules.rule")
	}
	return ruleLogger
}

// Schedule is using the clockwerk framework to schedule a rule job
func (r *Rule) Schedule(acc *account.Accounts) {
	if r.disabled {
		r.getLogger().Warn("disabled", r.name)
		return
	}
	r.getLogger().Debug("run",r.name)

	job := actions.NewJob(r.name, r.action, r.Args.GetArgs(),acc)
	interval := planning.ParseSchedule(r.schedule)
	if interval == planning.Invalid {
		return
	}
	job.GetLogger().Debug("new",job, interval)
	r.clock = cntl.NewClockwork()
	r.clock.Every(interval).Do(job)
	job.GetLogger().Debug("start",job, interval)
	r.clock.Start()
	job.GetLogger().Debug("started",job, interval)
}

// Stop is stopping a rule scheduled job
func (r Rule) Stop() {
	if r.clock == nil {
		return
	}
	r.clock.Stop()
}