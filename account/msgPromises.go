package account

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/charset"
	e "mailAssistant/errors"
	"strings"
)

func init() {
	message.CharsetReader = charset.Reader
}

// MsgPromises is collection of messages
type MsgPromises struct {
	*ImapPromise
	messages []*MsgPromise
	seqSet   *imap.SeqSet
}

// Each iterate over all messages
func (p MsgPromises) Each(callback func(promise *MsgPromise)) {
	for _, message := range p.messages {
		callback(message)
	}
}

// Message is returning one specific message
func (p MsgPromises) Message(i int) *MsgPromise {
	return p.messages[i]
}

// Move is moving all messages to a mailbox
func (p MsgPromises) Move(moveTo string) (int, error) {
	if p.seqSet.Empty() {
		return 0, e.NewEmpty()
	}
	path := strings.ReplaceAll(moveTo, "/", ".")
	if path == "" {
		path = "INBOX"
	}
	return countSeqSet(p.seqSet), p.client.Move(p.seqSet, path)
}

func countSeqSet(set *imap.SeqSet) int {
	sum := 0
	for _, k := range set.Set {
		sum += int(k.Stop-k.Start) + 1
	}
	return sum
}

// Delete is deleting all messages
func (p MsgPromises) Delete() (int, error) {
	storeItem := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.DeletedFlag}
	if p.seqSet.Empty() {
		return 0, e.NewEmpty()
	}
	if err := p.client.Store(p.seqSet, storeItem, flags, nil); err != nil {
		return 0, err
	} else if err := p.client.Expunge(nil); err != nil {
		return 0, err
	} else {
		return countSeqSet(p.seqSet), nil
	}
}

// SetSeen is setting the Seen flag to all messages
func (p MsgPromises) SetSeen() (int, error) {
	if p.seqSet.Empty() {
		return 0, e.NewEmpty()
	}
	storeItem := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.SeenFlag}
	if err := p.client.Store(p.seqSet, storeItem, flags, nil); err != nil {
		return 0, err
	}
	return countSeqSet(p.seqSet), nil
}

// GetAttachments is collecting all attachments with mimeType
func (p MsgPromises) GetAttachments(mimeType string) []*AttachmentPromise {
	result := make([]*AttachmentPromise, 0)
	p.Each(func(promise *MsgPromise) {
		result = append(result, promise.GetAttachment(mimeType)...)
	})
	return result
}

func (p *MsgPromises) add(msg *imap.Message) {
	p.messages = append(p.messages, newMsgPromise(msg, msg.SeqNum, p.client))
}

func (p *MsgPromises) addAll(messages chan *imap.Message) {
	for msg := range messages {
		p.add(msg)
	}
}

// Expunge is sending the expunge command to the IMAP server
func (p MsgPromises) Expunge() {
	_ = p.client.Expunge(nil)
}
