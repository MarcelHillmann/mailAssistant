package account

import (
	"errors"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"io"
	"io/ioutil"
	"mailAssistant/logging"
	"strings"
)

var errNextPartEOF = errors.New("multipart: NextPart: EOF")

func init() {
	message.CharsetReader = charset.Reader
}

// MsgPromise repesents a IMAP message and cover the underlying framework
type MsgPromise struct {
	IMessage
	IClient
	seqNum uint32
	logger logging.Logger
}

func newMsgPromise(msg IMessage, seqNum uint32, client IClient) *MsgPromise {
	log := logging.NewLogger()
	return &MsgPromise{msg, client,seqNum, log}
}

func (promise MsgPromise) getLogger() logging.Logger {
	return promise.logger
}

// GetLiteral returns the hole mail
func (promise *MsgPromise) GetLiteral() (literal imap.Literal){
	var section imap.BodySectionName
	literal = promise.GetBody(&section)
	if literal == nil {
		promise.getLogger().Warn("Server didn't returned message body")
	}
	return
}

// GetAttachment is searching for a attachment with mimeType
func (promise MsgPromise) GetAttachment(mimeTypeString string) []*AttachmentPromise {
	lowerMimeType := strings.TrimSpace(strings.ToLower(mimeTypeString))
	result := make([]*AttachmentPromise, 0)

	literal := promise.GetLiteral()
	if literal == nil {
		return result
	}

	reader, err := mail.CreateReader(literal)
	if err != nil {
		promise.getLogger().Severe(err)
		return result
	}

	for {
		part, err := nextPart(reader)
		if err != nil && (err == io.EOF || err.Error() == errNextPartEOF.Error()) {
			break
		} else if err != nil {
			panic(err)
		}

		switch header := part.Header.(type) {
		case *mail.AttachmentHeader:
			contentType, _, _ := header.ContentType()
			disposition, _, _ := header.ContentDisposition()
			filename, _ := header.Filename()

			if lowerMimeType != "" && //
				contentType == lowerMimeType || //
				mimeTypeString == "" {
				params := make(map[string]string)
				params["filename"] = filename
				params["content-disposition"] = disposition
				params["content-type"] = contentType
				content, _ := ioutil.ReadAll(part.Body)
				result = append(result, &AttachmentPromise{content, params})
			}
		}
	}

	return result
}

// DeletePromise is setting the delete flag to the represented message
func (promise MsgPromise) DeletePromise(callback func(err error)){
	callback(promise.IClient.Delete(promise.seqNum))
}

var nextPart  = internalNextPart
func internalNextPart(reader *mail.Reader) (*mail.Part, error) {
	return reader.NextPart()
}
