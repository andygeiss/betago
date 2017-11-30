package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/andygeiss/betago/application/betago"
	"github.com/andygeiss/betago/business/bot"
	"github.com/andygeiss/betago/infrastructure/udp"
)

var (
	// APPNAME ...
	APPNAME string
	// BUILD ...
	BUILD string
	// VERSION ...
	VERSION string
)

func main() {
	addr := flag.String("-addr", "", "MIA server address (host:port)")
	name := flag.String("-name", "BetaGo", "MIA Bot name")
	flag.Parse()
	if *addr == "" {
		printUsage()
		return
	}
	controller := udp.NewController(*addr)
	engine := betago.NewEngine(*name)
	bot := bot.NewDefaultBot(*name, controller, engine)
	if err := bot.Setup(); err != nil {
		log.Fatal(err)
	}
	if err := bot.Loop(); err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	fmt.Printf("%s %s (build %s)\n\n", APPNAME, VERSION, BUILD)
	fmt.Print("This MIA bot uses superhuman-powers like infrared-vision to beat others.\n\n")
	fmt.Print("Options:\n")
	flag.PrintDefaults()
	fmt.Print("\n")
	fmt.Print("Example:\n")
	fmt.Printf("\t%s -addr 172.17.0.2:9000 -name %s\n\n", APPNAME, APPNAME)
}
