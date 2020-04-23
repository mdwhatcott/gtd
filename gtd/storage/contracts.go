package storage

import (
	"io"

	"github.com/smartystreets/joyride/v2"
)

type NewEncoder func(io.Writer) Encoder

type Encoder interface {
	Encode(interface{}) error
}

type Writer func(AggregateRoot) io.WriteCloser

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
