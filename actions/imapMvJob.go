package actions

import (
	"mailAssistant/account"
	"mailAssistant/errors"
	"mailAssistant/logging"
	"math"
	"time"
)

const moveTo = "moveTo"

func init() {
	register("imap_mv", newImapMove)
}

func newImapMove(job Job, waitGroup *int32) {
	logger := logging.NewLogger()
	if isLockedElseLock(logger, waitGroup) {
		return
	}
	defer unlockAlways(logger, waitGroup)

	job.GetAccount(job.GetString("mail_account")).
		DialAndLoginPromise(func(promise *account.ImapPromise) {
			promise.SelectPromise(job.GetString("path"), false, func(promise *account.ImapPromise) {
				promise.FetchPromise(job.getSearchParameter(), false, func(promise *account.MsgPromises) {
					if num, err := promise.Move(job.GetString(moveTo)); errors.IsEmpty(err) {
						logger.Debug(moveTo, job.GetString(moveTo), "nothing to do")
					} else if err == nil {
						logger.Debug(moveTo, job.GetString(moveTo), "successfully", num)
						mod := time.Duration(math.RoundToEven(float64(num / 1000)))
						time.Sleep(mod * time.Second)
					} else {
						panic(err)
					}
				}) // FetchPromise
			}) // SelectPromise
		}) // DialAndLoginPromise
	// GetAccount
}
