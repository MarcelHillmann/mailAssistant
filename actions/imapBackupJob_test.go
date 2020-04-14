package actions

import (
	"crypto/tls"
	"errors"
	"github.com/emersion/go-imap"
	"github.com/stretchr/testify/require"
	"mailAssistant/account"
	"mailAssistant/arguments"
	"mailAssistant/logging"
	"testing"
	"time"
)

func TestImapBackupJobLocked(t *testing.T) {
	var wg int32 = 1
	newImapBackup(Job{Logger: logging.NewLogger()}, &wg, metricsDummy)
	require.Equal(t, Locked, wg)
}

func TestImapBackupJobSuccess(t *testing.T) {
	defer account.SetClientFactory(nil)

	sClient, tClient := sourceClient(t), targetClient(t)
	account.SetClientFactory(func(addr string, tlsConfig *tls.Config) (account.IClient, error) {
		require.NotEmpty(t, addr)
		require.NotNil(t, tlsConfig)
		var mock *account.MockClientPromise
		if addr == "bar.foo:20000" {
			require.False(t, tlsConfig.InsecureSkipVerify)
			mock = sClient
		} else {
			require.True(t, tlsConfig.InsecureSkipVerify)
			mock = tClient
		}
		return mock, nil
	})

	logging.SetLevel("unit", "all")
	job := Job{}
	job.Logger = logging.NewNamedLogger("unit.tests")
	job.Accounts = new(account.Accounts)
	job.Account = make(map[string]account.Account)
	job.Account["foo bar"] = account.NewAccountForTest(t, "foo bar", "foo", "bar", "bar.foo", false)
	job.Account["foo bar target"] = account.NewAccountForTest(t, "foo bar target", "foo", "bar", "target.local", true)
	job.Args = arguments.NewEmptyArgs()
	job.SetArg("mail_account", "foo bar")
	job.SetArg("target_account", "foo bar target")
	job.SetArg("path", "INBOX/foo/bar")
	job.SetArg("search",[]interface{}{})

	var wg int32
	newImapBackup(job, &wg, metricsDummy)
	require.Equal(t, Released, wg)
	require.Equal(t, "10110-00100-311", sClient.Assert())
	require.Equal(t,"10000-00030-001", tClient.Assert())
}

func TestImapBackupJobFailedAccountMissed(t *testing.T) {
	defer account.SetClientFactory(nil)

	account.SetClientFactory(func(addr string, tlsConfig *tls.Config) (account.IClient, error) {
		require.Fail(t, "never call this")
		return nil, nil
	})

	logging.SetLevel("unit", "all")
	job := Job{}
	job.Logger = logging.NewNamedLogger("unit.tests")
	job.Accounts = new(account.Accounts)
	job.Account = make(map[string]account.Account)
	job.Account["foo bar"] = account.NewAccountForTest(t, "foo bar", "foo", "bar", "bar.foo", false)
	job.Args = arguments.NewEmptyArgs()
	job.Args.SetArg("mail_account", "foo bar")
	job.Args.SetArg("target_account", "foo bar target")
	job.Args.SetArg("path", "INBOX/foo/bar")

	var wg int32
	newImapBackup(job, &wg, metricsDummy)
	require.Equal(t, Released, wg)
}

func TestImapBackupJobFailedLogin(t *testing.T) {
	defer account.SetClientFactory(nil)

	sClient, tClient := sourceClient(t), targetClient(t)

	account.SetClientFactory(func(addr string, tlsConfig *tls.Config) (account.IClient, error) {
		require.NotEmpty(t, addr)
		require.NotNil(t, tlsConfig)
		var mock *account.MockClientPromise
		if addr == "bar.foo:20000" {
			require.False(t, tlsConfig.InsecureSkipVerify)
			mock = sClient
		} else {
			require.True(t, tlsConfig.InsecureSkipVerify)
			mock = tClient
			mock.LoginCallback = func(username, password string) error {
				return errors.New("fail")
			}
		}
		return mock, nil
	})

	logging.SetLevel("unit", "all")
	job := Job{}
	job.Logger = logging.NewNamedLogger("unit.tests")
	job.Accounts = new(account.Accounts)
	job.Account = make(map[string]account.Account)
	job.Account["foo bar"] = account.NewAccountForTest(t, "foo bar", "foo", "bar", "bar.foo", false)
	job.Account["foo bar target"] = account.NewAccountForTest(t, "foo bar target", "foo", "bar", "target.local", true)
	job.Args = arguments.NewEmptyArgs()
	job.Args.SetArg("mail_account", "foo bar")
	job.Args.SetArg("target_account", "foo bar target")
	job.Args.SetArg("path", "INBOX/foo/bar")

	var wg int32
	newImapBackup(job, &wg, metricsDummy)
	require.Equal(t, Released, wg)
	require.Equal(t, "10000-00000-001", sClient.Assert())
	require.Equal(t, "10000-00000-001", tClient.Assert())
}

