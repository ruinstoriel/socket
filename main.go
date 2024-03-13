package main

import (
	S "socket/tcp_demo"
)

func main() {
	go udpServer()
	udpClient()

}

func udpServer() {
	S.UdpTlsServer("127.0.0.1:8080")
}
func udpClient() {
	S.UdpTlsClient("127.0.0.1:8080")
}
