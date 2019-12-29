package planning

import (
	"mailAssistant/logging"
	"strings"
	"time"
)

// ParseSchedule is parsing a string based duration to a time.Duration
func ParseSchedule(input string) time.Duration {
	_input := strings.ToLower(input)
	if strings.HasPrefix(_input, "@") {
		return constToDuration(_input[1:])
	} else if strings.HasPrefix(_input, "every") {
		return everyToDuration(strings.Replace(_input, "every", "", 1))
	} else if strings.HasPrefix(input, "P") || strings.HasPrefix(input, "-P") {
		return javaDuration(strings.TrimSpace(input))
	}else if goDuration, err := time.ParseDuration(input); err == nil{
		return goDuration
	}else{
		logging.NewLogger().Severe(err)
	}
	return Invalid
}
