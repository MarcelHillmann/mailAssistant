package actions

import (
	"bytes"
	"github.com/emersion/go-imap"
	"github.com/stretchr/testify/require"
	"log"
	"mailAssistant/logging"
	"os"
	"testing"
	"time"
)

func TestUtilIsLockedElseLock(t *testing.T) {
	defer func() {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags)
	}()

	buffer := bytes.NewBufferString("")
	log.SetOutput(buffer)
	log.SetFlags(0)

	logging.SetLevel("mailassistant.actions", "all")
	log := logging.NewNamedLogger("mailassistant.actions")

	var wg int32
	require.False(t, isLockedElseLock(log, &wg))
	require.Equal(t, "INFO    [mailassistant.actions#islockedelselock] lock\n", buffer.String())
	require.Equal(t, Locked, wg)
	buffer.Truncate(0)
	require.True(t, isLockedElseLock(log, &wg))
	require.Equal(t, "INFO    [mailassistant.actions#islockedelselock] is locked\n", buffer.String())
	require.Equal(t, Locked, wg)
}

func TestUtilUnlockAlways(t *testing.T) {
	defer func() {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags)
	}()

	buffer := bytes.NewBufferString("")
	log.SetOutput(buffer)
	log.SetFlags(0)

	logging.SetLevel("mailassistant.actions", "all")
	log := logging.NewNamedLogger("mailassistant.actions")

	var wg int32 = 400
	unlockAlways(log, &wg)
	require.Equal(t, "INFO    [mailassistant.actions#unlockalways] unlocked\n", buffer.String())
	require.Equal(t, Released, wg)
}

func TestUtilDeferUnlockAlways(t *testing.T) {
	buffer := bytes.NewBufferString("")
	var wg int32 = 400
	defer func() {
		require.Equal(t, "SEVERE  [mailassistant.actions#unlockalways] run recover()\n"+
			"INFO    [mailassistant.actions#unlockalways] unlocked\n", buffer.String())
		require.Equal(t, Released, wg)
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags)
	}()
	log.SetOutput(buffer)
	log.SetFlags(0)
	logging.SetLevel("mailassistant.actions", "all")
	log := logging.NewNamedLogger("mailassistant.actions")

	defer unlockAlways(log, &wg)
	panic("run recover()")
}

func createMessage(num uint32, body bool) *imap.Message {
	bodyMap := make(map[*imap.BodySectionName]imap.Literal)
	if body {
		bodyMap[new(imap.BodySectionName)] = bytes.NewBufferString("")
	}
	return &imap.Message{SeqNum: num,
		Items:         make(map[imap.FetchItem]interface{}),
		Envelope:      new(imap.Envelope),
		BodyStructure: new(imap.BodyStructure),
		Flags:         make([]string, 0),
		InternalDate:  time.Now(),
		Size:          0,
		Uid:           num,
		Body:          bodyMap,
	}
}

type called struct {
	login, selected, search, fetch, store, expunge int
}
