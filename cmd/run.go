package cmd

import (
	"github.com/urfave/cli/v2"
	"mailAssistant/account"
	"mailAssistant/cntl"
	"mailAssistant/rules"
	"os"
	"syscall"
)

func RunAssistant(c *cli.Context) error {
	if accounts, err := account.ImportAccounts(); err != nil {
		return err
	} else if err := rules.ImportAndLaunch(accounts); err != nil {
		return err
	} else {
		cntl.WaitForOsNotify(os.Interrupt, syscall.SIGTERM)
		cntl.WaitForNotify()
		cntl.StopAllClocks()
	}
	return nil
}