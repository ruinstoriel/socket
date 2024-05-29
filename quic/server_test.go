package quic_demo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"github.com/metacubex/quic-go"
)

func TestQuicServer(t *testing.T) {
	addr := "127.0.0.1:8080"
	quicServer(addr)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3s handshake timeout
	con, err := quic.DialAddr(ctx, addr, clientConfig(), &quic.Config{})
	if err != nil {
		panic(err)
	}
	defer cancel()

	b := []byte("你好啊？？？？？")
	str, err := con.OpenStream()
	if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
		panic(nerr)
	}
	str.Write(b)
	r := make([]byte, 100)

	str.Read(r)

	assert.Equal(t, strings.TrimSpace(string(r[:6])), "已阅", "错误")
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
