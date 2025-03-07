package server

type Address struct {
	IP   string
	Port string
}

func (addr *Address) String() string {
	return addr.IP + ":" + addr.Port
}
