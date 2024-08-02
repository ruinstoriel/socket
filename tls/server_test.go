package tls_package

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"gotest.tools/v3/assert"
)

func TestTlspServer(t *testing.T) {
	addr := "localhost:7890"
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
	ch := make(chan bool, 1)
	ch <- true
	wait(func() {

	}, ch)
}
func wait(cancel func(), end chan bool) {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	defer signal.Stop(sigs)
	<-sigs
	fmt.Println("收到退出信号")
	cancel()
	<-end
	fmt.Println("清理结束")

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
