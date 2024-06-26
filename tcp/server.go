package tcp_package

import (
	"fmt"
	"net"
)

func TcpServer(addr string) {
	listener, _ := net.Listen("tcp", addr)
	go accpet(listener)

}

func accpet(li net.Listener) {
	defer li.Close()
	for {
		conn, _ := li.Accept()
		go handleConnection(conn)
	}

}

func handleConnection(con net.Conn) {
	defer con.Close()
	// 如果没有读完就开始回写，会发生什么?
	for {
		b := make([]byte, 1024)
		r, err := con.Read(b)
		if err != nil {
			panic(err)
		}
		if r > 0 {
			fmt.Printf("读取了%d, 内容是%s \n", r, string(b[:r]))
		}

		w, _ := con.Write([]byte("bbbb\nbbbb"))
		if w > 0 {
			fmt.Printf("发送了%d \n", w)
		}

	}
}
