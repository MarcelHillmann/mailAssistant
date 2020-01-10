package main

import (
	"github.com/urfave/cli/v2"
	"log"
	cmd "mailAssistant/cmd"
	"mailAssistant/cmd/testDriver"
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
		{
			Name:   "test",
			Action: testDriver.TestTreiber,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Usage:    "Username for connection",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "password",
					Usage:    "Password for connection",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "server",
					Usage:    "server for connection",
					Required: true,
				},
				&cli.PathFlag{
					Name:     "file",
					Usage:    "rule to run",
					Required: true,
				},
				&cli.BoolFlag{
					Name:        "verbose",
					Required:    false,
					Hidden:      false,
					Value:       false,
				},
			},
		},
		{
			Name: "test_allMsg",
			Action: testDriver.TestDriverAllMsg,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Usage:    "Username for connection",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "password",
					Usage:    "Password for connection",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "server",
					Usage:    "server for connection",
					Required: true,
				},
				&cli.PathFlag{
					Name:     "select",
					Usage:    "select box",
					Required: true,
				},
				&cli.BoolFlag{
					Name:        "verbose",
					Required:    false,
					Hidden:      false,
					Value:       false,
				},
			},
		},
	}

	app.Run(os.Args)
}
