package tcpdemo

import "testing"

func TestTcpClient(t *testing.T) {
	tcpClient("localhost:8080")
}
func TestUdpClient(t *testing.T) {
	udpClient("localhost:8080")
}
