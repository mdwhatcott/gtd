package fake

import (
	"bytes"
	"strings"
)

type Writer struct {
	buffer   *bytes.Buffer
	closed   int
	writeErr error
	closeErr error
}

func NewWriter(_writeErr, _closeErr error) *Writer {
	return &Writer{
		buffer:   new(bytes.Buffer),
		writeErr: _writeErr,
		closeErr: _closeErr,
	}
}

func (this *Writer) SetWriteError(_err error) { this.writeErr = _err }
func (this *Writer) SetCloseError(_err error) { this.closeErr = _err }

func (this *Writer) Write(_p []byte) (_n int, _err error) {
	_n, _err = this.buffer.Write(_p)
	if _err != nil {
		return _n, _err
	}
	return _n, this.writeErr
}

func (this *Writer) Close() error {
	this.closed++
	return this.closeErr
}

func (this *Writer) Content() string {
	return this.buffer.String()
}

func (this *Writer) Lines() []string {
	return strings.Split(strings.TrimSpace(this.Content()), "\n")
}

func (this *Writer) CloseCount() int {
	return this.closed
}
