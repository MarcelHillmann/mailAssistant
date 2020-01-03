package conditions

import (
	"fmt"
	"github.com/emersion/go-imap"
	"log"
	"mailAssistant/planning"
	"strings"
	"time"
)

type Condition interface {
	Add(Condition)
	Get() []interface{}
	String() string
	ParseYaml(interface{})
}

// ParseYaml is reading a yaml stream and convert it to a Condition
func ParseYaml(item interface{}) Condition {
	cond := and{emptyConditions()}
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
	if v2, ok := item.(map[string]interface{}); ok {
		notAllowedKey(v2)
		field := strings.ToLower(v2["field"].(string))
		switch field {
		case "or":
			nCondition := or{emptyConditions()}
			nCondition.ParseYaml(v2["value"])
			condition.Add(nCondition)
		case "and":
			nCondition := and{emptyConditions()}
			nCondition.ParseYaml(v2["value"])
			condition.Add(nCondition)
		case "not":
			nCondition := not{emptyConditions()}
			nCondition.ParseYaml(v2["value"])
			condition.Add(nCondition)
		case "older":
			nCondition := pair{}
			nCondition.field = "since"
			duration := planning.ParseSchedule(v2["value"].(string))
			durSec := time.Now().Unix() - int64(duration.Seconds())
			nCondition.value = time.Unix(durSec, 0).Format(imap.DateLayout)
			condition.Add(nCondition)
		case "younger":
			nCondition := pair{}
			nCondition.field = "before"
			duration := planning.ParseSchedule(v2["value"].(string))
			durSec := time.Now().Unix() - int64(duration.Seconds())
			nCondition.value = time.Unix(durSec, 0).Format(imap.DateLayout)
			condition.Add(nCondition)
		default:
			if valid(field) {
				condition.Add(pair{field, v2["value"]})
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

func notAllowedKey(m map[string]interface{}) {
	for key := range m {
		if key == "field" || key == "value" {
			continue
		} else {
			panic(fmt.Errorf("invalid key '%s'", key))
		}
	}
}
