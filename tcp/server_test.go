package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestTcpServer(t *testing.T) {
	addr := "1.2.3.4:4433"
	//TcpServer(addr)
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

func TestTcpServerHttp(t *testing.T) {
	addr := ":8080"
	TcpServer(addr)
	resp, err := http.Get("http://127.0.0.1" + addr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	assert.Equal(t, string(body), "已阅", "错误")
}
