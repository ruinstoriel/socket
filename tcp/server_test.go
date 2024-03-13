package tcpdemo

import "testing"

func TestTcpServer(t *testing.T) {
	tcpServer("localhost:8080")
}
func TestUdpServer(t *testing.T) {
	UdpServer(":8080")
}
