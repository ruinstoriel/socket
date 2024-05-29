package hysteria

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"socket/log"
	"strings"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"github.com/metacubex/sing-quic/hysteria2"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/common/network"
)

func TestQuicServer(t *testing.T) {
	addr1 := "192.168.1.157:8888"
	addr2 := "192.168.1.157:9999"
	arr := []string{addr1, addr2}
	log.Infoln("%v", arr)
	gen := func() func(ctx context.Context) (*net.UDPAddr, error) {
		i := 0
		return func(ctx context.Context) (*net.UDPAddr, error) {
			log.Infoln("%s,index: %d", time.Now(), i)
			udpAddr, _ := net.ResolveUDPAddr("udp4", arr[i])
			i = i + 1
			i = i % 2
			return udpAddr, nil
		}
	}
	log.SetLevel(log.DEBUG)
	options := hysteria2.ClientOptions{

		Context: context.Background(),
		Logger:  log.SingLogger,
		Dialer: &network.DefaultDialer{
			net.Dialer{},
			net.ListenConfig{},
		},
		SalamanderPassword: "a",
		TLSConfig:          clientConfig(),
		ServerAddress:      gen(),
		HopInterval:        time.Duration(10) * time.Second,
		Password:           "x",
	}
	client, err := hysteria2.NewClient(options)

	if err != nil {
		fmt.Println(err)
	}
	udpAddr, _ := net.ResolveUDPAddr("udp4", "www.baidu.com:80")
	con, err := client.DialConn(context.Background(), M.SocksaddrFromNet(udpAddr))
	if err != nil {
		panic(err)
	}

	defer con.Close()
	for {
		time.Sleep(1 * time.Second)
		b := []byte("你好啊")
		con.Write(b)
		r := make([]byte, 1024)

		rb, err := con.Read(r)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, strings.TrimSpace(string(r[:rb])), "已阅", "错误")
	}

}
func clientConfig() *tls.Config {
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
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
