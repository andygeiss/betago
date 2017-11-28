package betago_test

import (
	"testing"
	"time"

	"github.com/andygeiss/miabot/application/betago"
)

func TestIfControllerCanConnectWithoutAnError(t *testing.T) {
	c := betago.NewController("172.20.0.3:9000")
	if err := c.Connect(); err != nil {
		t.Errorf("Connect should not return an error! %v", err.Error())
	}
}

func TestIfControllerCanDisconnectWithoutAnError(t *testing.T) {
	c := betago.NewController("172.17.0.3:9000")
	c.Connect()
	if err := c.Disconnect(); err != nil {
		t.Errorf("Disconnect should not return an error! %v", err.Error())
	}
}

func TestIfControllerCanWriteWithoutAnError(t *testing.T) {
	c := betago.NewController("172.17.0.3:9000")
	c.Connect()
	if err := c.Write("HEARTBEAT"); err != nil {
		t.Errorf("Write should not return an error! %v", err.Error())
	}
}

func TestIfControllerCanReadWithoutAnError(t *testing.T) {
	c := betago.NewController("172.17.0.3:9000")
	c.Connect()
	c.Write("REGISTER;AlphaGo")
	responses := make(chan string)
	go c.Read(responses)
	timeout := time.After(3 * time.Second)
	select {
	case <-responses:
	case <-timeout:
		t.Error("Read should not cause a timeout!")
	}
}
