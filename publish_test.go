package ems

import "testing"

func TestNewPublisher(t *testing.T) {

	ops := NewClientOptions().SetServerUrl("tcp://127.0.0.1:7222").SetUsername("admin").SetPassword("")

	c := NewPublisher(ops).(*client)

	assertNewClient(c, t)
}

func TestPublisher_Connect(t *testing.T) {

	ops := NewClientOptions().SetServerUrl("tcp://127.0.0.1:7222").SetUsername("admin").SetPassword("")

	c := NewPublisher(ops).(*client)

	err := c.Connect()
	if err != nil {
		t.Fatalf(err.Error())
	}

	c.Disconnect()
}

func TestPublisher_Send(t *testing.T) {

	ops := NewClientOptions().SetServerUrl("tcp://127.0.0.1:7222").SetUsername("admin").SetPassword("")

	c := NewPublisher(ops).(*client)

	err := c.Connect()
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = c.Send("queue.sample", "hello, world", 0, "non_persistent", 10000)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = c.Disconnect()
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestPublisher_SendReceive(t *testing.T) {

	ops := NewClientOptions().SetServerUrl("tcp://127.0.0.1:7222").SetUsername("admin").SetPassword("")

	c := NewPublisher(ops).(*client)

	err := c.Connect()
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = c.SendReceive("queue.sample", "hello, world", "non_persistent", 1000)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = c.Disconnect()
	if err != nil {
		t.Fatalf(err.Error())
	}
}
