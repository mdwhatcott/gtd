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

func NewFakeEncoder(_writer io.Writer, _err error) *FakeEncoder {
	return &FakeEncoder{
		writer: _writer,
		err:    _err,
	}
}

func (this *FakeEncoder) Encode(_v interface{}) error {
	_, ERR := fmt.Fprintln(this.writer, reflect.TypeOf(_v).Name())
	if ERR != nil {
		return ERR
	}
	return this.err
}
