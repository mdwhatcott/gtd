package json

import "errors"

var (
	ErrDecodingInvalidJSONTypeName = errors.New("invalid JSON for type name")
	ErrDecodingUnregisteredType    = errors.New("unregistered type name")
	ErrDecodingInvalidJSONValue    = errors.New("invalid JSON for value")
	ErrDecodingTypeMismatch        = errors.New("type/value mismatch")
)
