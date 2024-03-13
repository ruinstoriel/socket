package tcp_package

import (
	"net"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestTcpServer(t *testing.T) {
	addr := "localhost:8080"
	tcpServer(addr)
	con, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer con.Close()
	b := []byte("你好啊")
	con.Write(b)
	r := make([]byte, 1024)

	con.Read(r)
	assert.Equal(t, strings.TrimSpace(string(r[:6])), "已阅", "错误")

}
