package tunnels

import (
	"github.com/ecomgems/linkage/utils/config"
	"github.com/ecomgems/linkage/utils/runtime"
	"github.com/ecomgems/linkage/utils/tunnel"
	"github.com/urfave/cli/v2"
	"log"
)

// TunnelCmdHandler is the main execution function of the application.
// It opens tunnels using a configuration file and manages it
// until all tunnels will be closed.
func TunnelCmdHandler(c *cli.Context) error {
	runtime := runtime.ApplicationRuntime{
		IsDebugMode: c.Bool("debug"),
	}
	configFilePath := c.Path("config")
	configData, err := config.GetConfiguration(configFilePath)
	if err != nil {
		return err
	}

	var tunnels []*tunnel.Tunnel
	for _, serverConfig := range configData.Servers {
		for _, tunnelConfig := range serverConfig.Tunnels {
			newTunnel := tunnel.NewTunnel(runtime, serverConfig, tunnelConfig)

			go func() {
				//@todo Revise logs usage
				for {
					var logMessage string
					logMessage = <-newTunnel.LoggerCh
					if runtime.IsDebugMode {
						log.Println(logMessage)
					}
				}
			}()

			go func() {
				//@todo Revise stats usage
				for {
					var statsMessage *tunnel.Stats
					statsMessage = <-newTunnel.StatsCh
					if runtime.IsDebugMode {
						log.Println(newTunnel.GetTunnelId(), statsMessage)
					}
				}
			}()

			tunnels = append(tunnels, newTunnel)
		}
	}

	wait := make(chan bool)
	<-wait

	return nil
}
