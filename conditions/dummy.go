package conditions

import "fmt"

type dummy struct {

}

func (d dummy) Parent(*Condition) {
	panic(fmt.Errorf("not yet implemented"))
}

func (d dummy) SetCursor(){
	panic(fmt.Errorf("not yet implemented"))
}

func (d dummy) Add(Condition) {
	panic(fmt.Errorf("not yet implemented"))
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