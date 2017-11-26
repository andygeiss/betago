package main

import (
	"log"

	"github.com/andygeiss/miabot/application/mia"
)

func main() {
	name := "Alpha"
	controller := mia.NewController("172.17.0.2:9000")
	engine := mia.NewEngine(name)
	bot := mia.NewBot(name, controller, engine)
	if err := bot.Setup(); err != nil {
		log.Fatal(err)
	}
	if err := bot.Loop(); err != nil {
		log.Fatal(err)
	}
}
