package config

import (
	"reflect"
	"testing"
)

func TestGetConfiguration(t *testing.T) {
	testCases := []struct {
		name     string
		fileName string
		want     Configuration
		wantErr  bool
	}{
		{
			name:     "read file with multiple servers",
			fileName: "reader_test/fixture.1.many_servers.yml",
			want: Configuration{
				Servers: []Server{
					{
						Host:     "remote.server1.tld",
						Port:     22,
						User:     "dev",
						KeyFile:  "~/.ssh/id_rsa.pub",
						Password: "",
						Tunnels:  []Tunnel{
							{
								RemotePort: 80,
								RemoteHost: "web_app",
								LocalPort:  81,
							},
							{
								RemotePort: 443,
								RemoteHost: "web_app",
								LocalPort:  445,
							},
						},
					},
					{
						Host:     "remote.server2.tld",
						Port:     23,
						User:     "dev",
						KeyFile:  "",
						Password: "<password>",
						Tunnels:  []Tunnel{
							{
								RemotePort: 9200,
								RemoteHost: "elastic",
								LocalPort:  9201,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "read file with one servers",
			fileName: "reader_test/fixture.2.one_server.yml",
			want: Configuration{
				Servers: []Server{
					{
						Host:     "remote.server1.tld",
						Port:     22,
						User:     "dev",
						KeyFile:  "~/.ssh/id_rsa.pub",
						Password: "",
						Tunnels:  []Tunnel{
							{
								RemotePort: 80,
								RemoteHost: "web_app",
								LocalPort:  80,
							},
							{
								RemotePort: 443,
								RemoteHost: "web_app",
								LocalPort:  443,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "read file with broken syntax",
			fileName: "reader_test/fixture.3.broken.yml",
			want: Configuration{},
			wantErr: true,
		},
		{
			name:     "try to open non existent file",
			fileName: "reader_test/fixture.4.non_existent_file.yml",
			want: Configuration{},
			wantErr: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := GetConfiguration(testCase.fileName)
			if (err != nil) != testCase.wantErr {
				t.Errorf(
					"GetConfiguration() error = %v, wantErr %v",
					err,
					testCase.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf(
					"GetConfiguration() got = %v, want %v",
					got,
					testCase.want,
				)
			}
		})
	}
}
