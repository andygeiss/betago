package main

import (
	"log"

	"github.com/andygeiss/miabot/application/betago"
	"github.com/andygeiss/miabot/business/bot"
	"github.com/andygeiss/miabot/infrastructure/udp"
)

func main() {
	name := "BetaGo"
	controller := udp.NewController("172.17.0.3:9000")
	engine := betago.NewEngine(name)
	bot := bot.NewBot(name, controller, engine)
	if err := bot.Setup(); err != nil {
		log.Fatal(err)
	}
	if err := bot.Loop(); err != nil {
		log.Fatal(err)
	}
}
