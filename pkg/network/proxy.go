package network

import (
	"net"
	"io"
)

func proxy(svcConn net.Conn, cliConn net.Conn)  {
	go io.Copy(svcConn,  cliConn)
	io.Copy(cliConn, svcConn)
}