package betago

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
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
	DiceFreq           float32
	DiceFrequencyTable map[string][]int
	DiffFreq           float32
	DiffFrequenceTable map[string][]int
	PreviousAnnounced  string
	PreviousPlayer     string
	ShouldWeSee        bool
	Token              string
	ValueDiff          int
}

// NewEngine creates a new engine and returns its address.
func NewEngine(name string) engine.Engine {
	brain := &Brain{
		DiceFrequencyTable: make(map[string][]int, 0),
		DiffFrequenceTable: make(map[string][]int, 0),
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
	case "PLAYER LOST": // player, reason := fields[1], fields[2]
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
		regainEnergyForNextRound(e.brain)
		printStatistics(e.brain)
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

func printStatistics(brain *Brain) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	diceLen := len(dice.DiceTable)
	// Diff
	printLine()
	fmt.Printf("[%20s] ", "Diff-Frequency")
	for i := 0; i < diceLen; i++ {
		fmt.Printf("[%3d] ", i)
	}
	fmt.Print("\n")
	printLine()
	for player, playerTable := range brain.DiffFrequenceTable {
		fmt.Printf("[%20s] ", player)
		max := 0
		for i := 0; i < diceLen; i++ {
			max += playerTable[i]
		}
		for _, diff := range playerTable {
			freq := diff * 100 / max
			fmt.Printf("[%3d] ", freq)
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
	// Dice
	printLine()
	fmt.Printf("[%20s] ", "Dice-Frequency")
	for i := 0; i < diceLen; i++ {
		fmt.Printf("[%3s] ", dice.DiceTable[i])
	}
	fmt.Print("\n")
	printLine()
	for player, playerTable := range brain.DiceFrequencyTable {
		fmt.Printf("[%20s] ", player)
		max := 0
		for i := 0; i < diceLen; i++ {
			max += playerTable[i]
		}
		for _, diff := range playerTable {
			freq := diff * 100 / max
			fmt.Printf("[%3d] ", freq)
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func printLine() {
	fmt.Print("[--------------------] ")
	for i := 0; i <= 20; i++ {
		fmt.Print("[---] ")
	}
	fmt.Print("\n")

}

func regainEnergyForNextRound(brain *Brain) {
	brain.PreviousAnnounced = ""
	brain.ShouldWeSee = false
	brain.ValueDiff = 0
	brain.DiffFreq = 0.0
	brain.DiceFreq = 0.0
}

func shouldWeSee(brain *Brain) bool {
	return brain.ShouldWeSee
}

func storePlayersDiffFrequency(player string, diff int, brain *Brain) {
	if diff > 0 {
		playerTable := brain.DiffFrequenceTable[player]
		if playerTable == nil {
			playerTable = make([]int, 21)
		}
		playerTable[diff]++
		brain.DiffFrequenceTable[player] = playerTable
		max := 0
		for i := 0; i <= 20; i++ {
			max += playerTable[i]
		}
		cnt := playerTable[diff]
		freq := float32(cnt * 100 / max)
		brain.DiffFreq = freq
	}
}

func storePlayersDiceFrequency(player, announced string, brain *Brain) {
	playerTable := brain.DiceFrequencyTable[player]
	if playerTable == nil {
		playerTable = make([]int, 21)
	}
	value, _ := dice.Parse(announced)
	playerTable[value]++
	brain.DiceFrequencyTable[player] = playerTable
	max := 0
	for i := 0; i <= 20; i++ {
		max += playerTable[i]
	}
	cnt := playerTable[value]
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
			storePlayersDiceFrequency(player, announced, brain)
			storePlayersDiffFrequency(player, diff, brain)
			diceFreq := brain.DiceFreq
			diffFreq := brain.DiffFreq
			if diffFreq >= 20.0 || diceFreq >= 20.0 {
				brain.ShouldWeSee = true
			}
		}
	}
	brain.PreviousAnnounced = announced
	brain.PreviousPlayer = player
}

func weAreFirst(brain *Brain) bool {
	return brain.PreviousAnnounced == ""
}
