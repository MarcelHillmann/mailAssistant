package logging

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var (
	levelTestCases = [][]interface{}{
		{"SeVeRe", severe},
		{"WaRn", warn},
		{"inFO", info},
		{"DeBug", debug},
		{"aLL", all},
		{"nOnE", none},
		{"xXx", none},
		{"", notExists},
	}
)

func TestRegistryAux(t *testing.T) {
	t.Run("has children", registryAuxHasChildren)
	t.Run("has no level", registryAuxHasNoLevel)
	t.Run("get level", registryAuxGetLevel)
}

func registryAuxHasChildren(t *testing.T) {
	children := make([]*registryAux, 0)
	aux1 := registryAux{"", "", children}
	require.False(t, aux1.HasChildren())

	children = append(children, &registryAux{"a", "ALL", make([]*registryAux, 0)})
	aux2 := registryAux{"", "", children}
	require.True(t, aux2.HasChildren())
}

func registryAuxHasNoLevel(t *testing.T) {
	children := make([]*registryAux, 0)
	aux := registryAux{"", "", children}
	require.True(t, aux.HasNoLevel())
	aux.Level = "none"
	require.False(t, aux.HasNoLevel())
}

func registryAuxGetLevel(t *testing.T) {
	for _, test := range levelTestCases {
		key := test[0].(string)
		t.Run(strings.ToUpper(key), func(t *testing.T) {
			aux := registryAux{"", key, nil}
			require.Equal(t, test[1].(logLevel), aux.GetLevel())
		}) // t.Run
	} // for each
}

func TestStringToLogLevel(t *testing.T) {
	for _, test := range levelTestCases {
		key := test[0].(string)
		t.Run(strings.ToUpper(key), func(t *testing.T) {
			require.Equal(t, test[1].(logLevel), stringToLogLevel(key))
		})
	}
}
