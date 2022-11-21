package config

type Configuration struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	Host            string   `yaml:"host"`
	Port            int      `yaml:"port"`
	User            string   `yaml:"user"`
	KeyFile         string   `yaml:"key_file"`
	KeyFilePassword string   `yaml:"key_file_password"`
	Password        string   `yaml:"password"`
	Tunnels         []Tunnel `yaml:"tunnels"`
}

type Tunnel struct {
	RemotePort int    `yaml:"remote_port"`
	RemoteHost string `yaml:"remote_host"`
	LocalPort  int    `yaml:"local_port"`
}
