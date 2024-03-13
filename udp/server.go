package udp_package

import (
	"fmt"
	"net"
)

func UdpServer(addr string) {
	udpAddr, _ := net.ResolveUDPAddr("udp", addr)
	conn, _ := net.ListenUDP("udp", udpAddr)
	go udpHandle(conn)
}

func udpHandle(con *net.UDPConn) {
	defer con.Close()
	for {
		b := make([]byte, 1024)
		r, addr, err := con.ReadFrom(b)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(r)
		// 如果没有接收完，剩余的数据无法读取，也无法回写？
		if r > 0 {
			fmt.Printf("读取了%d, 内容是%s \n", r, string(b))
		}

		w, _ := con.WriteTo([]byte("已阅"), addr)
		if w > 0 {
			fmt.Printf("发送了%d \n", w)
		}

	}

}
