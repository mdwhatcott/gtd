package json

import (
	"encoding/json"
	"io"
	"reflect"
)

type Encoder struct{ inner *json.Encoder }

func NewEncoder(writer io.Writer) *Encoder {
	ENCODER := json.NewEncoder(writer)
	ENCODER.SetIndent("", "  ")
	return &Encoder{inner: ENCODER}
}

func (this *Encoder) Encode(_v interface{}) error {
	this.encodeValueTypeName(_v)
	return this.encodeValue(_v)
}

func (this *Encoder) encodeValueTypeName(_v interface{}) {
	TYPE := reflect.TypeOf(_v)
	NAME := TYPE.String()
	_ = this.inner.Encode(NAME)
}

func (this *Encoder) encodeValue(_v interface{}) error {
	return this.inner.Encode(_v)
}
