package conditions

import (
	"bytes"
	"fmt"
	"github.com/emersion/go-imap"
	"log"
	"mailAssistant/planning"
	"strings"
	"time"
)

const (
	stringEmpty       = ""
	stringOr          = " or "
	stringAnd         = " and "
	stringNot         = "not "
	stringFormat      = "( %s )"
	notYetImplemented = "not yet implemented"
	// CURSOR defines a not query request
	CURSOR = "CURSOR"
)

func conditionUnLocked() *bool {
	res := false
	return &res
}

// Condition represents a parsed YAML stream, for search on the IMAP server
type Condition interface {
	Add(Condition)
	Get() []interface{}
	String() string
	ParseYaml(interface{})
	Parent(c Condition)
	SetCursor()
}

// ParseYaml is reading a yaml stream and convert it to a Condition
func ParseYaml(item interface{}) Condition {
	cond := newAnd()
	cond.init()
	if item != nil {
		cond.ParseYaml(item)
	}
	return cond
}

func emptyConditions() *[]Condition {
	n := make([]Condition, 0)
	return &n
}

func parseYaml(item interface{}, condition Condition) {
	mapString, isMapString := item.(map[string]interface{})
	_, isMapInterface := item.(map[interface{}]interface{})

	if isMapString {
		mapIntf := make(map[interface{}]interface{})
		for key, val := range mapString {
			mapIntf[key] = val
		}
		item = mapIntf
		isMapInterface = true
	}

	if isMapInterface {
		v2 := item.(map[interface{}]interface{})
		allowedYamlKey(v2)
		field := strings.ToLower(v2["field"].(string))
		switch field {
		case "cursor":
			condition.SetCursor()
			return
		case "or":
			nCondition := newOr()
			nCondition.ParseYaml(v2["value"])
			condition.Add(nCondition)
		case "and":
			nCondition := newAnd()
			nCondition.ParseYaml(v2["value"])
			condition.Add(nCondition)
		case "not":
			nCondition := newNot()
			nCondition.ParseYaml(v2["value"])
			condition.Add(nCondition)
		case "older":
			duration := planning.ParseSchedule(v2["value"].(string))
			durSec := time.Now().Unix() - int64(duration.Seconds())
			value := time.Unix(durSec, 0).Format(imap.DateLayout)
			condition.Add(newPair("before", value))
		case "younger":
			duration := planning.ParseSchedule(v2["value"].(string))
			durSec := time.Now().Unix() - int64(duration.Seconds())
			value := time.Unix(durSec, 0).Format(imap.DateLayout)
			condition.Add(newPair("since", value))
		default:
			if validImapKeyword(field) {
				condition.Add(newPair(field, v2["value"]))
			} else {
				log.Panicf("invalid condition: %s", field)
			}
		}
	} else if v2, ok := item.([]interface{}); ok {
		for _, v3 := range v2 {
			parseYaml(v3, condition)
		}
	} else {
		log.Panicf("Value: %v", item)
	}
}

var keywords = []string{
	"all", "answered", "deleted", "draft", "flagged", "recent", "seen", "bcc", "cc", "from", "subject", "to",
	"before", "body", "header", "keyword", "larger", "new", "not", "old", "on", "or", "sentbefore", "senton", "sentsince",
	"since", "smaller", "text", "uid", "unanswered", "undeleted", "undraft", "unflagged", "unseen", "unkeyword",
}

func validImapKeyword(field string) bool {
	for _, k := range keywords {
		if k == field {
			return true
		}
	}
	return false
}

func allowedYamlKey(m map[interface{}]interface{}) {
	for key := range m {
		if key == "field" || key == "value" {
			continue
		} else {
			panic(fmt.Errorf("invalid yaml field: '%s'", key))
		}
	}
}

// MapToString converts a map to a string
func MapToString(m map[string]string) string {
	buf := bytes.NewBufferString("")
	for k, v := range m {
		buf.WriteString(k)
		buf.WriteString(": ")
		buf.WriteString(v)
		buf.WriteString(",")
	}
	return buf.String()
}

// ToString converts searchCriteria to string
func ToString(c *imap.SearchCriteria) string {
	s := searchToString{SearchCriteria: c}
	return s.String()
}

var p1d = 24 * time.Hour

type searchToString struct {
	*imap.SearchCriteria
	buffer *bytes.Buffer
}

