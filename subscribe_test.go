package ems

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getSubscriberClient(ops *ClientOptions) *client {
	return NewSubscriber(ops).(*client)
}

func TestNewSubscriber(t *testing.T) {
	ops := getClientOptions()

	client := getSubscriberClient(ops)
	assertNewClient(client, t)
}

func TestSubscriber_Receive(t *testing.T) {

	ops := getClientOptions()

	c := NewSubscriber(ops).(*client)

	err := c.Connect()
	if err != nil {
		t.Fatalf(err.Error())
	}

	receiver := make(chan *Message)
	go c.Listen("queue.sample", receiver)
	if err != nil {
		t.Fatalf(err.Error())
	}

	testMessage := "hello, world"
	go func() {
		ops := getClientOptions()

		client := getPublisherClient(ops)

		err := client.Connect()
		if err != nil {
			t.Fatalf(err.Error())
		}

		err = client.Send("queue.sample", testMessage, 0, "non_persistent", 10000)
		if err != nil {
			t.Fatalf(err.Error())
		}

		err = client.Disconnect()
		if err != nil {
			t.Fatalf(err.Error())
		}
	}() // send a message on the queue

	select {
	case received := <-receiver: // block until a message is received on the queue
		assert.Equal(t, received.MsgContent, testMessage)
	}

	err = c.Disconnect()
	if err != nil {
		t.Fatalf(err.Error())
	}

}
