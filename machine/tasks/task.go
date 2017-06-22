package tasks

import (
	"log"
	"os"
	"strings"

	"time"

	"github.com/garyburd/redigo/redis"
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

	client, err := GetRedisConn(peer)
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

const (
	TILE38_REDIS_IDLE_TIMEOUT = time.Duration(5) * time.Second
	TILE38_REDIS_NETWORK      = "tcp"
)

// GetRedisConn returns a raw redis connection
func GetRedisConn(address string) (redis.Conn, error) {
	dialOption := redis.DialConnectTimeout(TILE38_REDIS_IDLE_TIMEOUT)
	return redis.Dial(TILE38_REDIS_NETWORK, address, dialOption)
}
