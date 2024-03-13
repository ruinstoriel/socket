package tcp_package

import (
	"fmt"
	"net"
)

func tcpServer(addr string) {
	listener, _ := net.Listen("tcp", addr)
	defer listener.Close()
	for {
		conn, _ := listener.Accept()
		go handle(conn)
	}
}

func handle(con net.Conn) {
	defer con.Close()
	// 如果没有读完就开始回写，会发生什么?
	for {
		b := make([]byte, 1024)
		r, err := con.Read(b)
		if err != nil {
			panic(err)
		}
		if r > 0 {
			fmt.Printf("读取了%d, 内容是%s \n", r, string(b))
		}

		w, _ := con.Write([]byte("已阅"))
		if w > 0 {
			fmt.Printf("发送了%d \n", w)
		}

	}
}
