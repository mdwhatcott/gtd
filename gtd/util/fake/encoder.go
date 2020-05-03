package fake

import (
	"fmt"
	"io"
	"reflect"
)

type Encoder struct {
	writer io.Writer
	err    error
}

func NewEncoder(_writer io.Writer, _err error) *Encoder {
	return &Encoder{
		writer: _writer,
		err:    _err,
	}
}

func (this *Encoder) Encode(_v interface{}) error {
	_, ERR := fmt.Fprintln(this.writer, reflect.TypeOf(_v).Name())
	if ERR != nil {
		return ERR
	}
	return this.err
}
