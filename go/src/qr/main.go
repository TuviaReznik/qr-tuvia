package main

import (
	"net/http"

	"github.com/tuvirz/qr/go/src/qr/handler"
)

func main() {
	addr := "localhost:8888"
	mux := http.NewServeMux()
	mux.HandleFunc("/scan", handler.HandleScan)
	server := &http.Server{Addr: addr, Handler: mux}
	server.ListenAndServe()
}

// This implementation can detect and decode QR Code in an image.
