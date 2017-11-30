package protocol

import "fmt"

// Announce ...
func Announce(dice, token string, commands chan<- string) {
	commands <- fmt.Sprintf("ANNOUNCE;%s;%s", dice, token)
}

// Join ...
func Join(token string, commands chan<- string) {
	commands <- fmt.Sprintf("JOIN;%s", token)
}

// Roll ...
func Roll(token string, commands chan<- string) {
	commands <- fmt.Sprintf("ROLL;%s", token)
}

// See ...
func See(token string, commands chan<- string) {
	commands <- fmt.Sprintf("SEE;%s", token)
}
