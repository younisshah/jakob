package httpd

import (
	"log"
	"net/http"
	"os"

	"github.com/younisshah/jakob/cluster"
)

const ADDR = "localhost:30000"

var logger = log.New(os.Stderr, "[jakob-http-server] ", log.LstdFlags)

func StartHTTPServer() error {
	http.HandleFunc("/join", func(writer http.ResponseWriter, req *http.Request) {
		cluster.Join(writer, req)
	})
	http.HandleFunc("/init", func(writer http.ResponseWriter, req *http.Request) {
		cluster.Init(writer, req)
	})
	logger.Println("Running HTTP server on", ADDR)
	return http.ListenAndServe(ADDR, nil)
}
