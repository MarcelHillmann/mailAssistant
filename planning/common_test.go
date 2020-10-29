package planning

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type testCase struct {
	name   string
	value  string
	result time.Duration
}

func TestPlanning(t *testing.T) {
	t.Run("ConstToDuration", planningConstToDuration)
	t.Run("EveryToDuration", planningEveryToDuration)
	t.Run("JavaToDuration", planningJavaToDuration)
	t.Run("GoToDuration", planningGoToDuration)
	t.Run("Unknown", planningUnknown)
}

func planningUnknown(t *testing.T) {
	require.Equal(t, Invalid, ParseSchedule("??????"))

}
