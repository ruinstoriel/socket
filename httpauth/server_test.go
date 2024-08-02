package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestTcpServer(t *testing.T) {
	addr := "127.0.0.1:7890"
	go TcpServer(addr)

	ch := make(chan bool, 1)
	ch <- true
	wait(func() {

	}, ch)
}
func wait(cancel func(), end chan bool) {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	defer signal.Stop(sigs)
	<-sigs
	fmt.Println("收到退出信号")
	cancel()
	<-end
	fmt.Println("清理结束")

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
