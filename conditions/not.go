package conditions

type not struct {
	conditions *[]condition
}

func (n not) ParseYaml(item interface{}){
	parseYaml(item, n)
}
func (n not) Add(c condition) {
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