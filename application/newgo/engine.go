package newgo

import (
	"strings"

	"github.com/andygeiss/miabot/business/protocol"

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
	fields := strings.Split(message, ";")
	keyword := fields[0]
	switch keyword {
	case "ANNOUNCED":
	case "PLAYER LOST": //player, reason := fields[1], fields[2]
	case "PLAYER ROLLS": //player := fields[1]
	case "PLAYER WANTS TO SEE": //player := fields[1]
	case "ROLLED":
		dice, token := fields[1], fields[2]
		protocol.Announce(dice, token, commands)
	case "ROUND STARTED": // players := fields[1]
	case "ROUND STARTING":
		token := fields[1]
		protocol.Join(token, commands)
	case "SCORE": // players := fields[1]
	case "YOUR TURN":
		token := fields[1]
		protocol.Roll(token, commands)
	}
	return nil
}
