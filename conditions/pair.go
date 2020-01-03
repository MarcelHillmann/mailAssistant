package conditions

import (
	"fmt"
	"strings"
)

type pair struct {
	field string
	value interface{}
}

func(p pair) ParseYaml(interface{}) {
	panic(fmt.Errorf("not yet implemented"))
}

func (p pair) Add(condition){
	panic(fmt.Errorf("not yet implemented"))
}

func (p pair) Get() []interface{} {
	if p.value == nil {
		return []interface{}{p.field}
	} else if v, ok:= p.value.(int); ok {
		return []interface{}{p.field, uint32(v)}
	}
	return []interface{}{p.field, p.value}
}

func (p pair) String() string {
	return fmt.Sprintf("%s=%v",strings.ToUpper(p.field),p.value)
}