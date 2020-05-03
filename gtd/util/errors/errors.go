package errors

import (
	"errors"
	"fmt"
	"runtime"
)

var (
	New    = errors.New
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)

func Wrap(_outer error, _inner interface{}) error {
	_, FILE, LINE, _ := runtime.Caller(1)
	return fmt.Errorf("%w: [%v] (at %s:%d)", _outer, _inner, FILE, LINE)
}
