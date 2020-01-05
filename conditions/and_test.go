package conditions

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAnd(t *testing.T) {
	t.Run("init", andInit)
	t.Run("add", func(t *testing.T) {
		t.Run("notLocked", andAddUnLocked)
		t.Run("Locked", andAddLocked)
	})
	t.Run("get", andGet)
	t.Run("parseYaml", andParseYaml)
	t.Run("SetCursor", func(t *testing.T) {
		t.Run("no parent", andSetCursorNoParent)
		t.Run("with parent", andSetCursorWithParent)
	})
	t.Run("String", func(t *testing.T) {
		t.Run("no conditions", andString0)
		t.Run("one conditions", andString1)
		t.Run("some conditions", andString4)
	})
}

func andInit(t *testing.T) {
	a := newAnd()
	require.NotNil(t, a.conditions)
	require.False(t, *a.locked)
	require.NotNil(t, a.parent)
	require.Nil(t, a.parent.p)
}
func andAddUnLocked(t *testing.T) {
	a := newAnd()

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

func andAddLocked(t *testing.T) {
	a := newAnd()

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

func andGet(t *testing.T) {
	a := newAnd()
	require.Equal(t, []interface{}{}, a.Get())
	a.Add(newPair("a", "b"))
	require.Equal(t, []interface{}{"A", "b"}, a.Get())
}

func andParseYaml(t *testing.T) {
	in := map[string]interface{}{"field": "from", "value": "b"}
	a := newAnd()
	a.ParseYaml(in)

	require.NotNil(t, a.parent)
	require.Nil(t, a.parent.p)
	require.False(t, a.parent.HasParent())
	require.Len(t, *a.conditions, 1)
	require.False(t, *a.locked)
	get := a.Get()
	require.Equal(t, []interface{}{"FROM", "b"}, get)
	str := a.String()
	require.Equal(t, str, "FROM='b'")
}

func andSetCursorNoParent(t *testing.T) {
	a := newAnd()
	require.Equal(t, []interface{}{}, a.Get())
	a.Add(newPair("a", "b"))
	require.Equal(t, []interface{}{"A", "b"}, a.Get())
	a.SetCursor()
	require.Equal(t, []interface{}{"CURSOR"}, a.Get())
}

func andSetCursorWithParent(t *testing.T) {
	parent := newAnd()
	a := newAnd()
	parent.Add(a)
	a.Add(newPair("a", "b"))

	require.Equal(t, []interface{}{"A", "b"}, a.Get())
	require.Len(t, *parent.conditions, 1)
	require.Equal(t, (*parent.conditions)[0], a)

	a.SetCursor()
	require.Equal(t, []interface{}{"A", "b"}, a.Get())
	require.False(t, *a.locked)

	require.Len(t, *parent.conditions, 1)
	require.NotEqual(t, (*parent.conditions)[0], a)
	require.Equal(t, []interface{}{"CURSOR"}, parent.Get())
	require.True(t, *parent.locked)
}

func andString0(t *testing.T) {
	a := newAnd()
	require.Equal(t, "", a.String())
}

func andString1(t *testing.T) {
	a := newAnd()
	require.Equal(t, "", a.String())
	a.Add(newPair("a", "b"))
	require.Equal(t, "A='b'", a.String())
}

func andString4(t *testing.T) {
	a := newAnd()
	require.Equal(t, "", a.String())
	a.Add(newPair("a", "b"))
	a.Add(newPair("c", "d"))
	a.Add(newPair("e", "f"))
	a.Add(newPair("g", "h"))

	require.Equal(t, "( A='b' and C='d' and E='f' and G='h' )", a.String())
}