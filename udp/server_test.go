package udp_package

import (
	"net"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestUdpServer(t *testing.T) {
	addr := ":8080"
	UdpServer(addr)
	con, _ := net.Dial("udp", addr)

	defer con.Close()

	b := []byte("你好啊")
	con.Write(b)
	r := make([]byte, 100)

	con.Read(r)

	assert.Equal(t, strings.TrimSpace(string(r[:6])), "已阅", "错误")

}
