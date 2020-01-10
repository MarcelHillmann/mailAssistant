package actions

import (
	"crypto/tls"
	"errors"
	"github.com/emersion/go-imap"
	"github.com/stretchr/testify/require"
	"mailAssistant/account"
	"mailAssistant/arguments"
	e "mailAssistant/errors"
	"mailAssistant/logging"
	"testing"
)

func TestImapMvJob_Locked(t *testing.T){
	var wg int32 = 1
	newImapMove(Job{log: logging.NewLogger()}, &wg, metricsDummy)
	require.Equal(t, Locked, wg)
}

func TestImapMvJobSuccess(t *testing.T){
	defer account.SetClientFactory(nil)

	mock := account.NewMockClient()
	account.SetClientFactory(func(addr string, tlsConfig *tls.Config) (account.IClient, error) {
		require.Equal(t, "bar.foo:20000", addr)
		require.NotNil(t,  tlsConfig)
		require.True(t, tlsConfig.InsecureSkipVerify)

		mock.LoginCallback = func(u,p string) error {
			require.Equal(t, "foo", u)
			require.Equal(t, "bar", p)
			return nil
		}
		mock.SelectCallback = func(name string, readOnly bool) (*imap.MailboxStatus, error) {
			require.Equal(t, "INBOX.foo.bar",name)
			require.False(t, readOnly)
			return new(imap.MailboxStatus), nil
		}
		mock.SearchCallback = func(criteria *imap.SearchCriteria) ([]uint32,error) {
			require.NotNil(t, criteria)
			return []uint32{10,11,12}, nil
		}
		mock.FetchCallback = func(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
			require.NotNil(t, seqSet)
			require.NotNil(t, items)
			require.Len(t, items,3)
			require.NotNil(t, ch)

			ch <- createMessage(10, false)
			ch <- createMessage(11, false)
			ch <- createMessage(12, false)

			close(ch)
			return nil
		}
		mock.MoveCallback = func(seqSet *imap.SeqSet, dest string ) error {
			require.NotNil(t, seqSet)
			require.Equal(t, "INBOX", dest)
			if seqSet.Set[0].Start < 10 || seqSet.Set[0].Stop > 12 {
				require.Fail(t, "invalid seqSet")
			}
			return nil
		}

		return mock, nil
	})

	job:= Job{log: logging.NewLogger()}
	job.accounts = new(account.Accounts)
	job.accounts.Account = make(map[string]account.Account)
	job.accounts.Account["foo bar"] = account.NewAccountForTest(t,"foo bar", "foo","bar","bar.foo", true)
	job.Args = arguments.NewEmptyArgs()
	job.SetArg("mail_account", "foo bar")
	job.SetArg("path", "INBOX/foo/bar")
	job.SetArg("search",[]interface{}{})
	var wg int32
	newImapMove(job, &wg, metricsDummy)
	require.Equal(t, Released,wg)
	require.Equal(t,"10111-00100-001", mock.Assert())
}

func TestImapMvJobFailedLogin(t *testing.T){
	defer account.SetClientFactory(nil)

	mock := account.NewMockClient()
	account.SetClientFactory(func(addr string, tlsConfig *tls.Config) (account.IClient, error) {
		require.Equal(t, "bar.foo:20000", addr)
		require.NotNil(t,  tlsConfig)
		require.True(t, tlsConfig.InsecureSkipVerify)

		mock.LoginCallback = func(u,p string) error {
			require.Equal(t, "foo", u)
			require.Equal(t, "bar", p)
			return errors.New("Login failed")
		}
		mock.SelectCallback = func(name string, readOnly bool) (*imap.MailboxStatus, error) {
			require.Fail(t, "never call this")
			return new(imap.MailboxStatus), nil
		}
		mock.SearchCallback = func(criteria *imap.SearchCriteria) ([]uint32,error) {
			require.Fail(t, "never call this")
			return []uint32{}, nil
		}
		mock.FetchCallback = func(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
			require.Fail(t, "never call this")
			close(ch)
			return nil
		}
		mock.MoveCallback = func(seqSet *imap.SeqSet, dest string ) error {
			return nil
		}

		return mock, nil
	})

	job:= Job{log: logging.NewLogger()}
	job.accounts = new(account.Accounts)
	job.accounts.Account = make(map[string]account.Account)
	job.accounts.Account["foo bar"] = account.NewAccountForTest(t,"foo bar", "foo","bar","bar.foo",  true)
	job.Args = arguments.NewEmptyArgs()
	job.Args.SetArg("mail_account", "foo bar")
	job.Args.SetArg("path", "INBOX/foo/bar")

	var wg int32
	newImapMove(job, &wg, metricsDummy)
	require.Equal(t,Released, wg)
	require.Equal(t, "10000-00000-001", mock.Assert())
}

