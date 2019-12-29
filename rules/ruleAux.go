package rules

import "mailAssistant/arguments"

type ruleAux struct {
	fileName string
	Name     string                   `yaml:"name"`
	Schedule string                   `yaml:"schedule"`
	Action   string                   `yaml:"action"`
	Disabled bool                     `yaml:"disabled"`
	Args     []map[string]interface{} `yaml:"args"`
}

func (r ruleAux) convert() Rule {
	result := Rule{arguments.NewEmptyArgs(), r.Name, r.Schedule, r.Action, nil, r.Disabled}
	for arg := range r.Args {
		for key, value := range r.Args[arg] {
			result.SetArg(key, value)
		}
	}
	return result
}

func (r ruleAux) IsEmpty() bool {
	return r.Name == "" || r.Schedule == "" || r.Action == ""
}
