package eventstore

import (
	"fmt"
	"io"
)

type FakeDecoder struct {
	reader io.Reader
	err    error
}

func NewFakeDecoder(reader io.Reader, err error) *FakeDecoder {
	return &FakeDecoder{
		reader: reader,
		err:    err,
	}
}

func (this *FakeDecoder) Decode() (interface{}, error) {
	var n int
	_, err := fmt.Fscanln(this.reader, &n)
	return n, err
}
