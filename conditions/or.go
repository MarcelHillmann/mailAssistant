package conditions

import "fmt"

func newOr() (res or) {
	res = or{parent: &headParent{}}
	res.init()
	return
}

type or struct {
	parent     *headParent
	conditions *[]Condition
	locked     *bool
}

func (o *or) init() {
	o.conditions = emptyConditions()
	o.locked = conditionUnLocked()
}

func (o or) Parent(c Condition) {
	o.parent.Parent(c)
}

func (o or) SetCursor() {
	if o.parent != nil && o.parent.HasParent() {
		o.parent.SetCursor()
	} else {
		*o.conditions = *emptyConditions()
		o.Add(newCursor())
		*o.locked = true
	}
}

func (o or) ParseYaml(item interface{}) {
	parseYaml(item, o)
}

func (o or) Add(c Condition) {
	if *o.locked {
		return
	}

	c.Parent(o)
	*o.conditions = append(*o.conditions, c)
}

func (o or) Get() (res []interface{}) {
	last := len(*o.conditions) - 1
	res = make([]interface{}, 0)
	for i := range *o.conditions {
		if i < last || last == 0 && false == *o.locked {
			res = append(res, "or")
		}
		res = append(res, (*o.conditions)[i].Get()...)
	}
	return
}

func (o or) String() (str string) {
	len := len(*o.conditions)
	switch len {
	case 0:
		str = dummy{}.String()
	case 1:
		str = "or " + (*o.conditions)[0].String()
	default:
		builder := stringEmpty
		for _, c := range *o.conditions {
			if builder != stringEmpty {
				builder += stringOr
			}
			builder += c.String()
		}
		str = fmt.Sprintf(stringFormat, builder)
	}
	return
}
