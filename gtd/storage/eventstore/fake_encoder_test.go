package eventstore

import (
	"fmt"
	"io"
	"reflect"
)

type FakeEncoder struct {
	writer io.Writer
	err    error
}

func NewFakeEncoder(writer io.Writer, err error) *FakeEncoder {
	return &FakeEncoder{
		writer: writer,
		err:    err,
	}
}

func (this *FakeEncoder) Encode(v interface{}) error {
	_, err := fmt.Fprintln(this.writer, reflect.TypeOf(v).Name())
	if err != nil {
		return err
	}
	return this.err
}
