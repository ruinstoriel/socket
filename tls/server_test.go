package tls_package

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestTlspServer(t *testing.T) {
	addr := "localhost:8080"
	TlsServer(addr)
	con, err := net.Dial("tcp", addr)
	con = tls.Client(con, clientConfig())
	if err != nil {
		t.Fatal(err)
	}

	defer con.Close()
	b := []byte("你好啊sdfsdfsdfsdfsdfsdfsdfsdf")
	con.Write(b)
	r := make([]byte, 1024)

	i, _ := con.Read(r)
	assert.Equal(t, string(r[:i]), "已阅", "错误")
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
