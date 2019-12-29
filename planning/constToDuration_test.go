package planning

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var constTestCases = []testCase{ //
	{"Second", "@second", time.Second},                //
	{"Minute", "@minute", time.Minute},                //
	{"Hourly", "@hoUrLy", time.Hour},                  //
	{"Daily", "@daily", 24 * time.Hour},               //
	{"Weekly", "@weekly", 7 * 24 * time.Hour},         //
	{"Monthly", "@monthly", 30 * 24 * time.Hour},      //
	{"Yearly", "@yearly", 360 * 24 * time.Hour},       //
	{"Annually", "@annually", Invalid}, //
}

func planningConstToDuration(t *testing.T) {
	for _, testCase := range constTestCases {
		t.Run("Test"+testCase.name, func(t1 *testing.T) {
			got := ParseSchedule(testCase.value)
			require.Equal(t, testCase.result, got, testCase.value)
		})
	}
}