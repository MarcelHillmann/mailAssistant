package actions

import (
	"mailAssistant/account"
	e "mailAssistant/errors"
	"mailAssistant/logging"
)

func init() {
	register("junk", newJunkJob)
}

func newJunkJob(job Job, waitGroup *int32, result func(int)) {
	logger := logging.NewLogger()
	if isLockedElseLock(logger, waitGroup) {
		return
	}
	defer unlockAlways(logger, waitGroup)

	job.GetAccount(job.GetString("mail_account")).
		DialAndLoginPromise(func(promise *account.ImapPromise) {
			promise.SelectPromise(job.GetString("path"), false, func(promise *account.ImapPromise) {
				promise.FetchPromise(job.getSearchParameter(), false, func(promise *account.MsgPromises) {
					if deleted, err := promise.Delete(); err == nil {
						result(deleted)
						logger.Info("successfully deleted", deleted)
					} else if e.IsEmpty(err) {
						logger.Info("nothing to do")
					} else {
						panic(err)
					}
				}) // FetchPromise
			}) // SelectPromise
		}) // DialAndLoginPromise
	// GetAccount
}
