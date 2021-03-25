package main

import (
	"github.com/ecomgems/linkage/cmd/tunnels"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:        "linkage",
		Description: "Linkage is an SSH TCP Port Forwarding app to open multiple tunnels simultaneously.",
		Usage:       "linkage -c example.yml",
		Version:     "1.0.0",
		Action:      tunnels.TunnelCmdHandler,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Load configuration from `FILE`",
				Required: true,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
