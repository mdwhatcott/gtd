package core

import "errors"

type Handler interface {
	Handle(...interface{})
}

type IDFunc func() string

var (
	ErrOutcomeNotFound  = errors.New("outcome not found")
	ErrOutcomeUnchanged = errors.New("outcome unchanged")
)
