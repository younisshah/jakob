package health

import (
	"log"
	"time"

	"os"

	"github.com/younisshah/jakob/command"
	"github.com/younisshah/jakob/jfs"
	"golang.org/x/sync/errgroup"
)

var logger = log.New(os.Stderr, "[jakob-health-check] ", log.LstdFlags)

func Check() {
	ticker := time.NewTicker(time.Duration(10) * time.Minute)
	for {
		select {
		case <-ticker.C:
			var g errgroup.Group

			// setters
			g.Go(func() error {
				jyml := jfs.NewJYaml()
				jyml.Type = jfs.SETTER
				setters, err := jyml.Peers()
				if err != nil {
					logger.Println("couldn't get setter peers", err)
					return err
				}
				for i := range setters {
					c := command.PING(setters[i])
					c.Execute()
					if c.Error != nil || "PONG" != c.Result.(string) {
						logger.Printf("%s is NOT live\n", setters[i])
						logger.Println("removing from setter peers file")
						err = jyml.DeleteSetter(setters[i])
						if err != nil {
							logger.Printf("couldn't delete setter peer %s\n", setters[i])
						} else {
							logger.Printf("getter %s removed\n", setters[i])
						}
					} else {
						logger.Println(setters[i], " setter peer is healthy.")
					}
				}
				return nil
			})

			// getters
			g.Go(func() error {
				jyml := jfs.NewJYaml()
				jyml.Type = jfs.GETTER
				getters, err := jyml.Peers()
				if err != nil {
					logger.Println("couldn't get getter peers", err)
					return err
				}
				for i := range getters {
					c := command.PING(getters[i])
					c.Execute()
					if c.Error != nil || "PONG" != c.Result.(string) {
						logger.Printf("%s is NOT live\n", getters[i])
						logger.Println("removing from getter peers file")
						err = jyml.DeleteGetter(getters[i])
						if err != nil {
							logger.Printf("couldn't delete getter peer %s\n", getters[i])
						} else {
							logger.Printf("getter %s removed\n", getters[i])
						}
					} else {
						logger.Println(getters[i], " getter peer is healthy.")
					}
				}
				return nil
			})

			if err := g.Wait(); err != nil {
				logger.Println("an unexpected error occurred", err.Error())
			} else {
				logger.Println("health check done")
			}
		}
	}
}
