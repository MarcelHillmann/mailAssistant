package account

import (
	"bytes"
	"errors"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
	"github.com/stretchr/testify/require"
	"testing"
)

type MockMessage struct {
	callback func(section *imap.BodySectionName) imap.Literal
}

func (m MockMessage) GetBody(section *imap.BodySectionName) imap.Literal {
	return m.callback(section)
}

const mimeType = "application/pdf"

func TestMsgPromise(t *testing.T){
	t.Run("GetAttachment", func(t *testing.T) {
		t.Run("can't read", msgPromiseGetAttachmentCantRead)
		t.Run("found", msgPromiseGetAttachmentFound)
		t.Run("not found", msgPromiseGetAttachmentNotFound)
		t.Run("empty", msgPromiseGetAttachmentEmpty)
		t.Run("no GetLiteral", msgPromiseGetAttachmentNoLiteral)
	})
}

func TestMsgPromise_GetAttachmentFailedNextPart(t *testing.T) {
	nextPart = func(reader *mail.Reader) (part *mail.Part, err error) {
		return nil, errors.New("must fail")
	}

	defer func(){
		err := recover()
		require.NotNil(t,err)
		require.EqualError(t, err.(error),"must fail")
		nextPart = internalNextPart
	}()

	mock := new(MockMessage)
	mock.callback = func(section *imap.BodySectionName) imap.Literal {
		require.NotNil(t, section)
		return CreateMail([]string{mimeType,"pdf"})
	}

	msg := newMsgPromise(mock,0, nil)
	msg.GetAttachment(mimeType)
	require.Fail(t, "never call this")
}

func msgPromiseGetAttachmentCantRead(t *testing.T) {
	mock := new(MockMessage)
	mock.callback = func(section *imap.BodySectionName) imap.Literal {
		require.NotNil(t, section)
		buf := bytes.NewBufferString("")
		return buf
	}

	msg := newMsgPromise(mock,0, nil)
	attachments := msg.GetAttachment(mimeType)
	require.NotNil(t, attachments)
	require.Len(t, attachments, 0)
}

func msgPromiseGetAttachmentFound(t *testing.T) {
	mock := new(MockMessage)
	mock.callback = func(section *imap.BodySectionName) imap.Literal {
		require.NotNil(t, section)
		return CreateMail([]string{mimeType,"pdf"})
	}

	msg := newMsgPromise(mock,0, nil)
	attachments := msg.GetAttachment(mimeType)
	require.NotNil(t, attachments)
	require.Len(t, attachments, 1)
	require.NotNil(t, attachments[0])
	require.Equal(t, "test.pdf", attachments[0].GetFilename())
	require.Equal(t, mimeType, attachments[0].GetContentType())
	require.Equal(t, "attachment", attachments[0].GetContentDisposition())
}

func msgPromiseGetAttachmentNotFound(t *testing.T) {
	mock := new(MockMessage)
	mock.callback = func(section *imap.BodySectionName) imap.Literal {
		require.NotNil(t, section)
		return CreateMail([]string{"text/plain","txt"})
	}

	msg := newMsgPromise(mock,0, nil)
	attachments := msg.GetAttachment(mimeType)
	require.NotNil(t, attachments)
	require.Len(t, attachments, 0)
}

func msgPromiseGetAttachmentEmpty(t *testing.T) {
	mock := new(MockMessage)
	mock.callback = func(section *imap.BodySectionName) imap.Literal {
		require.NotNil(t, section)
		return CreateMail()
	}

	msg := newMsgPromise(mock,0, nil)
	attachments := msg.GetAttachment(mimeType)
	require.NotNil(t, attachments)
	require.Len(t, attachments, 0)
}

func msgPromiseGetAttachmentNoLiteral(t *testing.T) {
	mock := new(MockMessage)
	mock.callback = func(section *imap.BodySectionName) imap.Literal {
		require.NotNil(t, section)
		return nil
	}

	msg := newMsgPromise(mock,0, nil)
	attachments := msg.GetAttachment(mimeType)
	require.NotNil(t, attachments)
	require.Len(t, attachments, 0)
}

func CreateMail(mimeTypes ...[]string) imap.Literal {
	header := new(mail.Header)
	header.SetAddressList("From", []*mail.Address{{Name: "Foo Bar  DE", Address:"foo@bar.de"}})
	header.SetAddressList("To", []*mail.Address{{Name:"Foo Bar COM",  Address:"foo@bar.com"}})
	header.SetSubject("Foo.Bar")
	header.Set("MIME-Version", "1.0")
	header.Set("Content-Type", "TEXT/PLAIN; CHARSET=US-ASCII")
	header.Set("Date", "Wed, 17 Jul 1996 02:23:25 -0700 (PDT)")

	buf := bytes.NewBuffer([]byte{})
	mailW, _ := mail.CreateWriter(buf, *header)

	for _, mimeType := range mimeTypes {
		attHeader := mail.AttachmentHeader{}
		attHeader.SetContentType(mimeType[0], map[string]string{})
		attHeader.SetFilename("test."+mimeType[1])
		attWriter, _ := mailW.CreateAttachment(attHeader)
		_, _ = attWriter.Write([]byte{0, 0, 0, 0})
	}
	return buf
}
