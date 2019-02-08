package ems

/*
#include <tibems.h>
*/
import "C"
import (
	"errors"
	"sync"
	"sync/atomic"
	"unsafe"
)

type Client interface {
	IsConnected() bool
	Connect() error
	Disconnect() error
}

type client struct {
	conn         C.tibemsConnection
	cf           C.tibemsConnectionFactory
	errorContext C.tibemsErrorContext
	status       uint32
	options      ClientOptions
	sync.RWMutex
}

func (c *client) IsConnected() bool {

	c.RLock()
	defer c.RUnlock()

	return c.status == connected

}
func (c *client) Connect() error {

	c.RLock()
	defer c.RUnlock()

	status := C.tibemsErrorContext_Create(&c.errorContext)

	if status != TIBEMS_OK {
		return errors.New("failed to create error context")
	}

	c.cf = C.tibemsConnectionFactory_Create()

	url := c.options.GetServerUrl()

	status = C.tibemsConnectionFactory_SetServerURL(c.cf, C.CString(url.String()))
	if status != TIBEMS_OK {
		e, _ := c.getErrorContext()
		return errors.New(e)
	}

	// create the connection
	status = C.tibemsConnectionFactory_CreateConnection(c.cf, &c.conn, C.CString(c.options.username), C.CString(c.options.password))
	if status != TIBEMS_OK {
		e, _ := c.getErrorContext()
		return errors.New(e)
	}

	// start the connection
	status = C.tibemsConnection_Start(c.conn)
	if status != TIBEMS_OK {
		e, _ := c.getErrorContext()
		return errors.New(e)
	}

	c.setConnected(connected)

	return nil
}

func (c *client) Disconnect() error {

	c.RLock()
	defer c.RUnlock()

	if c.IsConnected() {

		status := C.tibemsConnection_Stop(c.conn)
		if status != TIBEMS_OK {
			return errors.New("failed to stop connection")
		}

		// close the connection
		status = C.tibemsConnection_Close(c.conn)
		if status != TIBEMS_OK {
			return errors.New("failed to close connection")
		}

		c.setConnected(disconnected)
	}

	return nil
}

func (c *client) connectionStatus() uint32 {
	c.RLock()
	defer c.RUnlock()
	status := atomic.LoadUint32(&c.status)
	return status
}

func (c *client) setConnected(status uint32) {
	c.RLock()
	defer c.RUnlock()
	atomic.StoreUint32(&c.status, status)
}

func (c *client) getErrorContext() (string, string) {

	var errorString, stackTrace = "", ""
	var buf *C.char
	defer C.free(unsafe.Pointer(buf))

	C.tibemsErrorContext_GetLastErrorString(c.errorContext, &buf)
	errorString = C.GoString(buf)

	C.tibemsErrorContext_GetLastErrorStackTrace(c.errorContext, &buf)
	stackTrace = C.GoString(buf)

	return errorString, stackTrace

}
