package conditions

type or struct {
	conditions *[]Condition
}

func (o or) ParseYaml(item interface{}){
	parseYaml(item, o)
}

func (o or) Add(c Condition) {
	*o.conditions = append(*o.conditions, c)
}

func (o or) Get() (res []interface{}) {
	last := len(*o.conditions) - 1
	res = make([]interface{}, 0)
	for i := range *o.conditions {
		if i < last || last == 0 {
			res = append(res, "or")
		}
		res = append(res, (*o.conditions)[i].Get()...)
	}
	return
}

func(o or)String() string {
	res:="("
	for _, c := range *o.conditions {
		if res != "(" {
			res +=" or "
		}
		res += " "+c.String()+" "
	}
	return res +")"
}