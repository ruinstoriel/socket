package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Proxy struct {
}

// 定义函数类型
type WriteFunc func(p []byte) (n int, err error)

// 为函数类型实现 Write 方法
func (f WriteFunc) Write(p []byte) (n int, err error) {
	return f(p)
}

type ReadFunc func(p []byte) (n int, err error)

func (f ReadFunc) Read(p []byte) (n int, err error) {
	return f(p)
}
func (p Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Method---------------" + req.Method)
	if req.Method == http.MethodConnect {
		if hijacker, ok := rw.(http.Hijacker); ok {

			conn, _, err := hijacker.Hijack()
			if err != nil {
				fmt.Printf("hijacker错误: %s", err)
			}
			conn.Write(tunnelEstablishedResponseLine)
			addr := req.URL.Host
			domain := strings.Split(addr, ":")[0]
			fmt.Println(domain)
			u := url.URL{Scheme: "wss", Host: "socket.weiliangcn.top", Path: "/"}
			//u := url.URL{Scheme: "ws", Host: "127.0.0.1:8787", Path: "/"}

			//earlyData := make([]byte, 4096)
			//i, err := conn.Read(earlyData)
			//fmt.Println("first package: ", i)
			if err != nil {
				panic(err)
			}

			fmt.Printf("connecting to %s \n", domain)
			h := http.Header{}
			hh := base64.StdEncoding.EncodeToString([]byte(domain))
			//earlyDataBase := base64.StdEncoding.EncodeToString(earlyData[:i])
			//h.Add("sec-websocket-protocol", earlyDataBase)

			// 0x01: IPv4 address
			// 0x03: Domain name
			// 0x04: IPv6 address
			h.Add("base", hh)
			h.Add("atype", "0x03")
			h.Add("token", "85aad405-ec1b-4332-b405-de08b8d53629")

			c, _, err := websocket.DefaultDialer.Dial(u.String(), h)
			if err != nil {
				panic(err)
			}
			defer func(c *websocket.Conn) {
				err = c.Close()
				if err != nil {
					fmt.Printf("close websocket err: %s", err)
				}
			}(c)

			copyErrChan := make(chan error, 2)
			go func() {
				myWriterFunc := WriteFunc(func(p []byte) (n int, err error) {
					return len(p), c.WriteMessage(websocket.BinaryMessage, p)
				})
				_, copyErr := io.Copy(myWriterFunc, conn)
				copyErrChan <- copyErr
			}()
			go func() {
				myReadFun := ReadFunc(func(p []byte) (n int, err error) {
					_, body, err := c.ReadMessage()
					if err != nil {
						return 0, err
					}
					copy(p, body)
					return len(body), err
				})

				_, copyErr := io.Copy(conn, myReadFun)
				copyErrChan <- copyErr
			}()
			err = <-copyErrChan
			if err != nil {
				fmt.Println(err)
			}

		} else {
			httpHandle(rw, req)

		}
	}
}
