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

func ToString(c *imap.SearchCriteria) string {
	buffer := bytes.NewBufferString("SearchCriteria {")

	if c.SeqNum != nil {
		_, _ = fmt.Fprintf(buffer, "SeqNum: %#v, ", *c.SeqNum)
	}
	if c.Uid != nil {
		_, _ = fmt.Fprintf(buffer, "UID: %#v, ", *c.Uid)
	}

	if !c.Since.IsZero() && !c.Before.IsZero() && c.Before.Sub(c.Since) == 24*time.Hour {
		_, _ = fmt.Fprintf(buffer, "ON: %s, ", c.Since.String())
	} else {
		if !c.Since.IsZero() {
			_, _ = fmt.Fprintf(buffer, "SINCE: %s, ", c.Since.String())
		}
		if !c.Before.IsZero() {
			_, _ = fmt.Fprintf(buffer, "BEFORE: %s, ", c.Before.String())
		}
	}
	if !c.SentSince.IsZero() && !c.SentBefore.IsZero() && c.SentBefore.Sub(c.SentSince) == 24*time.Hour {
		_, _ = fmt.Fprintf(buffer, "SENTON: %s, ", c.SentSince.String())
	} else {
		if !c.SentSince.IsZero() {
			_, _ = fmt.Fprintf(buffer, "SENTSINCE: %s, ", c.SentSince.String())
		}
		if !c.SentBefore.IsZero() {
			_, _ = fmt.Fprintf(buffer, "SENTBEFORE: %s, ", c.SentBefore.String())
		}
	}

	for key, values := range c.Header {
		switch key {
		case "Bcc", "Cc", "From", "Subject", "To":
			_, _ = fmt.Fprintf(buffer, "%s: ", strings.ToUpper(key))
		default:
			_, _ = fmt.Fprintf(buffer, "HEADER: %#v", key)
		}
		_, _ = fmt.Fprintf(buffer, "%s", values)
	}

	for _, value := range c.Body {
		_, _ = fmt.Fprintf(buffer, "BODY: %s, ", value)
	}
	for _, value := range c.Text {
		_, _ = fmt.Fprintf(buffer, "TEXT: %s, ", value)
	}

	for _, flag := range c.WithFlags {
		switch flag {
		case imap.AnsweredFlag, imap.DeletedFlag, imap.DraftFlag, imap.FlaggedFlag, imap.RecentFlag, imap.SeenFlag:
			_, _ = fmt.Fprintf(buffer, "Flag: %s, ", strings.ToUpper(strings.TrimPrefix(flag, "\\")))
		default:
			_, _ = fmt.Fprintf(buffer, "KEYWORD: %s, ", flag)
		}
	}
	for _, flag := range c.WithoutFlags {
		switch flag {
		case imap.AnsweredFlag, imap.DeletedFlag, imap.DraftFlag, imap.FlaggedFlag, imap.SeenFlag:
			_, _ = fmt.Fprintf(buffer, "UN: %s, ", strings.ToUpper(strings.TrimPrefix(flag, "\\")))
		case imap.RecentFlag:
			_, _ = fmt.Fprintf(buffer, "OLD, ")
		default:
			_, _ = fmt.Fprintf(buffer, "UNKEYWORD: %s, ", flag)
		}
	}

	if c.Larger > 0 {
		_, _ = fmt.Fprintf(buffer, "LARGER: %d, ", c.Larger)
	}
	if c.Smaller > 0 {
		_, _ = fmt.Fprintf(buffer, "SMALLER: %d, ", c.Smaller)
	}

	for _, not := range c.Not {
		_, _ = fmt.Fprintf(buffer, "NOT: %#v, ", not.Format())
	}

	for _, or := range c.Or {
		_, _ = fmt.Fprintf(buffer, "OR: %#v - %#v, ", or[0].Format(), or[1].Format())
	}

	buffer.WriteString("}")
	return buffer.String()
}
