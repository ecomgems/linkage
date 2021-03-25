package config

type Configuration struct {
	Servers []Server
}

type Server struct {
	Host     string
	Port     int
	User     string
	KeyFile  string
	Password string
	Tunnels  []Tunnel
}

type Tunnel struct {
	RemotePort int
	RemoteHost string
	LocalPort  int
}
