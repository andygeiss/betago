package main

import (
	"log"

	"github.com/andygeiss/miabot/application/alphago"
)

func main() {
	name := "BetaGo"
	controller := alphago.NewController("172.17.0.3:9000")
	engine := alphago.NewEngine(name)
	bot := alphago.NewBot(name, controller, engine)
	if err := bot.Setup(); err != nil {
		log.Fatal(err)
	}
	if err := bot.Loop(); err != nil {
		log.Fatal(err)
	}
}
