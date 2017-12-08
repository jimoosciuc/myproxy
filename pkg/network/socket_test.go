package network_test

import (
	"testing"
	"github.com/koko990/myproxy/pkg/network"
	"net"
	"fmt"
)

var ch = make(chan struct{})

func handle(conn net.Conn) {
	var buf = make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	var bufFormat []byte
	for k, v := range buf {
		if v != 0 {
			bufFormat = append(bufFormat, buf[k])
		}
	}
	str := fmt.Sprintf("read form socket done: {\n%s\n}", bufFormat)
	fmt.Println(str)
	ch <- struct{}{}
}
func TestNewSocket(t *testing.T) {
	var h = network.NewTcpHandle(handle)
	var cfg = network.Config{
		Proto:       network.Tcp,
		HostAddress: "127.0.0.1",
		HostPort:    "28088",
	}
	sock := network.NewSocket(cfg)
	go sock.Server(h)
	go writeSocket()
	<-ch
}

func writeSocket() {
	w, err := net.Dial("tcp", "127.0.0.1:28088")
	if err != nil {
		panic(err)
	}
	var buf = []byte("hello world !")
	i, err := w.Write(buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("write done, len is : %d\n", i)
}
