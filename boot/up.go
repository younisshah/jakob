package boot

import (
	"github.com/younisshah/jakob/health"
	"github.com/younisshah/jakob/httpd"
	"github.com/younisshah/jakob/machine/worker"
)

func Up() {

	// health check
	go func() {
		health.Check()
	}()

	// start Machinery worker
	go func() {
		err := worker.Worker()

	}()

	panic(httpd.StartHTTPServer())
}
