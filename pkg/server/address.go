package server

type ServerAddress struct {
	Ip   string
	Port string
}

func (addr *ServerAddress) String() string {
	return addr.Ip + ":" + addr.Port
}
