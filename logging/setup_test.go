package logging

import (
	"errors"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_ImportCnf_Ok(t *testing.T) {
	logCnf, err := filepath.Rel(".", "../testdata/logging/logging.yml")
	if err != nil {
		t.Fail()
	}

	loggerRegistry = make(map[string]logLevel)
	importCnf(logCnf)

	require.Equal(t, severe, loggerRegistry["global"])
	require.Equal(t, warn, loggerRegistry["mailAssistant"])
	require.Equal(t, all, loggerRegistry["mailAssistant.account"])
	require.Equal(t, warn, loggerRegistry["mailAssistant.actions"])
	require.Equal(t, debug, loggerRegistry["mailAssistant.actions.newImapMove"])
	require.Equal(t, info, loggerRegistry["mailAssistant.actions.newImapBackup"])
	require.Equal(t, severe, loggerRegistry["mailAssistant.rules"])
	require.Len(t, loggerRegistry, 7)
}

func Test_ImportCnf_FailedReadAll(t *testing.T) {
	logCnf, err := filepath.Rel(".", "../testdata/logging/logging.yml")
	if err != nil {
		t.Fail()
	}
	defer func() {
		setupReadAll = ioutil.ReadAll
	}()
	setupReadAll = func(r io.Reader) (bytes []byte, err error) {
		return []byte{}, errors.New("must fail")
	}
	loggerRegistry = make(map[string]logLevel)
	importCnf(logCnf)
	require.Len(t, loggerRegistry, 0)
}

func Test_ImportCnf_Failed(t *testing.T) {
	logCnf, err := filepath.Rel(".", "../testdata/logging/failed.yml")
	if err != nil {
		t.Fail()
	}

	loggerRegistry = make(map[string]logLevel)
	logReg := loggerRegistry
	importCnf(logCnf)
	require.True(t, reflect.ValueOf(logReg).Pointer() == reflect.ValueOf(loggerRegistry).Pointer())
	require.Len(t, loggerRegistry, 0)
}

func Test_ImportCnf_NoRootLogLevel(t *testing.T) {
	logCnf, err := filepath.Rel(".", "../testdata/logging/loggingNoLogLevel.yml")
	if err != nil {
		t.Fail()
	}

	logReg := loggerRegistry
	importCnf(logCnf)
	require.False(t, reflect.ValueOf(logReg).Pointer() == reflect.ValueOf(loggerRegistry).Pointer())
	require.Len(t, loggerRegistry, 2)
	require.Equal(t, none, loggerRegistry["global"])
	require.Equal(t, warn, loggerRegistry["mailAssistant"])
}

var testCases = [][]interface{}{
	{none, "NONE"},
	{all, "ALL"},
	{debug, "DEBUG"},
	{info, "INFO"},
	{warn, "WARN"},
	{severe, "SEVERE"},
	{notExists, "NOT EXISTS"},
}

func TestLogLevel(t *testing.T) {
	for _, testCase := range testCases {
		key := testCase[1].(string)
		t.Run(key, func(t *testing.T) {
			value := testCase[0].(logLevel)
			require.Equal(t, key, value.String())
		})
	}
}
