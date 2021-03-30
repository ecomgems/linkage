package tunnels

import (
	"fmt"
	"github.com/ecomgems/linkage/utils/config"
	"github.com/urfave/cli/v2"
	"log"
)

func TunnelCmdHandler(c *cli.Context) error {
	configFilePath := c.Path("config")
	configData, err := config.GetConfiguration(configFilePath)
	if err != nil {
		return err
	}

	for _, serverConfig := range configData.Servers {
		log.Println(serverConfig.Host)
	}

	fmt.Println(configData)

	return nil
}
