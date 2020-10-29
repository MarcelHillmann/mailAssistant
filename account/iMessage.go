package account

import "github.com/emersion/go-imap"

// IMessage represents the collection of all needed methods
type IMessage interface {
	GetBody(section *imap.BodySectionName) imap.Literal
}
