package main

import (
	"github.com/urfave/cli/v2"
	"log"
	cmd "mailAssistant/cmd"
	"os"
)

var version string

func main() {
	app := cli.NewApp()
	app.Authors = []*cli.Author{
		{
			Name: "Marcel Hillmann",
		},
	}
	app.Name = "mailAssistant"
	app.Version = version
	app.Copyright = "(c) 2020 mahillmann.de"
	app.Usage = "automation for my Mail Accounts, like Outlook rules"
	app.EnableBashCompletion = true
	app.ExitErrHandler = func(context *cli.Context, err error) {
		log.Panic(err)
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "configs",
			Value:    "",
			Usage:    "where to find the configs",
			EnvVars:  []string{"CONFIG_PATH"},
			Required: false,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "execute the mailAssistant",
			Action:  cmd.RunAssistant,
		},
		{
			Name:    "verify",
			Aliases: []string{},
			Usage:   "verify rules",
			Action:  cmd.RunConfigCheck,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "config",
					Usage:    "path to the rules",
					Required: true,
				},
			},
		},
	}

	app.Run(os.Args)
}
