package ems

import "testing"

func TestNewSubscriber(t *testing.T) {
	ops := NewClientOptions().SetServerUrl("tcp://127.0.0.1:7222").SetUsername("admin").SetPassword("")

	c := NewSubscriber(ops).(*client)
	assertNewClient(c, t)
}

func TestSubscriber_Receive(t *testing.T) {

	ops := NewClientOptions().SetServerUrl("tcp://127.0.0.1:7222").SetUsername("admin").SetPassword("")

	c := NewSubscriber(ops).(*client)

	err := c.Connect()
	if err != nil {
		t.Fatalf(err.Error())
	}

	go c.Listen("queue.sample")
	if err != nil {
		t.Fatalf(err.Error())
	}

	//select {} // # infinite loop

	err = c.Disconnect()
	if err != nil {
		t.Fatalf(err.Error())
	}

}
