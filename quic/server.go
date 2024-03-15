package quic_demo

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/quic-go/quic-go"
)

func quicServer(addr string) {
	udpAddr, _ := net.ResolveUDPAddr("udp4", addr)
	udpConn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		panic(err)
	}
	tr := quic.Transport{
		Conn: udpConn,
	}
	ln, err := tr.Listen(config(), &quic.Config{})
	if err != nil {
		panic(err)
	}
	// ... error handling
	go func() {
		for {
			conn, err := ln.Accept(context.Background())
			if err != nil {
				panic(err)
			}
			go quicHandle(conn)
		}
	}()
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
func quicHandle(con quic.Connection) {

	for {
		b := make([]byte, 1024)
		s, err := con.AcceptStream(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		r, err := s.Read(b)
		if err != nil {
			fmt.Println(err)
		}

		if r > 0 {
			fmt.Printf("读取了%d, 内容是%s \n", r, string(b))
		}

		w, _ := s.Write([]byte("已阅"))
		if w > 0 {
			fmt.Printf("发送了%d \n", w)
		}

	}

}
