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

func NewReader(_readerFunc storage.ReaderFunc, _decoderFunc storage.DecoderFunc) *Reader {
	return &Reader{
		reader:  _readerFunc,
		decoder: _decoderFunc,
	}
}

func (this *Reader) Read(_v ...interface{}) {
	for _, QUERY := range _v {
		switch QUERY := QUERY.(type) {
		case *storage.OutcomeEventStream:
			QUERY.Result.Stream = make(chan interface{})
			go this.stream(QUERY, QUERY.Result.Stream)
		default:
			panic(errors.Wrap(ErrUnrecognizedType, reflect.TypeOf(QUERY)))
		}
	}
}

func (this *Reader) stream(_id storage.Identifier, _stream chan interface{}) {
	defer close(_stream)

	READER := this.reader(_id)
	defer this.close(READER, _stream)

	DECODER := this.decoder(READER)
	this.decodeStream(DECODER, _stream)
}

func (this *Reader) decodeStream(_decoder storage.Decoder, _stream chan interface{}) {
	for {
		VALUE, ERR := _decoder.Decode()
		if ERR == io.EOF {
			return
		}
		if ERR != nil {
			_stream <- errors.Wrap(ErrUnexpectedReadError, ERR)
			return
		}
		_stream <- VALUE
	}
}

func (this *Reader) close(_reader io.ReadCloser, _stream chan interface{}) {
	ERR := _reader.Close()
	if ERR != nil {
		_stream <- errors.Wrap(ErrUnexpectedCloseError, ERR)
	}
}
