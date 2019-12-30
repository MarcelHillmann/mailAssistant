package rules

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"log"
	"mailAssistant/logging"
	"os"
	"strings"
	"testing"
)

func TestRulesParser(t *testing.T) {
	t.Run("failOpen", rulesParserFailOpen)
	t.Run("failNotYaml", rulesParserFailedNotYaml)
}

func rulesParserFailOpen(t *testing.T) {
	logging.SetLevel("mailAssistant.rules.parseYaml", "all")
	buffer := bytes.NewBufferString("")
	log.SetFlags(0)
	log.SetOutput(buffer)

	defer func() {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags)
		logging.SetLevel("mailAssistant.rules.parseYaml", "OFF")

		require.Condition(t, func() bool {
			return strings.HasPrefix(buffer.String(), "SEVERE  [mailAssistant.rules.parseYaml#lambda] open ../testdata/rules/notExists.yml:")
		}, buffer.String())
		err := recover()
		require.NotNil(t, err)
		sErr := err.(error).Error()
		require.Condition(t, func() bool {
			return strings.HasPrefix(sErr, "SEVERE  [mailAssistant.rules.parseYaml#lambda] open ../testdata/rules/notExists.yml:")
		}, sErr)
	}()

	_, _ = parseYaml("../testdata/rules", "", "../testdata/rules/notExists.yml")
	require.Fail(t, "never run")

}

func rulesParserFailedNotYaml(t *testing.T) {
	defer func() {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags)
	}()
	buffer := bytes.NewBufferString("")
	log.SetFlags(0)
	log.SetOutput(buffer)

	err, aux := parseYaml("../testdata/rules", "../testdata/rules/", "notExists.xxx")
	require.Nil(t, err)
	require.Nil(t, aux)
}

func TestParserFailedReadAll(t *testing.T) {
	defer func() {
		parserReadAll = ioutil.ReadAll
		err := recover()
		require.NotNil(t, err)
		require.EqualError(t, err.(error), "ReadAll must fail")
	}()
	parserReadAll = func(r io.Reader) (i []byte, err error) {
		return []byte{}, errors.New("must fail")
	}

	parseYaml("../testdata/rules", "../testdata/rules/", "fooBar.yml")
	require.Fail(t,"never call this")
}
