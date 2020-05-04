package fake

import (
	"bytes"
)

type Reader struct {
	data     *bytes.Buffer
	ReadErr  error
	CloseErr error
	Closed   int
}

func NewReader(_data string) *Reader {
	THIS := new(Reader)
	THIS.Initialize(_data)
	return THIS
}

func (this *Reader) Initialize(_data string) {
	this.data = bytes.NewBufferString(_data)
}

func (this *Reader) Read(_p []byte) (n_ int, err_ error) {
	if this.ReadErr != nil {
		return 0, this.ReadErr
	}
	n_, err_ = this.data.Read(_p)
	if err_ != nil {
		return n_, err_
	}
	return n_, nil
}

func (this *Reader) Close() error {
	this.Closed++
	return this.CloseErr
}

func (this *Reader) String() string {
	return this.data.String()
}
