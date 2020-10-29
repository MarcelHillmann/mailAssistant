package conditions

import "fmt"

type dummy struct {
}

func (d dummy) Parent(*Condition) {
	panic(fmt.Errorf(notYetImplemented))
}

func (d dummy) SetCursor() {
	panic(fmt.Errorf(notYetImplemented))

}

func (d dummy) Add(Condition) {
	panic(fmt.Errorf(notYetImplemented))
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
