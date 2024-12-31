package quic_demo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/metacubex/quic-go"
)

func TestQuicServer(t *testing.T) {
	addr := "127.0.0.1:8080"
	go quicServer_(addr)
	fmt.Println("server 加载完毕")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second) // 3s handshake timeout
	con, err := quic.DialAddr(ctx, addr, clientConfig(), &quic.Config{})
	if err != nil {
		panic(err)
	}
	defer cancel()
	count := 0
	for i := 0; i < 101; i++ {
		count++
		fmt.Print(count)
		b := make([]byte, 100)

		j := varintPut(b, FrameTypeTCPRequest)
		copy(b[j:], "are you ok")
		str, err2 := con.OpenStream()
		if err2 != nil {
			panic(err2)
		}
		if nerr, ok := err2.(net.Error); ok && nerr.Timeout() {
			panic(nerr)
		}
		_, err3 := str.Write(b)
		if err3 != nil {
			log.Fatalf(err3.Error())
		}
		r := make([]byte, 100)

		ss, err5 := str.Read(r)
		fmt.Println(string(r[:ss]))

		if err5 != io.EOF {
			//panic(err5)
			// io.Copy(io.Discard, str)
			//str.CancelRead(0)
		}
		//str.CancelRead(0)
		_ = str.Close()
		_, err4 := io.Copy(io.Discard, str)
		if err4 != nil {
			panic(err5)
		}

	}

}

func varintPut(b []byte, i uint64) int {
	if i <= maxVarInt1 {
		b[0] = uint8(i)
		return 1
	}
	if i <= maxVarInt2 {
		b[0] = uint8(i>>8) | 0x40
		b[1] = uint8(i)
		return 2
	}
	if i <= maxVarInt4 {
		b[0] = uint8(i>>24) | 0x80
		b[1] = uint8(i >> 16)
		b[2] = uint8(i >> 8)
		b[3] = uint8(i)
		return 4
	}
	if i <= maxVarInt8 {
		b[0] = uint8(i>>56) | 0xc0
		b[1] = uint8(i >> 48)
		b[2] = uint8(i >> 40)
		b[3] = uint8(i >> 32)
		b[4] = uint8(i >> 24)
		b[5] = uint8(i >> 16)
		b[6] = uint8(i >> 8)
		b[7] = uint8(i)
		return 8
	}
	panic(fmt.Sprintf("%#x doesn't fit into 62 bits", i))
}

func clientConfig() *tls.Config {
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		log.Println("Error reading CA certificate:", err)
		return nil
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	conf := &tls.Config{
		// InsecureSkipVerify: true,
		ServerName: "127.0.0.1",
		MinVersion: tls.VersionTLS13, // 指定最低版本为TLS 1.2
		RootCAs:    caCertPool,
	}
	return conf
}
