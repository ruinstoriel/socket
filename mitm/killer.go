package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Killer struct {
}

var CCS = []byte{20, 3, 3, 0, 1, 1}

func (p Killer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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
			u := url.URL{Scheme: "wss", Host: "socket.weiliangcn.top", Path: "/"}

			fmt.Printf("connecting to %s \n", domain)
			h := http.Header{}
			hh := base64.StdEncoding.EncodeToString([]byte(domain))
			h.Add("sec-websocket-protocol", hh)
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
			var mutex sync.Mutex

			uploading := false
			upCount := 0

			downloading := false
			downCount := 0

			copyErrChan := make(chan error, 2)
			go func() {
				//  客户端链接到远程复制
				myWriterFunc := WriteFunc(func(p []byte) (n int, err error) {
					return len(p), c.WriteMessage(websocket.BinaryMessage, p)
				})

				buf := make([]byte, 8192)
				for {
					// 从客户端的流读取
					n, err := conn.Read(buf)
					if err != nil {
						return
					}
					mutex.Lock()
					fmt.Println("upload---", upCount, n, buf)
					if upCount == 0 && n >= 6 && bytes.Equal(buf[:6], CCS) {
						fmt.Printf("%s 开始 upload  \n", req.URL.Host)
						uploading = true
					}
					if uploading {
						upCount += n
					}
					if downloading {
						fmt.Printf("%s 已经开始 download \n", req.URL.Host)
						downloading = false
						fmt.Printf("%v\tupCount %v\tdownCount %v\n", req.URL.Host, upCount, downCount)
						if upCount >= 650 && upCount <= 750 &&
							((downCount >= 170 && downCount <= 180) || (downCount >= 3000 && downCount <= 7500)) {
							fmt.Printf("%v is Trojan\n", req.URL.Host)
						}
					}
					mutex.Unlock()
					_, err = myWriterFunc.Write(buf[:n])
					if err != nil {
						copyErrChan <- err
					}
					if !downloading && downCount != 0 {
						_, copyErr := io.CopyBuffer(myWriterFunc, conn, buf)
						copyErrChan <- copyErr
					}
				}
			}()
			go func() {
				// 远程链接到客户端链接复制
				myReadFun := ReadFunc(func(p []byte) (n int, err error) {
					_, body, err := c.ReadMessage()
					if err != nil {
						return 0, err
					}
					copy(p, body)
					return len(body), err
				})
				buf := make([]byte, 8192)

				for {
					n, err := myReadFun.Read(buf)
					if err != nil {
						return
					}
					mutex.Lock()

					if uploading {
						fmt.Println("down---", downCount, n, buf)
						fmt.Printf("%s 已经upload \n", req.URL.Host)
						uploading = false
						downloading = true
					}
					if downloading {
						downCount += n
					}
					mutex.Unlock()
					_, err = conn.Write(buf[:n])
					if err != nil {
						return
					}
					if !downloading && downCount != 0 {
						_, copyErr := io.CopyBuffer(conn, myReadFun, buf)
						copyErrChan <- copyErr
					}
				}
			}()
			err = <-copyErrChan
			fmt.Println(err)
		} else {
			httpHandle(rw, req)

		}
	}
}
