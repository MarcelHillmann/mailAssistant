package planning

import (
	"time"
)

func constToDuration(keyword string) time.Duration {
	switch keyword {
	case second:
		return time.Second
	case minute:
		return time.Minute
	case hourly:
		return time.Hour
	case daily:
		return 24 * constToDuration(hourly)
	case weekly:
		return 7 * constToDuration(daily)
	case monthly:
		return 30 * constToDuration(daily)
	case yearly:
		return 12 * constToDuration(monthly)
	} // select keyword

	logger().Warn("Invalid const @", keyword)
	return Invalid
}
