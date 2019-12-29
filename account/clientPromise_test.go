package account

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClientPromise(t *testing.T) {
	t.Run("Append", clientPromiseAppend)
	t.Run("Delete", clientPromiseDelete)
	t.Run("Expunge", clientPromiseExpunge)
	t.Run("Fetch", clientPromiseFetch)
	t.Run("List", clientPromiseList)
	t.Run("Logout", clientPromiseLogout)
	t.Run("Move", clientPromiseMove)
	t.Run("Search", clientPromiseSearch)
	t.Run("Select", clientPromiseSelect)
	t.Run("State", clientPromiseState)
	t.Run("Store", clientPromiseStore)
}

func clientPromiseAppend(t *testing.T) {
	defer func() { _ = recover() }()
	c := clientPromise{}
	_ = c.Append("", []string{}, time.Now(), nil)
	require.Fail(t, "nil pointer???")
}

func clientPromiseDelete(t *testing.T) {
	defer func() { _ = recover() }()
	c := clientPromise{}
	_ = c.Delete(0)
	require.Fail(t, "nil pointer???")
}

func clientPromiseExpunge(t *testing.T) {
	defer func() { _ = recover() }()
	c := clientPromise{}
	_ = c.Expunge(nil)
	require.Fail(t, "nil pointer???")
}

func clientPromiseFetch(t *testing.T) {
	defer func() { _ = recover() }()
	c := clientPromise{}
	_ = c.Fetch(nil, nil, nil)
	require.Fail(t, "nil pointer???")
}
func clientPromiseList(t *testing.T) {
	defer func() { _ = recover() }()
	c := clientPromise{}
	_ = c.List("", "", nil)
	require.Fail(t, "nil pointer???")
}
func clientPromiseLogout(t *testing.T) {
	defer func() { _ = recover() }()
	c := clientPromise{}
	_ = c.Logout()
	require.Fail(t, "nil pointer???")
}
func clientPromiseMove(t *testing.T) {
	defer func() { _ = recover() }()
	c := clientPromise{}
	_ = c.Move(nil, "")
	require.Fail(t, "nil pointer???")
}
func clientPromiseSearch(t *testing.T) {
	defer func() { _ = recover() }()
	c := clientPromise{}
	_, _ = c.Search(nil)
	require.Fail(t, "nil pointer???")
}
func clientPromiseSelect(t *testing.T) {
	defer func() { _ = recover() }()
	c := clientPromise{}
	_, _ = c.Select("", false)
	require.Fail(t, "nil pointer???")
}
func clientPromiseState(t *testing.T) {
	defer func() { _ = recover() }()
	c := clientPromise{}
	_ = c.State()
	require.Fail(t, "nil pointer???")
}
func clientPromiseStore(t *testing.T) {
	defer func() { _ = recover() }()
	c := clientPromise{}
	_ = c.Store(nil, "", nil, nil)
	require.Fail(t, "nil pointer???")
}
