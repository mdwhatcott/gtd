package storage

import (
	"io"

	"github.com/smartystreets/joyride/v2"
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

type Projector interface {
	joyride.StorageWriter
}

type EventStore interface {
	joyride.StorageReader
	joyride.StorageWriter
}
