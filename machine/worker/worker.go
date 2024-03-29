package worker

import (
	"log"
	"os"

	"github.com/younisshah/jakob/machine/task-server"
)

const CONCURRENCY = 2

var logger = log.New(os.Stderr, "[jakob-worker] ", log.LstdFlags)

func Worker() error {

	server, err := task_server.StartTaskServer()

	if err != nil {
		logger.Println("couldn't start Jakob task server", err)
		return err
	}

	w := server.NewWorker("jakob_worker", CONCURRENCY)

	return w.Launch()
}
