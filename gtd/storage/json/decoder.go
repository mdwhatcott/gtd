package json

import (
	"encoding/json"
	"io"
	"reflect"
)

type Decoder struct {
	inner *json.Decoder
}

func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{inner: json.NewDecoder(reader)}
}

func (this *Decoder) DecodeType(registry map[string]reflect.Type) interface{} {
	var type_ string
	_ = this.inner.Decode(&type_)
	return reflect.New(registry[type_]).Elem()
}

func (this *Decoder) Decode(v interface{}) error {
	return this.inner.Decode(&v)
}
