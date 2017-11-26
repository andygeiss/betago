package engine

// Engine handles each message from a controller and
// creates an individual response or an error on failure.
type Engine interface {
	Handle(message string, commands chan<- string) error
}
