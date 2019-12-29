package cntl

import (
	"github.com/onatm/clockwerk"
	"sync"
)

var (
	clocks = make([]*clockwerk.Clockwerk,0)
	mutex = &sync.Mutex{}
)

// NewClockwork is a factory for the clockwerk framework, it's saving each instance
func NewClockwork() *clockwerk.Clockwerk {
	res := clockwerk.New()
	mutex.Lock()
	clocks = append(clocks, res)
	mutex.Unlock()
	return res
}

// StopAllClocks is stopping all saved clockwerk framework instances
// it is blocking new factory calls
func StopAllClocks(){
	mutex.Lock()
	for _, clock := range clocks {
		clock.Stop()
	}
	clocks = make([]*clockwerk.Clockwerk,0)
}