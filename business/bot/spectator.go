package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/andygeiss/betago/business/controller"
	"github.com/andygeiss/betago/business/engine"
)

// SpectatorBot is an example implementation of a MIA bot following the protocol.
// https://github.com/janernsting/maexchen/blob/master/protokoll.en.markdown
type SpectatorBot struct {
	name  string
	state int
	ctrl  controller.Controller
	eng   engine.Engine
}

// NewSpectatorBot creates a new bot and returns its address.
func NewSpectatorBot(name string, ctrl controller.Controller, eng engine.Engine) Bot {
	return &SpectatorBot{name, StateDisconnected, ctrl, eng}
}

// Loop ...
func (b *SpectatorBot) Loop() error {
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
func (b *SpectatorBot) Setup() error {
	// Early return if controller is not valid.
	if b.ctrl == nil {
		return fmt.Errorf(ErrorControllerAddressIsNil)
	}
	switch b.state {
	case StateDisconnected:
		if err := b.ctrl.Connect(); err != nil {
			return err
		}
		b.state = StateConnected
		// Create a channel to catch the responses.
		responses := make(chan string)
		go b.ctrl.Read(responses)
		// Now send the registration message.
		message := fmt.Sprintf("REGISTER_SPECTATOR;%s", b.name)
		b.ctrl.Write(message)
		// Handle the responses or timeout after 30 seconds.
		timeout := time.After(3 * time.Second)
		select {
		case response := <-responses:
			fields := strings.Split(response, ";")
			keyword := fields[0]
			switch keyword {
			case "REGISTERED":
				b.state = StateRegistered
			}
		case <-timeout:
		}
	}
	return nil
}

// State ...
func (b *SpectatorBot) State() int {
	return b.state
}
