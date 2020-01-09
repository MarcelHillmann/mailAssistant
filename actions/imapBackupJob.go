package actions

import (
	"mailAssistant/account"
	"mailAssistant/logging"
	"math"
	"time"
)

func init() {
	register("imap_backup", newImapBackup)
}

func newImapBackup(job Job, waitGroup *int32, result func(int)) {
	logger := logging.NewLogger()
	if isLockedElseLock(logger, waitGroup) {
		return
	}
	defer unlockAlways(logger, waitGroup)

	source := job.GetAccount(job.GetString("mail_account"))
	target := job.GetAccount(job.GetString("target_account"))
	pathFrom, pathTo := job.GetString("path"), job.GetString("saveTo")

	if source == nil || target == nil {
		return
	}
	source.DialAndLoginPromise(func(sourcePromise *account.ImapPromise) {
		target.DialAndLoginPromise(func(targetPromise *account.ImapPromise) {
			sourcePromise.SelectPromise(pathFrom, false, func(promise *account.ImapPromise) {
				promise.FetchPromise(job.getSearchParameter(), true, func(promise *account.MsgPromises) {
					targetPromise.UploadAndDelete(pathTo, promise, func(num int) {
						result(num)
						if num == 0 {
							logger.Debug("nothing to moveTo", moveTo)
						} else {
							promise.Expunge()
							logger.Debug("Successfully moveTo=", moveTo, "num=", num)

							mod := time.Duration(math.RoundToEven(float64(num / 1000)))
							time.Sleep(mod * time.Second)
						}
					})
				}) // source search
			}) // source select
		}) // target dial and login
	}) // source dial and login
} // imapBackup
