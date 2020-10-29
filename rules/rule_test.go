package rules

import (
	"github.com/onatm/clockwerk"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRule(t *testing.T) {
	t.Run("Logger", ruleGetLogger)
	t.Run("Stop", ruleStop)
	t.Run("Schedule", ruleSchedule)

}

func ruleSchedule(t *testing.T) {
	t.Run("passed", func(t *testing.T) {
		r := newRule("foo.bar", "1s", "dummy", false)
		r.Schedule(nil)
		require.Equal(t, "foo.bar", r.name)
		require.Equal(t, "1s", r.schedule)
		require.Equal(t, "dummy", r.action)
		require.NotNil(t, r.clock)
		r.Stop()
	})

	t.Run("invalid schedule", func(t *testing.T) {
		r := newRule("foo.bar", "", "dummy", false)
		r.Schedule(nil)
		require.Equal(t, "foo.bar", r.name)
		require.Equal(t, "", r.schedule)
		require.Equal(t, "dummy", r.action)
		require.Nil(t, r.clock)
	})

	t.Run("disabled", func(t *testing.T) {
		r := newRule("foo.bar", "1s", "dummy", true)
		r.clock = nil
		r.disabled = true
		r.Schedule(nil)
		require.Equal(t, "foo.bar", r.name)
		require.Equal(t, "1s", r.schedule)
		require.Equal(t, "dummy", r.action)
		require.Nil(t, r.clock)
	})

}

func ruleGetLogger(t *testing.T) {
	ruleLogger = nil
	r := Rule{}
	log := r.getLogger()
	require.NotNil(t, log)
	require.NotNil(t, ruleLogger)

	log2 := r.getLogger()
	require.Same(t, log, log2)
	require.Same(t, ruleLogger, log2)
	ruleLogger = nil
}

func ruleStop(_ *testing.T) {
	r := Rule{clock: nil}
	r.Stop()

	r.clock = &clockwerk.Clockwerk{}
	r.Stop()
}
