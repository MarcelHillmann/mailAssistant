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

var ruleLogger logging.Logger

// Rule represent a rule yaml file
type Rule struct {
	*arguments.Args
	job      *actions.Job
	name     string
	schedule string
	action   string
	clock    *clockwerk.Clockwerk
	disabled bool
}

func newRule(name, schedule, action string, disabled bool) (rule Rule) {
	rule = Rule{}
	rule.Args = arguments.NewEmptyArgs()
	rule.job = nil
	rule.clock = nil
	rule.name = name
	rule.schedule = schedule
	rule.action = action
	rule.disabled = disabled
	return
}
func (r Rule) getLogger() logging.Logger {
	if ruleLogger == nil {
		ruleLogger = logging.NewNamedLogger("${project}.rules.rule")
	}
	return ruleLogger
}

// Schedule is using the clockwork framework to schedule a rule job
func (r *Rule) Schedule(acc *account.Accounts) {
	r.getLogger().Debug("run", r.name)

	job := actions.NewJob(r.name, r.action, r.Args.GetArgs(), acc, r.disabled)
	r.job = &job
	if r.disabled {
		r.getLogger().Warn("disabled", r.name)
		return
	}
	interval := planning.ParseSchedule(r.schedule)
	if interval == planning.Invalid {
		return
	}
	job.GetLogger().Debug("new", job, interval)
	r.clock = cntl.NewClockwork()
	r.clock.Every(interval).Do(job)
	job.GetLogger().Debug("start", job, interval)
	r.clock.Start()
	job.GetLogger().Debug("started", job, interval)
}

// Stop is stopping a rule scheduled job
func (r Rule) Stop() {
	if r.job != nil {
		r.job.Stopped()
	}

	if r.clock == nil {
		return
	}
	r.clock.Stop()
}
