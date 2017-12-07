package network

import (
	"net"
	"fmt"
)

type Proto int

const (
	Tcp Proto = iota

)

type listenInterface interface {
	listener(handle)
	Server(handle)
} 
type handle interface {
	HandleFunc(net.Conn)

}

func (t *Config) getAddress() (addr string, err error) {
	addr = fmt.Sprintf("%s:%s", t.HostAddress, t.HostPort)
	return addr, nil
}

type Config struct {
	Proto       Proto
	HostAddress string
	HostPort    string
}
