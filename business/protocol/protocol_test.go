package protocol_test

import (
	"testing"

	"github.com/andygeiss/miabot/business/protocol"
)

func TestIfAnnounceIsValid(t *testing.T) {
	commands := make(chan string)
	go protocol.Announce("4,1", "TOKEN", commands)
	select {
	case command := <-commands:
		if command != "ANNOUNCE;4,1;TOKEN" {
			t.Errorf("Announce failed!")
		}
	}
}
func TestIfJoinIsValid(t *testing.T) {
	commands := make(chan string)
	go protocol.Join("TOKEN", commands)
	select {
	case command := <-commands:
		if command != "JOIN;TOKEN" {
			t.Errorf("Join failed!")
		}
	}
}
func TestIfRollIsValid(t *testing.T) {
	commands := make(chan string)
	go protocol.Roll("TOKEN", commands)
	select {
	case command := <-commands:
		if command != "ROLL;TOKEN" {
			t.Errorf("Roll failed!")
		}
	}
}

func TestIfSeeIsValid(t *testing.T) {
	commands := make(chan string)
	go protocol.See("TOKEN", commands)
	select {
	case command := <-commands:
		if command != "SEE;TOKEN" {
			t.Errorf("See failed!")
		}
	}
}
