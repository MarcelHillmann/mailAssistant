package account

import (
	"github.com/emersion/go-imap"
	"io"
	"time"
)

// IClient represents the collection of all needed methods
type IClient interface {
	Append(saveTo string, flags []string, date time.Time, msg imap.Literal) error
	// B
	// C
	// D
	Delete(num uint32) error
	// E
	Expunge(ch chan uint32) error
	// F
	Fetch(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error
	// G
	// H
	// I
	// J
	// K
	// L
	List(ref, name string, ch chan *imap.MailboxInfo) error
	Login(username string, password string) error
	Logout() error
	// M
	Move(seqSet *imap.SeqSet, dest string) error
	// N
	// O
	// P
	// Q
	// R
	// S
	Search(criteria *imap.SearchCriteria) (seqNums []uint32, err error)
	SetDebug(writer io.Writer)
	Select(name string, readOnly bool) (*imap.MailboxStatus, error)
	Store(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error
	State() imap.ConnState
	// T
	// U
	// V
	// W
	// X
	// Y
	// Z

}
