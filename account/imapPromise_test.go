package account

import (
	"bytes"
	"errors"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/stretchr/testify/require"
	"log"
	"mailAssistant/conditions"
	"mailAssistant/logging"
	"os"
	"testing"
	"time"
)

const (
	expectLog = "SEVERE  [mailAssistant.account.ImapPromise#listmailboxes] Mailboxes:\n" +
		"SEVERE  [mailAssistant.account.ImapPromise#listmailboxes] * foo\n"
	expectMustDie = "SEVERE  [mailAssistant.account.ImapPromise#listmailboxes] must die\n"
)

func TestImapPromise(t *testing.T) {
	t.Run("ListMailboxes", func(t *testing.T) {
		t.Run("Ok", imapPromiseListMailboxesOk)
		t.Run("Failed", imapPromiseListMailboxesFailed)
	})
	t.Run("Logout", imapPromiseLogout)
	t.Run("FetchPromise", func(t *testing.T) {
		t.Run("Ok", func(t *testing.T) {
			t.Run("Search", imapPromiseSearchPromiseOkSearch)
			t.Run("Cursor",imapPromiseSearchPromiseOkCursor)
			t.Run("Cursor empty",imapPromiseSearchPromiseOkCursorEmpty)
		})
		t.Run("Nothing", imapPromiseSearchPromiseNothing)
		t.Run("Failed", func(t *testing.T) {
			t.Run("search", imapPromiseSearchPromiseFailedSearch)
			t.Run("fetch", imapPromiseSearchPromiseFailedFetch)
		})
	})
	t.Run("SelectPromise", func(t *testing.T) {
		t.Run("OK", imapPromiseSelectPromiseOkWithPath)
		t.Run("OK without Path", imapPromiseSelectPromiseOkWithoutPath)
		t.Run("Failed", imapPromiseSelectPromiseFailed)
	})
	t.Run("Append", func(t *testing.T) {
		t.Run("OK", imapPromiseAppendOkWithPath)
		t.Run("OK without Path", imapPromiseAppendOkWithoutPath)
		t.Run("Failed", imapPromiseAppendFailed)
	})
	t.Run("Store", func(t *testing.T) {
		t.Run("OK", imapPromiseStoreOk)
		t.Run("Failed", imapPromiseStoreFailed)
	})

	t.Run("UploadAndDeleteTransaction", func(t *testing.T) {
		t.Run("OK", imapPromiseUploadAndDeleteOK)
		t.Run("no literal", imapPromiseUploadAndDeleteOKNoLiteral)
		t.Run("Empty", imapPromiseUploadAndDeleteOKEmpty)
		t.Run("Failed Append", imapPromiseUploadAndDeleteFailedAppend)
		t.Run("Failed Delete", imapPromiseUploadAndDeleteFailedDelete)
	})
}

func imapPromiseListMailboxesOk(t *testing.T) {
	called := 0
	buffer := bytes.NewBufferString("")
	log.SetFlags(0)
	log.SetOutput(buffer)
	logging.SetLevel("mailAssistant.account.ImapPromise", "all")

	defer require.Nil(t, recover())
	defer log.SetFlags(log.LstdFlags)
	defer log.SetOutput(os.Stderr)
	defer logging.SetLevel("mailAssistant.account.ImapPromise", "OFF")

	mock := NewMockClient()
	mock.ListCallback = func(ref, name string, ch chan *imap.MailboxInfo) error {
		require.NotNil(t, ch)
		require.Empty(t, ref)
		require.Equal(t, "*", name)
		called++
		ch <- &imap.MailboxInfo{Attributes: []string{"foo", "bar"},Delimiter: "bar", Name: "foo"}
		close(ch)
		return nil
	}

	promise := newImapPromise(mock)
	promise.ListMailboxes()
	require.Equal(t, 1, called)
	require.Equal(t, expectLog, buffer.String())
	require.Equal(t, "00000-01000-000", mock.Assert())
}

