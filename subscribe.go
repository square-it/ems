package ems

/*
#include <tibems.h>

typedef void* voidPointer;
typedef void (*OnMessageCallback)(tibemsMsgConsumer, tibemsMsg, voidPointer);

extern void OnMessage(tibemsMsgConsumer consumer, tibemsMsg msg, voidPointer closure);
*/
import "C"

import (
	"errors"
	"unsafe"
)

var (
	globalReceiver chan *Message
)

type Message struct {
	ReplyDestinationName string
	ReplyDestinationType string

	MsgContent string
}

type Subscriber interface {
	Client
	Listen(string, chan *Message) error
}

func NewSubscriber(o *ClientOptions) Subscriber {

	c := &client{}
	c.options = *o
	c.status = disconnected

	return c
}

func (client *client) Listen(destinationName string, receiver chan *Message) error {
	var destination C.tibemsDestination
	var session C.tibemsSession
	var msgConsumer C.tibemsMsgConsumer

	/* create the destination */
	//if (useTopic)
	//status := tibemsTopic_Create(&destination,name);
	//else
	status := C.tibemsDestination_Create(&destination, TIBEMS_QUEUE, C.CString(destinationName))
	if status != TIBEMS_OK {
		e, _ := client.getErrorContext()
		return errors.New(e)
	}

	/* create the session */
	status = C.tibemsConnection_CreateSession(client.conn, &session, TIBEMS_FALSE, TIBEMS_AUTO_ACKNOWLEDGE)
	if status != TIBEMS_OK {
		e, _ := client.getErrorContext()
		return errors.New(e)
	}

	/* create the consumer */
	status = C.tibemsSession_CreateConsumer(session, &msgConsumer, destination, nil, TIBEMS_FALSE)
	if status != TIBEMS_OK {
		e, _ := client.getErrorContext()
		return errors.New(e)
	}

	/* set the message listener */
	callback := C.OnMessageCallback(C.OnMessage)
	status = C.tibemsMsgConsumer_SetMsgListener(msgConsumer, callback, nil)
	if status != TIBEMS_OK {
		e, _ := client.getErrorContext()
		return errors.New(e)
	}

	globalReceiver = receiver
	return nil
}

//export OnMessage
func OnMessage(consumer C.tibemsMsgConsumer, msg C.tibemsMsg, closure C.voidPointer) {
	var buf *C.char
	buf = (*C.char)(C.calloc(32768, 1))
	defer C.free(unsafe.Pointer(buf))

	C.tibemsTextMsg_GetText(msg, &buf)
	msgContent := C.GoString(buf)

	var replyTo C.tibemsDestination
	status := C.tibemsMsg_GetReplyTo(msg, &replyTo)
	if status != TIBEMS_OK {
		return
	}

	var replyNameLength C.tibems_int = 128
	var replyName = (*C.char)(C.calloc(128, 1))
	defer C.free(unsafe.Pointer(replyName))
	status = C.tibemsDestination_GetName(replyTo, replyName, replyNameLength)

	defer C.tibemsMsg_Destroy(msg)

	result := &Message{
		ReplyDestinationName: C.GoString(replyName),
		MsgContent:           msgContent,
	}
	globalReceiver <- result
}
