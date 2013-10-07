package jsonrpc

import (
	"encoding/json"
	"io"
	"sync"
)

// A Connection is a JSON-RPC channel over an existing ReadWriter, such as a socket.
type Connection struct {
	sync.Mutex
	dialect Dialect
	err     error
	enc     *json.Encoder
	dec     *json.Decoder
	outbox  chan map[string]interface{}
	results map[int]chan Response
	nextId  int
}

// A Dialect is a specific wire encoding of requests and notifications.
type Dialect struct {
	Request      func(int, string, []interface{}) map[string]interface{}
	Notification func(string, []interface{}) map[string]interface{}
}

// A Response is a JSON-RPC 2.0 response as received by the server in response
// to a request.
type Response struct {
	Id     int
	Result interface{}
	Error  *struct {
		Code    int
		Message string
		Data    interface{}
	}
	JSONRPC string
}

// NewConnection creates a new JSON-RPC channel over the specified ReadWriter,
// using a specific Dialect to encode the transmitted data.
func NewConnection(rw io.ReadWriter, dialect Dialect) *Connection {
	enc := json.NewEncoder(rw)
	dec := json.NewDecoder(rw)
	conn := &Connection{
		dialect: dialect,
		enc:     enc,
		dec:     dec,
		outbox:  make(chan map[string]interface{}),
		results: make(map[int]chan Response),
	}
	go conn.writer()
	go conn.reader()
	return conn
}

// A Request function performs a request on the JSON-RPC server when called and
// returns a channel where the response may be read. The response object that
// is read from the channel may contain an error from the server. A nil channel
// and an error is returned if the request could not be sent due to
// communication problems. The meaning of the arguments to the request function
// is defined by the Dialect, but in the standard dialect these are simply
// sent verbatim as the method arguments.
type Request func(...interface{}) (<-chan Response, error)

// Request returns a new request function for use over the channel.
func (c *Connection) Request(method string) Request {
	return func(params ...interface{}) (<-chan Response, error) {
		if c.err != nil {
			return nil, c.err
		}

		res := make(chan Response, 1)

		c.Lock()
		id := c.nextId
		c.results[id] = res
		c.nextId++
		c.Unlock()

		c.outbox <- c.dialect.Request(id, method, params)
		return res, nil
	}
}

// A Notification function performs a notification on the JSON-RPC server when
// called.  An error is returned if the request could not be sent due to
// communication problems. The meaning of the arguments to the request function
// is defined by the Dialect, but in the standard dialect these are simply sent
// verbatim as the method arguments.
type Notification func(...interface{}) error

// Notification returns a new notification function for use over the channel.
func (c *Connection) Notification(method string) Notification {
	return func(params ...interface{}) error {
		if c.err != nil {
			return c.err
		}

		c.outbox <- c.dialect.Notification(method, params)
		return nil
	}
}

func (c *Connection) writer() {
	for wc := range c.outbox {
		c.err = c.enc.Encode(wc)
		if c.err != nil {
			break
		}
	}
}

func (c *Connection) reader() {
	for {
		var res Response
		c.err = c.dec.Decode(&res)
		if c.err != nil {
			break
		}

		c.Lock()
		ch, ok := c.results[res.Id]
		delete(c.results, res.Id)
		c.Unlock()

		if !ok {
			println("bug: result for unknown id", res.Id)
			continue
		}

		ch <- res
	}
}
