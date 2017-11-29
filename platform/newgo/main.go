package main

import (
	"log"

	"github.com/andygeiss/miabot/application/newgo"
	"github.com/andygeiss/miabot/business/bot"
	"github.com/andygeiss/miabot/infrastructure/udp"
)

func main() {
	name := "NewGo"
	controller := udp.NewController("172.17.0.3:9000")
	engine := newgo.NewEngine(name)
	bot := bot.NewBot(name, controller, engine)
	if err := bot.Setup(); err != nil {
		log.Fatal(err)
	}
	if err := bot.Loop(); err != nil {
		log.Fatal(err)
	}
}
