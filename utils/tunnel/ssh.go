package tunnel

import (
	"github.com/ecomgems/linkage/utils/config"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

func NewSshConfig(serverConfig config.Server) (*ssh.ClientConfig, error) {
	return nil, nil
}

func GetFullKeyPath(keyPath string) string {
	return strings.Replace(keyPath, "~", os.Getenv("HOME"), 1)
}
