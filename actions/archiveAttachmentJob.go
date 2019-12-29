package actions

import (
	"io/ioutil"
	"mailAssistant/account"
	"mailAssistant/logging"
	"os"
	"path/filepath"
)

func init() {
	register("archiveAttachment", newArchiveAttachment)
}

func newArchiveAttachment(job Job, waitGroup *int32) {
	logger := logging.NewLogger()
	if isLockedElseLock(logger, waitGroup) {
		return
	}
	defer unlockAlways(logger, waitGroup)

	if _, err := os.Stat(job.getSaveTo()); os.IsNotExist(err) {
		os.MkdirAll(job.getSaveTo(), 0)
	}
	job.GetAccount(job.GetString("mail_account")).
		DialAndLoginPromise(func(promise *account.ImapPromise) {
			promise.SelectPromise(job.GetString("path"), job.GetBool("readonly"), func(promise *account.ImapPromise) {
				promise.SearchPromise(job.getSearchParameter(), true, func(msgPromises *account.MsgPromises) {
					attachType := job.GetString("attachment_type")
					attachmentPromises := msgPromises.GetAttachments(attachType)
					for _, attachmentPromise := range attachmentPromises {
						saveTo := filepath.Join(job.getSaveTo(), attachmentPromise.GetFilename())
						if fileInfo, err := os.Stat(saveTo); os.IsExist(err) || fileInfo != nil {
							logger.Debug(saveTo, " ", err)
						} else if err := ioutil.WriteFile(saveTo, attachmentPromise.Body(), 0); err != nil {
							logger.Severe(err)
						} else {
							logger.Debug("saved ", saveTo)
						}
					}
				}) // SearchPromise
			}) // SelectPromise
		}) // DialAndLoginPromise
	// GetAccount
}
