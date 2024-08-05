package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/proxy"
)

// Function to perform DNS query through SOCKS5 proxy
func dnsQueryThroughSocks5Proxy(proxyAddress, dnsServer, query string) ([]net.IP, error) {
	// Create a SOCKS5 dialer
	dialer, err := proxy.SOCKS5("tcp", proxyAddress, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
	}

	// Set up a custom resolver using the SOCKS5 dialer
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.Dial(network, dnsServer)
		},
	}

	// Perform the DNS lookup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ips, err := resolver.LookupIP(ctx, "ip4", query)
	if err != nil {
		return nil, fmt.Errorf("DNS lookup failed: %w", err)
	}

	return ips, nil
}

func main() {

	tr := &http.Transport{
		Proxy: func(res *http.Request) (*url.URL, error) {
			return url.Parse("http://127.0.0.1:8080")
		},
	}

	// Create client
	myClient := &http.Client{
		Transport: tr,
	}
	res, err := myClient.Get("https://www.google.com")
	if err != nil {
		panic(err)
	}
	bufio.NewReader(res.Body).WriteTo(os.Stdout)
	res.Body.Close()
}
