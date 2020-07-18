package storage

import (
	"io"
)

const EventsDatabaseFilename = "events.csv"

type EncoderFunc func(io.Writer) Encoder
type DecoderFunc func(io.Reader) Decoder

type Encoder interface{ Encode(interface{}) error }
type Decoder interface{ Decode() (interface{}, error) }

type WriterFunc func() io.WriteCloser
type ReaderFunc func() io.ReadCloser

type Identifiable interface{ ID() string }

type Projection interface {
	Identifiable
	Apply(events ...interface{})
}
