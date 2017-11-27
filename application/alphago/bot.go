package alphago

import (
	"fmt"
	"strings"
	"time"

	"github.com/andygeiss/miabot/business/engine"
	"github.com/andygeiss/miabot/business/bot"
	"github.com/andygeiss/miabot/business/controller"
)

// Bot is an example implementation of a MIA bot following the protocol.
// https://github.com/janernsting/maexchen/blob/master/protokoll.en.markdown
type Bot struct {
	name  string
	state int
	ctrl  controller.Controller
	eng   engine.Engine
}

const (
	// ErrorControllerAddressIsNil ...
	ErrorControllerAddressIsNil = "Controller address should not be nil"
	// ErrorControllerCommunicationTimeout ...
	ErrorControllerCommunicationTimeout = "Controller communication timeout"
	// ErrorEngineAddressIsNil ...
	ErrorEngineAddressIsNil = "Engine address should not be nil"
)

// NewBot creates a new bot and returns its address.
func NewBot(name string, ctrl controller.Controller, eng engine.Engine) bot.Bot {
	return &Bot{name, bot.StateDisconnected, ctrl, eng}
}

// Loop ...
func (b *Bot) Loop() error {
	// Early return if engine is not valid.
	if b.eng == nil {
		return fmt.Errorf(ErrorEngineAddressIsNil)
	}
	// Share data between the application by communicating via channels.
	commands := make(chan string)
	responses := make(chan string)
	for {
		go b.ctrl.Read(responses)
		timeout := time.After(3 * time.Second)
		select {
		// If the server sends a response ...
		case response := <-responses:
			// Then handle each response by using the engines
			// to create the corresponding answer/command.
			go b.eng.Handle(response, commands)
		// Take each command created and write it using the controller.
		case command := <-commands:
			b.ctrl.Write(command)
		//
		case <-timeout:
			return fmt.Errorf(ErrorControllerCommunicationTimeout)
		}
	}
}

// Setup ...
func (b *Bot) Setup() error {
	// Early return if controller is not valid.
	if b.ctrl == nil {
		return fmt.Errorf(ErrorControllerAddressIsNil)
	}
	switch b.state {
	case bot.StateDisconnected:
		if err := b.ctrl.Connect(); err != nil {
			return err
		}
		b.state = bot.StateConnected
		// Create a channel to catch the responses.
		responses := make(chan string)
		go b.ctrl.Read(responses)
		// Now send the registration message.
		message := fmt.Sprintf("REGISTER;%s", b.name)
		b.ctrl.Write(message)
		// Handle the responses or timeout after 30 seconds.
		timeout := time.After(3 * time.Second)
		select {
		case response := <-responses:
			fields := strings.Split(response, ";")
			keyword := fields[0]
			switch keyword {
			case "REGISTERED":
				b.state = bot.StateRegistered
			}
		case <-timeout:
		}
	}
	return nil
}

// State ...
func (b *Bot) State() int {
	return b.state
}
