package actions

import (
	"bytes"
	"crypto/tls"
	"github.com/emersion/go-imap"
	"github.com/stretchr/testify/require"
	"log"
	"mailAssistant/account"
	"mailAssistant/arguments"
	"mailAssistant/logging"
	"os"
	"testing"
)

const listMailboxExpected = `INFO    [mailAssistant.actions.TestListMailbox#islockedelselock] lock
SEVERE  [mailAssistant.account.ImapPromise#listmailboxes] Mailboxes:
SEVERE  [mailAssistant.account.ImapPromise#listmailboxes] * TEST.A
INFO    [mailAssistant.actions.TestListMailbox#unlockalways] unlocked
`

func TestListMailbox(t *testing.T){
	mock := account.NewMockClient()
	mock.ListCallback = func(ref, name string, ch chan *imap.MailboxInfo) error {
		require.NotNil(t, ref)
		require.Equal(t,"", ref)
		require.NotNil(t, name)
		require.Equal(t,"*", name)
		require.NotNil(t, ch)

		ch <- &imap.MailboxInfo{make([]string,0),".","TEST.A"}
		close(ch)
		return nil
	}
	buffer := bytes.NewBufferString("")
	log.SetOutput(buffer)
	log.SetFlags(0)
	logging.SetLevel("*","")
	logging.SetLevel("global","ALL")
	defer logging.SetLevel("*","")
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.LstdFlags)

	acc := account.NewAccountForTest(t,"foo","foo","bar","a",false)
	accs := new(account.Accounts)
	accs.Account = map[string]account.Account{"foo":acc}
	account.SetClientFactory(func(addr string, tlsConfig *tls.Config) (client account.IClient, err error) {
		return mock, nil
	})
	job := Job{Args: arguments.NewEmptyArgs(),log:logging.NewLogger(),accounts:accs}
	job.SetArg("mail_account","foo")
	wg :=new(int32)

	newListMailbox(job, wg, nil)
	require.Empty(t, *wg)
	require.Equal(t, listMailboxExpected, buffer.String())
}

func TestListMailboxLocked(t *testing.T){
	job := Job{log: logging.NewLogger()}
	wg :=int32(1)
	newListMailbox(job, &wg, nil)
	require.NotEmpty(t,wg)
}