package main

import (
	"fmt"
	"net"
)

func main() {
	TcpServer(":8443")
}

func TcpServer(addr string) {
	listener, _ := net.Listen("tcp", addr)
	accpet(listener)
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
			fmt.Printf("读取了%d, 内容是 %s , remote address %s \n", r, string(b[:r]), con.RemoteAddr().String())
		}
		// 能不能另开一个conn 回复    不能
		// newReplay(con.RemoteAddr().String())
		b = []byte("HTTP/1.1 407 Proxy Authentication Required\r\nProxy-Authenticate: Basic realm=\"Proxy\"\r\n\r\n")
		w, _ := con.Write(b)
		if w > 0 {
			fmt.Printf("发送了%d \n", w)
		}

	}
}

func newReplay(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	conn.Write([]byte("hello world\n"))
}
