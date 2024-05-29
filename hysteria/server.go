package hysteria

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"socket/log"

	"github.com/metacubex/sing-quic/hysteria2"
	"github.com/sagernet/sing/common/buf"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/common/network"
)

type MyHandler struct {
}

func (h MyHandler) NewConnection(ctx context.Context, con net.Conn, metadata M.Metadata) error {
	defer con.Close()
	// 如果没有读完就开始回写，会发生什么?
	for {
		b := make([]byte, 1024)
		r, err := con.Read(b)
		if err != nil {
			panic(err)
		}
		if r > 0 {
			log.Infoln("tcp读取了%d, 内容是%s \n", r, string(b[:r]))
		}

		w, _ := con.Write([]byte("已阅"))
		if w > 0 {
			log.Infoln("发送了%d \n", w)
		}

	}

}

func (h MyHandler) NewPacketConnection(ctx context.Context, con network.PacketConn, metadata M.Metadata) error {
	defer con.Close()
	for {

		b := buf.New()
		addr, err := con.ReadPacket(b)
		if err != nil {
			log.Errorln("%s", err)
		}
		// 如果没有接收完，剩余的数据无法读取，也无法回写？
		if b.Len() > 0 {
			log.Infoln("udp读取了%d, 内容是%s \n", b.Len(), string(b.Bytes()[:b.Len()]))
		}
		w := buf.New()
		w.Write([]byte("已阅"))
		err = con.WritePacket(w, addr)
		if err == nil {
			log.Infoln("发送了%d \n", w.Len())
		}

	}
}

func HysteriaServer(addr string) {
	options := hysteria2.ServiceOptions{
		Context:               context.Background(),
		Logger:                log.SingLogger,
		BrutalDebug:           true,
		SendBPS:               1024,
		ReceiveBPS:            1024,
		IgnoreClientBandwidth: false,
		SalamanderPassword:    "a",
		TLSConfig:             config(),
		UDPDisabled:           false,
		Handler:               MyHandler{},
	}

	s, err := hysteria2.NewService[string](options)
	if err != nil {
		fmt.Println(err)
	}
	udpAddr, _ := net.ResolveUDPAddr("udp4", addr)
	udpConn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		panic(err)
	}
	s.UpdateUsers([]string{"yao", "a"}, []string{"b", "x"})
	s.Start(udpConn)
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
