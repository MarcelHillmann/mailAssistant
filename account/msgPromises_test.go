package account

import (
	"errors"
	"github.com/emersion/go-imap"
	"github.com/stretchr/testify/require"
	e "mailAssistant/errors"
	"testing"
)

func TestMsgPromises(t *testing.T) {
	t.Run("delete", func(t *testing.T) {
		t.Run("ok", msgPromisesDelete)
		t.Run("fail nothing marked", msgPromisesDeleteFailedNotMarked)
		t.Run("fail store", msgPromisesDeleteFailedStore)
		t.Run("fail expunge", msgPromisesDeleteFailedExpunge)
	})
	t.Run("get attachments", func(t *testing.T) {
		t.Run("ok", msgPromisesGetAttachmentsOk)
		t.Run("empty", msgPromisesGetAttachmentsEmpty)
		t.Run("not found", msgPromisesGetAttachmentsNotFound)
	})
	t.Run("move", func(t *testing.T) {
		t.Run("OK", msgPromisesMove)
		t.Run("Failed", msgPromisesMoveFailed)
		t.Run("SeqSet Empty", msgPromisesMoveSeqSetEmpty)
	})
	t.Run("messages", msgPromisesMessages)
	t.Run("setSeen", func(t *testing.T) {
		t.Run("OK", msgPromisesSetSeen)
		t.Run("Failed", msgPromisesSetSeenFailed)
		t.Run("SeSet Empty", msgPromisesSetSeenSeqSetEmpty)
	})
	t.Run("Expunge", func(t *testing.T) {
		t.Run("OK", msgPromisesExpunge)
		t.Run("Failed", msgPromisesExpungeFailed)
	})
}

func msgPromisesDelete(t *testing.T) {
	client := NewMockClientMinimal()
	client.StoreCallback = func(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
		return nil
	}
	client.ExpungeCallback = func(ch chan uint32) error {
		return nil
	}
	seqSet, _ := imap.ParseSeqSet("10,11,12")
	msg := MsgPromises{newImapPromise(client), make([]*MsgPromise, 0), seqSet}
	c, err := msg.Delete()
	require.NotNil(t, c)
	require.Nil(t, err)
	require.Equal(t, 3, c)
	require.Equal(t, "00000-00001-010", client.Assert())
}

func msgPromisesDeleteFailedNotMarked(t *testing.T) {
	client := NewMockClientMinimal()
	msg := MsgPromises{newImapPromise(client), make([]*MsgPromise, 0), &imap.SeqSet{Set: make([]imap.Seq, 0)}}
	c, err := msg.Delete()
	require.NotNil(t, c)
	require.NotNil(t, err)
	require.Empty(t, c)
	require.Equal(t, e.NewEmpty(), err)
	require.Equal(t, "00000-00000-000", client.Assert())
}

func msgPromisesDeleteFailedStore(t *testing.T) {
	client := NewMockClientMinimal()
	client.StoreCallback = func(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
		return errors.New("store must fail")
	}
	client.ExpungeCallback = func(ch chan uint32) error {
		return errors.New("expunge must fail")
	}
	seqSet, _ := imap.ParseSeqSet("10,11,12")
	msg := MsgPromises{newImapPromise(client), make([]*MsgPromise, 0), seqSet}
	c, err := msg.Delete()
	require.NotNil(t, c)
	require.NotNil(t, err)
	require.Empty(t, c)
	require.EqualError(t, err, "store must fail")
	require.Equal(t, "00000-00001-000", client.Assert())
}

func msgPromisesDeleteFailedExpunge(t *testing.T) {
	client := NewMockClientMinimal()
	client.StoreCallback = func(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
		return nil
	}
	client.ExpungeCallback = func(ch chan uint32) error {
		return errors.New("expunge must fail")
	}
	seqSet, _ := imap.ParseSeqSet("10,11,12")
	msg := MsgPromises{newImapPromise(client), []*MsgPromise{}, seqSet}
	c, err := msg.Delete()
	require.NotNil(t, c)
	require.NotNil(t, err)
	require.Empty(t, c)
	require.EqualError(t, err, "expunge must fail")
	require.Equal(t, "00000-00001-010", client.Assert())
}

