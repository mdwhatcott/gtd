package json

import (
	"encoding/json"
	"io"
	"reflect"
)

type Encoder struct{ inner *json.Encoder }

func NewEncoder(writer io.Writer) *Encoder {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return &Encoder{inner: encoder}
}
func (this *Encoder) Encode(v interface{}) error {
	type_ := reflect.TypeOf(v).String()
	_ = this.inner.Encode(type_)
	return this.inner.Encode(v)
}
