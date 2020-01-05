package conditions

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNot(t *testing.T) {
	t.Run("init", notInit)
	t.Run("add", func(t *testing.T) {
		t.Run("notLocked", notAddUnLocked)
		t.Run("Locked", notAddLocked)
	})
	t.Run("get", notGet)
	t.Run("parseYaml", notParseYaml)
	t.Run("SetCursor", func(t *testing.T) {
		t.Run("no parent", notSetCursorNoParent)
		t.Run("with parent", notSetCursorWithParent)
	})
	t.Run("String", func(t *testing.T) {
		t.Run("no conditions", notString0)
		t.Run("one conditions", notString1)
		t.Run("some conditions", notString4)
		t.Run("extended", notStringEx)
	})
}

func notInit(t *testing.T) {
	a := newNot()
	require.NotNil(t, a.conditions)
	require.False(t, *a.locked)
	require.NotNil(t, a.parent)
	require.Nil(t, a.parent.p)
}
func notAddUnLocked(t *testing.T) {
	a := newNot()

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

func notAddLocked(t *testing.T) {
	a := newNot()

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

func notGet(t *testing.T) {
	a := newNot()
	require.Equal(t, []interface{}{}, a.Get())
	a.Add(newPair("a", "b"))
	require.Equal(t, []interface{}{"not","A", "b"}, a.Get())
}

func notParseYaml(t *testing.T) {
	in := map[string]interface{}{"field": "from", "value": "b"}
	a := newNot()
	a.ParseYaml(in)

	require.NotNil(t, a.parent)
	require.Nil(t, a.parent.p)
	require.False(t, a.parent.HasParent())
	require.Len(t, *a.conditions, 1)
	require.False(t, *a.locked)
	require.Equal(t, []interface{}{"not", "FROM", "b"}, a.Get())
	require.Equal(t,  "not FROM='b'",a.String())
}

func notSetCursorNoParent(t *testing.T) {
	a := newNot()
	require.Equal(t, []interface{}{}, a.Get())
	a.Add(newPair("a", "b"))
	require.Equal(t, []interface{}{"not","A", "b"}, a.Get())
	a.SetCursor()
	require.Equal(t, []interface{}{"CURSOR"}, a.Get())
}

func notSetCursorWithParent(t *testing.T) {
	parent := newNot()
	a := newNot()
	parent.Add(a)
	a.Add(newPair("a", "b"))

	require.Equal(t, []interface{}{"not","A", "b"}, a.Get())
	require.Len(t, *parent.conditions, 1)
	require.Equal(t, (*parent.conditions)[0], a)

	a.SetCursor()
	require.Equal(t, []interface{}{"not","A", "b"}, a.Get())
	require.False(t, *a.locked)

	require.Len(t, *parent.conditions, 1)
	require.NotEqual(t, (*parent.conditions)[0], a)
	require.Equal(t, []interface{}{"CURSOR"}, parent.Get())
	require.True(t, *parent.locked)
}

func notString0(t *testing.T) {
	a := newNot()
	require.Equal(t, "", a.String())
}

func notString1(t *testing.T) {
	a := newNot()
	require.Equal(t, "", a.String())
	a.Add(newPair("a", "b"))
	require.Equal(t, "not A='b'", a.String())
}

func notString4(t *testing.T) {
	a := newNot()
	require.Equal(t, "", a.String())

	a.Add(newPair("a", "b"))
	a.Add(newPair("c", "d"))
	a.Add(newPair("e", "f"))
	a.Add(newPair("g", "h"))

	require.Equal(t, "not ( A='b' and C='d' and E='f' and G='h' )", a.String())
}

func notStringEx(t *testing.T) {
	a := newNot()
	require.Equal(t, "", a.String())

	and := newAnd()
	and.Add(newPair("a", "b"))
	and.Add(newPair("c", "d"))
	a.Add(and)

	or:=newOr()
	a.Add(or)
	or.Add(newPair("e", "f"))
	or.Add(newPair("g", "h"))

	require.Equal(t, "not ( ( A='b' and C='d' ) and ( E='f' or G='h' ) )", a.String())
}