package alphago

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	"github.com/andygeiss/miabot/business/engine"
)

// Engine ...
type Engine struct {
	Name      string
	announced Announcement
	statistic Statistic
	mutex     sync.Mutex
}

// Announcement ...
type Announcement struct {
	Pos    int
	Dice   string
	Player string
}

// Statistic ...
type Statistic struct {
	Lost    int
	Played  int
	Players int
}

// NewEngine creates a new engine and returns its address.
func NewEngine(name string) engine.Engine {
	return &Engine{name, Announcement{}, Statistic{0, 0, 0}, sync.Mutex{}}
}

// Handle ...
func (e *Engine) Handle(message string, commands chan<- string) error {
	// Following the protocol each message from the server contains
	// the keyword (with additional data and a token separated by a semicolon).
	fields := strings.Split(message, ";")
	keyword := fields[0]
	switch keyword {
	case "ANNOUNCED":
		e.mutex.Lock()
		announced := Announcement{Pos: e.announced.Pos + 1, Player: fields[1], Dice: fields[2]}
		e.announced = announced
		e.mutex.Unlock()
	case "PLAYER LOST":
		player := fields[1]
		reason := fields[2]
		e.mutex.Lock()
		// Show the reason why we lost!
		if reason != "MIA" {
			e.statistic.Played++
			if player == e.Name {
				e.statistic.Lost++
				log.Printf("[WE LOST! [ANNOUNCED %s @ %d FROM %s] [REASON %s]\n", e.announced.Dice, e.announced.Pos, e.announced.Player, reason)
			}
		}
		e.mutex.Unlock()
	case "PLAYER ROLLS":
		//player := fields[1]
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
		if e.announced.Pos == 0 {
			// Lets calc our bluffing dice which should create more pressure
			// to the following players by starting with a medium but not
			// to aggressive dice roll value.
			bluff := calcBluffingDice()
			// But if our actual dice is stil better, then we should use it instead.
			if isDiceBetter(dice, bluff) {
				command = fmt.Sprintf("ANNOUNCE;%s;%s", dice, token)
			} else {
				command = fmt.Sprintf("ANNOUNCE;%s;%s", bluff, token)
			}
		} else {
			if !isDiceBetter(dice, announced) {
				dice = calcBetterDice(announced)
			}
			command = fmt.Sprintf("ANNOUNCE;%s;%s", dice, token)
		}
		// Finnaly send the command
		commands <- command
	case "ROUND STARTED":
		e.mutex.Lock()
		e.announced = Announcement{0, "", ""}
		e.mutex.Unlock()
	case "ROUND STARTING":
		token := fields[1]
		commands <- fmt.Sprintf("JOIN;%s", token)
	case "SCORE":
		list := fields[1]
		players := strings.Split(list, ",")
		e.mutex.Lock()
		e.statistic.Players = len(players)
		e.mutex.Unlock()
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
	}
	return nil
}

func calcBetterDice(announced string) string {
	// Convert string representation into two Integer values.
	aparts := strings.Split(announced, ",")
	ap1, ap2 := aparts[0], aparts[1]
	ad1, _ := strconv.Atoi(ap1)
	ad2, _ := strconv.Atoi(ap2)
	d1 := ad1
	d2 := ad2
	// If a pair was announced then create a better pair.
	// This should increase our risk of caught bluffing
	// but we will lose if we announce a lower value.
	if d1 == d2 {
		d1++
		d2++
	} else {
		// Non-pair values are much safer.
		// We use a sneaky function
		// because it will be easy to detect bluffing
		// if we take the next higher dice EVERY single time.
		if d1-d2 == 1 {
			if d1 < 6 { // (3,2 => 4,3)  (4,3 => 5,4)  (5,4 => 6,5)
				d2 = d1
				d1++
			} else { // Exception: (6,5 => 1,1)
				d1 = 1
				d2 = 1
			}
		} else { // (3,1 => 4,3)  (4,1 4,2 => 5,4)  (5,1 5,2 5,3 => 6,5)
			d2 = d1
			d1++
		}
	}
	return fmt.Sprintf("%d,%d", d1, d2)
}

func calcBluffingDice() string {
	d1 := 4 + rand.Intn(1)
	d2 := 1 + rand.Intn(4)
	if d1 == d2 { // never return a pair because theres a high chance for caught bluffing
		d1++
		d2 -= rand.Intn(2)
	}
	return fmt.Sprintf("%d,%d", d1, d2)
}

func isBluffing(pos int, announced string) bool {
	aparts := strings.Split(announced, ",")
	ap1, ap2 := aparts[0], aparts[1]
	ad1, _ := strconv.Atoi(ap1)
	ad2, _ := strconv.Atoi(ap2)
	// With each player the chance is higher for bluffing.
	if (pos > 1 && ad1 == ad2) || (pos > 2 && ad1 >= 5) {
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
	// We have MIA
	if d1 == 2 && d2 == 1 {
		return true
	}
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
