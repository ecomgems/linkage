package tunnel

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"sync"
	"time"

	"github.com/ecomgems/linkage/utils/config"
	"github.com/ecomgems/linkage/utils/runtime"
)

type State int

const (
	Starting State = iota
	Started
	Stopped
)

type Tunnel struct {
	*sync.Mutex

	ServerConfig config.Server
	TunnelConfig config.Tunnel
	LoggerCh     chan string
	StatsCh      chan *Stats

	ctx        context.Context
	cancel     context.CancelFunc
	stats      *Stats
	timer      chan time.Time
	errChannel chan error
	sshConfig  *ssh.ClientConfig
}

func NewTunnel(appRuntime runtime.ApplicationRuntime, serverConfig config.Server, tunnelConfig config.Tunnel) *Tunnel {

	statsChannel := make(chan *Stats)
	stats := NewStats(Starting, func(stats *Stats) {
		statsChannel <- stats
	})

	sshConfig, err := NewSshConfig(serverConfig)
	if err != nil {
		panic(err)
	}

	t := Tunnel{
		ServerConfig: serverConfig,
		StatsCh:      statsChannel,
		TunnelConfig: tunnelConfig,
		LoggerCh:     make(chan string),
		stats:        stats,
		sshConfig:    sshConfig,
	}

	go t.Init(appRuntime)

	return &t
}

func (t *Tunnel) Init(runtime runtime.ApplicationRuntime) {
	t.LoggerCh <- fmt.Sprint("init:", t.GetTunnelId())

	go func() {
		for {
			if err := t.Start(); err != nil {
				t.LoggerCh <- fmt.Sprint("SSH tunnel stopped:", err.Error(), t.GetTunnelId())
				t.LoggerCh <- fmt.Sprint("wait for 1 second before restart...")

				time.Sleep(time.Second)
			}
		}
	}()
}

func (t *Tunnel) Start() error {
	t.Lock()

	localListener, err := net.Listen("tcp", fmt.Sprintf(":%d", t.TunnelConfig.LocalPort))
	if err != nil {
		return t.errorWhileNotStarted(err)
	}

	t.ctx, t.cancel = context.WithCancel(context.Background())
	t.errChannel = make(chan error)

	go func() {
		for {
			incomingConn, err := localListener.Accept()
			if err != nil {
				t.errorWhenStarted(fmt.Errorf("local accept on :%d failed: %s", t.TunnelConfig.LocalPort, err.Error()))
				break
			}

			t.LoggerCh <- fmt.Sprintf("accepted connection from %s", incomingConn.RemoteAddr().String())

			// Launch the forward
			go t.forward(incomingConn)
			t.stats.AddConnQty()
		}
	}()

	go func() {
		select {
		case <-t.ctx.Done():
			localListener.Close()
			t.stats.SetConnQty(0)
		}
	}()

	t.stats.UpdateState(Started)
	t.LoggerCh <- fmt.Sprintf("listening on :%d", t.TunnelConfig.LocalPort)

	t.Unlock()

	select {
	case err := <-t.errChannel:
		return err
	}
}

func (t *Tunnel) Stop() error {
	//TODO implement me
	panic("implement me")
}

func (t *Tunnel) isStarted() bool {
	return t.stats.State == Started
}

func (t *Tunnel) forward(conn net.Conn) {

}
