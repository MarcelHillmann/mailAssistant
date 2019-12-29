package errors

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmptyError_Error(t *testing.T) {
	err := NewEmpty()
	require.EqualError(t, err, "is empty")
	require.Equal(t, "is empty", err.Error())
}

func TestIsEmpty(t *testing.T) {
	err := NewEmpty()
	require.True(t, IsEmpty(err))
	ptr := &emptyError{}
	require.True(t, IsEmpty(ptr))
	require.False(t, IsEmpty(errors.New("noop")))
}

func TestNewEmpty(t *testing.T) {
	err := NewEmpty()
	require.NotNil(t, err)
	require.Error(t, err, "is empty")

}
