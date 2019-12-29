package arguments

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// NewEmptyArgs is a factory for a empty Args object
func NewEmptyArgs() *Args {
	return NewArgs(make(map[string]interface{}))
}

// NewArgs is a factory for a Args object
func NewArgs(args map[string]interface{}) *Args {
	return &Args{args}
}

// Args represents all arguments for jobs and rules
type Args struct {
	args map[string]interface{}
}

// GetArg returns a searched argument as interface
func (a Args) GetArg(key string) interface{} {
	if ret, ok := a.args[key]; ok {
		return ret
	}
	return nil
}

// GetArgKeys returns all argument keys
func (a Args) GetArgKeys() []string {
	result := make([]string, 0)
	for key := range a.args {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

// String shows all arguments as string
func (a Args) String() string {
	return fmt.Sprint(a.args)
}

// GetBool search for a argument and return it as bool
func (a Args) GetBool(key string) bool {
	value := a.GetArg(key)
	switch value.(type) {
	case bool:
		return value.(bool)
	case string:
		if value == "1" || strings.ToLower(value.(string)) == "true" {
			return true
		}
		return false
	default:
		return false
	}
}

// GetList search for a argument and return it as interface array
func (a Args) GetList(key string) []interface{} {
	value := a.GetArg(key)
	kind := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	if kind != nil && //
		(kind.Kind() == reflect.Array || kind.Kind() == reflect.Slice) {
		res := make([]interface{}, 0)
		for i := 0; i < v.Len(); i++ {
			res = append(res, v.Index(i).Interface())
		}
		return res
	}
	return []interface{}{}
}

// GetMap search for a argument and return it as map
func (a Args) GetMap(key string) map[string]interface{} {
	value := a.GetArg(key)
	kind := reflect.TypeOf(value)
	if kind != nil && kind.Kind() == reflect.Map {
		return value.(map[string]interface{})
	}
	return make(map[string]interface{})
}

// GetInt search for a argument and return it as int
func (a Args) GetInt(key string) int {
	value := a.GetArg(key)
	switch value.(type) {
	case int:
		return value.(int)
	case string:
		if res, err := strconv.Atoi(value.(string)); err == nil {
			return res
		}
		return 0
	default:
		return 0
	}
}

// HasArg is searching if a argument exists
func (a Args) HasArg(key string) bool {
	_, found := a.args[key]
	return found
}

// GetString search for a argument and return it as string
func (a Args) GetString(key string) string {
	value := a.GetArg(key)
	switch value.(type) {
	case string:
		return value.(string)
	case int:
		return strconv.Itoa(value.(int))
	case bool:
		if value.(bool) {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

// GetArgs is returning all arguments as map
func (a Args) GetArgs() map[string]interface{} {
	return a.args
}

// SetArg save a key value pair
func (a Args) SetArg(key string, value interface{}) {
	a.args[key] = value
}
