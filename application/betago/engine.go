package betago

import (
	"fmt"
	"log"
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
	// AMap stores the player announcements
	AMap map[string][]int
	// BMap stores the difference between the player announcement and the previous announcement
	BMap map[string][]int
	// RMap stores each round values
	RMap               map[int][]int
	RIndex             int
	RLastname          string
	RNextname          string
	Players            []string
	PIndex             int
	IsLastAnnounceFake bool
}

// NewEngine creates a new engine and returns its address.
func NewEngine(name string) engine.Engine {
	brain := &Brain{
		AMap:      make(map[string][]int, 0),
		BMap:      make(map[string][]int, 0),
		RMap:      make(map[int][]int, 0),
		RIndex:    0,
		RLastname: "",
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
		player := fields[1]
		dice := fields[2]
		e.mutex.Lock()
		addPlayersDice(player, dice, e.brain)
		e.brain.IsLastAnnounceFake = predictPlayersDiceIsFake(player, dice, e.brain)
		e.mutex.Unlock()
	case "PLAYER LOST":
		player, reason := fields[1], fields[2]
		if player == e.Name {
			log.Printf("[WE LOST! [%s]\n", reason)
		}
	case "PLAYER ROLLS":
		//player := fields[1]
	case "ROLLED":
		dice, token := fields[1], fields[2]
		e.mutex.Lock()
		dice = predictBestDiceAgainstPlayer(e.brain.RLastname, dice, e.brain)
		e.mutex.Unlock()
		commands <- fmt.Sprintf("ANNOUNCE;%s;%s", dice, token)
	case "ROUND STARTED":
		list := fields[1]
		players := strings.Split(list, ",")
		e.mutex.Lock()
		e.brain.RIndex++
		e.brain.RLastname = ""
		e.brain.Players = players
		e.brain.PIndex = 0
		e.mutex.Unlock()
	case "ROUND STARTING":
		token := fields[1]
		commands <- fmt.Sprintf("JOIN;%s", token)
	case "SCORE":
		//list := fields[1]
		//players := strings.Split(list, ",")
	case "YOUR TURN":
		token := fields[1]
		var command string
		e.mutex.Lock()
		if e.brain.RLastname == "" {
			command = fmt.Sprintf("ROLL;%s", token)
		} else if e.brain.IsLastAnnounceFake {
			command = fmt.Sprintf("SEE;%s", token)
		} else {
			command = fmt.Sprintf("ROLL;%s", token)
		}
		e.mutex.Unlock()
		commands <- command
	}
	return nil
}

var diceTable = []string{
	"3,1", "3,2",
	"4,1", "4,2", "4,3",
	"5,1", "5,2", "5,3", "5,4",
	"6,1", "6,2", "6,3", "6,4", "6,5",
	"1,1", "2,2", "3,3", "4,4", "5,5", "6,6",
	"2,1",
}

func addPlayersDice(player, dice string, brain *Brain) {
	amap := brain.AMap[player]
	bmap := brain.BMap[player]
	rindex := brain.RIndex
	rmap := brain.RMap[rindex]
	val := parseDice(dice)
	//
	brain.RLastname = player
	// AMap - Save the players announcement to its specific table
	if amap == nil {
		amap = make([]int, 0)
	}
	amap = append(amap, val)
	brain.AMap[player] = amap
	// RList - Save the announcement to the round table
	if rmap == nil {
		rmap = make([]int, 0)
	}
	rmap = append(rmap, val)
	brain.RMap[rindex] = rmap
	//
	var prev int
	index := len(rmap) - 2
	if index >= 0 {
		prev = rmap[index]
	} else {
		prev = val // diff will become 0
	}
	// DMap - Save the players difference to the previous announcement
	if bmap == nil {
		bmap = make([]int, 0)
	}
	bmap = append(bmap, val-prev) // Save the players difference to the previous announcement
	brain.BMap[player] = bmap
	/*
		for pl := range brain.AMap {
			fmt.Printf("[Player: %s]\n", pl)
			fmt.Printf("[Announcements: %v]\n", brain.AMap[pl])
			fmt.Printf("[Behaviour:     %v]\n", brain.BMap[pl])
		}
	*/
}

func nextPlayer(brain *Brain) string {
	index := brain.PIndex
	players := brain.Players
	if index == len(players)-1 {
		return "" // Player is the last player in the current round
	}
	return players[index+1]
}

func predictPlayersDiceIsFake(player, dice string, brain *Brain) bool {
	bmap := brain.BMap[player]
	val := parseDice(dice)
	freq := valueFrequency(val, bmap)
	//fmt.Printf("[Player: %20s] [Dice %s] is [%0.2f] a fake.\n", player, dice, freq)
	return freq > 50
}

func predictBestDiceAgainstPlayer(player, rolled string, brain *Brain) string {
	rmap := brain.RMap
	round := rmap[brain.RIndex]
	index := len(round)
	current := parseDice(rolled)
	if index == 0 {
		return toDice(current)
	}
	previous := round[index-1]
	best := previous + 1 // TODO: create a better/more random algorithm!
	if current > best {
		return toDice(current)
	}
	return toDice(best)
}

func valueFrequency(val int, list []int) float32 {
	freq := 0
	max := 0
	for _, stored := range list {
		if val > 0 {
			if val == stored {
				freq++
			}
			if stored != 0 {
				max++
			}
		}
	}
	if max == 0 {
		return 100
	}
	return float32(freq*100/max) * 20 // 20 values available
}

// parseDice returns the value of a given dice
// or -1 if its an invalid dice.
func parseDice(dice string) int {
	for val, str := range diceTable {
		if str == dice {
			return val
		}
	}
	return -1
}

// toDice returns the dice in string representation.
func toDice(val int) string {
	return diceTable[val]
}
