package main

import (
	"github.com/younisshah/jakob/httpd"
	"github.com/younisshah/jakob/machine/worker"
)

/**
*  Created by Galileo on 17/6/17.
 */

func main() {

	// start Machinery worker
	go func() {
		worker.Worker()
	}()

	panic(httpd.StartHTTPServer())
}
