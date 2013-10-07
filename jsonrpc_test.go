package jsonrpc_test

import (
	"bytes"
	"github.com/calmh/jsonrpc"
	"io"
	"testing"
)

func TestErrorResponse(t *testing.T) {
	var in = newChanReader(3)
	var out bytes.Buffer

	c := jsonrpc.NewConnection(newReadWriter(in, &out), jsonrpc.StandardDialect)
	login := c.Request("system.login")

	rc, err := login("admin", "test")
	if err != nil {
		t.Error(err)
	}

	in.ch <- []byte(`{"id":1,"error":{"code":-43,"message":"a no you"}}` + "\n")
	in.ch <- []byte(`{"id":0,"error":{"code":-42,"message":"no you"}}` + "\n")
	in.ch <- []byte(`{"id":2,"error":{"code":-44,"message":"no you b"}}` + "\n")

	res := <-rc

	if res.ID != 0 {
		t.Errorf("incorrect response ID %d != %d", res.ID, 0)
	}
	if res.Result != nil {
		t.Errorf("incorrect response Result %v != %v", res.Result, nil)
	}
	if res.Error == nil {
		t.Errorf("unexpected nil response Error")
	} else {
		if c := res.Error.Code; c != -42 {
			t.Errorf("incorrect response Error.Code %d != %d", c, -42)
		}
		if m := res.Error.Message; m != "no you" {
			t.Errorf("incorrect response Error.Message %q != %q", m, "no you")
		}
	}
}

func TestListResponse(t *testing.T) {
	var in = newChanReader(3)
	var out bytes.Buffer

	c := jsonrpc.NewConnection(newReadWriter(in, &out), jsonrpc.StandardDialect)
	login := c.Request("system.login")

	rc, err := login("admin", "test")
	if err != nil {
		t.Error(err)
	}

	in.ch <- []byte(`{"id":1,"error":{"code":-43,"message":"a no you"}}` + "\n")
	in.ch <- []byte(`{"id":0,"result":["a","list","of","strings"]}` + "\n")
	in.ch <- []byte(`{"id":2,"error":{"code":-44,"message":"no you b"}}` + "\n")

	res := <-rc

	if res.ID != 0 {
		t.Errorf("incorrect response ID %d != %d", res.ID, 0)
	}
	if res.Error != nil {
		t.Errorf("unexpected response Error %v", res.Error)
	}
	if l, ok := res.Result.([]interface{}); !ok {
		t.Errorf("incorrect response Result type %T", res.Result)
	} else {
		correct := []string{"a", "list", "of", "strings"}
		for i := range correct {
			if l[i].(string) != correct[i] {
				t.Errorf("incorrect result %d: %q != %q", i, l[i], correct[i])
			}
		}
	}
}

func TestMapResponse(t *testing.T) {
	var in = newChanReader(3)
	var out bytes.Buffer
	c := jsonrpc.NewConnection(newReadWriter(in, &out), jsonrpc.StandardDialect)
	login := c.Request("system.login")

	login("admin", "test")            // id:0
	login("admin", "test")            // id:1
	rc, err := login("admin", "test") // id:2
	if err != nil {
		t.Error(err)
	}

	in.ch <- []byte(`{"id":1,"error":{"code":-43,"message":"a no you"}}` + "\n")
	in.ch <- []byte(`{"id":0,"result":{"foo":"bar","baz":"quux"}}` + "\n")
	in.ch <- []byte(`{"id":2,"result":{"foo":"baz","baz":"quuax"}}` + "\n")

	res := <-rc

	if res.ID != 2 {
		t.Errorf("incorrect response ID %d != %d", res.ID, 2)
	}
	if res.Error != nil {
		t.Errorf("unexpected response Error %v", res.Error)
	}
	if m, ok := res.Result.(map[string]interface{}); !ok {
		t.Errorf("incorrect response Result type %T", res.Result)
	} else {
		correct := map[string]string{"foo": "baz", "baz": "quuax"}
		for k, v := range correct {
			if m[k].(string) != v {
				t.Errorf("incorrect result %q: %q != %q", k, m[k], v)
			}
		}
	}
}

func TestWriteError(t *testing.T) {
	var in = newChanReader(1)
	var out = newExpectWriter([]byte("contents doesn't matter"))

	c := jsonrpc.NewConnection(newReadWriter(in, out), jsonrpc.StandardDialect)
	login := c.Request("system.login")
	ping := c.Notification("ping")

	_, err := login("admin", "test")
	if err != nil {
		t.Error(err)
	}

	// Let that write happen
	<-out.done

	// This write will queue the request and then discover the error in the background
	_, err = login("admin", "test")
	if err != nil {
		t.Error(err)
	}

	// This write will hit the error set previously
	_, err = login("admin", "test")
	if err != io.EOF {
		t.Errorf("unexpected non-EOF error %#v", err)
	}

	// Likewise
	err = ping()
	if err != io.EOF {
		t.Errorf("unexpected non-EOF error %#v", err)
	}
}

func TestReadJSONError(t *testing.T) {
	var in = newChanReader(1)
	var out bytes.Buffer

	c := jsonrpc.NewConnection(newReadWriter(in, &out), jsonrpc.StandardDialect)
	login := c.Request("system.login")
	ping := c.Notification("ping")

	rc, err := login("admin", "test")
	if err != nil {
		t.Error(err)
	}

	in.ch <- []byte("parse this if you can\n")

	// Here we await the read error
	_, ok := <-rc
	if ok {
		t.Errorf("unexpected successful read")
	}

	// This write will hit the error set previously
	_, err = login("admin", "test")
	if err == nil {
		t.Errorf("unexpected nil error")
	}

	// Likewise
	err = ping()
	if err == nil {
		t.Errorf("unexpected nil error")
	}
}
