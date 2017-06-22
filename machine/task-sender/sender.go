package task_sender

import (
	"log"
	"os"
	"time"

	"bytes"
	"strings"

	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/younisshah/jakob/jfs"
	"github.com/younisshah/jakob/machine/task-server"
)

/**
*  Created by Galileo on 20/6/17.
 */

var logger = log.New(os.Stderr, "[jakob-sender] ", log.LstdFlags)

func Send(cmdName string, args interface{}) {

	server, err := task_server.StartTaskServer()
	if err != nil {
		logger.Println("couldn't start machinery server", err)
		return
	}
	logger.Println(server.GetBackend())

	jyml := jfs.NewJYaml()
	jyml.Type = jfs.GETTER
	peers, err := jyml.Peers()

	if err != nil {
		logger.Println("couldn't get getter peers", err)
		return
	}

	logger.Println("peers -", peers)

	jTasks := make([]*tasks.Signature, len(peers))

	for i := range peers {
		jTasks[i] = getTask(peers[i], cmdName, stringify(args))
	}

	group := tasks.NewGroup(jTasks...)

	asyncResults, err := server.SendGroup(group)
	if err != nil {
		logger.Println("server send an error while executing task group", err)
		return
	}
	for _, r := range asyncResults {
		results, err := r.Get(time.Duration(time.Millisecond * 5))
		if err != nil {
			logger.Printf("Getting task result failed with error: %s\n", err.Error())
		}
		logger.Println(results)
	}
}

func getTask(peer, cmdName string, args string) *tasks.Signature {
	return &tasks.Signature{
		Name: "Sync",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: peer,
			},
			{
				Type:  "string",
				Value: cmdName,
			},
			{
				Type:  "string",
				Value: args,
			},
		},
	}
}

func stringify(args interface{}) string {
	a := args.([]interface{})
	var buffer bytes.Buffer
	for _, v := range a {
		buffer.WriteString(" " + v.(string) + " ")
	}
	return strings.TrimSpace(buffer.String())
}
