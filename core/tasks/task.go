package tasks

import (
	"log"
	"os"
	"strings"

	"github.com/younisshah/jakob/network"
)

/**
*  Created by Galileo on 20/6/17.
 */
var logger = log.New(os.Stderr, "[jakob-task] ", log.LstdFlags)

// Sync sends T38 command to all getter peers
func Sync(args ...string) error {
	peer := args[0]
	cmdName := args[1]
	cmdArgs := pack(strings.Fields(args[2]))

	client, err := network.GetRedisConn(peer)
	if err != nil {
		return err
	}
	resp, err := client.Do(cmdName, cmdArgs...)
	if err != nil {
		return err
	}
	logger.Println("RESP", resp)
	return nil
}


func pack(args []string) []interface{} {
	a := make([]interface{}, len(args))
	for k, v := range args {
		a[k] = v
	}
	return a
}
