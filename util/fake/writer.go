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

func NewWriter(writeErr, closeErr error) *Writer {
	return &Writer{
		buffer:   new(bytes.Buffer),
		writeErr: writeErr,
		closeErr: closeErr,
	}
}

func (this *Writer) SetWriteError(err error) { this.writeErr = err }
func (this *Writer) SetCloseError(err error) { this.closeErr = err }

func (this *Writer) Write(p []byte) (n_ int, err_ error) {
	n_, err_ = this.buffer.Write(p)
	if err_ != nil {
		return n_, err_
	}
	return n_, this.writeErr
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
