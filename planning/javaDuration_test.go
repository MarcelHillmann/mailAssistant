package planning

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var javaTestCases = []testCase{ //
	{"parses as 20.345 seconds", "PT0.0S", 0},                                                             //
	{"parses as 20.345 seconds", "PT20.345S", 20*time.Second + 345*time.Millisecond},                      //
	{"parses as 15 minutes", "PT15M", 15 * time.Minute},                                                   //
	{"parses as 10 hours", "PT10H", 10 * time.Hour},                                                       //
	{"parses as 2 days", "P2D", 2 * 24 * time.Hour},                                                       //
	{"parses as 2 days, 3 hours and 4 minutes", "P2DT3H4M", 2*24*time.Hour + 3*time.Hour + 4*time.Minute}, //
	{"parses as -6 hours and +3 minutes", "PT-6H3M", 6*time.Hour*-1 + 3*time.Minute},                      //
	{"parses as -6 hours and -3 minutes", "-PT6H3M", (6*time.Hour + 3*time.Minute) * -1},                  //
	{"parses as +6 hours and -3 minutes", "-PT-6H+3M", 6*time.Hour + (3 * time.Minute * -1)},              //
	{"parses as -20.345 seconds", "PT-20.345S", (20*time.Second + 345*time.Millisecond) * -1},             //
	{"parses as -20.345 seconds", "PT-20.345S", (20*time.Second + 345*time.Millisecond) * -1},             //
	{"parses as -20.345 seconds", "PT-20.S", (20 * time.Second) * -1},                                     //
	{"parses as INVALID PRE1", "P", Invalid},                                                              //
	{"parses as INVALID PRE2", "PT", Invalid},                                                             //
	{"parses as INVALID POST", "P-6H+3M", Invalid},                                                        //
}

func planningJavaToDuration(t *testing.T) {
	for _, testCase := range javaTestCases {
		t.Run("Test"+testCase.name, func(t1 *testing.T) {
			got := ParseSchedule(testCase.value)
			require.Equal(t, testCase.result, got, testCase.value)
		})
	}
}
