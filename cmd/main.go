package main

import (
	"log"

	"github.com/aaguero_meli/W17-G6-Bootcamp/cmd/server"
)

func main() {

	// conf and env
	app, err := server.LoadServerConf()
	if err != nil {
		log.Fatal(err.Error())
	}

	// set up globalRouter and api config
	globalRouter, err := app.SetUp()
	if err != nil {
		log.Fatal(err.Error())
	}

	// - run
	if err := app.Run(globalRouter); err != nil {
		log.Fatal(err.Error())
		return
	}
}
