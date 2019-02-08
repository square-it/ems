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
	"fmt"
	"unsafe"
)

type Subscriber interface {
	Client
	Listen(destination string) error
}

func NewSubscriber(o *ClientOptions) Subscriber {

	c := &client{}
	c.options = *o
	c.status = disconnected

	return c
}

func (c *client) Listen(destinationName string) error {

	var destination C.tibemsDestination
	var session C.tibemsSession
	var msgConsumer C.tibemsMsgConsumer

	/* create the destination */
	//if (useTopic)
	//status := tibemsTopic_Create(&destination,name);
	//else
	status := C.tibemsDestination_Create(&destination, TIBEMS_QUEUE, C.CString(destinationName))
	if status != TIBEMS_OK {
		e, _ := c.getErrorContext()
		return errors.New(e)
	}

	/* create the session */
	status = C.tibemsConnection_CreateSession(c.conn, &session, TIBEMS_FALSE, TIBEMS_AUTO_ACKNOWLEDGE)
	if status != TIBEMS_OK {
		e, _ := c.getErrorContext()
		return errors.New(e)
	}

	/* create the consumer */
	status = C.tibemsSession_CreateConsumer(session, &msgConsumer, destination, nil, TIBEMS_FALSE)
	if status != TIBEMS_OK {
		e, _ := c.getErrorContext()
		return errors.New(e)
	}

	/* set the message listener */
	callback := C.OnMessageCallback(C.OnMessage)
	status = C.tibemsMsgConsumer_SetMsgListener(msgConsumer, callback, nil)
	if status != TIBEMS_OK {
		e, _ := c.getErrorContext()
		return errors.New(e)
	}

	return nil
}

//export OnMessage
func OnMessage(consumer C.tibemsMsgConsumer, msg C.tibemsMsg, closure C.voidPointer) {
	var buf *C.char
	buf = (*C.char)(C.calloc(32768, 1))
	defer C.free(unsafe.Pointer(buf))

	C.tibemsTextMsg_GetText(msg, &buf)

	result := C.GoString(buf)
	fmt.Println(result)

	C.tibemsMsg_Destroy(msg)
}
