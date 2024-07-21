package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"net"
	"unsafe"
)

func getLinger(conn *net.TCPConn) (*windows.Linger, error) {
	rawConn, err := conn.SyscallConn()
	if err != nil {
		return nil, err
	}

	var linger windows.Linger
	var size = int32(unsafe.Sizeof(linger))
	rawConn.Control(func(fd uintptr) {
		errno := windows.Getsockopt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_LINGER, (*byte)(unsafe.Pointer(&linger)), &size)
		if errno != nil {
			err = errno
		}
	})

	if err != nil {
		return nil, err
	}
	return &linger, nil
}

func setLinger(conn *net.TCPConn, onoff int32, linger int32) error {
	rawConn, err := conn.SyscallConn()
	if err != nil {
		return err
	}

	err = rawConn.Control(func(fd uintptr) {
		l := windows.Linger{
			Onoff:  onoff,
			Linger: linger,
		}
		err = windows.SetsockoptLinger(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_LINGER, &l)
		if err != nil {
			fmt.Println(err)
		}
	})

	return err
}

func main() {
	// 建立一个TCP连接（这里是示例，请使用适当的地址和端口）
	conn, err := net.Dial("tcp", "www.baidu.com:80")

	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	tcpConn := conn.(*net.TCPConn)

	// 设置 SO_LINGER 选项
	err = setLinger(tcpConn, 1, 0)
	if err != nil {
		fmt.Println("Error setting SO_LINGER:", err)
		return
	}

	// 获取 SO_LINGER 选项值
	linger, err := getLinger(tcpConn)
	if err != nil {
		fmt.Println("Error getting SO_LINGER:", err)
		return
	}

	fmt.Printf("SO_LINGER: onoff=%d, linger=%d\n", linger.Onoff, linger.Linger)
}
