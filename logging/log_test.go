package logging

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	logUndo := func() {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags)
	}
	inject := func() *bytes.Buffer {
		buf := bytes.NewBufferString("")
		log.SetFlags(0)
		log.SetOutput(buf)
		return buf
	}
	pre := loggerRegistry["global"]
	defer SetLevel("global", pre.String())
	SetLevel("global", "all")

	t.Run("debug", func(t *testing.T) {
		defer logUndo()
		l := inject()
		NewGlobalLogger().Debug("debug")
		require.Equal(t, "DEBUG   [global#lambda] debug\n", l.String())
	})
	t.Run("debugF", func(t *testing.T) {
		defer logUndo()
		l := inject()
		NewGlobalLogger().Debugf("%s", "debug")
		require.Equal(t, "DEBUG   [global#lambda] debug\n", l.String())
	})
	t.Run("debug Writer", func(t *testing.T) {
		defer logUndo()
		l := inject()
		NewNamedLogWriter("global").Write([]byte("debug"))
		require.Equal(t, "DEBUG   [global#lambda] \ndebug\n", l.String())
	})
	t.Run("INFO", func(t *testing.T) {
		defer logUndo()
		l := inject()
		NewGlobalLogger().Info("info")
		require.Equal(t, "INFO    [global#lambda] info\n", l.String())
	})
	t.Run("InfoF", func(t *testing.T) {
		defer logUndo()
		l := inject()
		NewGlobalLogger().Infof("%s", "info")
		require.Equal(t, "INFO    [global#lambda] info\n", l.String())
	})
	t.Run("warn", func(t *testing.T) {
		defer logUndo()
		l := inject()
		NewGlobalLogger().Warn("warning")
		require.Equal(t, "WARNING [global#lambda] warning\n", l.String())
	})
	t.Run("warnF", func(t *testing.T) {
		defer logUndo()
		l := inject()
		NewGlobalLogger().Warnf("%s", "warning")
		require.Equal(t, "WARNING [global#lambda] warning\n", l.String())
	})
	t.Run("severe", func(t *testing.T) {
		defer logUndo()
		l := inject()
		NewGlobalLogger().Severe("severe")
		require.Equal(t, "SEVERE  [global#lambda] severe\n", l.String())
	})
	t.Run("severeF", func(t *testing.T) {
		defer logUndo()
		l := inject()
		NewGlobalLogger().Severef("%s", "severe")
		require.Equal(t, "SEVERE  [global#lambda] severe\n", l.String())
	})
	t.Run("ENTER", func(t *testing.T) {
		defer logUndo()
		l := inject()
		NewGlobalLogger().Enter()
		require.Equal(t, "DEBUG   [global#lambda] >>\n", l.String())
	})
	t.Run("LEAVE", func(t *testing.T) {
		defer logUndo()
		l := inject()
		NewGlobalLogger().Leave()
		require.Equal(t, "DEBUG   [global#lambda] <<\n", l.String())
	})
	t.Run("PANIC", func(t *testing.T) {
		defer logUndo()
		l := inject()
		defer func() {
			require.Equal(t, "SEVERE  [global#lambda] panic test\n", l.String())
			err := recover()
			require.NotNil(t, err)
			require.EqualError(t, err.(error), "panic test")
		}()
		NewGlobalLogger().Panic("panic", "test")
	})
	t.Run("Name", func(t *testing.T) {
		require.Equal(t, "test", NewNamedLogger("test").Name())
	})

	t.Run("IsLogLevel", func(t *testing.T) {
		SetLevel("test.foo", "warn")
		l := NewNamedLogger("test.foo.bar").(*iLogger)
		require.True(t, l.isLogLevel(warn))
	})

	t.Run("NewLogger", func(t *testing.T) {
		l := NewLogger()
		require.NotNil(t, l)
		require.Equal(t, "mailAssistant.logging.TestLogger.func17", l.Name())
	})

	t.Run("IsLogLevel ROOT", func(t *testing.T) {
		SetLevel("global", "warn")
		l := NewNamedLogger("foo.bar").(*iLogger)
		require.True(t, l.isLogLevel(warn))
	})
}
