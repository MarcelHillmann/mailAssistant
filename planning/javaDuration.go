package planning

import (
	"regexp"
	"strconv"
	"time"
)

func javaDuration(javaString string) time.Duration {
	regEx := regexp.MustCompile("([-+]?)P(?:([-+]?[0-9]+)D)?" + //
		"(T(?:([-+]?[0-9]+)H)?(?:([-+]?[0-9]+)M)?(?:([-+]?[0-9]+)(?:[.,]([0-9]{0,9}))?S)?)?")
	parts := regEx.FindAllStringSubmatch(javaString, -1)
	if "T" != parts[0][3] {
		negate := "-" == parts[0][1]
		dayMatch := parts[0][2]
		hourMatch := parts[0][4]
		minuteMatch := parts[0][5]
		secondMatch := parts[0][6]
		fractionMatch := parts[0][7]
		if dayMatch != "" || hourMatch != "" || minuteMatch != "" || secondMatch != "" {
			daysAsSecs := parseNumber(javaString, dayMatch, secondsPerDay, "days")
			hoursAsSecs := parseNumber(javaString, hourMatch, secondsPerHour, "hours")
			minsAsSecs := parseNumber(javaString, minuteMatch, secondsPerMinute, "minutes")
			seconds := parseNumber(javaString, secondMatch, 1, "seconds")
			var nanos time.Duration
			if seconds < 0 {
				nanos = parseFraction(javaString, fractionMatch, -1)
			} else {
				nanos = parseFraction(javaString, fractionMatch, 1)
			}

			sum := addExact(daysAsSecs, addExact(hoursAsSecs, addExact(minsAsSecs, seconds)))
			if negate {
				return ofSeconds(sum, nanos) * -1
			}
			return ofSeconds(sum, nanos)
		}
	}
	return Invalid
}

func ofSeconds(seconds, nanos time.Duration) time.Duration {
	secs := addExact(seconds, floorDiv(nanos, nanosPerSecond))
	nos := floorMod(nanos, nanosPerSecond)

	if int(secs^nos) == 0 {
		return zero
	}
	return secs*time.Second + nos
}
func floorMod(x, y time.Duration) time.Duration {
	return x - floorDiv(x, y)*y
}
func floorDiv(x, y time.Duration) time.Duration {
	r := x / y
	if (x^y) < 0 && (r*y != x) {
		r--
	}
	return r
}
func addExact(x, y time.Duration) time.Duration {
	r := x + y
	if ((x ^ r) & (y ^ r)) < 0 {
		logger().Panic("long overflow")
	}
	return r
}
func parseFraction(text string, parsed string, negate int) time.Duration {
	_ = text
	_parsed := (parsed + "000000000")[0:9]
	iParsed, _ := strconv.Atoi(_parsed)
	return time.Duration(iParsed * negate)
}
func parseNumber(text, parsed string, multiplier int, field string) time.Duration {
	if parsed == "" {
		return zero
	}
	_, _ = text, field
	iP, _ := strconv.Atoi(parsed)
	return time.Duration(iP * multiplier)
}