package controller

// Controller connects to a specific MIA backend and
// handles the communication by reading/writing messages.
type Controller interface {
	Connect() error
	Disconnect() error
	Read(message chan<- string) error
	Write(message string) error
}