func TestImapBackupJobFailedSelect(t *testing.T) {
	defer account.SetClientFactory(nil)

	sClient, tClient := sourceClient(t), targetClient(t)
	account.SetClientFactory(func(addr string, tlsConfig *tls.Config) (account.IClient, error) {
		require.NotEmpty(t, addr)
		require.NotNil(t, tlsConfig)
		var mock *account.MockClientPromise
		if addr == "bar.foo:20000" {
			require.False(t, tlsConfig.InsecureSkipVerify)
			mock = sourceClient(t)
			mock.SelectCallback = func(name string, readOnly bool) (status *imap.MailboxStatus, err error) {
				return nil, errors.New("fail")
			}
		} else {
			require.True(t, tlsConfig.InsecureSkipVerify)
			mock = targetClient(t)
		}
		return mock, nil
	})

	logging.SetLevel("*", "")
	logging.SetLevel("global", "all")
	job := Job{}
	job.Logger = logging.NewNamedLogger("unit.tests")
	job.Accounts = new(account.Accounts)
	job.Account = make(map[string]account.Account)
	job.Account["foo bar"] = account.NewAccountForTest(t, "foo bar", "foo", "bar", "bar.foo", false)
	job.Account["foo bar target"] = account.NewAccountForTest(t, "foo bar target", "foo", "bar", "target.local", true)
	job.Args = arguments.NewEmptyArgs()
	job.Args.SetArg("mail_account", "foo bar")
	job.Args.SetArg("target_account", "foo bar target")
	job.Args.SetArg("path", "INBOX/foo/bar")

	var wg int32
	newImapBackup(job, &wg, metricsDummy)
	require.Equal(t, Released, wg)
	require.Equal(t, "00000-00000-000", sClient.Assert())
	require.Equal(t, "00000-00000-000", tClient.Assert())
}

func TestImapBackupJobFailedSearchEmpty(t *testing.T) {
	defer account.SetClientFactory(nil)

	sClient, tClient := sourceClient(t), targetClient(t)
	account.SetClientFactory(func(addr string, tlsConfig *tls.Config) (account.IClient, error) {
		require.NotEmpty(t, addr)
		require.NotNil(t, tlsConfig)
		var mock *account.MockClientPromise
		if addr == "bar.foo:20000" {
			require.False(t, tlsConfig.InsecureSkipVerify)
			mock = sourceClient(t)
			mock.SearchCallback = func(criteria *imap.SearchCriteria) (uint32s []uint32, err error) {
				require.NotNil(t, criteria)
				return []uint32{}, nil
			}
		} else {
			require.True(t, tlsConfig.InsecureSkipVerify)
			mock = targetClient(t)
		}
		return mock, nil
	})

	logging.SetLevel("unit", "all")
	job := Job{}
	job.Logger = logging.NewNamedLogger("unit.tests")
	job.Accounts = new(account.Accounts)
	job.Account = make(map[string]account.Account)
	job.Account["foo bar"] = account.NewAccountForTest(t, "foo bar", "foo", "bar", "bar.foo", false)
	job.Account["foo bar target"] = account.NewAccountForTest(t, "foo bar target", "foo", "bar", "target.local", true)
	job.Args = arguments.NewEmptyArgs()
	job.Args.SetArg("mail_account", "foo bar")
	job.Args.SetArg("target_account", "foo bar target")
	job.Args.SetArg("path", "INBOX/foo/bar")

	var wg int32
	newImapBackup(job, &wg, metricsDummy)
	require.Equal(t, Released, wg)
	require.Equal(t, "00000-00000-000", sClient.Assert())
	require.Equal(t, "00000-00000-000", tClient.Assert())
}

