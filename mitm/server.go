package main

import (
	"net/http"
)

type Proxy struct {
}

func (p Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("are you ok?"))
}

func New() Proxy {
	return Proxy{}
}
