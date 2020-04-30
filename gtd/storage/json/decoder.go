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

func NewDecoder(reader io.Reader, registry map[string]reflect.Type) *Decoder {
	return &Decoder{
		inner:    json.NewDecoder(reader),
		registry: registry,
	}
}

func (this *Decoder) Decode() (interface{}, error) {
	var NAME string
	ERR := this.inner.Decode(&NAME)
	if ERR != nil {
		return nil, errors.Wrap(errDecodingInvalidJSONTypeName, ERR)
	}

	TYPE := this.registry[NAME]
	if TYPE == nil {
		return nil, errors.Wrap(errDecodingUnregisteredType, NAME)
	}

	VALUE := reflect.Indirect(reflect.New(TYPE)).Interface()
	ERR = this.inner.Decode(&VALUE)
	if ERR != nil {
		return nil, errors.Wrap(errDecodingInvalidJSONValue, ERR)
	}

	if reflect.TypeOf(VALUE) != TYPE {
		return nil, errors.Wrap(errDecodingTypeMismatch, ERR)
	}
	return VALUE, nil
}

var (
	errDecodingInvalidJSONTypeName = errors.New("invalid JSON for type name")
	errDecodingUnregisteredType    = errors.New("unregistered type name")
	errDecodingInvalidJSONValue    = errors.New("invalid JSON for value")
	errDecodingTypeMismatch        = errors.New("type/value mismatch")
)
