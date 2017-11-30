package main

import (
	"log"

	"github.com/andygeiss/miabot/application/spectogo"
	"github.com/andygeiss/miabot/business/bot"
	"github.com/andygeiss/miabot/infrastructure/udp"
)

func main() {
	name := "SpectoGo"
	controller := udp.NewController("172.17.0.3:9000")
	engine := spectogo.NewEngine(name)
	bot := bot.NewSpectatorBot(name, controller, engine)
	if err := bot.Setup(); err != nil {
		log.Fatal(err)
	}
	if err := bot.Loop(); err != nil {
		log.Fatal(err)
	}
}
