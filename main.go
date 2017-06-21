package main

import (
	"log"

	"github.com/younisshah/jakob/core"
	"github.com/younisshah/jakob/core/worker"
	"github.com/younisshah/jakob/jfs"
	"github.com/younisshah/jakob/jring"
	"github.com/younisshah/jakob/network"
	"gopkg.in/alecthomas/kingpin.v2"
)

/**
*  Created by Galileo on 17/6/17.
 */

var (
	//_init         = kingpin.Flag("init", "Initialize the server [y,n]").Short('i').Required().String()
	setterAddress = kingpin.Flag("setter", "Tile38 setter address string").Required().Short('s').String()
	getterAddress = kingpin.Flag("getter", "Tile38 getter address string").Required().Short('g').String()
)

func main() {

	kingpin.Version("Jakob 1.0.0")
	kingpin.Parse()

	go func() {
		panic(worker.Worker())
	}()

	jyaml := jfs.NewJYaml()
	err := jyaml.Init(*setterAddress, *getterAddress)

	if err != nil {
		panic(err)
	}

	// check if setter and getter are live!
	cmd := core.PING(*setterAddress)
	if cmd.Error != nil && "PONG" != cmd.Result.(string) {
		log.Println("Setter peer:", *setterAddress, " is not live!")
	}
	cmd = core.PING(*getterAddress)
	if cmd.Error != nil && "PONG" != cmd.Result.(string) {
		log.Println("Getter peer:", *getterAddress, " is not live!")
	}

	jr := jring.New(jfs.SETTER)
	peer, err := jr.Get("driver_2323")
	if err != nil {
		log.Println(err)
	}

	cmd = core.NewCommand(peer, "SET", "my", "home", "POINT", "12.4", "23.45")
	cmd.Execute()
	if cmd.Error != nil {
		log.Println(cmd.Error)
	} else {
		log.Println(cmd.Result)
	}

	/*terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	log.Println("jakob exiting")*/

	panic(network.StartHTTPServer())

}
