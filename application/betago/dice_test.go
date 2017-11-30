package betago_test

import (
	"testing"

	"github.com/andygeiss/miabot/application/betago"
)

func TestIfParseDiceOf_3_1_Returns_0(t *testing.T) {
	val, _ := betago.ParseDice("3,1")
	if val != 0 {
		t.Errorf("Dice 3,1 should have a value of 0! But value is %d.", val)
	}
}

func TestIfParseDiceOf_4_1_Returns_2(t *testing.T) {
	val, _ := betago.ParseDice("4,1")
	if val != 2 {
		t.Errorf("Dice 4,1 should have a value of 2! But value is %d.", val)
	}
}

func TestIfParseDiceOf_1_1_Returns_14(t *testing.T) {
	val, _ := betago.ParseDice("1,1")
	if val != 14 {
		t.Errorf("Dice 1,1 should have a value of 14! But value is %d.", val)
	}
}
func TestIfParseDiceOf_2_1_Returns_20(t *testing.T) {
	val, _ := betago.ParseDice("2,1")
	if val != 20 {
		t.Errorf("Dice 2,1 (MIA) should have a value of 20! But value is %d.", val)
	}
}

func TestIfDiceValueOf_0_Returns_3_1(t *testing.T) {
	dice := betago.DiceToString(0)
	if dice != "3,1" {
		t.Errorf("Dice value of 0 should become 3,1! But is %s.", dice)
	}
}

func TestIfDiceValueOf_2_Returns_4_1(t *testing.T) {
	dice := betago.DiceToString(2)
	if dice != "4,1" {
		t.Errorf("Dice value of 2 should become 4,1! But is %s.", dice)
	}
}

func TestIfDiceValueOf_14_Returns_1_1(t *testing.T) {
	dice := betago.DiceToString(14)
	if dice != "1,1" {
		t.Errorf("Dice value of 14 should become 1,1! But is %s.", dice)
	}
}

func TestIfDiceValueOf_20_Returns_2_1(t *testing.T) {
	dice := betago.DiceToString(20)
	if dice != "2,1" {
		t.Errorf("Dice value of 20 should become 2,1! But is %s.", dice)
	}
}
