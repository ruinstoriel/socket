package tls_package

import (
	"crypto/tls"
	"fmt"
	"net"
)

func TlsServer(addr string) {
	listener, err := tls.Listen("tcp", addr, config())
	if err != nil {
		panic(err)
	}

	go accept(listener)

}
func accept(li net.Listener) {
	defer li.Close()
	for {
		conn, _ := li.Accept()
		go handleConnection(conn)
	}

}

func handleConnection(con net.Conn) {
	defer con.Close()
	// 如果没有读完就开始回写，会发生什么?
	for {
		b := make([]byte, 1024)
		r, err := con.Read(b)
		if err != nil {
			panic(err)
		}
		if r > 0 {
			fmt.Printf("读取了%d, 内容是%s \n", r, string(b[:r]))
		}

		w, _ := con.Write([]byte("已阅"))
		if w > 0 {
			fmt.Printf("发送了%d \n", w)
		}

	}
}
func config() *tls.Config {
	cer, err := tls.LoadX509KeyPair("domain.crt", "domain.key")
	if err != nil {
		panic(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cer},
		MinVersion:   tls.VersionTLS13, // 指定最低版本为TLS 1.2
	}
	return config
}
