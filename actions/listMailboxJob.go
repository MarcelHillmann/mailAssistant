package actions

import "mailAssistant/account"

func init() {
	register("list", newListMailbox)
}

func newListMailbox(job Job, waitGroup *int32, result func(int)) {
	logger := job.GetLogger()
	if isLockedElseLock(logger, waitGroup) {
		return
	}
	defer unlockAlways(logger, waitGroup)

	job.GetAccount(job.GetString("mail_account")).
		DialAndLoginPromise(func(promise *account.ImapPromise) {
			promise.ListMailboxes()
		})
} // newListMailbox