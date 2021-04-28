package printer

import (
	"github.com/ecomgems/linkage/utils/tunnel"
	"log"
)

func Print(tunnels []tunnel.Tunnel) error {
	for _, t := range tunnels {
		log.Println(t.GetTunnelId())
	}

	return nil
}
