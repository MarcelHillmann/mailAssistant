package account

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/require"
	"mailAssistant/cntl"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAccounts(t *testing.T) {
	t.Run("newAccount", accountsNewAccounts)
	t.Run("loadFromDisk", accountsLoadFromDisk)
	t.Run("GetAccount", accountsGetAccount)
	t.Run("HasAccount", accountsHasAccount)
	t.Run("importAccounts", func(t *testing.T) {
		t.Run("parseYaml", func(t *testing.T) {
			t.Run("Failed", accountsImportAccountsParseYamlFailed)
			t.Run("Empty", accountsImportAccountsParseYamlEmpty)
		})
		t.Run("Create", accountsImportAccountsCreate)
		t.Run("Modified", func(t *testing.T) {
			t.Run("existing", accountsImportAccountsModifiedExisting)
			t.Run("notExisting", accountsImportAccountsModifiedNotExisting)
		})
		t.Run("Removed", func(t *testing.T) {
			t.Run("existing", accountsImportAccountsRemovedExisting)
			t.Run("notExisting", accountsImportAccountsRemovedNotExisting)
		})
	})
	t.Run("startWatcher", func(t *testing.T) {
		t.Run("works", accountsStartWatcherWorks)
		t.Run("receive error", accountsStartWatcherError)
		t.Run("panic", accountsStartWatcherPanic)
	})
	t.Run("ImportAccounts", func(t *testing.T) {
		t.Run("OK", accountsImportAccountsOk)
		t.Run("ImportError", accountsImportAccountsImportError)
		t.Run("Error", accountsImportAccountsError)
	})

}

func accountsImportAccountsError(t *testing.T) {
	defer SetAccountImportPath("")
	SetAccountImportPath("")
	_, err := ImportAccounts()
	require.NotNil(t, err)
	sErr := err.Error()
	require.Condition(t, func() bool {
		return sErr == "open resources\\config\\accounts: The system cannot find the path specified." ||
			sErr == "open resources/config/accounts: no such file or directory"
	}, sErr)
}

func accountsImportAccountsImportError(t *testing.T) {
	defer SetAccountImportPath("")
	SetAccountImportPath("")
	importError = errors.New("import error")
	_, err := ImportAccounts()
	require.NotNil(t, err)
	require.Equal(t, "import error", err.Error())
}

func accountsImportAccountsOk(t *testing.T) {
	defer SetAccountImportPath("")
	SetAccountImportPath("../testdata/accounts/")

	acc, err := ImportAccounts()
	require.Nil(t, err)
	time.Sleep(500 * time.Millisecond)
	require.Equal(t, 1, cntl.ToNotify())
	cntl.Notify()
	require.Equal(t, 0, cntl.ToNotify())
	require.Equal(t, "../testdata/accounts/", acc.accountDir)
	require.Len(t, acc.Account, 1)
	require.Len(t, acc.files, 1)
}

func accountsStartWatcherWorks(t *testing.T) {
	defer func() {
		err := recover()
		require.Nil(t, err, "has not to panic")
	}()
	acc := newAccounts()
	acc.accountDir, _ = filepath.Abs("../testdata")

	wait := make(chan bool, 1)
	go func() {
		acc.startWatcher()
		wait <- true
	}()
	time.Sleep(500 * time.Millisecond)
	require.Equal(t, 1, cntl.ToNotify())
	cntl.Notify()
	<-wait
	require.Equal(t, 0, cntl.ToNotify())
}

func accountsStartWatcherError(t *testing.T) {
	defer func() {
		SetAccountWalker(nil)
		require.Nil(t, recover(), "has not to panic")
	}()
	acc := newAccounts()
	acc.accountDir, _ = filepath.Abs("../testdata")

	SetAccountWalker(func(watcher *fsnotify.Watcher) filepath.WalkFunc {
		return func(path string, info os.FileInfo, err error) error {
			watcher.Errors <- errors.New("test")
			return nil
		}
	})
	wait := make(chan bool, 1)
	go func() {
		acc.startWatcher()
		wait <- true
	}()
	time.Sleep(500 * time.Millisecond)
	require.Equal(t, 1, cntl.ToNotify())
	cntl.Notify()
	<-wait
	require.Equal(t, 0, cntl.ToNotify())
}

func accountsStartWatcherPanic(t *testing.T) {
	defer func() {
		err := recover()
		require.NotNil(t, err, "must panic")
		sErr := err.(error).Error()
		require.Condition(t, func() bool {
			return sErr == "startWalker - ERROR lstat : no such file or directory" ||
				sErr == "startWalker - ERROR Lstat : The system cannot find the path specified."
		}, sErr)
		require.Equal(t, 0, cntl.ToNotify())
	}()
	acc := newAccounts()
	acc.startWatcher()
	require.Fail(t, "never run")
}

func accountsNewAccounts(t *testing.T) {
	acc := newAccounts()
	require.Equal(t, "", acc.accountDir)
	require.NotNil(t, acc.Account)
	require.Len(t, acc.Account, 0)
	require.NotNil(t, acc.files)
	require.Len(t, acc.files, 0)
}

func accountsLoadFromDisk(t *testing.T) {
	acc := newAccounts()
	acc.accountDir, _ = filepath.Rel(".", "../testdata/accounts/")
	_ = acc.loadFromDisk()

	require.False(t, acc.HasAccount("test"))
	require.Nil(t, acc.GetAccount("test"))
	require.True(t, acc.HasAccount("muster@testcase.local"))
	require.NotNil(t, acc.GetAccount("muster@testcase.local"))
}

func accountsGetAccount(t *testing.T) {
	account := make(map[string]Account)
	account["test"] = Account{"test", "marcel", "geheim", "localhost", 1000, true, true}
	files := make(map[string]string)

	acc := Accounts{"testData", account, files}
	require.Nil(t, acc.GetAccount("NIL"))
	require.NotNil(t, acc.GetAccount("test"))
}

