package command

import (
	"github.com/hashicorp/go-multierror"
	"log"

	"os"

	"bytes"

	"fmt"

	"github.com/younisshah/jakob/jkafka"
	"github.com/younisshah/jakob/machine/task-sender"
	"github.com/younisshah/jakob/redisd"
)

// Command represents the T38 server (obtained from hash ring)
// on which the command (Name and Args) is to be executed
// setting the response of execution in Result and Error
type Command struct {
	redisAddress    string
	name            string
	pipelineCmdName []string
	pipelineArgs    []interface{}
	args            []interface{}
	Result          interface{}
	Error           error
}

var logger = log.New(os.Stderr, "[jakob-cmd] ", log.LstdFlags)

// NewCommand create a new command
func NewCommand(redisServer string, cmdName string, args ...interface{}) *Command {
	return &Command{redisAddress: redisServer, name: cmdName, args: args}
}

// NewPipelinedCommand creates a new pipelined command
func NewPipelinedCommand(redisServer string, commands []string, args []interface{}) *Command {
	return &Command{redisAddress: redisServer, pipelineCmdName: commands, pipelineArgs: args}
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
	if err == nil && (c.name != "GET" && c.name != "PING" && c.name != "NEARBY") {
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

// PipelinedExecute executes T38 commands using Redis pipeline
func (c *Command) PipelinedExecute() {
	var result *multierror.Error
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
	if len(c.pipelineCmdName) == len(c.pipelineArgs) {
		for i := range c.pipelineCmdName {
			if err := conn.Send(c.pipelineCmdName[i], c.pipelineArgs[i].([]interface{})...); err != nil {
				result = multierror.Append(result, err)
			}
		}
		if err := conn.Flush(); err != nil {
			result = multierror.Append(result, err)
		}
		resp, err := conn.Receive()
		if err != nil {
			c.Error = err
			c.Result = nil
			return
		}
		c.Error = result.ErrorOrNil()
		c.Result = resp
		if c.Error == nil {
			for i := range c.pipelineCmdName {
				go func(i int) {
					if c.pipelineCmdName[i] != "GET" && c.pipelineCmdName[i] != "PING" && c.pipelineCmdName[i] != "NEARBY" {
						var cmd string
						cmd = c.pipelineCmdName[i]
						a := c.pipelineArgs[i].([]interface{})
						for k := range a {
							cmd = cmd + " " + fmt.Sprint(a[k]) + " "
						}
						logger.Println("producing to kafka")
						logger.Println(" - COMMAND:", cmd)
						if err := jkafka.Produce(cmd); err != nil {
							logger.Println("couldn't produce to Kafka")
							logger.Println(" -ERROR", err)
							return
						}
						task_sender.Send(c.pipelineCmdName[i], c.pipelineArgs[i])
					} else {
						if c.pipelineCmdName[i] != "GET" && c.pipelineCmdName[i] != "PING" {
							logger.Println("[*] Error while producing cmd to machinery")
							logger.Println("Err", c.Error)
						}
					}
				}(i)
			}
		} else {
			logger.Println("couldn't produce to kafka", c.Error)
		}
	}
}

// Stringer interface
func (c *Command) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(c.name)
	if len(c.args) > 0 {
		for i := range c.args {
			buffer.WriteString(" " + fmt.Sprint(c.args[i]) + " ")
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
