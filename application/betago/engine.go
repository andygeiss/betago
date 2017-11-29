package betago

import (
	"fmt"
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
	RMap      map[int][]int
	RIndex    int
	RLastname string
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
		predictPlayersDiceIsFake(player, dice, e.brain)
		e.mutex.Unlock()
	case "PLAYER LOST":
		//player := fields[1]
		//reason := fields[2]
	case "PLAYER ROLLS":
		//player := fields[1]
	case "ROLLED":
		dice, token := fields[1], fields[2]
		commands <- fmt.Sprintf("ANNOUNCE;%s;%s", dice, token)
	case "ROUND STARTED":
		//list := fields[1]
		//players := strings.Split(list, ",")
		e.mutex.Lock()
		e.brain.RIndex++
		e.mutex.Unlock()
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

func predictPlayersDiceIsFake(player, dice string, brain *Brain) string {
	bmap := brain.BMap[player]
	val := parseDice(dice)
	freq := valueFrequency(val, bmap)
	fmt.Printf("[Player: %20s] [Dice %s] is [%0.2f] a fake.\n", player, dice, freq)
	return ""
}

func valueFrequency(val int, list []int) float32 {
	freq := 0
	max := 0
	for _, stored := range list {
		if val == stored {
			freq++
		}
		if stored != 0 {
			max++
		}
	}
	if max == 0 {
		return 100
	}
	return float32(freq * 100 / max)
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
