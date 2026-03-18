package quic_demo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/ruinstoriel/quic-go"
	"github.com/ruinstoriel/quic-go/quicvarint"
	"io"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
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
	for i := 0; i < 1; i++ {
		count++
		fmt.Printf("第 %d 个连接 \n", count)

		str, err2 := con.OpenStream()
		// err2 = WriteTCPRequest(str, "www.baidu.com")
		err2 = plaintextTransmission(str)
		if err2 != nil {
			panic(err2)
		}

		r := make([]byte, 100)

		_, err5 := str.Read(r)
		//fmt.Println(string(r[:ss]))

		if err5 != nil {
			// panic(err5)
		}

		fmt.Printf("client %d CancelRead() \n", str.StreamID())
		str.CancelRead(0)
		time.Sleep(10 * time.Millisecond)
		err = str.Close()
		fmt.Printf("client %v  \n", err)

	}

}

func plaintextTransmission(s *quic.Stream) error {
	buf := make([]byte, 1024)
	addrLen := 11
	i := varintPut(buf, FrameTypeTCPRequest)
	fmt.Println("addrLen: ", addrLen)
	fmt.Printf("i: %d\n", i)
	varintPut(buf[i:], uint64(addrLen))
	s.Write(buf)
	return nil
}

func WriteTCPRequest(w io.Writer, addr string) error {

	addrLen := len(addr)
	log.Printf("addrLen: %s", strconv.Itoa(addrLen))
	sz := int(quicvarint.Len(FrameTypeTCPRequest)) +
		int(quicvarint.Len(uint64(addrLen))) + addrLen
	buf := make([]byte, sz)
	i := varintPut(buf, FrameTypeTCPRequest)

	i += varintPut(buf[i:], uint64(addrLen))

	i += copy(buf[i:], addr)
	_, err := w.Write(buf)
	return err
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
		InsecureSkipVerify: true,
		ServerName:         "127.0.0.1",
		MinVersion:         tls.VersionTLS13, // 指定最低版本为TLS 1.2
		RootCAs:            caCertPool,
	}
	return conf
}