func imapPromiseListMailboxesFailed(t *testing.T) {
	buffer := bytes.NewBufferString("")

	called := 0
	log.SetFlags(0)
	log.SetOutput(buffer)
	logging.SetLevel("mailAssistant.account.ImapPromise", "all")

	mock := NewMockClient()
	defer func() {
		require.Equal(t, 1, called)
		require.Equal(t, expectLog+expectMustDie, buffer.String())
		err := recover()
		require.NotNil(t, err)
		require.EqualError(t, err.(error), "must die")
		require.Equal(t, "00000-01000-000", mock.Assert())
	}()
	defer logging.SetLevel("mailAssistant.account.ImapPromise", "OFF")
	defer log.SetFlags(log.LstdFlags)
	defer log.SetOutput(os.Stderr)


	mock.ListCallback = func(ref, name string, ch chan *imap.MailboxInfo) error {
		called++
		require.NotNil(t, ch)
		require.Empty(t, ref)
		require.Equal(t, "*", name)

		ch <- &imap.MailboxInfo{Attributes: []string{"foo", "bar"},Delimiter: "bar", Name: "foo"}
		close(ch)
		return errors.New("must die")
	}

	promise := newImapPromise(mock)
	promise.ListMailboxes()
	require.Fail(t, "never exec")
}

func imapPromiseLogout(t *testing.T) {
	called := 0
	defer require.Nil(t, recover())

	mock := NewMockClient()
	mock.LogoutCallback = func() error {
		called++
		return client.ErrAlreadyLoggedOut
	}

	promise := newImapPromise(mock)
	promise.Logout()
	require.True(t, true)
	require.Equal(t, 1, called)
	require.Equal(t, "00000-00000-001", mock.Assert())
}

func imapPromiseSearchPromiseOkSearch(t *testing.T) {
	injectPromise := 0
	defer func() {
		recover()
		require.NotEmpty(t, injectPromise)
	}()

	mock := NewMockClient()
	mock.SearchCallback = func(criteria *imap.SearchCriteria) (seqNums []uint32, err error) {
		require.NotNil(t, criteria)
		require.Len(t, criteria.WithFlags, 1)
		require.Equal(t, imap.SeenFlag, criteria.WithFlags[0])
		seqNums = []uint32{10, 11, 12}
		err = nil
		return
	}
	mock.FetchCallback = func(seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
		require.NotNil(t, seqset)
		require.Len(t, seqset.Set, 1)
		require.NotNil(t, items)
		require.Len(t, items, 9)
		require.NotNil(t, ch)

		ch <- &imap.Message{SeqNum: 1, Items: make(map[imap.FetchItem]interface{}, 0), Envelope: &imap.Envelope{}, BodyStructure: &imap.BodyStructure{}, Flags: []string{}, InternalDate: time.Now(), Size: 0, Uid: 0, Body: make(map[*imap.BodySectionName]imap.Literal)}
		close(ch)
		return nil
	}

	searchFor := make([]interface{}, 2)
	searchFor[0] = "KEYWORD"
	searchFor[1] = imap.SeenFlag

	promise := newImapPromise(mock)
	promise.FetchPromise(searchFor, true, func(promise *MsgPromises) {
		injectPromise++
		require.NotNil(t, promise)
	})

	require.Equal(t, "00010-00100-000", mock.Assert())
}

func imapPromiseSearchPromiseOkCursor(t *testing.T) {
	injectPromise := 0
	defer func() {
		recover()
		require.NotEmpty(t, injectPromise)
	}()

	mock := NewMockClient()
	mock.SearchCallback = func(criteria *imap.SearchCriteria) (seqNums []uint32, err error) {
		require.Fail(t, "never call this")
		return
	}
	mock.FetchCallback = func(seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
		require.NotNil(t, seqset)
		require.Len(t, seqset.Set, 1)
		require.NotNil(t, items)
		require.Len(t, items, 9)
		require.NotNil(t, ch)

		ch <- &imap.Message{SeqNum: 1, Items: make(map[imap.FetchItem]interface{}, 0), Envelope: &imap.Envelope{}, BodyStructure: &imap.BodyStructure{}, Flags: []string{}, InternalDate: time.Now(), Size: 0, Uid: 0, Body: make(map[*imap.BodySectionName]imap.Literal)}
		close(ch)
		return nil
	}

	searchFor := make([]interface{}, 1)
	searchFor[0] = conditions.CURSOR

	promise := newImapPromise(mock)
	promise.messages = 1
	promise.FetchPromise(searchFor, true, func(promise *MsgPromises) {
		injectPromise++
		require.NotNil(t, promise)
	})
	require.Equal(t, "00000-00100-000", mock.Assert())
}

