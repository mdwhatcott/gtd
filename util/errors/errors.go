package errors

import (
	"errors"
	"fmt"
	"runtime"
)

var New = errors.New

func Wrap(outer error, inner interface{}) error {
	_, FILE, LINE, _ := runtime.Caller(1)
	return fmt.Errorf("%w: [%v] (at %s:%d)", outer, inner, FILE, LINE)
}
