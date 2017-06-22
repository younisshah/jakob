// Package core deal with setting up redis client, parsing incoming T38 commands,
// executing commands and returning the resp and error and handling joining of new
// setter + getter peer and beanstalkd services
package command

import (
	"log"

	"os"

	"bytes"

	"github.com/younisshah/jakob/jkafka"
	"github.com/younisshah/jakob/machine/task-sender"
	"github.com/younisshah/jakob/redisd"
)

/**
*  Created by Galileo on 18/6/17.
 */

// Command represents the T38 server (obtained from hash ring)
// on which the command (Name and Args) is to be executed
// setting the response of execution in Result and Error
type Command struct {
	redisAddress string
	name         string
	args         []interface{}
	Result       interface{}
	Error        error
}

var logger = log.New(os.Stderr, "[jakob-cmd] ", log.LstdFlags)

// NewCommand create a new command
func NewCommand(redisServer string, cmdName string, args ...interface{}) *Command {
	return &Command{redisAddress: redisServer, name: cmdName, args: args}
}

// Executes the command
func (c *Command) Execute() {
	conn, err := redisd.GetRedisConn(c.redisAddress)
	if err != nil {
		c.Error = err
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Println("couldn't close Redis connection", err)
		}
	}()
	resp, err := conn.Do(c.name, c.args...)
	c.Error = err
	c.Result = resp
	if err == nil && (c.name != "GET" && c.name != "PING") {
		cmd := c.String()
		logger.Println("producing to kafka")
		logger.Println(" - COMMAND:", cmd)
		if err := jkafka.Produce(cmd); err != nil {
			logger.Println("couldn't produce to Kafka")
			logger.Println(" -ERROR", err)
			return
		}
		task_sender.Send(c.name, c.args)
	} else {
		if c.name != "GET" && c.name != "PING" {
			logger.Printf("[*] Error while producing cmd [%v] to machinery", c.String())
			logger.Println("Err", err)
		}
	}
}

// Stringer interface
func (c *Command) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(c.name)
	if len(c.args) > 0 {
		for i := range c.args {
			buffer.WriteString(" " + c.args[i].(string) + " ")
		}
		return buffer.String()
	}
	return c.name
}

// Pings the give address to check if it's live
func PING(address string) *Command {
	cmd := NewCommand(address, "PING", nil)
	cmd.Execute()
	return cmd
}
