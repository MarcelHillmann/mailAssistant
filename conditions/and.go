package conditions

import "fmt"

type and struct {
	conditions *[]Condition
}

func (a and) ParseYaml(item interface{}) {
	parseYaml(item, a)
}

func (a and) Add(c Condition) {
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
