package quic_demo

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/ruinstoriel/quic-go"
	"github.com/ruinstoriel/quic-go/http3"
	"github.com/ruinstoriel/quic-go/quicvarint"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
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
		StreamDispatcher: func(frameType http3.FrameType, stream *quic.Stream, err2 error) (hijacked bool, err error) {
			fmt.Println(frameType == FrameTypeTCPRequest)
			// 消耗掉
			bReader := quicvarint.NewReader(stream)
			_, err = quicvarint.Read(bReader)
			if err != nil {
				return false, err
			}
			
			go handleStream(stream)
			fmt.Println("---------------------------")
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
func quicHandle(con *quic.Conn) {

	for {
		s, err := con.AcceptStream(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		go handleStream(s)
	}

}

func handleStream(s *quic.Stream) {
	reqAddr, err := ReadTCPRequest(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("<UNK>", reqAddr)
}

func ReadTCPRequest(r io.Reader) (string, error) {
	bReader := quicvarint.NewReader(r)
	addrLen, err := quicvarint.Read(bReader)
	if err != nil {
		return "", err
	}
	log.Printf("addrLen: %s", strconv.Itoa(int(addrLen)))
	if addrLen == 0 || addrLen > MaxAddressLength {

	}

	addrBuf := make([]byte, addrLen)
	_, err = io.ReadFull(r, addrBuf)

	if err != nil {

		return "", err
	}
	paddingLen, err := quicvarint.Read(bReader)
	if err != nil {

		return "", err
	}
	if paddingLen > MaxPaddingLength {

	}
	if paddingLen > 0 {
		_, err = io.CopyN(io.Discard, r, int64(paddingLen))
		if err != nil {
			return "", err
		}
	}

	return string(addrBuf), nil
}
