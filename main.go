package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

func ProxyHandler(w http.ResponseWriter, req *http.Request) {
	myDial := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	myDialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		network = "tcp6"
		return myDial.DialContext(ctx, network, addr)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           myDialContext,
			MaxIdleConns:          20,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	resp, err := client.Get("https://repo.continuum.io" + req.URL.Path)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	fmt.Fprintln(w, string(body))
}

func main() {
	http.HandleFunc("/", ProxyHandler)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
