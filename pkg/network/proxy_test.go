package network_test

import (
	"net/http"
	"strings"
	"fmt"
	"testing"
	"net/url"
	"github.com/koko990/myproxy/pkg/network"
	"net"
	"io"
	"time"
	"runtime/pprof"
	"os"
)

var httpCh = make(chan struct{})

func handleProxy(conn net.Conn) {
	fmt.Println("handle start")
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
	fmt.Println(str)
	// 将请求通过socket复制到新的地址
	w, err := net.Dial("tcp", "127.0.0.1:27777")
	go func() {
		fmt.Println("将请求通过socket复制到新的地址")
		_, err = io.Copy(w, conn)
		fmt.Println("将请求通过socket复制到新的地址完毕")
	}()

	if err != nil {
		panic(err)
	}
	// 将收到的消息回拷
	go func() {
		fmt.Println("将收到的消息回拷")
		_, err = io.Copy(conn, w)
		fmt.Println("将收到的消息完毕")
	}()
	if err != nil {
		panic(err)
	}
	fmt.Println("handle done")
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
	//错误分析
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
	fmt.Println("listen 127.0.0.1:27777")
	http.ListenAndServe("127.0.0.1:27777", handleF)
}

var handleF http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
	var buf = make([]byte, 13)
	req := request.Body
	req.Read(buf)
	writer.Write(buf)
}

func writeHttpProxy() {
	fmt.Println("write to http://127.0.0.1:27777 proxy 28888")
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