func imapPromiseSearchPromiseOkCursorEmpty(t *testing.T) {
	injectPromise := 0
	defer func() {
		err := recover()
		require.Nil(t, err)
		require.Empty(t, injectPromise)
	}()

	mock := NewMockClient()
	mock.SearchCallback = func(criteria *imap.SearchCriteria) (seqNums []uint32, err error) {
		require.Fail(t, "never call this")
		return
	}
	mock.FetchCallback = func(seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
		require.NotNil(t, seqset)
		require.Len(t, seqset.Set, 1)
		require.NotNil(t, items)
		require.Len(t, items, 9)
		require.NotNil(t, ch)

		ch <- &imap.Message{SeqNum: 1, Items: make(map[imap.FetchItem]interface{}, 0), Envelope: &imap.Envelope{}, BodyStructure: &imap.BodyStructure{}, Flags: []string{}, InternalDate: time.Now(), Size: 0, Uid: 0, Body: make(map[*imap.BodySectionName]imap.Literal)}
		close(ch)
		return nil
	}

	searchFor := make([]interface{}, 1)
	searchFor[0] = conditions.CURSOR

	promise := newImapPromise(mock)
	promise.messages = 0
	promise.FetchPromise(searchFor, true, func(promise *MsgPromises) {
		injectPromise++
		require.NotNil(t, promise)
	})
	require.Empty(t, injectPromise)
	require.Equal(t, "00000-00000-000", mock.Assert())
}

func imapPromiseSearchPromiseNothing(t *testing.T) {
	injectPromise := 0
	mock := NewMockClient()

	defer func() {
		err := recover()
		require.Nil(t, err)
	}()


	mock.SearchCallback = func(criteria *imap.SearchCriteria) (seqNums []uint32, err error) {
		require.NotNil(t, criteria)
		require.Len(t, criteria.WithFlags, 1)
		require.Equal(t, imap.SeenFlag, criteria.WithFlags[0])
		seqNums = make([]uint32, 0)
		err = nil
		return
	}

	searchFor := make([]interface{}, 2)
	searchFor[0] = "KEYWORD"
	searchFor[1] = imap.SeenFlag

	promise := newImapPromise(mock)
	promise.FetchPromise(searchFor, true, func(promise *MsgPromises) {
		injectPromise++
		require.NotNil(t, promise)
	})
	require.Equal(t,1, injectPromise)
	require.Equal(t, "00010-00000-000", mock.Assert())
}

func imapPromiseSearchPromiseFailedSearch(t *testing.T) {
	injectPromise := 0
	mock := NewMockClient()
	defer func() {
		err := recover()
		require.NotNil(t, err)
		require.Equal(t, "search must fail []interface {}{\"KEYWORD\", \"\\\\Seen\"}", err.(error).Error())
		require.Empty(t, injectPromise)
		require.Equal(t, "00010-00000-000",mock.Assert())
	}()

	mock.SearchCallback = func(criteria *imap.SearchCriteria) (seqNums []uint32, err error) {
		require.NotNil(t, criteria)
		require.Len(t, criteria.WithFlags, 1)
		require.Equal(t, imap.SeenFlag, criteria.WithFlags[0])
		seqNums = []uint32{10, 11, 12}
		err = errors.New("search must fail")
		return
	}

	searchFor := make([]interface{}, 2)
	searchFor[0] = "KEYWORD"
	searchFor[1] = imap.SeenFlag

	promise := newImapPromise(mock)
	promise.FetchPromise(searchFor, true, func(promise *MsgPromises) {
		injectPromise++
		require.Nil(t, promise)
	})
}

