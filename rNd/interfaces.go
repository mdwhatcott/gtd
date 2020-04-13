package rNd

import (
	"fmt"
	"io"
)

type Handler interface {
	Handle(interface{})
}

type Renderer interface {
	io.Writer
	RenderStringer(fmt.Stringer)
	RenderJSON(interface{})
}

