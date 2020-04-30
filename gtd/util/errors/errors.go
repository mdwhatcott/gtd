package errors

import (
	"errors"
	"fmt"
)

var New = errors.New

func Wrap(outer error, inner interface{}) error {
	return fmt.Errorf("%w: [%v]", outer, inner)
}
