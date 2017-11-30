package bot_test

import (
	"testing"

	"github.com/andygeiss/miabot/business/bot"
)

func TestSpectatorIsDisconnectedAtStartup(t *testing.T) {
	b := bot.NewSpectatorBot("SpectatorBot", nil, nil)
	if b.State() != bot.StateDisconnected {
		t.Error("State should be Disconnected at startup!")
	}
}

func TestSpectatorIsRegisteredAfterSetup(t *testing.T) {
	c := &MockupController{}
	c.Inbound = "REGISTERED"
	b := bot.NewSpectatorBot("SpectatorBot", c, nil)
	if err := b.Setup(); err != nil {
		t.Errorf("Setup should not return an error! %v", err.Error())
	}
	if b.State() != bot.StateRegistered {
		t.Error("State should be Registered at startup!")
	}
}
