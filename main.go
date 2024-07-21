package main

import (
	"context"
	"fmt"
	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-box/outbound"
	route2 "github.com/sagernet/sing-box/route"
	M "github.com/sagernet/sing/common/metadata"
	"net/netip"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	vlessMain()
}

func vlessMain() {
	err := vlessInbound().Start()
	if err != nil {
		panic(err)
	}
	fmt.Println("--------create outbound socket")
	out := vlessOutbound()
	dialContext, err := out.DialContext(context.Background(), "tcp", M.ParseSocksaddr("www.baidu.com:80"))
	if err != nil {
		panic(err)
	}
	_, err = dialContext.Write([]byte("GET / HTTP/1.1\nHost: www.baidu.com\n\n"))
	if err != nil {
		panic(err)
	}
	bs := make([]byte, 1024)
	read, err := dialContext.Read(bs)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs[:read]))
	wait()
}
func wait() {
	exitChan := make(chan struct{})

	// 捕捉系统信号
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		<-sigs
		fmt.Println("收到退出信号")
		close(exitChan)
	}()

	// 阻塞主线程直到收到退出信号
	<-exitChan
	fmt.Println("程序退出")
}
func vlessInbound() *inbound.VLESS {
	route, _ := route2.NewRouter(context.Background(), log.NewNOPFactory(), option.RouteOptions{}, option.DNSOptions{}, option.NTPOptions{}, []option.Inbound{}, nil)

	vlessIn, err := inbound.NewVLESS(context.Background(), route, log.StdLogger(), "vless-in",
		option.VLESSInboundOptions{
			ListenOptions: option.ListenOptions{
				Listen:     option.NewListenAddress(netip.AddrFrom4([4]byte{127, 0, 0, 1})),
				ListenPort: 8888,
			},
			Users: []option.VLESSUser{
				{
					Name: "user",
					UUID: "775abd9a-5f06-4548-a0ca-3b887e36ce1a",
					Flow: "xtls-rprx-vision",
				},
			},
			InboundTLSOptionsContainer: option.InboundTLSOptionsContainer{
				TLS: &option.InboundTLSOptions{
					Enabled:    true,
					ServerName: "www.python.org",
					Reality: &option.InboundRealityOptions{
						Enabled: true,
						Handshake: option.InboundRealityHandshakeOptions{
							ServerOptions: option.ServerOptions{
								Server:     "www.python.org",
								ServerPort: 443,
							},
						},
						PrivateKey: "uO8nRwHcs8dPIVtLDVBVmA3wU0j1kpBRvQt5ZbW36Xc",
						ShortID:    []string{"6e1a647f0311592a"},
					},
				},
			},
		})
	if err != nil {
		panic(err)
	}
	return vlessIn
}

func vlessOutbound() *outbound.VLESS {
	vlessOut, err := outbound.NewVLESS(context.Background(), (adapter.Router)(nil), log.StdLogger(), "aa", option.VLESSOutboundOptions{
		ServerOptions: option.ServerOptions{
			Server:     "127.0.0.1",
			ServerPort: 8888,
		},
		UUID: "6f7c4066-cc28-48b1-a918-e1faeabeb6a1",
		Flow: "xtls-rprx-vision",

		OutboundTLSOptionsContainer: option.OutboundTLSOptionsContainer{
			TLS: &option.OutboundTLSOptions{
				Enabled:    true,
				ServerName: "www.python.org",
				UTLS: &option.OutboundUTLSOptions{
					Enabled:     true,
					Fingerprint: "chrome",
				},
				Reality: &option.OutboundRealityOptions{
					Enabled:   true,
					PublicKey: "AcIz8PMSh_FZrJvCuuUNC4QqF248ximK9MQF7j5ICyg",
					ShortID:   "6e1a647f0311592a",
				},
			},
		},
		Multiplex: &option.OutboundMultiplexOptions{
			Enabled: false,
		},
	})
	if err != nil {
		panic(err)
	}
	return vlessOut
}
