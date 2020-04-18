package core

import "errors"

var (
	ErrOutcomeNotFound  = errors.New("outcome not found")
	ErrOutcomeUnchanged = errors.New("outcome unchanged")
)
