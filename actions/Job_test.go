package actions

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"log"
	"mailAssistant/account"
	"mailAssistant/logging"
	"os"
	"runtime"
	"testing"
)

func TestNewJob(t *testing.T) {
	j := NewJob("dummy", "foo_bar", make(map[string]interface{}), nil, false)
	require.NotNil(t, j)
}

func TestJob_Run(t *testing.T) {
	j := NewJob("dummy", "foo_bar", make(map[string]interface{}), nil, false)
	logging.SetLevel("global", "all")
	defer func() {
		logging.SetLevel("*", "")
		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stderr)
	}()
	buffer := bytes.NewBufferString("")
	log.SetFlags(0)
	log.SetOutput(buffer)
	j.Run()
	require.NotEmpty(t, buffer.String())
	require.Equal(t, "DEBUG   [mailAssistant.actions.newDummy%dummy#job.run] >>\n"+
		"DEBUG   [mailAssistant.actions.newDummy%dummy#newdummy] map[]\n"+
		"DEBUG   [mailAssistant.actions.newDummy%dummy#job.run] <<\n", buffer.String())
}

func TestJob_GetAccount(t *testing.T) {
	acc := account.Accounts{}
	acc.Account = make(map[string]account.Account)
	acc.Account["test"] = account.NewAccountForTest(t, "test", "foo", "bar", "l", false)
	j := NewJob("dummy", "foo_bar", make(map[string]interface{}), &acc, false)
	logging.SetLevel("global", "all")
	defer func() {
		logging.SetLevel("*","")
		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stderr)
	}()
	buffer := bytes.NewBufferString("")
	log.SetFlags(0)
	log.SetOutput(buffer)
	require.NotNil(t, j.GetAccount("test"))
	require.Nil(t, j.GetAccount("test22"))
	require.NotEmpty(t, buffer.String())
	require.Equal(t, "SEVERE  [mailAssistant.actions.newDummy%dummy#job.getaccount] test22 is not defined\n", buffer.String())
}

func TestJob_GetLogger(t *testing.T) {
	j := NewJob("dummy", "foo_bar", make(map[string]interface{}), nil, false)
	logger := j.GetLogger()
	require.NotNil(t, logger)
	require.Equal(t, "mailAssistant.actions.newDummy%dummy", logger.Name())
}

func TestJob_getSaveTo(t *testing.T) {
	const (
		w = "bar\\foo"
		l = "foo/bar"
	)
	args := make(map[string]interface{})
	args["saveToWin"] = w
	args["saveTo"] = l

	j := NewJob("", "", args, nil, false)
	saveToV := j.getSaveTo()
	if runtime.GOOS == "windows" {
		require.Equal(t, w, saveToV)
		require.Equal(t, w, j.getSaveTo())
	} else {
		require.Equal(t, l, saveToV)
		require.Equal(t, l, j.getSaveTo())
	}

	j.saveTo = ""
	delete(args, "saveToWin")
	saveToV = j.getSaveTo()
	require.Equal(t, l, saveToV)
}

func TestJobGetSearchParameter(t *testing.T) {
	args := make(map[string]interface{})
	args["mail_account"] = "ignore"



	searchArg := make([]interface{},4)
	searchArg[0] = map[string] interface{} {"field":"ALL"}
	searchArg[1] = map[string] interface{} {"field":"CC", "value": "yang"}
	searchArg[2] = map[string] interface{} {"field":"older", "value": "every 5 seconds"}
	searchArg[3] = map[string] interface{} {"field":"or", "value": []interface{}{map[string]interface{}{"field": "from", "value": "foo@bar.org"}}}
	args["search"] = searchArg
	j := NewJob("dummy", "foo_bar", args, nil, false)
	search := j.getSearchParameter()

	require.NotNil(t, search)
	require.Len(t, search, 8)
	require.Equal(t, "ALL", search[0])
	require.Equal(t, "CC", search[1])
	require.Equal(t, "yang", search[2])
	require.Equal(t, "BEFORE", search[3])
	require.Equal(t, "or", search[5])
	require.Equal(t, "FROM", search[6])
	require.Equal(t, "foo@bar.org", search[7])
}
