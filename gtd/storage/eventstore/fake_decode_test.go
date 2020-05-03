package eventstore

import (
	"fmt"
	"io"
)

type FakeDecoder struct {
	reader io.Reader
	err    error
}

func NewFakeDecoder(_reader io.Reader, _err error) *FakeDecoder {
	return &FakeDecoder{
		reader: _reader,
		err:    _err,
	}
}

func (this *FakeDecoder) Decode() (interface{}, error) {
	var N int
	_, ERR := fmt.Fscanln(this.reader, &N)
	return N, ERR
}
