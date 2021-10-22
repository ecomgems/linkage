package tunnel

import (
	"fmt"
	"github.com/ecomgems/linkage/utils/config"
	"github.com/ecomgems/linkage/utils/runtime"
	"github.com/ecomgems/sshtun"
	"log"
	"os"
	"strings"
	"time"
)

type Tunnel struct {
	ServerConfig config.Server
	TunnelConfig config.Tunnel
	TxCount      int64
	RxCount      int64
	Error        error
	Debug        bool
}

func Create(appRuntime runtime.ApplicationRuntime, serverConfig config.Server, tunnelConfig config.Tunnel) *Tunnel {
	t := Tunnel{
		ServerConfig: serverConfig,
		TunnelConfig: tunnelConfig,
		TxCount:      0,
		RxCount:      0,
		Error:        nil,
	}

	go t.Open(appRuntime)

	return &t
}

func (t *Tunnel) Open(runtime runtime.ApplicationRuntime) {
	log.Println("open:", t.GetTunnelId())

	sshTun := sshtun.New(
		t.TunnelConfig.LocalPort,
		t.ServerConfig.Host,
		t.TunnelConfig.RemotePort,
	)

	sshTun.SetPort(t.ServerConfig.Port)
	sshTun.SetUser(t.ServerConfig.User)
	sshTun.SetPassword(t.ServerConfig.Password)
	sshTun.SetKeyFile(
		getFullKeyPath(t.ServerConfig.KeyFile),
	)
	sshTun.SetRemoteHost(t.TunnelConfig.RemoteHost)
	sshTun.SetTimeout(365 * 24 * time.Hour)

	sshTun.SetDebug(runtime.IsDebugMode)

	if runtime.IsDebugMode {
		sshTun.SetConnState(func(tun *sshtun.SSHTun, state sshtun.ConnState) {
			switch state {
			case sshtun.StateStarting:
				log.Println("STATE is Starting", t.GetTunnelId())
			case sshtun.StateStarted:
				log.Println("STATE is Started", t.GetTunnelId())
			case sshtun.StateStopped:
				log.Println("STATE is Stopped", t.GetTunnelId())
			}
		})
	}

	go func() {
		for {
			if err := sshTun.Start(); err != nil {
				log.Println("SSH tunnel stopped:", err.Error(), t.GetTunnelId())
				time.Sleep(time.Second)
			}
		}
	}()
}

func (t *Tunnel) GetTunnelId() string {
	return fmt.Sprintf(
		"%d:%s:%d over %s@%s:%d",
		t.TunnelConfig.RemotePort,
		t.TunnelConfig.RemoteHost,
		t.TunnelConfig.LocalPort,
		t.ServerConfig.User,
		t.ServerConfig.Host,
		t.ServerConfig.Port,
	)
}

func getFullKeyPath(keyPath string) string {
	return strings.Replace(keyPath, "~", os.Getenv("HOME"), 1)
}
