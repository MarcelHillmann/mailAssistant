package actions

import (
	"mailAssistant/account"
	"mailAssistant/errors"
	"mailAssistant/logging"
	"math"
	"time"
)

func init() {
	register("seen", newSeenJob)
}

func newSeenJob(job Job, waitGroup *int32) {
	logger := logging.NewLogger()
	if isLockedElseLock(logger, waitGroup) {
		return
	}
	defer unlockAlways(logger, waitGroup)
	job.GetAccount(job.GetString("mail_account")).
		DialAndLoginPromise(func(promise *account.ImapPromise) {
			promise.SelectPromise(job.GetString("path"), false, func(promise *account.ImapPromise) {
				promise.FetchPromise(job.getSearchParameter(), false, func(promise *account.MsgPromises) {
					if num, err := promise.SetSeen(); errors.IsEmpty(err) {
						logger.Debug("nothing to do")
					} else if err == nil {
						logger.Debug("successfully", num)
						mod := time.Duration(math.RoundToEven(float64(num / 1000)))
						time.Sleep(mod * time.Second)
					} else {
						panic(err)
					}
				}) // FetchPromise
			}) // SelectPromise
		}) // DialAndLoginPromise
	// job.GetAccount
}
