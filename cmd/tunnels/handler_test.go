package tunnels

import (
	"github.com/urfave/cli/v2"
	"testing"
)

func TestTunnelCmdHandler(t *testing.T) {
	type args struct {
		c *cli.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := TunnelCmdHandler(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("TunnelCmdHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
