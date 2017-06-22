package boot

import (
	"log"
	"os"

	"github.com/younisshah/jakob/httpd"
	"github.com/younisshah/jakob/machine/worker"
	"github.com/younisshah/jakob/health"
)

/**
*  Created by Galileo on 22/6/17.
 */

var logger = log.New(os.Stderr, "[jakob-boot-up] ", log.LstdFlags)

func Up() {

	// health check
	go func() {
		health.Check()
	}()

	// start Machinery worker
	go func() {
		worker.Worker()
	}()

	panic(httpd.StartHTTPServer())
}
