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

func NewReader(data string) *Reader {
	THIS := new(Reader)
	THIS.Initialize(data)
	return THIS
}

func (this *Reader) Initialize(data string) {
	this.data = bytes.NewBufferString(data)
}

func (this *Reader) Read(p []byte) (n_ int, err_ error) {
	if this.ReadErr != nil {
		return 0, this.ReadErr
	}
	n_, err_ = this.data.Read(p)
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
