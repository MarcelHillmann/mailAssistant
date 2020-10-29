package account

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAccountAux(t *testing.T) {
	t.Run("convert", accountAuxConvert)
	t.Run("isEmpty", accountAuxIsEmpty)
}

func accountAuxConvert(t *testing.T) {
	aux := accountAux{"test.yml", "test", "marcel", "geheim", "localhost", 1000, false, false}
	converted := aux.convert()

	require.Equal(t, "test", converted.name)
	require.Equal(t, "marcel", converted.username)
	require.Equal(t, "geheim", converted.password)
	require.Equal(t, "localhost", converted.hostname)
	require.Equal(t, 1000, converted.port)
	require.Equal(t, false, converted.debug)
	require.Equal(t, false, converted.skipVerify)
}

func accountAuxIsEmpty(t *testing.T) {
	aux := accountAux{"", "", "", "", "", 0, false, true}
	require.True(t, aux.IsEmpty())

	aux.fileName = "test"
	require.True(t, aux.IsEmpty())

	aux.Name = "test"
	require.True(t, aux.IsEmpty())

	aux.Hostname = "local"
	require.True(t, aux.IsEmpty())

	aux.Port = 1000
	require.False(t, aux.IsEmpty())
}