func msgPromisesGetAttachmentsEmpty(t *testing.T) {
	client := NewMockClientMinimal()
	m := new(MockMessage)
	m.callback = func(section *imap.BodySectionName) imap.Literal {
		return nil
	}
	messages := make([]*MsgPromise, 1)
	messages[0] = newMsgPromise(m, 0, nil)
	seqSet, _ := imap.ParseSeqSet("10,11")
	msg := MsgPromises{newImapPromise(client), messages, seqSet}
	promise := msg.GetAttachments("application/pdf")
	require.NotNil(t, promise)
	require.Len(t, promise, 0)
	require.Equal(t, "00000-00000-000", client.Assert())
}

func msgPromisesGetAttachmentsOk(t *testing.T) {
	client := NewMockClientMinimal()
	messages := make([]*MsgPromise, 2)
	m := new(MockMessage)
	m.callback = func(section *imap.BodySectionName) imap.Literal {
		return CreateMail([]string{"application/pdf", "pdf"})
	}
	messages[0] = newMsgPromise(m, 0, nil)
	m.callback = func(section *imap.BodySectionName) imap.Literal {
		return CreateMail([]string{"application/pdf", "pdf"})
	}
	messages[1] = newMsgPromise(m, 0, nil)
	seqSet, _ := imap.ParseSeqSet("10,11")
	msg := MsgPromises{newImapPromise(client), messages, seqSet}
	promise := msg.GetAttachments("application/pdf")
	require.NotNil(t, promise)
	require.Len(t, promise, 2)
	require.Equal(t, "00000-00000-000", client.Assert())
}

func msgPromisesGetAttachmentsNotFound(t *testing.T) {
	client := NewMockClientMinimal()
	messages := make([]*MsgPromise, 2)
	m := new(MockMessage)
	m.callback = func(section *imap.BodySectionName) imap.Literal {
		return CreateMail([]string{"text/html", "html"})
	}
	messages[0] = newMsgPromise(m, 0, nil)
	m.callback = func(section *imap.BodySectionName) imap.Literal {
		return CreateMail([]string{"text/plain", "txt"})
	}
	messages[1] = newMsgPromise(m, 0, nil)
	seqSet, _ := imap.ParseSeqSet("10,11")
	msg := MsgPromises{newImapPromise(client), messages, seqSet}
	promise := msg.GetAttachments("application/pdf")
	require.NotNil(t, promise)
	require.Len(t, promise, 0)
	require.Equal(t, "00000-00000-000", client.Assert())
}

func msgPromisesMove(t *testing.T) {
	client := NewMockClient()
	client.MoveCallback = func(seqSet *imap.SeqSet, dest string) error {
		require.NotNil(t, seqSet)
		require.Len(t, seqSet.Set, 1)
		require.Equal(t, uint32(10), seqSet.Set[0].Start)
		require.Equal(t, uint32(12), seqSet.Set[0].Stop)
		require.Equal(t, "INBOX", dest)
		return nil
	}
	seqSet, _ := imap.ParseSeqSet("10,11,12")
	msg := MsgPromises{newImapPromise(client), nil, seqSet}
	c, err := msg.Move("")
	require.NotNil(t, c)
	require.Nil(t, err)
	require.NotEmpty(t, c)
	require.Equal(t, 3, c)
	require.Equal(t, "00001-00000-000", client.Assert())
}

func msgPromisesMoveSeqSetEmpty(t *testing.T) {
	client := NewMockClient()
	seqSet, _ := imap.ParseSeqSet("")
	msg := MsgPromises{newImapPromise(client), nil, seqSet}
	c, err := msg.Move("")
	require.NotNil(t, c)
	require.NotNil(t, err)
	require.EqualError(t, err, "is empty")
	require.Equal(t, "00000-00000-000", client.Assert())
}

