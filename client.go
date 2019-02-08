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

func (client *client) IsConnected() bool {

	client.RLock()
	defer client.RUnlock()

	return client.status == connected

}
func (client *client) Connect() error {

	client.RLock()
	defer client.RUnlock()

	status := C.tibemsErrorContext_Create(&client.errorContext)

	if status != TIBEMS_OK {
		return errors.New("failed to create error context")
	}

	client.cf = C.tibemsConnectionFactory_Create()

	url := client.options.GetServerUrl()

	status = C.tibemsConnectionFactory_SetServerURL(client.cf, C.CString(url.String()))
	if status != TIBEMS_OK {
		e, _ := client.getErrorContext()
		return errors.New(e)
	}

	// create the connection
	status = C.tibemsConnectionFactory_CreateConnection(client.cf, &client.conn, C.CString(client.options.username), C.CString(client.options.password))
	if status != TIBEMS_OK {
		e, _ := client.getErrorContext()
		return errors.New(e)
	}

	// start the connection
	status = C.tibemsConnection_Start(client.conn)
	if status != TIBEMS_OK {
		e, _ := client.getErrorContext()
		return errors.New(e)
	}

	client.setConnected(connected)

	return nil
}

func (client *client) Disconnect() error {

	client.RLock()
	defer client.RUnlock()

	if client.IsConnected() {

		status := C.tibemsConnection_Stop(client.conn)
		if status != TIBEMS_OK {
			return errors.New("failed to stop connection")
		}

		// close the connection
		status = C.tibemsConnection_Close(client.conn)
		if status != TIBEMS_OK {
			return errors.New("failed to close connection")
		}

		client.setConnected(disconnected)
	}

	return nil
}

func (client *client) connectionStatus() uint32 {
	client.RLock()
	defer client.RUnlock()
	status := atomic.LoadUint32(&client.status)
	return status
}

func (client *client) setConnected(status uint32) {
	client.RLock()
	defer client.RUnlock()
	atomic.StoreUint32(&client.status, status)
}

func (client *client) getErrorContext() (string, string) {

	var errorString, stackTrace = "", ""
	var buf *C.char
	defer C.free(unsafe.Pointer(buf))

	C.tibemsErrorContext_GetLastErrorString(client.errorContext, &buf)
	errorString = C.GoString(buf)

	C.tibemsErrorContext_GetLastErrorStackTrace(client.errorContext, &buf)
	stackTrace = C.GoString(buf)

	return errorString, stackTrace

}
