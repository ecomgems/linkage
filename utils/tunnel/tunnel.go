package tunnel

import (
	"errors"
	"fmt"
	"github.com/ecomgems/linkage/utils/config"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

type Tunnel struct {
	ServerConfig config.Server
	TunnelConfig config.Tunnel
	TxCount      int64
	RxCount      int64
	Error        error
}

func Create(serverConfig config.Server, tunnelConfig config.Tunnel) Tunnel {
	t := Tunnel{
		ServerConfig: serverConfig,
		TunnelConfig: tunnelConfig,
		TxCount:      0,
		RxCount:      0,
		Error:        nil,
	}

	go t.Open()

	return t
}

func (t *Tunnel) Open()  {
	authMethods, err := t.getAuthMethods()
	if err != nil {
		t.Error = err
	}

	config := ssh.ClientConfig{
		User:            t.ServerConfig.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshAddress := fmt.Sprintf("%s:%d", t.ServerConfig.Host, t.ServerConfig.Port)
	sshConnection, err := ssh.Dial("tcp", sshAddress, &config)
	if err != nil {
		log.Fatalln(err)
	}
	defer sshConnection.Close()

	remoteAddress := fmt.Sprintf("%s:%d", t.TunnelConfig.RemoteHost, t.TunnelConfig.RemotePort)
	remoteConnection, err := sshConnection.Dial("tcp", remoteAddress)
	if err != nil {
		log.Fatalln(err)
	}

	localAddress := fmt.Sprintf("127.0.0.1:%d", t.TunnelConfig.LocalPort)
	localConnection, err := net.Listen("tcp", localAddress)
	if err != nil {
		log.Fatalln(err)
	}
	defer localConnection.Close()

	for {
		client, err := localConnection.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		t.handleClient(client, remoteConnection)
	}
}

func (t *Tunnel) handleClient(client net.Conn, remote net.Conn) {
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			log.Fatalln(err)
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

func (t *Tunnel) getAuthMethods() ([]ssh.AuthMethod, error) {
	var authMethods []ssh.AuthMethod

	if t.ServerConfig.KeyFile != "" {
		privateKey, err := parsePrivateKey(t.ServerConfig.KeyFile)
		if err != nil {
			return nil, err
		}

		authMethods = append(authMethods, ssh.PublicKeys(privateKey))
	}

	if t.ServerConfig.Password != "" {
		authMethods = append(authMethods, ssh.Password(t.ServerConfig.Password))
	}

	if len(authMethods) == 0 {
		return nil, errors.New(
			fmt.Sprintf(
				"at leat one auth method should be available for server %s:%d",
				t.ServerConfig.Host,
				t.ServerConfig.Port,
			),
		)
	}

	return authMethods, nil
}

func parsePrivateKey(keyPath string) (ssh.Signer, error) {
	keyFullPath := getFullKeyPath(keyPath)
	if _, err := os.Stat(keyFullPath); os.IsNotExist(err) {
		return nil, err
	}

	buff, _ := ioutil.ReadFile(keyFullPath)
	return ssh.ParsePrivateKey(buff)
}

func getFullKeyPath(keyPath string) string {
	return strings.Replace(keyPath, "~", os.Getenv("HOME"), 1)
}
