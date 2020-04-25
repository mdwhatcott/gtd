package storage

import (
	"io"

	"github.com/smartystreets/joyride/v2"
)

type EncoderFunc func(io.Writer) Encoder

type Encoder interface {
	Encode(interface{}) error
}

type WriterFunc func(AggregateRoot) io.WriteCloser

type AggregateRoot interface {
	ID() string
}

type Projector interface {
	joyride.StorageWriter
}

type EventStore interface {
	joyride.StorageReader
	joyride.StorageWriter
}
