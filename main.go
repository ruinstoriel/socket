package main

import (
	"fmt"
	"time"
)

func main() {
	d1 := make(chan int, 1)
	d2 := make(chan int, 1)
	go demo1(d1)
	go demo2(d2)
	select {
	case <-d1:
		fmt.Println("demo1")
	case <-d2:
		fmt.Println("demo2")
	}
	time.Sleep(4 * time.Second)
}

func demo1(c chan int) {
	time.Sleep(1 * time.Second)
	c <- 1
}

func demo2(c chan int) {
	time.Sleep(2 * time.Second)
	c <- 1
}
