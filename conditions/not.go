package conditions

func newNot() (res not) {
	res = not{parent: &headParent{}}
	res.init()
	return
}

type not struct {
	conditions *[]Condition
	locked     *bool
	parent     *headParent
}

func (n *not) init() {
	n.conditions = emptyConditions()
	n.locked = conditionUnLocked()
}

func (n not) Parent(c Condition) {
	n.parent.Parent(c)
}

func (n not) SetCursor() {
	if n.parent != nil && n.parent.HasParent() {
		n.parent.SetCursor()
	} else {
		*n.conditions = *emptyConditions()
		n.Add(newCursor())
		*n.locked = true
	}
}

func (n not) ParseYaml(item interface{}) {
	parseYaml(item, n)
}

func (n not) Add(c Condition) {
	if *n.locked {
		return
	}

	c.Parent(n)
	*n.conditions = append(*n.conditions, c)
}

func (n not) Get() (res []interface{}) {
	res = make([]interface{}, 0)
	if false == *n.locked && len(*n.conditions) > 0 {
		res = append(res, "not")
	}
	for _, c := range *n.conditions {
		res = append(res, c.Get()...)
	}
	return
}

func (n not) String() (str string) {
	len := len(*n.conditions)
	switch len {
	case 0:
		str = dummy{}.String()
	case 1:
		str = stringNot + (*n.conditions)[0].String()
	default:
		str = stringEmpty
		for _, c := range *n.conditions {
			if str != stringEmpty {
				str += stringAnd
			}
			str += c.String()
		}
		str = stringNot + "( " + str + " )"
	}
	return
}
