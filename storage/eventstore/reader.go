package eventstore

import (
	"io"
	"reflect"

	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/storage"
	"github.com/mdwhatcott/gtd/util/errors"
)

type Reader struct {
	reader  storage.ReaderFunc
	decoder storage.DecoderFunc
	log     *logging.Logger
}

func NewReader(_readerFunc storage.ReaderFunc, _decoderFunc storage.DecoderFunc) *Reader {
	return &Reader{
		reader:  _readerFunc,
		decoder: _decoderFunc,
	}
}

func (this *Reader) Read(_v ...interface{}) {
	for _, QUERY := range _v {
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

func (this *Reader) stream(_stream chan interface{}, _filter string) {
	defer close(_stream)

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

		if _filter == "" || IDENTIFIABLE.ID() == _filter {
			_stream <- VALUE
		}
	}
	this.log.Printf("Streamed %d events.", COUNT)
}

func (this *Reader) close(_reader io.ReadCloser) {
	ERR := _reader.Close()
	if ERR != nil {
		this.log.Println(errors.Wrap(ErrUnexpectedCloseError, ERR))
	}
}
