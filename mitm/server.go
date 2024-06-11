package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

var hopHeaders = []string{
	"Proxy-Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailer",
	"Transfer-Encoding",
}

// CopyHeader 浅拷贝Header
func CopyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

var badGateway = []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", http.StatusBadGateway, http.StatusText(http.StatusBadGateway)))
var tunnelEstablishedResponseLine = []byte("HTTP/1.1 200 Connection established\r\n\r\n")

type ConnBuffer struct {
	net.Conn
	buf *bufio.ReadWriter
}
type Proxy struct {
}

func httpHandle(rw http.ResponseWriter, req *http.Request) {
	newReq := new(http.Request)
	*newReq = *req
	newHeader := http.Header{}
	CloneHeader(newReq.Header, newHeader)
	newReq.Header = newHeader
	fmt.Println(newReq.Header)
	for _, item := range hopHeaders {
		if newReq.Header.Get(item) != "" {
			newReq.Header.Del(item)
		}
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConns:          100,
		MaxConnsPerHost:       10,
		IdleConnTimeout:       10 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	resp, err := tr.RoundTrip(newReq)
	if err == nil {
		for _, h := range hopHeaders {
			resp.Header.Del(h)
		}
	}
	if err != nil {
		fmt.Errorf("%s - HTTP请求错误: %s", req.URL, err)
		rw.WriteHeader(http.StatusBadGateway)
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	CopyHeader(rw.Header(), resp.Header)
	rw.WriteHeader(resp.StatusCode)
	buf := make([]byte, 32*1024)
	_, _ = io.CopyBuffer(rw, resp.Body, buf)
}

var cache = &CacheImp{}

func (p Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Method---------------" + req.Method)
	if req.Method == http.MethodConnect {
		if hijacker, ok := rw.(http.Hijacker); ok {

			conn, buf, err := hijacker.Hijack()
			if err != nil {
				fmt.Printf("hijacker错误: %s", err)
			}
			if buf == nil {
				buf = bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
			}
			// conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nContent-Type: text/html\r\n<html><head><title>Hello World</title></head><h1>Hello World</h1></html>"))
			conn.Write(tunnelEstablishedResponseLine)
			clientConn := &ConnBuffer{
				Conn: conn,
				buf:  buf,
			}
			defer func() {
				_ = clientConn.Close()
			}()

			cert := NewCertificate(cache, true)
			fmt.Println("host-------", req.URL.Host)
			tlsConfig, err := cert.GenerateTlsConfig(req.URL.Host)
			if err != nil {
				fmt.Printf("%s - HTTPS解密, 生成证书失败: %s", req.URL.Host, err)
				return
			}
			tlsClientConn := tls.Server(clientConn, tlsConfig)

			defer func() {
				_ = tlsClientConn.Close()
			}()
			if err := tlsClientConn.Handshake(); err != nil {
				fmt.Printf("%s - HTTPS解密, 握手失败: %s \n", req.URL.Host, err)
				return
			}
			buff := bufio.NewReader(tlsClientConn)
			tlsReq, err := http.ReadRequest(buff)
			if err != nil {
				if err != io.EOF {
					fmt.Printf("%s - HTTPS解密, 读取客户端请求失败: %s", req.URL.Host, err)
				}
				return
			}

			tr := &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				MaxIdleConns:          100,
				MaxConnsPerHost:       10,
				IdleConnTimeout:       10 * time.Second,
				TLSHandshakeTimeout:   5 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			}
			resp, err := tr.RoundTrip(tlsReq)
			if err == nil {
				for _, h := range hopHeaders {
					resp.Header.Del(h)
				}
			}
			if err != nil {
				fmt.Errorf("%s - HTTPS解密, 请求错误: %s", tlsReq.URL, err)
				_, _ = tlsClientConn.Write(badGateway)
				return
			}
			err = resp.Write(tlsClientConn)
			if err != nil {
				fmt.Errorf("%s - HTTPS解密, response写入客户端失败, %s", req.URL, err)
			}
			_ = resp.Body.Close()

		} else {
			httpHandle(rw, req)

		}
	}
}

func New() Proxy {
	return Proxy{}
}

// CloneHeader 深拷贝Header
func CloneHeader(h http.Header, h2 http.Header) {
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
}
