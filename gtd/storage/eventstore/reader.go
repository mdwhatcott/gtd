package eventstore

import (
	"io"
	"reflect"

	"github.com/mdwhatcott/gtd/gtd/storage"
	"github.com/mdwhatcott/gtd/gtd/util/errors"
)

type Reader struct {
	reader  storage.ReaderFunc
	decoder storage.DecoderFunc
}

func NewReader(readerFunc storage.ReaderFunc, decoderFunc storage.DecoderFunc) *Reader {
	return &Reader{
		reader:  readerFunc,
		decoder: decoderFunc,
	}
}

func (this *Reader) Read(v ...interface{}) {
	for _, query := range v {
		switch query := query.(type) {
		case *storage.OutcomeEventStream:
			query.Result.Stream = make(chan interface{})
			go this.stream(query, query.Result.Stream)
		default:
			panic(errors.Wrap(ErrUnrecognizedType, reflect.TypeOf(query)))
		}
	}
}

func (this *Reader) stream(id storage.Identifier, stream chan interface{}) {
	defer close(stream)

	reader := this.reader(id)
	defer this.close(reader, stream)

	decoder := this.decoder(reader)
	for {
		value, err := decoder.Decode()
		if err == io.EOF {
			break
		}
		if err != nil {
			stream <- errors.Wrap(ErrUnexpectedReadError, err)
			break
		}
		stream <- value
	}
}

func (this *Reader) close(reader io.ReadCloser, stream chan interface{}) {
	err := reader.Close()
	if err != nil {
		stream <- errors.Wrap(ErrUnexpectedCloseError, err)
	}
}
