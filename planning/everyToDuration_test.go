package planning

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var everyTestCases = []testCase{ //
	{"Every15Seconds", "EverY 15 seconds", 15 * time.Second},  //
	{"EveryYSeconds", "EverY Y seconds", Invalid},             //
	{"EverySecond", "Every second", time.Second},              //
	{"Every5Minutes", "EvErY 5 MinUtes", 5 * time.Minute},     //
	{"EveryYMinutes", "EvErY Y MinUtes", Invalid},             //
	{"EveryMinute", "Every minute", time.Minute},              //
	{"Every2Hours", "EverY 2 hours", 2 * time.Hour},           //
	{"EveryYHours", "EverY Y hours", Invalid},                 //
	{"EveryHour", "Every hour", time.Hour},                    //
	{"Every4Days", "EverY 4 dAyS", 4 * 24 * time.Hour},        //
	{"EveryYDays", "EverY Y dAyS", Invalid},                   //
	{"EveryDay", "Every day", 24 * time.Hour},                 //
	{"Every6Month", "EverY 6 mOnTh", 6 * 30 * 24 * time.Hour}, //
	{"EveryYMonth", "EverY Y mOnTh", Invalid},                 //
	{"EveryMonth", "Every MonTh", 30 * 24 * time.Hour},        //
	{"EveryInvalid", "Every Invalid", Invalid},                //
}

func planningEveryToDuration(t *testing.T) {
	for _, testCase := range everyTestCases {
		t.Run("Test"+testCase.name, func(t1 *testing.T) {
			got := ParseSchedule(testCase.value)
			require.Equal(t, testCase.result, got, testCase.value)
		})
	}
}