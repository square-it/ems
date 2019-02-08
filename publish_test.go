package ems

import (
	"testing"
)

func getPublisherClient(ops *ClientOptions) *client {
	return NewPublisher(ops).(*client)
}

func TestNewPublisher(t *testing.T) {
	ops := getClientOptions()

	client := getPublisherClient(ops)

	assertNewClient(client, t)
}

func TestPublisher_Connect(t *testing.T) {
	ops := getClientOptions()

	client := getPublisherClient(ops)

	err := client.Connect()
	if err != nil {
		t.Fatalf(err.Error())
	}

	client.Disconnect()
}

func TestPublisher_Send(t *testing.T) {
	ops := getClientOptions()

	client := getPublisherClient(ops)

	err := client.Connect()
	if err != nil {
		t.Fatalf(err.Error())
	}

	//err = client.Send("queue.sample", "hello, world", 0, "non_persistent", 10000)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = client.Disconnect()
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestPublisher_SendReceive(t *testing.T) {
	ops := getClientOptions()

	client := getPublisherClient(ops)

	err := client.Connect()
	if err != nil {
		t.Fatalf(err.Error())
	}

	receiver := make(chan *Message)
	subscribe(t, &receiver)

	go func() {
		_, err = client.SendReceive("queue.sample", "hello, world", "non_persistent", 1000)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	select {
	case received := <-receiver: // block until a message is received on the queue
		ops := getClientOptions()

		client := getPublisherClient(ops)

		err := client.Connect()
		if err != nil {
			t.Fatalf(err.Error())
		}

		err = client.Send(received.ReplyDestinationName, received.MsgContent, 0, "non_persistent", 10000)
		if err != nil {
			t.Fatalf(err.Error())
		}

		err = client.Disconnect()
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	//select {}

	err = client.Disconnect()
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func subscribe(t *testing.T, receiver *chan *Message) {
	ops := getClientOptions()

	c := NewSubscriber(ops).(*client)

	err := c.Connect()
	if err != nil {
		t.Fatalf(err.Error())
	}

	go c.Listen("queue.sample", *receiver)
}
