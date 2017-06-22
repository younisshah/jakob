package replicate

import (
	"log"
	"os"

	"github.com/younisshah/jakob/redisd"
)

/**
*  Created by Galileo on 22/6/17.
 */

var rlogger = log.New(os.Stderr, "[jakob-cluster-replicator] ", log.LstdFlags)

func Replicate(peer string, cmds map[string][]interface{}) error {
	rlogger.Println("replicating")
	rlogger.Println(" - PEER", peer)
	rlogger.Println(" - CMDs", cmds)
	conn, err := redisd.GetRedisConn(peer)
	defer func() {
		if err := conn.Close(); err != nil {
			rlogger.Println("couldn't close redis connection")
			rlogger.Println(" -PEER", peer)
		}
	}()
	if err != nil {
		rlogger.Println("couldn't connect get Redis connection for peer: ", peer)
		return err
	}
	for k, v := range cmds {
		if err := conn.Send(k, v...); err != nil {
			rlogger.Println("error while sending command ", err)
		}
	}
	if err := conn.Flush(); err != nil {
		rlogger.Println("error while flushing to Redis", err)
		return err
	}
	resp, err := conn.Receive()
	rlogger.Println("replication resp: ", resp)
	return err
}
