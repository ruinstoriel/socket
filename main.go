package main

import (
	S "socket/tcp"
	"time"
)

func main() {
	S.TcpServer("localhost:8080")
	time.Sleep(120 * time.Second)
}
