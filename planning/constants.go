package planning

import (
	"mailAssistant/logging"
	"time"
)

const (
	daily   = "daily"
	hourly  = "hourly"
	weekly  = "weekly"
	monthly = "monthly"
	yearly  = "yearly"
	seconds = "seconds"
	second  = "second"
	minutes = "minutes"
	minute  = "minute"
	hours   = "hours"
	hour    = "hour"
	days    = "days"
	day     = "day"
	month   = "month"
	// -----------------------------------------------------
	zero             = time.Duration(0)

	// Invalid represents an invalid schedule string
	Invalid          = time.Duration(-1)
	nanosPerSecond   = 1000000000;
	secondsPerMinute = 60
	minutesPerHour   = 60;
	hoursPerDay      = 24
	secondsPerHour   = secondsPerMinute * minutesPerHour
	secondsPerDay    = secondsPerHour * hoursPerDay;
	// -----------------------------------------------------
)

var durationLogger *logging.Logger

func logger() *logging.Logger {
	if durationLogger == nil{
		durationLogger=logging.NewLogger()
	}
	return durationLogger
}