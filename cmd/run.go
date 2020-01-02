package cmd

import (
	"github.com/urfave/cli/v2"
	"log"
	"mailAssistant/account"
	"mailAssistant/cntl"
	"mailAssistant/rules"
	"os"
	"syscall"
)

// RunAssistant is execute the main logic
func RunAssistant(c *cli.Context) error {
	log.Print(">> RunAssistant")
	_ = c.String("config")

	if accounts, err := account.ImportAccounts(); err != nil {
		log.Print("<< RunAssistant",err)
		return err
	} else if err := rules.ImportAndLaunch(accounts); err != nil {
		log.Print("<< RunAssistant",err)
		return err
	} else {
		cntl.WaitForOsNotify(os.Interrupt, syscall.SIGTERM)
		cntl.WaitForNotify()
		cntl.StopAllClocks()
	}
	log.Print("<< RunAssistant nil")
	return nil
}