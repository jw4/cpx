package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	localHost   string = "localhost"
	localPort   int    = 8080
	shouldCache bool
)

func init() {
	flag.StringVar(&localHost, "host", localHost, "Bind Host")
	flag.IntVar(&localPort, "port", localPort, "Bind Port")
	flag.BoolVar(&shouldCache, "cache", false, "Cache results of duplicate requests")
}

func main() {
	flag.Parse()
	handler := proxyHandler
	if shouldCache {
		handler = cacheHandler(handler)
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", localPort),
		Handler: http.HandlerFunc(handler),
	}
	log.Fatal(server.ListenAndServe())
}
