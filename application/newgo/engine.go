package newgo

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
	// Following the protocol each message from the server contains
	// the keyword (with additional data and a token separated by a semicolon).
	fields := strings.Split(message, ";")
	keyword := fields[0]
	switch keyword {
	case "ANNOUNCED":
	case "PLAYER LOST":
		player := fields[1]
		reason := fields[2]
		// Show why we lost!
		if player == e.Name {
			log.Printf("WE LOST! %s\n", reason)
		}
	case "PLAYER ROLLS":
		//player := fields[1]
	case "ROLLED":
		dice, token := fields[1], fields[2]
		commands <- fmt.Sprintf("ANNOUNCE;%s;%s", dice, token)
	case "ROUND STARTED":
	case "ROUND STARTING":
		token := fields[1]
		commands <- fmt.Sprintf("JOIN;%s", token)
	case "SCORE":
		//list := fields[1]
		//players := strings.Split(list, ",")
	case "YOUR TURN":
		token := fields[1]
		// commands <- fmt.Sprintf("SEE;%s", token)
		commands <- fmt.Sprintf("ROLL;%s", token)
	}
	return nil
}
