package betago

import (
	"fmt"
	"log"
	"strings"

	"github.com/andygeiss/miabot/business/engine"
)

// Engine ...
type Engine struct {
	Name string
}

// NewEngine creates a new engine and returns its address.
func NewEngine(name string) engine.Engine {
	return &Engine{name}
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
	case "SCORE":
	case "ROUND STARTING":
		token := fields[1]
		commands <- fmt.Sprintf("JOIN;%s", token)
	case "ROUND STARTED":
	case "ANNOUNCED":
	case "YOUR TURN":
		token := fields[1]
		// If we are first, then we cannot lose by rolling a low valued dice.
		var command string
		// command = fmt.Sprintf("SEE;%s", token)
		command = fmt.Sprintf("ROLL;%s", token)
		// Finally send the command
		commands <- command
	case "ROLLED":
		dice, token := fields[1], fields[2]
		// If your dice is higher than the last announced dice
		// then you should announce the truth
		// else you should lie ;-).
		var command string
		command = fmt.Sprintf("ANNOUNCE;%s;%s", dice, token)
		// Finnaly send the command
		commands <- command
	}
	return nil
}
