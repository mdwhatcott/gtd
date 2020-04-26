package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
)

type Decoder struct {
	inner *json.Decoder
	err   error
}

func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{inner: json.NewDecoder(reader)}
}

func (this *Decoder) Error() error {
	return this.err
}

func (this *Decoder) DecodeAll(registry map[string]reflect.Type) chan interface{} {
	all := make(chan interface{})
	go this.decodeAll(registry, all)
	return all
}

func (this *Decoder) decodeAll(registry map[string]reflect.Type, all chan interface{}) {
	defer close(all)
	for {
		var value interface{}
		value, this.err = this.decodeType(registry)
		if this.err == io.EOF {
			this.err = nil
			break
		}
		if this.err != nil {
			break
		}

		this.err = this.inner.Decode(&value)
		if this.err != nil {
			break
		}
		all <- value
	}
}

func (this *Decoder) decodeType(registry map[string]reflect.Type) (interface{}, error) {
	var name string
	err := this.inner.Decode(&name)
	if err == io.EOF {
		return nil, err
	}
	type_ := registry[name]
	if type_ == nil {
		return nil, fmt.Errorf("%w: [%s]", errUnregisteredType, name)
	}
	return reflect.New(type_).Elem(), nil
}

var errUnregisteredType = errors.New("unregistered type name")
