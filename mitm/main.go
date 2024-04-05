package main

import (
	"fmt"
	"net/http"
)

func main() {
	proxy := New()
	fmt.Print(proxy)
	server := http.Server{
		Addr:    ":8080",
		Handler: proxy,
	}
	fmt.Println("Server listening on port 8080...")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
