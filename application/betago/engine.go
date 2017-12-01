package betago

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/andygeiss/betago/business/dice"
	"github.com/andygeiss/betago/business/engine"
	"github.com/andygeiss/betago/business/protocol"
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
		storePlayersDiceFrequency(player, dice, e.brain)
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
	diceLen := len(dice.DiceTable)
	for i := 0; i < diceLen; i++ {
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
		diceLen := len(dice.DiceTable)
		playerTable := brain.DiffFrequenceTable[player]
		if playerTable == nil {
			playerTable = make([]int, diceLen)
		}
		playerTable[diff]++
		brain.DiffFrequenceTable[player] = playerTable
		max := 0
		for i := 0; i < diceLen; i++ {
			max += playerTable[i]
		}
		cnt := playerTable[diff]
		freq := float32(cnt * 100 / max)
		brain.DiffFreq = freq
	}
}

func storePlayersDiceFrequency(player, announced string, brain *Brain) {
	diceLen := len(dice.DiceTable)
	playerTable := brain.DiceFrequencyTable[player]
	if playerTable == nil {
		playerTable = make([]int, diceLen)
	}
	value, _ := dice.Parse(announced)
	playerTable[value]++
	brain.DiceFrequencyTable[player] = playerTable
	max := 0
	for i := 0; i < diceLen; i++ {
		max += playerTable[i]
	}
	cnt := playerTable[value]
	freq := float32(cnt * 100 / max)
	brain.DiceFreq = freq
}

func upgradeDiceWithSuperpower(announced string, brain *Brain) string {
	current, _ := dice.Parse(announced)
	previous, _ := dice.Parse(brain.PreviousAnnounced)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	diceLen := len(dice.DiceTable)
	mia := diceLen - 1
	// If the previous announced value is larger than our current value then calculate a better one.
	// Ensure that we create a larger value than the previous one.
	// Some bots will immediatily check for a value above 8 for bluffing.
	if previous >= current {
		// We need to randomize our value to prevent being caught by bots which uses DIFF FREQUENCY checks.
		// If bots only checks for DICE values then its optimal to increase the previous value by ONE.
		// But later if bots will check for DIFF frequency, we will immediately get caught.
		// So we need to balance that out to minimize the DIFF frequency per value below 20%.
		bluff := previous + 1 + rnd.Intn(6)
		if bluff < mia {
			return dice.ToString(bluff)
		}
		bluff = mia - rnd.Intn(6)
		return dice.ToString(bluff)
	}
	// If theres no previous announcment then we are first in this round!
	// Ensure that we create medium pressure to the next player by using the closest possible value below 9 (6,1).
	// Some bots will immediatily check for a value above 8 for bluffing.
	if previous == -1 {
		// We need to randomize our starting value to prevent being caught by bots which uses DICE FREQUENCY checks.
		// Lets try the full bandwith 4,1 ... 6,5.
		bluff := 4 + rnd.Intn(9)
		if bluff > current {
			return dice.ToString(bluff)
		}
	}
	return announced
}

func useInfraredToSeeThroughPlayersBluff(player, announced string, brain *Brain) {
	brain.ShouldWeSee = false
	// If there's a previous announcement then we could check for bluffing.
	// Out strategy is to analyse the DICE frequency and DIFF frequency of a specific player.
	// Some bots will use some dices more than other like a fixed starting dice.
	// With DICE frequency we will detect how often that dice was used.
	// Some bots will increase the previous value by a fixed or less random algorithm.
	// With DIFF frequency we will detect how often that algorithm like simply adding 1 to the value was used.
	if brain.PreviousAnnounced != "" {
		current, _ := dice.Parse(announced)
		previous, _ := dice.Parse(brain.PreviousAnnounced)
		// If the previous dice was invalid like 7,1 then we should see no matter what.
		if previous == -1 {
			brain.ShouldWeSee = true
		} else {
			// Get the difference between the last announcment and current value.
			// There are 20 possible dices (100%)
			// Each dice has an equal chance of 5% over time.
			// Thus the DIFF frequency of each dice should be around 5 percent.
			diff := current - previous
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
