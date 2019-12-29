package logging

import "strings"

type registryAux struct {
	Name     string         `yaml:"name"`
	Level    string         `yaml:"level"`
	Children []*registryAux `yaml:"children"`
}

func (aux registryAux) HasChildren() bool {
	return len(aux.Children) > 0
}

func (aux registryAux) HasNoLevel() bool {
	return aux.GetLevel() == notExists
}

func (aux registryAux) GetLevel() logLevel {
	return stringToLogLevel(aux.Level)
}

func stringToLogLevel(level string) logLevel {
	level = strings.ToUpper(level)
	switch {
	case level == "":
		return notExists
	case level == "SEVERE":
		return severe
	case level == "WARN":
		return warn
	case level == "INFO":
		return info
	case level == "DEBUG":
		return debug
	case level == "ALL":
		return all
	default:
		return none
	}
}
