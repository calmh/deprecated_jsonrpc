package jsonrpc_test

import (
	"bytes"
	"github.com/calmh/jsonrpc"
	"io"
	"testing"
)

type readWriter struct {
	io.Reader
	io.Writer
}

func NewReadWriter(r io.Reader, w io.Writer) io.ReadWriter {
	return &readWriter{r, w}
}

type chanReader chan []byte

func (ch chanReader) Read(bs []byte) (n int, err error) {
	tbs := <-ch
	n = copy(bs, tbs)
	return
}

func TestSend(t *testing.T) {
	var in, out bytes.Buffer
	c := jsonrpc.NewConnection(NewReadWriter(&in, &out), jsonrpc.StandardDialect)
	login := c.Request("system.login")
	login("admin", "test")

	correct := `{"id":0,"jsonrpc":"2.0","method":"system.login","params":["admin","test"]}` + "\n"
	if out.String() != correct {
		t.Errorf("incorrect command sent %q != %q", out.String(), correct)
	}
}

func TestErrorResponse(t *testing.T) {
	var in = make(chanReader, 3)
	var out bytes.Buffer
	c := jsonrpc.NewConnection(NewReadWriter(&in, &out), jsonrpc.StandardDialect)
	login := c.Request("system.login")

	rc, err := login("admin", "test")
	if err != nil {
		t.Error(err)
	}

	in <- []byte(`{"id":1,"error":{"code":-43,"message":"a no you"}}` + "\n")
	in <- []byte(`{"id":0,"error":{"code":-42,"message":"no you"}}` + "\n")
	in <- []byte(`{"id":2,"error":{"code":-44,"message":"no you b"}}` + "\n")

	res := <-rc

	if res.Id != 0 {
		t.Errorf("incorrect response Id %d != %d", res.Id, 0)
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
	var in = make(chanReader, 3)
	var out bytes.Buffer
	c := jsonrpc.NewConnection(NewReadWriter(&in, &out), jsonrpc.StandardDialect)
	login := c.Request("system.login")

	rc, err := login("admin", "test")
	if err != nil {
		t.Error(err)
	}

	in <- []byte(`{"id":1,"error":{"code":-43,"message":"a no you"}}` + "\n")
	in <- []byte(`{"id":0,"result":["a","list","of","strings"]}` + "\n")
	in <- []byte(`{"id":2,"error":{"code":-44,"message":"no you b"}}` + "\n")

	res := <-rc

	if res.Id != 0 {
		t.Errorf("incorrect response Id %d != %d", res.Id, 0)
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
	var in = make(chanReader, 3)
	var out bytes.Buffer
	c := jsonrpc.NewConnection(NewReadWriter(&in, &out), jsonrpc.StandardDialect)
	login := c.Request("system.login")

	login("admin", "test")            // id:0
	login("admin", "test")            // id:1
	rc, err := login("admin", "test") // id:2
	if err != nil {
		t.Error(err)
	}

	in <- []byte(`{"id":1,"error":{"code":-43,"message":"a no you"}}` + "\n")
	in <- []byte(`{"id":0,"result":{"foo":"bar","baz":"quux"}}` + "\n")
	in <- []byte(`{"id":2,"result":{"foo":"baz","baz":"quuax"}}` + "\n")

	res := <-rc

	if res.Id != 2 {
		t.Errorf("incorrect response Id %d != %d", res.Id, 2)
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
