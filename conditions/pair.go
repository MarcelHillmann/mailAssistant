package conditions

import (
	"fmt"
	"strings"
)

func newCursor() pair {
	return pair{parent: &headParent{}, keyval: &keyVal{field: CURSOR}}
}

func newPair(field string, value interface{}) pair {
	return pair{parent: &headParent{}, keyval: &keyVal{field: field,value: value}}
}

type keyVal struct {
	field string
	value interface{}
}

func (k *keyVal) Field(field string) {
	k.field = field
}

func (k *keyVal) Value(value interface{}) {
	k.value = value
}
type pair struct {
	parent *headParent
	keyval *keyVal
}

func (p pair) Parent(c Condition){
	p.parent.Parent(c)
}

func (p pair) SetCursor() {
	if p.parent != nil && p.parent.HasParent() {
		p.parent.SetCursor()
	} else {
		p.keyval.Field(CURSOR)
		p.keyval.Value(nil)
	}
}

func (p pair) ParseYaml(interface{}) {
	panic(fmt.Errorf("not yet implemented"))
}

func (p pair) Add(Condition) {
	panic(fmt.Errorf("not yet implemented"))
}

func (p pair) Get() []interface{} {
	field := strings.ToUpper(p.keyval.field)
	if p.keyval.value == nil {
		return []interface{}{field}
	} else if v, ok := p.keyval.value.(int); ok {
		return []interface{}{field, uint32(v)}
	}
	return []interface{}{field, p.keyval.value}
}

func (p pair) String() string {
	return fmt.Sprintf("%s='%v'", strings.ToUpper(p.keyval.field), p.keyval.value)
}
