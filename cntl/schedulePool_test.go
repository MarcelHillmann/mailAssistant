package cntl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test0NewClockwork(t *testing.T) {
	require.NotNil(t, NewClockwork())
	require.Len(t, clocks, 1)
}

func Test1StopAllClocks(t *testing.T) {
	require.Len(t, clocks, 1)
	StopAllClocks()
	require.Len(t, clocks, 0)
}
