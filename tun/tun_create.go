package main

import (
	"context"
	"net"
	"net/netip"
	"socket/log"
	"time"

	tun "github.com/metacubex/sing-tun"
	"github.com/sagernet/sing/common/buf"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/common/network"
)

type MyHandler struct {
}

func (h MyHandler) NewConnection(ctx context.Context, con net.Conn, metadata M.Metadata) error {
	defer con.Close()
	log.Infoln("tcp: %s", metadata.Destination.AddrString())
	return nil

}

func (h MyHandler) NewPacketConnection(ctx context.Context, con network.PacketConn, metadata M.Metadata) error {
	defer con.Close()
	log.Infoln("udp: %s", metadata.Destination.AddrString())
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
		httpResponse := "HTTP/1.1 200 OK\r\nContent-Length: 11\r\nContent-Type: text/plain\r\n\r\nHello World"
		w.Write([]byte(httpResponse))
		err = con.WritePacket(w, addr)
		if err == nil {
			log.Infoln("发送了%d \n", w.Len())
		}

	}
}
func (h MyHandler) NewError(ctx context.Context, err error) {
	log.Errorln("a", err)
}

func NewTun() {
	// 注册路由
	networkUpdateMonitor, err := tun.NewNetworkUpdateMonitor(log.SingLogger)
	if err != nil {
		panic(err)
	}
	err = networkUpdateMonitor.Start()
	if err != nil {
		panic(err)
	}
	// 观测默认出口的变化
	defaultInterfaceMonitor, err := tun.NewDefaultInterfaceMonitor(networkUpdateMonitor, log.SingLogger, tun.DefaultInterfaceMonitorOptions{OverrideAndroidVPN: true})
	if err != nil {
		panic(err)
	}
	err = defaultInterfaceMonitor.Start()
	if err != nil {
		panic(err)

	}
	pre, _ := netip.ParsePrefix("198.18.0.0/16")
	o := tun.Options{
		Name:             "meta1",
		Inet4Address:     []netip.Prefix{pre},
		MTU:              9000,
		GSO:              false,
		AutoRoute:        true,
		StrictRoute:      true,
		InterfaceMonitor: defaultInterfaceMonitor,
		TableIndex:       199,
	}
	tunIf, err := tun.New(o)
	if err != nil {
		panic(err)
	}

	stackOptions := tun.StackOptions{
		Context:                context.Background(),
		Tun:                    tunIf,
		TunOptions:             o,
		EndpointIndependentNat: false,
		UDPTimeout:             int64((5 * time.Minute).Seconds()),
		Handler:                MyHandler{},
		Logger:                 log.SingLogger,
	}
	tunStack, err := tun.NewStack("system", stackOptions)
	if err != nil {
		panic(err)
	}

	err = tunStack.Start()
	if err != nil {
		panic(err)

	}
}
