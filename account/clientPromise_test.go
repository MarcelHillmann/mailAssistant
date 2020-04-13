package account

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestClientPromise(t *testing.T) {
	t.Run("Factory", clientPromiseFactory)
	t.Run("Append", clientPromiseAppend)
	t.Run("Delete", clientPromiseDelete)
	t.Run("Expunge", clientPromiseExpunge)
	t.Run("Fetch", clientPromiseFetch)
	t.Run("Login", clientPromiseLogin)
	t.Run("List", clientPromiseList)
	t.Run("Logout", clientPromiseLogout)
	t.Run("Move", clientPromiseMove)
	t.Run("Search", clientPromiseSearch)
	t.Run("Select", clientPromiseSelect)
	t.Run("SetDebug", clientPromiseSetDebug)
	t.Run("State", clientPromiseState)
	t.Run("Store", clientPromiseStore)

}

func clientPromiseFactory(t *testing.T) {
	c := NewClientPromise(nil)
	require.NotNil(t, c)
	c2 := c.(clientPromise)
	require.Nil(t, c2.client)
	require.NotNil(t, c2.mvClient)
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

func clientPromiseLogin(t *testing.T) {
	defer func() { _ = recover() }()
	c := clientPromise{}
	_ = c.Login("","")
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

func clientPromiseSetDebug(t *testing.T) {
	defer func() { _ = recover() }()
	c := clientPromise{}
	c.SetDebug(os.Stdout)
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