func imapPromiseSearchPromiseFailedFetch(t *testing.T) {
	injectPromise := 0

	mock := NewMockClient()
	defer func() {
		err := recover()
		require.NotNil(t, err)
		require.EqualError(t, err.(error), "fetch must fail", )
		require.Empty(t, injectPromise)
		require.Equal(t, "00010-00100-000",mock.Assert())
	}()
	mock.SearchCallback = func(criteria *imap.SearchCriteria) (seqNums []uint32, err error) {
		require.NotNil(t, criteria)
		require.Len(t, criteria.WithFlags, 1)
		require.Equal(t, imap.SeenFlag, criteria.WithFlags[0])
		seqNums = []uint32{10, 11, 12}
		err = nil
		return
	}
	mock.FetchCallback = func(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
		require.NotNil(t, seqSet)
		require.Len(t, seqSet.Set, 1)
		require.NotNil(t, items)
		require.Len(t, items, 9)
		require.NotNil(t, ch)

		ch <- &imap.Message{SeqNum: 1, Items: make(map[imap.FetchItem]interface{}, 0), Envelope: &imap.Envelope{}, BodyStructure: &imap.BodyStructure{}, Flags: []string{}, InternalDate: time.Now(), Size: 0, Uid: 0, Body: make(map[*imap.BodySectionName]imap.Literal)}
		close(ch)
		return errors.New("fetch must fail")
	}

	searchFor := make([]interface{}, 2)
	searchFor[0] = "KEYWORD"
	searchFor[1] = imap.SeenFlag

	promise := newImapPromise(mock)
	promise.FetchPromise(searchFor, true, func(promise *MsgPromises) {
		injectPromise++
		require.NotNil(t, promise)
	})
}

func imapPromiseSelectPromiseOkWithoutPath(t *testing.T) {
	called := 0
	mock := NewMockClient()
	mock.SelectCallback = func(name string, readOnly bool) (status *imap.MailboxStatus, err error) {
		require.NotEmpty(t, name)
		require.True(t, readOnly)
		require.Equal(t, "INBOX", name)
		status = new(imap.MailboxStatus)
		return
	}
	promise := newImapPromise(mock)
	promise.SelectPromise("", true, func(promise *ImapPromise) {
		called++
		require.NotNil(t, promise)
	})
	require.NotEmpty(t, called)
	require.Equal(t, "00100-00000-000",mock.Assert())
}
func imapPromiseSelectPromiseOkWithPath(t *testing.T) {
	called := 0
	mock := NewMockClient()
	defer func() {
		require.NotEmpty(t, called)
		require.Equal(t, "00100-00000-000",mock.Assert())
	}()

	mock.SelectCallback = func(name string, readOnly bool) (status *imap.MailboxStatus, err error) {
		require.True(t, readOnly)
		require.NotEmpty(t, name)
		require.Equal(t, "INBOX.test.hugo", name)
		status = new(imap.MailboxStatus)
		return
	}
	promise := newImapPromise(mock)
	promise.SelectPromise("INBOX/test/hugo", true, func(promise *ImapPromise) {
		called++
		require.NotNil(t, promise)
	})
	require.Equal(t, "00100-00000-000", mock.Assert())
}

func imapPromiseSelectPromiseFailed(t *testing.T) {
	called := 0
	mock := NewMockClient()
	defer func() {
		require.Empty(t, called)
		err := recover()
		require.NotNil(t, err)
		require.Error(t, err.(error), "must fail")
		require.Equal(t, "00100-00000-000",mock.Assert())
	}()

	mock.SelectCallback = func(name string, readOnly bool) (status *imap.MailboxStatus, err error) {
		require.NotEmpty(t, name)
		require.True(t, readOnly)

		err = errors.New("must fail")
		return
	}
	promise := newImapPromise(mock)
	promise.SelectPromise("INBOX/test/hugo", true, func(promise *ImapPromise) {
		called++
		require.Fail(t, "never call this")
	})
}

