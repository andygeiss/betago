package mia

import (
	"fmt"
	"log"
	"strings"

	"github.com/andygeiss/miabot/business/engine"
)

// Engine ...
type Engine struct{}

// NewEngine creates a new engine and returns its address.
func NewEngine() engine.Engine {
	return &Engine{}
}

// Handle ...
func (e *Engine) Handle(message string, commands chan<- string) error {
	//
	log.Printf("Message [%s]\n", message)
	// Following the protocol each message from the server contains
	// the keyword (with additional data and a token separated by a semicolon).
	fields := strings.Split(message, ";")
	keyword := fields[0]
	switch keyword {
	case "ROUND STARTING":
		token := fields[1]
		commands <- fmt.Sprintf("JOIN;%s", token)
	case "YOUR TURN":
		token := fields[1]
		commands <- fmt.Sprintf("ROLL;%s", token)
		// If you don't trust the previous player
		// then you should call the player to show the dice
		// commands <- fmt.Sprintf("SEE;%s", token)
	case "ROLLED":
		dice, token := fields[1], fields[2]
		// If your dice is higher than the last announced dice
		// then you should announce the truth
		// else you should lie ;-).
		commands <- fmt.Sprintf("ANNOUNCE;%s;%s", dice, token)
	}
	return nil
}
