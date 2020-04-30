package eventstore

import (
	"fmt"
	"io"
	"reflect"

	"github.com/mdwhatcott/gtd/gtd/storage"
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
			panic(fmt.Errorf("unrecognized query type: <%v>", reflect.TypeOf(query))) // TODO: wrapped defined error type
		}
	}
}

func (this *Reader) stream(id storage.Identifier, stream chan interface{}) {
	defer close(stream)

	reader := this.reader(id)
	defer closing(reader)

	decoder := this.decoder(reader)
	for {
		value, err := decoder.Decode()
		if err == io.EOF {
			break
		}
		// TODO: panic if err != nil, wrapped in defined error
		stream <- value
	}
}

func closing(closer io.Closer) {
	_ = closer.Close() // TODO: panic, wrapped in defined error
}