func imapPromiseAppendOkWithPath(t *testing.T) {
	now := time.Now()
	mock := NewMockClient()
	mock.AppendCallback = func(mbox string, flags []string, date time.Time, msg imap.Literal) error {
		require.Equal(t, "INBOX.test.hugo", mbox)
		require.Nil(t, flags)
		require.Exactly(t, now, date)
		require.Nil(t, msg)
		return nil
	}
	exec := 0
	promise := newImapPromise(mock)
	promise.AppendPromise("INBOX/test/hugo", nil, now, nil, func() {
		exec++
	})
	require.NotEmpty(t, exec)
	require.Equal(t, "00000-00010-000", mock.Assert())
}

func imapPromiseAppendOkWithoutPath(t *testing.T) {
	now := time.Now()
	mock := NewMockClient()
	mock.AppendCallback = func(mbox string, flags []string, date time.Time, msg imap.Literal) error {
		require.Equal(t, "INBOX", mbox)
		require.Nil(t, flags)
		require.Exactly(t, now, date)
		require.Nil(t, msg)
		return nil
	}

	exec := 0
	promise := newImapPromise(mock)
	promise.AppendPromise("", nil, now, nil, func() {
		exec++
	})
	require.NotEmpty(t, exec)
	require.Equal(t, "00000-00010-000", mock.Assert())
}

func imapPromiseAppendFailed(t *testing.T) {
	now := time.Now()
	mock := NewMockClient()
	mock.AppendCallback = func(mbox string, flags []string, date time.Time, msg imap.Literal) error {
		require.Equal(t, "INBOX", mbox)
		require.Nil(t, flags)
		require.Exactly(t, now, date)
		require.Nil(t, msg)
		return errors.New("append must fail")
	}
	promise := newImapPromise(mock)
	promise.AppendPromise("", nil, now, nil, func() {
		require.Fail(t, "never call this")
	})
	require.Equal(t, "00000-00010-000", mock.Assert())
}

func imapPromiseStoreOk(t *testing.T) {
	mock := NewMockClient()
	mock.StoreCallback = func(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
		require.NotNil(t, seqSet)
		require.NotNil(t, item)
		require.Equal(t, "+FLAGS", string(item))
		require.NotNil(t, value)
		require.Equal(t, interface{}(1), value)
		require.Nil(t, ch)
		return nil
	}
	promise := newImapPromise(mock)
	require.Nil(t, promise.Store(&imap.SeqSet{}, imap.FormatFlagsOp(imap.AddFlags, false), 1, nil))
	require.Equal(t, "00000-00001-000", mock.Assert())
}

func imapPromiseStoreFailed(t *testing.T) {
	mock := NewMockClient()
	mock.StoreCallback = func(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
		require.NotNil(t, seqSet)
		require.NotNil(t, item)
		require.Equal(t, "+FLAGS", string(item))
		require.NotNil(t, value)
		require.Equal(t, interface{}(1), value)
		require.Nil(t, ch)
		return errors.New("store must fail")
	}
	promise := newImapPromise(mock)
	require.EqualError(t, promise.Store(&imap.SeqSet{}, imap.FormatFlagsOp(imap.AddFlags, false), 1, nil), "store must fail")
	require.Equal(t, "00000-00001-000", mock.Assert())
}

func imapPromiseUploadAndDeleteOK(t *testing.T) {
	mock := NewMockClient()
	mock.AppendCallback = func(mBox string, flags []string, date time.Time, msg imap.Literal) error {
		require.Equal(t, "INBOX", mBox)
		require.NotNil(t, flags)
		require.Len(t, flags,1)
		require.NotNil(t, date)
		require.NotNil(t, msg)
		return nil
	}
	mock.DeleteCallback = func(num uint32) error {
		require.Equal(t, uint32(10), num)
		return nil
	}
	mockMsg := MockMessage{}
	mockMsg.callback = func(section *imap.BodySectionName) imap.Literal {
		require.NotNil(t, section)
		buf := bytes.NewBufferString("")
		return buf
	}

	messages := make([]*MsgPromise, 1)
	messages[0] = newMsgPromise(mockMsg, 10, mock)

	msgPromises := MsgPromises{newImapPromise(mock), messages, &imap.SeqSet{}}

	called := 0
	promise := newImapPromise(mock)
	promise.UploadAndDelete("", &msgPromises, func(num int) {
		require.NotEmpty(t, num)
		called++
	})
	require.NotEmpty(t, called)
	require.Equal(t, "00000-00010-100", mock.Assert())
}

