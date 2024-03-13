package tcpdemo

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
)

func tcpClient(addr string) {
	con, _ := net.Dial("udp", addr)

	defer con.Close()

	b := []byte("你好啊")
	con.Write(b)
	r := make([]byte, 2)
	for {
		i, _ := con.Read(r)

		if i > 0 {
			fmt.Printf("client 收到:%s", string(r))
		}

	}
}
func UdpClient(addr string) {
	con, _ := net.Dial("udp", addr)

	defer con.Close()

	b := []byte("你好啊")
	con.Write(b)
	r := make([]byte, 100)
	for {
		i, _ := con.Read(r)

		if i > 0 {
			fmt.Printf("client 收到:%s", string(r))
		}

	}
}
func UdpTlsClient(addr string) {
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		log.Println("Error reading CA certificate:", err)
		return
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	conf := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS13, // 指定最低版本为TLS 1.2
		RootCAs:            caCertPool,
	}
	con, err := net.Dial("udp", addr)
	if err != nil {
		log.Println(err)
	}
	conn := tls.Client(con, conf)
	defer con.Close()
	conn.Handshake()
	b := []byte("你好啊")
	w, err := conn.Write(b)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(w)
	r := make([]byte, 100)
	for {
		i, _ := conn.Read(r)

		if i > 0 {
			fmt.Printf("client 收到:%s", string(r))
		}

	}
}
