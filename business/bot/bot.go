package bot

// Bot specifies a MIA bot.
type Bot interface {
	Loop() error
	Setup() error
	State() int
}

const (
	// StateDisconnected ...
	StateDisconnected = 0
	// StateConnected ..
	StateConnected = 1
	// StateUnregistered ...
	StateUnregistered = 2
	// StateRegistered ...
	StateRegistered = 3
)
