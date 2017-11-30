package betago

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"github.com/andygeiss/miabot/business/engine"
)

// Engine ...
type Engine struct {
	Name  string
	brain *Brain
	mutex *sync.Mutex
}

// Brain ...
type Brain struct {
	LastAnnounced  string
	LastAnnounced2 string
	LastPlayer     string
	ShouldWeSee    bool
}

// NewEngine creates a new engine and returns its address.
func NewEngine(name string) engine.Engine {
	brain := &Brain{
		LastAnnounced: "",
	}
	return &Engine{name, brain, &sync.Mutex{}}
}

// Handle ...
func (e *Engine) Handle(message string, commands chan<- string) error {
	// Following the protocol each message from the server contains
	// the keyword (with additional data and a token separated by a semicolon).
	fields := strings.Split(message, ";")
	keyword := fields[0]
	switch keyword {
	case "ANNOUNCED":
		player, dice := fields[1], fields[2]
		e.mutex.Lock()
		if e.brain.LastAnnounced != "" {
			last, _ := ParseDice(e.brain.LastAnnounced)
			if last == -1 {
				e.brain.ShouldWeSee = true
			} else {
				last2, _ := ParseDice(e.brain.LastAnnounced2)
				var diff int
				if last2 == -1 {
					diff = 0
				} else {
					diff = last - last2
				}
				if diff <= 4 {
					e.brain.ShouldWeSee = true
				}
			}
		}
		e.brain.LastAnnounced2 = e.brain.LastAnnounced
		e.brain.LastAnnounced = dice
		e.brain.LastPlayer = player
		e.mutex.Unlock()
	case "PLAYER LOST":
		player, reason := fields[1], fields[2]
		e.mutex.Lock()
		fmt.Printf("[PLAYER [%20s] LOST! [%s]\n", player, reason)
		e.mutex.Unlock()
	case "PLAYER ROLLS":
		//player := fields[1]
		e.mutex.Lock()
		e.mutex.Unlock()
	case "PLAYER WANTS TO SEE":
		//player := fields[1]
		e.mutex.Lock()
		e.mutex.Unlock()
	case "ROLLED":
		dice, token := fields[1], fields[2]
		var command string
		current, _ := ParseDice(dice)
		e.mutex.Lock()
		previous, _ := ParseDice(e.brain.LastAnnounced)
		if previous >= current {
			command = fmt.Sprintf("ANNOUNCE;%s;%s", DiceToString(previous+1+rand.Intn(3)), token)
		} else {
			command = fmt.Sprintf("ANNOUNCE;%s;%s", DiceToString(current), token)
		}
		e.mutex.Unlock()
		commands <- command
	case "ROUND STARTED":
		//list := fields[1]
		//players := strings.Split(list, ",")
		e.mutex.Lock()
		e.mutex.Unlock()
	case "ROUND STARTING":
		token := fields[1]
		e.mutex.Lock()
		e.mutex.Unlock()
		commands <- fmt.Sprintf("JOIN;%s", token)
	case "SCORE":
		//list := fields[1]
		//players := strings.Split(list, ",")
		e.mutex.Lock()
		e.brain.LastAnnounced = ""
		e.brain.LastAnnounced2 = ""
		e.brain.ShouldWeSee = false
		e.mutex.Unlock()
	case "YOUR TURN":
		token := fields[1]
		var command string
		e.mutex.Lock()
		if e.brain.LastAnnounced == "" {
			command = fmt.Sprintf("ROLL;%s", token)
		} else {
			if e.brain.ShouldWeSee {
				command = fmt.Sprintf("SEE;%s", token)
			} else {
				command = fmt.Sprintf("ROLL;%s", token)
			}
		}
		e.mutex.Unlock()
		commands <- command
	}
	return nil
}