func TestImapBackupJobFailedDelete(t *testing.T) {
	defer account.SetClientFactory(nil)

	sClient, tClient := sourceClient(t), targetClient(t)
	account.SetClientFactory(func(addr string, tlsConfig *tls.Config) (account.IClient, error) {
		require.NotEmpty(t, addr)
		require.NotNil(t, tlsConfig)
		var mock *account.MockClientPromise
		if addr == "bar.foo:20000" {
			require.False(t, tlsConfig.InsecureSkipVerify)
			mock = sourceClient(t)
			mock.DeleteCallback = func(num uint32) error {
				return errors.New("fail")
			}
		} else {
			require.True(t, tlsConfig.InsecureSkipVerify)
			mock = targetClient(t)
		}
		return mock, nil
	})

	logging.SetLevel("unit", "all")
	job := Job{}
	job.Logger = logging.NewNamedLogger("unit.tests")
	job.Accounts = new(account.Accounts)
	job.Account = make(map[string]account.Account)
	job.Account["foo bar"] = account.NewAccountForTest(t, "foo bar", "foo", "bar", "bar.foo", false)
	job.Account["foo bar target"] = account.NewAccountForTest(t, "foo bar target", "foo", "bar", "target.local", true)
	job.Args = arguments.NewEmptyArgs()
	job.Args.SetArg("mail_account", "foo bar")
	job.Args.SetArg("target_account", "foo bar target")
	job.Args.SetArg("path", "INBOX/foo/bar")

	var wg int32
	newImapBackup(job, &wg, metricsDummy)
	require.Equal(t,Released, wg)
	require.Equal(t, "00000-00000-000", sClient.Assert())
	require.Equal(t, "00000-00000-000", tClient.Assert())
}

func TestImapBackupJobFailedPanicUnlock(t *testing.T) {
	defer func() {
		account.SetClientFactory(nil)
		err := recover()
		require.Nil(t, err)
	}()

	sClient, tClient := sourceClient(t), targetClient(t)
	account.SetClientFactory(func(addr string, tlsConfig *tls.Config) (account.IClient, error) {
		require.NotEmpty(t, addr)
		require.NotNil(t, tlsConfig)
		var mock *account.MockClientPromise
		if addr == "bar.foo:20000" {
			require.False(t, tlsConfig.InsecureSkipVerify)
			mock = sourceClient(t)
		} else {
			require.True(t, tlsConfig.InsecureSkipVerify)
			mock = targetClient(t)
		}
		return mock, nil
	})

	logging.SetLevel("unit", "all")
	job := Job{}
	job.Logger = logging.NewNamedLogger("unit.tests")
	job.Accounts = new(account.Accounts)
	job.Account = make(map[string]account.Account)
	job.Account["foo bar"] = account.NewAccountForTest(t, "foo bar", "foo", "bar", "bar.foo", false)
	job.Account["foo bar target"] = account.NewAccountForTest(t, "foo bar target", "foo", "bar", "target.local", true)
	job.Args = arguments.NewEmptyArgs()
	job.Args.SetArg("mail_account", "foo bar")
	job.Args.SetArg("target_account", "foo bar target")
	job.Args.SetArg("path", "INBOX/foo/bar")

	var wg int32
	newImapBackup(job, &wg, metricsDummy)
	require.Equal(t,Released, wg)
	require.Equal(t, "00000-00000-000", sClient.Assert())
	require.Equal(t, "00000-00000-000", tClient.Assert())

}

func sourceClient(t *testing.T) *account.MockClientPromise {
	mock := account.NewMockClientMinimal()
	mock.DeleteCallback = func(num uint32) error {
		if num <10 || num > 12 {
			require.Fail(t, "invalid mail number")
		}
		return nil
	}
	mock.ExpungeCallback = func(ch chan uint32) error {
		require.Nil(t, ch)
		return nil
	}
	mock.LoginCallback = func(u, p string) error {
		require.Equal(t, "foo", u)
		require.Equal(t, "bar", p)
		return nil
	}
	mock.LogoutCallback = func() error {
		return nil
	}
	mock.SelectCallback = func(name string, readOnly bool) (*imap.MailboxStatus, error) {
		require.Equal(t, "INBOX.foo.bar", name)
		require.False(t, readOnly)
		return new(imap.MailboxStatus), nil
	}
	mock.SearchCallback = func(criteria *imap.SearchCriteria) ([]uint32, error) {
		require.NotNil(t, criteria)
		return []uint32{10, 11, 12}, nil
	}
	mock.FetchCallback = func(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
		require.NotNil(t, seqSet)
		require.NotNil(t, items)
		require.Len(t, items, 9)
		require.NotNil(t, ch)

		ch <- createMessage(10, true)
		ch <- createMessage(11, true)
		ch <- createMessage(12, true)

		close(ch)
		return nil
	}
	return mock
}

func targetClient(t *testing.T) *account.MockClientPromise {
	mock := account.NewMockClientMinimal()
	mock.LoginCallback = func(u, p string) error {
		require.Equal(t, "foo", u)
		require.Equal(t, "bar", p)
		return nil
	}
	mock.LogoutCallback = func() error {
		return nil
	}
	mock.AppendCallback = func(mBox string, flags []string, date time.Time, msg imap.Literal) error {
		require.Equal(t, "INBOX", mBox)
		return nil
	}
	return mock
}
