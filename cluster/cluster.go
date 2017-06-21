// Package core deal with setting up redis client, parsing incoming T38 commands,
// executing commands and returning the resp and error and handling joining of new
// setter + getter peer and beanstalkd/gearman services
package cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/younisshah/jakob/jfs"
	"golang.org/x/sync/errgroup"
	"github.com/younisshah/jakob/core"
)

/**
*  Created by Galileo on 18/6/17.
 */

var logger = log.New(os.Stderr, "[jakob-http-server] ", log.LstdFlags)

func Join(writer http.ResponseWriter, req *http.Request) {

	setter, getter, err := parseBody(writer, req)
	if err != nil {
		s := "couldn't read request body:" + err.Error()
		sendResp(&writer, s)
		return
	}

	var group errgroup.Group

	group.Go(func() error {
		jyml := jfs.NewJYaml()
		jyml.Type = jfs.SETTER
		jyml.Address = setter
		err := jyml.Append()
		if err != nil {
			return err
		}
		cmd := core.PING(setter)
		if cmd.Error != nil && "PONG" != cmd.Result.(string) {
			log.Printf("setter peer %s is not live", setter)
			return cmd.Error
		}
		return nil
	})

	group.Go(func() error {
		jyml := jfs.NewJYaml()
		jyml.Type = jfs.GETTER
		jyml.Address = getter
		err := jyml.Append()
		if err != nil {
			return err
		}
		cmd := core.PING(getter)
		if cmd.Error != nil && "PONG" != cmd.Result.(string) {
			log.Printf("getter peer %s is not live", getter)
			return cmd.Error
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		s := "couldn't join cluster:" + err.Error()
		sendResp(&writer, s)
		return
	}

	logger.Println("joined cluster")
	logger.Println("getter:", getter)
	logger.Println("setter:", setter)
	sendResp(&writer, "OK")
}

// parseBody parses the HTTP request body and returns getter and setter
func parseBody(writer http.ResponseWriter, req *http.Request) (string, string, error) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", "", err
	}
	defer req.Body.Close()
	var body map[string]interface{}
	err = json.Unmarshal(b, &body)
	if err != nil {
		s := "couldn't read request body:" + err.Error()
		sendResp(&writer, s)
	}

	setter, ok := body["setter"]
	if !ok {
		sendResp(&writer, "setter is missing")
		return "", "", err
	}
	getter, ok := body["getter"]
	if !ok {
		return "", "", err
	}
	return setter.(string), getter.(string), nil
}

func sendResp(writer *http.ResponseWriter, s string) {
	logger.Println(s)
	fmt.Fprintf(*writer, s)
}
