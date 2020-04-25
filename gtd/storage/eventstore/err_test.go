package eventstore

import "errors"

var errGophers = errors.New("GOPHERS")

func recovered(action func()) (result interface{}) {
	defer func() { result = recover() }()
	action()
	return result
}