func TestImapMvJobFailedSelect(t *testing.T){
	defer account.SetClientFactory(nil)

	mock := account.NewMockClient()
	account.SetClientFactory(func(addr string, tlsConfig *tls.Config) (account.IClient, error) {
		require.Equal(t, "bar.foo:20000", addr)
		require.NotNil(t, tlsConfig)
		require.True(t, tlsConfig.InsecureSkipVerify)

		mock.LoginCallback = func(u,p string) error {
			require.Equal(t, "foo", u)
			require.Equal(t, "bar", p)
			return nil
		}
		mock.SelectCallback = func(name string, readOnly bool) (*imap.MailboxStatus, error) {
			require.Equal(t, "INBOX.foo.bar",name)
			require.False(t, readOnly)
			return nil, errors.New("select failed")
		}
		mock.SearchCallback = func(criteria *imap.SearchCriteria) ([]uint32,error) {
			require.Fail(t, "never call this")
			return []uint32{}, nil
		}
		mock.FetchCallback = func(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
			require.Fail(t, "never call this")
			close(ch)
			return nil
		}
		mock.MoveCallback = func(seqSet *imap.SeqSet, dest string ) error {
			return nil
		}

		return mock, nil
	})

	job:= Job{log: logging.NewLogger()}
	job.accounts = new(account.Accounts)
	job.accounts.Account = make(map[string]account.Account)
	job.accounts.Account["foo bar"] = account.NewAccountForTest(t,"foo bar", "foo","bar","bar.foo",  true)
	job.Args = arguments.NewEmptyArgs()
	job.Args.SetArg("mail_account", "foo bar")
	job.Args.SetArg("path", "INBOX/foo/bar")

	var wg int32
	newImapMove(job, &wg, metricsDummy)
	require.Equal(t, Released, wg)
	require.Equal(t, "10100-00000-001", mock.Assert())
}

func TestImapMvJobFailedStoreEmpty(t *testing.T){
	defer account.SetClientFactory(nil)

	mock := account.NewMockClient()
	account.SetClientFactory(func(addr string, tlsConfig *tls.Config) (account.IClient, error) {
		require.Equal(t, "bar.foo:20000", addr)
		require.NotNil(t,  tlsConfig)
		require.True(t, tlsConfig.InsecureSkipVerify)

		mock.LoginCallback = func(u,p string) error {
			require.Equal(t, "foo", u)
			require.Equal(t, "bar", p)
			return nil
		}
		mock.SelectCallback = func(name string, readOnly bool) (*imap.MailboxStatus, error) {
			require.Equal(t, "INBOX.foo.bar",name)
			require.False(t, readOnly)
			return new(imap.MailboxStatus), nil
		}
		mock.SearchCallback = func(criteria *imap.SearchCriteria) ([]uint32,error) {
			require.NotNil(t, criteria)
			return []uint32{}, nil
		}
		mock.FetchCallback = func(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
			require.Fail(t, "never call this")
			close(ch)
			return nil
		}
		mock.MoveCallback = func(seqSet *imap.SeqSet, dest string ) error {
			return nil
		}

		return mock, nil
	})

	job:= Job{log: logging.NewLogger()}
	job.accounts = new(account.Accounts)
	job.accounts.Account = make(map[string]account.Account)
	job.accounts.Account["foo bar"] = account.NewAccountForTest(t,"foo bar", "foo","bar","bar.foo",  true)
	job.Args = arguments.NewEmptyArgs()
	job.SetArg("mail_account", "foo bar")
	job.SetArg("path", "INBOX/foo/bar")
	job.SetArg("search",[]interface{}{})

	var wg  int32
	newImapMove(job, &wg, metricsDummy)
	require.Equal(t,Released, wg)
	require.Equal(t, "10110-00000-001", mock.Assert())
}

