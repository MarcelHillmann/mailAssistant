package actions

import (
	"crypto/tls"
	"errors"
	"github.com/emersion/go-imap"
	"github.com/stretchr/testify/require"
	"mailAssistant/account"
	"mailAssistant/arguments"
	e "mailAssistant/errors"
	"testing"
)

func TestSeenJob_Locked(t *testing.T){
	var wg int32 = 1
	newSeenJob(Job{}, &wg)
}

func TestSeenJobSuccess(t *testing.T){
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
		mock.StoreCallback = func(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
			require.NotNil(t, seqSet)
			require.NotNil(t, item)
			require.NotNil(t, value)
			_v := value.([]interface{})
			require.Len(t, _v, 1)
			require.Equal(t, imap.SeenFlag, _v[0].(string))
			require.Nil(t, ch)
			if seqSet.Set[0].Start < 10 || seqSet.Set[0].Stop > 12 {
				require.Fail(t, "invalid seqSet")
			}
			return nil
		}

		return mock, nil
	})

	job:= Job{}
	job.accounts = new(account.Accounts)
	job.accounts.Account = make(map[string]account.Account)
	job.accounts.Account["foo bar"] = account.NewAccountForTest(t,"foo bar", "foo","bar","bar.foo", true)
	job.Args = arguments.NewEmptyArgs()
	job.Args.SetArg("mail_account", "foo bar")
	job.Args.SetArg("path", "INBOX/foo/bar")

	var wg int32
	newSeenJob(job, &wg)
	require.Equal(t,Released, wg)
	require.Equal(t, "10110-00101-001", mock.Assert())
}

func TestSeenJobFailedLogin(t *testing.T){
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
		mock.StoreCallback = func(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
			return nil
		}

		return mock, nil
	})

	job:= Job{}
	job.accounts = new(account.Accounts)
	job.accounts.Account = make(map[string]account.Account)
	job.accounts.Account["foo bar"] = account.NewAccountForTest(t,"foo bar", "foo","bar","bar.foo",  true)
	job.Args = arguments.NewEmptyArgs()
	job.Args.SetArg("mail_account", "foo bar")
	job.Args.SetArg("path", "INBOX/foo/bar")

	var wg int32
	newSeenJob(job, &wg)
	require.Equal(t,Released, wg)
	require.Equal(t, "10000-00000-001", mock.Assert())
}

func TestSeenJobFailedSelect(t *testing.T){
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
		mock.StoreCallback = func(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
			return nil
		}

		return mock, nil
	})

	job:= Job{}
	job.accounts = new(account.Accounts)
	job.accounts.Account = make(map[string]account.Account)
	job.accounts.Account["foo bar"] = account.NewAccountForTest(t,"foo bar", "foo","bar","bar.foo",  true)
	job.Args = arguments.NewEmptyArgs()
	job.Args.SetArg("mail_account", "foo bar")
	job.Args.SetArg("path", "INBOX/foo/bar")

	var wg int32
	newSeenJob(job, &wg)
	require.Equal(t,Released, wg)
	require.Equal(t, "10100-00000-001", mock.Assert())
}

func TestSeenJobFailedStoreEmpty(t *testing.T){
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
		mock.StoreCallback = func(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
			return nil
		}

		return mock, nil
	})

	job:= Job{}
	job.accounts = new(account.Accounts)
	job.accounts.Account = make(map[string]account.Account)
	job.accounts.Account["foo bar"] = account.NewAccountForTest(t,"foo bar", "foo","bar","bar.foo",  true)
	job.Args = arguments.NewEmptyArgs()
	job.Args.SetArg("mail_account", "foo bar")
	job.Args.SetArg("path", "INBOX/foo/bar")

	var wg int32
	newSeenJob(job, &wg)
	require.Equal(t,Released, wg)
	require.Equal(t, "10110-00000-001", mock.Assert())
}

func TestSeenJobFailedStoredEmpty(t *testing.T){
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
		mock.StoreCallback = func(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
			require.NotNil(t, seqSet)
			require.NotNil(t, item)
			require.NotNil(t, value)
			_v := value.([]interface{})
			require.Len(t, _v, 1)
			require.Equal(t, imap.SeenFlag, _v[0].(string))
			require.Nil(t, ch)
			if seqSet.Set[0].Start < 10 || seqSet.Set[0].Stop > 12 {
				require.Fail(t, "invalid seqSet")
			}
			return e.NewEmpty()
		}

		return mock, nil
	})

	job:= Job{}
	job.accounts = new(account.Accounts)
	job.accounts.Account = make(map[string]account.Account)
	job.accounts.Account["foo bar"] = account.NewAccountForTest(t,"foo bar", "foo","bar","bar.foo",  true)
	job.Args = arguments.NewEmptyArgs()
	job.Args.SetArg("mail_account", "foo bar")
	job.Args.SetArg("path", "INBOX/foo/bar")

	var wg int32
	newSeenJob(job, &wg)
	require.Equal(t,Released, wg)
	require.Equal(t, "10110-00101-001", mock.Assert())
}

func TestSeenJobFailedPanicUnlock(t *testing.T){
	defer func() {
		account.SetClientFactory(nil)
		err := recover()
		require.Nil(t,err)
	}()

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
		mock.StoreCallback = func(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
			require.NotNil(t, seqSet)
			require.NotNil(t, item)
			require.NotNil(t, value)
			_v := value.([]interface{})
			require.Len(t, _v, 1)
			require.Equal(t, imap.SeenFlag, _v[0].(string))
			require.Nil(t, ch)
			if seqSet.Set[0].Start < 10 || seqSet.Set[0].Stop > 12 {
				require.Fail(t, "invalid seqSet")
			}
			return errors.New("let me panic")
		}

		return mock, nil
	})

	job:= Job{}
	job.accounts = new(account.Accounts)
	job.accounts.Account = make(map[string]account.Account)
	job.accounts.Account["foo bar"] = account.NewAccountForTest(t,"foo bar", "foo","bar","bar.foo",  true)
	job.Args = arguments.NewEmptyArgs()
	job.Args.SetArg("mail_account", "foo bar")
	job.Args.SetArg("path", "INBOX/foo/bar")

	var wg int32
	newSeenJob(job, &wg)
	require.Equal(t,Released, wg)
	require.Equal(t, "10110-00101-001", mock.Assert())
}
