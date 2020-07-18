package eventstore

import "errors"

var errGophers = errors.New("GOPHERS")

func recovered(action func()) (result_ interface{}) {
	defer func() { result_ = recover() }()
	action()
	return result_
}
