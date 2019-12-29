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
	j := NewJob("dummy", "foo_bar",make(map[string]interface{}),nil)
	require.NotNil(t, j)
}

func TestJob_Run(t *testing.T) {
	j := NewJob("dummy", "foo_bar",make(map[string]interface{}),nil)
	logging.SetLevel(j.log.Name(),"all")
	defer func() {
		logging.SetLevel(j.log.Name(),"none")
		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stderr)
	}()
	buffer := bytes.NewBufferString("")
	log.SetFlags(0)
	log.SetOutput(buffer)
	j.Run()
	require.NotEmpty(t, buffer.String())
	require.Equal(t, "DEBUG   [mailAssistant.actions.newDummy#mailassistant.actions.job.run] >>\n" +
		"DEBUG   [mailAssistant.actions.newDummy#lambda] map[]\n" +
		"DEBUG   [mailAssistant.actions.newDummy#mailassistant.actions.job.run] <<\n", buffer.String())
}

func TestJob_GetAccount(t *testing.T) {
	acc := account.Accounts{}
	acc.Account = make(map[string]account.Account)
	acc.Account["test"] = account.NewAccountForTest(t,"test","foo","bar","l",false)
	j := NewJob("dummy", "foo_bar",make(map[string]interface{}),&acc)
	logging.SetLevel(j.log.Name(),"all")
	defer func() {
		logging.SetLevel(j.log.Name(),"none")
		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stderr)
	}()
	buffer := bytes.NewBufferString("")
	log.SetFlags(0)
	log.SetOutput(buffer)
	require.NotNil(t, j.GetAccount("test"))
	require.Nil(t, j.GetAccount("test22"))
	require.NotEmpty(t, buffer.String())
	require.Equal(t, "SEVERE  [mailAssistant.actions.newDummy#mailassistant.actions.job.getaccount] test22 is not defined\n",buffer.String())
}

func TestJob_GetLogger(t *testing.T) {
	j := NewJob("dummy", "foo_bar",make(map[string]interface{}),nil)
	log := j.GetLogger()
	require.NotNil(t,log)
	require.Equal(t, "mailAssistant.actions.newDummy", log.Name())
}

func TestJob_getSaveTo(t *testing.T){
	const (
		w = "bar\\foo"
		l = "foo/bar"
	)
	args := make(map[string]interface{})
	args["saveToWin"] = w
	args["saveTo"] = l

	j := NewJob("","",args,nil)
	saveToV := j.getSaveTo()
	if runtime.GOOS == "windows" {
		require.Equal(t, w, saveToV)
		require.Equal(t, w, j.getSaveTo())
	}else{
		require.Equal(t, l, saveToV)
		require.Equal(t, l, j.getSaveTo())
	}

	j.saveTo=""
	delete(args, "saveToWin")
	saveToV = j.getSaveTo()
	require.Equal(t, l, saveToV)
}

func TestJobParseRecursive(t *testing.T){
	list := make([]interface{},0)
	list = append(list, map[string]interface{}{"field":"foo", "value":"bar"})
	assertMe := make([]interface{},0)
	result := parseRecursive(assertMe,list)
	require.Len(t, result,2)
	require.Equal(t, "foo", result[0])
	require.Equal(t, "bar", result[1])

	list[0].(map[string]interface{})["value"] = []interface{}{map[string]interface{}{"field":"bar","value":"b"}}
	result2 := parseRecursive(assertMe,list)
	require.Len(t, result2,3)
	require.Equal(t, "foo", result2[0])
	require.Equal(t, "bar", result2[1])
	require.Equal(t, "b", result2[2])
}

func TestJobGetSearchParameter(t *testing.T){
	args := make(map[string]interface{})
	args["mail_account"] ="ignore"
	args["mail_foo"] ="bar"
	args["mail_ying"] ="yang"
	args["mail_or"] = []interface{}{map[string]interface{}{"field":"foo", "value": "bar"}}
	args["mail_older"]= "every 5 seconds"
	j := NewJob("dummy", "foo_bar",args,nil)
	search := j.getSearchParameter()

	require.NotNil(t, search)
	require.Len(t, search,4)
	require.Equal(t, []interface{}{"foo","bar"}, search[0])
	require.Equal(t, "before", search[1][0])
	require.Equal(t, []interface{}{"or","foo","bar"}, search[2])
	require.Equal(t, []interface{}{"ying","yang"}, search[3])
}