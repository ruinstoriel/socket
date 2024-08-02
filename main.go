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
	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/service"
	"github.com/sagernet/sing/service/pause"
	"net"
	"net/netip"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	addr := &net.UDPAddr{
		IP:   net.ParseIP("2409:8a62:164:1bb0:597f:821f:74fe:f4f0"),
		Port: 60949,
	}
	fmt.Println(addr)
	zoneID, err := strconv.ParseUint(addr.Zone, 10, 32)
	if err != nil {
		panic(err)
	}

	fmt.Println(uint32(zoneID))
}

func http2Direct() {
	lf := logFactory()
	ctx, cancel := context.WithCancel(context.Background())
	ctx = service.ContextWithDefaultRegistry(ctx)
	ctx = pause.WithDefaultManager(ctx)
	router := routeCreate(ctx, lf)
	in := httpInbound(ctx, router)
	err := in.Start()
	if err != nil {
		panic(err)
	}
	router.Initialize([]adapter.Inbound{in}, nil, func() adapter.Outbound {
		out, oErr := outbound.New(ctx, router, lf.NewLogger("outbound/direct"), "direct", option.Outbound{Type: "direct", Tag: "default"})
		common.Must(oErr)
		return out
	})
	end := make(chan bool, 1)
	go func() {
		select {
		case <-ctx.Done():
			in.Close()
			end <- true
		}
	}()
	wait(cancel, end)
}

func httpInbound(ctx context.Context, router adapter.Router) *inbound.HTTP {

	httpIn, err := inbound.NewHTTP(ctx, router, log.StdLogger(), "http-in",
		option.HTTPMixedInboundOptions{
			ListenOptions: option.ListenOptions{
				Listen:     option.NewListenAddress(netip.AddrFrom4([4]byte{127, 0, 0, 1})),
				ListenPort: 8888,
			},

			SetSystemProxy: true,
		})
	if err != nil {
		panic(err)
	}
	return httpIn
}

func vlessMain() {
	lf := logFactory()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	router := routeCreate(ctx, lf)
	in := vlessInbound(router)
	err := in.Start()
	if err != nil {
		panic(err)
	}
	fmt.Println("--------create outbound socket")
	out := vlessOutbound(router)
	router.Initialize([]adapter.Inbound{in}, nil, func() adapter.Outbound {
		return out
	})
	router.Start()
}

func logFactory() log.Factory {
	fa, err := log.New(log.Options{
		Context: context.Background(),
		Options: common.PtrValueOrDefault(&option.LogOptions{
			Disabled: false,
			Level:    "debug",
		}),
		Observable:     false,
		DefaultWriter:  os.Stdout,
		BaseTime:       time.Now(),
		PlatformWriter: nil,
	})
	if err != nil {
		panic(err)
	}
	return fa
}

func routeCreate(ctx context.Context, lf log.Factory) *route2.Router {

	router, err := route2.NewRouter(
		ctx,
		lf,
		common.PtrValueOrDefault(&option.RouteOptions{}),
		common.PtrValueOrDefault(&option.DNSOptions{}),
		common.PtrValueOrDefault(&option.NTPOptions{}),
		nil,
		nil,
	)
	if err != nil {
		panic(err)
	}
	return router
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
func vlessInbound(router adapter.Router) *inbound.VLESS {

	vlessIn, err := inbound.NewVLESS(context.Background(), router, log.StdLogger(), "vless-in",
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

func vlessOutbound(router adapter.Router) *outbound.VLESS {
	vlessOut, err := outbound.NewVLESS(context.Background(), router, log.StdLogger(), "aa", option.VLESSOutboundOptions{
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
