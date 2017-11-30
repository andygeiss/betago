package udp

import (
	"fmt"
	"net"

	"github.com/andygeiss/betago/business/controller"
)

// Controller ...
type Controller struct {
	address    string
	connection *net.UDPConn
}

const (
	// ErrorConnectionAddressIsNil ...
	ErrorConnectionAddressIsNil = "Connection address should not be nil"
)

// NewController creates a new controller and returns its address.
func NewController(address string) controller.Controller {
	return &Controller{address, nil}
}

// Connect ...
func (c *Controller) Connect() error {
	addr, err := net.ResolveUDPAddr("udp4", c.address)
	if err != nil {
		return fmt.Errorf("ResolveUDPAddr failed: %v", err.Error())
	}
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return fmt.Errorf("DialUDP failed: %v", err.Error())
	}
	c.connection = conn
	return nil
}

// Disconnect ...
func (c *Controller) Disconnect() error {
	// Early return if connection is not available.
	if c.connection == nil {
		return fmt.Errorf(ErrorConnectionAddressIsNil)
	}
	return c.connection.Close()
}

// Read ...
func (c *Controller) Read(message chan<- string) error {
	// Early return if connection is not available.
	if c.connection == nil {
		return fmt.Errorf(ErrorConnectionAddressIsNil)
	}
	buf := make([]byte, 1024)
	n, _, err := c.connection.ReadFromUDP(buf)
	if n == 0 || err != nil {
		return fmt.Errorf("ReadFromUDP failed: %v", err.Error())
	}
	buf = buf[:n]
	message <- string(buf)
	return nil
}

// Write ...
func (c *Controller) Write(message string) error {
	// Early return if connection is not available.
	if c.connection == nil {
		return fmt.Errorf(ErrorConnectionAddressIsNil)
	}
	n, err := c.connection.Write([]byte(message))
	if n == 0 || err != nil {
		return fmt.Errorf("Write failed: %v", err.Error())
	}
	return nil
}
