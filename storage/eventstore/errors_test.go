package eventstore

import "errors"

var errGophers = errors.New("GOPHERS")

func recovered(_action func()) (result_ interface{}) {
	defer func() { result_ = recover() }()
	_action()
	return result_
}
