package mia

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/andygeiss/miabot/business/engine"
)

// Engine ...
type Engine struct {
	Name      string
	announced Announcement
	mutex     sync.Mutex
}

// Announcement ...
type Announcement struct {
	Pos    int
	Dice   string
	Player string
}

// NewEngine creates a new engine and returns its address.
func NewEngine(name string) engine.Engine {
	return &Engine{name, Announcement{}, sync.Mutex{}}
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
		e.mutex.Lock()
		e.announced = Announcement{0, "", ""}
		e.mutex.Unlock()
	case "ANNOUNCED":
		e.mutex.Lock()
		announced := Announcement{Pos: e.announced.Pos + 1, Player: fields[1], Dice: fields[2]}
		e.mutex.Unlock()
		// Only save announcement if its not our bot.
		if announced.Player != e.Name {
			e.mutex.Lock()
			e.announced = announced
			e.mutex.Unlock()
		}
	case "YOUR TURN":
		token := fields[1]
		// If you don't trust the previous player
		// then you should call the player to show the dice
		var announced string
		var command string
		var pos int
		e.mutex.Lock()
		announced = e.announced.Dice
		pos = e.announced.Pos
		e.mutex.Unlock()
		// If we are first, then we cannot lose by rolling a low valued dice.
		if isDiceEmpty(announced) {
			command = fmt.Sprintf("ROLL;%s", token)
		} else {
			if !isDiceValid(announced) || isBluffing(pos, announced) {
				command = fmt.Sprintf("SEE;%s", token)
			} else {
				command = fmt.Sprintf("ROLL;%s", token)
			}
		}
		// Finally send the command
		commands <- command
	case "ROLLED":
		dice, token := fields[1], fields[2]
		// If your dice is higher than the last announced dice
		// then you should announce the truth
		// else you should lie ;-).
		var announced string
		var command string
		e.mutex.Lock()
		announced = e.announced.Dice
		e.mutex.Unlock()
		// If we are first, then we cannot lose by announcing a low valued dice.
		if isDiceEmpty(announced) {
			command = fmt.Sprintf("ANNOUNCE;%s;%s", dice, token)
		} else {
			// If we are not first then we need to calculate our chance.
			if !isDiceBetter(dice, announced) {
				// Calculate a better dice if we are not better.
				dice = calcBetterDice(announced)
			}
			command = fmt.Sprintf("ANNOUNCE;%s;%s", dice, token)
		}
		// Finnaly send the command
		commands <- command
	}
	return nil
}

func calcBetterDice(announced string) string {
	// Convert string representation into two Integer values.
	aparts := strings.Split(announced, ",")
	ap1, ap2 := aparts[0], aparts[1]
	ad1, _ := strconv.Atoi(ap1)
	ad2, _ := strconv.Atoi(ap2)
	// Create a better dice
	d1 := ad1 + 1
	d2 := ad2 + 1
	return fmt.Sprintf("%d,%d", d1, d2)
}

func isBluffing(pos int, announced string) bool {
	// If previous player as first then he must not lie!
	if pos == 1 {
		return false
	}
	// With each player the chance is higher for bluffing
	// Currently we think its highly possible for a bluff after 3 turns!
	possibility := float32(pos) * float32(0.4)
	if possibility > 1.0 {
		return true
	}
	return false
}

func isDiceBetter(dice, announced string) bool {
	// Convert string representation into two Integer values.
	aparts := strings.Split(announced, ",")
	ap1, ap2 := aparts[0], aparts[1]
	ad1, _ := strconv.Atoi(ap1)
	ad2, _ := strconv.Atoi(ap2)
	parts := strings.Split(dice, ",")
	p1, p2 := parts[0], parts[1]
	d1, _ := strconv.Atoi(p1)
	d2, _ := strconv.Atoi(p2)
	// We have the better pair
	if d1 == d2 && ad1 == ad2 && d1 > ad1 {
		return true
	}
	// We have the better non-pair
	if (d1 != d2 && ad1 != ad2) && (d1 > ad1 && d2 > ad2) {
		return true
	}
	return false
}

func isDiceEmpty(dice string) bool {
	if dice == "" {
		return true
	}
	return false
}

func isDiceValid(dice string) bool {
	// Convert string representation into two Integer values.
	parts := strings.Split(dice, ",")
	p1, p2 := parts[0], parts[1]
	d1, _ := strconv.Atoi(p1)
	d2, _ := strconv.Atoi(p2)
	// If dice values are INVALID then return FALSE
	if d1 > 6 || d2 > 6 {
		return false
	}
	return true
}
