package eventstore

import (
	"context"
	"io"
	"reflect"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/storage"
	"github.com/mdwhatcott/gtd/v3/util/errors"
)

type Reader struct {
	log     core.Logger
	reader  storage.ReaderFunc
	decoder storage.DecoderFunc
}

func NewReader(logger core.Logger, readerFunc storage.ReaderFunc, decoderFunc storage.DecoderFunc) *Reader {
	return &Reader{
		log:     logger,
		reader:  readerFunc,
		decoder: decoderFunc,
	}
}

func (this *Reader) Read(_ context.Context, v ...interface{}) {
	for _, QUERY := range v {
		switch QUERY := QUERY.(type) {
		case *storage.EventStream:
			QUERY.Result.Events = make(chan interface{})
			go this.stream(QUERY.Result.Events, "")
		case *storage.OutcomeEventStream:
			QUERY.Result.Events = make(chan interface{})
			go this.stream(QUERY.Result.Events, QUERY.OutcomeID)
		default:
			panic(errors.Wrap(ErrUnrecognizedType, reflect.TypeOf(QUERY)))
		}
	}
}

func (this *Reader) stream(stream chan interface{}, filter string) {
	defer close(stream)

	READER := this.reader()
	defer this.close(READER)

	DECODER := this.decoder(READER)
	COUNT := 0
	for ; ; COUNT++ {
		VALUE, ERR := DECODER.Decode()
		if ERR == io.EOF {
			break
		}

		if ERR != nil {
			this.log.Println(errors.Wrap(ErrUnexpectedReadError, ERR))
			break
		}

		IDENTIFIABLE, OK := VALUE.(storage.Identifiable)
		if !OK {
			this.log.Println(errors.Wrap(ErrUnidentifiableType, reflect.TypeOf(VALUE)))
			break
		}

		if filter == "" || IDENTIFIABLE.ID() == filter {
			stream <- VALUE
		}
	}
}

func (this *Reader) close(reader io.ReadCloser) {
	ERR := reader.Close()
	if ERR != nil {
		this.log.Println(errors.Wrap(ErrUnexpectedCloseError, ERR))
	}
}
