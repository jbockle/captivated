package main

import (
	"log"

	"github.com/jbockle/captivated/server"
	"github.com/jbockle/captivated/server/config"
	"github.com/jbockle/captivated/server/services"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalln("Exiting due to panic:", err)
		}
	}()

	config.Init()
	services.Init()
	go services.StartDeleteExpiredTask()

	server.Serve()
}
