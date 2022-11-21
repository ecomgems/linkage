package tunnel

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/ecomgems/linkage/utils/config"
	"golang.org/x/crypto/ssh"
)

func NewSshConfig(serverConfig config.Server) (*ssh.ClientConfig, error) {
	config := &ssh.ClientConfig{
		User: serverConfig.User,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 15 * time.Second,
	}

	authMethod, err := getSSHAuthMethod(serverConfig)
	if err != nil {
		return nil, err
	}

	config.Auth = []ssh.AuthMethod{authMethod}

	return config, nil
}

func getSSHAuthMethod(serverConfig config.Server) (ssh.AuthMethod, error) {
	if serverConfig.Password != "" {
		return ssh.Password(serverConfig.Password), nil
	}

	var keyFile string
	if serverConfig.KeyFile == "" {
		usr, _ := user.Current()
		if usr != nil {
			keyFile = usr.HomeDir + "/.ssh/id_rsa"
		} else {
			keyFile = "/root/.ssh/id_rsa"
		}

	} else {
		keyFile = GetFullKeyPath(serverConfig.KeyFile)
	}

	var key ssh.Signer
	var err error

	buf, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("error reading SSH key file %s: %s", serverConfig.KeyFile, err.Error())
	}

	encrypted := serverConfig.KeyFilePassword != ""
	if encrypted {
		key, err = ssh.ParsePrivateKeyWithPassphrase(buf, []byte(serverConfig.KeyFilePassword))
		if err != nil {
			return nil, fmt.Errorf("error parsing encrypted key: %s", err.Error())
		}
	} else {
		key, err = ssh.ParsePrivateKey(buf)
		if err != nil {
			return nil, fmt.Errorf("error parsing key: %s", err.Error())
		}
	}

	return ssh.PublicKeys(key), nil
}

func GetFullKeyPath(keyPath string) string {
	return strings.Replace(keyPath, "~", os.Getenv("HOME"), 1)
}
