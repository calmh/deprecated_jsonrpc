package jsonrpc_test

import (
	"github.com/calmh/jsonrpc"
	"testing"
)

func TestProceraRequest(t *testing.T) {
	correct := []byte(`{"id":0,"method":"system.login","params":["admin","test"],"tags":["foo","bar"]}` + "\n")
	var out = newExpectWriter(correct)
	var in = newChanReader(1)

	c := jsonrpc.NewConnection(newReadWriter(in, out), jsonrpc.ProceraDialect)
	login := c.Request("system.login")
	_, err := login([]string{"foo", "bar"}, "admin", "test")
	if err != nil {
		t.Error(err)
	}

	<-out.done
	if out.err != nil {
		t.Error(out.err)
	}
}

func TestProceraNotification(t *testing.T) {
	correct := []byte(`{"method":"system.login","params":["admin","test"],"tags":null}` + "\n")
	var out = newExpectWriter(correct)
	var in = newChanReader(1)

	c := jsonrpc.NewConnection(newReadWriter(in, out), jsonrpc.ProceraDialect)
	login := c.Notification("system.login")
	err := login(nil, "admin", "test")
	if err != nil {
		t.Error(err)
	}

	<-out.done
	if out.err != nil {
		t.Error(out.err)
	}
}
