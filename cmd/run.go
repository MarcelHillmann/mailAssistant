package cmd

import (
	"github.com/urfave/cli/v2"
	"mailAssistant/account"
	"mailAssistant/cntl"
	"mailAssistant/rules"
	"os"
	"syscall"
)

// RunAssistant is execute the main logic
func RunAssistant(c *cli.Context) error {
	_ = c.String("config")

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