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

func NewEncoder(writer io.Writer, err error) *Encoder {
	return &Encoder{
		writer: writer,
		err:    err,
	}
}

func (this *Encoder) Encode(v interface{}) error {
	_, ERR := fmt.Fprintln(this.writer, reflect.TypeOf(v).Name())
	if ERR != nil {
		return ERR
	}
	return this.err
}
