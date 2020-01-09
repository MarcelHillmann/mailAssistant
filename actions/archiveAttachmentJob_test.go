package actions

import (
	"crypto/tls"
	"github.com/emersion/go-imap"
	"github.com/stretchr/testify/require"
	"mailAssistant/account"
	"mailAssistant/arguments"
	"mailAssistant/logging"
	"testing"
)

func TestArchiveAttachmentLocked(t *testing.T) {
	var wg int32 = 1
	newArchiveAttachment(Job{}, &wg, metricsDummy)
	require.Equal(t, Locked, wg)
}

func TestArchiveAttachmentSuccess(t *testing.T) {
	defer account.SetClientFactory(nil)

	mock := account.NewMockClientMinimal()
	account.SetClientFactory(func(addr string, tlsConfig *tls.Config) (account.IClient, error) {
		require.Equal(t, "bar.foo:20000", addr)
		require.NotNil(t, tlsConfig)
		require.False(t, tlsConfig.InsecureSkipVerify)

		mock.LoginCallback = func(username, password string) error {
			require.Equal(t, "foo", username)
			require.Equal(t, "bar", password)
			return nil
		}
		mock.SelectCallback = func(name string, readOnly bool) (status *imap.MailboxStatus, err error) {
			require.Equal(t, "INBOX.foo.bar",name)
			require.True(t, readOnly)
			return new(imap.MailboxStatus),nil
		}
		mock.SearchCallback = func(criteria *imap.SearchCriteria) (uint32s []uint32, err error) {
			require.NotNil(t, criteria)
			return []uint32{10,11,12}, nil
		}
		mock.FetchCallback = func(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
			ch <- createMessage(10,false)
			ch <- createMessage(11,false)
			ch <- createMessage(12,false)
			close(ch)
			return nil
		}
		mock.LogoutCallback = func() error {
			return nil
		}
		return mock, nil
	})

	logging.SetLevel("unit", "all")
	job := Job{Args: arguments.NewEmptyArgs(),log: logging.NewNamedLogger("unit.tests"), accounts: new(account.Accounts)}
	job.accounts.Account = make(map[string]account.Account)
	job.accounts.Account["foo bar"] = account.NewAccountForTest(t, "foo bar", "foo", "bar", "bar.foo", false)
	job.accounts.Account["foo bar target"] = account.NewAccountForTest(t, "foo bar target", "foo", "bar", "target.local", true)
	job.SetArg("mail_account", "foo bar")
	job.SetArg("path", "INBOX/foo/bar")
	job.SetArg("readonly", true)
	job.SetArg("saveTo", "../../foo/bar")
	job.SetArg("attachment_type", "foo/bar")
	job.SetArg("search", []interface{}{})

	var wg int32
	newArchiveAttachment(job, &wg, metricsDummy)
	require.Equal(t, Released, wg)
	require.Equal(t, "10110-00100-001", mock.Assert())
}