func TestImapMvJobNotLockedEmpty(t *testing.T){
	defer account.SetClientFactory(nil)

	mock := account.NewMockClient()
	account.SetClientFactory(func(addr string, tlsConfig *tls.Config) (account.IClient, error) {
		require.Equal(t, "bar.foo:20000", addr)
		require.NotNil(t, tlsConfig)
		require.True(t, tlsConfig.InsecureSkipVerify)

		mock.LoginCallback = func(u,p string) error {
			require.Equal(t, "foo", u)
			require.Equal(t, "bar", p)
			return nil
		}
		mock.SelectCallback = func(name string, readOnly bool) (*imap.MailboxStatus, error) {
			require.Equal(t, "INBOX.foo.bar",name)
			require.False(t, readOnly)
			return new(imap.MailboxStatus), nil
		}
		mock.SearchCallback = func(criteria *imap.SearchCriteria) ([]uint32,error) {
			require.NotNil(t, criteria)
			return []uint32{10,11,12}, nil
		}
		mock.FetchCallback = func(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
			require.NotNil(t, seqSet)
			require.NotNil(t, items)
			require.Len(t, items,3)
			require.NotNil(t, ch)

			ch <- createMessage(10, false)
			ch <- createMessage(11, false)
			ch <- createMessage(12, false)

			close(ch)
			return nil
		}
		mock.MoveCallback = func(seqSet *imap.SeqSet, dest string ) error {
			return e.NewEmpty()
		}

		return mock, nil
	})

	job:= Job{log: logging.NewLogger()}
	job.accounts = new(account.Accounts)
	job.accounts.Account = make(map[string]account.Account)
	job.accounts.Account["foo bar"] = account.NewAccountForTest(t,"foo bar", "foo","bar","bar.foo",  true)
	job.Args = arguments.NewEmptyArgs()
	job.SetArg("mail_account", "foo bar")
	job.SetArg("path", "INBOX/foo/bar")
	job.SetArg("search",[]interface{}{})

	var wg int32
	newImapMove(job, &wg, metricsDummy)
	require.Equal(t,Released, wg)
	require.Equal(t, "10111-00100-001", mock.Assert())
}

func TestImapMvJobFailedPanicUnlocked(t *testing.T){
	defer func() {
		account.SetClientFactory(nil)
		err := recover()
		require.Nil(t,err)
	}()

	mock := account.NewMockClient()
	account.SetClientFactory(func(addr string, tlsConfig *tls.Config) (account.IClient, error) {
		require.Equal(t, "bar.foo:20000", addr)
		require.NotNil(t,  tlsConfig)
		require.True(t, tlsConfig.InsecureSkipVerify)

		mock.LoginCallback = func(u,p string) error {
			require.Equal(t, "foo", u)
			require.Equal(t, "bar", p)
			return nil
		}
		mock.SelectCallback = func(name string, readOnly bool) (*imap.MailboxStatus, error) {
			require.Equal(t, "INBOX.foo.bar",name)
			require.False(t, readOnly)
			return new(imap.MailboxStatus), nil
		}
		mock.SearchCallback = func(criteria *imap.SearchCriteria) ([]uint32,error) {
			require.NotNil(t, criteria)
			return []uint32{10,11,12}, nil
		}
		mock.FetchCallback = func(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
			require.NotNil(t, seqSet)
			require.NotNil(t, items)
			require.Len(t, items,3)
			require.NotNil(t, ch)

			ch <- createMessage(10, false)
			ch <- createMessage(11, false)
			ch <- createMessage(12, false)

			close(ch)
			return nil
		}
		mock.MoveCallback = func(seqSet *imap.SeqSet, dest string) error {
			require.NotNil(t, seqSet)
			require.Equal(t, "INBOX.foo_bar",dest)
			if seqSet.Set[0].Start < 10 || seqSet.Set[0].Stop > 12 {
				require.Fail(t, "invalid seqSet")
			}
			return errors.New("let me panic")
		}

		return mock, nil
	})

	job:= Job{log: logging.NewLogger()}
	job.accounts = new(account.Accounts)
	job.accounts.Account = make(map[string]account.Account)
	job.accounts.Account["foo bar"] = account.NewAccountForTest(t,"foo bar", "foo","bar","bar.foo",  true)
	job.Args = arguments.NewEmptyArgs()
	job.SetArg("mail_account", "foo bar")
	job.SetArg("path", "INBOX/foo/bar")
	job.SetArg(moveTo, "INBOX/foo_bar")
	job.SetArg("search", []interface{}{})

	var wg int32
	newImapMove(job, &wg, metricsDummy)
	require.Equal(t,Released, wg)
	require.Equal(t, "10111-00100-001", mock.Assert())
}
