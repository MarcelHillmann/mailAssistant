package planning

import (
	"github.com/bmizerany/assert"
	"testing"
	"time"
)

var goTestCases = []testCase{ //
	{"300ms", "300ms", 300 * time.Millisecond},       //
	{"1.5h", "1.5h", 1*time.Hour + 30*time.Minute},   //
	{"2h45m", "1h45m", 1*time.Hour + 45*time.Minute}, //
	{"2Days", "48h0m", 48 * time.Hour},               //
}

func planningGoToDuration(t *testing.T) {
	for _, testCase := range goTestCases {
		t.Run("Test"+testCase.name, func(t1 *testing.T) {
			got := ParseSchedule(testCase.value)
			assert.Equal(t, testCase.result, got, testCase.value)
		})
	}
}
