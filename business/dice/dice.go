package dice

import "fmt"

const (
	// ErrorDiceIsNotValid ...
	ErrorDiceIsNotValid = "Dice is not valid"
)

// DiceTable ...
var DiceTable = []string{
	"3,1", "3,2",
	"4,1", "4,2", "4,3",
	"5,1", "5,2", "5,3", "5,4",
	"6,1", "6,2", "6,3", "6,4", "6,5",
	"1,1", "2,2", "3,3", "4,4", "5,5", "6,6",
	"2,1",
}

// Parse ...
func Parse(dice string) (int, error) {
	for val, str := range DiceTable {
		if str == dice {
			return val, nil
		}
	}
	return -1, fmt.Errorf(ErrorDiceIsNotValid)
}

// ToString ...
func ToString(value int) string {
	if value > 20 {
		return DiceTable[20]
	} else if value < 0 {
		return DiceTable[0]
	}
	return DiceTable[value]
}
