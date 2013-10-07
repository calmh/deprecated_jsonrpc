jsonrpc [![Build Status](https://drone.io/github.com/calmh/jsonrpc/status.png)](https://drone.io/github.com/calmh/jsonrpc/latest)
=======

Package jsonrpc implements a standards compliant, asynchronous JSON-RPC client.
The JSON-RPC 2.0 standard as specified in http://www.jsonrpc.org/specification
is supported, while it is also possible to implement vendor specific dialects.
The client is thread safe, i.e. request and notification functions can be
called from any goroutine.

Documentation
-------------

http://godoc.org/github.com/calmh/jsonrpc

Example
-------

```go
// Connect to a remote JSON-RPC server

conn, err := net.Dial("tcp", "svr.example.com:3994")
if err != nil {
	panic(err)
}

// Set up a JSON-RPC channel on top of the socket, standard JSON-RPC 2.0
// dialect

rpc := jsonrpc.NewConnection(conn, jsonrpc.StandardDialect)

// Create two request functions, one for the method "hello" and one for the
// method "system.ping".

hello := rpc.Request("hello")
ping := rpc.Request("system.ping")

// Call the hello method with one string parameter.
// {"id": 0, "method": "hello", "params": ["world"], "jsonrpc": "2.0"}

helloRc, err := hello("world")
if err != nil {
	panic(err)
}

// Call the ping method with an empty parameter list. Note that we do not
// wait for the server to complete the hello request above before sending
// the ping request -- this is request pipelining.
// {"id": 1, "method": "ping", "params": [], "jsonrpc": "2.0"}

pingRc, err := ping()
if err != nil {
	panic(err)
}

// Await the response for the ping request and then the hello request. It
// does not matter in which order the server returned the responses.

var resp jsonrpc.Response
resp = <-pingRc
fmt.Printf("%v\n", resp)
resp = <-helloRc
fmt.Printf("%v\n", resp)
```

License
-------

MIT

