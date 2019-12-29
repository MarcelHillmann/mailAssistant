package actions

import (
	"mailAssistant/account"
	e "mailAssistant/errors"
	"mailAssistant/logging"
)

func init() {
	register("junk", newJunkJob)
}

func newJunkJob(job Job, waitGroup *int32) {
	logger := logging.NewLogger()
	if isLockedElseLock(logger, waitGroup) {
		return
	}
	defer unlockAlways(logger, waitGroup)

	job.GetAccount(job.GetString("mail_account")).
		DialAndLoginPromise(func(promise *account.ImapPromise) {
			promise.SelectPromise(job.GetString("path"), false, func(promise *account.ImapPromise) {
				promise.SearchPromise(job.getSearchParameter(), false, func(promise *account.MsgPromises) {
					if deleted, err := promise.Delete(); err == nil {
						logger.Info("successfully deleted", deleted)
					} else if e.IsEmpty(err) {
						logger.Info("nothing to do")
					} else {
						panic(err)
					}
				}) // SearchPromise
			}) // SelectPromise
		}) // DialAndLoginPromise
	// GetAccount
}
