package conditions

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDummy(t *testing.T) {
	t.Run("Add", dummyAdd)
	t.Run("Get", dummyGet)
	t.Run("Parent", dummyParent)
	t.Run("ParseYaml", dummyParseYaml)
	t.Run("SetCursor", dummySetCursor)
	t.Run("String", dummyString)
}

func dummyAdd(t *testing.T) {
	d := dummy{}
	require.Panics(t, func() { d.Add(nil) })
}

func dummyGet(t *testing.T) {
	d := dummy{}
	require.Equal(t, []interface{}{}, d.Get())
}

func dummyParent(t *testing.T) {
	d := dummy{}
	require.Panics(t, func() { d.Parent(nil) })
}

func dummyParseYaml(t *testing.T) {
	d := dummy{}
	require.Panics(t, func() { d.ParseYaml(nil) })
}

func dummySetCursor(t *testing.T) {
	d := dummy{}
	require.Panics(t, func() { d.SetCursor() })
}

func dummyString(t *testing.T) {
	d := dummy{}
	require.Equal(t, "", d.String())
}
