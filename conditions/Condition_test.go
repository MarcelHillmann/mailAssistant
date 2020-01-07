package conditions

import (
	"github.com/emersion/go-imap"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCondition(t *testing.T) {
	t.Run("ParseYaml", conditionParseYaml)
	t.Run("internal parseYaml", func(t *testing.T) {
		t.Run("cursor", conditionInternalParseYamlCursor)
		t.Run("or", conditionInternalParseYamlOr)
		t.Run("and", conditionInternalParseYamlAnd)
		t.Run("not", conditionInternalParseYamlNot)
		t.Run("older", conditionInternalParseYamlOlder)
		t.Run("younger", conditionInternalParseYamlYounger)
		t.Run("list", conditionInternalParseYamlList)
		t.Run("INVALID", conditionInternalParseYamlInvalid)
	})
	t.Run("allowedYamlKey", conditionAllowedYamlKey)
	t.Run("validImapKeyword", conditionValidImapKeyword)
}

func conditionParseYaml(t *testing.T) {
	a := ParseYaml(nil)
	require.NotNil(t, a)
	a2, ok := a.(and)
	require.True(t, ok)
	require.NotNil(t, a2.conditions)
	require.False(t, *a2.locked)
	require.NotNil(t, a2.parent)
	require.False(t, a2.parent.HasParent())

	require.Panics(t, func() {
		ParseYaml("a")
	})
}

func conditionInternalParseYamlCursor(t *testing.T) {
	cond := newAnd()
	item := make(map[string]interface{})
	item["field"] = CURSOR
	require.NotPanics(t, func() {
		parseYaml(item, cond)
	})
	require.True(t, *cond.locked)
	require.Len(t, *cond.conditions, 1)
	require.Equal(t, []interface{}{CURSOR}, cond.Get())
	require.Equal(t, CURSOR+"='<nil>'", cond.String())
}

func conditionInternalParseYamlOr(t *testing.T) {
	cond := newAnd()
	item := make(map[string]interface{})
	item["field"] = "or"
	item["value"] = []interface{}{map[interface{}]interface{}{"field": "from", "value": "test"}}
	require.NotPanics(t, func() {
		parseYaml(item, cond)
	})
	require.False(t, *cond.locked)
	require.Len(t, *cond.conditions, 1)
	require.Equal(t, []interface{}{"or", "FROM", "test"}, cond.Get())
	require.Equal(t, "or FROM='test'", cond.String())
}

func conditionInternalParseYamlAnd(t *testing.T) {
	cond := newAnd()
	item := make(map[string]interface{})
	item["field"] = "and"
	item["value"] = []interface{}{map[interface{}]interface{}{"field": "from", "value": "test"}}
	require.NotPanics(t, func() {
		parseYaml(item, cond)
	})
	require.False(t, *cond.locked)
	require.Len(t, *cond.conditions, 1)
	require.Equal(t, []interface{}{"FROM", "test"}, cond.Get())
	require.Equal(t, "FROM='test'", cond.String())
}

func conditionInternalParseYamlNot(t *testing.T) {
	cond := newAnd()
	item := make(map[string]interface{})
	item["field"] = "not"
	item["value"] = []interface{}{map[interface{}]interface{}{"field": "from", "value": "test"}}
	require.NotPanics(t, func() {
		parseYaml(item, cond)
	})
	require.False(t, *cond.locked)
	require.Len(t, *cond.conditions, 1)
	require.Equal(t, []interface{}{"not", "FROM", "test"}, cond.Get())
	require.Equal(t, "not FROM='test'", cond.String())
}

func conditionInternalParseYamlOlder(t *testing.T) {
	older := time.Now().UnixNano() - int64(24 * time.Hour)
	assertDate := time.Unix(0, older).Format(imap.DateLayout)
	cond := newAnd()
	item := make(map[string]interface{})
	item["field"] = "older"
	item["value"] = "P1D"
	require.NotPanics(t, func() {
		parseYaml(item, cond)
	})
	require.False(t, *cond.locked)
	require.Len(t, *cond.conditions, 1)
	require.Equal(t, []interface{}{"BEFORE", assertDate}, cond.Get())
	require.Equal(t, "BEFORE='"+assertDate+"'", cond.String())
}

func conditionInternalParseYamlYounger(t *testing.T) {
	older := time.Now().UnixNano() - int64(24*time.Hour)
	assertDate := time.Unix(0,older).Format(imap.DateLayout)
	cond := newAnd()
	item := make(map[string]interface{})
	item["field"] = "younger"
	item["value"] = "P1D"
	require.NotPanics(t, func() {
		parseYaml(item, cond)
	})
	require.False(t, *cond.locked)
	require.Len(t, *cond.conditions, 1)
	require.Equal(t, []interface{}{"SINCE", assertDate}, cond.Get())
	require.Equal(t, "SINCE='"+assertDate+"'", cond.String())
}

func conditionInternalParseYamlList(t *testing.T) {
	cond := newAnd()
	item := []interface{}{
		map[interface{}]interface{}{"field": "from", "value": "foo"},
		map[interface{}]interface{}{"field": "to", "value": "bar"},
	}
	require.NotPanics(t, func() {
		parseYaml(item, cond)
	})
	require.False(t, *cond.locked)
	require.Len(t, *cond.conditions, 2)
	require.Equal(t, []interface{}{"FROM", "foo", "TO", "bar"}, cond.Get())
	require.Equal(t, "( FROM='foo' and TO='bar' )", cond.String())
}

func conditionInternalParseYamlInvalid(t *testing.T) {
	cond := newAnd()
	item := []interface{}{
		map[interface{}]interface{}{"field": "from", "value": "foo"},
		map[interface{}]interface{}{"field": "XXXX", "value": "bar"},
	}
	require.Panics(t, func() {
		parseYaml(item, cond)
	})
}
func conditionAllowedYamlKey(t *testing.T) {
	for _, key := range []string{"field", "value"} {
		t.Run(key, func(t *testing.T) {
			require.NotPanics(t, func() {
				allowedYamlKey(map[interface{}]interface{}{key: ""})
			})
		})
	}

	t.Run("INVALID", func(t *testing.T) {
		require.Panics(t, func() {
			allowedYamlKey(map[interface{}]interface{}{"zzzzz": ""})
		})
	})
}
func conditionValidImapKeyword(t *testing.T) {
	for _, key := range keywords {
		t.Run(key, func(t *testing.T) {
			require.True(t, validImapKeyword(key))
		})
	}
	t.Run("INVALID", func(t *testing.T) {
		require.False(t, validImapKeyword("XXX"))
	})
}
