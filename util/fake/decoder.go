package fake

import (
	"fmt"
	"io"
)

type Decoder struct {
	reader io.Reader
	err    error
}

func NewDecoder(reader io.Reader, err error) *Decoder {
	return &Decoder{
		reader: reader,
		err:    err,
	}
}

func (this *Decoder) Decode() (interface{}, error) {
	var N int
	_, ERR := fmt.Fscanln(this.reader, &N)
	if ERR != nil {
		return nil, ERR
	}
	if N < 0 {
		return N, nil
	}
	return NewIdentifiable(N), nil
}
