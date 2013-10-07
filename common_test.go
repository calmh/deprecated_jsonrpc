package jsonrpc_test

import (
	"bytes"
	"errors"
	"io"
)

// The expectWriter expects exactly one write with the specified contents
type expectWriter struct {
	expected []byte
	done     chan struct{}
	err      error
}

func (e *expectWriter) Write(bs []byte) (int, error) {
	select {
	case <-e.done:
		return 0, io.EOF
	default:
		if bytes.Compare(bs, e.expected) != 0 {
			e.err = errors.New("incorrect write")
		}
		close(e.done)
		return len(bs), nil
	}
}

func newExpectWriter(expected []byte) *expectWriter {
	var done = make(chan struct{})
	return &expectWriter{expected, done, nil}
}

// The readWriter bundles an io.Reader and an io.Writer into an io.ReadWriter
type readWriter struct {
	io.Reader
	io.Writer
}

func newReadWriter(r io.Reader, w io.Writer) io.ReadWriter {
	return &readWriter{r, w}
}

// The chanReader answers n Read() calls with data from a channel, then returns EOF
type chanReader struct {
	ch chan []byte
	n  int
}

func (c *chanReader) Read(bs []byte) (n int, err error) {
	if c.n == 0 {
		return 0, io.EOF
	}

	tbs := <-c.ch
	n = copy(bs, tbs)
	c.n--
	return
}

func newChanReader(n int) *chanReader {
	ch := make(chan []byte, n)
	return &chanReader{ch, n}
}
