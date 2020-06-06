package eventstore

import (
	"io"
	"reflect"

	"github.com/mdwhatcott/gtd/storage"
	"github.com/mdwhatcott/gtd/util/errors"
)

type Reader struct {
	reader  storage.ReaderFunc
	decoder storage.DecoderFunc
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
			QUERY.Result.Events = this.stream()
		case *storage.OutcomeEventStream:
			QUERY.Result.Events = filter(this.stream(), QUERY.OutcomeID)
		default:
			panic(errors.Wrap(ErrUnrecognizedType, reflect.TypeOf(QUERY)))
		}
	}
}

func filter(_stream []interface{}, _id string) (filtered []interface{}) {
	for _, ITEM := range _stream {
		IDENTIFIABLE, OK := ITEM.(storage.Identifiable)
		if !OK {
			panic(errors.Wrap(ErrUnidentifiableType, reflect.TypeOf(ITEM)))
		}
		if IDENTIFIABLE.ID() == _id {
			filtered = append(filtered, ITEM)
		}
	}
	return filtered
}

func (this *Reader) stream() (events_ []interface{}) {
	READER := this.reader()
	defer this.close(READER)

	DECODER := this.decoder(READER)
	for {
		VALUE, ERR := DECODER.Decode()
		if ERR == io.EOF {
			return
		}
		if ERR != nil {
			panic(errors.Wrap(ErrUnexpectedReadError, ERR))
		}
		events_ = append(events_, VALUE)
	}
}

func (this *Reader) close(_reader io.ReadCloser) {
	ERR := _reader.Close()
	if ERR != nil {
		panic(errors.Wrap(ErrUnexpectedCloseError, ERR))
	}
}
