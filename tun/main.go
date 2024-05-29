package main

import (
	"socket/log"
	"sync"
)

var wg sync.WaitGroup // 创建 WaitGroup
func main() {
	log.SetLevel(log.DEBUG)
	wg.Add(1)
	NewTun()
	wg.Wait() // 等待所有goroutine完成
}
