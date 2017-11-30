package betago

import (
	"math/rand"
	"strings"
	"sync"

	"github.com/andygeiss/miabot/business/dice"
	"github.com/andygeiss/miabot/business/engine"
	"github.com/andygeiss/miabot/business/protocol"
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
	Token          string
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
	fields := strings.Split(message, ";")
	keyword := fields[0]
	switch keyword {
	case "ANNOUNCED":
		player, dice := fields[1], fields[2]
		e.mutex.Lock()
		useInfraredToSeeThroughPlayersBluff(player, dice, e.brain)
		e.mutex.Unlock()
	case "PLAYER LOST": //player, reason := fields[1], fields[2]
	case "PLAYER ROLLS": //player := fields[1]
	case "PLAYER WANTS TO SEE": //player := fields[1]
	case "ROLLED":
		dice, token := fields[1], fields[2]
		e.mutex.Lock()
		dice = upgradeDiceWithSuperpower(dice, e.brain)
		e.mutex.Unlock()
		protocol.Announce(dice, token, commands)
	case "ROUND STARTED": // players := fields[1]
	case "ROUND STARTING":
		token := fields[1]
		e.mutex.Lock()
		e.brain.Token = token
		e.mutex.Unlock()
		protocol.Join(token, commands)
	case "SCORE": // players := fields[1]
		e.mutex.Lock()
		e.brain.LastAnnounced = ""
		e.brain.LastAnnounced2 = ""
		e.brain.ShouldWeSee = false
		e.mutex.Unlock()
	case "YOUR TURN":
		token := fields[1]
		e.mutex.Lock()
		if weAreFirst(e.brain) {
			protocol.Roll(token, commands)
		} else {
			if shouldWeSee(e.brain) {
				protocol.See(token, commands)
			} else {
				protocol.Roll(token, commands)
			}
		}
		e.mutex.Unlock()
	}
	return nil
}

func shouldWeSee(brain *Brain) bool {
	return brain.ShouldWeSee
}

func upgradeDiceWithSuperpower(rolled string, brain *Brain) string {
	current, _ := dice.Parse(rolled)
	previous, _ := dice.Parse(brain.LastAnnounced)
	if previous >= current {
		return dice.ToString(previous + 1 + rand.Intn(3))
	}
	return rolled
}

func useInfraredToSeeThroughPlayersBluff(player, rolled string, brain *Brain) {
	if brain.LastAnnounced != "" {
		last, _ := dice.Parse(brain.LastAnnounced)
		if last == -1 {
			brain.ShouldWeSee = true
		} else {
			last2, _ := dice.Parse(brain.LastAnnounced2)
			var diff int
			if last2 == -1 {
				diff = 0
			} else {
				diff = last - last2
			}
			if diff <= 2 {
				brain.ShouldWeSee = true
			}
		}
	}
	brain.LastAnnounced2 = brain.LastAnnounced
	brain.LastAnnounced = rolled
	brain.LastPlayer = player
}

func weAreFirst(brain *Brain) bool {
	return brain.LastAnnounced == ""
}
