package conditions

import "fmt"

type and struct {
	conditions *[]Condition
	locked *bool
	parent *Condition
}

func(a *and) init(){
	a.conditions = emptyConditions()
	a.locked = conditionUnLocked
}

func (a and) ParseYaml(item interface{}) {
	parseYaml(item, a)
}

func (a and) SetCursor() {
	if a.parent != nil {
		(*a.parent).SetCursor()
	}else{
		*a.conditions = *emptyConditions()
		var x Condition = a
		a.Add(pair{"cursor", nil, &x})
		a.locked = conditionLocked
	}
}
func (a and) Add(c Condition) {
	if *a.locked {
		return
	}
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
		builder := ""
		for _, c := range *a.conditions {
			if builder != "" {
				builder += " and "
			}
			builder += " " + c.String() + " "
		}
		str = fmt.Sprintf("( %s )", builder)
	}
	return
}
