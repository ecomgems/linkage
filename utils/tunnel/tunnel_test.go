package tunnel

import (
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ecomgems/linkage/utils/config"
	"github.com/ecomgems/linkage/utils/runtime"
	"github.com/gliderlabs/ssh"
	"github.com/stretchr/testify/assert"
)

type TestServerHandler struct {
}

func (t TestServerHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	wr.Write([]byte("Hello!"))
}

func TestTunnel_Open_Send_A_Bites(t *testing.T) {
	// 1. Open server on port 58880
	handler := &TestServerHandler{}
	go func() {
		err := http.ListenAndServe(":58880", handler)
		if err != nil {
			t.Fail()
		}
	}()

	// 2. Open SSH server on port 58890
	go func() {
		forwardHandler := &ssh.ForwardedTCPHandler{}
		server := ssh.Server{
			LocalPortForwardingCallback: ssh.LocalPortForwardingCallback(func(ctx ssh.Context, dhost string, dport uint32) bool {
				return true
			}),
			Addr: ":58890",
			Handler: ssh.Handler(func(s ssh.Session) {
				io.WriteString(s, "Remote forwarding available...\n")
				select {}
			}),
			ReversePortForwardingCallback: ssh.ReversePortForwardingCallback(func(ctx ssh.Context, host string, port uint32) bool {
				return true
			}),
			RequestHandlers: map[string]ssh.RequestHandler{
				"tcpip-forward":        forwardHandler.HandleSSHRequest,
				"cancel-tcpip-forward": forwardHandler.HandleSSHRequest,
			},
			ChannelHandlers: map[string]ssh.ChannelHandler{
				"direct-tcpip": ssh.DirectTCPIPHandler,
				"session":      ssh.DefaultSessionHandler,
			},
		}
		log.Fatal(server.ListenAndServe())
	}()

	// 3. Open tunnel to 58880 through 58890 on port 58870
	go func() {
		tunnel := &Tunnel{
			ServerConfig: config.Server{
				Host: "localhost",
				Port: 58890,
				User: "root",
				//KeyFile:  "",
				Password: "pass",
			},
			TunnelConfig: config.Tunnel{
				RemotePort: 58880,
				RemoteHost: "127.0.0.1",
				LocalPort:  58870,
			},
		}

		tunnel.Open(runtime.ApplicationRuntime{})
	}()

	// 4. Connect to 58870 and send a bites
	timer := time.NewTimer(100 * time.Millisecond)
	<-timer.C

	resp, err := http.Get("http://127.0.0.1:58870/")
	if err != nil {
		panic(err)
		return
	}
	defer resp.Body.Close()

	// 5. Verify if bites were received on the 58880 side
	body, err := io.ReadAll(resp.Body)
	assert.Equal(t, body, []byte("Hello!"))
}

func TestTunnel_GetTunnelId(t *testing.T) {
	testCases := []struct {
		name   string
		tunnel *Tunnel
		want   string
	}{
		{
			"Tunnel with standard SSH port",
			&Tunnel{
				ServerConfig: config.Server{
					Host: "example.com",
					Port: 22,
					User: "root",
				},
				TunnelConfig: config.Tunnel{
					RemotePort: 80,
					RemoteHost: "web",
					LocalPort:  8080,
				},
			},
			"80:web:8080 over root@example.com",
		},
		{
			"Tunnel with custom SSH port",
			&Tunnel{
				ServerConfig: config.Server{
					Host: "example.com",
					Port: 22022,
					User: "root",
				},
				TunnelConfig: config.Tunnel{
					RemotePort: 80,
					RemoteHost: "web",
					LocalPort:  8080,
				},
			},
			"80:web:8080 over root@example.com:22022",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, testCase.tunnel.GetTunnelId())
		})
	}
}

func TestGetFullKeyPath(t *testing.T) {
	origHomePath := os.Getenv("HOME")

	os.Setenv("HOME", "/home/user")
	testCases := []struct {
		keyPath string
		want    string
	}{
		{
			"/home/user/.ssh/id_rsa",
			"/home/user/.ssh/id_rsa",
		},
		{
			"~/.ssh/id_rsa",
			"/home/user/.ssh/id_rsa",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.keyPath, func(t *testing.T) {
			assert.Equal(t, testCase.want, GetFullKeyPath(testCase.keyPath))
		})
	}

	os.Setenv("HOME", origHomePath)
}
