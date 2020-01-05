package conditions

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOr(t *testing.T) {
	t.Run("init", orInit)
	t.Run("add", func(t *testing.T) {
		t.Run("notLocked", orAddUnLocked)
		t.Run("Locked", orAddLocked)
	})
	t.Run("get", orGet)
	t.Run("parseYaml", orParseYaml)
	t.Run("SetCursor", func(t *testing.T) {
		t.Run("no parent", orSetCursorNoParent)
		t.Run("with parent", orSetCursorWithParent)
	})
	t.Run("String", func(t *testing.T) {
		t.Run("no conditions", orString0)
		t.Run("one conditions", orString1)
		t.Run("some conditions", orString4)
	})
}

func orInit(t *testing.T) {
	a := newOr()
	require.NotNil(t, a.conditions)
	require.False(t, *a.locked)
	require.NotNil(t, a.parent)
	require.Nil(t, a.parent.p)
}
func orAddUnLocked(t *testing.T) {
	a := newOr()

	require.NotNil(t, a.conditions)
	require.NotNil(t, a.locked)
	require.NotNil(t, a.parent)
	require.False(t, a.parent.HasParent())
	a.Add(newPair("hugo", "boss"))
	require.NotNil(t, a.parent)
	require.Nil(t, a.parent.p)
	require.NotNil(t, a.conditions)
	require.NotNil(t, a.locked)
	require.Len(t, *a.conditions, 1)
	pair := (*a.conditions)[0].(pair)
	require.NotNil(t, pair.parent)
	require.Equal(t, "hugo", pair.keyval.field)
	require.Equal(t, "boss", pair.keyval.value)
	require.True(t, a == pair.parent.p)
}

func orAddLocked(t *testing.T) {
	a := newOr()

	require.NotNil(t, a.conditions)
	require.False(t, *a.locked)
	require.NotNil(t, a.parent)
	require.False(t, a.parent.HasParent())

	a.SetCursor()

	require.NotNil(t, a.conditions)
	require.True(t, *a.locked)
	require.NotNil(t, a.parent)
	require.False(t, a.parent.HasParent())

	a.Add(newPair("hugo", "boss"))
	require.Len(t, *a.conditions, 1)
	pair := (*a.conditions)[0].(pair)
	require.NotNil(t, pair.parent)
	require.Equal(t, "CURSOR", pair.keyval.field)
	require.Nil(t,  pair.keyval.value)
}

func orGet(t *testing.T) {
	a := newOr()
	require.Equal(t, []interface{}{}, a.Get())
	a.Add(newPair("a", "b"))
	require.Equal(t, []interface{}{"or","A", "b"}, a.Get())
}

func orParseYaml(t *testing.T) {
	in := map[string]interface{}{"field": "from", "value": "b"}
	a := newOr()
	a.ParseYaml(in)

	require.NotNil(t, a.parent)
	require.Nil(t, a.parent.p)
	require.False(t, a.parent.HasParent())
	require.Len(t, *a.conditions, 1)
	require.False(t, *a.locked)
	get := a.Get()
	require.Equal(t, []interface{}{"or", "FROM", "b"}, get)
	str := a.String()
	require.Equal(t, str, "or FROM='b'")
}

func orSetCursorNoParent(t *testing.T) {
	a := newOr()
	require.Equal(t, []interface{}{}, a.Get())
	a.Add(newPair("a", "b"))
	require.Equal(t, []interface{}{"or","A", "b"}, a.Get())
	a.SetCursor()
	require.Equal(t, []interface{}{"CURSOR"}, a.Get())
}

func orSetCursorWithParent(t *testing.T) {
	parent := newOr()
	a := newOr()
	parent.Add(a)
	a.Add(newPair("a", "b"))

	require.Equal(t, []interface{}{"or","A", "b"}, a.Get())
	require.Len(t, *parent.conditions, 1)
	require.Equal(t, (*parent.conditions)[0], a)

	a.SetCursor()
	require.Equal(t, []interface{}{"or","A", "b"}, a.Get())
	require.False(t, *a.locked)

	require.Len(t, *parent.conditions, 1)
	require.NotEqual(t, (*parent.conditions)[0], a)
	require.Equal(t, []interface{}{"CURSOR"}, parent.Get())
	require.True(t, *parent.locked)
}

func orString0(t *testing.T) {
	a := newOr()
	require.Equal(t, "", a.String())
}

func orString1(t *testing.T) {
	a := newOr()
	require.Equal(t, "", a.String())
	a.Add(newPair("a", "b"))
	require.Equal(t, "or A='b'", a.String())
}

func orString4(t *testing.T) {
	a := newOr()
	require.Equal(t, "", a.String())
	a.Add(newPair("a", "b"))
	a.Add(newPair("c", "d"))
	a.Add(newPair("e", "f"))
	a.Add(newPair("g", "h"))

	require.Equal(t, "( A='b' or C='d' or E='f' or G='h' )", a.String())
}
