package task_server

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/younisshah/jakob/core/tasks"
)

/**
*  Created by Galileo on 20/6/17.
 */

// StartTaskServer gets a new machinery Server
func StartTaskServer() (server *machinery.Server, err error) {
	server, err = machinery.NewServer(LoadInMemConfig())
	if err != nil {
		return
	}
	err = server.RegisterTask("Sync", tasks.Sync)
	return
}

// TODO - get path from kingpin
func loadConfig() *config.Config {
	return config.NewFromYaml("./config.yml", true, true)
}
