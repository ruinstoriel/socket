package main

import (
	"log"

	"github.com/miekg/dns"
)

func main() {

	dns.ListenAndServe(":3053", "udp4", dns.HandlerFunc(forwardRequst))
}

func forwardRequst(w dns.ResponseWriter, m *dns.Msg) {
	log.Println(m.Question)
	r, err := dns.Exchange(m, "114.114.114.114:53")

	if err != nil {
		log.Println("从 114 读取数据失败", err)
	}
	if m.Question[0].Qtype == 1 {
		log.Println(r.Answer[0].(*dns.A).A)
	}
	w.WriteMsg(r)

}
