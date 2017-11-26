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
	case "ROUND STARTING":
		token := fields[1]
		commands <- fmt.Sprintf("JOIN;%s", token)
	case "ROUND STARTED":
		e.mutex.Lock()
		e.announced = Announcement{}
		e.mutex.Unlock()
	case "ANNOUNCED":
		announced := Announcement{Player: fields[1], Dice: fields[2]}
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
		// commands <- fmt.Sprintf("SEE;%s", token)

		var announced string
		e.mutex.Lock()
		announced = e.announced.Dice
		e.mutex.Unlock()
		// If we are the first player, then we can't SEE
		if announced == "" {
			commands <- fmt.Sprintf("ROLL;%s", token)
		} else {
			aparts := strings.Split(announced, ",")
			ap1, ap2 := aparts[0], aparts[1]
			ad1, _ := strconv.Atoi(ap1)
			ad2, _ := strconv.Atoi(ap2)
			// If announced dice is invalid then we should check for bluffing!
			if ad1 > 6 || ad2 > 6 {
				commands <- fmt.Sprintf("SEE;%s", token)
			} else {
				commands <- fmt.Sprintf("ROLL;%s", token)
			}
		}
	case "ROLLED":
		dice, token := fields[1], fields[2]
		// If your dice is higher than the last announced dice
		// then you should announce the truth
		// else you should lie ;-).

		var announced string
		e.mutex.Lock()
		announced = e.announced.Dice
		e.mutex.Unlock()

		// If we are the first player then we must announce because we cannot lose.
		if announced == "" {
			commands <- fmt.Sprintf("ANNOUNCE;%s;%s", dice, token)
		} else { // If we are not first then we need to calculate our chance.
			aparts := strings.Split(announced, ",")
			ap1, ap2 := aparts[0], aparts[1]
			ad1, _ := strconv.Atoi(ap1)
			ad2, _ := strconv.Atoi(ap2)
			parts := strings.Split(dice, ",")
			p1, p2 := parts[0], parts[1]
			d1, _ := strconv.Atoi(p1)
			d2, _ := strconv.Atoi(p2)
			// If our dice is better then announced we should not lie!
			if (d1 == d2 && d1 > ad1) || (d1 > ad1) || (d1 == ad1 && d2 > ad2) {
				commands <- fmt.Sprintf("ANNOUNCE;%s;%s", dice, token)
			} else { // If announced dice was better than we should do something about it
				if ad1 == ad2 && ad1 < 6 {
					// Take the next higher pair
					d1 = ad1 + 1
					d2 = ad2 + 2
				} else if (d1 == ad1 && d2 < ad2) || (d1 < ad1) {
					// Or take the next higher dice
					d1 = ad1 + 1
				}
				dice = fmt.Sprintf("%d,%d", d1, d2)
				commands <- fmt.Sprintf("ANNOUNCE;%s;%s", dice, token)
			}
		}
	}
	return nil
}
