package main

import (
	"github.com/ecomgems/linkage/cmd/tunnels"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:        "linkage.bak",
		Description: "Linkage is an SSH TCP Port Forwarding app to open multiple tunnels simultaneously.",
		Usage:       "linkage.bak -c example.yml",
		Version:     "1.0.0",
		Action:      tunnels.TunnelCmdHandler,
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Load configuration from `FILE`",
				Required: true,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "Print debug information into stdout",
				DefaultText: "yes",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
