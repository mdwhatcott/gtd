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

func (this *Encoder) Encode(v interface{}) error {
	this.encodeValueTypeName(v)
	return this.encodeValue(v)
}

func (this *Encoder) encodeValueTypeName(v interface{}) {
	TYPE := reflect.TypeOf(v)
	NAME := TYPE.String()
	_ = this.inner.Encode(NAME)
}

func (this *Encoder) encodeValue(v interface{}) error {
	return this.inner.Encode(v)
}
