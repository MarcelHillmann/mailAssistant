package conditions

import "fmt"

type dummy struct {

}

func (d dummy) Add(c Condition) {
	_ = c
}

func (d dummy) Get() []interface{} {
	return []interface{}{}
}

func (d dummy) String() string {
	return ""
}

func (d dummy) ParseYaml(interface{}) {
	panic(fmt.Errorf("never call this ;)"))
}