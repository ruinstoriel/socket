package tcp_package

import (
	"crypto/tls"
	"fmt"
	"log"
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

func UdpServer(addr string) {
	udpAddr, _ := net.ResolveUDPAddr("udp", addr)
	conn, _ := net.ListenUDP("udp", udpAddr)
	udpHandle(conn)
}
func UdpTlsServer(addr string) {
	udpAddr, _ := net.ResolveUDPAddr("udp", addr)
	conn, _ := net.ListenUDP("udp", udpAddr)
	tlsUdpConn(conn)
}
func tlsUdpConn(con net.Conn) {
	cer, err := tls.LoadX509KeyPair("domain.crt", "domain.key")
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cer},
		MinVersion:   tls.VersionTLS13, // 指定最低版本为TLS 1.2
	}
	server := tls.Server(con, config)
	tlsHandle(server)
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
func tlsHandle(con *tls.Conn) {
	defer con.Close()
	con.Handshake()
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
