package account

import (
	"github.com/emersion/go-imap"
	move "github.com/emersion/go-imap-move"
	def "github.com/emersion/go-imap/client"
	"io"
	"time"
)

type clientPromise struct {
	client   *def.Client
	mvClient *move.Client
}

func NewClientPromise(client *def.Client) IClient {
	return clientPromise{client, move.NewClient(client)}
}
func (promise clientPromise) Append(saveTo string, flags []string, date time.Time, msg imap.Literal) error {
	return promise.client.Append(saveTo, flags, date, msg)
}
func (promise clientPromise) Expunge(ch chan uint32) error {
	return promise.client.Expunge(ch)
}
func (promise clientPromise) Fetch(seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
	return promise.client.Fetch(seqset, items, ch)
}
func (promise clientPromise) List(ref, name string, ch chan *imap.MailboxInfo) error {
	return promise.client.List(ref, name, ch)
}
func (promise clientPromise) Login(username, password string) error {
	return promise.client.Login(username, password)
}
func (promise clientPromise) Logout() error {
	return promise.client.Logout()
}
func (promise clientPromise) Move(seqSet *imap.SeqSet, dest string) error {
	return promise.mvClient.MoveWithFallback(seqSet, dest)
}
func (promise clientPromise) Search(criteria *imap.SearchCriteria) (seqNums []uint32, err error) {
	return promise.client.Search(criteria)
}
func (promise clientPromise) SetDebug(w io.Writer) {
	promise.client.SetDebug(w)
}
func (promise clientPromise) Select(name string, readOnly bool) (*imap.MailboxStatus, error) {
	return promise.client.Select(name, readOnly)
}
func (promise clientPromise) Store(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
	return promise.client.Store(seqSet, item, value, ch)
}
func (promise clientPromise) State() imap.ConnState {
	return promise.client.State()
}
func (promise clientPromise) Delete(num uint32) error {
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(num)
	flags := []interface{}{imap.DeletedFlag}

	return promise.Store(seqSet,"+FLAGS.SILENT", flags, nil)
}