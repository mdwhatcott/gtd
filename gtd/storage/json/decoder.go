package json

import (
	"encoding/json"
	"io"
	"reflect"

	"github.com/mdwhatcott/gtd/gtd/util/errors"
)

type Decoder struct {
	inner    *json.Decoder
	registry map[string]reflect.Type
}

func NewDecoder(_reader io.Reader, _registry map[string]reflect.Type) *Decoder {
	return &Decoder{
		inner:    json.NewDecoder(_reader),
		registry: _registry,
	}
}

func (this *Decoder) Decode() (interface{}, error) {
	var NAME string
	ERR := this.inner.Decode(&NAME)
	if ERR != nil {
		return nil, errors.Wrap(ErrDecodingInvalidJSONTypeName, ERR)
	}

	TYPE := this.registry[NAME]
	if TYPE == nil {
		return nil, errors.Wrap(ErrDecodingUnregisteredType, NAME)
	}

	VALUE := reflect.Indirect(reflect.New(TYPE)).Interface()
	ERR = this.inner.Decode(&VALUE)
	if ERR != nil {
		return nil, errors.Wrap(ErrDecodingInvalidJSONValue, ERR)
	}

	if reflect.TypeOf(VALUE) != TYPE {
		return nil, errors.Wrap(ErrDecodingTypeMismatch, ERR)
	}
	return VALUE, nil
}
