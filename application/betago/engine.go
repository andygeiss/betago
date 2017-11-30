package betago

import (
	"fmt"
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
	DiceFreq            float32
	PlayerAnnouncements map[string][]string
	PlayerBehaviour     map[string][]int
	PreviousAnnounced   string
	PreviousPlayer      string
	ShouldWeSee         bool
	Token               string
	ValueDiff           int
	ValueDiffFreq       float32
}

// NewEngine creates a new engine and returns its address.
func NewEngine(name string) engine.Engine {
	brain := &Brain{
		PlayerAnnouncements: make(map[string][]string, 0),
		PlayerBehaviour:     make(map[string][]int, 0),
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
		storePlayersAnnouncement(player, dice, e.brain)
		useInfraredToSeeThroughPlayersBluff(player, dice, e.brain)
		e.mutex.Unlock()
	case "PLAYER LOST":
		player, reason := fields[1], fields[2]
		if player == e.Name && reason != "MIA" {
			fmt.Printf("[PLAYER %20s] [LOST! %s]\n", player, reason)
		}
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
		e.brain.PreviousAnnounced = ""
		e.brain.ShouldWeSee = false
		e.brain.ValueDiff = 0
		e.brain.ValueDiffFreq = 0.0
		e.brain.DiceFreq = 0.0
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

func storePlayersAnnouncement(player, dice string, brain *Brain) {
	playerAnnouncement := brain.PlayerAnnouncements[player]
	if playerAnnouncement == nil {
		playerAnnouncement = make([]string, 0)
	}
	playerAnnouncement = append(playerAnnouncement, dice)
	brain.PlayerAnnouncements[player] = playerAnnouncement
}

func storePlayersBehaviour(player string, diff int, brain *Brain) {
	playerBehaviour := brain.PlayerBehaviour[player]
	if playerBehaviour == nil {
		playerBehaviour = make([]int, 0)
	}
	playerBehaviour = append(playerBehaviour, diff)
	brain.PlayerBehaviour[player] = playerBehaviour
	brain.ValueDiff = diff
}

func storePlayersDiffFrequency(player string, diff int, brain *Brain) {
	if diff > 0 {
		playerBehaviour := brain.PlayerBehaviour[player]
		cnt, max := 0, 0
		for _, val := range playerBehaviour {
			if diff == val {
				cnt++
			}
			max++
		}
		freq := float32(cnt * 100 / max)
		brain.ValueDiffFreq = freq
	}
}

func storePlayersDiceFrequency(player, announced string, brain *Brain) {
	cnt, max := 0, 0
	playerAnnouncements := brain.PlayerAnnouncements[player]
	for _, val := range playerAnnouncements {
		if announced == val {
			cnt++
		}
		max++
	}
	freq := float32(cnt * 100 / max)
	brain.DiceFreq = freq
}

func upgradeDiceWithSuperpower(announced string, brain *Brain) string {
	current, _ := dice.Parse(announced)
	previous, _ := dice.Parse(brain.PreviousAnnounced)
	if previous >= current {
		return dice.ToString(previous + 1 + rand.Intn(2))
	}
	return announced
}

func useInfraredToSeeThroughPlayersBluff(player, announced string, brain *Brain) {
	brain.ShouldWeSee = false
	if brain.PreviousAnnounced != "" {
		current, _ := dice.Parse(announced)
		previous, _ := dice.Parse(brain.PreviousAnnounced)
		if previous == -1 {
			brain.ShouldWeSee = true
		} else {
			diff := current - previous
			storePlayersBehaviour(player, diff, brain)
			storePlayersDiceFrequency(player, announced, brain)
			storePlayersDiffFrequency(player, diff, brain)
			diceFreq := brain.DiceFreq
			diffFreq := brain.ValueDiffFreq
			if diffFreq >= 20.0 || diceFreq >= 20.0 {
				brain.ShouldWeSee = true
			}
		}
	}
	brain.PreviousAnnounced = announced
	brain.PreviousPlayer = player
	//fmt.Printf("[PLAYER %20s] [%s] [%.2d] [%.1f | %.1f]\n", player, dice, brain.ValueDiff, brain.ValueDiffFreq, brain.DiceFreq)
}

func weAreFirst(brain *Brain) bool {
	return brain.PreviousAnnounced == ""
}
