package main

import (
	"crypto/tls"
	"fmt"
	"io"
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

func (p Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodConnect {
		// 非隧道代理
		httpHandle(rw, req)
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
