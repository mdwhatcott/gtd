package eventstore

import (
	"io"
	"reflect"

	"github.com/mdwhatcott/gtd/v3/storage"
	"github.com/mdwhatcott/gtd/v3/util/errors"
)

type Writer struct {
	encoder storage.EncoderFunc
	writer  storage.WriterFunc
}

func NewWriter(_encoder storage.EncoderFunc, _writer storage.WriterFunc) *Writer {
	return &Writer{
		encoder: _encoder,
		writer:  _writer,
	}
}

func (this *Writer) Write(_events ...interface{}) {
	for _, EVENT := range _events {
		ROOT, OK := EVENT.(storage.Identifiable)
		if !OK {
			panic(errors.Wrap(ErrUnrecognizedType, reflect.TypeOf(EVENT)))
		}
		this.persist(ROOT)
	}
}

func (this *Writer) persist(_root storage.Identifiable) {
	WRITER := this.writer()
	defer this.close(WRITER)

	ENCODER := this.encoder(WRITER)
	ERR := ENCODER.Encode(_root)
	if ERR != nil {
		panic(errors.Wrap(ErrUnexpectedWriteError, ERR))
	}
}

func (this *Writer) close(_writer io.WriteCloser) {
	ERR := _writer.Close()
	if ERR != nil {
		panic(errors.Wrap(ErrUnexpectedCloseError, ERR))
	}
}
