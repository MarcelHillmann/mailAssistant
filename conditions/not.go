package conditions

type not struct {
	conditions *[]Condition
	locked *bool
	parent *Condition
}

func(n *not) init(){
	n.conditions = emptyConditions()
	n.locked = conditionUnLocked
}

func (n not) SetCursor(){
	if n.parent != nil {
		(*n.parent).SetCursor()
	}else{
		n.init()
		var x Condition = n
		n.Add(pair{"cursor", nil, &x})
		n.locked = conditionLocked
	}
}

func (n not) ParseYaml(item interface{}){
	parseYaml(item, n)
}

func (n not) Add(c Condition) {
	if *n.locked {
		return
	}
	*n.conditions = append(*n.conditions, c)
}

func (n not) Get() (res []interface{}) {
	res = make([]interface{},1)
	res[0] = "not"
	for _, c := range *n.conditions {
		res = append(res, c.Get()...)
	}
	return
}

func(n not)String() string {
	res:="not "
	for _, c := range *n.conditions {
		res += c.String()+" "
	}
	return res
}