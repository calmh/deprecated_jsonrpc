/*
Package jsonrpc implements a standards compliant, asynchronous JSON-RPC client.
The JSON-RPC 2.0 standard as specified in http://www.jsonrpc.org/specification
is supported, while it is also possible to implement vendor specific dialects.
The client is thread safe, i.e. request and notification functions can be
called from any goroutine.
*/
package jsonrpc
