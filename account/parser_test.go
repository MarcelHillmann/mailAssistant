package account

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"log"
	"mailAssistant/logging"
	"os"
	"testing"
)

func TestAccountParser(t *testing.T){
	t.Run("failOpen", accountParserFailOpen)
	t.Run("failNotYaml", accountParserFailNotYaml)
}

func accountParserFailOpen(t *testing.T){
	defer func(){
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags)
		logging.SetLevel("mailAssistant.account.parseYaml","OFF")
	} ()
	logging.SetLevel("mailAssistant.account.parseYaml","all")
	buffer := bytes.NewBufferString("")
	log.SetFlags(0)
	log.SetOutput(buffer)

	aux, err := parseYaml("","notExists.yml")
	require.NotNil(t, err)
	require.Nil(t, aux)

	str := buffer.String()
	require.Condition(t, func()bool {
		return str == "SEVERE  [mailAssistant.account.parseYaml#lambda] open notExists.yml: The system cannot find the file specified.\n" ||
			str == "SEVERE  [mailAssistant.account.parseYaml#lambda] open notExists.yml: no such file or directory\n"
	}, str)
}

func accountParserFailNotYaml(t *testing.T){
	defer func(){
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags)
	} ()
	buffer := bytes.NewBufferString("")
	log.SetFlags(0)
	log.SetOutput(buffer)

	aux, err := parseYaml("","notExists.xxx")
	require.Nil(t, err)
	require.Nil(t, aux)
}

func TestParserFailedReadAll(t *testing.T){
	defer func() {
		parserReadAll = ioutil.ReadAll
		err := recover()
		require.NotNil(t, err)
		require.EqualError(t, err.(error), "must fail")
	}()
	parserReadAll = func(r io.Reader) (i []byte, err error) {
		return []byte{}, errors.New("must fail")
	}

	parseYaml("","../testdata/accounts/muster@testcase.local.yml")
	require.Fail(t, "never called")
}