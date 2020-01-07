package actions

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"log"
	"mailAssistant/logging"
	"os"
	"testing"
)
const (
	registered = "DEBUG   [mailAssistant.actions.registry#register] register  foo.bar\n"
	dublicate = "SEVERE  [mailAssistant.actions.registry#register] duplicate action name foo.bar  =>  mailAssistant/actions.glob..func1\n"
)

var dummyJob = func(job Job, waitGroup *int32) {}

func TestRegistry_Add(t *testing.T){
	actions = make(map[string]jobCallBack)
	defer func() {
		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stderr)
	}()

	logging.SetLevel("mailAssistant.actions.registry","all")
	buffer := bytes.NewBufferString("")
	log.SetFlags(0)
	log.SetOutput(buffer)

	register("foo.bar", dummyJob)
	require.Len(t, actions,1)
	require.Equal(t,registered, buffer.String())
}

func TestRegistry_Twice(t *testing.T){
	actions = make(map[string]jobCallBack)
	defer func() {
		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stderr)
	}()

	logging.SetLevel("mailAssistant.actions.registry","all")
	buffer := bytes.NewBufferString("")
	log.SetFlags(0)
	log.SetOutput(buffer)

	register("foo.bar", dummyJob)
	require.Len(t, actions,1)
	require.Equal(t,registered, buffer.String())
	buffer.Truncate(0)
	register("foo.bar", dummyJob)
	require.Len(t, actions,1)
	require.Equal(t,registered+dublicate+dublicate, buffer.String())
}
