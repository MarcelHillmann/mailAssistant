package cmd

import (
	"github.com/urfave/cli/v2"
	"log"
	"mailAssistant/account"
	"mailAssistant/cntl"
	"mailAssistant/monitoring"
	"mailAssistant/rules"
	"os"
)

// RunAssistant is execute the main logic
func RunAssistant(c *cli.Context) error {
	log.Print(">> RunAssistant")
	zs := c.String("zipkin_server")

	if accounts, err := account.ImportAccounts(); err != nil {
		log.Print("<< RunAssistant",err)
		return err
	} else if err := rules.ImportAndLaunch(accounts); err != nil {
		log.Print("<< RunAssistant", err)
		return err
	}else if err := monitoring.StartServer(zs); err != nil {
		log.Print("<< RunAssistant", err)
		return err
	} else {
		cntl.WaitForOsNotify(os.Interrupt, os.Kill)
		cntl.WaitForNotify()
		cntl.StopAllClocks()
	}
	log.Print("<< RunAssistant nil")
	return nil
}