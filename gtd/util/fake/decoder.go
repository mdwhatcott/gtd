package fake

import (
	"fmt"
	"io"
)

type Decoder struct {
	reader io.Reader
	err    error
}

func NewDecoder(_reader io.Reader, _err error) *Decoder {
	return &Decoder{
		reader: _reader,
		err:    _err,
	}
}

func (this *Decoder) Decode() (interface{}, error) {
	var N int
	_, ERR := fmt.Fscanln(this.reader, &N)
	return N, ERR
}
