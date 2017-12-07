package network

import (
	"net"
	"log"
)

type tcpSocket struct {
	Config
}

// 轮询tcp network
func (t *tcpSocket) listener(h handle) {
	addr, err := t.getAddress()
	if err != nil {

	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		go h.HandleFunc(client)
	}
}

type handleTcp struct {
	fn func(con net.Conn)
}

func NewTcpHandle(f func(con net.Conn)) handle {
	return &handleTcp{
		fn: f,
	}

}

func (h *handleTcp) HandleFunc(con net.Conn) {
	h.fn(con)
}

func (t *tcpSocket) Server(h handle) {
	t.listener(h)

}
