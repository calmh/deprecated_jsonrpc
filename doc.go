/*

**DEPRECATION WARNING**

This package is unmaintained and probably not what you're looking for.

Package jsonrpc implements a standards compliant, asynchronous JSON-RPC client.
The JSON-RPC 2.0 standard as specified in http://www.jsonrpc.org/specification
is supported, while it is also possible to implement vendor specific dialects.
The client is thread safe, i.e. request and notification functions can be
called from any goroutine.

The RPC client is instantiated on top of an io.ReadWriter, for example a
net.Conn. You then create Request and Notification functions from the RPC
client. These are bound to the RPC client and return a chan Response (Request
only) and possibly an error when called. When the result is needed a read from
the channel will yield the Response object.

There are two stages where error handling is necessary:

- When calling the generated Request or Notification functions; these return an
error if the underlying ReadWriter is known bad (due to a previous read or
write error).

- When reading from the Response channel; the channel will be closed instead of
producing a value if there is an error in sending the request or receiving the
response.

Getting a nil error from a request function does not imply the request was
successfully sent. There are several more layers where errors may occur before
the request reaches the remote server, several of which are outside of the
control of this package. Only on receipt of a Response object can you be sure
that a request was processed by the server.

See the Connection example for details.
*/
package jsonrpc
