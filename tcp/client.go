package tcp_package

import (
	"fmt"
	"net"
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
