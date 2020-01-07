package conditions

import "fmt"

func newAnd() (res and) {
	res = and{parent: &headParent{}}
	res.init()
	return
}

type and struct {
	parent     *headParent
	conditions *[]Condition
	locked     *bool
}

func (a *and) init() {
	a.conditions = emptyConditions()
	a.locked = conditionUnLocked()
}

func (a and) ParseYaml(item interface{}) {
	parseYaml(item, a)
}

func (a and) SetCursor() {
	if a.parent != nil && a.parent.HasParent() {
		a.parent.SetCursor()
	} else {
		*a.conditions = *emptyConditions()
		a.Add(newCursor())
		*a.locked = true
	}
}

func (a and) Parent(c Condition) {
	a.parent.Parent(c)
}

func (a and) Add(c Condition) {
	if *a.locked {
		return
	}

	c.Parent(a)
	*a.conditions = append(*a.conditions, c)
}

func (a and) Get() (res []interface{}) {
	res = make([]interface{}, 0)
	for _, c := range *a.conditions {
		res = append(res, c.Get()...)
	}
	return
}

func (a and) String() (str string) {
	len := len(*a.conditions)
	switch len {
	case 0:
		str = dummy{}.String()
	case 1:
		str = (*a.conditions)[0].String()
	default:
		builder := stringEmpty
		for _, c := range *a.conditions {
			if builder != stringEmpty {
				builder += stringAnd
			}
			builder += c.String()
		}
		str = fmt.Sprintf(stringFormat, builder)
	}
	return
}