func imapPromiseUploadAndDeleteOKNoLiteral(t *testing.T) {
	mock := NewMockClient()
	mock.AppendCallback = func(mBox string, flags []string, date time.Time, msg imap.Literal) error {
		require.Fail(t, "never call append")
		return nil
	}
	mock.DeleteCallback = func(num uint32) error {
		require.Fail(t, "never call delete")
		return nil
	}
	mockMsg := MockMessage{}
	mockMsg.callback = func(section *imap.BodySectionName) imap.Literal {
		require.NotNil(t, section)
		return nil
	}

	messages := make([]*MsgPromise, 1)
	messages[0] = newMsgPromise(mockMsg, 10, mock)

	msgPromises := MsgPromises{newImapPromise(mock), messages, &imap.SeqSet{}}

	called := 0
	promise := newImapPromise(mock)
	promise.UploadAndDelete("", &msgPromises, func(num int) {
		require.Empty(t, num)
		called++
	})
	require.NotEmpty(t, called)
	require.Equal(t, "00000-00000-000", mock.Assert())
}

func imapPromiseUploadAndDeleteOKEmpty(t *testing.T) {
	mock := NewMockClient()
	mock.AppendCallback = func(mBox string, flags []string, date time.Time, msg imap.Literal) error {
		return nil
	}
	mock.DeleteCallback = func(num uint32) error {
		return nil
	}
	messages := make([]*MsgPromise, 0)
	msgPromises := MsgPromises{newImapPromise(mock), messages, &imap.SeqSet{}}

	called := 0
	promise := newImapPromise(mock)
	promise.UploadAndDelete("", &msgPromises, func(num int) {
		require.Empty(t, num)
		called++
	})
	require.NotEmpty(t, called)
	require.Equal(t, "00000-00000-000", mock.Assert())
}

func imapPromiseUploadAndDeleteFailedAppend(t *testing.T) {
	mock := NewMockClient()
	mock.AppendCallback = func(mBox string, flags []string, date time.Time, msg imap.Literal) error {
		return errors.New("append must fail")
	}
	mock.DeleteCallback = func(num uint32) error {
		return nil
	}

	mockMsg := MockMessage{}
	mockMsg.callback = func(section *imap.BodySectionName) imap.Literal {
		require.NotNil(t, section)
		return bytes.NewBufferString("")
	}

	messages := make([]*MsgPromise, 1)
	messages[0] = newMsgPromise(mockMsg, 10, mock)

	msgPromises := MsgPromises{newImapPromise(mock), messages, &imap.SeqSet{}}

	called := 0
	promise := newImapPromise(mock)
	promise.UploadAndDelete("", &msgPromises, func(num int) {
		require.Empty(t, num)
		called++
	})
	require.NotEmpty(t, called)
	require.Equal(t, "00000-00010-000", mock.Assert())
}

func imapPromiseUploadAndDeleteFailedDelete(t *testing.T) {
	mock := NewMockClient()
	mock.AppendCallback = func(mBox string, flags []string, date time.Time, msg imap.Literal) error {
		return nil
	}
	mock.DeleteCallback = func(num uint32) error {
		return errors.New("delete must fail")
	}

	mockMsg := MockMessage{}
	mockMsg.callback = func(section *imap.BodySectionName) imap.Literal {
		require.NotNil(t, section)
		return bytes.NewBufferString("")
	}

	messages := make([]*MsgPromise, 1)
	messages[0] = newMsgPromise(mockMsg, 10, mock)

	msgPromises := MsgPromises{newImapPromise(mock), messages, &imap.SeqSet{}}

	called := 0
	promise := newImapPromise(mock)
	promise.UploadAndDelete("", &msgPromises, func(num int) {
		require.Empty(t, num)
		called++
	})
	require.NotEmpty(t, called)
	require.Equal(t, "00000-00010-100", mock.Assert())
}
