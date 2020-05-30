package json

import (
	"encoding/json"
	"io"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/util/errors"
)

type Decoder struct {
	inner    *json.Decoder
	registry map[string]core.Transformer
}

func NewDecoder(_reader io.Reader, _registry map[string]core.Transformer) *Decoder {
	return &Decoder{
		inner:    json.NewDecoder(_reader),
		registry: _registry,
	}
}

func (this *Decoder) Decode() (interface{}, error) {
	var NAME string
	ERR := this.inner.Decode(&NAME)
	if ERR == io.EOF {
		return nil, io.EOF
	}

	if ERR != nil {
		return nil, errors.Wrap(ErrDecodingInvalidJSONTypeName, ERR)
	}

	TRANSFORMER := this.registry[NAME]
	if TRANSFORMER == nil {
		return nil, errors.Wrap(ErrDecodingUnregisteredType, NAME)
	}

	var VALUE map[string]interface{}
	ERR = this.inner.Decode(&VALUE)
	if ERR != nil {
		return nil, errors.Wrap(ErrDecodingInvalidJSONValue, ERR)
	}

	return TRANSFORMER(VALUE), nil
}
