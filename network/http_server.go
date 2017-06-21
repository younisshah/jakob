package network

import (
	"log"
	"net/http"
	"os"

	"github.com/younisshah/jakob/cluster"
)

/**
*  Created by Galileo on 21/6/17.
 */

const ADDR = "localhost:30000"

var logger = log.New(os.Stderr, "[jakob-http-server] ", log.LstdFlags)

func StartHTTPServer() error {
	http.HandleFunc("/join", func(writer http.ResponseWriter, req *http.Request) {
		cluster.Join(writer, req)
	})
	logger.Println("Starting HTTP server on", ADDR)
	return http.ListenAndServe(ADDR, nil)
}
