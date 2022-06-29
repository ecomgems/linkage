package tunnel

import "fmt"

func (t *Tunnel) GetTunnelId() string {
	tunnelId := fmt.Sprintf(
		"%d:%s:%d over %s@%s",
		t.TunnelConfig.RemotePort,
		t.TunnelConfig.RemoteHost,
		t.TunnelConfig.LocalPort,
		t.ServerConfig.User,
		t.ServerConfig.Host,
	)

	if t.ServerConfig.Port != 22 {
		tunnelId = fmt.Sprintf("%s:%d", tunnelId, t.ServerConfig.Port)
	}

	return tunnelId
}