func (s *searchToString) String() string {
	s.buffer = bytes.NewBufferString("SearchCriteria {")
	s.seqNum().uid().time().sent().header().body().text().flags().size().not().or()
	s.buffer.WriteString("}")
	return s.buffer.String()
}

func (s searchToString) seqNum() searchToString {
	if s.SearchCriteria.SeqNum != nil {
		s.addToBuffer("SeqNum: %#v, ", s.SearchCriteria.SeqNum)
	}
	return s
}

func (s searchToString) uid() searchToString {
	if s.SearchCriteria.Uid != nil {
		s.addToBuffer("UID: %#v, ", s.SearchCriteria.Uid)
	}
	return s
}

func (s searchToString) time() searchToString {
	c := s.SearchCriteria
	if !c.Since.IsZero() &&
		!c.Before.IsZero() &&
		c.Before.Sub(c.Since) == p1d {
		s.addToBuffer("ON: %s, ", c.Since.String())
	} else {
		if !c.Since.IsZero() {
			s.addToBuffer("SINCE: %s, ", c.Since.String())
		}
		if !c.Before.IsZero() {
			s.addToBuffer("BEFORE: %s, ", c.Before.String())
		}
	}
	return s
}

func (s searchToString) sent() searchToString {
	c := s.SearchCriteria
	if !c.SentSince.IsZero() &&
		!c.SentBefore.IsZero() &&
		c.SentBefore.Sub(c.SentSince) == p1d {
		s.addToBuffer("SENTON: %s, ", c.SentSince.String())
	} else {
		if !c.SentSince.IsZero() {
			s.addToBuffer("SENTSINCE: %s, ", c.SentSince.String())
		}
		if !c.SentBefore.IsZero() {
			s.addToBuffer("SENTBEFORE: %s, ", c.SentBefore.String())
		}
	}
	return s
}

func (s searchToString) header() searchToString {
	for key, values := range s.SearchCriteria.Header {
		switch key {
		case "Bcc", "Cc", "From", "Subject", "To":
			s.addToBuffer("%s: ", strings.ToUpper(key))
		default:
			s.addToBuffer("HEADER: %#v", key)
		}
		s.addToBuffer("%s", values)
	}
	return s
}

func (s searchToString) size() searchToString {
	if s.SearchCriteria.Larger > 0 {
		s.addToBuffer("LARGER: %d, ", s.SearchCriteria.Larger)
	}
	if s.SearchCriteria.Smaller > 0 {
		s.addToBuffer("SMALLER: %d, ", s.SearchCriteria.Smaller)
	}
	return s
}

func (s searchToString) body() searchToString {
	for _, value := range s.SearchCriteria.Body {
		s.addToBuffer("BODY: %s, ", value)
	}
	return s
}

func (s searchToString) text() searchToString {
	for _, value := range s.SearchCriteria.Text {
		s.addToBuffer("TEXT: %s, ", value)
	}
	return s
}

func (s searchToString) flags() searchToString {
	for _, flag := range s.SearchCriteria.WithFlags {
		switch flag {
		case imap.AnsweredFlag, imap.DeletedFlag, imap.DraftFlag, imap.FlaggedFlag, imap.RecentFlag, imap.SeenFlag:
			s.addToBuffer("Flag: %s, ", strings.ToUpper(strings.TrimPrefix(flag, "\\")))
		default:
			s.addToBuffer("KEYWORD: %s, ", flag)
		}
	}
	for _, flag := range s.SearchCriteria.WithoutFlags {
		switch flag {
		case imap.AnsweredFlag, imap.DeletedFlag, imap.DraftFlag, imap.FlaggedFlag, imap.SeenFlag:
			s.addToBuffer("UN: %s, ", strings.ToUpper(strings.TrimPrefix(flag, "\\")))
		case imap.RecentFlag:
			s.addToBuffer("OLD, ")
		default:
			s.addToBuffer("UNKEYWORD: %s, ", flag)
		}
	}
	return s
}

func (s searchToString) not() searchToString {
	for _, not := range s.SearchCriteria.Not {
		s.addToBuffer("NOT: %#v, ", not.Format())
	}
	return s
}

func (s searchToString) or() searchToString {
	for _, or := range s.SearchCriteria.Or {
		s.addToBuffer("OR: %#v - %#v, ", or[0].Format(), or[1].Format())
	}
	return s
}

func (s searchToString) addToBuffer(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(s.buffer, format, a...)
}
