package planning

import (
	"strconv"
	"strings"
	"time"
)

func everyToDuration(every string) time.Duration {
	switch {
	case strings.HasSuffix(every, days):
		finalEvery := strings.TrimSpace(strings.Replace(every, days, "", -1))
		iEvery, err := strconv.ParseInt(finalEvery, 0, 0)
		if err != nil {
			break
		}
		return time.Duration(int64(everyToDuration(day)) * iEvery)
	case strings.HasSuffix(every, hours):
		finalEvery := strings.TrimSpace(strings.Replace(every, hours, "", -1))
		iEvery, err := strconv.ParseInt(finalEvery, 0, 0)
		if err != nil {
			break
		}
		return time.Duration(int64(everyToDuration(hour)) * iEvery)
	case strings.HasSuffix(every, minutes):
		finalEvery := strings.TrimSpace(strings.Replace(every, minutes, "", -1))
		iEvery, err := strconv.ParseInt(finalEvery, 0, 0)
		if err != nil {
			break
		}
		return time.Duration(int64(everyToDuration(minute)) * iEvery)
	case strings.HasSuffix(every, seconds):
		finalEvery := strings.TrimSpace(strings.Replace(every, seconds, "", -1))
		iEvery, err := strconv.ParseInt(finalEvery, 0, 0)
		if err != nil {
			break
		}
		return time.Duration(int64(everyToDuration(second)) * iEvery)
	case strings.HasSuffix(every, month):
		finalEvery := strings.TrimSpace(strings.Replace(every, month, "", -1))
		if finalEvery == "" {
			return 30 * everyToDuration(day)
		}

		iEvery, err := strconv.ParseInt(finalEvery, 0, 0)
		if err != nil {
			break
		}
		return time.Duration(iEvery * int64(everyToDuration(month)))
	case strings.HasSuffix(every, day):
		return 24 * everyToDuration(hour)
	case strings.HasSuffix(every, hour):
		return time.Hour
	case strings.HasSuffix(every, minute):
		return time.Minute
	case strings.HasSuffix(every, second):
		return time.Second
	}

	logger().Severe("Invalid everyToDuration (", every, ")")
	return Invalid
}
