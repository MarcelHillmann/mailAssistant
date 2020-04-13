package actions

import (
	"mailAssistant/logging"
	"sync/atomic"
)

const (
	// Locked is the state of a locked job
	Locked = int32(1)
	// Released is the state of a released job
	Released = int32(0)
)


func isLockedElseLock(logger logging.Logger, waitGroup *int32) bool {
	if atomic.LoadInt32(waitGroup) > Released {
		logger.Info("is locked")
		return true
	}
	atomic.StoreInt32(waitGroup, Locked)
	logger.Info("lock")
	return false
}

func unlockAlways(logger logging.Logger, waitGroup *int32) {
	if err := recover(); err != nil {
		logger.Severe(err)
	}
	atomic.StoreInt32(waitGroup, Released)
	logger.Info("unlocked")
}

