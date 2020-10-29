package planning

import (
	"mailAssistant/logging"
	"strings"
	"time"
)

// ParseSchedule is parsing a string based duration to a time.Duration
func ParseSchedule(input string) time.Duration {
	newInput := strings.ToLower(input)
	if strings.HasPrefix(newInput, "@") {
		return constToDuration(newInput[1:])
	} else if strings.HasPrefix(newInput, "every") {
		return everyToDuration(strings.Replace(newInput, "every", "", 1))
	} else if strings.HasPrefix(input, "P") || strings.HasPrefix(input, "-P") {
		return javaDuration(strings.TrimSpace(input))
	} else if goDuration, err := time.ParseDuration(input); err == nil {
		return goDuration
	} else {
		logging.NewLogger().Severe(err)
	}
	return Invalid
}
