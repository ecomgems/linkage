package tunnels

import (
	"github.com/ecomgems/linkage/utils/config"
	"github.com/ecomgems/linkage/utils/printer"
	"github.com/ecomgems/linkage/utils/tunnel"
	"github.com/urfave/cli/v2"
)

// Function TunnelCmdHandler the execution of the application.
// It opens tunnels using a configuration file and manages it
// until all tunnels will be closed.
func TunnelCmdHandler(c *cli.Context) error {
	configFilePath := c.Path("config")
	configData, err := config.GetConfiguration(configFilePath)
	if err != nil {
		return err
	}

	tunnels := []tunnel.Tunnel{}
	for _, serverConfig := range configData.Servers {
		for _, tunnelConfig := range serverConfig.Tunnels {
			newTunnel := tunnel.Create(serverConfig, tunnelConfig)
			tunnels = append(tunnels, newTunnel)
		}
	}

	for {
		printer.Print(tunnels)
	}

	return nil
}
