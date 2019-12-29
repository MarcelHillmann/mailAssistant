package account

import (
	"fmt"
	"github.com/emersion/go-imap"
	"io"
	"time"
)

type callCollector struct {
	append,
	delete,
	expunge,
	fetch,
	list,
	login,
	logout,
	move,
	search,
	setdebug,
	sel,
	store,
	state int
}

type MockClientPromise struct {
	called           *callCollector
	AppendCallback   func(mBox string, flags []string, date time.Time, msg imap.Literal) error
	DeleteCallback   func(num uint32) error
	ExpungeCallback  func(ch chan uint32) error
	FetchCallback    func(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error
	ListCallback     func(ref, name string, ch chan *imap.MailboxInfo) error
	LoginCallback    func(username, password string) error
	LogoutCallback   func() error
	MoveCallback     func(seqSet *imap.SeqSet, dest string) error
	SearchCallback   func(*imap.SearchCriteria) ([]uint32, error)
	SetDebugCallback func(w io.Writer)
	SelectCallback   func(name string, readOnly bool) (*imap.MailboxStatus, error)
	StoreCallback    func(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error
	StateCallback    func() imap.ConnState
}

func (mock MockClientPromise) Append(mBox string, flags []string, date time.Time, msg imap.Literal) error {
	mock.called.append++
	return mock.AppendCallback(mBox, flags, date, msg)
}
func (mock MockClientPromise) Delete(num uint32) error {
	mock.called.delete++
	return mock.DeleteCallback(num)
}
func (mock MockClientPromise) Expunge(ch chan uint32) error {
	mock.called.expunge++
	return mock.ExpungeCallback(ch)
}
func (mock MockClientPromise) Fetch(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
	mock.called.fetch++
	return mock.FetchCallback(seqSet, items, ch)
}
func (mock MockClientPromise) List(ref, name string, ch chan *imap.MailboxInfo) error {
	mock.called.list++
	return mock.ListCallback(ref, name, ch)
}
func (mock MockClientPromise) Login(username, password string) error {
	mock.called.login++
	return mock.LoginCallback(username, password)
}
func (mock MockClientPromise) Logout() error {
	mock.called.logout++
	return mock.LogoutCallback()
}
func (mock MockClientPromise) Move(seqSet *imap.SeqSet, dest string) error {
	mock.called.move++
	return mock.MoveCallback(seqSet, dest)
}
func (mock MockClientPromise) Search(criteria *imap.SearchCriteria) (seqNums []uint32, err error) {
	mock.called.search++
	return mock.SearchCallback(criteria)
}
func (mock MockClientPromise) Select(name string, readOnly bool) (*imap.MailboxStatus, error) {
	mock.called.sel++
	return mock.SelectCallback(name, readOnly)
}
func (mock MockClientPromise) SetDebug(w io.Writer) {
	mock.called.setdebug++
	mock.SetDebugCallback(w)
}
func (mock MockClientPromise) Store(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
	mock.called.store++
	return mock.StoreCallback(seqSet, item, value, ch)
}
func (mock MockClientPromise) State() imap.ConnState {
	mock.called.state++
	return mock.StateCallback()
}

func (mock MockClientPromise) Assert() string {
	c := mock.called
	return fmt.Sprintf("%d%d%d%d%d%s%d%d%d%d%d%s%d%d%d",
		c.login, c.setdebug, c.sel, c.search, c.move, "-",
		c.state, c.list, c.fetch, c.append, c.store, "-",
		c.delete, c.expunge, c.logout)
}

func NewMockClientMinimal() *MockClientPromise {
	m := &MockClientPromise{}
	m.called = new(callCollector)
	return m
}
func NewMockClient() *MockClientPromise {
	m := NewMockClientMinimal()
	m.AppendCallback = defaultAppend
	m.DeleteCallback = defaultDelete
	m.ExpungeCallback = defaultExpunge
	m.FetchCallback = defaultFetch
	m.ListCallback = defaultList
	m.LoginCallback = defaultLogin
	m.LogoutCallback = func() error { return nil }
	m.MoveCallback = defaultMove
	m.SearchCallback = defaultSearch
	m.SetDebugCallback = func(w io.Writer) { _ = w }
	m.SelectCallback = defaultSelect
	m.StoreCallback = defaultStore
	m.StateCallback = func() imap.ConnState { return imap.ConnectedState }
	return m
}

func defaultExpunge(ch chan uint32) error {
	_ = ch
	return nil
}
func defaultFetch(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
	_, _, _ = seqSet, items, ch
	return nil
}
func defaultList(ref, name string, ch chan *imap.MailboxInfo) error {
	_, _, _ = ref, name, ch
	return nil
}
func defaultLogin(username, password string) error {
	return nil
}
func defaultMove(seqSet *imap.SeqSet, dest string) error {
	_, _ = seqSet, dest
	return nil
}
func defaultSearch(criteria *imap.SearchCriteria) (seqNums []uint32, err error) {
	_ = criteria
	return []uint32{}, nil
}
func defaultSelect(name string, readOnly bool) (*imap.MailboxStatus, error) {
	_, _ = name, readOnly
	return nil, nil
}
func defaultStore(seqSet *imap.SeqSet, item imap.StoreItem, value interface{}, ch chan *imap.Message) error {
	_, _, _, _ = seqSet, item, value, ch
	return nil
}
func defaultAppend(mBox string, flags []string, date time.Time, msg imap.Literal) error {
	_, _, _, _ = mBox, flags, date, msg
	return nil
}
func defaultDelete(num uint32) error {
	_ = num
	return nil
}
