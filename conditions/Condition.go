package conditions

import (
	"fmt"
	"github.com/emersion/go-imap"
	"log"
	"mailAssistant/planning"
	"strings"
	"time"
)

var (
	conditionLocked   *bool
	conditionUnLocked *bool
)

func init() {
	locked, unlocked := true, false
	conditionLocked = &locked
	conditionUnLocked = &unlocked
}

// Condition represents a parsed YAML stream, for search on the IMAP server
type Condition interface {
	Add(Condition)
	Get() []interface{}
	String() string
	ParseYaml(interface{})
	SetCursor()
}

// ParseYaml is reading a yaml stream and convert it to a Condition
func ParseYaml(item interface{}) Condition {
	cond := and{ parent: nil}
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
	_, mapString := item.(map[string]interface{})
	_, mapInterface := item.(map[interface{}]interface{})

	if mapString || mapInterface {
		v2 := item.(map[interface{}]interface{})
		notAllowedKey(v2)
		field := strings.ToLower(v2["field"].(string))
		switch field {
		case "cursor":
			condition.SetCursor()
			return
		case "or":
			nCondition := or{parent: &condition}
			nCondition.init()
			nCondition.ParseYaml(v2["value"])
			condition.Add(nCondition)
		case "and":
			nCondition := and{parent: &condition}
			nCondition.init()
			nCondition.ParseYaml(v2["value"])
			condition.Add(nCondition)
		case "not":
			nCondition := not{parent: &condition}
			nCondition.init()
			nCondition.ParseYaml(v2["value"])
			condition.Add(nCondition)
		case "older":
			nCondition := pair{field: "since", parent: &condition}
			duration := planning.ParseSchedule(v2["value"].(string))
			durSec := time.Now().Unix() - int64(duration.Seconds())
			nCondition.value = time.Unix(durSec, 0).Format(imap.DateLayout)
			condition.Add(nCondition)
		case "younger":
			nCondition := pair{field: "before", parent: &condition}
			duration := planning.ParseSchedule(v2["value"].(string))
			durSec := time.Now().Unix() - int64(duration.Seconds())
			nCondition.value = time.Unix(durSec, 0).Format(imap.DateLayout)
			condition.Add(nCondition)
		default:
			if valid(field) {
				condition.Add(pair{field: field, value: v2["value"], parent: &condition})
			} else {
				log.Panicf("invalid Field: %s", field)
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

func valid(field string) bool {
	for _, k := range keywords {
		if k == field {
			return true
		}
	}
	return false
}

func notAllowedKey(m map[interface{}]interface{}) {
	for key := range m {
		if key == "field" || key == "value" {
			continue
		} else {
			panic(fmt.Errorf("invalid key '%s'", key))
		}
	}
}
