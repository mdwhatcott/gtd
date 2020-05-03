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

func NewFakeWriter(_writeErr, _closeErr error) *FakeWriter {
	return &FakeWriter{
		buffer:   new(bytes.Buffer),
		writeErr: _writeErr,
		closeErr: _closeErr,
	}
}

func (this *FakeWriter) Write(_p []byte) (_n int, _err error) {
	_n, _err = this.buffer.Write(_p)
	if _err != nil {
		return _n, _err
	}
	return _n, this.writeErr
}

func (this *FakeWriter) Close() error {
	this.closed++
	return this.closeErr
}

func (this *FakeWriter) Lines() []string {
	return strings.Split(strings.TrimSpace(this.buffer.String()), "\n")
}
