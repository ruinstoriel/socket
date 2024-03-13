package tcp_package

import "testing"

func TestTcpServer(t *testing.T) {
	tcpServer("localhost:8080")
}
