package account

import (
	"github.com/emersion/go-imap"
	"mailAssistant/logging"
	"strings"
	"time"
)

var fetchFullMessage = []imap.FetchItem{
	imap.FetchBody,
	imap.FetchBodyStructure,
	imap.FetchEnvelope,
	imap.FetchFlags,
	imap.FetchInternalDate,
	imap.FetchRFC822,
	imap.FetchRFC822Header,
	imap.FetchRFC822Size,
	imap.FetchRFC822Text,
}

// ImapPromise is a promise obj to cover all client lib activities
type ImapPromise struct {
	client IClient
	logger *logging.Logger
}

func newImapPromise(connection IClient) *ImapPromise {
	prom := new(ImapPromise)
	prom.client = connection
	return prom
}

func (promise *ImapPromise) getLogger() *logging.Logger {
	if promise.logger == nil {
		promise.logger = logging.NewLogger()
	}
	return promise.logger
}

// ListMailboxes is listing all mailboxes on remote server
func (promise *ImapPromise) ListMailboxes() {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- promise.client.List("", "*", mailboxes)
	}()

	promise.getLogger().Severe("Mailboxes:")
	for m := range mailboxes {
		promise.getLogger().Severe("*", m.Name)
	}
	if err := <-done; err != nil {
		promise.getLogger().Panic(err)
	}
}

// AppendPromise adds a mail on the IMAP server
func (promise ImapPromise) AppendPromise(saveTo string, flags []string, date time.Time, msg imap.Literal, successfully func()) (err error) {
	saveTo = strings.ReplaceAll(saveTo, "/", ".")
	if saveTo == "" {
		saveTo = "INBOX"
	}
	if err = promise.client.Append(saveTo, flags, date, msg); err == nil {
		successfully()
	}
	return err
}

// Store sets flags on a mail
func (promise ImapPromise) Store(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
	return promise.client.Store(seqSet, item, value, ch)
}

// Logout disconnect from server
func (promise ImapPromise) Logout() {
	_ = promise.client.Logout()
}

// SelectPromise is selecting a mailbox, if successfully it calls a callback
func (promise ImapPromise) SelectPromise(path string, readOnly bool, callback func(promise *ImapPromise)) {
	if path == "" {
		path = "INBOX"
	}
	path = strings.Replace(path, "/", ".", -1)
	if _, err := promise.client.Select(path, readOnly); err != nil {
		panic(err)
	}
	callback(&promise)
}

// SearchPromise is searching on the IMAP server, if successfully it calls a callback
func (promise ImapPromise) SearchPromise(args [][]interface{}, fetchContent bool, callback func(promise *MsgPromises)) {
	searchCfg := imap.NewSearchCriteria()
	for _, arg := range args {
		_ = searchCfg.ParseWithCharset(arg, nil)
	}

	if seqNums, err := promise.client.Search(searchCfg); err != nil {
		panic(err)
	} else if len(seqNums) <= 0 {
		callback(&MsgPromises{ImapPromise: &promise, messages: make([]*MsgPromise, 0), seqSet: new(imap.SeqSet)})
	} else {
		seqSet := new(imap.SeqSet)
		seqSet.AddNum(seqNums...)

		fetchItems := imap.FetchFast.Expand()
		if fetchContent {
			fetchItems = fetchFullMessage
		}

		done := make(chan error, 1)
		messages := make(chan *imap.Message, 10)
		go func() {
			done <- promise.client.Fetch(seqSet, fetchItems, messages)
		}()

		msgPromise := MsgPromises{&promise, make([]*MsgPromise, 0), seqSet}
		msgPromise.addAll(messages)

		if err := <-done; err != nil {
			promise.getLogger().Panic(err)
		}
		if len(seqNums) > 0 {
			callback(&msgPromise)
		}
	}
}

// UploadAndDelete is uploading a message literal
// if successfully it is deleting the corresponding message
func (promise *ImapPromise) UploadAndDelete(saveTo string, messages *MsgPromises, callback func(num int)) {
	flags := []string{imap.SeenFlag}
	moved := 0
	var date time.Time

	messages.Each(func(message *MsgPromise) {
		if literal := message.GetLiteral(); literal != nil {
			promise.AppendPromise(saveTo, flags, date, literal, func() {
				message.DeletePromise(func(err error) {
					if err == nil {
						moved++
					} else {
						promise.getLogger().Severe("Delete", err)
					}
				})
			})
		}
	})
	callback(moved)
}