func accountsHasAccount(t *testing.T) {
	account := make(map[string]Account)
	account["test"] = Account{"test", "marcel", "geheim", "localhost", 1000, true, true}
	files := make(map[string]string)

	acc := Accounts{"testData", account, files}
	require.False(t, acc.HasAccount("NIL"))
	require.True(t, acc.HasAccount("test"))
}

func accountsImportAccountsParseYamlFailed(t *testing.T) {
	defer func() {
		err := recover()
		require.NotNil(t, err)
		e := err.(error).Error()
		if e == "ERROR yaml: line 7: could not find expected ':'" ||  e == "ERROR yaml: line 8: could not find expected ':'" {
			// passed
		} else {
			require.Fail(t, "Invalid error message: %s", e)
		}
	}()

	acc := newAccounts()
	acc.importAccount("", "../testdata/accounts_invalid/error.yml", fsnotify.Create)
	require.Fail(t, "never run this")
}

func accountsImportAccountsParseYamlEmpty(t *testing.T) {
	defer func() {
		err := recover()
		require.Nil(t, err)
	}()

	acc := newAccounts()
	acc.importAccount("", "../testdata/accounts_invalid/empty.yml", fsnotify.Create)
	require.True(t, true)
}

func accountsImportAccountsCreate(t *testing.T) {
	acc := newAccounts()
	acc.importAccount("", "../testdata/accounts/muster@testcase.local.yml", fsnotify.Create)

	require.Len(t, acc.files, 1)
	require.Len(t, acc.Account, 1)
	require.NotNil(t, acc.files["../testdata/accounts/muster@testcase.local.yml"])
	require.Equal(t, "muster@testcase.local", acc.files["../testdata/accounts/muster@testcase.local.yml"])
	require.NotNil(t, acc.Account["muster@testcase.local"])
	require.True(t, acc.HasAccount("muster@testcase.local"))
	require.NotNil(t, acc.GetAccount("muster@testcase.local"))

	account := acc.GetAccount("muster@testcase.local")
	require.Equal(t, "muster@testcase.local", account.name)
	require.Equal(t, "foo", account.username)
	require.Equal(t, "bar", account.password)
	require.Equal(t, "localhost", account.hostname)
	require.Equal(t, 993, account.port)
	require.Equal(t, "localhost:993", account.hostnamePort())
	require.False(t, account.debug)
}

func accountsImportAccountsModifiedExisting(t *testing.T) {
	acc := newAccounts()
	acc.importAccount("", "../testdata/accounts/muster@testcase.local.yml", fsnotify.Create)
	require.Len(t, acc.files, 1)
	require.Len(t, acc.Account, 1)

	acc.importAccount("", "../testdata/accounts/muster@testcase.local.yml", fsnotify.Write)
	require.Len(t, acc.files, 1)
	require.Len(t, acc.Account, 1)
	require.NotNil(t, acc.files["../testdata/accounts/muster@testcase.local.yml"])
	require.Equal(t, "muster@testcase.local", acc.files["../testdata/accounts/muster@testcase.local.yml"])
	require.NotNil(t, acc.Account["muster@testcase.local"])
	require.True(t, acc.HasAccount("muster@testcase.local"))
	require.NotNil(t, acc.GetAccount("muster@testcase.local"))

	account := acc.GetAccount("muster@testcase.local")
	require.Equal(t, "muster@testcase.local", account.name)
	require.Equal(t, "foo", account.username)
	require.Equal(t, "bar", account.password)
	require.Equal(t, "localhost", account.hostname)
	require.Equal(t, 993, account.port)
	require.Equal(t, "localhost:993", account.hostnamePort())
	require.False(t, account.debug)
}

func accountsImportAccountsModifiedNotExisting(t *testing.T) {
	acc := newAccounts()
	acc.importAccount("", "../testdata/accounts/muster@testcase.local.yml", fsnotify.Write)
	require.Len(t, acc.files, 1)
	require.Len(t, acc.Account, 1)
	require.NotNil(t, acc.files["../testdata/accounts/muster@testcase.local.yml"])
	require.Equal(t, "muster@testcase.local", acc.files["../testdata/accounts/muster@testcase.local.yml"])
	require.NotNil(t, acc.Account["muster@testcase.local"])
	require.True(t, acc.HasAccount("muster@testcase.local"))
	require.NotNil(t, acc.GetAccount("muster@testcase.local"))

	account := acc.GetAccount("muster@testcase.local")
	require.Equal(t, "muster@testcase.local", account.name)
	require.Equal(t, "foo", account.username)
	require.Equal(t, "bar", account.password)
	require.Equal(t, "localhost", account.hostname)
	require.Equal(t, 993, account.port)
	require.Equal(t, "localhost:993", account.hostnamePort())
	require.False(t, account.debug)
}

func accountsImportAccountsRemovedExisting(t *testing.T) {
	acc := newAccounts()
	acc.importAccount("", "../testdata/accounts/muster@testcase.local.yml", fsnotify.Create)
	require.Len(t, acc.files, 1)
	require.Len(t, acc.Account, 1)

	acc.importAccount("", "../testdata/accounts/muster@testcase.local.yml", fsnotify.Remove)
	require.Len(t, acc.files, 0)
	require.Len(t, acc.Account, 0)
}

func accountsImportAccountsRemovedNotExisting(t *testing.T) {
	acc := newAccounts()
	acc.importAccount("", "../testdata/accounts/muster@testcase.local.yml", fsnotify.Remove)
	require.Len(t, acc.files, 0)
	require.Len(t, acc.Account, 0)
}
