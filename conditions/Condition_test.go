package conditions

import (
	"github.com/emersion/go-imap"
	"github.com/stretchr/testify/require"
	"net/textproto"
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
	t.Run("MapToString", conditionMapToString)
	t.Run("ToString", func(t *testing.T) {
		t.Run("0", conditionToString0)
		t.Run("1", conditionToString1)
	})
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
	older := time.Now().UnixNano() - int64(24*time.Hour)
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
	assertDate := time.Unix(0, older).Format(imap.DateLayout)
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

func conditionMapToString(t *testing.T) {
	m := make(map[string]string)
	m["a"] = "2"
	m["b"] = "1"
	m["c"] = "0"

	str := MapToString(m)
	require.Equal(t, "a: 2, b: 1, c: 0", str)
}

func conditionToString0(t *testing.T) {
	sc := new(imap.SearchCriteria)
	sc.SeqNum = &imap.SeqSet{Set: []imap.Seq{{Start: 10, Stop: 11}}}
	sc.Uid = &imap.SeqSet{Set: []imap.Seq{{Start: 100, Stop: 101}}}
	sc.Since, _ = time.Parse(time.RFC3339, "2020-10-29T23:08:00Z")
	sc.Before, _ = time.Parse(time.RFC3339, "2020-10-28T23:08:00Z")
	sc.SentSince = sc.Since
	sc.SentBefore = sc.Before
	sc.Header = textproto.MIMEHeader{}
	sc.Header.Add("Subject", "xyz")
	sc.Header.Add("To", "z")
	sc.Header.Add("From", "x")
	sc.Header.Add("X", "y")
	sc.Header.Add("Cc", "0")
	sc.Header.Add("Bcc", "1")
	sc.Smaller = 1
	sc.Larger = 400
	sc.Body = []string{"Body1", "Body2"}
	sc.Text = []string{"Text1", "Text2"}
	sc.WithFlags = []string{imap.SeenFlag, imap.AnsweredFlag, imap.FlaggedFlag, imap.DeletedFlag, imap.DraftFlag, imap.RecentFlag, string(imap.SetFlags), imap.AddFlags, imap.RemoveFlags}
	sc.WithoutFlags = sc.WithFlags
	sc.Not = []*imap.SearchCriteria{}
	sc.Or = [][2]*imap.SearchCriteria{}
	require.Equal(t, "SearchCriteria {"+
		"SeqNum: 10:11, UID: 100:101, "+
		"ON: 2020-10-29 23:08:00 +0000 UTC, "+
		"SENTON: 2020-10-29 23:08:00 +0000 UTC, "+
		"BCC: [1] CC: [0] FROM: [x] SUBJECT: [xyz] TO: [z] HEADER: \"X\" [y] "+
		"BODY: Body1, BODY: Body2, "+
		"TEXT: Text1, TEXT: Text2, "+
		"Flag: SEEN, Flag: ANSWERED, Flag: FLAGGED, Flag: DELETED, Flag: DRAFT, Flag: RECENT, KEYWORD: FLAGS, KEYWORD: +FLAGS, KEYWORD: -FLAGS, "+
		"UN: SEEN, UN: ANSWERED, UN: FLAGGED, UN: DELETED, UN: DRAFT, OLD, UNKEYWORD: FLAGS, UNKEYWORD: +FLAGS, UNKEYWORD: -FLAGS, "+
		"LARGER: 400, "+
		"SMALLER: 1, "+
		"}",
		ToString(sc))
}

func conditionToString1(t *testing.T) {
	sc := new(imap.SearchCriteria)
	sc.Since, _ = time.Parse(time.RFC3339, "2020-10-29T23:08:00Z")
	sc.Before, _ = time.Parse(time.RFC3339, "2020-10-27T23:08:00Z")
	sc.SentSince = sc.Since
	sc.SentBefore = sc.Before
	sc.Header = textproto.MIMEHeader{}
	sc.Header.Add("Subject", "xyz")
	sc.Header.Add("To", "z")
	sc.Header.Add("From", "x")
	sc.Header.Add("X", "y")
	sc.Header.Add("Cc", "0")
	sc.Header.Add("Bcc", "1")
	sc.Smaller = 1
	sc.Larger = 400
	sc.Body = []string{"Body3"}
	sc.Text = []string{"Text3"}
	sc.WithFlags = []string{imap.SeenFlag, imap.AnsweredFlag, imap.FlaggedFlag, imap.DeletedFlag, imap.DraftFlag, imap.RecentFlag, string(imap.SetFlags), imap.AddFlags, imap.RemoveFlags}
	sc.WithoutFlags = sc.WithFlags
	sc.Not = []*imap.SearchCriteria{
		0: {SeqNum: &imap.SeqSet{Set: []imap.Seq{{20, 21}}}},
	}
	sc.Or = [][2]*imap.SearchCriteria{0: {
		0: {Uid: &imap.SeqSet{Set: []imap.Seq{{30, 31}}}},
		1: {Uid: &imap.SeqSet{Set: []imap.Seq{{40, 41}}}},
	}}
	require.Equal(t, "SearchCriteria {"+
		"SINCE: 2020-10-29 23:08:00 +0000 UTC, "+
		"BEFORE: 2020-10-27 23:08:00 +0000 UTC, "+
		"SENTSINCE: 2020-10-29 23:08:00 +0000 UTC, "+
		"SENTBEFORE: 2020-10-27 23:08:00 +0000 UTC, "+
		"BCC: [1] CC: [0] FROM: [x] SUBJECT: [xyz] TO: [z] HEADER: \"X\" [y] "+
		"BODY: Body3, "+
		"TEXT: Text3, "+
		"Flag: SEEN, Flag: ANSWERED, Flag: FLAGGED, Flag: DELETED, Flag: DRAFT, Flag: RECENT, KEYWORD: FLAGS, KEYWORD: +FLAGS, KEYWORD: -FLAGS, "+
		"UN: SEEN, UN: ANSWERED, UN: FLAGGED, UN: DELETED, UN: DRAFT, OLD, UNKEYWORD: FLAGS, UNKEYWORD: +FLAGS, UNKEYWORD: -FLAGS, "+
		"LARGER: 400, "+
		"SMALLER: 1, "+
		"NOT[]{ SearchCriteria {SeqNum: 20:21, }}NOT[], "+
		"OR[]{  "+
		"OR{ SearchCriteria {UID: 30:31, }, SearchCriteria {UID: 40:41, }}OR, "+
		"}OR[], "+
		"}",
		ToString(sc))
}
