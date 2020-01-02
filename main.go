package main

import (
	"github.com/urfave/cli"
	"mailAssistant/appCmd"
	"os"
)

var version string

func main() {
	app := cli.NewApp()
	app.Author = "Marcel Hillmann"
	app.Name = "mailAssistant"
	app.Version = version
	app.Copyright = "(c) 2020 mahillmann.de"
	app.Usage = "automation for my Mail Accounts, like Outlook rules"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "configs",
			Value:  "",
			Usage:  "where to find the configs",
			EnvVar: "CONFIG_PATH",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "execute the mailAssistant",
			Action:  appCmd.RunAssistant,
		},
	}

	app.Run(os.Args)

}
