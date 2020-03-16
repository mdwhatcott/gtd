package gtd

import (
	"fmt"
	"io"
)

type Handler interface {
	Handle(interface{})
}

type Renderer interface {
	io.Writer
	String(fmt.Stringer)
	JSON(interface{})
}

