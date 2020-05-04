package eventstore

import "errors"

var (
	ErrUnrecognizedType     = errors.New("unrecognized type")
	ErrUnidentifiableType   = errors.New("unidentifiable type")
	ErrUnexpectedReadError  = errors.New("unexpected read err")
	ErrUnexpectedWriteError = errors.New("unexpected write err")
	ErrUnexpectedCloseError = errors.New("unexpected close err")
)