func msgPromisesMoveFailed(t *testing.T) {
	client := NewMockClient()
	client.MoveCallback = func(seqSet *imap.SeqSet, dest string) error {
		require.Equal(t, "hugo.boss", dest)
		return errors.New("move must fail")
	}
	seqSet, _ := imap.ParseSeqSet("10,11,12")
	msg := MsgPromises{newImapPromise(client), nil, seqSet}
	c, err := msg.Move("hugo/boss")
	require.NotNil(t, c)
	require.NotNil(t, err)
	require.EqualError(t, err, "move must fail")
	require.Equal(t, "00001-00000-000", client.Assert())
}

func msgPromisesMessages(t *testing.T) {
	msg := MsgPromises{nil, make([]*MsgPromise, 2), nil}
	require.NotNil(t, msg)
	require.NotNil(t, msg.messages)
	require.NotEmpty(t, msg.messages)
	require.Len(t, msg.messages, 2)
	require.Nil(t, msg.Message(0))
	require.Nil(t, msg.Message(1))
}

func msgPromisesSetSeen(t *testing.T) {
	client := NewMockClient()
	client.StoreCallback = func(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
		require.NotNil(t, seqSet)
		require.NotNil(t, item)
		require.NotNil(t, value)
		require.Nil(t, ch)
		require.Len(t, seqSet.Set, 1)
		require.Equal(t, "+FLAGS.SILENT", string(item))
		require.Len(t, value, 1)
		return nil
	}
	seqSet, _ := imap.ParseSeqSet("10,11,12")
	msg := MsgPromises{newImapPromise(client), nil, seqSet}
	count, err := msg.SetSeen()
	require.Nil(t, err)
	require.NotNil(t, count)
	require.NotEmpty(t, count)
	require.Equal(t, count, 3)
	require.Equal(t, "00000-00001-000", client.Assert())
}

func msgPromisesSetSeenFailed(t *testing.T) {
	client := NewMockClient()
	client.StoreCallback = func(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
		require.NotNil(t, seqSet)
		require.NotNil(t, item)
		require.NotNil(t, value)
		require.Nil(t, ch)
		require.Len(t, seqSet.Set, 1)
		require.Equal(t, "+FLAGS.SILENT", string(item))
		require.Len(t, value, 1)
		return errors.New("store must fail")
	}
	seqSet, _ := imap.ParseSeqSet("10,11,12")
	msg := MsgPromises{newImapPromise(client), nil, seqSet}
	count, err := msg.SetSeen()
	require.NotNil(t, err)
	require.EqualError(t, err, "store must fail")
	require.NotNil(t, count)
	require.Empty(t, count)
	require.Equal(t, "00000-00001-000", client.Assert())
}

func msgPromisesSetSeenSeqSetEmpty(t *testing.T) {
	seqSet, _ := imap.ParseSeqSet("")
	msg := MsgPromises{nil, nil, seqSet}
	count, err := msg.SetSeen()
	require.NotNil(t, err)
	require.EqualError(t, err, "is empty")
	require.NotNil(t, count)
	require.Empty(t, count)
}

func msgPromisesExpunge(t *testing.T) {
	defer func() {
		err := recover()
		require.Nil(t, err)
	}()
	client := NewMockClient()
	client.ExpungeCallback = func(ch chan uint32) error {
		require.Nil(t, ch)
		return nil
	}
	msg := MsgPromises{newImapPromise(client), nil, nil}
	msg.Expunge()
	require.Equal(t, "00000-00000-010", client.Assert())
}

func msgPromisesExpungeFailed(t *testing.T) {
	defer func() {
		err := recover()
		require.Nil(t, err)
	}()
	client := NewMockClient()
	client.ExpungeCallback = func(ch chan uint32) error {
		require.Nil(t, ch)
		return nil
	}
	msg := MsgPromises{newImapPromise(client), nil, nil}
	msg.Expunge()
	require.Equal(t, "00000-00000-010", client.Assert())
}
