package network_test

import (
	"net/http"
	"strings"
	"fmt"
	"testing"
	"net/url"
	"github.com/koko990/myproxy/pkg/network"
	"net"
	"time"
	"os"
	"runtime/pprof"
	"log"
)

var httpCh = make(chan struct{})

func handleProxy(conn net.Conn) {
	var buf = make([]byte, 1024)
	// 读取proxy收到的请求
	_, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	// 打印proxy得到的请求
	var bufFormat []byte
	for k, v := range buf {
		if v != 0 {
			bufFormat = append(bufFormat, buf[k])
		}
	}
	str := fmt.Sprintf("read form socket done: {\n%s\n}", bufFormat)
	l := log.New(os.Stdout,"",3)
	l.Println(str)
	// 将请求通过socket复制到新的地址
	w, err := net.Dial("tcp", "127.0.0.1:27777")
	w.Write(buf)
	var ss = make([]byte, 128 )
	w.Read(ss)

	if err != nil {
		panic(err)
	}
	// 将收到的消息回拷
	conn.Write(ss)
	if err != nil {
		panic(err)
	}
	ch <- struct{}{}
}
func TestProxy(t *testing.T) {
	var h = network.NewTcpHandle(handleProxy)
	var cfg = network.Config{
		Proto:       network.Tcp,
		HostAddress: "127.0.0.1",
		HostPort:    "28888",
	}
	sock := network.NewSocket(cfg)
	go sock.Server(h)
	go httpSvc()
	go writeHttpProxy()
	go func() {
		time.Sleep(time.Second * 5)
		p := pprof.Profiles()
		for _, v := range p {
			v.WriteTo(os.Stdout, 1)
		}
		close(httpCh)
		close(ch)

	}()
	<-httpCh
	<-ch
}
func TestNoProxy(t *testing.T) {
	go httpSvc()
	go writeHttp()
	<-httpCh
}
func httpSvc() {
	http.ListenAndServe("127.0.0.1:27777", handleF)
}

var handleF http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("haddle")
	var buf = make([]byte, 1024)
	req := request.Body
	req.Read(buf)
	writer.Write(buf)
}

func writeHttpProxy() {
	proxyUrl, err := url.Parse("http://127.0.0.1:28888")
	if err != nil {
		panic(err)
	}
	var buf = make([]byte, 13)
	proxy := http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	cli := http.Client{Transport: &proxy}
	req, err := cli.Post("http://127.0.0.1:27777", "", strings.NewReader("hello world !"))
	if err != nil {
		panic(err)
	}
	read := req.Body
	_, err = read.Read(buf)
	if err != nil {
		panic(err)
	}
	str := fmt.Sprintf("read from http: %s\n", buf)
	fmt.Println(str)
	httpCh <- struct{}{}
}

func writeHttp() {
	var buf = make([]byte, 13)
	req, err := http.NewRequest("POST", "http://127.0.0.1:27777", strings.NewReader("hello world !"))
	if err != nil {
		panic(err)
	}
	read := req.Body
	_, err = read.Read(buf)
	if err != nil {
		panic(err)
	}
	str := fmt.Sprintf("read from http: %s\n", buf)
	fmt.Println(str)
	httpCh <- struct{}{}
}
