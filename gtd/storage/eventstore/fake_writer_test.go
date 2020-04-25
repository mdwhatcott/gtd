package eventstore

import (
	"bytes"
	"strings"
)

type FakeWriter struct {
	buffer   *bytes.Buffer
	closed   int
	writeErr error
	closeErr error
}

func NewFakeWriter(writeErr, closeErr error) *FakeWriter {
	return &FakeWriter{
		buffer:   new(bytes.Buffer),
		writeErr: writeErr,
		closeErr: closeErr,
	}
}

func (this *FakeWriter) Write(p []byte) (n int, err error) {
	n, err = this.buffer.Write(p)
	if err != nil {
		return n, err
	}
	return n, this.writeErr
}

func (this *FakeWriter) Close() error {
	this.closed++
	return this.closeErr
}

func (this *FakeWriter) Lines() []string {
	return strings.Split(strings.TrimSpace(this.buffer.String()), "\n")
}
