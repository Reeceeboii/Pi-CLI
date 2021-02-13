package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

var settings = Settings{
	PiHoleAddress: defaultAddress,
	PiHolePort:    defaultPort,
	RefreshS:      defaultRefreshS,
}

var app = cli.App{
	EnableBashCompletion: true,
	Name:                 "Pi-CLI",
	Usage:                "Third party program to retrieve and display Pi-Hole data right from your terminal.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "ipaddress",
			Value:       settings.PiHoleAddress,
			Required:    false,
			Aliases:     []string{"ip"},
			Destination: &settings.PiHoleAddress,
			Usage:       "The IP address of your Pi-Hole instance",
		},
		&cli.Int64Flag{
			Name:        "port",
			Value:       settings.PiHolePort,
			Required:    false,
			Aliases:     []string{"p"},
			Destination: &settings.PiHolePort,
			Usage:       "The port exposing your Pi-Hole's web interface",
		},
		&cli.Int64Flag{
			Name:        "refresh-rate",
			Value:       settings.RefreshS,
			Required:    false,
			Aliases:     []string{"r"},
			Destination: &settings.RefreshS,
			Usage:       "The rate (in seconds) at which the Pi-Hole API is polled for data",
		},
		&cli.BoolFlag{
			Name:     "save-to-file",
			Required: false,
			Aliases:  []string{"s"},
		},
	},
	Action: func(c *cli.Context) error {
		fmt.Println(settings)
		if len(c.FlagNames()) == 0 {
			fmt.Println("no flags provided")
			settings.loadFromFile()
		}
		return nil
	},
}
