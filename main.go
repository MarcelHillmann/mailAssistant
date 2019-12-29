package main

import (
	"mailAssistant/account"
	"mailAssistant/cntl"
	"mailAssistant/logging"
	"mailAssistant/rules"
	"os"
	"syscall"
)

func main() {
	if accounts, err := account.ImportAccounts(); err != nil {
		logging.NewGlobalLogger().Panic(err)
	} else if err := rules.ImportAndLaunch(accounts); err != nil {
		logging.NewGlobalLogger().Panic(err)
	} else {
		cntl.WaitForOsNotify(os.Interrupt, syscall.SIGTERM)
		cntl.WaitForNotify()
		cntl.StopAllClocks()
	}
}
