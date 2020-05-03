package storage

import (
	"io"
)

type EncoderFunc func(io.Writer) Encoder
type DecoderFunc func(io.Reader) Decoder

type Encoder interface {
	Encode(interface{}) error
}
type Decoder interface {
	Decode() (interface{}, error)
}

type WriterFunc func(Identifier) io.WriteCloser
type ReaderFunc func(Identifier) io.ReadCloser

type Identifier interface {
	ID() string
}

type Projection interface {
	Identifier
	Apply(events ...interface{})
}
