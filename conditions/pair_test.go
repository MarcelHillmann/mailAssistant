package conditions

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPair(t *testing.T) {
	t.Run("add", pairAdd)
	t.Run("get", func(t *testing.T) {
		t.Run("string", pairGetString)
		t.Run("int", pairGetInt)
		t.Run("nil", pairGetNil)
	})
	t.Run("parseYaml", pairParseYaml)
	t.Run("SetCursor", func(t *testing.T) {
		t.Run("no parent", pairSetCursorNoParent)
		t.Run("with parent", pairSetCursorWithParent)
	})
	t.Run("String", pairString)
}

func pairAdd(t *testing.T) {
	a := newPair("", "")
	require.Panics(t, func() {
		a.Add(newAnd())
	})
}

func pairGetString(t *testing.T) {
	a := newPair("a", "b")
	require.Equal(t, []interface{}{"A", "b"}, a.Get())
}

func pairGetInt(t *testing.T) {
	a := newPair("a", 32)
	require.Equal(t, []interface{}{"A", uint32(32)}, a.Get())
}

func pairGetNil(t *testing.T) {
	a := newPair("a", nil)
	require.Equal(t, []interface{}{"A"}, a.Get())
}

func pairParseYaml(t *testing.T) {
	in := map[string]interface{}{"field": "from", "value": "b"}
	a := newPair("", "")

	require.Panics(t, func() {
		a.ParseYaml(in)
	})

}

func pairSetCursorNoParent(t *testing.T) {
	a := newPair("a", "b")
	require.Equal(t, []interface{}{"A", "b"}, a.Get())
	a.SetCursor()
	require.Equal(t, []interface{}{CURSOR}, a.Get())
}

func pairSetCursorWithParent(t *testing.T) {
	parent := newAnd()
	a := newPair("a", "b")
	parent.Add(a)

	require.Equal(t, []interface{}{"A", "b"}, a.Get())
	require.Len(t, *parent.conditions, 1)
	require.Equal(t, (*parent.conditions)[0], a)

	a.SetCursor()
	require.Equal(t, []interface{}{"A", "b"}, a.Get())
	require.Len(t, *parent.conditions, 1)
	require.NotEqual(t, (*parent.conditions)[0], a)
	require.Equal(t, []interface{}{CURSOR}, parent.Get())
	require.True(t, *parent.locked)
}

func pairString(t *testing.T) {
	a := newPair("a", "b")
	require.Equal(t, "A='b'", a.String())
	a.keyval.value = nil
	require.Equal(t, "A='<nil>'", a.String())
}
