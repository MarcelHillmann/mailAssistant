package arguments

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var expected = map[string]interface{}{"a": "1", "b": false,"c": 4,"d":"","e":true,"f":uint(0)}

func TestArgs(t *testing.T) {
	t.Run("HasArg", argsHasArg)
	t.Run("GetArg", argsGetArg)
	t.Run("GetArgKeys", argsGetArgKey)
	t.Run("GetArgs", argsGetArgs)
	t.Run("GetBool", argsGetBool)
	t.Run("GetInt", argsGetInt)
	t.Run("GetList", argsGetList)
	t.Run("GetMap", argsGetMap)
	t.Run("GetString", argsGetString)
	t.Run("SetArg", argsSetArg)
	t.Run("String", argsString)
	t.Run("newArgs", newArgs)
	t.Run("newEmptyArgs", newEmptyArgs)
}

func argsHasArg(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		args := Args{map[string]interface{}{"test": true}}
		require.True(t, args.HasArg("test"))
	})
	t.Run("not found", func(t *testing.T) {
		args := Args{map[string]interface{}{"test": true}}
		require.False(t, args.HasArg("test2"))
	})
}

func argsGetArg(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		args := Args{map[string]interface{}{"test": true}}
		require.NotNil(t, args.GetArg("test"))
	})
	t.Run("not found", func(t *testing.T) {
		args := Args{map[string]interface{}{"test": true}}
		require.Nil(t, args.GetArg("test2"))
	})
}

func argsGetArgKey(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		args := Args{_expected()}
		keys := args.GetArgKeys()
		require.NotNil(t, keys)
		require.Len(t, keys, 6)
		require.Equal(t, []string{"a", "b", "c", "d","e","f"}, keys)
	})
	t.Run("empty", func(t *testing.T) {
		args := Args{map[string]interface{}{}}
		keys := args.GetArgKeys()
		require.NotNil(t, keys)
		require.Len(t, keys, 0)
		require.Equal(t, []string{}, keys)
	})
}

func argsGetArgs(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		args := Args{_expected()}
		actual := args.GetArgs()
		require.NotNil(t, actual)
		require.Equal(t, expected, actual)
		require.Len(t, actual,6)
	})
	t.Run("nilIsNil", func(t *testing.T) {
		args := Args{nil}
		require.Nil(t, args.GetArgs())
	})
}

func argsGetBool(t *testing.T) {
	args := Args{_expected()}
	t.Run("string_1", func(t *testing.T) {
		actual := args.GetBool("a")
		require.NotNil(t, actual)
		require.True(t,  actual)
	})
	t.Run("bool", func(t *testing.T) {
		actual := args.GetBool("b")
		require.NotNil(t, actual)
		require.False(t,  actual)
	})
	t.Run("int_4", func(t *testing.T) {
		actual := args.GetBool("c")
		require.NotNil(t, actual)
		require.False(t,  actual)
	})
	t.Run("string_blank", func(t *testing.T) {
		actual := args.GetBool("d")
		require.NotNil(t, actual)
		require.False(t, actual)
	})
}

func argsGetInt(t *testing.T) {
	args := Args{_expected()}

	t.Run("string_1", func(t *testing.T) {
		actual0 := args.GetInt("a")
		require.Equal(t, 1, actual0)
	})
	t.Run("bool_false", func(t *testing.T) {
		actual1 := args.GetInt("b")
		require.Equal(t, 0, actual1)
	})
	t.Run("int_4", func(t *testing.T) {
		actual2 := args.GetInt("c")
		require.Equal(t, 4, actual2)
	})
	t.Run("string_blank", func(t *testing.T) {
		actual3 := args.GetInt("d")
		require.Equal(t, 0, actual3)
	})
}

func argsGetList(t *testing.T) {
	args := Args{map[string]interface{}{"slice": []string{"a","b"},"array": [2]int{1,2},"else": 1}}
	t.Run("array", func(t *testing.T) {
		a :=args.GetList("array")
		require.NotNil(t,a)
		require.Len(t,a,2)
		require.Equal(t,[]interface{}{1,2},a)
	})
	t.Run("slice", func(t *testing.T) {
		a:=args.GetList("slice")
		require.NotNil(t,a)
		require.Len(t,a,2)
		require.Equal(t,[]interface{}{"a","b"},a)
	})
	t.Run("else", func(t *testing.T) {
		a:=args.GetList("else")
		require.NotNil(t,a)
		require.Len(t,a,0)
		require.Equal(t,[]interface{}{},a)
	})
	t.Run("???", func(t *testing.T) {
		a:=args.GetList("???")
		require.NotNil(t,a)
		require.Len(t,a,0)
		require.Equal(t,[]interface{}{},a)
	})
}

func argsGetMap(t *testing.T) {
	args := Args{map[string]interface{}{"test": nil,"map": map[string]interface{}{"a":"a","b":"b"}}}
	t.Run("map", func(t *testing.T) {
		a := args.GetMap("map")
		require.NotNil(t, a)
		require.Len(t, a, 2)
		require.Equal(t, map[string]interface{}{"a":"a","b":"b"},a)
	})
	t.Run("???", func(t *testing.T) {
		a := args.GetMap("test")
		require.NotNil(t, a)
		require.Len(t, a, 0)
		require.Equal(t, map[string]interface{}{},a)
	})
}

func argsGetString(t *testing.T) {
	args := Args{_expected()}
	actual0 := args.GetString("a")
	require.Equal(t, "1", actual0)
	actual1 := args.GetString("b")
	require.Equal(t, "false", actual1)
	actual2 := args.GetString("c")
	require.Equal(t, "4", actual2)
	actual3 := args.GetString("d")
	require.Equal(t, "", actual3)
	actual4 := args.GetString("e")
	require.Equal(t, "true", actual4)
	actual5 := args.GetString("f")
	require.Equal(t, "", actual5)
}

func argsSetArg(t *testing.T) {
	args := Args{_expected()}
	args.SetArg("foo", 1000)
	args.SetArg("bar", true)

	require.Len(t, args.args,8)
	require.Equal(t,1000, args.GetInt("foo"))
	require.True(t, args.GetBool("bar"))
}

func argsString(t *testing.T) {
	args := Args{_expected()}
	require.Equal(t, "map[a:1 b:false c:4 d: e:true f:0]", args.String())
}

func newArgs(t *testing.T) {
	arg := NewArgs(nil)
	require.Nil(t, arg.args)
}

func newEmptyArgs(t *testing.T) {
	arg := NewEmptyArgs()
	require.NotNil(t, arg.args)
	require.Len(t, arg.args,0)
}

func _expected()(res map[string]interface{}){
	res = make(map[string]interface{})
	for k, v:=range expected {
		res[k] = v
	}
	return
}