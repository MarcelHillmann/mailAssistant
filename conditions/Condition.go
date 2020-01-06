package conditions

import (
	"fmt"
	"github.com/emersion/go-imap"
	"log"
	"mailAssistant/planning"
	"strings"
	"time"
)

const (
	stringEmpty = ""
	stringOr = " or "
	stringAnd = " and "
	stringNot = "not "
	stringFormat = "( %s )"
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
