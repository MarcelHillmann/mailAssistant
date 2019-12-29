package rules

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRuleAux(t *testing.T){
	t.Run("convert",ruleAuxConvert)
	t.Run("IsEmpty", ruleAuxIsEmpty)
}

func ruleAuxConvert(t *testing.T){
	aux := ruleAux{"t.yml","t","1s","dummy", false,[]map[string]interface{}{map[string]interface{}{"test":"a", "bool":true,"int":0}}}
	converted := aux.convert()
	require.NotNil(t, converted)
	require.Equal(t, "t", converted.name)
	require.Equal(t, "1s", converted.schedule)
	require.Equal(t, "dummy", converted.action)
	require.Equal(t, false, converted.disabled)
	require.Equal(t, "a", converted.GetString("test"))
	require.Equal(t, true, converted.GetBool("bool"))
	require.Equal(t, 0, converted.GetInt("int"))
}

func ruleAuxIsEmpty(t *testing.T){
	aux := ruleAux{"","","","", false,make([]map[string]interface{},0)}
	require.True(t, aux.IsEmpty())
	aux.Name ="t"
	require.True(t, aux.IsEmpty())
	aux.Schedule ="t"
	require.True(t, aux.IsEmpty())
	aux.Action="t"
	require.False(t, aux.IsEmpty())
	aux.Name =""
	require.True(t, aux.IsEmpty())
	aux.Name ="t"
	aux.Schedule =""
	require.True(t, aux.IsEmpty())
}