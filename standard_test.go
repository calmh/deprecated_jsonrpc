package jsonrpc_test

import (
	"github.com/calmh/jsonrpc"
	"testing"
)

func TestStandardRequest(t *testing.T) {
	correct := []byte(`{"id":0,"jsonrpc":"2.0","method":"system.login","params":["admin","test"]}` + "\n")
	var out = newExpectWriter(correct)
	var in = newChanReader(1)

	c := jsonrpc.NewConnection(newReadWriter(in, out), jsonrpc.StandardDialect)
	login := c.Request("system.login")
	_, err := login("admin", "test")
	if err != nil {
		t.Error(err)
	}

	<-out.done
	if out.err != nil {
		t.Error(out.err)
	}
}

func TestStandardNotification(t *testing.T) {
	correct := []byte(`{"jsonrpc":"2.0","method":"system.login","params":["admin","test"]}` + "\n")
	var out = newExpectWriter(correct)
	var in = newChanReader(1)

	c := jsonrpc.NewConnection(newReadWriter(in, out), jsonrpc.StandardDialect)
	login := c.Notification("system.login")
	err := login("admin", "test")
	if err != nil {
		t.Error(err)
	}

	<-out.done
	if out.err != nil {
		t.Error(out.err)
	}
}
