package quic_demo

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/metacubex/quic-go"
	"github.com/metacubex/quic-go/http3"
	"io"
	"net"
	"net/http"
)

const (
	FrameTypeTCPRequest = 0x401

	// Max length values are for preventing DoS attacks

	MaxAddressLength = 2048
	MaxMessageLength = 2048
	MaxPaddingLength = 4096

	MaxUDPSize = 4096

	maxVarInt1 = 63
	maxVarInt2 = 16383
	maxVarInt4 = 1073741823
	maxVarInt8 = 4611686018427387903
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
	ln, err := tr.Listen(config(), &quic.Config{
		MaxIncomingStreams: 100,
	})
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

type MyHandler struct {
	http.Handler
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("请求了")
	data, _ := io.ReadAll(r.Body)
	fmt.Println(string(data))
	w.Write([]byte("已阅"))
}

func quicServer_(addr string) {
	udpAddr, _ := net.ResolveUDPAddr("udp4", addr)
	udpConn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		panic(err)
	}
	l, err := quic.Listen(udpConn, config(), &quic.Config{MaxIncomingStreams: 100})
	s := http3.Server{
		Handler: &MyHandler{},
		StreamHijacker: func(frameType http3.FrameType, connection quic.Connection, stream quic.Stream, err2 error) (hijacked bool, err error) {
			fmt.Println(frameType == FrameTypeTCPRequest)
			go handleStream(stream)
			return true, nil
		},
	}
	c, err := l.Accept(context.Background())
	err = s.ServeQUICConn(c)
	if err != nil {
		panic(err)
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
func quicHandle(con quic.Connection) {

	for {
		s, err := con.AcceptStream(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		go handleStream(s)
	}

}

func handleStream(s quic.Stream) {

	b := make([]byte, 1024)
	i, err := s.Read(b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("内容是%s \n", string(b[:i]))
	w, _ := s.Write([]byte("已阅"))
	if w > 0 {
		fmt.Printf("发送了%d \n", w)
	}
	//s.CancelRead(0)
	err = s.Close()
	_, err = io.Copy(io.Discard, s)

}
